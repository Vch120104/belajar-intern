package masteritementities

import "time"

var CreateBomDetailTable = "mtr_bom_detail"

type BomDetail struct {
	BomDetailId            int       `gorm:"column:bom_detail_id;size:30;not null;primaryKey" json:"bom_detail_id"`
	BomMasterId            int       `gorm:"column:bom_master_id;size:30;not null" json:"bom_master_id"`
	BomDetailSeq           int       `gorm:"column:bom_detail_seq;size:30;not null" json:"bom_detail_seq"`
	BomDetailQty           int       `gorm:"column:bom_detail_qty;size:30;not null" json:"bom_detail_qty"`
	BomDetailUom           string    `gorm:"column:bom_detail_uom;size:30;not null" json:"bom_detail_uom"`
	BomDetailEffectiveDate time.Time `gorm:"column:bom_detail_effective_date;not null" json:"bom_detail_effective_date"`
	BomDetailRemark        string    `gorm:"column:bom_detail_remark;size:30;not null" json:"bom_detail_remark"`
	BomDetailCostingPct    int       `gorm:"column:bom_detail_costing_percent;size:30;not null" json:"bom_detail_costing_percent"`
	BomDetailType          string    `gorm:"column:bom_detail_type;size:30;not null" json:"bom_detail_type"`
	BomDetailMaterialCode  string    `gorm:"column:bom_detail_material_code;size:30;not null" json:"bom_detail_material_code"`
	BomDetailMaterialName  string    `gorm:"column:bom_detail_material_name;size:30;not null" json:"bom_detail_material_name"`
	// Material                Material  `gorm:"foreignKey:MaterialId" json:"material"` //foreign key to mtr_item table
	// MaterialId              int       `gorm:"column:material_id;size:30;not null" json:"material_id"`
}

func (*BomDetail) TableName() string {

	return CreateBomDetailTable

}
