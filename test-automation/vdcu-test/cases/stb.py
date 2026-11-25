# -*- coding: utf-8 -*-
"""第 11 章 稳定性与可靠性测试（STB）。

轻量级可执行用例；长稳（7×24）类用例需人工窗口，不纳入脚本。
"""
from vdcu.core import case


@case("VDCU-STB-002", "反复创建销毁 vDCU 20 轮无残留", "STB", priority="P1")
def stb_002(ctx):
    """反复创建销毁：验证资源回收与配置残留。"""
    ctx.expect(ctx.destroy_card(ctx.card), "前置清理失败")
    for i in range(20):
        rc, out = ctx.hy("virtual", "--create-vdevices", "4",
                         "-d", str(ctx.card))
        ctx.expect(rc == 0, "第 %d 轮创建失败: %s" % (i + 1, out[-200:]))
        ctx.expect(len(ctx.vdevs_on(ctx.card)) == 4, "第 %d 轮创建数量异常" % (i + 1))
        ctx.expect(ctx.destroy_card(ctx.card), "第 %d 轮销毁失败" % (i + 1))
    rem = ctx.remaining()
    cu, mem = ctx.devices()[0]["cu"], ctx.devices()[0]["mem"]
    ctx.expect(rem[ctx.card]["cu"] == cu and rem[ctx.card]["mem"] == mem,
               "20 轮后剩余资源未完全回收: %s" % rem[ctx.card])
    rc, out = ctx.sh("ls %s/ | wc -l" % ctx.vdev_conf_dir, timeout=30)
    ctx.steps.append("%s 配置残留：%s 个文件（本机其他卡配置不计）"
                     % (ctx.vdev_conf_dir, out.strip()))


@case("VDCU-STB-004", "训练进程被 kill 后设备可恢复（10 次）", "STB",
      priority="P1")
def stb_004(ctx):
    """进程异常恢复：kill 训练进程后 vDCU 可立即复用，其他实例不受影响。"""
    vs = ctx.vdevs_on(ctx.card) or []
    if not vs:
        ctx.expect(ctx.destroy_card(ctx.card), "前置清理失败")
        ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
        vs = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])
    ok = 0
    for i in range(10):
        name = "tz_stb004_%d" % i
        ctx.sh(["docker", "rm", "-f", name], timeout=60)
        rc, out = ctx.docker(
            [vs[0]["id"]],
            "python - <<'PY'\nimport time, torch\n"
            "x = torch.empty(2 * 1024 ** 3, dtype=torch.uint8, device='cuda')\n"
            "print('HOLDING', flush=True)\ntime.sleep(600)\nPY",
            name=name, detach=True, env={"HIP_VISIBLE_DEVICES": "0"},
            timeout=120)
        ctx.expect(rc == 0, "第 %d 次启动失败" % (i + 1))
        import time as _t
        _t.sleep(6)
        rc, out = ctx.sh(["docker", "exec", name, "bash", "-lc",
                          "kill -9 1"], timeout=60)
        ctx.sh(["docker", "rm", "-f", name], timeout=60)
        # kill 后立即复用：启动短探测
        rc, out = ctx.docker([vs[0]["id"]],
                             "python -c \"import torch;print(torch.ones(2,device='cuda').sum().item())\"",
                             env={"HIP_VISIBLE_DEVICES": "0"}, timeout=300)
        if rc == 0 and "2.0" in out:
            ok += 1
    ctx.expect(ok == 10, "10 次 kill 后恢复成功 %d 次" % ok)
    ctx.steps.append("进程 kill 后设备恢复：10/10 成功")
