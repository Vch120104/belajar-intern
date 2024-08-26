package transactionjpcbentities

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"time"
)

const TableNameCarWash = "trx_car_wash"

type CarWash struct {
	CarWashId             int                                   `gorm:"column:car_wash_id;primaryKey;not null;size:30" json:"car_wash_id"`
	CompanyId             int                                   `gorm:"column:company_id;not null;size:30" json:"company_id"`
	WorkOrderSystemNumber int                                   `gorm:"column:work_order_system_number;unique;not null;size:30" json:"work_order_system_number"`
	BayId                 int                                   `gorm:"column:car_wash_bay_id;size:30" json:"car_wash_bay_id"`
	StatusId              int                                   `gorm:"column:car_wash_status_id;size:30" json:"car_wash_status_id"`
	PriorityId            int                                   `gorm:"column:car_wash_priority_id;not null;size:30" json:"car_wash_priority_id"`
	CarWashDate           time.Time                             `gorm:"column:car_wash_date" json:"car_wash_date"`
	StartTime             float32                               `gorm:"column:start_time;type:decimal(4,2)" json:"start_time"`
	EndTime               float32                               `gorm:"column:end_time;type:decimal(4,2)" json:"end_time"`
	ActualTime            float32                               `gorm:"column:actual_time;type:decimal(4,2)" json:"actual_time"`
	CarWashStatus         CarWashStatus                         `gorm:"foreignKey:StatusId" json:"car_wash_status"`
	CarWashPriority       CarWashPriority                       `gorm:"foreignKey:PriorityId" json:"priority_status"`
	CarWashBay            BayMaster                             `gorm:"foreignKey:BayId" json:"car_wash_bay"`
	WorkOrder             transactionworkshopentities.WorkOrder `gorm:"foreignKey:WorkOrderSystemNumber" json:"work_order"`
}

func (*CarWash) TableName() string {
	return TableNameCarWash
}
