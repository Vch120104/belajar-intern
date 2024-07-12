package masterpayloads

type DeductionDetailResponse struct {
	IsActive             bool    `json:"is_active"`
	DeductionDetailId    int     `json:"deduction_detail_id"`
	DeductionId          int     `json:"deduction_id"`
	LimitDays            int     `json:"limit_days"`
	DeductionDetailLevel int     `json:"deduction_detail_level"`
	DeductionPercent     float64 `json:"deduction_percent"`
}

type DeductionDetailPostResponse struct {
	DeductionId          int     `json:"deduction_id"`
	DeductionDetailLevel int     `json:"deduction_detail_level"`
	LimitDays            int     `json:"limit_days"`
	DeductionPercent     float64 `json:"deduction_percent"`
}

type DeductionDetailUpdate struct {
	DeductionPercent float64 `json:"deduction_percent"`
	LimitDays        int     `json:"limit_days"`
}
