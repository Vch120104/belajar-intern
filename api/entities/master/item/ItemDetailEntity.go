package masteritementities

var CreateItemDetailTable = "mtr_item_detail"

type ItemDetail struct {
	IsActive     bool    `gorm:"column:is_active;not null" json:"is_active"`
	ItemDetailId int     `gorm:"column:item_detail_id;size:30;not null;primaryKey" json:"item_detail_id"`
	ItemId       int     `gorm:"column:item_id;size:30;not null;uniqueindex:idx_item_detail" json:"item_id"`
	BrandId      int     `gorm:"column:brand_id;size:30;not null;uniqueindex:idx_item_detail" json:"brand_id"`
	ModelId      int     `gorm:"column:model_id;size:30;not null;uniqueindex:idx_item_detail" json:"model_id"`
	VariantId    int     `gorm:"column:variant_id;size:30;not null;uniqueindex:idx_item_detail" json:"variant_id"`
	MileageEvery float64 `gorm:"column:mileage_every" json:"mileage_every"`
	ReturnEvery  float64 `gorm:"column:return_every" json:"return_every"`
	Item         Item    `gorm:"foreignKey:ItemId;references:ItemId"`
}

func (*ItemDetail) TableName() string {

	return CreateItemDetailTable
}
