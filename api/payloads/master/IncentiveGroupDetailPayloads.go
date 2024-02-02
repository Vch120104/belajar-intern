package masterpayloads

type IncentiveGroupDetailResponse struct {
	IncentiveGroupDetailId int     `json:"incentive_group_detail_id"`
	IncentiveGroupId       int     `json:"incentive_group_id"`
	IncentiveGroupCode     string  `json:"incentive_group_code"`
	IncentiveLevel         float64 `json:"incentive_level"`
	TargetAmount           float64 `json:"target_amount"`
	TargetPercent          float64 `json:"target_percent"`
}
