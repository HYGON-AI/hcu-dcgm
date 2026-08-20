/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/golang/glog"
)

const (
	DiagLevel1 = 1
	DiagLevel2 = 2
	DiagLevel3 = 3
	DiagLevel4 = 4
)

func runDiag(level int) (diagResults DiagResults, err error) {
	// 保证诊断执行完后自动清理 stop 状态
	defer resetDiagStop()
	devCount, _ := listFilesInDevDri()
	numDevices, err := NumMonitorDevices()
	if numDevices == devCount {
		glog.V(5).Infof("HCU initialization is complete:%v", numDevices)
	} else {
		glog.V(5).Infof("HCU initialization is not complete: expected %v devices, but found %v devices", devCount, numDevices)
	}
	diagResults.DeviceNumber = numDevices

	// 硬件诊断
	diagHardwareResults, _ := diagHardware(level)
	diagResults.PerHCU = diagHardwareResults

	// 在硬件阶段后检查 stop 请求 —— 若已请求，则在返回前清标志并返回已有结果
	if IsDiagStopped() {
		glog.V(3).Infof("runDiag: stop requested after hardware checks, aborting further tests")
		resetDiagStop()
		return diagResults, nil
	}

	// 软件诊断
	diagSoftwareResults, _ := diagSoftware(level)
	diagResults.Software = diagSoftwareResults

	// 在软件阶段后检查 stop 请求
	if IsDiagStopped() {
		glog.V(3).Infof("runDiag: stop requested after software checks, aborting stress tests")
		resetDiagStop()
		return diagResults, nil
	}

	// 运行压力/稳定性测试并将结构化结果合并（runStressTests 不会并发启动每项）
	stressResults, serr := runStressTests(level)
	if serr != nil {
		// 记录错误但仍合并已得到的结果
		glog.Errorf("runStressTests error: %v", serr)
		diagResults.Software = append(diagResults.Software, DiagResult{
			Status:       DiagResultWarn,
			TestName:     "runStressTests",
			TestOutput:   "",
			ErrorCode:    -1,
			ErrorMessage: fmt.Sprintf("runStressTests error: %v", serr),
		})
	}

	// 把压力测试结果合并进主 diagResults
	mergeStressIntoDiag(&diagResults, stressResults)

	// run 完成前确保停止标志被清除（无论是正常结束还是被 stop）
	resetDiagStop()

	glog.V(5).Infof("diagResults:%v", diagResults)
	return diagResults, serr
}

// 硬件诊断
func diagHardware(level int) (hcuResults []HCUResult, err error) {
	numDevices, _ := NumMonitorDevices()
	for i := 0; i < numDevices; i++ {
		var hcuResult HCUResult
		hcuResult.HCU = i

		// ===== 水平一诊断 =====
		if level >= DiagLevel1 {
			// Memory Check
			memoryCap, _ := rsmiDevMemoryTotalGet(i, RSMI_MEM_TYPE_FIRST)
			memoryUsed, _ := rsmiDevMemoryUsageGet(i, RSMI_MEM_TYPE_FIRST)
			percent := float64(memoryUsed) / float64(memoryCap) * 100

			eccErrorDetails := ""
			eccHealthy := true
			if eccBlocksInfo, err := EccBlocksInfo(i); err == nil {
				for _, block := range eccBlocksInfo {
					if block.CE > 0 || block.UE > 0 {
						eccHealthy = false
						eccErrorDetails += fmt.Sprintf("Block: %s, State: %s, CE: %d, UE: %d; ", block.Block, block.State, block.CE, block.UE)
					}
				}
			} else {
				eccHealthy = false
				eccErrorDetails = fmt.Sprintf("Error fetching ECC blocks info: %v", err)
			}
			if eccHealthy {
				eccErrorDetails = "All ECC blocks are healthy."
			}

			status := DiagResultPass
			errorMsg := ""
			if memoryUsed >= memoryCap || !eccHealthy {
				status = DiagResultWarn
				errorMsg = fmt.Sprintf("Issues detected. Memory usage: %.2f%%. %s", percent, eccErrorDetails)
			}

			hcuResult.DiagResults = append(hcuResult.DiagResults, DiagResult{
				Status:       status,
				TestName:     "Memory Check",
				TestOutput:   fmt.Sprintf("Total: %.2fG, Used: %.2fG, Usage: %.2f%% | ECC Status: %s", bytesToGB(memoryCap), bytesToGB(memoryUsed), percent, eccErrorDetails),
				ErrorCode:    0,
				ErrorMessage: errorMsg,
			})

			// picBus Check
			if picBusInfo, err := PicBusInfo(i); err != nil {
				hcuResult.DiagResults = append(hcuResult.DiagResults, DiagResult{
					Status:       DiagResultFail,
					TestName:     "picBus Check",
					TestOutput:   "",
					ErrorCode:    -1,
					ErrorMessage: fmt.Sprintf("Error fetching bus info: %v", err),
				})
				hcuResult.RC = -1
			} else {
				hcuResult.DiagResults = append(hcuResult.DiagResults, DiagResult{
					Status:     DiagResultPass,
					TestName:   "picBus Check",
					TestOutput: picBusInfo,
					ErrorCode:  0,
				})
			}

			// Power Check
			maxPower, _ := MaxPower(i)
			power, _ := Power(i)
			status = DiagResultPass
			errorMsg = ""
			if power < 0 || power > maxPower {
				status = DiagResultWarn
				errorMsg = fmt.Sprintf("Power %d is out of range [0, %d]", power, maxPower)
			}
			hcuResult.DiagResults = append(hcuResult.DiagResults, DiagResult{
				Status:       status,
				TestName:     "Power Check",
				TestOutput:   fmt.Sprintf("Max Power: %d, Current Power: %d", maxPower, power),
				ErrorCode:    0,
				ErrorMessage: errorMsg,
			})
		}

		// ===== 水平二诊断 =====
		if level >= DiagLevel2 {
			// PCIe Bandwidth Check
			if rsmiPcieBandwidth, err := DevPciBandwidth(i); err != nil {
				hcuResult.DiagResults = append(hcuResult.DiagResults, DiagResult{
					Status:       DiagResultFail,
					TestName:     "PCIe Bandwidth Check",
					TestOutput:   "",
					ErrorCode:    -1,
					ErrorMessage: fmt.Sprintf("Error fetching PCIe bandwidth: %v", err),
				})
				hcuResult.RC = -1
			} else {
				hcuResult.DiagResults = append(hcuResult.DiagResults, DiagResult{
					Status:   DiagResultPass,
					TestName: "PCIe Bandwidth Check",
					TestOutput: fmt.Sprintf("PCIe Bandwidth Test Results | Supported Transfer Rates: %d, Current Transfer Rate: %d | Frequencies: %v | PCIe Lanes: %v",
						rsmiPcieBandwidth.TransferRate.NumSupported,
						rsmiPcieBandwidth.TransferRate.Current,
						rsmiPcieBandwidth.TransferRate.Frequency[:5],
						rsmiPcieBandwidth.Lanes[:5]),
					ErrorCode: 0,
				})
			}

		}

		hcuResults = append(hcuResults, hcuResult)
	}
	return
}

// 软件诊断
func diagSoftware(level int) (softwareDiagResults []DiagResult, err error) {
	numDevices, _ := NumMonitorDevices()

	if level >= DiagLevel1 {
		// DTK Version Check
		if dtkVersion, err := DTKVersion(); err != nil {
			softwareDiagResults = append(softwareDiagResults, DiagResult{
				Status:       DiagResultFail,
				TestName:     "DTK Version Check",
				TestOutput:   "",
				ErrorCode:    -1,
				ErrorMessage: fmt.Sprintf("Error fetching DTK version: %v", err),
			})
		} else {
			softwareDiagResults = append(softwareDiagResults, DiagResult{
				Status:     DiagResultPass,
				TestName:   "DTK Version Check",
				TestOutput: dtkVersion,
				ErrorCode:  0,
			})
		}

		// Driver Version Check
		if version, err := Version(RSMISwCompFirst); err != nil {
			softwareDiagResults = append(softwareDiagResults, DiagResult{
				Status:       DiagResultFail,
				TestName:     "Driver Version Check",
				TestOutput:   "",
				ErrorCode:    -1,
				ErrorMessage: fmt.Sprintf("Error fetching driver version: %v", err),
			})
		} else {
			softwareDiagResults = append(softwareDiagResults, DiagResult{
				Status:     DiagResultPass,
				TestName:   "Driver Version Check",
				TestOutput: version,
				ErrorCode:  0,
			})
		}

		// RSMI Version Check
		if rsmiVersion, err := HCUVersion(); err != nil {
			softwareDiagResults = append(softwareDiagResults, DiagResult{
				Status:       DiagResultFail,
				TestName:     "RSMI Version Check",
				TestOutput:   "",
				ErrorCode:    -1,
				ErrorMessage: fmt.Sprintf("Error fetching RSMI version: %v", err),
			})
		} else {
			output := fmt.Sprintf("RSMI Version: Major %d, Minor %d, Patch %d, Build %s",
				rsmiVersion.Major, rsmiVersion.Minor, rsmiVersion.Patch, rsmiVersion.Build)
			softwareDiagResults = append(softwareDiagResults, DiagResult{
				Status:     DiagResultPass,
				TestName:   "RSMI Version Check",
				TestOutput: output,
				ErrorCode:  0,
			})
		}

		// VBIOS & Compatibility Check
		for i := 0; i < numDevices; i++ {
			if vbios, err := VbiosVersion(i); err != nil {
				softwareDiagResults = append(softwareDiagResults, DiagResult{
					Status:       DiagResultFail,
					TestName:     fmt.Sprintf("Device %d VBIOS Version Check", i),
					TestOutput:   "",
					ErrorCode:    -1,
					ErrorMessage: fmt.Sprintf("Error fetching VBIOS version for device %d: %v", i, err),
				})
			} else {
				softwareDiagResults = append(softwareDiagResults, DiagResult{
					Status:     DiagResultPass,
					TestName:   fmt.Sprintf("Device %d VBIOS Version Check", i),
					TestOutput: fmt.Sprintf("VBIOS version for device %d: %v", i, vbios),
					ErrorCode:  0,
				})
			}

			devTypeId, _ := rsmiDevIdGet(i)
			devTypeName := type2name[fmt.Sprintf("%x", devTypeId)]
			dtkVersion, dtkErr := DTKVersion()
			driverVersion, driverErr := Version(RSMISwCompFirst)
			if dtkErr != nil || driverErr != nil {
				softwareDiagResults = append(softwareDiagResults, compatibilityDiagResult(devTypeName, CompatibilityResult{}, fmt.Errorf("获取兼容性版本失败: Driver=%v, DTK=%v", driverErr, dtkErr)))
				continue
			}
			compatibilityResult, compatibilityErr := Compatible(devTypeName, driverVersion, dtkVersion)
			softwareDiagResults = append(softwareDiagResults, compatibilityDiagResult(devTypeName, compatibilityResult, compatibilityErr))
		}
	}

	return
}

func compatibilityDiagResult(cardModel string, result CompatibilityResult, err error) DiagResult {
	testName := fmt.Sprintf("card %v Compatibility Check", cardModel)
	if err != nil {
		return DiagResult{
			Status:       DiagResultFail,
			TestName:     testName,
			ErrorCode:    -1,
			ErrorMessage: fmt.Sprintf("Compatibility check failed for card %v: %v", cardModel, err),
		}
	}
	if result.Status == CompatibilityWarn {
		return DiagResult{
			Status:   DiagResultWarn,
			TestName: testName,
			TestOutput: fmt.Sprintf(
				"Current Driver: %s\nCurrent DTK: %s\nRecommended Driver: %s\nRecommended DTK: %s",
				result.CurrentDriver,
				result.CurrentDTK,
				result.RecommendedDriver,
				result.RecommendedDTK,
			),
		}
	}
	return DiagResult{
		Status:     DiagResultPass,
		TestName:   testName,
		TestOutput: fmt.Sprintf("card %v is compatible", cardModel),
	}
}

// runStressTests 按等级依次执行压测项并返回结构化的 DiagResults。
// 特性：在每一项开始前检查 IsDiagStopped()，如果为 true 则不再启动该项及后续项。
// 在提前返回前会清除停止标志 resetDiagStop()。
func runStressTests(level int) (DiagResults, error) {
	var aggregated DiagResults
	numDevices, _ := NumMonitorDevices()
	deviceList := make([]int, numDevices)
	for i := 0; i < numDevices; i++ {
		deviceList[i] = i
	}

	var errs []string

	// ===== Level 2: Memory Bandwidth =====
	if level >= DiagLevel2 {
		if IsDiagStopped() {
			glog.V(3).Infof("runStressTests: stop requested before Memory Bandwidth, skipping remaining tests")
			resetDiagStop()
			return aggregated, nil
		}

		bwMap, err := BandwidthTestResult(deviceList)
		if err != nil {
			errs = append(errs, fmt.Sprintf("BandwidthTestResult: %v", err))
		}
		if bwMap != nil {
			for d, bw := range bwMap {
				status := DiagResultPass
				if bw <= 0 {
					status = DiagResultWarn
				}
				dr := DiagResult{
					Status:     status,
					TestName:   "Memory Bandwidth",
					TestOutput: fmt.Sprintf("Bandwidth: %.3f GB/s", bw),
					ErrorCode:  0,
				}
				merged := false
				for i := range aggregated.PerHCU {
					if aggregated.PerHCU[i].HCU == d {
						aggregated.PerHCU[i].DiagResults = append(aggregated.PerHCU[i].DiagResults, dr)
						merged = true
						break
					}
				}
				if !merged {
					aggregated.PerHCU = append(aggregated.PerHCU, HCUResult{
						HCU:         d,
						RC:          0,
						DiagResults: []DiagResult{dr},
					})
				}
			}
			aggregated.Software = append(aggregated.Software, DiagResult{
				Status:     DiagResultPass,
				TestName:   "Memory Bandwidth Summary",
				TestOutput: fmt.Sprintf("Parsed %d entries", len(bwMap)),
				ErrorCode:  0,
			})
		}
	}

	// ===== Level 3: PCIe / XHCL / TargetStress =====
	if level >= DiagLevel3 {
		// PCIe
		if IsDiagStopped() {
			glog.V(3).Infof("runStressTests: stop requested before PCIe, skipping remaining tests")
			resetDiagStop()
			return aggregated, nil
		}
		if pcieRes, err := PcieBandwidthTestResult(); err != nil {
			errs = append(errs, fmt.Sprintf("PcieBandwidthTestResult: %v", err))
		} else {
			for _, h := range pcieRes.HCUs {
				hcuIdx := h.DvInd
				sysToFb := h.SysToFb
				fbToSys := h.FbToSys
				status := DiagResultPass
				if sysToFb <= 0 || fbToSys <= 0 {
					status = DiagResultWarn
				}
				dr := DiagResult{
					Status:     status,
					TestName:   "PCIe Bandwidth",
					TestOutput: fmt.Sprintf("Sys->Fb: %.3f GB/s, Fb->Sys: %.3f GB/s", sysToFb, fbToSys),
					ErrorCode:  0,
				}
				merged := false
				for i := range aggregated.PerHCU {
					if aggregated.PerHCU[i].HCU == hcuIdx {
						aggregated.PerHCU[i].DiagResults = append(aggregated.PerHCU[i].DiagResults, dr)
						merged = true
						break
					}
				}
				if !merged {
					aggregated.PerHCU = append(aggregated.PerHCU, HCUResult{
						HCU:         hcuIdx,
						RC:          0,
						DiagResults: []DiagResult{dr},
					})
				}
			}
			aggregated.Software = append(aggregated.Software, DiagResult{
				Status:     DiagResultPass,
				TestName:   "PCIe Bandwidth Summary",
				TestOutput: fmt.Sprintf("Parsed devices: %d, log=%s", pcieRes.DeviceCount, pcieRes.LogFile),
				ErrorCode:  0,
			})
		}

		// XHCL
		if IsDiagStopped() {
			glog.V(3).Infof("runStressTests: stop requested before XHCL, skipping remaining tests")
			resetDiagStop()
			return aggregated, nil
		}
		if xhclPairs, err := XHCLTestResult(); err != nil {
			errs = append(errs, fmt.Sprintf("XHCLTestResult: %v", err))
		} else {
			for _, pair := range xhclPairs {
				srcIdx := pair.SrcHCUId
				dstIdx := pair.DstHCUId
				bw := pair.BandwidthGBs

				srcDiag := DiagResult{
					Status:     DiagResultPass,
					TestName:   fmt.Sprintf("XHCL to HCU%d", dstIdx),
					TestOutput: fmt.Sprintf("HCU%d <-> HCU%d XHCL: %.3f GB/s", srcIdx, dstIdx, bw),
					ErrorCode:  0,
				}
				dstDiag := DiagResult{
					Status:     DiagResultPass,
					TestName:   fmt.Sprintf("XHCL to HCU%d", srcIdx),
					TestOutput: fmt.Sprintf("HCU%d <-> HCU%d XHCL: %.3f GB/s", dstIdx, srcIdx, bw),
					ErrorCode:  0,
				}

				// 合并到 srcIdx
				mergedSrc := false
				for i := range aggregated.PerHCU {
					if aggregated.PerHCU[i].HCU == srcIdx {
						aggregated.PerHCU[i].DiagResults = append(aggregated.PerHCU[i].DiagResults, srcDiag)
						mergedSrc = true
						break
					}
				}
				if !mergedSrc {
					aggregated.PerHCU = append(aggregated.PerHCU, HCUResult{
						HCU:         srcIdx,
						RC:          0,
						DiagResults: []DiagResult{srcDiag},
					})
				}

				// 合并到 dstIdx
				mergedDst := false
				for i := range aggregated.PerHCU {
					if aggregated.PerHCU[i].HCU == dstIdx {
						aggregated.PerHCU[i].DiagResults = append(aggregated.PerHCU[i].DiagResults, dstDiag)
						mergedDst = true
						break
					}
				}
				if !mergedDst {
					aggregated.PerHCU = append(aggregated.PerHCU, HCUResult{
						HCU:         dstIdx,
						RC:          0,
						DiagResults: []DiagResult{dstDiag},
					})
				}
			}
			aggregated.Software = append(aggregated.Software, DiagResult{
				Status:     DiagResultPass,
				TestName:   "XHCL Summary",
				TestOutput: fmt.Sprintf("Parsed %d XHCL pairs", len(xhclPairs)),
				ErrorCode:  0,
			})
		}

		// TargetStress
		if IsDiagStopped() {
			glog.V(3).Infof("runStressTests: stop requested before TargetStress, skipping remaining tests")
			resetDiagStop()
			return aggregated, nil
		}
		if tsRes, err := TargetStressTestResult(); err != nil {
			errs = append(errs, fmt.Sprintf("TargetStressTestResult: %v", err))
		} else {
			for _, r := range tsRes.Results {
				status := DiagResultPass
				if r.Failed || r.Mean <= 0 {
					status = DiagResultWarn
				}
				dr := DiagResult{
					Status:     status,
					TestName:   fmt.Sprintf("TargetStress GEMM %s", r.GemmName),
					TestOutput: fmt.Sprintf("GEMM=%s, Mean=%.3f", r.GemmName, r.Mean),
					ErrorCode:  0,
				}
				merged := false
				for i := range aggregated.PerHCU {
					if aggregated.PerHCU[i].HCU == r.HCUId {
						aggregated.PerHCU[i].DiagResults = append(aggregated.PerHCU[i].DiagResults, dr)
						merged = true
						break
					}
				}
				if !merged {
					aggregated.PerHCU = append(aggregated.PerHCU, HCUResult{
						HCU:         r.HCUId,
						RC:          0,
						DiagResults: []DiagResult{dr},
					})
				}
			}
			aggregated.Software = append(aggregated.Software, DiagResult{
				Status:     DiagResultPass,
				TestName:   "TargetStress Summary",
				TestOutput: fmt.Sprintf("Parsed %d gemm entries", len(tsRes.Results)),
				ErrorCode:  0,
			})
		}
	}

	// ===== Level 4: MemtestCL / EDPp =====
	if level >= DiagLevel4 {
		if IsDiagStopped() {
			glog.V(3).Infof("runStressTests: stop requested before MemtestCL, skipping remaining tests")
			resetDiagStop()
			return aggregated, nil
		}

		if memRes, err := MemtestCLTestResult(deviceList); err != nil {
			errs = append(errs, fmt.Sprintf("MemtestCLTestResult: %v", err))
			for _, r := range memRes.Results {
				status := DiagResultPass
				if !r.Passed {
					status = DiagResultWarn
				}
				parts := make([]string, 0, len(r.Summary))
				for k, v := range r.Summary {
					parts = append(parts, fmt.Sprintf("%s: %s", k, v))
				}
				aggregated.PerHCU = append(aggregated.PerHCU, HCUResult{
					HCU:         r.HCUId,
					RC:          0,
					DiagResults: []DiagResult{{Status: status, TestName: "MemtestCL", TestOutput: strings.Join(parts, " ; "), ErrorCode: 0}},
				})
				aggregated.Software = append(aggregated.Software, DiagResult{
					Status:     DiagResultPass,
					TestName:   fmt.Sprintf("MemtestCL HCU%d", r.HCUId),
					TestOutput: fmt.Sprintf("log=%s", r.LogFile),
					ErrorCode:  0,
				})
			}
		} else {
			for _, r := range memRes.Results {
				status := DiagResultPass
				if !r.Passed {
					status = DiagResultWarn
				}
				parts := make([]string, 0, len(r.Summary))
				for k, v := range r.Summary {
					parts = append(parts, fmt.Sprintf("%s: %s", k, v))
				}
				aggregated.PerHCU = append(aggregated.PerHCU, HCUResult{
					HCU:         r.HCUId,
					RC:          0,
					DiagResults: []DiagResult{{Status: status, TestName: "MemtestCL", TestOutput: strings.Join(parts, " ; "), ErrorCode: 0}},
				})
				aggregated.Software = append(aggregated.Software, DiagResult{
					Status:     DiagResultPass,
					TestName:   fmt.Sprintf("MemtestCL HCU%d", r.HCUId),
					TestOutput: fmt.Sprintf("log=%s", r.LogFile),
					ErrorCode:  0,
				})
			}
		}

		if IsDiagStopped() {
			glog.V(3).Infof("runStressTests: stop requested after MemtestCL, skipping EDPp")
			resetDiagStop()
			return aggregated, nil
		}

		if edppRes, err := EDPpTestResult(); err != nil {
			errs = append(errs, fmt.Sprintf("EDPpTestResult: %v", err))
		} else {
			for _, d := range edppRes.HCUEdppResults {
				for _, p := range d.PatternResults {
					status := DiagResultPass
					if p.ECCCount > 0 || p.MemoryErrorCount > 0 || p.ComputeErrorCount > 0 {
						status = DiagResultWarn
					}
					aggregated.PerHCU = append(aggregated.PerHCU, HCUResult{
						HCU: d.HCUId,
						RC:  0,
						DiagResults: []DiagResult{
							{Status: status, TestName: fmt.Sprintf("EDPp Pattern %s", p.PatternName), TestOutput: fmt.Sprintf("Pattern=%s, ECC=%d, Mem=%d, Compute=%d", p.PatternName, p.ECCCount, p.MemoryErrorCount, p.ComputeErrorCount), ErrorCode: 0},
						},
					})
				}
			}
			aggregated.Software = append(aggregated.Software, DiagResult{
				Status:     DiagResultPass,
				TestName:   "EDPp Summary",
				TestOutput: fmt.Sprintf("Parsed %d HCU EDPp results, logdir=%s", len(edppRes.HCUEdppResults), edppRes.LogDir),
				ErrorCode:  0,
			})
		}
	}

	if len(errs) > 0 {
		return aggregated, fmt.Errorf("runStressTests encountered errors: %s", strings.Join(errs, " ; "))
	}
	return aggregated, nil
}

// mergeStressIntoDiag 将 runStressTests 返回的 DiagResults 合并到主 diagResults 中（按 HCU 索引合并）。
func mergeStressIntoDiag(main *DiagResults, stress DiagResults) {
	// 合并软件层面的摘要（直接追加）
	if len(stress.Software) > 0 {
		main.Software = append(main.Software, stress.Software...)
	}
	// 合并每个 HCU 的诊断结果：按 HCU 索引合并条目
	for _, stressHCU := range stress.PerHCU {
		merged := false
		for i := range main.PerHCU {
			if main.PerHCU[i].HCU == stressHCU.HCU {
				main.PerHCU[i].DiagResults = append(main.PerHCU[i].DiagResults, stressHCU.DiagResults...)
				merged = true
				break
			}
		}
		if !merged {
			// 如果主结果中不存在该 HCU，则直接追加整个 HCUResult
			main.PerHCU = append(main.PerHCU, stressHCU)
		}
	}
}

// stopDiagFlag 使用 int32 作原子布尔（0 = false, 1 = true）
var stopDiagFlag int32 = 0

// IsDiagStopped 返回是否已经请求停止
func IsDiagStopped() bool {
	return atomic.LoadInt32(&stopDiagFlag) == 1
}

// resetDiagStop 非导出函数：原子清除停止标志（0）
func resetDiagStop() {
	atomic.StoreInt32(&stopDiagFlag, 0)
}
