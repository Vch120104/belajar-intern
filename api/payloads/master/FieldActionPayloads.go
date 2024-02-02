package masterpayloads

import "time"

type FieldActionResponse struct {
	IsActive                bool      `json:"is_active"`
	FieldActionSystemNumber int       `json:"field_action_system_number"`
	ApprovalValue           int       `json:"approval_value"`
	BrandId                 int       `json:"brand_id"`
	FieldActionDocumentNo   string    `json:"field_action_document_no"`
	FieldActionName         string    `json:"field_action_name"`
	FieldActionPeriodFrom   time.Time `json:"field_action_period_from"`
	FieldActionPeriodTo     time.Time `json:"field_action_period_to"`
	IsNeverExpired          bool      `json:"is_never_expired"`
	RemarkPopup             string    `json:"remark_popup"`
	IsCritical              bool      `json:"is_critical"`
	RemarkInvoice           string    `json:"remark_invoice"`
}
