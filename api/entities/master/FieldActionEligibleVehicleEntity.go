package masterentities

import "time"

var CreateFieldActionEligibleVehicleTable = "mtr_field_action_eligible_vehicle"

type FieldActionEligibleVehicle struct {
	// FieldActionDetailSystemNumber int    	`gorm:"column:field_action_detail_system_number;size:30;not null;primaryKey"        json:"field_action_detail_system_number"`
	IsActive                               bool        `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	FieldActionEligibleVehicleSystemNumber int         `gorm:"column:field_action_eligible_vehicle_system_number;size:30;not null;primaryKey"        json:"field_action_eligible_vehicle_system_number"`
	FieldActionRecallLineNumber            int         `gorm:"column:field_action_recall_line_number;size:30;null"        json:"field_action_recall_line_number"`
	FieldActionSystemNumber                int         `gorm:"column:field_action_system_number;size:30;not null"        json:"field_action_system_number"`
	FieldAction                            FieldAction `gorm:"foreignKey:FieldActionSystemNumber"`
	VehicleId                              int         `gorm:"column:vehicle_id;size:30;not null"        json:"vehicle_id"`
	CompanyId                              int         `gorm:"column:company_id;size:30;not null"        json:"company_id"`
	FieldActionDate                        time.Time   `gorm:"column:field_action_date;null"        json:"field_action_date"`
	FieldActionHasTaken                    bool        `gorm:"column:field_action_has_taken;not null;default:false"        json:"field_action_has_taken"`
}

func (*FieldActionEligibleVehicle) TableName() string {
	return CreateFieldActionEligibleVehicleTable
}
