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
