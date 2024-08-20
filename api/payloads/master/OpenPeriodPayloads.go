package masterpayloads

type OpenPeriodPayloadResponse struct {
	PeriodYear        string `json:"period_year"`
	PeriodMonth       string `json:"period_month"`
	CurrentPeriodDate string `json:"current_period_date"`
}
