package masterentities

var CreateFieldActionEligibleVehicleItemOperationTable = "mtr_field_action_eligible_vehicle_item_operation"

type FieldActionEligibleVehicleItemOperation struct {
	IsActive                                            bool                       `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	FieldActionEligibleVehicleItemOperationSystemNumber int                        `gorm:"column:field_action_eligible_vehicle_item_operation_system_number;primaryKey"`
	FieldActionEligibleVehicleSystemNumber              int                        `gorm:"column:field_action_eligible_vehicle_system_number"`
	FieldActionEligibleVehicle                          FieldActionEligibleVehicle `gorm:"foreignKey:FieldActionEligibleVehicleSystemNumber;references:FieldActionEligibleVehicleSystemNumber;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ItemOperationId                                     int                        `gorm:"column:item_operation_id;not null;size:30"        json:"item_operation_id"`
	FieldActionFrt                                      float64                    `gorm:"column:field_action_frt;not null"        json:"field_action_frt"`
	ItemOperation                                       ItemOperation              `gorm:"foreignKey:ItemOperationId;references:ItemOperationId"`
}

func (*FieldActionEligibleVehicleItemOperation) TableName() string {
	return CreateFieldActionEligibleVehicleItemOperationTable
}
