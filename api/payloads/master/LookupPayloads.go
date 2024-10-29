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

type ItemListForPriceList struct {
	ItemId        int    `json:"item_id"`
	ItemCode      string `json:"item_code"`
	ItemName      string `json:"item_name"`
	ItemClassCode string `json:"item_class_code"`
	ItemTypeId    string `json:"item_type_id"`
	ItemLevel1    string `gorm:"column:item_level_1" json:"item_level_1"`
	ItemLevel2    string `gorm:"column:item_level_2" json:"item_level_2"`
	ItemLevel3    string `gorm:"column:item_level_3" json:"item_level_3"`
	ItemLevel4    string `gorm:"column:item_level_4" json:"item_level_4"`
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
