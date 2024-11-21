package masterpayloads

type DiscountResponse struct {
	IsActive                bool   `json:"is_active"`
	DiscountCodeId          int    `json:"discount_code_id"`
	DiscountCode            string `json:"discount_code"`
	DiscountDescription     string `json:"discount_description"`
	DiscountCodeDescription string `json:"discount_code_description"`
}

type DiscountUpdate struct {
	DiscountCodeDescription string `json:"discount_code_description"`
}
