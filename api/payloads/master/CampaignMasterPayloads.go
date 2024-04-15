package masterpayloads

import "time"

type CampaignMasterPost struct {
	IsActive           bool      `json:"is_active"`
	CampaignId         int       `json:"campaign_id"`
	CampaignCode       string    `json:"campaign_code"`
	CampaignName       string    `json:"campaign_name"`
	BrandId            int       `json:"brand_id"`
	ModelId            int       `json:"model_id"`
	CampaignPeriodFrom time.Time `json:"campaign_period_from"`
	CampaignPeriodTo   time.Time `json:"campaign_period_to"`
	Remark             string    `json:"remark"`
	AppointmentOnly    bool      `json:"appointment_only"`
}

type CampaignMasterResponse struct {
	IsActive           bool      `json:"is_active"`
	CampaignCode       string    `json:"campaign_code"`
	CampaignName       string    `json:"campaign_name"`
	BrandId            int       `json:"brand_id"`
	ModelId            int       `json:"model_id"`
	CampaignPeriodFrom time.Time `json:"campaign_period_from"`
	CampaignPeriodTo   time.Time `json:"campaign_period_to"`
	Remark             string    `json:"remark"`
	AppointmentOnly    bool      `json:"appointment_only"`
	Total              float64   `json:"total"`
	TotalVat           float64   `json:"total_vat"`
	TotalAfterVat      float64   `json:"total_after_vat"`
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
	ModelCode string `json:"model_code"`
	ModelName string `json:"model_name"`
}

type GetHistory struct {
	CampaignId   string `gorm:"campaign_id"`
	CampaignCode string `json:"campaign_code"`
	CampaignName string `json:"campaign_name"`
}
