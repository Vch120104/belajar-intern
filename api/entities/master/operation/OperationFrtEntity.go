package masteroperationentities

const TableNameOperationFrt = "mtr_operation_frt"

type OperationFrt struct {
	IsActive                bool    `gorm:"column:is_active;not null;default:true" json:"is_active"`
	OperationFrtId          int     `gorm:"column:operation_frt_id;not null;primaryKey" json:"operation_frt_id"`
	OperationModelMappingId int     `gorm:"column:operation_model_mapping_id;not null;" json:"operation_model_mapping_id"`
	VariantId               int     `gorm:"column:variant_id;size:30;not null" json:"variant_id"`
	FrtHour                 float64 `gorm:"column:frt_hour;null" json:"frt_hour"`
	FrtHourExpress          float64 `gorm:"column:frt_hour_2express;null" json:"frt_hour_2express"`
}

func (*OperationFrt) TableName() string {
	return TableNameOperationFrt
}
