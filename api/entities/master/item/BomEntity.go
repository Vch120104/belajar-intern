package masteritementities

import (
	"time"
)

var CreateBomTable = "mtr_bom"

type Bom struct {
	IsActive               bool      `gorm:"column:is_active;size:1;not null" json:"is_active"`
	BomMasterId            int       `gorm:"column:bom_master_id;size:30;not null;primaryKey" json:"bom_master_id"`
	BomDetailId            int       `gorm:"column:bom_detail_id;size:30;not null" json:"bom_detail_id"`
	BomMasterQty           int       `gorm:"column:bom_master_qty;size:30;not null" json:"bom_master_qty"`
	BomMasterUom           string    `gorm:"column:bom_master_uom;size:30;not null" json:"bom_master_uom"`
	BomMasterEffectiveDate time.Time `gorm:"column:bom_master_effective_date;not null" json:"bom_master_effective_date"`
	Item                   Item      `gorm:"foreignKey:ItemId" json:"item"` //foreign key to mtr_item table
	ItemId                 int       `gorm:"column:item_id;size:30;not null" json:"item_id"`
}

func (*Bom) TableName() string {

	return CreateBomTable

}
