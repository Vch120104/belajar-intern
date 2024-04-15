package masterentities

import masteritementities "after-sales/api/entities/master/item"

var CreateCampaignMasterDetailTable = "mtr_campaign_Detail"

type CampaignMasterDetail struct {
	IsActive           bool `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	CampaignDetailId   int  `gorm:"column:campaign_detail_id;size:30;primary key;not null; autoincrement:True" json:"campaign_detail_id"`
	CampaignId         int  `gorm:"column:campaign_id;size:30;not null" json:"campaign_id"`
	CampaignMaster     CampaignMaster
	LineTypeId         int                  `gorm:"column:line_type_id;size:30;not null" json:"line_type_id"`
	Quantity           float64              `gorm:"column:quantity" json:"quantity"`
	OperationItemCode  string               `gorm:"column:operation_item_code;size:15;unique" json:"operation_item_code"`
	OperationItemPrice float64              `gorm:"column:operation_item_price" json:"operation_item_price"`
	Description        string               `gorm:"column:description;size:60" json:"description"`
	PackageId          int                  `gorm:"column:package_id;size:30" json:"package_id"`
	ItemPackage        masteritementities.ItemPackage `gorm:"foreign_key:package_id" json:"item_package"`
	ShareBillTo        string               `gorm:"column:share_bill_to;size:10" json:"share_bill_to"`
	DiscountPercent    float64              `gorm:"column:discount_percent" json:"discount_percent"`
	SharePercent       float64              `gorm:"column:share_percent" json:"share_percent"`
}
