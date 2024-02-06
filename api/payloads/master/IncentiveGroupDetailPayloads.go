package masterpayloads

type IncentiveGroupDetailResponse struct {
	IsActive               bool    `json:"is_active"`
	IncentiveGroupDetailId int     `json:"incentive_group_detail_id"`
	IncentiveGroupId       int     `json:"incentive_group_id"`
	IncentiveLevel         float64 `json:"incentive_level"`
	TargetAmount           float64 `json:"target_amount"`
	TargetPercent          float64 `json:"target_percent"`
}

type IncentiveGroupDetailRequest struct {
	IsActive               bool    ` json:"is_active"`
	IncentiveGroupDetailId int     `json:"incentive_group_detail_id"`
	IncentiveGroupId       int     `json:"incentive_group_id"`
	IncentiveLevel         float64 `json:"incentive_level"`
	TargetAmount           float64 `json:"target_amount"`
	TargetPercent          float64 `json:"target_percent"`
}
