package masteroperationentities

const TableNameOperationEntries = "mtr_operation_entries"

type OperationEntries struct {
	IsActive                    bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
	OperationEntriesId          int    `gorm:"column:operation_entries_id;not null;primaryKey;size:30"  json:"operation_entries_id"`
	OperationEntriesCode        string `gorm:"column:operation_entries_code;not null"  json:"operation_entries_code"`
	OperationGroupId            int    `gorm:"column:operation_group_id;not null;size:30"  json:"operation_group_id"`
	OperationGroup              OperationGroup
	OperationSectionId          int `gorm:"column:operation_section_id;not null;size:30"  json:"operation_section_id"`
	OperationSection            OperationSection
	OperationKeyId              int `gorm:"column:operation_key_id;not null;size:30"  json:"operation_key_id"`
	OperationKey                OperationKey
	OperationLevel              OperationLevel
	OperationEntriesDescription string `gorm:"column:operation_entries_description;not null"  json:"operation_entries_description"`
}

func (*OperationEntries) TableName() string {
	return TableNameOperationEntries
}
