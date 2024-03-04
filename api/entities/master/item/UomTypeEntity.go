package masteritementities

var CreateUomTypeTable = "mtr_uom_type"

type UomType struct {
	IsActive    bool   `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	UomTypeId   int    `gorm:"column:uom_type_id;size:30;not null;primaryKey"        json:"uom_type_id"`
	UomTypeCode string `gorm:"column:uom_type_code;unique;size:10;not null"        json:"uom_type_code"`
	UomTypeDesc string `gorm:"column:uom_type_desc;size:100;not null"        json:"uom_type_desc"`
}

func (*UomType) TableName() string {
	return CreateUomTypeTable
}
