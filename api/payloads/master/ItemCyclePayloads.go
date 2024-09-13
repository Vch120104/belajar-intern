package masterpayloads

type ItemCycleInsertPayloads struct {
	CompanyId         int     `json:"company_id"`
	PeriodYear        string  `json:"period_year"`
	PeriodMonth       string  `json:"period_month"`
	ItemId            int     `json:"item_id"`
	OrderCycle        float64 `json:"order_cycle"`
	QuantityOnOrder   float64 `json:"quantity_on_order"`
	QuantityBackOrder float64 `json:"quantity_back_order"`
}
