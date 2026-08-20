// Copyright (c) 2026 Hygon Information Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package dcgm

import (
	"context"
	"strings"
	"testing"
)

func TestValidatePcieResult(t *testing.T) {
	valid := PcieResult{
		DeviceCount: 3,
		HCUs: []HCUBW{
			{DvInd: 0, SysToFb: 1, FbToSys: 1},
			{DvInd: 1, SysToFb: 2, FbToSys: 2},
		},
	}
	if err := validatePcieResult(valid); err != nil {
		t.Fatalf("valid PCIe result rejected: %v", err)
	}

	invalid := valid
	invalid.HCUs[1].FbToSys = 0
	if err := validatePcieResult(invalid); err == nil {
		t.Fatal("zero PCIe bandwidth was accepted")
	}
}

func TestRunBandwidthTestReturnsDeviceFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	results, err := runBandwidthTest(ctx, []int{3}, false)
	if err == nil || !strings.Contains(err.Error(), "HCU3") || !strings.Contains(err.Error(), context.Canceled.Error()) {
		t.Fatalf("expected HCU-specific cancellation error, got %v", err)
	}
	if got := results[3]; got != 0 {
		t.Fatalf("expected zero bandwidth for canceled HCU, got %v", got)
	}
}
