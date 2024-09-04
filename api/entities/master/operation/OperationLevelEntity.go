package masteroperationentities

const TableNameOperationLevel = "mtr_operation_level"

type OperationLevel struct {
	IsActive                bool `gorm:"column:is_active;not null;default:true" json:"is_active"`
	OperationLevelId        int  `gorm:"column:operation_level_id;size:30;not null;primaryKey;autoIncrement" json:"operation_level_id"`
	OperationModelMappingId int  `gorm:"column:operation_model_mapping_id;size:30;not null;" json:"operation_model_mapping_id"`
	OperationEntriesId      int  `gorm:"column:operation_entries_id;not null;size:30"  json:"operation_entries_id"`
}

func (*OperationLevel) TableName() string {
	return TableNameOperationLevel
}
