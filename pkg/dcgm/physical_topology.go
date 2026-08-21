/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var physicalTopologySources = []string{
	"/sys/class/hsd/hsd/link_topo",
	"/etc/hfm/topo.cfg",
}

var kfdTopologyRoots = []string{
	"/sys/class/kfd/kfd/topology/nodes",
	"/sys/devices/virtual/kfd/kfd/topology/nodes",
}

type PhysicalTopologyNode struct {
	NodeID      int               `json:"nodeId"`
	DvInd       *int              `json:"dvInd,omitempty"`
	GPUId       string            `json:"gpuId,omitempty"`
	LocationID  string            `json:"locationId,omitempty"`
	Domain      string            `json:"domain,omitempty"`
	RenderMinor string            `json:"drmRenderMinor,omitempty"`
	HiveID      string            `json:"hiveId,omitempty"`
	UniqueID    string            `json:"uniqueId,omitempty"`
	Properties  map[string]string `json:"properties"`
}

type PhysicalTopologyLink struct {
	Kind         string            `json:"kind"`
	Index        int               `json:"index"`
	NodeFrom     int               `json:"nodeFrom"`
	NodeTo       int               `json:"nodeTo"`
	Type         string            `json:"type,omitempty"`
	Weight       string            `json:"weight,omitempty"`
	MinBandwidth string            `json:"minBandwidth,omitempty"`
	MaxBandwidth string            `json:"maxBandwidth,omitempty"`
	Flags        string            `json:"flags,omitempty"`
	SwitchState  string            `json:"switchState,omitempty"`
	Properties   map[string]string `json:"properties"`
}

type PhysicalLinkTopology struct {
	Source    string                 `json:"source"`
	RawConfig string                 `json:"rawConfig,omitempty"`
	Nodes     []PhysicalTopologyNode `json:"nodes"`
	Links     []PhysicalTopologyLink `json:"links"`
}

// GetPhysicalLinkTopology 聚合 HSD 配置与 KFD 节点/链路信息。
// link_topo/topo.cfg 的厂商格式未知时保留原文，KFD properties 则按键值结构化。
func GetPhysicalLinkTopology() (PhysicalLinkTopology, error) {
	result := PhysicalLinkTopology{
		Nodes: make([]PhysicalTopologyNode, 0),
		Links: make([]PhysicalTopologyLink, 0),
	}

	for _, source := range physicalTopologySources {
		content, err := readSysfsFile(source)
		if err == nil {
			result.Source = source
			result.RawConfig = content
			break
		}
	}

	root, err := firstExistingDirectory(kfdTopologyRoots)
	if err == nil {
		nodes, links, readErr := readKFDPhysicalTopology(root)
		if readErr != nil {
			return PhysicalLinkTopology{}, readErr
		}
		result.Nodes = nodes
		result.Links = links
		if result.Source == "" {
			result.Source = root
		}
	}

	if result.Source == "" {
		return PhysicalLinkTopology{}, fmt.Errorf("physical topology sources are unavailable")
	}
	return result, nil
}

func firstExistingDirectory(candidates []string) (string, error) {
	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("none of the topology directories exist")
}

func readKFDPhysicalTopology(root string) ([]PhysicalTopologyNode, []PhysicalTopologyLink, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read KFD topology root %s: %w", root, err)
	}

	nodeIDs := make([]int, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		nodeID, err := strconv.Atoi(entry.Name())
		if err == nil {
			nodeIDs = append(nodeIDs, nodeID)
		}
	}
	sort.Ints(nodeIDs)

	nodes := make([]PhysicalTopologyNode, 0, len(nodeIDs))
	links := make([]PhysicalTopologyLink, 0)
	nextDvInd := 0
	for _, nodeID := range nodeIDs {
		nodeDir := filepath.Join(root, strconv.Itoa(nodeID))
		properties, err := readKeyValueProperties(filepath.Join(nodeDir, "properties"))
		if err != nil {
			continue
		}
		node := PhysicalTopologyNode{
			NodeID:      nodeID,
			GPUId:       strings.TrimSpace(readOptionalFile(filepath.Join(nodeDir, "gpu_id"))),
			LocationID:  properties["location_id"],
			Domain:      properties["domain"],
			RenderMinor: properties["drm_render_minor"],
			HiveID:      properties["hive_id"],
			UniqueID:    properties["unique_id"],
			Properties:  properties,
		}
		if node.GPUId != "" && node.GPUId != "0" {
			dvInd := nextDvInd
			node.DvInd = &dvInd
			nextDvInd++
		}
		nodes = append(nodes, node)

		for _, kind := range []string{"io_links", "p2p_links"} {
			nodeLinks, err := readKFDLinks(filepath.Join(nodeDir, kind), kind)
			if err == nil {
				links = append(links, nodeLinks...)
			}
		}
	}
	return nodes, links, nil
}

func readKFDLinks(dir, kind string) ([]PhysicalTopologyLink, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	links := make([]PhysicalTopologyLink, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		index, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		properties, err := readKeyValueProperties(filepath.Join(dir, entry.Name(), "properties"))
		if err != nil {
			continue
		}
		links = append(links, PhysicalTopologyLink{
			Kind:         kind,
			Index:        index,
			NodeFrom:     parsePropertyInt(properties["node_from"]),
			NodeTo:       parsePropertyInt(properties["node_to"]),
			Type:         properties["type"],
			Weight:       properties["weight"],
			MinBandwidth: properties["min_bandwidth"],
			MaxBandwidth: properties["max_bandwidth"],
			Flags:        properties["flags"],
			SwitchState:  properties["switch_state"],
			Properties:   properties,
		})
	}
	sort.Slice(links, func(i, j int) bool { return links[i].Index < links[j].Index })
	return links, nil
}

func readKeyValueProperties(path string) (map[string]string, error) {
	content, err := readSysfsFile(path)
	if err != nil {
		return nil, err
	}
	properties := make(map[string]string)
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			properties[fields[0]] = strings.Join(fields[1:], " ")
		}
	}
	return properties, nil
}

func readOptionalFile(path string) string {
	value, _ := readSysfsFile(path)
	return value
}

func parsePropertyInt(value string) int {
	parsed, _ := strconv.Atoi(value)
	return parsed
}
