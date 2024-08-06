package transactionjpcbentities

const TableNameBayMaster = "mtr_car_wash_bay"

type BayMaster struct {
	IsActive              bool   `gorm:"column:is_active;default:true;not null" json:"is_active"`
	CarWashBayId          int    `gorm:"column:car_wash_bay_id;primaryKey;not null" json:"car_wash_bay_id"`
	CarWashBayCode        string `gorm:"column:car_wash_bay_code;not null;type:varchar(10);unique" json:"car_wash_bay_code"`
	CarWashBayDescription string `gorm:"column:car_wash_bay_description;not null;type:varchar(50)" json:"car_wash_bay_description"`
}

func (*BayMaster) TableName() string {
	return TableNameBayMaster
}
