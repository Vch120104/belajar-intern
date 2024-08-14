package transactionjpcbentities

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"time"
)

const TableNameCarWash = "trx_car_wash"

type CarWash struct {
	CarWashId             int                                   `gorm:"column:car_wash_id;primaryKey;not null" json:"car_wash_id"`
	CompanyId             int                                   `gorm:"column:company_id;not null" json:"company_id"`
	WorkOrderSystemNumber int                                   `gorm:"column:work_order_system_number;unique;not null" json:"work_order_system_number"`
	WorkOrder             transactionworkshopentities.WorkOrder `gorm:"foreignKey:work_order_system_number;references:work_order_system_number" json:"work_order"`
	CarWashBayId          int                                   `gorm:"column:car_wash_bay_id" json:"car_wash_bay_id"`
	CarWashBay            BayMaster                             `gorm:"foreignKey:car_wash_bay_id;references:car_wash_bay_id" json:"car_wash_bay"`
	CarWashStatusId       int                                   `gorm:"column:car_wash_status_id" json:"car_wash_status_id"`
	CarWashStatus         CarWashStatus                         `gorm:"foreignKey:car_wash_status_id;references:car_wash_status_id" json:"car_wash_status"`
	CarWashDate           time.Time                             `gorm:"column:car_wash_date" json:"car_wash_date"`
	StartTime             float32                               `gorm:"column:start_time;type:decimal(4,2)" json:"start_time"`
	EndTime               float32                               `gorm:"column:end_time;type:decimal(4,2)" json:"end_time"`
	ActualTime            float32                               `gorm:"column:actual_time;type:decimal(4,2)" json:"actual_time"`
	CarWashPriorityId     int                                   `gorm:"column:car_wash_priority_id;not null" json:"car_wash_priority_id"`
	CarWashPriority       CarWashPriority                       `gorm:"foreignKey:car_wash_priority_id;references:car_wash_priority_id" json:"priority_status"`
}

func (*CarWash) TableName() string {
	return TableNameCarWash
}
