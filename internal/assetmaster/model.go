package assetmaster

type Item struct {
	ID, GroupID, ParentID, BrandID   uint64
	GroupName, ParentName, BrandName string
	Code, Name, Description          string
	BrandIDs, TypeIDs                string
	IsActive                         bool
}

type Option struct {
	ID, GroupID, ParentID uint64
	Name, GroupName       string
}

type Input struct {
	GroupID, ParentID, BrandID uint64
	Code, Name, Description    string
	IsActive                   bool
}
