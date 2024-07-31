package masteritementities

var (
	CreatePurchasePriceDetailTable = "mtr_purchase_price_detail"
)

type PurchasePriceDetail struct {
	IsActive              bool    `gorm:"column:is_active;size:1;not null" json:"is_active"`
	PurchasePriceDetailId int     `gorm:"column:purchase_price_detail_id;size:30;primaryKey" json:"purchase_price_detail_id"`
	PurchasePriceId       int     `gorm:"column:purchase_price_id;size:30;not null" json:"purchase_price_id"`
	ItemId                int     `gorm:"column:item_id;size:30;not null" json:"item_id"`
	PurchasePrice         float64 `gorm:"column:purchase_price;size:30;not null" json:"purchase_price"`
}

func (*PurchasePriceDetail) TableName() string {
	return CreatePurchasePriceDetailTable
}
