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
	WorkOrder             transactionworkshopentities.WorkOrder `gorm:"foreignKey:WorkOrderSystemNumber"`
	CarWashBayId          int                                   `gorm:"column:car_wash_bay_id" json:"car_wash_bay_id"`
	CarWashBay            []BayMaster                           `gorm:"foreignKey:CarWashBayId"`
	CarWashStatusId       int                                   `gorm:"column:car_wash_status_id" json:"car_wash_status_id"`
	CarWashStatus         CarWashStatus                         `gorm:"foreignKey:CarWashStatusId"`
	CarWashDate           time.Time                             `gorm:"column:car_wash_date" json:"car_wash_date"`
	StartTime             float32                               `gorm:"column:start_time;type:decimal(4,2)" json:"start_time"`
	EndTime               float32                               `gorm:"column:end_time;type:decimal(4,2)" json:"end_time"`
	ActualTime            float32                               `gorm:"column:actual_time;type:decimal(4,2)" json:"actual_time"`
	PriorityStatusId      int                                   `gorm:"column:priority_status_id;not null" json:"priority_status_id"`
	PriorityStatus        CarWashPriority                       `gorm:"foreignKey:PriorityStatusId"`
}

func (*CarWash) TableName() string {
	return TableNameCarWash
}
