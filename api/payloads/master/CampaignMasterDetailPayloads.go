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

type CampaignMasterDetailGetPayloads struct {
	IsActive         bool    `json:"is_active"`
	CampaignId       int     `json:"campaign_id"`
	CampaignDetailId int     `json:"campaign_detail_id"`
	PackageId        int     `json:"package_id"`
	LineTypeId       int     `json:"line_type_id"`
	ItemOperationId  int     `json:"item_operation_id"`
	Quantity         float64 `json:"quantity"`
	Price            float64 `json:"price"`
	DiscountPercent  float64 `json:"discount_percent"`
	SharePercent     float64 `json:"share_percent"`
	ShareBillTo      string  `json:"share_bill_to"`
}

type CampaignMasterDetailPostFromPackageRequest struct {
	CampaignId int `json:"campaign_id"`
	CompanyId  int `json:"company_id"`
	BrandId    int `json:"brand_id"`
	ModelId    int `json:"model_id"`
	PackageId  int `json:"package_id"`
}

type TaxFarePercentResponse struct {
	TaxPercent float64 `json:"tax_percent"`
}
