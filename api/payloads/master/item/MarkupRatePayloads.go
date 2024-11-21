package masteritempayloads

type MarkupRateResponse struct {
	IsActive       bool    `json:"is_active"`
	MarkupRateId   int     `json:"markup_rate_id"`
	MarkupMasterId int     `json:"markup_master_id"`
	OrderTypeId    int     `json:"order_type_id"`
	MarkupRate     float64 `json:"markup_rate"`
}

type MarkupRateListResponse struct {
	IsActive                bool    `json:"is_active" parent_entity:"mtr_markup_rate"`
	MarkupRateId            int     `json:"markup_rate_id" parent_entity:"mtr_markup_rate" main_table:"mtr_markup_rate"`
	MarkupMasterId          int     `json:"markup_master_id" parent_entity:"mtr_markup_master" references:"mtr_markup_master"`
	MarkupMasterCode        string  `json:"markup_code" parent_entity:"mtr_markup_master"`
	MarkupMasterDescription string  `json:"markup_description" parent_entity:"mtr_markup_master"`
	OrderTypeId             int     `json:"order_type_id" parent_entity:"mtr_markup_rate"`
	MarkupRate              float64 `json:"markup_rate" parent_entity:"mtr_markup_rate"`
}

type MarkupRateRequest struct {
	MarkupRateId   int     `json:"markup_rate_id"`
	MarkupMasterId int     `json:"markup_master_id"`
	OrderTypeId    int     `json:"order_type_id"`
	MarkupRate     float64 `json:"markup_rate"`
}
