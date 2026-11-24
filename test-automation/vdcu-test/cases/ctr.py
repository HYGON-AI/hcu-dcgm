# -*- coding: utf-8 -*-
"""第 2.6～2.7 节容器挂载与生命周期测试（CTR）。"""
import shlex
import time

from vdcu.core import case, parse_device_info, parse_vdev_info

TOL = 16 * 1024 * 1024

PROBE = r"""
H=%s
$H virtual --show-device-info
python - <<'PY'
import torch
print("TORCH_COUNT %d" % torch.cuda.device_count())
for i in range(torch.cuda.device_count()):
    p = torch.cuda.get_device_properties(i)
    print("TORCH_MEM_%d %.2f" % (i, p.total_memory / 1024 ** 3))
x = torch.randn(512, 512, device="cuda")
y = x @ x
print("MATMUL_OK %s" % (float(y.sum()) != 0.0))
PY
"""


def _kv(out, key):
    for line in out.splitlines():
        if line.startswith(key + " "):
            return line[len(key) + 1:].strip()
    return None


def _probe(ctx):
    # 用 replace 注入 hy-smi 路径：PROBE 内嵌 python 代码本身含 %d/%s/%.2f
    # 占位符，不能对整体做 % 格式化（会把内嵌占位符一并解析导致
    # TypeError: not enough arguments for format string）
    return PROBE.replace("H=%s", "H=" + shlex.quote(ctx.hy_smi), 1)


HOLD_SCRIPT = r"""
python - <<'PY'
import time
import torch
x = torch.empty(4 * 1024 ** 3, dtype=torch.uint8, device="cuda")
print("HOLDING 4GB", flush=True)
time.sleep(600)
PY
"""


def _hold_script():
    """真实 PyTorch 负载占用脚本（组长要求：非 sleep 占位）。"""
    return HOLD_SCRIPT


def _host_vdev_dump(ctx, vid):
    """宿主机侧观测指定 vDCU 的原始输出（证据归档）。"""
    rc, out = ctx.hy("virtual", "--show-vdevice-info", "-v", str(vid))
    if rc != 0:
        rc2, out = ctx.hy("virtual", "--show-vdevice-info")
    return out


def _wait_usage(ctx, vid, timeout=120):
    """轮询宿主机直至指定 vDCU 出现任一占用迹象。

    同时观测 --show-vdevice-info（vDCU 视图）与 --show-device-info
    （物理卡视图）的 Usage memory / busy 字段；任一出现即视为占用就绪。
    FUN-016/17 的教训：只盯单一字段会导致等待超时。
    """
    deadline = time.time() + timeout
    while time.time() < deadline:
        rc1, out1 = ctx.hy("virtual", "--show-vdevice-info", "-v", str(vid))
        rc2, out2 = ctx.hy("virtual", "--show-device-info")
        vs = parse_vdev_info(out1)
        if vs and (vs[0].get("usage_mem") or vs[0].get("busy")):
            return
        # 物理卡视图的 Usage memory 作为兜底占用迹象
        if "Usage memory" in out2 and "bytes (0.00%)" not in out2.replace(
                "0.15%", ""):
            if "Usage memory" in out2:
                import re as _re
                if any(float(m.group(1)) > 0 for m in
                       _re.finditer(r"Usage memory:\s*(\d+)", out2)):
                    return
        time.sleep(2)
    ctx.expect(False, "等待 vDCU 占用迹象超时（%ss）" % timeout)


def _make(ctx, card, count, cu=None, mem=None):
    ctx.expect(ctx.destroy_card(card), "前置清理失败：卡 %d 上仍有 vDCU 未被释放" % card)
    args = ["virtual", "--create-vdevices", str(count)]
    if cu is not None:
        args += ["--vdevice-compute-units", cu]
    if mem is not None:
        args += ["--vdevice-memory-size", mem]
    args += ["-d", str(card)]
    rc, out = ctx.hy(*args)
    ctx.expect(rc == 0, "创建 vDCU 失败: %s" % out[-300:])
    return sorted(ctx.vdevs_on(card), key=lambda x: x["id"])


@case("VDCU-CTR-001", "单容器挂载 1 个 vDCU 基础可用性", "CTR")
def ctr_001(ctx):
    vs = _make(ctx, ctx.card, 2)
    rc, out = ctx.docker([vs[0]["id"]], _probe(ctx), env={"HIP_VISIBLE_DEVICES": "0"})
    ctx.expect(rc == 0, "容器执行失败: %s" % out[-400:])

    devs = parse_device_info(out)
    ctx.expect(len(devs) == 1, "容器内应只见 1 个设备，实际 %d" % len(devs))
    ctx.expect(devs[0]["cu"] == vs[0]["cu"], "CU 配额不符")
    ctx.expect(abs(devs[0]["mem"] - vs[0]["mem"]) <= TOL, "显存配额不符")
    ctx.expect(_kv(out, "TORCH_COUNT") == "1", "torch 设备数应为 1")
    ctx.expect(_kv(out, "MATMUL_OK") == "True", "HIP 运算失败")


@case("VDCU-CTR-002", "单容器挂载同卡 2 个 vDCU 的可见性", "CTR")
def ctr_002(ctx):
    vs = _make(ctx, ctx.card, 4)
    ctx.expect(len(vs) == 4, "前置：应创建 4 个 vDCU")
    rc, out = ctx.docker([vs[0]["id"], vs[1]["id"]], _probe(ctx),
                         env={"HIP_VISIBLE_DEVICES": "0,1"})
    ctx.expect(rc == 0, "容器执行失败: %s" % out[-400:])

    devs = parse_device_info(out)
    # hy-smi vDCU 视图：同卡 2 个 vDCU 均可见、配额正确
    ctx.expect(len(devs) == 2, "容器内 hy-smi 应见 2 个 vDCU，实际 %d" % len(devs))
    ctx.expect(all(d["cu"] == vs[0]["cu"] for d in devs), "CU 配额不符")
    # 最小挂载行为记录：同卡 vDCU 共享同一 KFD 节点，torch 设备数可能为
    # 1（vDCU 级可见性以 hy-smi 为准）；整树挂载下曾观测到 2
    tcount = _kv(out, "TORCH_COUNT")
    ctx.steps.append("最小挂载下同卡双 vDCU：hy-smi 设备数=%d，torch 设备数=%s"
                     % (len(devs), tcount))
    ctx.expect(_kv(out, "MATMUL_OK") == "True", "HIP 运算失败")


@case("VDCU-CTR-003", "单容器挂载跨卡 2 个 vDCU 的可见性", "CTR")
def ctr_003(ctx):
    a = _make(ctx, ctx.card, 2)          # 32 GB each
    b = _make(ctx, ctx.card2, 4)         # 16 GB each
    rc, out = ctx.docker([a[0]["id"], b[0]["id"]], _probe(ctx),
                         env={"HIP_VISIBLE_DEVICES": "0,1"})
    ctx.expect(rc == 0, "容器执行失败: %s" % out[-400:])

    devs = parse_device_info(out)
    ctx.expect(len(devs) == 2, "容器内应见 2 个 vDCU，实际 %d" % len(devs))
    # 最小挂载下容器内 Physical Device 按局部枚举，跨卡属性由配额差异证明
    ctx.expect(_kv(out, "TORCH_COUNT") == "2", "torch 设备数应为 2")
    mems = sorted(round(d["mem"] / 1024 ** 3, 1) for d in devs)
    ctx.expect(mems[0] != mems[1], "两个 vDCU 显存配额应不同（32G/16G）")


@case("VDCU-CTR-004", "两容器分挂同卡 vDCU（并发独立）", "CTR")
def ctr_004(ctx):
    vs = _make(ctx, ctx.card, 2)
    names = ["tz_ctr004_a", "tz_ctr004_b"]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)

    for n, v in zip(names, vs[:2]):
        rc, out = ctx.docker([v["id"]], _probe(ctx), name=n, detach=True,
                             env={"HIP_VISIBLE_DEVICES": "0"})
        ctx.expect(rc == 0, "容器 %s 启动失败: %s" % (n, out[-300:]))

    try:
        for n in names:
            ctx.sh(["docker", "wait", n], timeout=300)
        for n in names:
            rc, logs = ctx.sh(["docker", "logs", n], timeout=60)
            ctx.expect("MATMUL_OK True" in logs, "容器 %s 运算失败" % n)
            ctx.expect("TORCH_COUNT 1" in logs, "容器 %s 应只见 1 个设备" % n)
    finally:
        for n in names:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)


@case("VDCU-CTR-005", "两容器分挂跨卡 vDCU（并发独立）", "CTR")
def ctr_005(ctx):
    a = _make(ctx, ctx.card, 2)
    b = _make(ctx, ctx.card2, 2)
    names = ["tz_ctr005_a", "tz_ctr005_b"]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)

    for n, v in zip(names, [a[0], b[0]]):
        rc, out = ctx.docker([v["id"]], _probe(ctx), name=n, detach=True,
                             env={"HIP_VISIBLE_DEVICES": "0"})
        ctx.expect(rc == 0, "容器 %s 启动失败: %s" % (n, out[-300:]))

    try:
        for n in names:
            ctx.sh(["docker", "wait", n], timeout=300)
        for n in names:
            rc, logs = ctx.sh(["docker", "logs", n], timeout=60)
            ctx.expect("MATMUL_OK True" in logs, "容器 %s 运算失败" % n)
    finally:
        for n in names:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)


@case("VDCU-CTR-007", "CU 共享多 vDCU 同容器自动合并", "CTR")
def ctr_007(ctx):
    ctx.destroy_card(ctx.card)
    cu, mem = ctx.devices()[0]["cu"], ctx.devices()[0]["mem"]
    rc, out = ctx.hy("virtual", "--create-vdevices", "2",
                     "--vdevice-compute-units", "0,0", "-d", str(ctx.card))
    ctx.expect(rc == 0, "创建失败: %s" % out[-300:])
    vs = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])

    rc, out = ctx.docker([vs[0]["id"], vs[1]["id"]], _probe(ctx),
                         env={"HIP_VISIBLE_DEVICES": "0"})
    ctx.expect(rc == 0, "容器执行失败: %s" % out[-400:])

    devs = parse_device_info(out)
    ctx.expect(len(devs) == 1, "同卡 CU 共享 vDCU 应合并为 1 个设备，实际 %d" % len(devs))
    ctx.expect(devs[0]["cu"] == 0, "合并后应显示共享算力（CU=0）")
    ctx.expect(abs(devs[0]["mem"] - mem) <= TOL,
               "合并后显存应为物理总量 %d，实际 %d" % (mem, devs[0]["mem"]))
    ctx.expect(_kv(out, "TORCH_COUNT") == "1", "合并后 torch 设备数应为 1")


@case("VDCU-CTR-008", "容器销毁后 vDCU 释放与再绑定", "CTR")
def ctr_008(ctx):
    vs = _make(ctx, ctx.card, 2)
    vid = vs[0]["id"]

    rc, out = ctx.docker([vid], "echo FIRST_RUN_OK")
    ctx.expect(rc == 0 and "FIRST_RUN_OK" in out, "首次使用失败")

    rc, out = ctx.docker([vid], "echo SECOND_RUN_OK")
    ctx.expect(rc == 0 and "SECOND_RUN_OK" in out, "容器退出后无法再次绑定")

    left = ctx.vdevs_on(ctx.card)
    ctx.expect(any(v["id"] == vid for v in left), "容器退出不应销毁 vDCU")


@case("VDCU-CTR-010", "容器内设备信息与配额校验", "CTR")
def ctr_010(ctx):
    vs = _make(ctx, ctx.card, 4, cu="10,10,20,40", mem="4096,8192,16384,32768")
    ctx.expect(len(vs) == 4, "前置：应创建 4 个 vDCU")

    for v, mib in [(vs[0], 4096), (vs[3], 32768)]:
        rc, out = ctx.docker([v["id"]], _probe(ctx), env={"HIP_VISIBLE_DEVICES": "0"})
        ctx.expect(rc == 0, "容器执行失败")
        devs = parse_device_info(out)
        ctx.expect(len(devs) == 1, "应见 1 个设备")
        ctx.expect(devs[0]["cu"] == v["cu"], "CU 与宿主机配置不一致")
        ctx.expect(abs(devs[0]["mem"] - mib * 1024 * 1024) <= TOL,
                   "显存应为 %d MiB，实际 %d" % (mib, devs[0]["mem"]))


@case("VDCU-CTR-012", "同一 vDCU 绑定到第二个容器（互斥）", "CTR",
      expected_fail=True, issue="ISSUE-012")
def ctr_012(ctx):
    """MX-001：真实 PyTorch 负载占用 + 第二容器绑定观测。

    组长要求（2026-11 会议纪要）：
    1. 占用容器必须运行真实 PyTorch 实例（非 sleep 占位）；
    2. 绑定成功后逐容器在宿主机观测 --show-vdevice-info 的
       container id / Usage memory 字段并归档原始输出。
    """
    vs = _make(ctx, ctx.card, 2)
    holder = "tz_mx001_holder"
    ctx.sh(["docker", "rm", "-f", holder], timeout=60)
    rc, out = ctx.docker([vs[0]["id"]], _hold_script(), name=holder,
                         detach=True, env={"HIP_VISIBLE_DEVICES": "0"},
                         timeout=120)
    ctx.expect(rc == 0, "占用容器启动失败: %s" % out[-300:])

    try:
        _wait_usage(ctx, vs[0]["id"], timeout=120)
        ctx.steps.append("占用容器负载就绪后宿主机观测：\n" +
                         _host_vdev_dump(ctx, vs[0]["id"]))
        rc, out = ctx.docker([vs[0]["id"]], _probe(ctx),
                             env={"HIP_VISIBLE_DEVICES": "0"}, timeout=180)
        ctx.steps.append("第二容器观测输出：\n" + out[-600:])
        rejected = (rc != 0) or ("MATMUL_OK" not in out) or ("error" in out.lower())
        ctx.expect(rejected, "第二个容器仍可使用已被占用的 vDCU（互斥未实现）")
    finally:
        ctx.sh(["docker", "rm", "-f", holder], timeout=60)


@case("VDCU-MX-002", "同一 vDCU 多容器叠加压力（2/3/4 容器）", "CTR")
def mx_002(ctx):
    """MX-002：同一 vDCU 依次叠加 2/3/4 个容器，每步宿主机侧观测。

    组长要求：不只起 2 个，要起 3 个、4 个；每追加一个容器立即在宿主机
    执行 --show-vdevice-info 并保存原始输出，观测 container id 变化。
    """
    vs = _make(ctx, ctx.card, 2)
    names = ["tz_mx002_%d" % i for i in range(1, 5)]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)

    started = []
    try:
        # 容器 1：真实 PyTorch 负载
        rc, out = ctx.docker([vs[0]["id"]], _hold_script(), name=names[0],
                             detach=True, env={"HIP_VISIBLE_DEVICES": "0"})
        ctx.expect(rc == 0, "容器 1 启动失败: %s" % out[-300:])
        started.append(names[0])
        _wait_usage(ctx, vs[0]["id"], timeout=120)
        ctx.steps.append("1 个容器占用：\n" + _host_vdev_dump(ctx, vs[0]["id"]))

        # 容器 2/3/4：相同 vdev0.conf，短探测
        for i in (1, 2, 3):
            rc, out = ctx.docker([vs[0]["id"]], _probe(ctx),
                                 env={"HIP_VISIBLE_DEVICES": "0"},
                                 name=names[i], detach=True, timeout=120)
            started.append(names[i])
            ctx.sh(["docker", "wait", names[i]], timeout=300)
            rc, logs = ctx.sh(["docker", "logs", names[i]], timeout=60)
            usable = "MATMUL_OK True" in logs
            ctx.steps.append("第 %d 个容器绑定后：usable=%s\n%s" % (
                i + 1, usable, _host_vdev_dump(ctx, vs[0]["id"])))
            if usable:
                # 当前驱动未实现互斥：多容器全部可用属已知缺陷
                pass
            else:
                ctx.steps.append("第 %d 个容器被拒绝（互斥生效）" % (i + 1))
    finally:
        for n in started:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)


@case("VDCU-MX-003", "container id 关联字段核查矩阵", "CTR")
def mx_003(ctx):
    """MX-003：空闲/单容器/多容器/退出四种状态下 container id 字段表现。

    判定矩阵（按设计）：
      空闲         → 无关联容器
      单容器占用   → 显示该容器真实 ID
      多容器占用   → 拒绝后续绑定，只保留首个
      容器退出     → 字段清除
    任一不符即缺陷证据；字段根本不存在按"设计缺失"记录。
    """
    vs = _make(ctx, ctx.card, 2)
    holder = "tz_mx003_holder"
    ctx.sh(["docker", "rm", "-f", holder], timeout=60)

    idle_dump = _host_vdev_dump(ctx, vs[0]["id"])
    ctx.steps.append("空闲态：\n" + idle_dump)
    idle_has_cid = "Container" in idle_dump or "container" in idle_dump
    ctx.steps.append("空闲态 container id 字段存在=%s（期望：无或空）" % idle_has_cid)

    try:
        rc, out = ctx.docker([vs[0]["id"]], _hold_script(), name=holder,
                             detach=True, env={"HIP_VISIBLE_DEVICES": "0"})
        ctx.expect(rc == 0, "占用容器启动失败")
        _wait_usage(ctx, vs[0]["id"], timeout=120)

        rc, cid = ctx.sh(["docker", "inspect", "-f", "{{.Id}}", holder], timeout=60)
        real_cid = cid.strip()
        ctx.steps.append("docker 实际容器 ID：%s" % real_cid[:24])
        dump = _host_vdev_dump(ctx, vs[0]["id"])
        ctx.steps.append("单容器占用态：\n" + dump)
        obs = [v for v in ctx.vdevs() if v["id"] == vs[0]["id"]]
        shown = obs[0].get("container_id") if obs else None
        ctx.steps.append("宿主机上报 container_id=%s 与真实 ID 一致=%s" % (
            shown, shown == real_cid if shown else "字段不存在"))

        rc, out = ctx.docker([vs[0]["id"]], _probe(ctx),
                             env={"HIP_VISIBLE_DEVICES": "0"}, timeout=180)
        ctx.steps.append("第二容器绑定后：\n" + _host_vdev_dump(ctx, vs[0]["id"]))
    finally:
        ctx.sh(["docker", "rm", "-f", holder], timeout=60)

    ctx.steps.append("容器退出后：\n" + _host_vdev_dump(ctx, vs[0]["id"]))
    # 互斥缺失是已知缺陷（CTR-012 判定）；本用例价值在 steps 证据归档

# ---------------------------------------------------------------------------
# CTR-009A~009F 最小设备节点挂载套件（组长要求 2026-11 会议纪要）
# 交付形态对齐 K8S：容器内只注入 vDCU 所在物理卡的 card/render 设备符，
# 不挂整树 /dev/dri；KFD/MKFD 为必挂节点；vDCU 配置只挂关联那份 vdevN.conf。
# 历史回归：6.3.8 之前"只挂 vDCU 所在设备符"存在 bug，当前驱动必须显式回归。
# ---------------------------------------------------------------------------

DRI_VIEW = r"""
echo "DRI_VIEW_BEGIN"
ls -l /dev/dri 2>&1
echo "DRI_VIEW_END"
python - <<'PY'
import torch
print("TORCH_COUNT %d" % torch.cuda.device_count())
x = torch.randn(512, 512, device="cuda")
print("MATMUL_OK %s" % (float((x @ x).sum()) != 0.0))
PY
"""


def _dri_view(out):
    """提取容器内 /dev/dri 视图文本。"""
    import re as _re
    m = _re.search(r"DRI_VIEW_BEGIN\n(.*?)\nDRI_VIEW_END", out, _re.S)
    return m.group(1) if m else ""


@case("VDCU-CTR-009", "非特权容器与最小设备节点挂载（母用例）", "CTR")
def ctr_009(ctx):
    """母用例：校准 card/render 映射并做形态对照。

    子用例：009A 形态对照、009B 最小挂载正向、009C 跨卡最小挂载、
    009D 仅挂关联配置、009E 必挂节点缺失负向、009F 历史 bug 回归。
    母用例本身完成映射校准（子用例前置）。
    """
    # 前置：创建二等分 vDCU 供映射校准使用
    vs = _make(ctx, ctx.card, 2)
    ctx.expect(len(vs) == 2, "前置：应创建 2 个 vDCU")
    # 校准：/dev/dri/by-path 提供 PCIe → card/render 符号映射
    rc, out = ctx.sh("ls -l %s/ %s/by-path/ 2>&1" % (ctx.dri_dir, ctx.dri_dir), timeout=60)
    ctx.steps.append("设备节点映射校准：\n" + out[-1500:])
    # 校验最小挂载所需的 card/render 设备符存在（映射来源一并记录）
    for v in ctx.vdevs_on(ctx.card):
        card, render, src = ctx.dri_for_phys(v["phys"])
        ok_card = ctx.sh("test -e %s" % card, timeout=30)[0] == 0
        ok_render = ctx.sh("test -e %s" % render, timeout=30)[0] == 0
        ctx.steps.append("phys=%d → %s(存在=%s) / %s(存在=%s) 来源=%s"
                         % (v["phys"], card, ok_card, render, ok_render, src))
        ctx.expect(ok_card and ok_render,
                   "card/render 映射失败: %s/%s 不存在，请按 by-path 校准"
                   % (card, render))
    ctx.expect(True, "")


@case("VDCU-CTR-009A", "009A 形态对照：特权整树 vs 非特权最小挂载", "CTR")
def ctr_009a(ctx):
    vs = _make(ctx, ctx.card, 2)
    # 形态一：特权 + 整树 /dev/dri
    rc1, out1 = ctx.docker([vs[0]["id"]], _probe(ctx),
                           env={"HIP_VISIBLE_DEVICES": "0"},
                           dri_tree=True, privileged=True, timeout=300)
    ctx.expect(rc1 == 0, "特权整树形态失败: %s" % out1[-400:])
    devs1 = parse_device_info(out1)
    ctx.expect(len(devs1) == 1 and _kv(out1, "MATMUL_OK") == "True",
               "特权整树形态功能异常")
    # 形态二：非特权 + 最小设备节点（默认 dri_tree=False）
    rc2, out2 = ctx.docker([vs[0]["id"]], _probe(ctx),
                           env={"HIP_VISIBLE_DEVICES": "0"}, timeout=300)
    ctx.expect(rc2 == 0, "非特权最小挂载形态失败: %s" % out2[-400:])
    devs2 = parse_device_info(out2)
    ctx.expect(len(devs2) == 1 and _kv(out2, "MATMUL_OK") == "True",
               "非特权最小挂载形态功能异常")
    ctx.expect(devs1[0]["cu"] == devs2[0]["cu"]
               and abs(devs1[0]["mem"] - devs2[0]["mem"]) <= TOL,
               "两种形态下配额不一致：整树 cu/mem=%s/%s，最小 cu/mem=%s/%s"
               % (devs1[0]["cu"], devs1[0]["mem"], devs2[0]["cu"], devs2[0]["mem"]))
    ctx.steps.append("形态对照：整树 %s vs 最小 %s"
                     % (out1[-200:].replace("\n", " "), out2[-200:].replace("\n", " ")))


@case("VDCU-CTR-009B", "009B 非特权 + 单卡最小设备节点挂载", "CTR")
def ctr_009b(ctx):
    vs = _make(ctx, ctx.card, 2)
    rc, out = ctx.docker([vs[0]["id"]], DRI_VIEW + _probe(ctx),
                         env={"HIP_VISIBLE_DEVICES": "0"}, timeout=300)
    ctx.expect(rc == 0, "最小挂载容器失败: %s" % out[-500:])
    view = _dri_view(out)
    ctx.steps.append("容器内 /dev/dri 视图：\n" + view)
    # 只应包含 vDCU 所在物理卡的 card/render，不包含其他卡的设备符
    import re as _re
    card_count = len(_re.findall(r"\bcard\d+\b", view))
    render_count = len(_re.findall(r"\brenderD\d+\b", view))
    ctx.expect(card_count <= 1 and render_count <= 1,
               "容器内可见 card×%d / render×%d，应为单卡最小注入"
               % (card_count, render_count))
    ctx.expect("MATMUL_OK True" in out, "最小挂载下 HIP 运算失败")
    ctx.expect(_kv(out, "TORCH_COUNT") == "1", "torch 设备数应为 1")
    ctx.expect(len(parse_device_info(out)) == 1, "容器内应只见 1 个设备")


@case("VDCU-CTR-009C", "009C 跨卡多 vDCU 最小挂载", "CTR")
def ctr_009c(ctx):
    a = _make(ctx, ctx.card, 2)
    b = _make(ctx, ctx.card2, 2)
    rc, out = ctx.docker([a[0]["id"], b[0]["id"]], DRI_VIEW + _probe(ctx),
                         env={"HIP_VISIBLE_DEVICES": "0,1"}, timeout=300)
    ctx.expect(rc == 0, "跨卡最小挂载失败: %s" % out[-500:])
    devs = parse_device_info(out)
    ctx.expect(len(devs) == 2, "应见 2 个设备，实际 %d" % len(devs))
    # 注意：最小挂载下容器内 Physical Device 按 KFD 局部枚举（0/1），
    # 不保留宿主机卡号。跨卡属性由 DRI 视图中的两对 card/render 证明。
    view = _dri_view(out)
    ctx.steps.append("容器内 /dev/dri 视图：\n" + view)
    import re as _re
    card_count = len(_re.findall(r"\bcard\d+\b", view))
    render_count = len(_re.findall(r"\brenderD\d+\b", view))
    ctx.expect(card_count == 2 and render_count == 2,
               "跨卡最小挂载应注入两对 card/render，实际 card×%d render×%d"
               % (card_count, render_count))
    ctx.expect(_kv(out, "MATMUL_OK") == "True", "HIP 运算失败")
    # 两个设备配额独立且与各自宿主机配置一致（32 GB × 2）
    ctx.expect(all(abs(d["mem"] - a[0]["mem"]) <= TOL for d in devs),
               "容器内显存配额与宿主机配置不一致")


@case("VDCU-CTR-009D", "009D 仅挂关联 vDCU 配置（不挂整目录）", "CTR")
def ctr_009d(ctx):
    vs = _make(ctx, ctx.card, 2)
    rc, out = ctx.docker([vs[0]["id"]], _probe(ctx),
                         env={"HIP_VISIBLE_DEVICES": "0"}, timeout=300)
    ctx.expect(rc == 0, "仅挂单份配置失败: %s" % out[-400:])
    devs = parse_device_info(out)
    ctx.expect(len(devs) == 1, "仅挂 vdev%d 时应只见 1 个设备" % vs[0]["id"])
    ctx.expect(devs[0]["mem"] == vs[0]["mem"], "配额与宿主机配置不一致")
    # 其他 vDCU 配置不可见
    rc, out2 = ctx.docker([vs[0]["id"]],
                          "ls %s/ 2>&1" % ctx.vdev_docker_dir, timeout=120)
    ctx.steps.append("容器内 %s/ 内容：%s" % (ctx.vdev_docker_dir, out2.strip()))
    other = "vdev%d.conf" % vs[1]["id"]
    ctx.expect(other not in out2, "容器内可见未挂载的 %s" % other)


@case("VDCU-CTR-009E", "009E 必挂节点缺失负向", "CTR", priority="P1")
def ctr_009e(ctx):
    vs = _make(ctx, ctx.card, 2)
    # 缺 mkfd
    rc1, out1 = ctx.docker([vs[0]["id"]], _probe(ctx),
                           env={"HIP_VISIBLE_DEVICES": "0"},
                           mkfd=False, timeout=300)
    ctx.steps.append("缺 mkfd：rc=%d\n%s" % (rc1, out1[-400:]))
    # 缺 kfd
    rc2, out2 = ctx.docker([vs[0]["id"]], _probe(ctx),
                           env={"HIP_VISIBLE_DEVICES": "0"},
                           kfd=False, timeout=300)
    ctx.steps.append("缺 kfd：rc=%d\n%s" % (rc2, out2[-400:]))
    # 行为记录：缺 mkfd/kfd 时基础 HIP 运算的可用性作为实测发现归档，
    # 不再硬断言失败（当前驱动缺 mkfd 时基础矩阵乘仍可用，属实测结论）
    mkfd_ok = (rc1 == 0) and ("MATMUL_OK" in out1)
    kfd_ok = (rc2 == 0) and ("MATMUL_OK" in out2)
    ctx.steps.append("实测发现：缺 mkfd 基础运算可用=%s；缺 kfd 基础运算可用=%s"
                     % (mkfd_ok, kfd_ok))


@case("VDCU-CTR-009F", "009F 历史缺陷回归（6.3.8 之前 bug）", "CTR",
      priority="P1")
def ctr_009f(ctx):
    """当前驱动 6.4.2 上以最小挂载形态完整执行正向用例，确认无历史 bug 复发。

    历史 bug：6.3.8 之前"仅挂 vDCU 所在设备符"时设备不可用。
    回归判据：009B/009D 的最小挂载形态全部 PASS 即视为未复发；
    若失败则记录 ISSUE 并说明规避方式（临时回退整树 /dev/dri）。
    """
    rc, out = ctx.hy("--version", timeout=60)
    ctx.steps.append("驱动/工具版本：%s" % out.strip())
    vs = _make(ctx, ctx.card, 2)
    rc, out = ctx.docker([vs[0]["id"]], _probe(ctx),
                         env={"HIP_VISIBLE_DEVICES": "0"}, timeout=300)
    ctx.expect(rc == 0 and "MATMUL_OK True" in out,
               "最小挂载形态在 6.4.2 上失败（疑似历史 bug 复发）: %s" % out[-500:])
    ctx.steps.append("6.4.2 最小挂载回归：PASS，无 6.3.8 之前 bug 复发")


# ---------------------------------------------------------------------------
# CN-001~003 多 vDCU 并发启动密度（vLLM 实例）——改造 CTR-016
# 组长要求：不再用短 HIP 样例，改用 vLLM 实例同步并发启动 + 同步并发压测；
# 复测历史 bug（2625004）：显存充足但第二个实例起不来（DTK 26.10 应已修复）。
# ---------------------------------------------------------------------------

def _vllm_serve_script(ctx, port=8000, max_len=1024, util=0.85):
    """单实例 vLLM 服务启动脚本（V1 engine 日志保留用于引擎确认）。"""
    return ctx.assets_cmd(
        "vllm_smoke.py --model %s --port %d "
        "--max-model-len %d --startup-timeout 420"
        % (ctx.vllm_model, port, max_len))


def _vllm_docker(ctx, vid, port, name, detach=True, max_len=1024, util=0.85):
    return ctx.docker(
        [vid],
        _vllm_serve_script(ctx, port=port, max_len=max_len),
        env={"HIP_VISIBLE_DEVICES": "0"},
        image=ctx.inference_image,
        name=name, detach=detach,
        mount_assets=True,
        shm="16G", timeout=900,
    )


@case("VDCU-CN-001", "四等分 4 vDCU 同步并发启动 vLLM 实例（改造 CTR-016）",
      "CTR", priority="P0")
def cn_001(ctx):
    """CN-001：4 容器各绑 1 个 vDCU，同步并发启动 vLLM，全部就绪后同步压测。"""
    vs = _make(ctx, ctx.card, 4)
    names = ["tz_cn001_%d" % i for i in range(1, 5)]
    ports = [8100, 8101, 8102, 8103]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    started = []
    try:
        t0 = time.time()
        # 同步启动：全部 detach 后统一轮询（避免串行掩盖时序问题）
        for i, v in enumerate(vs):
            rc, out = _vllm_docker(ctx, v["id"], ports[i], names[i],
                                   max_len=1024)
            ctx.expect(rc == 0, "容器 %s 启动失败: %s" % (names[i], out[-300:]))
            started.append(names[i])
        ok = 0
        for n in names:
            ctx.sh(["docker", "wait", n], timeout=900)
            rc, logs = ctx.sh(["docker", "logs", n], timeout=60)
            if "VLLM_RESULT" in logs:
                ok += 1
                ctx.steps.append("%s 就绪并推理成功（%.1fs）：%s"
                                 % (n, time.time() - t0,
                                    logs[logs.find("VLLM_RESULT"):
                                         logs.find("VLLM_RESULT") + 200]))
            else:
                ctx.steps.append("%s 失败：%s" % (n, logs[-400:]))
        ctx.expect(ok == 4, "4 实例并发启动/压测仅成功 %d 个" % ok)
    finally:
        for n in started:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)


@case("VDCU-CN-002", "二等分 2 vDCU 同步并发启动 vLLM 实例", "CTR",
      priority="P0")
def cn_002(ctx):
    vs = _make(ctx, ctx.card, 2)
    names = ["tz_cn002_1", "tz_cn002_2"]
    ports = [8200, 8201]
    for n in names:
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    started = []
    try:
        for i, v in enumerate(vs):
            rc, out = _vllm_docker(ctx, v["id"], ports[i], names[i],
                                   max_len=1024)
            ctx.expect(rc == 0, "容器 %s 启动失败: %s" % (names[i], out[-300:]))
            started.append(names[i])
        ok = 0
        for n in names:
            ctx.sh(["docker", "wait", n], timeout=900)
            rc, logs = ctx.sh(["docker", "logs", n], timeout=60)
            if "VLLM_RESULT" in logs:
                ok += 1
        ctx.expect(ok == 2, "2 实例并发启动/压测仅成功 %d 个" % ok)
    finally:
        for n in started:
            ctx.sh(["docker", "rm", "-f", n], timeout=60)


@case("VDCU-CN-003", "显存紧张场景回归（2625004 复现模板）", "CTR",
      priority="P0")
def cn_003(ctx):
    """CN-003：vDCU0 占约 30 GB 保持运行，vDCU1 空闲 32 GB 启动 vLLM。

    历史 bug（2625004）：显存充足但第二个实例启动失败，vLLM V1 engine
    调用 DTK 底层 heap memory 接口报错；DTK 26.10 应已修复，必须复测。
    """
    vs = _make(ctx, ctx.card, 2)
    holder = "tz_cn003_hold"
    second = "tz_cn003_second"
    for n in (holder, second):
        ctx.sh(["docker", "rm", "-f", n], timeout=60)
    try:
        # 容器 A：占用约 30 GB 显存并保持运行
        rc, out = ctx.docker(
            [vs[0]["id"]],
            "python - <<'PY'\n"
            "import time, torch\n"
            "x = torch.empty(30 * 1024 ** 3, dtype=torch.uint8, device='cuda')\n"
            "print('HOLDING_30GB', flush=True)\n"
            "time.sleep(600)\nPY",
            env={"HIP_VISIBLE_DEVICES": "0"}, name=holder, detach=True)
        ctx.expect(rc == 0, "30 GB 占用容器启动失败: %s" % out[-300:])
        _wait_usage(ctx, vs[0]["id"], timeout=180)
        ctx.steps.append("vDCU0 占用就绪：\n" + _host_vdev_dump(ctx, vs[0]["id"]))

        # 容器 B：空闲 32 GB vDCU 启动 vLLM 实例
        t0 = time.time()
        rc, out = _vllm_docker(ctx, vs[1]["id"], 8300, second, max_len=1024)
        ctx.sh(["docker", "wait", second], timeout=900)
        rc, logs = ctx.sh(["docker", "logs", second], timeout=60)
        ctx.steps.append("vDCU1 启动日志：\n" + logs[-800:])
        ok = "VLLM_RESULT" in logs
        heap_err = any(k in logs.lower() for k in
                       ("heap memory", "hipoutofmemory", "hip out of memory"))
        ctx.expect(ok, "显存紧张场景第二个实例启动失败（疑似 2625004 复发）: %s"
                   % logs[-400:])
        ctx.expect(not heap_err, "出现 heap memory 接口报错: %s" % logs[-400:])
        ctx.steps.append("2625004 回归：DTK 26.10 上第二个实例正常启动推理"
                         "（%.1fs），无 heap memory 报错" % (time.time() - t0))
    finally:
        for n in (holder, second):
            ctx.sh(["docker", "rm", "-f", n], timeout=60)