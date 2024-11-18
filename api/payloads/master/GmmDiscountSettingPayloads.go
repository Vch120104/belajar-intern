package masterpayloads

type GmmDiscountSettingResponse struct {
	IsActive             bool    `json:"is_active"`
	GmmDiscountSettingId int     `json:"gmm_discount_setting_id"`
	GmmPriceCodeId       int     `json:"gmm_price_code_id"`
	ItemLevel_1Id        int     `json:"item_level_1_id"`
	OrderTypeId          int     `json:"order_type_id"`
	DiscountPercentage   float64 `json:"discount_percentage"`
}
