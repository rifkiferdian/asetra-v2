package assetregistry

type Asset struct {
	ID                                                               uint64
	AssetNumber, Name, SerialNumber, Barcode, Description            string
	CompanyName, BranchName, LocationName, LocationPath              string
	CategoryName, TypeName, BrandName, ModelName, VariantName        string
	ConditionCode, ConditionName, StatusCode, StatusName, VendorName string
	AcquisitionDate, AcquisitionCost, WarrantyUntil, Notes           string
	CreatedAt, UpdatedAt                                             string
}

type Summary struct {
	Total, Active, Maintenance, Inactive int
}
