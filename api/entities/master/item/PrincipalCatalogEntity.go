package masteritementities

var CreatePrincipalCatalogTable = "mtr_principal_catalog"

type PrincipalCatalog struct {
	IsActive                    bool   `gorm:"column:is_active;not null"`
	PrincipalCatalogId          int    `gorm:"column:principal_catalog_id;size:30;primaryKey" json:"principal_catalog_id"`
	PrincipalCatalogCode        string `gorm:"column:principal_catalog_code;unique;not null" json:"principal_catalog_code"`
	PrincipalCatalogDescription string `gorm:"column:principal_catalog_description;not null" json:"principal_catalog_description"`
}

func (*PrincipalCatalog) TableName() string {
	return CreatePrincipalCatalogTable
}
