/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// 嵌入的可执行文件 bytes
//
//go:embed resources/rocm-bandwidth-test
var rocmBandwidthBytes []byte

// 日志目录常量
const PCIELogDir = "logs/rocm-bandwidth-test/pcie"

// 将嵌入的 rocm 可执行文件写到临时文件并设置可执行权限
// UNCHANGED
func extractRocmBandwidth() (string, error) {
	if len(rocmBandwidthBytes) == 0 {
		return "", fmt.Errorf("embedded rocm-bandwidth-test binary is empty")
	}
	tmpFile, err := os.CreateTemp("", "rocm-bandwidth-*")
	if err != nil {
		return "", fmt.Errorf("无法创建临时文件: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(rocmBandwidthBytes); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("写入临时文件失败: %w", err)
	}

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("无法设置执行权限: %w", err)
	}
	return tmpFile.Name(), nil
}

// ----------------------------
// 新增的类型和函数（NEW）
// ----------------------------

// HCUBW 表示每个 HCU 的 Sys->Fb 和 Fb->Sys 带宽信息。 // NEW
type HCUBW struct {
	DvInd   int
	SysToFb float64
	FbToSys float64
}

// PcieResult 表示 PCIe 拓扑与带宽测试的结构化解析结果。
type PcieResult struct {
	Matrix      [][]float64 // 带宽矩阵（GB/s）：Matrix[i][j] = 设备 i 到 j 的单向带宽
	DeviceCount int         // 检测到的 HCU 数量，用于确认 Matrix 的有效维度
	CPUCount    int         // 参与拓扑识别的 CPU root 数量（如多 socket 系统）
	HCUs        []HCUBW     // 每个 HCU 的带宽和关联信息，仅在存在 HCU 时有效
	LogFile     string      // 原始日志文件路径
}

func validatePcieResult(result PcieResult) error {
	if result.DeviceCount <= 0 || len(result.HCUs) == 0 {
		return fmt.Errorf("PCIe bandwidth result contains no HCU data")
	}
	for _, hcu := range result.HCUs {
		if hcu.SysToFb <= 0 || hcu.FbToSys <= 0 {
			return fmt.Errorf("HCU %d PCIe bandwidth is invalid: Sys->Fb=%.3f GB/s, Fb->Sys=%.3f GB/s", hcu.DvInd, hcu.SysToFb, hcu.FbToSys)
		}
	}
	return nil
}

// runPcieBandwidthTestWithResult 执行可执行文件、写日志并解析为结构化结果。 // NEW
func runPcieBandwidthTestWithResult() (PcieResult, error) {
	var res PcieResult

	totalHCU, _ := rsmiNumMonitorDevices()
	if totalHCU <= 0 {
		return res, fmt.Errorf("获取 HCU 总数失败")
	}

	execPath, err := extractRocmBandwidth()
	if err != nil {
		return res, fmt.Errorf("提取 rocm-bandwidth-test 失败: %w", err)
	}
	// 保证临时文件被删除
	defer os.Remove(execPath)

	// 确保日志目录存在
	if err := os.MkdirAll(PCIELogDir, 0755); err != nil {
		return res, fmt.Errorf("创建日志目录失败: %w", err)
	}
	logFile := filepath.Join(PCIELogDir, "pcie_bandwidth_all.log")
	res.LogFile = logFile
	// 为了避免文件打开冲突，先创建（或截断）再写入
	outf, err := os.Create(logFile)
	if err != nil {
		return res, fmt.Errorf("无法创建日志文件: %w", err)
	}
	// 需要在结束时关闭
	defer outf.Close()

	cmd := exec.Command(execPath) // 不带参数
	cmd.Stdout = outf
	cmd.Stderr = outf
	if err := cmd.Run(); err != nil {
		return res, fmt.Errorf("rocm-bandwidth-test 执行失败: %w", err)
	}

	// 解析日志为结构化结果
	parsed, err := parsePcieLogAllToResult(logFile, totalHCU)
	if err != nil {
		return res, fmt.Errorf("解析日志失败: %w", err)
	}
	parsed.LogFile = logFile
	if err := validatePcieResult(parsed); err != nil {
		return res, err
	}
	return parsed, nil
}

// parsePcieLogAllToResult 将日志文件解析成 PcieResult（原 parsePcieLogAll 的逻辑改为返回数据）。 // NEW
func parsePcieLogAllToResult(logFile string, totalHCU int) (PcieResult, error) {
	var res PcieResult

	data, err := os.ReadFile(logFile)
	if err != nil {
		return res, fmt.Errorf("无法读取日志文件 %s: %w", logFile, err)
	}
	lines := strings.Split(string(data), "\n")

	// 统计设备总数
	deviceCount := 0
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "Device:") {
			deviceCount++
		}
	}
	if deviceCount == 0 {
		return res, fmt.Errorf("Device 数量为0，无法解析")
	}
	res.DeviceCount = deviceCount

	cpuCount := deviceCount - totalHCU
	if cpuCount <= 0 {
		cpuCount = 1
	}
	res.CPUCount = cpuCount

	// 只保留 PCIe 单向矩阵
	M := make([][]float64, deviceCount)
	for i := 0; i < deviceCount; i++ {
		M[i] = make([]float64, deviceCount)
	}

	inUni := false
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Unidirectional copy peak bandwidth") {
			inUni = true
			continue
		}
		if strings.HasPrefix(line, "Bidirectional copy peak bandwidth") {
			inUni = false
			continue
		}
		if !inUni {
			continue
		}
		if strings.HasPrefix(line, "D/D") || len(line) < 2 {
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
			} else {
				v = 0
			}
			M[rowIdx][col] = v
		}
	}

	res.Matrix = M

	// 计算每个 HCU 的 Sys->Fb 和 Fb->Sys
	hcus := make([]HCUBW, 0, totalHCU)
	for dvInd := 0; dvInd < totalHCU; dvInd++ {
		hcuRow := cpuCount + dvInd
		if hcuRow < 0 || hcuRow >= deviceCount {
			// 越界则返回 0 值
			hcus = append(hcus, HCUBW{DvInd: dvInd, SysToFb: 0, FbToSys: 0})
			continue
		}

		sysToFb := 0.0
		for r := 0; r < cpuCount && r < deviceCount; r++ {
			if M[r][hcuRow] > sysToFb {
				sysToFb = M[r][hcuRow]
			}
		}

		fbToSys := 0.0
		for c := 0; c < cpuCount && c < deviceCount; c++ {
			if M[hcuRow][c] > fbToSys {
				fbToSys = M[hcuRow][c]
			}
		}

		hcus = append(hcus, HCUBW{
			DvInd:   dvInd,
			SysToFb: sysToFb,
			FbToSys: fbToSys,
		})
	}
	res.HCUs = hcus

	return res, nil
}

// runPcieBandwidthTest 之前是直接执行并 parse -> print，现在改为调用带结果的版本再打印（CHANGED）
func runPcieBandwidthTest() bool {
	res, err := runPcieBandwidthTestWithResult()
	if err != nil {
		glog.Errorf("rocm-bandwidth-test 失败: %v", err)
		return false
	}
	// 保持原来控制台打印格式以兼容现有 CLI 行为
	printMatrixWithDeviceLabels("Unidirectional copy peak bandwidth GB/s (PCIe)", res.Matrix, "device")
	for _, h := range res.HCUs {
		fmt.Printf("HCU%d PCIe Sys->Fb: %.3f GB/s, Fb->Sys: %.3f GB/s\n", h.DvInd, h.SysToFb, h.FbToSys)
	}
	return true
}

// printMatrixWithDeviceLabels 打印矩阵（UNCHANGED）
func printMatrixWithDeviceLabels(title string, M [][]float64, labelPrefix string) {
	fmt.Println()
	fmt.Printf("  %s\n\n", title)

	fmt.Printf("      ")
	for j := range M {
		fmt.Printf("%12s", fmt.Sprintf("%s%d", labelPrefix, j))
	}
	fmt.Println()

	for i := range M {
		fmt.Printf("%4s ", fmt.Sprintf("%s%d", labelPrefix, i))
		for j := range M {
			if i == j {
				fmt.Printf("%12s", "N/A")
			} else {
				fmt.Printf("%12.3f", M[i][j])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
