package dcgm

import "testing"

func TestCopyProcessGPUIndices(t *testing.T) {
	got, err := copyProcessGPUIndices([]int32{5, 1, 4, 20, 6, 3}, 6)
	if err != nil {
		t.Fatalf("copyProcessGPUIndices returned error: %v", err)
	}

	want := []int{5, 1, 4, 20, 6, 3}
	if len(got) != len(want) {
		t.Fatalf("got %d indices, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestCopyProcessGPUIndicesRejectsInvalidCount(t *testing.T) {
	for _, count := range []int{-1, 3} {
		if _, err := copyProcessGPUIndices([]int32{0, 1}, count); err == nil {
			t.Errorf("copyProcessGPUIndices(%d) returned nil error", count)
		}
	}
}
