package masterpayloads

import "time"

type FieldActionRequest struct {
	IsActive                  bool      `json:"is_active" parent_entity:"mtr_field_action"`
	FieldActionSystemNumber   int       `json:"field_action_system_number" parent_entity:"mtr_field_action"`
	ApprovalStatusId          int       `json:"approval_status_id" parent_entity:"mtr_field_action"`
	BrandId                   int       `json:"brand_id" parent_entity:"mtr_field_action"`
	FieldActionDocumentNumber string    `json:"field_action_document_number" parent_entity:"mtr_field_action"`
	FieldActionName           string    `json:"field_action_name" parent_entity:"mtr_field_action"`
	FieldActionPeriodFrom     time.Time `json:"field_action_period_from" parent_entity:"mtr_field_action"`
	FieldActionPeriodTo       time.Time `json:"field_action_period_to" parent_entity:"mtr_field_action"`
	IsNeverExpired            bool      `json:"is_never_expired" parent_entity:"mtr_field_action"`
	RemarkPopup               string    `json:"remark_popup" parent_entity:"mtr_field_action"`
	IsCritical                bool      `json:"is_critical" parent_entity:"mtr_field_action"`
	RemarkInvoice             string    `json:"remark_invoice" parent_entity:"mtr_field_action"`
}

type FieldActionResponse struct {
	IsActive                  bool   `json:"is_active" parent_entity:"mtr_field_action"`
	FieldActionSystemNumber   int    `json:"field_action_system_number" parent_entity:"mtr_field_action" main_table:"mtr_field_action"`
	ApprovalStatusId          int    `json:"approval_status_id" parent_entity:"mtr_field_action"`
	BrandId                   int    `json:"brand_id" parent_entity:"mtr_field_action"`
	FieldActionDocumentNumber string `json:"field_action_document_number" parent_entity:"mtr_field_action"`
	FieldActionName           string `json:"field_action_name" parent_entity:"mtr_field_action"`
	FieldActionPeriodFrom     string `json:"field_action_period_from" parent_entity:"mtr_field_action"`
	FieldActionPeriodTo       string `json:"field_action_period_to" parent_entity:"mtr_field_action"`
	IsNeverExpired            bool   `json:"is_never_expired" parent_entity:"mtr_field_action"`
	RemarkPopup               string `json:"remark_popup" parent_entity:"mtr_field_action"`
	IsCritical                bool   `json:"is_critical" parent_entity:"mtr_field_action"`
	RemarkInvoice             string `json:"remark_invoice" parent_entity:"mtr_field_action"`
	VehicleId                 int    `json:"vehicle_id" parent_entity:"mtr_field_action_eligible_vehicle" main_table:"mtr_field_action"`
}

type FieldActionDetailResponse struct {
	IsActive                               bool      `json:"is_active"`
	FieldActionEligibleVehicleSystemNumber int       `json:"field_action_eligible_vehicle_system_number"`
	FieldActionRecallLineNumber            int       `json:"field_action_recall_line_number"`
	FieldActionSystemNumber                int       `json:"field_action_system_number"`
	VehicleId                              int       `json:"vehicle_id"`
	CompanyId                              int       `json:"company_id"`
	FieldActionDate                        time.Time `json:"field_action_date"`
	FieldActionHasTaken                    bool      `json:"field_action_has_taken"`
}

type FieldActionItemDetailResponse struct {
	IsActive                                   bool    `json:"is_active"`
	FieldActionEligibleVehicleItemSystemNumber int     `json:"field_action_eligible_vehicle_item_system_number"`
	FieldActionEligibleVehicleSystemNumber     int     `json:"field_action_eligible_vehicle_system_number"`
	LineTypeId                                 int     `json:"line_type_id"`
	FieldActionEligibleVehicleItemLineNumber   float64 `json:"field_action_eligible_vehicle_item_line_number"`
	ItemOperationCode                          int     `json:"item_operation_code"`
	FieldActionFrt                             float64 `json:"field_action_frt"`
}

type FieldActionOperationDetailResponse struct {
	IsActive                                   bool    `json:"is_active"`
	FieldActionEligibleVehicleItemSystemNumber int     `json:"field_action_eligible_vehicle_operation_system_number"`
	FieldActionEligibleVehicleSystemNumber     int     `json:"field_action_eligible_vehicle_system_number"`
	LineTypeId                                 int     `json:"line_type_id"`
	FieldActionEligibleVehicleItemLineNumber   float64 `json:"field_action_eligible_vehicle_operation_line_number"`
	ItemOperationCode                          int     `json:"item_operation_code"`
	FieldActionFrt                             float64 `json:"field_action_frt"`
}

type ApprovalStatusResponse struct {
	ApprovalStatusId   int    `json:"approval_status_id"`
	ApprovalStatusName string `json:"approval_status_description"`
	ApprovalStatusCode int    `json:"approval_status_code"`
}

type VehicleChassisResponse struct {
	VehicleId            int    `json:"vehicle_id"`
	VehicleChassisNumber string `json:"vehicle_chassis_number"`
}

type FieldActionEligibleVehicleItem struct {
	IsActive                                   bool    `json:"is_active"`
	FieldActionEligibleVehicleItemSystemNumber int     `json:"field_action_eligible_vehicle_item_system_number"`
	FieldActionEligibleVehicleSystemNumber     int     `json:"field_action_eligible_vehicle_system_number"`
	LineTypeId                                 int     `json:"line_type_id"`
	FieldActionEligibleVehicleItemLineNumber   float64 `json:"field_action_eligible_vehicle_item_line_number"`
	ItemId                                     int     `json:"item_id"`
	ItemName                                   string  `json:"item_name"`
	FieldActionFrt                             float64 `json:"field_action_frt"`
}

type FieldActionEligibleVehicleOperation struct {
	IsActive                                   bool    `json:"is_active"`
	FieldActionEligibleVehicleItemSystemNumber int     `json:"field_action_eligible_vehicle_item_system_number"`
	FieldActionEligibleVehicleSystemNumber     int     `json:"field_action_eligible_vehicle_system_number"`
	LineTypeId                                 int     `json:"line_type_id"`
	FieldActionEligibleVehicleItemLineNumber   float64 `json:"field_action_eligible_vehicle_item_line_number"`
	OperationModelMappingId                    int     `json:"operation_model_mapping_id"`
	OperationName                              string  `json:"operation_name"`
	FieldActionFrt                             float64 `json:"field_action_frt"`
}

type FieldActionEligibleVehicleItemOperationResp struct {
	IsActive                                            bool    `json:"is_active"`
	FieldActionEligibleVehicleItemOperationSystemNumber int     `json:"field_action_eligible_vehicle_item_operation_system_number"`
	FieldActionEligibleVehicleSystemNumber              int     `json:"field_action_eligible_vehicle_system_number"`
	LineTypeId                                          int     `json:"line_type_id"`
	ItemId                                              int     `json:"item_id"`
	ItemName                                            string  `json:"item_name"`
	ItemCode                                            string  `json:"item_code"`
	OperationId                                         int     `json:"operation_id"`
	OperationCode                                       string  `json:"operation_code"`
	OperationName                                       string  `json:"operation_name"`
	FieldActionFrt                                      float64 `json:"field_action_frt"`
}

type FieldActionEligibleVehicleItemOperationRequest struct {
	IsActive                                            bool    `json:"is_active"`
	FieldActionEligibleVehicleItemOperationSystemNumber int     `json:"field_action_eligible_vehicle_item_operation_system_number"`
	FieldActionEligibleVehicleSystemNumber              int     `json:"field_action_eligible_vehicle_system_number"`
	LineTypeId                                          int     `json:"line_type_id"`
	ItemOperationId                                     int     `json:"item_operation_id"`
	FieldActionFrt                                      float64 `json:"field_action_frt"`
}
