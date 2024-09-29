package masterpayloads

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
