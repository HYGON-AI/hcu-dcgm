/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package main

import "C"
import (
	"flag"
	"fmt"
	"strings"

	"github.com/goccy/go-json"
	"github.com/golang/glog"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

const defaultSampleDurationMs = 1000

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	dcgm.Init()
	defer dcgm.ShutDown()

	//gpuCount, err := dcgm.DeviceGetCount()
	//if err != nil {
	//	log.Fatalf("获取GPU数量失败: %v", err)
	//}
	//fmt.Printf("deviceCount: %d\n", gpuCount)
	//
	////maxMigPerGpu := 4 // 每卡最大MIG分区数
	//globalMigId := 0 // 全局递增MIG id
	//
	//for gpuIdx := 0; gpuIdx < gpuCount; gpuIdx++ {
	//	foundAny := false
	//	maxMigPerGpu, _ := dcgm.DeviceGetMaxMigDeviceCountByIndex(gpuIdx)
	//	glog.V(5).Infof("MaxMigDeviceCount:%v", maxMigPerGpu)
	//	for local := 0; local < maxMigPerGpu; local++ {
	//		migDev, err := dcgm.DeviceGetMigDeviceHandleByIndex(gpuIdx, globalMigId)
	//		if err != nil || migDev == nil {
	//			break // 当前GPU的MIG设备查完
	//		}
	//		fmt.Printf("[INFO] 物理GPU=%d, 全局MIG id=%d 存在\n", gpuIdx, globalMigId)
	//		globalMigId++
	//		foundAny = true
	//	}
	//	if !foundAny {
	//		fmt.Printf("[INFO] 物理GPU=%d 没有有效MIG设备\n", gpuIdx)
	//	}
	//}
	//mode, pendingMode, _ := dcgm.SystemMigMode()
	//glog.V(5).Infof("SystemMigMode currentMode:%v , pendingMode:%v ", mode, pendingMode)
	//migs, err := dcgm.MigInfos()
	//if err != nil {
	//	log.Fatalf("遍历MIG失败: %v", err)
	//}
	//glog.V(5).Infof("MigInfos:%v ", dataToJson(migs))
	//glog.V(5).Infof("----------------------------------------------------------------------------------------")

	//mig, err := dcgm.MigInfoByDvInd(0) // 查询第0号卡的所有MIG
	//if err != nil {
	//	log.Fatalf("遍历失败: %v", err)
	//}
	//glog.V(5).Infof("HCU %d MigInfoByDvInd:%v ", 0, dataToJson(mig))
	//glog.V(5).Infof("----------------------------------------------------------------------------------------")
	//
	//mig1, err := dcgm.MigInfoByDvInd(1) // 查询第0号卡的所有MIG
	//if err != nil {
	//	log.Fatalf("遍历失败: %v", err)
	//}
	//glog.V(5).Infof("HCU %d MigInfoByDvInd:%v ", 1, dataToJson(mig1))
	//info, err := dcgm.MigInfoByUUID("MIG-60224abc-7f66-42cd-9cca-45447756db55")
	//glog.V(5).Infof("MigInfoByUUID:%v ", dataToJson(info))
	//glog.V(5).Infof("----------------------------------------------------------------------------------------")
	//info1, err := dcgm.MigInfoByUUID("MIG-26d20f73-fb65-4abf-80a9-062d653aedef")
	//glog.V(5).Infof("MigInfoByUUID:%v ", dataToJson(info1))
	//name, err := dcgm.DeviceGetName(0, 0)
	//glog.V(5).Infof("HCU %d MIGName:%v ", 0, name)
	//configs, err := dcgm.MIGConfigs()
	//glog.V(5).Infof("MIGConfigs:%v", dataToJson(configs))
	//// 1) 自动计算 seCount
	//seCount, err := dcgm.MIGSECount()
	//if err != nil {
	//	log.Fatalf("获取 seCount 失败: %v", err)
	//}
	//fmt.Printf("seCount = %d\n", seCount)
	//
	//// 2) 调用格式化函数
	//gpuSliceCount := 2
	//memorySizeMB := uint64(32760) // 示例：16 GB
	//migName := dcgm.MIGName(seCount, gpuSliceCount, 1, memorySizeMB)
	//
	//fmt.Println("生成的 MIG 名称:", migName)

	/*---------------------*/
	//hip-stream
	// 执行测试
	//if passed := dcgm.RunBandwidthTests([]int{0}); passed {
	//	log.Println("所有GPU带宽测试通过!")
	//} else {
	//	log.Println("部分GPU带宽测试未达标!")
	//}

	/*---------------------*/
	//rocm-bandwidth-test

	//fmt.Println("===== 开始 PCIe/XHCL 带宽测试 =====")
	//
	//// 调用主测试函数
	//allPass := dcgm.RunPcieBandwidthTests()
	//
	//if allPass {
	//	fmt.Println("✅ 所有 GPU 测试通过")
	//	os.Exit(0)
	//} else {
	//	fmt.Println("❌ 存在 GPU 测试未通过")
	//	os.Exit(1)
	//}
	//
	//result := dcgm.HCUXHCLTest()
	//log.Printf("XHCL test completed with result: %d", result)

	// 请确保 resources/gemmPower 可执行且 resources/int8_cu_check.co 存在
	// 示例 runInfo（按实际硬件/环境调整）
	// 获取当前工作目录，用于相对路径
	//currentDir, err := os.Getwd()
	//if err != nil {
	//	fmt.Printf("获取当前目录失败: %v\n", err)
	//	return
	//}
	//
	//// 设置日志目录（相对于当前工作目录）
	//logDir := filepath.Join(currentDir, "logs")
	//
	//run := &dcgm.PowerInfo{
	//	TotalHCU:        2, // 两张卡
	//	BusID:           []string{"0000:01:00.0", "0000:02:00.0"},
	//	ToDriverID:      []int{0, 1},
	//	ToOAMID:         []int{0, 1},
	//	NumAID:          []int{0, 1},
	//	TargetPowerTime: 1,
	//	LogDir:          logDir,
	//}
	//
	//// 确保日志目录存在
	//if err := os.MkdirAll(logDir, 0755); err != nil {
	//	fmt.Printf("创建日志目录失败: %v\n", err)
	//	return
	//}
	//
	//// 直接调用 TargetPower，它会使用内嵌的资源
	//if err := dcgm.TargetPower(run); err != nil {
	//	fmt.Println("targetPower err:", err)
	//} else {
	//	fmt.Println("targetPower finished OK")
	//}

	/*---------------------*/
	//gemmPerf
	//dcgm.TargetStress()
	//
	//// 假设有两个 GPU: 0 和 1
	//gpuIDs := []int{0}
	//if err := dcgm.MemtestCL(gpuIDs); err != nil {
	//	fmt.Println("整体结果: 测试失败 ❌:", err)
	//} else {
	//	fmt.Println("整体结果: 测试成功 ✅")
	//}

	// 日志保存目录
	//logDir := filepath.Join(".", "edpp_logs")
	//if err := os.MkdirAll(logDir, 0755); err != nil {
	//	fmt.Printf("创建日志目录失败: %v\n", err)
	//	return
	//}
	//
	//// 创建 EdppInfo 对象
	//info := &dcgm.EdppInfo{
	//	LogDir:         logDir,
	//	TotalHCU:       2,  // 系统 GPU 数量
	//	EdppStressTime: 10, // 每个频率跑 1 分钟，可根据需要调整
	//}
	//
	//// 运行 EDPp 测试
	//dcgm.EDPpTest()
	//
	//fmt.Println("===== Demo 测试完成 =====")
	//dcgm.UMCBandwidth(0,0,10)
	//cards, err := dcgm.CardSeriesList()
	//if err != nil {
	//	fmt.Println("error:", err)
	//	return
	//}
	//glog.V(5).Infof("cards: %v", cards)
	//for _, card := range cards {
	//	fmt.Printf("卡 %d 型号: %s\n", card.DvInd, card.SeriesName)
	//}
	//dvIds := []int{0, 1}
	//bwMap, err := dcgm.BandwidthTestAPI(dvIds)
	//if err != nil {
	//	log.Printf("带宽测试失败: %v", err)
	//	return
	//}
	//glog.V(5).Infof("bwMap:%v", bwMap)
	//fmt.Println("===== 带宽测试结果 =====")
	//for id, bw := range bwMap {
	//	fmt.Printf("设备 %d: %.2f GB/s\n", id, bw)
	//}
	//
	//fmt.Println("所有设备测试完成 ✅")
	//bandwidthInfo, err := dcgm.PcieBw(0)
	//if err != nil {
	//	log.Printf("带宽测试失败: %v", err)
	//	return
	//}
	//glog.V(5).Infof("bandwidthInfo:%v", bandwidthInfo)

	//tests := []string{
	//	"GPU-AX100,rev2",
	//	"GPU-BX200",
	//	",leading-comma",
	//	"",
	//	"no,comma,more",
	//	"no,comma,more",
	//	"BW200_H, UBB BW100_H",
	//}
	//
	//for _, t := range tests {
	//	fmt.Printf("in: %q -> out: %q\n", t, dcgm.NormalizeDevTypeName(t))
	//}
	//dcgm.ProcessInfoByPid()
	//id, _ := dcgm.GetDeviceUniqueId(0)
	//glog.V(5).Infof("dcInd : 0, id: %v", id)
	//uniqueId, _ := dcgm.GetDeviceUniqueId(1)
	//glog.V(5).Infof("dcInd : 1, id: %v", uniqueId)
	//虚拟设备百分比
	// 捕获中断信号，安全退出
	//stop := make(chan os.Signal, 1)
	//signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	//ticker := time.NewTicker(2 * time.Second) // 每隔 2 秒
	//defer ticker.Stop()
	//
	//for {
	//	select {
	//	case <-ticker.C:
	//		glog.Info("VDevBusyPercent for vdevice 0")
	//		percent, _ := dcgm.VDevBusyPercent(0)
	//		glog.V(5).Infof("VDevBusyPercent: %v", percent)
	//
	//	case <-stop:
	//		glog.Info("Received stop signal, exiting...")
	//		return
	//	}
	//}
	//devices, err := dcgm.NumMonitorDevices()
	//if err != nil {
	//	return
	//}
	//
	//// 构建邻接表
	//matrix := make([][]string, devices)
	//for i := 0; i < devices; i++ {
	//	matrix[i] = make([]string, devices)
	//	for j := 0; j < devices; j++ {
	//		matrix[i][j] = "-" // 默认
	//	}
	//}
	//
	//// 填充远端 BDFID
	//for src := 0; src < devices; src++ {
	//	links, err := dcgm.DumpXhclRemoteBdfids(src)
	//	if err != nil {
	//		glog.Warningf("DumpXhclRemoteBdfids failed for HCU %d: %v", src, err)
	//		continue
	//	}
	//
	//	for _, link := range links {
	//		remoteBdf := link.BdfID
	//		bus := (remoteBdf >> 8) & 0xff
	//		device := (remoteBdf >> 3) & 0x1f
	//		function := remoteBdf & 0x7
	//
	//		// 用 LinkID 填入矩阵
	//		matrix[src][link.LinkID] = fmt.Sprintf("%02x:%02x.%x", bus, device, function)
	//	}
	//}
	//
	//// 打印矩阵（格式对齐）
	//fmt.Println("\nDCU ↔ HCU XHCL Neighbor Table:\n")
	//
	//// 打印表头
	//fmt.Printf("%10s", "")
	//for j := 0; j < devices; j++ {
	//	fmt.Printf("  HCU[%d]", j)
	//}
	//fmt.Println()
	//
	//// 打印每行
	//for i := 0; i < devices; i++ {
	//	fmt.Printf("HCU[%d]", i)
	//	for j := 0; j < devices; j++ {
	//		fmt.Printf("  %8s", matrix[i][j])
	//	}
	//	fmt.Println()
	//}
	//
	//return

	// ---------- Demo 参数 ----------
	numDevices, err := dcgm.NumMonitorDevices()
	if err != nil {
		glog.Errorf("Failed to get number of devices: %v", err)
		return
	}

	glog.V(5).Infof("Found %d devices", numDevices)

	// 遍历每个 GPU
	//for dvInd := 0; dvInd < numDevices; dvInd++ {
	//	// 调用 PicBusInfo 获取 PCI 总线信息
	//	picID, err := dcgm.PicBusInfo(dvInd)
	//	if err != nil {
	//		glog.V(5).Infof("Device %d: Failed to get PCI bus info: %v", dvInd, err)
	//	} else {
	//		glog.V(5).Infof("Device %d: PCI Bus ID: %s", dvInd, picID)
	//	}
	//
	//	glog.V(5).Infof("----------------------------------------")
	//
	//	// 调用 DumpXhclRemoteBdfids 获取 XHCL 远端 BDFID
	//	bdfList, err := dcgm.DumpXhclRemoteBdfids(dvInd)
	//	if err != nil {
	//		glog.V(5).Infof("Device %d: Failed to dump XHCL remote BDF IDs: %v", dvInd, err)
	//	} else {
	//		glog.V(5).Infof("Device %d: XHCL remote BDF count: %d", dvInd, len(bdfList))
	//		for _, bdf := range bdfList {
	//			glog.V(5).Infof("✈️ LinkID: %v  BdfID: 0x%x", bdf.LinkID, bdf.BdfID)
	//		}
	//	}
	//
	//	glog.V(5).Infof("----------------------------------------")
	//
	//	// 构造一个不包含当前设备索引的切片
	//	otherIndices := make([]int, 0, numDevices-1)
	//	for j := 0; j < numDevices; j++ {
	//		if j != dvInd {
	//			otherIndices = append(otherIndices, j)
	//		}
	//	}
	//
	//	glog.V(5).Infof("[DEMO] Start testing GetGpuInterconnectInfo, dvInd=%v", dvInd)
	//
	//	// ---------- 调用核心函数 ----------
	//	info, err := dcgm.GetGpuInterconnectInfo(dvInd, otherIndices)
	//	if err != nil {
	//		glog.Errorf("[DEMO] GetGpuInterconnectInfo failed, dvInd=%v, err=%v", dvInd, err)
	//		continue
	//	}
	//	glog.V(5).Infof("info: %v", dataToJson(info))
	//
	//	// ---------- 打印结果（结构化输出） ----------
	//	glog.V(5).Infof("[DEMO] GetGpuInterconnectInfo success, dvInd=%v", dvInd)
	//
	//	fmt.Println("========== GPU Interconnect Info ==========")
	//	fmt.Printf("GPU Index     : %d\n", dvInd)
	//	fmt.Printf("Card Count    : %d\n", info.CardCount)
	//	fmt.Printf("LinkType      : %s\n", info.LinkType)
	//	fmt.Printf("Direct Links  : %d\n", len(info.Links))
	//	fmt.Println("-------------------------------------------")
	//
	//	for idx, link := range info.Links {
	//		fmt.Printf("[%d] RemoteDvInd: %d, RemoteBDFID: 0x%x, Weight: %d, LinkType: %s\n",
	//			idx, link.RemoteDvInd, link.RemoteBdfID, link.Weight, link.LinkType)
	//	}
	//
	//	fmt.Println("===========================================")
	//
	//	glog.V(5).Infof("[DEMO] Test finished for dvInd=%v", dvInd)
	//	glog.V(5).Infof("info: %v", dataToJson(info))
	//}
	//dvIdList := make([]int, numDevices)
	//for i := 0; i < numDevices; i++ {
	//	dvIdList[i] = i
	//}
	//infos, err := dcgm.ShowNumaTopology(dvIdList)
	//glog.V(5).Infof("ShowNumaTopology: %v", dataToJson(infos))

	//demoDeviceIdentityProbe(numDevices)
	//demoCardSeriesList(numDevices)
	//demoDeviceUtilization(numDevices)

	//demoGetDeviceId()
	//demoHCUHealthCheck()
	demoDeviceTemperatureInfo(numDevices)
	demoDeviceUtilizationInfo(numDevices)

	// fmt.Println("==== HCU Interconnect Topology Demo ====")
	//
	// matrix, err := dcgm.DiscoverInterconnectTopology()
	// if err != nil {
	// 	fmt.Printf("DiscoverInterconnectTopology failed: %v\n", err)
	// 	return
	// }
	// glog.V(5).Infof("DiscoverInterconnectTopology: %v", dataToJson(matrix))
	//
	// hcuCount := matrix.DeviceCount
	// fmt.Printf("Total HCU count: %d\n\n", hcuCount)
	//
	// for src := 0; src < hcuCount; src++ {
	// 	fmt.Printf("From HCU %d:\n", src)
	// 	for dst := 0; dst < hcuCount; dst++ {
	// 		info := matrix.Matrix[src][dst]
	//
	// 		fmt.Printf(
	// 			"  -> HCU %-2d | LinkType: %-12s | Weight: %d\n",
	// 			info.DstDvInd,
	// 			info.LinkType,
	// 			info.Weight,
	// 		)
	// 	}
	// 	fmt.Println()
	// }
	//
	// fmt.Println("==== Demo Finished ====")
	//
	// if hcuCount > 3 {
	// 	src, dst := 0, 3
	// 	linkInfo := matrix.Matrix[src][dst]
	// 	fmt.Printf("HCU %d -> HCU %d : LinkType=%s, Weight=%d, Hops=%d, PciID=%s\n",
	// 		src, dst, linkInfo.LinkType, linkInfo.Weight, linkInfo.Hops, linkInfo.PciID)
	// }
}

// demoDeviceIdentityProbe 打印驱动用于识别卡型的关键字段。
func demoDeviceIdentityProbe(numDevices int) {
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("  Device Identity Probe")
	fmt.Println("============================================================")
	for dvInd := 0; dvInd < numDevices; dvInd++ {
		devTypeName, cuCount, err := dcgm.DevTypeName(dvInd)
		if err != nil {
			fmt.Printf("HCU %d: DevTypeName ERROR: %v\n", dvInd, err)
		} else {
			fmt.Printf("HCU %d: DevTypeName=%q, ComputeUnitCount=%.0f\n", dvInd, devTypeName, cuCount)
		}

		devTypeID, err := dcgm.DevTypeID(dvInd)
		if err != nil {
			fmt.Printf("HCU %d: DevTypeID ERROR: %v\n", dvInd, err)
		} else {
			fmt.Printf("HCU %d: DevTypeID=%s\n", dvInd, devTypeID)
		}

		deviceInfo, err := dcgm.GetDeviceInfo(dvInd)
		if err != nil {
			fmt.Printf("HCU %d: GetDeviceInfo ERROR: %v\n", dvInd, err)
		} else {
			fmt.Printf("HCU %d: GetDeviceInfo.Name=%q, ComputeUnitCount=%d\n", dvInd, deviceInfo.Name, deviceInfo.ComputeUnitCount)
		}

		uniqueID, err := dcgm.GetDeviceUniqueId(dvInd)
		if err != nil {
			fmt.Printf("HCU %d: GetDeviceUniqueId ERROR: %v\n", dvInd, err)
		} else {
			fmt.Printf("HCU %d: UniqueId=%s\n", dvInd, uniqueID)
		}

		pciID, err := dcgm.PicBusInfo(dvInd)
		if err != nil {
			fmt.Printf("HCU %d: PicBusInfo ERROR: %v\n", dvInd, err)
		} else {
			fmt.Printf("HCU %d: PCI Bus ID=%s\n", dvInd, pciID)
		}
		fmt.Println("------------------------------------------------------------")
	}
	fmt.Println("============================================================")
}

// demoCardSeriesList 打印 EDPp 后端识别依赖的卡型信息。
func demoCardSeriesList(numDevices int) {
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("  CardSeriesList Probe")
	fmt.Println("============================================================")
	fmt.Printf("NumMonitorDevices: %d\n", numDevices)

	cards, err := dcgm.CardSeriesList()
	if err != nil {
		fmt.Printf("CardSeriesList ERROR: %v\n", err)
		glog.Errorf("CardSeriesList failed: %v", err)
		return
	}

	fmt.Printf("CardSeriesList count: %d\n", len(cards))
	fmt.Printf("CardSeriesList raw: %s\n", dataToJson(cards))
	for _, card := range cards {
		devTypeID, idErr := dcgm.DevTypeID(card.DvInd)
		if idErr != nil {
			fmt.Printf("HCU %d: SeriesName=%q, DevTypeID ERROR=%v, expectedEDPpBackend=%s\n", card.DvInd, card.SeriesName, idErr, expectedEDPpBackend(card.SeriesName, ""))
			continue
		}
		backend := expectedEDPpBackend(card.SeriesName, devTypeID)
		fmt.Printf("HCU %d: SeriesName=%q, DevTypeID=%s, expectedEDPpBackend=%s\n", card.DvInd, card.SeriesName, devTypeID, backend)
	}

	if len(cards) != numDevices {
		fmt.Printf("WARNING: CardSeriesList count (%d) != NumMonitorDevices (%d)\n", len(cards), numDevices)
	}
	fmt.Println("============================================================")
}

func expectedEDPpBackend(seriesName, devTypeID string) string {
	if backend := expectedEDPpBackendBySeries(seriesName); backend != "legacy" {
		return backend
	}
	return expectedEDPpBackendByDevTypeID(devTypeID)
}

func expectedEDPpBackendBySeries(seriesName string) string {
	series := strings.ToUpper(strings.TrimSpace(seriesName))
	switch {
	case strings.Contains(series, "BW1100"),
		strings.Contains(series, "BW1102"),
		strings.Contains(series, "BW1200"):
		return "nmz-gfx938"
	case strings.Contains(series, "BW100"),
		strings.Contains(series, "BW101"),
		strings.Contains(series, "BW150"),
		strings.Contains(series, "BW151"),
		strings.Contains(series, "BW200"),
		strings.Contains(series, "BW1000"):
		return "bmz"
	default:
		return "legacy"
	}
}

func expectedEDPpBackendByDevTypeID(devTypeID string) string {
	id := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(devTypeID)), "0x")
	switch id {
	case "6430":
		return "nmz-gfx938"
	case "6320":
		return "bmz"
	default:
		return "legacy"
	}
}

// demoDeviceUtilization 演示五个设备利用率相关 API 的调用与输出。
func demoDeviceUtilization(numDevices int) {
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("  Device Utilization APIs Demo")
	fmt.Println("  (HCUCuUsage / HCUSampledUsage / HCUCUSampledUsage / HCUWaveSampledUsage / HCUSEUsage)")
	fmt.Println("============================================================")
	fmt.Printf("Detected %d HCU(s), sample window = %d ms\n\n", numDevices, defaultSampleDurationMs)

	//for dvInd := 0; dvInd < numDevices; dvInd++ {
	//	fmt.Printf("-------------------- HCU [%d] --------------------\n", dvInd)
	//
	//	// 1. HCUCuUsage — 瞬时 HCU 占用率
	//	fmt.Println("[1/5] HCUCuUsage")
	//	fmt.Println("      API   : dcgm.HCUCuUsage(dvInd)")
	//	fmt.Println("      RSMI  : rsmi_dev_cu_usage_get")
	//	fmt.Println("      hy-smi: -u (瞬时 HCU 占用率，活跃 CU 数 / CU 总数)")
	//	rate, err := dcgm.HCUCuUsage(dvInd)
	//	if err != nil {
	//		fmt.Printf("      ERROR : %v\n", err)
	//		glog.Errorf("HCU %d HCUCuUsage failed: %v", dvInd, err)
	//	} else {
	//		fmt.Printf("      RESULT: %.2f%%\n", rate)
	//		glog.V(5).Infof("HCU %d HCUCuUsage = %.2f%%", dvInd, rate)
	//	}
	//	fmt.Println()
	//
	//	// 2. HCUSampledUsage — 采样窗口内 HCU 活跃占比
	//	fmt.Println("[2/5] HCUSampledUsage")
	//	fmt.Printf("      API   : dcgm.HCUSampledUsage(dvInd, %d)\n", defaultSampleDurationMs)
	//	fmt.Println("      RSMI  : rsmi_dev_hcu_util_get")
	//	fmt.Println("      hy-smi: --showhcuutil (采样窗口内 HCU 活跃次数占比)")
	//	rate, err = dcgm.HCUSampledUsage(dvInd, defaultSampleDurationMs)
	//	if err != nil {
	//		fmt.Printf("      ERROR : %v\n", err)
	//		glog.Errorf("HCU %d HCUSampledUsage failed: %v", dvInd, err)
	//	} else {
	//		fmt.Printf("      RESULT: %.2f%% (over %d ms)\n", rate, defaultSampleDurationMs)
	//		glog.V(5).Infof("HCU %d HCUSampledUsage = %.2f%%", dvInd, rate)
	//	}
	//	fmt.Println()
	//
	//	// 3. HCUCUSampledUsage — 采样窗口内 CU 活跃占比（全 CU 平均）
	//	fmt.Println("[3/5] HCUCUSampledUsage")
	//	fmt.Printf("      API   : dcgm.HCUCUSampledUsage(dvInd, %d)\n", defaultSampleDurationMs)
	//	fmt.Println("      RSMI  : rsmi_dev_cu_util_get")
	//	fmt.Println("      hy-smi: --showcuutil (各 CU 至少 1 个 wave 的周期占比，取平均)")
	//	rate, err = dcgm.HCUCUSampledUsage(dvInd, defaultSampleDurationMs)
	//	if err != nil {
	//		fmt.Printf("      ERROR : %v\n", err)
	//		glog.Errorf("HCU %d HCUCUSampledUsage failed: %v", dvInd, err)
	//	} else {
	//		fmt.Printf("      RESULT: %.2f%% (over %d ms)\n", rate, defaultSampleDurationMs)
	//		glog.V(5).Infof("HCU %d HCUCUSampledUsage = %.2f%%", dvInd, rate)
	//	}
	//	fmt.Println()
	//
	//	// 4. HCUWaveSampledUsage — 采样窗口内 wave 驻留占比（全 CU 平均）
	//	fmt.Println("[4/5] HCUWaveSampledUsage")
	//	fmt.Printf("      API   : dcgm.HCUWaveSampledUsage(dvInd, %d)\n", defaultSampleDurationMs)
	//	fmt.Println("      RSMI  : rsmi_dev_wave_util_get")
	//	fmt.Println("      hy-smi: --showwaveutil (各 CU 上 wave 驻留数量占比，取平均)")
	//	rate, err = dcgm.HCUWaveSampledUsage(dvInd, defaultSampleDurationMs)
	//	if err != nil {
	//		fmt.Printf("      ERROR : %v\n", err)
	//		glog.Errorf("HCU %d HCUWaveSampledUsage failed: %v", dvInd, err)
	//	} else {
	//		fmt.Printf("      RESULT: %.2f%% (over %d ms)\n", rate, defaultSampleDurationMs)
	//		glog.V(5).Infof("HCU %d HCUWaveSampledUsage = %.2f%%", dvInd, rate)
	//	}
	//	fmt.Println()
	//
	//	// 5. HCUSEUsage — 各 SE 瞬时 CU 占用率
	//	fmt.Println("[5/5] HCUSEUsage")
	//	fmt.Println("      API   : dcgm.HCUSEUsage(dvInd)")
	//	fmt.Println("      RSMI  : rsmi_dev_se_util_get")
	//	fmt.Println("      hy-smi: --showseuse (各 Shader Engine 活跃 CU 占比，瞬时值)")
	//	seUsage, err := dcgm.HCUSEUsage(dvInd)
	//	if err != nil {
	//		fmt.Printf("      ERROR : %v\n", err)
	//		glog.Errorf("HCU %d HCUSEUsage failed: %v", dvInd, err)
	//	} else {
	//		fmt.Println("      RESULT: SE utilization (percent per Shader Engine):")
	//		for i := 0; i < dcgm.MAX_SE_CNT; i++ {
	//			fmt.Printf("        SE[%d] = %.2f%%\n", i, seUsage.Percent[i])
	//		}
	//		glog.V(5).Infof("HCU %d HCUSEUsage = %s", dvInd, dataToJson(seUsage))
	//	}
	//	fmt.Println()
	//}

	// HCUMclk: 读取指定卡的 memory clock 频率
	//fmt.Println("============================================================")
	//fmt.Println("  HCUMclk - Memory Clock Frequency")
	//fmt.Println("============================================================")
	//for _, dvInd := range []int{0, 1} {
	//	fmt.Printf("  HCU %d:\n", dvInd)
	//	fmt.Println("      API   : dcgm.HCUMclk(dvInd)")
	//	fmt.Println("      RSMI  : rsmi_dev_gpu_clk_freq_get(RSMI_CLK_TYPE_MEM)")
	//	mclk, err := dcgm.HCUMclk(dvInd)
	//	if err != nil {
	//		fmt.Printf("      ERROR : %v\n", err)
	//		glog.Errorf("HCU %d HCUMclk failed: %v", dvInd, err)
	//	} else {
	//		fmt.Printf("      RESULT: mclk = %.2f MHz\n", mclk)
	//		glog.V(5).Infof("HCU %d HCUMclk = %.2f MHz", dvInd, mclk)
	//	}
	//	fmt.Println()
	//}

}

func demoGetDeviceId() {
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("  GetDeviceId Demo (HCU 0–7)")
	fmt.Println("============================================================")
	for dvInd := 0; dvInd < 2; dvInd++ {
		deviceId, err := dcgm.GetDeviceId(dvInd)
		if err != nil {
			fmt.Printf("HCU %d: ERROR: %v\n", dvInd, err)
			glog.Errorf("HCU %d GetDeviceId failed: %v", dvInd, err)
		} else {
			fmt.Printf("HCU %d: DeviceId = %s\n", dvInd, deviceId)
			glog.V(5).Infof("HCU %d GetDeviceId = %s", dvInd, deviceId)
		}
	}
	fmt.Println("============================================================")
}

func dataToJson(data any) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error serializing to JSON:", err)
	}
	return string(jsonData)
}

func demoDeviceTemperatureInfo(numDevices int) {
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("  GetDeviceTemperatureInfo Demo")
	fmt.Println("============================================================")

	for dvInd := 0; dvInd < numDevices; dvInd++ {
		temperatures, err := dcgm.GetDeviceTemperatureInfo(dvInd)
		fmt.Printf("HCU[%d] Edge=%.2f°C  Junction=%.2f°C  Memory=%.2f°C  Core=%.2f°C\n",
			dvInd,
			temperatures.Edge,
			temperatures.Junction,
			temperatures.Memory,
			temperatures.Core,
		)
		if err != nil {
			fmt.Printf("HCU[%d] PARTIAL ERROR: %v\n", dvInd, err)
			glog.Warningf("GetDeviceTemperatureInfo partially failed for HCU %d: %v", dvInd, err)
			continue
		}
		glog.V(5).Infof("HCU %d temperature info: %s", dvInd, dataToJson(temperatures))
	}

	fmt.Println("============================================================")
}

func demoHCUHealthCheck() {
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("  HCUHealthCheck - hy-smi --healthcheck")
	fmt.Println("============================================================")
	statuses, err := dcgm.HCUHealthCheck()
	if err != nil {
		fmt.Printf("  ERROR: %v\n", err)
		glog.Errorf("HCUHealthCheck failed: %v", err)
	} else {
		for _, s := range statuses {
			fmt.Printf("  HCU[%d]  BusId: %-20s  Status: %s\n", s.HCU, s.BusId, s.Status)
		}
		glog.V(5).Infof("HCUHealthCheck result: %s", dataToJson(statuses))
	}
	fmt.Println("============================================================")
}

func demoDeviceUtilizationInfo(numDevices int) {
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("  DeviceUtilizationInfoGet Demo")
	fmt.Println("============================================================")

	for dvInd := 0; dvInd < numDevices; dvInd++ {
		info, err := dcgm.DeviceUtilizationInfoGet(dvInd)
		if err != nil {
			fmt.Printf("HCU[%d] ERROR: %v\n", dvInd, err)
			glog.Errorf("DeviceUtilizationInfoGet failed for HCU %d: %v", dvInd, err)
			continue
		}
		if !info.Supported {
			fmt.Printf("HCU[%d] 粗粒度利用率计数器不支持（硬件/驱动未实现）\n", dvInd)
			continue
		}
		fmt.Printf("HCU[%d] GfxActivity=%d  MemActivity=%d  Timestamp=%d ns\n",
			dvInd, info.GfxActivity, info.MemActivity, info.Timestamp)
		glog.V(5).Infof("HCU %d utilization info: %s", dvInd, dataToJson(info))
	}

	fmt.Println("============================================================")
}
