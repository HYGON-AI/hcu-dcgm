#!/bin/bash
cd "$(dirname "$0")" || exit
exec python3 run.py "$@"
