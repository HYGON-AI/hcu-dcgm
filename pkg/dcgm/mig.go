/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L/opt/hyhal/lib -Wl,-rpath,/opt/hyhal/lib -lrocm_smi64 -lhydmi -lhydmi_mig -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdio.h>
#include <stdlib.h>
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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"github.com/golang/glog"
)

const (
	// maxGpuSlicesPerDev 对应 C 里的 #define MAX_GPU_SLICES_PER_DEV (4)
	maxGpuSlicesPerDev = 4

	// topologyNodesDir 是 sysfs 中的 topology 路径
	topologyNodesDir = "/sys/devices/virtual/kfd/kfd/topology/nodes"
)

func nvmlDeviceGetCount() (deviceCount int, err error) {
	var count C.uint
	// 调用 C 的 nvmlDeviceGetCount 传递指针
	ret := C.nvmlDeviceGetCount(&count)
	if err = migErrorString(ret); err != nil {
		return 0, fmt.Errorf("clock type not supported: %s", err)
	}
	deviceCount = int(count)
	glog.V(5).Infof("deviceCount: %v", count)
	return
}

// 根据索引号获取特定设备的句柄。
func nvmlDeviceGetHandleByIndex(dvInd int) (device MIGDevice, err error) {
	var dev C.nvmlDevice_t
	ret := C.nvmlDeviceGetHandleByIndex(C.uint(dvInd), &dev)
	if err = migErrorString(ret); err != nil {
		return MIGDevice(nil), fmt.Errorf("nvmlDeviceGetHandleByIndex failed: %v", err)
	}
	device = MIGDevice(dev)
	glog.V(5).Infof("nvmlDeviceGetHandleByIndex:%v", device)
	return
}

func nvmlDeviceGetMigDeviceHandleByIndex(dvInd int, migId int) (migDevice MIGDevice, err error) {
	// 第一步：通过索引获取设备句柄
	device, err := nvmlDeviceGetHandleByIndex(dvInd)
	if err != nil {
		return migDevice, fmt.Errorf("nvmlDeviceGetHandleByIndex(%d) failed: %v", dvInd, err)
	}
	var migDev C.nvmlDevice_t
	ret := C.nvmlDeviceGetMigDeviceHandleByIndex(C.nvmlDevice_t(device), C.uint(migId), &migDev)
	if err = migErrorString(ret); err != nil {
		return migDevice, fmt.Errorf("nvmlDeviceGetMigDeviceHandleByIndex failed: %v", err)
	}
	migDevice = MIGDevice(migDev)
	return
}

// 传设备索引，返回最大MIG数量
func nvmlDeviceGetMaxMigDeviceCountByIndex(dvInd int) (count int, err error) {
	// 第一步：通过索引获取设备句柄
	device, err := nvmlDeviceGetHandleByIndex(dvInd)
	if err != nil {
		return 0, fmt.Errorf("nvmlDeviceGetHandleByIndex(%d) failed: %v", dvInd, err)
	}

	// 第二步：拿到句柄后获取 MIG 设备数量
	var cCount C.uint
	ret := C.nvmlDeviceGetMaxMigDeviceCount(C.nvmlDevice_t(device), &cCount)
	if err = migErrorString(ret); err != nil {
		return 0, fmt.Errorf("nvmlDeviceGetMaxMigDeviceCount failed: %v", err)
	}
	count = int(cCount)
	glog.V(5).Infof("nvmlDeviceGetMaxMigDeviceCount: %v", count)
	return
}

// 根据索引获得设备属性（attributes）
func nvmlDeviceGetAttributesByIndex(dvInd int, migId int) (attr NvmlDeviceAttributes, err error) {
	// 第一步 获取设备句柄
	device, err := nvmlDeviceGetMigDeviceHandleByIndex(dvInd, migId)
	if err != nil {
		return attr, fmt.Errorf("nvmlDeviceGetMigDeviceHandleByIndex(%d) failed: %v", dvInd, err)
	}

	// 第二步 调用C函数并获取C结构体
	var cAttr C.nvmlDeviceAttributes_t
	ret := C.nvmlDeviceGetAttributes(C.nvmlDevice_t(device), &cAttr)
	if err = migErrorString(ret); err != nil {
		return attr, fmt.Errorf("nvmlDeviceGetAttributes failed: %v", err)
	}
	// 第三步：C结构体转Go结构体
	attr = NvmlDeviceAttributes{
		Index:                     uint32(cAttr.index),
		CUCount:                   uint32(cAttr.cu_count),
		MemorySizeMB:              uint64(cAttr.memory_size_MB),
		UUID:                      C.GoString(&cAttr.uuid[0]),
		Name:                      C.GoString(&cAttr.name[0]),
		GPUInstanceSliceCount:     uint32(cAttr.gpu_instance_slice_count),
		ComputeInstanceSliceCount: uint32(cAttr.compute_instance_slice_count),
	}
	return
}

// nvmlDeviceSetMigMode 设置指定设备的 MIG（Multi-Instance GPU）模式。
// dvInd        —— 设备的索引ID（从0开始）；
// mode            —— 要设置的模式，NVML_DEVICE_MIG_ENABLE 或 NVML_DEVICE_MIG_DISABLE；
// 返回 activationStatus（激活结果，成功为NVML_SUCCESS）。
func nvmlDeviceSetMigModeByIndex(dvInd int, mode uint32) (activationStatus uint32, err error) {
	// 获取设备句柄
	device, err := nvmlDeviceGetHandleByIndex(dvInd)
	if err != nil {
		return 0, fmt.Errorf("nvmlDeviceGetHandleByIndex(%d) failed: %v", dvInd, err)
	}

	// C函数调用
	var cActivationStatus C.nvmlReturn_t
	ret := C.nvmlDeviceSetMigMode(
		C.nvmlDevice_t(device),
		C.uint(mode),
		&cActivationStatus,
	)
	if err = migErrorString(ret); err != nil {
		return uint32(cActivationStatus), fmt.Errorf("nvmlDeviceSetMigMode failed: %v", err)
	}
	return uint32(cActivationStatus), nil
}

// nvmlDeviceGetMigModeByIndex 获取指定GPU设备的当前和待激活（pending）MIG模式。
//
// 参数：
//
//	dvInd: 设备在本机的索引（从0开始）
//
// 返回：
//
//	currentMode:   当前MIG模式（NVML_DEVICE_MIG_ENABLE或NVML_DEVICE_MIG_DISABLE）
//	pendingMode:   待激活的MIG模式（等待设备重启等操作后生效）
//	err:           错误信息，若调用成功则为nil
func nvmlDeviceGetMigModeByIndex(dvInd int) (currentMode uint32, pendingMode uint32, err error) {
	// 获取设备句柄
	device, err := nvmlDeviceGetHandleByIndex(dvInd)
	if err != nil {
		return
	}

	var cCurrentMode C.uint
	var cPendingMode C.uint

	ret := C.nvmlDeviceGetMigMode(
		C.nvmlDevice_t(device),
		&cCurrentMode,
		&cPendingMode,
	)
	if goErr := migErrorString(ret); goErr != nil {
		err = fmt.Errorf("nvmlDeviceGetMigMode failed: %v", goErr)
		return
	}
	currentMode = uint32(cCurrentMode)
	pendingMode = uint32(cPendingMode)
	return
}

// nvmlDeviceGetGpuInstanceRemainingCapacity 查询指定GPU和profile的MIG实例剩余容量
// dvInd:     GPU卡的索引号（0、1...）
// profileId: GPU instance profile 的ID
// 返回值: 该profile还能创建多少个instance
func nvmlDeviceGetGpuInstanceRemainingCapacity(dvInd int, profileId uint32) (count uint32, err error) {
	// 获取GPU设备句柄
	device, err := nvmlDeviceGetHandleByIndex(dvInd)
	if err != nil {
		return
	}

	var cCount C.uint
	ret := C.nvmlDeviceGetGpuInstanceRemainingCapacity(
		C.nvmlDevice_t(device),
		C.uint(profileId),
		&cCount,
	)

	if errC := migErrorString(ret); errC != nil {
		err = fmt.Errorf("nvmlDeviceGetGpuInstanceRemainingCapacity failed: %w", errC)
		return
	}

	count = uint32(cCount)
	return
}

// nvmlDeviceGetGpuInstanceProfileInfo 获取指定GPU、指定profile编号的MIG GPU实例profile详细信息。
//
// dvInd:      物理GPU的索引号（从0开始，依次递增）
// profileId: profile编号（一般为 0~NVML_GPU_INSTANCE_PROFILE_COUNT-1）
// 返回值:
//   - NvmlGpuInstanceProfileInfo：profile的详细信息结构体（包含id、名称、cu数量、内存等）
//   - error：调用失败时返回详细错误信息
func nvmlDeviceGetGpuInstanceProfileInfo(dvInd int, profileId uint32) (info NvmlGpuInstanceProfileInfo, err error) {
	// 获取GPU设备句柄
	device, err := nvmlDeviceGetHandleByIndex(dvInd)
	if err != nil {
		return
	}
	var cInfo C.nvmlGpuInstanceProfileInfo_t
	ret := C.nvmlDeviceGetGpuInstanceProfileInfo(
		C.nvmlDevice_t(device),
		C.uint(profileId),
		&cInfo,
	)
	if errC := migErrorString(ret); errC != nil {
		err = fmt.Errorf("nvmlDeviceGetGpuInstanceProfileInfo failed: %w", errC)
		return
	}

	// 将C结构体转成Go结构体
	info = NvmlGpuInstanceProfileInfo{
		ID:          uint32(cInfo.id),
		GiCountMax:  uint32(cInfo.gi_count_max),
		CuCount:     uint32(cInfo.cu_count),
		GpuSliceCnt: uint32(cInfo.gpu_slice_count),
		MemSizeMB:   uint64(cInfo.memory_size_MB),
		Name:        C.GoString(&cInfo.name[0]),
	}
	return
}

// nvmlDeviceGetGpuInstances
// 获取指定物理GPU和profileId下的所有已存在MIG Instance的详细信息列表。
//
// dvInd:     物理GPU的索引号（从0开始递增）
// profileId: MIG实例profile类型编号（NVML_GPU_INSTANCE_PROFILE_1_SLICE 等）
//
// 返回值：
//   - []GpuInstanceInfo：GPU上、该profileId类型下所有已存在的MIG实例信息切片（每个元素包含Instance的giId、profileId等关键字段）
//   - error：调用失败时返回详细错误信息，正常时为nil
func nvmlDeviceGetGpuInstances(dvInd int, profileId uint32) ([]GpuInstanceInfo, error) {
	device, err := nvmlDeviceGetHandleByIndex(dvInd)
	if err != nil {
		return nil, err
	}

	const maxInstance = 8
	handles := make([]C.nvmlGpuInstance_t, maxInstance)
	count := C.uint(maxInstance)

	ret := C.nvmlDeviceGetGpuInstances(
		C.nvmlDevice_t(device),
		C.uint(profileId),
		(*C.nvmlGpuInstance_t)(unsafe.Pointer(&handles[0])),
		&count,
	)
	if errC := migErrorString(ret); errC != nil {
		return nil, errC
	}

	var out []GpuInstanceInfo
	for i := 0; i < int(count); i++ {
		var info C.nvmlGpuInstanceInfo_t
		ret = C.nvmlGpuInstanceGetInfo(handles[i], &info)
		if errC := migErrorString(ret); errC != nil {
			continue
		}
		out = append(out, GpuInstanceInfo{
			Device:    uintptr(unsafe.Pointer(info.device)), // 设备句柄，以 uintptr 存储
			ID:        uint32(info.id),                      // 唯一实例ID（giId）
			ProfileID: uint32(info.profile_id),              // Profile编号
			Placement: GpuInstancePlacement{ // 占用 compute slice 信息
				Start: uint32(info.placement.start),
				Size:  uint32(info.placement.size),
			},
		})
	}
	return out, nil
}

// nvmlGpuInstanceGetComputeInstanceProfileInfo
// 获取指定 GPU instance 的指定 compute profile 和引擎 profile 的 ComputeInstanceProfile 详细信息。
//
// 参数：
//
//	dvInd       - 物理 GPU 的索引号（如 0、1、2...）
//	giId        - GPU Instance ID（MIG实例唯一标识；可用 nvmlDeviceGetGpuInstances 获取，也可由 API 创建新的 Instance 后获得）
//	profileId   - Compute Profile 编号（取值范围见 NVML_COMPUTE_INSTANCE_PROFILE_* 常量，通常为 0~NVML_COMPUTE_INSTANCE_PROFILE_COUNT-1）
//	engProfileId- 引擎类型编号（通常传 0，即 NVML_COMPUTE_INSTANCE_ENGINE_PROFILE_SHARED）
//
// 返回值：
//
//	info  - NvmlComputeInstanceProfileInfo 结构体，包含该 profile 的 id、最大实例数、CU 数量、slice 数、profile 名称等关键信息
//	err   - 调用过程中出现的错误（若无错则为 nil）
//
// 典型场景：本函数用于查询某 MIG instance 下，指定 compute profile 在当前硬件上的详细支持参数。
// 使用前需确保该 GPU 处于已开启 MIG 模式，且目标 MIG instance 已存在。
func nvmlGpuInstanceGetComputeInstanceProfileInfo(dvInd int, giId uint32, profileId uint32, engProfileId uint32) (info NvmlComputeInstanceProfileInfo, err error) {
	// 步骤一：获取物理 GPU 句柄
	device, err := nvmlDeviceGetHandleByIndex(dvInd)
	if err != nil {
		return
	}

	// 步骤二：根据 giId 获取 GPU instance 句柄
	var giHandle C.nvmlGpuInstance_t
	ret := C.nvmlDeviceGetGpuInstanceById(
		C.nvmlDevice_t(device),
		C.uint(giId),
		&giHandle,
	)
	if errC := migErrorString(ret); errC != nil {
		err = fmt.Errorf("nvmlDeviceGetGpuInstanceById failed: %w", errC)
		return
	}

	// 步骤三：查询该 instance 下指定 profileId/engProfileId 的详细 ComputeInstanceProfile 信息
	var cInfo C.nvmlComputeInstanceProfileInfo_t
	ret = C.nvmlGpuInstanceGetComputeInstanceProfileInfo(
		giHandle,
		C.uint(profileId),
		C.uint(engProfileId),
		&cInfo,
	)
	if errC := migErrorString(ret); errC != nil {
		err = fmt.Errorf("nvmlGpuInstanceGetComputeInstanceProfileInfo failed: %w", errC)
		return
	}

	// 步骤四：转换为Go结构体输出
	info = NvmlComputeInstanceProfileInfo{
		ID:          uint32(cInfo.id),              // profile唯一编号
		CiCountMax:  uint32(cInfo.ci_count_max),    // 当前profile最大支持的Compute Instance数
		CuCount:     uint32(cInfo.cu_count),        // 可用CU核心数
		GpuSliceCnt: uint32(cInfo.gpu_slice_count), // Slice数
		Name:        C.GoString(&cInfo.name[0]),    // profile名称
	}
	return
}

// allComputeInstanceProfileInfo
// 查询指定物理GPU、指定MIG profileId下，所有已存在MIG Instance对应的compute instance profile信息。
//
// 参数：
//
//	dvInd     - 物理GPU索引号（如 0, 1, 2 ...）
//	profileId - MIG instance profile类型编号（如 NVML_GPU_INSTANCE_PROFILE_1_SLICE 等）
//	ciProfileId - compute profile编号（如 0~NVML_COMPUTE_INSTANCE_PROFILE_COUNT-1）
//	engProfileId- 引擎编号（一般传0）
//
// 返回值：
//   - []NvmlComputeInstanceProfileInfo：每个MIG Instance下，该profile对应的详细信息切片
//   - error：如有任何NVML调用失败则返回
func allComputeInstanceProfileInfo(dvInd int, profileId, ciProfileId, engProfileId uint32) ([]NvmlComputeInstanceProfileInfo, error) {
	// 第一步：查出所有已存在MIG Instance的giId
	instances, err := DeviceGetGpuInstancesInfo(dvInd, profileId)
	if err != nil {
		return nil, err
	}

	var out []NvmlComputeInstanceProfileInfo
	for _, inst := range instances {
		info, err := GpuInstanceGetComputeInstanceProfileInfo(dvInd, inst.ID, ciProfileId, engProfileId)
		if err != nil {
			continue // 或log失败的giId及错误，再决定继续或return
		}
		out = append(out, info)
	}
	return out, nil
}

// nvmlGpuInstanceGetComputeInstanceRemainingCapacity
// 查询指定物理GPU上某个已存在MIG实例(giId)，以及指定compute profile类型下，
// 还可创建的compute instance数量。
//
// dvInd    : 物理GPU的索引号（从0开始递增，对应nvidia-smi中的编号）
// giId     : MIG实例的唯一ID（giId，在nvmlDeviceGetGpuInstances返回的 GpuInstanceInfo 中获取）
// profileId: 要查询的compute instance profile类型编号
//
// 返回值：
//   - uint32：该MIG实例下，指定compute profile还能新建的compute instance数量
//   - error ：调用失败时返回详细错误信息，正常时为nil
func nvmlGpuInstanceGetComputeInstanceRemainingCapacity(
	dvInd int,
	giId uint32,
	profileId uint32,
) (uint32, error) {
	// 第一步：获取目标物理GPU设备句柄
	device, err := nvmlDeviceGetHandleByIndex(dvInd)
	if err != nil {
		return 0, err
	}

	// 第二步：枚举该设备上的全部MIG实例，查找giId对应的句柄
	const maxInstance = 8
	handles := make([]C.nvmlGpuInstance_t, maxInstance)
	count := C.uint(maxInstance)

	ret := C.nvmlDeviceGetGpuInstances(
		C.nvmlDevice_t(device),
		0, // 0表示查找全部profileId
		(*C.nvmlGpuInstance_t)(unsafe.Pointer(&handles[0])),
		&count,
	)
	if errC := migErrorString(ret); errC != nil {
		return 0, errC
	}

	// 匹配giId对应的gpuInstance句柄
	var giHandle C.nvmlGpuInstance_t
	found := false
	for i := 0; i < int(count); i++ {
		var info C.nvmlGpuInstanceInfo_t
		ret = C.nvmlGpuInstanceGetInfo(handles[i], &info)
		if errC := migErrorString(ret); errC != nil {
			continue // 某个实例失败可跳过
		}
		if uint32(info.id) == giId {
			giHandle = handles[i]
			found = true
			break
		}
	}
	if !found {
		return 0, fmt.Errorf("giId=%d does not exist on GPU%d", giId, dvInd)
	}

	// 第三步：查该MIG实例还可新建指定profile的compute instance数量
	var remain C.uint
	ret = C.nvmlGpuInstanceGetComputeInstanceRemainingCapacity(
		giHandle,
		C.uint(profileId),
		&remain,
	)
	if errC := migErrorString(ret); errC != nil {
		return 0, errC
	}
	return uint32(remain), nil
}

// nvmlAllGpuInstancesComputeInstanceRemainingCapacity
// 获取指定物理GPU、指定MIG实例profile类型下，所有已存在MIG实例目前还可新建的compute instance数量。
// （不区分compute profile类型，只统计NVML允许的新建数量）
//
// dvInd    : 物理GPU的索引号（从0开始递增，对应nvidia-smi中的编号）
// profileId: MIG实例profile类型编号（如 NVML_GPU_INSTANCE_PROFILE_1_SLICE 等，具体枚举值根据显卡支持为准）
//
// 返回值：
//   - []ComputeInstanceRemainInfo：对每个MIG实例，返回giId、profileId和当前剩余可新建compute instance数
//   - error：调用失败时返回详细错误信息，正常时为nil
func nvmlAllGpuInstancesComputeInstanceRemainingCapacity(
	dvInd int,
	profileId uint32,
) ([]ComputeInstanceRemainInfo, error) {

	// 1. 获取已存在的目标profile类型MIG实例信息（包含giId）
	gis, err := nvmlDeviceGetGpuInstances(dvInd, profileId)
	if err != nil {
		return nil, err
	}
	if len(gis) == 0 {
		return []ComputeInstanceRemainInfo{}, nil
	}

	// 2. 获取设备句柄
	device, err := nvmlDeviceGetHandleByIndex(dvInd)
	if err != nil {
		return nil, err
	}

	// 3. 枚举MIG实例（这里profileId已筛选，后面只需遍历符合要求的实例即可）
	const maxInstance = 8
	handles := make([]C.nvmlGpuInstance_t, maxInstance)
	count := C.uint(maxInstance)
	ret := C.nvmlDeviceGetGpuInstances(
		C.nvmlDevice_t(device),
		C.uint(profileId),
		(*C.nvmlGpuInstance_t)(unsafe.Pointer(&handles[0])),
		&count,
	)
	if errC := migErrorString(ret); errC != nil {
		return nil, errC
	}

	var out []ComputeInstanceRemainInfo
	for i := 0; i < int(count); i++ {
		var info C.nvmlGpuInstanceInfo_t
		ret = C.nvmlGpuInstanceGetInfo(handles[i], &info)
		if errC := migErrorString(ret); errC != nil {
			continue
		}
		giId := uint32(info.id)

		// 不指定profileId，默认为0，即查所有CI profile的整体剩余
		var remain C.uint
		ret = C.nvmlGpuInstanceGetComputeInstanceRemainingCapacity(
			handles[i],
			0, // 0表示查所有compute profile剩余（查“整体”）
			&remain,
		)
		if errC := migErrorString(ret); errC != nil {
			out = append(out, ComputeInstanceRemainInfo{
				GiID:          giId,
				ProfileID:     uint32(info.profile_id),
				CiRemainCount: 0,
			})
			continue
		}
		out = append(out, ComputeInstanceRemainInfo{
			GiID:          giId,
			ProfileID:     uint32(info.profile_id),
			CiRemainCount: uint32(remain),
		})
	}
	return out, nil
}

// nvmlDeviceGetGpuInstanceId 获取设备对应的GPU实例ID
func nvmlDeviceGetGpuInstanceId(dvInd int, migId int) (gpuInstanceId uint32, err error) {
	// 获取设备句柄
	device, err := nvmlDeviceGetMigDeviceHandleByIndex(dvInd, migId)
	if err != nil {
		return 0, fmt.Errorf("nvmlDeviceGetHandleByIndex(%d) failed: %v", dvInd, err)
	}

	var giId C.uint
	ret := C.nvmlDeviceGetGpuInstanceId(
		C.nvmlDevice_t(device), // 传递MIG设备句柄
		&giId,                  // 获取GPU实例ID的指针
	)

	if errC := migErrorString(ret); errC != nil {
		err = fmt.Errorf("nvmlDeviceGetGpuInstanceId failed: %w", errC)
		return 0, err
	}
	gpuInstanceId = uint32(giId)
	return
}

// nvmlDeviceGetComputeInstanceId
// 根据物理GPU索引 dvInd 获取 Compute Instance ID
func nvmlDeviceGetComputeInstanceId(dvInd int, migId int) (computeInstanceId uint32, err error) {
	// 1. 获取设备句柄
	device, err := nvmlDeviceGetMigDeviceHandleByIndex(dvInd, migId)
	if err != nil {
		err = fmt.Errorf("get device handle failed: %w", err)
		return
	}

	// 2. 调用 C API 获取 Compute Instance ID
	var ciId C.uint
	ret := C.nvmlDeviceGetComputeInstanceId(
		C.nvmlDevice_t(device),
		&ciId,
	)
	if errC := migErrorString(ret); errC != nil {
		err = fmt.Errorf("nvmlDeviceGetComputeInstanceId failed: %w", errC)
		return
	}

	computeInstanceId = uint32(ciId)
	return
}

// availableMigDeviceIds 返回指定GPU下所有可用MIG device的id
func availableMigDeviceIds(dvInd int) (int, []int, error) {
	log.Printf("[INFO] 开始查询 GPU 索引 %d 的 MIG device 信息...", dvInd)

	device, err := nvmlDeviceGetHandleByIndex(dvInd)
	if err != nil {
		log.Printf("[ERROR] 获取 GPU %d 句柄失败: %v", dvInd, err)
		return dvInd, nil, fmt.Errorf("nvmlDeviceGetHandleByIndex(%d) failed: %v", dvInd, err)
	}

	maxMigCount, err := nvmlDeviceGetMaxMigDeviceCountByIndex(dvInd)
	if err != nil {
		log.Printf("[ERROR] 获取 GPU %d 最大 MIG 数失败: %v", dvInd, err)
		return dvInd, nil, fmt.Errorf("nvmlDeviceGetMaxMigDeviceCountByIndex(%d) failed: %v", dvInd, err)
	}
	log.Printf("[DEBUG] GPU %d 最大支持的 MIG 数量: %d", dvInd, maxMigCount)

	var ids []int
	for id := 0; id < maxMigCount; id++ {
		var migDev C.nvmlDevice_t
		ret := C.nvmlDeviceGetMigDeviceHandleByIndex(C.nvmlDevice_t(device), C.uint(id), &migDev)
		if NvmlReturn(ret) == NVML_ERROR_NOT_FOUND {
			log.Printf("[INFO] GPU %d 在第 %d 个索引处没有更多 MIG device", dvInd, id)
			break
		}
		if err := migErrorString(ret); err != nil {
			log.Printf("[WARN] GPU %d 获取 MIG device 索引 %d 失败: %v", dvInd, id, err)
			continue
		}
		log.Printf("[DEBUG] GPU %d 获取到有效 MIG device: id=%d", dvInd, id)
		ids = append(ids, id)
	}

	log.Printf("[INFO] GPU %d 可用的 MIG device 索引: %v", dvInd, ids)
	return dvInd, ids, nil
}

// 假设 MIGDevice 本质上可被转为 C.nvmlDevice_t 类型（你可根据实际做类型转换）
func nvmlDeviceGetPciInfo(device MIGDevice) (pciInfo NvmlPciInfo, err error) {
	// 1. 声明C结构体用于接收结果
	var cPciInfo C.nvmlPciInfo

	// 2. 转换MIGDevice到C.nvmlDevice_t
	// 这里假定MIGDevice本质就是C.nvmlDevice_t，如果不是请自行补充转换逻辑
	cDevice := C.nvmlDevice_t(device)

	// 3. 调用NVML C API
	ret := C.nvmlDeviceGetPciInfo(cDevice, &cPciInfo)
	if errC := migErrorString(ret); errC != nil {
		err = fmt.Errorf("nvmlDeviceGetPciInfo failed: %w", errC)
		return
	}

	// 4. Go结构体赋值
	pciInfo = NvmlPciInfo{
		Domain:   uint32(cPciInfo.pci_domain),
		Bus:      uint32(cPciInfo.pci_bus),
		Device:   uint32(cPciInfo.pci_device),
		Function: uint32(cPciInfo.pci_function),
		BusID:    C.GoString(&cPciInfo.bus_id[0]),
	}
	return
}

func nvmlGetSystemMigMode() (currentMode, pendingMode int, err error) {
	// 1. 声明C类型变量用于接收
	var cCurrent C.uint
	var cPending C.uint

	// 2. 调用C API
	ret := C.nvmlGetSystemMigMode(&cCurrent, &cPending)
	if errC := migErrorString(ret); errC != nil {
		err = fmt.Errorf("nvmlGetSystemMigMode failed: %w", errC)
		return
	}

	// 3. 转为Go int类型
	currentMode = int(cCurrent)
	pendingMode = int(cPending)
	return
}

// migConfigs 解析所有MIG配置文件
// 返回: MigConfig切片和错误信息
func migConfigs() ([]MigConfig, error) {
	const (
		baseDir = "/etc/dmi_mig_config" // 基础配置目录
		ciDir   = "/ci" // CI配置子目录
	)

	// 检查基础目录是否存在
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("dmi_mig_config directory not found in /etc")
	}

	// 构建完整CI目录路径
	fullCiDir := filepath.Join(baseDir, ciDir)

	// 检查CI目录是否存在
	if _, err := os.Stat(fullCiDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("ci directory not found in %s", baseDir)
	}

	// 查找所有.conf文件
	files, err := filepath.Glob(filepath.Join(fullCiDir, "*.conf"))
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	// 如果没有找到配置文件
	if len(files) == 0 {
		return nil, fmt.Errorf("no .conf files found in %s", fullCiDir)
	}

	// 解析所有配置文件
	var configs []MigConfig
	for _, file := range files {
		// 从文件名提取GPU ID
		filename := filepath.Base(file)
		dvInd, err := extractGpuId(filename)
		if err != nil {
			fmt.Printf("Warning: skipping file %s - %v\n", filename, err)
			continue
		}

		// 解析配置文件内容
		config, err := parseSingleConfFile(file)
		if err != nil {
			fmt.Printf("Warning: error parsing %s - %v\n", filename, err)
			continue
		}

		// 设置GPU ID
		config.DvInd = dvInd
		configs = append(configs, config)
	}

	// 检查是否成功解析了任何文件
	if len(configs) == 0 {
		return nil, fmt.Errorf("no valid .conf files parsed")
	}

	return configs, nil
}

// parseSingleConfFile 解析单个配置文件
func parseSingleConfFile(path string) (MigConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return MigConfig{}, err
	}
	defer file.Close()

	var config MigConfig
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		lineNum++

		// 解析前13行 (GI部分)
		if lineNum <= 13 {
			switch key {
			case "pci":
				config.Gi.Pci = value
			case "id":
				config.Gi.Id, _ = strconv.Atoi(value)
			case "pipe_mask":
				config.Gi.PipeMask = value
			case "gpu_slice_mask":
				config.Gi.GpuSliceMask = value
			case "cu_mask":
				if config.Gi.CuMask1 == "" {
					config.Gi.CuMask1 = value
				} else {
					config.Gi.CuMask2 = value
				}
			case "profile_id":
				config.Gi.ProfileId, _ = strconv.Atoi(value)
			case "gi_count_max":
				config.Gi.GiCountMax, _ = strconv.Atoi(value)
			case "cu_count":
				config.Gi.CuCount, _ = strconv.Atoi(value)
			case "gpu_slice_count":
				config.Gi.GpuSliceCount, _ = strconv.Atoi(value)
			case "memory_size_MB":
				config.Gi.MemorySizeMB, _ = strconv.Atoi(value)
			case "placement_start":
				config.Gi.PlacementStart, _ = strconv.Atoi(value)
			case "placement_size":
				config.Gi.PlacementSize, _ = strconv.Atoi(value)
			}
		} else { // 解析后14行 (CI部分)
			switch key {
			case "pci":
				config.Ci.Pci = value
			case "gi_id":
				config.Ci.GiId, _ = strconv.Atoi(value)
			case "id":
				config.Ci.Id, _ = strconv.Atoi(value)
			case "pipe_mask":
				config.Ci.PipeMask = value
			case "gpu_slice_mask":
				config.Ci.GpuSliceMask = value
			case "cu_mask":
				if config.Ci.CuMask1 == "" {
					config.Ci.CuMask1 = value
				} else {
					config.Ci.CuMask2 = value
				}
			case "profile_id":
				config.Ci.ProfileId, _ = strconv.Atoi(value)
			case "ci_count_max":
				config.Ci.CiCountMax, _ = strconv.Atoi(value)
			case "cu_count":
				config.Ci.CuCount, _ = strconv.Atoi(value)
			case "gpu_slice_count":
				config.Ci.GpuSliceCount, _ = strconv.Atoi(value)
			case "placement_start":
				config.Ci.PlacementStart, _ = strconv.Atoi(value)
			case "placement_size":
				config.Ci.PlacementSize, _ = strconv.Atoi(value)
			case "mig_uuid":
				config.Ci.MigUUID = value
			}
		}
	}

	return config, scanner.Err()
}

// extractGpuId 从文件名提取GPU设备ID（无正则表达式版本）
func extractGpuId(filename string) (int, error) {
	// 基本格式验证
	if !strings.HasPrefix(filename, "dev") ||
		!strings.Contains(filename, "gi") ||
		!strings.HasSuffix(filename, ".conf") {
		return 0, fmt.Errorf("invalid filename format: %s", filename)
	}

	// 找到"dev"之后"gi"之前的部分
	start := 3 // 跳过"dev"
	end := strings.Index(filename, "gi")
	if end <= start {
		return 0, fmt.Errorf("missing gi in filename: %s", filename)
	}

	// 提取数字部分
	numStr := filename[start:end]
	if numStr == "" {
		return 0, fmt.Errorf("empty device ID in filename: %s", filename)
	}

	// 转换为整数
	return strconv.Atoi(numStr)
}

// getSECount 从 topologyNodesDir 里自动
// 1) 找到第一个 gpu_id != 0 的节点目录
// 2) 读取它的 array_count 和 simd_arrays_per_engine
// 3) 返回 array_count/simd_arrays_per_engine 计算结果
func getSECount() (int, error) {
	entries, err := os.ReadDir(topologyNodesDir)
	if err != nil {
		return 0, fmt.Errorf("读取节点目录失败: %w", err)
	}

	for _, e := range entries {
		if !e.IsDir() || !isNumeric(e.Name()) {
			continue
		}
		nodeDir := filepath.Join(topologyNodesDir, e.Name())

		// 读取 gpu_id，跳过 0
		gpuID, err := readIntFromFile(filepath.Join(nodeDir, "gpu_id"))
		if err != nil {
			log.Printf("节点 %s: 读取 gpu_id 失败: %v\n", e.Name(), err)
			continue
		}
		if gpuID == 0 {
			continue
		}

		// 读取 properties 中的 array_count 与 simd_arrays_per_engine
		propsPath := filepath.Join(nodeDir, "properties")
		f, err := os.Open(propsPath)
		if err != nil {
			return 0, fmt.Errorf("打开 %s 失败: %w", propsPath, err)
		}
		defer f.Close()

		var arrayCount, simdPerEngine int
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) != 2 {
				continue
			}
			key, val := fields[0], fields[1]
			switch key {
			case "array_count":
				arrayCount, _ = strconv.Atoi(val)
			case "simd_arrays_per_engine":
				simdPerEngine, _ = strconv.Atoi(val)
			}
			if arrayCount > 0 && simdPerEngine > 0 {
				// 计算并返回 seCount
				return arrayCount / simdPerEngine, nil
			}
		}
		if err := scanner.Err(); err != nil {
			return 0, fmt.Errorf("扫描 %s 失败: %w", propsPath, err)
		}
	}

	return 0, fmt.Errorf("未找到有效 GPU 节点")
}

// readIntFromFile 读取一个只包含整数的文件
func readIntFromFile(path string) (int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(b)))
}

// isNumeric 判断字符串是否全数字
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// formatMIGName 构造 MIG 实例名称。
// 参数：
//   - seCount:     总的 Streaming Engine 数（SE count）
//   - giSliceCount:GI 实例所占用的 GPU slice 数量
//   - ciSliceCount:CI 实例所占用的 GPU slice 数量（若为 0，则表示只关注 GI）
//   - memorySizeMB:GI 实例总显存，单位 MB
//
// 返回值示例：
//   - 当 ciSliceCount == 0 时，返回 "MIG {totalGISlices}g.{memoryGB}gb"
//   - 当 ciSliceCount > 0 时，返回 "MIG {totalCISlices}c.{totalGISlices}g.{memoryGB}gb"
func formatMIGName(
	seCount int,
	giSliceCount int,
	ciSliceCount int,
	memorySizeMB uint64,
) string {
	// 将显存从 MB 转为 GB，向下取整
	memoryGB := int(memorySizeMB / 1024)

	// 每个 GPU slice 对应多少 SE unit
	sePerSlice := seCount / maxGpuSlicesPerDev

	// 如果 giSliceCount 和 ciSliceCount 相等，只生成 GI 部分名称："{totalSlices}g.{memory}gb"
	if giSliceCount == ciSliceCount {
		totalSlices := sePerSlice * giSliceCount
		return fmt.Sprintf("MIG %dg.%dgb", totalSlices, memoryGB)
	}

	// 如果两者不相等，同时生成 CI (c) 和 GI (g) 部分名称："{totalCISlices}c.{totalGISlices}g.{memory}gb"
	totalCISlices := sePerSlice * ciSliceCount
	totalGISlices := sePerSlice * giSliceCount
	return fmt.Sprintf(
		"MIG %dc.%dg.%dgb",
		totalCISlices,
		totalGISlices,
		memoryGB,
	)
}

