#!/bin/bash
# 第 -1 -1 -1 -1 -1 -1 -1 -1 0 -1 -1 -1 -1 1 方向：CTR 测试（vDCU 虚拟化测试）
cd "$(dirname "$0")" || exit
exec python3 run.py --category CTR "$@"
