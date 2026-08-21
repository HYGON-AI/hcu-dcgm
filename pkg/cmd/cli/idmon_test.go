// Copyright (c) 2026 Hygon Information Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package cli

import "testing"

func TestParseIDMonDevices(t *testing.T) {
	devices, err := parseIDMonDevices("0,2,2")
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) != 2 || devices[0] != 0 || devices[1] != 2 {
		t.Fatalf("devices = %v", devices)
	}

	if _, err := parseIDMonDevices("-1"); err == nil {
		t.Fatal("negative device index must fail")
	}
}

func TestSelectIDMonMetrics(t *testing.T) {
	metrics, err := selectIDMonMetrics("upc")
	if err != nil {
		t.Fatal(err)
	}
	if len(metrics) != 3 || metrics[0].Key != 'c' || metrics[1].Key != 'p' || metrics[2].Key != 'u' {
		t.Fatalf("metrics = %+v", metrics)
	}

	if _, err := selectIDMonMetrics("x"); err == nil {
		t.Fatal("unknown metric must fail")
	}
}
