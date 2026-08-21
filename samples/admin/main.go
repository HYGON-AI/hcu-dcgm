/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package main

import (
	"flag"
	"log"

	"github.com/golang/glog"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

// 添加注释以描述 server 信息
//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/router/v1

// @securityDefinitions.basic	BasicAuth
func main() {
	glog.Infof("go-dcgm start ...")
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	//初始化dcgm服务
	dcgm.Init()
	log.Println("go-dcgm Init ...")
	numDevices, _ := dcgm.NumMonitorDevices()
	log.Println("HCU number of devices: ", numDevices)
	//glog.Infof("HCU number of devices: %d", numDevices)
	defer dcgm.ShutDown()
	log.Println("go-dcgm ShutDown ...")
}
