package masteritementities

var CreateItemDetailTable = "mtr_item_detail"

type ItemDetail struct {
	ItemDetailId int  `gorm:"column:item_detail_id;size:30;not null;primaryKey" json:"item_detail_id"`
	IsActive     bool `gorm:"column:is_active;not null" json:"is_active"`
	ItemId       int  `gorm:"column:item_id;size:30;not null" json:"item_id"`
	Item         Item
	BrandId      int     `gorm:"column:brand_id;size:30;not null" json:"brand_id"`
	ModelId      int     `gorm:"column:model_id;size:30;not null" json:"model_id"`
	VariantId    int     `gorm:"column:variant_id;size:30;not null" json:"variant_id"`
	MillageEvery float64 `gorm:"column:millage_every" json:"millage_every"`
	ReturnEvery  float64 `gorm:"column:return_every" json:"return_every"`
}

func (*ItemDetail) TableName() string {

	return CreateItemDetailTable
}
