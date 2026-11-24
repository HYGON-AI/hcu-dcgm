# -*- coding: utf-8 -*-
"""第 2 章虚拟化功能与资源约束测试（FUN）。"""
import time

from vdcu.core import case, Context

TOL = 16 * 1024 * 1024  # 显存字节容差（驱动保留开销）

HOLD_DEVICE = r"""
python - <<'PY'
import time
import torch

torch.cuda.set_device(0)
a = torch.randn(1024, 1024, device="cuda")
b = torch.randn(1024, 1024, device="cuda")
while True:
    _ = a @ b
    torch.cuda.synchronize()
    time.sleep(0.2)
PY
"""


def baseline(ctx):
    """读取物理卡规格：返回 (cu, mem)。"""
    devs = ctx.devices()
    ctx.expect(len(devs) > 0, "未查询到任何物理设备")
    selected = next((d for d in devs
                     if d.get("id") == ctx.card or d.get("phys") == ctx.card), devs[0])
    return selected["cu"], selected["mem"]


def near(a, b, tol=TOL):
    return a is not None and b is not None and abs(a - b) <= tol


@case("VDCU-FUN-001", "物理设备与剩余资源查询基线", "FUN")
def fun_001(ctx):
    ctx.destroy_card(ctx.card)
    cu, mem = baseline(ctx)
    devs = ctx.devices()

    ctx.expect(len(devs) == ctx.physical_cards,
               "物理设备数应为 %d，实际 %d" % (ctx.physical_cards, len(devs)))
    ctx.expect(all(d["cu"] == cu for d in devs), "各卡 CU 数不一致")
    ctx.expect(all(d["mem"] == mem for d in devs), "各卡显存不一致")

    rem = ctx.remaining()
    ctx.expect(rem[ctx.card]["cu"] == cu, "剩余 CU 应等于物理 CU %d，实际 %s"
               % (cu, rem[ctx.card]["cu"]))
    ctx.expect(rem[ctx.card]["mem"] == mem, "剩余显存应等于物理显存 %d" % mem)


@case("VDCU-FUN-002", "创建 2 个 CU 与显存等分 vDCU", "FUN")
def fun_002(ctx):
    ctx.destroy_card(ctx.card)
    cu, mem = baseline(ctx)
    rc, out = ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    ctx.expect(rc == 0 and "successfully" in out.lower(), "创建失败: %s" % out[-300:])

    vs = ctx.vdevs_on(ctx.card)
    ctx.expect(len(vs) == 2, "应创建 2 个 vDCU，实际 %d" % len(vs))
    for v in vs:
        ctx.expect(v["cu"] == cu // 2, "CU 应为 %d，实际 %s" % (cu // 2, v["cu"]))
        ctx.expect(near(v["mem"], mem // 2), "显存应为 %d，实际 %s" % (mem // 2, v["mem"]))


@case("VDCU-FUN-003", "创建 4 个等分 vDCU（实例上限）", "FUN")
def fun_003(ctx):
    ctx.destroy_card(ctx.card)
    cu, mem = baseline(ctx)
    rc, out = ctx.hy("virtual", "--create-vdevices", "4", "-d", str(ctx.card))
    ctx.expect(rc == 0, "创建失败: %s" % out[-300:])

    vs = ctx.vdevs_on(ctx.card)
    ctx.expect(len(vs) == 4, "应创建 4 个 vDCU，实际 %d" % len(vs))
    for v in vs:
        ctx.expect(v["cu"] == cu // 4, "CU 应为 %d，实际 %s" % (cu // 4, v["cu"]))
        ctx.expect(near(v["mem"], mem // 4), "显存应为 %d，实际 %s" % (mem // 4, v["mem"]))

    rem = ctx.remaining()
    ctx.expect(rem[ctx.card]["cu"] == 0, "4 等分后剩余 CU 应为 0")


@case("VDCU-FUN-004", "创建数量超限（5 个）应被拒绝", "FUN")
def fun_004(ctx):
    ctx.destroy_card(ctx.card)
    rc, out = ctx.hy("virtual", "--create-vdevices", "5", "-d", str(ctx.card))
    ctx.expect(rc != 0, "创建 5 个应返回非 0")
    ctx.expect("1 to 4" in out or "failed" in out.lower(), "缺少明确错误提示: %s" % out[-300:])
    ctx.expect(len(ctx.vdevs_on(ctx.card)) == 0, "失败后不应产生任何 vDCU")


@case("VDCU-FUN-005", "自定义 CU 与显存组合创建", "FUN")
def fun_005(ctx):
    ctx.destroy_card(ctx.card)
    rc, out = ctx.hy("virtual", "--create-vdevices", "4",
                     "--vdevice-compute-units", "10,10,20,40",
                     "--vdevice-memory-size", "4096,8192,16384,32768",
                     "-d", str(ctx.card))
    ctx.expect(rc == 0, "创建失败: %s" % out[-300:])

    vs = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])
    ctx.expect(len(vs) == 4, "应创建 4 个 vDCU")
    ctx.expect([v["cu"] for v in vs] == [10, 10, 20, 40],
               "CU 组合不符: %s" % [v["cu"] for v in vs])
    for v, mib in zip(vs, [4096, 8192, 16384, 32768]):
        ctx.expect(near(v["mem"], mib * 1024 * 1024), "显存不符: %s != %d MiB" % (v["mem"], mib))


@case("VDCU-FUN-006", "CU 共享模式创建（CU=0）", "FUN")
def fun_006(ctx):
    ctx.destroy_card(ctx.card)
    cu, mem = baseline(ctx)
    rc, out = ctx.hy("virtual", "--create-vdevices", "2",
                     "--vdevice-compute-units", "0,0", "-d", str(ctx.card))
    ctx.expect(rc == 0, "创建失败: %s" % out[-300:])

    vs = ctx.vdevs_on(ctx.card)
    ctx.expect(len(vs) == 2, "应创建 2 个 vDCU")
    ctx.expect(all(v["cu"] == 0 for v in vs), "CU 共享模式应显示 Compute units: 0")
    for v in vs:
        ctx.expect(near(v["mem"], mem // 2), "显存应为物理的一半")


@case("VDCU-FUN-007", "显存总和超物理显存应被拒绝", "FUN")
def fun_007(ctx):
    ctx.destroy_card(ctx.card)
    rc, out = ctx.hy("virtual", "--create-vdevices", "2",
                     "--vdevice-memory-size", "40960,40960", "-d", str(ctx.card))
    ctx.expect(rc != 0, "显存超配应返回非 0")
    ctx.expect("exceed" in out.lower() or "failed" in out.lower(), "缺少明确错误: %s" % out[-300:])
    ctx.expect(len(ctx.vdevs_on(ctx.card)) == 0, "失败后不应产生 vDCU")


@case("VDCU-FUN-008", "CU 总和超物理 CU 应被拒绝", "FUN")
def fun_008(ctx):
    ctx.destroy_card(ctx.card)
    cu, _ = baseline(ctx)
    over = cu // 2 + 10
    rc, out = ctx.hy("virtual", "--create-vdevices", "2",
                     "--vdevice-compute-units", "%d,%d" % (over, over),
                     "-d", str(ctx.card))
    ctx.expect(rc != 0, "CU 超配应返回非 0")
    ctx.expect(len(ctx.vdevs_on(ctx.card)) == 0, "失败后不应产生 vDCU")


@case("VDCU-FUN-009", "更新 vDCU 的 CU 与显存", "FUN")
def fun_009(ctx):
    ctx.destroy_card(ctx.card)
    cu, mem = baseline(ctx)
    rc, out = ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    ctx.expect(rc == 0, "前置创建失败")
    vs = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])
    vid = vs[0]["id"]

    rc, out = ctx.hy("virtual", "--update-vdevice", "30,*", "-v", str(vid))
    ctx.expect(rc == 0 and "successfully" in out.lower(), "更新 CU 失败: %s" % out[-300:])
    cur = [v for v in ctx.vdevs() if v["id"] == vid][0]
    ctx.expect(cur["cu"] == 30, "更新后 CU 应为 30，实际 %s" % cur["cu"])
    ctx.expect(near(cur["mem"], mem // 2), "仅更新 CU 时显存不应变化")

    rc, out = ctx.hy("virtual", "--update-vdevice", "*,16384", "-v", str(vid))
    ctx.expect(rc == 0, "更新显存失败: %s" % out[-300:])
    cur = [v for v in ctx.vdevs() if v["id"] == vid][0]
    ctx.expect(near(cur["mem"], 16384 * 1024 * 1024), "更新后显存应为 16 GiB")
    ctx.expect(cur["cu"] == 30, "仅更新显存时 CU 不应变化")

    rc, out = ctx.hy("virtual", "--update-vdevice", "20,8192", "-v", str(vid))
    ctx.expect(rc == 0, "同时更新失败: %s" % out[-300:])
    cur = [v for v in ctx.vdevs() if v["id"] == vid][0]
    ctx.expect(cur["cu"] == 20 and near(cur["mem"], 8192 * 1024 * 1024),
               "同时更新结果不符: cu=%s mem=%s" % (cur["cu"], cur["mem"]))


@case("VDCU-FUN-010", "更新后资源超限应被拒绝", "FUN")
def fun_010(ctx):
    ctx.destroy_card(ctx.card)
    cu, mem = baseline(ctx)
    ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    vs = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])
    vid = vs[0]["id"]

    rc, out = ctx.hy("virtual", "--update-vdevice", "%d,*" % (cu - 10), "-v", str(vid))
    ctx.expect(rc != 0, "CU 超限更新应失败（合计 %d > %d）" % (cu - 10 + cu // 2, cu))

    rc, out = ctx.hy("virtual", "--update-vdevice", "*,40960", "-v", str(vid))
    ctx.expect(rc != 0, "显存超限更新应失败")

    cur = [v for v in ctx.vdevs() if v["id"] == vid][0]
    ctx.expect(cur["cu"] == cu // 2 and near(cur["mem"], mem // 2),
               "失败的更新不应改变原配置")


def _assert_in_use_rejected(ctx, args, description):
    rc, out = ctx.hy(*args)
    text = out.lower()
    rejected = rc != 0 or any(word in text for word in (
        "in use", "may be in use", "being used", "exist and may",
    ))
    ctx.expect(rejected, "%s 未被拒绝: rc=%d output=%s" %
               (description, rc, out[-300:]))


@case("VDCU-FUN-016", "使用中 vDCU 销毁应被拒绝", "FUN", priority="P1")
def fun_016(ctx):
    ctx.destroy_card(ctx.card)
    rc, out = ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    ctx.expect(rc == 0, "前置创建失败: %s" % out[-300:])
    vdev = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])[0]["id"]
    holder = "tz_fun016_holder"
    ctx.sh(["docker", "rm", "-f", holder], timeout=60)
    started = False
    try:
        rc, out = ctx.docker([vdev], HOLD_DEVICE, name=holder, detach=True,
                             env={"HIP_VISIBLE_DEVICES": "0"}, timeout=60)
        ctx.expect(rc == 0, "占用容器启动失败: %s" % out[-300:])
        started = True
        rc, state = ctx.sh(["docker", "inspect", "-f", "{{.State.Running}}", holder], timeout=60)
        ctx.expect(rc == 0 and state.strip().lower() == "true", "占用容器未保持运行: %s" % state[-200:])
        # 等待显存占用被驱动记录（HIP 上下文创建完成，宿主机侧轮询）
        used = False
        for _ in range(40):
            rc2, info = ctx.hy("virtual", "--show-vdevice-info", "-d", str(ctx.card))
            body = [b for b in info.split("\n") if "Usage memory" in b]
            if rc2 == 0 and body and "0 bytes" not in body[0]:
                used = True
                break
            time.sleep(2)
        ctx.expect(used, "等待容器占用设备超时（宿主机 Usage memory 未出现）")
        _assert_in_use_rejected(
            ctx, ["virtual", "--destroy-vdevices", "-v", str(vdev)],
            "使用中的 vDCU 销毁")
        ctx.expect(any(v["id"] == vdev for v in ctx.vdevs_on(ctx.card)),
                   "拒绝销毁后 vDCU 不应消失")
    finally:
        if started:
            ctx.sh(["docker", "rm", "-f", holder], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-FUN-017", "使用中 vDCU 更新应被拒绝且保持原配置", "FUN", priority="P1")
def fun_017(ctx):
    ctx.destroy_card(ctx.card)
    rc, out = ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    ctx.expect(rc == 0, "前置创建失败: %s" % out[-300:])
    vdev = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])[0]
    holder = "tz_fun017_holder"
    ctx.sh(["docker", "rm", "-f", holder], timeout=60)
    started = False
    try:
        rc, out = ctx.docker([vdev["id"]], HOLD_DEVICE, name=holder, detach=True,
                             env={"HIP_VISIBLE_DEVICES": "0"}, timeout=60)
        ctx.expect(rc == 0, "占用容器启动失败: %s" % out[-300:])
        started = True
        rc, state = ctx.sh(["docker", "inspect", "-f", "{{.State.Running}}", holder], timeout=60)
        ctx.expect(rc == 0 and state.strip().lower() == "true", "占用容器未保持运行: %s" % state[-200:])
        # 等待显存占用被驱动记录（HIP 上下文创建完成，宿主机侧轮询）
        used = False
        for _ in range(40):
            rc2, info = ctx.hy("virtual", "--show-vdevice-info", "-d", str(ctx.card))
            body = [b for b in info.split("\n") if "Usage memory" in b]
            if rc2 == 0 and body and "0 bytes" not in body[0]:
                used = True
                break
            time.sleep(2)
        ctx.expect(used, "等待容器占用设备超时（宿主机 Usage memory 未出现）")
        _assert_in_use_rejected(
            ctx, ["virtual", "--update-vdevice", "30,*", "-v", str(vdev["id"])],
            "使用中的 vDCU 更新")
        cur = [v for v in ctx.vdevs_on(ctx.card) if v["id"] == vdev["id"]][0]
        ctx.expect(cur["cu"] == vdev["cu"] and cur["mem"] == vdev["mem"],
                   "拒绝更新后 vDCU 配置被改变")
    finally:
        if started:
            ctx.sh(["docker", "rm", "-f", holder], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-FUN-014", "使用中的 vDCU 禁止销毁与更新", "FUN", priority="P0")
def fun_014(ctx):
    """v1.2 FUN-014：占用容器保持期间，销毁与更新均被拒绝（提示设备使用中）；
    容器销毁后操作成功、资源回收（R1 明确：使用中的虚拟设备无法销毁或更新）。"""
    ctx.destroy_card(ctx.card)
    rc, out = ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    ctx.expect(rc == 0, "前置创建失败: %s" % out[-300:])
    vdev = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])[0]["id"]
    holder = "tz_fun014_holder"
    ctx.sh(["docker", "rm", "-f", holder], timeout=60)
    try:
        rc, out = ctx.docker([vdev], HOLD_DEVICE, name=holder, detach=True,
                             env={"HIP_VISIBLE_DEVICES": "0"}, timeout=60)
        ctx.expect(rc == 0, "占用容器启动失败: %s" % out[-300:])
        rc, state = ctx.sh(["docker", "inspect", "-f",
                            "{{.State.Running}}", holder], timeout=60)
        ctx.expect(rc == 0 and state.strip().lower() == "true",
                   "占用容器未保持运行: %s" % state[-200:])
        # 等待 HIP 上下文建立、驱动记录占用（查询路径损坏时以命令拒绝结果为准）
        time.sleep(10)
        _assert_in_use_rejected(
            ctx, ["virtual", "--destroy-vdevices", "-v", str(vdev)],
            "使用中的 vDCU 销毁")
        _assert_in_use_rejected(
            ctx, ["virtual", "--update-vdevice", "40,8192", "-v", str(vdev)],
            "使用中的 vDCU 更新")
        ctx.expect(any(v["id"] == vdev for v in ctx.vdevs_on(ctx.card)),
                   "拒绝销毁后 vDCU 不应消失")
        # 容器销毁后应可正常销毁，资源回收
        ctx.sh(["docker", "rm", "-f", holder], timeout=60)
        time.sleep(3)   # 容器退出后引用释放延迟（实测约 1~3 秒）
        ctx.expect(ctx.destroy_card(ctx.card),
                   "容器销毁后 vDCU 应可销毁（资源回收）")
        ctx.steps.append("容器销毁后资源回收验证通过：卡 %d 已无 vDCU" % ctx.card)
    finally:
        ctx.sh(["docker", "rm", "-f", holder], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-FUN-021", "同一 vDCU 重复绑定到多个容器", "FUN",
      expected_fail=True, issue="ISSUE-012")
def fun_021(ctx):
    """v1.2 FUN-021：同一 vDCU 应仅允许一个容器绑定（互斥）。
    实测：驱动未做互斥，两个容器均可绑定并运行（ISSUE-012，同 CTR-012）。"""
    ctx.destroy_card(ctx.card)
    rc, out = ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    ctx.expect(rc == 0, "前置创建失败: %s" % out[-300:])
    vdev = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])[0]["id"]
    names = ["tz_fun021_a", "tz_fun021_b"]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    try:
        rc, out = ctx.docker([vdev], HOLD_DEVICE, name=names[0], detach=True,
                             env={"HIP_VISIBLE_DEVICES": "0"}, timeout=60)
        ctx.expect(rc == 0, "容器 A 启动失败: %s" % out[-300:])
        time.sleep(10)
        rc, out = ctx.docker(
            [vdev],
            "python - <<'PY'\nimport torch\n"
            "print('DEVICE_COUNT', torch.cuda.device_count())\n"
            "x = torch.randn(64, 64, device='cuda')\n"
            "y = x @ x\n"
            "torch.cuda.synchronize()\n"
            "print('MATMUL_OK')\nPY",
            name=names[1], detach=True, env={"HIP_VISIBLE_DEVICES": "0"},
            timeout=120)
        ctx.expect(rc == 0, "容器 B 启动失败: %s" % out[-300:])
        rc, out = ctx.sh(["docker", "wait", names[1]], timeout=180)
        ctx.expect(rc == 0, "容器 B 等待失败")
        rc, logs = ctx.sh(["docker", "logs", names[1]], timeout=60)
        ctx.steps.append("容器 B（复用同一 vDCU 配置）输出：\n" + logs[-400:])
        rc, state = ctx.sh(["docker", "inspect", "-f",
                            "{{.State.Running}}", names[0]], timeout=60)
        ctx.steps.append("容器 A 负载状态：%s" % state.strip())
        rejected = (rc != 0) or ("MATMUL_OK" not in logs)
        ctx.expect(rejected,
                   "同一 vDCU 可被两个容器同时绑定（互斥未实现，ISSUE-012）")
    finally:
        for n in names:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-FUN-011", "销毁指定 vDCU", "FUN")
def fun_011(ctx):
    ctx.destroy_card(ctx.card)
    ctx.hy("virtual", "--create-vdevices", "4", "-d", str(ctx.card))
    vs = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])
    ctx.expect(len(vs) == 4, "前置：应创建 4 个 vDCU")
    victim = vs[1]["id"]

    rc, out = ctx.hy("virtual", "--destroy-vdevices", "-v", str(victim))
    ctx.expect(rc == 0, "销毁失败: %s" % out[-300:])
    left = ctx.vdevs_on(ctx.card)
    ctx.expect(len(left) == 3, "应剩余 3 个 vDCU，实际 %d" % len(left))
    ctx.expect(all(v["id"] != victim for v in left), "被销毁的 vDCU 仍存在")


@case("VDCU-FUN-012", "销毁指定物理卡的全部 vDCU", "FUN")
def fun_012(ctx):
    ctx.destroy_card(ctx.card)
    ctx.destroy_card(ctx.card2)
    ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card2))
    ctx.expect(len(ctx.vdevs_on(ctx.card)) == 2, "前置：卡 A 应有 2 个 vDCU")
    ctx.expect(len(ctx.vdevs_on(ctx.card2)) == 2, "前置：卡 B 应有 2 个 vDCU")

    rc, out = ctx.hy("virtual", "--destroy-vdevices", "-d", str(ctx.card))
    ctx.expect(rc == 0, "按卡销毁失败: %s" % out[-300:])
    ctx.expect(len(ctx.vdevs_on(ctx.card)) == 0, "卡 A 的 vDCU 未清空")
    ctx.expect(len(ctx.vdevs_on(ctx.card2)) == 2, "卡 B 的 vDCU 不应受影响")
    ctx.destroy_card(ctx.card2)


@case("VDCU-FUN-013", "销毁全部 vDCU", "FUN", priority="P2")
def fun_013(ctx):
    if not ctx.allow_destructive:
        ctx.skip("全局销毁会清除生产卡上的既有 vDCU，需 --allow-destructive 后执行")
    rc, out = ctx.hy("virtual", "--destroy-vdevices")
    ctx.expect(rc == 0, "销毁全部失败: %s" % out[-300:])
    ctx.expect(len(ctx.vdevs()) == 0, "仍存在 vDCU")


@case("VDCU-FUN-015", "非 root 用户执行管理命令应被拒绝", "FUN", priority="P1")
def fun_015(ctx):
    ctx.destroy_card(ctx.card)
    rc, out = ctx.sh(["runuser", "-u", "tangzhao", "--", ctx.hy_smi,
                      "virtual", "--create-vdevices", "2", "-d", str(ctx.card)])
    ctx.expect(rc != 0, "非 root 创建应失败")
    ctx.expect("root" in out.lower(), "应提示需要 root 权限: %s" % out[-200:])
    ctx.expect(len(ctx.vdevs_on(ctx.card)) == 0, "非 root 不应创建出 vDCU")


@case("VDCU-FUN-018", "加密虚拟机状态查询", "FUN", priority="P1")
def fun_018(ctx):
    rc, out = ctx.hy("virtual", "--show-encrypted-status")
    ctx.expect(rc == 0, "查询失败: %s" % out[-200:])
    ctx.expect("Encrypted status" in out, "输出缺少状态字段: %s" % out[-200:])


@case("VDCU-FUN-019", "加密虚拟机开关可切换并恢复", "FUN", priority="P2")
def fun_019(ctx):
    rc, before = ctx.hy("virtual", "--show-encrypted-status")
    ctx.expect(rc == 0, "读取初始加密状态失败: %s" % before[-200:])
    enabled = "enabled" in before.lower() or "开启" in before
    rc, out = ctx.hy("virtual", "--disable-encrypted")
    ctx.expect(rc == 0, "关闭加密虚拟机失败: %s" % out[-300:])
    rc, out = ctx.hy("virtual", "--show-encrypted-status")
    ctx.expect(rc == 0 and "Encrypted status" in out,
               "关闭后状态不正确: %s" % out[-200:])
    try:
        rc, out = ctx.hy("virtual", "--enable-encrypted")
        ctx.expect(rc == 0, "恢复加密虚拟机失败: %s" % out[-300:])
    finally:
        ctx.hy("virtual", "--enable-encrypted")
        if not enabled:
            ctx.hy("virtual", "--disable-encrypted")


@case("VDCU-FUN-020", "多物理卡并行切分与编号连续性", "FUN", priority="P1")
def fun_020(ctx):
    ctx.destroy_card(ctx.card)
    ctx.destroy_card(ctx.card2)
    rc1, _ = ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    rc2, _ = ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card2))
    ctx.expect(rc1 == 0 and rc2 == 0, "多卡创建失败")

    vs = ctx.vdevs()
    on_a = [v for v in vs if v["phys"] == ctx.card]
    on_b = [v for v in vs if v["phys"] == ctx.card2]
    ctx.expect(len(on_a) == 2 and len(on_b) == 2, "两卡各应有 2 个 vDCU")
    ids = [v["id"] for v in vs]
    ctx.expect(len(ids) == len(set(ids)), "vDCU 编号出现重复: %s" % ids)
    ctx.expect(min(ids) >= 0 and max(ids) < 64, "vDCU 编号越界: %s" % ids)

    ctx.destroy_card(ctx.card)
    ctx.destroy_card(ctx.card2)


@case("VDCU-FUN-022", "版本查询与帮助信息", "FUN", priority="P2")
def fun_022(ctx):
    rc, out = ctx.hy("--version")
    ctx.expect(rc == 0 and len(out.strip()) > 0, "版本查询失败")

    rc, out = ctx.hy("virtual", "-h")
    ctx.expect(rc == 0, "帮助查询失败")
    for opt in ["--create-vdevices", "--destroy-vdevices", "--show-vdevice-info"]:
        ctx.expect(opt in out, "帮助信息缺少 %s" % opt)
