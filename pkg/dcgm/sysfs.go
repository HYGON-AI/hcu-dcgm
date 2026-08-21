/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// readSysfsFile 读取 sysfs 文件内容（去除首尾空白符）
func readSysfsFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// writeSysfsFile 写入 sysfs 文件（需 root 权限）
func writeSysfsFile(path, value string) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer f.Close()

	if _, err := f.WriteString(value); err != nil {
		return fmt.Errorf("failed to write to %s: %w", path, err)
	}
	return nil
}

// listSysfsDir 列出目录下的子目录或文件
func listSysfsDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list %s: %w", path, err)
	}

	var results []string
	for _, entry := range entries {
		results = append(results, filepath.Join(path, entry.Name()))
	}
	return results, nil
}

// dvIndexToCardIndex 将 dvInd 映射到 DRM cardX 编号。
//
// 通过读取 KFD topology 节点的 drm_render_minor，再在 /sys/class/drm/
// 中找到包含对应 renderD{minor} 的 cardX 目录。
// 若查找失败回退到 dvInd+1（card0 通常是 boot_vga，HCU 从 card1 起）。
func dvIndexToCardIndex(dvInd int) int {
	// KFD topology：node 0 = CPU，GPU 从 node 1 起，因此 dvInd 对应 node dvInd+1
	propPath := fmt.Sprintf("/sys/class/kfd/kfd/topology/nodes/%d/properties", dvInd+1)
	content, err := readSysfsFile(propPath)
	if err != nil {
		glog.V(3).Infof("dvIndexToCardIndex: cannot read %s: %v, fallback=%d", propPath, err, dvInd)
		return dvInd
	}

	renderMinor := -1
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "drm_render_minor ") {
			if n, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "drm_render_minor "))); err == nil {
				renderMinor = n
			}
			break
		}
	}
	if renderMinor < 0 {
		glog.V(3).Infof("dvIndexToCardIndex: drm_render_minor not found in %s, fallback=%d", propPath, dvInd)
		return dvInd
	}

	// 在 /sys/class/drm/cardX/device/drm/ 下找包含 renderD{minor} 的 cardX
	renderName := fmt.Sprintf("renderD%d", renderMinor)
	entries, err := os.ReadDir("/sys/class/drm")
	if err != nil {
		glog.V(3).Infof("dvIndexToCardIndex: cannot read /sys/class/drm: %v, fallback=%d", err, dvInd)
		return dvInd
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, "card") || strings.Contains(name, "-") {
			continue
		}
		checkPath := filepath.Join("/sys/class/drm", name, "device", "drm", renderName)
		if _, err := os.Stat(checkPath); err == nil {
			if cardIdx, err := strconv.Atoi(strings.TrimPrefix(name, "card")); err == nil {
				glog.V(5).Infof("dvIndexToCardIndex: dvInd=%d -> renderD%d -> %s -> cardIdx=%d", dvInd, renderMinor, name, cardIdx)
				return cardIdx
			}
		}
	}

	glog.V(3).Infof("dvIndexToCardIndex: renderD%d not matched, fallback=%d", renderMinor, dvInd)
	return dvInd
}

// writeDebugfs 写入 debugfs 文件，写入前检查路径是否可访问。
// debugfs 通常需要 root 权限且须已挂载：mount -t debugfs none /sys/kernel/debug
func writeDebugfs(path, value string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("debugfs path not found (is debugfs mounted?): %s", path)
		}
		return fmt.Errorf("cannot access debugfs path %s: %w", path, err)
	}
	return writeSysfsFile(path, value)
}
