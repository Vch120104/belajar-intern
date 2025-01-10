package masterentities

import "time"

var CreateFieldActionTable = "mtr_field_action"

type FieldAction struct {
	IsActive                   bool                         `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	FieldActionSystemNumber    int                          `gorm:"column:field_action_system_number;size:30;not null;primaryKey"        json:"field_action_system_number"`
	ApprovalStatusId           int                          `gorm:"column:approval_status_id;size:30;not null"        json:"approval_status_id"`
	BrandId                    int                          `gorm:"column:brand_id;size:30;not null"        json:"brand_id"`
	FieldActionDocumentNumber  string                       `gorm:"column:field_action_document_number;size:30;not null"        json:"field_action_document_number"`
	FieldActionName            string                       `gorm:"column:field_action_name;size:100;not null"        json:"field_action_name"`
	FieldActionPeriodFrom      time.Time                    `gorm:"column:field_action_period_from;not null"        json:"field_action_period_from"`
	FieldActionPeriodTo        time.Time                    `gorm:"column:field_action_period_to;not null"        json:"field_action_period_to"`
	IsNeverExpired             bool                         `gorm:"column:is_never_expired;null"        json:"is_never_expired"`
	RemarkPopup                string                       `gorm:"column:remark_popup;size:255;null"        json:"remark_popup"`
	IsCritical                 bool                         `gorm:"column:is_critical;null"        json:"is_critical"`
	RemarkInvoice              string                       `gorm:"column:remark_invoice;size:255;null"        json:"remark_invoice"`
	FieldActionEligibleVehicle []FieldActionEligibleVehicle `gorm:"foreignKey:FieldActionSystemNumber;references:FieldActionSystemNumber" json:"field_action_eligible_vehicle"`
}

func (*FieldAction) TableName() string {
	return CreateFieldActionTable
}
