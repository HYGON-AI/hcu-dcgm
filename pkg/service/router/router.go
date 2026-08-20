/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 注册路由
	router.GET("/devicename/:dvInd", GetDevName)
	router.GET("/NumMonitorDevices", GetNumMonitorDevices)
	router.GET("/DevSku/:dvInd", GetDevSku)
	router.GET("/DevBrand/:dvInd", DevBrand)
	router.GET("/DevVendorName/:dvInd", DevVendorName)
	router.GET("/DevVramVendor/:dvInd", DevVramVendor)
	router.GET("/DevPciBandwidth/:dvInd", DevPciBandwidth)
	router.POST("/DevPciBandwidthSet", DevPciBandwidthSet)

	router.GET("/MemoryPercent/:dvInd", MemoryPercent)
	// 路由注册
	router.GET("/PerfLevel/:dvInd", PerfLevel)
	router.POST("/DevPerfLevelSet/:dvInd", DevPerfLevelSet)
	router.GET("/DevGpuMetricsInfo/:dvInd", DevGpuMetricsInfo)
	router.GET("/CollectDeviceMetrics", CollectDeviceMetrics)
	router.GET("/DeviceInfo/:dvInd", GetDeviceByDvInd)
	router.GET("/VDeviceByDvInd/:dvInd", VDeviceByDvInd)
	// 路由
	router.GET("/PicbusInfo/:dvInd", PicBusInfo)
	router.GET("/DFBandwidth/:dvInd", DFBandwidth)
	router.GET("/HCUCuUsage/:dvInd", HCUCuUsage)
	router.GET("/HCUSampledUsage/:dvInd", HCUSampledUsage)
	router.GET("/HCUCUSampledUsage/:dvInd", HCUCUSampledUsage)
	router.GET("/HCUWaveSampledUsage/:dvInd", HCUWaveSampledUsage)
	router.GET("/HCUSEUsage/:dvInd", HCUSEUsage)
	// 路由(K100 AI不支持)
	router.GET("/FanSpeedInfo/:dvInd", FanSpeedInfo)
	// 路由
	router.GET("/HCUUse/:dvInd", HCUUse)
	// 路由
	router.GET("/DevTypeID", DevTypeID)
	router.GET("/DevTypeName", DevTypeName)
	router.GET("/DevSubsystemId", DevSubsystemId)
	router.GET("/DevSubsystemName", DevSubsystemName)
	// 路由
	router.GET("/MaxPower", GetMaxPower)
	// 路由
	router.GET("/MemInfo", GetMemInfo)
	router.GET("/AllDeviceInfos", AllDeviceInfos)
	// 路由
	router.GET("/DeviceInfos", DeviceInfos)

	router.GET("/DeviceStatus", DeviceStatus)
	// 路由
	router.GET("/ProcessName", GetProcessName)
	router.GET("/ProcessInfo", GetProcessInfo)
	// 路由注册
	router.GET("/Power/:dvInd", Power)

	//设置 GPU 时钟频率
	router.POST("/DevGpuClkFreqSet", DevGpuClkFreqSet)
	// 路由注册
	router.GET("/EccStatus/:dvInd", EccStatus)
	// 路由注册
	router.GET("/Temperature/:dvInd", Temperature)
	// 路由注册
	router.GET("/VbiosVersion/:dvInd", VbiosVersion)
	// 路由注册
	router.GET("/Version", Version)
	// 重置设备时钟(K100 AI不支持)
	router.POST("/ResetClocks", ResetClocks)
	router.POST("/ResetFans", ResetFans)

	router.POST("/ResetProfile", ResetProfile)
	//(K100 AI不支持)
	router.POST("/ResetHSLErr", ResetHSLErr)
	//(K100 AI不支持)
	router.GET("/HSLErrorStatus", HSLErrorStatus)
	router.GET("/HSLHiveId", HSLHiveIdGet)
	router.POST("/ResetPerfDeterminism", ResetPerfDeterminism)
	// 路由(K100 AI不支持)
	router.POST("/SetClockRange", SetClockRange)
	// 路由（K100_AI卡不支持）
	router.POST("/SetPowerPlayTableLevel", SetPowerPlayTableLevel)
	// 路由（sudo权限)
	router.POST("/SetClockOverDrive", SetClockOverDrive)
	// 路由（K100_AI卡不支持）
	router.POST("/SetPerfDeterminism", SetPerfDeterminism)
	// 设置风扇速度(K100 AI不支持)
	router.POST("/SetFanSpeed", SetFanSpeed)
	// 获取设备风扇转速(K100 AI不支持)
	router.GET("/DevFanRpms/:dvInd", DevFanRpms)
	// 设置设备性能等级
	router.POST("/SetPerformanceLevel", SetPerformanceLevel)
	// 设置功率配置文件（K100_AI卡不支持该操作,CUSTOM不支持，剩余几个类型超出安全范围）
	router.POST("/SetProfile", SetProfile)
	// 设置设备功率配置文件（K100_AI卡不支持该操作）
	router.POST("/DevPowerProfileSet/:dvInd", DevPowerProfileSet)
	// 获取设备总线信息
	router.GET("/GetBus/:dvInd", GetBus)
	// 显示设备硬件信息
	router.POST("/ShowAllConciseHw", ShowAllConciseHw)
	// 显示设备时钟信息
	router.POST("/ShowClocks", ShowClocks)
	router.POST("/fans/current", ShowCurrentFans)
	router.POST("/temps/current", ShowCurrentTemps)
	router.GET("/firmware/info", ShowFwInfo)
	router.GET("/process/list", PidList)
	router.POST("/utilization/coarse", GetCoarseGrainUtil)
	router.POST("/HCU/use", ShowHCUUse)
	router.POST("/energy", ShowEnergy)
	router.POST("/memory/info", ShowMemInfo)
	router.POST("/memory/use", ShowMemUse)
	router.POST("/memory/vendor", ShowMemVendor)
	router.POST("/pcie/bandwidth", ShowPcieBw)
	router.POST("/pcie/replaycount", ShowPcieReplayCount)
	router.GET("/process/name/:pid", GetProcessName)
	router.GET("/processInfo", ProcessInfo)
	router.GET("/processHCUInfo", ProcessHCUInfo)
	router.POST("/device/power", GetDevicePower)
	router.GET("/device/showPids", ShowPids)
	//（K100_AI卡不支持该操作）
	router.POST("/device/powerplay", GetDevicePowerPlayTable)
	router.POST("/device/product", GetDeviceProductName)
	router.POST("/device/profile", GetDeviceProfile)

	router.POST("/device/retiredpages", GetDeviceRetiredPages)
	router.POST("/device/serialnumber", GetDeviceSerialNumber)
	router.POST("/showUId", ShowUId)
	router.POST("/showVbiosVersion", ShowVbiosVersion)
	router.POST("/showVoltage", ShowVoltage)
	//（K100_AI卡不支持该操作）
	router.POST("/showVoltageCurve", ShowVoltageCurve)
	//（K100_AI卡不支持该操作）
	router.POST("/showHSLErr", ShowHSLErr)

	router.POST("/showWeightTopology", ShowWeightTopology)
	router.POST("/showHopsTopology", ShowHopsTopology)
	router.POST("/showTypeTopology", ShowTypeTopology)
	router.POST("/showNumaTopology", ShowNumaTopology)
	router.POST("/showHwTopology", ShowHwTopology)
	router.GET("/discoverInterconnectTopology", DiscoverInterconnectTopology)
	router.GET("/deviceCount", DeviceCount)
	router.GET("/VDeviceSingleInfo/:vDvInd", VDeviceSingleInfo)
	router.GET("/vDeviceCount", VDeviceCount)
	router.GET("/deviceRemainingInfo/:dvInd", DeviceRemainingInfo)
	router.POST("/CreateVDevices", CreateVDevices)
	router.DELETE("/DestroyVDevice", DestroyVDevice)
	router.DELETE("/DestroySingleVDevice", DestroySingleVDevice)
	router.PUT("/UpdateSingleVDevice", UpdateSingleVDevice)
	// 启动指定的虚拟设备
	router.GET("/StartVDevice/:vDvInd", StartVDevice)
	// 停止指定的虚拟设备
	router.GET("/StopVDevice/:vDvInd", StopVDevice)
	// 设置虚拟机加密状态
	router.POST("/SetEncryptionVMStatus", SetEncryptionVMStatus)
	// 获取加密虚拟机状态
	router.GET("/EncryptionVMStatus", EncryptionVMStatus)
	// 打印设备的事件列表
	router.GET("/PrintEventList/:device", PrintEventList)
	router.GET("/device/info/:dvInd", GetDeviceInfo)
	router.GET("/VDeviceInfos", VDeviceInfos)
	// 路由
	router.POST("/device/control", DeviceControl)

	router.GET("/EccBlocksInfo", EccBlocksInfo)
	router.GET("/RSMIVersion", RsmiVersion)
	// 诊断路由
	router.GET("/diag/run/:level", RunDiag)
	router.POST("/diag/stop", StopDiag)

	router.POST("/diag/start/:level", startDiag)
	router.GET("/diag/status/:id", statusDiag)
	router.POST("/bandwidth/start", startBandwidth)
	router.POST("/pcie/start", startPcie)
	router.POST("/xhcl/start", startXHCL)
	router.POST("/targetstress/start", startTargetStress)
	router.POST("/memtestcl/start", startMemtestCL)
	router.POST("/edpp/start", startEdpp)

	// 健康检查路由
	router.POST("/health/checkConfig", SetHealthCheckConfig)
	router.GET("/health/checkConfig", GetHealthCheckConfig)
	router.DELETE("/health/checkConfig", DeleteHealthCheckConfig)
	router.POST("/health/healthcheck", HealthCheckById)
	router.POST("/health/healthcheck/group", HealthCheckByGroup)
	router.GET("/deviceModelInfos", GetDeviceModelInfos)
	router.GET("/compatible", Compatible)
	router.GET("/migInfos", MigInfos)
	router.POST("/bandwidth/umc", UMCBandwidth)
	router.POST("/bandwidth/xhcl", XHCLBandwidth)
	router.GET("/GetHyLinkLinkStatus", HyLinkLinkStatus)
	router.GET("/HyLinkStatusByDcuId/:dvInd", HyLinkStatusByDcuId)
	router.GET("/GetHyUMCStatus", HyUMCStatus)
	// HSL 链路通断状态（P0，对应 DCGM_FI_DEV_NVLINK_LINK_DOWN_COUNT_TOTAL）
	router.GET("/hsl/linkstate/:dvInd", HSLLinkState)
	// 降频状态（P1，对应 THERMAL/POWER/SLOWDOWN/BOARD_LIMIT_THROTTLE）
	router.GET("/throttle/flags/:dvInd", GetThrottleFlags)
	// P2-A：温度 / 显存可用量 / 时钟
	router.GET("/memory/temperature/:dvInd", GetMemoryTemperature)
	router.GET("/board/temperature/:dvInd", GetBoardTemperature)
	router.GET("/memory/available/:dvInd", GetMemoryAvailable)
	router.GET("/clock/mem/:dvInd", GetMemClock)
	// P2-B：最大额定频率
	router.GET("/clock/gfx/max/:dvInd", GetGfxClockMax)
	router.GET("/clock/mem/max/:dvInd", GetMemClockMax)
	// P2-C：静态规格 + 拓扑
	router.GET("/device/cunum/:dvInd", GetDevCuNum)
	router.GET("/device/nodeid/:dvInd", GetDevNodeId)
	router.GET("/device/memoverdrive/:dvInd", GetDevMemOverdrive)
	router.GET("/topology/p2p/:srcDvInd/:dstDvInd", GetP2PAccessible)
	router.GET("/version/hy", GetHyVersion)
	router.GET("/groups", ListAllGroups)
	router.GET("/groups/:groupId", GetGroupInfo)
	router.POST("/groups", CreateGroup)
	router.DELETE("/groups/:groupId", DestroyGroup)
	router.POST("/groups/:groupId/hcus", AddToGroup)
	router.POST("groups/:groupId/entities", AddEntityToGroup)
	router.DELETE("groups/:groupId/entities", RemoveEntityFromGroup)
	router.GET("/fieldgroups", ListAllFieldGroups)
	router.GET("/fieldgroups/:fieldGroupId", GetFieldGroupInfo)
	router.POST("/fieldgroups", CreateFieldGroup)
	router.DELETE("/fieldgroups/:fieldGroupId", DestroyFieldGroup)
	router.POST("/policies/hcus/:dvInd", SetPolicy)
	router.DELETE("/policies/hcus/:dvInd", ClearPolicy)
	router.GET("/policies/hcus/:dvInd", GetPolicy)
	router.GET("/policies", ListPolicies)
	router.GET("/policies/result", JudgePolicyConditions)
	router.POST("/policies/action/:dvInd", TakePolicyAction)
	router.GET("/fields/:fieldId", GetFieldMetaById)
	router.GET("/fields", ListFieldMeta)
	router.POST("/fields/hcus/:dvInd", WatchFields)
	router.POST("/fields/entities/:entityGroupId/:entityId", WatchFieldsWithEntity)
	router.POST("/fields/groups/:groupId", WatchFieldsWithGroup)
	router.DELETE("/fields/hcus/:dvInd", UnWatchFields)
	router.DELETE("/fields/entities/:entityGroupId/:entityId", UnWatchFieldsWithEntity)
	router.DELETE("/fields/groups/:groupId", UnWatchFieldsWithGroup)
	router.GET("/fields/hcus/:dvInd/values", GetLatestValuesForFields)
	router.GET("/fields/entities/:entityGroupId/:entityId/values", EntityGetLatestValues)
	router.GET("/fields/groups/:groupId/values", GetLatestValuesForFieldsWithGroup)
	router.GET("/metricgroups/hcus/:dvInd", GetSupportedMetricGroups)

	// 设备级 wave 统计
	router.GET("/device/cuutil/:dvInd", GetDevCuUtil)
	router.GET("/device/waveutil/:dvInd", GetDevWaveUtil)

	// 仪电近似指标：vDevice compute busy（≈MIG_ENC_UTIL/DEC_UTIL）
	router.GET("/vdevice/busypercent/:vDvInd", GetVDevBusyPercent)
	// 仪电近似指标：进程 GPU compute 占用率（≈PROC_ENC_UTIL）
	router.GET("/process/gpuusage/:pid", GetProcessGpuUsage)

	// 功率上限设置与复位（对应 hy-smi --setpoweroverdrive / --resetpoweroverdrive）
	router.POST("/powercap/:dvInd", SetPowerCap)
	router.POST("/powercap/:dvInd/reset", ResetPowerCap)
	// 性能等级设置 v1（对应 hy-smi --setperflevel，使用新版 rsmi 接口）
	router.POST("/perflevel/v1/:dvInd", SetPerfLevelV1)

	// sysfs 只读功能（第2批：HSW 固件、boost、频率历史、显存 ECC）
	router.GET("/hswfw/:dvInd", GetHswFwVersion)
	router.GET("/boost/:dvInd", GetGfxBoostInfo)
	router.GET("/fdrhistory/:dvInd", GetFdrHistory)
	router.GET("/memeccinfo/:dvInd", GetMemBankEccInfo)

	// sysfs/debugfs 写入功能（第3批）
	router.POST("/boost/:dvInd", SetGfxBoostLevel)
	router.POST("/lowpower/:dvInd", SetLowPower)
	router.POST("/ras/:dvInd", SetRasControl)

	// 复杂聚合功能（第4批）
	router.GET("/hslstatus/:dvInd", GetHSLLogStatus)
	router.GET("/linktopo", GetPhysicalLinkTopology)
	router.POST("/config/save", SaveDeviceConfig)
	router.POST("/config/load", LoadDeviceConfig)

	return router
}
