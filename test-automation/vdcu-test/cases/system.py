# -*- coding: utf-8 -*-
"""监控、兼容性与性能开销类用例。"""
import re
import time

from vdcu.core import case, parse_device_info, parse_vdev_info


@case("VDCU-MON-001", "hy-smi 字段完整性与语义", "MON")
def mon_001(ctx):
    devs = ctx.devices()
    ctx.expect(len(devs) == ctx.physical_cards,
               "物理设备数应为 %d" % ctx.physical_cards)
    cu = devs[0]["cu"]
    ctx.expect(all(d["phys"] is not None and d["cu"] is not None and
                   d["mem"] is not None for d in devs),
               "物理设备字段不完整")

    ctx.expect(ctx.destroy_card(ctx.card), "前置清理失败")
    rc, out = ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    ctx.expect(rc == 0, "创建失败: %s" % out[-300:])

    rc, out = ctx.hy("virtual", "--show-vdevice-info")
    vs = [v for v in parse_vdev_info(out) if v["phys"] == ctx.card]
    ctx.expect(len(vs) == 2, "本卡应查询到 2 个 vDCU，实际 %d" % len(vs))
    for v in vs:
        ctx.expect(v["phys"] == ctx.card, "缺少或错误的 Physical Device 字段")
        ctx.expect(v["cu"] == cu // 2, "缺少或错误的 Compute units 字段")
        ctx.expect(v["mem"] is not None, "缺少 Global memory 字段")


@case("VDCU-CMP-001", "基础镜像兼容（DTK/PyTorch/HIP 可用）", "CMP")
def cmp_001(ctx):
    ctx.expect(ctx.destroy_card(ctx.full_card), "兼容性冒烟前整卡清理失败")
    script = r"""
python - <<'PY'
import torch
print("TORCH", torch.__version__)
print("HIP_VERSION", getattr(torch.version, "hip", ""))
print("HIP_AVAILABLE", torch.cuda.is_available())
x = torch.ones(8, 8, device="cuda")
print("HIP_OP_OK", float(x.sum()) == 64.0)
PY
"""
    rc, out = ctx.docker([], script,
                         env={"HIP_VISIBLE_DEVICES": "0"},
                         timeout=300)
    ctx.expect(rc == 0, "镜像内执行失败: %s" % out[-400:])
    ctx.expect("TORCH " in out, "未报告 PyTorch 版本: %s" % out[-300:])
    # 版本以被测宿主机实际镜像为准。当前环境使用 bert:v1，可能是
    # PyTorch 2.13/DTK 26.10，也允许现场用兼容的 2.x 镜像覆盖。
    m = re.search(r"TORCH\s+([0-9]+\.[0-9]+)", out)
    ctx.expect(m is not None and int(m.group(1).split(".")[0]) >= 2,
               "PyTorch 主版本不满足要求: %s" % out[-300:])
    ctx.expect("HIP_VERSION " in out, "未报告 HIP/DTK 版本")
    ctx.expect("HIP_AVAILABLE True" in out, "HIP 不可用")
    ctx.expect("HIP_OP_OK True" in out, "HIP 基本运算失败")


@case("VDCU-PERF-009", "vDCU 创建/销毁操作开销", "PERF", priority="P2")
def perf_009(ctx):
    ctx.destroy_card(ctx.card)
    rounds = 20
    costs = []
    for _ in range(rounds):
        t0 = time.time()
        rc1, o1 = ctx.hy("virtual", "--create-vdevices", "4", "-d", str(ctx.card))
        rc2, o2 = ctx.hy("virtual", "--destroy-vdevices", "-d", str(ctx.card))
        dt = time.time() - t0
        ctx.expect(rc1 == 0 and rc2 == 0, "第 %d 轮创建/销毁失败" % (len(costs) + 1))
        ctx.expect(ctx.destroy_card(ctx.card, retries=3, delay=0.1),
                   "第 %d 轮结束后 vDCU 未及时释放" % (len(costs) + 1))
        costs.append(dt)

    ctx.expect(len(ctx.vdevs_on(ctx.card)) == 0, "循环结束后仍有残留 vDCU")
    avg = sum(costs) / len(costs)
    print("        %d 轮：平均 %.2fs，最大 %.2fs" % (rounds, avg, max(costs)))
    ctx.expect(max(costs) < 10, "单轮创建+销毁耗时 %.2fs 超过 10s" % max(costs))
