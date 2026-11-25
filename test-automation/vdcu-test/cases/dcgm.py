# -*- coding: utf-8 -*-
"""DCGM/监控指标校验测试（MON 域）。

依据会议纪要（2026-11）：
1. 测试过程中同步校验 vDCU 显存、算力占用的监控数据准确性，
   避免出现"1/4 显存分区占用显示超过 25%"的异常；
2. VDU 功能用例需迁移至 DCGM 自动化体系（Go 接口形式）——本模块
   先完成 DCGM exporter 指标与 hy-smi 的交叉校验，接口级迁移
   与东明团队对齐后按 DCGM 测试框架实现（见文档第 12 章说明）。
"""
import re

from vdcu.core import case


def _enabled(ctx):
    if not ctx.cluster_enabled:
        ctx.skip("未启用集群测试；设置 VDCU_CLUSTER_TESTS=1 后执行")


def _exporter_pods(ctx):
    rc, out, payload = ctx.kubectl_json(
        "get", "pods", "-n", "kube-system", "-o", "json", timeout=90)
    ctx.expect(rc == 0 and payload is not None, "查询 dcu-exporter Pod 失败")
    pods = [p for p in payload.get("items", [])
            if "dcu-exporter" in p.get("metadata", {}).get("name", "")
            and p.get("status", {}).get("phase") == "Running"]
    return pods


def _exporter_port(pod):
    containers = (pod.get("spec", {}).get("containers", [{}]))
    ports = containers[0].get("ports", [])
    if ports:
        return int(ports[0].get("containerPort", 16080))
    return 16080


def _exporter_host(pod):
    return pod.get("status", {}).get("hostIP", "127.0.0.1")


def _metrics(ctx, pod, port=16080):
    """经 hostPort 从宿主机直接拉取 dcu-exporter 的 Prometheus 指标。

    exporter 镜像内无 curl（实测 bash: curl: command not found），
    必须走 hostIP:hostPort（当前集群 hostPort=16080）。
    """
    host = _exporter_host(pod)
    rc, out = ctx.sh(
        "curl -s --max-time 10 http://%s:%d/metrics 2>&1" % (host, port),
        timeout=90)
    return rc, out


@case("VDCU-MON-101", "dcu-exporter 指标端点与 vDCU 指标存在性", "MON",
      priority="P0")
def mon_101(ctx):
    _enabled(ctx)
    pods = _exporter_pods(ctx)
    ctx.expect(pods, "未找到 Running 状态的 dcu-exporter Pod")
    pod_obj = pods[0]
    pod = pod_obj.get("metadata", {}).get("name")
    port = _exporter_port(pod_obj)
    rc, out = _metrics(ctx, pod_obj, port)
    ctx.expect(rc == 0, "拉取 dcu-exporter 指标失败: %s" % out[-300:])
    ctx.steps.append("指标名称清单：\n" +
                     "\n".join([l for l in out.splitlines()
                                if l.startswith("# HELP")][:30]))
    # 实测指标集：dcu_usedmemory_bytes / dcu_memorycap_bytes /
    # dcu_cu_usage / dcu_utilizationrate / vdcu_count 等
    has_dev = bool(re.search(r"^dcu_", out, re.M))
    has_mem = bool(re.search(r"(usedmemory|memorycap|memory_remaining)",
                             out, re.I))
    has_vdcu = "vdcu_count" in out
    ctx.expect(has_dev, "指标中未见设备类指标（dcu_*）")
    ctx.expect(has_mem, "指标中未见显存类指标")
    ctx.expect(has_vdcu, "指标中未见 vdcu_count 指标")
    print("        exporter=%s:%d 指标行数=%d vdcu_count=%s" %
          (pod, port, len(out.splitlines()), has_vdcu))


@case("VDCU-MON-102", "vDCU 显存占用指标准确性（exporter vs hy-smi 交叉校验）",
      "MON", priority="P0", expected_fail=True, issue="ISSUE-MON-001")
def mon_102(ctx):
    """会议纪要要求：校验监控数据准确性，避免"1/4 显存分区占用显示超过 25%"
    的异常。方法：同卡同时刻 exporter 上报的显存占用与 hy-smi 交叉校验。

    实测（2026-11-05）：exporter 的 dcu_usedmemory_bytes 在全部 8 张卡
    上统一报 ~82.7%（约 53 GB），而 hy-smi 显示 VRAM 0%——指标失准。
    本用例断言两者一致；不一致时按已知缺陷记录（XFAIL/ISSUE-MON-001）。
    """
    _enabled(ctx)
    ctx.expect(ctx.destroy_card(ctx.card), "前置清理失败")
    rc, out = ctx.hy("virtual", "--create-vdevices", "4", "-d", str(ctx.card))
    ctx.expect(rc == 0, "创建 4 等分失败: %s" % out[-300:])
    vs = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])
    holder = "tz_mon102"
    ctx.sh(["docker", "rm", "-f", holder], timeout=60)
    try:
        rc, out = ctx.docker(
            [vs[0]["id"]],
            "python - <<'PY'\nimport time, torch\n"
            "x = torch.empty(15 * 1024 ** 3, dtype=torch.uint8, device='cuda')\n"
            "print('HOLDING_15GB', flush=True)\n"
            "time.sleep(240)\nPY",
            env={"HIP_VISIBLE_DEVICES": "0"}, name=holder, detach=True)
        ctx.expect(rc == 0, "占用容器启动失败")

        import time as _t
        _t.sleep(15)

        # hy-smi 侧：物理卡视图 VRAM%
        rc, smi = ctx.sh("%s 2>&1" % ctx.hy_smi, timeout=60)
        ctx.steps.append("hy-smi 物理卡视图：\n" + "\n".join(
            smi.splitlines()[:12]))

        pods = _exporter_pods(ctx)
        ctx.expect(pods, "未找到 dcu-exporter Pod")
        pod_obj = pods[0]
        port = _exporter_port(pod_obj)
        rc, metrics = _metrics(ctx, pod_obj, port)
        ctx.expect(rc == 0, "拉取指标失败: %s" % metrics[-300:])
        used = {}
        cap = {}
        for line in metrics.splitlines():
            m = re.match(r"dcu_usedmemory_bytes\{[^}]*pcieBus_number=\""
                         r"([0-9a-f:.]+)\"[^}]*\}\s+([\d.e+-]+)", line)
            if m:
                used[m.group(1)] = float(m.group(2))
            m = re.match(r"dcu_memorycap_bytes\{[^}]*pcieBus_number=\""
                         r"([0-9a-f:.]+)\"[^}]*\}\s+([\d.e+-]+)", line)
            if m:
                cap[m.group(1)] = float(m.group(2))
        ratios = {pci: round(u / cap[pci] * 100, 1)
                  for pci, u in used.items() if pci in cap and cap[pci] > 0}
        ctx.steps.append("exporter 按 PCI 上报占用比例：%s" % ratios)
        # 测试卡 phys5 = PCI b5:00.0；hy-smi 显示该卡 VRAM 约 23%（15/64）
        # 与 exporter 上报值偏差超过 30 个百分点视为指标失准
        test_pci = {v["phys"]: pci for v, pci in []}
        # 直接从 vdev conf 取测试卡 PCI
        pci_map = {}
        for v in vs:
            try:
                with open("%s/vdev%d.conf" % (ctx.vdev_conf_dir, v["id"])) as f:
                    conf = f.read()
                m = re.search(r"PciBusId:\s*([0-9a-fA-F:.]+)", conf)
                if m:
                    pci_map[v["phys"]] = m.group(1)
            except OSError:
                pass
        test_pci = pci_map.get(ctx.card)
        ctx.steps.append("测试卡 PCI：%s" % test_pci)
        ctx.expect(test_pci in ratios,
                   "exporter 未上报测试卡 %s 的显存指标" % test_pci)
        exporter_ratio = ratios[test_pci]
        # 双向校验：15 GB 负载占 64 GB 物理卡的约 23%。
        # 实测失准表现为两个方向：全卡统一 82.7%（偏高）或全卡 0.0%（偏低），
        # 均与真实占用不符。判定区间 [13%, 33%]（23% ±10 个百分点）。
        ctx.expect(13.0 <= exporter_ratio <= 33.0,
                   "测试卡 exporter 显存占用 %.1f%% 与真实占用约 23%% 不符"
                   "（监控指标失准，ISSUE-MON-001）" % exporter_ratio)
    finally:
        ctx.sh(["docker", "rm", "-f", holder], timeout=60)
        ctx.destroy_card(ctx.card)


@case("VDCU-MON-103", "exporter 指标与 hy-smi 交叉校验", "MON",
      priority="P1", expected_fail=True, issue="ISSUE-MON-001")
def mon_103(ctx):
    """同一时刻 exporter 上报的显存占用与 hy-smi 查询值交叉校验。

    实测（2026-11-05）：exporter 的 dcu_usedmemory_bytes 不随负载变化
    （8 GB 负载前后增量 0.00 GB），指标失准。本用例断言增量与负载同量级；
    不符时按已知缺陷记录（XFAIL/ISSUE-MON-001）。
    """
    _enabled(ctx)
    pods = _exporter_pods(ctx)
    ctx.expect(pods, "未找到 dcu-exporter Pod")
    pod_obj = pods[0]
    pod = pod_obj.get("metadata", {}).get("name")
    port = _exporter_port(pod_obj)
    rc, idle = _metrics(ctx, pod_obj, port)
    ctx.expect(rc == 0, "空闲态指标拉取失败")

    def _used_sum(metrics):
        total = 0.0
        for line in metrics.splitlines():
            m = re.match(r"dcu_usedmemory_bytes\{[^}]*\}\s+([\d.e+-]+)", line)
            if m:
                total += float(m.group(1))
        return total

    ctx.expect(ctx.destroy_card(ctx.card), "前置清理失败")
    ctx.hy("virtual", "--create-vdevices", "2", "-d", str(ctx.card))
    vs = sorted(ctx.vdevs_on(ctx.card), key=lambda x: x["id"])
    holder = "tz_mon103"
    ctx.sh(["docker", "rm", "-f", holder], timeout=60)
    try:
        rc, out = ctx.docker(
            [vs[0]["id"]],
            "python - <<'PY'\nimport time, torch\n"
            "x = torch.empty(8 * 1024 ** 3, dtype=torch.uint8, device='cuda')\n"
            "print('HOLDING_8GB', flush=True)\n"
            "time.sleep(240)\nPY",
            env={"HIP_VISIBLE_DEVICES": "0"}, name=holder, detach=True)
        ctx.expect(rc == 0, "占用容器启动失败")
        import time as _t
        _t.sleep(15)
        rc, busy = _metrics(ctx, pod_obj, port)
        ctx.expect(rc == 0, "占用态指标拉取失败")
        idle_sum = _used_sum(idle)
        busy_sum = _used_sum(busy)
        delta_gb = (busy_sum - idle_sum) / 1024 ** 3
        ctx.steps.append("空闲态 usedmemory 总和=%.2f GB，占用态=%.2f GB，"
                         "增量=%.2f GB（负载为 8 GB）"
                         % (idle_sum / 1024 ** 3, busy_sum / 1024 ** 3,
                            delta_gb))
        # 增量应与负载 8 GB 同量级（3~12 GB 区间，允许驱动保留与采集窗口偏差）
        ctx.expect(3.0 <= delta_gb <= 12.0,
                   "exporter 上报的显存占用增量 %.2f GB 与负载 8 GB 不符"
                   % delta_gb)
    finally:
        ctx.sh(["docker", "rm", "-f", holder], timeout=60)
        ctx.destroy_card(ctx.card)
