/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/golang/glog"
)

// --- 嵌入 hip-stream 二进制 ---
//
//go:embed resources/hip-stream
var hipStreamBytes []byte

const (
	DiagLogDirBandwidth = "logs/hip-stream"
	logDir              = DiagLogDirBandwidth
)

// bandResult 用于通过 channel 返回带 dvInd 的结果
type bandResult struct {
	DvInd int
	BW    float64 // GB/s
	Err   error
}

// extractHipStream 将嵌入的二进制写到临时文件并返回路径
func extractHipStream() (string, error) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "hip-stream-*")
	if err != nil {
		return "", fmt.Errorf("无法创建临时文件: %w", err)
	}
	defer tmpFile.Close()

	// 写入嵌入内容
	if _, err := tmpFile.Write(hipStreamBytes); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("无法写入临时文件: %w", err)
	}

	// 设置可执行权限
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("无法设置可执行权限: %w", err)
	}

	return tmpFile.Name(), nil
}

// testHCUBandwidth 在单个设备上运行 hip-stream，并把结果通过 results channel 返回。
// ctx 用于取消，函数会记录详细的 debug 日志到 logDir。
func testHCUBandwidth(ctx context.Context, dvInd int, wg *sync.WaitGroup, results chan<- bandResult) {
	defer wg.Done()

	// 早期取消检查
	select {
	case <-ctx.Done():
		results <- bandResult{DvInd: dvInd, BW: 0, Err: fmt.Errorf("HCU%d: %w", dvInd, ctx.Err())}
		return
	default:
	}

	hipStreamPath, err := extractHipStream()
	if err != nil {
		failure := fmt.Errorf("HCU%d: 提取 hip-stream 失败: %w", dvInd, err)
		glog.Error(failure)
		results <- bandResult{DvInd: dvInd, BW: 0, Err: failure}
		return
	}
	// 在函数结束时删除临时文件
	defer func() {
		if err := os.Remove(hipStreamPath); err != nil {
			glog.Errorf("无法删除临时文件 %s: %v", hipStreamPath, err)
		}
	}()

	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		failure := fmt.Errorf("HCU%d: 无法创建带宽日志目录 %s: %w", dvInd, logDir, err)
		glog.Error(failure)
		results <- bandResult{DvInd: dvInd, BW: 0, Err: failure}
		return
	}

	// 这里的 cmdStr 根据你先前讨论设置：可按需调整 precision/numtimes/arraysize
	// 例子（你可以修改这些参数）：
	// precision=1/2, numtimes=100/200, arraysize=...
	// 下面用你之前最终推荐的组合示例（注意精度/arraysize 需你确认）
	cmdStr := fmt.Sprintf(
		"export HIP_VISIBLE_DEVICES=%d; %s --precision %d --numtimes %d --arraysize %d",
		dvInd,
		hipStreamPath,
		bandwidthPrecision,
		bandwidthNumtimes,
		bandwidthArraysize,
	)

	logFile := fmt.Sprintf("%s/hbm_bandwidth_oam%d.log", logDir, dvInd)
	debugFile := fmt.Sprintf("%s/hbm_bandwidth_oam%d.debug.log", logDir, dvInd)

	glog.Infof("OAM%d 开始带宽测试，cmd=%s", dvInd, cmdStr)
	start := time.Now()

	// 使用 CommandContext，这样 ctx 取消会杀掉子进程
	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)

	// 捕获全部输出
	out, err := cmd.CombinedOutput()
	duration := time.Since(start)

	// 把输出写到 logFile
	if writeErr := os.WriteFile(logFile, out, 0644); writeErr != nil {
		glog.Errorf("写入 logFile %s 失败: %v", logFile, writeErr)
	}

	// debug 日志：包含元信息与输出（可能较大）
	debugBuf := bytes.NewBuffer(nil)
	fmt.Fprintf(debugBuf, "cmd=%s\nstart=%s\nduration=%s\n", cmdStr, start.Format(time.RFC3339), duration)
	if cmd.Process != nil {
		fmt.Fprintf(debugBuf, "pid=%d\n", cmd.Process.Pid)
	}
	fmt.Fprintf(debugBuf, "error=%v\n\n--- output ---\n", err)
	debugBuf.Write(out)
	if writeErr := os.WriteFile(debugFile, debugBuf.Bytes(), 0644); writeErr != nil {
		glog.Errorf("写入 debugFile %s 失败: %v", debugFile, writeErr)
	}

	// 如果子进程返回错误，保留 exit code/signal 并将设备级错误返回给调用方。
	if err != nil {
		failure := fmt.Errorf("HCU%d: hip-stream 执行失败: %w；日志: %s, %s", dvInd, err, logFile, debugFile)
		if exitErr, ok := err.(*exec.ExitError); ok {
			if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitCode := ws.ExitStatus()
				if ws.Signaled() {
					failure = fmt.Errorf("HCU%d: hip-stream 非零退出, exit=%d signal=%v: %w；日志: %s, %s", dvInd, exitCode, ws.Signal(), err, logFile, debugFile)
				} else {
					failure = fmt.Errorf("HCU%d: hip-stream 非零退出, exit=%d: %w；日志: %s, %s", dvInd, exitCode, err, logFile, debugFile)
				}
				if exitCode == 137 {
					failure = fmt.Errorf("%w；进程可能被 SIGKILL，请检查 dmesg 或 kernel 日志是否存在 OOM/外部终止", failure)
				}
			}
		}
		glog.Error(failure)

		outStr := string(out)
		if len(outStr) > 4096 {
			outStr = outStr[:4096] + "\n...(truncated)"
		}
		glog.Errorf("OAM%d output (truncated):\n%s", dvInd, outStr)

		results <- bandResult{DvInd: dvInd, BW: 0, Err: failure}
		return
	}

	// 解析输出，寻找包含带宽信息的行（示例匹配 "Read <num>"）
	// 注意：不同版本工具输出格式不同，请根据实际输出调整解析逻辑
	scanner := bufio.NewScanner(bytes.NewReader(out))
	found := false
	for scanner.Scan() {
		// 运行时也检查 ctx
		select {
		case <-ctx.Done():
			results <- bandResult{DvInd: dvInd, BW: 0, Err: fmt.Errorf("HCU%d: %w", dvInd, ctx.Err())}
			return
		default:
		}

		line := scanner.Text()
		// 常见格式示例："Read 770.823 MB/s ..." 或 "Read 770.823"
		if strings.Contains(line, "Read") {
			parts := strings.Fields(line)
			// 尝试取 parts[1] 为数字（兼容多种输出）
			if len(parts) >= 2 {
				if bw, err := strconv.ParseFloat(parts[1], 64); err == nil {
					// 假定输出单位为 MB/s，转换为 GB/s
					results <- bandResult{DvInd: dvInd, BW: bw / 1000.0}
					found = true
					break
				}
			}
		}
		// 另外也检测可能是 "Bandwidth: 770.823 GB/s" 这种
		if strings.Contains(line, "Bandwidth") && strings.Contains(line, "GB/s") {
			// 从行中抽取第一个数字
			fields := strings.Fields(line)
			for _, f := range fields {
				if val, err := strconv.ParseFloat(strings.TrimSuffix(f, "GB/s"), 64); err == nil {
					results <- bandResult{DvInd: dvInd, BW: val}
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}
	if !found {
		failure := fmt.Errorf("HCU%d: hip-stream 输出中未找到带宽结果；日志: %s, %s", dvInd, logFile, debugFile)
		glog.Warning(failure)
		results <- bandResult{DvInd: dvInd, BW: 0, Err: failure}
	}
}

// runBandwidthTest 运行多卡带宽测试，verbose 控制是否打印到 stdout（banner）。
func runBandwidthTest(ctx context.Context, dvIdList []int, verbose bool) (map[int]float64, error) {
	var wg sync.WaitGroup
	results := make(chan bandResult, len(dvIdList))

	for _, dvInd := range dvIdList {
		wg.Add(1)
		go testHCUBandwidth(ctx, dvInd, &wg, results)
	}

	// 等待所有 goroutine 完成
	wg.Wait()
	close(results)

	// 收集结果
	bwMap := make(map[int]float64, len(dvIdList))
	failures := make([]string, 0)
	for res := range results {
		bwMap[res.DvInd] = res.BW
		if res.Err != nil {
			failures = append(failures, res.Err.Error())
		}
	}

	if len(failures) > 0 {
		return bwMap, fmt.Errorf("带宽测试失败: %s", strings.Join(failures, "; "))
	}

	// verbose 控制是否打印 banner（仅 CLI 场景希望看到）
	if verbose {
		fmt.Printf("===== 带宽测试结果 =====\n")
		for _, dvInd := range dvIdList {
			bw := bwMap[dvInd]
			fmt.Printf("HCU%d: %.2f GB/s\n", dvInd, bw)
		}
	}

	// respect ctx cancel
	select {
	case <-ctx.Done():
		return bwMap, ctx.Err()
	default:
	}

	return bwMap, nil
}

// BandwidthTest (CLI)：打印 banner，返回 bool（成功/失败）
func bandwidthTest(dvIdList []int) bool {
	ctx := context.Background()
	bwMap, err := runBandwidthTest(ctx, dvIdList, true) // verbose=true -> 打印 banner
	if err != nil {
		glog.Errorf("带宽测试被取消或出错: %v", err)
		return false
	}
	// 检查是否所有结果都是 0（认为失败）
	allZero := true
	for _, bw := range bwMap {
		if bw > 0 {
			allZero = false
			break
		}
	}
	if allZero {
		glog.Errorf("带宽测试结果为空或全部为0")
		return false
	}
	return true
}

// BandwidthTestResult (API/diag)：不打印 banner，返回 map 和 error
func bandwidthTestResult(dvIdList []int) (map[int]float64, error) {
	ctx := context.Background()
	bwMap, err := runBandwidthTest(ctx, dvIdList, false) // verbose=false -> 不打印
	if err != nil {
		glog.Errorf("带宽测试被取消或出错: %v", err)
		return nil, err
	}
	// 检查是否所有结果都是 0（认为失败）
	allZero := true
	for _, bw := range bwMap {
		if bw > 0 {
			allZero = false
			break
		}
	}
	if allZero {
		err = fmt.Errorf("带宽测试结果为空或全部为0")
		glog.Error(err)
		return bwMap, err
	}
	return bwMap, nil
}
