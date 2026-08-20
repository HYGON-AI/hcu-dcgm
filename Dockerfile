# Copyright (c) 2026 Hygon Information Technology Co., Ltd.
# SPDX-License-Identifier: Apache-2.0

# 1. 使用轻量级的基础镜像
FROM ubuntu:20.04

# 2. 创建工作目录
RUN mkdir -p /home/dcgm
WORKDIR /home/dcgm

# 3. 更新包列表并安装 kmod
RUN apt update && apt install -y kmod pciutils

# 4. 复制已编译好的二进制文件
COPY hcu-dcgm /home/dcgm/hcu-dcgm
RUN chmod +x /home/dcgm/hcu-dcgm

# 设置 LD_LIBRARY_PATH=/opt/hyhal/lib

# 5. 暴露端口
EXPOSE 16081

# 6. 启动命令
CMD ["/bin/sh", "-c", "/home/dcgm/hcu-dcgm -logtostderr -v=5 > /home/dcgm/dcgm.log 2>&1"]
