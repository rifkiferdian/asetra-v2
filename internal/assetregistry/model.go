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

type FormOption struct {
	ID        uint64
	ParentID  uint64
	GroupID   uint64
	RelatedID uint64
	Code      string
	Name      string
	Label     string
}

type DepreciationPreview struct {
	CompanyID, GroupID, TypeID, UsefulLifeMonths           uint64
	MethodName, Rate, Residual, EffectiveFrom, EffectiveTo string
}

type AssetFormOptions struct {
	Companies, Branches, Departments, Locations []FormOption
	Types, Brands, Models, Variants             []FormOption
	Conditions, Statuses, Vendors               []FormOption
	DepreciationPolicies                        []DepreciationPreview
}

type AssetInput struct {
	CompanyID, BranchID, DepartmentID, LocationID uint64
	TypeID, BrandID, ModelID, VariantID           uint64
	ConditionID, StatusID, VendorID               uint64
	AssetNumber, Name, SerialNumber, Barcode      string
	Description, AcquisitionDate, AcquisitionCost string
	WarrantyUntil, Notes                          string
}
