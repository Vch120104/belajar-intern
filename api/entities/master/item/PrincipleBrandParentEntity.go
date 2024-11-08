package masteritementities

var CreatePrincipalBrandParentTable = "mtr_principal_brand_parent"

type PrincipalBrandParent struct {
	IsActive                        bool   `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	PrincipalBrandParentId          int    `gorm:"column:principal_brand_parent_id;size:30;primaryKey"        json:"principal_brand_parent_id"`
	PrincipalBrandParentCode        string `gorm:"column:principal_brand_parent_code;unique;size:10;not null"        json:"principal_brand_parent_code"`
	PrincipalBrandParentDescription string `gorm:"column:principal_brand_parent_description;size:50;not null"        json:"principal_brand_parent_description"`
	PrincipalCatalogId              int    `gorm:"column:principal_catalog_id;size:30;unique" json:"principal_catalog_id"`
	PrincipalCatalog                PrincipalCatalog
}

func (*PrincipalBrandParent) TableName() string {
	return CreatePrincipalBrandParentTable
}
