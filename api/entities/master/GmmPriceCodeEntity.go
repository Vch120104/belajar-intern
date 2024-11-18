package masterentities

const TableNameGmmPriceCode = "mtr_gmm_price_code"

type GmmPriceCode struct {
	IsActive       bool   `gorm:"column:is_active;default:false;not null" json:"is_active"`
	GmmPriceCodeId int    `gorm:"column:gmm_price_code_id;size:30;primaryKey" json:"gmm_price_code_id"`
	GmmPriceCode   string `gorm:"column:gmm_price_code;size:10;unique;not null" json:"gmm_price_code"`
	GmmPriceDesc   string `gorm:"column:gmm_price_desc;size:50" json:"gmm_price_desc"`
}

func (*GmmPriceCode) TableName() string {
	return TableNameGmmPriceCode
}
