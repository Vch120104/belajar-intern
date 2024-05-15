package masterpayloads

type PackageMasterResponse struct {
	IsActive       bool    `json:"is_active"`
	PackageCode    string  `json:"package_code"`
	PackageName    string  `json:"package_name"`
	PackageId      int     `json:"package_id"`
	ProfitCenterId int     `json:"profit_center_id"`
	ItemGroupId    int     `json:"item_group_id"`
	BrandId        int     `json:"brand_id"`
	ModelId        int     `json:"model_id"`
	VariantId      int     `json:"variant_id"`
	PackageSet     bool    `json:"package_set"`
	CurrencyId     int     `json:"currency_id"`
	PackagePrice   float64 `json:"package_price"`
	TaxTypeId      int     `json:"tax_type_id"`
	PackageRemark  string  `json:"package_remark"`
}

type PackageMasterListResponse struct {
	PackageId      int     `json:"package_id" parent_entity:"mtr_package" main_table:"mtr_package"`
	PackageCode    string  `json:"package_code" parent_entity:"mtr_package"`
	PackageName    string  `json:"package_name" parent_entity:"mtr_package"`
	ProfitCenterId int     `json:"profit_center_id" parent_entity:"mtr_package"`
	ModelId        int     `json:"model_id" parent_entity:"mtr_package"`
	VariantId      int     `json:"variant_id" parent_entity:"mtr_package"`
	PackagePrice   float64 `json:"package_price" parent_entity:"mtr_package"`
	IsActive       bool    `json:"is_active" parent_entity:"mtr_package"`
}

type CurrencyResponse struct {
	CurrencyId   int    `json:"currency_id"`
	CurrencyCode string `json:"currency_code"`
}

type GetVariantResponse struct {
	VariantId   int    `json:"variant_id"`
	VariantCode string `json:"variant_code"`
	VariantDesc string `json:"variant_description"`
}

type GetProfitMaster struct {
	ProfitCenterId   int    `json:"profit_center_id"`
	ProfitCenterName string `json:"profit_center_name"`
}

type PackageMasterForCampaignMaster struct {
	PackageCode string `json:"package_code"`
	PackageName string `json:"package_name"`
}
