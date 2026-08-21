/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang/glog"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router"
)

// 【修改】增加 listenFlag 支持绑定 IP，可选
var (
	portFlag = flag.Int("port", 16081, "Port number for the DCGM")
	//ip 监听
	listenFlag = flag.String("listen", "", "Optional listen IP address, default 0.0.0.0")
)

func main() {
	// 解析命令行标志
	flag.Parse()

	// 确保程序退出时刷新 glog 缓存
	defer glog.Flush()

	// 初始化 DCGM 服务
	err := dcgm.Init()
	if err != nil {
		glog.Errorf("DCGM 初始化失败: %v", err)
		return
	}
	defer dcgm.ShutDown()

	log.Println("服务启动中...")

	// 初始化路由
	r := router.InitRouter()

	//处理端口和监听 IP
	port := *portFlag

	// 优先使用命令行参数 -listen
	ip := *listenFlag
	if ip == "" { // 如果命令行没传，再查环境变量
		ip = os.Getenv("HCU_DCGM_LISTEN")
	}
	// 如果命令行和环境变量都没设置，默认全零监听
	if ip == "" {
		ip = "0.0.0.0"
	}

	//拼成 ip:port
	addr := fmt.Sprintf("%s:%d", ip, port)
	glog.Infof("DCGM 服务监听地址: %s", addr)

	// 【修改】启动 Gin 服务，使用指定的 IP 和端口
	if err := r.Run(addr); err != nil {
		glog.Fatalf("启动服务失败: %v", err)
	}
}
