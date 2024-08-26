package transactionjpcbentities

const TableNameCarWashStatus = "mtr_car_wash_status"

type CarWashStatus struct {
	IsActive                 bool   `gorm:"column:is_active;default:true;not null" json:"is_active"`
	CarWashStatusId          int    `gorm:"column:car_wash_status_id;primaryKey;not null;size:30" json:"car_wash_status_id"`
	CarWashStatusCode        string `gorm:"column:car_wash_status_code;not null;type:varchar(10);unique" json:"car_wash_status_code"`
	CarWashStatusDescription string `gorm:"column:car_wash_status_description;not null;type:varchar(50)" json:"car_wash_status_description"`
}

func (*CarWashStatus) TableName() string {
	return TableNameCarWashStatus
}
