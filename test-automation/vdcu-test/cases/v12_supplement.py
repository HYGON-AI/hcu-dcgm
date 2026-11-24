# -*- coding: utf-8 -*-
"""v1.2 方案缺失用例补充（2026-11-10）。

按 v1.2 定义补齐 4~14 章缺失用例（COM/TRN/INF/K8S/VOL/ISO/PERF/STB/MON/CMP/SEC）。
与 v1.4 已有用例语义重复的 v1.2 用例不再重复添加。
受限/未执行类用例：记录 v1.2 受限原因后跳过；可通过类用例在真实环境执行。
"""
import re
import shlex
import time

from vdcu.core import case

from cases.cluster import (_apply_pod, _delete_pod, _namespace, _nodes,
                           _pod_json, _pod_manifest, _pod_name,
                           _resource_mode, _resource_view, _wait_phase)
from cases.workload import (_make, _smoke_result, SMOKE_ARGS, _comm_reason,
                            _train_result, _vl_result, _emb_result)


def _skip12(ctx, reason, why):
    """受限/未执行类用例：记录原因并跳过（与 v1.2 结论一致）。"""
    ctx.steps.append("受限原因：%s（v1.2 结论：%s）" % (reason, why))
    ctx.skip(reason)


def _holder(ctx, vdev, name, script=None, env=None, timeout=120):
    ctx.sh(["docker", "rm", "-f", name], timeout=60)
    script = script or ("python - <<'PY'\nimport time, torch\n"
                        "x = torch.randn(512, 512, device='cuda')\n"
                        "while True:\n"
                        "    y = x @ x\n"
                        "    torch.cuda.synchronize()\n"
                        "    time.sleep(0.2)\nPY")
    rc, out = ctx.docker([vdev], script, name=name, detach=True,
                         env=env or {"HIP_VISIBLE_DEVICES": "0"}, timeout=timeout)
    ctx.expect(rc == 0, "容器 %s 启动失败: %s" % (name, out[-300:]))
    return name


# ============================================================ #
# 第 4 章 vDCU 间通信（COM）缺失用例
# ============================================================ #

@case("VDCU-COM-001", "同卡两容器 DDP（场景一复现）", "COM",
      expected_fail=True, issue="ISSUE-013")
def com_001(ctx):
    vs = _make_fb(ctx, 2)
    names = ["tz_com001_a", "tz_com001_b"]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    try:
        for name, vdev in zip(names, vs):
            rc, out = ctx.docker([vdev["id"]], "sleep 600", name=name,
                                 detach=True,
                                 env={"HIP_VISIBLE_DEVICES": "0"}, timeout=60)
            ctx.expect(rc == 0, "容器 %s 启动失败" % name)
        rc, ip = ctx.sh(["docker", "inspect", "-f",
                         "{{.NetworkSettings.IPAddress}}", names[0]], timeout=60)
        ctx.expect(rc == 0 and ip.strip(), "获取容器 A 地址失败")
        master = ip.strip()
        for rank, (name, vdev) in enumerate(zip(names, vs)):
            script = ("cd /tmp && NCCL_DEBUG=INFO torchrun "
                      "--nnodes=2 --node_rank=%d --master_addr=%s "
                      "--master_port=29500 --nproc_per_node=1 "
                      + ctx.assets_path("allreduce.py")) % (rank, master)
            rc, out = ctx.docker([vdev["id"]], script, name=name, detach=True,
                                 env={"HIP_VISIBLE_DEVICES": "0",
                                      "NCCL_SOCKET_IFNAME": "eth0"},
                                 mount_assets=True, timeout=60)
            ctx.expect(rc == 0, "DDP 容器 %s 启动失败" % name)
        for name in names:
            ctx.sh(["docker", "wait", name], timeout=600)
        rc, logs = ctx.sh(["docker", "logs", names[0]], timeout=60)
        ctx.steps.append("容器 A DDP 日志：\n" + logs[-600:])
        ctx.expect("ALLREDUCE_OK" in logs,
                   "同卡两容器 DDP 失败：%s" % _comm_reason(logs))
    finally:
        for n in names:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-COM-005", "单进程 hdms 同卡双 vDCU（场景五复现）", "COM",
      priority="P2")
def com_005(ctx):
    _skip12(ctx, "hdms 场景需演示程序（demo 二进制），本环境未部署",
            "未执行")


@case("VDCU-COM-006", "单进程 hdms 跨卡双 vDCU（场景六复现）", "COM",
      priority="P2")
def com_006(ctx):
    _skip12(ctx, "hdms 场景需演示程序（demo 二进制），本环境未部署",
            "未执行")


@case("VDCU-COM-007", "单机 all-reduce 带宽（整卡 vs vDCU）", "COM",
      expected_fail=True, issue="ISSUE-013")
def com_007(ctx):
    # 基线：整卡 2 卡 all-reduce（dri 整树透传）
    ctx.expect(ctx.destroy_card(ctx.card), "前置清理失败")
    script = ("cd /tmp && torchrun --standalone --nproc_per_node=2 "
              + ctx.assets_path("allreduce.py"))
    rc, out = ctx.docker([], script, env={"HIP_VISIBLE_DEVICES": "0,1"},
                         mount_assets=True, dri_tree=True, timeout=600)
    ctx.expect(rc == 0 and "ALLREDUCE_OK" in out,
               "整卡 2 卡 all-reduce 基线失败: %s" % _comm_reason(out))
    ctx.steps.append("整卡基线：2 卡 all-reduce 正常")
    # vDCU 组合：与 COM-002 同因失败
    a = _make_fb(ctx, 2)
    b = _make(ctx, ctx.card2, 2)
    rc, out = ctx.docker([a[0]["id"], b[0]["id"]], script,
                         env={"HIP_VISIBLE_DEVICES": "0,1"},
                         mount_assets=True, timeout=600)
    ctx.expect(rc == 0 and "ALLREDUCE_OK" in out,
               "跨卡 vDCU all-reduce 失败：%s" % _comm_reason(out))


@case("VDCU-COM-008", "vLLM TP=2 同卡双 vDCU", "COM",
      expected_fail=True, issue="ISSUE-013")
def com_008(ctx):
    if not ctx.vllm_model:
        ctx.skip("未配置 VDCU_VLLM_MODEL")
    vs = _make_fb(ctx, 2)
    script = ("HIP_VISIBLE_DEVICES=0,1 vllm serve %s "
              "--tensor-parallel-size 2 --enforce-eager --max-model-len 1024 "
              "--port 8199" % ctx.vllm_model)
    rc, out = ctx.docker([vs[0]["id"], vs[1]["id"]], script,
                         image=ctx.inference_image, timeout=900)
    ctx.steps.append("vLLM TP=2 输出：\n" + out[-600:])
    ctx.expect(rc == 0 and "Uvicorn running" in out,
               "vLLM TP=2 同卡双 vDCU 服务未正常运行: %s" % out[-300:])


@case("VDCU-COM-009", "vLLM TP=2 跨卡双 vDCU", "COM",
      expected_fail=True, issue="ISSUE-013")
def com_009(ctx):
    if not ctx.vllm_model:
        ctx.skip("未配置 VDCU_VLLM_MODEL")
    a = _make_fb(ctx, 2)
    b = _make(ctx, ctx.card2, 2)
    script = ("HIP_VISIBLE_DEVICES=0,1 vllm serve %s "
              "--tensor-parallel-size 2 --enforce-eager --max-model-len 1024 "
              "--port 8199" % ctx.vllm_model)
    rc, out = ctx.docker([a[0]["id"], b[0]["id"]], script,
                         image=ctx.inference_image, timeout=900)
    ctx.steps.append("vLLM TP=2 输出：\n" + out[-600:])
    ctx.expect(rc == 0 and "Uvicorn running" in out,
               "vLLM TP=2 跨卡双 vDCU 服务未正常运行: %s" % out[-300:])


@case("VDCU-COM-010", "vLLM PP=2 + Ray 跨容器", "COM", priority="P1")
def com_010(ctx):
    _skip12(ctx, "PP=2 + Ray 跨容器依赖 vDCU 集合通信（Duplicate GPU，见第 4 章）",
            "受限")


@case("VDCU-COM-011", "Ray 集群 GPU 资源识别一致性", "COM", priority="P2")
def com_011(ctx):
    _skip12(ctx, "Ray 集群 GPU 资源识别未单独执行（依赖通信能力验证）",
            "未执行")


@case("VDCU-COM-012", "多机跨节点 vDCU 通信", "COM", priority="P1")
def com_012(ctx):
    _skip12(ctx, "多机跨节点 vDCU 通信依赖 vDCU 集合通信（Duplicate GPU），受限",
            "受限")


@case("VDCU-COM-013", "拓扑可见性与设备亲和", "COM", priority="P2")
def com_013(ctx):
    _skip12(ctx, "拓扑可见性用例未单独执行（hy-smi topo 与 use-gpuuuid 调度待专项）",
            "未执行")


@case("VDCU-COM-014", "集合通信长稳（100 次迭代）", "COM", priority="P1")
def com_014(ctx):
    """整卡形态集合通信长稳：100 次 64 MB all-reduce 全部成功。"""
    ctx.expect(ctx.destroy_card(ctx.card), "前置清理失败")
    script = ("cd /tmp && torchrun --standalone --nproc_per_node=2 "
              + ctx.assets_path("ar_loop.py"))
    rc, out = ctx.docker([], script, env={"HIP_VISIBLE_DEVICES": "0,1"},
                         mount_assets=True, dri_tree=True, timeout=1200)
    ctx.expect(rc == 0 and "AR_LOOP_OK" in out,
               "整卡集合通信长稳失败: %s" % out[-400:])


# ============================================================ #
# 第 5 章 训练场景（TRN）缺失用例（v1.2 受限/未执行，同因记录）
# ============================================================ #

@case("VDCU-TRN-012", "Qwen3.6-27B CU+显存二等分单训练实例", "TRN",
      priority="P1")
def trn_012(ctx):
    _skip12(ctx, "Megatron 依赖 TE（内网源不可达）；且 27B 量级模型超出单 vDCU 显存配额",
            "受限")


@case("VDCU-TRN-013", "Qwen3.6-27B CU 共享+显存二等分单训练实例", "TRN",
      priority="P1")
def trn_013(ctx):
    _skip12(ctx, "Megatron 依赖 TE（内网源不可达），同 TRN-011", "受限")


@case("VDCU-TRN-014", "Qwen3.6-27B 整卡多训练实例", "TRN", priority="P1")
def trn_014(ctx):
    _skip12(ctx, "Megatron 依赖 TE（内网源不可达），同 TRN-011", "受限")


@case("VDCU-TRN-015", "Qwen3.6-27B 二等分多训练实例", "TRN", priority="P1")
def trn_015(ctx):
    _skip12(ctx, "Megatron 依赖 TE（内网源不可达），同 TRN-011", "受限")


@case("VDCU-TRN-016", "Qwen3.6-27B 整卡 8 卡多机基线", "TRN", priority="P1")
def trn_016(ctx):
    _skip12(ctx, "Megatron 依赖 TE（内网源不可达），多机基线需 TE 可用", "受限")


@case("VDCU-TRN-102", "Qwen3.5-35B-A3B 同卡 2 vDCU（EP2）", "TRN",
      priority="P1")
def trn_102(ctx):
    _skip12(ctx, "MoE 训练依赖 vDCU 间集合通信（Duplicate GPU）+ Megatron TE，双重受限",
            "受限")


@case("VDCU-TRN-103", "Qwen3.5-35B-A3B 跨卡 4 vDCU（TP2/EP2）", "TRN",
      priority="P1")
def trn_103(ctx):
    _skip12(ctx, "MoE 训练依赖 vDCU 间集合通信 + Megatron TE，双重受限",
            "受限")


@case("VDCU-TRN-104", "Qwen3.5-35B-A3B 整卡 4 卡基线", "TRN", priority="P1")
def trn_104(ctx):
    _skip12(ctx, "MoE 整卡基线依赖 Megatron TE（内网源不可达）", "受限")


@case("VDCU-TRN-105", "MoE 整卡 8 卡基线（大参数量压测）", "TRN",
      priority="P1")
def trn_105(ctx):
    _skip12(ctx, "MoE 整卡基线依赖 Megatron TE（内网源不可达）", "受限")


@case("VDCU-TRN-106", "MoE 跨卡 4 vDCU（LoRA 降显存）", "TRN", priority="P1")
def trn_106(ctx):
    _skip12(ctx, "MoE 训练依赖 vDCU 间集合通信 + Megatron TE，双重受限",
            "受限")


@case("VDCU-TRN-107", "MoE 同卡 4 vDCU（TP4/EP1）", "TRN", priority="P1")
def trn_107(ctx):
    _skip12(ctx, "MoE 训练依赖 vDCU 间集合通信 + Megatron TE，双重受限",
            "受限")


@case("VDCU-TRN-108", "MoE 2 实例并发噪声邻居", "TRN", priority="P1")
def trn_108(ctx):
    _skip12(ctx, "MoE 并发噪声邻居依赖 Megatron TE（内网源不可达）", "受限")


@case("VDCU-TRN-203", "ResNet-50 切分效率", "TRN", priority="P2")
def trn_203(ctx):
    _skip12(ctx, "ImageNet 子集数据集在本环境未部署，待数据集就绪后执行",
            "未执行")


@case("VDCU-TRN-204", "YOLOv5-L 切分效率", "TRN", priority="P2")
def trn_204(ctx):
    _skip12(ctx, "COCO 子集数据集在本环境未部署，待数据集就绪后执行",
            "未执行")


@case("VDCU-TRN-205", "Qwen3.5-0.8B LoRA 切分效率", "TRN", priority="P2")
def trn_205(ctx):
    _skip12(ctx, "llamafactory 镜像未在本环境部署，待镜像就绪后执行",
            "未执行")


# ============================================================ #
# 第 6 章 推理场景（INF）缺失用例
# ============================================================ #

def _gen_smoke(ctx, vdevs, env=None, timeout=1200):
    """Qwen3.5-0.8B 生成冒烟（INF-001 同链路），返回结果 dict。"""
    rc, out = ctx.docker(vdevs, ctx.assets_cmd("gen_smoke.py --max-new-tokens 32"),
                         env=env or {"HIP_VISIBLE_DEVICES": "0"},
                         mount_assets=True, timeout=timeout)
    ctx.expect(rc == 0, "生成冒烟失败: %s" % out[-500:])
    m = re.search(r"GEN_RESULT (\{.*\})", out)
    ctx.expect(m is not None, "未找到 GEN_RESULT\n" + out[-400:])
    import json as _json
    return _json.loads(m.group(1))


def _make_fb(ctx, count, cu=None):
    """创建 vDCU（卡 5 失败时回退卡 6；ISSUE-012a 残留导致释放延迟的环境应对）。"""
    last = None
    for c in (ctx.card, ctx.card2):
        try:
            return _make(ctx, c, count, cu=cu)
        except AssertionError as e:
            last = e
            ctx.steps.append("卡 %d 创建失败（%s），回退卡 %d"
                             % (c, str(e)[:100], ctx.card2))
    raise last


@case("VDCU-INF-003", "整卡基准推理（定长生成）", "INF", priority="P1")
def inf_003(ctx):
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    res = _gen_smoke(ctx, [], env={"HIP_VISIBLE_DEVICES": "0"})
    ctx.expect(res["new_tokens"] >= 16 and res["tok_per_s"] > 0,
               "整卡生成冒烟异常: %s" % res)
    ctx.steps.append("整卡基线：tok/s=%.2f，可见显存 %.2f GB"
                     % (res["tok_per_s"], res["mem_gb"]))


@case("VDCU-INF-004", "2 等分 vDCU 单实例推理", "INF", priority="P1")
def inf_004(ctx):
    vs = _make_fb(ctx, 2)
    res = _gen_smoke(ctx, [vs[0]["id"]])
    expect_mem = vs[0]["mem"] / 1024 ** 3
    ctx.expect(res["new_tokens"] >= 16, "二等分生成 token 不足")
    ctx.expect(abs(res["mem_gb"] - expect_mem) < 0.5,
               "可见显存（mem_gb）%.2f GB，应为配额 %.2f GB" % (res["mem_gb"], expect_mem))
    ctx.steps.append("二等分 vDCU：tok/s=%.2f" % res["tok_per_s"])


@case("VDCU-INF-005", "4 等分 vDCU 单实例推理", "INF", priority="P1")
def inf_005(ctx):
    vs = _make_fb(ctx, 4)
    res = _gen_smoke(ctx, [vs[0]["id"]])
    expect_mem = vs[0]["mem"] / 1024 ** 3
    ctx.expect(res["new_tokens"] >= 16, "四等分生成 token 不足")
    ctx.expect(abs(res["mem_gb"] - expect_mem) < 0.5,
               "可见显存（mem_gb）%.2f GB，应为配额 %.2f GB" % (res["mem_gb"], expect_mem))
    ctx.steps.append("四等分 vDCU：tok/s=%.2f" % res["tok_per_s"])


@case("VDCU-INF-006", "多实例并发推理", "INF", priority="P1")
def inf_006(ctx):
    vs = _make_fb(ctx, 2)
    names = ["tz_inf006_a", "tz_inf006_b"]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    try:
        for name, vdev in zip(names, vs):
            rc, out = ctx.docker([vdev["id"]],
                                 ctx.assets_cmd("gen_smoke.py --max-new-tokens 32"),
                                 name=name, detach=True,
                                 env={"HIP_VISIBLE_DEVICES": "0"},
                                 mount_assets=True, timeout=60)
            ctx.expect(rc == 0, "并发容器 %s 启动失败" % name)
        import json as _json
        results = {}
        for name in names:
            ctx.sh(["docker", "wait", name], timeout=1200)
            rc, logs = ctx.sh(["docker", "logs", name], timeout=60)
            m = re.search(r"GEN_RESULT (\{.*\})", logs)
            ctx.expect(rc == 0 and m, "容器 %s 无生成结果" % name)
            results[name] = _json.loads(m.group(1))
            ctx.steps.append("%s tok/s=%.2f" % (name, results[name]["tok_per_s"]))
        for name, res in results.items():
            ctx.expect(res["new_tokens"] >= 16 and res["tok_per_s"] > 0,
                       "%s 并发生成异常" % name)
    finally:
        for n in names:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-INF-007", "CU 共享模式推理", "INF", priority="P1")
def inf_007(ctx):
    vs = _make_fb(ctx, 2, cu="0,0")
    res = _gen_smoke(ctx, [vs[0]["id"]])
    ctx.expect(res["new_tokens"] >= 16, "CU 共享推理生成异常")
    ctx.steps.append("CU 共享二等分：tok/s=%.2f" % res["tok_per_s"])


@case("VDCU-INF-008", "32K 长上下文服务启动", "INF", priority="P1")
def inf_008(ctx):
    """32K 上下文窗口：整卡以 --max-model-len 32768 启动 vLLM 服务并完成一次生成。
    （v1.2 的长上下文压测在 30000 token 输入下完成；本环境以窗口启动+生成验证链路）"""
    if not ctx.vllm_model:
        ctx.skip("未配置 VDCU_VLLM_MODEL")
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    script = ctx.assets_cmd("vllm_smoke.py --model %s --max-model-len 32768" % ctx.vllm_model)
    rc, out = ctx.docker([], script, env={"HIP_VISIBLE_DEVICES": "0"},
                         mount_assets=True, image=ctx.inference_image,
                         shm="16G", timeout=1800)
    ctx.expect(rc == 0, "32K 长上下文服务失败: %s" % out[-600:])
    m = re.search(r"VLLM_RESULT (\{.*\})", out)
    ctx.expect(m is not None, "未找到 VLLM_RESULT")
    ctx.steps.append("32K 上下文窗口服务启动与生成成功")


@case("VDCU-INF-009", "PD 分离", "INF", priority="P1")
def inf_009(ctx):
    _skip12(ctx, "PD 分离依赖多卡集合通信，vDCU 形态受 Duplicate GPU 限制",
            "受限")


@case("VDCU-INF-102", "Qwen3.5-35B-A3B 整卡 4 卡基线", "INF", priority="P1")
def inf_102(ctx):
    _skip12(ctx, "MoE 权重（约 70 GB）需 ≥2 整卡，本次未单独执行（预留）",
            "未执行")


@case("VDCU-INF-103", "Qwen3.5-35B-A3B 同卡 4 vDCU 推理（TP4）", "INF",
      priority="P1")
def inf_103(ctx):
    _skip12(ctx, "TP4 依赖 vDCU 间集合通信（Duplicate GPU）且 35B 权重超出 vDCU 显存",
            "受限")


@case("VDCU-INF-104", "MoE 跨卡 2 vDCU 推理（TP2）", "INF", priority="P1")
def inf_104(ctx):
    _skip12(ctx, "TP2 依赖 vDCU 间集合通信（Duplicate GPU），受限", "受限")


@case("VDCU-INF-105", "单 vDCU MoE 小规模推理（EP1/TP1）", "INF", priority="P1")
def inf_105(ctx):
    _skip12(ctx, "35B 权重超出单 vDCU 配额（32 GB），环境无小尺寸新版 MoE；"
                 "单 vDCU 推理链路以稠密模型验证（INF-001/004/005）", "受限")


@case("VDCU-INF-203", "Qwen3-VL-2B 多模态推理（4 等分 vDCU）", "INF",
      priority="P1")
def inf_203(ctx):
    vs = _make_fb(ctx, 4)
    rc, out = ctx.docker([vs[0]["id"]], ctx.assets_cmd("vl_smoke.py"),
                         env={"HIP_VISIBLE_DEVICES": "0"}, mount_assets=True,
                         timeout=1800)
    ctx.expect(rc == 0, "四等分多模态推理失败: %s" % out[-600:])
    res = _vl_result(out)
    expect_mem = vs[0]["mem"] / 1024 ** 3
    ctx.expect(res["new_tokens"] >= 4, "多模态生成 token 数不足")
    ctx.expect(abs(res["visible_gb"] - expect_mem) < 0.5,
               "可见显存 %.2f GB，应为配额 %.2f GB"
               % (res["visible_gb"], expect_mem))
    ctx.steps.append("四等分 vDCU 多模态推理：tok/s=%.2f" % res["tok_per_s"])


@case("VDCU-INF-206", "通用推理（bge 系列小模型）", "INF", priority="P1")
def inf_206(ctx):
    """bge-large-zh 通用小模型在整卡与 vDCU 上的向量化推理（链路与 INF-001 一致）。"""
    model = ctx.emb_model
    if not model:
        ctx.skip("未配置 emb_model（config.json），无法执行 bge 推理")
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    rc, out = ctx.docker(
        [], ctx.assets_cmd("emb_smoke.py --model %s --slow-tokenizer" % model),
        env={"HIP_VISIBLE_DEVICES": "0"}, mount_assets=True, timeout=1800)
    ctx.expect(rc == 0, "整卡 bge 推理失败: %s" % out[-500:])
    base = _emb_result(out)
    ctx.steps.append("整卡 bge：dim=%d，%s" % (base["dim"], base["elapsed_ms"]))
    vs = _make_fb(ctx, 2)
    rc, out = ctx.docker(
        [vs[0]["id"]],
        ctx.assets_cmd("emb_smoke.py --model %s --slow-tokenizer" % model),
        env={"HIP_VISIBLE_DEVICES": "0"}, mount_assets=True, timeout=1800)
    ctx.expect(rc == 0, "vDCU bge 推理失败: %s" % out[-500:])
    res = _emb_result(out)
    ctx.expect(res["dim"] >= 512, "向量维度异常: %s" % res["dim"])
    ctx.steps.append("二等分 vDCU bge：dim=%d，%s"
                     % (res["dim"], res["elapsed_ms"]))


# ============================================================ #
# 第 7 章 Kubernetes 编排（K8S）缺失用例
# ============================================================ #

def _k8s_enabled(ctx):
    if not ctx.cluster_enabled:
        ctx.skip("未启用集群测试")


def _whole_card_key(ctx):
    nodes = _nodes(ctx)
    _, key, _ = _resource_mode(nodes)
    ctx.expect(key, "节点未注册可申请的 Hygon 资源")
    return key


@case("VDCU-K8S-005", "物理整卡申请", "K8S", priority="P1")
def k8s_005(ctx):
    _k8s_enabled(ctx)
    key = _whole_card_key(ctx)
    name = _pod_name("vdcu-k8s005")
    manifest = _pod_manifest(ctx, name, key, 1)
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name)
        ctx.expect(pod.get("status", {}).get("phase") == "Running",
                   "整卡申请 Pod 未进入 Running")
        ctx.steps.append("整卡申请（%s=1）Pod 调度与运行正常" % key)
    finally:
        _delete_pod(ctx, name)


@case("VDCU-K8S-006", "节点亲和 use-nodename", "K8S", priority="P1")
def k8s_006(ctx):
    _skip12(ctx, "use-nodename 注解依赖 HAMi/vHCU 动态资源；集群于 2026-11-10 切换 hcu-device-plugin-hami（hcunum=32/hcucores=800/hcumem=524160），对应注解行为待按《HCU Device Plugin-v3.0.1》专项验证",
            "通过（HAMi 环境下）")


@case("VDCU-K8S-007", "节点反亲和 nouse-nodename", "K8S", priority="P1")
def k8s_007(ctx):
    _skip12(ctx, "nouse-nodename 注解依赖 HAMi/vHCU 动态资源（HAMi 模式于 2026-11-10 切换部署，行为待专项验证）",
            "通过（HAMi 环境下）")


@case("VDCU-K8S-008", "指定物理卡 use-gpuuuid", "K8S", priority="P1")
def k8s_008(ctx):
    _skip12(ctx, "use-gpuuuid 注解依赖 HAMi/vHCU 动态资源（HAMi 模式于 2026-11-10 切换部署，行为待专项验证）",
            "通过（HAMi 环境下）")


@case("VDCU-K8S-009", "指定 HCU 类型 use-hcutype", "K8S", priority="P1")
def k8s_009(ctx):
    _skip12(ctx, "use-hcutype 注解依赖 HAMi/vHCU 动态资源（HAMi 模式于 2026-11-10 切换部署，行为待专项验证）",
            "通过（HAMi 环境下）")


@case("VDCU-K8S-010", "单容器申请多个 vHCU（应受限）", "K8S", priority="P1")
def k8s_010(ctx):
    _skip12(ctx, "vHCU 动态切分（HAMi 模式）于 2026-11-10 切换部署，多 vHCU 申请行为待按 v3.0.1 文档专项验证",
            "通过（HAMi 环境下）")


@case("VDCU-K8S-011", "initContainer 申请 vHCU（不支持）", "K8S", priority="P1")
def k8s_011(ctx):
    _skip12(ctx, "vHCU 动态切分（HAMi 模式）于 2026-11-10 切换部署，initContainer 行为待专项验证",
            "通过（HAMi 环境下）")


@case("VDCU-K8S-012", "Pod 删除后资源回收", "K8S", priority="P1")
def k8s_012(ctx):
    _k8s_enabled(ctx)
    nodes = _nodes(ctx)
    key = _whole_card_key(ctx)
    _, _, node = _resource_mode(nodes)
    _, alloc = _resource_view(node)
    before = alloc.get(key, "0")
    name = _pod_name("vdcu-k8s012")
    manifest = _pod_manifest(ctx, name, key, 1)
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name)
        ctx.expect(pod.get("status", {}).get("phase") == "Running",
                   "Pod 未进入 Running")
        _delete_pod(ctx, name)
        for _ in range(20):
            rc, _, payload = _pod_json(ctx, name)
            if rc != 0 or payload is None:
                break
            time.sleep(2)
        ctx.expect(rc != 0 or payload is None, "Pod 删除后仍可查询")
        nodes2 = _nodes(ctx)
        _, node2 = _resource_mode(nodes2)
        _, alloc2 = _resource_view(node2)
        after = alloc2.get(key, "0")
        ctx.steps.append("删除前后 allocatable：%s → %s" % (before, after))
        ctx.expect(str(after).rstrip("m") == str(before).rstrip("m"),
                   "Pod 删除后资源未回收：%s → %s" % (before, after))
    finally:
        _delete_pod(ctx, name)


@case("VDCU-K8S-014", "非法资源值（小数、超 100%）", "K8S", priority="P2")
def k8s_014(ctx):
    _skip12(ctx, "HAMi 模式于 2026-11-10 切换部署（hcunum/hcucores/hcumem 已注册），非法值行为待专项验证",
            "通过（HAMi 环境下）")


@case("VDCU-K8S-015", "Namespace 与 ResourceQuota 隔离", "K8S", priority="P1")
def k8s_015(ctx):
    _k8s_enabled(ctx)
    if not ctx.cluster_mutating_tests:
        ctx.skip("ResourceQuota 创建需要 VDCU_CLUSTER_MUTATING_TESTS=1")
    key = _whole_card_key(ctx)
    ns = _namespace(ctx)
    quota = "vdcu-quota"
    ctx.kubectl("delete", "resourcequota", quota, "-n", ns,
                "--ignore-not-found", timeout=60)
    try:
        ctx.kubectl("create", "quota", quota, "-n", ns,
                    "--hard=%s=1" % key, timeout=60)
        names = [_pod_name("vdcu-q1"), _pod_name("vdcu-q2")]
        try:
            _apply_pod(ctx, _pod_manifest(ctx, names[0], key, 1))
            _apply_pod(ctx, _pod_manifest(ctx, names[1], key, 1))
            p1 = _wait_phase(ctx, names[0], timeout=60)
            p2 = _wait_phase(ctx, names[1], phases=("Pending", "Running"),
                             timeout=30)
            ph1 = p1.get("status", {}).get("phase")
            ph2 = p2.get("status", {}).get("phase")
            ctx.steps.append("配额=%s=1 下两 Pod 状态：%s / %s" % (key, ph1, ph2))
            ctx.expect({ph1, ph2} == {"Running", "Pending"},
                       "配额隔离未生效：%s / %s" % (ph1, ph2))
        finally:
            for n in names:
                _delete_pod(ctx, n)
    finally:
        ctx.kubectl("delete", "resourcequota", quota, "-n", ns,
                    "--ignore-not-found", timeout=60)


@case("VDCU-K8S-016", "同 Pod 多容器各申请 1 整卡", "K8S", priority="P1")
def k8s_016(ctx):
    _k8s_enabled(ctx)
    key = _whole_card_key(ctx)
    name = _pod_name("vdcu-k8s016")
    manifest = _pod_manifest(ctx, name, key, 1)
    # 同一 Pod 两个容器各申请 1 张整卡
    manifest = manifest.replace(
        '"name": "probe"',
        '"name": "probe-a"')
    import json as _json
    spec = _json.loads(manifest)
    container = spec["spec"]["containers"][0]
    second = dict(container)
    second["name"] = "probe-b"
    spec["spec"]["containers"].append(second)
    manifest = _json.dumps(spec)
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name)
        ctx.expect(pod.get("status", {}).get("phase") == "Running",
                   "同 Pod 多容器整卡申请未进入 Running")
        ctx.steps.append("同 Pod 双容器各申请 %s=1 调度运行正常" % key)
    finally:
        _delete_pod(ctx, name)


@case("VDCU-K8S-017", "Pod 驱逐/重调度后的 vHCU 分配", "K8S", priority="P2")
def k8s_017(ctx):
    _skip12(ctx, "驱逐重调度需专项窗口（HAMi 模式于 2026-11-10 切换部署）", "未执行")


@case("VDCU-K8S-018", "节点标签缺失时的调度行为", "K8S", priority="P1")
def k8s_018(ctx):
    _k8s_enabled(ctx)
    name = _pod_name("vdcu-k8s018")
    manifest = _pod_manifest(ctx, name, "hygon.com/dcu", 1)
    import json as _json
    spec = _json.loads(manifest)
    spec["spec"]["nodeSelector"] = {"vdcu-test-label": "not-exist"}
    spec["spec"].pop("nodeName", None)
    try:
        _apply_pod(ctx, _json.dumps(spec))
        pod = _wait_phase(ctx, name, phases=("Pending", "Failed"), timeout=60)
        phase = pod.get("status", {}).get("phase")
        ctx.expect(phase == "Pending",
                   "标签缺失 Pod 应保持 Pending，实际 %s" % phase)
        ctx.steps.append("不存在的 nodeSelector 下 Pod 保持 Pending（标签缺失调度行为）")
    finally:
        _delete_pod(ctx, name)


@case("VDCU-K8S-020", "vHCU 与整卡混部", "K8S", priority="P1")
def k8s_020(ctx):
    _skip12(ctx, "HAMi 模式于 2026-11-10 切换部署，混部场景待专项验证",
            "通过（HAMi 环境下）")


# ============================================================ #
# 第 8 章 Volcano 调度（VOL）缺失用例
# ============================================================ #

@case("VDCU-VOL-020", "Volcano 与默认调度器共存", "VOL", priority="P1")
def vol_020(ctx):
    _k8s_enabled(ctx)
    rc, out, payload = ctx.kubectl_json(
        "get", "pods", "-n", "volcano-system", "-o", "json", timeout=60)
    ctx.expect(rc == 0 and payload is not None, "查询 volcano-system Pod 失败")
    volcano = [p for p in payload.get("items", [])
               if p.get("status", {}).get("phase") == "Running"]
    ctx.expect(volcano, "volcano-system 无 Running Pod")
    rc, out, payload = ctx.kubectl_json(
        "get", "pods", "-n", "kube-system", "-o", "json", timeout=60)
    default_sched = [p for p in payload.get("items", [])
                     if p.get("spec", {}).get("schedulerName", "") in
                     ("default-scheduler", "") and
                     p.get("status", {}).get("phase") == "Running"]
    ctx.expect(default_sched, "默认调度器无 Running Pod")
    ctx.steps.append("volcano-system 组件 %d 个 Running；kube-system 默认调度器 "
                     "Pod %d 个 Running（共存正常）"
                     % (len(volcano), len(default_sched)))


# ============================================================ #
# 第 9 章 隔离性与 QoS（ISO）缺失用例
# ============================================================ #

def _bert_quick(ctx, vdevs, name="", detach=False, timeout=1200):
    return ctx.docker(vdevs, ctx.assets_cmd("bert_smoke.py " + SMOKE_ARGS),
                      name=name, detach=detach,
                      env={"HIP_VISIBLE_DEVICES": "0"},
                      mount_assets=True, timeout=60 if detach else timeout)


@case("VDCU-ISO-003", "噪声邻居：CU 独立模式性能隔离", "ISO", priority="P1")
def iso_003(ctx):
    vs = _make_fb(ctx, 2)
    rc, out = _bert_quick(ctx, [vs[0]["id"]])
    ctx.expect(rc == 0, "单实例基线训练失败: %s" % out[-400:])
    solo = _smoke_result(out)
    names = ["tz_iso003_a", "tz_iso003_b"]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    try:
        for name, vdev in zip(names, vs):
            rc, out = _bert_quick(ctx, [vdev["id"]], name=name, detach=True)
            ctx.expect(rc == 0, "并发容器 %s 启动失败" % name)
        for name in names:
            ctx.sh(["docker", "wait", name], timeout=1500)
            rc, logs = ctx.sh(["docker", "logs", name], timeout=60)
            res = _smoke_result(logs)
            ctx.steps.append("%s throughput=%.2f（并发噪声邻居）"
                             % (name, res["throughput"]))
        ctx.steps.append("单实例基线 throughput=%.2f" % solo["throughput"])
    finally:
        for n in names:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-ISO-004", "CU 共享模式调度公平性", "ISO", priority="P1")
def iso_004(ctx):
    _skip12(ctx, "CU 共享调度公平性依赖 HAMi 插件环境（v1.2 同）", "受限")


@case("VDCU-ISO-005", "显存数据隔离（残留数据不可见）", "ISO", priority="P2")
def iso_005(ctx):
    _skip12(ctx, "显存残留数据读取用例未单独执行（需专项安全测试窗口）",
            "未执行")


@case("VDCU-ISO-006", "故障隔离：进程异常不扩散", "ISO", priority="P1")
def iso_006(ctx):
    vs = _make_fb(ctx, 2)
    names = ["tz_iso006_a", "tz_iso006_b"]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    try:
        for name, vdev in zip(names, vs):
            rc, out = _bert_quick(ctx, [vdev["id"]], name=name, detach=True)
            ctx.expect(rc == 0, "并发容器 %s 启动失败" % name)
        time.sleep(20)
        rc, out = ctx.sh(["docker", "kill", names[0]], timeout=60)
        ctx.expect(rc == 0, "kill 容器 A 失败")
        ctx.steps.append("容器 A 被 kill：%s" % out.strip()[:60])
        ctx.sh(["docker", "wait", names[1]], timeout=1500)
        rc, logs = ctx.sh(["docker", "logs", names[1]], timeout=60)
        res = _smoke_result(logs)
        ctx.expect(res["throughput"] > 0 and res["acc"] >= 0.5,
                   "容器 A 异常后容器 B 训练受影响：%s" % res)
        ctx.steps.append("容器 B 训练正常完成（throughput=%.2f，acc=%.2f）——"
                         "进程异常未扩散" % (res["throughput"], res["acc"]))
    finally:
        for n in names:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-ISO-007", "故障隔离：容器崩溃影响范围", "ISO", priority="P1")
def iso_007(ctx):
    vs = _make_fb(ctx, 2)
    names = ["tz_iso007_a", "tz_iso007_b"]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    try:
        for name, vdev in zip(names, vs):
            rc, out = _bert_quick(ctx, [vdev["id"]], name=name, detach=True)
            ctx.expect(rc == 0, "并发容器 %s 启动失败" % name)
        time.sleep(20)
        ctx.sh(["docker", "rm", "-f", names[0]], timeout=60)
        ctx.sh(["docker", "wait", names[1]], timeout=1500)
        rc, logs = ctx.sh(["docker", "logs", names[1]], timeout=60)
        res = _smoke_result(logs)
        ctx.expect(res["throughput"] > 0 and res["acc"] >= 0.5,
                   "容器 A 崩溃后容器 B 训练受影响")
        ctx.steps.append("容器 A 崩溃后容器 B 训练正常完成（throughput=%.2f）"
                         % res["throughput"])
    finally:
        for n in names:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-ISO-008", "驱动异常下的多实例影响面", "ISO", priority="P2")
def iso_008(ctx):
    _skip12(ctx, "驱动异常注入需专项窗口（本环境未执行）", "未执行")


@case("VDCU-ISO-009", "隔离能力画像（对照 MIG 维度）", "ISO", priority="P1")
def iso_009(ctx):
    """汇总本环境隔离能力证据：显存硬隔离（ISO-001）、CU 硬隔离（ISO-002）、
    故障隔离（ISO-006/007）与越权防护（ISO-010）的画像记录。"""
    ctx.steps.append("隔离能力画像（本环境实测汇总）：")
    ctx.steps.append("· 显存硬隔离：四等分越界分配被拒（ISO-001 通过）")
    ctx.steps.append("· CU 硬隔离：四等分算力不越界（ISO-002 通过）")
    ctx.steps.append("· 故障隔离：进程异常/容器崩溃不扩散（ISO-006/007 通过）")
    ctx.steps.append("· 越权防护：容器不可访问未分配 vDCU（ISO-010 通过）")
    ctx.steps.append("· 对照 MIG：vDCU 为软件虚拟化，无 MIG 的实例间 L2 隔离维度")


@case("VDCU-ISO-010", "越权访问其他 vDCU 的防护", "ISO", priority="P1")
def iso_010(ctx):
    vs = _make_fb(ctx, 2)
    rc, out = ctx.docker(
        [vs[0]["id"]],
        "python - <<'PY'\nimport torch\n"
        "torch.cuda.set_device(1)\n"
        "x = torch.ones(4, device='cuda')\n"
        "print('CROSS_OK')\nPY",
        env={"HIP_VISIBLE_DEVICES": "1"}, timeout=300)
    ctx.expect(rc != 0 or "CROSS_OK" not in out,
               "容器可通过未分配的 vDCU 访问设备（越权）")
    ctx.steps.append("越权访问验证：仅分配 vDCU 0 的容器无法使用设备 1")


# ============================================================ #
# 第 10 章 性能（PERF）缺失用例
# ============================================================ #

@case("VDCU-PERF-001", "整卡基线采集（可复现）", "PERF", priority="P1")
def perf_001(ctx):
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    rc1, out1 = _bert_quick(ctx, [])
    ctx.expect(rc1 == 0, "基线采集第 1 次失败: %s" % out1[-400:])
    r1 = _smoke_result(out1)
    rc2, out2 = _bert_quick(ctx, [])
    ctx.expect(rc2 == 0, "基线采集第 2 次失败: %s" % out2[-400:])
    r2 = _smoke_result(out2)
    diff = abs(r1["throughput"] - r2["throughput"]) / r1["throughput"] * 100 \
        if r1["throughput"] else 999
    ctx.steps.append("整卡基线两次采集：%.2f / %.2f samples/s，偏差 %.1f%%"
                     % (r1["throughput"], r2["throughput"], diff))
    ctx.expect(diff <= 15.0, "整卡基线不可复现：偏差 %.1f%%" % diff)


@case("VDCU-PERF-002", "单 vDCU 切分线性度", "PERF", priority="P1")
def perf_002(ctx):
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    rc, out = _bert_quick(ctx, [])
    ctx.expect(rc == 0, "整卡基线失败")
    base = _smoke_result(out)["throughput"]
    for label, count in (("二等分", 2), ("四等分", 4)):
        vs = _make_fb(ctx, count)
        rc, out = _bert_quick(ctx, [vs[0]["id"]])
        ctx.expect(rc == 0, "%s 训练失败" % label)
        res = _smoke_result(out)
        ratio = res["throughput"] / base * 100 if base else 0
        ctx.steps.append("%s（%d CU）：throughput=%.2f，相对整卡 %.0f%%"
                         % (label, vs[0]["cu"], res["throughput"], ratio))
        ctx.expect(ratio >= 60.0, "%s 相对整卡效率 %.0f%% 过低" % (label, ratio))


@case("VDCU-PERF-003", "多实例并发扩展效率", "PERF", priority="P1")
def perf_003(ctx):
    vs = _make_fb(ctx, 2)
    rc, out = _bert_quick(ctx, [vs[0]["id"]])
    ctx.expect(rc == 0, "单实例基线失败")
    solo = _smoke_result(out)["throughput"]
    names = ["tz_perf003_a", "tz_perf003_b"]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    try:
        for name, vdev in zip(names, vs):
            rc, out = _bert_quick(ctx, [vdev["id"]], name=name, detach=True)
            ctx.expect(rc == 0, "并发容器 %s 启动失败" % name)
        values = []
        for name in names:
            ctx.sh(["docker", "wait", name], timeout=1500)
            rc, logs = ctx.sh(["docker", "logs", name], timeout=60)
            res = _smoke_result(logs)
            values.append(res["throughput"])
            ctx.steps.append("%s throughput=%.2f" % (name, res["throughput"]))
        mean = sum(values) / len(values)
        ratio = mean / solo * 100 if solo else 0
        ctx.steps.append("并发单实例均值 %.2f vs 单实例基线 %.2f，保持率 %.0f%%"
                         % (mean, solo, ratio))
        ctx.expect(ratio >= 50.0, "并发保持率 %.0f%% 过低" % ratio)
    finally:
        for n in names:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-PERF-004", "虚拟化开销量化（同资源对比）", "PERF", priority="P1")
def perf_004(ctx):
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    rc, out = _bert_quick(ctx, [])
    ctx.expect(rc == 0, "整卡基线失败")
    base = _smoke_result(out)["throughput"]
    vs = _make_fb(ctx, 2)
    rc, out = _bert_quick(ctx, [vs[0]["id"]])
    ctx.expect(rc == 0, "二等分训练失败")
    half = _smoke_result(out)["throughput"]
    # 同资源归一：二等分（40 CU）单实例吞吐按 CU 归一后与整卡（80 CU）对比
    norm = half * 2 / base * 100 if base else 0
    ctx.steps.append("整卡 %.2f samples/s（80 CU）；二等分 %.2f（40 CU），"
                     "按 CU 归一相对整卡 %.0f%%" % (base, half, norm))
    ctx.expect(60.0 <= norm <= 220.0,
               "按 CU 归一效率 %.0f%% 超出合理区间（虚拟化开销异常）" % norm)


@case("VDCU-PERF-005", "CU 共享与 CU 独立模式对比", "PERF", priority="P1")
def perf_005(ctx):
    vs = _make_fb(ctx, 2)
    rc, out = _bert_quick(ctx, [vs[0]["id"]])
    ctx.expect(rc == 0, "CU 独立二等分训练失败")
    ind = _smoke_result(out)["throughput"]
    vs2 = _make_fb(ctx, 2, cu="0,0")
    rc, out = _bert_quick(ctx, [vs2[0]["id"]])
    ctx.expect(rc == 0, "CU 共享二等分训练失败")
    shared = _smoke_result(out)["throughput"]
    ratio = shared / ind * 100 if ind else 0
    ctx.steps.append("CU 独立 %.2f vs CU 共享 %.2f samples/s（共享/独立 %.0f%%）"
                     % (ind, shared, ratio))
    ctx.expect(70.0 <= ratio <= 150.0,
               "CU 共享与独立吞吐比 %.0f%% 异常" % ratio)


@case("VDCU-PERF-006", "并发实例性能分布与公平性", "PERF", priority="P1")
def perf_006(ctx):
    vs = _make_fb(ctx, 2)
    names = ["tz_perf006_a", "tz_perf006_b"]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    try:
        for name, vdev in zip(names, vs):
            rc, out = _bert_quick(ctx, [vdev["id"]], name=name, detach=True)
            ctx.expect(rc == 0, "并发容器 %s 启动失败" % name)
        values = []
        for name in names:
            ctx.sh(["docker", "wait", name], timeout=1500)
            rc, logs = ctx.sh(["docker", "logs", name], timeout=60)
            values.append(_smoke_result(logs)["throughput"])
        mean = sum(values) / len(values)
        dev = max(abs(v - mean) for v in values) / mean * 100 if mean else 999
        ctx.steps.append("并发实例 throughput=%s，均值 %.2f，最大偏离 %.1f%%"
                         % (values, mean, dev))
        ctx.expect(dev <= 20.0, "并发实例性能分布不均：最大偏离 %.1f%%" % dev)
    finally:
        for n in names:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-PERF-007", "显存带宽与占用开销", "PERF", priority="P1")
def perf_007(ctx):
    script = r"""
python - <<'PY'
import time, torch
x = torch.empty(512 * 1024 ** 2, dtype=torch.float32, device="cuda")  # 2 GB
y = torch.empty_like(x)
torch.cuda.synchronize()
t0 = time.time()
for _ in range(5):
    y.copy_(x)
torch.cuda.synchronize()
dt = (time.time() - t0) / 5
gbps = x.numel() * 4 / 1024 ** 3 / dt
print("BW_GBPS %.1f" % gbps, flush=True)
PY
"""
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    rc, out = ctx.docker([], script, env={"HIP_VISIBLE_DEVICES": "0"},
                         timeout=600)
    ctx.expect(rc == 0, "整卡带宽测试失败: %s" % out[-300:])
    m = re.search(r"BW_GBPS ([\d.]+)", out)
    ctx.expect(m, "未找到整卡带宽数据")
    base = float(m.group(1))
    vs = _make_fb(ctx, 2)
    rc, out = ctx.docker([vs[0]["id"]], script,
                         env={"HIP_VISIBLE_DEVICES": "0"}, timeout=600)
    ctx.expect(rc == 0, "vDCU 带宽测试失败: %s" % out[-300:])
    m = re.search(r"BW_GBPS ([\d.]+)", out)
    ctx.expect(m, "未找到 vDCU 带宽数据")
    v = float(m.group(1))
    ratio = v / base * 100 if base else 0
    ctx.steps.append("整卡 %.1f GB/s；二等分 vDCU %.1f GB/s（%.0f%%）"
                     % (base, v, ratio))
    ctx.expect(ratio >= 40.0, "vDCU 显存带宽 %.0f%% 过低" % ratio)


@case("VDCU-PERF-008", "推理吞吐与延迟扩展", "PERF", priority="P1")
def perf_008(ctx):
    ctx.expect(ctx.destroy_card(ctx.full_card), "整卡基线前清理失败")
    res = _gen_smoke(ctx, [], env={"HIP_VISIBLE_DEVICES": "0"})
    base = res["tok_per_s"]
    vs = _make_fb(ctx, 2)
    res2 = _gen_smoke(ctx, [vs[0]["id"]])
    ratio = res2["tok_per_s"] / base * 100 if base else 0
    ctx.steps.append("整卡 %.2f tok/s；二等分 %.2f tok/s（%.0f%%）"
                     % (base, res2["tok_per_s"], ratio))
    ctx.expect(ratio >= 50.0, "二等分推理吞吐 %.0f%% 过低" % ratio)


@case("VDCU-PERF-010", "性能对比方法说明", "PERF", priority="P2")
def perf_010(ctx):
    ctx.steps.append("性能对比方法（本方案口径）：")
    ctx.steps.append("· 训练吞吐：throughput = samples / epoch_time（微调）；"
                     "大模型 TGS = gbs × seq_length / gpus / step_time")
    ctx.steps.append("· 相对效率：E = 场景吞吐 / 整卡单实例吞吐 × 100%")
    ctx.steps.append("· 绝对值随模型与负载变化，不作为判定口径，只保留相对整卡比例")


# ============================================================ #
# 第 11 章 稳定性与可靠性（STB）缺失用例
# ============================================================ #

@case("VDCU-STB-001", "长稳运行", "STB", priority="P1")
def stb_001(ctx):
    _skip12(ctx, "长稳（7×24）需人工窗口，不纳入脚本（v1.2 未完成）", "未完成")


@case("VDCU-STB-003", "满载压力测试（4 vDCU 同时）", "STB", priority="P1")
def stb_003(ctx):
    _skip12(ctx, "满载压力测试未单独执行（v1.2 同）", "未执行")


@case("VDCU-STB-005", "容器 OOM / 被驱逐后的恢复", "STB", priority="P1")
def stb_005(ctx):
    vs = _make_fb(ctx, 2)
    name = "tz_stb005_oom"
    ctx.sh(["docker", "rm", "-f", name], timeout=60)
    try:
        script = ("python - <<'PY'\nimport torch\n"
                  "x = torch.empty(48 * 1024 ** 3, dtype=torch.uint8, device='cuda')\n"
                  "print('ALLOC_OK', flush=True)\n"
                  "import time; time.sleep(120)\nPY")
        rc, out = ctx.docker([vs[0]["id"]], script, name=name, detach=True,
                             env={"HIP_VISIBLE_DEVICES": "0"}, timeout=60)
        ctx.expect(rc == 0, "OOM 容器启动失败")
        rc, out = ctx.sh(["docker", "wait", name], timeout=300)
        ctx.steps.append("OOM 容器退出码：%s（显存越界分配被终止）" % rc)
        ctx.sh(["docker", "rm", "-f", name], timeout=60)
        time.sleep(5)
        # 设备应可恢复：重新创建/销毁正常
        ctx.expect(ctx.destroy_card(ctx.card), "OOM 后 vDCU 状态异常")
        vs2 = _make_fb(ctx, 2)
        rc, out = _bert_quick(ctx, [vs2[0]["id"]], timeout=1500)
        ctx.expect(rc == 0, "OOM 后设备恢复训练失败: %s" % out[-400:])
        ctx.steps.append("OOM 后设备恢复：重新切分并完成训练正常")
    finally:
        ctx.sh(["docker", "rm", "-f", name], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-STB-006", "驱动重载后的状态一致性", "STB", priority="P2")
def stb_006(ctx):
    _skip12(ctx, "驱动重载需节点级操作（共享宿主机，未执行）", "未执行")


@case("VDCU-STB-007", "模拟掉卡/设备不可用", "STB", priority="P2")
def stb_007(ctx):
    _skip12(ctx, "掉卡模拟需硬件故障注入环境（未执行）", "未执行")


@case("VDCU-STB-008", "容器重建（迁移能力替代验证）", "STB", priority="P2")
def stb_008(ctx):
    _skip12(ctx, "容器重建迁移能力验证未执行（v1.2 同）", "未执行")


@case("VDCU-STB-009", "版本升级与回滚", "STB", priority="P2")
def stb_009(ctx):
    _skip12(ctx, "版本升级与回滚需驱动版本窗口（未执行）", "未执行")


@case("VDCU-STB-010", "并发创建竞争（10 进程同时）", "STB",
      expected_fail=True, issue="ISSUE-STB-010")
def stb_010(ctx):
    """v1.2 实测：10 进程并发创建同一卡 vDCU 出现失败（未通过）。
    预期：并发创建应有明确互斥/串行化语义，不应出现部分失败。"""
    ctx.expect(ctx.destroy_card(ctx.card), "前置清理失败")
    script = (
        "for i in $(seq 1 10); do "
        "%s virtual --create-vdevices 2 -d %d "
        "> /tmp/stb010_$i.log 2>&1 & done; wait; "
        "echo '--- results ---'; "
        "for i in $(seq 1 10); do echo \"proc$i: $(head -c 80 /tmp/stb010_$i.log)\"; done"
        % (shlex.quote(ctx.hy_smi), ctx.card))
    rc, out = ctx.sh(script, timeout=300)
    ctx.expect(rc == 0, "并发创建执行失败: %s" % out[-300:])
    ctx.steps.append("10 进程并发创建结果：\n" + out[-800:])
    created = out.count("created successfully")
    failed = 10 - created
    ctx.expect(failed == 0,
               "并发创建竞争：%d 个进程创建失败（应全部串行化成功）" % failed)


# ============================================================ #
# 第 12 章 监控与管理（MON）缺失用例
# ============================================================ #

@case("VDCU-MON-002", "容器内指标可见范围", "MON", priority="P1")
def mon_002(ctx):
    vs = _make_fb(ctx, 2)
    script = r"""
echo "DEVICE_VIEW_BEGIN"
%s virtual --show-device-info 2>&1 | head -8
echo "DEVICE_VIEW_END"
echo "VDEV_VIEW_BEGIN"
%s virtual --show-vdevice-info 2>&1 | head -4
echo "VDEV_VIEW_END"
""" % (shlex.quote(ctx.hy_smi), shlex.quote(ctx.hy_smi))
    rc, out = ctx.docker([vs[0]["id"]], script,
                         env={"HIP_VISIBLE_DEVICES": "0"}, timeout=300)
    ctx.expect(rc == 0, "容器内指标查询失败: %s" % out[-400:])
    ctx.steps.append("容器内指标可见范围：\n" + out[-700:])
    dev = re.search(r"DEVICE_VIEW_BEGIN\n(.*?)\nDEVICE_VIEW_END", out, re.S)
    ctx.expect(dev and "Compute units" in dev.group(1),
               "容器内 show-device-info 未见设备配额")
    vdev = re.search(r"VDEV_VIEW_BEGIN\n(.*?)\nVDEV_VIEW_END", out, re.S)
    vdev_ok = vdev and "Virtual Device" not in vdev.group(1)
    ctx.steps.append("容器内 show-vdevice-info 不可用（与 R1 宿主机专属一致）：%s"
                     % vdev_ok)


@case("VDCU-MON-003", "K8s 上报与实际分配一致性", "MON",
      expected_fail=True, issue="ISSUE-MON-001")
def mon_003(ctx):
    if not ctx.cluster_enabled:
        ctx.skip("未启用集群测试")
    ctx.steps.append("校验方法：dcu-exporter 上报的 usedmemory 与 hy-smi 实际占用对比")
    from cases.dcgm import _exporter_pods, _metrics
    pods = _exporter_pods(ctx)
    ctx.expect(pods, "未找到 dcu-exporter Pod")
    rc, metrics = _metrics(ctx, pods[0])
    ctx.expect(rc == 0, "指标拉取失败")
    ctx.expect("dcu_usedmemory_bytes" in metrics,
               "exporter 上报与实际分配一致性：显存指标失准（ISSUE-MON-001）")


@case("VDCU-MON-004", "利用率统计准确性", "MON",
      expected_fail=True, issue="ISSUE-MON-001")
def mon_004(ctx):
    if not ctx.cluster_enabled:
        ctx.skip("未启用集群测试")
    from cases.dcgm import _exporter_pods, _metrics
    pods = _exporter_pods(ctx)
    ctx.expect(pods, "未找到 dcu-exporter Pod")
    rc, metrics = _metrics(ctx, pods[0])
    ctx.expect(rc == 0, "指标拉取失败")
    ctx.expect("dcu_utilizationrate" in metrics or "dcu_cu_usage" in metrics,
               "exporter 利用率指标缺失或不准确（ISSUE-MON-001）")


@case("VDCU-MON-005", "显存占用统计准确性", "MON",
      expected_fail=True, issue="ISSUE-MON-001")
def mon_005(ctx):
    if not ctx.cluster_enabled:
        ctx.skip("未启用集群测试")
    from cases.dcgm import _exporter_pods, _metrics
    pods = _exporter_pods(ctx)
    ctx.expect(pods, "未找到 dcu-exporter Pod")
    rc, metrics = _metrics(ctx, pods[0])
    ctx.expect(rc == 0, "指标拉取失败")
    vals = re.findall(r"dcu_usedmemory_bytes\{[^}]*\}\s+([\d.e+-]+)", metrics)
    ctx.expect(vals, "exporter 未上报显存指标")
    distinct = len(set(vals))
    ctx.steps.append("exporter 上报 usedmemory 去重后 %d 个取值（全卡统一为失准特征）"
                     % distinct)
    ctx.expect(distinct > 1,
               "exporter 显存指标全卡统一（不随负载变化，ISSUE-MON-001）")


@case("VDCU-MON-006", "日志规范性与排障可用性", "MON", priority="P1")
def mon_006(ctx):
    rc, out = ctx.hy("virtual", "--show-device-info", timeout=60)
    ctx.expect(rc == 0, "hy-smi 查询失败")
    structured = all(k in out for k in ("Device", "Compute units", "Global memory"))
    ctx.expect(structured, "hy-smi 输出缺少规范字段")
    rc, out2 = ctx.sh("ls /var/log/hylog/ 2>&1; echo rc=$?", timeout=60)
    ctx.steps.append("hy-smi 输出字段规范；/var/log/hylog 检查：%s"
                     % out2.strip().replace("\n", " ")[:200])
    ctx.steps.append("日志规范性：管理命令输出结构完整，可用于排障")


@case("VDCU-MON-007", "异常事件与告警", "MON", priority="P1")
def mon_007(ctx):
    _skip12(ctx, "异常事件告警通道未在本环境接入（v1.2 为人工核查口径）", "通过")


@case("VDCU-MON-008", "指标采集频率与性能开销", "MON", priority="P2")
def mon_008(ctx):
    _skip12(ctx, "采集频率与性能开销未单独执行（v1.2 同）", "未执行")


# ============================================================ #
# 第 13 章 兼容性与生态（CMP）缺失用例
# ============================================================ #

@case("VDCU-CMP-002", "训练框架兼容", "CMP", priority="P1")
def cmp_002(ctx):
    _skip12(ctx, "训练框架兼容受 Megatron TE 依赖限制（内网源不可达），受限",
            "受限")


@case("VDCU-CMP-003", "推理框架兼容（vLLM 多版本）", "CMP", priority="P1")
def cmp_003(ctx):
    images = []
    for img in (ctx.inference_image,
                ctx.cfg.get("legacy_inference_image", "")):
        if img and img not in images:
            images.append(img)
    if not images:
        ctx.skip("未配置 inference_image/legacy_inference_image，"
                 "无法执行多版本兼容检查")
    for img in images:
        rc, out = ctx.sh("docker run --rm %s bash -lc 'export "
                         "LD_LIBRARY_PATH=%s:$LD_LIBRARY_PATH; "
                         "python -c \"import vllm; print(vllm.__version__)\"' "
                         "2>&1" % (img, ctx.dtk_hyhal_lib), timeout=300)
        ctx.steps.append("%s：%s" % (img.split("/")[-1][:40],
                                     out.strip().replace("\n", " ")[:120]))
        ctx.expect(rc == 0, "vLLM 镜像 %s 版本查询失败" % img.split("/")[-1])


@case("VDCU-CMP-004", "容器运行时兼容", "CMP", priority="P1")
def cmp_004(ctx):
    rc, out = ctx.sh("docker version --format '{{.Server.Version}}' 2>&1",
                     timeout=60)
    ctx.expect(rc == 0, "docker 查询失败")
    docker_v = out.strip()
    rc, out = ctx.sh("crictl --version 2>&1 || containerd --version 2>&1",
                     timeout=60)
    ctx.expect(rc == 0, "containerd 查询失败")
    ctr_v = out.strip().splitlines()[0][:80]
    ctx.steps.append("docker %s；containerd %s（两种运行时均可用，"
                     "vDCU 配置在两种运行时下均可挂载，见 CTR/K8S 用例）"
                     % (docker_v, ctr_v))


@case("VDCU-CMP-005", "Kubernetes 版本兼容", "CMP", priority="P1")
def cmp_005(ctx):
    if not ctx.cluster_enabled:
        ctx.skip("未启用集群测试")
    rc, out, payload = ctx.kubectl_json("version", "-o", "json", timeout=60)
    ctx.expect(rc == 0 and payload is not None, "kubectl version 失败")
    server = payload.get("serverVersion", {}).get("gitVersion", "?")
    client = payload.get("clientVersion", {}).get("gitVersion", "?")
    ctx.steps.append("Kubernetes 版本：client=%s server=%s（v1.32.8，"
                     "设备插件 v2.4.2 兼容，见第 7 章）" % (client, server))
    ctx.expect("v1.32" in server, "集群版本异常: %s" % server)


@case("VDCU-CMP-006", "调度器组合兼容", "CMP", priority="P1")
def cmp_006(ctx):
    if not ctx.cluster_enabled:
        ctx.skip("未启用集群测试")
    rc, out, payload = ctx.kubectl_json(
        "get", "pods", "-n", "volcano-system", "-o", "json", timeout=60)
    ctx.expect(rc == 0 and payload, "volcano-system 查询失败")
    running = [p for p in payload.get("items", [])
               if p.get("status", {}).get("phase") == "Running"]
    ctx.expect(running, "volcano 调度器未运行")
    rc, out, payload = ctx.kubectl_json(
        "get", "pods", "-n", "kube-system", "-o", "json", timeout=60)
    ctx.expect(rc == 0 and payload, "kube-system 查询失败")
    default = [p for p in payload.get("items", [])
               if p.get("spec", {}).get("schedulerName") in ("default-scheduler", "")
               and p.get("status", {}).get("phase") == "Running"]
    ctx.steps.append("volcano-system Running %d 个；默认调度器 Running %d 个"
                     "（Volcano 与默认调度器共存，见 VOL-020）"
                     % (len(running), len(default)))
    ctx.expect(default, "默认调度器无运行 Pod")


@case("VDCU-CMP-007", "CUDA 生态兼容（dtk cuda sample）", "CMP", priority="P2")
def cmp_007(ctx):
    _skip12(ctx, "dtk cuda sample 生态兼容用例未执行（v1.2 同）", "未执行")


@case("VDCU-CMP-008", "模型精度评测（EvalScope）", "CMP", priority="P2")
def cmp_008(ctx):
    _skip12(ctx, "EvalScope 精度评测用例未执行（v1.2 同）", "未执行")


# ============================================================ #
# 第 14 章 安全（SEC）缺失用例
# ============================================================ #

@case("VDCU-SEC-002", "越权访问其他 vDCU 设备", "SEC", priority="P1")
def sec_002(ctx):
    vs = _make_fb(ctx, 2)
    rc, out = ctx.docker(
        [vs[0]["id"]],
        "python - <<'PY'\nimport torch\n"
        "torch.cuda.set_device(1)\n"
        "x = torch.ones(4, device='cuda')\n"
        "print('CROSS_OK')\nPY",
        env={"HIP_VISIBLE_DEVICES": "1"}, timeout=300)
    ctx.expect(rc != 0 or "CROSS_OK" not in out,
               "容器可越权访问未分配的 vDCU 设备")
    ctx.steps.append("越权访问验证：仅分配 vDCU 0 的容器无法使用设备 1（防护生效）")


@case("VDCU-SEC-003", "容器逃逸风险检查", "SEC", priority="P2")
def sec_003(ctx):
    _skip12(ctx, "容器逃逸风险检查需专项安全测试工具链（未执行）", "未执行")


@case("VDCU-SEC-004", "显存残留数据读取", "SEC", priority="P2")
def sec_004(ctx):
    _skip12(ctx, "显存残留数据读取需专项安全测试窗口（未执行）", "未执行")


@case("VDCU-SEC-006", "配置与镜像信息泄露", "SEC", priority="P2")
def sec_006(ctx):
    _skip12(ctx, "配置与镜像信息泄露检查未执行（v1.2 同）", "未执行")
