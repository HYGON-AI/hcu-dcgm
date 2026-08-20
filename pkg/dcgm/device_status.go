/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L/opt/hyhal/lib -Wl,-rpath,/opt/hyhal/lib -lrocm_smi64 -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdint.h>
#include <kfd_ioctl.h>
#include <rocm_smi64Config.h>
#include <rocm_smi.h>
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/golang/glog"
)

// rsmiDevTempMetricGet 获取设备的温度度量值 *
func rsmiDevTempMetricGet(dvInd int, sensorType int, metric RSMITemperatureMetric) (temp int64, err error) {
	var temperature C.int64_t
	ret := C.rsmi_dev_temp_metric_get(C.uint32_t(dvInd), C.uint32_t(sensorType), C.rsmi_temperature_metric_t(metric), &temperature)
	//glog.V(5).Infof("rsmi_dev_temp_metric_get ret:%v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("rsmiDevTempMetricGet:%s", err)
	}
	temp = int64(temperature)
	return
}

// rsmiDevVoltMetricGet 获取设备的电压度量值
func rsmiDevVoltMetricGet(dvInd int, voltageType RSMIVoltageType, metric RSMIVoltageMetric) int64 {
	var voltage C.int64_t
	C.rsmi_dev_volt_metric_get(C.uint32_t(dvInd), C.rsmi_voltage_type_t(voltageType), C.rsmi_voltage_metric_t(metric), &voltage)
	return int64(voltage)
}

// rsmiDevFanSpeedSet 设置设备风扇转速，以rpm为单位
func rsmiDevFanSpeedSet(dvInd, sensorInd int, speed int64) (err error) {
	ret := C.rsmi_dev_fan_speed_set(C.uint32_t(dvInd), C.uint32_t(sensorInd), C.uint64_t(speed))
	glog.V(5).Infof("rsmi_dev_fan_speed_set_ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_fan_speed_set: %s", err)
	}
	return nil
}

// rsmiDevBusyPercentGet 获取设备设备忙碌时间百分比
func rsmiDevBusyPercentGet(dvInd int) (busyPercent int, err error) {
	var cbusyPercent C.uint32_t
	ret := C.rsmi_dev_busy_percent_get(C.uint32_t(dvInd), &cbusyPercent)
	//glog.V(5).Infof("rsmi_dev_busy_percent_get ret:%v ,retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("Error rsmi_dev_busy_percent_get:%s", err)
	}
	busyPercent = int(cbusyPercent)
	return busyPercent, nil
}

// rsmiUtilizationCountGet 获取设备利用率计数器
func rsmiUtilizationCountGet(dvInd int, utilizationCounters []UtilizationCounter, count int) (timestamp int64, err error) {
	// 转换 Go 结构体数组到 C 结构体数组
	// 注意：rsmi_utilization_count_get 的参数是 [inout]：
	// - [in]  type 字段：调用者必须设置想查询的计数器类型
	// - [out] value/fine_value/fine_value_count：驱动填充，必须零初始化
	cUtilizationCounters := make([]C.rsmi_utilization_counter_t, len(utilizationCounters))
	for i, uc := range utilizationCounters {
		cUtilizationCounters[i] = C.rsmi_utilization_counter_t{
			_type: C.RSMI_UTILIZATION_COUNTER_TYPE(uc.Type),
			// value、fine_value、fine_value_count 保持零初始化，等待驱动填充
		}
	}

	var ctimestamp C.uint64_t
	// 调用 C 函数
	ret := C.rsmi_utilization_count_get(
		C.uint32_t(dvInd),
		&cUtilizationCounters[0],
		C.uint32_t(count),
		&ctimestamp,
	)
	glog.V(5).Infof("rsmi_utilization_count_get ret:%v ,retstr:%v ", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("Error rsmi_utilization_count_get:%s", err)
	}
	// 更新 Go 结构体数组中的值
	for i := range utilizationCounters {
		utilizationCounters[i].Value = uint64(cUtilizationCounters[i].value)
	}
	glog.V(5).Infof("utilizationCounters:%v,timestamp:%v", utilizationCounters, int64(ctimestamp))

	return int64(ctimestamp), nil
}

// rsmiDevPerfLevelGet 获取设备的性能级别
func rsmiDevPerfLevelGet(dvInd int) (perf DevPerfLevel, err error) {
	var cPerfLevel C.rsmi_dev_perf_level_t
	ret := C.rsmi_dev_perf_level_get(C.uint32_t(dvInd), &cPerfLevel)
	if err = errorString(ret); err != nil {
		return DevPerfLevel(cPerfLevel), fmt.Errorf("Error rsmi_dev_perf_level_get:%s", err)
	}
	perf = DevPerfLevel(cPerfLevel)
	glog.V(5).Infof("dev_perf_level:", perf)
	return perf, nil
}

// rsmiPerfDeterminismModeSet 设置设备的性能确定性模式
func rsmiPerfDeterminismModeSet(dvInd int, clkValue int64) (err error) {
	ret := C.rsmi_perf_determinism_mode_set(C.uint32_t(dvInd), C.uint64_t(clkValue))
	glog.V(5).Infof("dev_perf_determinism_mode ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_perf_determinism_mode_set:%s", err)
	}
	return
}

// rsmiDevOverdriveLevelGet 获取设备的超速百分比
func rsmiDevOverdriveLevelGet(dvInd int) (od int, err error) {
	var cod C.uint32_t
	ret := C.rsmi_dev_overdrive_level_get(C.uint32_t(dvInd), &cod)
	glog.V(5).Infof("rsmi_dev_overdrive_level_get:ret:%v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return int(cod), fmt.Errorf("Error rsmi_dev_overdrive_level_get:%s", err)
	}
	od = int(cod)
	glog.V(5).Infof("rsmiDevOverdriveLevelGet od:%v", od)
	return
}

// rsmiDevGpuClkFreqGet 获取设备系统时钟速度列表
func rsmiDevGpuClkFreqGet(dvInd int, clkType RSMIClkType) (freq Frequencies, err error) {
	var cfreq C.rsmi_frequencies_t
	// 获取实际数据
	ret := C.rsmi_dev_gpu_clk_freq_get(
		C.uint32_t(dvInd),
		C.rsmi_clk_type_t(clkType),
		&cfreq,
	)
	glog.V(5).Infof("rsmi_dev_gpu_clk_freq_get:ret:%v, retStr:%v", ret, errorString(ret))
	glog.V(5).Infof("cfreq: %v ", cfreq)
	if err = errorString(ret); err != nil {
		return freq, fmt.Errorf("clock type not supported: %s", errorString(ret))
	}
	// 类型转换
	freq = Frequencies{
		HasDeepSleep: bool(cfreq.has_deep_sleep),
		NumSupported: uint32(cfreq.num_supported),
		Current:      uint32(cfreq.current),
		Frequency:    *(*[33]uint64)(unsafe.Pointer(&cfreq.frequency)),
	}
	glog.V(5).Infof("freq: %v", freq)
	//
	if freq.NumSupported == 0 {
		return freq, fmt.Errorf("no supported frequencies (num_supported=0)")
	}
	if freq.Current >= freq.NumSupported {
		return freq, fmt.Errorf("invalid current index %d (num_supported=%d)",
			freq.Current, freq.NumSupported)
	}
	if int(freq.Current) >= len(freq.Frequency) {
		return freq, fmt.Errorf("current index %d exceeds array size %d",
			freq.Current, len(freq.Frequency))
	}
	glog.V(5).Infof("freq: %v", freq)
	return freq, nil
}

// rsmiDevOdVoltInfoGet 获取设备电压/频率曲线信息
func rsmiDevOdVoltInfoGet(dvInd int) (odv OdVoltFreqData, err error) {
	var codv C.rsmi_od_volt_freq_data_t
	ret := C.rsmi_dev_od_volt_info_get(C.uint32_t(dvInd), &codv)
	glog.V(5).Infof("rsmi_dev_od_volt_info_get ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return odv, fmt.Errorf("Error rsmi_dev_od_volt_info_get:%s", err)
	}
	odv = OdVoltFreqData{
		CurrMclkRange: RSMIRange{
			LowerBound: uint64(codv.curr_sclk_range.lower_bound),
			UpperBound: uint64(codv.curr_sclk_range.upper_bound),
		},
		CurrSclkRange: RSMIRange{
			LowerBound: uint64(codv.curr_mclk_range.lower_bound),
			UpperBound: uint64(codv.curr_mclk_range.upper_bound),
		},
		SclkFreqLimits: RSMIRange{
			LowerBound: uint64(codv.sclk_freq_limits.lower_bound),
			UpperBound: uint64(codv.sclk_freq_limits.upper_bound),
		},
		MclkFreqLimits: RSMIRange{
			LowerBound: uint64(codv.mclk_freq_limits.lower_bound),
			UpperBound: uint64(codv.mclk_freq_limits.upper_bound),
		},
		Curve:      RSMIOdVoltCurve{},
		NumRegions: uint32(codv.num_regions),
	}
	for i := 0; i < len(codv.curve.vc_points); i++ {
		odv.Curve.VcPoints[i] = RSMIOdVddcPoint{
			Frequency: uint64(codv.curve.vc_points[i].frequency),
			Voltage:   uint64(codv.curve.vc_points[i].voltage),
		}
	}
	return
}

// rsmiDevGpuMetricsInfoGet 获取gpu度量信息
func rsmiDevGpuMetricsInfoGet(dvInd int) (gpuMetrics RSMIGPUMetrics, err error) {
	var cgpuMetrics C.rsmi_gpu_metrics_t
	ret := C.rsmi_dev_gpu_metrics_info_get(C.uint32_t(dvInd), &cgpuMetrics)
	if err = errorString(ret); err != nil {
		return gpuMetrics, fmt.Errorf("Error rsmi_dev_gpu_metrics_info_get:%s", err)
	}
	gpuMetrics = RSMIGPUMetrics{
		CommonHeader: MetricsTableHeader{
			StructureSize:   uint16(cgpuMetrics.common_header.structure_size),
			FormatRevision:  uint8(cgpuMetrics.common_header.format_revision),
			ContentRevision: uint8(cgpuMetrics.common_header.content_revision),
		},
		TemperatureEdge:        uint16(cgpuMetrics.temperature_edge),
		TemperatureHotspot:     uint16(cgpuMetrics.temperature_hotspot),
		TemperatureMem:         uint16(cgpuMetrics.temperature_mem),
		TemperatureVRGfx:       uint16(cgpuMetrics.temperature_vrgfx),
		TemperatureVRSoc:       uint16(cgpuMetrics.temperature_vrsoc),
		TemperatureVRMem:       uint16(cgpuMetrics.temperature_vrmem),
		AverageGfxActivity:     uint16(cgpuMetrics.average_gfx_activity),
		AverageUmcActivity:     uint16(cgpuMetrics.average_umc_activity),
		AverageMmActivity:      uint16(cgpuMetrics.average_mm_activity),
		AverageSocketPower:     uint16(cgpuMetrics.average_socket_power),
		EnergyAccumulator:      uint64(cgpuMetrics.energy_accumulator),
		SystemClockCounter:     uint64(cgpuMetrics.system_clock_counter),
		AverageGfxclkFrequency: uint16(cgpuMetrics.average_gfxclk_frequency),
		AverageSocclkFrequency: uint16(cgpuMetrics.average_socclk_frequency),
		AverageUclkFrequency:   uint16(cgpuMetrics.average_uclk_frequency),
		AverageVclk0Frequency:  uint16(cgpuMetrics.average_vclk0_frequency),
		AverageDclk0Frequency:  uint16(cgpuMetrics.average_dclk0_frequency),
		AverageVclk1Frequency:  uint16(cgpuMetrics.average_vclk1_frequency),
		AverageDclk1Frequency:  uint16(cgpuMetrics.average_dclk1_frequency),
		CurrentGfxclk:          uint16(cgpuMetrics.current_gfxclk),
		CurrentSocclk:          uint16(cgpuMetrics.current_socclk),
		CurrentUclk:            uint16(cgpuMetrics.current_uclk),
		CurrentVclk0:           uint16(cgpuMetrics.current_vclk0),
		CurrentDclk0:           uint16(cgpuMetrics.current_dclk0),
		CurrentVclk1:           uint16(cgpuMetrics.current_vclk1),
		CurrentDclk1:           uint16(cgpuMetrics.current_dclk1),
		ThrottleStatus:         uint32(cgpuMetrics.throttle_status),
		CurrentFanSpeed:        uint16(cgpuMetrics.current_fan_speed),
		PcieLinkWidth:          uint16(cgpuMetrics.pcie_link_width),
		PcieLinkSpeed:          uint16(cgpuMetrics.pcie_link_speed),
		//Padding:                uint16(cgpuMetrics.padding),
		GfxActivityAcc: uint32(cgpuMetrics.gfx_activity_acc),
		MemActivityAcc: uint32(cgpuMetrics.mem_activity_acc),
	}
	const hbmInstances = 4 // 确保跟 C 端的 RSMI_NUM_HBM_INSTANCES 一致
	for i := 0; i < hbmInstances; i++ {
		gpuMetrics.TempetureHBM[i] = uint16(cgpuMetrics.temperature_hbm[C.int(i)])
	}
	glog.V(5).Infof("rsmi_dev_gpu_metrics_info_get:%s", gpuMetrics)
	return
}

// rsmiDevEccStatusGet 获取GPU块的ECC状态
func rsmiDevEccStatusGet(dvInd int, block RSMIGpuBlock) (state RSMIRasErrState, err error) {
	glog.V(5).Infof("rsmiDevEccStatusGet: %d,%d", dvInd, block)
	var sstate C.rsmi_ras_err_state_t
	ret := C.rsmi_dev_ecc_status_get(C.uint32_t(dvInd), C.rsmi_gpu_block_t(block), &sstate)
	glog.V(5).Infof("rsmi_dev_ecc_status_get ret:%v", ret)
	if err = errorString(ret); err != nil {
		return state, fmt.Errorf("Error rsmi_dev_ecc_status_get:%s", err)
	}
	state = RSMIRasErrState(sstate)
	glog.V(5).Infof("rsmiDevEccStatusGet:%v", sstate)
	return
}

// rsmiDevEccCountGet 获取GPU块的错误计数
func rsmiDevEccCountGet(dvInd int, gpuBlock RSMIGpuBlock) (errorCount RSMIErrorCount, err error) {
	glog.V(5).Infof("dvInd:%v , RSMIGpuBlock:%v", dvInd, gpuBlock)
	var cerrorCount C.rsmi_error_count_t
	ret := C.rsmi_dev_ecc_count_get(C.uint32_t(dvInd), C.rsmi_gpu_block_t(gpuBlock), &cerrorCount)
	glog.V(5).Infof("rsmiDevEccCountGet:%v,ret retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return errorCount, fmt.Errorf("Error rsmi_dev_ecc_count_get:%s", err)
	}
	errorCount = RSMIErrorCount{
		CorrectableErr:   uint64(cerrorCount.correctable_err),
		UncorrectableErr: uint64(cerrorCount.uncorrectable_err),
	}
	glog.V(5).Infof("HCUBlockType:%v, DevEccCount:%v", gpuBlock, errorCount)
	return
}

// rsmiDevEccEnabledGet 获取已启用的ECC位掩码
func rsmiDevEccEnabledGet(dvInd int) (enabledBlocks int64, err error) {
	var cenabledBlocks C.uint64_t
	ret := C.rsmi_dev_ecc_enabled_get(C.uint32_t(dvInd), &cenabledBlocks)
	if err = errorString(ret); err != nil {
		return enabledBlocks, fmt.Errorf("Error rsmi_dev_ecc_enabled_get:%s", err)
	}
	enabledBlocks = int64(cenabledBlocks)
	glog.V(5).Infof("HCUBlockType:%v", enabledBlocks)
	return
}

// rsmiDevMemOverdriveLevelGet 获取显存 overdrive 百分比（0-20）
// 对应 C 接口 rsmi_dev_mem_overdrive_level_get
func rsmiDevMemOverdriveLevelGet(dvInd int) (od int, err error) {
	var cod C.uint32_t
	ret := C.rsmi_dev_mem_overdrive_level_get(C.uint32_t(dvInd), &cod)
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("rsmi_dev_mem_overdrive_level_get: %s", err)
	}
	od = int(cod)
	return
}

// rsmiDevNodeIdGet 获取 KFD node ID
// 对应 C 接口 rsmi_dev_node_id_get
func rsmiDevNodeIdGet(dvInd int) (nodeId uint32, err error) {
	var cid C.uint32_t
	ret := C.rsmi_dev_node_id_get(C.uint32_t(dvInd), &cid)
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("rsmi_dev_node_id_get: %s", err)
	}
	nodeId = uint32(cid)
	return
}

// rsmiDevCuNumGet 获取 CU（Compute Unit）数量
// 对应 C 接口 rsmi_dev_cu_num_get；cgo 编译走 C 分支（指针参数）
func rsmiDevCuNumGet(dvInd int) (cuNum int, err error) {
	var ccnt C.int
	ret := C.rsmi_dev_cu_num_get(C.uint32_t(dvInd), &ccnt)
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("rsmi_dev_cu_num_get: %s", err)
	}
	cuNum = int(ccnt)
	return
}

// rsmiIsP2PAccessible 查询两个设备之间的 P2P 可达性
// 对应 C 接口 rsmi_is_P2P_accessible
func rsmiIsP2PAccessible(srcDvInd, dstDvInd int) (accessible bool, err error) {
	var cval C.bool
	ret := C.rsmi_is_P2P_accessible(C.uint32_t(srcDvInd), C.uint32_t(dstDvInd), &cval)
	if err = errorString(ret); err != nil {
		return false, fmt.Errorf("rsmi_is_P2P_accessible: %s", err)
	}
	accessible = bool(cval)
	return
}

// rsmiHyVersionGet 获取海光扩展版本（hy_major/hy_minor）
// 对应 C 接口 rsmi_hy_version_get；返回 HyVersionInfo
func rsmiHyVersionGet() (ver HyVersionInfo, err error) {
	var cv C.rsmi_hy_version_t
	ret := C.rsmi_hy_version_get(&cv)
	if err = errorString(ret); err != nil {
		return ver, fmt.Errorf("rsmi_hy_version_get: %s", err)
	}
	ver = HyVersionInfo{
		Major: uint32(cv.hy_major),
		Minor: uint32(cv.hy_minor),
	}
	return
}
