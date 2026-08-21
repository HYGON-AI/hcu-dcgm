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
	"regexp"
	"strconv"
	"strings"
)

// ---------------- GEMM 类型及期望性能 ----------------
var rel = []string{
	"dgemm", "sgemm", "hgemm", "i8gemm", "i4gemm",
	"i16gemm", "i32gemm", "hgemm(non-hpa)", "uint8gemm", "uint4gemm",
	"bf16gemm", "tf32gemm",
}

var expectedPerf = []float64{
	20.6, 40.6, 200.0, 330.1, 600.75, 0, 0, 0, 0, 0, 200.12, 110.1,
}

const (
	GEMMLogDir    = "logs/gemm"
	CurrentDevice = "BOWEN88"
)

// ---------------- 嵌入 gemmPerf ----------------

//go:embed resources/gemmPerf
var gemmPerfBytes []byte

func extractGemmPerf() (string, error) {
	if len(gemmPerfBytes) == 0 {
		return "", fmt.Errorf("embedded gemmPerf binary is empty")
	}

	tmpFile, err := os.CreateTemp("", "gemmPerf-*")
	if err != nil {
		return "", fmt.Errorf("无法创建临时文件: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(gemmPerfBytes); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("写入临时文件失败: %w", err)
	}

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("无法设置执行权限: %w", err)
	}

	return tmpFile.Name(), nil
}

// ---------------- 依赖检查 ----------------
func checkDependencies(binaryPath string) error {
	cmd := exec.Command("ldd", binaryPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行 ldd 失败: %v", err)
	}

	missing := []string{}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "not found") {
			missing = append(missing, strings.TrimSpace(line))
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("缺少依赖的 so 文件:\n%s", strings.Join(missing, "\n"))
	}
	return nil
}

// ---------------- 启动单个 GEMM 压测（增强版） ----------------
func runGemmTest(gemmPerfPath string, devInd int, gemmIdx int, iterations int, logfile string, mValue int) error {
	args := []string{
		"-m", strconv.Itoa(mValue),
		"-n", "4096",
		"-k", "12032",
		"-g", strconv.Itoa(gemmIdx),
		"-t", "1",
		"-i", strconv.Itoa(iterations),
	}

	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("无法打开日志文件: %w", err)
	}
	defer f.Close()

	cmd := exec.Command(gemmPerfPath, args...)
	cmd.Stdout = f
	cmd.Stderr = f // stderr 也写入日志

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("执行 gemmPerf 失败: %v，请查看日志: %s", err, logfile)
	}

	// 确认日志大小
	//info, err := os.Stat(logfile)
	//if err == nil {
	//	//fmt.Printf("Log written: %s (%d bytes)\n", logfile, info.Size())
	//}

	return nil
}

// ---------------- 解析 GEMM 日志 ----------------
func parseGemmLog(logfile string, devInd int) (mean float64, fail bool, err error) {
	f, err := os.Open(logfile)
	if err != nil {
		return 0, false, fmt.Errorf("无法打开日志文件: %w", err)
	}
	defer f.Close()

	re := regexp.MustCompile(`HCU` + strconv.Itoa(devInd) + `:.*mean:\s*([0-9]*\.?[0-9]+)`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); len(matches) == 2 {
			mean, err = strconv.ParseFloat(matches[1], 64)
			if err != nil {
				return 0, false, fmt.Errorf("解析 mean 失败: %w", err)
			}
		}
		if strings.Contains(line, "FAIL") {
			fail = true
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, false, fmt.Errorf("扫描日志失败: %w", err)
	}
	return mean, fail, nil
}

// ---------------- GEMM 压测入口 ----------------
func targetStressTest() {
	//InitGemmConfig(cfg)
	gemmPerfPath, err := extractGemmPerf()
	if err != nil {
		fmt.Printf("Failed to extract gemmPerf: %v\n", err)
		return
	}
	defer os.Remove(gemmPerfPath)

	// ----运行前检查依赖 ----
	if err := checkDependencies(gemmPerfPath); err != nil {
		fmt.Printf("依赖检查失败: %v\n", err)
		return
	}

	totalHCU, _ := rsmiNumMonitorDevices()
	//fmt.Printf("Total HCU: %d\n", totalHCU)
	if totalHCU <= 0 {
		fmt.Println("No HCU devices found")
		return
	}

	iterations := TargetStressTimeSeconds * 200 / (totalHCU * 4)
	if iterations < 1 {
		iterations = 1
	}

	if err := os.MkdirAll(GEMMLogDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		return
	}

	for devInd := 0; devInd < totalHCU; devInd++ {
		for _, gemm := range gemmList {
			logfile := filepath.Join(GEMMLogDir, fmt.Sprintf("%s_hcu%d.log", gemm.Name, devInd))
			mValue := 5632

			if err := runGemmTest(gemmPerfPath, devInd, gemm.Idx, iterations, logfile, mValue); err != nil {
				fmt.Printf("[gemmperf] device %d, gemm %s failed: %v\n", devInd, gemm.Name, err)
				continue
			}

			mean, _, err := parseGemmLog(logfile, devInd)
			if err != nil {
				fmt.Printf("[gemmperf parse] device %d, gemm %s failed: %v\n", devInd, gemm.Name, err)
				continue
			}

			fmt.Printf("TargetStress: HCU%d, %s mean %.2f\n", devInd, gemm.Name, mean)
		}
	}
	fmt.Printf("TargetStress completed!")
}

// ---------------- 结构化 API ----------------

// GemmTestResult 表示单条 GEMM 测试结果
type GemmTestResult struct {
	HCUId    int     // HCU 索引（设备编号）
	GemmName string  // GEMM 测试名称（如 "hgemm", "sgemm"）
	Mean     float64 // 测得的平均性能值（单位同你日志/工具，比如 GFLOPS 或 MB/s）
	Failed   bool    // 测试是否失败（true=失败/出错，false=成功）
}

// TargetStressResult 汇总整个测试的结构化结果
type TargetStressResult struct {
	Results []GemmTestResult // 每个 HCU/每个 GEMM 的测试结果列表
	LogDir  string           // 保存日志的目录路径
}

// runTargetStressTestWithResult 内部实现：复现 targetStressTest 的流程，但收集并返回结构化数据（不打印）
func runTargetStressTestWithResult() (TargetStressResult, error) {
	//InitGemmConfig(cfg)
	var res TargetStressResult

	gemmPerfPath, err := extractGemmPerf()
	if err != nil {
		return res, fmt.Errorf("Failed to extract gemmPerf: %v", err)
	}
	defer os.Remove(gemmPerfPath)

	// 依赖检查
	if err := checkDependencies(gemmPerfPath); err != nil {
		return res, fmt.Errorf("依赖检查失败: %v", err)
	}

	//fmt.Println("=== TargetStressTest start ===")

	totalHCU, _ := rsmiNumMonitorDevices()
	if totalHCU <= 0 {
		return res, fmt.Errorf("No HCU devices found")
	}

	//fmt.Printf("Will run GEMM stress test for %d devices\n", totalHCU)

	iterations := TargetStressTimeSeconds * 200 / (totalHCU * 4)
	if iterations < 1 {
		iterations = 1
	}

	if err := os.MkdirAll(GEMMLogDir, 0755); err != nil {
		return res, fmt.Errorf("Failed to create log directory: %v", err)
	}
	res.LogDir = GEMMLogDir

	results := make([]GemmTestResult, 0, totalHCU*len(gemmList))

	for devInd := 0; devInd < totalHCU; devInd++ {
		for _, gemm := range gemmList {
			logfile := filepath.Join(GEMMLogDir, fmt.Sprintf("%s_hcu%d.log", gemm.Name, devInd))
			mValue := 5632

			if err := runGemmTest(gemmPerfPath, devInd, gemm.Idx, iterations, logfile, mValue); err != nil {
				results = append(results, GemmTestResult{
					HCUId:    devInd,
					GemmName: gemm.Name,
					Mean:     0,
					Failed:   true,
				})
				continue
			}

			mean, fail, err := parseGemmLog(logfile, devInd)
			if err != nil {
				results = append(results, GemmTestResult{
					HCUId:    devInd,
					GemmName: gemm.Name,
					Mean:     0,
					Failed:   true,
				})
				continue
			}

			results = append(results, GemmTestResult{
				HCUId:    devInd,
				GemmName: gemm.Name,
				Mean:     mean,
				Failed:   fail,
			})
		}
	}

	res.Results = results
	return res, nil
}
