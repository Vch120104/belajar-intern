package masterpayloads

import "time"

type ItemOprCodeResponse struct {
	ItemOprCodeId int    `json:"item_opr_code_id"`
	ItemOprCode   string `json:"item_opr_code"`
	ItemOprDesc   string `json:"item_opr_desc"`
}

type CampaignDiscount struct {
	OprItemPrice       float64 `gorm:"column:operation_item_price" json:"operation_item_price"`
	OprItemDiscPercent float64 `gorm:"column:operation_item_discount_percent" json:"operation_item_discount_percent"`
	OprItemDiscAmount  float64 `gorm:"column:operation_item_discount_amount" json:"operation_item_discount_amount"`
	TrxTypeId          int     `gorm:"column:transaction_type_id" json:"transaction_type_id"`
}

type WarehouseMasterForItemLookupResponse struct {
	WarehouseId        int    `json:"warehouse_id"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	WarehouseGroupName string `json:"warehouse_group_name"`
	WarehouseCode      string `json:"warehouse_code"`
	WarehouseName      string `json:"warehouse_name"`
}

type WarehouseGroupByCompanyResponse struct {
	WarehouseGroupId       int    `json:"warehouse_group_id"`
	WarehouseGroupCodeName string `json:"warehouse_group_code_name"`
}

type ItemListTransResponse struct {
	ItemId           int    `json:"item_id"`
	ItemCode         string `json:"item_code"`
	ItemName         string `json:"item_name"`
	ItemClassCode    string `json:"item_class_code"`
	ItemTypeCode     string `json:"item_type"`
	ItemLevel_1_Code string `json:"item_level_1_code"`
	ItemLevel_2_Code string `json:"item_level_2_code"`
	ItemLevel_3_Code string `json:"item_level_3_code"`
	ItemLevel_4_Code string `json:"item_level_4_code"`
}

type ItemListTransPLResponse struct {
	ItemId           int    `json:"item_id"`
	ItemCode         string `json:"item_code"`
	ItemName         string `json:"item_name"`
	ItemClassCode    string `json:"item_class_code"`
	ItemTypeCode     string `json:"item_type"`
	ItemLevel_1_Code string `json:"item_level_1_code"`
	ItemLevel_2_Code string `json:"item_level_2_code"`
	ItemLevel_3_Code string `json:"item_level_3_code"`
	ItemLevel_4_Code string `json:"item_level_4_code"`
}

type GetPriceListCodeResponse struct {
	IsActive          bool   `json:"is_active"`
	PriceListCodeId   int    `json:"price_list_code_id"`
	PriceListCodeName string `json:"price_list_code_name"`
	PriceListCode     string `json:"price_list_code"`
}

type GetCurrentPeriodResponse struct {
	PeriodYear        string    `json:"period_year"`
	PeriodMonth       string    `json:"period_month"`
	CurrentPeriodDate time.Time `json:"current_period_date"`
}

type LocationAvailableResponse struct {
	IsActive              bool   `json:"is_active"`
	WarehouseLocationId   int    `json:"warehouse_location_id"`
	WarehouseLocationCode string `json:"warehouse_location_code"`
	WarehouseLocationName string `json:"warehouse_location_name"`
}

type ItemDetailForItemInquiryResponse struct {
	ModelId     int    `json:"model_id"`
	ModelCode   string `json:"model_code"`
	ModelName   string `json:"model_name"`
	VariantId   int    `json:"variant_id"`
	VariantCode string `json:"variant_code"`
	VariantName string `json:"variant_name"`
}

type ItemDetailForItemInquiryPayload struct {
	ModelId   int `json:"model_id"`
	VariantId int `json:"variant_id"`
}

type ItemSubstituteDetailForItemInquiryResponse struct {
	ItemSubstituteDetailId int     `json:"item_substitute_detail_id"`
	IsActive               bool    `json:"is_active"`
	ItemId                 int     `json:"item_id"`
	ItemName               string  `json:"item_name"`
	Quantity               float64 `json:"quantity"`
	Sequence               int     `json:"sequence"`
}
