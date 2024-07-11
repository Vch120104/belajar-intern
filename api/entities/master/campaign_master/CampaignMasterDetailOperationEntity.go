package mastercampaignmasterentities

var CreateCampaignMasterDetailOperationTable = "mtr_campaign_detail_operation"

type CampaignMasterOperationDetail struct {
	IsActive                  bool    `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CampaignDetailOperationId int     `gorm:"column:campaign_detail_operation_id;primaryKey;autoIncrement;size:30" json:"campaign_detail_operation_id"`
	CampaignId                int     `gorm:"column:campaign_id;not null;index:idx_campaign_operationitemcode;size:30" json:"campaign_id"`
	LineTypeId                int     `gorm:"column:line_type_id;not null;size:30" json:"line_type_id"`
	Quantity                  float64 `gorm:"column:quantity" json:"quantity"`
	OperationModelMappingId   int     `gorm:"column:operation_model_mapping_id;index:idx_campaign_operationitemcode;size:30" json:"operation_model_mapping_id"`
	ShareBillTo               string  `gorm:"column:share_bill_to;size:10" json:"share_bill_to"`
	DiscountPercent           float64 `gorm:"column:discount_percent" json:"discount_percent"`
	SharePercent              float64 `gorm:"column:share_percent" json:"share_percent"`
	Price                     float64 `gorm:"column:operation_price" json:"operation_price"`
}

func (*CampaignMasterOperationDetail) TableName() string {
	return "campaign_master_operation_details"
}