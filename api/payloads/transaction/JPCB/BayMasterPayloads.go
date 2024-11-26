package transactionjpcbpayloads

type CarWashBayGetAllResponse struct {
	CarWashBayId          int    `json:"car_wash_bay_id"`
	CarWashBayCode        string `json:"car_wash_bay_code"`
	CarWashBayDescription string `json:"car_wash_bay_description"`
	IsActive              bool   `json:"is_active"`
}

type CarWashBayUpdateRequest struct {
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

type CarWashBayPostRequest struct {
	CompanyId             int    `json:"company_id"`
	CarWashBayCode        string `json:"car_wash_bay_code"`
	CarWashBayDescription string `json:"car_wash_bay_description"`
}

type CarWashBayPutRequest struct {
	CarWashBayID          int    `json:"car_wash_bay_id"`
	CarWashBayCode        string `json:"car_wash_bay_code"`
	CarWashBayDescription string `json:"car_wash_bay_description"`
}
