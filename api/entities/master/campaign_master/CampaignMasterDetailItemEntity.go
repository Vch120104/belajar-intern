package mastercampaignmasterentities

import masteritementities "after-sales/api/entities/master/item"

var CreateCampaignMasterDetailItemTable = "mtr_campaign_detail_item"

type CampaignMasterDetailItem struct {
	IsActive         bool `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	CampaignDetailId int  `gorm:"column:campaign_detail_id;size:30;primary key;not null; autoincrement:True" json:"campaign_detail_id"`
	CampaignId       int  `gorm:"column:campaign_id;size:30;not null;uniqueIndex:idx_campaign_operationitemcode" json:"campaign_id"`
	Campaign         *CampaignMaster
	LineTypeId       int     `gorm:"column:line_type_id;size:30;not null" json:"line_type_id"`
	Quantity         float64 `gorm:"column:quantity" json:"quantity"`
	ItemId           int     `gorm:"item_id" json:"item_id"`
	Item             *masteritementities.Item
	ShareBillTo      string  `gorm:"column:share_bill_to;size:10" json:"share_bill_to"`
	DiscountPercent  float64 `gorm:"column:discount_percent" json:"discount_percent"`
	SharePercent     float64 `gorm:"column:share_percent" json:"share_percent"`
	Price            float64 `gorm:"price" json:"price"`
}

func (*CampaignMasterDetailItem) TableName() string {
	return CreateCampaignMasterDetailItemTable
}
