package transactionunitpayloads

import (
	"time"
)

type PdiRequest struct {
	BrandID              int     `json:"brand_id"`
	PdiDocumentNumber    string  `json:"pdi_document_number"`
	ModelID              int     `json:"model_id"`
	VariantID            int     `json:"variant_id"`
	VehicleID            int     `json:"vehicle_id"`
	CompanyID            int     `json:"company_id"`
	OperationNumber      string  `json:"operation_number"`
	Frt                  float64 `json:"frt"`
	ContractSystemNumber int     `json:"contract_system_number"`
}

type PdiRequestDetail struct {
	PdiRequestSystemNumber             int       `json:"pdi_request_system_number"`
	PdiRequestDetailLineNumber         int       `json:"pdi_request_detail_line_number"`
	OperationNumberId                  int       `json:"operation_number_id"`
	VehicleId                          int       `json:"vehicle_id"`
	EstimatedDelivery                  time.Time `json:"estimated_delivery"`
	PdiRequestDetailLineRemark         int       `json:"pdi_request_detail_line_remark"`
	Frt                                float64   `json:"frt"`
	ServiceDate                        time.Time `json:"service_date"`
	ServiceTime                        time.Time `json:"service_time"`
	PdiRequestDetailLineStatusId       int       `json:"pdi_request_detail_line_status_id"`
	BookingSystemNumber                int       `json:"booking_system_number"`
	WorkOrderSystemNumber              int       `json:"work_order_system_number"`
	InvoicePayableSystemNumber         int       `json:"invoice_payable_system_number"`
	VehicleRegistrationCertificateTnkb string    `json:"vehicle_registration_certificate_tnkb"`
}

type ProfitCenterResponse struct {
	ProfitCenterId   int    `json:"profit_center_id"`
	ProfitCenterCode string `json:"profit_center_code"`
}

type ApprovalStatus struct {
	ApprovalStatusId          int    `json:"approval_status_id"`
	ApprovalStatusDescription string `json:"approval_status_description"`
}

type ContractService struct {
	ContractServiceId int `json:"contract_service_id"`
}
