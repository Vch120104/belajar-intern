package masteroperationentities

const TableNameOperationCode = "mtr_operation_code"

type OperationCode struct {
	IsActive                bool             `gorm:"column:is_active;not null;default:true" json:"is_active"`
	OperationId             int              `gorm:"column:operation_id;size:30;not null;primaryKey;"  json:"operation_id"`
	OperationCode           string           `gorm:"column:operation_code;unique;not null"  json:"operation_code"`
	OperationName           string           `gorm:"column:operation_name;null"  json:"operation_name"`
	OperationUsingIncentive bool             `gorm:"column:operation_using_incentive;default:false"  json:"operation_using_incentive"`
	OperationUsingActual    bool             `gorm:"column:operation_using_actual;default:false"  json:"operation_using_actual"`
	OperationGroup          OperationGroup   `gorm:"foreignkey:operation_id" references:"operation_id"`
	OperationSection        OperationSection `gorm:"foreignkey:operation_id" references:"operation_id"`
	OperationKey            OperationKey     `gorm:"foreignkey:operation_id" references:"operation_id"`
	OperationEntries        OperationEntries `gorm:"foreignkey:operation_id" references:"operation_id"`
}

func (*OperationCode) TableName() string {
	return TableNameOperationCode
}
