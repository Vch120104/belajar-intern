package masterpayloads

type IncentiveMasterResponse struct {
	IncentiveLevelId      int     `json:"incentive_level_id"`
	IsActive              bool    `json:"is_active"`
	IncentiveLevelCode    int     `json:"incentive_level_code"`
	JobPositionId         int     `json:"job_position_id"`
	IncentiveLevelPercent float64 `json:"incentive_level_percent"`
}

type IncentiveMasterListResponse struct {
	IsActive              bool    `json:"is_active" parent_entity:"mtr_aftersales_incentive"`
	IncentiveLevelId      int     `json:"incentive_level_id" parent_entity:"mtr_aftersales_incentive" main_table:"mtr_aftersales_incentive"`
	IncentiveLevelCode    int     `json:"incentive_level_code" parent_entity:"mtr_aftersales_incentive"`
	JobPositionId         int     `json:"job_position_id" parent_entity:"mtr_aftersales_incentive"`
	IncentiveLevelPercent float64 `json:"incentive_level_percent" parent_entity:"mtr_aftersales_incentive"`
}

type IncentiveMasterRequest struct {
	IncentiveLevelId      int     `json:"incentive_level_id"`
	IncentiveLevelCode    int     `json:"incentive_level_code"`
	JobPositionId         int     `json:"job_position_id"`
	IncentiveLevelPercent float64 `json:"incentive_level_percent"`
}

type JobPositionResponse struct {
	JobPositionId   int    `json:"job_position_id"`
	JobPositionName string `json:"job_position_name"`
	JobPositionCode string `json:"job_position_code"`
}
