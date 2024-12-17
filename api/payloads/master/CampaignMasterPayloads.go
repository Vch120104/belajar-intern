package masterpayloads

import "time"

type CampaignMasterPost struct {
	CampaignId         int       `json:"campaign_id"`
	CampaignCode       string    `json:"campaign_code"`
	CampaignName       string    `json:"campaign_name"`
	BrandId            int       `json:"brand_id"`
	ModelId            int       `json:"model_id"`
	CampaignPeriodFrom time.Time `json:"campaign_period_from"`
	CampaignPeriodTo   time.Time `json:"campaign_period_to"`
	Remark             string    `json:"remark"`
	AppointmentOnly    bool      `json:"appointment_only"`
	TaxId              int       `json:"tax_id"`
	CompanyId          int       `json:"company_id"`
	WarehouseGroupId   int       `json:"warehouse_group_id"`
}

type CampaignMasterResponse struct {
	IsActive           bool    `json:"is_active"`
	CampaignCode       string  `json:"campaign_code"`
	CampaignName       string  `json:"campaign_name"`
	CampaignId         int     `json:"campaign_id"`
	BrandId            int     `json:"brand_id"`
	BrandCode          string  `json:"brand_code"`
	BrandName          string  `json:"brand_name"`
	ModelId            int     `json:"model_id"`
	ModelCode          string  `json:"model_code"`
	ModelDescription   string  `json:"model_description"`
	CampaignPeriodFrom string  `json:"campaign_period_from"`
	CampaignPeriodTo   string  `json:"campaign_period_to"`
	Remark             string  `json:"remark"`
	AppointmentOnly    bool    `json:"appointment_only"`
	Total              float64 `json:"total"`
	TotalVat           float64 `json:"total_vat"`
	TotalAfterVat      float64 `json:"total_after_vat"`
	CompanyId          int     `json:"company_id"`
	CompanyCode        string  `json:"company_code"`
	CompanyName        string  `json:"company_name"`
}

type CampaignMasterListReponse struct {
	IsActive           bool      `json:"is_active"`
	CampaignCode       string    `json:"campaign_code"`
	CampaignName       string    `json:"campaign_name"`
	ModelId            int       `json:"model_id"`
	CampaignPeriodFrom time.Time `json:"campaign_period_from"`
	CampaignPeriodTo   time.Time `json:"campaign_period_to"`
}

type GetModelResponse struct {
	ModelId          int    `json:"model_id"`
	ModelCode        string `json:"model_code"`
	ModelDescription string `json:"model_description"`
}
type GetBrandResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
}
type GetHistory struct {
	CampaignId   string `json:"campaign_id"`
	CampaignCode string `json:"campaign_code"`
	CampaignName string `json:"campaign_name"`
}

type CampaignMasterTaxAndTotal struct {
	TaxId float64 `json:"tax_id"`
	Total int     `json:"total"`
}

type CampaignMasterJobTypeResponse struct {
	IsActive    bool   `json:"is_active"`
	JobTypeId   int    `json:"job_type_id"`
	JobTypeCode string `json:"job_type_code"`
	JobTypeName string `json:"job_type_name"`
}

type CampaignMasterCompanyResponse struct {
	CompanyId    int `json:"company_id"`
	CompanyName  int `json:"company_name"`
	VatCompanyId int `json:"vat_company_id"`
}

type CampaignMasterCompanyReferenceResponse struct {
	CurrencyId int `json:"currency_id"`
}

type CampaignMasterWOTransactionResponse struct {
	IsActive                     bool   `json:"is_active"`
	WorkOrderTransactionTypeId   int    `json:"work_order_transaction_type_id"`
	WorkOrderTransactionTypeName string `json:"work_order_transaction_type_name"`
	WorkOrderTransactionTypeCode string `json:"work_order_transaction_type_code"`
}
