#!/bin/bash
# 第 0 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 1 方向：FUN 测试（vDCU 虚拟化测试）
cd "$(dirname "$0")" || exit
exec python3 run.py --category FUN "$@"
