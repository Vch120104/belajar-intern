package masteritementities

import "time"

var (
	CreatePurchasePriceTable = "mtr_purchase_price"
)

type PurchasePrice struct {
	PurchasePriceId            int                   `gorm:"column:purchase_price_id;size:30;primaryKey" json:"purchase_price_id"`
	IsActive                   bool                  `gorm:"column:is_active;size:1;not null" json:"is_active"`
	SupplierId                 int                   `gorm:"column:supplier_id;size:30;not null" json:"supplier_id"`
	CurrencyId                 int                   `gorm:"column:currency_id;size:30;not null" json:"currency_id"`
	PurchasePriceEffectiveDate time.Time             `gorm:"column:purchase_price_effective_date;size:30;not null;type:datetime" json:"purchase_price_effective_date"`
	PurchasePriceDetail        []PurchasePriceDetail `gorm:"foreignKey:PurchasePriceId" json:"detail_purchase_price"`
}

func (*PurchasePrice) TableName() string {
	return CreatePurchasePriceTable
}
