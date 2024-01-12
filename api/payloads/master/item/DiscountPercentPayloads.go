package masteritempayloads

type DiscountPercentResponse struct {
	IsActive          bool    `json:"is_active"`
	DiscountPercentId int     `json:"discount_percent_id"`
	DiscountCodeId    int     `json:"discount_code_id"`
	OrderTypeId       int     `json:"order_type_id"` //FK with mtr_order_type from general service
	Discount          float64 `json:"discount"`
}

type DiscountPercentListResponse struct {
	IsActive                bool    `json:"is_active" parent_entity:"mtr_discount_percent"`
	DiscountPercentId       int     `json:"discount_percent_id" parent_entity:"mtr_discount_percent" main_table:"mtr_discount_percent"`
	DiscountCodeId          int     `json:"discount_code_id" parent_entity:"mtr_discount" references:"mtr_discount"`
	DiscountCodeValue       string  `json:"discount_code_value" parent_entity:"mtr_discount"`
	DiscountCodeDescription string  `json:"discount_code_description" parent_entity:"mtr_discount"`
	OrderTypeId             int     `json:"order_type_id" parent_entity:"mtr_discount_percent"`
	Discount                float64 `json:"discount" parent_entity:"mtr_discount_percent"`
}

type OrderTypeResponse struct {
	OrderTypeId   int    `json:"order_type_id"`
	OrderTypeName string `json:"order_type_name"`
}
