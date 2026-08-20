/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed resources/memtestCL
var memtestCLBytes []byte

const MemtestCLLogDir = "logs"

// 将嵌入的 memtestCL 写到临时文件并返回可执行路径
func extractMemtestCL() (string, error) {
	if len(memtestCLBytes) == 0 {
		return "", fmt.Errorf("embedded memtestCL binary is empty")
	}

	tmpFile, err := os.CreateTemp("", "memtestCL-*")
	if err != nil {
		return "", fmt.Errorf("无法创建临时文件: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(memtestCLBytes); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("写入临时文件失败: %w", err)
	}

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("无法设置执行权限: %w", err)
	}

	return tmpFile.Name(), nil
}

// 执行 memtestCL
func runMemtestCL(memtestCLPath string, dvId int, logDir string) (string, error) {
	logFile := filepath.Join(logDir, fmt.Sprintf("oam%d_memtestcl.log", dvId))
	cmd := exec.Command(memtestCLPath,
		"--hcu", fmt.Sprintf("%d", dvId),
		"-maxTestMem", "1", "20000", "1",
	)

	// 输出到日志文件
	f, err := os.Create(logFile)
	if err != nil {
		return "", fmt.Errorf("无法创建日志文件: %w", err)
	}
	defer f.Close()

	cmd.Stdout = f
	cmd.Stderr = f
	if err := cmd.Run(); err != nil {
		return logFile, fmt.Errorf("执行 memtestCL 失败: %w", err)
	}
	return logFile, nil
}

// 检查 memtestCL 日志结果并解析 summary
func checkMemtestCLResult(logFile string) (bool, map[string]string, error) {
	file, err := os.Open(logFile)
	if err != nil {
		return false, nil, fmt.Errorf("无法打开日志文件 %s: %w", logFile, err)
	}
	defer file.Close()

	summaryStarted := false
	results := make(map[string]string)
	finalPass := true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 进入 summary 部分
		if strings.HasPrefix(line, "Test summary:") {
			summaryStarted = true
			continue
		}
		if summaryStarted {
			// 解析每个测试项
			if strings.Contains(line, "failed iterations") {
				// 示例: "Moving inversions (ones and zeros): 0 failed iterations"
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					testName := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					results[testName] = value
					if !strings.Contains(value, "0 failed iterations") {
						finalPass = false
					}
				}
			}
			// Final error count
			if strings.HasPrefix(line, "Final error count") {
				if !strings.Contains(line, "0 errors") {
					finalPass = false
				}
				results["Final"] = line
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return false, nil, fmt.Errorf("读取日志文件失败: %w", err)
	}

	return finalPass, results, nil
}

// Demo: 批量运行多个 HCU 的 memtestCL
func memtestCL(dvIdList []int) error {
	memtestCLPath, err := extractMemtestCL()
	if err != nil {
		return fmt.Errorf("提取 memtestCL 失败: %w", err)
	}
	defer os.Remove(memtestCLPath)

	logDir := MemtestCLLogDir
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("无法创建日志目录: %w", err)
	}

	allPass := true
	for _, dvId := range dvIdList {
		logFile, err := runMemtestCL(memtestCLPath, dvId, logDir)
		if err != nil {
			fmt.Printf("[HCU %d] 运行失败: %v\n", dvId, err)
			allPass = false
			continue
		}

		ok, summary, err := checkMemtestCLResult(logFile)
		if err != nil {
			fmt.Printf("[HCU %d] 检查日志失败: %v\n", dvId, err)
			allPass = false
			continue
		}

		if !ok {
			fmt.Printf("[HCU %d] ❌ 测试失败\n", dvId)
			allPass = false
		}

		// 打印 summary 表
		fmt.Println("------ Test Summary ------")
		for k, v := range summary {
			fmt.Printf("%-35s : %s\n", k, v)
		}
		fmt.Println("--------------------------")
	}

	if !allPass {
		return fmt.Errorf("至少有一个 HCU 的 memtestCL 测试未通过")
	}
	//fmt.Println("所有 HCU memtestCL 测试均通过 ✅")
	return nil
}

// -------------------- 新增：结构化 API（NEW） --------------------

// MemtestCLResult 表示单个 HCU 的 memtestCL 运行与解析结果
type MemtestCLResult struct {
	HCUId   int               // HCU 索引（设备编号）
	Passed  bool              // 测试是否通过（true=通过，false=存在错误/失败）
	Summary map[string]string // 每个测试项的解析结果（键为测试项名称，值为统计信息，如 "0 failed iterations"）
	LogFile string            // 该 HCU 对应日志文件的路径
}

// MemtestCLAllResult 汇总多个 HCU 的 memtestCL 结果
type MemtestCLAllResult struct {
	Results []MemtestCLResult // 每个 HCU 的测试结果列表
	LogDir  string            // 存放所有日志文件的目录
}

// runMemtestCLWithResult 内部实现：复用原流程，但收集并返回结构化数据（不打印）
func runMemtestCLWithResult(dvIdList []int) (MemtestCLAllResult, error) {
	var out MemtestCLAllResult

	memtestCLPath, err := extractMemtestCL()
	if err != nil {
		return out, fmt.Errorf("提取 memtestCL 失败: %w", err)
	}
	defer os.Remove(memtestCLPath)

	logDir := MemtestCLLogDir
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return out, fmt.Errorf("无法创建日志目录: %w", err)
	}
	out.LogDir = logDir

	allPass := true
	results := make([]MemtestCLResult, 0, len(dvIdList))

	for _, dvId := range dvIdList {
		logFile, err := runMemtestCL(memtestCLPath, dvId, logDir)
		if err != nil {
			// record failure with empty summary
			results = append(results, MemtestCLResult{
				HCUId:   dvId,
				Passed:  false,
				Summary: map[string]string{"error": err.Error()},
				LogFile: logFile,
			})
			allPass = false
			continue
		}

		ok, summary, err := checkMemtestCLResult(logFile)
		if err != nil {
			results = append(results, MemtestCLResult{
				HCUId:   dvId,
				Passed:  false,
				Summary: map[string]string{"error": err.Error()},
				LogFile: logFile,
			})
			allPass = false
			continue
		}

		results = append(results, MemtestCLResult{
			HCUId:   dvId,
			Passed:  ok,
			Summary: summary,
			LogFile: logFile,
		})
		if !ok {
			allPass = false
		}
	}

	out.Results = results
	if !allPass {
		return out, fmt.Errorf("至少有一个 HCU 的 memtestCL 测试未通过")
	}
	return out, nil
}
