package masteroperationentities

const TableNameOperationKey = "mtr_operation_key"

type OperationKey struct {
	IsActive                bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
	OperationKeyId          int    `gorm:"column:operation_key_id;size:30;not null;primaryKey"  json:"operation_key_id"`
	OperationKeyCode        string `gorm:"column:operation_key_code;uniqueIndex:idx_code_id;not null;type:char(5)"  json:"operation_key_code"`
	OperationGroupId        int    `gorm:"column:operation_group_id;size:30;uniqueIndex:idx_code_id;not null"  json:"operation_group_id"` //udah
	OperationGroup          *OperationGroup
	OperationSectionId      int `gorm:"column:operation_section_id;size:30;uniqueIndex:idx_code_id;not null"  json:"operation_section_id"` //udah
	OperationSection        *OperationSection
	OperationKeyDescription string           `gorm:"column:operation_key_description;not null;size:50"  json:"operation_key_description"`
	OperationEntries        OperationEntries `gorm:"foreignKey:OperationKeyId;references:OperationKeyId"`
}

func (*OperationKey) TableName() string {
	return TableNameOperationKey
}
