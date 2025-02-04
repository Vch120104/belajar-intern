package transactionworkshoppayloads

type PrintGatePassResponse struct {
	WorkOrderSystemNumber   int    `json:"work_order_system_number"`
	WorkOrderDocumentNumber string `json:"work_order_document_number"`
	WorkOrderDate           string `json:"work_order_date"`
	CustomerId              int    `json:"customer_id"`
	CustomerName            string `json:"customer_name"`
	VehicleId               int    `json:"vehicle_id"`
	VehicleBrandId          int    `json:"vehicle_brand_id"`
	ModelId                 int    `json:"model_id"`
	GatePassSystemNumber    int    `json:"gate_pass_system_number"`
	GatePassDocumentNumber  string `json:"gate_pass_document_number"`
	GatePassDate            string `json:"gate_pass_date"`
	DeliveryName            string `json:"delivery_name"`
	DeliveryAddress         string `json:"delivery_address"`
}

type PrintGatePassRequest struct {
	WorkOrderSystemNumber int `json:"work_order_system_number"`
}
