// Package dcgm 提供一个用于对 GPU 卡进行功耗压测的简单工具集。
// 功能包括：启动卡的监控进程 (hy-smi)、启动 workload（通过 numactl 启动
// 内嵌的 gemm 压测程序）、等待 workload 完成并解析监控日志来计算每卡的最大功耗与平均功耗。
/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// -------------------- 内嵌资源 --------------------

//go:embed resources/gemmPower resources/int8_cu_check.co
var embeddedFiles embed.FS

// extractGemmPower 把内嵌的 gemmPower 可执行文件写到临时目录并返回路径。
func extractGemmPower() (string, error) {
	data, err := embeddedFiles.ReadFile("resources/gemmPower")
	if err != nil {
		return "", fmt.Errorf("读取内嵌 gemmPower 失败: %w", err)
	}
	tmpFile, err := os.CreateTemp("", "gemmPower-*")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(data); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("写入临时文件失败: %w", err)
	}
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("设置执行权限失败: %w", err)
	}
	return tmpFile.Name(), nil
}

// extractKernel 把内嵌的 int8_cu_check.co 写到临时目录并返回路径。
func extractKernel() (string, error) {
	data, err := embeddedFiles.ReadFile("resources/int8_cu_check.co")
	if err != nil {
		return "", fmt.Errorf("读取内嵌 kernel 失败: %w", err)
	}
	tmpFile, err := os.CreateTemp("", "int8_cu_check-*.co")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(data); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("写入临时文件失败: %w", err)
	}
	return tmpFile.Name(), nil
}

// -------------------- 结构体定义 --------------------

// PowerInfo 描述一次压测所需的基本信息与资源分配。
type PowerInfo struct {
	TotalHCU        int
	BusID           []string // 目前未使用，保留将来扩展
	ToDriverID      []int
	ToOAMID         []int
	NumAID          []int
	CurrentHCU      int    // goroutine 中会写入，记录当前卡索引
	TargetPowerTime int    // minutes
	LogDir          string // 日志目录
}

// -------------------- 辅助函数 --------------------

// startMonitor 启动 hy-smi 监控进程
func startMonitor(driverID int, logPath string) (*exec.Cmd, error) {
	cmd := exec.Command("hy-smi", "-idmon", strconv.Itoa(driverID),
		"-sdmon", "p", "-odmon", "DT", "-fdmon", logPath)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("启动 hy-smi 失败: %w", err)
	}
	return cmd, nil
}

// startGemmPower 启动 gemmPower workload
func startGemmPower(exePath, kernelPath string, driverID, numaID, durationSec int, logfile string, extraArgs []string) (*exec.Cmd, error) {
	commonArgs := []string{
		"-m", "2560",
		"-n", "2048",
		"-k", "12032",
		"-a", "5632",
		"-b", "4096",
		"-g", "3",
		"-t", "1",
		"-f", kernelPath,
		"-d", strconv.Itoa(durationSec),
	}
	allExeArgs := append(commonArgs, extraArgs...)

	fullArgs := []string{"--cpubind", strconv.Itoa(numaID),
		"--membind", strconv.Itoa(numaID),
		exePath}
	fullArgs = append(fullArgs, allExeArgs...)

	cmd := exec.Command("numactl", fullArgs...)
	cmd.Env = append(os.Environ(), "HIP_VISIBLE_DEVICES="+strconv.Itoa(driverID))

	lf, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %w", err)
	}
	cmd.Stdout = lf
	cmd.Stderr = lf

	if err := cmd.Start(); err != nil {
		lf.Close()
		return nil, fmt.Errorf("启动 gemmPower 失败: %w", err)
	}

	go func() {
		cmd.Wait()
		lf.Close()
	}()

	return cmd, nil
}

// processPowerData 解析 hy-smi 生成的日志
func processPowerData(cardnum int, logDir string) (float64, float64, error) {
	path := fmt.Sprintf("%s/OAM%d_power.log", logDir, cardnum)
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, fmt.Errorf("打开日志失败: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var sum, max float64
	var count int
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 6 {
			continue
		}
		p, err := strconv.ParseFloat(parts[5], 64)
		if err != nil || p <= 300 {
			continue
		}
		if p > max {
			max = p
		}
		sum += p
		count++
	}
	if count == 0 {
		return 0, 0, fmt.Errorf("日志 %s 没有有效数据", path)
	}
	return max, sum / float64(count), nil
}

// -------------------- 对外导出函数 --------------------

// TargetPower 启动压测流程：监控 + workload + 日志解析
// 注意：exePath 已经内嵌，不再需要外部传入。
func TargetPower(powerInfo *PowerInfo) error {
	if powerInfo.TargetPowerTime == 0 {
		return nil
	}

	// 确保日志目录存在
	if err := os.MkdirAll(powerInfo.LogDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 提取内嵌可执行文件
	exePath, err := extractGemmPower()
	if err != nil {
		return err
	}
	kernelPath, err := extractKernel()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, powerInfo.TotalHCU)

	for card := 0; card < powerInfo.TotalHCU; card++ {
		wg.Add(1)
		go func(cardnum int) {
			defer wg.Done()
			powerInfo.CurrentHCU = cardnum

			driverID := powerInfo.ToDriverID[cardnum]
			oamID := powerInfo.ToOAMID[cardnum]
			numa := powerInfo.NumAID[cardnum]
			logfile := fmt.Sprintf("%s/OAM%d_power.log", powerInfo.LogDir, oamID)

			monCmd, err := startMonitor(driverID, logfile)
			if err != nil {
				errCh <- fmt.Errorf("card %d monitor 启动失败: %w", cardnum, err)
				return
			}
			defer func() {
				_ = monCmd.Process.Signal(os.Interrupt)
				time.Sleep(500 * time.Millisecond)
				_ = monCmd.Process.Kill()
				monCmd.Process.Wait()
			}()

			durationSec := powerInfo.TargetPowerTime * 60
			gemmCmd, err := startGemmPower(exePath, kernelPath, driverID, numa, durationSec, logfile, nil)
			if err != nil {
				errCh <- fmt.Errorf("card %d workload 启动失败: %w", cardnum, err)
				return
			}
			_ = gemmCmd.Wait()
			errCh <- nil
		}(card)
	}

	wg.Wait()
	close(errCh)

	var anyErr bool
	for err := range errCh {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			anyErr = true
		}
	}

	// 解析结果
	for card := 0; card < powerInfo.TotalHCU; card++ {
		maxp, avg, err := processPowerData(powerInfo.ToOAMID[card], powerInfo.LogDir)
		if err != nil {
			fmt.Printf("OAM%d 解析失败: %v\n", powerInfo.ToOAMID[card], err)
			continue
		}
		fmt.Printf("targeted_power: OAM%d max: %.1f W avg: %.1f W\n",
			powerInfo.ToOAMID[card], maxp, avg)
	}

	if anyErr {
		return fmt.Errorf("部分卡运行失败，详情见日志")
	}
	return nil
}
