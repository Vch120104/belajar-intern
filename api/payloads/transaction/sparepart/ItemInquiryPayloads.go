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
	WarehouseGroupId       int     `json:"warehouse_group_id"`
	WarehouseId            int     `json:"warehouse_id"`
	WarehouseLocationId    int     `json:"warehouse_location_id"`
	WarehouseGroupCode     string  `json:"warehouse_group_code"`
	WarehouseCode          string  `json:"warehouse_code"`
	WarehouseLocationCode  string  `json:"warehouse_location_code"`
	SalesPrice             float64 `json:"sales_price"`
	QuantityAvailable      float64 `json:"quantity_available"`
	ItemSubstitute         string  `json:"item_substitute"`
	MovingCode             string  `json:"moving_code"`
	AvailableInOtherDealer string  `json:"available_in_other_dealer"`
}

type ItemInquiryGetAllResponse struct {
	ItemDetailId           int                      `json:"item_detail_id"`
	ItemId                 int                      `json:"item_id"`
	ItemCode               string                   `json:"item_code"`
	ItemName               string                   `json:"item_name"`
	ItemClassCode          string                   `json:"item_class_code"`
	BrandId                int                      `json:"brand_id"`
	BrandCode              string                   `json:"brand_code"`
	ModelCode              string                   `json:"model_code"`
	WarehouseGroupId       int                      `json:"warehouse_group_id"`
	WarehouseGroupCode     string                   `json:"warehouse_group_code"`
	WarehouseId            int                      `json:"warehouse_id"`
	WarehouseCode          string                   `json:"warehouse_code"`
	WarehouseLocationId    int                      `json:"warehouse_location_id"`
	WarehouseLocationCode  string                   `json:"warehouse_location_code"`
	SalesPrice             interface{}              `json:"sales_price"`
	QuantityAvailable      interface{}              `json:"quantity_available"`
	ItemSubstitute         string                   `json:"item_substitute"`
	MovingCode             string                   `json:"moving_code"`
	AvailableInOtherDealer string                   `json:"available_in_other_dealer"`
	Tooltip                []map[string]interface{} `json:"tooltip"`
}

type ItemInquiryBrandResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
}

type ItemInquiryModelResponse struct {
	ModelId   int    `json:"model_id"`
	ModelCode string `json:"model_code"`
}

type ItemInquiryToolTip struct {
	BrandId int `json:"brand_id"`
	ModelId int `json:"model_id"`
}

type ItemInquiryGetAllToolTip struct {
	ItemId  int                      `json:"ItemId"`
	Tooltip []map[string]interface{} `json:"tooltip"`
}

type ItemInquiryGetByIdFilter struct {
	ItemId              int
	CompanyId           int
	WarehouseId         int
	WarehouseLocationId int
	BrandId             int
	CurrencyId          int
	CompanySessionId    int
}

type ItemInquiryGetByIdResponse struct {
	CompanyId              int     `json:"company_id"`
	PeriodYear             string  `json:"period_year"`
	PeriodMonth            string  `json:"period_month"`
	ItemId                 int     `json:"item_id"`
	ItemCode               string  `json:"item_code"`
	ItemName               string  `json:"item_name"`
	BrandId                int     `json:"brand_id"`
	WarehouseGroupId       int     `json:"warehouse_group_id"`
	WarehouseGroupCode     string  `json:"warehouse_group_code"`
	WarehouseGroupName     string  `json:"warehouse_group_name"`
	WarehouseId            int     `json:"warehouse_id"`
	WarehouseCode          string  `json:"warehouse_code"`
	WarehouseName          string  `json:"warehouse_name"`
	WarehouseLocationId    int     `json:"warehouse_location_id"`
	WarehouseLocationCode  string  `json:"warehouse_location_code"`
	WarehouseLocationName  string  `json:"warehouse_location_name"`
	PriceListAmount        float64 `json:"price_list_amount"`
	QuantityAvailable      float64 `json:"quantity_available"`
	QuantityBegin          float64 `json:"quantity_begin"`
	QuantitySales          float64 `json:"quantity_sales"`
	QuantitySalesReturn    float64 `json:"quantity_sales_return"`
	QuantityPurchase       float64 `json:"quantity_purchase"`
	QuantityPurchaseReturn float64 `json:"quantity_purchase_return"`
	QuantityTransferIn     float64 `json:"quantity_transfer_in"`
	QuantityTransferOut    float64 `json:"quantity_transfer_out"`
	QuantityInTransit      float64 `json:"quantity_in_transit"`
	QuantityClaimIn        float64 `json:"quantity_claim_in"`
	QuantityClaimOut       float64 `json:"quantity_claim_out"`
	QuantityAdjustment     float64 `json:"quantity_adjustment"`
	QuantityAllocated      float64 `json:"quantity_allocated"`
	QuantityOnHand         float64 `json:"quantity_on_hand"`
	QuantityOnOrder        float64 `json:"quantity_on_order"`
	PriceCurrent           float64 `json:"price_current"`
	MovingCode             string  `json:"moving_code"`
	QuantityBackOrder      float64 `json:"quantity_back_order"`
	QuantityMax            float64 `json:"quantity_max"`
	QuantityMin            float64 `json:"quantity_min"`
	DiscountId             int     `json:"discount_id"`
	DiscountCode           string  `json:"discount_code"`
	DiscountEmergency      string  `json:"discount_emergency"`
	DiscountRegular        string  `json:"discount_regular"`
	ItemClassName          string  `json:"item_class_name"`
	IsTechnicalDefect      bool    `json:"is_technical_defect"`
}
