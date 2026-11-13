package assetdepreciation

import "testing"

func TestMonthlyAmount(t *testing.T) {
	basis := postingBasis{UsefulLifeMonths: 48, AcquisitionCost: 120000000, ResidualPercent: 20}
	if got := monthlyAmount(basis, 0); got != 2000000 {
		t.Fatalf("monthlyAmount() = %v, want 2000000", got)
	}
	if got := monthlyAmount(basis, 95000000); got != 1000000 {
		t.Fatalf("last monthlyAmount() = %v, want remaining 1000000", got)
	}
	if got := monthlyAmount(basis, 96000000); got != 0 {
		t.Fatalf("fully depreciated monthlyAmount() = %v, want 0", got)
	}
}
