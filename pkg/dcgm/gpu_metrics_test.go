package dcgm

import "testing"

func TestRSMIGPUMetricsABILayout(t *testing.T) {
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"sizeof(rsmi_gpu_metrics_t)", int(rsmiGPUMetricsStructSize), 4544},
		{"sizeof(struct gpu_xcp_metrics_t)", int(rsmiGPUXCPMetricsStructSize), 504},
	}

	for _, check := range checks {
		if check.got != check.want {
			t.Errorf("%s = %d, want %d", check.name, check.got, check.want)
		}
	}
}

func TestRSMIGPUMetricsArrayDimensions(t *testing.T) {
	var metrics RSMIGPUMetrics
	checks := []struct {
		name string
		got  int
		want int
	}{
		{"TempetureHBM", len(metrics.TempetureHBM), 4},
		{"VcnActivity", len(metrics.VcnActivity), 4},
		{"XgmiReadDataAcc", len(metrics.XgmiReadDataAcc), 8},
		{"XgmiWriteDataAcc", len(metrics.XgmiWriteDataAcc), 8},
		{"CurrentGfxclks", len(metrics.CurrentGfxclks), 8},
		{"CurrentSocclks", len(metrics.CurrentSocclks), 4},
		{"CurrentVclk0s", len(metrics.CurrentVclk0s), 4},
		{"CurrentDclk0s", len(metrics.CurrentDclk0s), 4},
		{"JpegActivity", len(metrics.JpegActivity), 32},
		{"XcpStats", len(metrics.XcpStats), 8},
		{"XgmiLinkStatus", len(metrics.XgmiLinkStatus), 8},
		{"XcpStats.GfxBusyInst", len(metrics.XcpStats[0].GfxBusyInst), 8},
		{"XcpStats.JpegBusy", len(metrics.XcpStats[0].JpegBusy), 40},
		{"XcpStats.VcnBusy", len(metrics.XcpStats[0].VcnBusy), 4},
		{"XcpStats.GfxBusyAcc", len(metrics.XcpStats[0].GfxBusyAcc), 8},
		{"XcpStats.GfxBelowHostLimitAcc", len(metrics.XcpStats[0].GfxBelowHostLimitAcc), 8},
		{"XcpStats.GfxBelowHostLimitPptAcc", len(metrics.XcpStats[0].GfxBelowHostLimitPptAcc), 8},
		{"XcpStats.GfxBelowHostLimitThmAcc", len(metrics.XcpStats[0].GfxBelowHostLimitThmAcc), 8},
		{"XcpStats.GfxLowUtilizationAcc", len(metrics.XcpStats[0].GfxLowUtilizationAcc), 8},
		{"XcpStats.GfxBelowHostLimitTotalAcc", len(metrics.XcpStats[0].GfxBelowHostLimitTotalAcc), 8},
	}

	for _, check := range checks {
		if check.got != check.want {
			t.Errorf("len(%s) = %d, want %d", check.name, check.got, check.want)
		}
	}
}
