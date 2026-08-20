/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

const XHCLLogDir = "logs/rocm-bandwidth-test/xhcl"

// -------------------------------------------
// 原有逻辑（保持不变）
// -------------------------------------------

func hcuXHCLTest() bool {
	// 获取 HCU 总数（占位函数）
	totalHCU, _ := rsmiNumMonitorDevices()
	if totalHCU <= 0 {
		glog.Errorf("获取 HCU 总数失败")
		return false
	}

	// 提取可执行并写到临时文件
	execPath, err := extractRocmBandwidth()
	if err != nil {
		log.Printf("提取 rocm-bandwidth-test 失败: %v", err)
		return false
	}
	defer os.Remove(execPath)

	// 确保日志目录存在
	if err := os.MkdirAll(XHCLLogDir, 0755); err != nil {
		log.Printf("创建日志目录失败: %v", err)
		return false
	}
	logFile := filepath.Join(XHCLLogDir, "xhcl_bandwidth_all.log")

	// 运行可执行文件
	lf, err := os.Create(logFile)
	if err != nil {
		log.Printf("无法创建日志文件: %v", err)
		return false
	}
	defer lf.Close()

	cmd := exec.Command(execPath) // 不带参数
	cmd.Stdout = lf
	cmd.Stderr = lf
	if err := cmd.Run(); err != nil {
		glog.Errorf("rocm-bandwidth-test 执行失败: %v", err)
		return false
	}

	// 解析并打印表格
	parseXHCLLogAll(logFile, totalHCU)

	return true
}

func parseXHCLLogAll(logFile string, totalHCU int) {
	data, err := os.ReadFile(logFile)
	if err != nil {
		glog.Errorf("无法读取日志文件 %s: %v", logFile, err)
		return
	}
	lines := strings.Split(string(data), "\n")

	deviceCount := 0
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "Device:") {
			deviceCount++
		}
	}
	if deviceCount == 0 {
		glog.Errorf("Device 数量为0，无法解析")
		return
	}

	cpuCount := deviceCount - totalHCU
	if cpuCount <= 0 {
		cpuCount = 1
	}

	B := make([][]float64, deviceCount)
	for i := 0; i < deviceCount; i++ {
		B[i] = make([]float64, deviceCount)
	}

	inBi := false
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Bidirectional copy peak bandwidth") {
			inBi = true
			continue
		}
		if !inBi || strings.HasPrefix(line, "D/D") || len(line) < 2 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0][0] < '0' || fields[0][0] > '9' {
			continue
		}
		rowIdx, err := strconv.Atoi(fields[0])
		if err != nil || rowIdx < 0 || rowIdx >= deviceCount {
			continue
		}

		for col := 0; col < deviceCount && col+1 < len(fields); col++ {
			f := fields[col+1]
			var v float64
			if f == "N/A" {
				v = 0
			} else {
				if parsed, err := strconv.ParseFloat(f, 64); err == nil {
					v = parsed
				}
			}
			B[rowIdx][col] = v
		}
	}

	printMatrixWithDeviceLabels("Bidirectional copy peak bandwidth GB/s (XHCL)", B, "device")

	for dvInd := 0; dvInd < totalHCU; dvInd++ {
		hcuRow := cpuCount + dvInd
		if hcuRow < 0 || hcuRow >= deviceCount {
			glog.Errorf("hcuRow 越界: hcuRow=%d deviceCount=%d", hcuRow, deviceCount)
			continue
		}

		for i := 0; i < totalHCU; i++ {
			if i == dvInd {
				continue
			}
			otherRow := cpuCount + i
			val := B[hcuRow][otherRow]
			fmt.Printf("HCU%d <-> HCU%d XHCL: %.3f GB/s\n", dvInd, i, val)
		}
	}
}

// -------------------------------------------
// 新增：结构化 API 接口（NEW）
// -------------------------------------------

// XHCLResult 表示一次 XHCL 带宽测试的结构化结果
type XHCLResult struct {
	SrcHCUId     int     // 源 HCU（起点）的索引
	DstHCUId     int     // 目标 HCU（终点）的索引
	BandwidthGBs float64 // 带宽（GB/s）
}

// -------------------------------------------
// 内部实现（NEW）
// -------------------------------------------

func runHcuXHCLTestWithResult() ([]XHCLResult, error) {
	// 获取 HCU 总数
	totalHCU, _ := rsmiNumMonitorDevices()
	if totalHCU <= 0 {
		return nil, fmt.Errorf("获取 HCU 总数失败")
	}

	// 提取 rocm-bandwidth-test 可执行文件
	execPath, err := extractRocmBandwidth()
	if err != nil {
		return nil, fmt.Errorf("提取 rocm-bandwidth-test 失败: %v", err)
	}
	defer os.Remove(execPath)

	// 确保日志目录存在
	if err := os.MkdirAll(XHCLLogDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}
	logFile := filepath.Join(XHCLLogDir, "xhcl_bandwidth_all.log")

	// 执行测试
	lf, err := os.Create(logFile)
	if err != nil {
		return nil, fmt.Errorf("无法创建日志文件: %v", err)
	}
	defer lf.Close()

	cmd := exec.Command(execPath)
	cmd.Stdout = lf
	cmd.Stderr = lf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("rocm-bandwidth-test 执行失败: %v", err)
	}

	// 解析日志文件为结构化结果
	results, err := parseXHCLLogAllToResult(logFile, totalHCU)
	if err != nil {
		return nil, fmt.Errorf("解析日志文件失败: %v", err)
	}

	return results, nil
}

// parseXHCLLogAllToResult 解析日志文件并返回结构化带宽数据（不打印）
func parseXHCLLogAllToResult(logFile string, totalHCU int) ([]XHCLResult, error) {
	data, err := os.ReadFile(logFile)
	if err != nil {
		return nil, fmt.Errorf("无法读取日志文件 %s: %v", logFile, err)
	}
	lines := strings.Split(string(data), "\n")

	deviceCount := 0
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "Device:") {
			deviceCount++
		}
	}
	if deviceCount == 0 {
		return nil, fmt.Errorf("Device 数量为0，无法解析")
	}

	cpuCount := deviceCount - totalHCU
	if cpuCount <= 0 {
		cpuCount = 1
	}

	B := make([][]float64, deviceCount)
	for i := 0; i < deviceCount; i++ {
		B[i] = make([]float64, deviceCount)
	}

	inBi := false
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Bidirectional copy peak bandwidth") {
			inBi = true
			continue
		}
		if !inBi || strings.HasPrefix(line, "D/D") || len(line) < 2 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0][0] < '0' || fields[0][0] > '9' {
			continue
		}
		rowIdx, err := strconv.Atoi(fields[0])
		if err != nil || rowIdx < 0 || rowIdx >= deviceCount {
			continue
		}

		for col := 0; col < deviceCount && col+1 < len(fields); col++ {
			f := fields[col+1]
			var v float64
			if f == "N/A" {
				v = 0
			} else if parsed, err := strconv.ParseFloat(f, 64); err == nil {
				v = parsed
			}
			B[rowIdx][col] = v
		}
	}

	// 生成结构化结果
	var results []XHCLResult
	for src := 0; src < totalHCU; src++ {
		srcRow := cpuCount + src
		for dst := 0; dst < totalHCU; dst++ {
			if src == dst {
				continue
			}
			dstRow := cpuCount + dst
			val := B[srcRow][dstRow]
			results = append(results, XHCLResult{
				SrcHCUId:     src,
				DstHCUId:     dst,
				BandwidthGBs: val,
			})
		}
	}

	return results, nil
}
