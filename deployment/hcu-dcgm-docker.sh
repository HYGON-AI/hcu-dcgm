#!/bin/bash
# Copyright (c) 2026 Hygon Information Technology Co., Ltd.
# SPDX-License-Identifier: Apache-2.0

mkdir -p /etc/vdev

docker run --name hcu-dcgm -d --privileged \
  --device=/dev/kfd \
  --device=/dev/mkfd \
  --device=/dev/dri \
  -v /etc/vdev:/etc/vdev \
  -v /etc/hostname:/etc/hostname \
  -v /etc/vdev:/etc/vdev \
  -v /opt/hyhal:/opt/hyhal \
  -v /home/chengdm/config:/home/dcgm/config \
  -p 16081:16081 \
  -e LD_LIBRARY_PATH="/opt/hyhal/lib" \
  hcu-dcgm:v2.0.0
