package masteritementities

import "time"

var (
	CreatePurchasePriceTable = "mtr_purchase_price"
)

type PurchasePrice struct {
	PurchasePriceId            int       `gorm:"column:purchase_price_id;size:30;primaryKey" json:"purchase_price_id"`
	SupplierId                 int       `gorm:"column:supplier_id;size:30;not null" json:"supplier_id"`
	CurrencyId                 int       `gorm:"column:currency_id;size:30;not null" json:"currency_id"`
	PurchasePriceEffectiveDate time.Time `json:"purchase_price_effective_date"`
}

func (*PurchasePrice) TableName() string {
	return CreatePurchasePriceTable
}
