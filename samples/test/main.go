/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/golang/glog"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := dcgm.Init(); err != nil {
		glog.Errorf("init dcgm failed: %v", err)
		return
	}
	defer func() {
		if err := dcgm.ShutDown(); err != nil {
			glog.Errorf("shutdown dcgm failed: %v", err)
		}
	}()

	testDiscoverInterconnectTopology()

	num, err := dcgm.NumMonitorDevices()
	if err != nil {
		glog.Errorf("get monitor number failed: %v", err)
		return
	}
	glog.Infof("card num: %d", num)

	processes, err := dcgm.ProcessHCUInfo()
	if err != nil {
		// ProcessHCUInfo 可能返回部分成功结果；先记录错误，继续打印已有结果。
		glog.Errorf("get process info failed: %v", err)
	}

	for _, process := range processes {
		glog.Infof("process id: %d, gpu index: %v", process.ProcessID, process.MinorNumbers)
	}
}

func testDiscoverInterconnectTopology() {
	fmt.Println("===== DiscoverInterconnectTopology demo =====")

	matrix, err := dcgm.DiscoverInterconnectTopology()
	if err != nil {
		fmt.Printf("DiscoverInterconnectTopology ERROR: %v\n", err)
		return
	}
	fmt.Printf("HCU count: %d\n", matrix.DeviceCount)
	if matrix.DeviceCount == 0 {
		fmt.Println("No HCU devices found.")
		return
	}

	// 行为源设备，列为目标设备；单元格直接展示接口返回的链路类型和权重。
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprint(w, "src/dst")
	for dst := 0; dst < matrix.DeviceCount; dst++ {
		fmt.Fprintf(w, "\tHCU%d", dst)
	}
	fmt.Fprintln(w)
	for src := 0; src < matrix.DeviceCount; src++ {
		fmt.Fprintf(w, "HCU%d", src)
		for dst := 0; dst < matrix.DeviceCount; dst++ {
			if src == dst {
				fmt.Fprint(w, "\tX")
				continue
			}
			link := matrix.Matrix[src][dst]
			fmt.Fprintf(w, "\t%s(%d)", link.LinkType, link.Weight)
		}
		fmt.Fprintln(w)
	}
	if err := w.Flush(); err != nil {
		glog.Errorf("print interconnect topology matrix failed: %v", err)
	}
	fmt.Println("Rows: source HCU; columns: destination HCU; cells: LinkType(Weight); X: self.")
}

// 以下为原 samples/test 中与本次进程信息验证无关的设备指标测试，暂时保留并注释。
//
// import "C"
// import (
// 	"fmt"
//
// 	"github.com/goccy/go-json"
//
// 	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
// )
//
// func oldDeviceMetricsDemo() {
// 	fmt.Println("===== DevGpuMetricsInfo demo =====")
// 	fmt.Println("Testing devices: 0, 1")
//
// 	if err := dcgm.Init(); err != nil {
// 		fmt.Printf("dcgm.Init ERROR: %v\n", err)
// 		return
// 	}
// 	defer func() {
// 		if err := dcgm.ShutDown(); err != nil {
// 			fmt.Printf("dcgm.ShutDown ERROR: %v\n", err)
// 		}
// 	}()
//
// 	for _, dvInd := range []int{0, 1} {
// 		testDevThrottleFlags(dvInd)
// 		testDevGpuMetricsInfo(dvInd)
// 	}
//
// 	fmt.Println("===== DCGM demo finished =====")
// }
//
// func testDevGpuMetricsInfo(dvInd int) {
// 	fmt.Println()
// 	fmt.Printf("[DevGpuMetricsInfo] device=%d\n", dvInd)
// 	fmt.Println("  RSMI: rsmi_dev_gpu_metrics_info_get")
//
// 	metrics, err := dcgm.DevGpuMetricsInfo(dvInd)
// 	if err != nil {
// 		fmt.Printf("  RESULT: ERROR: %v\n", err)
// 		return
// 	}
//
// 	data, err := json.MarshalIndent(metrics, "  ", "  ")
// 	if err != nil {
// 		fmt.Printf("  RESULT: received metrics, JSON serialization failed: %v\n", err)
// 		return
// 	}
//
// 	fmt.Println("  RESULT: OK")
// 	fmt.Println(string(data))
// }
//
// func testDevThrottleFlags(dvInd int) {
// 	fmt.Println()
// 	fmt.Printf("[DevThrottleFlags] device=%d\n", dvInd)
// 	fmt.Println("  RSMI: rsmi_dev_gpu_metrics_info_get")
//
// 	flags, err := dcgm.DevThrottleFlags(dvInd)
// 	if err != nil {
// 		fmt.Printf("  RESULT: ERROR: %v\n", err)
// 		return
// 	}
//
// 	data, err := json.MarshalIndent(flags, "  ", "  ")
// 	if err != nil {
// 		fmt.Printf("  RESULT: received flags, JSON serialization failed: %v\n", err)
// 		return
// 	}
//
// 	fmt.Println("  RESULT: OK")
// 	fmt.Println(string(data))
// }
