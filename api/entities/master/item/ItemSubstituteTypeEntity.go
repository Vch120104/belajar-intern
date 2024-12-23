package masteritementities

const CreateItemSubstituteType = "mtr_substitute_type"

type ItemSubstituteType struct {
	ItemSubstituteTypeId          int    `gorm:"column:item_substitute_type_id;not null;primaryKey;size:30" json:"item_substitute_type_id"`
	ItemSubstituteTypeCode        string `gorm:"column:item_substitute_type_code;not null" json:"item_substitute_type_code"`
	ItemSubstituteTypeDescription string `gorm:"column:item_substitute_type_description;not null" json:"item_substitute_type_description"`
}

func (*ItemSubstituteType) TableName() string {
	return CreateItemSubstituteType
}
