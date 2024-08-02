package masterentities

var CreateFieldActionEligibleVehicleOperationTable = "mtr_field_action_eligible_vehicle_operation"

type FieldActionEligibleVehicleOperation struct {
	IsActive                                        bool                       `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	FieldActionEligibleVehicleOperationSystemNumber int                        `gorm:"column:field_action_eligible_vehicle_item_system_number;size:30;not null;primaryKey"        json:"field_action_eligible_vehicle_operation_system_number"`
	FieldActionEligibleVehicleSystemNumber          int                        `gorm:"column:field_action_eligible_vehicle_system_number;size:30;not null"        json:"field_action_eligible_vehicle_system_number"`
	FieldActionEligibleVehicle                      FieldActionEligibleVehicle `gorm:"foreignKey:FieldActionEligibleVehicleSystemNumber"`
	LineTypeId                                      int                        `gorm:"column:line_type_id;size:30;not null"        json:"line_type_id"`
	FieldActionEligibleVehicleItemLineNumber        float64                    `gorm:"column:field_action_eligible_vehicle_item_line_number;null"        json:"field_action_eligible_vehicle_item_line_number"`
	OperationModelMappingId                         int                        `gorm:"column:operation_model_mapping_id;size:30;not null" json:"operation_model_mapping_id"`
	FieldActionFrt                                  float64                    `gorm:"column:field_action_frt;not null"        json:"field_action_frt"`
	// FieldActionHasTaken                        bool    `gorm:"column:field_action_has_taken;null"        json:"field_action_has_taken"`
}

func (*FieldActionEligibleVehicleOperation) TableName() string {
	return CreateFieldActionEligibleVehicleOperationTable
}
