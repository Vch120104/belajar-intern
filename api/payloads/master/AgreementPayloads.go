package masterpayloads

import "time"

type AgreementRequest struct {
	AgreementId       int       `json:"agreement_id" parent_entity:"mtr_agreement" main_table:"mtr_agreement"`
	AgreementCode     string    `json:"agreement_code" parent_entity:"mtr_agreement"`
	IsActive          bool      `json:"is_active" parent_entity:"mtr_agreement"`
	BrandId           int       `json:"brand_id" parent_entity:"mtr_agreement"`
	CustomerId        int       `json:"customer_id" parent_entity:"mtr_agreement"`
	ProfitCenterId    int       `json:"profit_center_id"  parent_entity:"mtr_agreement"`
	AgreementDateFrom time.Time `json:"agreement_date_from" parent_entity:"mtr_agreement"`
	AgreementDateTo   time.Time `json:"agreement_date_to" parent_entity:"mtr_agreement"`
	DealerId          int       `json:"company_id" parent_entity:"mtr_agreement"`
	TopId             int       `json:"top_id" parent_entity:"mtr_agreement"`
	AgreementRemark   string    `json:"agreement_remark" parent_entity:"mtr_agreement"`
}

type AgreementResponse struct {
	AgreementId       int       `json:"agreement_id"`
	AgreementCode     string    `json:"agreement_code"`
	IsActive          bool      `json:"is_active"`
	BrandId           int       `json:"brand_id"`
	CustomerId        int       `json:"customer_id"`
	CustomerCode      string    `json:"customer_code"`
	CustomerName      string    `json:"customer_name"`
	CustomerType      string    `json:"customer_type"`
	ProfitCenterId    int       `json:"profit_center_id"`
	AgreementDateFrom time.Time `json:"agreement_date_from"`
	AgreementDateTo   time.Time `json:"agreement_date_to"`
	DealerId          int       `json:"company_id"`
	DealerName        string    `json:"company_name"`
	DealerCode        string    `json:"company_code"`
	TopId             int       `json:"top_id"`
	AgreementRemark   string    `json:"agreement_remark"`
}

type AgreementCustomerResponse struct {
	CustomerId   int    `json:"customer_id"`
	CustomerCode string `json:"customer_code"`
	CustomerName string `json:"customer_name"`
	CustomerType string `json:"customer_type"`
}

type AgreementCompanyResponse struct {
	CompanyId   int    `json:"company_id"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
	CompanyType string `json:"company_type"`
}
