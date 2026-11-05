package assetorganization

type CompanyGroup struct {
	ID           uint64
	Code         string
	Name         string
	Description  string
	Address      string
	Phone        string
	Email        string
	IsActive     bool
	CompanyCount int
}

type Company struct {
	ID          uint64
	GroupID     uint64
	GroupName   string
	Code        string
	Name        string
	LegalName   string
	TaxNumber   string
	Address     string
	Phone       string
	Email       string
	IsActive    bool
	BranchCount int
}

type Branch struct {
	ID             uint64
	CompanyID      uint64
	CompanyName    string
	Code           string
	Name           string
	BranchType     string
	BranchTypeName string
	Address        string
	City           string
	Province       string
	PostalCode     string
	Phone          string
	IsActive       bool
	LocationCount  int
}

type Department struct {
	ID          uint64
	GroupID     uint64
	GroupName   string
	ParentID    uint64
	ParentName  string
	Code        string
	Name        string
	Description string
	IsActive    bool
	ChildCount  int
}

type AssetLocation struct {
	ID               uint64
	BranchID         uint64
	BranchName       string
	ParentID         uint64
	ParentName       string
	Code             string
	Name             string
	LocationType     string
	LocationTypeName string
	Description      string
	IsActive         bool
	ChildCount       int
}

type GroupInput struct {
	Code, Name, Description, Address, Phone, Email string
	IsActive                                       bool
}
type CompanyInput struct {
	GroupID                                                 uint64
	Code, Name, LegalName, TaxNumber, Address, Phone, Email string
	IsActive                                                bool
}
type BranchInput struct {
	CompanyID                                                          uint64
	Code, Name, BranchType, Address, City, Province, PostalCode, Phone string
	IsActive                                                           bool
}
type DepartmentInput struct {
	GroupID, ParentID       uint64
	Code, Name, Description string
	IsActive                bool
}
type LocationInput struct {
	BranchID, ParentID                    uint64
	Code, Name, LocationType, Description string
	IsActive                              bool
}
