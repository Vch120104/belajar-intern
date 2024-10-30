package transactionunitpayloads

type PdiRequest struct {
	BrandID           int    `json:"brand_id"`
	PdiDocumentNumber string `json:"pdi_request_document_number"`
	CompanyID         int    `json:"company_id"`
}

type PdiRequestDetail struct {
	PdiRequestSystemNumber       int `json:"pdi_request_system_number"`
	PdiRequestDetailSystemNumber int `json:"pdi_request_detail_system_number"`
	PdiRequestDetailLineNumber   int `json:"pdi_request_detail_line_number"`
}

type PdiRequestDetailById struct {
	PdiRequestDetailSystemNumber int     `json:"pdi_request_detail_system_number"`
	PdiRequestDetailLineNumber   int     `json:"pdi_request_detail_line_number"`
	PdiRequestSystemNumber       int     `json:"pdi_request_system_number"`
	OperationNumberId            int     `json:"operation_number_id"`
	VehicleId                    int     `json:"vehicle_id"`
	PdiRequestDetailLineRemark   string  `json:"pdi_request_detail_line_remark"`
	Frt                          float64 `json:"frt"`
	WorkOrderSystemNumber        int     `json:"work_order_system_number"`
	BookingSystemNumber          int     `json:"booking_system_number"`
	InvoicePayableSystemNumber   int     `json:"invoice_payable_system_no"`
	PdiStatusDescription         string  `json:"pdi_status_description"`
	VehicleEngineNumber          string  `json:"vehicle_engine_number"`
	VehicleChassisNumber         string  `json:"vehicle_chassis_number"`
	ModelId                      int     `json:"model_id"`
	VariantId                    int     `json:"variant_id"`
	ColourId                     int     `json:"colour_id"`
}

type ProfitCenterResponse struct {
	ProfitCenterId   int    `json:"profit_center_id"`
	ProfitCenterCode string `json:"profit_center_code"`
}

type ApprovalStatus struct {
	ApprovalStatusId   string `json:"approval_status_id"`
	ApprovalStatusCode string `json:"approval_status_code"`
}

type ContractService struct {
	ContractServiceId int `json:"contract_service_id"`
}
