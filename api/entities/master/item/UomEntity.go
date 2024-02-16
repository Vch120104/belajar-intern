package masteritementities

var CreateUomTable = "mtr_uom"

type Uom struct {
	IsActive       bool    `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	UomId          int     `gorm:"column:uom_id;size:30;not null;primaryKey"        json:"uom_id"`
	UomTypeId      int     `gorm:"column:uom_type_id;size:30"        json:"uom_type_id"`
	UomType        UomType `gorm:"foreignKey:UomTypeId"`
	UomCode        string  `gorm:"column:uom_code;size:10;not null"        json:"uom_code"`
	UomDescription string  `gorm:"column:uom_description;size:50;not null"        json:"uom_description"`
}

func (*Uom) TableName() string {
	return CreateUomTable
}
