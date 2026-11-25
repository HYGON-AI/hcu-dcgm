# -*- coding: utf-8 -*-
"""vDCU 自动化测试框架核心：用例注册、上下文、设备解析、执行与报告。"""
from __future__ import annotations

import glob
import json
import os
import re
import shlex
import subprocess
import time
from dataclasses import dataclass, field
from typing import Callable, Dict, List, Optional


# --------------------------------------------------------------------------
# 报告脱敏（PR#9 审查意见：工具侧默认剥离 IP/机器标识/内部域名，
# 避免运行报告携带内网环境信息入库）
# --------------------------------------------------------------------------
_IPV4_RE = re.compile(r"\b(?:\d{1,3}\.){3}\d{1,3}\b")
_SENSITIVE_KEY_RE = re.compile(
    r'"(machineID|bootID|systemUUID|containerID|container_id|podIP|hostIP)"'
    r'\s*:\s*"[^"]*"')


def sanitize_text(text: str, tokens=()) -> str:
    """剥离内网 IP、机器唯一标识与配置的内部主机名/域名。"""
    if not isinstance(text, str):
        return text
    text = _IPV4_RE.sub("<ip>", text)
    text = _SENSITIVE_KEY_RE.sub(r'"\1": "<redacted>"', text)
    for tok in tokens or ():
        if tok:
            text = text.replace(tok, "<masked>")
    return text


def mask_registry(image: str) -> str:
    """镜像名中的仓库地址（host[:port]）替换为 <registry>，避免泄露内部镜像源。"""
    if not isinstance(image, str) or "/" not in image:
        return image
    head, rest = image.split("/", 1)
    if ":" in head or "." in head:
        return "<registry>/" + rest
    return image


# --------------------------------------------------------------------------
# 用例定义与注册
# --------------------------------------------------------------------------
@dataclass
class Case:
    cid: str
    name: str
    category: str
    priority: str
    fn: Callable
    expected_fail: bool = False
    issue: str = ""
    skip: str = ""


class SkipCase(Exception):
    """运行时条件不满足时跳过用例。"""


REGISTRY: List[Case] = []


def case(cid: str, name: str, category: str, priority: str = "P0",
         expected_fail: bool = False, issue: str = "", skip: str = ""):
    """注册一条用例。expected_fail 用于已知缺陷（结果记 XFAIL）。"""
    def deco(fn):
        REGISTRY.append(Case(cid, name, category, priority, fn,
                             expected_fail, issue, skip))
        return fn
    return deco


@dataclass
class Result:
    cid: str
    name: str
    category: str
    priority: str
    status: str = "FAIL"
    duration: float = 0.0
    message: str = ""
    steps: List[str] = field(default_factory=list)
    issue: str = ""

    def to_dict(self):
        return {
            "id": self.cid, "name": self.name, "category": self.category,
            "priority": self.priority, "status": self.status,
            "duration_s": round(self.duration, 2), "message": self.message,
            "steps": self.steps, "issue": self.issue,
        }


# --------------------------------------------------------------------------
# 运行上下文
# --------------------------------------------------------------------------
class Context:
    def __init__(self, config: Dict[str, str], allow_destructive: bool = False,
                 verbose: bool = False):
        self.cfg = config
        self.allow_destructive = allow_destructive
        self.verbose = verbose
        self.hy_smi = config["hy_smi"]
        self.image = config["image"]
        self.inference_image = config.get("inference_image", self.image)
        self.full_card = int(config.get("full_card", config.get("card", "0")))
        self.card = int(config["card"])
        self.card2 = int(config["card2"])
        self.physical_cards = int(config.get("physical_cards", "8"))
        self.work = config["work"]
        self.vllm_model = config.get("vllm_model", "")
        self.cluster_enabled = str(config.get("cluster_tests", "0")).lower() in (
            "1", "true", "yes", "on",
        )
        self.asset_dir = os.path.join(os.path.dirname(os.path.dirname(
            os.path.abspath(__file__))), "assets")
        # 命令执行记录：每条命令与输出自动归档到 report（任务要求 3）
        self.steps: List[str] = []
        self._record = str(config.get("record_commands", "1")).lower() in (
            "1", "true", "yes", "on",
        )
        self._record_limit = int(config.get("record_output_limit", "1500"))
        # 报告脱敏开关（PR#9 审查意见）：默认剥离 IP/机器标识/内部域名，
        # 内部调试可设 sanitize=0 关闭
        self._sanitize = str(config.get("sanitize", "1")).lower() in (
            "1", "true", "yes", "on",
        )
        self._sanitize_tokens = [
            t.strip()
            for t in str(config.get("sanitize_tokens", "")).split(",")
            if t.strip()
        ]
        # ---- 路径/挂载配置（默认值均为产品标准路径，可按环境修改）----
        # hyhal 安装根目录：容器需挂载整目录（含 lib/，torch 运行库依赖）；
        # 默认由 hy_smi 路径推导（/opt/hyhal/bin/hy-smi → /opt/hyhal）
        self.hyhal_dir = config.get("hyhal_dir", "").strip()
        if not self.hyhal_dir:
            hy_smi_abs = os.path.abspath(self.hy_smi)
            parent = os.path.dirname(hy_smi_abs)
            root = os.path.dirname(parent)
            self.hyhal_dir = parent if root in ("", "/") else root
        self.vdev_conf_dir = config.get("vdev_conf_dir", "/etc/vdev").rstrip("/")
        # 容器内/挂载参数均为 Linux 路径，固定正斜杠拼接
        self.vdev_docker_dir = "%s/docker" % self.vdev_conf_dir
        self.dri_dir = config.get("dri_dir", "/dev/dri")
        self.kfd_device = config.get("kfd_device", "/dev/kfd")
        self.mkfd_device = config.get("mkfd_device", "/dev/mkfd")
        self.dtk_hyhal_lib = config.get(
            "dtk_hyhal_lib", "/opt/dtk/.hyhal/rocm_smi/lib:/opt/dtk/.hyhal/lib")
        self.assets_container_dir = config.get(
            "assets_container_dir", "/opt/assets")
        self.workspace_container_dir = config.get(
            "workspace_container_dir", "/workspace")
        self.emb_model = config.get("emb_model", "")
        # ---- 镜像仓库凭据（可选；默认空表示不需要登录）----
        self.registry_user = config.get("registry_user", "")
        self.registry_password = config.get("registry_password", "")

    # ---------------- 基础命令 ----------------
    def sh(self, cmd, timeout: int = 120, check: bool = False,
           input_text: Optional[str] = None, record: Optional[bool] = None):
        if isinstance(cmd, str):
            cmd = ["bash", "-lc", cmd]
        if self.verbose:
            print("    $ " + " ".join(cmd))
        p = subprocess.run(cmd, capture_output=True, text=True, timeout=timeout,
                           input=input_text)
        out = (p.stdout or "") + (p.stderr or "")
        if check and p.returncode != 0:
            raise AssertionError("命令失败(rc=%d): %s\n%s" % (p.returncode, cmd, out[-800:]))
        # 记录真实执行的命令与输出（写入 report JSON 的 steps 字段）
        do_record = self._record if record is None else record
        if do_record and hasattr(self, "steps") and self.steps is not None:
            limit = self._record_limit
            shown = out[-limit:] if len(out) > limit else out
            tail_mark = "\n...（输出共 %d 字符，已截断至末 %d）" % (len(out), limit) \
                if len(out) > limit else ""
            step = "$ %s\n[rc=%d] %s%s" % (
                " ".join(cmd), p.returncode, shown, tail_mark)
            if self._sanitize:
                step = sanitize_text(step, self._sanitize_tokens)
            self.steps.append(step)
        return p.returncode, out

    def hy(self, *args, timeout: int = 120):
        return self.sh([self.hy_smi] + list(args), timeout=timeout)

    # ---------------- Kubernetes ----------------
    @property
    def cluster_namespace(self):
        return self.cfg.get("cluster_namespace", "default")

    @property
    def k8s_image(self):
        return self.cfg.get("k8s_image", self.image)

    @property
    def cluster_mutating_tests(self):
        return str(self.cfg.get("cluster_mutating_tests", "0")).lower() in (
            "1", "true", "yes", "on",
        )

    def kubectl(self, *args, timeout: int = 120, input_text: Optional[str] = None):
        return self.sh(["kubectl"] + list(args), timeout=timeout,
                       input_text=input_text)

    def kubectl_json(self, *args, timeout: int = 120):
        rc, out = self.kubectl(*args, timeout=timeout)
        if rc != 0:
            return rc, out, None
        try:
            return rc, out, json.loads(out)
        except json.JSONDecodeError:
            return rc, out, None

    def apply_manifest(self, manifest: str, timeout: int = 120):
        return self.kubectl("apply", "-f", "-", timeout=timeout,
                            input_text=manifest)

    def delete_manifest(self, manifest: str, timeout: int = 120):
        return self.kubectl("delete", "-f", "-", "--ignore-not-found",
                            timeout=timeout, input_text=manifest)

    # ---------------- vDCU 管理 ----------------
    def destroy_card(self, card: int, retries: int = 12, delay: float = 1.0) -> bool:
        """销毁指定卡上的全部 vDCU。

        容器退出后驱动侧存在约 1~3 秒的引用释放延迟，此时销毁命令会返回
        "The device is exist and may be in use" 但退出码仍为 0，因此必须以
        实际剩余状态作为成功判据并重试。
        """
        for _ in range(retries):
            if not self.vdevs_on(card):
                return True
            self.hy("virtual", "--destroy-vdevices", "-d", str(card))
            time.sleep(delay)
        return not self.vdevs_on(card)

    def destroy_vdev(self, vdev: int, retries: int = 12, delay: float = 1.0) -> bool:
        for _ in range(retries):
            if not any(v["id"] == vdev for v in self.vdevs()):
                return True
            self.hy("virtual", "--destroy-vdevices", "-v", str(vdev))
            time.sleep(delay)
        return not any(v["id"] == vdev for v in self.vdevs())

    def vdevs(self) -> List[dict]:
        """查询全部 vDCU。

        实测发现（2026-11-09）：重度创建/销毁压测后宿主机
        show-vdevice-info 查询路径可能损坏（generic error），而
        create/destroy/容器侧视图均正常（ISSUE-K8S-002 类运行时损坏）。
        此时回退解析 vdev_conf_dir/*.conf 获取 vDCU 状态。
        """
        rc, out = self.hy("virtual", "--show-vdevice-info")
        if "Virtual Device" in out:
            return parse_vdev_info(out)
        fallback = self._vdevs_from_conf()
        if (rc != 0 or "generic error" in out.lower()) \
                and not getattr(self, "_vdev_fallback_noted", False):
            self._vdev_fallback_noted = True
            self.steps.append(
                "注：show-vdevice-info 查询异常（%s），已回退解析 "
                "%s/*.conf 获取 vDCU 状态（ISSUE-K8S-002 类运行时损坏）"
                % (out[-120:].strip().replace("\n", " ") or "rc=%d" % rc,
                   self.vdev_conf_dir))
        return fallback

    def _vdevs_from_conf(self) -> List[dict]:
        """从 vdev_conf_dir/vdev*.conf 解析 vDCU 状态（查询路径损坏时的回退源）。"""
        out = []
        for path in sorted(glob.glob("%s/vdev*.conf" % self.vdev_conf_dir)):
            if not os.path.isfile(path):
                continue
            try:
                with open(path) as f:
                    conf = f.read()
            except OSError:
                continue
            if not re.search(r"^enable:\s*1", conf, re.M):
                continue

            def _num(key, default=0):
                m = re.search(r"^%s:\s*(\d+)" % key, conf, re.M)
                return int(m.group(1)) if m else default

            m = re.search(r"^vdev_id:\s*(\d+)", conf, re.M)
            out.append({
                "id": int(m.group(1)) if m else 0,
                "phys": _num("device_id"),
                "cu": _num("cu_count"),
                "mem": _num("mem") * 1024 ** 2,
                "container_id": "",
                "usage_mem": 0,
                "busy": 0,
                "conf_fallback": True,
            })
        return out

    def vdevs_on(self, card: int) -> List[dict]:
        return [v for v in self.vdevs() if v["phys"] == card]

    def devices(self) -> List[dict]:
        rc, out = self.hy("virtual", "--show-device-info")
        return parse_device_info(out)

    def remaining(self) -> Dict[int, dict]:
        rc, out = self.hy("virtual", "--show-remaining-info")
        return parse_remaining(out)

    def conf_path(self, vdev: int) -> str:
        return "%s/vdev%d.conf" % (self.vdev_conf_dir, vdev)

    # ---------------- 容器 ----------------
    def dri_for_phys(self, phys: int) -> tuple:
        """物理卡 phys → (card, renderD) 设备符路径。

        映射不写死：优先按 vDCU 配置文件中的 PciBusId 与
        /dev/dri/by-path 符号链接校准；PCI 枚举顺序与物理卡号并非
        线性对应（本机实测 card(P+1) 公式不成立）。校准失败时回退
        公式 card=phys+1 / renderD=128+phys 并在返回值中附带来源。
        """
        cache = getattr(self, "_dri_cache", None)
        if cache is None:
            cache = self._build_dri_map()
            self._dri_cache = cache
        if phys in cache:
            card, render, src = cache[phys]
            return (card, render, src)
        # 回退公式
        return (os.path.join(self.dri_dir, "card%d" % (phys + 1)),
                os.path.join(self.dri_dir, "renderD%d" % (128 + phys)),
                "formula")

    def _build_dri_map(self) -> Dict[int, tuple]:
        """构建 phys→(card,render,来源) 映射。

        两个来源互补：
        1) vdev 配置中的 PciBusId（有 vDCU 时最直接）；
        2) KFD 拓扑节点（sysfs）：node 8+i 对应 phys i，
           properties.location_id 为 PCI BDF（十进制），转十六进制后
           经 /dev/dri/by-path 解析 card/render（整卡形态无 vDCU 时使用）。
        """
        mapping = {}
        # 来源 1：vdev 配置
        vlist = self.vdevs()
        for v in vlist:
            vid = v["id"]
            try:
                with open(self.conf_path(vid)) as f:
                    conf = f.read()
            except OSError:
                continue
            m = re.search(r"PciBusId:\s*([0-9a-fA-F:.]+)", conf)
            if not m:
                continue
            pci = m.group(1)
            card = self._bypath_to_dev(pci, "card")
            render = self._bypath_to_dev(pci, "render")
            if card and render:
                mapping[v["phys"]] = (card, render, "vdev-conf:%s" % pci)
        # 来源 2：KFD 拓扑
        for i in range(8):
            if i in mapping:
                continue
            loc_file = ("/sys/class/kfd/kfd/topology/nodes/%d/properties"
                        % (8 + i))
            try:
                with open(loc_file) as f:
                    props = f.read()
            except OSError:
                continue
            m = re.search(r"location_id\s+(\d+)", props)
            if not m:
                continue
            loc = int(m.group(1))
            bus = (loc >> 8) & 0xFF
            dev = (loc >> 3) & 0x1F
            fn = loc & 0x7
            pci = "0000:%02x:%02x.%d" % (bus, dev, fn)
            card = self._bypath_to_dev(pci, "card")
            render = self._bypath_to_dev(pci, "render")
            if card and render:
                mapping[i] = (card, render, "kfd-topo:%s" % pci)
        return mapping

    def _bypath_to_dev(self, pci: str, kind: str) -> Optional[str]:
        """dri_dir/by-path/pci-<pci>-<kind> → 实际设备路径。"""
        link = os.path.join(self.dri_dir, "by-path",
                            "pci-%s-%s" % (pci, kind))
        if os.path.exists(link):
            target = os.readlink(link)
            name = os.path.basename(target)
            return os.path.join(self.dri_dir, name)
        return None

    def dri_paths_for(self, vdevs: List[int]) -> List[str]:
        """给定 vDCU 列表，返回其所在物理卡的最小 card/render 设备符集合。

        vDCU 全部位于同一物理卡时只挂一对；跨卡时按物理卡去重后逐对挂载。
        """
        vlist = self.vdevs()
        phys = sorted({v["phys"] for v in vlist if v["id"] in vdevs})
        if not phys:
            # 未挂 vDCU：整卡形态，按 full_card 所在物理卡透传
            phys = [self.full_card]
        paths = []
        for p in phys:
            card, render, src = self.dri_for_phys(p)
            paths += ["--device=%s" % card, "--device=%s" % render]
        return paths

    def docker(self, vdevs: List[int], script: str, env: Optional[Dict] = None,
               name: str = "", detach: bool = False, shm: str = "16G",
               timeout: int = 900, extra_mounts: Optional[List[str]] = None,
               mount_assets: bool = False, image: Optional[str] = None,
               dri_tree: bool = False, privileged: bool = False,
               kfd: bool = True, mkfd: bool = True,
               mount_vdev_conf: bool = True):
        """启动测试容器。

        - dri_tree=False（默认）：按 K8s 交付形态，仅透传 vDCU 所在物理卡的
          card/render 设备符（最小设备节点挂载）；
        - dri_tree=True：整树 --device=/dev/dri（对照形态）；
        - kfd/mkfd 为透传必挂节点，负向用例可单独关闭；
        - mount_vdev_conf=False：不挂 vDCU 配置（负向用例）。
        """
        args = ["docker", "run"]
        if detach:
            args += ["-d", "--name", name]
        else:
            args += ["--rm"]
        if privileged:
            args += ["--privileged"]
        args += ["--shm-size", shm]
        if kfd:
            args += ["--device=%s" % self.kfd_device]
        if mkfd:
            args += ["--device=%s" % self.mkfd_device]
        if dri_tree:
            args += ["--device=%s" % self.dri_dir]
        else:
            args += self.dri_paths_for(vdevs)
        # 所有挂载点均来自配置：hyhal 安装根目录（必须挂整目录而不是仅 bin/——
        # 容器内 torch 等运行时依赖 <hyhal>/lib 下的 librocm_smi64.so.2 等库，
        # 只挂 bin/ 会导致容器内 torch 导入失败）、work 工作目录、
        # vdev 配置目录、assets 目录；容器侧路径同样由配置指定
        hyhal_dir = self.hyhal_dir
        args += ["-v", "%s:%s:ro" % (hyhal_dir, hyhal_dir),
                 "-v", "%s:%s:ro" % (self.work, self.workspace_container_dir)]
        if mount_vdev_conf:
            for v in vdevs:
                args += ["-v", "%s:%s/vdev%d.conf:ro"
                         % (self.conf_path(v), self.vdev_docker_dir, v)]
        if mount_assets:
            args += ["-v", "%s:%s:ro"
                     % (self.assets_dir(), self.assets_container_dir)]
        for m in (extra_mounts or []):
            args += ["-v", m]
        for k, v in (env or {}).items():
            args += ["-e", "%s=%s" % (k, v)]
        args += [image or self.image, "bash", "-lc", script]
        return self.sh(args, timeout=timeout)

    def assets_dir(self) -> str:
        return self.asset_dir

    def assets_path(self, script: str) -> str:
        """容器内 assets 目录（来自配置 assets_container_dir）下的脚本路径。

        容器内为 Linux 环境，始终使用正斜杠拼接。
        """
        return "%s/%s" % (self.assets_container_dir.rstrip("/"), script)

    def assets_cmd(self, args: str) -> str:
        """容器内 assets 脚本调用命令：python <assets_container_dir>/<args>。"""
        return "python %s" % self.assets_path(args)

    # ---------------- 断言辅助 ----------------
    @staticmethod
    def expect(cond, msg):
        if not cond:
            raise AssertionError(msg)

    @staticmethod
    def skip(reason):
        raise SkipCase(reason)


# --------------------------------------------------------------------------
# 输出解析
# --------------------------------------------------------------------------
_BLOCK_RE = re.compile(r"^(Virtual Device|Device)\s+(\d+):\s*$", re.M)


def _blocks(text: str) -> List[tuple]:
    """把 hy-smi 输出切成 (kind, id, body) 块。"""
    matches = list(_BLOCK_RE.finditer(text))
    result = []
    for i, m in enumerate(matches):
        start = m.end()
        end = matches[i + 1].start() if i + 1 < len(matches) else len(text)
        result.append((m.group(1), int(m.group(2)), text[start:end]))
    return result


def _field(body: str, key: str) -> Optional[int]:
    # hy-smi has used both singular and plural labels across driver releases;
    # accept whitespace and an optional unit suffix while retaining bytes.
    m = re.search(r"%s\s*:\s*([\d]+)" % re.escape(key), body, re.I)
    return int(m.group(1)) if m else None


def _field_any(body: str, keys) -> Optional[int]:
    for key in keys:
        value = _field(body, key)
        if value is not None:
            return value
    return None


def _field_str(body: str, keys) -> Optional[str]:
    """字符串字段：container id 等非数值信息。"""
    for key in keys:
        m = re.search(r"%s\s*:\s*(\S+)" % re.escape(key), body, re.I)
        if m:
            return m.group(1).strip()
    return None


def parse_vdev_info(text: str) -> List[dict]:
    out = []
    for kind, vid, body in _blocks(text):
        if kind != "Virtual Device":
            continue
        out.append({
            "id": vid,
            "phys": _field(body, "Physical Device"),
            "cu": _field(body, "Compute units"),
            "mem": _field(body, "Global memory"),
            # 占用观测字段（MX-003 排查对象）：驱动 6.4.2 的输出可能不含
            # container id 字段，None 即"字段不存在"证据，不视为解析错误。
            "container_id": _field_str(body, ("Container id", "Container ID",
                                              "ContainerID")),
            "usage_mem": _field(body, "Usage memory"),
            "busy": _field(body, "Dev busy percent"),
        })
    return out


def parse_device_info(text: str) -> List[dict]:
    out = []
    for kind, vid, body in _blocks(text):
        if kind != "Device":
            continue
        out.append({
            "id": vid,
            "phys": _field(body, "Physical Device"),
            "cu": _field(body, "Compute units"),
            "mem": _field(body, "Global memory"),
            "used": _field(body, "Usage memory"),
        })
    return out


def parse_remaining(text: str) -> Dict[int, dict]:
    """返回 {physical_device: {cu, mem}}"""
    out = {}
    for kind, did, body in _blocks(text):
        if kind != "Device":
            continue
        out[did] = {
            "cu": _field_any(body, ("Remaining compute units", "Remaining CU")),
            "mem": _field_any(body, (
                "Remaining memories", "Remaining memory", "Remaining global memory",
                "Remaining global memories",
            )),
        }
    return out


# --------------------------------------------------------------------------
# 执行器与报告
# --------------------------------------------------------------------------
def run_cases(cases: List[Case], ctx: Context) -> List[Result]:
    results = []
    for c in cases:
        r = Result(c.cid, c.name, c.category, c.priority, issue=c.issue)
        t0 = time.time()
        print("[%-10s] %-11s %s" % (c.category, c.cid, c.name), flush=True)

        if c.skip:
            r.status = "SKIP"
            r.message = c.skip
            r.duration = 0.0
            print("             -> SKIP (%s)" % c.skip)
            results.append(r)
            continue

        case_ctx = Context(ctx.cfg, ctx.allow_destructive, ctx.verbose)
        try:
            c.fn(case_ctx)
            r.steps = list(case_ctx.steps)
            r.status = "XPASS" if c.expected_fail else "PASS"
            r.message = "expected failure not reproduced" if c.expected_fail else "ok"
        except AssertionError as e:
            r.steps = list(getattr(case_ctx, "steps", []) or [])
            if c.expected_fail:
                r.status = "XFAIL"
                r.message = str(e)
            else:
                r.status = "FAIL"
                r.message = str(e)
        except SkipCase as e:
            r.status = "SKIP"
            r.message = str(e)
        except subprocess.TimeoutExpired as e:
            r.steps = list(getattr(case_ctx, "steps", []) or [])
            r.status = "XFAIL" if c.expected_fail else "FAIL"
            r.message = "timeout after %ss" % getattr(e, "timeout", "?")
        except Exception as e:  # noqa: BLE001
            r.steps = list(getattr(case_ctx, "steps", []) or [])
            r.status = "XFAIL" if c.expected_fail else "FAIL"
            r.message = "%s: %s" % (type(e).__name__, e)
        finally:
            r.duration = time.time() - t0
        print("             -> %-5s (%.1fs) %s" % (r.status, r.duration,
                                                   "" if r.status in ("PASS", "XFAIL") else r.message[:160]))
        results.append(r)
    return results


def summarize(results: List[Result]) -> Dict[str, int]:
    s = {"PASS": 0, "FAIL": 0, "XFAIL": 0, "XPASS": 0, "SKIP": 0}
    for r in results:
        s[r.status] = s.get(r.status, 0) + 1
    return s


def _sanitize_env(env: Dict[str, str], tokens=()) -> Dict[str, str]:
    """报告 environment 脱敏：镜像仓库地址、主机工作目录与节点名等。"""
    out = {}
    for k, v in env.items():
        if not isinstance(v, str):
            out[k] = v
            continue
        if k in ("image", "inference_image", "k8s_image"):
            v = mask_registry(v)
        elif k == "work":
            v = "<local-path>"
        elif k == "k8s_node" and v:
            v = "<node>"
        out[k] = sanitize_text(v, tokens)
    return out


def write_reports(results: List[Result], outdir: str, env: Dict[str, str],
                  sanitize: bool = True,
                  sanitize_tokens=()) -> List[str]:
    os.makedirs(outdir, exist_ok=True)
    ts = time.strftime("%Y%m%d-%H%M%S")

    if sanitize:
        env = _sanitize_env(env, sanitize_tokens)

    data = {
        "generated_at": time.strftime("%Y-%m-%d %H:%M:%S"),
        "environment": env,
        "summary": summarize(results),
        "results": [r.to_dict() for r in results],
    }
    jpath = os.path.join(outdir, "report-%s.json" % ts)
    with open(jpath, "w", encoding="utf-8") as f:
        json.dump(data, f, ensure_ascii=False, indent=2)

    s = data["summary"]
    total = sum(s.values())
    lines = [
        "# vDCU 自动化测试报告",
        "",
        "- 生成时间：%s" % data["generated_at"],
        "- 执行环境：%s" % json.dumps(env, ensure_ascii=False),
        "- 用例总数：%d　通过 %d　失败 %d　预期失败 %d　跳过 %d" % (
            total, s["PASS"], s["FAIL"], s["XFAIL"], s["SKIP"]),
        "",
        "| 用例编号 | 用例名称 | 类别 | 优先级 | 状态 | 耗时(s) | 说明 |",
        "| --- | --- | --- | --- | --- | --- | --- |",
    ]
    for r in results:
        msg = r.message.replace("|", "/").replace("\n", " ")[:120]
        lines.append("| %s | %s | %s | %s | **%s** | %.1f | %s |" % (
            r.cid, r.name, r.category, r.priority, r.status, r.duration, msg))
    mpath = os.path.join(outdir, "report-%s.md" % ts)
    with open(mpath, "w", encoding="utf-8") as f:
        f.write("\n".join(lines) + "\n")

    return [jpath, mpath]
