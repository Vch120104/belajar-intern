package masterentities

var CreatePackageMasterTable = "mtr_package"

type PackageMaster struct {
	IsActive            bool    `gorm:"column:is_active;not null;default:true" json:"is_active"`
	PackageId           int     `gorm:"column:package_id;size:30;primarykey;not null" json:"package_id"`
	PackageCode         string  `gorm:"column:package_code;size:15;not null;unique:true" json:"package_code"`
	ItemGroupId         int     `gorm:"column:item_group_id;size:30;not null" json:"item_group-id"`
	PackageName         string  `gorm:"column:package_name;size:40;not null" json:"package_name"`
	BrandId             int     `gorm:"column:brand_id;size:30;not null" json:"brand_id"`
	ModelId             int     `gorm:"column:model_id;size:30;not null" json:"model_id"`
	VariantId           int     `gorm:"column:variant_id;size:30;not null" json:"variant_id"`
	ProfitCenterId      int     `gorm:"column:profit_center_id;size:30;not null" json:"profit_center_id"`
	PackageSet          bool    `gorm:"column:package_set;not null" json:"package_set"`
	CurrencyId          int     `gorm:"column:currency_id;size:30;not null" json:"currency_id"`
	PackagePrice        float64 `gorm:"column:package_price;not null" json:"package_price"`
	TaxTypeId           int     `gorm:"column:tax_type_id;not null;size:30" json:"tax_type_id"`
	PackageRemark       string  `gorm:"column:package_remark;size:256" json:"package_remark"`
	PackageMasterDetail []PackageMasterDetail `gorm:"foreignKey:PackageId;references:PackageId" json:"package_master_detail"`
}

func (*PackageMaster) TableName() string {
	return CreatePackageMasterTable
}
