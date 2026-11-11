package assetregistry

type Asset struct {
	ID                                                                  uint64
	AssetNumber, Name, SerialNumber, Barcode, Description               string
	CompanyName, BranchName, DepartmentName, LocationName, LocationPath string
	CategoryName, TypeName, BrandName, ModelName, VariantName           string
	ConditionCode, ConditionName, StatusCode, StatusName, VendorName    string
	AcquisitionDate, AcquisitionCost, WarrantyUntil, Notes              string
	CreatedAt, UpdatedAt                                                string
}

type Summary struct {
	Total, Active, Maintenance, Inactive int
}

type DepreciationPolicy struct {
	ID, UsefulLifeMonths                                           uint64
	AssetTypeName, MethodCode, MethodName, BookType, GLAccount     string
	ResidualPercent, AnnualRatePercent, EffectiveFrom, EffectiveTo string
}

type DepreciationDetail struct {
	Available                                                         bool
	PolicyName, MethodName, BookType, GLAccount                       string
	UsefulLife, AnnualRate, ResidualPercent, EffectivePeriod          string
	MonthlyDepreciation, AccumulatedDepreciation, BookValue, AsOfDate string
	ElapsedMonths                                                     uint64
}
