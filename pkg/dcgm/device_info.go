/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L/opt/hyhal/lib -Wl,-rpath,/opt/hyhal/lib -lrocm_smi64 -lhydmi -Wl,--unresolved-symbols=ignore-in-object-files
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
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/golang/glog"
)

// rsmiNumMonitorDevices 获取gpu数量 *
func rsmiNumMonitorDevices() (gpuNum int, err error) {
	var p C.uint
	ret := C.rsmi_num_monitor_devices(&p)
	glog.V(5).Infof("rsmi_num_monitor_devices_ret:%v, retStr : %v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return 0, fmt.Errorf("Error go_rsmi_num_monitor_devices_ret: %s", err)
	}
	gpuNum = int(p)
	glog.V(5).Infof("rsmiNumMonitorDevices:%v", gpuNum)
	return gpuNum, nil
}

// rsmiDevSkuGet 获取设备sku
func rsmiDevSkuGet(dvInd int) (sku int, err error) {
	var csku C.uint16_t
	ret := C.rsmi_dev_sku_get(C.uint32_t(dvInd), &csku)
	if err = errorString(ret); err != nil {
		return sku, err
	}
	sku = int(csku)
	glog.V(5).Infof("rsmiDevSkuGet:", sku)
	return
}

// rsmiDevVendorIdGet 获取设备供应商id
func rsmiDevVendorIdGet(dvInd int) uint {
	var vid C.uint16_t
	C.rsmi_dev_vendor_id_get(C.uint32_t(dvInd), &vid)
	return uint(vid)
}

// rsmiDevIdGet 获取设备类型标识id
func rsmiDevIdGet(dvInd int) (id int, err error) {
	var cid C.uint16_t
	ret := C.rsmi_dev_id_get(C.uint32_t(dvInd), &cid)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmiDevIdGet:%v,retStr:%v", err, errorString(ret))
		return 0, fmt.Errorf("Error rsmiDevIdGet:%v", err)
	}
	glog.V(5).Infof("rsmiDevIdGet cid:%v", cid)
	id = int(cid)
	//glog.V(5).Infof("rsmiDevIdGet: %v", id)
	return
}

// rsmiDevNameGet 获取设备名称
func rsmiDevNameGet(dvInd int) (nameStr string, err error) {
	name := make([]C.char, uint32(256))
	ret := C.rsmi_dev_name_get(C.uint32_t(dvInd), &name[0], 256)
	if err = errorString(ret); err != nil {
		return nameStr, fmt.Errorf("Error go_rsmi_dev_name_get: %s", err)
	}
	nameStr = C.GoString(&name[0])
	//glog.V(5).Infof("rsmiDevNameGet:", nameStr)
	return
}

// rsmiDevBrandGet 获取设备品牌名称
func rsmiDevBrandGet(dvInd int) (brand string, err error) {
	brands := make([]C.char, uint32(256))
	C.rsmi_dev_brand_get(C.uint32_t(dvInd), &brands[0], 256)
	brand = C.GoString(&brands[0])
	glog.V(5).Infof("rsmiDevBrandGet:", brand)
	return
}

// rsmiDevVendorNameGet 获取设备供应商名称
func rsmiDevVendorNameGet(dvInd int) (bname string, err error) {
	cbname := make([]C.char, uint32(256))
	ret := C.rsmi_dev_vendor_name_get(C.uint32_t(dvInd), &cbname[0], 80)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_vendor_name_get:%v", err)
		return bname, fmt.Errorf("Error rsmi_dev_vendor_name_get:%v", err)
	}
	bname = C.GoString(&cbname[0])
	//glog.V(5).Infof("rsmiDevVendorNameGet:%v", bname)
	return
}

// rsmiDevVramVendorGet 获取设备显存供应商名称
func rsmiDevVramVendorGet(dvInd int) (result string, err error) {
	bname := make([]C.char, uint32(256))
	ret := C.rsmi_dev_vram_vendor_get(C.uint32_t(dvInd), &bname[0], 80)
	if err = errorString(ret); err != nil {
		return "", fmt.Errorf("Error rsmi_dev_vram_vendor_get:%s", err)
	}
	result = C.GoString(&bname[0])
	glog.V(5).Infof("rsmiDevVramVendorGet: %v", result)
	return
}

// rsmiDevSerialNumberGet 获取设备序列号
func rsmiDevSerialNumberGet(dvInd int) (serialNumber string, err error) {
	cserialNumber := make([]C.char, uint32(256))
	ret := C.rsmi_dev_serial_number_get(C.uint32_t(dvInd), &cserialNumber[0], 256)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_serial_number_get:%v, errstr:%v", err, errorString(ret))
		return "", fmt.Errorf("Error rsmi_dev_serial_number_get:%s", err)
	}
	serialNumber = C.GoString(&cserialNumber[0])
	//glog.V(5).Infof("Serial number: %v", serialNumber)
	return
}

// rsmiDevSubsystemIdGet 获取设备子系统id
func rsmiDevSubsystemIdGet(dvInd int) (subSystemId int, err error) {
	var id C.uint16_t
	ret := C.rsmi_dev_subsystem_id_get(C.uint32_t(dvInd), &id)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_subsystem_id_get:%v", err)
		return subSystemId, fmt.Errorf("Error rsmi_dev_subsystem_id_get:%s", err)
	}
	glog.V(5).Infof("rsmi_dev_subsystem_id_get:%v", id)
	subSystemId = int(id)
	return
}

// rsmiDevSubsystemNameGet 获取设备子系统名称
func rsmiDevSubsystemNameGet(dvInd int) (subSystemName string, err error) {
	csubSystemName := make([]C.char, uint32(256))
	ret := C.rsmi_dev_subsystem_name_get(C.uint32_t(dvInd), &csubSystemName[0], 256)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_subsystem_name_get:%v", err)
		return subSystemName, fmt.Errorf("Error rsmi_dev_subsystem_name_get:%s", err)
	}
	subSystemName = C.GoString(&csubSystemName[0])
	glog.V(5).Infof("rsmiDevSubsystemNameGet:%v", subSystemName)
	return
}

// rsmiDevDrmRenderMinorGet 获取设备drm次编号
func rsmiDevDrmRenderMinorGet(dvInd int) int {
	var id C.uint32_t
	C.rsmi_dev_drm_render_minor_get(C.uint32_t(dvInd), &id)
	return int(id)
}

// rsmiDevUniqueIdGet 获取设备唯一id
func rsmiDevUniqueIdGet(dvInd int) (uniqueId int64, err error) {
	var cuniqueId C.uint64_t
	ret := C.rsmi_dev_unique_id_get(C.uint32_t(dvInd), &cuniqueId)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_unique_id_get:%v, retstr:%v", ret, errorString(ret))
		return uniqueId, fmt.Errorf("Error rsmi_dev_unique_id_get:%s", err)
	}
	uniqueId = int64(cuniqueId)
	//glog.V(5).Infof("rsmiDevUniqueIdGet:%v", uint64(uniqueId))
	return
}

// rsmiDevSubsystemVendorIdGet 获取设备子系统供应商id
func rsmiDevSubsystemVendorIdGet(dvInd int) int {
	var id C.uint16_t
	C.rsmi_dev_subsystem_vendor_id_get(C.uint32_t(dvInd), &id)
	return int(id)
}

/****************************************** PCIe *********************************************/

// rsmiDevPciBandwidthGet 获取可用的pcie带宽列表
func rsmiDevPciBandwidthGet(dvInd int) (pcieBandwidth PcieBandwidth, err error) {
	var bandwidth C.rsmi_pcie_bandwidth_t
	ret := C.rsmi_dev_pci_bandwidth_get(C.uint32_t(dvInd), &bandwidth)
	if err = errorString(ret); err != nil {
		return pcieBandwidth, fmt.Errorf("Error rsmi_dev_pci_bandwidth_get%s", err)
	}
	// 先读取 transfer_rate 的基本字段
	tr := bandwidth.transfer_rate
	pcieBandwidth.TransferRate.NumSupported = uint32(tr.num_supported)
	pcieBandwidth.TransferRate.Current = uint32(tr.current)

	// frequency 逐项复制
	for i := 0; i < 33; i++ {
		pcieBandwidth.TransferRate.Frequency[i] = uint64(tr.frequency[i])
	}

	// lanes 逐项复制（长度是 33）
	for i := 0; i < 33; i++ {
		pcieBandwidth.Lanes[i] = uint32(bandwidth.lanes[i])
	}

	glog.V(5).Infof("RSMIPcieBandwidth (num_supported=%d): %v",
		pcieBandwidth.TransferRate.NumSupported,
		pcieBandwidth)

	return
}

// rsmiDevPciIdGet 获取唯一pci设备标识符
func rsmiDevPciIdGet(dvInd int) (bdfid int64, err error) {
	var cbdfid C.uint64_t
	ret := C.rsmi_dev_pci_id_get(C.uint32_t(dvInd), &cbdfid)
	glog.V(5).Infof("rsmi_dev_pci_id_get ret:%v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		glog.Errorf("rsmi_dev_pci_id_get err:%v", err.Error())
		return bdfid, err
	}
	bdfid = int64(cbdfid)
	glog.V(5).Infof("🚀🚀🚀dvInd: %v  rsmiDevPciIdGet bdfid:%v", dvInd, bdfid)
	return
}

// rsmiTopoNumaAffinityGet 获取与设备关联的numa节点
func rsmiTopoNumaAffinityGet(dvInd int) (namaNode int, err error) {
	var cnamaNode C.int32_t
	ret := C.rsmi_topo_numa_affinity_get(C.uint32_t(dvInd), &cnamaNode)
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_topo_numa_affinity_get ret:%v, retstr:%v", ret, errorString(ret))
		return namaNode, fmt.Errorf("Error rsmi_topo_numa_affinity_get:%s", err)
	}
	namaNode = int(cnamaNode)
	return
}

// rsmiDevPciThroughputGet 获取pcie流量信息
func rsmiDevPciThroughputGet(dvInd int) (sent int64, received int64, maxPktSz int64, err error) {
	var csent, creceived, cmaxpktsz C.uint64_t
	ret := C.rsmi_dev_pci_throughput_get(C.uint32_t(dvInd), &csent, &creceived, &cmaxpktsz)
	//glog.V(5).Infof("rsmi_dev_pci_throughput_get ret:%v ,retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return 0, 0, 0, fmt.Errorf("Error rsmi_dev_pci_throughput_get:%s", err)
	}
	//glog.V(5).Infof("csent: %v, creceived: %v, cmaxpktsz: %v", csent, creceived, cmaxpktsz)
	sent = int64(csent)
	received = int64(creceived)
	maxPktSz = int64(cmaxpktsz)
	//glog.V(5).Infof("sent: %v, received: %v, maxPktSz: %v", sent, received, maxPktSz)
	return
}

// rsmiDevPciReplayCounterGet 获取pcie重放计数
func rsmiDevPciReplayCounterGet(dvInd int) (counter int64, err error) {
	var ccounter C.uint64_t
	ret := C.rsmi_dev_pci_replay_counter_get(C.uint32_t(dvInd), &ccounter)
	if err = errorString(ret); err != nil {
		return counter, fmt.Errorf("Error rsmi_dev_pci_replay_counter_get:%s", err)
	}
	counter = int64(ccounter)
	glog.V(5).Infof("counter:%v", ccounter)
	return
}

// rsmiDevPciBandwidthSet 设置可使用的pcie带宽集
func rsmiDevPciBandwidthSet(dvInd int, bwBitmask int64) (err error) {
	ret := C.rsmi_dev_pci_bandwidth_set(C.uint32_t(dvInd), C.uint64_t(bwBitmask))
	glog.V(5).Infof("rsmiDevPciBandwidthSet, ret:%v ,retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmiDevPciBandwidthSet:%v", err)
	}
	return
}

/****************************************** Power *********************************************/

// rsmiDevPowerAveGet 获取设备平均功耗
func rsmiDevPowerAveGet(dvInd int, senserId int) (power int64, err error) {
	var cpower C.uint64_t
	ret := C.rsmi_dev_power_ave_get(C.uint32_t(dvInd), C.uint32_t(senserId), &cpower)
	//glog.V(5).Infof("rsmi_dev_power_ave_get, ret:%v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return power, fmt.Errorf("Error rsmiDevPowerAveGet:%v", err)
	}
	power = int64(cpower)
	return
}

// rsmiDevEnergyCountGet 获取设备的能量累加计数
func rsmiDevEnergyCountGet(dvInd int) (power uint64, counterResolution float32, timestamp uint64, err error) {
	var cPower C.uint64_t
	var cCounterResolution C.float
	var cTimestamp C.uint64_t
	ret := C.rsmi_dev_energy_count_get(C.uint32_t(dvInd), &cPower, &cCounterResolution, &cTimestamp)
	if ret != C.RSMI_STATUS_SUCCESS {
		return 0, 0, 0, fmt.Errorf("Error in rsmi_dev_energy_count_get: %s", errorString(ret))
	}
	return uint64(cPower), float32(cCounterResolution), uint64(cTimestamp), nil
}

// rsmiDevPowerCapGet 获取设备功率上限
func rsmiDevPowerCapGet(dvInd int, senserId int) (power int64, err error) {
	var cpower C.uint64_t
	ret := C.rsmi_dev_power_cap_get(C.uint32_t(dvInd), C.uint32_t(senserId), &cpower)
	//glog.V(5).Infof("rsmi_dev_power_cap_get ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return power, fmt.Errorf("Error rsmiDevPowerCapGet:%s", err)
	}
	power = int64(cpower)
	return
}

// rsmiDevPowerCapRangeGet 获取设备功率有效值范围
func rsmiDevPowerCapRangeGet(dvInd int, senserId int) (max, min int64, err error) {
	var cmax, cmin C.uint64_t
	ret := C.rsmi_dev_power_cap_range_get(C.uint32_t(dvInd), C.uint32_t(senserId), &cmax, &cmin)
	glog.V(5).Infof("rsmiDevPowerCapRangeGet ret:%v ,retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return max, min, fmt.Errorf("Error rsmiDevPowerCapRangeGet:%s", err)
	}
	glog.V(5).Infof("cmin:%v cmax:%v", cmin, cmax)
	max, min = int64(cmax), int64(cmin)
	glog.V(5).Infof("rsmiDevPowerCapRangeGet max:%v, min:%v", max, min)
	return
}

/****************************************** Memory *********************************************/

// rsmiDevMemoryTotalGet 获取设备内存总量 *
func rsmiDevMemoryTotalGet(dvInd int, memoryType RSMIMemoryType) (total int64, err error) {
	var ctotal C.uint64_t
	ret := C.rsmi_dev_memory_total_get(C.uint32_t(dvInd), C.rsmi_memory_type_t(memoryType), &ctotal)
	//glog.V(5).Infof("rsmi_dev_memory_total_get ret:%v ,retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return total, fmt.Errorf("Error rsmiDevMemoryTotalGet:%s", err)
	}
	total = int64(ctotal)
	//glog.V(5).Infof("memory_total:", total)
	return
}

// rsmiDevMemoryUsageGet 获取当前设备内存使用情况 *
func rsmiDevMemoryUsageGet(dvInd int, memoryType RSMIMemoryType) (used int64, err error) {
	var cused C.uint64_t
	ret := C.rsmi_dev_memory_usage_get(C.uint32_t(dvInd), C.rsmi_memory_type_t(memoryType), &cused)
	//glog.V(5).Infof("rsmi_dev_memory_usage_get ret:%v ,retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return used, fmt.Errorf("Error rsmiDevMemoryUsageGet:%s", err)
	}
	used = int64(cused)
	//glog.V(5).Infof("memory_used:", used)
	return
}

// rsmiDevMemoryBusyPercentGet 获取设备内存使用的百分比
func rsmiDevMemoryBusyPercentGet(dvInd int) (busyPercent int, err error) {
	var cbusyPercent C.uint32_t
	ret := C.rsmi_dev_memory_busy_percent_get(C.uint32_t(dvInd), &cbusyPercent)
	if err = errorString(ret); err != nil {
		return busyPercent, fmt.Errorf("Error rsmi_dev_memory_busy_percent_get:%s", err)
	}
	busyPercent = int(cbusyPercent)
	glog.V(5).Infof("busy_percent:", busyPercent)
	return
}

// rsmiDevMemoryReservedPagesGet 获取有关保留的(“已退休”)内存页的信息
func rsmiDevMemoryReservedPagesGet(dvInd int) (numPages int, records []RetiredPageRecord, err error) {
	var cnumPages C.uint32_t
	ret := C.rsmi_dev_memory_reserved_pages_get(C.uint32_t(dvInd), &cnumPages, nil)
	if ret != 0 {
		return 0, nil, fmt.Errorf("failed to get the number of pages, error code: %d", ret)
	}
	glog.V(5).Infof("cnumPages:", int(cnumPages))
	numPages = int(cnumPages)
	if numPages == 0 {
		return 0, nil, nil // No pages to retrieve
	}
	cRecords := make([]C.rsmi_retired_page_record_t, numPages)
	ret = C.rsmi_dev_memory_reserved_pages_get(C.uint32_t(dvInd), &cnumPages, (*C.rsmi_retired_page_record_t)(unsafe.Pointer(&cRecords[0])))
	if ret != 0 {
		return 0, nil, fmt.Errorf("failed to get the page records, error code: %d", ret)
	}

	records = make([]RetiredPageRecord, numPages)
	for i, rec := range cRecords {
		records[i] = RetiredPageRecord{
			PageAddress: uint64(rec.page_address),
			PageSize:    uint64(rec.page_size),
			Status:      MemoryPageStatus(rec.status),
		}
	}
	indent, _ := json.MarshalIndent(records, "", "  ")
	glog.V(5).Infof("records:", indent)
	return
}

// rsmiDevFanRpmsGet 获取设备的风扇速度，实际转速
func rsmiDevFanRpmsGet(dvInd, sensorInd int) (speed int64, err error) {
	var cspeed C.int64_t
	ret := C.rsmi_dev_fan_rpms_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &cspeed)
	glog.V(5).Infof("rsmi_dev_fan_rpms_get: ret:%v ,retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return speed, fmt.Errorf("Error rsmi_dev_fan_rpms_get:%s", err)
	}
	speed = int64(cspeed)
	glog.V(5).Infof("rsmi_dev_fan_rpms_get speed value: %v", speed)
	return
}

// rsmiDevFanSpeedGet 获取设备的风扇速度，相对速度值
func rsmiDevFanSpeedGet(dvInd, sensorInd int) (speed int64, err error) {
	var cspeed C.int64_t
	ret := C.rsmi_dev_fan_speed_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &cspeed)
	glog.V(5).Infof("rsmi_dev_fan_speed_get ret:%v ,retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return speed, fmt.Errorf("Error rsmiDevFanSpeedGet:%s", err)
	}
	speed = int64(cspeed)
	return
}

// rsmiDevFanSpeedMaxGet 获取设备的风扇速度，最大风速
func rsmiDevFanSpeedMaxGet(dvInd, sensorInd int) (maxSpeed int64, err error) {
	var cmaxSpeed C.uint64_t
	ret := C.rsmi_dev_fan_speed_max_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &cmaxSpeed)
	glog.V(5).Infof("rsmi_dev_fan_speed_max_get ret:%v ,retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return maxSpeed, fmt.Errorf("Error rsmiDevFanSpeedMaxGet:%s", err)
	}
	maxSpeed = int64(cmaxSpeed)
	return
}

// rsmiDevOdVoltCurveRegionsGet
func rsmiDevOdVoltCurveRegionsGet(dvInd int) (numRegions int, regions []RSMIFreqVoltRegion, err error) {
	var cnumRegions C.uint32_t
	ret := C.rsmi_dev_od_volt_curve_regions_get(C.uint32_t(dvInd), &cnumRegions, nil)
	if err = errorString(ret); err != nil {
		return 0, nil, fmt.Errorf("Error dev_od_volt_curve_regions_get: %v", err)
	}

	cbuffer := make([]C.rsmi_freq_volt_region_t, cnumRegions)
	ret = C.rsmi_dev_od_volt_curve_regions_get(C.uint32_t(dvInd), &cnumRegions, &cbuffer[0])
	if err = errorString(ret); err != nil {
		return 0, nil, fmt.Errorf("Error dev_od_volt_curve_regions_get: %v", err)
	}

	regions = make([]RSMIFreqVoltRegion, cnumRegions)
	for i := 0; i < int(cnumRegions); i++ {
		regions[i] = RSMIFreqVoltRegion{
			FreqRange: RSMIRange{
				LowerBound: uint64(cbuffer[i].freq_range.lower_bound),
				UpperBound: uint64(cbuffer[i].freq_range.upper_bound),
			},
			VoltRange: RSMIRange{
				LowerBound: uint64(cbuffer[i].volt_range.lower_bound),
				UpperBound: uint64(cbuffer[i].volt_range.upper_bound),
			},
		}
	}
	numRegions = int(cnumRegions)
	return
}

// rsmiDevPowerProfilePresetsGet 获取可用预设电源配置文件列表并指示当前活动的配置文件
func rsmiDevPowerProfilePresetsGet(dvInd, sensorInd int) (powerProfileStatus PowerProfileStatus, err error) {
	var cpowerProfileStatus C.rsmi_power_profile_status_t
	ret := C.rsmi_dev_power_profile_presets_get(C.uint32_t(dvInd), C.uint32_t(sensorInd), &cpowerProfileStatus)
	glog.V(5).Infof("rsmi_dev_power_profile_presets_get ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return powerProfileStatus, fmt.Errorf("Error dev_power_profile_presets_get:%s", err)
	}
	powerProfileStatus = PowerProfileStatus{
		AvailableProfiles: BitField(cpowerProfileStatus.available_profiles),
		Current:           PowerProfilePresetMasks(cpowerProfileStatus.current),
		NumProfiles:       uint32(cpowerProfileStatus.num_profiles),
	}
	glog.V(5).Infof("powerProfileStatus value: %v", powerProfileStatus)
	return
}

// rsmiVersionGet 获取当前运行的RSMI版本
func rsmiVersionGet() (version DevVersion, err error) {

	var cVersion C.rsmi_version_t
	ret := C.rsmi_version_get(&cVersion)
	if err = errorString(ret); err != nil {
		return version, fmt.Errorf("Error to get version: %s", err)
	}
	version = DevVersion{
		Major: uint32(cVersion.major),
		Minor: uint32(cVersion.minor),
		Patch: uint32(cVersion.patch),
		Build: C.GoString(cVersion.build),
	}
	glog.V(5).Infof("rsmiVersionGet:%v", version)
	return
}

// rsmiVersionStrGet 获取当前系统的驱动程序版本
func rsmiVersionStrGet(component SwComponent, len int) (varStr string, err error) {
	cvarStr := make([]C.char, len)
	ret := C.rsmi_version_str_get(C.rsmi_sw_component_t(component), (*C.char)(unsafe.Pointer(&cvarStr[0])), C.uint32_t(len))
	if err = errorString(ret); err != nil {
		return "", fmt.Errorf("Error rsmi_version_str_get:%s", err)
	}
	varStr = C.GoString(&cvarStr[0])
	return
}

// rsmiDevVbiosVersionGet 获取VBIOS版本
func rsmiDevVbiosVersionGet(dvInd, len int) (vbios string, err error) {
	cvbios := make([]C.char, len)
	ret := C.rsmi_dev_vbios_version_get(C.uint32_t(dvInd), &cvbios[0], C.uint32_t(len))
	if err = errorString(ret); err != nil {
		return vbios, fmt.Errorf("Error rsmi_dev_vbios_version_get:%s", err)
	}
	vbios = C.GoString(&cvbios[0])
	return
}

// rsmiDevFirmwareVersionGet 获取设备的固件版本
func rsmiDevFirmwareVersionGet(dvInd int, fwBlock RSMIFwBlock) (fwVersion int64, err error) {
	var cfwBlock C.uint64_t
	ret := C.rsmi_dev_firmware_version_get(C.uint32_t(dvInd), C.rsmi_fw_block_t(fwBlock), &cfwBlock)
	if err = errorString(ret); err != nil {
		return fwVersion, fmt.Errorf("Error rsmi_dev_firmware_version_get:%s", err)
	}
	fwVersion = int64(cfwBlock)
	return
}

// 获取DF带宽信息
func dfBandwidth(dvInd int, bandwidthType int) (dfBandwidthInfo DFBandwidthInfo, err error) {
	var dfBandwidth C.rsmi_df_bandwidth_info_t

	// 调用 C 函数
	ret := C.rsmi_dev_df_bandwidth_get(C.uint32_t(dvInd), C.RSMI_DF_BW_TYPE(bandwidthType), &dfBandwidth)
	glog.V(5).Infof("rsmi_dev_df_bandwidth_get:%v", ret)
	if err = errorString(ret); err != nil {
		err = fmt.Errorf("failed to get bandwidth info, error code: %d", int(ret))
		return
	}

	dfBandwidthInfo = DFBandwidthInfo{
		ReadBW:      float64(dfBandwidth.read_bw),
		WriteBW:     float64(dfBandwidth.write_bw),
		ReadWriteBW: float64(dfBandwidth.read_write_bw),
	}
	glog.V(5).Infof("bandwidth info: %v", dfBandwidthInfo)
	return
}

// 获取 xhcl 带宽信息
// dvInd: 设备索引 (int)
// linkId: 指定 link id (int)
// direction: 方向枚举值
// delay: 测量延迟
func xhclBandwidth(dvInd int, linkId int, direction int, delay int) (info XhclBandwidthInfo, err error) {
	glog.V(5).Infof("xhclBandwidth, dvInd:%v, linkId:%v, direction:%v, delay:%v ", dvInd, linkId, direction, delay)
	if err = ensureXhclBandwidth(); err != nil {
		return info, err
	}
	if dvInd < 0 {
		return info, fmt.Errorf("invalid dvInd: %d", dvInd)
	}
	//if linkId < 0 || linkId > linkAllID {
	//	return info, fmt.Errorf("invalid linkId: %d (must be 0..%d for single link, or %d for all links)",
	//		linkId, MAX_XHCL_LINK_NUM-1, linkAllID)
	//}
	var cinfo C.rsmi_xhcl_bandwidth_info_t

	ret := C.rsmi_dev_xhcl_bandwidth_get(
		C.uint32_t(dvInd),
		C.uint32_t(linkId),
		C.uint8_t(direction),
		C.int(delay),
		&cinfo,
	)
	glog.V(5).Infof("rsmi_dev_xhcl_bandwidth_get returned: %v", ret)
	glog.V(5).Infof("rsmi_xhcl_bandwidth_info_t: %v", cinfo)

	if RSMIStatus(ret) == RSMI_STATUS_NOT_SUPPORTED {
		glog.V(4).Infof(
			"xhcl bandwidth not supported: device=%d link=%d direction=%d",
			dvInd,
			linkId,
			direction,
		)
		return info, nil
	}

	if err := errorString(ret); err != nil {
		glog.Errorf("failed to get xhcl bandwidth info: %v", err)
		return info, fmt.Errorf("failed to get xhcl bandwidth info: %w", err)
	}

	// 逐个读取 cinfo.bw[i]
	for i := 0; i < MAX_XHCL_LINK_NUM; i++ {
		info.Bw[i] = float64(cinfo.bw[i])
	}

	glog.V(5).Infof("xhcl bandwidth info: %v", info)
	return
}

// 获取 UMC 带宽信息（
// dvInd: 设备索引
// chanId: channel 索引（0..MAX_UMC_CHAN_NUM-1）
// delay: 测量延迟
func umcBandwidth(dvInd int, chanId int, delay int) (info UMCBandwidthInfo, err error) {
	if err = ensureUmcBandwidth(); err != nil {
		return info, err
	}
	// 非负校验
	if dvInd < 0 {
		return info, fmt.Errorf("invalid dvInd: %d", dvInd)
	}
	if chanId < 0 || chanId > MAX_UMC_CHAN_NUM {
		return info, fmt.Errorf("invalid chanId: %d (must be 0..%d)", chanId, MAX_UMC_CHAN_NUM)
	}

	var cinfo C.rsmi_umc_bandwidth_info_t

	ret := C.rsmi_dev_umc_bandwidth_get(
		C.uint32_t(dvInd),
		C.uint32_t(chanId),
		C.int(delay),
		&cinfo,
	)
	glog.V(5).Infof("rsmi_dev_umc_bandwidth_get returned: %v", errorString(ret))
	glog.V(5).Infof("rsmi_dev_umc_bandwidth_get: %v", cinfo)
	if err = errorString(ret); err != nil {
		err = fmt.Errorf("failed to get umc bandwidth info: %w", err)
		return
	}

	// 逐个读取 C 结构数组字段并转换
	for i := 0; i < MAX_UMC_CHAN_NUM; i++ {
		info.ReadBW[i] = float64(cinfo.read_bw[i])
		info.WriteBW[i] = float64(cinfo.write_bw[i])
		info.ReadWriteBW[i] = float64(cinfo.read_write_bw[i])
	}
	glog.V(5).Infof("umc bandwidth info: %v", info)
	return
}

/*************************************VHCU******************************************/
// 设备数量
func dmiGetDeviceCount() (count int, err error) {
	var ccount C.int
	ret := C.dmiGetDeviceCount(&ccount)
	glog.V(5).Infof("dmiGetDeviceCount:%v,retmessage:%v", ret, dmiErrorString(ret))
	if err = dmiErrorString(ret); err != nil {
		return 0, fmt.Errorf("Error vDeviceCount:%s", err)
	}
	count = int(ccount)
	glog.V(5).Infof("dmiDeviceCount:%v", count)
	return
}

// 设备信息
func dmiGetDeviceInfo(dvInd int) (deviceInfo DMIDeviceInfo, err error) {
	var cdeviceInfo C.dmiDeviceInfo
	ret := C.dmiGetDeviceInfo(C.int(dvInd), &cdeviceInfo)
	glog.V(5).Infof("dmiDeviceInfo ret:%v,cdeviceInfo:%v", ret, cdeviceInfo)
	if err = dmiErrorString(ret); err != nil {
		return deviceInfo, fmt.Errorf("Error dmiGetDeviceInfo:%s", err)
	}
	// 创建一个新的变量来存储 name 字段
	var deviceName [DMI_NAME_SIZE]byte
	for i := 0; i < DMI_NAME_SIZE; i++ {
		deviceName[i] = byte(cdeviceInfo.name[i])
	}
	glog.V(5).Infof("deviceName:%v", deviceName)
	deviceInfo = DMIDeviceInfo{
		ComputeUnitCount: int(cdeviceInfo.compute_unit_count),
		GlobalMemSize:    uintptr(cdeviceInfo.global_mem_size),
		UsageMemSize:     uintptr(cdeviceInfo.usage_mem_size),
		DeviceID:         int(cdeviceInfo.device_id),
		Name:             ConvertASCIIToString(deviceName[:]),
	}

	glog.V(5).Infof("DeviceInfo: %v", deviceInfo)
	return
}

// 物理设备支持最大虚拟化设备数量
func dmiGetMaxVDeviceCount() (count int, err error) {
	var ccount C.int
	ret := C.dmiGetMaxVDeviceCount(&ccount)
	if err = dmiErrorString(ret); err != nil {
		return 0, fmt.Errorf("Error dmiGetMaxVDeviceCount:%s", err)
	}
	count = int(ccount)
	return
}

// 虚拟设备数量
func dmiGetVDeviceCount() (count int, err error) {
	var ccount C.int
	ret := C.dmiGetVDeviceCount(&ccount)
	glog.V(5).Infof("dmiGetVDeviceCount ret:%v,retmessage:%v", ret, dmiErrorString(ret))
	if err = dmiErrorString(ret); err != nil {
		return 0, fmt.Errorf("Error dmiGetVDeviceCount:%s", err)
	}
	count = int(ccount)
	glog.V(5).Infof("dmiGetVDeviceCount:%v", count)
	return
}

// 虚拟设备信息
func dmiGetVDeviceInfo(vDvInd int) (vDeviceInfo VDeviceInfo, err error) {
	var cvDeviceInfo C.dmiDeviceInfo
	ret := C.dmiGetVDeviceInfo(C.int(vDvInd), &cvDeviceInfo)
	/*glog.V(5).Infof("dmiGetVDeviceInfo ret:%v", ret)
	glog.V(5).Infof("name: %s", C.GoString(&cvDeviceInfo.name[0]))
	glog.V(5).Infof("compute_unit_count: %d", cvDeviceInfo.compute_unit_count)
	glog.V(5).Infof("global_mem_size: %d bytes", cvDeviceInfo.global_mem_size)
	glog.V(5).Infof("usage_mem_size: %d bytes", cvDeviceInfo.usage_mem_size)
	glog.V(5).Infof("container_id: %d", cvDeviceInfo.container_id)
	glog.V(5).Infof("device_id: %d", cvDeviceInfo.device_id)*/
	if err = dmiErrorString(ret); err != nil {
		glog.V(5).Infof("dmiGetVDeviceInfo vDvInd :%v, ret=%d, msg=%v", vDvInd, ret, err)
		return vDeviceInfo, err
	}
	//glog.V(5).Infof("dmiGetVDeviceInfo ret:%v,retmessage:%v", ret, dmiErrorString(ret))
	// 创建一个新的变量来存储 name 字段
	var deviceName [DMI_NAME_SIZE]byte
	for i := 0; i < DMI_NAME_SIZE; i++ {
		deviceName[i] = byte(cvDeviceInfo.name[i])
	}
	percent, _ := dmiGetVDevBusyPercent(vDvInd)
	vDeviceInfo = VDeviceInfo{
		VComputeUnitCount: int(cvDeviceInfo.compute_unit_count),
		VMemoryTotal:      uintptr(cvDeviceInfo.global_mem_size),
		VMemoryUsed:       uintptr(cvDeviceInfo.usage_mem_size),
		ContainerID:       uint64(cvDeviceInfo.container_id),
		DvInd:             int(cvDeviceInfo.device_id),
		VPercent:          percent,
		SubsystemTypeName: ConvertASCIIToString(deviceName[:]),
	}
	glog.V(5).Infof("vDeviceInfo: %v", vDeviceInfo)
	return
}

// 指定物理设备剩余的CU和内存
func dmiGetDeviceRemainingInfo(dvInd int) (cus, memories uint64, err error) {
	var ccus, cmemories C.size_t
	ret := C.dmiGetDeviceRemainingInfo(C.int(dvInd), &ccus, &cmemories)
	glog.V(5).Infof("dmiGetDeviceRemainingInfo ret:%v, retstr:%v", ret, dmiErrorString(ret))
	if err = dmiErrorString(ret); err != nil {
		return cus, memories, fmt.Errorf("Error dmiGetDeviceRemainingInfo:%s", err)
	}
	cus = uint64(ccus)
	memories = uint64(cmemories)
	glog.V(5).Infof("cus:%v,memories:%v", cus, memories)
	return
}

// 创建指定数量的虚拟设备
//
//	deviceID := 0
//	vdevCount := 2
//	vdevCUs := []int{4, 4}
//	vdevMemSize := []int{1024, 2048}
//
// 物理设备 ID: 0
//
//	├── 虚拟设备 1
//	│    ├── 计算单元: 4
//	│    └── 内存大小: 1024 字节
//	└── 虚拟设备 2
//	     ├── 计算单元: 4
//	     └── 内存大小: 2048 字节
func dmiCreateVDevices(dvInd int, vDevCount int, vDevCUs []int, vDevMemSize []int) (vdevIDs []int, err error) {
	if len(vDevCUs) != vDevCount || len(vDevMemSize) != vDevCount {
		return vdevIDs, fmt.Errorf("Invalid args")
	}

	fmt.Printf("deviceID: %d, vDevCount: %d, vDevCUs: %v, vDevMemSize: %v\n", dvInd, vDevCount, vDevCUs, vDevMemSize)

	// 获取调用前的配置文件列表
	beforeFiles, err := getConfigFiles("/etc/vdev")
	if err != nil {
		return vdevIDs, fmt.Errorf("Failed to get config files: %v", err)
	}
	fmt.Println("Before processing, the files in /etc/vdev are:")
	for _, file := range beforeFiles {
		fmt.Println(" -", file.Name())
	}

	// n := len(vDevCUs) (已经和 vDevCount 校验过)
	n := len(vDevCUs)

	// 分配 c 内存（按 C.int 大小）
	cVdevCusRaw := C.malloc(C.size_t(n) * C.size_t(C.sizeof_int))
	if cVdevCusRaw == nil {
		return vdevIDs, fmt.Errorf("Memory allocation failed for cVdevCus")
	}
	defer C.free(cVdevCusRaw)

	// 第二块内存
	cVdevMemRaw := C.malloc(C.size_t(n) * C.size_t(C.sizeof_int))
	if cVdevMemRaw == nil {
		// 释放第一块并返回，避免泄漏
		C.free(cVdevCusRaw)
		return vdevIDs, fmt.Errorf("Memory allocation failed for cVdevMemSize")
	}
	defer C.free(cVdevMemRaw)

	// 把 raw 指针转换成可索引的 C 数组切片视图（长度 n）
	cVdevCusArr := (*[1 << 28]C.int)(cVdevCusRaw)[:n:n] // 1<<28 是个大上限
	cVdevMemArr := (*[1 << 28]C.int)(cVdevMemRaw)[:n:n]

	// 逐项拷贝
	for i := 0; i < n; i++ {
		cVdevCusArr[i] = C.int(vDevCUs[i])
		cVdevMemArr[i] = C.int(vDevMemSize[i])
	}

	// 调用 C 接口
	ret := C.dmiCreateVDevices(
		C.int(dvInd),
		C.int(vDevCount),
		(*C.int)(cVdevCusRaw),
		(*C.int)(cVdevMemRaw),
	)
	glog.V(5).Infof("dmiCreateVDevices ret:%v ,err:%v", ret, dmiErrorString(ret))
	if err = dmiErrorString(ret); err != nil {
		return vdevIDs, fmt.Errorf("Error dmiCreateVDevices:%s", err)
	}

	// 获取调用后的配置文件列表
	afterFiles, err := getConfigFiles("/etc/vdev")
	if err != nil {
		return vdevIDs, fmt.Errorf("Failed to get config files: %v", err)
	}
	fmt.Println("After processing, the files in /etc/vdev are:")
	for _, file := range afterFiles {
		fmt.Println(" -", file.Name())
	}
	// 找出新增的配置文件
	newFiles := map[string]os.DirEntry{}
	for _, af := range afterFiles {
		found := false
		for _, bf := range beforeFiles {
			if af.Name() == bf.Name() {
				found = true
				break
			}
		}
		if !found {
			newFiles[af.Name()] = af
		}
	}
	// 处理新增的配置文件，提取vdev_id并返回
	for fileName := range newFiles {
		filePath := "/etc/vdev/" + fileName
		config, err := parseConfigFile(filePath)
		if err != nil {
			return vdevIDs, fmt.Errorf("Failed to parse config file %s: %v", filePath, err)
		}
		fmt.Printf("New config file: %s, content: %v\n", fileName, config)

		if vdevIDStr, ok := config["vdev_id"]; ok {
			var vdevID int
			fmt.Sscanf(vdevIDStr, "%d", &vdevID)
			vdevIDs = append(vdevIDs, vdevID)
		}
	}
	glog.V(5).Infof("vdevIDs:%v", vdevIDs)
	return
}

// 销毁指定物理设备上的所有虚拟设备
func dmiDestroyVDevices(dvInd int) (err error) {
	glog.V(5).Infof("dmiDestroyVDevices: %v", dvInd)
	ret := C.dmiDestroyVDevices(C.int(dvInd))
	glog.V(5).Infof("dmiDestroyVDevices ret:%v", ret)
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiDestroyVDevices:%s", err)
	}
	return
}

// 销毁指定虚拟设备
func dmiDestroySingleVDevice(vDvInd int) (err error) {
	glog.V(5).Infof("dmiDestroySingleVDevice:%v", vDvInd)
	ret := C.dmiDestroySingleVDevice(C.int(vDvInd))
	glog.V(5).Infof("dmiDestroySingleVDevice ret:%v", ret)
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiDestroySingleVDevice:%s", err)
	}
	return
}

// 更新指定设备资源大小，vDevCUs和vDevMemSize为-1是不更改
func dmiUpdateSingleVDevice(vDvInd int, vDevCUs int, vDevMemSize int) (err error) {
	ret := C.dmiUpdateSingleVDevice(C.int(vDvInd), C.int(vDevCUs), C.int(vDevMemSize))
	glog.V(5).Infof("dmiUpdateSingleVDevice ret:%v, retstr:%v", ret, dmiErrorString(ret))
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiUpdateSingleVDevice:%s", err)
	}
	return
}

// 启动虚拟设备
func dmiStartVDevice(vDvInd int) (err error) {
	ret := C.dmiStartVDevice(C.int(vDvInd))
	glog.V(5).Infof("StartVDevice ret:%v", ret)
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiStartVDevice:%s", err)
	}
	return
}

// 停止虚拟设备
func dmiStopVDevice(vDvInd int) (err error) {
	ret := C.dmiStopVDevice(C.int(vDvInd))
	glog.V(5).Infof("dmiStopVDevice ret:%v,retmessage:%v", ret, dmiErrorString(ret))
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiStopVDevice:%s", err)
	}
	return
}

// 返回物理设备使用百分比
func dmiGetDevBusyPercent(dvInd int) (percent int, err error) {
	var cpercent C.int
	ret := C.dmiGetDevBusyPercent(C.int(dvInd), &cpercent)
	glog.V(5).Infof("dmiGetDevBusyPercent ret:%v,retmessage:%v", ret, dmiErrorString(ret))
	if err = dmiErrorString(ret); err != nil {
		return percent, fmt.Errorf("Error dmiGetDevBusyPercent:%s", err)
	}
	percent = int(cpercent)
	glog.V(5).Infof("dmiGetDevBusyPercent: %v", percent)
	return
}

// 返回虚拟设备使用百分比
func dmiGetVDevBusyPercent(vDvInd int) (percent int, err error) {
	// ---------- 根据 vDvInd 获取 conf 文件 ----------
	confPath := filepath.Join("/etc/vdev", fmt.Sprintf("vdev%d.conf", vDvInd))
	file, err := os.Open(confPath)
	if err != nil {
		glog.Warningf("Failed to open conf file %s: %v", confPath, err)
		// 文件不存在或无法读取，假设设备不可用
		return 0, nil
	}
	defer file.Close()

	// ---------- 读取 cu_count 和 device_id ----------
	cuCount := 0
	dvInd := -1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// cu_count
		if strings.HasPrefix(line, "cu_count:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				valStr := strings.TrimSpace(parts[1])
				cuCount, err = strconv.Atoi(valStr)
				if err != nil {
					glog.Warningf("Invalid cu_count value in %s: %s", confPath, valStr)
					cuCount = 0
				}
			}
			continue
		}

		// device_id
		if strings.HasPrefix(line, "device_id:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				valStr := strings.TrimSpace(parts[1])
				dvInd, err = strconv.Atoi(valStr)
				if err != nil {
					glog.Warningf("Invalid device_id value in %s: %s", confPath, valStr)
					dvInd = -1
				}
			}
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		glog.Warningf("Error reading conf file %s: %v", confPath, err)
		return 0, nil
	}

	// ---------- cu_count == 0：使用物理设备使用率 ----------
	if cuCount == 0 {
		if dvInd < 0 {
			glog.Warningf("vdev%d cu_count is 0 but device_id is invalid", vDvInd)
			return 0, nil
		}

		glog.V(5).Infof(
			"vdev %d cu_count is 0, fallback to dmiGetDevBusyPercent(device_id=%d)",
			vDvInd, dvInd,
		)
		return dmiGetDevBusyPercent(dvInd)
	}

	// ---------- cu_count > 0：使用虚拟设备使用率 ----------
	var cpercent C.int
	ret := C.dmiGetVDevBusyPercent(C.int(vDvInd), &cpercent)
	glog.V(5).Infof(
		"dmiGetVDevBusyPercent: ret=%v, retmessage=%v",
		ret, dmiErrorString(ret),
	)

	if err = dmiErrorString(ret); err != nil {
		return 0, fmt.Errorf("Error dmiGetVDevBusyPercent: %s", err)
	}

	percent = int(cpercent)
	glog.V(5).Infof("dmiGetVDevBusyPercent(vdev=%d): %v", vDvInd, percent)
	return
}

// 设置虚拟机加密状态 status为true，则开启加密虚拟机，否则关闭
func dmiSetEncryptionVMStatus(status bool) (err error) {
	ret := C.dmiSetEncryptionVMStatus(C.bool(status))
	if err = dmiErrorString(ret); err != nil {
		return fmt.Errorf("Error dmiSetEncryptionVMStatus:%s", err)
	}
	return
}

// 获取加密虚拟机状态
func dmiGetEncryptionVMStatus() (status bool, err error) {
	var cstatus C.bool
	ret := C.dmiGetEncryptionVMStatus(&cstatus)
	if err = dmiErrorString(ret); err != nil {
		return false, fmt.Errorf("Error dmiGetEncryptionVMStatus:%s", err)
	}
	status = bool(cstatus)
	glog.V(5).Infof("DmiGetEncryptionVMStatus: %v", status)
	return
}

func rsmiDevGpuReset(dvInd int) (err error) {
	ret := C.rsmi_dev_gpu_reset(C.uint32_t(dvInd))
	glog.V(5).Infof("rsmi_dev_gpu_reset ret:%v, retStr: %v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_gpu_reset:%s", err)
	}
	return
}

// rsmiTopoGetLinkType 查询两个设备之间的互联链路类型及拓扑跳数（hops）。
// 参数说明：
//   - srcDev：源设备索引（GPU index）
//   - dstDev：目标设备索引（GPU index；若查询 GPU→CPU，
//             需传入 CPU_NODE_INDEX = 0xFFFFFFFF）
//
// 返回值说明：
//   - hops：设备之间的拓扑跳数，数值越小表示拓扑距离越近
//   - linkType：设备之间的链路类型
//   - err：调用失败时返回的错误信息

//func rsmiTopoGetLinkType(srcDvInd, dstDvInd int) (hops uint64, linkType RSMI_IO_LINK_TYPE, err error) {
//	var chops C.uint64_t
//	var ctype C.RSMI_IO_LINK_TYPE
//	ret := C.rsmi_topo_get_link_type(
//		C.uint32_t(srcDvInd),
//		C.uint32_t(dstDvInd),
//		&chops,
//		&ctype,
//	)
//	glog.V(5).Infof(
//		"rsmi_topo_get_link_type ret:%v, hops:%v, linkType:%v",
//		ret, chops, ctype,
//	)
//
//	if err = errorString(ret); err != nil {
//		return 0, 0, fmt.Errorf("rsmi_topo_get_link_type 调用失败: %s", err)
//	}
//
//	hops = uint64(chops)
//	linkType = RSMI_IO_LINK_TYPE(ctype)
//	return
//}

/****************************************** Utilization *********************************************/

// rsmiDevCuUsageGet 查询设备 CU 使用率。
//
// 接口含义：
//   - 对应 C 接口 rsmi_dev_cu_usage_get，返回指定设备的 CU 使用占比。
//   - 与 rsmiDevCuUtilGet 不同：本接口无需 duration，由驱动内部 UsageManager 统计；
//     cu_util_get 则需在指定时间窗口内统计 wave 占用周期。
//   - 调用前需已完成 rsmi_init 初始化。
//
// 参数说明：
//   - dvInd：设备索引，与 rsmi_num_monitor_devices 枚举顺序一致，从 0 开始
//
// 返回值说明：
//   - percent：CU 使用率，浮点数，范围通常为 0~1；设备空闲时通常为 0
//   - err：非 nil 表示调用失败，常见为 RSMI_STATUS_INVALID_ARGS（参数无效）或
//     RSMI_STATUS_INIT_ERROR（UsageManager 初始化失败）
func rsmiDevCuUsageGet(dvInd int) (percent float32, err error) {
	if err = ensureDevCuUsage(); err != nil {
		return percent, err
	}
	var cpercent C.float
	ret := C.rsmi_dev_cu_usage_get(C.uint32_t(dvInd), &cpercent)
	glog.V(5).Infof("rsmi_dev_cu_usage_get ret:%v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return percent, fmt.Errorf("Error rsmi_dev_cu_usage_get:%s", err)
	}
	percent = float32(cpercent)
	glog.V(5).Infof("rsmiDevCuUsageGet dvInd:%v percent:%v", dvInd, percent)
	return
}

// rsmiDevHcuUtilGet 查询设备 HCU（Hygon Compute Unit）在采样窗口内的活跃时间占比。
//
// 接口含义：
//   - 对应 C 接口 rsmi_dev_hcu_util_get，反映 GPU 计算单元（HCU）在指定时间段内
//     处于活跃状态的时间比例，可用于衡量设备整体计算繁忙程度。
//   - 调用前需已完成 rsmi_init 初始化。
//
// 参数说明：
//   - dvInd：设备索引，与 rsmi_num_monitor_devices 枚举顺序一致，从 0 开始
//   - duration：采样时间窗口；头文件未注明具体单位，实测可传 1000 作为常用值
//
// 返回值说明：
//   - percent：HCU 活跃时间占比，浮点数；设备空闲时通常为 0，负载越高越接近 1
//   - err：非 nil 表示调用失败，常见为 RSMI_STATUS_INVALID_ARGS（参数无效）或
//     RSMI_STATUS_NOT_SUPPORTED（当前硬件/驱动不支持）
func rsmiDevHcuUtilGet(dvInd int, duration int) (percent float32, err error) {
	if err = ensureDevHcuUtil(); err != nil {
		return percent, err
	}
	var cpercent C.float
	ret := C.rsmi_dev_hcu_util_get(C.uint32_t(dvInd), C.uint32_t(duration), &cpercent)
	glog.V(5).Infof("rsmi_dev_hcu_util_get ret:%v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return percent, fmt.Errorf("Error rsmi_dev_hcu_util_get:%s", err)
	}
	percent = float32(cpercent)
	glog.V(5).Infof("rsmiDevHcuUtilGet dvInd:%v duration:%v percent:%v", dvInd, duration, percent)
	return
}

// rsmiDevCuUtilGet 查询设备 CU 在采样窗口内的 wave 占用周期占比（全 CU 平均）。
//
// 接口含义：
//   - 对应 C 接口 rsmi_dev_cu_util_get，统计在 duration 时间窗口内，
//     每个 CU 至少分配了 1 个 wave 的时钟周期占比，再对所有 CU 取平均值。
//   - 反映计算单元被 kernel 实际占用的程度，比 HCU 粒度更细。
//
// 参数说明：
//   - dvInd：设备索引，从 0 开始
//   - duration：采样时间窗口，含义同 rsmiDevHcuUtilGet
//
// 返回值说明：
//   - percent：CU wave 占用周期占比（全 CU 平均），范围通常为 0~1
//   - err：调用失败时返回错误信息
func rsmiDevCuUtilGet(dvInd int, duration int) (percent float32, err error) {
	if err = ensureDevCuUtil(); err != nil {
		return percent, err
	}
	var cpercent C.float
	ret := C.rsmi_dev_cu_util_get(C.uint32_t(dvInd), C.uint32_t(duration), &cpercent)
	glog.V(5).Infof("rsmi_dev_cu_util_get ret:%v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return percent, fmt.Errorf("Error rsmi_dev_cu_util_get:%s", err)
	}
	percent = float32(cpercent)
	glog.V(5).Infof("rsmiDevCuUtilGet dvInd:%v duration:%v percent:%v", dvInd, duration, percent)
	return
}

// rsmiDevWaveUtilGet 查询设备 CU 在采样窗口内的 wave 驻留数量占比（全 CU 平均）。
//
// 接口含义：
//   - 对应 C 接口 rsmi_dev_wave_util_get，统计 duration 窗口内各 CU 上
//     驻留 wave 数量占最大可驻留 wave 数量的比例，再对所有 CU 取平均值。
//   - 与 rsmiDevCuUtilGet 互补：cu_util 看“有没有 wave”，wave_util 看“wave 填了多少”。
//
// 参数说明：
//   - dvInd：设备索引，从 0 开始
//   - duration：采样时间窗口，含义同 rsmiDevHcuUtilGet
//
// 返回值说明：
//   - percent：wave 驻留数量占比（全 CU 平均），范围通常为 0~1
//   - err：调用失败时返回错误信息
func rsmiDevWaveUtilGet(dvInd int, duration int) (percent float32, err error) {
	if err = ensureDevWaveUtil(); err != nil {
		return percent, err
	}
	var cpercent C.float
	ret := C.rsmi_dev_wave_util_get(C.uint32_t(dvInd), C.uint32_t(duration), &cpercent)
	glog.V(5).Infof("rsmi_dev_wave_util_get ret:%v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return percent, fmt.Errorf("Error rsmi_dev_wave_util_get:%s", err)
	}
	percent = float32(cpercent)
	glog.V(5).Infof("rsmiDevWaveUtilGet dvInd:%v duration:%v percent:%v", dvInd, duration, percent)
	return
}

// rsmiDevSeUtilGet 查询设备各 Shader Engine（SE）上活跃 CU 的利用率。
//
// 接口含义：
//   - 对应 C 接口 rsmi_dev_se_util_get，按 SE 维度返回活跃 CU 占比，
//     用于观察负载在不同 SE 之间的分布是否均衡。
//   - 本接口不需要 duration 参数，为即时/短周期统计。
//
// 参数说明：
//   - dvInd：设备索引，从 0 开始
//
// 返回值说明：
//   - seUsage：SE 利用率结构体，Percent[i] 为第 i 个 SE 的活跃 CU 占比（0~1）；
//     数组长度 MAX_SE_CNT（8），超出设备实际 SE 数的槽位一般为 0
//   - err：调用失败时返回错误信息
func rsmiDevSeUtilGet(dvInd int) (seUsage SEUsageInfo, err error) {
	if err = ensureDevSeUtil(); err != nil {
		return seUsage, err
	}
	var cinfo C.rsmi_se_usage_info_t
	ret := C.rsmi_dev_se_util_get(C.uint32_t(dvInd), &cinfo)
	glog.V(5).Infof("rsmi_dev_se_util_get ret:%v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return seUsage, fmt.Errorf("Error rsmi_dev_se_util_get:%s", err)
	}
	for i := 0; i < MAX_SE_CNT; i++ {
		seUsage.Percent[i] = float32(cinfo.percent[i])
	}
	glog.V(5).Infof("rsmiDevSeUtilGet dvInd:%v seUsage:%v", dvInd, seUsage)
	return
}

// rsmiDevPowerCapSet 设置设备功率上限（单位：微瓦）
func rsmiDevPowerCapSet(dvInd int, sensorInd int, cap uint64) error {
	ret := C.rsmi_dev_power_cap_set(C.uint32_t(dvInd), C.uint32_t(sensorInd), C.uint64_t(cap))
	if err := errorString(ret); err != nil {
		return fmt.Errorf("rsmi_dev_power_cap_set: %s", err)
	}
	return nil
}

// rsmiDevPowerCapDefaultGet 读取设备默认功率上限（单位：微瓦）
func rsmiDevPowerCapDefaultGet(dvInd int) (uint64, error) {
	var cDefault C.uint64_t
	ret := C.rsmi_dev_power_cap_default_get(C.uint32_t(dvInd), &cDefault)
	if err := errorString(ret); err != nil {
		return 0, fmt.Errorf("rsmi_dev_power_cap_default_get: %s", err)
	}
	return uint64(cDefault), nil
}

// rsmiDevPerfLevelSetV1 使用新版接口设置设备性能等级
func rsmiDevPerfLevelSetV1(dvInd int, devPerfLevel DevPerfLevel) error {
	ret := C.rsmi_dev_perf_level_set_v1(C.uint32_t(dvInd), C.rsmi_dev_perf_level_t(devPerfLevel))
	if err := errorString(ret); err != nil {
		return fmt.Errorf("rsmi_dev_perf_level_set_v1: %s", err)
	}
	return nil
}

// hcuHealthCheck 执行 hy-smi --healthcheck 并解析每张卡的健康状态。
// 输出格式示例：
//
//	HCU[0]          : Bus Id : 0000:09:00.0         : Healthy
func hcuHealthCheck() ([]HCUHealthStatus, error) {
	cmd := exec.Command("/opt/hyhal/bin/hy-smi", "--healthcheck")
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run hy-smi --healthcheck: %w, output: %s", err, outBytes)
	}

	re := regexp.MustCompile(`(?m)HCU\[(\d+)\]\s*:\s*Bus Id\s*:\s*(\S+)\s*:\s*(\w+)`)
	matches := re.FindAllStringSubmatch(string(outBytes), -1)

	result := make([]HCUHealthStatus, 0, len(matches))
	for _, m := range matches {
		idx, _ := strconv.ParseUint(m[1], 10, 64)
		result = append(result, HCUHealthStatus{
			HCU:    uint(idx),
			BusId:  m[2],
			Status: m[3],
		})
	}
	glog.V(5).Infof("hcuHealthCheck: %v", result)
	return result, nil
}
