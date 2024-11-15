package masterpayloads

type GmmPriceCodeResponse struct {
	IsActive       bool   `json:"is_active"`
	GmmPriceCodeId int    `json:"gmm_price_code_id"`
	GmmPriceCode   string `json:"gmm_price_code"`
	GmmPriceDesc   string `json:"gmm_price_desc"`
}

type GmmPriceCodeDropdownResponse struct {
	IsActive         bool   `json:"is_active"`
	GmmPriceCodeId   int    `json:"gmm_price_code_id"`
	GmmPriceCode     string `json:"gmm_price_code"`
	GmmPriceDesc     string `json:"gmm_price_desc"`
	GmmPriceCodeDesc string `json:"gmm_price_code_desc"`
}

type GmmPriceCodeSaveRequest struct {
	GmmPriceCode string `json:"gmm_price_code" validate:"required"`
	GmmPriceDesc string `json:"gmm_price_desc"`
}

type GmmPriceCodeUpdateRequest struct {
	GmmPriceCode string `json:"gmm_price_code" validate:"required"`
	GmmPriceDesc string `json:"gmm_price_desc"`
}
