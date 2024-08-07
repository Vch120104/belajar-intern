package masteroperationentities

var CreateLabourSellingPriceDetailTable = "mtr_labour_selling_price_detail"

type LabourSellingPriceDetail struct {
	IsActive                   bool    `gorm:"column:is_active;not null;default:true" json:"is_active"`
	LabourSellingPriceDetailId int     `gorm:"column:labour_selling_price_detail_id;size:30;not null;primaryKey" json:"labour_selling_price_detail_id"`
	LabourSellingPriceId       int     `gorm:"column:labour_selling_price_id;size:30;not null" json:"labour_selling_price_id"`
	LabourSellingPrice         LabourSellingPrice
	ModelId                    int     `gorm:"column:model_id;size:30;not null" json:"model_id"`     //Fk with mtr_unit_model on sales service
	VariantId                  int     `gorm:"column:variant_id;size:30;not null" json:"variant_id"` //Fk with mtr_unit_variant on sales service
	SellingPrice               float64 `gorm:"column:selling_price;not null" json:"selling_price"`
}

func (*LabourSellingPriceDetail) TableName() string {
	return CreateLabourSellingPriceDetailTable
}
