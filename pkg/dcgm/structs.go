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

// PcieBandwidth 表示设备的 PCIe 带宽信息
// swagger:model PcieBandwidth
type PcieBandwidth struct {

	// TransferRate 表示传输速率的频率信息
	TransferRate Frequencies

	// lanes 表示 PCIe 通道的配置
	Lanes [33]uint32
}

// Frequencies 表示设备支持的频率信息
// swagger:model RSMIFrequencies
type Frequencies struct {
	HasDeepSleep bool

	// NumSupported 表示设备支持的频率数量
	NumSupported uint32

	// Current 表示当前使用的频率
	Current uint32

	// Frequency 表示设备支持的频率列表
	Frequency [33]uint64
}

type PowerProfilePresetMasks C.rsmi_power_profile_preset_masks_t

const (
	RSMI_PWR_PROF_PRST_CUSTOM_MASK       PowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_CUSTOM_MASK       //!< Custom Power Profile
	RSMI_PWR_PROF_PRST_VIDEO_MASK        PowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VIDEO_MASK        //!< Video Power Profile
	RSMI_PWR_PROF_PRST_POWER_SAVING_MASK PowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_POWER_SAVING_MASK //!< Power Saving Profile
	RSMI_PWR_PROF_PRST_COMPUTE_MASK      PowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_COMPUTE_MASK      //!< Compute Saving Profile
	RSMI_PWR_PROF_PRST_VR_MASK           PowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VR_MASK           //!< VR Power Profile

	//!< 3D Full Screen Power Profile
	RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK PowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK
	RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT   PowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT //!< Default Boot Up Profile
	RSMI_PWR_PROF_PRST_LAST             PowerProfilePresetMasks = RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT

	//!< Invalid power profile
	RSMI_PWR_PROF_PRST_INVALID PowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_INVALID
)

type RetiredPageRecord struct {
	PageAddress uint64           //!< Start address of page
	PageSize    uint64           //!< Page size
	Status      MemoryPageStatus //!< Page "reserved" status
}

type MemoryPageStatus C.rsmi_memory_page_status_t

const (
	RSMI_MEM_PAGE_STATUS_RESERVED     MemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_RESERVED
	RSMI_MEM_PAGE_STATUS_PENDING      MemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_PENDING
	RSMI_MEM_PAGE_STATUS_UNRESERVABLE MemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_UNRESERVABLE
)

type RSMIFreqVoltRegion struct {
	FreqRange RSMIRange
	VoltRange RSMIRange
}

type RSMITemperatureMetric C.rsmi_temperature_metric_t

const (
	RSMI_TEMP_CURRENT        RSMITemperatureMetric = C.RSMI_TEMP_CURRENT
	RSMI_TEMP_FIRST          RSMITemperatureMetric = C.RSMI_TEMP_FIRST
	RSMI_TEMP_MAX            RSMITemperatureMetric = C.RSMI_TEMP_MAX
	RSMI_TEMP_MIN            RSMITemperatureMetric = C.RSMI_TEMP_MIN
	RSMI_TEMP_MAX_HYST       RSMITemperatureMetric = C.RSMI_TEMP_MAX_HYST
	RSMI_TEMP_MIN_HYST       RSMITemperatureMetric = C.RSMI_TEMP_MIN_HYST
	RSMI_TEMP_CRITICAL       RSMITemperatureMetric = C.RSMI_TEMP_CRITICAL
	RSMI_TEMP_CRITICAL_HYST  RSMITemperatureMetric = C.RSMI_TEMP_CRITICAL_HYST
	RSMI_TEMP_EMERGENCY      RSMITemperatureMetric = C.RSMI_TEMP_EMERGENCY
	RSMI_TEMP_EMERGENCY_HYST RSMITemperatureMetric = C.RSMI_TEMP_EMERGENCY_HYST
	RSMI_TEMP_CRIT_MIN       RSMITemperatureMetric = C.RSMI_TEMP_CRIT_MIN
	RSMI_TEMP_CRIT_MIN_HYST  RSMITemperatureMetric = C.RSMI_TEMP_CRIT_MIN_HYST
	RSMI_TEMP_OFFSET         RSMITemperatureMetric = C.RSMI_TEMP_OFFSET
	RSMI_TEMP_LOWEST         RSMITemperatureMetric = C.RSMI_TEMP_LOWEST
	RSMI_TEMP_HIGHEST        RSMITemperatureMetric = C.RSMI_TEMP_HIGHEST
	RSMI_TEMP_LAST           RSMITemperatureMetric = C.RSMI_TEMP_LAST
)

// DeviceUtilizationInfo 表示指定设备的 GPU 与显存粗粒度利用率计数器快照。
// 底层调用 rsmi_utilization_count_get；若硬件不支持该功能，对应字段为 0，Supported 为 false。
type DeviceUtilizationInfo struct {
	// GfxActivity GPU 图形引擎活动计数器（RSMI_COARSE_GRAIN_GFX_ACTIVITY）。
	GfxActivity uint64
	// MemActivity 显存控制器活动计数器（RSMI_COARSE_GRAIN_MEM_ACTIVITY）。
	MemActivity uint64
	// Timestamp 计数器采样时间戳，分辨率 1ns。
	Timestamp int64
	// Supported 为 false 表示当前硬件/驱动不支持该功能（计数器值均为 0xFFFFFFFF）。
	Supported bool
}

// DeviceTemperatureInfo 表示指定设备四个常用传感器的当前温度，单位为摄氏度。
// 各传感器不可读时，对应字段保持零值，具体失败原因由 GetDeviceTemperatureInfo 的 error 返回。
type DeviceTemperatureInfo struct {
	// Edge GPU 边缘温度，对应 RSMI_TEMP_TYPE_EDGE(0)。
	Edge float64
	// Junction GPU 核心热点温度，对应 RSMI_TEMP_TYPE_JUNCTION(1)。
	Junction float64
	// Memory 显存温度，对应 RSMI_TEMP_TYPE_MEMORY(2)。
	Memory float64
	// Core HCU 核心温度，对应 RSMI_TEMP_TYPE_CORE(11)。
	Core float64
}

type RSMIVoltageType C.rsmi_voltage_type_t

const (
	RSMI_VOLT_TYPE_FIRST   RSMIVoltageType = C.RSMI_VOLT_TYPE_FIRST
	RSMI_VOLT_TYPE_VDDGFX  RSMIVoltageType = C.RSMI_VOLT_TYPE_VDDGFX
	RSMI_VOLT_TYPE_LAST    RSMIVoltageType = C.RSMI_VOLT_TYPE_LAST
	RSMI_VOLT_TYPE_INVALID RSMIVoltageType = C.RSMI_VOLT_TYPE_INVALID
)

type RSMIVoltageMetric C.rsmi_voltage_metric_t

const (
	RSMI_VOLT_CURRENT  RSMIVoltageMetric = C.RSMI_VOLT_CURRENT //!< Voltage current value.
	RSMI_VOLT_FIRST    RSMIVoltageMetric = C.RSMI_VOLT_FIRST
	RSMI_VOLT_MAX      RSMIVoltageMetric = C.RSMI_VOLT_MAX      //!< Voltage max value.
	RSMI_VOLT_MIN_CRIT RSMIVoltageMetric = C.RSMI_VOLT_MIN_CRIT //!< Voltage critical min value.
	RSMI_VOLT_MIN      RSMIVoltageMetric = C.RSMI_VOLT_MIN      //!< Voltage min value.
	RSMI_VOLT_MAX_CRIT RSMIVoltageMetric = C.RSMI_VOLT_MAX_CRIT //!< Voltage critical max value.
	RSMI_VOLT_AVERAGE  RSMIVoltageMetric = C.RSMI_VOLT_AVERAGE  //!< Average voltage.
	RSMI_VOLT_LOWEST   RSMIVoltageMetric = C.RSMI_VOLT_LOWEST   //!< Historical minimum voltage.
	RSMI_VOLT_HIGHEST  RSMIVoltageMetric = C.RSMI_VOLT_HIGHEST  //!< Historical maximum voltage.
	RSMI_VOLT_LAST                       = C.RSMI_VOLT_LAST
)

type UtilizationCounterType C.RSMI_UTILIZATION_COUNTER_TYPE

const (
	RSMI_UTILIZATION_COUNTER_FIRST UtilizationCounterType = C.RSMI_UTILIZATION_COUNTER_FIRST
	RSMI_COARSE_GRAIN_GFX_ACTIVITY UtilizationCounterType = C.RSMI_COARSE_GRAIN_GFX_ACTIVITY
	RSMI_COARSE_GRAIN_MEM_ACTIVITY UtilizationCounterType = C.RSMI_COARSE_GRAIN_MEM_ACTIVITY
	RSMI_UTILIZATION_COUNTER_LAST  UtilizationCounterType = C.RSMI_UTILIZATION_COUNTER_LAST
)

// @swagignore
type UtilizationCounter struct {

	// Type 表示利用率计数器的类型
	Type UtilizationCounterType

	// Value 表示计数器的值
	Value uint64
}

type RSMIClkType C.rsmi_clk_type_t

const (
	// sclk clock level
	RSMI_CLK_TYPE_SYS  RSMIClkType = C.RSMI_CLK_TYPE_SYS
	RSMI_CLK_TYPE_DF   RSMIClkType = C.RSMI_CLK_TYPE_DF
	RSMI_CLK_TYPE_DCEF RSMIClkType = C.RSMI_CLK_TYPE_DCEF
	// socclk clock level
	RSMI_CLK_TYPE_SOC  RSMIClkType = C.RSMI_CLK_TYPE_SOC
	RSMI_CLK_TYPE_MEM  RSMIClkType = C.RSMI_CLK_TYPE_MEM
	RSMI_CLK_TYPE_PCIE RSMIClkType = C.RSMI_CLK_TYPE_PCIE
	RSMI_CLK_INVALID   RSMIClkType = C.RSMI_CLK_INVALID
)

type OdVoltFreqData struct {
	CurrSclkRange  RSMIRange
	CurrMclkRange  RSMIRange
	SclkFreqLimits RSMIRange
	MclkFreqLimits RSMIRange
	Curve          RSMIOdVoltCurve
	NumRegions     uint32
}

type RSMIRange struct {
	LowerBound uint64
	UpperBound uint64
}

type RSMIOdVoltCurve struct {
	VcPoints [3]RSMIOdVddcPoint
}

type RSMIOdVddcPoint struct {
	Frequency uint64
	Voltage   uint64
}

// MetricsTableHeader  度量表头信息
// swagger:model MetricsTableHeader
type MetricsTableHeader struct {
	// StructureSize   结构体大小
	StructureSize uint16
	// FormatRevision   格式版本
	FormatRevision uint8
	// ContentRevision   内容版本
	ContentRevision uint8
}

// RSMIGPUMetrics  表示设备的度量信息
// swagger:model RSMIGPUMetrics
type RSMIGPUMetrics struct {
	// CommonHeader   公共表头
	CommonHeader MetricsTableHeader
	// TemperatureEdge   边缘温度
	TemperatureEdge uint16
	// TemperatureHotspot   热点温度
	TemperatureHotspot uint16
	// TemperatureMem   内存温度
	TemperatureMem uint16
	// TemperatureVRGfx   VR图形温度
	TemperatureVRGfx uint16
	// TemperatureVRSoc   VRSoC温度
	TemperatureVRSoc uint16
	// TemperatureVRMem   VR内存温度
	TemperatureVRMem uint16
	// AverageGfxActivity   平均图形活动
	AverageGfxActivity uint16
	// AverageUmcActivity   平均内存控制器活动
	AverageUmcActivity uint16
	// AverageMmActivity   平均多媒体活动
	AverageMmActivity uint16
	// AverageSocketPower   平均插座功率
	AverageSocketPower uint16
	// EnergyAccumulator   能量累加器
	EnergyAccumulator uint64
	// SystemClockCounter   系统时钟计数器
	SystemClockCounter uint64
	// AverageGfxclkFrequency   平均图形时钟频率
	AverageGfxclkFrequency uint16
	// AverageSocclkFrequency   平均SoC时钟频率
	AverageSocclkFrequency uint16
	// AverageUclkFrequency   平均内存时钟频率
	AverageUclkFrequency uint16
	// AverageVclk0Frequency   平均视频时钟0频率
	AverageVclk0Frequency uint16
	// AverageDclk0Frequency   平均显示时钟0频率
	AverageDclk0Frequency uint16
	// AverageVclk1Frequency   平均视频时钟1频率
	AverageVclk1Frequency uint16
	// AverageDclk1Frequency   平均显示时钟1频率
	AverageDclk1Frequency uint16
	// CurrentGfxclk   当前图形时钟
	CurrentGfxclk uint16
	// CurrentSocclk   当前SoC时钟
	CurrentSocclk uint16
	// CurrentUclk   当前内存时钟
	CurrentUclk uint16
	// CurrentVclk0   当前视频时钟0
	CurrentVclk0 uint16
	// CurrentDclk0   当前显示时钟0
	CurrentDclk0 uint16
	// CurrentVclk1   当前视频时钟1
	CurrentVclk1 uint16
	// CurrentDclk1   当前显示时钟1
	CurrentDclk1 uint16
	// ThrottleStatus   节流状态
	ThrottleStatus uint32
	// CurrentFanSpeed   当前风扇速度
	CurrentFanSpeed uint16
	// PcieLinkWidth   PCIe链路宽度
	PcieLinkWidth uint16
	// PcieLinkSpeed   PCIe链路速度（0.1 GT/s）
	PcieLinkSpeed uint16
	// Padding   填充
	Padding uint16
	// GfxActivityAcc   图形活动累加器
	GfxActivityAcc uint32
	// MemActivityAcc   内存活动累加器
	MemActivityAcc uint32

	// TempetureHBM   高带宽内存温度
	TempetureHBM [4]uint16
}

// ThrottleFlags 将 RSMIGPUMetrics.ThrottleStatus 位掩码拆分为各类降频状态
// bit 偏移基于海光 GPU gpu_metrics 驱动惯例；如厂商更新定义请同步修正
type ThrottleFlags struct {
	Thermal    bool // bit0: 芯片温度超限降频   → DCGM_FI_DEV_THERMAL_THROTTLE
	Power      bool // bit1: 功耗达到上限降频   → DCGM_FI_DEV_POWER_THROTTLE
	Slowdown   bool // bit2: 软硬件强制限频     → DCGM_FI_DEV_SLOWDOWN_THROTTLE
	BoardLimit bool // bit3: 板卡硬件限制降频   → DCGM_FI_DEV_BOARD_LIMIT_THROTTLE
}

type DevPerfLevel C.rsmi_dev_perf_level_t

const (
	RSMI_DEV_PERF_LEVEL_AUTO            DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_AUTO
	RSMI_DEV_PERF_LEVEL_FIRST           DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_FIRST
	RSMI_DEV_PERF_LEVEL_LOW             DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_LOW
	RSMI_DEV_PERF_LEVEL_HIGH            DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_HIGH
	RSMI_DEV_PERF_LEVEL_MANUAL          DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_MANUAL
	RSMI_DEV_PERF_LEVEL_STABLE_STD      DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_STD
	RSMI_DEV_PERF_LEVEL_STABLE_PEAK     DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_PEAK
	RSMI_DEV_PERF_LEVEL_STABLE_MIN_MCLK DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_MIN_MCLK
	RSMI_DEV_PERF_LEVEL_STABLE_MIN_SCLK DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_MIN_SCLK
	RSMI_DEV_PERF_LEVEL_DETERMINISM     DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_DETERMINISM
	RSMI_DEV_PERF_LEVEL_LAST            DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_LAST
	RSMI_DEV_PERF_LEVEL_UNKNOWN         DevPerfLevel = C.RSMI_DEV_PERF_LEVEL_UNKNOWN
)

// 系统支持的配置文件
type BitField C.rsmi_bit_field_t

// PowerProfileStatus  电源配置文件状态信息
type PowerProfileStatus struct {
	// AvailableProfiles  哪些配置文件被系统支持
	AvailableProfiles BitField
	// Current 当前激活的电源配置文件
	Current PowerProfilePresetMasks
	//  NumProfiles 可用的电源配置文件数量
	NumProfiles uint32
}

// DevVersion RSMI版本信息
type DevVersion struct {
	// Major 主版本号
	Major uint32
	// Minor 次版本号
	Minor uint32
	// Patch 补丁版本号
	Patch uint32
	// Build 构建字符串，包含构建的额外信息
	Build string
}

type SwComponent C.rsmi_sw_component_t

const (
	RSMISwCompFirst  SwComponent = C.RSMI_SW_COMP_FIRST
	RSMISwCompDriver SwComponent = C.RSMI_SW_COMP_DRIVER
	RSMISwCompLast   SwComponent = C.RSMI_SW_COMP_LAST
)

// 用于识别各种固
type RSMIFwBlock C.rsmi_fw_block_t

const (
	RSMIFwBlockFirst    RSMIFwBlock = C.RSMI_FW_BLOCK_FIRST
	RSMIFwBlockASD      RSMIFwBlock = C.RSMI_FW_BLOCK_ASD
	RSMIFwBlockCE       RSMIFwBlock = C.RSMI_FW_BLOCK_CE
	RSMIFwBlockDMCU     RSMIFwBlock = C.RSMI_FW_BLOCK_DMCU
	RSMIFwBlockMC       RSMIFwBlock = C.RSMI_FW_BLOCK_MC
	RSMIFwBlockME       RSMIFwBlock = C.RSMI_FW_BLOCK_ME
	RSMIFwBlockMEC      RSMIFwBlock = C.RSMI_FW_BLOCK_MEC
	RSMIFwBlockMEC2     RSMIFwBlock = C.RSMI_FW_BLOCK_MEC2
	RSMIFwBlockPFP      RSMIFwBlock = C.RSMI_FW_BLOCK_PFP
	RSMIFwBlockRLC      RSMIFwBlock = C.RSMI_FW_BLOCK_RLC
	RSMIFwBlockRLC_SRLC RSMIFwBlock = C.RSMI_FW_BLOCK_RLC_SRLC
	RSMIFwBlockRLC_SRLG RSMIFwBlock = C.RSMI_FW_BLOCK_RLC_SRLG
	RSMIFwBlockRLC_SRLS RSMIFwBlock = C.RSMI_FW_BLOCK_RLC_SRLS
	RSMIFwBlockSDMA     RSMIFwBlock = C.RSMI_FW_BLOCK_SDMA
	RSMIFwBlockSDMA2    RSMIFwBlock = C.RSMI_FW_BLOCK_SDMA2
	RSMIFwBlockSMC      RSMIFwBlock = C.RSMI_FW_BLOCK_SMC
	RSMIFwBlockSOS      RSMIFwBlock = C.RSMI_FW_BLOCK_SOS
	RSMIFwBlockTA_RAS   RSMIFwBlock = C.RSMI_FW_BLOCK_TA_RAS
	RSMIFwBlockTA_HSL   RSMIFwBlock = C.RSMI_FW_BLOCK_TA_XGMI
	RSMIFwBlockUVD      RSMIFwBlock = C.RSMI_FW_BLOCK_UVD
	RSMIFwBlockVCE      RSMIFwBlock = C.RSMI_FW_BLOCK_VCE
	RSMIFwBlockVCN      RSMIFwBlock = C.RSMI_FW_BLOCK_VCN
	RSMIFwBlockLast     RSMIFwBlock = C.RSMI_FW_BLOCK_LAST
)

// 保存错误计
type RSMIErrorCount struct {
	CorrectableErr   uint64
	UncorrectableErr uint64
}

// 用于标识不同的GPU
type RSMIGpuBlock C.rsmi_gpu_block_t

const (
	RSMIGpuBlockInvalid  RSMIGpuBlock = C.RSMI_GPU_BLOCK_INVALID
	RSMIGpuBlockFirst    RSMIGpuBlock = C.RSMI_GPU_BLOCK_FIRST
	RSMIGpuBlockUMC      RSMIGpuBlock = C.RSMI_GPU_BLOCK_UMC
	RSMIGpuBlockSDMA     RSMIGpuBlock = C.RSMI_GPU_BLOCK_SDMA
	RSMIGpuBlockGFX      RSMIGpuBlock = C.RSMI_GPU_BLOCK_GFX
	RSMIGpuBlockMMHUB    RSMIGpuBlock = C.RSMI_GPU_BLOCK_MMHUB
	RSMIGpuBlockATHUB    RSMIGpuBlock = C.RSMI_GPU_BLOCK_ATHUB
	RSMIGpuBlockPCIEBIF  RSMIGpuBlock = C.RSMI_GPU_BLOCK_PCIE_BIF
	RSMIGpuBlockHDP      RSMIGpuBlock = C.RSMI_GPU_BLOCK_HDP
	RSMIGpuBlockHSLWAFL  RSMIGpuBlock = C.RSMI_GPU_BLOCK_XGMI_WAFL
	RSMIGpuBlockDF       RSMIGpuBlock = C.RSMI_GPU_BLOCK_DF
	RSMIGpuBlockSMN      RSMIGpuBlock = C.RSMI_GPU_BLOCK_SMN
	RSMIGpuBlockSEM      RSMIGpuBlock = C.RSMI_GPU_BLOCK_SEM
	RSMIGpuBlockMP0      RSMIGpuBlock = C.RSMI_GPU_BLOCK_MP0
	RSMIGpuBlockMP1      RSMIGpuBlock = C.RSMI_GPU_BLOCK_MP1
	RSMIGpuBlockFuse     RSMIGpuBlock = C.RSMI_GPU_BLOCK_FUSE
	RSMIGpuBlockMCA      RSMIGpuBlock = C.RSMI_GPU_BLOCK_MCA
	RSMIGpuBlockLast     RSMIGpuBlock = C.RSMI_GPU_BLOCK_LAST
	RSMIGpuBlockReserved RSMIGpuBlock = C.RSMI_GPU_BLOCK_RESERVED
)

// 当前ECC状态
type RSMIRasErrState C.rsmi_ras_err_state_t

const (
	RSMIRasErrStateNone     RSMIRasErrState = C.RSMI_RAS_ERR_STATE_NONE
	RSMIRasErrStateDisabled RSMIRasErrState = C.RSMI_RAS_ERR_STATE_DISABLED
	RSMIRasErrStateParity   RSMIRasErrState = C.RSMI_RAS_ERR_STATE_PARITY
	RSMIRasErrStateSingC    RSMIRasErrState = C.RSMI_RAS_ERR_STATE_SING_C
	RSMIRasErrStateMultUC   RSMIRasErrState = C.RSMI_RAS_ERR_STATE_MULT_UC
	RSMIRasErrStatePoison   RSMIRasErrState = C.RSMI_RAS_ERR_STATE_POISON
	RSMIRasErrStateEnabled  RSMIRasErrState = C.RSMI_RAS_ERR_STATE_ENABLED
	RSMIRasErrStateLast     RSMIRasErrState = C.RSMI_RAS_ERR_STATE_LAST
	RSMIRasErrStateInvalid  RSMIRasErrState = C.RSMI_RAS_ERR_STATE_INVALID
)

// 事件组枚举值
type RSMIEventGroup C.rsmi_event_group_t

const (
	RSMI_EVNT_GRP_HSL          RSMIEventGroup = C.RSMI_EVNT_GRP_XGMI
	RSMI_EVNT_GRP_HSL_DATA_OUT RSMIEventGroup = C.RSMI_EVNT_GRP_XGMI_DATA_OUT
	RSMI_EVNT_GRP_INVALID      RSMIEventGroup = C.RSMI_EVNT_GRP_INVALID
)

type RSMIEventType C.rsmi_event_type_t

const (
	RSMIEventFirst RSMIEventType = C.RSMI_EVNT_FIRST

	RSMIEventHSLFirst       RSMIEventType = C.RSMI_EVNT_XGMI_FIRST
	RSMIEventHSL0NopTx      RSMIEventType = C.RSMI_EVNT_XGMI_0_NOP_TX
	RSMIEventHSL0RequestTx  RSMIEventType = C.RSMI_EVNT_XGMI_0_REQUEST_TX
	RSMIEventHSL0ResponseTx RSMIEventType = C.RSMI_EVNT_XGMI_0_RESPONSE_TX
	RSMIEventHSL0BeatsTx    RSMIEventType = C.RSMI_EVNT_XGMI_0_BEATS_TX
	RSMIEventHSL1NopTx      RSMIEventType = C.RSMI_EVNT_XGMI_1_NOP_TX
	RSMIEventHSL1RequestTx  RSMIEventType = C.RSMI_EVNT_XGMI_1_REQUEST_TX
	RSMIEventHSL1ResponseTx RSMIEventType = C.RSMI_EVNT_XGMI_1_RESPONSE_TX
	RSMIEventHSL1BeatsTx    RSMIEventType = C.RSMI_EVNT_XGMI_1_BEATS_TX

	RSMIEventHSLLast RSMIEventType = C.RSMI_EVNT_XGMI_LAST

	RSMIEventHSLDataOutFirst RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_FIRST

	RSMIEventHSLDataOut0    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_0
	RSMIEventHSLDataOut1    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_1
	RSMIEventHSLDataOut2    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_2
	RSMIEventHSLDataOut3    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_3
	RSMIEventHSLDataOut4    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_4
	RSMIEventHSLDataOut5    RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_5
	RSMIEventHSLDataOutLast RSMIEventType = C.RSMI_EVNT_XGMI_DATA_OUT_LAST

	RSMIEventLast RSMIEventType = C.RSMI_EVNT_LAST
)

type EventHandle C.rsmi_event_handle_t

type RSMICounterCommand C.rsmi_counter_command_t

const (
	RSMI_CNTR_CMD_START RSMICounterCommand = C.RSMI_CNTR_CMD_START
	RSMI_CNTR_CMD_STOP  RSMICounterCommand = C.RSMI_CNTR_CMD_STOP
)

// 计数器值
type RSMICounterValue struct {
	Value       uint64
	TimeEnabled uint64
	TimeRunning uint64
}

// ProcessInfos 进程的计算资源使用信息
type ProcessInfos struct {
	// ProcessID 进程 ID
	ProcessID uint32

	// Pasid 进程地址空间 ID
	Pasid uint32

	// VramUsage 显存使用量
	VramUsage uint64

	// SdmaUsage SDMA 使用量
	SdmaUsage uint64

	// CuOccupancy CU 占用率
	CuOccupancy uint32
}

type RsmiProcessInfoV2 struct {
	ProcessID     uint32    // Process ID
	VRAMUsageSize uint64    // VRAM usage size in MiB
	VRAMUsageRate float32   // VRAM usage rate as a percentage
	GPUCount      int       // Number of GPUs used
	GPUIndex      []int     // GPU index as a slice
	GPUUsageRate  []float32 // GPU usage rate as a percentage
}

// 进程的信息
type Process struct {
	// ProcessID 进程id
	ProcessID uint32
	//ProcessName 进程名称
	ProcessName string
	// Pasid 进程地址空间
	Pasid uint32
	// VramUsage 显存使用量
	VramUsage uint64
	// SdmaUsage SDMA使用量
	SdmaUsage uint64
	// CuOccupancy CU占用率
	CuOccupancy uint32
	// MinorNumbers 设备的索引号
	MinorNumbers []int
}

// RSMIHSLStatus HSL状态
type RSMIHSLStatus C.rsmi_xgmi_status_t

const (
	// RSMIHSLStatus 0
	RSMIHSLStatusNoErrors RSMIHSLStatus = C.RSMI_XGMI_STATUS_NO_ERRORS
	// RSMIHSLStatusError 1
	RSMIHSLStatusError RSMIHSLStatus = C.RSMI_XGMI_STATUS_ERROR
	// RSMIHSLStatusMultipleErrors 2
	RSMIHSLStatusMultipleErrors RSMIHSLStatus = C.RSMI_XGMI_STATUS_MULTIPLE_ERRORS
)

// 表示两个设备（GPU↔GPU / GPU↔CPU）之间的互联链路类型。
type RSMIIOLinkType C.RSMI_IO_LINK_TYPE

const (
	// RSMI_IOLINK_TYPE_UNDEFINED 表示链路类型未知或设备之间不存在直连关系，查询失败、设备未连通
	RSMIIOLinkTypeUndefined RSMIIOLinkType = C.RSMI_IOLINK_TYPE_UNDEFINED
	// RSMI_IOLINK_TYPE_PCIEXPRESS 表示设备之间通过 PCI Express 进行连接，GPU↔CPU，或通过 PCIe Switch/Bridge 的 GPU↔GPU 连接
	RSMIIOLinkTypePCIExpress RSMIIOLinkType = C.RSMI_IOLINK_TYPE_PCIEXPRESS
	// RSMI_IOLINK_TYPE_HSL 表示设备之间通过 HSL 进行直连，HSL 是 GPU↔GPU 的高速、低延迟互联方式
	RSMIIOLinkTypeHSL RSMIIOLinkType = C.RSMI_IOLINK_TYPE_XGMI
	// RSMI_IOLINK_TYPE_NUMIOLINKTYPES 表示当前支持的链路类型数量。
	RSMIIOLinkTypeNumIOLinkTypes RSMIIOLinkType = C.RSMI_IOLINK_TYPE_NUMIOLINKTYPES
	// RSMI_IOLINK_TYPE_SIZE 用于强制枚举大小为 32 位，保证 ABI 兼容性。
	RSMIIOLinkTypeSize RSMIIOLinkType = C.RSMI_IOLINK_TYPE_SIZE
)

type RSMIFuncIDIterHandle C.rsmi_func_id_iter_handle_t

type RSMIMemoryType C.rsmi_memory_type_t

const (
	RSMI_MEM_TYPE_FIRST    RSMIMemoryType = C.RSMI_MEM_TYPE_FIRST
	RSMI_MEM_TYPE_VRAM     RSMIMemoryType = C.RSMI_MEM_TYPE_VRAM
	RSMI_MEM_TYPE_VIS_VRAM RSMIMemoryType = C.RSMI_MEM_TYPE_VIS_VRAM
	RSMI_MEM_TYPE_GTT      RSMIMemoryType = C.RSMI_MEM_TYPE_GTT
	RSMI_MEM_TYPE_LAST     RSMIMemoryType = C.RSMI_MEM_TYPE_LAST
)

type RSMIFuncIDValue struct {
	ID         uint64
	Name       string
	MemoryType RSMIMemoryType
	TempMetric RSMITemperatureMetric
	EventType  RSMIEventType
	EventGroup RSMIEventGroup
	ClkType    RSMIClkType
	FwBlock    RSMIFwBlock
	GpuBlock   RSMIGpuBlock
}

type RSMIEvtNotificationType C.rsmi_evt_notification_type_t

const (
	RSMI_EVT_NOTIF_VMFAULT          RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_VMFAULT
	RSMI_EVT_NOTIF_FIRST            RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_FIRST
	RSMI_EVT_NOTIF_THERMAL_THROTTLE RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_THERMAL_THROTTLE
	RSMI_EVT_NOTIF_GPU_PRE_RESET    RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_GPU_PRE_RESET
	RSMI_EVT_NOTIF_GPU_POST_RESET   RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_GPU_POST_RESET
	RSMI_EVT_NOTIF_LAST             RSMIEvtNotificationType = C.RSMI_EVT_NOTIF_LAST
)

type RSMIEEvtNotificationData struct {
	DvInd   uint32
	Event   RSMIEvtNotificationType
	Message [64]byte
}

type RSMIStatus C.rsmi_status_t

const (
	RSMI_STATUS_SUCCESS             RSMIStatus = C.RSMI_STATUS_SUCCESS             //!< Operation was successful
	RSMI_STATUS_INVALID_ARGS        RSMIStatus = C.RSMI_STATUS_INVALID_ARGS        //!< Passed in arguments are not valid
	RSMI_STATUS_NOT_SUPPORTED       RSMIStatus = C.RSMI_STATUS_NOT_SUPPORTED       //!< The requested information or
	RSMI_STATUS_FILE_ERROR          RSMIStatus = C.RSMI_STATUS_FILE_ERROR          //!< Problem accessing a file. This
	RSMI_STATUS_PERMISSION          RSMIStatus = C.RSMI_STATUS_PERMISSION          //!< Permission denied/EACCESS file
	RSMI_STATUS_OUT_OF_RESOURCES    RSMIStatus = C.RSMI_STATUS_OUT_OF_RESOURCES    //!< Unable to acquire memory or other
	RSMI_STATUS_INTERNAL_EXCEPTION  RSMIStatus = C.RSMI_STATUS_INTERNAL_EXCEPTION  //!< An internal exception was caught
	RSMI_STATUS_INPUT_OUT_OF_BOUNDS RSMIStatus = C.RSMI_STATUS_INPUT_OUT_OF_BOUNDS //!< The provided input is out of
	RSMI_STATUS_INIT_ERROR          RSMIStatus = C.RSMI_STATUS_INIT_ERROR          //!< An error occurred when rsmi
	RSMI_INITIALIZATION_ERROR       RSMIStatus = C.RSMI_INITIALIZATION_ERROR
	RSMI_STATUS_NOT_YET_IMPLEMENTED RSMIStatus = C.RSMI_STATUS_NOT_YET_IMPLEMENTED //!< The requested function has not
	RSMI_STATUS_NOT_FOUND           RSMIStatus = C.RSMI_STATUS_NOT_FOUND           //!< An item was searched for but not
	RSMI_STATUS_INSUFFICIENT_SIZE   RSMIStatus = C.RSMI_STATUS_INSUFFICIENT_SIZE   //!< Not enough resources were
	RSMI_STATUS_INTERRUPT           RSMIStatus = C.RSMI_STATUS_INTERRUPT           //!< An interrupt occurred during
	RSMI_STATUS_UNEXPECTED_SIZE     RSMIStatus = C.RSMI_STATUS_UNEXPECTED_SIZE     //!< An unexpected amount of data
	RSMI_STATUS_NO_DATA             RSMIStatus = C.RSMI_STATUS_NO_DATA             //!< No data was found for a given
	RSMI_STATUS_UNEXPECTED_DATA     RSMIStatus = C.RSMI_STATUS_UNEXPECTED_DATA     //!< The data read or provided to
	RSMI_STATUS_BUSY                RSMIStatus = C.RSMI_STATUS_BUSY
	RSMI_STATUS_REFCOUNT_OVERFLOW   RSMIStatus = C.RSMI_STATUS_REFCOUNT_OVERFLOW   //!< An internal reference counter
	RSMI_STATUS_SETTING_UNAVAILABLE RSMIStatus = C.RSMI_STATUS_SETTING_UNAVAILABLE //!< Requested setting is unavailable
	RSMI_STATUS_AMDGPU_RESTART_ERR  RSMIStatus = C.RSMI_STATUS_AMDGPU_RESTART_ERR  //!< Could not successfully restart
	RSMI_STATUS_UNKNOWN_ERROR       RSMIStatus = C.RSMI_STATUS_UNKNOWN_ERROR
)

// MonitorInfo 设备监控信息
// swagger:model MonitorInfo
type MonitorInfo struct {
	//  MinorNumber 设备索引号
	MinorNumber int
	//  PciBusNumber PCI ID
	PciBusNumber string
	//  DeviceId 设备序列号
	DeviceId string
	//  SubSystemName 型号名称
	SubSystemName string
	// Temperature 设备温度
	Temperature float64
	//  PowerUsage 设备平均功耗
	PowerUsage float64
	//  PowerCap 设备功率上限
	PowerCap float64
	//  MemoryCap 设备内存总量
	MemoryCap float64
	//  MemoryUsed 设备内存使用量
	MemoryUsed float64
	//  UtilizationRate 设备忙碌时间百分比
	UtilizationRate float64
	//  PcieBwMb pcie流量信息
	PcieBwMb float64
	// Clk 系统时钟速度
	Clk float64
	// SclkFrequency 系统时钟频率列表
	SclkFrequency []string
	// Socclk socclk时钟
	Socclk float64
	// SocclkFrequency Soc时钟频率列表
	SocclkFrequency []string
	// PerfLevel 性能水平
	PerfLevel string
}

// DeviceInfo 设备信息结构体
type DeviceInfo struct {
	// DvInd 设备索引
	DvInd int
	// DeviceId 设备ID
	DeviceId string
	// DevType 设备类型
	DevType string
	// DevTypeName 设备类型名称
	DevTypeName string
	//SubsystemTypeId 设备子系统类型名称
	SubsystemTypeId string
	//SubsystemTypeName 设备子系统名称
	SubsystemTypeName string
	// PciBusNumber 设备的总线号
	PciBusNumber string
	// PowerTotal 设备电源总功率
	PowerTotal float64
	// MemoryTotal 设备的内存总量
	MemoryTotal float64
	// ComputeUnit 设备的计算单元数量
	ComputeUnit float64
	// VDeviceCount 虚拟卡数量
	VDeviceCount int
}

// DeviceInfo 设备信息结构体
type DeviceStatusInfo struct {
	// DvInd 设备索引
	DvInd int
	// Temperature 设备当前的温度
	Temperature float64
	// PowerUsed 设备当前的功耗
	PowerUsed float64
	// MemoryUsed 设备已使用的内存
	MemoryUsed float64
	// UtilizationRate 设备时间忙碌百分比
	UtilizationRate float64
	// PcieBwMb 设备的PCIe带宽 (MB/s)
	PcieBwMb float64
	// PcieSent 设备的PCIe发送数据量
	PcieSent float64
	// PcieReceived 设备的PCIe接收数据量
	PcieReceived float64
	// Clk 设备的当前时钟频率
	Clk float64
	// Percent 物理设备使用百分比
	Percent int
	// ComputeUnitRemainingCount 设备剩余可用的计算单元数量
	ComputeUnitRemainingCount uint64
	// MemoryRemaining 设备剩余可用的内存量
	MemoryRemaining uint64
	// BlocksInfo 设备的block信息
	BlocksInfos []BlocksInfo
}

var type2name = map[string]string{
	// Z 系列
	"55b7": "Z100L",
	"55c7": "Z200SM_71",
	"55d7": "Z200SM_71_S",

	// K 系列
	"61a7": "K500SM_B",
	"61b7": "K500SM",
	"61c7": "K500SM_B",
	"61d7": "K500SM",
	"61e7": "K500SM",
	"61f7": "HG Design HCU K500",

	"6210": "K100_AI",
	"6211": "K100_LC_E_AI",
	"6212": "K100_LC_AI",
	"6213": "K100_AI_i",
	"6214": "HG Design HCU K500",

	"62a0": "K500SM_AI",
	"62b0": "K500SM_AI",
	"62b7": "K100",
	"62c7": "K100-LC",

	// BW 系列
	"6310": "BW200",
	"6311": "BW200_LC",
	"631a": "BW1000_H",

	"6320": "BW150",
	"6321": "BW100_LC",
	"632a": "BW151",

	"6330": "BW500SM",

	"6360": "BW10",
	"636a": "BW11",

	"6370": "BW100",
	"637a": "BW101",

	"6423": "BW1500B",
	"6430": "BW1100",
	"6431": "BW1000_LC",
	"6436": "BW3000B",
	"6437": "BW3000B_LC",
}

var computeUnitType = map[string]float64{
	"K100_AI":      120,
	"K100_LC_AI":   120,
	"K100_LC_E_AI": 128,
	"K100":         120,
	"Z100":         60,
	"Z100L":        60,
	"BW200":        80,
	"BW200_LC":     80,
	"BW100":        80,
	"BW100_LC":     80,
}
var memoryType = map[string]float64{
	"K100_AI":      64,
	"K100_LC_AI":   64,
	"K100_LC_E_AI": 64,
	"K100":         64,
	"Z100":         32,
	"Z100L":        32,
	"BW200":        80,
}

var memoryTypeL = []string{"VRAM", "VIS_VRAM", "GTT"}

var memoryTypeMap = map[string]RSMIMemoryType{
	"VRAM":     RSMI_MEM_TYPE_VRAM,
	"VIS_VRAM": RSMI_MEM_TYPE_VIS_VRAM,
	"GTT":      RSMI_MEM_TYPE_GTT,
}

var memTypeMapReverse = map[RSMIMemoryType]string{
	RSMI_MEM_TYPE_VRAM:     "VRAM",
	RSMI_MEM_TYPE_VIS_VRAM: "VIS_VRAM",
	RSMI_MEM_TYPE_GTT:      "GTT",
}

const DMI_NAME_SIZE = 256

// @swagignore
type DMIDeviceInfo struct {
	Name             string
	ComputeUnitCount int
	// @swagignore
	ComputeUnitRemainingCount uintptr
	// @swagignore
	MemoryRemaining uintptr
	// @swagignore
	GlobalMemSize uintptr
	// @swagignore
	UsageMemSize    uintptr
	DeviceID        int
	Percent         int
	MaxVDeviceCount int
}

// VDeviceInfo 虚拟设备信息
type VDeviceInfo struct {
	// Name 虚拟设备的名称
	Name string

	//SubsystemTypeName 设备子系统名称
	SubsystemTypeName string

	// ComputeUnitCount 虚拟设备的计算单元数量
	VComputeUnitCount int

	// VMemoryTotal 虚拟设备的全局内存大小
	VMemoryTotal uintptr

	// VMemoryUsed 虚拟设备的已使用内存大小
	VMemoryUsed uintptr

	// ContainerID 虚拟设备的容器ID
	ContainerID uint64

	// DvInd 物理设备的设备ID
	DvInd int

	// VPercent 虚拟设备的使用百分比
	VPercent int

	// VdvInd 虚拟设备的索引号
	VdvInd int

	// PciBusNumber 虚拟设备的总线编号
	PciBusNumber string
}

type DMIStatus C.dmiStatus

const (
	DMI_STATUS_SUCCESS                DMIStatus = C.DMI_STATUS_SUCCESS
	DMI_STATUS_ERROR                  DMIStatus = C.DMI_STATUS_ERROR
	DMI_STATUS_NO_MEMORY              DMIStatus = C.DMI_STATUS_NO_MEMORY
	DMI_STATUS_OPEN_MKFD_FAILED       DMIStatus = C.DMI_STATUS_OPEN_MKFD_FAILED
	DMI_STATUS_MKFD_ALREADY_OPENED    DMIStatus = C.DMI_STATUS_MKFD_ALREADY_OPENED
	DMI_STATUS_SYS_NODE_NOT_EXIST     DMIStatus = C.DMI_STATUS_SYS_NODE_NOT_EXIST
	DMI_STATUS_NOT_SUPPORTED          DMIStatus = C.DMI_STATUS_NOT_SUPPORTED
	DMI_STATUS_MKFD_NOT_OPENED        DMIStatus = C.DMI_STATUS_MKFD_NOT_OPENED
	DMI_STATUS_CREATE_VDEV_FAILED     DMIStatus = C.DMI_STATUS_CREATE_VDEV_FAILED
	DMI_STATUS_DESTROY_VDEV_FAILED    DMIStatus = C.DMI_STATUS_DESTROY_VDEV_FAILED
	DMI_STATUS_INVALID_ARGUMENTS      DMIStatus = C.DMI_STATUS_INVALID_ARGUMENTS
	DMI_STATUS_OUT_OF_RESOURCES       DMIStatus = C.DMI_STATUS_OUT_OF_RESOURCES
	DMI_STATUS_QUERY_VDEV_INFO_FAILED DMIStatus = C.DMI_STATUS_QUERY_VDEV_INFO_FAILED
	DMI_STATUS_ERROR_NOT_INITIALIZED  DMIStatus = C.DMI_STATUS_ERROR_NOT_INITIALIZED
	DMI_STATUS_DEVICE_NOT_SUPPORT     DMIStatus = C.DMI_STATUS_DEVICE_NOT_SUPPORT
	DMI_STATUS_VDEV_NOT_EXIST         DMIStatus = C.DMI_STATUS_VDEV_NOT_EXIST
	DMI_STATUS_INIT_DEVICE_FAILED     DMIStatus = C.DMI_STATUS_INIT_DEVICE_FAILED
	DMI_STATUS_DEVICE_BUSY            DMIStatus = C.DMI_STATUS_DEVICE_BUSY
	DMI_STATUS_FILE_ERROR             DMIStatus = C.DMI_STATUS_FILE_ERROR
	DMI_STATUS_PERMISSION             DMIStatus = C.DMI_STATUS_PERMISSION
	DMI_STATUS_INTERNAL_EXCEPTION     DMIStatus = C.DMI_STATUS_INTERNAL_EXCEPTION
	DMI_STATUS_INPUT_OUT_OF_BOUNDS    DMIStatus = C.DMI_STATUS_INPUT_OUT_OF_BOUNDS
	DMI_STATUS_SMI_INIT_ERROR         DMIStatus = C.DMI_STATUS_SMI_INIT_ERROR
	DMI_STATUS_NOT_FOUND              DMIStatus = C.DMI_STATUS_NOT_FOUND
	DMI_STATUS_INSUFFICIENT_SIZE      DMIStatus = C.DMI_STATUS_INSUFFICIENT_SIZE
	DMI_STATUS_INTERRUPT              DMIStatus = C.DMI_STATUS_INTERRUPT
	DMI_STATUS_UNEXPECTED_SIZE        DMIStatus = C.DMI_STATUS_UNEXPECTED_SIZE
	DMI_STATUS_NO_DATA                DMIStatus = C.DMI_STATUS_NO_DATA
	DMI_STATUS_UNEXPECTED_DATA        DMIStatus = C.DMI_STATUS_UNEXPECTED_DATA
	DMI_STATUS_SMI_BUSY               DMIStatus = C.DMI_STATUS_SMI_BUSY
	DMI_STATUS_REFCOUNT_OVERFLOW      DMIStatus = C.DMI_STATUS_REFCOUNT_OVERFLOW
	DMI_STATUS_NOT_YET_IMPLEMENTED    DMIStatus = C.DMI_STATUS_NOT_YET_IMPLEMENTED
	DMI_STATUS_UNKNOWN_ERROR          DMIStatus = C.DMI_STATUS_UNKNOWN_ERROR
)

// Device 物理设备的详细信息
type Device struct {
	// DvInd 设备的索引号
	DvInd int

	// PciBusNumber 设备的总线编号
	PciBusNumber string

	// DeviceId 设备的唯一标识符
	DeviceId string

	// DevTypeId 设备类型id
	DevTypeId string

	// DevTypeName 设备类型名称
	DevTypeName string

	//SubsystemTypeId 设备子系统类型名称
	SubsystemTypeId string

	//SubsystemTypeName 设备子系统名称
	SubsystemTypeName string

	// Temperature 设备当前的温度
	Temperature float64

	// PowerUsed 设备当前的功耗
	PowerUsed float64

	// PowerTotal 设备的功耗上限
	PowerTotal float64

	// MemoryTotal 设备的内存容量
	MemoryTotal float64

	// MemoryUsed 设备已使用的内存
	MemoryUsed float64

	// UtilizationRate 设备的利用率
	UtilizationRate float64

	// PcieBwMb 设备的PCIe带宽 (MB/s)
	PcieBwMb float64

	// PcieSent 设备的PCIe发送数据量
	PcieSent float64

	// PcieReceived 设备的PCIe接收数据量
	PcieReceived float64

	// Clk 设备的当前时钟频率
	Clk float64

	// ComputeUnitCount 设备的计算单元总数
	ComputeUnitCount float64

	// ComputeUnitRemainingCount 设备剩余可用的计算单元数量
	ComputeUnitRemainingCount uint64

	// MemoryRemaining 设备剩余可用的内存量
	MemoryRemaining uint64

	// Percent 物理设备使用百分比
	Percent int

	// MaxVDeviceCount 物理设备上支持的最大虚拟设备数量
	MaxVDeviceCount int

	// VDeviceCount 虚拟设备数量
	VDeviceCount int

	// BlocksInfo 设备的block信息
	BlocksInfos []BlocksInfo
	// DFBandwidthInfo DF内存带宽信息
	DFBandwidthInfo DFBandwidthInfo
}

// PhysicalDeviceInfo 物理设备信息
type PhysicalDeviceInfo struct {
	// Device 物理设备的详细信息
	Device Device
	// VirtualDevices 该物理设备上关联的虚拟设备信息列表
	VirtualDevices []VDeviceInfo
}

// 定义事件通知类型名称
var notificationTypeNames = []string{"VM_FAULT", "THERMAL_THROTTLE", "GPU_RESET"}

// 设备结构体
type DeviceId struct {
	id uint32
}

// FailedMessage 重置clock错误信息
// @Description 包含重置clock操作失败时的设备ID和错误信息
type FailedMessage struct {
	// ID 设备ID
	ID int
	// ErrorMsg 错误信息
	ErrorMsg string
}

// 时钟类型映射
var rsmiClkNamesDict = map[string]RSMIClkType{
	"sclk":    RSMI_CLK_TYPE_SYS,
	"fclk":    RSMI_CLK_TYPE_DF,
	"dcefclk": RSMI_CLK_TYPE_DCEF,
	"socclk":  RSMI_CLK_TYPE_SOC,
	"mclk":    RSMI_CLK_TYPE_MEM,
}

var validLevels = map[string]DevPerfLevel{
	"auto":   RSMI_DEV_PERF_LEVEL_AUTO,
	"low":    RSMI_DEV_PERF_LEVEL_LOW,
	"high":   RSMI_DEV_PERF_LEVEL_HIGH,
	"manual": RSMI_DEV_PERF_LEVEL_MANUAL,
}

// 定义RAS错误状态字符串映射
var rasErrStaleMachine = []string{
	"NONE", "DISABLED", "UNKNOWN ERROR",
	"SING", "MULT", "POSITION", "ENABLED",
}

// RSMI 温度传感器类型常量
const (
	SENSOR_EDGE     = 0
	SENSOR_JUNCTION = 1
	SENSOR_MEMORY   = 2
	SENSOR_HBM0     = 3
	SENSOR_HBM1     = 4
	SENSOR_HBM2     = 5
	SENSOR_HBM3     = 6
)

// RSMI 温度传感器类型名称列表
var tempTypeList = []struct {
	Name string
	Type int
}{
	{"edge", SENSOR_EDGE},
	{"junction", SENSOR_JUNCTION},
	{"memory", SENSOR_MEMORY},
	{"HBM 0", SENSOR_HBM0},
	{"HBM 1", SENSOR_HBM1},
	{"HBM 2", SENSOR_HBM2},
	{"HBM 3", SENSOR_HBM3},
}

// TemperatureInfo 结构体表示一个设备的温度信息
type TemperatureInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	//  SensorTemps 传感器名称到温度的映射
	SensorTemps map[string]float64
}

// 固件块名称列表
var fwBlockNames = []string{
	"ASD", "CE", "DMCU", "MC", "ME", "MEC", "MEC2", "PFP",
	"RLC", "RLC SRLC", "RLC SRLG", "RLC SRLS", "SDMA", "SDMA2",
	"SMC", "SOS", "TA RAS", "TA XGMI", "UVD", "VCE", "VCN",
}

// FirmwareInfo 设备的固件信息
type FirmwareInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	//  FirmwareVer 固件块名称到版本信息的映射
	FirmwareVer map[string]string
}

var utilizationCounterName = []string{"GFX Activity", "Memory Activity"}

// DeviceUseInfo 设备使用信息列表
type DeviceUseInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	//  GPUUsePercent 设备使用率
	GPUUsePercent int
	//  利用率
	Utilization map[string]uint64
}

// DeviceMemVendorInfo 设备供应商信息
type DeviceMemVendorInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	//  Vendor 供应商信息
	Vendor string
}

// PcieBandwidthInfo 设备PCIe带宽信息
type PcieBandwidthInfo struct {
	//  DvInd 设备索引号
	DvInd int
	//  Sent 发送
	Sent float64
	//  Received 接收
	Received float64

	//  Bw bw
	Bw float64
}

// PcieReplayCountInfo 设备的PCIe重放计数信息
type PcieReplayCountInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	// 重放总数
	Count int64
}

// DevicePowerInfo 设备的功率信息
type DevicePowerInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	//  Power 设备功率
	Power int64
}

// DevicePowerPlayInfo 设备的GPU时钟频率和电压信息
type DevicePowerPlayInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	//  SCLK SCLK
	SCLK []string
	//  MCLK MCLK
	MCLK string
	//  DDC_CURVE DDC_CURVE
	DDC_CURVE []string
	//  OD_RANGE RANGE
	RANGE []string
}

// DeviceproductInfo 设备的产品信息列表
type DeviceproductInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	//  CardSeries 设备名称
	CardSeries string
	//  CardModel 设备子系统名称
	CardModel string
	//  CardVendor 设备供应商名称
	CardVendor string
	//  CardSKU SKU
	CardSKU string
}

// DeviceProfile 设备的电源配置文件信息
type DeviceProfile struct {
	//  DeviceID 设备索引号
	DeviceID int
	// Profiles 文件信息
	Profiles []string
}

var MemoryPageStatusStr = map[MemoryPageStatus]string{
	RSMI_MEM_PAGE_STATUS_RESERVED:     "reserved",
	RSMI_MEM_PAGE_STATUS_PENDING:      "pending",
	RSMI_MEM_PAGE_STATUS_UNRESERVABLE: "unreservable",
}

// DeviceSerialInfo 设备的序列号信息
type DeviceSerialInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	//  SerialNumber 设备序列号
	SerialNumber string
}

// DeviceUIdInfo 设备的唯一ID信息
type DeviceUIdInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	//  UId 设备唯一id
	UId string
}

// DeviceVBIOSInfo 设备的VBIOS版本信息
type DeviceVBIOSInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	//  VBIOS 版本信息
	VBIOS string
}

// DeviceVoltageInfo 设备电压信息
type DeviceVoltageInfo struct {
	//  DeviceID 设备索引号
	DeviceID int
	//  Voltage 电压信息
	Voltage int64 // 电压以毫伏为单位
}

// LinkType 表示 GPU 间互联类型
const (
	LinkTypePCIE        = "PCIE"     // 仅 PCIe 互联
	LinkTypeHSL         = "HSL"      // 仅 HSL 互联（非 Hyswitch）
	LinkTypeHSLHyswitch = "HYSWITCH" // Hyswitch
	LinkTypeHybrid      = "HYBRID"   // PCIe + HSL 混合互联
	LinkTypeUnknown     = "UNKNOWN"  // 无法识别
	LinkTypeNONE        = "NONE"
)

// BlocksInfo block信息
type BlocksInfo struct {
	// Block 类型
	Block string
	// State 状态
	State string
	// CE 错误数
	CE int64
	// UE 错误数
	UE int64
}

// NumaInfo 设备的Numa信息
type NumaInfo struct {
	// DeviceID 设备索引号
	DeviceID int
	// NumaNode numaNode值
	NumaNode int
	// NumaAffinity 关联信息
	NumaAffinity int
}

// 定义常量表示诊断结果
const (
	DiagResultPass   = "pass"    // 通过
	DiagResultSkip   = "skipped" // 跳过
	DiagResultWarn   = "warn"    // 警告
	DiagResultFail   = "fail"    // 失败
	DiagResultNotRun = "notrun"  // 未运行
)

// 定义常量表示健康检查结果
const (
	HealthStatusHealthy = "Healthy" // 健康
	HealthStatusWarning = "Warning" // 警告
	HealthStatusFailure = "Failure" // 失败
	HealthStatusSkipped = "Skipped" // 跳过检查
	HealthStatusUnknown = "unknown" // 未知
)

// 定义 Type 与数字的对应关系
var HealthType = map[int]string{
	1: "NumaTopology Health",
	2: "PcieBandwidth Health",
	3: "Power Health",
	4: "Memory Health",
	5: "Temperature Health",
	6: "Performance Health",
	7: "EccBlocks Health",
	8: "HCUUsage Health",
}

const (
	NumaTopologyHealth  = "NumaTopology Health"
	PcieBandwidthHealth = "PcieBandwidth Health"
	PowerHealth         = "Power Health"
	MemoryHealth        = "Memory Health"
	TemperatureHealth   = "Temperature Health"
	PerformanceHealth   = "Performance Health"
	EccBlocksHealth     = "EccBlocks Health"
	HCUUsageHealth      = "HCUUsage Health"
)

// 定义 DeviceModelInfo 结构体
type DeviceModelInfo struct {
	// Model HCU类型
	Model string
	// CUCount CU数量
	CUCount float64
	// MemorySize 内存大小
	MemorySize float64
}

// 定义设备型号的“枚举”类型
type DeviceModel int

const (
	K100_AI DeviceModel = iota
	K100_AI_Liquid
	K100_AI_Eco
	K100
	Z100
	Z100L
	BW200
)

// 型号到名称的映射
var modelName = map[DeviceModel]string{
	K100_AI:        "K100_AI",
	K100_AI_Liquid: "K100 AI Liquid",
	K100_AI_Eco:    "K100 AI eco",
	K100:           "K100",
	Z100:           "Z100",
	Z100L:          "Z100L",
	BW200:          "BW200",
}

// 型号到 CU 数量的映射
var computeUnit = map[DeviceModel]float64{
	K100_AI:        120,
	K100_AI_Liquid: 120,
	K100_AI_Eco:    128,
	K100:           120,
	Z100:           60,
	Z100L:          60,
	BW200:          80,
}

// 型号到内存大小的映射
var memorySize = map[DeviceModel]float64{
	K100_AI:        64,
	K100_AI_Liquid: 64,
	K100_AI_Eco:    64,
	K100:           64,
	Z100:           32,
	Z100L:          32,
	BW200:          64,
}

// DiagResult 诊断测试的结果
type DiagResult struct {
	// Status 诊断的结果状态
	Status string //诊断的结果状态
	// TestName 测试的名称
	TestName string //测试的名称，描述了这个测试的目的。
	// TestOutput 测试执行后生成的详细结果
	TestOutput string //存储测试执行后生成的详细结果.
	// ErrorCode 错误码
	ErrorCode int //错误码，如果测试失败，这里保存失败原因的代码。
	// ErrorMessage 错误信息
	ErrorMessage string //错误信息，描述测试失败时发生的错误。
}

// HCUResult HCU相关的诊断结果
type HCUResult struct {
	// HCU HCU 的编号
	HCU int //HCU 的编号或 ID
	//  RC 返回代码
	RC int //返回代码，指示诊断的总体成功或失败状态。
	// DiagResults HCU 上的所有单独诊断测试的结果
	DiagResults []DiagResult //该HCU上的所有单独诊断测试的结果（由多个 DiagResult 组成的数组）
}

// DiagResults 保存总体诊断结果，包括软件测试结果和每个 GPU 的测试结果。
type DiagResults struct {
	// DeviceNumber 设备数量信息
	DeviceNumber int //设备数量信息
	// Software 软件层面的诊断结果
	Software []DiagResult //保存软件层面的诊断结果，例如驱动程序、环境变量等
	// PerHCU  HCU 的诊断结果
	PerHCU []HCUResult //保存每个 HCU 的诊断结果
}

// HealthCheckConfig 用于存储健康检查的状态和选项
type HealthCheckConfig struct {
	// Enabled 是否开启健康检查
	Enabled bool `json:"enabled"`
	// Options 康检查项
	Options []string `json:"options"`
}

// 检查详情
type SystemWatch struct {
	//Type 检查项的类型
	Type string
	// Status 状态
	Status string
	// Error 错误的详细信息
	Error string
	// Result 结果
	Result interface{}
}

// DeviceHealth 设备健康结果
type DeviceHealth struct {
	// HCU HCU编号
	HCU uint
	// Status 状态
	Status string
	// Watches 检查详情
	Watches []SystemWatch
}

// HCUHealthStatus 单张 HCU 卡的健康检查结果，对应 hy-smi --healthcheck 输出的一行记录
type HCUHealthStatus struct {
	// HCU HCU 卡编号（从 0 开始）
	HCU uint `json:"hcu" example:"0"`
	// BusId PCIe Bus ID，格式如 0000:09:00.0
	BusId string `json:"busId" example:"0000:09:00.0"`
	// Status 健康状态，通常为 Healthy 或 Unhealthy
	Status string `json:"status" example:"Healthy"`
}

// DFBandwidthInfo DF 带宽信息
type DFBandwidthInfo struct {
	// ReadBW DF 读取带宽（GB/s）
	ReadBW float64 `json:"readBW" example:"120.5"`

	// WriteBW DF 写入带宽（GB/s）
	WriteBW float64 `json:"writeBW" example:"98.3"`

	// ReadWriteBW DF 读写综合带宽（GB/s）
	ReadWriteBW float64 `json:"readWriteBW" example:"210.8"`
}

const MAX_XHCL_LINK_NUM = 7

// XhclBandwidthInfo XHCL 带宽信息
type XhclBandwidthInfo struct {
	// Bw 每条 XHCL 链路的带宽值
	// 数组下标对应链路 ID，长度为 MAX_XHCL_LINK_NUM
	Bw [MAX_XHCL_LINK_NUM]float64 `json:"bw"`
}

const MAX_UMC_CHAN_NUM = 32

// MAX_SE_CNT 每个设备最大 SE 数量
const MAX_SE_CNT = 8

// SEUsageInfo 各 Shader Engine（SE）瞬时占用率信息。
// 对应 rsmi_dev_se_util_get / hy-smi --showseuse；Percent[i] 为第 i 个 SE 的占用率。
type SEUsageInfo struct {
	// Percent 各 SE 占用率，下标 i 对应 SE 索引，长度为 MAX_SE_CNT（8）
	Percent [MAX_SE_CNT]float32 `json:"percent"`
}

// UMCBandwidthInfo UMC 带宽信息（按通道统计）
type UMCBandwidthInfo struct {
	// ReadBW 每个 UMC 通道的读取带宽
	ReadBW [MAX_UMC_CHAN_NUM]float64

	// WriteBW 每个 UMC 通道的写入带宽
	WriteBW [MAX_UMC_CHAN_NUM]float64

	// ReadWriteBW 每个 UMC 通道的读写混合带宽
	ReadWriteBW [MAX_UMC_CHAN_NUM]float64
}

// Go 等价的枚举类型
const (
	RSMI_DF_BW_TYPE_R   = C.RSMI_DF_BW_TYPE_R
	RSMI_DF_BW_TYPE_W   = C.RSMI_DF_BW_TYPE_W
	RSMI_DF_BW_TYPE_R_W = C.RSMI_DF_BW_TYPE_R_W
	RSMI_DF_BW_TYPE_ALL = C.RSMI_DF_BW_TYPE_ALL
)

type NvmlReturn C.nvmlReturn_t

const (
	NVML_SUCCESS                         NvmlReturn = C.NVML_SUCCESS                         // 操作成功
	NVML_ERROR_UNINITIALIZED             NvmlReturn = C.NVML_ERROR_UNINITIALIZED             // NVML 未初始化（未调用 nvmlInit）
	NVML_ERROR_INVALID_ARGUMENT          NvmlReturn = C.NVML_ERROR_INVALID_ARGUMENT          // 参数无效
	NVML_ERROR_NOT_SUPPORTED             NvmlReturn = C.NVML_ERROR_NOT_SUPPORTED             // 当前设备不支持该操作
	NVML_ERROR_NO_PERMISSION             NvmlReturn = C.NVML_ERROR_NO_PERMISSION             // 当前用户无操作权限
	NVML_ERROR_ALREADY_INITIALIZED       NvmlReturn = C.NVML_ERROR_ALREADY_INITIALIZED       // 已初始化（已弃用，保留）
	NVML_ERROR_NOT_FOUND                 NvmlReturn = C.NVML_ERROR_NOT_FOUND                 // 查询对象未找到
	NVML_ERROR_INSUFFICIENT_SIZE         NvmlReturn = C.NVML_ERROR_INSUFFICIENT_SIZE         // 输入参数长度不足
	NVML_ERROR_INSUFFICIENT_POWER        NvmlReturn = C.NVML_ERROR_INSUFFICIENT_POWER        // 外部电源连接异常
	NVML_ERROR_DRIVER_NOT_LOADED         NvmlReturn = C.NVML_ERROR_DRIVER_NOT_LOADED         // NVIDIA 驱动未加载
	NVML_ERROR_TIMEOUT                   NvmlReturn = C.NVML_ERROR_TIMEOUT                   // 用户传入超时时间已过
	NVML_ERROR_IRQ_ISSUE                 NvmlReturn = C.NVML_ERROR_IRQ_ISSUE                 // NVIDIA 内核检测到GPU的中断问题
	NVML_ERROR_LIBRARY_NOT_FOUND         NvmlReturn = C.NVML_ERROR_LIBRARY_NOT_FOUND         // NVML共享库未找到/加载失败
	NVML_ERROR_FUNCTION_NOT_FOUND        NvmlReturn = C.NVML_ERROR_FUNCTION_NOT_FOUND        // NVML本地版本未实现此功能
	NVML_ERROR_CORRUPTED_INFOROM         NvmlReturn = C.NVML_ERROR_CORRUPTED_INFOROM         // infoROM损坏
	NVML_ERROR_GPU_IS_LOST               NvmlReturn = C.NVML_ERROR_GPU_IS_LOST               // GPU已掉线或不可访问
	NVML_ERROR_RESET_REQUIRED            NvmlReturn = C.NVML_ERROR_RESET_REQUIRED            // GPU需重置后才能再次使用
	NVML_ERROR_OPERATING_SYSTEM          NvmlReturn = C.NVML_ERROR_OPERATING_SYSTEM          // 操作系统/控制组阻止了GPU控制
	NVML_ERROR_LIB_RM_VERSION_MISMATCH   NvmlReturn = C.NVML_ERROR_LIB_RM_VERSION_MISMATCH   // 驱动/库版本不匹配
	NVML_ERROR_IN_USE                    NvmlReturn = C.NVML_ERROR_IN_USE                    // GPU正在被使用，无法执行操作
	NVML_ERROR_MEMORY                    NvmlReturn = C.NVML_ERROR_MEMORY                    // 内存不足
	NVML_ERROR_NO_DATA                   NvmlReturn = C.NVML_ERROR_NO_DATA                   // 无数据
	NVML_ERROR_VGPU_ECC_NOT_SUPPORTED    NvmlReturn = C.NVML_ERROR_VGPU_ECC_NOT_SUPPORTED    // vGPU 操作不支持 ECC
	NVML_ERROR_INSUFFICIENT_RESOURCES    NvmlReturn = C.NVML_ERROR_INSUFFICIENT_RESOURCES    // 关键资源不足（非内存）
	NVML_ERROR_FREQ_NOT_SUPPORTED        NvmlReturn = C.NVML_ERROR_FREQ_NOT_SUPPORTED        // 不支持指定的频率
	NVML_ERROR_ARGUMENT_VERSION_MISMATCH NvmlReturn = C.NVML_ERROR_ARGUMENT_VERSION_MISMATCH // 版本参数无效或不支持
	NVML_ERROR_DEPRECATED                NvmlReturn = C.NVML_ERROR_DEPRECATED                // 所请求功能已废弃
	NVML_ERROR_NOT_READY                 NvmlReturn = C.NVML_ERROR_NOT_READY                 // 系统未准备好处理请求
	NVML_ERROR_UNKNOWN                   NvmlReturn = C.NVML_ERROR_UNKNOWN                   // 内部驱动未知错误
)

var nvmlErrorCodeMap = map[NvmlReturn]string{
	NVML_SUCCESS:                         "SUCCESS",
	NVML_ERROR_UNINITIALIZED:             "UNINITIALIZED",
	NVML_ERROR_INVALID_ARGUMENT:          "INVALID_ARGUMENT",
	NVML_ERROR_NOT_SUPPORTED:             "NOT_SUPPORTED",
	NVML_ERROR_NO_PERMISSION:             "NO_PERMISSION",
	NVML_ERROR_ALREADY_INITIALIZED:       "ALREADY_INITIALIZED",
	NVML_ERROR_NOT_FOUND:                 "NOT_FOUND",
	NVML_ERROR_INSUFFICIENT_SIZE:         "INSUFFICIENT_SIZE",
	NVML_ERROR_INSUFFICIENT_POWER:        "INSUFFICIENT_POWER",
	NVML_ERROR_DRIVER_NOT_LOADED:         "DRIVER_NOT_LOADED",
	NVML_ERROR_TIMEOUT:                   "TIMEOUT",
	NVML_ERROR_IRQ_ISSUE:                 "IRQ_ISSUE",
	NVML_ERROR_LIBRARY_NOT_FOUND:         "LIBRARY_NOT_FOUND",
	NVML_ERROR_FUNCTION_NOT_FOUND:        "FUNCTION_NOT_FOUND",
	NVML_ERROR_CORRUPTED_INFOROM:         "CORRUPTED_INFOROM",
	NVML_ERROR_GPU_IS_LOST:               "GPU_IS_LOST",
	NVML_ERROR_RESET_REQUIRED:            "RESET_REQUIRED",
	NVML_ERROR_OPERATING_SYSTEM:          "OPERATING_SYSTEM",
	NVML_ERROR_LIB_RM_VERSION_MISMATCH:   "LIB_RM_VERSION_MISMATCH",
	NVML_ERROR_IN_USE:                    "IN_USE",
	NVML_ERROR_MEMORY:                    "MEMORY",
	NVML_ERROR_NO_DATA:                   "NO_DATA",
	NVML_ERROR_VGPU_ECC_NOT_SUPPORTED:    "VGPU_ECC_NOT_SUPPORTED",
	NVML_ERROR_INSUFFICIENT_RESOURCES:    "INSUFFICIENT_RESOURCES",
	NVML_ERROR_FREQ_NOT_SUPPORTED:        "FREQ_NOT_SUPPORTED",
	NVML_ERROR_ARGUMENT_VERSION_MISMATCH: "ARGUMENT_VERSION_MISMATCH",
	NVML_ERROR_DEPRECATED:                "DEPRECATED",
	NVML_ERROR_NOT_READY:                 "NOT_READY",
	NVML_ERROR_UNKNOWN:                   "UNKNOWN",
}

// nvmlErrorCodeMap 是NVML返回码到可读字符串的映射
//var nvmlErrorCodeMap = map[NvmlReturn]string{
//	NVML_SUCCESS:                         "SUCCESS",
//	NVML_ERROR_UNINITIALIZED:             "UNINITIALIZED (未初始化)",
//	NVML_ERROR_INVALID_ARGUMENT:          "INVALID_ARGUMENT (参数无效)",
//	NVML_ERROR_NOT_SUPPORTED:             "NOT_SUPPORTED (不支持该操作)",
//	NVML_ERROR_NO_PERMISSION:             "NO_PERMISSION (无权限)",
//	NVML_ERROR_ALREADY_INITIALIZED:       "ALREADY_INITIALIZED (已初始化)",
//	NVML_ERROR_NOT_FOUND:                 "NOT_FOUND (未找到对象)",
//	NVML_ERROR_INSUFFICIENT_SIZE:         "INSUFFICIENT_SIZE (长度不足)",
//	NVML_ERROR_INSUFFICIENT_POWER:        "INSUFFICIENT_POWER (电源不足)",
//	NVML_ERROR_DRIVER_NOT_LOADED:         "DRIVER_NOT_LOADED (驱动未加载)",
//	NVML_ERROR_TIMEOUT:                   "TIMEOUT (超时)",
//	NVML_ERROR_IRQ_ISSUE:                 "IRQ_ISSUE (中断异常)",
//	NVML_ERROR_LIBRARY_NOT_FOUND:         "LIBRARY_NOT_FOUND (未找到NVML库)",
//	NVML_ERROR_FUNCTION_NOT_FOUND:        "FUNCTION_NOT_FOUND (库不支持此函数)",
//	NVML_ERROR_CORRUPTED_INFOROM:         "CORRUPTED_INFOROM (infoROM损坏)",
//	NVML_ERROR_GPU_IS_LOST:               "GPU_IS_LOST (GPU丢失)",
//	NVML_ERROR_RESET_REQUIRED:            "RESET_REQUIRED (需重启/重置)",
//	NVML_ERROR_OPERATING_SYSTEM:          "OPERATING_SYSTEM (操作系统错误)",
//	NVML_ERROR_LIB_RM_VERSION_MISMATCH:   "LIB_RM_VERSION_MISMATCH (库版本不匹配)",
//	NVML_ERROR_IN_USE:                    "IN_USE (GPU被占用)",
//	NVML_ERROR_MEMORY:                    "MEMORY (内存错误)",
//	NVML_ERROR_NO_DATA:                   "NO_DATA (无数据)",
//	NVML_ERROR_VGPU_ECC_NOT_SUPPORTED:    "VGPU_ECC_NOT_SUPPORTED (vGPU ECC不支持)",
//	NVML_ERROR_INSUFFICIENT_RESOURCES:    "INSUFFICIENT_RESOURCES (资源不足)",
//	NVML_ERROR_FREQ_NOT_SUPPORTED:        "FREQ_NOT_SUPPORTED (频率不支持)",
//	NVML_ERROR_ARGUMENT_VERSION_MISMATCH: "ARGUMENT_VERSION_MISMATCH (版本不支持)",
//	NVML_ERROR_DEPRECATED:                "DEPRECATED (已废弃)",
//	NVML_ERROR_NOT_READY:                 "NOT_READY (未就绪)",
//	NVML_ERROR_UNKNOWN:                   "UNKNOWN (未知错误)",
//}

// NvmlDeviceAttributes 描述了单个 GPU 或 MIG 设备的硬件属性信息。
type NvmlDeviceAttributes struct {
	// Index 设备索引，表示 GPU 或 MIG 设备在系统中的编号。
	Index uint32

	// CUCount 计算单元（Compute Unit，CU）数量。
	CUCount uint32

	// MemorySizeMB 设备的显存容量，单位为 MB。
	MemorySizeMB uint64

	// UUID 设备的全局唯一标识符（UUID），用于唯一标识 GPU 或 MIG 设备。
	UUID string

	// Name 设备的名称字符串。
	Name string

	// GPUInstanceSliceCount GPU 实例中的 slice 数量，仅对 MIG 设备有效。若为 GPU 物理设备，则该字段为 0。
	GPUInstanceSliceCount uint32

	// ComputeInstanceSliceCount Compute 实例中的 slice 数量，仅对 MIG 设备有效。若为 GPU 物理设备，则该字段为 0。
	ComputeInstanceSliceCount uint32
}

type MIGDevice C.nvmlDevice_t

type DMIGpuInstance C.nvmlGpuInstance_t

type DMIComputeInstance C.nvmlComputeInstance_t

const (
	NVML_DEVICE_MIG_DISABLE = C.NVML_DEVICE_MIG_DISABLE
	NVML_DEVICE_MIG_ENABLE  = C.NVML_DEVICE_MIG_ENABLE
)

// NvmlGpuInstanceProfileInfo 映射C结构体的Go定义
type NvmlGpuInstanceProfileInfo struct {
	ID          uint32 // profile的唯一ID（在本GPU内唯一）
	GiCountMax  uint32 // 本profile最多支持的实例数量
	CuCount     uint32 // 计算单元数量
	GpuSliceCnt uint32 // 本profile分配的GPU Slice数量
	MemSizeMB   uint64 // 本profile的显存容量（单位MB）
	Name        string // profile名称
}

// NvmlComputeInstanceProfileInfo 结构体
type NvmlComputeInstanceProfileInfo struct {
	ID          uint32 // profile唯一id
	CiCountMax  uint32 // 能创建几个该profile的compute instance
	CuCount     uint32 // CU数量
	GpuSliceCnt uint32 // slice数量
	Name        string
}

// GpuInstancePlacement 描述 MIG 实例在物理 GPU 上的计算单元摆放
type GpuInstancePlacement struct {
	Start uint32 // 起始 compute slice 的索引
	Size  uint32 // 占用的 compute slice 数量
}

// GpuInstanceInfo 包含 MIG 实例的详细信息
type GpuInstanceInfo struct {
	Device    uintptr              // 父物理设备句柄（可选，uintptr 方便在 Go/CGO 传递句柄）
	ID        uint32               // 实例唯一 ID（giId）
	ProfileID uint32               // Profile 类型编号
	Placement GpuInstancePlacement // 当前实例在设备上占用的资源位置信息
}

type ComputeInstanceRemainInfo struct {
	GiID          uint32 // MIG实例唯一ID
	ProfileID     uint32 // MIG实例profile类型
	CiRemainCount uint32 // 剩余可分配的compute instance数
}

// MigInfo 表示单个MIG分区的信息
type MigInfo struct {
	DvInd             int
	MigId             int
	Name              string
	UUID              string
	ComputeUnit       uint32
	MemoryTotal       uint64
	GpuInstanceId     uint32
	ComputeInstanceId uint32
	PciBusNumber      string
	GiProfileId       int
	CiProfileId       int
}

type NvmlPciInfo struct {
	Domain   uint32
	Bus      uint32
	Device   uint32
	Function uint32
	BusID    string
}

// GiInfo 存储GI配置信息
type GiInfo struct {
	Pci            string // PCI设备地址
	Id             int    // GI实例ID
	PipeMask       string // 管道掩码
	GpuSliceMask   string // GPU切片掩码
	CuMask1        string // 第一个计算单元掩码
	CuMask2        string // 第二个计算单元掩码
	ProfileId      int    // 配置文件ID
	GiCountMax     int    // 最大GI实例数
	CuCount        int    // 计算单元数量
	GpuSliceCount  int    // GPU切片数量
	MemorySizeMB   int    // 内存大小(MB)
	PlacementStart int    // 放置起始位置
	PlacementSize  int    // 放置大小
}

// CiInfo 存储CI配置信息
type CiInfo struct {
	Pci            string // PCI设备地址
	GiId           int    // 所属GI实例ID
	Id             int    // CI实例ID
	PipeMask       string // 管道掩码
	GpuSliceMask   string // GPU切片掩码
	CuMask1        string // 第一个计算单元掩码
	CuMask2        string // 第二个计算单元掩码
	ProfileId      int    // 配置文件ID
	CiCountMax     int    // 最大CI实例数
	CuCount        int    // 计算单元数量
	GpuSliceCount  int    // GPU切片数量
	PlacementStart int    // 放置起始位置
	PlacementSize  int    // 放置大小
	MigUUID        string // MIG实例的唯一标识符
}

// MigConfig MIG配置信息
type MigConfig struct {
	DvInd int    // GPU设备ID
	Gi    GiInfo // GI配置信息
	Ci    CiInfo // CI配置信息
}

// Info 保存最终要输出的 MIG 信息
type Info struct {
	GpuSliceCount int    // 申请的 slice 数
	MemorySizeMB  uint64 // 申请的显存（MB）
	Name          string // 格式化后的 MIG 名称
}

const (
	linkAllID = 0xff             // 表示一次查询所有 link
	chanAllID = MAX_UMC_CHAN_NUM // 传 32 获取全部 channel（按你的说明）
	delay     = 10               // 采样间隔默认为 10ms
)

// DeviceLinkSum 表示单个设备（DvInd）的 Link 带宽汇总结果
type DeviceLinkSum struct {
	// DvInd 设备索引
	DvInd int

	// Recv 接收方向（direction=0）所有 Link 的带宽总和
	Recv float64

	// Send 发送方向（direction=1）所有 Link 的带宽总和
	Send float64

	// Err 该设备查询过程中的错误信息（成功时为空字符串）
	Err string

	// Links 每个 Link 的接收/发送带宽明细
	Links []LinkBandwidth
}

// LinkBandwidth 表示单个 Link 的接收/发送带宽信息
type LinkBandwidth struct {
	// LinkId Link 索引
	LinkId int

	// Recv 接收方向（direction=0）的带宽
	Recv float64

	// Send 发送方向（direction=1）的带宽
	Send float64
}

// DeviceUmcSum 表示单个设备（device index）的 UMC 带宽求和结果。
type DeviceUmcSum struct {
	DvInd     int     // DvInd: 设备索引（从 0 开始）
	Read      float64 // Read: 所有 32 个 channel 的读带宽之和，保留两位小数；失败则为 0
	Write     float64 // Write: 所有 32 个 channel 的写带宽之和，保留两位小数；失败则为 0
	ReadWrite float64 // ReadWrite: 所有 32 个 channel 的读写带宽之和，保留两位小数；失败则为 0
	Err       string  // Err: 如果该设备任一方向或查询发生错误，记录错误信息；成功时为空串
}

const (
	// XhclLinkDown 表示 XHCL 链路不可用 / 未连接
	XhclLinkDown uint32 = 0

	// XhclLinkUp 表示 XHCL 链路正常工作
	XhclLinkUp uint32 = 1
)

// XhclLinkState 表示单条 XHCL（HSL）物理链路的状态信息
type XhclLinkState struct {
	LinkID int    // 链路编号（从 0 开始）
	State  uint32 // 驱动返回的原始链路状态值
	Up     bool   // 是否为 UP（State == XhclLinkUp）
}

// XhclRemoteBdf 表示一条 XHCL 链路及其远端设备 BDF 信息
type XhclRemoteBdf struct {
	LinkID int    // XHCL 链路 ID
	BdfID  uint64 // 远端设备 BDF ID
}

const (
	pciIDsPath   = "/usr/share/hwdata/pci.ids"
	targetVendor = "1d94"
)

var pciDeviceID2Name = make(map[string]string)

const updateIDsPath = "/opt/hyhal/bin/update_ids"

var updateIDsMap map[string]string

// DirectLinkInfo 描述当前 GPU 与某一远端 GPU 的直连关系
type DirectLinkInfo struct {
	RemoteDvInd int    // 远端 GPU 的设备索引（-1 表示未匹配到）
	RemoteBdfID uint64 // 远端 GPU 的 BDFID
	LinkType    string // 链路类型：PCIE / HSL
	Weight      int    // 与该 GPU 的直连链路数量（权重）
}

// GpuInterconnectInfo 描述单张 GPU 的互联拓扑信息
type GpuInterconnectInfo struct {
	LinkType  string           // 总体互联类型：PCIE / HSL / HSL_HYSWITCH / HYBRID
	CardCount int              // 参与互联的 GPU 数量（包含自身）
	Links     []DirectLinkInfo // 与当前 GPU 直连的 GPU 信息
}

// HCULinkInfo 描述两张 HCU 之间的互联关系
type HCULinkInfo struct {
	SrcDvInd int    // 源 HCU 索引
	DstDvInd int    // 目标 HCU 索引
	BdfID    uint64 // 目标 HCU 的 BDFID
	PciID    string // 目标 HCU 的 PCI ID
	LinkType string // PCIE / HSL / HYSWITCH / NONE
	Weight   int    // 链路权重
	Hops     int    // 跳数（目前为 0）
}

// HCUInterconnectMatrix 描述整机 HCU 互联矩阵
type HCUInterconnectMatrix struct {
	DeviceCount int
	Matrix      [][]HCULinkInfo // [src][dst]
}

// HyVersionInfo 海光扩展版本号（对应 rsmi_hy_version_t）
// 与 Version() 返回的 ROCm 版本字符串不同，这里是 Hygon 专有版本的结构化格式
type HyVersionInfo struct {
	Major uint32 `json:"major"`
	Minor uint32 `json:"minor"`
}
