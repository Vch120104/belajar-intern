package masteroperationentities

const TableNameOperationCode = "mtr_operation_code"

type OperationCode struct {
	IsActive                bool             `gorm:"column:is_active;not null;default:true" json:"is_active"`
	OperationId             int32            `gorm:"column:operation_id;not null;primaryKey"  json:"operation_id"`
	OperationCode           string           `gorm:"column:operation_code;unique;not null"  json:"operation_code"`
	OperationName           string           `gorm:"column:operation_name;null"  json:"operation_name"`
	OperationGroupId        int32            `gorm:"column:operation_group_id;not null"  json:"operation_group_id"`
	OperationGroup          OperationGroup  
	OperationSectionId      int32            `gorm:"column:operation_section_id;not null"  json:"operation_section_id"`
	OperationSection        OperationSection 
	OperationKeyId          int32            `gorm:"column:operation_key_id;not null"  json:"operation_key_id"`
	OperationKey            OperationKey    
	OperationEntriesId      int32            `gorm:"column:operation_entries_id;not null"  json:"operation_entries_id"`
	OperationEntries        OperationEntries 
	OperationUsingIncentive bool             `gorm:"column:operation_using_incentive;null"  json:"operation_using_incentive"`
	OperationUsingActual    bool             `gorm:"column:operation_using_actual;null"  json:"operation_using_actual"`
}

func (*OperationCode) TableName() string {
	return TableNameOperationCode
}
