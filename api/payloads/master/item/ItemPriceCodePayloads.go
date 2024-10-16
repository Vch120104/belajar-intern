package masteritempayloads

type SaveItemPriceCode struct {
	ItemPriceCodeId   int    `json:"item_price_code_id"`
	IsActive          bool   `json:"is_active"`
	ItemPriceCode     string `json:"item_price_code"`
	ItemPriceCodeName string `json:"item_price_code_name"`
}

type UpdateItemPriceCode struct {
	ItemPriceCode     string `json:"item_price_code"`
	ItemPriceCodeName string `json:"item_price_code_name"`
}
