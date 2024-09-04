package transactionjpcbpayloads

type BayMasterGetAllResponse struct {
	CarWashBayId          int    `json:"car_wash_bay_id"`
	CarWashBayCode        string `json:"car_wash_bay_code"`
	CarWashBayDescription string `json:"car_wash_bay_description"`
	IsActive              bool   `json:"is_active"`
}

type BayMasterUpdateRequest struct {
	CompanyId    int  `json:"company_id"`
	CarWashBayId int  `json:"car_wash_bay_id"`
	IsActive     bool `json:"is_active"`
}

type CarWashBayDropDownResponse struct {
	CarWashBayId          int    `json:"car_wash_bay_id"`
	CarWashBayDescription string `json:"car_wash_bay_description"`
	CarWashBayCode        string `json:"car_wash_bay_code"`
	IsActive              bool   `json:"is_active"`
}
