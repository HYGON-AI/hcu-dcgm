#!/bin/bash
# Copyright (c) 2026 Hygon Information Technology Co., Ltd.
# SPDX-License-Identifier: Apache-2.0

mkdir -p /etc/vdev

docker run --name hcu-dcgm -d --privileged \
  --device=/dev/kfd \
  --device=/dev/mkfd \
  --device=/dev/dri \
  -v /etc/vdev:/etc/vdev \
  -v /opt/hyhal:/opt/hyhal \
  -p 16081:16081 \
  image.sourcefind.cn:5000/dcu/admin/base/dcu-dcgm:v2.1.0 \
  /bin/bash -c "/usr/local/bin/start-dcgm.sh"
