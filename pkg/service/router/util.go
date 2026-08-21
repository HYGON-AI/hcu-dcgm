/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package router

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

// 将字符串转换为 RSMIDevPerfLevel 类型
func ConvertToRSMIDevPerfLevel(level string) (dcgm.DevPerfLevel, error) {
	switch level {
	case "AUTO":
		return dcgm.RSMI_DEV_PERF_LEVEL_AUTO, nil
	case "FIRST":
		return dcgm.RSMI_DEV_PERF_LEVEL_FIRST, nil
	case "LOW":
		return dcgm.RSMI_DEV_PERF_LEVEL_LOW, nil
	case "HIGH":
		return dcgm.RSMI_DEV_PERF_LEVEL_HIGH, nil
	case "MANUAL":
		return dcgm.RSMI_DEV_PERF_LEVEL_MANUAL, nil
	case "STABLE_STD":
		return dcgm.RSMI_DEV_PERF_LEVEL_STABLE_STD, nil
	case "STABLE_PEAK":
		return dcgm.RSMI_DEV_PERF_LEVEL_STABLE_PEAK, nil
	case "STABLE_MIN_MCLK":
		return dcgm.RSMI_DEV_PERF_LEVEL_STABLE_MIN_MCLK, nil
	case "STABLE_MIN_SCLK":
		return dcgm.RSMI_DEV_PERF_LEVEL_STABLE_MIN_SCLK, nil
	case "DETERMINISM":
		return dcgm.RSMI_DEV_PERF_LEVEL_DETERMINISM, nil
	case "LAST":
		return dcgm.RSMI_DEV_PERF_LEVEL_LAST, nil
	case "UNKNOWN":
		return dcgm.RSMI_DEV_PERF_LEVEL_UNKNOWN, nil
	default:
		return dcgm.RSMI_DEV_PERF_LEVEL_UNKNOWN, fmt.Errorf("invalid level string: %s", level)
	}
}

// ConvertToRSMIGpuBlock 函数定义
func ConvertToRSMIGpuBlock(block string) (dcgm.RSMIGpuBlock, error) {
	switch block {
	case "INVALID":
		return dcgm.RSMIGpuBlockInvalid, nil
	case "FIRST":
		return dcgm.RSMIGpuBlockFirst, nil
	case "UMC":
		return dcgm.RSMIGpuBlockUMC, nil
	case "SDMA":
		return dcgm.RSMIGpuBlockSDMA, nil
	case "GFX":
		return dcgm.RSMIGpuBlockGFX, nil
	case "MMHUB":
		return dcgm.RSMIGpuBlockMMHUB, nil
	case "ATHUB":
		return dcgm.RSMIGpuBlockATHUB, nil
	case "PCIEBIF":
		return dcgm.RSMIGpuBlockPCIEBIF, nil
	case "HDP":
		return dcgm.RSMIGpuBlockHDP, nil
	case "HSLWAFL":
		return dcgm.RSMIGpuBlockHSLWAFL, nil
	case "DF":
		return dcgm.RSMIGpuBlockDF, nil
	case "SMN":
		return dcgm.RSMIGpuBlockSMN, nil
	case "SEM":
		return dcgm.RSMIGpuBlockSEM, nil
	case "MP0":
		return dcgm.RSMIGpuBlockMP0, nil
	case "MP1":
		return dcgm.RSMIGpuBlockMP1, nil
	case "FUSE":
		return dcgm.RSMIGpuBlockFuse, nil
	case "MCA":
		return dcgm.RSMIGpuBlockMCA, nil
	case "LAST":
		return dcgm.RSMIGpuBlockLast, nil
	case "RESERVED":
		return dcgm.RSMIGpuBlockReserved, nil
	default:
		return dcgm.RSMIGpuBlockInvalid, fmt.Errorf("invalid block string: %s", block)
	}
}

// ConvertToRSMISwComponent 函数定义
func ConvertToRSMISwComponent(component string) (dcgm.SwComponent, error) {
	switch component {
	case "FIRST":
		return dcgm.RSMISwCompFirst, nil
	case "DRIVER":
		return dcgm.RSMISwCompDriver, nil
	case "LAST":
		return dcgm.RSMISwCompLast, nil
	default:
		return dcgm.RSMISwCompFirst, fmt.Errorf("invalid component string: %s", component)
	}
}

// ConvertFrequencyToSclkClock 将sclk频率值转换为对应的十进制值
func ConvertFrequencyToSclkClock(freq string) (int64, error) {
	switch freq {
	case "600":
		return 1, nil
	case "700":
		return 2, nil
	case "750":
		return 4, nil
	case "800":
		return 8, nil
	case "900":
		return 16, nil
	case "1000":
		return 32, nil
	case "1106":
		return 64, nil
	case "1200":
		return 128, nil
	case "1270":
		return 256, nil
	case "1319":
		return 512, nil
	case "1400":
		return 1024, nil
	case "1500":
		return 2048, nil
	default:
		return 0, fmt.Errorf("invalid frequency value: %s", freq)
	}
}

// ConvertFrequencyToSocclkClock 将socclk频率值转换为对应的十进制值
func ConvertFrequencyToSocclkClock(freq string) (int64, error) {
	switch freq {
	case "309":
		return 1, nil
	case "523":
		return 2, nil
	case "566":
		return 4, nil
	case "618":
		return 8, nil
	case "680":
		return 16, nil
	case "755":
		return 32, nil
	case "850":
		return 64, nil
	case "971":
		return 128, nil
	default:
		return 0, fmt.Errorf("invalid frequency value: %s", freq)
	}
}

// ----------------------------
// 辅助：日志/结果文件存放位置
// ----------------------------

// getLogsDir 返回用于存放 diag 结果的 logs 目录路径（不会创建目录）。
// 逻辑：尝试在可执行文件所在目录及其上级目录查找 dcgm-hcu.log，若找到，取该目录作为基准；
// 否则回退到当前工作目录。
// 最终结果是 <基准目录>/logs
func getLogsDir() string {
	// 尝试基于可执行文件位置查找
	if exePath, err := os.Executable(); err == nil {
		dir := filepath.Dir(exePath)
		// 往上查两层（exe dir, parent, grandparent）
		for i := 0; i < 3; i++ {
			cand := filepath.Join(dir, "dcgm-hcu.log")
			if _, err := os.Stat(cand); err == nil {
				return filepath.Join(dir, "logs")
			}
			dir = filepath.Dir(dir)
		}
	}

	// 否则使用当前工作目录
	if wd, err := os.Getwd(); err == nil {
		return filepath.Join(wd, "logs")
	}

	// 最后兜底：相对路径 logs
	return "logs"
}

// ensureLogsDir 确保 logs 目录存在（会创建目录），返回目录路径或错误
func ensureLogsDir() (string, error) {
	d := getLogsDir()
	if err := os.MkdirAll(d, 0755); err != nil {
		return "", err
	}
	return d, nil
}

// saveJobToFile 将 job（包含 Result）序列化为 JSON 写入 logs/<job.ID>.json
// 注意：job 指针在并发场景下可能被修改，建议传入已经更新完成的 job 副本或在调用处做好同步。
func saveJobToFile(job *Job) error {
	dir, err := ensureLogsDir()
	if err != nil {
		return err
	}
	fpath := filepath.Join(dir, fmt.Sprintf("%s.json", job.ID))

	// 创建临时文件然后重命名，避免并发读到半写入文件
	tmp := fpath + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(job); err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return err
	}
	f.Close()
	if err := os.Rename(tmp, fpath); err != nil {
		return err
	}
	return nil
}

// loadJobFromFile 从 logs/<jobID>.json 反序列化 Job 并返回
func loadJobFromFile(jobID string) (*Job, error) {
	dir := getLogsDir()
	fpath := filepath.Join(dir, fmt.Sprintf("%s.json", jobID))
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var job Job
	dec := json.NewDecoder(f)
	if err := dec.Decode(&job); err != nil {
		return nil, err
	}
	return &job, nil
}

// -----------------------------
// Job / Job Store 类型与全局变量
// -----------------------------

// Job 状态常量
const (
	JobStatusPending  = "pending"
	JobStatusRunning  = "running"
	JobStatusDone     = "done"
	JobStatusCanceled = "canceled"
	JobStatusFailed   = "failed"
)

// jobStoreMu 保护 jobStore 的互斥锁
var jobStoreMu sync.Mutex

// jobStore 内存存储已提交的 job（key = job ID）
var jobStore = map[string]*Job{}

// jobIDCounter 原子计数器，用于生成唯一 job id 的一部分
var jobIDCounter uint64

// currentJobID 存储当前正在运行的 job id（用于 Stop 接口给出快速反馈）
// 使用 atomic.Value 以保证并发安全读取/写入。
var currentJobID atomic.Value

// 将 dcgm.DiagResults 转换为 API 层的 DiagResults
func convertDiagResults(src *dcgm.DiagResults) *DiagResults {
	if src == nil {
		return nil
	}

	// 转换每个 HCU 的诊断结果
	perHCU := make([]HCUResult, len(src.PerHCU))
	for i, d := range src.PerHCU {
		diagResults := make([]DiagResult, len(d.DiagResults))
		for j, r := range d.DiagResults {
			diagResults[j] = DiagResult{
				Status:       r.Status,
				TestName:     r.TestName,
				TestOutput:   r.TestOutput,
				ErrorCode:    r.ErrorCode,
				ErrorMessage: r.ErrorMessage,
			}
		}
		perHCU[i] = HCUResult{
			HCU:         d.HCU,
			RC:          d.RC,
			DiagResults: diagResults,
		}
	}

	// 转换软件层诊断结果
	software := make([]DiagResult, len(src.Software))
	for i, r := range src.Software {
		software[i] = DiagResult{
			Status:       r.Status,
			TestName:     r.TestName,
			TestOutput:   r.TestOutput,
			ErrorCode:    r.ErrorCode,
			ErrorMessage: r.ErrorMessage,
		}
	}

	return &DiagResults{
		DeviceNumber: src.DeviceNumber,
		Software:     software,
		PerHCU:       perHCU,
	}
}
