package assetlocationreport

type Asset struct {
	ID         uint64
	Number     string
	Name       string
	TypeName   string
	BrandName  string
	ModelName  string
	Condition  string
	Status     string
	BranchName string
	Location   string
	Area       string
}

type Area struct {
	Name   string
	Assets []Asset
}

type Location struct {
	BranchName string
	Name       string
	Total      int
	Areas      []Area
}

type Report struct {
	TotalAssets    int
	TotalLocations int
	Locations      []Location
}
