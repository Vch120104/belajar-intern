package masterpayloads

type CampaignMasterDetailPayloads struct {
	CampaignId      int     `json:"campaign_id"`
	LineTypeId      int     `json:"line_type_id"`
	Quantity        float64 `json:"quantity"`
	OperationItemId int     `json:"operation_item_id"`
	DiscountPercent float64 `json:"discount_percent"`
	SharePercent    float64 `json:"share_percent"`
	ShareBillTo     string  `json:"share_bill_to"`
}

type CampaignMasterDetailSearchPayloads struct {
	IsActive    bool   `json:"is_active"`
	PackageCode string `json:"package_code"`
	LineTypeId  int    `json:"line_type_id"`
}

type CampaignMasterDetailOperationPayloads struct {
	IsActive        bool    `json:"is_active"`
	PackageCode     int     `json:"package_code"`
	PackageId       int     `json:"package_id"`
	LineTypeId      int     `json:"line_type_id"`
	OperationId     int     `json:"operation_id"`
	OperationCode   string  `json:"operation_code"`
	OperationName   string  `json:"operation_name"`
	Quantity        float64 `json:"quantity"`
	Price           float64 `json:"price"`
	DiscountPercent float64 `json:"discount_percent"`
	SharePercent    float64 `json:"share_percent"`
}

type CampaignMasterDetailItemPayloads struct {
	IsActive        bool    `json:"is_active"`
	PackageCode     int     `json:"package_code"`
	PackageId       int     `json:"package_id"`
	LineTypeId      int     `json:"line_type_id"`
	ItemOperationId int     `json:"item_operation_id"`
	ItemCode        string  `json:"item_code"`
	ItemName        string  `json:"item_name"`
	Quantity        float64 `json:"quantity"`
	Price           float64 `json:"price"`
	DiscountPercent float64 `json:"discount_percent"`
	SharePercent    float64 `json:"share_percent"`
}
