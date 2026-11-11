package assetregistry

import (
	"testing"
	"time"
)

func TestCalculateStraightLineDepreciation(t *testing.T) {
	asset := Asset{AcquisitionDate: "2025-01-15", AcquisitionCost: "10000000"}
	policy := DepreciationPolicy{
		AssetTypeName: "Laptop", MethodCode: "SL", MethodName: "Garis Lurus",
		BookType: "ACCOUNTING", UsefulLifeMonths: 48, ResidualPercent: "0.00",
		AnnualRatePercent: "25.00", EffectiveFrom: "2025-01-01",
	}

	got, err := calculateDepreciation(asset, policy, time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if got.ElapsedMonths != 12 {
		t.Fatalf("elapsed months = %d, want 12", got.ElapsedMonths)
	}
	if got.MonthlyDepreciation != "Rp208.333" {
		t.Fatalf("monthly depreciation = %q", got.MonthlyDepreciation)
	}
	if got.AccumulatedDepreciation != "Rp2.500.000" {
		t.Fatalf("accumulated depreciation = %q", got.AccumulatedDepreciation)
	}
	if got.BookValue != "Rp7.500.000" {
		t.Fatalf("book value = %q", got.BookValue)
	}
}

func TestCalculateDepreciationStopsAtResidualValue(t *testing.T) {
	asset := Asset{AcquisitionDate: "2020-01-01", AcquisitionCost: "10000000"}
	policy := DepreciationPolicy{
		AssetTypeName: "Laptop", MethodName: "Garis Lurus", UsefulLifeMonths: 48,
		ResidualPercent: "10", AnnualRatePercent: "25", EffectiveFrom: "2020-01-01",
	}

	got, err := calculateDepreciation(asset, policy, time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if got.ElapsedMonths != 48 || got.BookValue != "Rp1.000.000" {
		t.Fatalf("elapsed=%d book value=%q, want 48 and residual value", got.ElapsedMonths, got.BookValue)
	}
}

func TestFormatPercent(t *testing.T) {
	tests := map[string]string{
		"25.0000": "25%",
		"12.5000": "12,5%",
		"12.3456": "12,35%",
		"0.00":    "0%",
	}
	for input, want := range tests {
		if got := formatPercent(input); got != want {
			t.Errorf("formatPercent(%q) = %q, want %q", input, got, want)
		}
	}
}
