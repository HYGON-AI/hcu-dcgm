/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

// -------------------- 常量定义 --------------------

// EDPPLogDir 定义 EDPp 日志目录（相对可执行文件目录）
const EDPPLogDir = "logs/edpp"

// EDPpStressTime 定义每个模式的测试时间（秒）
const EDPpStressTime = 10

// -------------------- 数据结构 --------------------

// ErrorSummary 保存每个 HCU 的硬件错误统计（原有类型，文件中已有使用）
// 保留：用于内部复用 checkHardwareError 返回值
type ErrorSummary struct {
	HCU        int
	ECCErr     int
	MemErr     int
	ComputeErr int
}

// HCUStatus 用于记录 HCU 的实时采集信息（保留原有定义）
type HCUStatus struct {
	Utilization string
	Power       string
	GFXClock    string
	Temperature string
}

// -------------------- EDPp 可执行文件嵌入 --------------------

const (
	edppBackendLegacy    = "legacy"
	edppBackendNMZGFX938 = "nmz-gfx938"
	edppBackendBMZ       = "bmz"

	edppDevTypeIDGFX936 = "6320"
	edppDevTypeIDGFX938 = "6430"

	edppStressLoopCount = "1900000000"
	edppStressDuration  = 10 * time.Second

	// BW1000 空闲功耗约 86-105W，满载约 1000W+；设 300W 为最低门槛，
	// 压测启动后功耗应迅速超过该值，空闲状态绝对无法触及。
	edppBMZMinAvgPowerW = 300
	// BW1100/NMZ 实测空闲约 170W，EDPp 压力约 800W且不超过 1000W；
	// 以约 400W 作为进入计算状态的保守门槛，避免依赖 loader 不输出的 TFLOPS 文本。
	edppNMZMinAvgPowerW = 400
)

type edppBackendSpec struct {
	name         string
	loader       []byte
	loaderName   string
	co           []byte
	coName       string
	matrixArgs   []string
	timingFlag   string
	stressMode   bool
	minAvgPowerW int64
}

func edppBackendSpecFor(backend string) (edppBackendSpec, error) {
	switch backend {
	case edppBackendNMZGFX938:
		return edppBackendSpec{
			name:         edppBackendNMZGFX938,
			loader:       edppNMZLoaderBytes,
			loaderName:   "gemmPower_edpp",
			co:           edppNMZCOBytes,
			coName:       "fp8_nmz_edpp.co",
			matrixArgs:   []string{"-m", "2048", "-n", "2048", "-k", "12032", "-g", "22", "-t", "1"},
			timingFlag:   "-vs",
			stressMode:   true,
			minAvgPowerW: edppNMZMinAvgPowerW,
		}, nil
	case edppBackendBMZ:
		return edppBackendSpec{
			name:         edppBackendBMZ,
			loader:       edppBMZLoaderBytes,
			loaderName:   "gemmPower_edpp",
			co:           edppBMZCOBytes,
			coName:       "fp16_bmz_edpp.co",
			matrixArgs:   []string{"-m", "2048", "-n", "3840", "-k", "12032", "-g", "2", "-t", "1"},
			timingFlag:   "-vs",
			stressMode:   true,
			minAvgPowerW: edppBMZMinAvgPowerW,
		}, nil
	default:
		return edppBackendSpec{}, fmt.Errorf("unsupported EDPp backend %q", backend)
	}
}

//go:embed resources/EDPp
var edppBytes []byte

//go:embed resources/gemmpoweredpp/bmz/gemmPower_edpp
var edppBMZLoaderBytes []byte

//go:embed resources/gemmpoweredpp/bmz/fp16_bmz_edpp.co
var edppBMZCOBytes []byte

//go:embed resources/gemmpoweredpp/nmz/gemmPower_edpp
var edppNMZLoaderBytes []byte

//go:embed resources/gemmpoweredpp/nmz/fp8_nmz_edpp.co
var edppNMZCOBytes []byte

// extractEDPp 将嵌入的 EDPp 二进制写入临时文件，并返回可执行路径
func extractEDPp() (string, error) {
	if len(edppBytes) == 0 {
		return "", fmt.Errorf("embedded EDPp binary is empty")
	}

	tmpFile, err := os.CreateTemp("", "EDPp-*")
	if err != nil {
		return "", fmt.Errorf("无法创建临时文件: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(edppBytes); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("写入临时文件失败: %w", err)
	}

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("无法设置执行权限: %w", err)
	}

	return tmpFile.Name(), nil
}

func edppResourcesForBackend(backend string) (edppBackendSpec, error) {
	return edppBackendSpecFor(backend)
}

func writeEmbeddedFile(path string, data []byte, mode os.FileMode) error {
	if len(data) == 0 {
		return fmt.Errorf("embedded EDPp resource %s is empty", filepath.Base(path))
	}
	if err := os.WriteFile(path, data, mode); err != nil {
		return fmt.Errorf("write embedded EDPp resource %s: %w", filepath.Base(path), err)
	}
	if err := os.Chmod(path, mode); err != nil {
		return fmt.Errorf("set EDPp resource permissions %s: %w", filepath.Base(path), err)
	}
	return nil
}

func extractEdppResources(backend string) (loaderPath string, coPath string, cleanup func(), err error) {
	resources, err := edppResourcesForBackend(backend)
	if err != nil {
		return "", "", nil, err
	}

	tempDir, err := os.MkdirTemp("", "dcgm-edpp-*")
	if err != nil {
		return "", "", nil, fmt.Errorf("create EDPp temp directory: %w", err)
	}
	cleanup = func() { _ = os.RemoveAll(tempDir) }

	loaderPath = filepath.Join(tempDir, resources.loaderName)
	if err := writeEmbeddedFile(loaderPath, resources.loader, 0o755); err != nil {
		cleanup()
		return "", "", nil, err
	}
	coPath = filepath.Join(tempDir, resources.coName)
	if err := writeEmbeddedFile(coPath, resources.co, 0o644); err != nil {
		cleanup()
		return "", "", nil, err
	}

	return loaderPath, coPath, cleanup, nil
}

// GemmPowerEdppResult 表示新 EDPp 后端在单张 HCU 上的执行结果。
// BMZ/NMZ 都采用持续 GEMM 施压，并通过 RSMI 功耗/时钟采样判断 GPU 是否进入计算状态。
type GemmPowerEdppResult struct {
	HCU        int
	Backend    string
	TFLOPS     float64 // 保留兼容字段；当前 BMZ/NMZ stress 模式不依赖该值。
	AvgPowerW  int64
	PeakPowerW int64
	AvgGFXMHz  int64
	ECCDeltaCE int64 // 压测期间新增 correctable ECC 错误（各 block 合计 delta）
	ECCDeltaUE int64 // 压测期间新增 uncorrectable ECC 错误（各 block 合计 delta）
	Passed     bool
	Output     string
	Error      string
}

func (r GemmPowerEdppResult) Summary() string {
	status := "PASS"
	if !r.Passed {
		status = "FAIL"
	}
	s := fmt.Sprintf("Avg Power: %d W, Peak Power: %d W, GFX Clock: %d MHz, %s", r.AvgPowerW, r.PeakPowerW, r.AvgGFXMHz, status)
	if r.ECCDeltaCE > 0 || r.ECCDeltaUE > 0 {
		s += fmt.Sprintf(", ECC CE: %d, UE: %d", r.ECCDeltaCE, r.ECCDeltaUE)
	}
	return s
}

type limitedBuffer struct {
	mu    sync.Mutex
	buf   bytes.Buffer
	limit int
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.limit <= 0 || b.buf.Len() >= b.limit {
		return len(p), nil
	}
	remaining := b.limit - b.buf.Len()
	if len(p) > remaining {
		_, _ = b.buf.Write(p[:remaining])
		return len(p), nil
	}
	_, _ = b.buf.Write(p)
	return len(p), nil
}

func (b *limitedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.String()
}

func visibleDeviceEnv(hcu int) []string {
	filtered := make([]string, 0, len(os.Environ())+1)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "HIP_VISIBLE_DEVICES=") ||
			strings.HasPrefix(entry, "ROCR_VISIBLE_DEVICES=") ||
			strings.HasPrefix(entry, "GPU_DEVICE_ORDINAL=") ||
			strings.HasPrefix(entry, "CUDA_VISIBLE_DEVICES=") ||
			strings.HasPrefix(entry, "HSA_VISIBLE_DEVICES=") {
			continue
		}
		filtered = append(filtered, entry)
	}
	return append(filtered, fmt.Sprintf("HIP_VISIBLE_DEVICES=%d", hcu))
}

type edppStressSamples struct {
	totalPower   int64
	peakPower    int64
	totalGFX     int64
	powerSamples int64
	clockSamples int64
	lastPowerErr error
	lastClockErr error
}

// sumEccCounts 将 EccBlocksInfo 返回的所有 block CE/UE 合计。
// 失败时静默忽略（ECC 读取失败不阻断压测流程）。
func sumEccCounts(hcu int) (ce, ue int64) {
	blocks, err := EccBlocksInfo(hcu)
	if err != nil {
		return 0, 0
	}
	for _, b := range blocks {
		ce += b.CE
		ue += b.UE
	}
	return ce, ue
}

func finalizeStressEdppResult(result *GemmPowerEdppResult, spec edppBackendSpec, samples edppStressSamples) error {
	if samples.powerSamples > 0 {
		result.AvgPowerW = samples.totalPower / samples.powerSamples
		result.PeakPowerW = samples.peakPower
	}
	if samples.clockSamples > 0 {
		result.AvgGFXMHz = samples.totalGFX / samples.clockSamples
	}
	if samples.powerSamples == 0 {
		result.Error = "no valid power samples"
		if samples.lastPowerErr != nil {
			result.Error += ": " + samples.lastPowerErr.Error()
		}
		return fmt.Errorf("%s EDPp HCU %d: %s", spec.name, result.HCU, result.Error)
	}
	if result.AvgPowerW < spec.minAvgPowerW {
		result.Error = fmt.Sprintf("avg power %dW < threshold %dW, GPU may not be computing", result.AvgPowerW, spec.minAvgPowerW)
		if samples.lastClockErr != nil && samples.clockSamples == 0 {
			result.Error += "; no valid GFX clock samples: " + samples.lastClockErr.Error()
		}
		return fmt.Errorf("%s EDPp HCU %d: %s", spec.name, result.HCU, result.Error)
	}
	result.Passed = true
	return nil
}

// runStressEdpp 以长循环方式启动 gemmPower_edpp，并用 RSMI 采样功耗/GFX 时钟判断是否进入计算状态。
func runStressEdpp(spec edppBackendSpec, loaderPath, coPath string, hcu int) (GemmPowerEdppResult, error) {
	result := GemmPowerEdppResult{HCU: hcu, Backend: spec.name}
	output := &limitedBuffer{limit: 32 * 1024}

	cmd := exec.Command(loaderPath, gemmPowerEdppArgs(spec, coPath)...)
	cmd.Dir = filepath.Dir(loaderPath)
	cmd.Env = visibleDeviceEnv(hcu)
	cmd.Stdout = output
	cmd.Stderr = output
	if err := cmd.Start(); err != nil {
		result.Error = err.Error()
		return result, fmt.Errorf("%s stress start HCU %d: %w", spec.name, hcu, err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	processDone := false
	stopStress := func() {
		if cmd.Process != nil && !processDone {
			_ = cmd.Process.Kill()
			<-done
			processDone = true
		}
	}
	defer stopStress()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	deadline := time.NewTimer(edppStressDuration)
	defer deadline.Stop()

	// 压测前 ECC 基线快照（用于计算 delta）
	eccBaseCE, eccBaseUE := sumEccCounts(hcu)

	var samples edppStressSamples
	for {
		select {
		case waitErr := <-done:
			processDone = true
			result.Output = strings.TrimSpace(output.String())
			result.Error = "stress process exited before sampling completed"
			if waitErr != nil {
				result.Error += ": " + waitErr.Error()
			}
			return result, fmt.Errorf("%s EDPp HCU %d: %s", spec.name, hcu, result.Error)
		case <-deadline.C:
			stopStress()
			result.Output = strings.TrimSpace(output.String())
			if err := finalizeStressEdppResult(&result, spec, samples); err != nil {
				return result, err
			}
			// 压测后 ECC 快照，计算压测期间新增的错误数
			eccEndCE, eccEndUE := sumEccCounts(hcu)
			result.ECCDeltaCE = eccEndCE - eccBaseCE
			result.ECCDeltaUE = eccEndUE - eccBaseUE
			return result, nil
		case <-ticker.C:
			pw, err := Power(hcu)
			if err != nil {
				samples.lastPowerErr = err
			} else {
				samples.totalPower += pw
				samples.powerSamples++
				if pw > samples.peakPower {
					samples.peakPower = pw
				}
			}
			freqs, cur, err := GetClocksByType(hcu, RSMI_CLK_TYPE_SYS)
			if err != nil {
				samples.lastClockErr = err
			} else if int(cur) < len(freqs) {
				samples.totalGFX += int64(freqs[cur])
				samples.clockSamples++
			}
		}
	}
}

func detectEdppBackendBySeries(seriesName string) string {
	series := strings.ToUpper(strings.TrimSpace(seriesName))
	switch {
	case strings.Contains(series, "BW1100"),
		strings.Contains(series, "BW1102"),
		strings.Contains(series, "BW1200"):
		return edppBackendNMZGFX938
	case strings.Contains(series, "BW100"),
		strings.Contains(series, "BW101"),
		strings.Contains(series, "BW150"),
		strings.Contains(series, "BW151"),
		strings.Contains(series, "BW200"),
		strings.Contains(series, "BW1000"):
		return edppBackendBMZ
	default:
		return edppBackendLegacy
	}
}

func detectEdppBackendByDevTypeID(devTypeID string) string {
	id := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(devTypeID)), "0x")
	switch id {
	case edppDevTypeIDGFX938:
		return edppBackendNMZGFX938
	case edppDevTypeIDGFX936:
		return edppBackendBMZ
	default:
		return edppBackendLegacy
	}
}

func detectEdppBackend(cards []CardSeriesInfo, hcu int) string {
	return detectEdppBackendWithDevTypeID(cards, hcu, DevTypeID)
}

func detectEdppBackendWithDevTypeID(cards []CardSeriesInfo, hcu int, devTypeID func(int) (string, error)) string {
	for _, card := range cards {
		if card.DvInd != hcu {
			continue
		}
		if backend := detectEdppBackendBySeries(card.SeriesName); backend != edppBackendLegacy {
			return backend
		}
		break
	}

	if devTypeID == nil {
		return edppBackendLegacy
	}
	id, err := devTypeID(hcu)
	if err != nil {
		glog.Warningf("EDPp backend fallback DevTypeID failed, hcu=%d, err=%v", hcu, err)
		return edppBackendLegacy
	}
	return detectEdppBackendByDevTypeID(id)
}

func gemmPowerEdppArgs(spec edppBackendSpec, coPath string) []string {
	args := append([]string{}, spec.matrixArgs...)
	loopCount := "4"
	if spec.stressMode {
		loopCount = edppStressLoopCount
	}
	return append(args, "-f", coPath, "-w", "-l", loopCount, spec.timingFlag)
}

// -------------------- 硬件错误统计 --------------------

// checkHardwareError 按模式名读取每个 HCU 日志并统计错误（原有函数，保留不变）
func checkHardwareError(logDir string, hcus []int, pattern string) ([]ErrorSummary, error) {
	var results []ErrorSummary
	for _, hcu := range hcus {
		logFile := filepath.Join(logDir, fmt.Sprintf("edpp%s_hcu%d.log", pattern, hcu))
		data, err := os.ReadFile(logFile)
		if err != nil {
			return nil, fmt.Errorf("无法读取日志 %s: %w", logFile, err)
		}

		content := string(data)
		summary := ErrorSummary{HCU: hcu}
		summary.ECCErr = strings.Count(content, "ECC error")
		summary.MemErr = strings.Count(content, "Memory error")
		summary.ComputeErr = strings.Count(content, "Compute error")

		results = append(results, summary)
	}
	return results, nil
}

// -------------------- 实时监控 goroutine --------------------

// monitorEdpp 采集指定 HCU 的实时信息并写入日志文件（保留 legacy 日志格式）
func monitorEdpp(logFile string, hcu int, stop chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("无法打开日志文件 %s: %v\n", logFile, err)
		return
	}
	defer f.Close()

	// 写入 CSV 表头
	fmt.Fprintln(f, "index,timestamp,utilization.hcu [%],power.draw [W],clocks.current.gfx [MHz],temperature.hcu")

	for {
		select {
		case <-stop:
			return
		default:
			now := time.Now()
			timestamp := now.Format("2006-01-02 15:04:05")

			util, power, temp, gfx := "0", "0", "0", "0"

			// hy-smi 命令获取利用率、功耗、温度
			cmdTemp := exec.Command("sh", "-c",
				fmt.Sprintf("hy-smi -d %d | grep 'Normal' | awk '{print $7, $3, $2}'", hcu))
			out, _ := cmdTemp.Output()
			fields := strings.Fields(string(out))
			if len(fields) >= 3 {
				util, power, temp = fields[0], fields[1], fields[2]
			}

			// hy-smi 命令获取 GFX Clock
			cmdGfx := exec.Command("sh", "-c",
				fmt.Sprintf("hy-smi -d %d -c | grep sclk | awk -F'[()]' '{print $2}'", hcu))
			out, _ = cmdGfx.Output()
			if len(out) > 0 {
				gfx = strings.TrimSpace(string(out))
			}

			fmt.Fprintf(f, "%d,%s,%s,%s,%s,%s\n", hcu, timestamp, util, power, gfx, temp)
			time.Sleep(1 * time.Second) // 采集间隔
		}
	}
}

func appendEdppLog(logFile string, output []byte) error {
	if len(output) == 0 {
		return nil
	}
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(output); err != nil {
		return err
	}
	if output[len(output)-1] != '\n' {
		_, err = f.WriteString("\n")
	}
	return err
}

// -------------------- 测试模式列表 --------------------

var edppNames = []string{"10KHz", "5KHz", "4KHz", "3KHz", "2KHz", "1.5KHz", "1KHz", "800Hz", "500Hz", "200Hz", "100Hz", "50Hz", "100ms"}

// -------------------- EDPp 测试主函数（原有、不变） --------------------

// RunEDPpTest 主函数，运行指定 HCU 列表的 EDPp 测试
// dvIdList: 要测试的 HCU ID 列表
func edpppTest() {
	result, err := runEdppTestWithResult()
	for _, hcuResult := range result.HCUEdppResults {
		fmt.Printf("[HCU %d] backend=%s\n", hcuResult.HCUId, hcuResult.Backend)
		for _, pattern := range hcuResult.PatternResults {
			fmt.Printf("  Pattern=%s, ECC=%d, Mem=%d, Compute=%d\n", pattern.PatternName, pattern.ECCCount, pattern.MemoryErrorCount, pattern.ComputeErrorCount)
		}
		for _, gemmPower := range hcuResult.GemmPowerResults {
			fmt.Printf("  %s\n", gemmPower.Summary())
			if gemmPower.Error != "" {
				fmt.Printf("  Error: %s\n", gemmPower.Error)
			}
		}
	}
	if err != nil {
		fmt.Printf("EDPp test failed: %v\n", err)
	}
}

// -------------------- 新增：结构化 API（保留原实现不变） --------------------

// PatternResult 表示某个测试模式下某个 HCU 的错误统计（字段名规范）
type PatternResult struct {
	PatternName       string // 测试模式名称（例如 "10KHz", "5KHz" 等）
	ECCCount          int    // 该模式下检测到的 ECC 错误数
	MemoryErrorCount  int    // 该模式下检测到的内存错误数
	ComputeErrorCount int    // 该模式下检测到的计算相关错误数
}

// HCUEdppResult 表示某个 HCU 的 EDPp 结果。旧后端填充 PatternResults，
// gemmPower 后端填充 GemmPowerResults，两个字段保持独立语义。
type HCUEdppResult struct {
	HCUId            int
	Backend          string
	PatternResults   []PatternResult
	GemmPowerResults []GemmPowerEdppResult
}

// EDPPResult 汇总整个 EDPp 测试的结构化结果。
type EDPPResult struct {
	HCUEdppResults []HCUEdppResult
	LogDir         string
}

func runLegacyEdppWithResult(dvIdList []int) (map[int][]PatternResult, error) {
	results := make(map[int][]PatternResult, len(dvIdList))
	for _, hcu := range dvIdList {
		results[hcu] = make([]PatternResult, 0, len(edppNames))
	}
	if len(dvIdList) == 0 {
		return results, nil
	}

	edppPath, err := extractEDPp()
	if err != nil {
		return results, fmt.Errorf("提取 legacy EDPp 失败: %w", err)
	}
	defer os.Remove(edppPath)

	for patternIndex, name := range edppNames {
		for _, hcu := range dvIdList {
			logFile := filepath.Join(EDPPLogDir, fmt.Sprintf("edpp%s_hcu%d.log", name, hcu))
			stop := make(chan struct{})
			var wg sync.WaitGroup
			wg.Add(1)
			go monitorEdpp(logFile, hcu, stop, &wg)

			cmd := exec.Command(edppPath,
				"-d", strconv.Itoa(hcu),
				"-t", strconv.Itoa(EDPpStressTime),
				"-f", strconv.Itoa(patternIndex),
			)
			output, runErr := cmd.CombinedOutput()
			close(stop)
			wg.Wait()
			if logErr := appendEdppLog(logFile, output); logErr != nil {
				return results, fmt.Errorf("写入 legacy EDPp HCU %d pattern %s 日志失败: %w", hcu, name, logErr)
			}
			if runErr != nil {
				return results, fmt.Errorf("legacy EDPp HCU %d pattern %s failed: %w, output=%s", hcu, name, runErr, strings.TrimSpace(string(output)))
			}
		}

		hardwareErrors, err := checkHardwareError(EDPPLogDir, dvIdList, name)
		if err != nil {
			return results, fmt.Errorf("检查 legacy EDPp pattern %s 硬件错误失败: %w", name, err)
		}
		for _, hardwareError := range hardwareErrors {
			results[hardwareError.HCU] = append(results[hardwareError.HCU], PatternResult{
				PatternName:       name,
				ECCCount:          hardwareError.ECCErr,
				MemoryErrorCount:  hardwareError.MemErr,
				ComputeErrorCount: hardwareError.ComputeErr,
			})
		}
	}

	return results, nil
}

func runGemmPowerEdppWithResult(backend string, dvIdList []int) (map[int][]GemmPowerEdppResult, error) {
	results := make(map[int][]GemmPowerEdppResult, len(dvIdList))
	if len(dvIdList) == 0 {
		return results, nil
	}

	loaderPath, coPath, cleanup, err := extractEdppResources(backend)
	if err != nil {
		return results, err
	}
	defer cleanup()

	var executionErrors []string
	for _, hcu := range dvIdList {
		spec, err := edppBackendSpecFor(backend)
		if err != nil {
			return results, err
		}
		result, runErr := runStressEdpp(spec, loaderPath, coPath, hcu)

		results[hcu] = []GemmPowerEdppResult{result}
		logPath := filepath.Join(EDPPLogDir, fmt.Sprintf("gemmPower_edpp_%s_hcu%d.log", backend, hcu))
		logContent := fmt.Sprintf("avgPower=%dW peakPower=%dW avgGFX=%dMHz passed=%t err=%s\n%s",
			result.AvgPowerW, result.PeakPowerW, result.AvgGFXMHz, result.Passed, result.Error, result.Output)
		if writeErr := os.WriteFile(logPath, []byte(logContent+"\n"), 0o644); writeErr != nil {
			executionErrors = append(executionErrors, fmt.Sprintf("HCU %d 写日志失败: %v", hcu, writeErr))
		}
		if runErr != nil {
			executionErrors = append(executionErrors, fmt.Sprintf("HCU %d: %v", hcu, runErr))
		}
	}
	if len(executionErrors) > 0 {
		return results, fmt.Errorf("%s EDPp failed: %s", backend, strings.Join(executionErrors, "; "))
	}
	return results, nil
}

func edppTargetDevices(totalHCU int, dvIdList []int) ([]int, error) {
	if len(dvIdList) == 0 {
		devices := make([]int, 0, totalHCU)
		for hcu := range totalHCU {
			devices = append(devices, hcu)
		}
		return devices, nil
	}

	seen := make(map[int]bool, len(dvIdList))
	devices := make([]int, 0, len(dvIdList))
	for _, hcu := range dvIdList {
		if hcu < 0 || hcu >= totalHCU {
			return nil, fmt.Errorf("invalid HCU ID %d: valid range is 0-%d", hcu, totalHCU-1)
		}
		if seen[hcu] {
			continue
		}
		seen[hcu] = true
		devices = append(devices, hcu)
	}
	return devices, nil
}

func runEdppTestWithResult(dvIdList ...int) (EDPPResult, error) {
	result := EDPPResult{LogDir: EDPPLogDir}
	totalHCU, err := rsmiNumMonitorDevices()
	if err != nil {
		return result, fmt.Errorf("获取 HCU 总数失败: %w", err)
	}
	if totalHCU <= 0 {
		return result, fmt.Errorf("获取 HCU 总数失败: device count=%d", totalHCU)
	}
	targetDevices, err := edppTargetDevices(totalHCU, dvIdList)
	if err != nil {
		return result, err
	}
	if err := os.MkdirAll(EDPPLogDir, 0o755); err != nil {
		return result, fmt.Errorf("无法创建 EDPp 日志目录: %w", err)
	}

	cards, cardErr := CardSeriesList()
	legacyDevices := make([]int, 0, len(targetDevices))
	nmzDevices := make([]int, 0, len(targetDevices))
	bmzDevices := make([]int, 0, len(targetDevices))
	backends := make(map[int]string, len(targetDevices))
	for _, hcu := range targetDevices {
		backend := edppBackendLegacy
		if cardErr == nil {
			backend = detectEdppBackend(cards, hcu)
		}
		backends[hcu] = backend
		switch backend {
		case edppBackendNMZGFX938:
			nmzDevices = append(nmzDevices, hcu)
		case edppBackendBMZ:
			bmzDevices = append(bmzDevices, hcu)
		default:
			legacyDevices = append(legacyDevices, hcu)
		}
	}

	legacyResults, legacyErr := runLegacyEdppWithResult(legacyDevices)
	nmzResults, nmzErr := runGemmPowerEdppWithResult(edppBackendNMZGFX938, nmzDevices)
	bmzResults, bmzErr := runGemmPowerEdppWithResult(edppBackendBMZ, bmzDevices)
	for _, hcu := range targetDevices {
		gemmPowerResults := nmzResults[hcu]
		if gemmPowerResults == nil {
			gemmPowerResults = bmzResults[hcu]
		}
		result.HCUEdppResults = append(result.HCUEdppResults, HCUEdppResult{
			HCUId:            hcu,
			Backend:          backends[hcu],
			PatternResults:   legacyResults[hcu],
			GemmPowerResults: gemmPowerResults,
		})
	}

	var executionErrors []string
	if legacyErr != nil {
		executionErrors = append(executionErrors, legacyErr.Error())
	}
	if nmzErr != nil {
		executionErrors = append(executionErrors, nmzErr.Error())
	}
	if bmzErr != nil {
		executionErrors = append(executionErrors, bmzErr.Error())
	}
	if len(executionErrors) > 0 {
		return result, fmt.Errorf("EDPp execution failed: %s", strings.Join(executionErrors, "; "))
	}
	return result, nil
}
