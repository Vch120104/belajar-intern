package transactionsparepartpayloads

type ItemInquiryCompanyResponse struct {
	CompanyId   int    `json:"company_id"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
}

type ItemInquiryCompanyReferenceResponse struct {
	CurrencyId int `json:"currency_id"`
}

type ItemInquiryPriceListCodeResponse struct {
	PriceListCodeId int `json:"price_list_code_id"`
}

type ItemInquiryCurrentPeriodResponse struct {
	PeriodYear  string `json:"period_year"`
	PeriodMonth string `json:"period_month"`
}

type ItemInquiryItemGroupResponse struct {
	IsActive      bool   `json:"is_active"`
	ItemGroupId   int    `json:"item_group_id"`
	ItemGroupCode string `json:"item_group_code"`
	ItemGroupName string `json:"item_group_name"`
}

type ItemInquiryCompanyBrandResponse struct {
	CompanyBrandId int    `json:"company_brand_id"`
	BrandId        int    `json:"brand_id"`
	BrandCode      string `json:"brand_code"`
	BrandName      string `json:"brand_name"`
}

type ItemInquiryGetAllPayloads struct {
	ItemDetailId           int     `json:"item_detail_id"`
	ItemId                 int     `json:"item_id"`
	ItemCode               string  `json:"item_code"`
	ItemName               string  `json:"item_name"`
	ItemClassCode          string  `json:"item_class_code"`
	BrandId                int     `json:"brand_id"`
	ModelCode              string  `json:"model_code"`
	WarehouseGroupCode     string  `json:"warehouse_group_code"`
	WarehouseCode          string  `json:"warehouse_code"`
	WarehouseLocationCode  string  `json:"warehouse_location_code"`
	PriceListAmount        float64 `json:"price_list_amount"`
	QuantityAvailable      float64 `json:"quantity_available"`
	ItemSubstitute         string  `json:"item_substitute"`
	MovingCode             string  `json:"moving_code"`
	AvailableInOtherDealer string  `json:"available_in_other_dealer"`
}

type ItemInquiryBrandResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
}
