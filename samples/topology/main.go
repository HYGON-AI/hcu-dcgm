/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/golang/glog"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	//初始化dcgm服务
	dcgm.Init()
	defer dcgm.ShutDown()

	//硬件拓扑信息(支持json打印信息)
	//dcgm.ShowWeightTopology([]int{0, 1, 2}, true)
	//dcgm.ShowWeightTopology([]int{0, 1, 2}, false)
	//基于跳数显示硬件拓扑信息(支持json打印信息)
	//dcgm.ShowHopsTopology([]int{0, 1, 2}, false)
	//dcgm.ShowHopsTopology([]int{0, 1, 2}, true)
	//基于链接类型的硬件拓扑信息(支持json打印信息)
	//dcgm.ShowTypeTopology([]int{0, 1, 2}, true)
	//dcgm.ShowTypeTopology([]int{0, 1, 2}, false)
	//numa节点HW拓扑信息
	//dcgm.ShowNumaTopology([]int{0, 1})
	//显示硬件拓扑信息,包括权重、跳数、链接类型以及NUMA节点信息
	//dcgm.ShowHwTopology([]int{0, 1, 2})

	fmt.Println("==== HCU Interconnect Topology Demo ====")

	matrix, err := dcgm.DiscoverInterconnectTopology()
	if err != nil {
		fmt.Printf("DiscoverInterconnectTopology failed: %v\n", err)
		return
	}
	glog.V(5).Infof("DiscoverInterconnectTopology: %v", dataToJson(matrix))

	hcuCount := matrix.DeviceCount
	fmt.Printf("Total HCU count: %d\n\n", hcuCount)

	for src := 0; src < hcuCount; src++ {
		fmt.Printf("From HCU %d:\n", src)
		for dst := 0; dst < hcuCount; dst++ {
			info := matrix.Matrix[src][dst]

			fmt.Printf(
				"  -> HCU %-2d | LinkType: %-12s | Weight: %-2d | PciID: %s\n",
				info.DstDvInd,
				info.LinkType,
				info.Weight,
				info.PciID,
			)
		}
		fmt.Println()
	}

	fmt.Println("==== Demo Finished ====")

	src, dst := 0, 3
	if hcuCount > dst {
		linkInfo := matrix.Matrix[src][dst]
		fmt.Printf(
			"HCU %d -> HCU %d : LinkType=%s, Weight=%d, Hops=%d, PciID=%s\n",
			src, dst,
			linkInfo.LinkType,
			linkInfo.Weight,
			linkInfo.Hops,
			linkInfo.PciID,
		)
	}
}

func dataToJson(data any) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error serializing to JSON:", err)
	}
	return string(jsonData)
}
