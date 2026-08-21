/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

func TestFormatDiagFailure(t *testing.T) {
	got := formatDiagFailure("Bandwidth test", errors.New("HCU0: exit status 98"), "logs/hip-stream")

	for _, want := range []string{
		"Bandwidth test failed: HCU0: exit status 98",
		"Logs: logs/hip-stream",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output %q does not contain %q", got, want)
		}
	}
}

func TestFormatDiagFailureSkipsEmptyLogPath(t *testing.T) {
	got := formatDiagFailure("XHCL stress test", errors.New("command failed"), "")
	if strings.Contains(got, "Logs:") {
		t.Fatalf("output %q unexpectedly contains a log line", got)
	}
}

func TestFormatEDPPResultOmitsRawOutput(t *testing.T) {
	output := formatEDPPResult(dcgm.EDPPResult{
		HCUEdppResults: []dcgm.HCUEdppResult{{
			HCUId: 5,
			GemmPowerResults: []dcgm.GemmPowerEdppResult{{
				HCU:        5,
				AvgPowerW:  671,
				PeakPowerW: 800,
				AvgGFXMHz:  1265,
				Passed:     true,
				Output:     "iterator: 0\niterator: 1",
			}},
		}},
	})

	if !strings.Contains(output, "HCU 5  Avg Power: 671 W, Peak Power: 800 W, GFX Clock: 1265 MHz, PASS") {
		t.Fatalf("output %q does not contain the EDPp summary", output)
	}
	if strings.Contains(output, "iterator:") {
		t.Fatalf("output %q unexpectedly contains raw EDPp output", output)
	}
}

func TestFormatPCIESummary(t *testing.T) {
	output := formatPCIESummary(dcgm.PcieResult{
		HCUs: []dcgm.HCUBW{
			{DvInd: 1, SysToFb: 25.10, FbToSys: 24.02},
			{DvInd: 0, SysToFb: 25.34, FbToSys: 23.86},
		},
	})

	for _, want := range []string{
		"=== PCIe Bandwidth Summary ===",
		"HCU 0  System memory->HCU: 25.34 GB/s  HCU->System memory: 23.86 GB/s",
		"HCU 1  System memory->HCU: 25.10 GB/s  HCU->System memory: 24.02 GB/s",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output %q does not contain %q", output, want)
		}
	}
	if strings.Index(output, "HCU 0") > strings.Index(output, "HCU 1") {
		t.Fatalf("output %q is not sorted by HCU", output)
	}
}

func TestFormatGEMMSummary(t *testing.T) {
	output := formatGEMMSummary(dcgm.TargetStressResult{
		Results: []dcgm.GemmTestResult{
			{HCUId: 1, GemmName: "hgemm", Failed: true},
			{HCUId: 0, GemmName: "sgemm", Mean: 45.20},
			{HCUId: 1, GemmName: "dgemm", Mean: 23.90},
			{HCUId: 0, GemmName: "dgemm", Mean: 24.10},
		},
	})

	for _, want := range []string{
		"=== GEMM Performance Summary ===",
		"HCU 0  PASS 2/2  mean: sgemm=45.20, dgemm=24.10",
		"HCU 1  WARN 1/2  mean: hgemm=FAIL, dgemm=23.90",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output %q does not contain %q", output, want)
		}
	}
	if strings.Index(output, "HCU 0") > strings.Index(output, "HCU 1") {
		t.Fatalf("output %q is not sorted by HCU", output)
	}
}

func TestDiagnosticStatusLabel(t *testing.T) {
	if got := diagnosticStatusLabel(dcgm.DiagResultWarn); got != "Warning" {
		t.Fatalf("warning label = %q, want Warning", got)
	}
	if got := diagnosticStatusLabel(dcgm.DiagResultPass); got != dcgm.DiagResultPass {
		t.Fatalf("pass label = %q, want %q", got, dcgm.DiagResultPass)
	}
}
