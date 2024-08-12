package masterpayloads

import "time"

type CampaignMasterPost struct {
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
	ModelId            int     `json:"model_id"`
	CampaignPeriodFrom string  `json:"campaign_period_from"`
	CampaignPeriodTo   string  `json:"campaign_period_to"`
	Remark             string  `json:"remark"`
	AppointmentOnly    bool    `json:"appointment_only"`
	Total              float64 `json:"total"`
	TotalVat           float64 `json:"total_vat"`
	TotalAfterVat      float64 `json:"total_after_vat"`
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
	Brandname string `json:"brand_name"`
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
