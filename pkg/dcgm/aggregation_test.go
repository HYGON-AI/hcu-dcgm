// Copyright (c) 2026 Hygon Information Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package dcgm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHSLLogEntry(t *testing.T) {
	entry := parseHSLLogEntry("[20260805_120000.000001 42 error hfm/link.cpp:10] HCU 3 link down")
	if entry.Level != "ERROR" {
		t.Fatalf("level = %q, want ERROR", entry.Level)
	}
	if entry.Message != "HCU 3 link down" {
		t.Fatalf("message = %q", entry.Message)
	}
	if entry.Header == "" || entry.Raw == "" {
		t.Fatalf("entry must preserve header and raw line: %+v", entry)
	}
}

func TestReadKeyValueProperties(t *testing.T) {
	path := filepath.Join(t.TempDir(), "properties")
	content := "type 2\nnode_from 1\nnode_to 3\nmax_bandwidth 128000\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	properties, err := readKeyValueProperties(path)
	if err != nil {
		t.Fatal(err)
	}
	if properties["node_from"] != "1" || properties["node_to"] != "3" {
		t.Fatalf("unexpected properties: %+v", properties)
	}
	if properties["max_bandwidth"] != "128000" {
		t.Fatalf("max_bandwidth = %q", properties["max_bandwidth"])
	}
}
