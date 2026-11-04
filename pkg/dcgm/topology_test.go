/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestCollectXhclRemoteBdfidsHSDTopology(t *testing.T) {
	// bw150-114 的 hy-smi -R 和 HSD link_topo：每张卡有 7 条连接及 1 条空链路。
	bdfs := []uint64{0x2a00, 0x2d00, 0x3200, 0x3900, 0xaa00, 0xad00, 0xb200, 0xb900}
	remoteNodes := [][]int{
		{1, 1, 1, 3, 3, 2, 2},
		{0, 0, 0, 2, 2, 3, 3},
		{3, 3, 3, 1, 1, 0, 0},
		{2, 2, 2, 0, 0, 1, 1},
		{5, 5, 5, 7, 7, 6, 6},
		{4, 4, 4, 6, 6, 7, 7},
		{7, 7, 7, 5, 5, 4, 4},
		{6, 6, 6, 4, 4, 5, 5},
	}
	wantWeights := [8][8]int{
		{0, 3, 2, 2, 0, 0, 0, 0},
		{3, 0, 2, 2, 0, 0, 0, 0},
		{2, 2, 0, 3, 0, 0, 0, 0},
		{2, 2, 3, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 3, 2, 2},
		{0, 0, 0, 0, 3, 0, 2, 2},
		{0, 0, 0, 0, 2, 2, 0, 3},
		{0, 0, 0, 0, 2, 2, 3, 0},
	}

	for src, peers := range remoteNodes {
		t.Run(fmt.Sprintf("HCU%d", src), func(t *testing.T) {
			state := uint32(1)
			if src >= 4 {
				state = 2
			}
			states := make([]XhclLinkState, 8)
			for linkID := range states {
				states[linkID].LinkID = linkID
				if linkID < len(peers) {
					states[linkID].State = state
					states[linkID].Up = state == XhclLinkUp
				}
			}

			var queried []int
			links := collectXhclRemoteBdfids(src, states, func(dvInd, linkID int) (uint64, error) {
				if dvInd != src || linkID < 0 || linkID >= len(peers) {
					t.Fatalf("unexpected remote query: device=%d link=%d", dvInd, linkID)
				}
				queried = append(queried, linkID)
				return bdfs[peers[linkID]], nil
			})
			if !reflect.DeepEqual(queried, []int{0, 1, 2, 3, 4, 5, 6}) {
				t.Fatalf("queried links = %v, want links 0 through 6", queried)
			}
			if len(links) != len(peers) {
				t.Fatalf("got %d links, want %d", len(links), len(peers))
			}

			counts := make(map[uint64]int)
			for i, link := range links {
				if link.LinkID != i || link.BdfID != bdfs[peers[i]] {
					t.Errorf("link %d = %+v, want remote BDF %#x", i, link, bdfs[peers[i]])
				}
				counts[link.BdfID]++
			}
			for dst, bdf := range bdfs {
				if got := counts[bdf]; got != wantWeights[src][dst] {
					t.Errorf("HCU%d -> HCU%d weight = %d, want %d", src, dst, got, wantWeights[src][dst])
				}
			}
		})
	}
}

func TestCollectXhclRemoteBdfidsSkipsEmptyAndFailedQueries(t *testing.T) {
	states := []XhclLinkState{
		{LinkID: 0, State: 0},
		{LinkID: 1, State: 2},
		{LinkID: 2, State: 2},
		{LinkID: 3, State: 3},
	}
	var queried []int
	links := collectXhclRemoteBdfids(4, states, func(dvInd, linkID int) (uint64, error) {
		if dvInd != 4 {
			t.Fatalf("device = %d, want 4", dvInd)
		}
		queried = append(queried, linkID)
		if linkID == 2 {
			return 0xad00, nil
		}
		return 0, errors.New("remote BDF unavailable")
	})
	if !reflect.DeepEqual(queried, []int{1, 2, 3}) {
		t.Errorf("queried links = %v, want [1 2 3]", queried)
	}
	want := []XhclRemoteBdf{{LinkID: 2, BdfID: 0xad00}}
	if !reflect.DeepEqual(links, want) {
		t.Errorf("links = %v, want %v", links, want)
	}
}
