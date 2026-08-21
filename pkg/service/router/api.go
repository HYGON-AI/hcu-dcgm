/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package router

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

// GetDevName 获取指定设备的名称
// @Summary 获取设备名称
// @Description 根据设备索引 ID 获取设备名称
// @Tags Device
// @Param dvInd path int true "设备索引ID"
// @Success 200 {object} DevNameResp "成功返回设备名称"
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /devices/{dvInd}/name [get]
func GetDevName(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	name, err := dcgm.DevName(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := DevNameResp{
		DeviceName: name,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// GetNumMonitorDevices 获取 GPU 数量
// @Tags Device
// @Summary 获取 GPU 数量
// @Description 获取监视的 GPU 数量
// @Produce json
// @Success 200 {object} NumMonitorDevicesResp "GPU 数量"
// @Failure 500 {object} Response "获取 GPU 数量失败"
// @Router /NumMonitorDevices [get]
func GetNumMonitorDevices(c *gin.Context) {
	num, err := dcgm.NumMonitorDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse("获取 GPU 数量失败"))
		return
	}

	resp := NumMonitorDevicesResp{
		GpuCount: num,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// GetDevSku 根据设备索引获取设备 SKU
// @Tags Device
// @Summary 获取设备SKU
// @Description 根据设备索引获取SKU
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} DevSkuResp "返回设备SKU"
// @Failure 400 {object} Response "获取设备SKU失败"
// @Router /DevSku/{dvInd} [get]
func GetDevSku(c *gin.Context) {
	dvIndStr := c.Param("dvInd")
	dvInd, err := strconv.Atoi(dvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}

	sku, err := dcgm.DevSku(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("获取设备SKU失败"))
		return
	}

	resp := DevSkuResp{
		Sku: sku,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// DevBrand 根据设备索引获取设备品牌名称
// @Tags Device
// @Summary 获取设备品牌名称
// @Description 根据设备索引获取品牌名称
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} DevBrandResp "返回设备品牌名称"
// @Failure 400 {object} Response "请求失败"
// @Router /DevBrand/{dvInd} [get]
func DevBrand(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}

	brand, err := dcgm.DevBrand(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	resp := DevBrandResp{
		Brand: brand,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// DevVendorName 根据设备索引获取设备供应商名称
// @Tags Device
// @Summary 获取设备供应商名称
// @Description 根据设备索引获取供应商名称
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} DevVendorNameResp "返回设备供应商名称"
// @Failure 400 {object} Response "请求失败"
// @Router /DevVendorName/{dvInd} [get]
func DevVendorName(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}

	bname, err := dcgm.DevVendorName(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	resp := DevVendorNameResp{
		BName: bname,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// DevVramVendor 根据设备索引获取显存供应商名称
// @Tags Device
// @Summary 获取设备显存供应商名称
// @Description 根据设备索引获取显存供应商名称
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} DevVramVendorResp "返回显存供应商名称"
// @Failure 400 {object} Response "请求失败"
// @Router /DevVramVendor/{dvInd} [get]
func DevVramVendor(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}

	name, err := dcgm.DevVramVendor(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	resp := DevVramVendorResp{
		VendorName: name,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// DevPciBandwidth 获取设备的可用 PCIe 带宽列表
// @Tags PCIe
// @Summary 获取可用的 PCIe 带宽列表
// @Description 根据设备 ID 获取设备的可用 PCIe 带宽列表
// @Produce json
// @Param dvInd path int true "设备 ID"
// @Success 200 {object} DevPciBandwidthResp "PCIe 带宽列表"
// @Failure 400 {object} Response "请求错误"
// @Failure 404 {object} Response "设备未找到"
// @Router /DevPciBandwidth/{dvInd} [get]
func DevPciBandwidth(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备 ID 无效"))
		return
	}

	// 从 dcgm 获取原始数据
	raw, err := dcgm.DevPciBandwidth(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	// 转换为 router 包中的类型
	resp := DevPciBandwidthResp{
		PcieBandwidth: PcieBandwidth{
			TransferRate: Frequencies{
				NumSupported: raw.TransferRate.NumSupported,
				Current:      raw.TransferRate.Current,
				Frequency:    raw.TransferRate.Frequency,
			},
			Lanes: raw.Lanes,
		},
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags PCIe
// DevPciBandwidthSet 设置设备允许的 PCIe 带宽
// @Summary 设置设备允许的 PCIe 带宽
// @Description 根据设备索引和带宽掩码限制设备允许的 PCIe 带宽
// @Produce json
// @Param dvInd query int true "设备索引"
// @Param bwBitmask query int64 true "带宽掩码"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "服务器内部错误"
// @Router /DevPciBandwidthSet [post]
func DevPciBandwidthSet(c *gin.Context) {
	var dvInd int
	var bwBitmask int64

	// 获取 query 中的 dvInd 和 bwBitmask 参数
	if err := c.ShouldBindQuery(&dvInd); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid dvInd parameter")
		return
	}
	if err := c.ShouldBindQuery(&bwBitmask); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid bwBitmask parameter")
		return
	}

	// 调用已有的 DevPciBandwidthSet 函数
	if err := dcgm.DevPciBandwidthSet(dvInd, bwBitmask); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Memory
// @Summary 获取内存使用百分比
// @Description 根据设备 ID 获取设备内存的使用百分比。
// @Param dvInd path int true "设备 ID"
// @Success 200 {object} MemoryPercentResp "内存使用百分比"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /MemoryPercent/{dvInd} [get]
func MemoryPercent(c *gin.Context) {
	dvInd, _ := strconv.Atoi(c.Param("dvInd"))

	busyPercent, err := dcgm.MemoryPercent(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	resp := MemoryPercentResp{
		BusyPercent: busyPercent,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Power
// 设置设备 PowerPlay 性能级别
// @Summary 设置设备 PowerPlay 性能级别
// @Description 根据设备 ID 设置 PowerPlay 性能级别。
// @Param dvInd path int true "设备 ID"
// @Param level query string true "要设置的性能级别"
// @Success 200 {string} string "操作成功"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevPerfLevelSet/{dvInd} [post]
func DevPerfLevelSet(c *gin.Context) {
	dvInd, _ := strconv.Atoi(c.Param("dvInd"))
	level := c.Query("level")

	// 将 level 字符串转换为 RSMIDevPerfLevel 类型
	levelConverted, err := ConvertToRSMIDevPerfLevel(level)
	if err != nil {
		// 如果转换失败，返回错误响应
		c.JSON(http.StatusBadRequest, ErrorResponse("无效的性能级别"))
		return
	}

	// 调用 dcgm.DevPerfLevelSet 并传入转换后的 level
	err = dcgm.DevPerfLevelSet(dvInd, levelConverted)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Device
// 获取 GPU 度量信息
// @Summary 获取 GPU 度量信息
// @Description 根据设备 ID 获取 GPU 的度量信息。
// @Param dvInd path int true "设备 ID"
// @Success 200 {object} RSMIGPUMetrics "GPU 度量信息"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevGpuMetricsInfo/{dvInd} [get]
func DevGpuMetricsInfo(c *gin.Context) {
	dvInd, _ := strconv.Atoi(c.Param("dvInd"))
	gpuMetrics, err := dcgm.DevGpuMetricsInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gpuMetrics)
}

// @Tags Device
// 获取设备监控中的指标
// @Summary 获取设备监控中的指标
// @Description 收集所有设备的监控指标信息。
// @Success 200 {array} MonitorInfo "设备监控指标信息列表"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /CollectDeviceMetrics [get]
func CollectDeviceMetrics(c *gin.Context) {
	monitorInfos, err := dcgm.CollectDeviceMetrics()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求失败"))
		return
	}
	c.JSON(http.StatusOK, monitorInfos)
}

// @Tags Device
// @Summary 获取设备信息
// @Description 根据设备 ID 获取物理设备的详细信息
// @Accept  json
// @Produce  json
// @Param   dvInd     path   int     true  "Device ID"
// @Success 200 {object} PhysicalDeviceInfo "设备信息"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /deviceinfo/{dvInd} [get]
func GetDeviceByDvInd(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	deviceInfo, err := dcgm.GetDeviceByDvInd(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, deviceInfo)
}

// VDeviceByDvInd 获取指定物理设备下的虚拟设备索引
// @Tags Device
// @Summary 获取物理设备的虚拟设备索引
// @Description 根据物理设备索引 dvInd，返回该物理设备下的虚拟设备数量及虚拟设备索引列表
// @Accept  json
// @Produce  json
// @Param   dvInd   path   int  true  "物理设备索引（Device Index）"
// @Success 200 {object} VDeviceByDvIndResp "虚拟设备索引信息"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /device/{dvInd}/vdevices [get]
func VDeviceByDvInd(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	vDeviceCount, vDevInds, err := dcgm.VDeviceByDvInd(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := VDeviceByDvIndResp{
		VDeviceCount: vDeviceCount,
		VDevInds:     vDevInds,
	}
	c.JSON(http.StatusOK, resp)
}

// @Tags Device
// @Summary 获取所有物理设备信息
// @Description 该接口返回所有物理设备的详细信息。
// @Produce json
// @Success 200 {array} PhysicalDeviceInfo "所有设备的详细信息"
// @Failure 400 {object} error "无效的请求参数"
// @Failure 500 {object} error "服务器内部错误"
// @Router /AllDeviceInfos [get]
func AllDeviceInfos(c *gin.Context) {
	deviceInfos, err := dcgm.AllDeviceInfos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, deviceInfos)
}

// @Tags PCIe
// @Summary 获取总线信息
// @Description 获取设备的总线信息 (BDF格式)
// @Accept  json
// @Produce  json
// @Param   dvInd     path   int     true  "Device ID"
// @Success 200 {object} PicBusInfoResp "总线信息"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /picbusinfo/{dvInd} [get]
func PicBusInfo(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	busInfo, err := dcgm.PicBusInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := PicBusInfoResp{
		BusInfo: busInfo,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Physical State
// @Summary 获取风扇转速
// @Description 获取指定设备的风扇转速及其占最大转速的百分比
// @Accept  json
// @Produce  json
// @Param   dvInd     path   int     true  "Device ID"
// @Success 200 {object} FanSpeedInfoResp "风扇转速信息"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /fanspeedinfo/{dvInd} [get]
func FanSpeedInfo(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	fanLevel, fanPercentage, err := dcgm.FanSpeedInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := FanSpeedInfoResp{
		FanLevel:      fanLevel,
		FanPercentage: fanPercentage,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Physical State
// @Summary 获取HCU使用率
// @Description 获取指定设备的HCU当前使用百分比
// @Accept  json
// @Produce  json
// @Param   dvInd     path   int     true  "Device ID"
// @Success 200 {object} HCUUseResp "HCU 使用率信息"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /gpuuse/{dvInd} [get]
func HCUUse(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	gpuUsage, err := dcgm.HCUUse(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := HCUUseResp{
		GPUUsage: gpuUsage,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Device
// @Summary 获取设备ID的十六进制值
// @Description 根据设备索引返回设备ID的十六进制值
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {object} DevTypeIDResp "设备ID十六进制值"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DevID [get]
func DevTypeID(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	id, err := dcgm.DevTypeID(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := DevTypeIDResp{
		ID: id,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// DevTypeName 获取设备类型名称
// @Tags Device
// @Summary 获取设备类型名称
// @Description 根据物理设备索引 dvInd，返回设备类型名称及单位
// @Accept  json
// @Produce  json
// @Param   dvInd   query   int  true  "物理设备索引（Device Index）"
// @Success 200 {object} DevTypeNameResp "设备类型名称信息"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /device/type/name [get]
func DevTypeName(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	devTypeName, unit, err := dcgm.DevTypeName(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := DevTypeNameResp{
		DevTypeName: devTypeName,
		Unit:        unit,
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// DevSubsystemId 获取设备子系统 ID
// @Tags Device
// @Summary 获取设备子系统 ID
// @Description 根据物理设备索引 dvInd，返回设备子系统 ID
// @Accept  json
// @Produce  json
// @Param   dvInd   query   int  true  "物理设备索引（Device Index）"
// @Success 200 {object} DevSubsystemIdResp "设备子系统 ID"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /device/subsystem/id [get]
func DevSubsystemId(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	subsystemId, err := dcgm.DevSubsystemId(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := DevSubsystemIdResp{
		SubsystemId: subsystemId,
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// DevSubsystemName 获取设备子系统名称
// @Tags Device
// @Summary 获取设备子系统名称
// @Description 根据物理设备索引 dvInd，返回设备子系统名称
// @Accept  json
// @Produce  json
// @Param   dvInd   query   int  true  "物理设备索引（Device Index）"
// @Success 200 {object} DevSubsystemNameResp "设备子系统名称"
// @Failure 400 {object} error "Invalid device ID"
// @Failure 500 {object} error "Internal Server Error"
// @Router /device/subsystem/name [get]
func DevSubsystemName(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	subsystemName, err := dcgm.DevSubsystemName(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := DevSubsystemNameResp{
		SubsystemName: subsystemName,
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Power
// @Summary 获取设备的最大功率
// @Description 根据设备索引返回设备的最大功率（以瓦特为单位）
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {object} MaxPowerResp "设备最大功率"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /MaxPower [get]
func GetMaxPower(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	power, err := dcgm.MaxPower(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := MaxPowerResp{
		Power: power,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Memory
// @Summary 获取设备的指定内存使用情况
// @Description 根据设备索引和内存类型返回内存的使用量和总量
// @Produce json
// @Param dvInd query int true "设备索引"
// @Param memType query string true "内存类型（可选值: vram, vis_vram, gtt）"
// @Success 200 {object} MemInfoResp "返回指定内存类型的使用量和总量"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /MemInfo [get]
func GetMemInfo(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	memType := c.Query("memType")
	memUsed, memTotal, err := dcgm.MemInfo(dvInd, memType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := MemInfoResp{
		MemUsed:  memUsed,
		MemTotal: memTotal,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Device
// @Summary 获取设备信息列表
// @Description 返回所有设备的详细信息列表
// @Produce json
// @Success 200 {object} DeviceInfo "返回设备信息列表"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DeviceInfos [get]
func DeviceInfos(c *gin.Context) {
	deviceInfos, err := dcgm.DeviceInfos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, deviceInfos)
}

// DeviceStatus 获取设备状态信息
// @Summary 获取设备状态信息
// @Description 返回所有设备设备状态信息
// @Produce json
// @Success 200 {array} DeviceStatusInfo "返回设备设备状态信息"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DeviceStatus [get]
func DeviceStatus(c *gin.Context) {
	deviceStatusInfos, err := dcgm.DeviceStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, deviceStatusInfos)
}

// @Tags Device
// @Summary 获取 DF 带宽信息
// @Description 根据物理设备索引 dvInd 获取 DF（Data Fabric）带宽信息
// @Accept  json
// @Produce  json
// @Param   dvInd path int true "物理设备索引（Device Index）"
// @Success 200 {object} DFBandwidthResp "DF 带宽信息"
// @Failure 400 {object} FailedMessage "Invalid device ID"
// @Failure 500 {object} FailedMessage "Internal Server Error"
// @Router /dfbandwidth/{dvInd} [get]
func DFBandwidth(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	info, err := dcgm.DFBandwidth(dvInd, dcgm.RSMI_DF_BW_TYPE_ALL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// 手动转换
	resp := DFBandwidthResp{
		DFBandwidthInfo: DFBandwidthInfo{
			ReadBW:      info.ReadBW,
			WriteBW:     info.WriteBW,
			ReadWriteBW: info.ReadWriteBW,
		},
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Utilization
// @Summary 获取 HCU 瞬时占用率
// @Description 查询指定设备当前的 HCU 占用率（瞬时值）。只要某个 CU 内存在至少一个活跃 wave，即认为该 CU 活跃。对应 hy-smi -u / rsmi_dev_cu_usage_get。
// @Accept json
// @Produce json
// @Param dvInd path int true "物理设备索引"
// @Success 200 {object} HCUCuUsageResp "HCU 瞬时占用率"
// @Failure 400 {object} FailedMessage "Invalid device ID"
// @Failure 500 {object} FailedMessage "Internal Server Error"
// @Router /HCUCuUsage/{dvInd} [get]
func HCUCuUsage(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	utilizationRate, err := dcgm.HCUCuUsage(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(HCUCuUsageResp{UtilizationRate: utilizationRate}))
}

// @Tags Utilization
// @Summary 获取 HCU 采样占用情况
// @Description 在采样窗口内周期性统计 HCU 活跃状态占比。对应 hy-smi --showhcuutil / rsmi_dev_hcu_util_get；默认采样窗口 1000ms（1s）。
// @Accept json
// @Produce json
// @Param dvInd path int true "物理设备索引"
// @Param sampleDurationMs query int false "采样时间窗口（毫秒），默认 1000"
// @Success 200 {object} DevUtilSampleResp "HCU 采样占用率"
// @Failure 400 {object} FailedMessage "Invalid parameters"
// @Failure 500 {object} FailedMessage "Internal Server Error"
// @Router /HCUSampledUsage/{dvInd} [get]
func HCUSampledUsage(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	sampleDurationMs, err := parseSampleDurationMs(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid sampleDurationMs"))
		return
	}

	utilizationRate, err := dcgm.HCUSampledUsage(dvInd, sampleDurationMs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(DevUtilSampleResp{
		UtilizationRate:  utilizationRate,
		SampleDurationMs: sampleDurationMs,
	}))
}

// @Tags Utilization
// @Summary 获取 CU 采样占用情况
// @Description 在采样窗口内统计各 CU 活跃占比并取平均值。对应 hy-smi --showcuutil / rsmi_dev_cu_util_get；默认采样窗口 1000ms。
// @Accept json
// @Produce json
// @Param dvInd path int true "物理设备索引"
// @Param sampleDurationMs query int false "采样时间窗口（毫秒），默认 1000"
// @Success 200 {object} DevUtilSampleResp "CU 平均采样占用率"
// @Failure 400 {object} FailedMessage "Invalid parameters"
// @Failure 500 {object} FailedMessage "Internal Server Error"
// @Router /HCUCUSampledUsage/{dvInd} [get]
func HCUCUSampledUsage(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	sampleDurationMs, err := parseSampleDurationMs(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid sampleDurationMs"))
		return
	}

	utilizationRate, err := dcgm.HCUCUSampledUsage(dvInd, sampleDurationMs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(DevUtilSampleResp{
		UtilizationRate:  utilizationRate,
		SampleDurationMs: sampleDurationMs,
	}))
}

// @Tags Utilization
// @Summary 获取 Wave 采样占用情况
// @Description 在采样窗口内统计各 CU 上活跃 wave 占比并取平均值。对应 hy-smi --showwaveutil / rsmi_dev_wave_util_get；默认采样窗口 1000ms。
// @Accept json
// @Produce json
// @Param dvInd path int true "物理设备索引"
// @Param sampleDurationMs query int false "采样时间窗口（毫秒），默认 1000"
// @Success 200 {object} DevUtilSampleResp "Wave 平均采样占用率"
// @Failure 400 {object} FailedMessage "Invalid parameters"
// @Failure 500 {object} FailedMessage "Internal Server Error"
// @Router /HCUWaveSampledUsage/{dvInd} [get]
func HCUWaveSampledUsage(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	sampleDurationMs, err := parseSampleDurationMs(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid sampleDurationMs"))
		return
	}

	utilizationRate, err := dcgm.HCUWaveSampledUsage(dvInd, sampleDurationMs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(DevUtilSampleResp{
		UtilizationRate:  utilizationRate,
		SampleDurationMs: sampleDurationMs,
	}))
}

// @Tags Utilization
// @Summary 获取 SE 瞬时占用率
// @Description 按 SE 维度返回活跃 CU 占比（瞬时值）。对应 hy-smi --showseuse / rsmi_dev_se_util_get。
// @Accept json
// @Produce json
// @Param dvInd path int true "物理设备索引"
// @Success 200 {object} HCUSEUsageResp "各 SE 占用率"
// @Failure 400 {object} FailedMessage "Invalid device ID"
// @Failure 500 {object} FailedMessage "Internal Server Error"
// @Router /HCUSEUsage/{dvInd} [get]
func HCUSEUsage(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	usage, err := dcgm.HCUSEUsage(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := HCUSEUsageResp{
		ShaderEngineUsage: SEUsageInfo{Percent: usage.Percent},
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// parseSampleDurationMs 解析采样窗口 query 参数，缺省为 1000ms（与 hy-smi 默认 1s 一致）。
func parseSampleDurationMs(c *gin.Context) (int, error) {
	if v := c.Query("sampleDurationMs"); v != "" {
		return strconv.Atoi(v)
	}
	return 1000, nil
}

// @Tags Device
// @Summary 获取 UMC 带宽信息
// @Description 根据物理设备索引和 UMC 通道 ID 获取 UMC 带宽信息
// @Accept json
// @Produce json
// @Param body body UMCBandwidthReq true "UMC 带宽查询参数"
// @Success 200 {object} UMCBandwidthResp "UMC 带宽信息"
// @Failure 400 {object} FailedMessage "Invalid parameters"
// @Failure 500 {object} FailedMessage "Internal Server Error"
// @Router /umc/bandwidth [post]
func UMCBandwidth(c *gin.Context) {
	var params UMCBandwidthReq

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("参数无效"))
		return
	}

	info, err := dcgm.UMCBandwidth(params.DvInd, params.ChanId, params.Delay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// 转换类型
	resp := UMCBandwidthResp{
		UMCBandwidth: UMCBandwidthInfo{
			ReadBW:      info.ReadBW[:],
			WriteBW:     info.WriteBW[:],
			ReadWriteBW: info.ReadWriteBW[:],
		},
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// XHCLBandwidth godoc
// @Tags Bandwidth
// @Summary 查询 XHCL 带宽信息
// @Description 根据物理设备索引和链路参数获取 XHCL 带宽信息
// @Accept json
// @Produce json
// @Param req body XHCLBandwidthReq true "XHCL 带宽查询参数"
// @Success 200 {object} XHCLBandwidthResp
// @Failure 400 {object} Response "参数错误"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /bandwidth/xhcl [post]
func XHCLBandwidth(c *gin.Context) {
	var params XHCLBandwidthReq

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("参数无效"))
		return
	}

	info, err := dcgm.XHCLBandwidth(params.DvInd, params.LinkId, params.Direction, params.Delay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// 转换类型
	resp := XHCLBandwidthResp{
		XhclBandwidth: XhclBandwidthInfo{
			Bw: info.Bw[:],
		},
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// HyLinkLinkStatus godoc
// @Summary      查询所有 HCU 的 HyLink 带宽状态
// @Description  返回所有监控设备的 HyLink 带宽汇总及每条 Link 的收发带宽明细。
// @Tags         HyLink
// @Produce      json
// @Success      200 {object} HyLinkStatusResp "查询成功"
// @Failure      500 {object} Response "服务器内部错误"
// @Router       /GetHyLinkLinkStatus [get]
func HyLinkLinkStatus(c *gin.Context) {
	dcgmDeviceSums, err := dcgm.GetHyLinkStatus()
	if err != nil {
		glog.Errorf("GetHyLinkStatus failed: %v", err)
		c.JSON(http.StatusInternalServerError,
			ErrorResponse("获取设备 link 带宽失败"))
		return
	}

	// 转换类型
	var deviceSums []DeviceLinkSum
	for _, d := range dcgmDeviceSums {
		links := make([]LinkBandwidth, len(d.Links))
		for i, l := range d.Links {
			links[i] = LinkBandwidth{
				LinkId: l.LinkId,
				Recv:   l.Recv,
				Send:   l.Send,
			}
		}

		deviceSums = append(deviceSums, DeviceLinkSum{
			DvInd: d.DvInd,
			Recv:  d.Recv,
			Send:  d.Send,
			Err:   d.Err,
			Links: links,
		})
	}

	resp := HyLinkStatusResp{
		Devices: deviceSums,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// HyLinkStatusByHcuId godoc
// @Summary      查询指定 HCU 的 HyLink 带宽状态
// @Description  根据物理设备索引（DvInd）查询该 HCU 的 HyLink 带宽信息。
// @Description  返回该设备所有 Link 的接收 / 发送带宽汇总及明细。
// @Tags         HyLink
// @Produce      json
// @Param        dvInd path int true "HCU 设备索引（从 0 开始）"
// @Success      200 {object} HyLinkStatusByHcuIdResp "查询成功"
// @Failure      400 {object} Response "参数错误（非法的 dvInd）"
// @Failure      500 {object} Response "服务器内部错误"
// @Router       /HyLinkStatusByHcuId/{dvInd} [get]
func HyLinkStatusByHcuId(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {

		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	info, err := dcgm.HyLinkStatusByHcuId(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	links := make([]LinkBandwidth, len(info.Links))
	for i, l := range info.Links {
		links[i] = LinkBandwidth{
			LinkId: l.LinkId,
			Recv:   l.Recv,
			Send:   l.Send,
		}
	}

	resp := HyLinkStatusByHcuIdResp{
		DeviceLinkBandwidth: DeviceLinkSum{
			DvInd: info.DvInd,
			Recv:  info.Recv,
			Send:  info.Send,
			Err:   info.Err,
			Links: links,
		},
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))

}

// HyUMCStatus godoc
// @Summary      查询所有设备的 UMC 带宽状态
// @Description  查询当前所有监控设备的 UMC 带宽聚合信息。
// @Description  每个设备返回 Read、Write 和 ReadWrite 指标。
// @Tags         HyLink
// @Produce      json
// @Success      200 {object} HyUMCStatusResp "查询成功"
// @Failure      500 {object} Response "服务器内部错误"
// @Router       /GetHyUMCStatus [get]
func HyUMCStatus(c *gin.Context) {
	deviceSums, err := dcgm.GetHyUMCStatus()
	if err != nil {
		glog.Errorf("GetHyUMCStatus failed: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse("获取设备 UMC 带宽失败"))
		return
	}
	deviceUmcBandwidth := make([]DeviceUmcSum, len(deviceSums))
	for i, info := range deviceSums {
		deviceUmcBandwidth[i] = DeviceUmcSum{
			DvInd:     info.DvInd,
			Read:      info.Read,
			Write:     info.Write,
			ReadWrite: info.ReadWrite,
			Err:       info.Err,
		}
	}
	resp := HyUMCStatusResp{
		DeviceUmcBandwidth: deviceUmcBandwidth,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags System
// @Summary 获取指定PID的进程名
// @Description 根据进程ID（PID）返回对应的进程名称
// @Produce json
// @Param pid query int true "进程ID"
// @Success 200 {object} GetProcessNameResp "返回进程名称"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ProcessName [get]
func GetProcessName(c *gin.Context) {
	pid, err := strconv.Atoi(c.Query("pid"))
	if err != nil || pid < 1 {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求参数错误"))
		return
	}

	pName := dcgm.ProcessName(pid)

	resp := GetProcessNameResp{
		ProcessName: pName,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// ProcessInfo godoc
// @Summary      获取指定进程的计算信息
// @Description  根据进程 PID 查询该进程的计算资源使用情况，包括显存使用量、SDMA 使用量和 CU 占用率
// @Tags         Process
// @Accept       json
// @Produce      json
// @Param        pid   query   int   true  "进程 PID"
// @Success      200   {object} ProcessInfoResp
// @Failure      400   {object} Response
// @Failure      500   {object} Response
// @Router       /process/info [get]
func ProcessInfo(c *gin.Context) {
	pid, err := strconv.Atoi(c.Query("pid"))
	if err != nil || pid < 1 {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求参数错误"))
		return
	}

	proc, err := dcgm.ProcessInfo(pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// 直接映射单个结构体
	resp := ProcessInfoResp{
		ProcessInfo: ProcessInfos{
			ProcessID:   proc.ProcessID,
			Pasid:       proc.Pasid,
			VramUsage:   proc.VramUsage,
			SdmaUsage:   proc.SdmaUsage,
			CuOccupancy: proc.CuOccupancy,
		},
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Performance
// @Summary 获取设备的当前性能水平
// @Description 返回指定设备的当前性能等级
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} PerfLevelResp "返回当前性能水平"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /PerfLevel/{dvInd} [get]
func PerfLevel(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}

	perf, err := dcgm.PerfLevel(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := PerfLevelResp{
		PerfLevel: perf,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Power
// @Summary 获取设备的平均功耗
// @Description 返回指定设备的平均功耗（瓦特）
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} PowerResp "返回平均功耗"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Power/{dvInd} [get]
func Power(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}

	power, err := dcgm.Power(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := PowerResp{
		Power: power,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Error
// @Summary 获取GPU块的ECC状态
// @Description 返回指定GPU块的ECC状态
// @Produce json
// @Param dvInd path int true "设备索引"
// @Param block query string true "GPU块"
// @Success 200 {object} EccStatusResp "返回ECC状态"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /EccStatus/{dvInd} [get]
func EccStatus(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}

	block := c.Query("block")
	if block == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("缺少 block 参数"))
		return
	}

	// 将 block 字符串转换为 RSMIGpuBlock 类型
	blockConverted, err := ConvertToRSMIGpuBlock(block)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("无效的 block 参数"))
		return
	}

	eccStatus, err := dcgm.EccStatus(dvInd, blockConverted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := EccStatusResp{
		EccStatus: eccStatus,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Physical State
// @Summary 获取设备温度
// @Description 返回指定设备的当前温度
// @Produce json
// @Param dvInd path int true "设备索引"
// @Param sensorType query int true "传感器类型: 0: Edge GPU temperature; 1: Junction/hotspot temperature; 2: VRAM temperature"
// @Success 200 {object} TemperatureResp "返回设备温度"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Temperature/{dvInd} [get]
func Temperature(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}

	temp, err := dcgm.Temperature(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := TemperatureResp{
		Temp: temp,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Version
// @Summary 获取RSMI版本
// @Description 返回当前系统的RSMI版本
// @Produce json
// @Success 200 {object} RsmiVersionResp "返回RSMI版本信息"
// @Failure 500 {object} error "服务器内部错误"
// @Router /RsmiVersion [get]
func RsmiVersion(c *gin.Context) {
	dcgmVersion, err := dcgm.HCUVersion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	// 构建 router 层的返回结构体
	resp := RsmiVersionResp{
		RsmiVersion: DevVersion{
			Major: dcgmVersion.Major,
			Minor: dcgmVersion.Minor,
			Patch: dcgmVersion.Patch,
			Build: dcgmVersion.Build,
		},
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Version
// @Summary 获取设备的VBIOS版本
// @Description 返回指定设备的VBIOS版本
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} VbiosVersionResp "返回设备VBIOS版本信息"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /VbiosVersion/{dvInd} [get]
func VbiosVersion(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}

	vbios, err := dcgm.VbiosVersion(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := VbiosVersionResp{
		Vbios: vbios,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Performance
// @Summary 设置 GPU 时钟频率
// @Description 设置 GPU 上指定时钟的允许频率。clkType 设置为默认值，无需传递。
// @Accept json
// @Produce json
// @Param dvInd query int true "设备索引"
// @Param freqBitmask query int64 true "频率掩码"
// @Success 200 {object} DevGpuClkFreqSetResp "操作结果"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DevGpuClkFreqSet [post]
func DevGpuClkFreqSet(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("无效的 dvInd 参数"))
		return
	}

	freqBitmask, err := strconv.ParseInt(c.Query("freqBitmask"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("无效的 freqBitmask 参数"))
		return
	}

	// 调用 dcgm 包函数设置 GPU 时钟频率
	err = dcgm.DevGpuClkFreqSet(dvInd, dcgm.RSMI_CLK_TYPE_SYS, freqBitmask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := DevGpuClkFreqSetResp{
		Message: "设置成功",
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Version
// Version 获取当前系统的驱动程序版本
// @Summary 获取当前系统的驱动程序版本
// @Description 返回指定组件的驱动程序版本
// @Produce json
// @Param component query string true "驱动组件:FIRST、DRIVER、LAST"
// @Success 200 {object} VersionResp "返回驱动程序版本信息"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Version [get]
func Version(c *gin.Context) {
	component := c.Query("component")
	if component == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("Missing component parameter"))
		return
	}

	// 将 component 字符串转换为 RSMISwComponent 类型
	componentConverted, err := ConvertToRSMISwComponent(component)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid component parameter"))
		return
	}

	version, err := dcgm.Version(componentConverted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := VersionResp{
		Version: version,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Performance
// ResetClocks 将设备的时钟重置为默认值
// @Summary 重置设备时钟
// @Description 重置指定设备的时钟和性能等级为默认值
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {object} FailedResp "返回失败消息列表"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ResetClocks [post]
func ResetClocks(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON data"))
		return
	}

	dcgmFailed := dcgm.ResetClocks(dvIdList)

	// 转换类型
	var failedMessages []FailedMessage
	for _, f := range dcgmFailed {
		failedMessages = append(failedMessages, FailedMessage{
			ID:       f.ID,
			ErrorMsg: f.ErrorMsg,
		})
	}

	resp := FailedResp{
		FailedMessages: failedMessages,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Physical State
// ResetFans 复位风扇驱动控制
// @Summary 复位风扇控制
// @Description 重置指定设备的风扇控制为默认值
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {string} string "复位成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ResetFans [post]
func ResetFans(c *gin.Context) {
	var dvIdList []int
	if err := c.ShouldBindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	if err := dcgm.ResetFans(dvIdList); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Performance
// ResetProfile 重置设备的配置文件
// @Summary 重置指定设备的电源配置文件和性能级别
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {object} FailedResp "返回失败的设备及其错误信息"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ResetProfile [post]
func ResetProfile(c *gin.Context) {
	var dvIdList []int
	if err := c.ShouldBindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	dcgmFailed := dcgm.ResetProfile(dvIdList)

	// 转换类型
	var failedMessages []FailedMessage
	for _, f := range dcgmFailed {
		failedMessages = append(failedMessages, FailedMessage{
			ID:       f.ID,
			ErrorMsg: f.ErrorMsg,
		})
	}

	resp := FailedResp{
		FailedMessages: failedMessages,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Error
// ResetHSLErr 重置设备的HSL错误状态
// @Summary 重置指定设备的HSL错误状态
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {object} FailedResp "返回失败的设备及其错误信息"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ResetHSLErr [post]
func ResetHSLErr(c *gin.Context) {
	var dvIdList []int
	if err := c.ShouldBindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	dcgmFailed := dcgm.ResetHSLErr(dvIdList)

	// 转换类型
	var failedMessages []FailedMessage
	for _, f := range dcgmFailed {
		failedMessages = append(failedMessages, FailedMessage{
			ID:       f.ID,
			ErrorMsg: f.ErrorMsg,
		})
	}

	resp := FailedResp{
		FailedMessages: failedMessages,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Error
// HSLErrorStatus 获取HSL错误状态
// @Summary 获取HSL错误状态
// @Description 获取指定物理设备的HSL（高速互连链路）错误状态。
// @Param dvInd query int true "物理设备的索引"
// @Success 200 {object} HSLErrorStatusResp "返回HSL错误状态码"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HSLErrorStatus [get]
func HSLErrorStatus(c *gin.Context) {
	dvIndStr := c.Query("dvInd")
	dvInd, err := strconv.Atoi(dvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	status, err := dcgm.HSLErrorStatus(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := HSLErrorStatusResp{
		Status: int(status),
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Error
// HSLHiveIdGet 获取设备的HSL hive id
// @Summary 获取指定设备的HSL hive id
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {object} HSLHiveIdResp "返回设备的HSL hive id"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HSLHiveIdGet [get]
func HSLHiveIdGet(c *gin.Context) {
	dvIndStr := c.Query("dvInd")
	dvInd, err := strconv.Atoi(dvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	hiveId, err := dcgm.HSLHiveIdGet(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := HSLHiveIdResp{
		HiveId: hiveId,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Performance
// ResetPerfDeterminism 处理重置Performance Determinism
// @Summary 重置Performance Determinism
// @Description 该接口用于重置指定设备的性能决定性设置。请求体中需要包含设备ID列表。
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {object} ResetPerfDeterminismResp "成功或失败信息"
// @Failure 400 {object} error "无效的请求体或部分设备失败"
// @Router /ResetPerfDeterminism [post]
func ResetPerfDeterminism(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	dcgmFailedMessages := dcgm.ResetPerfDeterminism(dvIdList)

	// 直接在这里转换类型
	failedMessages := make([]FailedMessage, len(dcgmFailedMessages))
	for i, msg := range dcgmFailedMessages {
		failedMessages[i] = FailedMessage{
			ID:       msg.ID,
			ErrorMsg: msg.ErrorMsg,
		}
	}

	resp := ResetPerfDeterminismResp{
		FailedMessages: failedMessages,
	}

	if len(resp.FailedMessages) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse(resp))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// @Tags Performance
// SetClockRange 处理设置时钟频率范围
// @Summary 设置设备的时钟频率范围
// @Description 设置设备的时钟频率范围（sclk 或 mclk）
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Param clkType query string true "时钟类型（sclk 或 mclk）"
// @Param minvalue query string true "最小值（MHz）"
// @Param maxvalue query string true "最大值（MHz）"
// @Param autoRespond query bool false "自动响应超出规格的警告"
// @Success 200 {object} SetClockRangeResp "返回失败信息或成功"
// @Failure 400 {object} error "无效的请求参数或无法设置时钟范围"
// @Router /SetClockRange [post]
func SetClockRange(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	clkType := c.Query("clkType")
	minvalue := c.Query("minvalue")
	maxvalue := c.Query("maxvalue")
	autoRespond, _ := strconv.ParseBool(c.Query("autoRespond"))

	dcgmFailedMessages := dcgm.SetClockRange(dvIdList, clkType, minvalue, maxvalue, autoRespond)

	// 类型转换
	failedMessages := make([]FailedMessage, len(dcgmFailedMessages))
	for i, msg := range dcgmFailedMessages {
		failedMessages[i] = FailedMessage{
			ID:       msg.ID,
			ErrorMsg: msg.ErrorMsg,
		}
	}

	resp := SetClockRangeResp{
		FailedMessages: failedMessages,
	}

	if len(resp.FailedMessages) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse(resp))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// SetPowerPlayTableLevel 处理设置PowerPlay表级别
// @Tags Performance
// @Summary 设置设备的PowerPlay表级别
// @Description 设置设备的PowerPlay表级别
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Param clkType query string true "时钟类型（sclk 或 mclk）"
// @Param point query string true "电压点"
// @Param clk query string true "时钟值（MHz）"
// @Param volt query string true "电压值（mV）"
// @Param autoRespond query bool false "自动响应超出规格的警告"
// @Success 200 {object} PowerPlayResp "成功设置PowerPlay表级别"
// @Failure 400 {object} PowerPlayResp "返回失败消息列表"
// @Router /SetPowerPlayTableLevel [post]
func SetPowerPlayTableLevel(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	clkType := c.Query("clkType")
	point := c.Query("point")
	clk := c.Query("clk")
	volt := c.Query("volt")
	autoRespond, _ := strconv.ParseBool(c.Query("autoRespond"))

	dcgmFailedMessages := dcgm.SetPowerPlayTableLevel(dvIdList, clkType, point, clk, volt, autoRespond)

	// 类型转换
	failedMessages := make([]FailedMessage, len(dcgmFailedMessages))
	for i, msg := range dcgmFailedMessages {
		failedMessages[i] = FailedMessage{
			ID:       msg.ID,
			ErrorMsg: msg.ErrorMsg,
		}
	}

	resp := PowerPlayResp{
		FailedMessages: failedMessages,
	}

	if len(resp.FailedMessages) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse(resp))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// SetClockOverDrive 处理设置时钟OverDrive
// @Tags Performance
// @Summary 设置设备的时钟OverDrive
// @Description 设置设备的时钟OverDrive
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Param clktype query string true "时钟类型（sclk 或 mclk）"
// @Param value query string true "OverDrive值，表示为百分比（0-20%）"
// @Param autoRespond query bool false "自动响应超出规格的警告"
// @Success 200 {object} ClockOverDriveResp "成功设置时钟OverDrive"
// @Failure 400 {object} ClockOverDriveResp "返回失败消息列表"
// @Router /SetClockOverDrive [post]
func SetClockOverDrive(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	clktype := c.Query("clktype")
	value := c.Query("value")
	autoRespond, _ := strconv.ParseBool(c.Query("autoRespond"))

	dcgmFailedMessages := dcgm.SetClockOverDrive(dvIdList, clktype, value, autoRespond)

	// 类型转换
	failedMessages := make([]FailedMessage, len(dcgmFailedMessages))
	for i, msg := range dcgmFailedMessages {
		failedMessages[i] = FailedMessage{
			ID:       msg.ID,
			ErrorMsg: msg.ErrorMsg,
		}
	}

	resp := ClockOverDriveResp{
		FailedMessages: failedMessages,
	}

	if len(resp.FailedMessages) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse(resp))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// @Tags Performance
// setPerfDeterminism 处理设置性能确定性
// @Summary 设置设备的性能确定性
// @Description 设置设备的性能确定性
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Param clkvalue query string true "时钟频率值"
// @Success 200 {array} FailedMessage "返回失败的设备及其错误信息"
// @Failure 400 {object} error "无效的请求体或无法设置性能确定性"
// @Router /SetPerfDeterminism [post]
func SetPerfDeterminism(c *gin.Context) {
	var dvIdList []int
	clkvalue := c.Query("clkvalue")

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	failedMessages, err := dcgm.SetPerfDeterminism(dvIdList, clkvalue)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	if len(failedMessages) > 0 {
		c.JSON(http.StatusOK, failedMessages)
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// @Tags Physical State
// SetFanSpeed 设置风扇转速
// @Summary 设置风扇转速
// @Description 根据设备ID列表和给定的风扇速度，设置设备的风扇速度
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param fan query string true "风扇速度值（0-255,单位:RPM）"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /SetFanSpeed [post]
func SetFanSpeed(c *gin.Context) {
	var dvIdList []int
	fan := c.Query("fan")

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	dcgm.SetFanSpeed(dvIdList, fan)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// DevFanRpms 获取设备的风扇速度
// @Tags Physical State
// @Summary 获取设备的风扇速度
// @Description 获取指定设备的风扇速度（RPM）
// @Accept json
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} FanSpeedResp "返回风扇速度"
// @Failure 400 {object} error "请求参数错误或获取失败"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DevFanRpms/{dvInd} [get]
func DevFanRpms(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	speed, err := dcgm.DevFanRpms(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := FanSpeedResp{
		Speed: speed,
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Performance
// SetPerformanceLevel 设置设备性能等级
// @Summary 设置设备性能等级
// @Description 根据设备ID列表和给定的性能等级，设置设备的性能等级
// @Accept  json
// @Produce  json
// @Param deviceList body []int true "设备 ID 列表"
// @Param level query string true "性能等级 (auto, low, high, normal)"
// @Success 200 {array} FailedMessage
// @Failure 400 {object} FailedMessage
// @Router /SetPerformanceLevel [post]
func SetPerformanceLevel(c *gin.Context) {
	var deviceList []int
	level := c.Query("level")

	if err := c.BindJSON(&deviceList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	failedMessages := dcgm.SetPerformanceLevel(deviceList, level)
	if len(failedMessages) > 0 {
		c.JSON(http.StatusOK, failedMessages)
	} else {
		c.JSON(http.StatusOK, SuccessResponse(nil))
	}
}

// @Tags Power
// SetProfile 设置功率配置
// @Summary 设置功率配置
// @Description 根据设备ID列表和给定的功率配置文件，设置设备的功率配置
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param profile query string true "功率配置文件名称:CUSTOM、VIDEO、POWER SAVING、COMPUTE、VR、3D FULL SCREEN、BOOTUP DEFAULT"
// @Success 200 {array} FailedMessage "设置成功的消息列表"
// @Failure 400 {object} FailedMessage "失败的消息列表"
// @Router /SetProfile [post]
func SetProfile(c *gin.Context) {
	var dvIdList []int
	profile := c.Query("profile")

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	failedMessages := dcgm.SetProfile(dvIdList, profile)
	c.JSON(http.StatusOK, failedMessages)
}

// @Tags Power
// DevPowerProfileSet 设置设备功率配置文件
// @Summary 设置设备功率配置文件
// @Description 设置指定设备的功率配置文件
// @Accept  json
// @Produce  json
// @Param dvInd path int true "设备索引"
// @Param reserved query int true "保留参数，通常为0"
// @Param profile query int true "功率配置文件的枚举值"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /DevPowerProfileSet [post]
func DevPowerProfileSet(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	reserved, err := strconv.Atoi(c.Query("reserved"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid reserved value"))
		return
	}

	profileEnum, err := strconv.Atoi(c.Query("profile"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid profile value"))
		return
	}

	err = dcgm.DevPowerProfileSet(dvInd, reserved, dcgm.PowerProfilePresetMasks(profileEnum))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// GetBus 获取设备总线信息
// @Tags System
// @Summary 获取设备总线信息
// @Description 获取指定设备的总线信息
// @Accept json
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} BusInfoResp "返回设备总线信息"
// @Failure 400 {object} error "请求参数错误或获取失败"
// @Failure 500 {object} error "服务器内部错误"
// @Router /GetBus/{dvInd} [get]
func GetBus(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	picId, err := dcgm.GetBus(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	resp := BusInfoResp{
		PicID: picId,
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Device
// ShowAllConciseHw 显示设备硬件信息
// @Summary 显示设备硬件信息
// @Description 显示指定设备列表的简要硬件信息
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /ShowAllConciseHw [post]
func ShowAllConciseHw(c *gin.Context) {
	var dvIdList []int

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	dcgm.ShowAllConciseHw(dvIdList)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Performance
// ShowClocks 显示时钟信息
// @Summary 显示时钟信息
// @Description 显示指定设备的时钟信息
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /ShowClocks [post]
func ShowClocks(c *gin.Context) {
	var dvIdList []int

	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	dcgm.ShowClocks(dvIdList)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Physical State
// ShowCurrentFans 展示风扇转速和风扇级别
// @Summary 展示风扇转速和风扇级别
// @Description 显示指定设备的当前风扇转速和风扇级别
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /fans/current [post]
func ShowCurrentFans(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	dcgm.ShowCurrentFans(dvIdList, true)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// ShowCurrentTemps 显示设备温度传感器数据
// @Tags Physical State
// @Summary 显示设备温度传感器数据
// @Description 获取指定设备列表的温度传感器数据
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} ShowCurrentTempsResp "返回设备温度信息列表"
// @Failure 400 {object} error "请求参数错误或获取温度失败"
// @Router /temps/current [post]
func ShowCurrentTemps(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid request body"))
		return
	}

	// 调用 dcgm 包获取温度信息
	dcgmTemps, err := dcgm.ShowCurrentTemps(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// 类型转换：dcgm.TemperatureInfo -> router.TemperatureInfo
	temps := make([]TemperatureInfo, 0, len(dcgmTemps))
	for _, t := range dcgmTemps {
		temps = append(temps, TemperatureInfo{
			DeviceID:    t.DeviceID,
			SensorTemps: t.SensorTemps,
		})
	}

	// 返回结构体
	resp := ShowCurrentTempsResp{
		TemperatureInfos: temps,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// ShowFwInfo 显示设备固件版本信息
// @Tags Version
// @Summary 显示设备固件版本信息
// @Description 获取指定设备列表和固件类型列表的固件版本信息
// @Accept json
// @Produce json
// @Param dvIdList query []int true "设备 ID 列表"
// @Param fwType query []string true "固件类型列表"
// @Success 200 {object} ShowFwInfoResp "固件版本信息列表"
// @Failure 400 {object} error "请求参数错误或获取固件信息失败"
// @Router /firmware/info [get]
func ShowFwInfo(c *gin.Context) {
	var dvIdList []int
	var fwType []string

	if err := c.BindQuery(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid query parameters for dvIdList"))
		return
	}
	if err := c.BindQuery(&fwType); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid query parameters for fwType"))
		return
	}

	dcgmFwInfos, err := dcgm.ShowFwInfo(dvIdList, fwType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	fwInfos := make([]FirmwareInfoResp, 0, len(dcgmFwInfos))
	for _, f := range dcgmFwInfos {
		blocks := make([]FirmwareBlock, 0, len(f.FirmwareVer))
		for name, ver := range f.FirmwareVer {
			blocks = append(blocks, FirmwareBlock{
				BlockName: name,
				Version:   ver,
			})
		}
		fwInfos = append(fwInfos, FirmwareInfoResp{
			DeviceID:    f.DeviceID,
			FirmwareVer: blocks,
		})
	}

	resp := ShowFwInfoResp{
		FwInfos: fwInfos,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// PidList 获取计算进程列表
// @Tags System
// @Summary 获取计算进程列表
// @Description 返回当前系统计算进程的 PID 列表
// @Produce json
// @Success 200 {object} PidListResp "进程 ID 列表"
// @Failure 400 {object} error "获取进程列表失败"
// @Router /process/list [get]
func PidList(c *gin.Context) {
	pids, err := dcgm.PidList()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	resp := PidListResp{
		PidList: pids,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// ProcessHCUInfo 获取进程信息及其使用的 HCU 设备信息
// @Tags System
// @Summary 进程列表信息
// @Description 获取系统中计算进程的信息及其使用的 HCU 设备索引
// @Produce json
// @Success 200 {object} ProcessHCUInfoResp "成功返回进程信息"
// @Failure 400 {object} error "请求错误或获取进程信息失败"
// @Router /processHCUInfo [get]
func ProcessHCUInfo(c *gin.Context) {
	dcgmProcesses, err := dcgm.ProcessHCUInfo()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	// 类型转换：dcgm.Process -> router.Process
	var processes []Process
	for _, p := range dcgmProcesses {
		processes = append(processes, Process{
			ProcessID:    p.ProcessID,
			ProcessName:  p.ProcessName,
			Pasid:        p.Pasid,
			VramUsage:    p.VramUsage,
			SdmaUsage:    p.SdmaUsage,
			CuOccupancy:  p.CuOccupancy,
			MinorNumbers: p.MinorNumbers,
		})
	}

	resp := ProcessHCUInfoResp{
		ProcessInfo: processes,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Performance
// GetCoarseGrainUtil 获取设备粗粒度利用率
// @Summary 获取设备粗粒度利用率
// @Accept json
// @Produce json
// @Param body body GetCoarseGrainUtilReq true "请求参数"
// @Success 200 {object} GetCoarseGrainUtilResp "利用率计数器列表"
// @Failure 400 {object} error "错误信息"
// @Router /utilization/coarse [post]
func GetCoarseGrainUtil(c *gin.Context) {
	var req GetCoarseGrainUtilReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmCounters, err := dcgm.GetCoarseGrainUtil(req.Device, req.TypeName)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	// 类型转换：dcgm.UtilizationCounter -> router.UtilizationCounter
	counters := make([]UtilizationCounter, 0, len(dcgmCounters))
	for _, c := range dcgmCounters {
		counters = append(counters, UtilizationCounter{
			Type:  UtilizationCounterType(c.Type),
			Value: c.Value,
		})
	}

	resp := GetCoarseGrainUtilResp{
		UtilizationCounters: counters,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Device
// ShowGpuUse 显示设备的 HCU 使用率
// @Summary 显示设备的 HCU 使用率
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} ShowHCUUseResp "设备使用信息列表"
// @Failure 400 {object} error "错误信息"
// @Router /gpu/use [post]
func ShowHCUUse(c *gin.Context) {
	var dvIdList []int
	if err := c.ShouldBindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmInfos, err := dcgm.ShowHCUUse(dvIdList)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	// dcgm.DeviceUseInfo -> router.DeviceUseInfo
	infos := make([]DeviceUseInfo, 0, len(dcgmInfos))
	for _, info := range dcgmInfos {
		infos = append(infos, DeviceUseInfo{
			DeviceID:      info.DeviceID,
			GPUUsePercent: info.GPUUsePercent,
			Utilization:   info.Utilization,
		})
	}

	resp := ShowHCUUseResp{
		DeviceUseInfos: infos,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Power
// ShowEnergy 展示设备消耗的能量
// @Summary 展示设备的能量消耗
// @Description 获取并展示指定设备的能量消耗情况。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功返回设备的能量消耗信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /energy [post]
func ShowEnergy(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowEnergy(dvIdList)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Memory
// ShowMemInfo 展示设备的内存信息
// @Summary 展示设备内存信息
// @Description 获取并展示指定设备的内存使用情况，包括不同类型的内存。
// @Param dvIdList body []int true "设备 ID 列表"
// @Param memTypes body []string true "内存类型列表，如 'all' 或指定类型"
// @Success 200 {string} string "成功返回设备的内存信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /memory/info [post]
func ShowMemInfo(c *gin.Context) {
	var request struct {
		DvIdList []int    `json:"dvIdList"`
		MemTypes []string `json:"memTypes"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowMemInfo(request.DvIdList, request.MemTypes)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Memory
// ShowMemUse 展示设备的内存使用情况
// @Summary 展示设备内存使用情况
// @Description 获取并展示指定设备的当前内存使用百分比和其他相关的利用率数据。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功返回设备的内存使用信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /memory/use [post]
func ShowMemUse(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowMemUse(dvIdList)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Memory
// ShowMemVendor 展示设备内存供应商信息
// @Summary 展示设备的内存供应商信息
// @Description 获取并展示指定设备的内存供应商信息
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} ShowMemVendorResp "设备内存供应商信息"
// @Failure 400 {object} error "请求参数错误"
// @Router /memory/vendor [post]
func ShowMemVendor(c *gin.Context) {
	var dvIdList []int
	if err := c.ShouldBindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmInfos, err := dcgm.ShowMemVendor(dvIdList)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	// dcgm.DeviceMemVendorInfo -> router.DeviceMemVendorInfo
	infos := make([]DeviceMemVendorInfo, 0, len(dcgmInfos))
	for _, info := range dcgmInfos {
		infos = append(infos, DeviceMemVendorInfo{
			DeviceID: info.DeviceID,
			Vendor:   info.Vendor,
		})
	}

	resp := ShowMemVendorResp{
		DeviceMemVendorInfos: infos,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags PCIe
// ShowPcieBw 展示设备的 PCIe 带宽使用情况
// @Summary 展示设备的 PCIe 带宽使用情况
// @Description 获取并展示指定设备的 PCIe 带宽使用情况，包括发送和接收带宽
// @Accept json
// @Produce json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} ShowPcieBwResp "设备 PCIe 带宽信息"
// @Failure 400 {object} error "请求参数错误"
// @Router /pcie/bandwidth [post]
func ShowPcieBw(c *gin.Context) {
	var dvIdList []int
	if err := c.ShouldBindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmInfos, err := dcgm.ShowPcieBw(dvIdList)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	// dcgm.PcieBandwidthInfo -> router.PcieBandwidthInfo
	infos := make([]PcieBandwidthInfo, 0, len(dcgmInfos))
	for _, info := range dcgmInfos {
		infos = append(infos, PcieBandwidthInfo{
			DvInd:    info.DvInd,
			Sent:     info.Sent,
			Received: info.Received,
			Bw:       info.Bw,
		})
	}
	resp := ShowPcieBwResp{
		PcieBandwidthInfos: infos,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags PCIe
// @Summary 展示设备的 PCIe 重放计数
// @Description 获取并展示指定设备的 PCIe 重放计数
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} ShowPcieReplayCountResponse "设备的 PCIe 重放计数信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /pcie/replaycount [post]
func ShowPcieReplayCount(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	// 从 dcgm 层获取数据（业务层结构体）
	dcgmInfos, err := dcgm.ShowPcieReplayCount(dvIdList)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	// 转换为 router 层结构体，解决跨包类型不一致问题
	pcieReplayCountInfos := make([]PcieReplayCountInfo, 0, len(dcgmInfos))
	for _, info := range dcgmInfos {
		pcieReplayCountInfos = append(pcieReplayCountInfos, PcieReplayCountInfo{
			DeviceID: info.DeviceID,
			Count:    info.Count,
		})
	}

	c.JSON(http.StatusOK, SuccessResponse(pcieReplayCountInfos))
}

// @Tags System
// ShowPids 展示进程信息
// @Summary 展示系统中正在运行的KFD进程信息
// @Description 获取并展示当前系统中运行的KFD进程的详细信息。
// @Success 200 {string} string "成功返回进程信息"
// @Failure 400 {string} string "请求错误"
// @Router /pids [get]
func ShowPids(c *gin.Context) {
	err := dcgm.ShowPids()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Power
// GetDevicePower 展示设备的平均功率消耗
// @Summary 展示设备的平均功率消耗
// @Description 获取并展示指定设备的平均图形功率消耗
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} []DevicePowerInfo "设备的功率信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "服务器内部错误"
// @Router /device/power [post]
func GetDevicePower(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	// dcgm 层数据
	dcgmInfos, err := dcgm.ShowPower(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// 转换为 router 层结构体
	devicePowerInfos := make([]DevicePowerInfo, 0, len(dcgmInfos))
	for _, info := range dcgmInfos {
		devicePowerInfos = append(devicePowerInfos, DevicePowerInfo{
			DeviceID: info.DeviceID,
			Power:    info.Power,
		})
	}

	c.JSON(http.StatusOK, SuccessResponse(devicePowerInfos))
}

// @Tags Power
// @Summary 展示设备的HCU内存时钟频率和电压
// @Description 获取并展示指定设备的HCU内存时钟频率和电压表
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {array} DevicePowerPlayInfo "设备的HCU时钟频率和电压信息"
// @Failure 400 {object} error "请求参数错误"
// @Router /device/powerplay [post]
func GetDevicePowerPlayTable(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmInfos, err := dcgm.ShowPowerPlayTable(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	devicePowerPlayInfos := make([]DevicePowerPlayInfo, 0, len(dcgmInfos))
	for _, info := range dcgmInfos {
		devicePowerPlayInfos = append(devicePowerPlayInfos, DevicePowerPlayInfo{
			DeviceID:  info.DeviceID,
			SCLK:      info.SCLK,
			MCLK:      info.MCLK,
			DDC_CURVE: info.DDC_CURVE,
			RANGE:     info.RANGE,
		})
	}

	c.JSON(http.StatusOK, SuccessResponse(devicePowerPlayInfos))
}

// @Tags Device
// @Summary 显示设备的产品名称
// @Description 获取并显示指定设备的产品名称、供应商、系列、型号和 SKU 信息
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {array} DeviceProductInfo "设备的产品信息列表"
// @Failure 400 {object} error "请求参数错误"
// @Router /device/product [post]
func GetDeviceProductName(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmInfos, err := dcgm.ShowProductName(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	deviceProductInfos := make([]DeviceProductInfo, 0, len(dcgmInfos))
	for _, info := range dcgmInfos {
		deviceProductInfos = append(deviceProductInfos, DeviceProductInfo{
			DeviceID:   info.DeviceID,
			CardSeries: info.CardSeries,
			CardModel:  info.CardModel,
			CardVendor: info.CardVendor,
			CardSKU:    info.CardSKU,
		})
	}

	c.JSON(http.StatusOK, SuccessResponse(deviceProductInfos))
}

// @Tags Power
// @Summary 显示设备的电源配置文件
// @Description 获取并显示指定设备的电源配置文件，包括可用的电源配置选项
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {array} DeviceProfile "设备的电源配置文件信息列表"
// @Failure 400 {object} error "请求参数错误"
// @Router /device/profile [post]
func GetDeviceProfile(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmProfiles, err := dcgm.ShowProfile(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	deviceProfiles := make([]DeviceProfile, 0, len(dcgmProfiles))
	for _, p := range dcgmProfiles {
		deviceProfiles = append(deviceProfiles, DeviceProfile{
			DeviceID: p.DeviceID,
			Profiles: p.Profiles,
		})
	}

	c.JSON(http.StatusOK, SuccessResponse(deviceProfiles))
}

// @Tags Power
// GetDeviceRange 显示设备的电流或电压范围
// @Summary 显示设备的电流或电压范围（K100_AI卡不支持该操作）
// @Description 获取并显示指定设备的有效电流或电压范围
// @Param dvIdList body []int true "设备ID列表"
// @Param rangeType body string true "范围类型 (sclk, mclk, voltage)"
// @Success 200 {string} string "设备的电流或电压范围信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /device/range [post]
func GetDeviceRange(c *gin.Context) {
	var request struct {
		DvIdList  []int  `json:"dvIdList"`
		RangeType string `json:"rangeType"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowRange(request.DvIdList, request.RangeType)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Memory
// @Summary 显示设备的退役页信息
// @Description 获取并显示指定设备的退役内存页信息
// @Param dvIdList body []int true "设备ID列表"
// @Param retiredType body string false "退役类型 (默认为'all')"
// @Success 200 {string} string "设备的退役页信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /device/retiredpages [post]
func GetDeviceRetiredPages(c *gin.Context) {
	var request struct {
		DvIdList    []int  `json:"dvIdList"`
		RetiredType string `json:"retiredType"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	dcgm.ShowRetiredPages(request.DvIdList, request.RetiredType)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Device
// @Summary 显示设备的序列号
// @Description 获取并显示指定设备的序列号信息
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} DeviceSerialInfo "设备的序列号信息列表"
// @Failure 400 {object} error "请求参数错误"
// @Router /device/serialnumber [post]
func GetDeviceSerialNumber(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmSerialInfos, err := dcgm.ShowSerialNumber(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	serialInfos := make([]DeviceSerialInfo, 0, len(dcgmSerialInfos))
	for _, info := range dcgmSerialInfos {
		serialInfos = append(serialInfos, DeviceSerialInfo{
			DeviceID:     info.DeviceID,
			SerialNumber: info.SerialNumber,
		})
	}

	c.JSON(http.StatusOK, SuccessResponse(serialInfos))
}

// @Tags Device
// @Summary 显示设备的唯一ID
// @Description 获取并显示指定设备的唯一ID信息。
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} DeviceUIdInfo "设备的唯一ID信息列表"
// @Failure 400 {object} error "请求参数错误"
// @Router /showUId [post]
func ShowUId(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmUIdInfos, err := dcgm.ShowUId(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	uIdInfos := make([]DeviceUIdInfo, 0, len(dcgmUIdInfos))
	for _, info := range dcgmUIdInfos {
		uIdInfos = append(uIdInfos, DeviceUIdInfo{
			DeviceID: info.DeviceID,
			UId:      info.UId,
		})
	}

	c.JSON(http.StatusOK, SuccessResponse(uIdInfos))
}

// @Tags Version
// @Summary 显示设备的VBIOS版本
// @Description 获取并显示指定设备的VBIOS版本信息。
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} DeviceVBIOSInfo "设备的VBIOS版本信息列表"
// @Failure 400 {object} error "请求参数错误"
// @Router /showVbiosVersion [post]
func ShowVbiosVersion(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmVBIOSInfos, err := dcgm.ShowVbiosVersion(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	vbiosInfos := make([]DeviceVBIOSInfo, 0, len(dcgmVBIOSInfos))
	for _, info := range dcgmVBIOSInfos {
		vbiosInfos = append(vbiosInfos, DeviceVBIOSInfo{
			DeviceID: info.DeviceID,
			VBIOS:    info.VBIOS,
		})
	}

	c.JSON(http.StatusOK, SuccessResponse(vbiosInfos))
}

// @Tags Event
// ShowEvents 显示设备的事件
// @Summary 显示设备的事件
// @Description 获取并显示指定设备的事件信息。
// @Param dvIdList body []int true "设备ID列表"
// @Param eventTypes body []string true "事件类型列表"
// @Success 200 {string} string "成功返回设备的事件信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showEvents [post]
func ShowEvents(c *gin.Context) {
	var requestData struct {
		DvIdList   []int    `json:"dvIdList"`
		EventTypes []string `json:"eventTypes"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.String(http.StatusBadRequest, "请求参数错误")
		return
	}

	dcgm.ShowEvents(requestData.DvIdList, requestData.EventTypes)
	c.String(http.StatusOK, "设备事件信息已显示")
}

// @Tags Power
// @Summary 显示设备的电压信息
// @Description 获取并显示指定设备的当前电压信息。
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} DeviceVoltageInfo "设备的电压信息列表"
// @Failure 400 {object} error "请求参数错误"
// @Router /showVoltage [post]
func ShowVoltage(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmVoltageInfos, err := dcgm.ShowVoltage(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	voltageInfos := make([]DeviceVoltageInfo, 0, len(dcgmVoltageInfos))
	for _, info := range dcgmVoltageInfos {
		voltageInfos = append(voltageInfos, DeviceVoltageInfo{
			DeviceID: info.DeviceID,
			Voltage:  info.Voltage,
		})
	}

	c.JSON(http.StatusOK, SuccessResponse(voltageInfos))
}

// @Tags Power
// ShowVoltageCurve 显示设备的电压曲线点
// @Summary 显示设备的电压曲线点
// @Description 获取并显示指定设备的电压曲线点信息。
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {string} string "设备的电压曲线点信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showVoltageCurve [post]
func ShowVoltageCurve(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowVoltageCurve(dvIdList)
	c.String(http.StatusOK, "设备电压曲线点信息已显示")
}

// @Tags Error
// ShowHSLErr 显示 HSL 错误状态
// @Summary 显示 HSL 错误状态
// @Description 显示一组 GPU 设备的 HSL 错误状态。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "HSL 错误状态信息"
// @Router /showHSLErr [post]
func ShowHSLErr(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowHSLErr(dvIdList, true)
	c.String(http.StatusOK, "HSL 错误状态信息已显示")
}

// @Tags Topo
// ShowWeightTopology 显示 GPU 拓扑权重
// @Summary 显示 GPU 拓扑权重
// @Description 显示 GPU 设备之间的权重信息。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "GPU 拓扑权重信息"
// @Router /showWeightTopology [post]
func ShowWeightTopology(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	dcgm.ShowWeightTopology(dvIdList, true)
	c.String(http.StatusOK, "GPU 拓扑权重信息已显示")
}

// @Tags Topo
// ShowHopsTopology 显示 GPU 拓扑跳数
// @Summary 显示 GPU 拓扑跳数
// @Description 显示 GPU 设备之间的跳数信息。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "GPU 拓扑跳数信息"
// @Router /showHopsTopology [post]
func ShowHopsTopology(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowHopsTopology(dvIdList, true)
	c.String(http.StatusOK, "GPU 拓扑跳数信息已显示")
}

// @Tags Topo
// ShowTypeTopology 显示 GPU 拓扑中两台设备之间的链接类型。
// @Summary 显示 GPU 拓扑链接类型
// @Description 显示 GPU 设备之间的链接类型信息。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "GPU 拓扑链接类型信息"
// @Router /showTypeTopology [post]
func ShowTypeTopology(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}
	dcgm.ShowTypeTopology(dvIdList, true)
	c.String(http.StatusOK, " GPU拓扑中两台设备之间的链接类型已展示")
}

// @Tags Topo
// @Summary 显示 NUMA 节点信息
// @Description 显示一组 GPU 设备的 NUMA 节点和关联信息。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {array} NumaInfo "NUMA 节点信息列表"
// @Failure 400 {object} error "请求参数错误"
// @Router /showNumaTopology [post]
func ShowNumaTopology(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmInfos, err := dcgm.ShowNumaTopology(dvIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	numaInfos := make([]NumaInfo, 0, len(dcgmInfos))
	for _, info := range dcgmInfos {
		numaInfos = append(numaInfos, NumaInfo{
			DeviceID:     info.DeviceID,
			NumaNode:     info.NumaNode,
			NumaAffinity: info.NumaAffinity,
		})
	}

	c.JSON(http.StatusOK, SuccessResponse(numaInfos))
}

// @Tags Topo
// ShowHwTopology 显示指定设备的完整硬件拓扑信息。
// @Summary 显示完整的硬件拓扑信息
// @Description 显示一组 GPU 设备的权重、跳数、链接类型和 NUMA 节点信息。
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "完整的硬件拓扑信息"
// @Router /showHwTopology [post]
func ShowHwTopology(c *gin.Context) {
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgm.ShowHwTopology(dvIdList)
	c.String(http.StatusOK, "指定设备的完整硬件拓扑信息已显示")
}

// @Tags Topology
// @Summary 显示整机 HCU 互联矩阵信息
// @Description 枚举整机 HCU 的互联关系，包括链路类型（PCIe / HSL / HYSWITCH / NONE）及对应权重。
//
//	返回 HCU × HCU 的互联矩阵，可用于拓扑分析。
//
// @Success 200 {object} HCUInterconnectMatrix "HCU 互联矩阵信息"
// @Failure 500 {object} error "查询 HCU 互联信息失败"
// @Router /discoverInterconnectTopology [get]
func DiscoverInterconnectTopology(c *gin.Context) {
	matrix, err := dcgm.DiscoverInterconnectTopology()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	deviceMatrix := HCUInterconnectMatrix{
		DeviceCount: matrix.DeviceCount,
		Matrix:      make([][]HCULinkInfo, matrix.DeviceCount),
	}

	for i := 0; i < matrix.DeviceCount; i++ {
		deviceMatrix.Matrix[i] = make([]HCULinkInfo, matrix.DeviceCount)
		for j := 0; j < matrix.DeviceCount; j++ {
			link := matrix.Matrix[i][j]
			deviceMatrix.Matrix[i][j] = HCULinkInfo{
				SrcDvInd: link.SrcDvInd,
				DstDvInd: link.DstDvInd,
				PciID:    link.PciID,
				LinkType: link.LinkType,
				Weight:   link.Weight,
				Hops:     link.Hops,
			}
		}
	}

	c.JSON(http.StatusOK, SuccessResponse(deviceMatrix))
}

// @Tags Device
// @Summary 获取设备数量
// @Description 获取当前系统中的设备数量
// @Success 200 {object} DeviceCountInfo "设备数量信息"
// @Failure 500 {object} string "内部服务器错误"
// @Router /deviceCount [get]
func DeviceCount(c *gin.Context) {
	count, err := dcgm.DeviceCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(DeviceCountInfo{
		Count: count,
	}))
}

// @Tags VDevice
// @Summary 获取单个虚拟设备的信息
// @Description 根据设备索引获取对应的虚拟设备信息
// @Param vDvInd query int true "设备索引"
// @Success 200 {object} VDeviceInfo "虚拟设备信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /VDeviceSingleInfo [get]
func VDeviceSingleInfo(c *gin.Context) {
	vDvIndStr := c.Param("vDvInd") // 路由注册为 /:vDvInd，用 Param 而非 Query
	vDvInd, err := strconv.Atoi(vDvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "无效的设备索引")
		return
	}

	vDeviceInfo, err := dcgm.VDeviceSingleInfo(vDvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// 直接返回结构体，不用 map
	c.JSON(http.StatusOK, SuccessResponse(vDeviceInfo))
}

// @Tags VDevice
// @Summary 获取虚拟设备数量
// @Description 获取当前系统中的虚拟设备数量
// @Success 200 {object} VDeviceCountResp "虚拟设备数量"
// @Failure 500 {string} string "内部服务器错误"
// @Router /vDeviceCount [get]
func VDeviceCount(c *gin.Context) {
	vcount, err := dcgm.VDeviceCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	resp := VDeviceCountResp{
		VDeviceCount: vcount,
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// VDeviceInfos godoc
// @Summary      获取所有虚拟设备信息
// @Description  查询当前系统中所有虚拟设备（VDevice）的详细信息，包括所属物理设备、使用率及 PCI 信息
// @Tags         VDevice
// @Accept       json
// @Produce      json
// @Success      200 {object} VDeviceInfosResp
// @Failure      500 {object} Response
// @Router       /vdevice/infos [get]
func VDeviceInfos(c *gin.Context) {
	dcgmInfos, err := dcgm.VDeviceInfos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// 显式转换成 router 包的类型
	var infos []VDeviceInfo
	for _, info := range dcgmInfos {
		infos = append(infos, VDeviceInfo{
			Name:              info.Name,
			SubsystemTypeName: info.SubsystemTypeName,
			VComputeUnitCount: info.VComputeUnitCount,
			VMemoryTotal:      uint64(info.VMemoryTotal),
			VMemoryUsed:       uint64(info.VMemoryUsed),
			ContainerID:       info.ContainerID,
			DvInd:             info.DvInd,
			VPercent:          info.VPercent,
			VdvInd:            info.VdvInd,
			PciBusNumber:      info.PciBusNumber,
		})
	}
	resp := VDeviceInfosResp{
		VDeviceInfos: infos,
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags VDevice
// @Summary 获取设备剩余信息
// @Description 获取指定物理设备的剩余计算单元和内存信息
// @Param dvInd path int true "物理设备索引"
// @Success 200 {object} DeviceRemainingInfoResp "设备剩余信息"
// @Failure 400 {string} string "无效的设备索引"
// @Failure 500 {string} string "内部服务器错误"
// @Router /deviceRemainingInfo/{dvInd} [get]
func DeviceRemainingInfo(c *gin.Context) {
	dvIndStr := c.Param("dvInd")
	dvInd, err := strconv.Atoi(dvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "无效的设备索引")
		return
	}

	cus, memories, err := dcgm.DeviceRemainingInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	resp := DeviceRemainingInfoResp{
		CUs:      cus,
		Memories: memories,
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags VDevice
// @Summary 创建虚拟设备
// @Description 在指定的物理设备上创建指定数量的虚拟设备，返回创建的虚拟设备ID集合
// @Param dvInd query int true "物理设备的索引"
// @Param vDevCount query int true "要创建的虚拟设备数量"
// @Param vDevCUs query []int true "每个虚拟设备的计算单元数量，多个值使用多个 vDevCUs 参数传递，例如：vDevCUs=10&vDevCUs=20"
// @Param vDevMemSize query []int true "每个虚拟设备的内存大小，多个值使用多个 vDevMemSize 参数传递，例如：vDevMemSize=1024&vDevMemSize=2048"
// @Success 200 {object} CreateVDevicesResp "虚拟设备创建成功，返回虚拟设备ID集合"
// @Failure 400 {string} string "请求参数无效"
// @Failure 500 {string} string "创建虚拟设备失败"
// @Router /CreateVDevices [post]
func CreateVDevices(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Query("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "无效的物理设备索引")
		return
	}

	vDevCount, err := strconv.Atoi(c.Query("vDevCount"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "无效的虚拟设备数量")
		return
	}

	vDevCUsStr := c.QueryArray("vDevCUs")
	vDevMemSizeStr := c.QueryArray("vDevMemSize")

	vDevCUs := make([]int, len(vDevCUsStr))
	vDevMemSize := make([]int, len(vDevMemSizeStr))

	for i, v := range vDevCUsStr {
		vDevCUs[i], err = strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, "无效的计算单元数量")
			return
		}
	}

	for i, v := range vDevMemSizeStr {
		vDevMemSize[i], err = strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, "无效的内存大小")
			return
		}
	}

	vdevIDs, err := dcgm.CreateVDevices(dvInd, vDevCount, vDevCUs, vDevMemSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "创建虚拟设备失败")
		return
	}

	resp := CreateVDevicesResp{
		VDevIDs: vdevIDs,
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags VDevice
// DestroyVDevice 销毁指定物理设备上的所有虚拟设备
// @Summary 销毁所有虚拟设备
// @Description 销毁指定物理设备上的所有虚拟设备
// @Param dvInd query int true "物理设备的索引"
// @Success 200 {string} string "虚拟设备销毁成功"
// @Failure 400 {string} string "虚拟设备销毁失败"
// @Router /DestroyVDevice [delete]
func DestroyVDevice(c *gin.Context) {
	/*dvIndStr := c.Param("dvInd")
	dvInd, err := strconv.Atoi(dvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "无效的设备索引")
		return
	}
	if err := c.BindQuery(&dvInd); err != nil {
		c.JSON(http.StatusBadRequest, "虚拟设备销毁失败")
		return
	}*/
	// 获取查询参数 vDvInd，默认为空字符串
	dvIndStr := c.DefaultQuery("dvInd", "")
	if dvIndStr == "" {
		c.JSON(http.StatusBadRequest, "虚拟设备索引不能为空")
		return
	}

	// 将字符串转换为整数
	dvInd, err := strconv.Atoi(dvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "无效的虚拟设备索引")
		return
	}

	// 打印查询参数，确保 vDvInd 正确
	glog.Infof("Received dvInd: %d\n", dvInd)
	if err := dcgm.DestroyVDevice(dvInd); err != nil {
		c.JSON(http.StatusInternalServerError, "虚拟设备销毁失败:"+err.Error())
		return
	}
	c.JSON(http.StatusOK, "虚拟设备销毁成功")
}

// @Tags VDevice
// DestroySingleVDevice 销毁指定虚拟设备
// @Summary 销毁单个虚拟设备
// @Description 销毁指定索引的虚拟设备
// @Param vDvInd query int true "虚拟设备的索引"
// @Success 200 {string} string "虚拟设备销毁成功"
// @Failure 400 {string} string "虚拟设备销毁失败"
// @Router /DestroySingleVDevice [delete]
func DestroySingleVDevice(c *gin.Context) {
	// 获取查询参数 vDvInd，默认为空字符串
	vDvIndStr := c.DefaultQuery("vDvInd", "")
	if vDvIndStr == "" {
		c.JSON(http.StatusBadRequest, "虚拟设备索引不能为空")
		return
	}

	// 将字符串转换为整数
	vDvInd, err := strconv.Atoi(vDvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "无效的虚拟设备索引")
		return
	}

	// 打印查询参数，确保 vDvInd 正确
	glog.Infof("Received vDvInd: %d\n", vDvInd)

	if err := dcgm.DestroySingleVDevice(vDvInd); err != nil {
		c.JSON(http.StatusInternalServerError, "虚拟设备销毁失败:"+err.Error())
		return
	}
	c.JSON(http.StatusOK, "虚拟设备销毁成功")
}

// @Tags VDevice
// UpdateSingleVDevice 更新指定设备资源大小
// @Summary 更新虚拟设备资源
// @Description 更新指定虚拟设备的计算单元和内存大小
// @Param vDvInd query int true "虚拟设备的索引"
// @Param vDevCUs query int true "更新后的计算单元数量"
// @Param vDevMemSize query int true "更新后的内存大小"
// @Success 200 {string} string "虚拟设备更新成功"
// @Failure 400 {string} string "虚拟设备更新失败"
// @Router /UpdateSingleVDevice [put]
func UpdateSingleVDevice(c *gin.Context) {
	var vDvInd, vDevCUs, vDevMemSize int
	if err := c.ShouldBindQuery(&vDvInd); err != nil {
		c.JSON(http.StatusBadRequest, "虚拟设备更新失败")
		return
	}
	if err := c.ShouldBindQuery(&vDevCUs); err != nil {
		c.JSON(http.StatusBadRequest, "虚拟设备更新失败")
		return
	}
	if err := c.ShouldBindQuery(&vDevMemSize); err != nil {
		c.JSON(http.StatusBadRequest, "虚拟设备更新失败")
		return
	}
	if err := dcgm.UpdateSingleVDevice(vDvInd, vDevCUs, vDevMemSize); err != nil {
		c.JSON(http.StatusInternalServerError, "虚拟设备更新失败")
		return
	}
	c.JSON(http.StatusOK, "虚拟设备更新成功")
}

// @Tags VDevice
// 启动虚拟设备
// @Summary 启动指定的虚拟设备
// @Description 启动虚拟设备，指定设备索引
// @Param vDvInd path int true "虚拟设备索引"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /StartVDevice/{vDvInd} [get]
func StartVDevice(c *gin.Context) {
	vDvInd, err := strconv.Atoi(c.Param("vDvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "请求参数错误")
		return
	}

	if err := dcgm.StartVDevice(vDvInd); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "启动成功")
}

// @Tags VDevice
// 停止虚拟设备
// @Summary 停止指定的虚拟设备
// @Description 停止虚拟设备，指定设备索引
// @Param vDvInd path int true "虚拟设备索引"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /StopVDevice/{vDvInd} [get]
func StopVDevice(c *gin.Context) {
	vDvInd, err := strconv.Atoi(c.Param("vDvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "请求参数错误")
		return
	}

	if err := dcgm.StopVDevice(vDvInd); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "停止成功")
}

// @Tags VDevice
// 设置虚拟机加密状态
// @Summary 设置虚拟机加密状态
// @Description 根据提供的状态开启或关闭虚拟机加密
// @Param status query bool true "加密状态"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /SetEncryptionVMStatus [post]
func SetEncryptionVMStatus(c *gin.Context) {
	var status bool
	if err := c.BindQuery(&status); err != nil {
		c.JSON(http.StatusBadRequest, "请求参数错误")
		return
	}

	if err := dcgm.SetEncryptionVMStatus(status); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "设置成功")
}

// @Tags VDevice
// @Summary 获取当前虚拟机的加密状态
// @Description 返回虚拟机是否处于加密状态
// @Success 200 {object} EncryptionVMStatusResp "加密状态"
// @Failure 400 {string} string "操作失败"
// @Router /EncryptionVMStatus [get]
func EncryptionVMStatus(c *gin.Context) {
	status, err := dcgm.EncryptionVMStatus()
	if err != nil {
		c.JSON(http.StatusBadRequest, "操作失败")
		return
	}

	resp := EncryptionVMStatusResp{
		Status: status,
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Event
// 打印事件列表
// @Summary 打印设备的事件列表
// @Description 打印指定设备的事件列表，并设置延迟
// @Param device path int true "设备索引"
// @Param delay query int true "延迟时间（秒）"
// @Param eventList query []string true "事件列表"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /PrintEventList/{device} [get]
func PrintEventList(c *gin.Context) {
	device, err := strconv.Atoi(c.Param("device"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "请求参数错误")
		return
	}

	var delay int
	var eventList []string
	if err := c.BindQuery(&delay); err != nil || c.BindQuery(&eventList) != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	dcgm.PrintEventList(device, delay, eventList)
	c.JSON(http.StatusOK, "操作成功")
}

// @Tags Device
// @Summary 获取设备信息
// @Description 根据设备索引获取对应的设备信息
// @Param dvInd path int true "设备索引"
// @Success 200 {object} DeviceInfoResp "设备信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /device/info/{dvInd} [get]
func GetDeviceInfo(c *gin.Context) {
	dvIndStr := c.Param("dvInd")
	dvInd, err := strconv.Atoi(dvIndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "请求参数错误")
		return
	}

	dcgmInfo, err := dcgm.GetDeviceInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "内部服务器错误")
		return
	}

	// 👇 dcgm -> router 显式转换
	deviceInfo := DMIDeviceInfo{
		Name:             dcgmInfo.Name,
		ComputeUnitCount: dcgmInfo.ComputeUnitCount,
		DeviceID:         dcgmInfo.DeviceID,
		Percent:          dcgmInfo.Percent,
		MaxVDeviceCount:  dcgmInfo.MaxVDeviceCount,
	}

	resp := DeviceInfoResp{
		DeviceInfo: deviceInfo,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Device
// DeviceControl 处理设备控制
// @Summary 控制设备的性能级别、时钟频率和风扇重置
// @Description 根据传入的设备控制信息，设置设备的性能级别、时钟频率，并可选择性重置风扇
// @Accept json
// @Produce json
// @Param deviceControl body DeviceControlInfo true "设备控制信息"
// @Success 200 {object} DeviceControlResp "操作成功"
// @Failure 400 {object} Response "无效的请求参数或操作失败"
// @Failure 500 {object} Response "内部服务器错误"
// @Router /device/control [post]
func DeviceControl(c *gin.Context) {
	var deviceInfo DeviceControlInfo
	var validationErrors []string

	if err := c.ShouldBindJSON(&deviceInfo); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	// 参数校验
	dvInd := deviceInfo.DvInd
	if dvInd < 0 {
		validationErrors = append(validationErrors, "无效的 DvInd")
	}

	var levelConverted dcgm.DevPerfLevel
	if deviceInfo.PerfLevel != "" {
		var err error
		levelConverted, err = ConvertToRSMIDevPerfLevel(deviceInfo.PerfLevel)
		if err != nil {
			validationErrors = append(validationErrors, "无效的性能级别 PerfLevel")
		}
	}

	var sclkClock int64
	if deviceInfo.SclkClock != "" {
		var err error
		sclkClock, err = ConvertFrequencyToSclkClock(deviceInfo.SclkClock)
		if err != nil {
			validationErrors = append(validationErrors, "无效的 SclkClock："+err.Error())
		}
	}

	var socclkClock int64
	if deviceInfo.SocclkClock != "" {
		var err error
		socclkClock, err = ConvertFrequencyToSocclkClock(deviceInfo.SocclkClock)
		if err != nil {
			validationErrors = append(validationErrors, "无效的 SocclkClock："+err.Error())
		}
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest,
			ErrorResponse("参数校验失败: "+strings.Join(validationErrors, ", ")),
		)
		return
	}

	// 执行阶段
	var executionErrors []string

	if deviceInfo.PerfLevel != "" {
		if err := dcgm.DevPerfLevelSet(dvInd, levelConverted); err != nil {
			executionErrors = append(executionErrors, "设置性能级别失败: "+err.Error())
		}
	}

	if deviceInfo.SclkClock != "" {
		if err := dcgm.DevGpuClkFreqSet(dvInd, dcgm.RSMI_CLK_TYPE_SYS, sclkClock); err != nil {
			executionErrors = append(executionErrors, "设置 SCLK 失败: "+err.Error())
		}
	}

	if deviceInfo.SocclkClock != "" {
		if err := dcgm.DevGpuClkFreqSet(dvInd, dcgm.RSMI_CLK_TYPE_SOC, socclkClock); err != nil {
			executionErrors = append(executionErrors, "设置 SOCCLK 失败: "+err.Error())
		}
	}

	if deviceInfo.ResetFan {
		if err := dcgm.ResetFans([]int{dvInd}); err != nil {
			executionErrors = append(executionErrors, "重置风扇失败: "+err.Error())
		}
	}

	if len(executionErrors) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse(DeviceControlResp{
			Success: false,
			Errors:  executionErrors,
		}))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(DeviceControlResp{
		Success: true,
	}))
}

// @Tags Error
// @Summary 获取 ECC block 信息
// @Description 根据设备索引获取 ECC block 信息
// @Accept json
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {array} BlocksInfo "ECC block 信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /EccBlocksInfo [get]
func EccBlocksInfo(c *gin.Context) {
	var dvInd int
	if err := c.BindQuery(&dvInd); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("参数无效"))
		return
	}

	// dcgm 层返回的是 []dcgm.BlocksInfo
	dcgmInfos, err := dcgm.EccBlocksInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// 转换为 router 层结构体，避免类型冲突
	resp := make([]BlocksInfo, 0, len(dcgmInfos))
	for _, info := range dcgmInfos {
		resp = append(resp, BlocksInfo{
			Block: info.Block,
			State: info.State,
			CE:    info.CE,
			UE:    info.UE,
		})
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Device
// RunDiag 诊断设备
// @Summary 运行设备诊断
// @Description 运行设备的诊断工具，返回诊断结果
// @Produce json
// @Param level query int true "诊断级别"
// @Success 200 {object} DiagResults "返回诊断结果"
// @Failure 400 {string} string "参数无效"
// @Failure 500 {object} error "服务器内部错误"
// @Router /RunDiag [get]
func RunDiag(c *gin.Context) {
	levelStr := c.Param("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("参数无效"))
		return
	}

	// dcgm 层结果
	dcgmResult, err := dcgm.RunDiag(level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// ===== 显式转换 dcgm -> router =====
	resp := DiagResults{
		DeviceNumber: dcgmResult.DeviceNumber,
	}

	// Software 结果
	for _, s := range dcgmResult.Software {
		resp.Software = append(resp.Software, DiagResult{
			Status:       s.Status,
			TestName:     s.TestName,
			TestOutput:   s.TestOutput,
			ErrorCode:    s.ErrorCode,
			ErrorMessage: s.ErrorMessage,
		})
	}

	// PerHCU 结果
	for _, hcu := range dcgmResult.PerHCU {
		hcuResult := HCUResult{
			HCU: hcu.HCU,
			RC:  hcu.RC,
		}

		for _, r := range hcu.DiagResults {
			hcuResult.DiagResults = append(hcuResult.DiagResults, DiagResult{
				Status:       r.Status,
				TestName:     r.TestName,
				TestOutput:   r.TestOutput,
				ErrorCode:    r.ErrorCode,
				ErrorMessage: r.ErrorMessage,
			})
		}

		resp.PerHCU = append(resp.PerHCU, hcuResult)
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// startDiag godoc
// @Summary      启动异步诊断任务
// @Description  启动指定级别（1~4）的异步诊断任务，立即返回 jobId；同一时间只允许一个诊断任务运行
// @Tags         Diag
// @Accept       json
// @Produce      json
// @Param        level path int true "诊断级别（1~4）"
// @Success      200 {object} Response
// @Failure      400 {object} Response
// @Failure      409 {object} Response "已有诊断任务正在运行"
// @Failure      500 {object} Response "服务器内部错误"
// @Router       /diag/start/{level} [post]
func startDiag(c *gin.Context) {
	levelStr := c.Param("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil || level < dcgm.DiagLevel1 || level > dcgm.DiagLevel4 {
		c.JSON(http.StatusBadRequest, ErrorResponse("invalid level; must be 1..4"))
		return
	}

	// 如果当前已有正在运行的 job，立即拒绝
	if v := currentJobID.Load(); v != nil {
		if s, ok := v.(string); ok && s != "" {
			c.JSON(http.StatusConflict, ErrorResponse("another job is running"))
			return
		}
	}

	// 生成 job id
	id := fmt.Sprintf("job-%d-%d", time.Now().UnixNano(), atomic.AddUint64(&jobIDCounter, 1))
	job := &Job{
		ID:     id,
		Level:  level,
		Status: JobStatusPending,
	}

	jobStoreMu.Lock()
	jobStore[id] = job
	jobStoreMu.Unlock()

	currentJobID.Store(id)

	go func(j *Job) {
		jobStoreMu.Lock()
		j.Status = JobStatusRunning
		j.StartedAt = time.Now()
		jobStoreMu.Unlock()

		res, rErr := dcgm.RunDiag(j.Level)

		jobStoreMu.Lock()
		if rErr != nil {
			j.Status = JobStatusFailed
			j.ErrorMessage = rErr.Error()
		} else {
			if dcgm.IsDiagStopped() {
				j.Status = JobStatusCanceled
			} else {
				j.Status = JobStatusDone
			}
		}
		j.EndedAt = time.Now()

		var persistJob Job
		persistJob.ID = j.ID
		persistJob.Level = j.Level
		persistJob.Status = j.Status
		persistJob.ErrorMessage = j.ErrorMessage
		persistJob.StartedAt = j.StartedAt
		persistJob.EndedAt = j.EndedAt
		persistJob.Result = convertDiagResults(&res)
		jobStoreMu.Unlock()

		if err := saveJobToFile(&persistJob); err != nil {
			fmt.Printf("saveJobToFile failed for %s: %v\n", j.ID, err)
		}

		jobStoreMu.Lock()
		delete(jobStore, j.ID)
		jobStoreMu.Unlock()

		if v := currentJobID.Load(); v != nil {
			if s, ok := v.(string); ok && s == j.ID {
				currentJobID.Store("")
			}
		}
	}(job)

	c.JSON(http.StatusAccepted, SuccessResponse(DiagResp{
		JobID:     id,
		Status:    job.Status,
		StatusURL: fmt.Sprintf("/diag/status/%s", id),
	}))
}

// statusDiag godoc
// @Summary      查询诊断作业状态
// @Description  根据 job id 查询诊断作业当前状态及结果。
// @Description  优先从内存中查询运行中或刚完成的作业；
// @Description  若内存中不存在，则从磁盘 logs/<job-id>.json 中读取历史结果。
// @Tags         Diag
// @Accept       json
// @Produce      json
// @Param        id path string true "作业 ID"
// @Success      200 {object} Job "返回作业状态及结果（若已完成）"
// @Failure      404 {object} Response "作业不存在"
// @Router       /diag/status/{id} [get]
func statusDiag(c *gin.Context) {
	id := c.Param("id")

	// 1) 先查内存（运行中的 job）
	jobStoreMu.Lock()
	j, ok := jobStore[id]
	jobStoreMu.Unlock()
	if ok {
		c.JSON(http.StatusOK, j)
		return
	}

	// 2) 尝试从磁盘读（logs/<id>.json）
	if persisted, err := loadJobFromFile(id); err == nil {
		c.JSON(http.StatusOK, persisted)
		return
	}

	// 3) 都没有，404
	c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
}

// StopDiag godoc
// @Summary      停止当前诊断作业
// @Description  请求停止当前正在运行的诊断任务（非强制杀死）。
// @Description  后台任务会在安全检查点检测到停止标志后终止后续步骤。
// @Tags         Diag
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "停止请求已接收"
// @Router       /diag/stop [post]
func StopDiag(c *gin.Context) {
	dcgm.StopDiag()
	// 快速反馈：把当前 running job 标记为 canceling（若存在）
	if v := currentJobID.Load(); v != nil {
		if id, ok := v.(string); ok && id != "" {
			jobStoreMu.Lock()
			if j, exists := jobStore[id]; exists && j.Status == JobStatusRunning {
				// 仅给出快速反馈，后台 goroutine 最终会把状态改为 canceled/failed/done
				j.Status = JobStatusCanceled
			}
			jobStoreMu.Unlock()
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "stop requested"})
}

// startBandwidth godoc
// @Summary      启动 Memory Bandwidth 压测
// @Description  异步启动内存带宽测试，立即返回 job id。
// @Description  测试结果会被封装为 DiagResults 并持久化到 logs/<job-id>.json。
// @Tags         Diag
// @Accept       json
// @Produce      json
// @Success      200 {object} DiagResp
// @Failure      409 {object} Response "已有作业正在运行"
// @Router       /bandwidth/start [post]
func startBandwidth(c *gin.Context) {
	// 并发保护：如果已有正在运行的 job，返回 409
	if v := currentJobID.Load(); v != nil {
		if s, ok := v.(string); ok && s != "" {
			c.JSON(http.StatusConflict, gin.H{
				"error":      "another job is running",
				"running_id": s,
			})
			return
		}
	}

	// 生成 job id
	id := fmt.Sprintf("bw-%d-%d", time.Now().UnixNano(), atomic.AddUint64(&jobIDCounter, 1))
	job := &Job{
		ID:     id,
		Level:  0,
		Status: JobStatusPending,
	}

	// 存到内存并标记为当前 job
	jobStoreMu.Lock()
	jobStore[id] = job
	jobStoreMu.Unlock()
	currentJobID.Store(id)

	// 后台执行测试
	go func(j *Job) {
		// 标记 running
		jobStoreMu.Lock()
		j.Status = JobStatusRunning
		j.StartedAt = time.Now()
		jobStoreMu.Unlock()

		// 构造 dvIdList（所有设备）
		numDevices, _ := dcgm.NumMonitorDevices()
		dvIdList := make([]int, 0, numDevices)
		for i := 0; i < numDevices; i++ {
			dvIdList = append(dvIdList, i)
		}

		// 调用带返回值的 API（不会向 stdout 打印）
		bwMap, err := dcgm.BandwidthTestResult(dvIdList)

		// 更新内存 job 状态与时间
		jobStoreMu.Lock()
		if err != nil {
			j.Status = JobStatusFailed
			j.ErrorMessage = err.Error()
		} else {
			// 如果全局 stop 标志被设置，则视为用户取消
			if dcgm.IsDiagStopped() {
				j.Status = JobStatusCanceled
			} else {
				j.Status = JobStatusDone
			}
		}
		j.EndedAt = time.Now()
		// 构造持久化副本
		var persist Job
		persist.ID = j.ID
		persist.Level = j.Level
		persist.Status = j.Status
		persist.ErrorMessage = j.ErrorMessage
		persist.StartedAt = j.StartedAt
		persist.EndedAt = j.EndedAt
		// 构造 dcgm.DiagResults 并填充
		var diag dcgm.DiagResults
		diag.DeviceNumber = 0
		if bwMap != nil {
			diag.DeviceNumber = len(bwMap)
			// 保证顺序：按 HCU 索引遍历
			keys := make([]int, 0, len(bwMap))
			for k := range bwMap {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			for _, d := range keys {
				bw := bwMap[d]
				status := DiagResultPass
				if bw <= 0 {
					status = DiagResultWarn
				}
				dr := dcgm.DiagResult{
					Status:       status,
					TestName:     "Memory Bandwidth",
					TestOutput:   fmt.Sprintf("Bandwidth: %.3f GB/s", bw),
					ErrorCode:    0,
					ErrorMessage: "",
				}
				diag.PerHCU = append(diag.PerHCU, dcgm.HCUResult{
					HCU:         d,
					RC:          0,
					DiagResults: []dcgm.DiagResult{dr},
				})
			}
			diag.Software = append(diag.Software, dcgm.DiagResult{
				Status:     DiagResultPass,
				TestName:   "Memory Bandwidth Summary",
				TestOutput: fmt.Sprintf("Parsed %d entries", len(bwMap)),
				ErrorCode:  0,
			})
		} else {
			diag.Software = append(diag.Software, dcgm.DiagResult{
				Status:       DiagResultFail,
				TestName:     "Memory Bandwidth",
				TestOutput:   "",
				ErrorCode:    -1,
				ErrorMessage: fmt.Sprintf("BandwidthTestResult returned nil and no error"),
			})
		}
		persist.Result = convertDiagResults(&diag)

		jobStoreMu.Unlock()

		// 持久化到文件 logs/<job-id>.json
		if err := saveJobToFile(&persist); err != nil {
			fmt.Printf("saveJobToFile failed for bandwidth job %s: %v\n", j.ID, err)
		}

		// 从内存中删除该 job（只保留磁盘持久化）
		jobStoreMu.Lock()
		delete(jobStore, j.ID)
		jobStoreMu.Unlock()

		// 清理 currentJobID（仅当仍指向本 job）
		if v := currentJobID.Load(); v != nil {
			if s, ok := v.(string); ok && s == j.ID {
				currentJobID.Store("")
			}
		}
	}(job)

	// 立即返回 job id
	c.JSON(http.StatusAccepted, SuccessResponse(DiagResp{
		JobID:     id,
		Status:    job.Status,
		StatusURL: fmt.Sprintf("/diag/status/%s", id),
	}))
}

// startPcie godoc
// @Summary      启动 PCIe 带宽压测
// @Description  异步启动 PCIe 带宽测试（Sys<->Fb）。
// @Description  返回 job id，可通过 /diag/status/{id} 查询结果。
// @Tags         Diag
// @Accept       json
// @Produce      json
// @Success      200 {object} DiagResp
// @Failure      409 {object} Response "已有作业正在运行"
// @Router       /pcie/start [post]
func startPcie(c *gin.Context) {
	// 并发保护：若已有运行 job，返回 409
	if v := currentJobID.Load(); v != nil {
		if s, ok := v.(string); ok && s != "" {
			c.JSON(http.StatusConflict, gin.H{
				"error":      "another job is running",
				"running_id": s,
			})
			return
		}
	}

	// 生成 job id
	id := fmt.Sprintf("pcie-%d-%d", time.Now().UnixNano(), atomic.AddUint64(&jobIDCounter, 1))
	job := &Job{
		ID:     id,
		Level:  0,
		Status: JobStatusPending,
	}

	jobStoreMu.Lock()
	jobStore[id] = job
	jobStoreMu.Unlock()
	currentJobID.Store(id)

	// 后台执行
	go func(j *Job) {
		jobStoreMu.Lock()
		j.Status = JobStatusRunning
		j.StartedAt = time.Now()
		jobStoreMu.Unlock()

		// 调用 PCIe 带返回值 API（不打印）
		pcieRes, err := dcgm.PcieBandwidthTestResult()

		// 更新内存 job 状态
		jobStoreMu.Lock()
		if err != nil {
			j.Status = JobStatusFailed
			j.ErrorMessage = err.Error()
		} else {
			if dcgm.IsDiagStopped() {
				j.Status = JobStatusCanceled
			} else {
				j.Status = JobStatusDone
			}
		}
		j.EndedAt = time.Now()

		// 构造持久化 Job 副本
		var persist Job
		persist.ID = j.ID
		persist.Level = j.Level
		persist.Status = j.Status
		persist.ErrorMessage = j.ErrorMessage
		persist.StartedAt = j.StartedAt
		persist.EndedAt = j.EndedAt

		// 构造 dcgm.DiagResults 并填充 PCIe 数据
		var diag dcgm.DiagResults
		diag.DeviceNumber = pcieRes.DeviceCount

		// 假设 pcieRes.HCUs 是 []struct{ DvInd int; SysToFb, FbToSys float64 }
		for _, h := range pcieRes.HCUs {
			hcuIdx := h.DvInd
			sysToFb := h.SysToFb
			fbToSys := h.FbToSys

			status := DiagResultPass
			if sysToFb <= 0 || fbToSys <= 0 {
				status = DiagResultWarn
			}

			dr := dcgm.DiagResult{
				Status:       status,
				TestName:     "PCIe Bandwidth",
				TestOutput:   fmt.Sprintf("Sys->Fb: %.3f GB/s, Fb->Sys: %.3f GB/s", sysToFb, fbToSys),
				ErrorCode:    0,
				ErrorMessage: "",
			}
			diag.PerHCU = append(diag.PerHCU, dcgm.HCUResult{
				HCU:         hcuIdx,
				RC:          0,
				DiagResults: []dcgm.DiagResult{dr},
			})
		}

		// 将一个 summary 放入 Software 段（便于快速查看）
		diag.Software = append(diag.Software, dcgm.DiagResult{
			Status:     DiagResultPass,
			TestName:   "PCIe Bandwidth Summary",
			TestOutput: fmt.Sprintf("Parsed devices: %d, log=%s", pcieRes.DeviceCount, pcieRes.LogFile),
			ErrorCode:  0,
		})

		persist.Result = convertDiagResults(&diag)
		jobStoreMu.Unlock()

		// 写盘
		if err := saveJobToFile(&persist); err != nil {
			fmt.Printf("saveJobToFile failed for pcie job %s: %v\n", j.ID, err)
		}

		// 从内存删除 job，保留磁盘
		jobStoreMu.Lock()
		delete(jobStore, j.ID)
		jobStoreMu.Unlock()

		// 清理 currentJobID（仅当仍指向本 job）
		if v := currentJobID.Load(); v != nil {
			if s, ok := v.(string); ok && s == j.ID {
				currentJobID.Store("")
			}
		}
	}(job)

	// 立即返回
	c.JSON(http.StatusAccepted, SuccessResponse(DiagResp{
		JobID:     id,
		Status:    job.Status,
		StatusURL: fmt.Sprintf("/diag/status/%s", id),
	}))
}

// startXHCL godoc
// @Summary      启动 XHCL 带宽测试
// @Description  异步启动 XHCL（跨 HCU）互联带宽测试。
// @Description  测试结果按 HCU 聚合并持久化。
// @Tags         Diag
// @Accept       json
// @Produce      json
// @Success      200 {object} DiagResp
// @Failure      409 {object} Response "已有作业正在运行"
// @Router       /xhcl/start [post]
func startXHCL(c *gin.Context) {
	// 若已有运行 job，则拒绝（避免并行压测）
	if v := currentJobID.Load(); v != nil {
		if s, ok := v.(string); ok && s != "" {
			c.JSON(http.StatusConflict, gin.H{
				"error":      "another job is running",
				"running_id": s,
			})
			return
		}
	}

	// 生成 job id（前缀 xhcl-）
	id := fmt.Sprintf("xhcl-%d-%d", time.Now().UnixNano(), atomic.AddUint64(&jobIDCounter, 1))
	job := &Job{
		ID:     id,
		Level:  0,
		Status: JobStatusPending,
	}

	// 保存到内存并标记为当前 job
	jobStoreMu.Lock()
	jobStore[id] = job
	jobStoreMu.Unlock()
	currentJobID.Store(id)

	// 在后台执行具体测试并持久化结果
	go func(j *Job) {
		// 标记为 running
		jobStoreMu.Lock()
		j.Status = JobStatusRunning
		j.StartedAt = time.Now()
		jobStoreMu.Unlock()

		// 调用 dcgm API（不打印）
		xhclPairs, err := dcgm.XHCLTestResult()

		// 更新 job 状态与时间
		jobStoreMu.Lock()
		if err != nil {
			j.Status = JobStatusFailed
			j.ErrorMessage = err.Error()
		} else {
			if dcgm.IsDiagStopped() {
				j.Status = JobStatusCanceled
			} else {
				j.Status = JobStatusDone
			}
		}
		j.EndedAt = time.Now()

		// 构造要持久化的 Job 副本
		var persist Job
		persist.ID = j.ID
		persist.Level = j.Level
		persist.Status = j.Status
		persist.ErrorMessage = j.ErrorMessage
		persist.StartedAt = j.StartedAt
		persist.EndedAt = j.EndedAt

		// 把 xhclPairs 转成 dcgm.DiagResults（PerHCU 和 Software）
		var diag dcgm.DiagResults

		// 如果成功拿到结果，按 pair 将信息写入 PerHCU
		if err == nil && xhclPairs != nil {
			// 对于每个 pair (src->dst, bw)，把一条记录写入 src 的 PerHCU，一条写入 dst 的 PerHCU
			for _, p := range xhclPairs {
				src := p.SrcHCUId
				dst := p.DstHCUId
				bw := p.BandwidthGBs

				status := DiagResultPass
				if bw <= 0 {
					status = DiagResultWarn
				}

				srcDiag := dcgm.DiagResult{
					Status:       status,
					TestName:     fmt.Sprintf("XHCL to HCU%d", dst),
					TestOutput:   fmt.Sprintf("HCU%d <-> HCU%d XHCL: %.3f GB/s", src, dst, bw),
					ErrorCode:    0,
					ErrorMessage: "",
				}
				dstDiag := dcgm.DiagResult{
					Status:       status,
					TestName:     fmt.Sprintf("XHCL to HCU%d", src),
					TestOutput:   fmt.Sprintf("HCU%d <-> HCU%d XHCL: %.3f GB/s", dst, src, bw),
					ErrorCode:    0,
					ErrorMessage: "",
				}

				// 合并到 src 的 PerHCU
				found := false
				for i := range diag.PerHCU {
					if diag.PerHCU[i].HCU == src {
						diag.PerHCU[i].DiagResults = append(diag.PerHCU[i].DiagResults, srcDiag)
						found = true
						break
					}
				}
				if !found {
					diag.PerHCU = append(diag.PerHCU, dcgm.HCUResult{
						HCU:         src,
						RC:          0,
						DiagResults: []dcgm.DiagResult{srcDiag},
					})
				}

				// 合并到 dst 的 PerHCU
				found = false
				for i := range diag.PerHCU {
					if diag.PerHCU[i].HCU == dst {
						diag.PerHCU[i].DiagResults = append(diag.PerHCU[i].DiagResults, dstDiag)
						found = true
						break
					}
				}
				if !found {
					diag.PerHCU = append(diag.PerHCU, dcgm.HCUResult{
						HCU:         dst,
						RC:          0,
						DiagResults: []dcgm.DiagResult{dstDiag},
					})
				}
			}

			diag.Software = append(diag.Software, dcgm.DiagResult{
				Status:     DiagResultPass,
				TestName:   "XHCL Summary",
				TestOutput: fmt.Sprintf("Parsed %d XHCL pairs", len(xhclPairs)),
				ErrorCode:  0,
			})
		} else {
			// 如果没有得到结果或出错，写一条 Software 失败条目以便排查
			diag.Software = append(diag.Software, dcgm.DiagResult{
				Status:       DiagResultFail,
				TestName:     "XHCL Test",
				TestOutput:   "",
				ErrorCode:    -1,
				ErrorMessage: fmt.Sprintf("XHCLTestResult error: %v", err),
			})
		}

		persist.Result = convertDiagResults(&diag)
		jobStoreMu.Unlock()

		// 写盘 logs/<job-id>.json
		if err := saveJobToFile(&persist); err != nil {
			fmt.Printf("saveJobToFile failed for xhcl job %s: %v\n", j.ID, err)
		}

		// 从内存移除 job（持久化后）
		jobStoreMu.Lock()
		delete(jobStore, j.ID)
		jobStoreMu.Unlock()

		// 清理 currentJobID（仅当仍指向本 job）
		if v := currentJobID.Load(); v != nil {
			if s, ok := v.(string); ok && s == j.ID {
				currentJobID.Store("")
			}
		}
	}(job)

	// 立即返回 job id
	c.JSON(http.StatusAccepted, SuccessResponse(DiagResp{
		JobID:     id,
		Status:    job.Status,
		StatusURL: fmt.Sprintf("/diag/status/%s", id),
	}))
}

// startTargetStress godoc
// @Summary      启动 TargetStress（GEMM）压测
// @Description  异步执行多种 GEMM 模式的 TargetStress 压测。
// @Description  返回 job id，结果通过状态接口获取。
// @Tags         Diag
// @Accept       json
// @Produce      json
// @Success      200 {object} DiagResp
// @Failure      409 {object} Response "已有作业正在运行"
// @Router       /targetstress/start [post]
func startTargetStress(c *gin.Context) {
	// 如果已有正在运行的 job，则拒绝（避免并发压测）
	if v := currentJobID.Load(); v != nil {
		if s, ok := v.(string); ok && s != "" {
			c.JSON(http.StatusConflict, gin.H{
				"error":      "another job is running",
				"running_id": s,
			})
			return
		}
	}

	// 生成 job id
	id := fmt.Sprintf("target-%d-%d", time.Now().UnixNano(), atomic.AddUint64(&jobIDCounter, 1))
	job := &Job{
		ID:     id,
		Level:  0,
		Status: JobStatusPending,
	}

	// 保存到内存 store 并标记为当前运行 job
	jobStoreMu.Lock()
	jobStore[id] = job
	jobStoreMu.Unlock()
	currentJobID.Store(id)

	// 后台执行实际测试
	go func(j *Job) {
		// 标记为 running
		jobStoreMu.Lock()
		j.Status = JobStatusRunning
		j.StartedAt = time.Now()
		jobStoreMu.Unlock()

		// 调用 dcgm 的 API（应为不会打印的版本）
		tsRes, err := dcgm.TargetStressTestResult()

		// 更新 job 状态与结束时间（先在内存中更新）
		jobStoreMu.Lock()
		if err != nil {
			j.Status = JobStatusFailed
			j.ErrorMessage = err.Error()
		} else {
			if dcgm.IsDiagStopped() {
				j.Status = JobStatusCanceled
			} else {
				j.Status = JobStatusDone
			}
		}
		j.EndedAt = time.Now()

		// 构造持久化副本（避免并发引用问题）
		var persist Job
		persist.ID = j.ID
		persist.Level = j.Level
		persist.Status = j.Status
		persist.ErrorMessage = j.ErrorMessage
		persist.StartedAt = j.StartedAt
		persist.EndedAt = j.EndedAt

		// 把 TargetStressResult 转换为 dcgm.DiagResults
		var diag dcgm.DiagResults

		if err == nil && tsRes.Results != nil {
			// 收集 HCU 列表并按 HCU 聚合结果
			// 使用 map[HCU] -> []dcgm.DiagResult，然后转成 diag.PerHCU（保证稳定顺序）
			perMap := make(map[int][]dcgm.DiagResult)
			for _, r := range tsRes.Results {
				status := DiagResultPass
				if r.Failed || r.Mean <= 0 {
					status = DiagResultWarn
				}
				dr := dcgm.DiagResult{
					Status:       status,
					TestName:     fmt.Sprintf("TargetStress GEMM %s", r.GemmName),
					TestOutput:   fmt.Sprintf("GEMM=%s, Mean=%.3f", r.GemmName, r.Mean),
					ErrorCode:    0,
					ErrorMessage: "",
				}
				perMap[r.HCUId] = append(perMap[r.HCUId], dr)
			}

			// 把 map 转为 diag.PerHCU（按 HCU 索引升序）
			keys := make([]int, 0, len(perMap))
			for k := range perMap {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			for _, d := range keys {
				diag.PerHCU = append(diag.PerHCU, dcgm.HCUResult{
					HCU:         d,
					RC:          0,
					DiagResults: perMap[d],
				})
			}

			// 写入一个 Software 层面的 summary（便于快速查看）
			diag.Software = append(diag.Software, dcgm.DiagResult{
				Status:     DiagResultPass,
				TestName:   "TargetStress Summary",
				TestOutput: fmt.Sprintf("Parsed %d gemm entries, logdir=%s", len(tsRes.Results), tsRes.LogDir),
				ErrorCode:  0,
			})
			// DeviceNumber 以参与 HCU 数为准
			diag.DeviceNumber = len(keys)
		} else {
			// 失败或无结果：把错误写入 Software，用于诊断
			diag.Software = append(diag.Software, dcgm.DiagResult{
				Status:       DiagResultFail,
				TestName:     "TargetStress",
				TestOutput:   "",
				ErrorCode:    -1,
				ErrorMessage: fmt.Sprintf("TargetStressTestResult error: %v", err),
			})
		}

		persist.Result = convertDiagResults(&diag)
		jobStoreMu.Unlock()

		// 持久化到文件 logs/<job-id>.json
		if err := saveJobToFile(&persist); err != nil {
			fmt.Printf("saveJobToFile failed for targetstress job %s: %v\n", j.ID, err)
		}

		// 从内存中删除该 job（持久化成功或失败后都删除内存条目）
		jobStoreMu.Lock()
		delete(jobStore, j.ID)
		jobStoreMu.Unlock()

		// 清理 currentJobID（仅当仍指向本 job 时）
		if v := currentJobID.Load(); v != nil {
			if s, ok := v.(string); ok && s == j.ID {
				currentJobID.Store("")
			}
		}
	}(job)

	// 立即返回 job id 与状态查询 URL
	c.JSON(http.StatusAccepted, SuccessResponse(DiagResp{
		JobID:     id,
		Status:    job.Status,
		StatusURL: fmt.Sprintf("/diag/status/%s", id),
	}))
}

// @Tags MemtestCL
// @Summary 异步启动 MemtestCL
// @Description 异步启动 MemtestCL 对所有 HCU 进行测试，测试结果会落盘，可通过 /diag/status/<job-id> 查询作业状态和结果
// @Accept json
// @Produce json
// @Success 200 {object} DiagResp "返回作业 ID 和状态查询 URL"
// @Failure 409 {object} Response "已有正在运行的作业"
// @Failure 500 {object} Response "启动作业失败"
// @Router /memtestcl/start [post]
func startMemtestCL(c *gin.Context) {
	// 若已有运行 job，则拒绝（避免并发压测）
	if v := currentJobID.Load(); v != nil {
		if s, ok := v.(string); ok && s != "" {
			c.JSON(http.StatusConflict, gin.H{
				"error":      "another job is running",
				"running_id": s,
			})
			return
		}
	}

	// 生成 job id（前缀 memtest-）
	id := fmt.Sprintf("memtest-%d-%d", time.Now().UnixNano(), atomic.AddUint64(&jobIDCounter, 1))
	job := &Job{
		ID:     id,
		Level:  0,
		Status: JobStatusPending,
	}

	// 保存到内存 store 并标记为当前运行 job
	jobStoreMu.Lock()
	jobStore[id] = job
	jobStoreMu.Unlock()
	currentJobID.Store(id)

	// 后台执行实际测试
	go func(j *Job) {
		// 标记为 running
		jobStoreMu.Lock()
		j.Status = JobStatusRunning
		j.StartedAt = time.Now()
		jobStoreMu.Unlock()

		// 构造 dvIdList（全部设备）
		numDevices, _ := dcgm.NumMonitorDevices()
		dvIdList := make([]int, 0, numDevices)
		for i := 0; i < numDevices; i++ {
			dvIdList = append(dvIdList, i)
		}

		// 调用 dcgm API（应为不会打印的版本）
		memAllRes, err := dcgm.MemtestCLTestResult(dvIdList)

		// 更新内存 job 状态与时间
		jobStoreMu.Lock()
		if err != nil {
			j.Status = JobStatusFailed
			j.ErrorMessage = err.Error()
		} else {
			if dcgm.IsDiagStopped() {
				j.Status = JobStatusCanceled
			} else {
				j.Status = JobStatusDone
			}
		}
		j.EndedAt = time.Now()

		// 构造持久化副本
		var persist Job
		persist.ID = j.ID
		persist.Level = j.Level
		persist.Status = j.Status
		persist.ErrorMessage = j.ErrorMessage
		persist.StartedAt = j.StartedAt
		persist.EndedAt = j.EndedAt

		// 把 memAllRes 转换为 dcgm.DiagResults
		var diag dcgm.DiagResults
		if err == nil && memAllRes.Results != nil {
			// 保证按 HCU 索引顺序写入
			// 使用 map 临时聚合，然后按键排序写入 PerHCU
			perMap := make(map[int][]dcgm.DiagResult)
			for _, r := range memAllRes.Results {
				status := DiagResultPass
				if !r.Passed {
					status = DiagResultWarn
				}
				// 将 summary map 转为 "k: v" 拼接字符串
				parts := make([]string, 0, len(r.Summary))
				for kk, vv := range r.Summary {
					parts = append(parts, fmt.Sprintf("%s: %s", kk, vv))
				}
				sort.Strings(parts)
				out := strings.Join(parts, " ; ")
				dr := dcgm.DiagResult{
					Status:       status,
					TestName:     "MemtestCL",
					TestOutput:   out,
					ErrorCode:    0,
					ErrorMessage: "",
				}
				perMap[r.HCUId] = append(perMap[r.HCUId], dr)
			}
			// 把 map 转成 PerHCU 切片（按升序）
			keys := make([]int, 0, len(perMap))
			for k := range perMap {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			for _, d := range keys {
				diag.PerHCU = append(diag.PerHCU, dcgm.HCUResult{
					HCU:         d,
					RC:          0,
					DiagResults: perMap[d],
				})
			}
			// Software summary
			diag.Software = append(diag.Software, dcgm.DiagResult{
				Status:     DiagResultPass,
				TestName:   "MemtestCL Summary",
				TestOutput: fmt.Sprintf("Parsed %d HCU results, logdir=%s", len(memAllRes.Results), memAllRes.LogDir),
				ErrorCode:  0,
			})
			diag.DeviceNumber = len(keys)
		} else {
			diag.Software = append(diag.Software, dcgm.DiagResult{
				Status:       DiagResultFail,
				TestName:     "MemtestCL",
				TestOutput:   "",
				ErrorCode:    -1,
				ErrorMessage: fmt.Sprintf("MemtestCLTestResult error: %v", err),
			})
		}

		persist.Result = convertDiagResults(&diag)
		jobStoreMu.Unlock()

		// 持久化到文件 logs/<job-id>.json
		if err := saveJobToFile(&persist); err != nil {
			fmt.Printf("saveJobToFile failed for memtest job %s: %v\n", j.ID, err)
		}

		// 从内存中删除该 job（持久化后）
		jobStoreMu.Lock()
		delete(jobStore, j.ID)
		jobStoreMu.Unlock()

		// 清理 currentJobID（仅当仍指向本 job）
		if v := currentJobID.Load(); v != nil {
			if s, ok := v.(string); ok && s == j.ID {
				currentJobID.Store("")
			}
		}
	}(job)

	// 立即返回 job id 与状态查询 url
	c.JSON(http.StatusAccepted, SuccessResponse(DiagResp{
		JobID:     id,
		Status:    job.Status,
		StatusURL: fmt.Sprintf("/diag/status/%s", id),
	}))
}

func summarizeEdppOutput(output string) string {
	const maxLength = 512

	summary := strings.Join(strings.Fields(output), " ")
	if len(summary) <= maxLength {
		return summary
	}
	return summary[:maxLength] + "..."
}

func convertEdppResult(edppResult dcgm.EDPPResult) dcgm.DiagResults {
	var result dcgm.DiagResults
	result.DeviceNumber = len(edppResult.HCUEdppResults)
	summaryStatus := DiagResultPass

	for _, hcuResult := range edppResult.HCUEdppResults {
		perHCU := dcgm.HCUResult{HCU: hcuResult.HCUId}
		for _, pattern := range hcuResult.PatternResults {
			status := DiagResultPass
			if pattern.ECCCount > 0 || pattern.MemoryErrorCount > 0 || pattern.ComputeErrorCount > 0 {
				status = DiagResultWarn
				if summaryStatus == DiagResultPass {
					summaryStatus = DiagResultWarn
				}
			}
			perHCU.DiagResults = append(perHCU.DiagResults, dcgm.DiagResult{
				Status:     status,
				TestName:   fmt.Sprintf("EDPp Pattern %s", pattern.PatternName),
				TestOutput: fmt.Sprintf("Pattern=%s, ECC=%d, Mem=%d, Compute=%d", pattern.PatternName, pattern.ECCCount, pattern.MemoryErrorCount, pattern.ComputeErrorCount),
			})
		}

		for _, gemmPower := range hcuResult.GemmPowerResults {
			status := DiagResultPass
			errorCode := 0
			if !gemmPower.Passed {
				status = DiagResultFail
				errorCode = -1
				summaryStatus = DiagResultFail
			}
			testOutput := fmt.Sprintf("Backend=%s, %s", gemmPower.Backend, gemmPower.Summary())
			if outputSummary := summarizeEdppOutput(gemmPower.Output); outputSummary != "" {
				testOutput += ", Output=" + outputSummary
			}
			perHCU.DiagResults = append(perHCU.DiagResults, dcgm.DiagResult{
				Status:       status,
				TestName:     fmt.Sprintf("EDPp %s", gemmPower.Backend),
				TestOutput:   testOutput,
				ErrorCode:    errorCode,
				ErrorMessage: gemmPower.Error,
			})
		}
		if len(perHCU.DiagResults) == 0 {
			perHCU.RC = -1
			perHCU.DiagResults = append(perHCU.DiagResults, dcgm.DiagResult{
				Status:       DiagResultFail,
				TestName:     fmt.Sprintf("EDPp %s", hcuResult.Backend),
				ErrorCode:    -1,
				ErrorMessage: "EDPp backend returned no result",
			})
			summaryStatus = DiagResultFail
		}
		result.PerHCU = append(result.PerHCU, perHCU)
	}

	result.Software = append(result.Software, dcgm.DiagResult{
		Status:     summaryStatus,
		TestName:   "EDPp Summary",
		TestOutput: fmt.Sprintf("Parsed %d HCU EDPp results, logdir=%s", len(edppResult.HCUEdppResults), edppResult.LogDir),
	})
	return result
}

// startEdpp godoc
// @Summary      启动 EDPp 测试
// @Description  异步启动 EDPp 稳定性与错误注入测试。
// @Description  按 HCU 和 Pattern 聚合结果并持久化。
// @Tags         Diag
// @Accept       json
// @Produce      json
// @Success      200 {object} DiagResp
// @Failure      409 {object} Response "已有作业正在运行"
// @Router       /edpp/start [post]
func startEdpp(c *gin.Context) {
	// 如果已有正在运行的 job，则拒绝（避免并发压测）
	if v := currentJobID.Load(); v != nil {
		if s, ok := v.(string); ok && s != "" {
			c.JSON(http.StatusConflict, gin.H{
				"error":      "another job is running",
				"running_id": s,
			})
			return
		}
	}

	// 生成 job id（前缀 edpp-）
	id := fmt.Sprintf("edpp-%d-%d", time.Now().UnixNano(), atomic.AddUint64(&jobIDCounter, 1))
	job := &Job{
		ID:     id,
		Level:  0,
		Status: JobStatusPending,
	}

	// 存入内存 store 并标记为当前运行 job
	jobStoreMu.Lock()
	jobStore[id] = job
	jobStoreMu.Unlock()
	currentJobID.Store(id)

	// 在后台执行实际测试
	go func(j *Job) {
		// 标记为 running
		jobStoreMu.Lock()
		j.Status = JobStatusRunning
		j.StartedAt = time.Now()
		jobStoreMu.Unlock()

		// 调用 dcgm API（应为不会向 stdout 打印的版本）
		edppRes, err := dcgm.EDPpTestResult()

		// 更新内存 job 状态与结束时间
		jobStoreMu.Lock()
		if err != nil {
			j.Status = JobStatusFailed
			j.ErrorMessage = err.Error()
		} else {
			if dcgm.IsDiagStopped() {
				j.Status = JobStatusCanceled
			} else {
				j.Status = JobStatusDone
			}
		}
		j.EndedAt = time.Now()

		// 构造持久化副本
		var persist Job
		persist.ID = j.ID
		persist.Level = j.Level
		persist.Status = j.Status
		persist.ErrorMessage = j.ErrorMessage
		persist.StartedAt = j.StartedAt
		persist.EndedAt = j.EndedAt

		var diag dcgm.DiagResults
		if len(edppRes.HCUEdppResults) > 0 {
			diag = convertEdppResult(edppRes)
		}
		if err != nil {
			diag.Software = append(diag.Software, dcgm.DiagResult{
				Status:       DiagResultFail,
				TestName:     "EDPp Test",
				ErrorCode:    -1,
				ErrorMessage: fmt.Sprintf("EDPpTestResult error: %v", err),
			})
		}
		if len(diag.Software) == 0 {
			diag.Software = append(diag.Software, dcgm.DiagResult{
				Status:       DiagResultFail,
				TestName:     "EDPp Test",
				ErrorCode:    -1,
				ErrorMessage: "EDPpTestResult returned no result",
			})
		}

		persist.Result = convertDiagResults(&diag)
		jobStoreMu.Unlock()

		// 写盘 logs/<job-id>.json
		if err := saveJobToFile(&persist); err != nil {
			fmt.Printf("saveJobToFile failed for edpp job %s: %v\n", j.ID, err)
		}

		// 从内存移除该 job（持久化后）
		jobStoreMu.Lock()
		delete(jobStore, j.ID)
		jobStoreMu.Unlock()

		// 清理 currentJobID（仅当仍指向本 job）
		if v := currentJobID.Load(); v != nil {
			if s, ok := v.(string); ok && s == j.ID {
				currentJobID.Store("")
			}
		}
	}(job)

	// 立即返回 job id 与状态查询 url
	c.JSON(http.StatusAccepted, SuccessResponse(DiagResp{
		JobID:     id,
		Status:    job.Status,
		StatusURL: fmt.Sprintf("/diag/status/%s", id),
	}))
}

/************** health ****************/

// @Tags Device
// SetHealthCheckConfig 设置是否开启健康检查
// @Summary 设置是否开启健康检查
// @Description 配置是否启用健康检查，并选择检查项
// @Accept json
// @Produce json
// @Param enabled body bool true "是否开启"
// @Param options body []int true "检查项"
// @Success 200 {object} string "设置成功"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /SetHealthCheckConfig [post]
func SetHealthCheckConfig(c *gin.Context) {
	var params struct {
		Enabled bool  `json:"enabled"`
		Options []int `json:"options"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, "参数无效")
		return
	}

	err := dcgm.SetHealthCheckConfig(params.Enabled, params.Options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// @Tags Device
// GetHealthCheckConfig 获取健康检查项信息
// @Summary 获取健康检查项的配置
// @Description 获取当前健康检查项的配置信息
// @Produce json
// @Success 200 {object} HealthCheckConfig "返回健康检查项信息"
// @Failure 500 {object} error "服务器内部错误"
// @Router /GetHealthCheckConfig [get]
func GetHealthCheckConfig(c *gin.Context) {
	// dcgm 层返回
	dcgmConfig, err := dcgm.GetHealthCheckConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// ===== 显式转换 dcgm -> router =====
	resp := HealthCheckConfig{
		Enabled: dcgmConfig.Enabled,
		Options: dcgmConfig.Options,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Device
// DeleteHealthCheckConfig 删除健康检查项信息
// @Summary 删除健康检查项的配置
// @Description 删除当前的健康检查项配置
// @Produce json
// @Success 200 {string} string "删除成功"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DeleteHealthCheckConfig [delete]
func DeleteHealthCheckConfig(c *gin.Context) {
	err := dcgm.DeleteHealthCheckConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, "删除成功")
}

// @Tags Device
// HealthCheckById 设备健康检查
// @Summary 获取指定设备的健康检查结果
// @Description 根据设备 ID 列表返回对应设备的健康检查状态信息
// @Produce json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {object} HealthCheckByGroupResp "返回指定设备的健康检查结果"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /HealthCheckById [get]
func HealthCheckById(c *gin.Context) {
	// 解析请求体
	var dvIdList []int
	if err := c.BindJSON(&dvIdList); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	// dcgm 层返回结果
	dcgmDeviceHealths, err := dcgm.HealthCheckById(dvIdList, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// ===== dcgm -> router 结构体转换 =====
	deviceHealths := make([]DeviceHealth, 0, len(dcgmDeviceHealths))
	for _, dh := range dcgmDeviceHealths {
		watches := make([]SystemWatch, 0, len(dh.Watches))
		for _, w := range dh.Watches {
			watches = append(watches, SystemWatch{
				Type:   w.Type,
				Status: w.Status,
				Error:  w.Error,
				Result: w.Result,
			})
		}

		deviceHealths = append(deviceHealths, DeviceHealth{
			HCU:     dh.HCU,
			Status:  dh.Status,
			Watches: watches,
		})
	}

	resp := HealthCheckByGroupResp{
		DeviceHealths: deviceHealths,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// HealthCheckByGroup 通过组ID获取设备健康状态
// @Summary 获取设备健康状态
// @Description 根据 groupId 获取设备健康状态，可选择是否检查 HealthConfig
// @Tags Health
// @Accept json
// @Produce json
// @Produce json
// @Param request body HealthCheckByGroupReq true "请求参数"
// @Success 200 {object} HealthCheckByGroupResp
// @Failure 400 {object} Response "请求参数错误"
// @Failure 500 {object} Response "内部错误"
// @Router /health/group [post]
func HealthCheckByGroup(c *gin.Context) {
	var req HealthCheckByGroupReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid JSON body"))
		return
	}

	dcgmDeviceHealths, err := dcgm.HealthCheckByGroupId(req.GroupId, req.CheckHealthConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// 转换类型
	var routerDeviceHealths []DeviceHealth
	for _, d := range dcgmDeviceHealths {
		var watches []SystemWatch
		for _, w := range d.Watches {
			watches = append(watches, SystemWatch{
				Type:   w.Type,
				Status: w.Status,
				Error:  w.Error,
				Result: w.Result,
			})
		}
		routerDeviceHealths = append(routerDeviceHealths, DeviceHealth{
			HCU:     d.HCU,
			Status:  d.Status,
			Watches: watches,
		})
	}

	resp := HealthCheckByGroupResp{
		DeviceHealths: routerDeviceHealths,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags Device
// GetDeviceModelInfos 获取设备型号信息
// @Summary 获取所有设备的型号信息
// @Description 返回系统中所有设备的型号信息列表
// @Produce json
// @Success 200 {object} GetDeviceModelInfosResp "设备型号信息列表"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /GetDeviceModelInfos [get]
func GetDeviceModelInfos(c *gin.Context) {
	dcgmDevices := dcgm.GetDeviceModelInfos()

	routerDevices := make([]DeviceModelInfo, 0, len(dcgmDevices))
	for _, d := range dcgmDevices {
		routerDevices = append(routerDevices, DeviceModelInfo{
			Model:      d.Model,
			CUCount:    d.CUCount,
			MemorySize: d.MemorySize,
		})
	}

	resp := GetDeviceModelInfosResp{
		Devices: routerDevices,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

// @Tags System
// @Summary 获取进程列表信息
// @Description 获取并返回进程信息和使用的 HCU 设备信息
// @Produce json
// @Success 200 {object} GetProcessInfoResp "进程信息列表"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /process-info [get]
func GetProcessInfo(c *gin.Context) {
	dcgmProcesses, err := dcgm.ProcessHCUInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// 没有进程时，返回空数组（比 nil 更友好）
	if len(dcgmProcesses) == 0 {
		c.JSON(http.StatusOK, SuccessResponse(GetProcessInfoResp{
			Processes: []Process{},
		}))
		return
	}

	routerProcesses := make([]Process, 0, len(dcgmProcesses))
	for _, p := range dcgmProcesses {
		routerProcesses = append(routerProcesses, Process{
			ProcessID:    p.ProcessID,
			ProcessName:  p.ProcessName,
			Pasid:        p.Pasid,
			VramUsage:    p.VramUsage,
			SdmaUsage:    p.SdmaUsage,
			CuOccupancy:  p.CuOccupancy,
			MinorNumbers: p.MinorNumbers,
		})
	}

	resp := GetProcessInfoResp{
		Processes: routerProcesses,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

func toCompatibilityResponse(result dcgm.CompatibilityResult) CompatibilityResult {
	return CompatibilityResult{
		Status:            string(result.Status),
		CardModel:         result.CardModel,
		CurrentDriver:     result.CurrentDriver,
		CurrentDTK:        result.CurrentDTK,
		RecommendedDriver: result.RecommendedDriver,
		RecommendedDTK:    result.RecommendedDTK,
	}
}

// Compatible 校验卡型号、驱动版本和 DTK 版本的兼容性
// @Summary 校验兼容性
// @Description 根据卡型号、驱动版本和 DTK 版本，检查是否符合兼容性要求
// @Produce json
// @Param cardModel query string true "卡型号，例如 K100, K100_AI"
// @Param driverVersion query string true "驱动版本，例如 5.16.18, 6.2.26"
// @Param dtkVersion query string true "DTK 版本，例如 23.10, 24.04"
// @Success 200 {object} CompatibilityResult "兼容性检查结果，版本不匹配时 Status 为 warn"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /compatible [get]
func Compatible(c *gin.Context) {
	// 从查询参数中获取输入
	cardModel := c.Query("cardModel")
	driverVersion := c.Query("driverVersion")
	dtkVersion := c.Query("dtkVersion")
	// 校验必填参数
	if cardModel == "" || driverVersion == "" || dtkVersion == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求参数错误: 缺少必填参数"))
		return
	}
	// 调用业务逻辑函数
	result, err := dcgm.Compatible(cardModel, driverVersion, dtkVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(toCompatibilityResponse(result)))
}

// @Tags Device
// @Summary 获取 MIG 分区信息
// @Description 获取系统中所有 MIG 分区的详细信息
// @Produce json
// @Success 200 {object} MigInfosResp "MIG 分区信息列表"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /mig/infos [get]
func MigInfos(c *gin.Context) {
	dcgmMigInfos, err := dcgm.MigInfos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	migInfos := make([]MigInfo, 0, len(dcgmMigInfos))
	for _, m := range dcgmMigInfos {
		migInfos = append(migInfos, MigInfo{
			DvInd:             m.DvInd,
			MigId:             m.MigId,
			Name:              m.Name,
			UUID:              m.UUID,
			ComputeUnit:       m.ComputeUnit,
			MemoryTotal:       m.MemoryTotal,
			GpuInstanceId:     m.GpuInstanceId,
			ComputeInstanceId: m.ComputeInstanceId,
			PciBusNumber:      m.PciBusNumber,
			GiProfileId:       m.GiProfileId,
			CiProfileId:       m.CiProfileId,
		})
	}

	resp := MigInfosResp{
		MigInfos: migInfos,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}

func ListAllGroups(c *gin.Context) {
	groups, err := dcgm.ListAllGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	response := map[string]interface{}{
		"groups": groups,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func GetGroupInfo(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid group ID"))
		return
	}

	groupInfo, err := dcgm.GetGroupInfo(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	response := map[string]interface{}{
		"groupInfo": groupInfo,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func CreateGroup(c *gin.Context) {
	var req CreateGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	groupId, err := dcgm.CreateGroup(req.GroupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	response := map[string]interface{}{
		"groupId": groupId,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func DestroyGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid group ID"))
		return
	}
	err = dcgm.DestroyGroup(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func AddToGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid group ID"))
		return
	}

	var req AddHcuToGroupRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	err = dcgm.AddToGroup(groupId, req.DvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func AddEntityToGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid group ID"))
		return
	}

	var req EntityListRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	err = dcgm.AddEntityToGroup(groupId, req.EntityList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func RemoveEntityFromGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid group ID"))
		return
	}

	var req EntityListRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	err = dcgm.RemoveEntityFromGroup(groupId, req.EntityList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func CreateFieldGroup(c *gin.Context) {
	var req CreateFieldGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	fieldGroupId, err := dcgm.CreateFieldGroup(req.FieldGroupName, req.FieldIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	response := map[string]interface{}{
		"fieldGroupId": fieldGroupId,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func DestroyFieldGroup(c *gin.Context) {
	fieldGroupId, err := strconv.Atoi(c.Param("fieldGroupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid field group ID"))
		return
	}
	err = dcgm.DestroyFieldGroup(fieldGroupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func ListAllFieldGroups(c *gin.Context) {
	fieldGroups, err := dcgm.ListAllFieldGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	response := map[string]interface{}{
		"fieldGroups": fieldGroups,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func GetFieldGroupInfo(c *gin.Context) {
	fieldGroupId, err := strconv.Atoi(c.Param("fieldGroupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid field group ID"))
		return
	}

	fieldGroupInfo, err := dcgm.GetFieldGroupInfo(fieldGroupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	response := map[string]interface{}{
		"fieldGroupInfo": fieldGroupInfo,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func SetPolicy(c *gin.Context) {
	hcuIndex, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid HCU ID"))
		return
	}

	var req dcgm.Policy
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	err = dcgm.SetPolicy(req, hcuIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func ClearPolicy(c *gin.Context) {
	hcuIndex, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid HCU ID"))
		return
	}
	err = dcgm.ClearPolicy([]int{hcuIndex})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func ListPolicies(c *gin.Context) {
	numHcus, err := dcgm.NumMonitorDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	hcuList := make([]int, numHcus)
	for i := 0; i < numHcus; i++ {
		hcuList[i] = i
	}
	policyList, err := dcgm.GetPolicy(hcuList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"policyList": policyList,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func GetPolicy(c *gin.Context) {
	hcuIndex, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid HCU ID"))
		return
	}
	policyList, err := dcgm.GetPolicy([]int{hcuIndex})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"policyList": policyList,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func JudgePolicyConditions(c *gin.Context) {
	var req DevIdListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	hcuIndex, err := dcgm.JudgePolicyConditions(req.DevIdList)
	response := map[string]interface{}{
		"errorDvInd": hcuIndex,
		"error":      err,
	}

	if err == nil {
		response["errorDvInd"] = nil
	}

	c.JSON(http.StatusOK, SuccessResponse(response))
}

func TakePolicyAction(c *gin.Context) {
	hcuIndex, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid HCU ID"))
		return
	}
	err = dcgm.TakePolicyAction(hcuIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func GetFieldMetaById(c *gin.Context) {
	fieldId, err := strconv.Atoi(c.Param("fieldId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid field ID"))
		return
	}
	fieldMeta := dcgm.GetFieldMetaById(fieldId)
	response := map[string]interface{}{
		"fieldMeta": fieldMeta,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func WatchFields(c *gin.Context) {
	hcuIndex, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid HCU ID"))
		return
	}
	var req FieldIdListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	err = dcgm.WatchFields(hcuIndex, req.FieldIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func WatchFieldsWithEntity(c *gin.Context) {
	entityGroupId, err := strconv.Atoi(c.Param("entityGroupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid Entity Group ID"))
		return
	}
	entityId, err := strconv.Atoi(c.Param("entityId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid Entity ID"))
		return
	}
	var req FieldIdListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	err = dcgm.WatchFieldsWithEntity(dcgm.Field_Entity_Group(entityGroupId), entityId, req.FieldIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func WatchFieldsWithGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid group ID"))
		return
	}

	var req WatchFieldsWithGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	updateFreq := dcgm.DefaultUpdateFreq
	if req.UpdateFreq > 0 {
		updateFreq = time.Duration(req.UpdateFreq * float64(time.Second))
	}

	maxKeepAge := time.Duration(dcgm.DefaultMaxKeepAge)
	if req.MaxKeepAge > 0 {
		maxKeepAge = time.Duration(req.MaxKeepAge * float64(time.Second))
	}

	maxKeepSamples := int32(dcgm.DefaultMaxKeepSamples)
	if req.MaxKeepSamples > 0 {
		maxKeepSamples = int32(req.MaxKeepSamples)
	}

	err = dcgm.WatchFieldsWithEntityGroup(req.FieldIdList, groupId, updateFreq, maxKeepAge, maxKeepSamples)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func ListFieldMeta(c *gin.Context) {
	fieldMetaList := dcgm.ListFieldMeta()
	response := map[string]interface{}{
		"fieldMetaList": fieldMetaList,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func UnWatchFields(c *gin.Context) {
	hcuIndex, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid HCU ID"))
		return
	}
	dcgm.UnWatchFields(hcuIndex)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func UnWatchFieldsWithEntity(c *gin.Context) {
	entityGroupId, err := strconv.Atoi(c.Param("entityGroupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid Entity Group ID"))
		return
	}
	entityId, err := strconv.Atoi(c.Param("entityId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid Entity ID"))
		return
	}
	dcgm.UnWatchFieldsWithEntity(dcgm.Field_Entity_Group(entityGroupId), entityId)
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func UnWatchFieldsWithGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid group ID"))
		return
	}

	err = dcgm.UnWatchFieldsWithGroup(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func GetLatestValuesForFields(c *gin.Context) {
	hcuIndex, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid HCU ID"))
		return
	}
	var req FieldIdListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	fieldValueList, err := dcgm.GetLatestValuesForFields(hcuIndex, req.FieldIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"fieldValueList": fieldValueList,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func EntityGetLatestValues(c *gin.Context) {
	entityGroupId, err := strconv.Atoi(c.Param("entityGroupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid Entity Group ID"))
		return
	}
	entityId, err := strconv.Atoi(c.Param("entityId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid Entity ID"))
		return
	}
	var req FieldIdListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	fieldValueList, err := dcgm.EntityGetLatestValues(dcgm.Field_Entity_Group(entityGroupId), entityId, req.FieldIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"fieldValueList": fieldValueList,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func GetLatestValuesForFieldsWithGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid group ID"))
		return
	}

	var req FieldIdListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	groupInfo, err := dcgm.GetGroupInfo(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	entities := groupInfo.EntityList
	fieldValueList, err := dcgm.EntitiesGetLatestValues(entities, req.FieldIdList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	entityList := make([]EntityListWithFieldValuesResp, 0, len(entities))

	type key struct {
		g dcgm.Field_Entity_Group
		i int
	}

	byEntity := make(map[key][]dcgm.FieldValue_v1, len(entities))
	for _, fv := range fieldValueList {
		k := key{g: fv.EntityGroupId, i: fv.EntityId}
		byEntity[k] = append(byEntity[k], dcgm.FieldValue_v1{
			FieldId:   fv.FieldId,
			Timestamp: fv.Timestamp,
			Value:     fv.Value,
			Err:       fv.Err,
			Tag:       fv.Tag,
		})
	}

	for _, e := range entities {
		k := key{g: e.EntityGroupId, i: e.EntityId}
		fvs := byEntity[k]
		entityList = append(entityList, EntityListWithFieldValuesResp{
			EntityGroupId:  k.g,
			EntityId:       k.i,
			FieldValueList: fvs,
		})
	}

	response := map[string]interface{}{
		"entityList": entityList,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

func GetSupportedMetricGroups(c *gin.Context) {
	hcuIndex, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}
	metricGroupList, err := dcgm.GetSupportedMetricGroups(hcuIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	response := map[string]interface{}{
		"metricGroupList": metricGroupList,
	}
	c.JSON(http.StatusOK, SuccessResponse(response))
}

// HSLLinkItem 单条 HSL/XHCL 链路的通断状态
type HSLLinkItem struct {
	LinkID int    `json:"linkId"`
	State  uint32 `json:"state"` // 驱动原始状态值
	Up     bool   `json:"up"`    // true = 链路正常
}

// HSLLinkStateResp /hsl/linkstate/:dvInd 响应体
type HSLLinkStateResp struct {
	DvInd int           `json:"dvInd"`
	Links []HSLLinkItem `json:"links"`
}

// HSLLinkState godoc
// @Summary      查询指定设备的 HSL/XHCL 链路通断状态
// @Description  返回该设备所有 HSL（XHCL/HSL）链路的当前通断状态（Up/Down）。
// @Description  对应 DCGM_FI_DEV_NVLINK_LINK_DOWN_COUNT_TOTAL 的 HCU 近似语义。
// @Description  注意：HCU 驱动不提供历史断链计数，仅反映当前时刻通断状态。
// @Tags         Topology
// @Produce      json
// @Param        dvInd  path  int  true  "设备索引"
// @Success      200 {object} HSLLinkStateResp
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /hsl/linkstate/{dvInd} [get]
func HSLLinkState(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	states, err := dcgm.XhclLinkStates(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	links := make([]HSLLinkItem, len(states))
	for i, s := range states {
		links[i] = HSLLinkItem{
			LinkID: s.LinkID,
			State:  s.State,
			Up:     s.Up,
		}
	}

	c.JSON(http.StatusOK, SuccessResponse(HSLLinkStateResp{
		DvInd: dvInd,
		Links: links,
	}))
}

func toThrottleFlags(flags dcgm.ThrottleFlags) ThrottleFlags {
	return ThrottleFlags{
		Thermal:    flags.Thermal,
		Power:      flags.Power,
		Slowdown:   flags.Slowdown,
		BoardLimit: flags.BoardLimit,
	}
}

// GetThrottleFlags godoc
// @Summary      查询设备降频状态
// @Description  返回设备当前各类降频保护的触发状态（true = 正在降频）。
// @Description  对应 NVIDIA DCGM 的 THERMAL/POWER/SLOWDOWN/BOARD_LIMIT_THROTTLE 四项指标。
// @Description  不支持项：LOW_UTIL_VIOLATION / RELIABILITY_VIOLATION / SYNC_BOOST_VIOLATION（NV 专有机制）。
// @Tags         Monitor
// @Produce      json
// @Param        dvInd  path  int  true  "设备索引"
// @Success      200 {object} ThrottleFlags
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /throttle/flags/{dvInd} [get]
func GetThrottleFlags(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}

	flags, err := dcgm.DevThrottleFlags(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(toThrottleFlags(flags)))
}

// GetMemoryTemperature godoc
// @Summary      查询显存温度
// @Description  返回显存颗粒当前温度（摄氏度）。对应 DCGM_FI_DEV_MEMORY_TEMP。
// @Tags         Monitor
// @Produce      json
// @Param        dvInd  path  int  true  "设备索引"
// @Success      200 {object} float64
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /memory/temperature/{dvInd} [get]
func GetMemoryTemperature(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	temp, err := dcgm.MemoryTemperature(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"temperature": temp}))
}

// GetBoardTemperature godoc
// @Summary      查询板卡边缘温度
// @Description  返回板卡边缘传感器当前温度（摄氏度）。对应 DCGM_FI_DEV_BOARD_TEMP。
// @Tags         Monitor
// @Produce      json
// @Param        dvInd  path  int  true  "设备索引"
// @Success      200 {object} float64
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /board/temperature/{dvInd} [get]
func GetBoardTemperature(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	temp, err := dcgm.BoardTemperature(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"temperature": temp}))
}

// GetMemoryAvailable godoc
// @Summary      查询可用显存量
// @Description  返回驱动当前可向应用分配的近似剩余显存量（字节）。
// @Description  以 total_vram - used_vram 近似，对应 DCGM_FI_DEV_FB_RESERVED 语义。
// @Tags         Monitor
// @Produce      json
// @Param        dvInd  path  int  true  "设备索引"
// @Success      200 {object} int64
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /memory/available/{dvInd} [get]
func GetMemoryAvailable(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	available, err := dcgm.MemoryAvailable(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"available": available}))
}

// GetMemClock godoc
// @Summary      查询显存当前频率
// @Description  返回显存当前运行频率（MHz）。对应 DCGM_FI_DEV_MEM_CLOCK。
// @Tags         Monitor
// @Produce      json
// @Param        dvInd  path  int  true  "设备索引"
// @Success      200 {object} float64
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /clock/mem/{dvInd} [get]
func GetMemClock(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	clk, err := dcgm.HCUMclk(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"clockMhz": clk}))
}

// GetGfxClockMax godoc
// @Summary      查询 GFX 最大额定频率
// @Description  返回 GFX（shader）时钟的最大额定频率（MHz）。对应 DCGM_FI_DEV_SM_CLOCK_MAX。
// @Tags         Monitor
// @Produce      json
// @Param        dvInd  path  int  true  "设备索引"
// @Success      200 {object} uint64
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /clock/gfx/max/{dvInd} [get]
func GetGfxClockMax(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	maxClk, err := dcgm.DevGfxClockMax(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"clockMhz": maxClk}))
}

// GetMemClockMax godoc
// @Summary      查询显存最大额定频率
// @Description  返回显存时钟的最大额定频率（MHz）。对应 DCGM_FI_DEV_MEM_CLOCK_MAX。
// @Tags         Monitor
// @Produce      json
// @Param        dvInd  path  int  true  "设备索引"
// @Success      200 {object} uint64
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /clock/mem/max/{dvInd} [get]
func GetMemClockMax(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	maxClk, err := dcgm.DevMemClockMax(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"clockMhz": maxClk}))
}

// GetDevCuNum godoc
// @Summary      查询 CU 数量
// @Description  返回设备的 Compute Unit（CU）数量，属于静态设备规格信息。
// @Tags         DevInfo
// @Produce      json
// @Param        dvInd  path  int  true  "设备索引"
// @Success      200 {object} int
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /device/cunum/{dvInd} [get]
func GetDevCuNum(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	num, err := dcgm.DevCuNum(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"cuNum": num}))
}

// GetDevNodeId godoc
// @Summary      查询 KFD node ID
// @Description  返回设备的 KFD（Kernel Fusion Driver）node ID，用于拓扑关联。
// @Tags         DevInfo
// @Produce      json
// @Param        dvInd  path  int  true  "设备索引"
// @Success      200 {object} uint32
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /device/nodeid/{dvInd} [get]
func GetDevNodeId(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	nodeId, err := dcgm.DevNodeId(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"nodeId": nodeId}))
}

// GetDevMemOverdrive godoc
// @Summary      查询显存 overdrive 百分比
// @Description  返回显存当前 overdrive 百分比（0-20）。对应 hy-smi --showmemoverdrive。
// @Tags         Monitor
// @Produce      json
// @Param        dvInd  path  int  true  "设备索引"
// @Success      200 {object} int
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /device/memoverdrive/{dvInd} [get]
func GetDevMemOverdrive(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid device ID"))
		return
	}
	od, err := dcgm.DevMemOverdriveLevelGet(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"memOverdrive": od}))
}

// GetP2PAccessible godoc
// @Summary      查询两设备间 P2P 可达性
// @Description  返回 srcDvInd → dstDvInd 方向的 P2P 直接访问可达性。对应 hy-smi --showtopoaccess。
// @Tags         Topology
// @Produce      json
// @Param        srcDvInd  path  int  true  "源设备索引"
// @Param        dstDvInd  path  int  true  "目标设备索引"
// @Success      200 {object} bool
// @Failure      400 {object} Response "参数错误"
// @Failure      500 {object} Response "查询失败"
// @Router       /topology/p2p/{srcDvInd}/{dstDvInd} [get]
func GetP2PAccessible(c *gin.Context) {
	srcDvInd, err := strconv.Atoi(c.Param("srcDvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid srcDvInd"))
		return
	}
	dstDvInd, err := strconv.Atoi(c.Param("dstDvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid dstDvInd"))
		return
	}
	accessible, err := dcgm.IsP2PAccessible(srcDvInd, dstDvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"accessible": accessible}))
}

func toHyVersionInfo(version dcgm.HyVersionInfo) HyVersionInfo {
	return HyVersionInfo{
		Major: version.Major,
		Minor: version.Minor,
	}
}

// GetHyVersion godoc
// @Summary      查询海光扩展版本
// @Description  返回 RSMI 库的海光扩展版本（hy_major/hy_minor）。对应 hy-smi -libversion 中 hy 段。
// @Tags         DevInfo
// @Produce      json
// @Success      200 {object} HyVersionInfo
// @Failure      500 {object} Response "查询失败"
// @Router       /version/hy [get]
func GetHyVersion(c *gin.Context) {
	ver, err := dcgm.HyVersion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(toHyVersionInfo(ver)))
}

// GetDevCuUtil godoc
// @Summary      查询设备 CU wave 占用周期占比
// @Description  返回设备在采样窗口内每个 CU 至少驻留 1 个 wave 的时钟周期占比（全 CU 平均）。
// @Description  duration 为采样窗口（驱动单位，实测传 1000 约等于 1s）；传 0 使用驱动默认值。
// @Tags         DevInfo
// @Produce      json
// @Param        dvInd     path  int  true  "设备索引"
// @Param        duration  query int  false "采样窗口，默认 0"
// @Success      200 {object} Response
// @Failure      400 {object} Response
// @Failure      500 {object} Response
// @Router       /device/cuutil/{dvInd} [get]
func GetDevCuUtil(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid dvInd"))
		return
	}
	duration, _ := strconv.Atoi(c.DefaultQuery("duration", "0"))
	percent, err := dcgm.DevCuUtil(dvInd, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"dvInd": dvInd, "cuUtil": percent}))
}

// GetDevWaveUtil godoc
// @Summary      查询设备 wave 驻留数量占比
// @Description  返回设备在采样窗口内各 CU 上驻留 wave 数量占最大可驻留 wave 数量的比例（全 CU 平均）。
// @Description  与 /device/cuutil 互补：cuUtil 看"有没有 wave"，waveUtil 看"wave 填了多满"。
// @Tags         DevInfo
// @Produce      json
// @Param        dvInd     path  int  true  "设备索引"
// @Param        duration  query int  false "采样窗口，默认 0"
// @Success      200 {object} Response
// @Failure      400 {object} Response
// @Failure      500 {object} Response
// @Router       /device/waveutil/{dvInd} [get]
func GetDevWaveUtil(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid dvInd"))
		return
	}
	duration, _ := strconv.Atoi(c.DefaultQuery("duration", "0"))
	percent, err := dcgm.DevWaveUtil(dvInd, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{"dvInd": dvInd, "waveUtil": percent}))
}

// GetVDevBusyPercent godoc
// @Summary      查询虚拟设备 compute 忙碌占比（近似 MIG_ENC_UTIL/DEC_UTIL）
// @Description  返回指定 vDevice 的整体 compute busy 百分比。
// @Description  注意：HCU 驱动不区分 enc/dec，本字段为整体 compute 占用率，用作 DCGM_FI_MIG_ENC_UTIL/DEC_UTIL 的近似值。
// @Tags         VirtualDevice
// @Produce      json
// @Param        vDvInd  path  int  true  "虚拟设备索引"
// @Success      200 {object} Response
// @Failure      400 {object} Response
// @Failure      500 {object} Response
// @Router       /vdevice/busypercent/{vDvInd} [get]
func GetVDevBusyPercent(c *gin.Context) {
	vDvInd, err := strconv.Atoi(c.Param("vDvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid vDvInd"))
		return
	}
	percent, err := dcgm.VDevBusyPercent(vDvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"vDvInd":       vDvInd,
		"busyPercent":  percent,
		"semanticNote": "compute busy percent, approximation for MIG_ENC_UTIL/DEC_UTIL",
	}))
}

// GetProcessGpuUsage godoc
// @Summary      查询进程各 GPU 占用率（近似 DCGM_FI_PROC_ENC_UTIL）
// @Description  返回指定进程在各设备上的 GPU compute 利用率（gpuUsageRate）及 VRAM 用量。
// @Description  注意：HCU 驱动无进程级 enc 专项统计，本字段为整体 compute 占用率，用作 DCGM_FI_PROC_ENC_UTIL 的近似值。
// @Tags         Process
// @Produce      json
// @Param        pid  path  int  true  "进程 ID"
// @Success      200 {object} Response
// @Failure      400 {object} Response
// @Failure      500 {object} Response
// @Router       /process/gpuusage/{pid} [get]
func GetProcessGpuUsage(c *gin.Context) {
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("Invalid pid"))
		return
	}
	info, err := dcgm.ProcessInfoByPid(uint32(pid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"processId":    info.ProcessID,
		"vramUsageMiB": info.VRAMUsageSize,
		"vramUsagePct": info.VRAMUsageRate,
		"gpuCount":     info.GPUCount,
		"gpuIndex":     info.GPUIndex,
		"gpuUsageRate": info.GPUUsageRate,
		"semanticNote": "gpu compute usage rate per device, approximation for DCGM_FI_PROC_ENC_UTIL",
	}))
}

// SetPowerCap 设置设备功率上限
// @Tags Power
// @Summary 设置设备功率上限
// @Description 设置指定设备的功率上限（单位：微瓦）。有效范围可通过 GET /DevPowerCapRange/:dvInd 查询。
// @Accept json
// @Produce json
// @Param dvInd path int true "设备索引"
// @Param body body object true "功率上限（微瓦）" example({"cap": 300000000})
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /powercap/{dvInd} [post]
func SetPowerCap(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	var req struct {
		Cap uint64 `json:"cap"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求体解析失败：需要 {\"cap\": <微瓦>}"))
		return
	}
	if err := dcgm.DevPowerCapSet(dvInd, req.Cap); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// ResetPowerCap 复位设备功率上限至默认值
// @Tags Power
// @Summary 复位设备功率上限
// @Description 将指定设备的功率上限恢复为出厂默认值。
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /powercap/{dvInd}/reset [post]
func ResetPowerCap(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	if err := dcgm.DevPowerCapReset(dvInd); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// SetPerfLevelV1 使用新版接口设置设备性能等级
// @Tags Device
// @Summary 设置设备性能等级（v1）
// @Description 使用 rsmi_dev_perf_level_set_v1 新版接口设置性能等级。接受的 level 值与旧版相同。
// @Accept json
// @Produce json
// @Param dvInd path int true "设备索引"
// @Param level query string true "性能等级 (auto, low, high, normal, manual, stable_std, stable_peak, stable_min_mclk, stable_min_sclk)"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /perflevel/v1/{dvInd} [post]
func SetPerfLevelV1(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	level := c.Query("level")
	levelConverted, err := ConvertToRSMIDevPerfLevel(level)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("无效的性能等级"))
		return
	}
	if err := dcgm.DevPerfLevelSetV1(dvInd, levelConverted); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// ── sysfs 只读 handler（第2批）────────────────────────────────────────────────

// GetHswFwVersion 获取指定设备的 HSW 固件版本
// @Tags Firmware
// @Summary 获取 HSW 固件版本
// @Description 读取 /sys/class/hsd/hsd{dvInd}/build_info，对应 hy-smi --showfwinfo HSW
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} Response{data=map[string]string}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /hswfw/{dvInd} [get]
func GetHswFwVersion(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	ver, err := dcgm.HswFwVersion(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]string{"version": ver}))
}

// GetGfxBoostInfo 获取指定设备的 GFX boost 等级信息
// @Tags Device
// @Summary 获取 GFX boost 状态
// @Description 读取 pp_gfx_boost sysfs 文件，对应 hy-smi --showboost
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} Response{data=map[string]string}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /boost/{dvInd} [get]
func GetGfxBoostInfo(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	info, err := dcgm.GfxBoostInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]string{"boost": info}))
}

// GetFdrHistory 获取指定设备的 GFX 频率限制历史
// @Tags Device
// @Summary 获取频率限制历史
// @Description 读取 pp_gfx_limit_info sysfs 文件，对应 hy-smi --showfdrhistory
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} Response{data=map[string]string}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /fdrhistory/{dvInd} [get]
func GetFdrHistory(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	history, err := dcgm.FdrHistory(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(map[string]string{"history": history}))
}

// GetMemBankEccInfo 获取指定设备的显存 bank 级 ECC 信息
// @Tags Memory
// @Summary 获取显存 ECC 信息（bank 级）
// @Description 遍历 mem_banks sysfs 目录，对应 hy-smi --showmemeccinfo
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} Response{data=map[string]string}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /memeccinfo/{dvInd} [get]
func GetMemBankEccInfo(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	info, err := dcgm.MemBankEccInfo(dvInd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(info))
}

// SetGfxBoostLevel 设置指定设备的 GFX boost 等级。
// @Tags Device
// @Summary 设置 GFX boost 等级
// @Accept json
// @Produce json
// @Param dvInd path int true "设备索引"
// @Param request body object{level=int} true "boost 等级"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 403 {object} Response
// @Failure 500 {object} Response
// @Router /boost/{dvInd} [post]
func SetGfxBoostLevel(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	var req struct {
		Level *int `json:"level" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求体解析失败：需要 {\"level\": <整数>}"))
		return
	}
	if err := dcgm.SetGfxBoostLevel(dvInd, *req.Level); err != nil {
		c.JSON(sysfsWriteStatus(err), ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// SetLowPower 设置运行时低功耗模式及自动挂起延时。
// @Tags Power
// @Summary 设置低功耗控制
// @Accept json
// @Produce json
// @Param dvInd path int true "设备索引"
// @Param request body object{enabled=bool,delayMs=int} true "至少提供一个字段"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 403 {object} Response
// @Failure 500 {object} Response
// @Router /lowpower/{dvInd} [post]
func SetLowPower(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	var req struct {
		Enabled *bool `json:"enabled"`
		DelayMs *int  `json:"delayMs"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || (req.Enabled == nil && req.DelayMs == nil) {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求体至少需要 enabled 或 delayMs 字段"))
		return
	}
	// 先设置延时再启用自动挂起，避免使用旧延时立即进入低功耗。
	if req.DelayMs != nil {
		if err := dcgm.SetLowPowerDelay(dvInd, *req.DelayMs); err != nil {
			c.JSON(sysfsWriteStatus(err), ErrorResponse(err.Error()))
			return
		}
	}
	if req.Enabled != nil {
		if err := dcgm.SetLowPowerMode(dvInd, *req.Enabled); err != nil {
			c.JSON(sysfsWriteStatus(err), ErrorResponse(err.Error()))
			return
		}
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

// SetRasControl 执行 RAS enable、disable 或 inject 命令。
// @Tags Memory
// @Summary 设置 RAS 控制
// @Accept json
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 403 {object} Response
// @Failure 500 {object} Response
// @Router /ras/{dvInd} [post]
func SetRasControl(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	var req struct {
		Action     string   `json:"action" binding:"required"`
		Block      string   `json:"block"`
		Parameters []string `json:"parameters"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求体解析失败"))
		return
	}

	switch strings.ToLower(req.Action) {
	case "enable", "disable":
		if req.Block == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse("enable/disable 操作需要 block 字段"))
			return
		}
		err = dcgm.SetRasFeature(dvInd, req.Block, strings.EqualFold(req.Action, "enable"))
	case "inject":
		err = dcgm.InjectRasError(dvInd, req.Parameters)
	default:
		c.JSON(http.StatusBadRequest, ErrorResponse("action 仅支持 enable、disable 或 inject"))
		return
	}
	if err != nil {
		c.JSON(sysfsWriteStatus(err), ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func sysfsWriteStatus(err error) int {
	if errors.Is(err, os.ErrPermission) {
		return http.StatusForbidden
	}
	return http.StatusInternalServerError
}

// GetHSLLogStatus 返回指定设备最近的 HSL 管理日志状态。
// @Tags Topology
// @Summary 获取 HSL 链路日志状态
// @Produce json
// @Param dvInd path int true "设备索引"
// @Param limit query int false "最多返回的日志条数" default(100)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /hslstatus/{dvInd} [get]
func GetHSLLogStatus(c *gin.Context) {
	dvInd, err := strconv.Atoi(c.Param("dvInd"))
	if err != nil || dvInd < 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse("设备索引无效"))
		return
	}
	limit := 100
	if value := c.Query("limit"); value != "" {
		limit, err = strconv.Atoi(value)
		if err != nil || limit <= 0 || limit > 1000 {
			c.JSON(http.StatusBadRequest, ErrorResponse("limit 必须是 1 到 1000 的整数"))
			return
		}
	}
	status, err := dcgm.GetHSLLogStatus(dvInd, limit)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.JSON(http.StatusNotFound, ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(status))
}

// GetPhysicalLinkTopology 返回 HSD 与 KFD 聚合的物理链路拓扑。
// @Tags Topology
// @Summary 获取物理链接拓扑
// @Produce json
// @Success 200 {object} Response
// @Failure 503 {object} Response
// @Router /linktopo [get]
func GetPhysicalLinkTopology(c *gin.Context) {
	topology, err := dcgm.GetPhysicalLinkTopology()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(topology))
}

// SaveDeviceConfig 保存设备可恢复配置快照。
// @Tags Configuration
// @Summary 保存设备配置
// @Accept json
// @Produce json
// @Param request body object{name=string,deviceIds=[]int} true "配置名称与设备列表；设备列表为空表示全部设备"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 403 {object} Response
// @Failure 500 {object} Response
// @Router /config/save [post]
func SaveDeviceConfig(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		DeviceIDs []int  `json:"deviceIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求体需要 name，可选 deviceIds"))
		return
	}
	document, err := dcgm.SaveDeviceConfig(req.Name, req.DeviceIDs)
	if err != nil {
		c.JSON(deviceConfigErrorStatus(err), ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(document))
}

// LoadDeviceConfig 加载并应用设备配置快照。
// @Tags Configuration
// @Summary 加载设备配置
// @Accept json
// @Produce json
// @Param request body object{name=string} true "配置名称"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 403 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /config/load [post]
func LoadDeviceConfig(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("请求体需要 name"))
		return
	}
	result, err := dcgm.LoadDeviceConfig(req.Name)
	if err != nil {
		c.JSON(deviceConfigErrorStatus(err), ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(result))
}

func deviceConfigErrorStatus(err error) int {
	switch {
	case errors.Is(err, dcgm.ErrInvalidDeviceConfigName):
		return http.StatusBadRequest
	case errors.Is(err, os.ErrNotExist):
		return http.StatusNotFound
	case errors.Is(err, os.ErrPermission):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
