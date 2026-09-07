/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseGemmLog_ValidMean(t *testing.T) {
	dir := t.TempDir()
	logfile := filepath.Join(dir, "hgemm.log")
	content := "fp16(h)-HCU0: min: 90.0\tmax: 99.0\tmean: 98.500\tstdev: 1.0\n"
	if err := os.WriteFile(logfile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	mean, fail, err := parseGemmLog(logfile, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fail {
		t.Fatal("expected fail=false")
	}
	if mean != 98.500 {
		t.Fatalf("mean = %v, want 98.500", mean)
	}
}

func TestParseGemmLog_MissingMean(t *testing.T) {
	dir := t.TempDir()
	logfile := filepath.Join(dir, "hgemm.log")
	// 只有 HCU0，没有 HCU1
	content := "fp16(h)-HCU0: min: 90.0\tmax: 99.0\tmean: 98.500\tstdev: 1.0\n"
	if err := os.WriteFile(logfile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, _, err := parseGemmLog(logfile, 1)
	if err == nil {
		t.Fatal("expected error for missing mean, got nil")
	}
}

func TestParseGemmLog_FailFlag(t *testing.T) {
	dir := t.TempDir()
	logfile := filepath.Join(dir, "hgemm.log")
	content := "fp16(h)-HCU0: min: 40.0\tmax: 55.0\tmean: 50.000\tstdev: 2.0\nFAIL: verification error\n"
	if err := os.WriteFile(logfile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	mean, fail, err := parseGemmLog(logfile, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fail {
		t.Fatal("expected fail=true")
	}
	if mean != 50.000 {
		t.Fatalf("mean = %v, want 50.000", mean)
	}
}
