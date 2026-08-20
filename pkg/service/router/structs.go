/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package router

/*
#cgo CFLAGS: -Wall -I../../dcgm/include
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
	"time"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

// Response represents a basic structure for API responses.
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type CompatibilityResult struct {
	Status            string `json:"status"`
	CardModel         string `json:"cardModel"`
	CurrentDriver     string `json:"currentDriver"`
	CurrentDTK        string `json:"currentDTK"`
	RecommendedDriver string `json:"recommendedDriver"`
	RecommendedDTK    string `json:"recommendedDTK"`
}

type ThrottleFlags struct {
	Thermal    bool
	Power      bool
	Slowdown   bool
	BoardLimit bool
}

type HyVersionInfo struct {
	Major uint32 `json:"major"`
	Minor uint32 `json:"minor"`
}

// SuccessResponse creates a successful response payload.
func SuccessResponse(data interface{}) Response {
	return Response{
		Message: "成功",
		Data:    data,
	}
}

// ErrorResponse creates an error response payload.
func ErrorResponse(err interface{}) Response {
	return Response{
		Message: "失败",
		Error:   err,
	}
}

// PcieBandwidth 表示设备的 PCIe 带宽信息
// swagger:model PcieBandwidth
type PcieBandwidth struct {
	// TransferRate 表示传输速率的频率信息
	TransferRate Frequencies `json:"transferRate"`

	// Lanes 表示 PCIe 通道的配置
	Lanes [33]uint32 `json:"lanes"`
}

// Frequencies 表示设备支持的频率信息
// swagger:model RSMIFrequencies
type Frequencies struct {
	// NumSupported 表示设备支持的频率数量
	NumSupported uint32 `json:"numSupported"`

	// Current 表示当前使用的频率
	Current uint32 `json:"current"`

	// Frequency 表示设备支持的频率列表
	Frequency [33]uint64 `json:"frequency"`
}

type RSNIPowerProfilePresetMasks C.rsmi_power_profile_preset_masks_t

const (
	RSMI_PWR_PROF_PRST_CUSTOM_MASK       RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_CUSTOM_MASK       //!< Custom Power Profile
	RSMI_PWR_PROF_PRST_VIDEO_MASK        RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VIDEO_MASK        //!< Video Power Profile
	RSMI_PWR_PROF_PRST_POWER_SAVING_MASK RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_POWER_SAVING_MASK //!< Power Saving Profile
	RSMI_PWR_PROF_PRST_COMPUTE_MASK      RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_COMPUTE_MASK      //!< Compute Saving Profile
	RSMI_PWR_PROF_PRST_VR_MASK           RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VR_MASK           //!< VR Power Profile

	//!< 3D Full Screen Power Profile
	RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK
	RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT   RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT //!< Default Boot Up Profile
	RSMI_PWR_PROF_PRST_LAST             RSNIPowerProfilePresetMasks = RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT

	//!< Invalid power profile
	RSMI_PWR_PROF_PRST_INVALID RSNIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_INVALID
)

type RSMIRetiredPageRecord struct {
	PageAddress uint64               //!< Start address of page
	PageSize    uint64               //!< Page size
	Status      RSMIMemoryPageStatus //!< Page "reserved" status
}

type RSMIMemoryPageStatus C.rsmi_memory_page_status_t

const (
	RSMI_MEM_PAGE_STATUS_RESERVED     RSMIMemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_RESERVED
	RSMI_MEM_PAGE_STATUS_PENDING      RSMIMemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_PENDING
	RSMI_MEM_PAGE_STATUS_UNRESERVABLE RSMIMemoryPageStatus = C.RSMI_MEM_PAGE_STATUS_UNRESERVABLE
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

type RSMIUtilizationCounterType uint32

const (
	RSMI_UTILIZATION_COUNTER_FIRST RSMIUtilizationCounterType = C.RSMI_UTILIZATION_COUNTER_FIRST
	RSMI_COARSE_GRAIN_GFX_ACTIVITY RSMIUtilizationCounterType = C.RSMI_COARSE_GRAIN_GFX_ACTIVITY
	RSMI_COARSE_GRAIN_MEM_ACTIVITY RSMIUtilizationCounterType = C.RSMI_COARSE_GRAIN_MEM_ACTIVITY
	RSMI_UTILIZATION_COUNTER_LAST  RSMIUtilizationCounterType = C.RSMI_UTILIZATION_COUNTER_LAST
)

// @swagignore
type RSMIUtilizationCounter struct {

	// Type 表示利用率计数器的类型
	Type RSMIUtilizationCounterType

	// Value 表示计数器的值
	Value uint64
}

type RSMIClkType C.rsmi_clk_type_t

const (
	RSMI_CLK_TYPE_SYS  RSMIClkType = C.RSMI_CLK_TYPE_SYS
	RSMI_CLK_TYPE_DF   RSMIClkType = C.RSMI_CLK_TYPE_DF
	RSMI_CLK_TYPE_DCEF RSMIClkType = C.RSMI_CLK_TYPE_DCEF
	RSMI_CLK_TYPE_SOC  RSMIClkType = C.RSMI_CLK_TYPE_SOC
	RSMI_CLK_TYPE_MEM  RSMIClkType = C.RSMI_CLK_TYPE_MEM
	RSMI_CLK_TYPE_PCIE RSMIClkType = C.RSMI_CLK_TYPE_PCIE
	RSMI_CLK_INVALID   RSMIClkType = C.RSMI_CLK_INVALID
)

type RSMIOdVoltFreqData struct {
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

type RSMIDevPerfLevel C.rsmi_dev_perf_level_t

const (
	RSMI_DEV_PERF_LEVEL_AUTO            RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_AUTO
	RSMI_DEV_PERF_LEVEL_FIRST           RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_FIRST
	RSMI_DEV_PERF_LEVEL_LOW             RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_LOW
	RSMI_DEV_PERF_LEVEL_HIGH            RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_HIGH
	RSMI_DEV_PERF_LEVEL_MANUAL          RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_MANUAL
	RSMI_DEV_PERF_LEVEL_STABLE_STD      RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_STD
	RSMI_DEV_PERF_LEVEL_STABLE_PEAK     RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_PEAK
	RSMI_DEV_PERF_LEVEL_STABLE_MIN_MCLK RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_MIN_MCLK
	RSMI_DEV_PERF_LEVEL_STABLE_MIN_SCLK RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_STABLE_MIN_SCLK
	RSMI_DEV_PERF_LEVEL_DETERMINISM     RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_DETERMINISM
	RSMI_DEV_PERF_LEVEL_LAST            RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_LAST
	RSMI_DEV_PERF_LEVEL_UNKNOWN         RSMIDevPerfLevel = C.RSMI_DEV_PERF_LEVEL_UNKNOWN
)

// 系统支持的配置文件
type RSMIBitField C.rsmi_bit_field_t

// 当前激活的电源配置文件
type RSMIPowerProfilePresetMasks C.rsmi_power_profile_preset_masks_t

// 定义 power profile preset masks 的枚举类型
const (
	RSMIPowerProfPrstCustomMask      RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_CUSTOM_MASK       // Custom Power Profile
	RSMIPowerProfPrstVideoMask       RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VIDEO_MASK        // Video Power Profile
	RSMIPowerProfPrstPowerSavingMask RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_POWER_SAVING_MASK // Power Saving Profile
	RSMIPowerProfPrstComputeMask     RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_COMPUTE_MASK      // Compute Saving Profile
	RSMIPowerProfPrstVRMask          RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_VR_MASK           // VR Power Profile
	RSMIPowerProfPrst3DFullScrMask   RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_3D_FULL_SCR_MASK  // 3D Full Screen Power Profile
	RSMIPowerProfPrstBootupDefault   RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT    // Default Boot Up Profile
	RSMIPowerProfPrstLast            RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_LAST              // Last Profile (same as Bootup Default)
	RSMIPowerProfPrstInvalid         RSMIPowerProfilePresetMasks = C.RSMI_PWR_PROF_PRST_INVALID           // Invalid power profile
)

// RSMPowerProfileStatus  电源配置文件状态信息
type RSMPowerProfileStatus struct {
	// AvailableProfiles  哪些配置文件被系统支持
	AvailableProfiles RSMIBitField
	// Current 当前激活的电源配置文件
	Current RSMIPowerProfilePresetMasks
	//  NumProfiles 可用的电源配置文件数量
	NumProfiles uint32
}

type RSMIVersion struct {
	Major uint32
	Minor uint32
	Patch uint32
	Build string
}

type RSMISwComponent C.rsmi_sw_component_t

const (
	RSMISwCompFirst  RSMISwComponent = C.RSMI_SW_COMP_FIRST
	RSMISwCompDriver RSMISwComponent = C.RSMI_SW_COMP_DRIVER
	RSMISwCompLast   RSMISwComponent = C.RSMI_SW_COMP_LAST
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

// 进程的信息
type RSMIProcessInfo struct {
	ProcessID   uint32
	Pasid       uint32
	VramUsage   uint64
	SdmaUsage   uint64
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

// IO链路类型
type RSMIIOLinkType C.RSMI_IO_LINK_TYPE

const (
	RSMIIOLinkTypeUndefined      RSMIIOLinkType = C.RSMI_IOLINK_TYPE_UNDEFINED
	RSMIIOLinkTypePCIExpress     RSMIIOLinkType = C.RSMI_IOLINK_TYPE_PCIEXPRESS
	RSMIIOLinkTypeHSL            RSMIIOLinkType = C.RSMI_IOLINK_TYPE_XGMI
	RSMIIOLinkTypeNumIOLinkTypes RSMIIOLinkType = C.RSMI_IOLINK_TYPE_NUMIOLINKTYPES
	RSMIIOLinkTypeSize           RSMIIOLinkType = C.RSMI_IOLINK_TYPE_SIZE
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
	// Clk 备系统时钟速度列表
	Clk float64
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
	// PciBusNumber 设备的总线号
	PciBusNumber string
	// MemoryTotal 设备的内存总量
	MemoryTotal float64
	// MemoryUsed 设备的已使用内存量
	MemoryUsed float64
	// ComputeUnit 设备的计算单元数量
	ComputeUnit float64
}

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

// BlocksInfo 设备 Block 状态、CE 和 UE 信息
type BlocksInfo struct {
	// Block Block 名称
	Block string `json:"block" example:"UMC"`
	// State Block 状态
	State string `json:"state" example:"OK"`

	// CE Correctable Error 数量
	CE int64 `json:"ce" example:"0"`

	// UE Uncorrectable Error 数量
	UE int64 `json:"ue" example:"0"`
}

// Device 物理设备的详细信息
type Device struct {
	MinorNumber               int          `json:"minorNumber"`               // 设备的索引号
	PciBusNumber              string       `json:"pciBusNumber"`              // 设备的总线编号
	DeviceId                  string       `json:"deviceId"`                  // 设备的唯一标识符
	SubSystemName             string       `json:"subSystemName"`             // 设备的子系统名称
	Temperature               float64      `json:"temperature"`               // 当前温度（摄氏度）
	PowerUsage                float64      `json:"powerUsage"`                // 当前功耗（瓦特）
	PowerCap                  float64      `json:"powerCap"`                  // 功耗上限（瓦特）
	MemoryCap                 float64      `json:"memoryCap"`                 // 显存容量（MB）
	MemoryUsed                float64      `json:"memoryUsed"`                // 已使用显存（MB）
	UtilizationRate           float64      `json:"utilizationRate"`           // 设备利用率（百分比）
	PcieBwMb                  float64      `json:"pcieBwMb"`                  // PCIe带宽（MB/s）
	Clk                       float64      `json:"clk"`                       // 当前时钟频率（MHz）
	ComputeUnitCount          float64      `json:"computeUnitCount"`          // 计算单元总数
	ComputeUnitRemainingCount uint64       `json:"computeUnitRemainingCount"` // 剩余可用计算单元数量
	MemoryRemaining           uint64       `json:"memoryRemaining"`           // 剩余可用显存（MB）
	MaxVDeviceCount           int          `json:"maxVDeviceCount"`           // 最大虚拟设备数量
	VDeviceCount              int          `json:"vDeviceCount"`              // 当前虚拟设备数量
	BlocksInfos               []BlocksInfo `json:"blocksInfos"`               // 设备Block信息列表
}

// DMIVDeviceInfo 虚拟设备信息
type DMIVDeviceInfo struct {
	Name             string `json:"name"`             // 虚拟设备名称
	ComputeUnitCount int    `json:"computeUnitCount"` // 虚拟设备计算单元数量
	GlobalMemSize    uint64 `json:"globalMemSize"`    // 全局内存大小（MB） @swagignore
	UsageMemSize     uint64 `json:"usageMemSize"`     // 已使用内存大小（MB） @swagignore
	ContainerID      uint64 `json:"containerID"`      // 所属容器ID
	DeviceID         int    `json:"deviceID"`         // 设备ID
	Percent          int    `json:"percent"`          // 使用百分比
	VMinorNumber     int    `json:"vminorNumber"`     // 虚拟设备索引号
	PciBusNumber     string `json:"pciBusNumber"`     // 虚拟设备总线编号
}

// PhysicalDeviceInfo 物理设备信息
type PhysicalDeviceInfo struct {
	// Device 物理设备的详细信息
	// @swagger:allOf
	Device `json:",inline"`

	// VirtualDevices 该物理设备上关联的虚拟设备信息列表
	VirtualDevices []DMIVDeviceInfo `json:"virtualDevices"` // 虚拟设备列表
}

// 定义事件通知类型名称
var notificationTypeNames = []string{"VM_FAULT", "THERMAL_THROTTLE", "GPU_RESET"}

// 设备结构体
type DeviceId struct {
	id uint32
}

// VDeviceByDvIndResp 虚拟设备索引返回结构
type VDeviceByDvIndResp struct {
	// VDeviceCount 虚拟设备数量
	VDeviceCount int `json:"vDeviceCount"`
	// VDevInds 虚拟设备索引列表
	VDevInds []int `json:"vDevInds"`
}

// 时钟类型映射
var rsmiClkNamesDict = map[string]RSMIClkType{
	"sclk":    RSMI_CLK_TYPE_SYS,
	"fclk":    RSMI_CLK_TYPE_DF,
	"dcefclk": RSMI_CLK_TYPE_DCEF,
	"socclk":  RSMI_CLK_TYPE_SOC,
	"mclk":    RSMI_CLK_TYPE_MEM,
}

var validLevels = map[string]RSMIDevPerfLevel{
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

// 固件块名称列表
var fwBlockNames = []string{
	"ASD", "CE", "DMCU", "MC", "ME", "MEC", "MEC2", "PFP",
	"RLC", "RLC SRLC", "RLC SRLG", "RLC SRLS", "SDMA", "SDMA2",
	"SMC", "SOS", "TA RAS", "TA XGMI", "UVD", "VCE", "VCN",
}

var utilizationCounterName = []string{"GFX Activity", "Memory Activity"}

var MemoryPageStatus = map[RSMIMemoryPageStatus]string{
	RSMI_MEM_PAGE_STATUS_RESERVED:     "reserved",
	RSMI_MEM_PAGE_STATUS_PENDING:      "pending",
	RSMI_MEM_PAGE_STATUS_UNRESERVABLE: "unreservable",
}

// 定义常量表示链接类型
const (
	LinkTypePCIE    = "PCIE"
	LinkTypeHSL     = "HSL"
	LinkTypeUnknown = "XXXX"
)

// DeviceControlInfo 控制设备信息
type DeviceControlInfo struct {
	// DvInd 设备索引号
	DvInd int
	// PerfLevel 性能水平
	PerfLevel string
	// SclkClock sclk时钟频率 600、700、750、800、900、1000、1106、1200、1270、1319、1400、1500、1600
	SclkClock string
	// SocclkClock soclk时钟频率 309、523、566、618、680、755、850、971
	SocclkClock string
	// ResetFan 是否重置风扇控制
	ResetFan bool
}

const (
	Freq600Mhz  = 1    // 0b000000000001
	Freq700Mhz  = 2    // 0b000000000010
	Freq750Mhz  = 4    // 0b000000000100
	Freq800Mhz  = 8    // 0b000000001000
	Freq900Mhz  = 16   // 0b000000010000
	Freq1000Mhz = 32   // 0b000000100000
	Freq1106Mhz = 64   // 0b000001000000
	Freq1200Mhz = 128  // 0b000010000000
	Freq1270Mhz = 256  // 0b000100000000
	Freq1319Mhz = 512  // 0b001000000000
	Freq1400Mhz = 1024 // 0b010000000000
	Freq1600Mhz = 2048 // 0b100000000000
)

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

// 定义设备型号的“枚举”类型
type DeviceModel int

const (
	K100_AI DeviceModel = iota
	K100_AI_Liquid
	K100_AI_Eco
	K100
	Z100
	Z100L
)

// DiagResult 单项诊断测试结果
type DiagResult struct {
	// Status 诊断结果状态
	Status string `json:"status" example:"PASS"`

	// TestName 测试名称
	TestName string `json:"testName" example:"PCIe Bandwidth Test"`

	// TestOutput 测试输出信息
	TestOutput string `json:"testOutput" example:"Bandwidth is within expected range"`

	// ErrorCode 错误码
	ErrorCode int `json:"errorCode" example:"0"`

	// ErrorMessage 错误信息
	ErrorMessage string `json:"errorMessage" example:""`
}

// HCUResult 单个 HCU 的诊断结果
type HCUResult struct {
	// HCU HCU 编号
	HCU int `json:"hcu" example:"0"`

	// RC 返回码
	RC int `json:"rc" example:"0"`

	// DiagResults HCU 下的诊断结果列表
	DiagResults []DiagResult `json:"diagResults"`
}

// DiagResults 总体诊断结果
type DiagResults struct {
	// DeviceNumber 设备数量信息
	DeviceNumber int `json:"deviceNumber" example:"8"`

	// Software 软件层诊断结果
	Software []DiagResult `json:"software"`

	// PerHCU 每个 HCU 的诊断结果
	PerHCU []HCUResult `json:"perHCU"`
}

// HealthCheckConfig 健康检查配置返回结构
type HealthCheckConfig struct {
	// Enabled 是否开启健康检查
	Enabled bool `json:"enabled" example:"true"`

	// Options 健康检查项列表
	Options []string `json:"options" example:"power,temperature,pcie"`
}

type CreateGroupRequest struct {
	GroupName string `json:"groupName"`
}

type AddDcuToGroupRequest struct {
	DvInd int `json:"dvInd"`
}

type EntityListRequest struct {
	EntityList []dcgm.GroupEntityPair `json:"entityList"`
}

type DevIdListRequest struct {
	DevIdList []int `json:"hcuList"`
}

type CreateFieldGroupRequest struct {
	FieldGroupName string `json:"fieldGroupName"`
	FieldIds       []int  `json:"fieldIds"`
}

type FieldIdListRequest struct {
	FieldIdList []int `json:"fieldIdList"`
}

type WatchFieldsWithGroupRequest struct {
	FieldIdList    []int   `json:"fieldIdList"`
	UpdateFreq     float64 `json:"updateFreq"`
	MaxKeepAge     float64 `json:"maxKeepAge"`
	MaxKeepSamples int     `json:"maxKeepSamples"`
}

type EntityListWithFieldValuesResp struct {
	EntityGroupId  dcgm.Field_Entity_Group `json:"entityGroupId"`
	EntityId       int                     `json:"entityId"`
	FieldValueList interface{}             `json:"fieldValueList"`
}

// DevTypeNameResp 设备类型名称返回
type DevTypeNameResp struct {
	// DevTypeName 设备类型名称
	DevTypeName string `json:"devTypeName"`
	// Unit 单位
	Unit float64 `json:"unit"`
}

// DevSubsystemIdResp 设备子系统 ID 返回
type DevSubsystemIdResp struct {
	// SubsystemId 子系统 ID
	SubsystemId string `json:"subsystemId"`
}

// DevSubsystemNameResp 设备子系统名称返回
type DevSubsystemNameResp struct {
	// SubsystemName 子系统名称
	SubsystemName string `json:"subsystemName"`
}

// UMCBandwidthReq UMC 带宽查询请求参数
type UMCBandwidthReq struct {
	// DvInd 物理设备索引
	DvInd int `json:"dvInd" example:"0"`

	// ChanId UMC 通道 ID（0 ~ 31）
	ChanId int `json:"chanId" example:"0"`

	// Delay 采样延迟（单位：秒）
	Delay int `json:"delay" example:"1"`
}

// UMCBandwidthResp UMC 带宽接口返回结构
type UMCBandwidthResp struct {
	// UMCBandwidth UMC 带宽信息
	UMCBandwidth UMCBandwidthInfo `json:"umcBandwidth"`
}

type UMCBandwidthInfo struct {
	// ReadBW 各通道读带宽
	ReadBW []float64 `json:"readBW"`

	// WriteBW 各通道写带宽
	WriteBW []float64 `json:"writeBW"`

	// ReadWriteBW 各通道读写带宽
	ReadWriteBW []float64 `json:"readWriteBW"`
}

// XHCLBandwidthReq XHCL 带宽查询请求参数
type XHCLBandwidthReq struct {
	// DvInd 物理设备索引
	DvInd int

	// LinkId XHCL 链路 ID
	LinkId int

	// Direction 带宽方向
	Direction int

	// Delay 采样延迟（秒）
	Delay int
}

// XHCLBandwidthResp XHCL 带宽响应
type XHCLBandwidthResp struct {
	// XhclBandwidth XHCL 带宽信息
	XhclBandwidth XhclBandwidthInfo
}

const MAX_XHCL_LINK_NUM = 7

// XhclBandwidthInfo XHCL 带宽信息
type XhclBandwidthInfo struct {
	// Bw 每条 XHCL 链路的带宽值
	// 数组下标对应链路 ID，长度为 MAX_XHCL_LINK_NUM
	Bw []float64 `json:"bw"`
}

// HyLinkStatusResp HyLink Link 状态返回数据
type HyLinkStatusResp struct {
	// Devices 各设备的 Link 带宽聚合结果
	Devices []DeviceLinkSum `json:"devices"`
}

// DeviceLinkSum 表示单个设备（DvInd）的 Link 带宽汇总结果
type DeviceLinkSum struct {
	// DvInd 设备索引
	DvInd int `json:"dvInd"`

	// Recv 接收方向（direction=0）所有 Link 的带宽总和
	Recv float64 `json:"recv"`

	// Send 发送方向（direction=1）所有 Link 的带宽总和
	Send float64 `json:"send"`

	// Err 该设备查询过程中的错误信息（成功时为空字符串）
	Err string `json:"err"`

	// Links 每个 Link 的接收/发送带宽明细
	Links []LinkBandwidth `json:"links"`
}

// LinkBandwidth 表示单个 Link 的接收/发送带宽信息
type LinkBandwidth struct {
	// LinkId Link 索引
	LinkId int `json:"linkId"`

	// Recv 接收方向（direction=0）的带宽
	Recv float64 `json:"recv"`

	// Send 发送方向（direction=1）的带宽
	Send float64 `json:"send"`
}

// HyLinkStatusByDcuIdResp  HCU 设备的 HyLink 带宽
type HyLinkStatusByDcuIdResp struct {
	// DeviceLinkBandwidth 设备的 HyLink 带宽汇总结果
	DeviceLinkBandwidth DeviceLinkSum `json:"deviceLinkBandwidth"`
}

// HyUMCStatusResp 所有设备的 UMC 带宽
type HyUMCStatusResp struct {
	// DeviceUmcBandwidth 所有设备的 UMC 带宽汇总信息列表
	DeviceUmcBandwidth []DeviceUmcSum `json:"deviceUmcBandwidth"`
}

// DeviceUmcSum 表示单个设备（device index）的 UMC 带宽聚合结果
type DeviceUmcSum struct {
	// DvInd 设备索引（从 0 开始）
	DvInd int `json:"dvInd"`

	// Read 所有 UMC channel 的读带宽之和（失败时为 0）
	Read float64 `json:"read"`

	// Write 所有 UMC channel 的写带宽之和（失败时为 0）
	Write float64 `json:"write"`

	// ReadWrite 所有 UMC channel 的读写带宽之和（失败时为 0）
	ReadWrite float64 `json:"readWrite"`

	// Err 错误信息；成功时为空字符串
	Err string `json:"err"`
}

// ProcessInfoResp 进程信息接口返回结构
type ProcessInfoResp struct {
	// ProcessInfo 指定进程的计算资源使用信息
	ProcessInfo ProcessInfos `json:"processInfo"`
}

// ProcessInfos 进程的计算资源使用信息
type ProcessInfos struct {
	// ProcessID 进程 ID
	ProcessID uint32 `json:"processID"`

	// Pasid 进程地址空间 ID
	Pasid uint32 `json:"pasid"`

	// VramUsage 显存使用量
	VramUsage uint64 `json:"vramUsage"`

	// SdmaUsage SDMA 使用量
	SdmaUsage uint64 `json:"sdmaUsage"`

	// CuOccupancy CU 占用率
	CuOccupancy uint32 `json:"cuOccupancy"`
}

// VDeviceInfosResp 虚拟设备信息列表返回结构
type VDeviceInfosResp struct {
	// VDeviceInfos 虚拟设备信息列表
	VDeviceInfos []VDeviceInfo `json:"vDeviceInfos"`
}

// VDeviceInfo 虚拟设备信息
type VDeviceInfo struct {
	// Name 虚拟设备的名称
	Name string `json:"name"`

	// SubsystemTypeName 设备子系统名称
	SubsystemTypeName string `json:"subsystemTypeName"`

	// VComputeUnitCount 虚拟设备的计算单元数量
	VComputeUnitCount int `json:"vComputeUnitCount"`

	// VMemoryTotal 虚拟设备的全局内存大小
	VMemoryTotal uint64 `json:"vMemoryTotal"`

	// VMemoryUsed 虚拟设备的已使用内存大小
	VMemoryUsed uint64 `json:"vMemoryUsed"`

	// ContainerID 虚拟设备的容器ID
	ContainerID uint64 `json:"containerID"`

	// DvInd 物理设备的设备ID
	DvInd int `json:"dvInd"`

	// VPercent 虚拟设备的使用百分比
	VPercent int `json:"vPercent"`

	// VdvInd 虚拟设备的索引号
	VdvInd int `json:"vdvInd"`

	// PciBusNumber 虚拟设备的总线编号
	PciBusNumber string `json:"pciBusNumber"`
}

// DiagResp 诊断任务返回结构
type DiagResp struct {
	// JobID 诊断任务 ID
	JobID string `json:"jobId"`

	// Status 当前任务状态（pending / running）
	Status string `json:"status"`

	// StatusURL 查询任务状态的接口地址
	StatusURL string `json:"statusUrl"`
}

// HealthCheckByGroupReq 获取组设备健康状态请求参数
type HealthCheckByGroupReq struct {
	// GroupId 设备组ID
	GroupId int `json:"groupId" binding:"required"`
	// CheckHealthConfig 是否检查 HealthConfig
	CheckHealthConfig bool `json:"checkHealthConfig"`
}

// HealthCheckByGroupResp 获取组设备健康状态响应
type HealthCheckByGroupResp struct {
	// DeviceHealths 设备健康结果列表
	DeviceHealths []DeviceHealth `json:"deviceHealths"`
}

// SystemWatch 检查详情
type SystemWatch struct {
	// Type 检查项的类型
	Type string `json:"type"`

	// Status 状态
	Status string `json:"status"`

	// Error 错误的详细信息
	Error string `json:"error"`

	// Result 结果
	Result interface{} `json:"result"`
}

// DeviceHealth 设备健康结果
type DeviceHealth struct {
	// HCU HCU编号
	HCU uint `json:"hcu"`

	// Status 状态
	Status string `json:"status"`

	// Watches 检查详情
	Watches []SystemWatch `json:"watches"`
}

// DevNameResp 设备名称返回结构体
type DevNameResp struct {
	DeviceName string `json:"deviceName"`
}

// NumMonitorDevicesResp GPU 数量返回结构体
type NumMonitorDevicesResp struct {
	// GpuCount GPU 数量
	GpuCount int `json:"gpuCount"`
}

// DevSkuResp 返回设备 SKU 结构体
type DevSkuResp struct {
	// Sku 设备 SKU
	Sku int `json:"sku"`
}

// DevBrandResp 返回设备品牌名称
type DevBrandResp struct {
	// Brand 设备品牌名称
	Brand string `json:"brand"`
}

// DevVendorNameResp 返回设备供应商名称
type DevVendorNameResp struct {
	// BName 设备供应商名称
	BName string `json:"bname"`
}

// DevVramVendorResp 返回显存供应商名称
type DevVramVendorResp struct {
	// VendorName 显存供应商名称
	VendorName string `json:"vendorName"`
}

// DevPciBandwidthResp PCIe 带宽响应结构体
type DevPciBandwidthResp struct {
	// RsmiPcieBandwidth PCIe 带宽信息
	PcieBandwidth PcieBandwidth `json:"PcieBandwidth"`
}

// MemoryPercentResp 内存使用百分比响应
type MemoryPercentResp struct {
	// BusyPercent 内存使用百分比
	BusyPercent int `json:"busyPercent"`
}

// PicBusInfoResp 总线信息响应
type PicBusInfoResp struct {
	// BusInfo BDF 格式的总线信息
	BusInfo string `json:"busInfo"`
}

// FanSpeedInfoResp 风扇转速信息响应
type FanSpeedInfoResp struct {
	// FanLevel 风扇转速
	FanLevel int64 `json:"fanLevel"`

	// FanPercentage 风扇转速占最大转速的百分比
	FanPercentage float64 `json:"fanPercentage"`
}

// HCUUseResp HCU 使用率信息响应
type HCUUseResp struct {
	// GPUUsage HCU 当前使用百分比
	GPUUsage int `json:"gpuUsage"`
}

// MaxSECount 每个设备最大 SE 数量
const MaxSECount = 8

// SEUsageInfo 各 Shader Engine（SE）瞬时占用率
type SEUsageInfo struct {
	// Percent 各 SE 占用率，下标对应 SE 索引
	Percent [MaxSECount]float32 `json:"percent"`
}

// HCUCuUsageResp HCU 瞬时占用率响应（对应 hy-smi -u）
type HCUCuUsageResp struct {
	// UtilizationRate HCU 瞬时占用率
	UtilizationRate float64 `json:"utilizationRate"`
}

// DevUtilSampleResp 采样类占用率响应（HCU/CU/Wave）
type DevUtilSampleResp struct {
	// UtilizationRate 采样窗口内的占用率
	UtilizationRate float64 `json:"utilizationRate"`

	// SampleDurationMs 采样时间窗口（毫秒）
	SampleDurationMs int `json:"sampleDurationMs"`
}

// HCUSEUsageResp SE 瞬时占用率响应（对应 hy-smi --showseuse）
type HCUSEUsageResp struct {
	// ShaderEngineUsage 各 SE 占用率
	ShaderEngineUsage SEUsageInfo `json:"shaderEngineUsage"`
}

// DevTypeIDResp 设备ID十六进制值响应
type DevTypeIDResp struct {
	// ID 设备ID的十六进制值
	ID string `json:"id"`
}

// MaxPowerResp 设备最大功率响应
type MaxPowerResp struct {
	// Power 设备的最大功率（瓦特）
	Power int64 `json:"power"`
}

// MemInfoResp 设备内存信息响应
type MemInfoResp struct {
	// MemUsed 已使用的内存量
	MemUsed int64 `json:"memUsed"`

	// MemTotal 内存总量
	MemTotal int64 `json:"memTotal"`
}

// DFBandwidthResp DF 带宽响应
type DFBandwidthResp struct {
	// DFBandwidthInfo DF 带宽信息
	DFBandwidthInfo DFBandwidthInfo `json:"dfBandwidthInfo"`
}

// DFBandwidthInfo DF 带宽信息
type DFBandwidthInfo struct {
	// ReadBW DF 读带宽
	ReadBW float64 `json:"readBW"`

	// WriteBW DF 写带宽
	WriteBW float64 `json:"writeBW"`

	// ReadWriteBW DF 读写带宽
	ReadWriteBW float64 `json:"readWriteBW"`
}

// GetProcessNameResp 返回进程名称
type GetProcessNameResp struct {
	// ProcessName 进程名称
	ProcessName string `json:"processName"`
}

// PerfLevelResp 设备性能等级响应
type PerfLevelResp struct {
	// PerfLevel 当前性能等级
	PerfLevel string `json:"perfLevel"`
}

// PowerResp 设备平均功耗响应
type PowerResp struct {
	// Power 平均功耗，单位瓦特
	Power int64 `json:"power"`
}

// EccStatusResp GPU块ECC状态响应
type EccStatusResp struct {
	// EccStatus ECC状态，字符串类型
	EccStatus string `json:"eccStatus"`
}

// TemperatureResp 设备温度响应
type TemperatureResp struct {
	// Temp 当前温度（摄氏度）
	Temp float64 `json:"temp"`
}

// RsmiVersionResp RSMI版本响应
type RsmiVersionResp struct {
	// RsmiVersion 当前RSMI版本
	RsmiVersion DevVersion `json:"rsmiVersion"`
}

// DevVersion RSMI版本信息
type DevVersion struct {
	// Major 主版本号
	Major uint32 `json:"major"`
	// Minor 次版本号
	Minor uint32 `json:"minor"`
	// Patch 补丁版本号
	Patch uint32 `json:"patch"`
	// Build 构建字符串，包含构建的额外信息
	Build string `json:"build"`
}

// VbiosVersionResp VBIOS版本响应
type VbiosVersionResp struct {
	// Vbios VBIOS版本号
	Vbios string `json:"vbios"`
}

// DevGpuClkFreqSetResp GPU 时钟频率设置响应
type DevGpuClkFreqSetResp struct {
	// Message 提示信息
	Message string `json:"message"`
}

// VersionResp 驱动程序版本响应
type VersionResp struct {
	// Version 驱动程序版本
	Version string `json:"version"`
}

// FailedResp 重置时钟响应
type FailedResp struct {
	// FailedMessages 包含重置clock操作失败时的设备ID和错误信息
	FailedMessages []FailedMessage `json:"failedMessages"`
}

// FailedMessage 重置clock错误信息
// @Description 包含重置clock操作失败时的设备ID和错误信息
type FailedMessage struct {
	// ID 设备ID
	ID int `json:"id"`
	// ErrorMsg 错误信息
	ErrorMsg string `json:"errorMsg"`
}

// HSLErrorStatusResp HSL错误状态响应
type HSLErrorStatusResp struct {
	// Status HSL状态码
	Status int `json:"status"`
}

// HSLHiveIdResp HSL Hive ID 响应
type HSLHiveIdResp struct {
	// HiveId 设备的HSL hive id
	HiveId int64 `json:"hiveId"`
}

// ResetPerfDeterminismResp 重置性能决定性响应
type ResetPerfDeterminismResp struct {
	// FailedMessages 包含操作失败的设备及其错误信息
	FailedMessages []FailedMessage `json:"failedMessages"`
}

// SetClockRangeResp 设置时钟范围响应
type SetClockRangeResp struct {
	// FailedMessages 包含操作失败的设备及其错误信息
	FailedMessages []FailedMessage `json:"failedMessages"`
}

// PowerPlayResp PowerPlay表操作响应
type PowerPlayResp struct {
	// FailedMessages 包含操作失败的设备及其错误信息
	FailedMessages []FailedMessage `json:"failedMessages"`
}

// ClockOverDriveResp 时钟OverDrive操作响应
type ClockOverDriveResp struct {
	// FailedMessages 包含操作失败的设备及其错误信息
	FailedMessages []FailedMessage `json:"failedMessages"`
}

// FanSpeedResp 风扇速度响应
type FanSpeedResp struct {
	// Speed 风扇速度，单位 RPM
	Speed int64 `json:"speed"`
}

// BusInfoResp 设备总线信息响应
type BusInfoResp struct {
	// PicID 设备总线ID
	PicID string `json:"picId"`
}

// TemperatureInfo 表示一个设备的温度信息
type TemperatureInfo struct {
	// DeviceID 设备索引号
	DeviceID int `json:"deviceID"`

	// SensorTemps 传感器名称到温度的映射
	// key 为传感器名称（如 "Edge", "Junction", "VRAM"）
	// value 为对应温度（摄氏度）
	SensorTemps map[string]float64 `json:"sensorTemps"`
}

// ShowCurrentTempsResp 封装接口返回体
type ShowCurrentTempsResp struct {
	// TemperatureInfos 温度信息列表
	TemperatureInfos []TemperatureInfo `json:"temperatureInfos"`
}

// FirmwareBlock 表示单个固件块的名称和版本
type FirmwareBlock struct {
	// BlockName 固件块名称
	BlockName string `json:"blockName"`

	// Version 固件版本
	Version string `json:"version"`
}

// FirmwareInfoResp 表示设备的固件信息
type FirmwareInfoResp struct {
	// DeviceID 设备索引号
	DeviceID int `json:"deviceID"`

	// FirmwareVer 固件块列表
	FirmwareVer []FirmwareBlock `json:"firmwareVer"`
}

// ShowFwInfoResp 封装接口返回体
type ShowFwInfoResp struct {
	// FwInfos 设备固件信息列表
	FwInfos []FirmwareInfoResp `json:"fwInfos"`
}

// PidListResp 表示计算进程列表返回体
type PidListResp struct {
	// PidList 计算进程 ID 列表
	PidList []string `json:"pidList"`
}

// ProcessHCUInfoResp 表示进程及其使用的 HCU 设备信息
type ProcessHCUInfoResp struct {
	// ProcessInfo 进程信息列表
	ProcessInfo []Process `json:"processInfo"`
}

// GetCoarseGrainUtilReq 获取粗粒度利用率请求体
type GetCoarseGrainUtilReq struct {
	// Device 设备 ID
	Device int `json:"device"`

	// TypeName 利用率计数器类型名称（可选）
	TypeName string `json:"typeName,omitempty"`
}

// GetCoarseGrainUtilResp 获取粗粒度利用率响应体
type GetCoarseGrainUtilResp struct {
	// UtilizationCounters 利用率计数器列表
	UtilizationCounters []UtilizationCounter `json:"utilizationCounters"`
}

// UtilizationCounterType 利用率计数器类型
type UtilizationCounterType uint32

// UtilizationCounter 利用率计数器信息
type UtilizationCounter struct {
	// Type 利用率计数器类型
	Type UtilizationCounterType `json:"type"`

	// Value 利用率值
	Value uint64 `json:"value"`
}

// ShowHCUUseResp HCU 使用率响应体
type ShowHCUUseResp struct {
	// DeviceUseInfos 设备使用信息列表
	DeviceUseInfos []DeviceUseInfo `json:"deviceUseInfos"`
}

// DeviceUseInfo 设备使用信息
type DeviceUseInfo struct {
	// DeviceID 设备索引号
	DeviceID int `json:"deviceId"`

	// GPUUsePercent GPU 使用率（百分比）
	GPUUsePercent int `json:"gpuUsePercent"`

	// Utilization 各类利用率统计
	Utilization map[string]uint64 `json:"utilization"`
}

// ShowMemVendorResp 内存供应商信息响应体
type ShowMemVendorResp struct {
	// DeviceMemVendorInfos 设备内存供应商信息列表
	DeviceMemVendorInfos []DeviceMemVendorInfo `json:"deviceMemVendorInfos"`
}

// DeviceMemVendorInfo 设备内存供应商信息
type DeviceMemVendorInfo struct {
	// DeviceID 设备索引号
	DeviceID int `json:"deviceId"`

	// Vendor 内存供应商信息
	Vendor string `json:"vendor"`
}

// ShowPcieBwResp PCIe 带宽响应体
type ShowPcieBwResp struct {
	// PcieBandwidthInfos PCIe 带宽信息列表
	PcieBandwidthInfos []PcieBandwidthInfo `json:"pcieBandwidthInfos"`
}

// PcieBandwidthInfo 设备 PCIe 带宽信息
type PcieBandwidthInfo struct {
	// DeviceID 设备索引号
	DvInd int `json:"dvInd"`

	// Sent PCIe 发送带宽
	Sent float64 `json:"sent"`

	// Received PCIe 接收带宽
	Received float64 `json:"received"`

	// Bw 总带宽
	Bw float64 `json:"bw"`
}

// PcieReplayCountInfo 设备的 PCIe 重放计数信息
type PcieReplayCountInfo struct {
	// DeviceID 设备索引号
	DeviceID int `json:"deviceId"`
	// Count PCIe 重放总数
	Count int64 `json:"count"`
}

// ShowPcieReplayCountResponse PCIe 重放计数返回体
type ShowPcieReplayCountResponse struct {
	// PcieReplayCountInfos 各设备 PCIe 重放计数信息
	PcieReplayCountInfos []PcieReplayCountInfo `json:"pcieReplayCountInfos"`
}

// DevicePowerInfo 设备的功率信息
type DevicePowerInfo struct {
	// DeviceID 设备索引号
	DeviceID int `json:"deviceId"`

	// Power 设备功率
	Power int64 `json:"power"`
}

// DevicePowerPlayInfo 设备的GPU时钟频率和电压信息
// @Description 设备的GPU时钟频率和电压信息
type DevicePowerPlayInfo struct {
	// DeviceID 设备索引号
	DeviceID int `json:"deviceId"`

	// SCLK SCLK 时钟列表
	SCLK []string `json:"sclk"`

	// MCLK MCLK 时钟值
	MCLK string `json:"mclk"`

	// DDC_CURVE 电压曲线
	DDC_CURVE []string `json:"ddcCurve"`

	// RANGE 时钟范围
	RANGE []string `json:"range"`
}

// DeviceProductInfo 设备的产品信息列表
// @Description 设备的产品信息，包括产品系列、型号、供应商和 SKU
type DeviceProductInfo struct {
	// DeviceID 设备索引号
	// example: 0
	DeviceID int `json:"deviceId"`

	// CardSeries 设备系列名称
	// example: "Hygon HCU"
	CardSeries string `json:"cardSeries"`

	// CardModel 设备型号
	// example: "HDCU-1000"
	CardModel string `json:"cardModel"`

	// CardVendor 设备供应商
	// example: "Hygon"
	CardVendor string `json:"cardVendor"`

	// CardSKU SKU 信息
	// example: "SKU-001"
	CardSKU string `json:"cardSku"`
}

// DeviceProfile 设备的电源配置文件信息
// @Description 设备的电源配置文件信息，包括可用的电源配置选项
type DeviceProfile struct {
	// DeviceID 设备索引号
	// example: 0
	DeviceID int `json:"deviceId"`

	// Profiles 文件信息列表
	// example: ["Profile1", "Profile2"]
	Profiles []string `json:"profiles"`
}

// DeviceSerialInfo 设备的序列号信息
// @Description 设备的序列号信息，包括设备索引号和序列号
type DeviceSerialInfo struct {
	// DeviceID 设备索引号
	// example: 0
	DeviceID int `json:"deviceId"`

	// SerialNumber 设备序列号
	// example: "SN123456789"
	SerialNumber string `json:"serialNumber"`
}

// DeviceUIdInfo 设备的唯一ID信息
// @Description 设备的唯一ID信息，包括设备索引号和唯一ID
type DeviceUIdInfo struct {
	// DeviceID 设备索引号
	// example: 0
	DeviceID int `json:"deviceId"`

	// UId 设备唯一ID
	// example: "UID123456789"
	UId string `json:"uid"`
}

// DeviceVBIOSInfo 设备的VBIOS版本信息
// @Description 设备的VBIOS版本信息，包括设备索引号和VBIOS版本
type DeviceVBIOSInfo struct {
	// DeviceID 设备索引号
	// example: 0
	DeviceID int `json:"deviceId"`

	// VBIOS 版本信息
	// example: "113-C56801-102"
	VBIOS string `json:"vbios"`
}

// DeviceVoltageInfo 设备的电压信息
// @Description 设备的电压信息，包括设备索引号和当前电压
// @example
type DeviceVoltageInfo struct {
	// DeviceID 设备索引号
	// example: 0
	DeviceID int `json:"deviceId"`

	// Voltage 电压信息
	// example: 120
	Voltage int64 `json:"voltage"`
}

// NumaInfo 设备的 NUMA 信息
// @Description 设备的 NUMA 节点和关联信息
// @example
type NumaInfo struct {
	// DeviceID 设备索引号
	// example: 0
	DeviceID int `json:"deviceId"`

	// NumaNode NUMA 节点
	// example: 1
	NumaNode int `json:"numaNode"`

	// NumaAffinity NUMA 关联信息
	// example: 255
	NumaAffinity int `json:"numaAffinity"`
}

// DeviceCountInfo 设备数量信息
type DeviceCountInfo struct {
	// Count 设备数量
	// example: 4
	Count int `json:"count"`
}

// VDeviceCountResp 虚拟设备数量返回结构
type VDeviceCountResp struct {
	// VDeviceCount 虚拟设备数量
	VDeviceCount int `json:"vDeviceCount"`
}

// DeviceRemainingInfoResp 设备剩余信息返回结构
type DeviceRemainingInfoResp struct {
	// CUs 剩余计算单元数量
	CUs uint64 `json:"cus"`
	// Memories 剩余内存大小（单位：字节）
	Memories uint64 `json:"memories"`
}

// CreateVDevicesResp 虚拟设备创建响应
type CreateVDevicesResp struct {
	// VDevIDs 创建成功的虚拟设备ID集合
	VDevIDs []int `json:"vdevIDs"`
}

// EncryptionVMStatusResp 加密虚拟机状态返回结构
type EncryptionVMStatusResp struct {
	// Status 是否处于加密状态
	Status bool `json:"status"`
}

// DeviceInfoResp 设备信息返回结构
// @Description 获取设备信息接口的返回结构
type DeviceInfoResp struct {
	// DeviceInfo 设备信息
	DeviceInfo DMIDeviceInfo `json:"deviceInfo"`
}

// DMIDeviceInfo 设备信息
// @Description 单个物理设备的详细信息
type DMIDeviceInfo struct {
	// Name 设备名称
	// @example HCU-0
	Name string `json:"name"`

	// ComputeUnitCount 计算单元总数
	// @example 64
	ComputeUnitCount int `json:"computeUnitCount"`

	// ComputeUnitRemainingCount 剩余计算单元数
	// @swagignore
	ComputeUnitRemainingCount uint64

	// MemoryRemaining 剩余内存
	// @swagignore
	MemoryRemaining uint64

	// GlobalMemSize 全局内存大小
	// @swagignore
	GlobalMemSize uint64

	// UsageMemSize 已使用内存大小
	// @swagignore
	UsageMemSize uint64

	// DeviceID 设备索引号
	// @example 0
	DeviceID int `json:"deviceID"`

	// Percent 使用率（百分比）
	// @example 75
	Percent int `json:"percent"`

	// MaxVDeviceCount 最大虚拟设备数量
	// @example 8
	MaxVDeviceCount int `json:"maxVDeviceCount"`
}

// DeviceControlResp 设备控制响应
type DeviceControlResp struct {
	// Success 是否执行成功
	Success bool `json:"success"`

	// Errors 执行失败的错误信息
	Errors []string `json:"errors,omitempty"`
}

// GetDeviceModelInfosResp 获取设备型号信息响应
type GetDeviceModelInfosResp struct {
	// Devices 设备型号信息列表
	Devices []DeviceModelInfo `json:"devices"`
}

// DeviceModelInfo 设备型号信息
type DeviceModelInfo struct {
	// Model HCU 类型
	Model string `json:"model"`

	// CUCount CU 数量
	CUCount float64 `json:"cuCount"`

	// MemorySize 内存大小（单位：MB 或 GB，视实现而定）
	MemorySize float64 `json:"memorySize"`
}

// GetProcessInfoResp 获取进程信息响应
type GetProcessInfoResp struct {
	// Processes 进程信息列表
	Processes []Process `json:"processes"`
}

// Process 进程信息
type Process struct {
	// ProcessID 进程 ID
	ProcessID uint32 `json:"processId"`

	// ProcessName 进程名称
	ProcessName string `json:"processName"`

	// Pasid 进程地址空间 ID
	Pasid uint32 `json:"pasid"`

	// VramUsage 显存使用量（单位：Bytes）
	VramUsage uint64 `json:"vramUsage"`

	// SdmaUsage SDMA 使用量
	SdmaUsage uint64 `json:"sdmaUsage"`

	// CuOccupancy CU 占用率
	CuOccupancy uint32 `json:"cuOccupancy"`

	// MinorNumbers 使用的设备索引号
	MinorNumbers []int `json:"minorNumbers"`
}

// MigInfosResp MIG 信息响应
type MigInfosResp struct {
	// MigInfos MIG 分区信息列表
	MigInfos []MigInfo `json:"migInfos"`
}

// MigInfo MIG 分区信息
type MigInfo struct {
	// DvInd 物理设备索引
	DvInd int `json:"dvInd"`

	// MigId MIG 分区 ID
	MigId int `json:"migId"`

	// Name MIG 分区名称
	Name string `json:"name"`

	// UUID MIG 分区 UUID
	UUID string `json:"uuid"`

	// ComputeUnit 计算单元数量
	ComputeUnit uint32 `json:"computeUnit"`

	// MemoryTotal 分区总显存（Bytes）
	MemoryTotal uint64 `json:"memoryTotal"`

	// GpuInstanceId GPU Instance ID
	GpuInstanceId uint32 `json:"gpuInstanceId"`

	// ComputeInstanceId Compute Instance ID
	ComputeInstanceId uint32 `json:"computeInstanceId"`

	// PciBusNumber PCI Bus 号
	PciBusNumber string `json:"pciBusNumber"`

	// GiProfileId GI Profile ID
	GiProfileId int `json:"giProfileId"`

	// CiProfileId CI Profile ID
	CiProfileId int `json:"ciProfileId"`
}

// Job 表示一次异步诊断任务（job）的元数据与执行结果
// @Description Job represents an asynchronous diagnostic task with its metadata and results
type Job struct {
	// ID 唯一作业标识符（由服务生成，例如 job-<unixnano>-<seq>）
	// example: job-1698493923-01
	ID string `json:"id"`

	// Level 要执行的诊断等级（1..4）
	// example: 2
	Level int `json:"level"`

	// Status 作业当前状态
	// example: pending
	Status string `json:"status"`

	// Result 作业完成后保存的结构化诊断结果
	Result *DiagResults `json:"result,omitempty"`

	// ErrorMessage 如果作业执行失败或中止，这里保存错误信息或说明
	// example: GPU not found
	ErrorMessage string `json:"errorMessage,omitempty"`

	// StartedAt 作业实际开始执行的时间戳（零值表示尚未开始）
	// example: 2025-12-29T10:00:00Z
	StartedAt time.Time `json:"startedAt"`

	// EndedAt 作业结束（成功/失败/取消）的时间戳（零值表示尚未结束）
	// example: 2025-12-29T10:30:00Z
	EndedAt time.Time `json:"endedAt"`
}

// HCULinkInfo 描述两张 HCU 之间的互联关系
type HCULinkInfo struct {
	// SrcDvInd 源 HCU 索引
	// example: 0
	SrcDvInd int `json:"srcDvInd"`

	// DstDvInd 目标 HCU 索引
	// example: 1
	DstDvInd int `json:"dstDvInd"`

	// PciID 目标 HCU 的 PCI ID
	// example: 0000:07:00.0
	PciID string `json:"PciID"`

	// LinkType 链路类型: PCIE / HSL / HYSWITCH / NONE
	// example: HSL
	LinkType string `json:"linkType"`

	// Weight 链路权重
	// example: 2
	Weight int `json:"weight"`

	// Hops 跳数（目前可置 0 或 1）
	// example: 1
	Hops int `json:"hops"`
}

// HCUInterconnectMatrix 描述整机 HCU 互联矩阵
type HCUInterconnectMatrix struct {
	// DeviceCount HCU 总数
	// example: 8
	DeviceCount int `json:"deviceCount"`

	// Matrix 互联矩阵: [src][dst] 对应 HCULinkInfo
	Matrix [][]HCULinkInfo `json:"matrix"`
}
