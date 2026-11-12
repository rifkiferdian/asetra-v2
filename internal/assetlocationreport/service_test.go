package assetlocationreport

import "testing"

func TestBuildReportGroupsAssetsByLocationAndArea(t *testing.T) {
	report := BuildReport([]Asset{
		{BranchName: "MK5", Location: "OFFICE", Area: "KINASIH", Name: "Proyektor"},
		{BranchName: "MK5", Location: "OFFICE", Area: "KINASIH", Name: "PC"},
		{BranchName: "MK5", Location: "OFFICE", Area: "PROGRAMMER", Name: "AC"},
	})
	if report.TotalAssets != 3 || report.TotalLocations != 1 || report.Locations[0].Total != 3 {
		t.Fatalf("unexpected report totals: %#v", report)
	}
	if len(report.Locations[0].Areas) != 2 || len(report.Locations[0].Areas[0].Assets) != 2 {
		t.Fatalf("unexpected area grouping: %#v", report.Locations[0].Areas)
	}
}

func TestBuildReportIncludesAreaWithoutAssets(t *testing.T) {
	report := BuildReport([]Asset{{LocationID: 9, BranchName: "MK5", Location: "OFFICE", Area: "GUDANG"}})
	if report.TotalAssets != 0 || report.TotalLocations != 1 {
		t.Fatalf("unexpected report totals: %#v", report)
	}
	if len(report.Locations[0].Areas) != 1 || report.Locations[0].Areas[0].ID != 9 || len(report.Locations[0].Areas[0].Assets) != 0 {
		t.Fatalf("empty area was not retained: %#v", report.Locations[0].Areas)
	}
}
