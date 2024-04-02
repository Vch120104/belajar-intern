package masteroperationentities

const TableNameOperationModelMapping = "mtr_operation_model_mapping"

type OperationModelMapping struct {
	IsActive                bool          `gorm:"column:is_active;not null;default:true" json:"is_active"`
	OperationModelMappingId int           `gorm:"column:operation_model_mapping_id;not null;primaryKey" json:"operation_model_mapping_id"`
	BrandId                 int           `gorm:"column:brand_id;size:30;not null" json:"brand_id"` //fk luar
	ModelId                 int           `gorm:"column:model_id;size:30;not null" json:"model_id"` //fk luar
	OperationId             int           `gorm:"column:operation_id;size:30;not null" json:"operation_id"`
	OperationCode           OperationCode `gorm:"foreignKey:OperationId"`
	OperationUsingIncentive bool          `gorm:"column:operation_using_incentive;null" json:"operation_using_incentive"`
	OperationUsingActual    bool          `gorm:"column:operation_using_actual;null" json:"operation_using_actual"`
	OperationPdi            bool          `gorm:"column:operation_pdi;null" json:"operation_pdi"`
}

func (*OperationModelMapping) TableName() string {
	return TableNameOperationModelMapping
}
