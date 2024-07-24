package mastercampaignmasterentities

import "time"

var CreateCampaignMasterTable = "mtr_campaign"

type CampaignMaster struct {
	IsActive                      bool                            `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CampaignId                    int                             `gorm:"column:campaign_id;primaryKey;size:30" json:"campaign_id"`
	CampaignCode                  string                          `gorm:"column:campaign_code;size:20;not null" json:"campaign_code"`
	CampaignName                  string                          `gorm:"column:campaign_name;size:40;not null" json:"campaign_name"`
	BrandId                       int                             `gorm:"column:brand_id;not null;size:30" json:"brand_id"`
	ModelId                       int                             `gorm:"column:model_id;not null;size:30" json:"model_id"`
	CampaignPeriodFrom            time.Time                       `gorm:"column:campaign_period_from" json:"campaign_period_from"`
	CampaignPeriodTo              time.Time                       `gorm:"column:campaign_period_to" json:"campaign_period_to"`
	Remark                        string                          `gorm:"column:remark;size:512" json:"remark"`
	Total                         float64                         `gorm:"column:total" json:"total"`
	TotalVat                      float64                         `gorm:"column:total_vat" json:"total_vat"`
	TotalAfterVat                 float64                         `gorm:"column:total_after_vat" json:"total_after_vat"`
	TaxId                         int                             `gorm:"column:tax_id;not null;size:30" json:"tax_id"`
	WarehouseGroupId              int                             `gorm:"column:warehouse_group_id;not null;size:30" json:"warehouse_group_id"`
	CompanyId                     int                             `gorm:"column:company_id;not null;size:30" json:"company_id"`
	AppointmentOnly               bool                            `gorm:"column:appointment_only;not null;default:false" json:"appointment_only"`
	CampaignMasterDetailItem      []CampaignMasterDetailItem      `gorm:"foreignKey:CampaignId;references:CampaignId"`
	CampaignMasterOperationDetail []CampaignMasterOperationDetail `gorm:"foreignKey:CampaignId;references:CampaignId"`
}

func (*CampaignMaster) TableName() string {
	return CreateCampaignMasterTable
}