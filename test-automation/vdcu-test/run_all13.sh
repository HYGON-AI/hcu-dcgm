#!/bin/bash
# 13 个测试方向串行执行，每个方向独立生成 report。
# 配置来源：同目录 config.json（默认值已内置，用户按环境修改后使用）；
# 环境变量可覆盖：IMAGE（K8s 测试镜像）、HS（hy-smi 路径）、
# CLEANUP_CARDS（结束清理的物理卡号）、NS（K8s 命名空间）。
cd "$(dirname "$0")" || exit 2
set -euo pipefail

# 从 config.json 读取默认值；解析失败（无 python3/文件缺失）时回退内置值
_cfg() {
  # 解析失败（无 python3/config.json）时返回空值，交由下方内置回退处理；
  # 不加 || true 会在 set -e 下直接退出，导致回退值不可达
  python3 -c 'import json,sys;print(json.load(open("config.json",encoding="utf-8")).get(sys.argv[1],""))' "$1" 2>/dev/null || true
}
IMAGE="${IMAGE:-$(_cfg image)}";      IMAGE="${IMAGE:-bert:v1}"
HS="${HS:-$(_cfg hy_smi)}";           HS="${HS:-/opt/hyhal/bin/hy-smi}"
CLEANUP_CARDS="${CLEANUP_CARDS:-$(_cfg cleanup_cards)}"; CLEANUP_CARDS="${CLEANUP_CARDS:-5 6 7}"
NS="${NS:-$(_cfg cluster_namespace)}"; NS="${NS:-vdcu-test}"
LOG=/tmp/all13.log
: > "$LOG"

# K8s 方向前置：镜像导入 containerd（节点 /data 高占用会 GC，执行前导入）。
# 保护：SKIP_IMPORT=1 强制跳过；/data 剩余不足 50G 时自动跳过，避免写爆数据盘
FREE_GB=$(df -BG --output=avail /data 2>/dev/null | tail -1 | tr -dc '0-9' || true)
if [ "${SKIP_IMPORT:-0}" = "1" ]; then
  echo "提示：SKIP_IMPORT=1，跳过 containerd 镜像导入" | tee -a "$LOG"
elif [ -n "$FREE_GB" ] && [ "$FREE_GB" -lt 50 ]; then
  echo "警告：/data 剩余 ${FREE_GB}G，自动跳过 containerd 镜像导入（K8s 方向可能受影响）" | tee -a "$LOG"
elif ! docker save "$IMAGE" | ctr -n k8s.io images import - >/dev/null 2>&1; then
  echo "警告：containerd 导入 $IMAGE 失败，K8s 用例可能受影响" >&2
fi
echo "containerd $IMAGE: $(crictl images 2>/dev/null | grep -c "$IMAGE")"

DIRS="fun ctr com trn inf k8s vol iso perf stb mon cmp sec"
for d in $DIRS; do
  echo "===== 方向 $d 开始 $(date +%H:%M:%S) =====" | tee -a "$LOG"
  # 单方向失败不中断后续方向（run.py 用例失败会返回 1），rc 记入日志
  if ./run_$d.sh >> "$LOG" 2>&1; then
    rc=0
  else
    rc=$?
  fi
  echo "===== 方向 $d 结束 rc=$rc $(date +%H:%M:%S) =====" | tee -a "$LOG"
done

# 测试卡清理（清理失败不影响脚本退出码，现场信息留在日志）
for c in $CLEANUP_CARDS; do
  "$HS" virtual --destroy-vdevices -d "$c" >/dev/null 2>&1 || true
done
kubectl delete pods -n "$NS" -l vdcu-test=true --ignore-not-found >/dev/null 2>&1 || true
echo "ALL_13_DONE $(date +%H:%M:%S)" | tee -a "$LOG"
