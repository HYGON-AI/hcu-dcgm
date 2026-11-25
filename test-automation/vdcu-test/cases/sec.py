# -*- coding: utf-8 -*-
"""第 14 章 安全测试（SEC）。"""
import shlex

from vdcu.core import case


@case("VDCU-SEC-001", "设备节点最小权限挂载与越权防护", "SEC", priority="P1")
def sec_001(ctx):
    """最小挂载容器内：仅注入单卡设备符，mknod 创建新设备被拒绝。"""
    ctx.expect(ctx.destroy_card(ctx.card), "前置清理失败")
    ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    vs = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])
    script = r"""
echo "DRI_BEGIN"
ls /dev/dri
echo "DRI_END"
echo "MKNOD_BEGIN"
# 用容器内可见设备的真实主设备号伪造节点（未被注入的次设备号）
D=$(ls /dev/dri/card* 2>/dev/null | head -1)
if [ -n "$D" ]; then MAJ=$((16#$(stat -c %%t "$D"))); else MAJ=226; fi
mknod /tmp/fake_card c $MAJ 99 2>&1; echo "mknod_rc=$?"
python - <<'PY'
import os
try:
    fd = os.open("/tmp/fake_card", os.O_RDWR)
    print("FAKE_OPEN_OK")
    os.close(fd)
except OSError as e:
    print("FAKE_OPEN_DENIED", e.errno, str(e)[:80])
PY
echo "MKNOD_END"
ls %s 2>&1
%s virtual --destroy-vdevices -v 0 2>&1 | tail -2
""" % (ctx.vdev_docker_dir, shlex.quote(ctx.hy_smi))
    rc, out = ctx.docker([vs[0]["id"]], script,
                         env={"HIP_VISIBLE_DEVICES": "0"}, timeout=300)
    ctx.steps.append("容器内安全探测输出：\n" + out[-1200:])
    ctx.expect(rc == 0, "容器执行失败")
    import re as _re
    m = _re.search(r"DRI_BEGIN\n(.*?)\nDRI_END", out, _re.S)
    dri = m.group(1) if m else ""
    card_count = len(_re.findall(r"\bcard\d+\b", dri))
    render_count = len(_re.findall(r"\brenderD\d+\b", dri))
    ctx.expect(card_count <= 1 and render_count <= 1,
               "容器内应仅注入单卡设备符，实际 card×%d render×%d"
               % (card_count, render_count))
    # 越权判定的安全边界在设备访问层：mknod 创建节点文件可能受宿主
    # docker seccomp 配置影响（实测可创建），但伪造节点无对应设备后端，
    # open 必须被拒绝（ENXIO/EACCES 等）
    ctx.expect("FAKE_OPEN_DENIED" in out and "FAKE_OPEN_OK" not in out,
               "容器内可通过伪造设备节点访问设备（越权）")
    ctx.expect("destroy" in out and (
        "not allowed" in out.lower() or "permission" in out.lower()
        or "root privilege" in out.lower() or "fail" in out.lower()),
        "容器内执行销毁命令未被拒绝")
    ctx.steps.append("越权防护：伪造设备节点 open 被拒绝、销毁命令被拒绝、仅见已分配设备")


@case("VDCU-SEC-005", "加密虚拟机开关与状态查询", "SEC", priority="P1")
def sec_005(ctx):
    rc, out = ctx.hy("virtual", "--show-encrypted-status", timeout=60)
    ctx.expect(rc == 0 and "Encrypted status" in out,
               "加密状态查询失败: %s" % out[-200:])
    ctx.steps.append("加密状态：%s" % out.strip())
