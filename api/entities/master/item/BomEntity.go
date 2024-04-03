package masteritementities

import (
	"time"
)

var CreateBomTable = "mtr_bom"

type Bom struct {
	BomMasterId            int         `gorm:"column:bom_master_id;size:30;not null;primaryKey" json:"bom_master_id"`
	IsActive               bool        `gorm:"column:is_active;size:1;not null" json:"is_active"`
	BomMasterSeq           int         `gorm:"column:bom_master_seq;size:30;not null" json:"bom_master_seq"`
	BomMasterQty           int         `gorm:"column:bom_master_qty;size:30;not null" json:"bom_master_qty"`
	BomMasterEffectiveDate time.Time   `gorm:"column:bom_master_effective_date;not null" json:"bom_master_effective_date"`
	BomMasterChangeNumber  int         `gorm:"column:bom_master_change_number;size:30;default:0" json:"bom_master_change_number"`
	ItemId                 int         `gorm:"column:item_id;size:30;not null" json:"item_id"`
	Item                   Item        `gorm:"foreignKey:ItemId" json:"item"` //foreign key to mtr_item table
	BomDetail              []BomDetail `gorm:"foreignKey:BomMasterId" json:"detail_bom"`
	//UomId                  int       `gorm:"column:uom_id;size:30;not null" json:"uom_id"`
	//Uom                    Uom       `gorm:"foreignKey:UomID" json:"uom"` //foreign key to mtr_item to Uom table
}

func (*Bom) TableName() string {

	return CreateBomTable

}
