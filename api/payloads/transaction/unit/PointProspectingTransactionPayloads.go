package transactionunitpayloads

type ProcessRequest struct {
	CompanyCode    float64 `json:"company_code"`
	PeriodMonth    string  `json:"period_month"`
	PeriodYear     string  `json:"period_year"`
	CreationUserId string  `json:"creation_user_id"`
}

type PointProspectingResponse struct {
	CompanyCode float64 `gorm:"column:company_code" json:"company_code"`
	CompanyName string  `gorm:"column:company_name" json:"company_name"`
	Status      string  `gorm:"column:record_status" json:"status"`
}

type GetAllSalesRepresentativeResponse struct {
	EmployeeNo   string `json:"sales_code"`
	EmployeeName string `json:"sales_name"`
	CompanyName  string `json:"company_name"`
}
