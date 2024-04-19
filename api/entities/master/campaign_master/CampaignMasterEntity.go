package mastercampaignmasterentities

import "time"

var CreateCampaignMasterTable = "mtr_campaign"

type CampaignMaster struct {
	IsActive             bool                 `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	CampaignId           int                  `gorm:"column:campaign_id;size:30;not null;primaryKey"        json:"campaign_id"`
	CampaignCode         string               `gorm:"column:campaign_code;size:20;not null" json:"campaign_code"`
	CampaignName         string               `gorm:"column:campaign_name;size:40;not nill" json:"campaign_name"`
	BrandId              int                  `gorm:"column:brand_id;size:30;not null" json:"brand_id"`
	ModelId              int                  `gorm:"column:model_id;size30;not null" json:"model_id"`
	CampaignPeriodFrom   time.Time            `gorm:"column:campaign_period_from;null" json:"campign_period_from"`
	CampaignPeriodTo     time.Time            `gorm:"column:campaign_period_to;null" json:"campaign_period_to"`
	Remark               string               `gorm:"column:remark;size:512;null" json:"remark"`
	Total                float64              `gorm:"column:total;null" json:"total"`
	TotalVat             float64              `gorm:"column:total_vat;null" json:"total_vat"`
	TotalAfterVat        float64              `gorm:"column:total_after_vat;null" json:"total_after_vat"`
	TaxId                int                  `gorm:"column:tax_id;size:30;not null" json:"tax_id"`
	WarehouseGroupId     int                  `gorm:"column:warehouse_group_id;size:30;not null" json:"warehouse_group_id"`
	CompanyId            int                  `gorm:"column:company_id;size:30;not null" json:"company_id" `
	AppointmentOnly      bool                 `gorm:"column:appointment_only;not null;default:false" json:"appointment_only"`
	CampaignMasterDetail []CampaignMasterDetail `gorm:"foreignkey:CampaignId;references:CampaignId"`
}

func (*CampaignMaster) TableName() string {

	return CreateCampaignMasterTable

}