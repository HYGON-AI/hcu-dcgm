#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""vDCU 虚拟化自动化测试入口。

用法：
    python3 run.py --list                     # 列出全部用例
    python3 run.py                            # 执行全部用例
    python3 run.py --category FUN,CTR         # 按类别执行
    python3 run.py --id VDCU-FUN-002          # 执行单条用例
    python3 run.py --allow-destructive        # 允许执行全局销毁类用例
"""
from __future__ import annotations

import argparse
import json
import os
import shutil
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)

import cases.ctr       # noqa: E402,F401
import cases.cluster   # noqa: E402,F401
import cases.dcgm      # noqa: E402,F401
import cases.fun       # noqa: E402,F401
import cases.sec       # noqa: E402,F401
import cases.stb       # noqa: E402,F401
import cases.system    # noqa: E402,F401
import cases.workload  # noqa: E402,F401
import cases.v12_supplement  # noqa: E402,F401
from vdcu.core import REGISTRY, Context, parse_device_info, run_cases, summarize, write_reports  # noqa: E402

DEFAULT_CONFIG = {
    # 以下为 config.json 缺失时的最小回退默认值；
    # 完整默认值（含注释）见同目录 config.json，用户按环境修改 config.json
    # 或使用 VDCU_<键名大写> 环境变量覆盖。
    "hy_smi": "/opt/hyhal/bin/hy-smi",
    # hyhal 安装根目录（容器挂载整目录，含 torch 运行库）；空=由 hy_smi 推导
    "hyhal_dir": "",
    "image": "bert:v1",
    # 内部镜像仓库地址不进代码，统一由 config.json / VDCU_ 环境变量提供
    "inference_image": "",
    "legacy_inference_image": "",
    "card": "5",
    "card2": "6",
    "full_card": "7",
    "physical_cards": "8",
    "work": "/root/tz",
    "vllm_model": "",
    "emb_model": "",
    "cluster_tests": "0",
    "cluster_mutating_tests": "0",
    "cluster_namespace": "vdcu-test",
    "k8s_image": "",
    "k8s_node": "",
    "resource_prefix": "",
    "card_offset": "1",
    "render_offset": "0",
    "vdev_conf_dir": "/etc/vdev",
    "dri_dir": "/dev/dri",
    "kfd_device": "/dev/kfd",
    "mkfd_device": "/dev/mkfd",
    "dtk_hyhal_lib": "/opt/dtk/.hyhal/rocm_smi/lib:/opt/dtk/.hyhal/lib",
    "assets_container_dir": "/opt/assets",
    "workspace_container_dir": "/workspace",
    "registry": "",
    "registry_user": "",
    "registry_password": "",
    "cleanup_cards": "5 6 7",
    "record_commands": "1",
    "record_output_limit": "1500",
    "sanitize": "1",
    "sanitize_tokens": "",
}


def load_config():
    cfg = dict(DEFAULT_CONFIG)
    path = os.path.join(HERE, "config.json")
    if os.path.exists(path):
        with open(path, encoding="utf-8") as f:
            cfg.update(json.load(f))
    for k in cfg:
        env_key = "VDCU_" + k.upper()
        if env_key in os.environ:
            cfg[k] = os.environ[env_key]
    return cfg


def _image_exists(image):
    try:
        p = subprocess.run(["docker", "image", "inspect", image],
                           stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
                           timeout=60)
        return p.returncode == 0
    except (OSError, subprocess.TimeoutExpired):
        return False


def _registry_login(cfg):
    """私有镜像仓库登录（可选）。

    registry_user/registry_password 来自 config.json（默认空=不登录）。
    密码仅通过 stdin 传给 docker login，不会出现在命令行与报告中；
    登录失败只告警不阻断（镜像可能已在本地，由预检继续判断）。
    """
    user = cfg.get("registry_user", "")
    password = cfg.get("registry_password", "")
    if not user:
        return
    host = cfg.get("registry", "") or ""
    if not host:
        image = cfg.get("inference_image", "")
        if "/" in image:
            head = image.split("/")[0]
            if ":" in head or "." in head:
                host = head
    if not host:
        print("警告：已配置 registry_user，但 registry 为空且无法从 "
              "inference_image 解析仓库地址，跳过 docker login")
        return
    try:
        p = subprocess.run(["docker", "login", host, "-u", user,
                            "--password-stdin"],
                           input=password, capture_output=True, text=True,
                           timeout=120)
    except (OSError, subprocess.TimeoutExpired) as e:
        print("警告：docker login %s 执行失败：%s" % (host, e))
        return
    if p.returncode != 0:
        print("警告：docker login %s 失败：%s"
              % (host, (p.stderr or p.stdout)[-200:].strip()))
    else:
        print("docker login %s 成功（用户 %s）" % (host, user))


def preflight(cfg, cases=None):
    cases = cases or []
    needs_local = any(c.category.upper() not in ("K8S", "VOL", "HAMI")
                      for c in cases)
    needs_k8s = any(c.category.upper() in ("K8S", "VOL", "HAMI")
                    for c in cases)
    problems = []
    if needs_local and not os.path.exists(cfg["hy_smi"]):
        problems.append("hy-smi 不存在: %s" % cfg["hy_smi"])
    if needs_local and not shutil.which("docker"):
        problems.append("docker 不可用")
    if needs_local and not os.path.isdir(cfg["work"]):
        try:
            os.makedirs(cfg["work"], exist_ok=True)
            print("工作目录 %s 不存在，已自动创建（模型/数据文件请自行放入）"
                  % cfg["work"])
        except OSError as e:
            problems.append("工作目录不存在且无法创建: %s (%s)"
                            % (cfg["work"], e))
    if needs_k8s and not shutil.which("kubectl"):
        problems.append("kubectl 不可用")
    if problems:
        print("环境预检失败：")
        for p in problems:
            print("  - " + p)
        return False
    if needs_local and not _image_exists(cfg["image"]):
        print("环境预检失败：镜像 %s 不存在" % cfg["image"])
        return False
    if any(c.cid == "VDCU-INF-002" and cfg.get("vllm_model") for c in cases):
        if not cfg.get("inference_image"):
            print("环境预检失败：未配置 inference_image")
            return False
        if not _image_exists(cfg["inference_image"]):
            print("环境预检失败：推理镜像 %s 不存在" % cfg["inference_image"])
            return False

    # 私有仓库登录（可选，失败不阻断）
    if needs_local and cfg.get("registry_user"):
        _registry_login(cfg)

    if not needs_local:
        return True

    try:
        p = subprocess.run([cfg["hy_smi"], "virtual", "--show-device-info"],
                           capture_output=True, text=True, timeout=60)
        physical = parse_device_info((p.stdout or "") + (p.stderr or ""))
        ids = {d["id"] for d in physical if d.get("id") is not None}
        required = {int(cfg["card"]), int(cfg["card2"]),
                    int(cfg.get("full_card", cfg["card"]))}
        if p.returncode != 0 or not required.issubset(ids):
            print("环境预检失败：测试卡不在 hy-smi 列表中 %s (found=%s)" %
                  (sorted(required), sorted(ids)))
            return False
    except (OSError, subprocess.TimeoutExpired) as e:
        print("环境预检失败：无法查询物理设备: %s" % e)
        return False
    return True


def select(args):
    items = list(REGISTRY)
    if args.category:
        wanted = {c.strip().upper() for c in args.category.split(",") if c.strip()}
        items = [c for c in items if c.category.upper() in wanted]
    if args.id:
        wanted = {i.strip() for i in args.id.split(",") if i.strip()}
        items = [c for c in items if c.cid in wanted]
    if args.priority:
        wanted = {p.strip().upper() for p in args.priority.split(",") if p.strip()}
        items = [c for c in items if c.priority.upper() in wanted]
    return items


def main():
    ap = argparse.ArgumentParser(description="vDCU 虚拟化自动化测试")
    ap.add_argument("--list", action="store_true", help="列出用例")
    ap.add_argument("--category", help="按类别过滤，如 FUN,CTR")
    ap.add_argument("--id", help="按用例编号过滤，逗号分隔")
    ap.add_argument("--priority", help="按优先级过滤，如 P0,P1")
    ap.add_argument("--allow-destructive", action="store_true",
                    help="允许执行会清除全局状态的用例")
    ap.add_argument("--verbose", action="store_true", help="打印执行的命令")
    ap.add_argument("--no-report", action="store_true", help="不生成报告文件")
    ap.add_argument("--cleanup", action="store_true", help="结束后清理测试卡上的 vDCU")
    args = ap.parse_args()

    cfg = load_config()
    cases = select(args)

    if args.list:
        print("%-16s %-6s %-6s %s" % ("用例编号", "类别", "优先级", "名称"))
        print("-" * 90)
        for c in cases:
            tail = "  [预期失败 %s]" % c.issue if c.expected_fail else ""
            skip = "  [跳过] %s" % c.skip if c.skip else ""
            print("%-16s %-6s %-6s %s%s%s" %
                  (c.cid, c.category, c.priority, c.name, tail, skip))
        print("\n共 %d 条用例" % len(cases))
        return 0

    if not cases:
        print("没有匹配的用例")
        return 2
    if not preflight(cfg, cases):
        return 2

    print("=" * 78)
    print(" vDCU 虚拟化自动化测试")
    print(" 镜像 %s | 测试卡 %s/%s | 整卡基线 %s | 共 %d 条用例" %
          (cfg["image"], cfg["card"], cfg["card2"], cfg.get("full_card", cfg["card"]), len(cases)))
    print("=" * 78)

    ctx = Context(cfg, allow_destructive=args.allow_destructive,
                  verbose=args.verbose)
    results = run_cases(cases, ctx)

    if args.cleanup:
        ctx.destroy_card(int(cfg["card"]))
        ctx.destroy_card(int(cfg["card2"]))
        full_card = int(cfg.get("full_card", cfg["card"]))
        if full_card not in {int(cfg["card"]), int(cfg["card2"])}:
            ctx.destroy_card(full_card)

    s = summarize(results)
    print("-" * 78)
    print(" 汇总：PASS %d | FAIL %d | XFAIL %d | XPASS %d | SKIP %d （共 %d）" % (
        s["PASS"], s["FAIL"], s["XFAIL"], s["XPASS"], s["SKIP"], len(results)))

    if not args.no_report:
        sanitize = str(cfg.get("sanitize", "1")).lower() in (
            "1", "true", "yes", "on",
        )
        sanitize_tokens = [t.strip() for t in
                           str(cfg.get("sanitize_tokens", "")).split(",")
                           if t.strip()]
        # 报告 environment 不携带注释类键（_comments 等）与密码
        env = {k: v for k, v in cfg.items()
               if not k.startswith("_") and k != "registry_password"}
        paths = write_reports(results, os.path.join(HERE, "reports"), env,
                              sanitize=sanitize, sanitize_tokens=sanitize_tokens)
        for p in paths:
            print(" 报告：" + p)

    failed = [r for r in results if r.status in ("FAIL", "XPASS")]
    if failed:
        print("\n未通过用例：")
        for r in failed:
            print("  - %s %s：%s" % (r.cid, r.name, r.message[:200]))
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
