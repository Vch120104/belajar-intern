package masteroperationentities

const TableNameOperationEntries = "mtr_operation_entries"

type OperationEntries struct {
	IsActive             bool             `gorm:"column:is_active;not null;default:true" json:"is_active"`
	OperationEntriesId   int32            `gorm:"column:operation_entries_id;not null;primaryKey"  json:"operation_entries_id"`
	OperationEntriesCode string           `gorm:"column:operation_entries_code;not null"  json:"operation_entries_code"`
	OperationGroupId     int32            `gorm:"column:operation_group_id;not null"  json:"operation_group_id"`
	OperationGroup       OperationGroup   
	OperationSectionId   int32            `gorm:"column:operation_section_id;not null"  json:"operation_section_id"`
	OperationSection     OperationSection 
	OperationKeyId       int32            `gorm:"column:operation_key_id;not null"  json:"operation_key_id"`
	OperationKey         OperationKey     
	OperationEntriesDesc string           `gorm:"column:operation_entries_desc;not null"  json:"operation_entries_desc"`
}

func (*OperationEntries) TableName() string {
	return TableNameOperationEntries
}
