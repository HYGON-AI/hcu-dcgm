// Copyright (c) 2026 Hygon Information Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package dcgm

import (
	"errors"
	"strings"
	"testing"
)

func TestEdppTargetDevices(t *testing.T) {
	tests := []struct {
		name    string
		total   int
		input   []int
		want    []int
		wantErr string
	}{
		{name: "default all", total: 4, want: []int{0, 1, 2, 3}},
		{name: "single", total: 4, input: []int{2}, want: []int{2}},
		{name: "dedupe", total: 4, input: []int{1, 1, 3}, want: []int{1, 3}},
		{name: "invalid", total: 4, input: []int{4}, wantErr: "valid range is 0-3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := edppTargetDevices(tt.total, tt.input)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want contains %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("err = %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestEdppBackendSpecsUseStressMode(t *testing.T) {
	nmz, err := edppBackendSpecFor(edppBackendNMZGFX938)
	if err != nil {
		t.Fatalf("nmz spec error: %v", err)
	}
	if !nmz.stressMode {
		t.Fatalf("nmz stressMode = false, want true")
	}
	if nmz.coName != "fp8_nmz_edpp.co" || nmz.minAvgPowerW != edppNMZMinAvgPowerW {
		t.Fatalf("nmz spec = %+v, want fp8 co and nmz threshold", nmz)
	}
	if got := strings.Join(gemmPowerEdppArgs(nmz, "./fp8_nmz_edpp.co"), " "); !strings.Contains(got, "-l "+edppStressLoopCount) {
		t.Fatalf("nmz args = %q, want stress loop count", got)
	}

	bmz, err := edppBackendSpecFor(edppBackendBMZ)
	if err != nil {
		t.Fatalf("bmz spec error: %v", err)
	}
	if !bmz.stressMode {
		t.Fatalf("bmz stressMode = false, want true")
	}
	if bmz.coName != "fp16_bmz_edpp.co" || bmz.minAvgPowerW != edppBMZMinAvgPowerW {
		t.Fatalf("bmz spec = %+v, want fp16 co and bmz threshold", bmz)
	}
}

func TestFinalizeStressEdppResult(t *testing.T) {
	spec := edppBackendSpec{name: edppBackendNMZGFX938, minAvgPowerW: edppNMZMinAvgPowerW}
	result := GemmPowerEdppResult{HCU: 0, Backend: spec.name}
	err := finalizeStressEdppResult(&result, spec, edppStressSamples{
		totalPower:   1600,
		peakPower:    800,
		totalGFX:     3600,
		powerSamples: 2,
		clockSamples: 2,
	})
	if err != nil {
		t.Fatalf("finalize success error: %v", err)
	}
	if !result.Passed || result.AvgPowerW != 800 || result.PeakPowerW != 800 || result.AvgGFXMHz != 1800 {
		t.Fatalf("result = %+v, want pass with averaged samples", result)
	}

	result = GemmPowerEdppResult{HCU: 1, Backend: spec.name}
	err = finalizeStressEdppResult(&result, spec, edppStressSamples{lastPowerErr: errors.New("rsmi power failed")})
	if err == nil || !strings.Contains(result.Error, "no valid power samples") || !strings.Contains(result.Error, "rsmi power failed") {
		t.Fatalf("no sample err = %v, result = %+v", err, result)
	}

	result = GemmPowerEdppResult{HCU: 2, Backend: spec.name}
	err = finalizeStressEdppResult(&result, spec, edppStressSamples{totalPower: 740, peakPower: 370, powerSamples: 2})
	if err == nil || !strings.Contains(result.Error, "avg power 370W < threshold 400W") {
		t.Fatalf("low power err = %v, result = %+v", err, result)
	}
}
