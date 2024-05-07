package masterpayloads

type MovingCodeResponse struct {
	IsActive              bool    `json:"is_active"`
	MovingCodeId          int     `json:"moving_code_id"`
	MovingCode            string  `json:"moving_code"`
	CompanyId             int     `json:"company_id"`
	MovingCodeDescription string  `json:"moving_code_description"`
	MinimumQuantityDemand float64 `json:"minimum_quantity_demand"`
	Priority              float64 `json:"priority"`
	AgingMonthFrom        float64 `json:"aging_month_from"`
	AgingMonthTo          float64 `json:"aging_month_to"`
	DemandExistMonthFrom  float64 `json:"demand_exist_month_from"`
	DemandExistMonthTo    float64 `json:"demand_exist_month_to"`
	LastMovingMonthFrom   float64 `json:"last_moving_month_from"`
	LastMovingMonthTo     float64 `json:"last_moving_month_to"`
	Remark                string  `json:"remark"`
}

type MovingCodeRequest struct {
	IsActive              bool    `json:"is_active"`
	MovingCodeId          int     `json:"moving_code_id"`
	CompanyId             int     `json:"company_id"`
	MovingCodeDescription string  `json:"moving_code_description"`
	MovingCode            string  `json:"moving_code"`
	DemandExistMonthFrom  float64 `json:"demand_exist_month_from"`
	DemandExistMonthTo    float64 `json:"demand_exist_month_to"`
	AgingMonthFrom        float64 `json:"aging_month_from"`
	AgingMonthTo          float64 `json:"aging_month_to"`
	LastMovingMonthFrom   float64 `json:"last_moving_month_from"`
	LastMovingMonthTo     float64 `json:"last_moving_month_to"`
	MinimumQuantityDemand float64 `json:"minimum_quantity_demand"`
	Remark                string  `json:"remark"`
}
