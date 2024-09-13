package masterentities

var CreateCampaignMasterDetailTable = "mtr_campaign_master_detail"

type CampaignMasterDetail struct {
	IsActive            bool           `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CampaignDetailId    int            `gorm:"column:campaign_detail_id;primaryKey;autoIncrement;size:30" json:"campaign_detail_id"`
	CampaignId          int            `gorm:"column:campaign_id;not null;index:idx_campaign_operationitemcode;size:30" json:"campaign_id"`
	LineTypeId          int            `gorm:"column:line_type_id;not null;size:30" json:"line_type_id"`
	Quantity            float64        `gorm:"column:quantity" json:"quantity"`
	ItemOperationId     int            `gorm:"column:item_operation_id;index:idx_campaign_operationitemcode;size:30" json:"item_operation_id"`
	ShareBillTo         string         `gorm:"column:share_bill_to;size:10" json:"share_bill_to"`
	DiscountPercent     float64        `gorm:"column:discount_percent" json:"discount_percent"`
	SharePercent        float64        `gorm:"column:share_percent" json:"share_percent"`
	Price               float64        `gorm:"column:price" json:"price"`
	ItemOperationEntity *ItemOperation `gorm:"foreignKey:ItemOperationId;references:ItemOperationId"`
}

func (*CampaignMasterDetail) TableName() string {
	return CreateCampaignMasterDetailTable
}
