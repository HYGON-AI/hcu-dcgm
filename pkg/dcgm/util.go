/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L/opt/hyhal/lib -Wl,-rpath,/opt/hyhal/lib -lrocm_smi64 -lhydmi -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdint.h>
#include <kfd_ioctl.h>
#include <rocm_smi64Config.h>
#include <rocm_smi.h>
#include <dmi_virtual.h>
#include <dmi_error.h>
#include <dmi.h>
#include <dmi_mig.h>
*/
import "C"
import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/golang/glog"
)

func errorString(result C.rsmi_status_t) error {
	if RSMIStatus(result) == RSMI_STATUS_SUCCESS {
		return nil
	}
	var cStatusString *C.char
	statusCode := C.rsmi_status_string(result, (**C.char)(unsafe.Pointer(&cStatusString)))
	if RSMIStatus(statusCode) != RSMI_STATUS_SUCCESS {
		return fmt.Errorf("error: %d", statusCode)
	}
	goStatusString := C.GoString(cStatusString)
	return fmt.Errorf("%s", goStatusString)
}

func dmiErrorString(result C.dmiStatus) error {
	if DMIStatus(result) == DMI_STATUS_SUCCESS {
		return nil
	}
	var cStatusString *C.char
	statusCode := C.dmiGetStatusString(result, (**C.char)(unsafe.Pointer(&cStatusString)))
	if DMIStatus(statusCode) != DMI_STATUS_SUCCESS {
		return fmt.Errorf("error: %d", statusCode)
	}
	goStatusString := C.GoString(cStatusString)
	return fmt.Errorf("%s", goStatusString)
}
func migErrorString(result C.nvmlReturn_t) error {
	code := NvmlReturn(result)
	if code == NVML_SUCCESS {
		return nil
	}
	if msg, ok := nvmlErrorCodeMap[code]; ok {
		return fmt.Errorf("MIG error: %s (code=%d)", msg, int(code))
	}
	return fmt.Errorf("MIG error: UNKNOWN_ERROR_CODE_%d", int(code))
}

// 获取所提供的RSMI错误状态的描述
func go_rsmi_status_string(status RSMIStatus) (statusStr string, err error) {
	var cstatusStr *C.char
	ret := C.rsmi_status_string(C.rsmi_status_t(status), (**C.char)(unsafe.Pointer(&cstatusStr)))
	if err = errorString(ret); err != nil {
		return statusStr, fmt.Errorf("Error go_rsmi_status_string:%s", err)
	}
	statusStr = C.GoString(cstatusStr)
	return
}
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

func perfLevelString(i int) string {
	switch i {
	case 0:
		return "AUTO"
	case 1:
		return "LOW"
	case 2:
		return "HIGH"
	case 3:
		return "MANUAL"
	case 4:
		return "STABLE_STD"
	case 5:
		return "STABLE_PEAK"
	case 6:
		return "STABLE_MIN_MCLK"
	case 7:
		return "STABLE_MIN_SCLK"
	default:
		return "UNKNOWN"
	}
}

func ConvertASCIIToString(asciiCodes []byte) string {
	var result []rune
	for _, code := range asciiCodes {
		// Stop at the first null character
		if code == 0 {
			break
		}
		// Filter out non-ASCII characters
		if code > 127 {
			continue
		}
		result = append(result, rune(code))
	}
	return string(result)
}

// parseConfig 解析配置文件内容为DMIVDeviceInfo结构体
func parseConfig(filePath string) (*VDeviceInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &VDeviceInfo{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) < 2 {
			continue
		}
		key := parts[0]
		value := parts[1]

		switch key {
		case "cu_count":
			config.VComputeUnitCount, _ = strconv.Atoi(value)
		case "mem":
			// 解析内存大小，例如 "4096 MiB"
			memParts := strings.Fields(value)
			if len(memParts) == 2 {
				memSize, err := strconv.Atoi(memParts[0])
				if err == nil {
					// 转换为字节数（假设单位是 MiB）
					config.VMemoryTotal = uintptr(memSize * 1024 * 1024)
				}
			}
		case "device_id":
			config.DvInd, _ = strconv.Atoi(value)
		case "vdev_id":
			config.VdvInd, _ = strconv.Atoi(value)
		case "PciBusId":
			config.PciBusNumber = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

// 打印二维数组
func print2DArray(data [][]string) {
	for _, row := range data {
		fmt.Println(strings.Join(row, "\t"))
	}
}

// 打印超出规格运行的警告，并提示用户接受条款
func confirmOutOfSpecWarning(autoRespond bool) {
	warning := `
          ******WARNING******

          Operating your HCU device outside of official AMD specifications or outside of
          factory settings, including but not limited to the conducting of overclocking,
          over-volting or under-volting (including use of this interface software,
          even if such software has been directly or indirectly provided by AMD or otherwise
          affiliated in any way with AMD), may cause damage to your HCU device, system components
          and/or result in system failure, as well as cause other problems.
          DAMAGES CAUSED BY USE OF YOUR HCU device OUTSIDE OF OFFICIAL AMD SPECIFICATIONS OR
          OUTSIDE OF FACTORY SETTINGS ARE NOT COVERED UNDER ANY AMD PRODUCT WARRANTY AND
          MAY NOT BE COVERED BY YOUR BOARD OR SYSTEM MANUFACTURER'S WARRANTY.
          Please use this utility with caution.
          `

	fmt.Println(warning)

	var userInput string
	if !autoRespond {
		fmt.Print("Do you accept these terms? [y/N] ")
		fmt.Scanln(&userInput)
	} else {
		userInput = "y"
	}

	userInput = strings.ToLower(userInput)
	if userInput == "yes" || userInput == "y" {
		return
	} else {
		fmt.Println("Confirmation not given. Exiting without setting value")
		os.Exit(1)
	}
}
func profileString(profile interface{}) string {
	dictionary := map[int]string{
		1:  "CUSTOM",
		2:  "VIDEO",
		4:  "POWER SAVING",
		8:  "COMPUTE",
		16: "VR",
		32: "3D FULL SCREEN",
		64: "BOOTUP DEFAULT",
	}

	switch v := profile.(type) {
	case int:
		if name, ok := dictionary[v]; ok {
			return name
		}
	case string:
		if num, err := strconv.Atoi(v); err == nil {
			if name, ok := dictionary[num]; ok {
				return name
			}
		} else {
			for key, val := range dictionary {
				if val == v {
					return strconv.Itoa(key)
				}
			}
		}
	}
	return "UNKNOWN"
}

func profileEnum(profile string) PowerProfilePresetMasks {
	dictionary := map[string]PowerProfilePresetMasks{
		"CUSTOM":         RSMI_PWR_PROF_PRST_CUSTOM_MASK,
		"VIDEO":          RSMI_PWR_PROF_PRST_VIDEO_MASK,
		"POWER SAVING":   RSMI_PWR_PROF_PRST_POWER_SAVING_MASK,
		"COMPUTE":        RSMI_PWR_PROF_PRST_COMPUTE_MASK,
		"VR":             RSMI_PWR_PROF_PRST_VR_MASK,
		"3D FULL SCREEN": RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK,
		"BOOTUP DEFAULT": RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT,
	}

	if val, ok := dictionary[profile]; ok {
		return val
	}
	return RSMI_PWR_PROF_PRST_INVALID
}

func printTableLog(headers []string, data [][]string, device int, title string) {
	fmt.Printf("Device: %d - %s\n", device, title)
	fmt.Println(headers)
	for _, row := range data {
		fmt.Println(row)
	}
	fmt.Println()
}

func formatMatrixToJSON(deviceList []int, matrix [][]int64, metricName string) {
	for rowIndx := 0; rowIndx < len(deviceList); rowIndx++ {
		for colInd := rowIndx + 1; colInd < len(deviceList); colInd++ {
			valueStr := matrix[deviceList[rowIndx]][deviceList[colInd]]
			fmt.Printf(metricName+"\n", deviceList[rowIndx], deviceList[colInd])
			fmt.Println(valueStr)
		}
	}
}

func formatMatrixToStrJSON(deviceList []int, matrix [][]string, metricName string) {
	for rowIndx := 0; rowIndx < len(deviceList); rowIndx++ {
		for colInd := rowIndx + 1; colInd < len(deviceList); colInd++ {
			valueStr := matrix[deviceList[rowIndx]][deviceList[colInd]]
			fmt.Printf(metricName+"\n", deviceList[rowIndx], deviceList[colInd])
			fmt.Println(valueStr)
		}
	}
}

func printTableRow(format string, displayString interface{}) {
	if format != "" {
		fmt.Printf(format, displayString)
	} else {
		fmt.Print(displayString)
	}
	fmt.Print(" ")
}

// 获取指定目录下的文件列表，如果目录不存在或为空，返回空切片
func getConfigFiles(dir string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果目录不存在，返回空切片
			return []os.DirEntry{}, nil
		}
		return nil, err
	}
	return files, nil
}

// 解析配置文件内容
func parseConfigFile(filePath string) (map[string]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	config := make(map[string]string)
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return config, nil
}

// 执行并行任务
func executeInParallel(wg *sync.WaitGroup, task func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		task()
	}()
}

// 用安全锁更新monitorInfo的字段值
func updateMonitorInfo(mu *sync.Mutex, updateFunc func()) {
	mu.Lock()
	defer mu.Unlock()
	updateFunc()
}

// Helper function to perform a task and update monitorInfo
func fetchAndUpdate[T any](mu *sync.Mutex, wg *sync.WaitGroup, fetchFunc func() T, updateFunc func(T)) {
	executeInParallel(wg, func() {
		result := fetchFunc()
		updateMonitorInfo(mu, func() {
			updateFunc(result)
		})
	})
}

var blockToStringMap = map[RSMIGpuBlock]string{
	RSMIGpuBlockInvalid:  "INVALID",  // 无效模块，用于初始化或表示错误状态。
	RSMIGpuBlockUMC:      "UMC",      // 统一内存控制器（Unified Memory Controller），管理 GPU 内存的分配和访问。
	RSMIGpuBlockSDMA:     "SDMA",     // 单指令多数据引擎（Single Data Multiple Access），处理 GPU 和 CPU 之间的内存传输。
	RSMIGpuBlockGFX:      "GFX",      // 图形处理模块（Graphics Processing Unit），负责图形渲染和计算。
	RSMIGpuBlockMMHUB:    "MMHUB",    // 内存管理中心（Memory Management Hub），协调内存请求。
	RSMIGpuBlockATHUB:    "ATHUB",    // 高速互连模块（AT Hub），连接 GPU 与其他设备的通信。
	RSMIGpuBlockPCIEBIF:  "PCIEBIF",  // PCI Express 总线接口（PCIe Bus Interface），管理 PCIe 数据传输。
	RSMIGpuBlockHDP:      "HDP",      // 主机数据路径（Host Data Path），处理 GPU 和 CPU 的通信。
	RSMIGpuBlockHSLWAFL:  "HSLWAFL",  // 高速 GPU 间互连（HSL），支持多 GPU 间快速通信。
	RSMIGpuBlockDF:       "DF",       // 数据结构控制器（Data Fabric），管理 CPU 和 GPU 数据流。
	RSMIGpuBlockSMN:      "SMN",      // 系统管理网络（System Management Network），提供硬件监控功能。
	RSMIGpuBlockSEM:      "SEM",      // 安全管理模块（Security Management Module），负责硬件安全。
	RSMIGpuBlockMP0:      "MP0",      // 微处理器 0（Microprocessor 0），用于低功耗状态的管理。
	RSMIGpuBlockMP1:      "MP1",      // 微处理器 1（Microprocessor 1），支持高性能的硬件管理。
	RSMIGpuBlockFuse:     "FUSE",     // 熔丝模块（Fuse Block），负责硬件配置和特性设置。
	RSMIGpuBlockMCA:      "MCA",      // 机器检查架构（Machine Check Architecture），检测硬件错误。
	RSMIGpuBlockReserved: "RESERVED", // 预留值，不用于实际模块。
}

func ConvertFromRSMIGpuBlock(block RSMIGpuBlock) string {
	if str, exists := blockToStringMap[block]; exists {
		return str
	}
	return "UNKNOWN"
}

func listFilesInDevDri() (int, error) {
	foundCounter := 0
	baseDir := "/sys/devices"

	if _, err := os.Stat(baseDir); err != nil {
		return foundCounter, fmt.Errorf("基础目录 %s 不存在或无法访问: %v", baseDir, err)
	}

	if err := processDir(baseDir, &foundCounter); err != nil {
		glog.Errorf("处理目录失败: %v", err)
		return foundCounter, fmt.Errorf("处理目录 %s 失败: %v", baseDir, err)
	}
	return foundCounter, nil
}

func processDir(dirPath string, foundCounter *int) error {
	// 先 stat 判断是否为目录（跟随 symlink）
	fi, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("无法 stat 路径 %s: %v", dirPath, err)
	}
	if !fi.Mode().IsDir() {
		// 如果不是目录，但是 device 文件则处理之（兼容误传）
		if filepath.Base(dirPath) == "device" {
			return processDeviceFile(dirPath, foundCounter)
		}
		// 不是目录也不是 device，忽略
		glog.V(4).Infof("跳过非目录路径 %s", dirPath)
		return nil
	}

	// 打开目录并读取项名
	dir, err := os.Open(dirPath)
	if err != nil {
		return fmt.Errorf("无法打开目录 %s: %v", dirPath, err)
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("读取目录 %s 文件失败: %v", dirPath, err)
	}

	var errorsAcc []string
	for _, name := range names {
		fullPath := filepath.Join(dirPath, name)

		// 处理以 pci0000 开头的目录（递归）
		if strings.HasPrefix(name, "pci0000") {
			// 先 stat 判断是否目录
			if fi2, err := os.Stat(fullPath); err == nil && fi2.Mode().IsDir() {
				if err := processDir(fullPath, foundCounter); err != nil {
					glog.Warningf("处理目录 %s 失败: %v", fullPath, err)
					errorsAcc = append(errorsAcc, fmt.Sprintf("处理目录 %s 失败: %v", fullPath, err))
				}
			} else {
				glog.V(4).Infof("跳过非目录或无法 stat 的 %s", fullPath)
			}
			continue
		}

		// 处理以 0000 开头的项（可能是目录或 symlink）
		if strings.HasPrefix(name, "0000") {
			// stat 跟随 symlink 决定是否递归
			if fi2, err := os.Stat(fullPath); err == nil && fi2.Mode().IsDir() {
				if err := process0000Dir(fullPath, foundCounter); err != nil {
					glog.Warningf("处理目录 %s 失败: %v", fullPath, err)
					errorsAcc = append(errorsAcc, fmt.Sprintf("处理目录 %s 失败: %v", fullPath, err))
				}
			} else {
				glog.V(4).Infof("跳过非目录的 0000 项 %s", fullPath)
			}
			continue
		}
		// 其它项忽略
	}

	if len(errorsAcc) > 0 {
		return fmt.Errorf("处理目录 %s 时遇到错误:\n%s", dirPath, strings.Join(errorsAcc, "\n"))
	}
	return nil
}

// 单独处理 device 文件的函数（便于复用）
func processDeviceFile(fullPath string, foundCounter *int) error {
	//if err := loadPCIDeviceNames(); err != nil {
	//	glog.Errorf("loadPCIDeviceNames failed: %v", err)
	//}
	if err := loadUpdateIDsMap(); err != nil {
		glog.Warningf("loadUpdateIDsMap failed: %v", err)
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("无法读取 device 文件 %s: %v", fullPath, err)
	}
	deviceValue := strings.TrimSpace(string(data))
	if len(deviceValue) >= 2 && (strings.HasPrefix(deviceValue, "0x") || strings.HasPrefix(deviceValue, "0X")) {
		deviceValue = deviceValue[2:]
	}
	deviceValue = strings.ToLower(deviceValue)

	//modelName, found := type2name[deviceValue]
	//modelName, found := pciDeviceID2Name[deviceValue]
	// 查 update_ids 优先
	modelName, found := updateIDsMap[deviceValue]
	// update_ids 没有，再用手动维护的 type2name
	if !found {
		modelName, found = type2name[deviceValue]
	}

	if !found {
		modelName = "未知型号"
	}

	glog.V(5).Infof("device 文件: %s, 值: %s, 型号: %s", fullPath, deviceValue, modelName)
	if modelName != "未知型号" {
		(*foundCounter)++
	}
	return nil
}

func process0000Dir(dirPath string, foundCounter *int) error {
	// 判断是否目录或是 device 文件
	fi, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("无法 stat %s: %v", dirPath, err)
	}

	// 不是目录是 device 文件，直接处理
	if !fi.Mode().IsDir() {
		if filepath.Base(dirPath) == "device" {
			return processDeviceFile(dirPath, foundCounter)
		}
		// 否则忽略
		glog.V(4).Infof("process0000Dir: 跳过非目录路径 %s", dirPath)
		return nil
	}

	// dirPath 确认是目录
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("无法打开目录 %s: %v", dirPath, err)
	}

	var errs []string
	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(dirPath, name)

		// 遇到以 "0000" 开头的项，递归（需 stat 跟随 symlink）
		if strings.HasPrefix(name, "0000") {
			if fi2, err := os.Stat(fullPath); err == nil && fi2.Mode().IsDir() {
				if err := process0000Dir(fullPath, foundCounter); err != nil {
					glog.Warningf("处理目录 %s 失败: %v", fullPath, err)
					errs = append(errs, fmt.Sprintf("递归 %s: %v", fullPath, err))
				}
			} else {
				glog.V(4).Infof("process0000Dir: 跳过非目录项 %s", fullPath)
			}
			continue
		}

		// 处理 device 文件
		if name == "device" {
			if err := processDeviceFile(fullPath, foundCounter); err != nil {
				glog.Warningf("无法打开文件 %s: %v", fullPath, err)
				errs = append(errs, fmt.Sprintf("打开 %s: %v", fullPath, err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("处理 %s 时出现 %d 个错误: %s", dirPath, len(errs), strings.Join(errs, "; "))
	}
	return nil
}

func getDTKVersionByReadFile() (string, error) {
	filePath := "/opt/dtk/.info/version-dev"
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var result string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	glog.V(5).Infof("result: ", result)
	return result, nil
}

// int 转 []byte (固定转为 int64)
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(int64(v)))
	return b
}

// []byte 转 int (固定转为 int64 再转 int)
func btoi(b []byte) int {
	if b == nil || len(b) < 8 {
		return 0
	}
	return int(int64(binary.BigEndian.Uint64(b)))
}

func bytesToGB(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}

// extractInt32Array 从 C 的 uint32_t* 复制出一个 Go 的 []int32 副本。
func extractInt32Array(ptr *C.uint32_t, length int) []int32 {
	if ptr == nil || length <= 0 {
		return nil
	}
	const maxLen = 1 << 24
	if length > maxLen {
		return nil
	}
	carr := (*[1 << 28]C.uint32_t)(unsafe.Pointer(ptr))[:length:length]
	out := make([]int32, length)
	for i, v := range carr {
		out[i] = int32(v)
	}
	return out
}

// 从 C 数组提取 float32 切片
func extractFloat32Array(ptr *C.float, length int) []float32 {
	if ptr == nil || length <= 0 {
		return nil
	}
	slice := make([]float32, length)
	for i := 0; i < length; i++ {
		slice[i] = float32(*(*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(i)*unsafe.Sizeof(C.float(0)))))
	}
	return slice
}

// CheckLive 执行两个命令，并检查输出是否包含 "live"。
// 如果两个命令的输出均为空或都不包含 "live"，则等待 10 秒后重试，最多尝试 6 次。
// 如果在任意一次中检测到 "live"，则立即结束并返回 true，否则返回 false。
func checkDriverLive() bool {
	cmd1 := "cat /sys/module/hycu/initstate"
	cmd2 := "cat /sys/module/hydcu/initstate"

	for i := 0; i < 16; i++ {
		// 执行命令
		out1, _ := exec.Command("bash", "-c", cmd1).Output()
		out2, _ := exec.Command("bash", "-c", cmd2).Output()
		output1 := strings.TrimSpace(string(out1))
		output2 := strings.TrimSpace(string(out2))

		glog.V(5).Infof("Attempt %d: Command 1 Output: %s, Command 2 Output: %s", i+1, output1, output2)

		// 判断是否含有 "live"
		if strings.Contains(output1, "live") || strings.Contains(output2, "live") {
			glog.V(5).Infof("Detected 'live', exiting loop.")
			return true
		}

		// 如果还未到达最后一次，则休眠 10 秒
		if i < 15 {
			glog.V(5).Infof("No 'live' detected, sleeping for 10 seconds...")
			time.Sleep(10 * time.Second)
		}
	}

	glog.V(5).Infof("Maximum attempts reached, exiting.")
	return false
}

// ExtractBinaryToTemp 将嵌入的二进制写入临时文件，并设置可执行权限
// 参数:
//   - binBytes: 嵌入的二进制内容
//   - prefix: 生成临时文件名前缀，例如 "rocm-bandwidth-" 或 "hip-stream-"
//
// 返回值:
//   - string: 临时文件路径
//   - error: 错误信息
func ExtractBinaryToTemp(binBytes []byte, prefix string) (string, error) {
	if len(binBytes) == 0 {
		return "", fmt.Errorf("embedded binary is empty")
	}

	tmpFile, err := os.CreateTemp("", prefix+"*")
	if err != nil {
		return "", fmt.Errorf("无法创建临时文件: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(binBytes); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("写入临时文件失败: %w", err)
	}

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("无法设置可执行权限: %w", err)
	}

	return tmpFile.Name(), nil
}

// -------------------- 辅助函数：打印 ASCII 表格 --------------------
func printAsciiTable(title string, headers []string, rows [][]string) {
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, col := range row {
			if len(col) > colWidths[i] {
				colWidths[i] = len(col)
			}
		}
	}

	printLine := func() {
		for _, w := range colWidths {
			fmt.Print("+", strings.Repeat("-", w+2))
		}
		fmt.Println("+")
	}

	fmt.Println()
	fmt.Printf("+%s+\n", strings.Repeat("=", sum(colWidths)+3*len(colWidths)-1))
	fmt.Printf("| %-*s |\n", sum(colWidths)+3*len(colWidths)-1, title)
	fmt.Printf("+%s+\n", strings.Repeat("=", sum(colWidths)+3*len(colWidths)-1))

	printLine()
	// 打印表头
	fmt.Print("|")
	for i, h := range headers {
		fmt.Printf(" %-*s |", colWidths[i], h)
	}
	fmt.Println()
	printLine()
	// 打印数据行
	for _, row := range rows {
		fmt.Print("|")
		for i, col := range row {
			fmt.Printf(" %-*s |", colWidths[i], col)
		}
		fmt.Println()
	}
	printLine()
	fmt.Println()
}

// 计算列总宽度
func sum(arr []int) int {
	total := 0
	for _, v := range arr {
		total += v
	}
	return total
}

// CardSeriesInfo 表示单张卡的系列信息
// DvInd: 卡索引 (如 HCU[0] -> 0)
// SeriesName: 卡系列名称 (如 K100_AI)
type CardSeriesInfo struct {
	DvInd      int    // 卡索引
	SeriesName string // 卡系列名称
}

// CardSeriesList 执行 hy-smi 命令并解析返回卡系列信息
// 命令：/opt/hyhal/bin/hy-smi --showproductname
//
// 返回:
//   - []CardSeriesInfo: 所有卡的系列信息
//   - error: 若命令执行失败或解析失败
func CardSeriesList() ([]CardSeriesInfo, error) {

	cmd := exec.Command("/opt/hyhal/bin/hy-smi", "--showproductname")
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to run hy-smi: %w, output: %s", err, output)
	}

	// 匹配：HCU[0] : Card Series: K100_AI
	re := regexp.MustCompile(`(?m)HCU\[(\d+)\]\s*:\s*Card Series:\s*(.+)$`)
	matches := re.FindAllStringSubmatch(output, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no card series info found in hy-smi output")
	}

	var result []CardSeriesInfo
	seen := make(map[int]bool)

	for _, match := range matches {
		index, _ := strconv.Atoi(strings.TrimSpace(match[1]))
		seriesName := strings.TrimSpace(match[2])

		if seen[index] {
			continue
		}
		seen[index] = true

		result = append(result, CardSeriesInfo{
			DvInd:      index,
			SeriesName: seriesName,
		})
	}

	sort.Slice(result, func(a, b int) bool {
		return result[a].DvInd < result[b].DvInd
	})

	return result, nil
}

// NormalizeCardSeriesName 根据驱动主版本判断是否需要使用卡系列名覆盖传入的型号名称。
// 逻辑：
//   - 若驱动版本属于 6.3.x 则调用 GetCardSeriesList()，取第一张卡的 SeriesName 返回。
//   - 若获取失败或版本不是 6.3.x，则返回原始 typeName。
//
// 适用场景：用于在不同驱动版本下统一卡型号展示。
func NormalizeCardSeriesName(typeName string) string {
	ver, _ := Version(RSMISwCompFirst)
	glog.V(5).Infof("version:%v", ver)
	if IsDriverMajorVersion63(ver) {
		if cards, err := CardSeriesList(); err == nil && len(cards) > 0 {
			series := strings.TrimSpace(cards[0].SeriesName)
			if series != "" {
				return series
			}
		}
	}
	// 默认返回传入的名称
	return typeName
}

/*
IsDriverMajorVersion63 判断驱动版本是否属于 6.3 主版本系列。
示例：

	"6.3.8-V1.9.2" -> true
	"v6.3.1"       -> true
	"6.3"          -> true
*/
func IsDriverMajorVersion63(version string) bool {
	v := strings.TrimSpace(version)
	if v == "" {
		return false
	}
	// 去掉可能的 v/V 前缀
	if v[0] == 'v' || v[0] == 'V' {
		v = v[1:]
	}
	return v == "6.3" || strings.HasPrefix(v, "6.3.")
}

func NormalizeDevTypeName(devTypeName string) string {
	if i := strings.IndexByte(devTypeName, ','); i >= 0 {
		return devTypeName[:i]
	}
	return devTypeName
}

func loadPCIDeviceNames() error {
	file, err := os.Open(pciIDsPath)
	if err != nil {
		return fmt.Errorf("无法打开 pci.ids (%s): %v", pciIDsPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 只处理 device 行（必须有缩进）
		if !strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, "    ") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		vendorID := strings.ToLower(fields[0])
		deviceID := strings.ToLower(fields[1])
		deviceName := strings.Join(fields[2:], " ")

		if vendorID == targetVendor {
			pciDeviceID2Name[deviceID] = deviceName
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取 pci.ids 出错: %v", err)
	}

	glog.V(5).Infof("loaded %d pci device names for vendor %s",
		len(pciDeviceID2Name), targetVendor)

	return nil
}

func loadUpdateIDsMap() error {
	file, err := os.Open(updateIDsPath)
	if err != nil {
		return fmt.Errorf("无法打开 update_ids (%s): %v", updateIDsPath, err)
	}
	defer file.Close()

	updateIDsMap = make(map[string]string)

	re := regexp.MustCompile(`^map\["([0-9a-fA-F]+)"\]="([^"]+)"`)

	scanner := bufio.NewScanner(file)
	inMap := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "declare -A map" {
			inMap = true
			continue
		}

		if inMap && strings.HasPrefix(line, "declare -A ") {
			break
		}

		if !inMap {
			continue
		}

		if m := re.FindStringSubmatch(line); len(m) == 3 {
			updateIDsMap[strings.ToLower(m[1])] = m[2]
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取 update_ids 出错: %v", err)
	}

	glog.V(5).Infof(
		"loaded %d entries from update_ids",
		len(updateIDsMap),
	)

	return nil
}

func pcieThroughputMB(sent, received, maxPktSz int64) (sentMB, receivedMB, bwMB float64) {
	sentMB, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(sent)*float64(maxPktSz)/(1024*1024)), 64)
	receivedMB, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(received)*float64(maxPktSz)/(1024*1024)), 64)
	bwMB, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", (sentMB+receivedMB)/(1024*1024)), 64)
	return
}

// formatBDFID converts a uint64 BDFID to standard PCI format (dddd:bb:dd.f).
func formatBDFID(bdfid uint64) string {
	domain := (bdfid >> 32) & 0xffffffff
	bus := (bdfid >> 8) & 0xff
	device := (bdfid >> 3) & 0x1f
	function := bdfid & 0x7

	return fmt.Sprintf("%04x:%02x:%02x.%x",
		domain,
		bus,
		device,
		function,
	)
}
