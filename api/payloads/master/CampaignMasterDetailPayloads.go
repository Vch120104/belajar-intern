package masterpayloads

type CampaignMasterDetailPayloads struct {
	CampaignDetailId   int     `json:"campaign_detail_id"`
	CampaignId         int     `json:"campaign_id"`
	LineTypeId         int     `json:"line_type_id"`
	Quantity           float64 `json:"quantity"`
	OperationItemCode  string  `json:"operation_item_code"`
	OperationItemPrice float64 `json:"operation_item_price"`
	DiscountPercent    float64 `json:"discount_percent"`
	SharePercent       float64 `json:"share_percent"`
	ShareBillTo        string  `json:"share_bill_to"`
}

type CampaignMasterDetailSearchPayloads struct {
	IsActive    bool   `json:"is_active"`
	PackageCode string `json:"package_code"`
	LineTypeId  int    `json:"line_type_id"`
}
