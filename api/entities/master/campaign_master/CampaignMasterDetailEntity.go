package mastercampaignmasterentities

var CreateCampaignMasterDetailTable = "mtr_campaign_detail"

type CampaignMasterDetail struct {
	IsActive           bool    `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	CampaignDetailId   int     `gorm:"column:campaign_detail_id;size:30;primary key;not null; autoincrement:True" json:"campaign_detail_id"`
	CampaignId         int     `gorm:"column:campaign_id;size:30;not null;uniqueIndex:idx_campaign_operationitemcode" json:"campaign_id"`
	LineTypeId         int     `gorm:"column:line_type_id;size:30;not null" json:"line_type_id"`
	Quantity           float64 `gorm:"column:quantity" json:"quantity"`
	OperationItemCode  string  `gorm:"column:operation_item_code;size:15;uniqueIndex:idx_campaign_operationitemcode" json:"operation_item_code"`
	OperationItemPrice float64 `gorm:"column:operation_item_price" json:"operation_item_price"`
	PackageId          int     `gorm:"column:package_id;size:30" json:"package_id"`
	ShareBillTo        string  `gorm:"column:share_bill_to;size:10" json:"share_bill_to"`
	DiscountPercent    float64 `gorm:"column:discount_percent" json:"discount_percent"`
	SharePercent       float64 `gorm:"column:share_percent" json:"share_percent"`
}

func (*CampaignMasterDetail) TableName() string {
	return CreateCampaignMasterDetailTable
}
