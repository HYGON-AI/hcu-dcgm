#!/bin/bash
# 13 个测试方向串行执行，每个方向独立生成 report
cd /root/tz/vdcu-test || exit 2
LOG=/tmp/all13.log
: > $LOG

# K8s 方向前置：bert:v2 导入 containerd（节点 /data 高占用会 GC，执行前导入）
docker save bert:v2 | ctr -n k8s.io images import - >/dev/null 2>&1
echo "containerd bert:v2: $(crictl images 2>/dev/null | grep -c bert)"

DIRS="fun ctr com trn inf k8s vol iso perf stb mon cmp sec"
for d in $DIRS; do
  echo "===== 方向 $d 开始 $(date +%H:%M:%S) =====" | tee -a $LOG
  ./run_$d.sh >> $LOG 2>&1
  echo "===== 方向 $d 结束 rc=$? $(date +%H:%M:%S) =====" | tee -a $LOG
done

# 测试卡清理
HS=/opt/hyhal/bin/hy-smi
for c in 5 6 7; do $HS virtual --destroy-vdevices -d $c >/dev/null 2>&1; done
kubectl delete pods -n vdcu-test -l vdcu-test=true --ignore-not-found >/dev/null 2>&1
echo "ALL_13_DONE $(date +%H:%M:%S)" | tee -a $LOG
