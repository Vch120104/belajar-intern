package masteritementities

var CreatePrincipleBrandParentTable = "mtr_principle_brand_parent"

type PrincipleBrandParent struct {
	IsActive                        bool   `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	PrincipalBrandParentId          int    `gorm:"column:principal_brand_parent_id;primaryKey"        json:"principal_brand_parent_id"`
	PrincipalBrandParentCode        string `gorm:"column:principal_brand_parent_code;unique;size:10;not null"        json:"principal_brand_parent_code"`
	PrincipalBrandParentDescription string `gorm:"column:principal_brand_parent_description;size:50;not null"        json:"principal_brand_parent_description"`
	CatalogueCode                   string `gorm:"column:catalogue_code;size:10;not null"        json:"catalogue_code"`
}

func (*PrincipleBrandParent) TableName() string {
	return CreatePrincipleBrandParentTable
}
