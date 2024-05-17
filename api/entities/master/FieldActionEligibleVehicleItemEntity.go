package masterentities

var CreateFieldActionEligibleVehicleItemTable = "mtr_field_action_eligible_vehicle_item"

type FieldActionEligibleVehicleItem struct {
	IsActive                                   bool                       `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	FieldActionEligibleVehicleItemSystemNumber int                        `gorm:"column:field_action_eligible_vehicle_item_system_number;not null;primaryKey"        json:"field_action_eligible_vehicle_item_system_number"`
	FieldActionEligibleVehicleSystemNumber     int                        `gorm:"column:field_action_eligible_vehicle_system_number;not null"        json:"field_action_eligible_vehicle_system_number"`
	FieldActionEligibleVehicle                 FieldActionEligibleVehicle `gorm:"foreignKey:FieldActionEligibleVehicleSystemNumber"`
	LineTypeId                                 int                        `gorm:"column:line_type_id;not null"        json:"line_type_id"`
	FieldActionEligibleVehicleItemLineNumber   float64                    `gorm:"column:field_action_eligible_vehicle_item_line_number;null"        json:"field_action_eligible_vehicle_item_line_number"`
	ItemOperationCode                          int                        `gorm:"column:item_operation_code;not null"        json:"item_operation_code"`
	FieldActionFrt                             float64                    `gorm:"column:field_action_frt;not null"        json:"field_action_frt"`
	// FieldActionHasTaken                        bool    `gorm:"column:field_action_has_taken;null"        json:"field_action_has_taken"`
}

func (*FieldActionEligibleVehicleItem) TableName() string {
	return CreateFieldActionEligibleVehicleItemTable
}
