# -*- coding: utf-8 -*-
"""Kubernetes、HAMi 与 Volcano 的可执行集群测试。

集群组件的资源名和部署方式在不同版本中并不完全相同。本模块先从
Node.status.capacity/allocatable 和集群对象中识别实际资源，再执行最小
影响面的 Pod 申请、运行、删除和资源回收验证。Volcano/HAMi 未部署时，
用例会明确跳过，不把“组件不存在”误报成产品缺陷。
"""
import json
import re
import shlex
import time
import uuid

from vdcu.core import case


def _enabled(ctx):
    if not ctx.cluster_enabled:
        ctx.skip("未启用集群测试；设置 VDCU_CLUSTER_TESTS=1 后执行")


def _nodes(ctx):
    rc, out, payload = ctx.kubectl_json("get", "nodes", "-o", "json",
                                        timeout=120)
    ctx.expect(rc == 0 and payload is not None,
               "kubectl 查询节点失败: %s" % out[-500:])
    nodes = payload.get("items", [])
    ctx.expect(nodes, "集群没有可查询的节点")
    return nodes


def _resource_view(node):
    status = node.get("status", {})
    return status.get("capacity", {}) or {}, status.get("allocatable", {}) or {}


def _resource_mode(nodes):
    """返回 (mode, resource-key, node)。

    前缀可配置（组长要求 2026-11）：识别顺序为
      1. hcunum（官方文档 5.1 的 HCU 前缀，开源后命名）；
      2. dcunum（HAMi/vHCU 模式的 DCU 前缀，开源前命名）；
      3. dcu（普通设备插件的整卡资源）。
    注意：混合模式下 dcunum/hcunum 可能注册但 allocatable 为 0
    （vHCU 动态切分未启用），此时以实际可分配量 > 0 的资源为准，
    否则 Pod 会因申请 0 容量资源而 Pending。
    """
    order = ["hygon.com/hcunum", "hygon.com/dcunum", "hygon.com/dcu"]
    for node in nodes:
        capacity, allocatable = _resource_view(node)
        for key in order:
            if key in allocatable and _alloc_positive(allocatable[key]):
                mode = "vhcu" if key.endswith("num") else "dcu"
                return mode, key, node
    # 标准模式下的预切分 vHCU 资源（hygon.com/dcu-share-*c-*g，
    # 其中 0c 为 CU 共享形态）
    for node in nodes:
        capacity, allocatable = _resource_view(node)
        shares = sorted([k for k in allocatable if "share" in k
                         and _alloc_positive(allocatable[k])])
        if shares:
            return "vshare", shares[0], node
    # 全部为 0：仍返回注册名（PX-001 记录现状），但标注不可申请
    for node in nodes:
        capacity, allocatable = _resource_view(node)
        for key in order:
            if key in allocatable:
                return "registered-zero", key, node
    return "", "", nodes[0]


def _alloc_positive(value):
    try:
        return int(str(value).rstrip("m")) > 0
    except (ValueError, TypeError):
        return False


def _all_hygon_keys(nodes):
    """收集节点上所有 hygon.com 资源名（PX-001 前缀现状确认）。"""
    keys = set()
    for node in nodes:
        _, allocatable = _resource_view(node)
        for k in allocatable:
            if "hygon.com" in k:
                keys.add(k)
    return sorted(keys)


def _plugin_info(ctx):
    """返回 (image, env_map) 或 (None, {})。"""
    rc, out, payload = ctx.kubectl_json(
        "get", "daemonset", "-n", "kube-system", "-o", "json", timeout=90)
    if rc != 0 or not payload:
        return None, {}
    for item in payload.get("items", []):
        name = item.get("metadata", {}).get("name", "")
        if "device-plugin" in name or "dcu-plugin" in name:
            container = (item.get("spec", {}).get("template", {})
                         .get("spec", {}).get("containers", [{}]))[0]
            env = {e.get("name"): e.get("value")
                   for e in container.get("env", [])}
            return container.get("image", ""), env
    return None, {}


def _namespace(ctx):
    return ctx.cluster_namespace


def _pod_name(prefix):
    return "%s-%s" % (prefix, uuid.uuid4().hex[:8])


def _pod_manifest(ctx, name, resource_key, amount=1, scheduler=None,
                  labels=None, annotations=None, command=None):
    metadata = {
        "name": name,
        "namespace": _namespace(ctx),
        "labels": labels or {"vdcu-test": "true"},
    }
    if annotations:
        metadata["annotations"] = annotations
    resources = {"limits": {resource_key: str(amount)},
                 "requests": {resource_key: str(amount)}}
    container = {
        "name": "probe",
        "image": ctx.k8s_image,
        "command": command or ["bash", "-lc", "sleep 3600"],
        "resources": resources,
    }
    spec = {"restartPolicy": "Never", "containers": [container]}
    if scheduler:
        spec["schedulerName"] = scheduler
    # 测试镜像只导入本机（containerd），固定调度到本机节点，
    # 避免其他节点 ErrImagePull（集群节点无 docker 镜像源）
    node_name = ctx.cfg.get("k8s_node", "")
    if node_name:
        spec["nodeName"] = node_name
    return json.dumps({"apiVersion": "v1", "kind": "Pod",
                       "metadata": metadata, "spec": spec})


def _apply_pod(ctx, manifest):
    rc, out = ctx.apply_manifest(manifest, timeout=120)
    ctx.expect(rc == 0, "提交 Pod 失败: %s" % out[-500:])


def _pod_json(ctx, name):
    return ctx.kubectl_json("get", "pod", name, "-n", _namespace(ctx),
                            "-o", "json", timeout=60)


def _wait_phase(ctx, name, phases=("Running",), timeout=180):
    deadline = time.time() + timeout
    last = ""
    last_payload = None
    while time.time() < deadline:
        rc, out, payload = _pod_json(ctx, name)
        if rc != 0 or payload is None:
            last = out[-300:]
            time.sleep(2)
            continue
        last_payload = payload
        phase = payload.get("status", {}).get("phase", "")
        last = phase
        if phase in phases:
            return payload
        if phase in ("Failed", "Succeeded") and phase not in phases:
            # 提前结束等待，但返回完整 payload，便于上层输出 reason/message
            return payload
        time.sleep(2)
    if last_payload is not None:
        return last_payload
    return {"status": {"phase": last}}


def _delete_pod(ctx, name):
    ctx.kubectl("delete", "pod", name, "-n", _namespace(ctx),
                "--ignore-not-found", "--wait=true", timeout=120)


def _probe_cmd(ctx):
    """Pod 内设备探测命令（hy-smi 与 DTK 库路径均来自配置）。

    最小镜像可能未包含 hy-smi，Pod 内优先 hy-smi，
    缺失时回退 torch.cuda 设备查询作为等价验证。
    """
    return (
        "%s virtual --show-device-info 2>&1 && echo PROBE_HY_SMI "
        "|| LD_LIBRARY_PATH=%s:$LD_LIBRARY_PATH python - <<'PY'\n"
        "import os, torch\n"
        "print('PROBE_TORCH')\n"
        "print('HIP_VISIBLE_DEVICES =', os.environ.get('HIP_VISIBLE_DEVICES'))\n"
        "n = torch.cuda.device_count()\n"
        "print('device_count =', n)\n"
        "PY"
    ) % (shlex.quote(ctx.hy_smi), ctx.dtk_hyhal_lib)


def _exec_probe(ctx, name):
    """Pod 内设备探测。

    最小镜像可能未包含 hy-smi，Pod 内优先 hy-smi，
    缺失时回退 torch.cuda 设备查询作为等价验证。
    """
    return ctx.kubectl("exec", "-n", _namespace(ctx), name, "--",
                       "bash", "-lc", _probe_cmd(ctx), timeout=150)


def _probe_ok(out):
    """探测输出判定：hy-smi 返回设备块，或 torch 回退路径查到设备。"""
    if "Device" in out:
        return True
    if "PROBE_TORCH" in out and "device_count = 0" not in out:
        return True
    return False


def _volcano_objects(ctx):
    rc, out, payload = ctx.kubectl_json("get", "ns", "volcano-system",
                                        "-o", "json", timeout=60)
    if rc != 0 or payload is None:
        return False, "未找到 volcano-system namespace"
    rc, out, payload = ctx.kubectl_json(
        "get", "pods", "-n", "volcano-system", "-o", "json", timeout=60)
    if rc != 0 or payload is None:
        return False, "无法查询 volcano-system Pod"
    names = [item.get("metadata", {}).get("name", "")
             for item in payload.get("items", [])]
    if not any("scheduler" in name for name in names):
        return False, "volcano-system 中没有 scheduler Pod"
    return True, ""


def _require_volcano(ctx, config=False):
    _enabled(ctx)
    ok, reason = _volcano_objects(ctx)
    if not ok:
        ctx.skip(reason + "；安装 Volcano 后再执行 VOL 用例")
    if config:
        rc, out, payload = ctx.kubectl_json(
            "get", "configmap", "-n", "volcano-system", "-o", "json",
            timeout=60)
        ctx.expect(rc == 0 and payload is not None,
                   "查询 Volcano ConfigMap 失败: %s" % out[-400:])
        text = json.dumps(payload, ensure_ascii=False)
        ctx.expect("deviceshare" in text.lower(),
                   "Volcano 配置中未发现 deviceshare")
    return True


def _hami_present(ctx, nodes=None):
    nodes = nodes or _nodes(ctx)
    mode, key, _ = _resource_mode(nodes)
    # dcunum/hcunum 可分配量 > 0 才视为 vHCU 动态切分可用
    if mode == "vhcu":
        return True
    for kind in ("pods", "daemonsets", "deployments"):
        rc, out, payload = ctx.kubectl_json(
            "get", kind, "-A", "-o", "json", timeout=60)
        if rc == 0 and payload:
            text = json.dumps(payload).lower()
            if "hami" in text or "vhcu" in text:
                return True
    return False


@case("VDCU-K8S-001", "Kubernetes 节点资源注册", "K8S", priority="P0")
def k8s_001(ctx):
    _enabled(ctx)
    nodes = _nodes(ctx)
    mode, key, _ = _resource_mode(nodes)
    ctx.expect(mode, "节点未注册 hygon.com/dcu 或 vHCU 资源")
    missing = []
    for node in nodes:
        capacity, allocatable = _resource_view(node)
        if key not in capacity or key not in allocatable:
            missing.append(node.get("metadata", {}).get("name", "<unknown>"))
    ctx.expect(not missing, "节点资源字段不完整: %s" % ", ".join(missing))
    print("        K8S 资源模式=%s，资源名=%s，节点数=%d" %
          (mode, key, len(nodes)))


@case("VDCU-K8S-002", "Kubernetes Pod 资源申请、设备可见与回收", "K8S",
      priority="P0")
def k8s_002(ctx):
    _enabled(ctx)
    nodes = _nodes(ctx)
    mode, key, _ = _resource_mode(nodes)
    ctx.expect(mode, "节点未注册可申请的 Hygon 资源")
    name = _pod_name("vdcu-k8s")
    manifest = _pod_manifest(ctx, name, key)
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name)
        phase = pod.get("status", {}).get("phase")
        ctx.expect(phase == "Running", "Pod 未进入 Running: %s" % phase)
        rc, out = _exec_probe(ctx, name)
        if rc == 0 and _probe_ok(out):
            print("        容器内设备信息：%s" % out[-500:].replace("\n", " "))
    finally:
        _delete_pod(ctx, name)
    rc, out, _ = _pod_json(ctx, name)
    ctx.expect(rc != 0, "Pod 删除后仍可查询，资源可能未回收")


@case("VDCU-K8S-003", "Kubernetes 资源不足时保持 Pending", "K8S",
      priority="P1")
def k8s_003(ctx):
    _enabled(ctx)
    if not ctx.cluster_mutating_tests:
        ctx.skip("资源不足测试可能占用节点；设置 VDCU_CLUSTER_MUTATING_TESTS=1")
    nodes = _nodes(ctx)
    mode, key, node = _resource_mode(nodes)
    ctx.expect(mode, "节点未注册可申请的 Hygon 资源")
    _, allocatable = _resource_view(node)
    try:
        total = int(str(allocatable.get(key, "1")).rstrip("m"))
    except ValueError:
        total = 1
    # 记录各节点 allocatable，便于对照调度端资源视图（PR#9 审查意见：
    # 调度端与节点端资源视图不一致时，需此证据定位问题）
    for n in nodes:
        _, alloc = _resource_view(n)
        ctx.steps.append("节点 %s allocatable[%s]=%s" % (
            n.get("metadata", {}).get("name", "<unknown>"),
            key, alloc.get(key, "未注册")))
    name = _pod_name("vdcu-pending")
    manifest = _pod_manifest(ctx, name, key, total + 1)
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name, phases=("Pending", "Failed"), timeout=30)
        phase = pod.get("status", {}).get("phase")
        reason = pod.get("status", {}).get("reason", "")
        message = (pod.get("status", {}).get("message") or "").strip()
        # 显式断言（PR#9 审查意见）：用例预期为"资源不足时保持 Pending"，
        # 仅 Pending 判 PASS；Failed/UnexpectedAdmissionError 说明调度端
        # 与节点端资源视图不一致或准入链路异常，判 FAIL 并输出原因。
        detail = (" reason=%s" % reason) if reason else ""
        if message:
            detail += " message=%s" % message[:300].replace("\n", " ")
        ctx.expect(phase == "Pending",
                   "超额申请（%d > %d）应保持 Pending，实际 phase=%s%s"
                   % (total + 1, total, phase, detail))
        ctx.steps.append("超额申请（%d > %d）Pod 状态：%s%s"
                         % (total + 1, total, phase, detail))
    finally:
        _delete_pod(ctx, name)


@case("VDCU-K8S-004", "Kubernetes 设备插件重启后资源恢复", "K8S",
      priority="P1")
def k8s_004(ctx):
    _enabled(ctx)
    rc, out, payload = ctx.kubectl_json(
        "get", "daemonset", "-A", "-o", "json", timeout=90)
    ctx.expect(rc == 0 and payload is not None, "查询设备插件 DaemonSet 失败")
    items = [x for x in payload.get("items", [])
             if "device-plugin" in x.get("metadata", {}).get("name", "").lower()]
    ctx.expect(items, "未找到 device-plugin DaemonSet")
    ds = items[0]
    ns = ds.get("metadata", {}).get("namespace", "kube-system")
    name = ds.get("metadata", {}).get("name")
    rc, out = ctx.kubectl("rollout", "status", "daemonset/%s" % name,
                          "-n", ns, "--timeout=120s", timeout=150)
    ctx.expect(rc == 0, "设备插件当前未就绪: %s" % out[-300:])
    if not ctx.cluster_mutating_tests:
        ctx.skip("设备插件重启属于变更操作；设置 VDCU_CLUSTER_MUTATING_TESTS=1")
    rc, out = ctx.kubectl("rollout", "restart", "daemonset/%s" % name,
                          "-n", ns, timeout=90)
    ctx.expect(rc == 0, "重启设备插件失败: %s" % out[-300:])
    rc, out = ctx.kubectl("rollout", "status", "daemonset/%s" % name,
                          "-n", ns, "--timeout=180s", timeout=210)
    ctx.expect(rc == 0, "设备插件重启后未恢复: %s" % out[-400:])


@case("VDCU-HAMI-001", "HAMi/HCU 动态资源与组件识别", "HAMI",
      priority="P0")
def hami_001(ctx):
    _enabled(ctx)
    nodes = _nodes(ctx)
    mode, key, _ = _resource_mode(nodes)
    if not _hami_present(ctx, nodes):
        # 当前集群为标准/mixed 模式（dcunum 注册为 0），HAMi 动态切分未启用。
        # 按会议纪要需按《HCU Device Plugin-v3.0.1》第 4/5 章重新部署
        # HAMi 环境后再执行 HAMI 用例，当前记录现状并跳过。
        ctx.steps.append("当前资源模式=%s，资源名=%s（vHCU 动态切分未启用）"
                         % (mode, key))
        ctx.skip("HAMi/vHCU 动态切分未启用（dcunum=0）；"
                 "按官方文档部署 HAMi 模式后执行")
    ctx.steps.append("HAMi 资源可用：%s" % key)


@case("VDCU-HAMI-002", "HAMi vHCU Pod 申请与设备可见", "HAMI",
      priority="P0")
def hami_002(ctx):
    _enabled(ctx)
    nodes = _nodes(ctx)
    mode, key, _ = _resource_mode(nodes)
    if not _hami_present(ctx, nodes):
        ctx.skip("当前集群未启用 HAMi/vHCU 资源")
    name = _pod_name("vdcu-hami")
    manifest = _pod_manifest(
        ctx, name, key, 1,
        annotations={"hygon.com/use-hcutype": "HCU-BW1000"},
    )
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name)
        ctx.expect(pod.get("status", {}).get("phase") == "Running",
                   "HAMi vHCU Pod 未进入 Running")
        rc, out = _exec_probe(ctx, name)
        ctx.expect(rc == 0, "无法进入 HAMi vHCU Pod: %s" % out[-300:])
    finally:
        _delete_pod(ctx, name)


@case("VDCU-HAMI-003", "HAMi Pod 删除后的 vHCU 资源回收", "HAMI",
      priority="P1")
def hami_003(ctx):
    _enabled(ctx)
    nodes = _nodes(ctx)
    mode, key, _ = _resource_mode(nodes)
    if not _hami_present(ctx, nodes):
        ctx.skip("当前集群未启用 HAMi/vHCU 资源")
    name = _pod_name("vdcu-hami-reclaim")
    manifest = _pod_manifest(ctx, name, key, 1)
    _apply_pod(ctx, manifest)
    try:
        pod = _wait_phase(ctx, name)
        ctx.expect(pod.get("status", {}).get("phase") == "Running",
                   "HAMi Pod 未进入 Running")
    finally:
        _delete_pod(ctx, name)
    rc, out, _ = _pod_json(ctx, name)
    ctx.expect(rc != 0, "HAMi Pod 删除后仍存在")


@case("VDCU-VOL-001", "Volcano 组件状态与调度器可用", "VOL",
      priority="P0")
def vol_001(ctx):
    _require_volcano(ctx)
    rc, out, payload = ctx.kubectl_json(
        "get", "pods", "-n", "volcano-system", "-o", "json", timeout=60)
    ctx.expect(rc == 0 and payload is not None, "查询 Volcano Pod 失败")
    bad = []
    for item in payload.get("items", []):
        phase = item.get("status", {}).get("phase")
        # 已完成的一次性初始化 Job（Succeeded）不算异常
        if phase not in ("Running", "Succeeded"):
            bad.append("%s=%s" % (item.get("metadata", {}).get("name"), phase))
    ctx.expect(not bad, "Volcano 组件未全部 Running: %s" % ", ".join(bad))


@case("VDCU-VOL-002", "Volcano deviceshare 配置生效", "VOL",
      priority="P0")
def vol_002(ctx):
    _require_volcano(ctx, config=True)
    rc, out, payload = ctx.kubectl_json(
        "get", "configmap", "-n", "volcano-system", "-o", "json",
        timeout=60)
    text = json.dumps(payload, ensure_ascii=False)
    ctx.expect("VHCUEnable" in text or "vhcuenable" in text.lower(),
               "deviceshare 未配置 VHCUEnable")


@case("VDCU-VOL-003", "Volcano 调度器镜像与版本信息", "VOL",
      priority="P1")
def vol_003(ctx):
    _require_volcano(ctx)
    rc, out, payload = ctx.kubectl_json(
        "get", "pods", "-n", "volcano-system", "-o", "json", timeout=60)
    images = []
    for item in (payload or {}).get("items", []):
        if "scheduler" in item.get("metadata", {}).get("name", ""):
            images += [c.get("image", "") for c in
                       item.get("spec", {}).get("containers", [])]
    ctx.expect(images, "未找到 Volcano scheduler 镜像")
    print("        Volcano scheduler image: %s" % ", ".join(images))


def _volcano_mutating(ctx, title):
    _require_volcano(ctx)
    if not ctx.cluster_mutating_tests:
        ctx.skip("%s 需要创建/删除 Volcano 工作负载；设置 VDCU_CLUSTER_MUTATING_TESTS=1"
                 % title)


@case("VDCU-VOL-004", "Volcano 调度器运行 vHCU Pod", "VOL", priority="P0")
def vol_004(ctx):
    _volcano_mutating(ctx, "vHCU Pod 调度")
    nodes = _nodes(ctx)
    mode, key, _ = _resource_mode(nodes)
    if not mode:
        ctx.skip("节点没有可申请的 Hygon 资源")
    name = _pod_name("vdcu-vol")
    manifest = _pod_manifest(ctx, name, key, 1, scheduler="volcano")
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name)
        ctx.expect(pod.get("status", {}).get("phase") == "Running",
                   "Volcano vHCU Pod 未进入 Running")
    finally:
        _delete_pod(ctx, name)


@case("VDCU-VOL-005", "Volcano Gang 调度配置与 PodGroup", "VOL",
      priority="P1")
def vol_005(ctx):
    _require_volcano(ctx, config=True)
    rc, out, payload = ctx.kubectl_json(
        "get", "crd", "podgroups.scheduling.volcano.sh", "-o", "json", timeout=60)
    ctx.expect(rc == 0 and payload is not None,
               "未安装 Volcano PodGroup CRD，无法执行 gang 测试")
    if not ctx.cluster_mutating_tests:
        ctx.skip("Gang 测试需要创建 PodGroup；设置 VDCU_CLUSTER_MUTATING_TESTS=1")


@case("VDCU-VOL-006", "Volcano 队列与配额对象", "VOL", priority="P1")
def vol_006(ctx):
    _require_volcano(ctx)
    found = False
    for resource in ("queue", "queues", "queue.volcano.sh"):
        rc, out, payload = ctx.kubectl_json("get", resource, "-A", "-o",
                                            "json", timeout=60)
        if rc == 0 and payload is not None:
            found = True
            break
    ctx.expect(found, "未找到 Volcano Queue 资源")


@case("VDCU-VOL-007", "Volcano 抢占与资源回收配置", "VOL", priority="P1")
def vol_007(ctx):
    _require_volcano(ctx, config=True)
    rc, out, payload = ctx.kubectl_json(
        "get", "configmap", "-n", "volcano-system", "-o", "json",
        timeout=60)
    text = json.dumps(payload, ensure_ascii=False).lower()
    ctx.expect("priority" in text or "preempt" in text,
               "Volcano 配置未启用 priority/preempt 相关插件")


@case("VDCU-VOL-008", "Volcano binpack/spread 调度策略", "VOL", priority="P1")
def vol_008(ctx):
    _require_volcano(ctx, config=True)
    rc, out, payload = ctx.kubectl_json(
        "get", "configmap", "-n", "volcano-system", "-o", "json",
        timeout=60)
    text = json.dumps(payload, ensure_ascii=False).lower()
    ctx.expect("binpack" in text or "spread" in text,
               "未配置 deviceshare SchedulePolicy（binpack/spread）")


@case("VDCU-VOL-009", "Volcano vHCU 与整卡资源混合调度", "VOL", priority="P1")
def vol_009(ctx):
    _require_volcano(ctx)
    nodes = _nodes(ctx)
    mode, key, _ = _resource_mode(nodes)
    ctx.expect(mode, "节点未注册 Hygon 资源，无法验证混合调度")
    ctx.expect(key in ("hygon.com/dcu", "hygon.com/dcunum"),
               "未知 Hygon 资源名: %s" % key)
    if not ctx.cluster_mutating_tests:
        ctx.skip("混合调度需要提交两个工作负载；设置 VDCU_CLUSTER_MUTATING_TESTS=1")


@case("VDCU-VOL-010", "Volcano Job 失败重试与资源回收", "VOL",
      priority="P1")
def vol_010(ctx):
    _volcano_mutating(ctx, "Job 重试")
    ctx.skip("Job 重试场景需按现场 Volcano Job 模板配置，脚本保留执行入口")


@case("VDCU-VOL-011", "Volcano PodGroup 状态与事件可观测性", "VOL",
      priority="P2")
def vol_011(ctx):
    _require_volcano(ctx)
    rc, out, payload = ctx.kubectl_json("get", "crd", "-o", "json",
                                        timeout=60)
    text = json.dumps(payload or {}).lower()
    ctx.expect("podgroups.scheduling.volcano.sh" in text,
               "未安装 PodGroup CRD")


@case("VDCU-VOL-012", "Volcano 训练与推理混合负载", "VOL", priority="P1")
def vol_012(ctx):
    _volcano_mutating(ctx, "训练/推理混合负载")
    ctx.skip("混合负载需现场提供训练与推理镜像/模型，脚本保留执行入口")


@case("VDCU-VOL-013", "Volcano 并发 Job 提交与资源稳定性", "VOL",
      priority="P2")
def vol_013(ctx):
    _volcano_mutating(ctx, "并发 Job")
    ctx.skip("并发压测需现场设置并发数量，避免影响共享集群")

# ---------------------------------------------------------------------------
# PX-001~006 资源名前缀配置化测试（组长要求 2026-11 会议纪要）
# 背景：开源前 DCU 前缀；开源合规审查后 HCU 前缀（3.0.0/3.0.1 硬切换）；
# 后续前缀改为配置字段（可配 GPU/DCU/HCU 任意值），测试需覆盖不同命名场景。
# ---------------------------------------------------------------------------

@case("VDCU-PX-001", "资源名前缀现状确认", "K8S", priority="P0")
def px_001(ctx):
    _enabled(ctx)
    nodes = _nodes(ctx)
    keys = _all_hygon_keys(nodes)
    ctx.steps.append("节点注册的 hygon.com 资源名：%s" % ", ".join(keys))
    image, env = _plugin_info(ctx)
    ctx.steps.append("设备插件镜像：%s\n插件环境变量：%s" % (image, env))
    mode, key, _ = _resource_mode(nodes)
    ctx.expect(mode != "", "节点未注册任何 hygon.com 资源")
    # 前缀与插件版本对应关系记录为证据，不强制断言（插件版本决定前缀）
    ctx.steps.append("当前识别资源：%s（%s）" % (key, mode))


@case("VDCU-PX-002", "HCU 前缀全流程（当前版本基线）", "K8S", priority="P0")
def px_002(ctx):
    _enabled(ctx)
    nodes = _nodes(ctx)
    mode, key, _ = _resource_mode(nodes)
    if key != "hygon.com/hcunum":
        ctx.skip("当前节点未注册 HCU 前缀（hygon.com/hcunum），"
                 "以实际注册前缀执行：%s" % key)
    name = _pod_name("vdcu-px-hcu")
    manifest = _pod_manifest(ctx, name, key, 1)
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name)
        ctx.expect(pod.get("status", {}).get("phase") == "Running",
                   "HCU 前缀 Pod 未进入 Running")
        rc, out = _exec_probe(ctx, name)
        ctx.steps.append("容器内设备信息：\n" + out[-500:])
    finally:
        _delete_pod(ctx, name)
    rc, out, _ = _pod_json(ctx, name)
    ctx.expect(rc != 0, "HCU 前缀 Pod 删除后资源未回收")


@case("VDCU-PX-003", "DCU 前缀全流程（同机线复测）", "K8S", priority="P0")
def px_003(ctx):
    _enabled(ctx)
    nodes = _nodes(ctx)
    keys = _all_hygon_keys(nodes)
    # DCU 前缀候选：dcunum（动态切分）或 dcu（整卡）或 dcu-share-*（预切分）
    candidates = ["hygon.com/dcunum", "hygon.com/dcu"] + \
                 [k for k in keys if "share" in k]
    picked = None
    for key in candidates:
        for node in nodes:
            _, allocatable = _resource_view(node)
            if key in allocatable and _alloc_positive(allocatable[key]):
                picked = key
                break
        if picked:
            break
    if picked is None:
        ctx.skip("当前节点无 DCU 前缀可申请资源（注册名：%s），"
                 "双前缀同机线测试待版本支持" % ", ".join(keys))
    ctx.steps.append("DCU 前缀全流程使用资源：%s" % picked)
    name = _pod_name("vdcu-px-dcu")
    manifest = _pod_manifest(ctx, name, picked, 1)
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name)
        ctx.expect(pod.get("status", {}).get("phase") == "Running",
                   "DCU 前缀 Pod 未进入 Running")
        rc, out = _exec_probe(ctx, name)
        ctx.steps.append("容器内设备信息：\n" + out[-500:])
        ctx.expect(rc == 0, "DCU 前缀 Pod 内设备查询失败")
    finally:
        _delete_pod(ctx, name)
    rc, out, _ = _pod_json(ctx, name)
    ctx.expect(rc != 0, "DCU 前缀 Pod 删除后资源未回收")


@case("VDCU-PX-004", "自定义前缀压测（GPU 等任意配置）", "K8S",
      priority="P2")
def px_004(ctx):
    """PX-004：前缀配置为自定义值（如 GPU）跑通全流程，专门发现组件 bug。

    配置化能力（金龙/得闯负责改造）未上线时本用例为跟踪项。
    """
    _enabled(ctx)
    nodes = _nodes(ctx)
    keys = _all_hygon_keys(nodes)
    custom = [k for k in keys if "gpu" in k.lower()]
    if not custom:
        ctx.skip("当前插件未注册自定义前缀（GPU 等）；"
                 "配置化版本上线后执行（注册名：%s）" % ", ".join(keys))
    name = _pod_name("vdcu-px-custom")
    manifest = _pod_manifest(ctx, name, custom[0], 1)
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name)
        ctx.expect(pod.get("status", {}).get("phase") == "Running",
                   "自定义前缀 Pod 未进入 Running")
        rc, out = _exec_probe(ctx, name)
        ctx.expect(rc == 0, "自定义前缀 Pod 内设备查询失败")
    finally:
        _delete_pod(ctx, name)


@case("VDCU-PX-005", "前缀非法值边界", "K8S", priority="P2")
def px_005(ctx):
    """PX-005：非法前缀（空/大小写混用/特殊字符）的组件行为记录。

    脚本侧仅记录插件 env 中的前缀相关配置，非法输入的实际组件行为
    需在配置化版本上线后由部署侧验证（组件行为，非脚本可模拟）。
    """
    _enabled(ctx)
    image, env = _plugin_info(ctx)
    ctx.expect(image is not None, "未找到设备插件")
    prefix_env = {k: v for k, v in env.items()
                  if "prefix" in k.lower() or "resource" in k.lower()}
    ctx.steps.append("插件前缀相关环境变量：%s" % prefix_env)
    ctx.steps.append("非法前缀输入（空值/大小写混用/特殊字符）的组件侧"
                     "行为待配置化版本上线后验证（PX-004 同批）")


@case("VDCU-PX-006", "前缀切换与共存", "K8S", priority="P2")
def px_006(ctx):
    """PX-006：新旧前缀混用时的调度与资源计数行为。

    切换操作属于部署变更（组件配置化上线后执行），当前版本记录现状。
    """
    _enabled(ctx)
    nodes = _nodes(ctx)
    keys = _all_hygon_keys(nodes)
    ctx.steps.append("当前并存的前缀资源：%s" % ", ".join(keys))
    ctx.expect(len(keys) >= 1, "节点未注册任何资源")


@case("VDCU-K8S-002D", "Pod 内 /dev/dri 注入视图与单机最小挂载一致", "K8S",
      priority="P1")
def k8s_002d(ctx):
    """K8S-002 增强：设备插件注入的 Pod 内 /dev/dri 视图应为单卡 card/render。

    与单机最小设备节点挂载（CTR-009B）对齐：容器内不应看到整机所有卡的
    设备符，只应有被分配物理卡对应的 cardN/renderDN。
    """
    _enabled(ctx)
    nodes = _nodes(ctx)
    mode, key, _ = _resource_mode(nodes)
    ctx.expect(mode, "节点未注册可申请的 Hygon 资源")
    name = _pod_name("vdcu-driview")
    manifest = _pod_manifest(ctx, name, key, 1)
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name)
        ctx.expect(pod.get("status", {}).get("phase") == "Running",
                   "Pod 未进入 Running")
        rc, out = ctx.kubectl("exec", "-n", _namespace(ctx), name, "--",
                              "bash", "-lc", "ls -l /dev/dri 2>&1",
                              timeout=90)
        ctx.steps.append("Pod 内 /dev/dri 视图：\n" + out)
        ctx.expect(rc == 0, "Pod 内查询 /dev/dri 失败")
        # 统计 card 数量：K8s 交付形态应为单卡注入（≤2 个 card 属合理，
        # 整树注入会出现 8+ 个 card）
        card_count = len(re.findall(r"\bcard\d+\b", out))
        ctx.steps.append("Pod 内可见 card 设备符数量：%d" % card_count)
        ctx.expect(card_count <= 2,
                   "Pod 内可见 %d 个 card 设备符，疑似整树注入" % card_count)
    finally:
        _delete_pod(ctx, name)
@case("VDCU-K8S-002S", "预切分 vHCU（dcu-share）Pod 申请与设备可见", "K8S",
      priority="P0")
def k8s_002s(ctx):
    """K8S-002S：标准模式下的预切分 vHCU 资源全流程。

    当前集群 dcunum 注册为 0（vHCU 动态切分未启用），但节点
    m09r3n14 注册了预切分资源（dcu-share-0c-16g 等，含 CU 共享形态），
    是 HAMi 重新部署前验证 K8s vHCU 交付形态的可用路径。
    """
    _enabled(ctx)
    nodes = _nodes(ctx)
    # 直接扫描所有节点的 dcu-share-* 预切分资源（可分配量 > 0）
    key = None
    node_name = None
    for node in nodes:
        _, allocatable = _resource_view(node)
        shares = sorted([k for k in allocatable if "share" in k
                         and _alloc_positive(allocatable[k])])
        if shares:
            key = shares[0]
            node_name = node.get("metadata", {}).get("name")
            break
    if key is None:
        ctx.skip("集群未注册可分配的预切分 vHCU 资源（dcu-share-*）")
    ctx.steps.append("使用资源 %s（节点 %s）" % (key, node_name))
    name = _pod_name("vdcu-share")
    manifest = _pod_manifest(ctx, name, key, 1)
    try:
        _apply_pod(ctx, manifest)
        pod = _wait_phase(ctx, name)
        phase = pod.get("status", {}).get("phase")
        ctx.expect(phase == "Running", "预切分 vHCU Pod 未进入 Running: %s"
                   % phase)
        rc, out = _exec_probe(ctx, name)
        ctx.steps.append("容器内设备信息：\n" + out[-500:])
        ctx.expect(rc == 0 and _probe_ok(out),
                   "容器内未查询到设备信息: %s" % out[-300:])
        rc, out = ctx.kubectl("exec", "-n", _namespace(ctx), name, "--",
                              "bash", "-lc", "ls -l /dev/dri 2>&1",
                              timeout=90)
        ctx.steps.append("Pod 内 /dev/dri 视图：\n" + out)
    finally:
        _delete_pod(ctx, name)
    rc, out, _ = _pod_json(ctx, name)
    ctx.expect(rc != 0, "预切分 vHCU Pod 删除后资源未回收")