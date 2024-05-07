package masterpayloads

import "time"

type FieldActionResponse struct {
	IsActive                  bool      `json:"is_active"`
	FieldActionSystemNumber   int       `json:"field_action_system_number"`
	ApprovalValue             int       `json:"approval_value"`
	BrandId                   int       `json:"brand_id"`
	FieldActionDocumentNumber string    `json:"field_action_document_number"`
	FieldActionName           string    `json:"field_action_name"`
	FieldActionPeriodFrom     time.Time `json:"field_action_period_from"`
	FieldActionPeriodTo       time.Time `json:"field_action_period_to"`
	IsNeverExpired            bool      `json:"is_never_expired"`
	RemarkPopup               string    `json:"remark_popup"`
	IsCritical                bool      `json:"is_critical"`
	RemarkInvoice             string    `json:"remark_invoice"`
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

// type FieldActionMultiVehicleRequest struct {
// 	IsActive                bool   `json:"is_active"`
// 	CompanyId               int    `json:"company_id"`
// 	FieldActionSystemNumber int    `json:"field_action_system_number"`
// 	VehicleIdArray          string `json:"multiple_vehicle_id"`
// }

// type FieldActionItemDetailRequest struct {
// 	LineTypeId                                 int     `json:"line_type_id"`
// 	ItemOperationCode                          int     `json:"item_operation_code"`
// 	FieldActionFrt                             float64 `json:"field_action_frt"`
// }

// type FieldActionByIdResponse struct {

// 	Data       interface{} `json:"data"`
// 	detail []ResponsePagination

// 	RemarkInvoice           string    `json:"remark_invoice"`
// }
