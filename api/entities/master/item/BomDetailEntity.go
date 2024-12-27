package masteritementities

var CreateBomDetailTable = "mtr_bom_detail"

type BomDetail struct {
	IsActive          bool    `gorm:"column:is_active;size:1;not null" json:"is_active"` // Naturally, this will always be `true`
	BomDetailId       int     `gorm:"column:bom_detail_id;size:30;not null;primaryKey" json:"bom_detail_id"`
	BomId             int     `gorm:"column:bom_id;size:30;not null;index:,unique,composite:un" json:"bom_id"`
	Seq               int     `gorm:"column:seq;size:30;not null" json:"seq"`
	ItemId            int     `gorm:"column:item_id;size:30;not null;index:,unique,composite:un" json:"item_id"`
	Qty               float64 `gorm:"column:qty;size:30;not null" json:"qty"`
	Remark            string  `gorm:"column:remark;size:512" json:"remark"`
	CostingPercentage float64 `gorm:"column:costing_percentage;size:30;not null" json:"costing_percentage"`
	//Bom               Bom     `gorm:"foreignKey:BomId;references:BomId" json:"bom"` // not needed
	//Item              Item    `gorm:"foreignKey:ItemId;references:ItemId" json:"item"` // not needed
}

func (*BomDetail) TableName() string {
	return CreateBomDetailTable
}
