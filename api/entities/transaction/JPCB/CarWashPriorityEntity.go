package transactionjpcbentities

const TableNameCarWashPriority = "mtr_car_wash_priority"

type CarWashPriority struct {
	IsActive                   bool   `gorm:"column:is_active;default:true;not null" json:"is_active"`
	CarWashPriorityId          int    `gorm:"column:car_wash_priority_id;primaryKey;not null" json:"car_wash_priority_id"`
	CarWashPriorityCode        string `gorm:"column:car_wash_priority_code;not null;type:varchar(10);unique" json:"car_wash_priority_code"`
	CarWashPriorityDescription string `gorm:"column:car_wash_priority_description;not null;type:varchar(50)" json:"car_wash_priority_description"`
}

func (*CarWashPriority) TableName() string {
	return TableNameCarWashPriority
}
