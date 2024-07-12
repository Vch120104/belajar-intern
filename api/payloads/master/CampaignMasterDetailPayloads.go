package masterpayloads

type CampaignMasterDetailPayloads struct {
	CampaignId       int     `json:"campaign_id"`
	LineTypeId       int     `json:"line_type_id"`
	Quantity         float64 `json:"quantity"`
	OperationItemId  int     `json:"operation_item_id"`
	DiscountPercent  float64 `json:"discount_percent"`
	SharePercent     float64 `json:"share_percent"`
	ShareBillTo      string  `json:"share_bill_to"`
}

type CampaignMasterDetailSearchPayloads struct {
	IsActive    bool   `json:"is_active"`
	PackageCode string `json:"package_code"`
	LineTypeId  int    `json:"line_type_id"`
}

type CampaignMasterDetailOperationPayloads struct {
	IsActive        int     `json:"is_active"`
	PackageCode     int     `json:"package_code"`
	PackageId       int     `json:"package_id"`
	LineTypeId      int     `json:"line_type_id"`
	OperationCode   float64 `json:"operation_code"`
	OperationName   int     `json:"operation_name"`
	Quantity        float64 `json:"quantity"`
	Price           float64 `json:"price"`
	DiscountPercent float64 `json:"discount_percent"`
	SharePercent    float64 `json:"share_percent"`
}

type CampaignMasterDetailItemPayloads struct {
	IsActive        int     `json:"is_active"`
	PackageCode     int     `json:"package_code"`
	PackageId       int     `json:"package_id"`
	LineTypeId      int     `json:"line_type_id"`
	ItemCode        float64 `json:"item_code"`
	ItemName        int     `json:"item_name"`
	Quantity        float64 `json:"quantity"`
	Price           float64 `json:"price"`
	DiscountPercent float64 `json:"discount_percent"`
	SharePercent    float64 `json:"share_percent"`
}
