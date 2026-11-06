package supportmaster

type Input struct {
	ID, CompanyGroupID, CompanyID, ParentID, AssetTypeID, MethodID, AccountingGroupID    uint64
	AssetGLID, AccumulatedGLID, ExpenseGLID                                              uint64
	Code, Name, Description, VendorType, TaxNumber, Address, Phone, Email, ContactPerson string
	AccountType, ExternalAccountID, BookType, EffectiveFrom, EffectiveTo                 string
	UsefulLifeMonths                                                                     uint64
	ResidualValuePercent                                                                 string
	IsActive                                                                             bool
}

type Option struct {
	ID, GroupID, CompanyID uint64
	Code, Name, Extra      string
}
