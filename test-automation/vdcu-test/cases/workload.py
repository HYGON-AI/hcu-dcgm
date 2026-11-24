# -*- coding: utf-8 -*-
"""第 6/7/9 章：训练、推理、集合通信与隔离性用例。"""
import json
import re
import shlex
import time

from vdcu.core import case, parse_device_info

AR_SCRIPT = r"""
python - <<'PY'
import os, torch, torch.distributed as dist
rank = int(os.environ["RANK"]); world = int(os.environ["WORLD_SIZE"])
dist.init_process_group("nccl", rank=rank, world_size=world)
t = torch.ones(4, device="cuda")
dist.all_reduce(t)
if rank == 0:
    print("ALLREDUCE_OK", t.tolist(), flush=True)
dist.destroy_process_group()
PY
"""

SMOKE_ARGS = "--samples 1200 --epochs 1 --batch 8 --max-len 64"


def _make(ctx, card, count, cu=None, mem=None, tries=4):
    ctx.expect(ctx.destroy_card(card), "前置清理失败：卡 %d 上仍有 vDCU 未被释放" % card)
    args = ["virtual", "--create-vdevices", str(count)]
    if cu is not None:
        args += ["--vdevice-compute-units", cu]
    if mem is not None:
        args += ["--vdevice-memory-size", mem]
    args += ["-d", str(card)]
    rc, out = ctx.hy(*args)
    for _ in range(tries - 1):
        if rc == 0:
            break
        # 运行时段性 generic error（ISSUE-K8S-002 宿主形态）：销毁后立即创建偶发失败，稍候重试
        time.sleep(3)
        rc, out = ctx.hy(*args)
    ctx.expect(rc == 0, "创建 vDCU 失败: %s" % out[-300:])
    return sorted(ctx.vdevs_on(card), key=lambda x: x["id"])


def _smoke_result(out):
    m = re.search(r"SMOKE_RESULT (\{.*\})", out)
    ctx_msg = "未找到 SMOKE_RESULT，训练未完成"
    if not m:
        raise AssertionError(ctx_msg + "\n" + out[-600:])
    return json.loads(m.group(1))


def _run_bert(ctx, vdev, expected_mem=None, label="vDCU", timeout=1200):
    rc, out = ctx.docker(
        [vdev], ctx.assets_cmd("bert_smoke.py " + SMOKE_ARGS),
        env={"HIP_VISIBLE_DEVICES": "0"}, mount_assets=True, timeout=timeout,
    )
    ctx.expect(rc == 0, "%s 训练失败: %s" % (label, out[-500:]))
    res = _smoke_result(out)
    ctx.expect(res["acc"] >= 0.5, "%s 训练精度异常: %s" % (label, res["acc"]))
    ctx.expect(res["throughput"] > 0, "%s 吞吐异常" % label)
    if expected_mem is not None:
        ctx.expect(abs(res["mem_gb"] - expected_mem) < 0.5,
                   "%s 可见显存 %.2f GB，应为配额 %.2f GB" %
                   (label, res["mem_gb"], expected_mem))
    print("        %s: %s" % (label, json.dumps(res, ensure_ascii=False)))
    return res


@case("VDCU-TRN-001", "Bert 整卡单实例训练基线（冒烟）", "TRN")
def trn_001(ctx):
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    rc, out = ctx.docker([], ctx.assets_cmd("bert_smoke.py " + SMOKE_ARGS),
                         env={"HIP_VISIBLE_DEVICES": "0"},
                         mount_assets=True, timeout=1200)
    ctx.expect(rc == 0, "整卡训练失败: %s" % out[-500:])
    res = _smoke_result(out)
    ctx.expect(res["acc"] >= 0.5, "整卡训练精度异常: %s" % res["acc"])
    ctx.expect(res["throughput"] > 0, "整卡吞吐异常: %s" % res["throughput"])
    print("        baseline: %s" % json.dumps(res, ensure_ascii=False))


@case("VDCU-TRN-206", "Bert 短跑冒烟（vDCU 32GB）", "TRN")
def trn_206(ctx):
    vs = _make(ctx, ctx.card, 2)
    expect_mem = vs[0]["mem"] / 1024 ** 3
    _run_bert(ctx, vs[0]["id"], expect_mem, "二等分 vDCU")


@case("VDCU-TRN-003", "Bert 四等分 vDCU 单实例训练", "TRN")
def trn_003(ctx):
    vs = _make(ctx, ctx.card, 4)
    _run_bert(ctx, vs[0]["id"], vs[0]["mem"] / 1024 ** 3, "四等分 vDCU")


@case("VDCU-TRN-004", "Bert CU 共享 vDCU 单实例训练", "TRN")
def trn_004(ctx):
    vs = _make(ctx, ctx.card, 2, cu="0,0")
    _run_bert(ctx, vs[0]["id"], vs[0]["mem"] / 1024 ** 3, "CU 共享 vDCU")


@case("VDCU-TRN-005", "同卡二等分 vDCU 双容器并发训练", "TRN", priority="P1")
def trn_005(ctx):
    vs = _make(ctx, ctx.card, 2)
    names = ["tz_trn005_a", "tz_trn005_b"]
    for name in names:
        ctx.sh(["docker", "rm", "-f", name], timeout=60)
    try:
        for name, vdev in zip(names, vs):
            rc, out = ctx.docker(
                [vdev["id"]], ctx.assets_cmd("bert_smoke.py " + SMOKE_ARGS),
                name=name, detach=True, env={"HIP_VISIBLE_DEVICES": "0"},
                mount_assets=True, timeout=60,
            )
            ctx.expect(rc == 0, "并发容器 %s 启动失败: %s" % (name, out[-300:]))
        for name in names:
            rc, out = ctx.sh(["docker", "wait", name], timeout=1500)
            ctx.expect(rc == 0, "并发容器 %s 等待失败" % name)
            rc, logs = ctx.sh(["docker", "logs", name], timeout=60)
            ctx.expect(rc == 0, "并发容器 %s 日志读取失败" % name)
            res = _smoke_result(logs)
            ctx.expect(res["acc"] >= 0.5 and res["throughput"] > 0,
                       "并发容器 %s 训练结果异常" % name)
    finally:
        for name in names:
            ctx.sh(["docker", "rm", "-f", name], timeout=60)


@case("VDCU-INF-001", "单 vDCU 模型加载与生成（冒烟）", "INF")
def inf_001(ctx):
    vs = _make(ctx, ctx.card, 2)   # 32 GB
    rc, out = ctx.docker([vs[0]["id"]],
                         ctx.assets_cmd("gen_smoke.py --max-new-tokens 32"),
                         env={"HIP_VISIBLE_DEVICES": "0"}, mount_assets=True,
                         timeout=1200)
    ctx.expect(rc == 0, "推理执行失败: %s" % out[-600:])
    m = re.search(r"GEN_RESULT (\{.*\})", out)
    ctx.expect(m is not None, "未找到 GEN_RESULT，推理未完成\n" + out[-500:])
    res = json.loads(m.group(1))

    ctx.expect(res["new_tokens"] >= 16, "生成 token 数不足: %s" % res["new_tokens"])
    ctx.expect(res["tok_per_s"] > 0, "生成吞吐异常")
    expect_mem = vs[0]["mem"] / 1024 ** 3
    ctx.expect(abs(res["mem_gb"] - expect_mem) < 0.5,
               "推理进程可见显存 %.2f GB，应为配额 %.2f GB" % (res["mem_gb"], expect_mem))
    print("        gen: %s" % json.dumps(res, ensure_ascii=False))


@case("VDCU-INF-002", "vLLM 单实例服务与 OpenAI 接口冒烟", "INF", priority="P1")
def inf_002(ctx):
    if not ctx.vllm_model:
        ctx.skip("未配置 VDCU_VLLM_MODEL，按方案第 7 章配置模型路径后执行")
    vs = _make(ctx, ctx.card, 2)
    script = ctx.assets_cmd("vllm_smoke.py --model %s --max-model-len 2048" % shlex.quote(ctx.vllm_model))
    rc, out = ctx.docker(
        [vs[0]["id"]], script, env={"HIP_VISIBLE_DEVICES": "0"},
        mount_assets=True, image=ctx.inference_image, shm="16G", timeout=1500,
    )
    ctx.expect(rc == 0, "vLLM 服务冒烟失败: %s" % out[-800:])
    m = re.search(r"VLLM_RESULT (\{.*\})", out)
    ctx.expect(m is not None, "未找到 VLLM_RESULT，服务未完成启动或接口失败")
    result = json.loads(m.group(1))
    ctx.expect(result.get("model"), "vLLM 未返回模型列表")
    ctx.expect(result.get("text"), "vLLM completions 接口未返回文本")
    print("        vllm: %s" % json.dumps(result, ensure_ascii=False))


def _comm_reason(out):
    """从 torchrun/NCCL 输出中提取可读的失败原因。"""
    for pat in (r"Duplicate GPU detected[^\n]*",
                r"ncclInvalidUsage[^\n]*",
                r"HIP error[^\n]*",
                r"RuntimeError[^\n]*"):
        m = re.search(pat, out)
        if m:
            return m.group(0).strip()[:160]
    return "未出现 ALLREDUCE_OK（train 未完成）"


@case("VDCU-COM-002", "单容器同卡双 vDCU 集合通信", "COM",
      expected_fail=True, issue="ISSUE-013")
def com_002(ctx):
    vs = _make(ctx, ctx.card, 2)
    script = "cd /tmp && torchrun --standalone --nproc_per_node=2 " + ctx.assets_path("allreduce.py")
    rc, out = ctx.docker([vs[0]["id"], vs[1]["id"]], script,
                         env={"HIP_VISIBLE_DEVICES": "0,1"}, mount_assets=True,
                         timeout=600)
    ctx.expect(rc == 0 and "ALLREDUCE_OK" in out,
               "同卡双 vDCU 集合通信失败：%s" % _comm_reason(out))


@case("VDCU-COM-003", "单容器跨卡双 vDCU 集合通信", "COM",
      expected_fail=True, issue="ISSUE-013")
def com_003(ctx):
    a = _make(ctx, ctx.card, 2)
    b = _make(ctx, ctx.card2, 2)
    script = "cd /tmp && torchrun --standalone --nproc_per_node=2 " + ctx.assets_path("allreduce.py")
    rc, out = ctx.docker([a[0]["id"], b[0]["id"]], script,
                         env={"HIP_VISIBLE_DEVICES": "0,1"}, mount_assets=True,
                         timeout=600)
    ctx.expect(rc == 0 and "ALLREDUCE_OK" in out,
               "跨卡双 vDCU 集合通信失败：%s" % _comm_reason(out))


@case("VDCU-ISO-001", "显存配额硬隔离（越界分配）", "ISO")
def iso_001(ctx):
    vs = _make(ctx, ctx.card, 4)   # 16 GB each
    quota_gb = vs[0]["mem"] / 1024 ** 3
    script = r"""
python - <<'PY'
import torch
try:
    x = torch.empty(int(24 * 1024 ** 3), dtype=torch.uint8, device="cuda")
    print("ALLOC_OK", x.numel())
except Exception as e:
    print("ALLOC_FAIL", type(e).__name__, str(e)[:200])
PY
"""
    rc, out = ctx.docker([vs[0]["id"]], script,
                         env={"HIP_VISIBLE_DEVICES": "0"}, timeout=300)
    ctx.expect("ALLOC_FAIL" in out, "越界分配未被拒绝（配额 %.1f GB）" % quota_gb)
    ctx.expect("out of memory" in out.lower() or "OutOfMemory" in out,
               "缺少 OOM 错误信息: %s" % out[-300:])
    ctx.expect("ALLOC_OK" not in out, "越界分配意外成功")


@case("VDCU-ISO-002", "CU 配额硬隔离（算力不越界）", "ISO", priority="P1")
def iso_002(ctx):
    vs = _make(ctx, ctx.card, 4)   # 20 CU each
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡性能基线前清理失败")
    timer = r"""
python - <<'PY'
import torch, time
torch.cuda.set_device(0)
a = torch.randn(2048, 2048, device="cuda")
b = torch.randn(2048, 2048, device="cuda")
for _ in range(5):
    a @ b
torch.cuda.synchronize()
t = time.time()
for _ in range(40):
    a @ b
torch.cuda.synchronize()
print("ELAPSED %.3f" % (time.time() - t))
PY
"""
    rc, out_v = ctx.docker([vs[0]["id"]], timer, env={"HIP_VISIBLE_DEVICES": "0"},
                           timeout=600)
    ctx.expect(rc == 0, "vDCU 计时失败: %s" % out_v[-300:])

    rc, out_f = ctx.docker([], timer,
                           env={"HIP_VISIBLE_DEVICES": "0"},
                           timeout=600)
    ctx.expect(rc == 0, "整卡计时失败: %s" % out_f[-300:])
    match_v = re.search(r"ELAPSED ([\d.]+)", out_v)
    match_f = re.search(r"ELAPSED ([\d.]+)", out_f)
    ctx.expect(match_v is not None and match_f is not None,
               "计时输出缺少 ELAPSED 字段")
    t_vdcu = float(match_v.group(1))
    t_full = float(match_f.group(1))

    ratio = t_vdcu / t_full if t_full > 0 else 999
    print("        4等分耗时 %.2fs，整卡 %.2fs，比值 %.2f" % (t_vdcu, t_full, ratio))
    ctx.expect(ratio >= 2.0,
               "4 等分 vDCU 耗时仅为整卡的 %.2f 倍，CU 配额疑似未生效" % ratio)

# ---------------------------------------------------------------------------
# TRN-006U~010U + TRN-CMP 多实例并发性能对比（组长要求 2026-11 会议纪要）
# 测试目的：回答"整卡多 stream 共享"与"vDCU 隔离（等分/CU 共享等分）"
# 两种并发机制下单个实例的性能差异。
# 测量纪律：N 个相同实例全部拉起、全部运行后再比较，不允许串行测量。
# 预期结论（建强文档）：stream 与 CU 共享四等分无明显差异（≤10%）。
# 模型：ModernBERT-base（main_modernbert.py，weibo_sentiment 4800 样本）。
# ---------------------------------------------------------------------------

MODERNBERT_CMD = "cd %s && EPOCHS=1 python main_modernbert.py"


def _concurrent_train(ctx, mounts, label, timeout=1800):
    """同步拉起 N 个相同训练实例，全部运行完成后汇总各实例吞吐。

    mounts: [(容器名, vdev_id 或 None, env), ...]
      - vdev_id=None 表示整卡实例（多 stream 形态，最小挂载整卡设备符）
    返回 {容器名: throughput}。
    """
    names = [m[0] for m in mounts]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    started = []
    try:
        # 阶段一：全部 detach 同步拉起（不允许先跑完一个再起下一个）
        for name, vdev, env in mounts:
            vdevs = [vdev] if vdev is not None else []
            rc, out = ctx.docker(vdevs, MODERNBERT_CMD % ctx.workspace_container_dir, name=name,
                                 detach=True, env=env, timeout=120)
            ctx.expect(rc == 0, "并发容器 %s 启动失败: %s" % (name, out[-300:]))
            started.append(name)
        # 阶段二：全部运行完成后统一采集
        results = {}
        for name in names:
            rc, out = ctx.sh(["docker", "wait", name], timeout=timeout)
            ctx.expect(rc == 0, "容器 %s 等待失败" % name)
            rc, logs = ctx.sh(["docker", "logs", name], timeout=60)
            m = re.search(r"throughput=([\d.]+) samples/s", logs)
            ctx.expect(m is not None, "容器 %s 无吞吐输出:\n%s"
                       % (name, logs[-600:]))
            results[name] = float(m.group(1))
            acc = re.search(r"acc=([\d.]+)", logs)
            ctx.steps.append("%s throughput=%.2f acc=%s"
                             % (name, results[name],
                                acc.group(1) if acc else "?"))
        mean = sum(results.values()) / len(results)
        print("        %s: 各实例 %s, 均值 %.2f" %
              (label, {k: round(v, 2) for k, v in results.items()}, mean))
        return results
    finally:
        for n in started:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)


@case("VDCU-TRN-006U", "整卡多实例（stream 基线）4 实例并发", "TRN",
      priority="P1")
def trn_006u(ctx):
    """TRN-006U：整卡 4 实例 = 多 stream 共享整卡（无隔离），作为对比基线。"""
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡前清理失败")
    mounts = [("tz_006u_%d" % i, None, {"HIP_VISIBLE_DEVICES": "0"})
              for i in range(1, 5)]
    _concurrent_train(ctx, mounts, "整卡 4 stream")


@case("VDCU-TRN-008U", "四等分 4 实例并发（vDCU 硬隔离）", "TRN",
      priority="P1")
def trn_008u(ctx):
    """TRN-008U：CU+显存四等分，4 实例各自独立 vDCU，验证算力硬隔离。"""
    vs = _make(ctx, ctx.card, 4)
    mounts = [("tz_008u_%d" % i, v["id"], {"HIP_VISIBLE_DEVICES": "0"})
              for i, v in enumerate(vs, 1)]
    _concurrent_train(ctx, mounts, "四等分 4 实例")


@case("VDCU-TRN-010U", "CU 共享四等分 4 实例并发", "TRN", priority="P1")
def trn_010u(ctx):
    """TRN-010U：CU 共享四等分（仅显存隔离），4 实例并发。"""
    vs = _make(ctx, ctx.card, 4, cu="0,0,0,0")
    mounts = [("tz_010u_%d" % i, v["id"], {"HIP_VISIBLE_DEVICES": "0"})
              for i, v in enumerate(vs, 1)]
    _concurrent_train(ctx, mounts, "CU 共享四等分 4 实例")


@case("VDCU-TRN-CMP", "stream vs vDCU 隔离并发性能对比结论", "TRN",
      priority="P1")
def trn_cmp(ctx):
    """TRN-CMP：整卡多 stream 与 CU 共享四等分并排对比。

    判据（建强文档结论）：两者单实例吞吐差异 ≤ 10% 判定"无明显差异"；
    CU 独立四等分（TRN-008U）显著低于 stream 组，体现算力硬隔离。
    """
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡前清理失败")
    stream_mounts = [("tz_cmp_s%d" % i, None, {"HIP_VISIBLE_DEVICES": "0"})
                     for i in range(1, 5)]
    stream = _concurrent_train(ctx, stream_mounts, "整卡 4 stream")

    vs = _make(ctx, ctx.card, 4, cu="0,0,0,0")
    share_mounts = [("tz_cmp_c%d" % i, v["id"], {"HIP_VISIBLE_DEVICES": "0"})
                    for i, v in enumerate(vs, 1)]
    share = _concurrent_train(ctx, share_mounts, "CU 共享四等分 4 实例")

    mean_s = sum(stream.values()) / len(stream)
    mean_c = sum(share.values()) / len(share)
    diff = abs(mean_s - mean_c) / mean_s * 100 if mean_s else 999
    ctx.steps.append("stream 均值 %.2f vs CU 共享均值 %.2f，差异 %.1f%%"
                     % (mean_s, mean_c, diff))
    ctx.expect(diff <= 10.0,
               "stream 与 CU 共享四等分差异 %.1f%% 超过 10%% 阈值" % diff)
    print("        结论：整卡多 stream 与 CU 共享四等分并发单实例性能差异 "
          "%.1f%%（≤10%% 判无明显差异）" % diff)


# ---------------------------------------------------------------- #
# v1.2 方案 5.4 节（大模型与小型训练）/ 6.2 节（推理）对应新增用例
# ---------------------------------------------------------------- #

def _train_result(out):
    m = re.search(r"TRAIN_RESULT (\{.*\})", out)
    if not m:
        raise AssertionError("未找到 TRAIN_RESULT，训练未完成\n" + out[-600:])
    return json.loads(m.group(1))


def _vl_result(out):
    m = re.search(r"VL_RESULT (\{.*\})", out)
    if not m:
        raise AssertionError("未找到 VL_RESULT，多模态推理未完成\n" + out[-600:])
    return json.loads(m.group(1))


def _emb_result(out):
    m = re.search(r"EMB_RESULT (\{.*\})", out)
    if not m:
        raise AssertionError("未找到 EMB_RESULT，Embedding 推理未完成\n" + out[-600:])
    return json.loads(m.group(1))


@case("VDCU-TRN-201", "Qwen3.5-0.8B 小模型切分效率（整卡/二等分/四等分）",
      "TRN", priority="P1")
def trn_201(ctx):
    """v1.2 5.4 TRN-201：小模型训练链路以 transformers 框架验证
    （Megatron 版 0.8B 预训练同受 TE 依赖限制，见 TRN-011）。"""
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    rc, out = ctx.docker(
        [], ctx.assets_cmd("qwen_train_smoke.py --steps 20"),
        env={"HIP_VISIBLE_DEVICES": "0"}, mount_assets=True, timeout=1800)
    ctx.expect(rc == 0, "整卡 0.8B 训练失败: %s" % out[-500:])
    base = _train_result(out)
    ctx.steps.append("整卡基线：tok/s=%.1f loss=%.4f" %
                     (base["tok_per_s"], base["loss"]))
    ctx.expect(base["loss"] < 10.0, "整卡训练 loss 异常: %s" % base["loss"])
    for label, count in (("二等分", 2), ("四等分", 4)):
        vs = None
        for c in (ctx.card, ctx.card2):
            try:
                vs = _make(ctx, c, count)
                break
            except AssertionError as e:
                ctx.steps.append("卡 %d 创建 %s 失败：%s"
                                 % (c, label, str(e)[:140]))
        ctx.expect(vs is not None, "%s 0.8B 训练前置创建失败" % label)
        rc, out = ctx.docker(
            [vs[0]["id"]], ctx.assets_cmd("qwen_train_smoke.py --steps 20"),
            env={"HIP_VISIBLE_DEVICES": "0"}, mount_assets=True, timeout=1800)
        ctx.expect(rc == 0, "%s 0.8B 训练失败: %s" % (label, out[-500:]))
        res = _train_result(out)
        eff = res["tok_per_s"] / base["tok_per_s"] * 100 \
            if base["tok_per_s"] else 0
        ctx.steps.append("%s vDCU：tok/s=%.1f loss=%.4f，相对整卡效率 %.0f%%"
                         % (label, res["tok_per_s"], res["loss"], eff))
        ctx.expect(res["loss"] < 10.0, "%s 训练 loss 异常" % label)
    print("        0.8B 小模型训练：整卡/二等分/四等分全部完成")


@case("VDCU-TRN-202", "Qwen3.5-0.8B 二等分 2 实例并发训练", "TRN",
      priority="P1")
def trn_202(ctx):
    """v1.2 5.4 TRN-202：多实例保持率（0.8B 小模型，二等分双实例）。"""
    vs = _make(ctx, ctx.card, 2)
    names = ["tz_trn202_a", "tz_trn202_b"]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    try:
        for name, vdev in zip(names, vs):
            rc, out = ctx.docker(
                [vdev["id"]], ctx.assets_cmd("qwen_train_smoke.py --steps 20"),
                name=name, detach=True, env={"HIP_VISIBLE_DEVICES": "0"},
                mount_assets=True, timeout=60)
            ctx.expect(rc == 0, "并发容器 %s 启动失败: %s" % (name, out[-300:]))
        for name in names:
            rc, out = ctx.sh(["docker", "wait", name], timeout=1800)
            ctx.expect(rc == 0, "并发容器 %s 等待失败" % name)
            rc, logs = ctx.sh(["docker", "logs", name], timeout=60)
            ctx.expect(rc == 0, "并发容器 %s 日志读取失败" % name)
            res = _train_result(logs)
            ctx.steps.append("%s tok/s=%.1f loss=%.4f" %
                             (name, res["tok_per_s"], res["loss"]))
            ctx.expect(res["loss"] < 10.0, "%s 训练 loss 异常" % name)
        print("        0.8B 二等分 2 实例并发训练完成")
    finally:
        for n in names:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-TRN-011", "Qwen3.6-27B 整卡单训练实例（Dense 基线）", "TRN",
      priority="P1")
def trn_011(ctx):
    """v1.2 5.4 TRN-011：dcu_megatron（core13.0）硬依赖 transformer_engine，
    TE 内网安装源在本环境不可达，Megatron 预训练链路受限
    （与 v1.2 实测一致）。数据集与 tokenizer 已就绪（oscar-1GB / Qwen3.5
    tokenizer），待 TE 可用后按 5.3 节公共步骤执行。"""
    ctx.steps.append("受限原因：dcu_megatron core13.0 flux 层规格硬依赖 "
                     "transformer_engine；TE 内网安装源"
                     "不可达，本地实现路径不可用。27B 模型（52 GB）已就绪。")
    ctx.skip("TE 安装源不可达，Megatron 大模型预训练受限（待 TE 可用后执行）")


@case("VDCU-TRN-101", "Qwen3.5-35B-A3B 跨卡 2 vDCU（EP2 + LoRA）", "TRN",
      priority="P1")
def trn_101(ctx):
    """v1.2 5.4 TRN-101：MoE 训练双重受限——vDCU 间集合通信不可用
    （ISSUE-013，见第 4 章）+ Megatron 依赖 TE（见 TRN-011）。"""
    ctx.steps.append("双重受限：① vDCU 间 all-to-all 集合通信不可用"
                     "（Duplicate GPU，ISSUE-013）；② Megatron 依赖 TE。"
                     "MoE 模型（67 GB）已就绪，基线待整卡 4 卡环境与 TE "
                     "就绪后执行。")
    ctx.skip("MoE 训练依赖 vDCU 间通信与 TE，双重受限（与 v1.2 实测一致）")


@case("VDCU-INF-101", "Qwen3.5-35B-A3B vDCU 推理（TP2+EP）", "INF",
      priority="P1")
def inf_101(ctx):
    """v1.2 6.2 INF-101/105：MoE 推理受限——vDCU 间通信不可用（ISSUE-013）
    且 35B 权重（约 70 GB）超出单 vDCU 配额（二等分 32 GB）。"""
    ctx.steps.append("双重受限：① TP/EP 依赖 vDCU 间集合通信（Duplicate GPU，"
                     "ISSUE-013）；② 35B 权重超出单 vDCU 显存配额。"
                     "单 vDCU 推理链路以 MiniCPM5/Qwen3.5-0.8B 稠密模型验证"
                     "（INF-001/002）。")
    ctx.skip("MoE 推理受通信与显存双重限制（与 v1.2 实测一致）")


@case("VDCU-INF-201", "Qwen3-VL-2B 多模态推理（整卡基线）", "INF",
      priority="P1")
def inf_201(ctx):
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    rc, out = ctx.docker([], ctx.assets_cmd("vl_smoke.py"),
                         env={"HIP_VISIBLE_DEVICES": "0"}, mount_assets=True,
                         timeout=1800)
    ctx.expect(rc == 0, "整卡多模态推理失败: %s" % out[-600:])
    res = _vl_result(out)
    ctx.expect(res["new_tokens"] >= 4, "多模态生成 token 数不足: %s"
               % res["new_tokens"])
    ctx.expect(res["tok_per_s"] > 0, "多模态生成吞吐异常")
    print("        vl-整卡: %s" % json.dumps(res, ensure_ascii=False))


@case("VDCU-INF-202", "Qwen3-VL-2B 多模态推理（二等分 vDCU）", "INF",
      priority="P1")
def inf_202(ctx):
    vs = _make(ctx, ctx.card, 2)
    rc, out = ctx.docker([vs[0]["id"]], ctx.assets_cmd("vl_smoke.py"),
                         env={"HIP_VISIBLE_DEVICES": "0"}, mount_assets=True,
                         timeout=1800)
    ctx.expect(rc == 0, "vDCU 多模态推理失败: %s" % out[-600:])
    res = _vl_result(out)
    expect_mem = vs[0]["mem"] / 1024 ** 3
    ctx.expect(res["new_tokens"] >= 4, "多模态生成 token 数不足: %s"
               % res["new_tokens"])
    ctx.expect(abs(res["visible_gb"] - expect_mem) < 0.5,
               "可见显存 %.2f GB，应为配额 %.2f GB"
               % (res["visible_gb"], expect_mem))
    print("        vl-vDCU: %s" % json.dumps(res, ensure_ascii=False))


@case("VDCU-INF-204", "Qwen3-Embedding-0.6B Embedding 推理（整卡基线）",
      "INF", priority="P1")
def inf_204(ctx):
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    rc, out = ctx.docker([], ctx.assets_cmd("emb_smoke.py"),
                         env={"HIP_VISIBLE_DEVICES": "0"}, mount_assets=True,
                         timeout=1800)
    ctx.expect(rc == 0, "整卡 Embedding 推理失败: %s" % out[-600:])
    res = _emb_result(out)
    ctx.expect(res["dim"] >= 512, "向量维度异常: %s" % res["dim"])
    print("        emb-整卡: %s" % json.dumps(res, ensure_ascii=False))


@case("VDCU-INF-205", "Qwen3-Embedding-0.6B Embedding 推理（二等分 vDCU）",
      "INF", priority="P1")
def inf_205(ctx):
    vs = _make(ctx, ctx.card, 2)
    rc, out = ctx.docker([vs[0]["id"]], ctx.assets_cmd("emb_smoke.py"),
                         env={"HIP_VISIBLE_DEVICES": "0"}, mount_assets=True,
                         timeout=1800)
    ctx.expect(rc == 0, "vDCU Embedding 推理失败: %s" % out[-600:])
    res = _emb_result(out)
    expect_mem = vs[0]["mem"] / 1024 ** 3
    ctx.expect(res["dim"] >= 512, "向量维度异常: %s" % res["dim"])
    ctx.expect(abs(res["visible_gb"] - expect_mem) < 0.5,
               "可见显存 %.2f GB，应为配额 %.2f GB"
               % (res["visible_gb"], expect_mem))
    print("        emb-vDCU: %s" % json.dumps(res, ensure_ascii=False))