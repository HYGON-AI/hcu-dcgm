/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/golang/glog"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Info("go-dcgm start ...")
	dcgm.Init()
	defer dcgm.ShutDown()
	//dcgm.RunDiag()
	/*dvlistFlag := flag.String("dvlist", "0", "comma-separated device ids to test (e.g. 0,1,2)")
		skipPCIe := flag.Bool("skip-pcie", false, "skip PCIe API test")
		skipXHCL := flag.Bool("skip-xhcl", false, "skip XHCL API test")
		skipEDPp := flag.Bool("skip-edpp", false, "skip EDPp API test")
		skipTarget := flag.Bool("skip-target", false, "skip TargetStress API test")
		skipMem := flag.Bool("skip-mem", false, "skip MemtestCL API test")
		skipBw := flag.Bool("skip-bw", false, "skip BandwidthTest API test")
		timeout := flag.Duration("timeout", 0, "optional global timeout (e.g. 30s); 0 means no timeout (not enforced in this demo)")
		flag.Parse()

		dvlist, err := parseDvList(*dvlistFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid -dvlist: %v\n", err)
			os.Exit(2)
		}

		fmt.Printf("Demo start @ %s\n", time.Now().Format(time.RFC3339))
		fmt.Printf("Devices: %v\n\n", dvlist)

		// 1) PCIe (PcieBandwidthTestAPI)
		if !*skipPCIe {
			fmt.Println("==> Calling PcieBandwidthTestAPI() ...")
			res, err := dcgm.PcieBandwidthTestResult()
			if err != nil {
				fmt.Printf("PcieBandwidthTestAPI error: %v\n\n", err)
			} else {
				fmt.Println("PcieBandwidthTestAPI result:")
				prettyPrint(res)
				fmt.Println()
			}
		} else {
			fmt.Println("==> Skipping PcieBandwidthTestAPI")
		}

		// 2) XHCL (XHCLTestAPI)
		if !*skipXHCL {
			fmt.Println("==> Calling XHCLTestAPI() ...")
			res, err := dcgm.XHCLTestResult()
			if err != nil {
				fmt.Printf("XHCLTestAPI error: %v\n\n", err)
			} else {
				fmt.Println("XHCLTestAPI result:")
				prettyPrint(res)
				fmt.Println()
			}
		} else {
			fmt.Println("==> Skipping XHCLTestAPI")
		}

		// 3) EDPp (EDPpTestAPI)
		if !*skipEDPp {
			fmt.Println("==> Calling EDPpTestAPI() ...")
			res, err := dcgm.EDPpTestResult()
			if err != nil {
				fmt.Printf("EDPpTestAPI error: %v\n\n", err)
			} else {
				fmt.Println("EDPpTestAPI result:")
				prettyPrint(res)
				fmt.Println()
			}
		} else {
			fmt.Println("==> Skipping EDPpTestAPI")
		}

		// 4) TargetStress (TargetStressTestAPI)
		if !*skipTarget {
			fmt.Println("==> Calling TargetStressTestAPI() ...")
			res, err := dcgm.TargetStressTestResult()
			if err != nil {
				fmt.Printf("TargetStressTestAPI error: %v\n\n", err)
			} else {
				fmt.Println("TargetStressTestAPI result:")
				prettyPrint(res)
				fmt.Println()
			}
		} else {
			fmt.Println("==> Skipping TargetStressTestAPI")
		}

		// 5) MemtestCL (MemtestCLAPI requires dv list)
		if !*skipMem {
			fmt.Println("==> Calling MemtestCLAPI(dvlist) ...")
			res, err := dcgm.MemtestCLTestResult(dvlist)
			if err != nil {
				fmt.Printf("MemtestCLAPI error: %v\n\n", err)
				// still print partial result if any
				prettyPrint(res)
				fmt.Println()
			} else {
				fmt.Println("MemtestCLAPI result:")
				prettyPrint(res)
				fmt.Println()
			}
		} else {
			fmt.Println("==> Skipping MemtestCLAPI")
		}

		// 6) BandwidthTestAPI (example: map[int]float64) requires dvlist
		if !*skipBw {
			fmt.Println("==> Calling BandwidthTestAPI(dvlist) ...")
			bwRes, err := dcgm.BandwidthTestResult(dvlist)
			if err != nil {
				fmt.Printf("BandwidthTestAPI error: %v\n\n", err)
			} else {
				fmt.Println("BandwidthTestAPI result (map[int]float64):")
				prettyPrint(bwRes)
				fmt.Println()
			}
		} else {
			fmt.Println("==> Skipping BandwidthTestAPI")
		}

		fmt.Printf("Demo finished @ %s\n", time.Now().Format(time.RFC3339))
		if *timeout != 0 {
			fmt.Printf("Note: demo run used timeout flag: %v (not enforced here)\n", *timeout)
		}
	}

	func parseDvList(s string) ([]int, error) {
		if s == "" {
			return []int{0}, nil
		}
		parts := strings.Split(s, ",")
		out := make([]int, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			n, err := strconv.Atoi(p)
			if err != nil {
				return nil, fmt.Errorf("invalid dv id %q: %w", p, err)
			}
			out = append(out, n)
		}
		if len(out) == 0 {
			out = []int{0}
		}
		return out, nil*/

	//fmt.Println("Starting synchronous RunDiag(level=3) ... (this will block until finished)")
	//
	//start := time.Now()
	//// 直接同步调用（会阻塞当前 goroutine）
	//res, err := dcgm.RunDiag(3)
	//elapsed := time.Since(start)
	//
	//fmt.Printf("RunDiag finished in %v, err=%v\n", elapsed, err)
	//
	//// 打印结构化结果（JSON）
	//b, jerr := json.MarshalIndent(res, "", "  ")
	//if jerr != nil {
	//	log.Fatalf("json marshal failed: %v", jerr)
	//}
	//fmt.Println("DiagResults:\n", string(b))
	//fmt.Println("Starting async RunDiag(level=3) in background...")
	//
	//// 结果通道
	//type result struct {
	//	res dcgm.DiagResults
	//	err error
	//}
	//resultCh := make(chan result, 1)
	//
	//// 在后台 goroutine 启动 RunDiag（不阻塞主 goroutine）
	//go func() {
	//	res, err := dcgm.RunDiag(3)
	//	resultCh <- result{res: res, err: err}
	//}()
	//
	//// 等待一段时间再发 Stop 请求（模拟用户在压测中断）
	//sleepBeforeStop := 10 * time.Second
	//fmt.Printf("Will call StopDiag() after %v ...\n", sleepBeforeStop)
	//time.Sleep(sleepBeforeStop)
	//
	//fmt.Println("Calling StopDiag() ...")
	//dcgm.StopDiag()
	//
	//// 等待 RunDiag 返回，但我们希望不要无限等 —— 设一个最大等待超时
	//waitTimeout := 2 * time.Minute
	//fmt.Printf("Waiting up to %v for RunDiag to finish...\n", waitTimeout)
	//
	//select {
	//case r := <-resultCh:
	//	fmt.Println("RunDiag returned.")
	//	if r.err != nil {
	//		fmt.Printf("error: %v\n", r.err)
	//	}
	//	b, _ := json.MarshalIndent(r.res, "", "  ")
	//	fmt.Println("DiagResults:\n", string(b))
	//case <-time.After(waitTimeout):
	//	// 超时后你可以选择继续后台等待或报警
	//	fmt.Println("Timed out waiting for RunDiag to finish.")
	//	// 如果你想强制结束进程，可在这里退出或执行其它补救
	//}
	//
	//fmt.Println("demo done")

	// 参数：diagnose level、等待多少秒再发 stop、等待 runDiag 返回的超时时间
	//level := flag.Int("level", 3, "diag level to run (1-4)")
	//waitBeforeStop := flag.Int("wait", 10, "seconds to wait before calling StopDiag()")
	//waitForResult := flag.Int("timeout", 300, "seconds to wait for RunDiag to return before giving up")
	//flag.Parse()
	//
	//fmt.Printf("Demo: start RunDiag(level=%d) in background\n", *level)
	//
	//type result struct {
	//	res dcgm.DiagResults
	//	err error
	//}
	//resultCh := make(chan result, 1)
	//
	//startTime := time.Now()
	//
	//// 启动后台 goroutine 运行诊断（不阻塞当前主 goroutine）
	//go func() {
	//	// 注意：如果你的导出函数名不是 RunDiag，请替换为实际导出名
	//	res, err := dcgm.RunDiag(*level)
	//	resultCh <- result{res: res, err: err}
	//}()
	//
	//// 等待一段时间再触发 StopDiag（模拟用户在运行中发起停止）
	//fmt.Printf("Will call StopDiag() after %d seconds...\n", *waitBeforeStop)
	//time.Sleep(time.Duration(*waitBeforeStop) * time.Second)
	//
	////fmt.Println("Calling StopDiag() now.")
	////dcgm.StopDiag()
	//
	//// 等待 RunDiag 返回（带超时）
	//timeout := time.After(time.Duration(*waitForResult) * time.Second)
	//select {
	//case r := <-resultCh:
	//	elapsed := time.Since(startTime)
	//	fmt.Printf("RunDiag returned after %v\n", elapsed)
	//	if r.err != nil {
	//		fmt.Printf("RunDiag returned error: %v\n", r.err)
	//	} else {
	//		fmt.Println("RunDiag returned nil error")
	//	}
	//
	//	// 打印判定信息：IsDiagStopped() 应该为 false（因为 runDiag 返回前有 reset）
	//	fmt.Printf("IsDiagStopped() (current flag) = %v\n", dcgm.IsDiagStopped())
	//
	//	// 输出 DiagResults 简短摘要（JSON 也同时打印）
	//	summary := struct {
	//		PerHCUCount int `json:"per_hcu_count"`
	//		SoftwareCnt int `json:"software_count"`
	//	}{
	//		PerHCUCount: len(r.res.PerHCU),
	//		SoftwareCnt: len(r.res.Software),
	//	}
	//	fmt.Printf("DiagResults summary: %+v\n", summary)
	//
	//	b, err := json.MarshalIndent(r.res, "", "  ")
	//	if err != nil {
	//		fmt.Printf("failed to marshal DiagResults to JSON: %v\n", err)
	//	} else {
	//		fmt.Println("Full DiagResults JSON:")
	//		fmt.Println(string(b))
	//	}
	//
	//case <-timeout:
	//	fmt.Printf("Timed out waiting for RunDiag to return after %d seconds.\n", *waitForResult)
	//	// 注意：runDiag 可能仍在后台运行。如果需要强制杀死进程或子进程，需在相应测试函数中支持 context cancellation。
	//}
	//
	//fmt.Println("Demo finished.")

	num, err := dcgm.NumMonitorDevices()
	if err != nil {
		log.Fatalf("NumMonitorDevices failed: %v", err)
	}

	dvIdList := make([]int, 0, num)
	for i := 0; i < num; i++ {
		dvIdList = append(dvIdList, i)
	}

	fmt.Println("Start MemtestCL on:", dvIdList)

	res, err := dcgm.MemtestCLTestResult(dvIdList)

	fmt.Println("err =", err)
	fmt.Printf("logdir = %s\n", res.LogDir)

	for _, r := range res.Results {
		fmt.Printf(
			"HCU %d passed=%v log=%s summary=%v\n",
			r.HCUId, r.Passed, r.LogFile, r.Summary,
		)
	}
}

func prettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("marshal error: %v\n", err)
		return
	}
	fmt.Println(string(b))
}
