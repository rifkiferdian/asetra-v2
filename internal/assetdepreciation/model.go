package assetdepreciation

import "time"

type Company struct {
	ID   uint64
	Name string
}

type Row struct {
	AssetID, PolicyID, UsefulLifeMonths uint64
	AssetNumber, AssetName, CompanyName string
	AcquisitionDate, AcquisitionCost    string
	ResidualPercent, Accumulated        string
	BookValue, MonthlyAmount            string
	PostedAmount, PostedAt              string
	CanPost                             bool
}

type Page struct {
	Companies                      []Company
	Rows                           []Row
	Period, PeriodLabel, MaxPeriod string
	CompanyID                      uint64
	TotalAssets, PostedAssets      int
	TotalPosted, TotalBookValue    string
}

type postingBasis struct {
	PolicyID, UsefulLifeMonths uint64
	AcquisitionDate            time.Time
	AcquisitionCost            float64
	ResidualPercent            float64
}
