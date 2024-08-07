package transactionjpcbpayloads

type BayMasterResponse struct {
	CompanyId             int    `json:"company_id"`
	CarWashId             int    `json:"car_wash_id"`
	CarWashBayId          int    `json:"car_wash_bay_id"`
	CarWashBayCode        string `json:"car_wash_bay_code"`
	CarWashBayDescription string `json:"car_wash_bay_description"`
	CarWashStatusId       int    `json:"car_wash_status_id"`
	WorkOrderSystemNumber int    `json:"work_order_system_number"`
	IsActive              bool   `json:"is_active"`
}

type BayMasterUpdateRequest struct {
	CompanyId    int
	CarWashBayId int
	RecordStatus bool
}
