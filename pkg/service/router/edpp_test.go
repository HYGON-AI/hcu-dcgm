// Copyright (c) 2026 Hygon Information Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package router

import (
	"strings"
	"testing"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

func TestConvertEdppResult(t *testing.T) {
	input := dcgm.EDPPResult{
		LogDir: "logs/edpp",
		HCUEdppResults: []dcgm.HCUEdppResult{
			{
				HCUId:   0,
				Backend: "legacy",
				PatternResults: []dcgm.PatternResult{
					{PatternName: "10KHz", ECCCount: 1},
				},
			},
			{
				HCUId:   1,
				Backend: "nmz-gfx938",
				GemmPowerResults: []dcgm.GemmPowerEdppResult{
					{
						HCU:        1,
						Backend:    "nmz-gfx938",
						AvgPowerW:  800,
						PeakPowerW: 900,
						AvgGFXMHz:  1800,
						Passed:     true,
						Output:     "[CPU CHECK] compute time: 589.815 ms\niterator: 0\n",
					},
				},
			},
		},
	}

	result := convertEdppResult(input)
	if result.DeviceNumber != 2 {
		t.Fatalf("DeviceNumber = %d, want 2", result.DeviceNumber)
	}
	if len(result.PerHCU) != 2 {
		t.Fatalf("len(PerHCU) = %d, want 2", len(result.PerHCU))
	}
	if got := result.PerHCU[0].DiagResults[0].Status; got != DiagResultWarn {
		t.Fatalf("legacy status = %q, want %q", got, DiagResultWarn)
	}
	nmz := result.PerHCU[1].DiagResults[0]
	if nmz.Status != DiagResultPass {
		t.Fatalf("nmz status = %q, want %q", nmz.Status, DiagResultPass)
	}
	if !strings.Contains(nmz.TestOutput, "Avg Power: 800 W") {
		t.Fatalf("nmz output = %q, want power summary", nmz.TestOutput)
	}
	if !strings.Contains(nmz.TestOutput, "[CPU CHECK] compute time") {
		t.Fatalf("nmz output = %q, want raw output summary", nmz.TestOutput)
	}
	if strings.Contains(nmz.TestOutput, "TFLOPS") {
		t.Fatalf("nmz output = %q, should not contain TFLOPS", nmz.TestOutput)
	}
	if strings.Contains(nmz.TestOutput, "\n") {
		t.Fatalf("nmz output = %q, want single-line summary", nmz.TestOutput)
	}
	if len(result.Software) != 1 || !strings.Contains(result.Software[0].TestOutput, "2 HCU") {
		t.Fatalf("summary = %+v, want 2 HCU", result.Software)
	}
}

func TestConvertEdppResultMarksFailedGemmPower(t *testing.T) {
	input := dcgm.EDPPResult{
		HCUEdppResults: []dcgm.HCUEdppResult{
			{
				HCUId:   0,
				Backend: "bmz",
				GemmPowerResults: []dcgm.GemmPowerEdppResult{
					{HCU: 0, Backend: "bmz", Passed: false, Error: "CPU CHECK failed"},
				},
			},
		},
	}

	result := convertEdppResult(input)
	diag := result.PerHCU[0].DiagResults[0]
	if diag.Status != DiagResultFail || diag.ErrorCode != -1 {
		t.Fatalf("failed result = %+v, want fail/-1", diag)
	}
	if diag.ErrorMessage != "CPU CHECK failed" {
		t.Fatalf("ErrorMessage = %q", diag.ErrorMessage)
	}
}

func TestConvertEdppResultProjectsBMZSuccess(t *testing.T) {
	input := dcgm.EDPPResult{
		HCUEdppResults: []dcgm.HCUEdppResult{
			{
				HCUId:   0,
				Backend: "bmz",
				GemmPowerResults: []dcgm.GemmPowerEdppResult{
					{HCU: 0, Backend: "bmz", AvgPowerW: 876, PeakPowerW: 1024, AvgGFXMHz: 1800, Passed: true},
				},
			},
		},
	}

	result := convertEdppResult(input)
	diag := result.PerHCU[0].DiagResults[0]
	if diag.Status != DiagResultPass {
		t.Fatalf("bmz status = %q, want %q", diag.Status, DiagResultPass)
	}
	if !strings.Contains(diag.TestName, "bmz") || !strings.Contains(diag.TestOutput, "Avg Power: 876 W") {
		t.Fatalf("bmz result = %+v, want backend and power summary", diag)
	}
	if strings.Contains(diag.TestOutput, "TFLOPS") {
		t.Fatalf("bmz output = %q, should not contain TFLOPS", diag.TestOutput)
	}
}
