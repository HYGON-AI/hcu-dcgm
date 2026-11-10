/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L/opt/hyhal/lib -Wl,--enable-new-dtags -Wl,-rpath,/opt/hyhal/lib -lrocm_smi64 -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdint.h>
#include <kfd_ioctl.h>
#include <rocm_smi64Config.h>
#include <rocm_smi.h>
*/
import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/golang/glog"
)

// rsmiComputeProcessInfoGet 获取当前使用GPU的所有进程信息
func rsmiComputeProcessInfoGet() (processInfo []ProcessInfos, numItems int, err error) {
	var cnumItems C.uint32_t
	// 第一次调用获取进程数量
	ret := C.rsmi_compute_process_info_get(nil, &cnumItems)
	glog.V(5).Infof("rsmiComputeProcessInfoGet:%v, cnumItems:%v", ret, cnumItems)
	if err := errorString(ret); err != nil {
		return nil, 0, fmt.Errorf("Error rsmiComputeProcessInfoGet (initial call): %s", err)
	}
	// 如果数量为零，返回空
	if cnumItems == 0 {
		return nil, 0, nil
	}
	// 创建一个大小为cnumItems的切片
	processInfo = make([]ProcessInfos, int(cnumItems))
	// 第二次调用以获取实际的数据
	ret = C.rsmi_compute_process_info_get((*C.rsmi_process_info_t)(unsafe.Pointer(&processInfo[0])), &cnumItems)
	if err := errorString(ret); err != nil {
		return nil, 0, fmt.Errorf("Error rsmiComputeProcessInfoGet: %s", err)
	}
	numItems = int(cnumItems)
	glog.V(5).Infof("numItems:%v,processInfo:%v", numItems, processInfo)
	return
}

// copyProcessGPUIndices copies the 32-bit C GPU indices into Go ints.
// The RSMI ABI declares gpuIndex as an array of C int values. On 64-bit
// systems, C int is 4 bytes while Go int is 8 bytes, so the array must not
// be reinterpreted as []int with unsafe.Pointer.
func copyProcessGPUIndices(values []int32, count int) ([]int, error) {
	if count < 0 || count > len(values) {
		return nil, fmt.Errorf("invalid GPU count %d, expected 0..%d", count, len(values))
	}

	indices := make([]int, count)
	for i := 0; i < count; i++ {
		indices[i] = int(values[i])
	}
	return indices, nil
}

// getProcessInfoByPID 根据进程 ID 获取进程的计算信息
func getProcessInfoByPID(pid uint32) (rsmiProcessInfoV2 RsmiProcessInfoV2, err error) {
	// 定义一个 C 语言的结构体，用于存储进程信息
	var proc C.rsmi_process_info_v2_t

	// 调用 C 函数获取进程信息
	ret := C.rsmi_compute_process_info_by_pid_get_v2(C.uint32_t(pid), &proc)
	if err = errorString(ret); err != nil {
		return rsmiProcessInfoV2, fmt.Errorf("failed to get process info")
	}

	gpuCount := int(proc.usedGpus)
	if gpuCount < 0 || gpuCount > len(proc.gpuIndex) {
		return rsmiProcessInfoV2, fmt.Errorf("invalid GPU count %d, expected 0..%d", gpuCount, len(proc.gpuIndex))
	}

	// The RSMI ABI declares gpuIndex as an array of 32-bit C int values.
	rawGPUIndices := make([]int32, gpuCount)
	for i := 0; i < gpuCount; i++ {
		rawGPUIndices[i] = int32(proc.gpuIndex[i])
	}
	gpuIndex, err := copyProcessGPUIndices(rawGPUIndices, gpuCount)
	if err != nil {
		return rsmiProcessInfoV2, err
	}

	gpuUsageRate := make([]float32, gpuCount)
	for i := 0; i < gpuCount; i++ {
		gpuUsageRate[i] = float32(proc.gpuUsageRate[i])
	}

	// 构建并返回 Go 语言的结构体
	rsmiProcessInfoV2 = RsmiProcessInfoV2{
		ProcessID:     uint32(proc.processId),      // 进程 ID
		VRAMUsageSize: uint64(proc.vramUsageSize),  // VRAM 使用量（单位 MiB）
		VRAMUsageRate: float32(proc.vramUsageRate), // VRAM 使用率（百分比）
		GPUCount:      int(proc.usedGpus),          // 使用的 GPU 数量
		GPUIndex:      gpuIndex,                    // GPU 索引切片
		GPUUsageRate:  gpuUsageRate,                // GPU 使用率切片
	}
	glog.V(5).Infof("rsmiProcessInfoV2:%v", rsmiProcessInfoV2)
	return
}

// rsmiComputeProcessInfoByPidGet 获取指定进程的进程信息
func rsmiComputeProcessInfoByPidGet(pid int) (proc ProcessInfos, err error) {
	var cproc C.rsmi_process_info_t
	glog.V(5).Infof("rsmiComputeProcessInfoByPid: %v", pid)
	ret := C.rsmi_compute_process_info_by_pid_get(C.uint32_t(pid), &cproc)
	if err = errorString(ret); err != nil {
		return proc, fmt.Errorf("Error rsmiComputeProcessInfoByPidGet: %s", err)
	}
	glog.V(5).Infof("rsmi_compute_process_info_by_pid_get:%v,errorString: %v, cproc:%v", ret, nil, cproc)

	proc = ProcessInfos{
		ProcessID:   uint32(cproc.process_id),
		Pasid:       uint32(cproc.pasid),
		VramUsage:   uint64(cproc.vram_usage),
		SdmaUsage:   uint64(cproc.sdma_usage),
		CuOccupancy: uint32(cproc.cu_occupancy),
	}
	glog.V(5).Infof("RsmiComputeProcessInfoByPidGet: %v", proc)
	return
}

// GetProcessInfoByDevice 是封装的函数，用于调用 rsmi_compute_process_info_by_device_get
func rsmiProcessInfoByDevice(pid int, dvInd int) (proc ProcessInfos, err error) {
	var cproc C.rsmi_process_info_t

	ret := C.rsmi_compute_process_info_by_device_get(C.uint32_t(pid), C.uint32_t(dvInd), &cproc)
	if err = errorString(ret); err != nil {
		return proc, fmt.Errorf("Error GetProcessInfoByDevice:%v", err)
	}

	proc = ProcessInfos{
		ProcessID:   uint32(cproc.process_id),
		Pasid:       uint32(cproc.pasid),
		VramUsage:   uint64(cproc.vram_usage),
		SdmaUsage:   uint64(cproc.sdma_usage),
		CuOccupancy: uint32(cproc.cu_occupancy),
	}
	glog.V(5).Infof("GetProcessInfoByDevice: %v", proc)
	return
}

// rsmiComputeProcessGpusGet 获取进程当前正在使用的设备索引
func rsmiComputeProcessGpusGet(pid int) (dvIndices []int, err error) {
	var cnumDevices C.uint32_t
	cPid := C.uint32_t(pid)

	// 第一次调用以获取numDevices的值
	ret := C.rsmi_compute_process_gpus_get(cPid, nil, &cnumDevices)
	glog.V(5).Infof(
		"rsmi_compute_process_gpus_get initial: pid=%d, ret=%d, num_devices=%d",
		pid, ret, cnumDevices,
	)
	if err := errorString(ret); err != nil {
		return dvIndices, fmt.Errorf("Error in RSMIComputeProcessGPUsGet (initial call): %s", err)
	}

	// 没有正在使用设备时不需要进行第二次调用，避免访问空切片的首元素。
	if cnumDevices == 0 {
		glog.V(5).Infof("rsmi_compute_process_gpus_get result: pid=%d, devices=[]", pid)
		return []int{}, nil
	}

	// 创建一个大小为numDevices的切片
	requestedDevices := cnumDevices
	dvIndicesC := make([]C.uint32_t, int(requestedDevices))
	// 第二次调用以获取实际的数据
	ret = C.rsmi_compute_process_gpus_get(cPid, &dvIndicesC[0], &cnumDevices)
	glog.V(5).Infof(
		"rsmi_compute_process_gpus_get data: pid=%d, ret=%d, requested_devices=%d, returned_devices=%d, raw_indices=%v",
		pid, ret, requestedDevices, cnumDevices, dvIndicesC,
	)
	if err := errorString(ret); err != nil {
		return nil, fmt.Errorf("Error in RSMIComputeProcessGPUsGet: %s", err)
	}
	if cnumDevices > requestedDevices {
		return nil, fmt.Errorf(
			"invalid device count returned for PID %d: %d, requested %d",
			pid, cnumDevices, requestedDevices,
		)
	}

	// 将C数组转换为Go切片
	dvIndices = make([]int, int(cnumDevices))
	for i := 0; i < int(cnumDevices); i++ {
		dvIndices[i] = int(dvIndicesC[i])
	}
	glog.V(5).Infof("RSMIComputeProcessGPUsGet: pid=%d, devices=%v", pid, dvIndices)
	return
}

// rsmiDevSupportedFuncIteratorOpen 获取设备支持RSMI函数的函数名迭代器
func rsmiDevSupportedFuncIteratorOpen(dvInd int) (iterHandle RSMIFuncIDIterHandle, err error) {
	var handle C.rsmi_func_id_iter_handle_t
	ret := C.rsmi_dev_supported_func_iterator_open(C.uint32_t(dvInd), &handle)
	if err = errorString(ret); err != nil {
		return iterHandle, fmt.Errorf("Error rsmiDevSupportedFuncIteratorOpen: %s", err)
	}
	iterHandle = RSMIFuncIDIterHandle(handle)
	return
}

// rsmiDevSupportedVariantIteratorOpen 获取给定句柄的变体迭代器
func rsmiDevSupportedVariantIteratorOpen(iterHandle RSMIFuncIDIterHandle) (handle RSMIFuncIDIterHandle, err error) {
	var chandle C.rsmi_func_id_iter_handle_t
	ret := C.rsmi_dev_supported_variant_iterator_open(C.rsmi_func_id_iter_handle_t(iterHandle), &chandle)
	if err = errorString(ret); err != nil {
		return iterHandle, fmt.Errorf("Error rsmiDevSupportedVariantIteratorOpen: %s", err)
	}
	handle = RSMIFuncIDIterHandle(chandle)
	return
}

// rsmiFuncIterNext 推进函数标识符迭代器
func rsmiFuncIterNext(handle RSMIFuncIDIterHandle) (err error) {
	ret := C.rsmi_func_iter_next(C.rsmi_func_id_iter_handle_t(handle))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmiFuncIterNext:%s", err)
	}
	return
}

// rsmiDevSupportedFuncIteratorClose 关闭变量迭代器句柄
func rsmiDevSupportedFuncIteratorClose(handle RSMIFuncIDIterHandle) (err error) {
	cHandle := C.rsmi_func_id_iter_handle_t(handle)
	ret := C.rsmi_dev_supported_func_iterator_close(&cHandle)
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmiDevSupportedFuncIteratorClose:%s", err)
	}
	return
}

// rsmiFuncIterValueGet 获取与函数/变量迭代器相关联的值
//func rsmiFuncIterValueGet(handle RSMIFuncIDIterHandle) (value RSMIFuncIDValue, err error) {
//	var cvalue C.rsmi_func_id_value_t
//	// 调用C函数
//	ret := C.rsmi_func_iter_value_get(C.rsmi_func_id_iter_handle_t(handle), &cvalue)
//	if err = errorString(ret); err != nil {
//		return value, fmt.Errorf("Error rsmiFuncIterValueGet:%s", err)
//	}
//	value.ID = uint64(cvalue.id)
//	value.Name = C.GoString((*C.char)(unsafe.Pointer(cvalue.name)))
//	value.MemoryType = RSMIMemoryType(*(*C.rsmi_memory_type_t)(unsafe.Pointer(&cvalue)))
//	value.TempMetric = RSMITemperatureMetric(*(*C.rsmi_temperature_metric_t)(unsafe.Pointer(&cvalue)))
//	value.EventType = RSMIEventType(*(*C.rsmi_event_type_t)(unsafe.Pointer(&cvalue)))
//	value.EventGroup = RSMIEventGroup(*(*C.rsmi_event_group_t)(unsafe.Pointer(&cvalue)))
//	value.ClkType = RSMIClkType(*(*C.rsmi_clk_type_t)(unsafe.Pointer(&cvalue)))
//	value.FwBlock = RSMIFwBlock(*(*C.rsmi_fw_block_t)(unsafe.Pointer(&cvalue)))
//	value.GpuBlock = RSMIGpuBlock(*(*C.rsmi_gpu_block_t)(unsafe.Pointer(&cvalue)))
//	return
//}
/*************事件************/
// rsmiEventNotificationInit 准备收集GPU事件通知 初始化事件通知
func rsmiEventNotificationInit(deInd int) (err error) {
	ret := C.rsmi_event_notification_init(C.uint32_t(deInd))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Rrror rsmiEventNotificationInit:%s", err)
	}
	return
}

// rsmiEventNotificationMaskSet 设置设备指定要收集的事件。设置事件通知掩码
func rsmiEventNotificationMaskSet(dvInd int, mask int64) (err error) {
	ret := C.rsmi_event_notification_mask_set(C.uint32_t(dvInd), C.uint64_t(mask))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Rrror rsmiEventNotificationMaskSet:%s", err)
	}
	return
}

// rsmiEventNotificationGet 收集事件通知，等待指定时间
func rsmiEventNotificationGet(timeoutMs int) (numElem int, datas []RSMIEEvtNotificationData, err error) {
	var cnumElen C.uint32_t
	ret := C.rsmi_event_notification_get(C.int(timeoutMs), &cnumElen, nil)
	if err = errorString(ret); err != nil {
		return 0, nil, fmt.Errorf("Error rsmiEventNotificationGet,numElem:%s", err)
	}
	numElem = int(cnumElen)
	cdatas := make([]C.rsmi_evt_notification_data_t, numElem)
	ret = C.rsmi_event_notification_get(C.int(timeoutMs), &cnumElen, (*C.rsmi_evt_notification_data_t)(unsafe.Pointer(&cdatas[0])))
	if err = errorString(ret); err != nil {
		return numElem, nil, fmt.Errorf("Error rsmiEventNotificationGet,datas:%s", err)
	}
	datas = make([]RSMIEEvtNotificationData, numElem)
	for i, data := range cdatas {
		datas[i] = RSMIEEvtNotificationData{
			DvInd:   uint32(data.dv_ind),
			Event:   RSMIEvtNotificationType(data.event),
			Message: *(*[64]byte)(unsafe.Pointer(&data.message)),
		}
	}
	return
}

// rsmiEventNotificationStop 关闭任何文件句柄并释放由GPU事件通知使用的任何资源。
func rsmiEventNotificationStop(dvInd int) (err error) {
	ret := C.rsmi_event_notification_stop(C.uint32_t(dvInd))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmiEventNotificationStop:%s", err)
	}
	return
}

// 打印事件列表方法
func printEventList(device int, delay int, eventList []string) {
	print2DArray([][]string{{"DEVICE", "TIME", "TYPE", "DESCRIPTION"}})
	mask := int64(0)

	if err := rsmiEventNotificationInit(device); err != nil {
		glog.Error(device, "Unable to initialize event notifications.")
		return
	}

	for _, eventType := range eventList {
		for i, name := range notificationTypeNames {
			if strings.ToUpper(eventType) == name {
				mask |= 1 << uint(i)
			}
		}
	}

	if err := rsmiEventNotificationMaskSet(device, mask); err != nil {
		glog.Error(device, "Unable to set event notification mask.")
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				_, datas, err := rsmiEventNotificationGet(delay)
				if err != nil {
					continue
				}
				for _, data := range datas {
					if len(data.Message) > 0 {
						print2DArray([][]string{
							{fmt.Sprintf("GPU[%d]", data.DvInd), time.Now().Format("2006-01-02 15:04:05"), notificationTypeNames[data.Event-1], string(data.Message[:])},
						})
					}
				}
				time.Sleep(time.Millisecond * time.Duration(delay))
			}
		}
	}()

	<-stop
	fmt.Println("Exiting...")
}
