package transactionworkshopentities

import "time"

const CreateWorkOrderBlockingChassis = "mtr_work_order_blocking_chassis"

type WorkOrderMasterBlockingChassis struct {
	IsActive             bool      `gorm:"column:is_active;" json:"is_active"`
	BlockingSystemNumber int       `gorm:"column:blocking_system_number;size:30;not null;primaryKey;" json:"blocking_system_number"`
	BlockingCode         string    `gorm:"column:blocking_code;size:30;not null" json:"blocking_code"`
	BlockingName         string    `gorm:"column:blocking_name;size:100;not null" json:"blocking_name"`
	PeriodFrom           time.Time `gorm:"column:period_from;not null" json:"period_from"`
	PeriodTo             time.Time `gorm:"column:period_to;not null" json:"period_to"`
	NeverExpired         bool      `gorm:"column:never_expired;" json:"never_expired"`
	VehicleId            int       `gorm:"column:vehicle_id;size:30;" json:"vehicle_id"`
}

func (*WorkOrderMasterBlockingChassis) TableName() string {
	return CreateWorkOrderBlockingChassis
}
