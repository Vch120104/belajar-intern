package masterpayloads

type DeductionDetailResponse struct {
	IsActive             bool    `json:"is_active"`
	DeductionDetailId    int     `json:"deduction_detail_id"`
	DeductionDetailCode  string  `json:"deduction_detail_code"`
	DeductionListId      int     `json:"deduction_list_id"`
	DeductionDetailLevel int     `json:"deduction_detail_level"`
	DeductionPercent     float64 `json:"deduction_percent"`
}

type DeductionDetailPostResponse struct {
	DeductionDetailCode  string  `json:"deduction_detail_code"`
	DeductionListId      int     `json:"deduction_list_id"`
	DeductionDetailLevel int     `json:"deduction_detail_level"`
	DeductionPercent     float64 `json:"deduction_percent"`
}
