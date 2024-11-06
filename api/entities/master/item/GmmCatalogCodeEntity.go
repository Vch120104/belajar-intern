package masteritementities

var CreateGmmCatalogCodeTable = "mtr_gmm_catalog_code"

type GmmCatalogCode struct {
	IsActive              bool   `gorm:"column:is_active;not null;default:true"`
	GmmCatalogId          int    `gorm:"column:gmm_catalog_id;size:30;primaryKey" json:"gmm_catalog_id"`
	GmmCatalogCode        string `gorm:"column:gmm_catalog_code;unique;not null" json:"gmm_catalog_code"`
	GmmCatalogDescription string `gorm:"column:gmm_catalog_description;not null" json:"gmm_catalog_description"`
}

func (*GmmCatalogCode) TableName() string {
	return CreateGmmCatalogCodeTable
}
