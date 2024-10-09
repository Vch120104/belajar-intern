package masterpayloads

type DiscountResponse struct {
	IsActive                     bool   `json:"is_active"`
	DiscountCodeId               int    `json:"discount_code_id"`
	DiscountCodeValue            string `json:"discount_code_value"`
	DiscountCodeDescription      string `json:"discount_code_description"`
	DiscountCodeValueDescription string `json:"discount_code_value_description"`
}
