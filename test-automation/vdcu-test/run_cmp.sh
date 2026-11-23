#!/bin/bash
# 第 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 0 -1 1 方向：CMP 测试（vDCU 虚拟化测试）
cd "$(dirname "$0")" || exit
exec python3 run.py --category CMP "$@"
