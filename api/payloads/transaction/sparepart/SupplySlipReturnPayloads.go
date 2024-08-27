package transactionsparepartpayloads

type SupplySlipReturnSearchResponse struct {
	SupplyReturnSystemNumber   int    `json:"supply_return_system_number"`
	SupplyReturnDocumentNumber string `json:"supply_return_document_number"`
	SupplyReturnDate           string `json:"supply_return_date"`
	SupplyDocumentNumber       string `json:"supply_document_number"`
	WorkOrderDocumentNumber    string `json:"work_order_document_number"`
	CustomerId                 int    `json:"customer_id"`             //external - general (attribute from trx_work_order)
	SupplyReturnStatusId       int    `json:"supply_return_status_id"` //external - general
}

type SupplyReturnStatusResponse struct {
	SupplyReturnStatusId          int    `json:"approval_status_id"`
	SupplyReturnStatusDescription string `json:"approval_status_description"`
}

type SupplySlipReturnResponse struct {
	SupplyReturnSystemNumber   int    `json:"supply_return_system_number"`
	SupplyReturnDocumentNumber string `json:"supply_return_document_number"`
	SupplyReturnDate           string `json:"supply_return_date"`
	SupplyReturnStatusId       int    `json:"supply_return_status_id"` //external - general
	SupplySystemNumber         int    `json:"supply_system_number"`
}
