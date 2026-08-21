/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	// 初始化 DCGM 服务。
	if err := dcgm.Init(); err != nil {
		glog.Errorf("dcgm init failed: %v", err)
		return
	}
	defer dcgm.ShutDown()

	processes, err := dcgm.ProcessHCUInfo()
	if err != nil {
		glog.Errorf("ProcessHCUInfo failed: %v", err)
	}

	fmt.Printf("ProcessHCUInfo returned %d process(es):\n", len(processes))
	for index, process := range processes {
		fmt.Printf("\n[Process %d]\n", index+1)
		fmt.Printf("  PID             : %d\n", process.ProcessID)
		fmt.Printf("  Process Name    : %s\n", process.ProcessName)
		fmt.Printf("  PASID           : %d\n", process.Pasid)
		fmt.Printf("  VRAM Usage      : %d bytes (%.2f MiB)\n",
			process.VramUsage, float64(process.VramUsage)/(1024*1024))
		fmt.Printf("  SDMA Usage      : %d\n", process.SdmaUsage)
		fmt.Printf("  CU Occupancy    : %d\n", process.CuOccupancy)
		fmt.Printf("  HCU Device IDs  : %v\n", process.MinorNumbers)
	}

	for _, pid := range []uint32{61756, 61527} {
		printProcessInfoByPID(pid)
	}
}

func printProcessInfoByPID(pid uint32) {
	info, err := dcgm.ProcessInfoByPid(pid)
	if err != nil {
		fmt.Printf("\n[ProcessInfoByPid: PID %d]\n  Query Error     : %v\n", pid, err)
		return
	}

	fmt.Printf("\n[ProcessInfoByPid: PID %d]\n", pid)
	fmt.Printf("  Process ID      : %d\n", info.ProcessID)
	fmt.Printf("  VRAM Usage      : %d MiB\n", info.VRAMUsageSize)
	fmt.Printf("  VRAM Usage Rate : %.2f%%\n", info.VRAMUsageRate)
	fmt.Printf("  HCU Device Count: %d\n", info.GPUCount)
	if info.GPUCount == 0 {
		fmt.Println("  HCU Devices     : none")
		return
	}

	fmt.Println("  HCU Devices:")
	for index, deviceID := range info.GPUIndex {
		usageRate := float32(0)
		if index < len(info.GPUUsageRate) {
			usageRate = info.GPUUsageRate[index]
		}
		fmt.Printf("    - Device ID %d, Usage Rate %.2f%%\n", deviceID, usageRate)
	}
}
