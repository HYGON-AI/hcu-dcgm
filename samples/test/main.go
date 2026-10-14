/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package main

import "C"
import (
	"fmt"

	"github.com/goccy/go-json"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

func main() {
	fmt.Println("===== DevGpuMetricsInfo demo =====")
	fmt.Println("Testing devices: 0, 1")

	if err := dcgm.Init(); err != nil {
		fmt.Printf("dcgm.Init ERROR: %v\n", err)
		return
	}
	defer func() {
		if err := dcgm.ShutDown(); err != nil {
			fmt.Printf("dcgm.ShutDown ERROR: %v\n", err)
		}
	}()

	for _, dvInd := range []int{0, 1} {
		testDevThrottleFlags(dvInd)
		testDevGpuMetricsInfo(dvInd)

	}

	fmt.Println("===== DCGM demo finished =====")
}

func testDevGpuMetricsInfo(dvInd int) {
	fmt.Println()
	fmt.Printf("[DevGpuMetricsInfo] device=%d\n", dvInd)
	fmt.Println("  RSMI: rsmi_dev_gpu_metrics_info_get")

	metrics, err := dcgm.DevGpuMetricsInfo(dvInd)
	if err != nil {
		fmt.Printf("  RESULT: ERROR: %v\n", err)
		return
	}

	data, err := json.MarshalIndent(metrics, "  ", "  ")
	if err != nil {
		fmt.Printf("  RESULT: received metrics, JSON serialization failed: %v\n", err)
		return
	}

	fmt.Println("  RESULT: OK")
	fmt.Println(string(data))
}

func testDevThrottleFlags(dvInd int) {
	fmt.Println()
	fmt.Printf("[DevThrottleFlags] device=%d\n", dvInd)
	fmt.Println("  RSMI: rsmi_dev_gpu_metrics_info_get")

	flags, err := dcgm.DevThrottleFlags(dvInd)
	if err != nil {
		fmt.Printf("  RESULT: ERROR: %v\n", err)
		return
	}

	data, err := json.MarshalIndent(flags, "  ", "  ")
	if err != nil {
		fmt.Printf("  RESULT: received flags, JSON serialization failed: %v\n", err)
		return
	}

	fmt.Println("  RESULT: OK")
	fmt.Println(string(data))
}

// 原有的利用率、ECC、温度、拓扑、带宽和压测 demo 不在本次测试中执行。
