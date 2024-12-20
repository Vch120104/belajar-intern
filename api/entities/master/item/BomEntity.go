package masteritementities

import (
	"time"
)

var CreateBomTable = "mtr_bom"

type Bom struct {
	IsActive      bool      `gorm:"column:is_active;size:1;not null" json:"is_active"`
	BomId         int       `gorm:"column:bom_id;size:30;not null;primaryKey" json:"bom_id"`
	EffectiveDate time.Time `gorm:"column:effective_date;not null;type:datetime" json:"effective_date"`
	ItemId        int       `gorm:"column:item_id;size:30;not null" json:"item_id"`
	Qty           float64   `gorm:"column:qty;size:30;not null" json:"qty"`
	Item          Item      `gorm:"foreignKey:ItemId;references:ItemId"`
}

func (*Bom) TableName() string {
	return CreateBomTable
}
