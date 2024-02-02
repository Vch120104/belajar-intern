package masteroperationentities

const TableNameOperationGroup = "mtr_operation_group"

type OperationGroup struct {
	IsActive                  bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
	OperationGroupId          int    `gorm:"column:operation_group_id;not null;primaryKey;size:30" json:"operation_group_id"`
	OperationGroupCode        string `gorm:"column:operation_group_code;unique;size:2;type:char(2);not null" json:"operation_group_code"`
	OperationGroupDescription string `gorm:"column:operation_group_description;not null;size:50" json:"operation_group_description"`
}

func (*OperationGroup) TableName() string {
	return TableNameOperationGroup
}
