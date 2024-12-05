package masteritempayloads

type MarkupRateResponse struct {
	IsActive       bool    `json:"is_active"`
	MarkupRateId   int     `json:"markup_rate_id"`
	MarkupMasterId int     `json:"markup_master_id"`
	OrderTypeId    int     `json:"order_type_id"`
	MarkupRate     float64 `json:"markup_rate"`
}

type MarkupRateListResponse struct {
	IsActive          bool    `json:"is_active"`
	MarkupRateId      int     `json:"markup_rate_id"`
	MarkupMasterId    int     `json:"markup_master_id"`
	OrderTypeId       int     `json:"order_type_id"`
	MarkupRate        float64 `json:"markup_rate"`
	OrderTypeName     string  `json:"order_type_name"`
	MarkupCode        string  `json:"markup_code"`        // Ensure correct field mapping
	MarkupDescription string  `json:"markup_description"` // Ensure correct field mapping
}

type MarkupRateRequest struct {
	MarkupRateId   int     `json:"markup_rate_id"`
	MarkupMasterId int     `json:"markup_master_id"`
	OrderTypeId    int     `json:"order_type_id"`
	MarkupRate     float64 `json:"markup_rate"`
}
