package transactionsparepartpayloads

type GetAllItemqueryAllCompanyResponse struct {
	CompanyId      int    `json:"company_id"`
	ItemId         int    `json:"item_id"`
	ItemName       string `json:"item_name"`
	QuantityOnHand int    `json:"quantity_on_hand"`
	MovingCodeId   int    `json:"moving_code_id"`
	MovingCode     string `json:"moving_code"`
	PeriodYear     int    `json:"period_year"`
	PeriodMonth    string `json:"period_month"`
}

type GetItemQueryAllCompanyDownloadResponse struct {
	CompanyId      int     `json:"company_id"`
	CompanyCode    string  `json:"company_code"`
	CompanyName    string  `json:"company_name"`
	ItemCode       string  `json:"item_code"`
	ItemName       string  `json:"item_name"`
	MovingCode     string  `json:"moving_code"`
	QuantityOnHand float64 `json:"quantity_on_hand"`
}
