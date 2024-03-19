package masterpayloads

type IncentiveMasterResponse struct {
	IsActive                   bool    `json:"is_active"`
	IncentiveMasterId          int     `json:"incentive_master_id"`
	IncentiveMasterLevel       int     `json:"incentive_master_level"`
	IncentiveMasterValue       string  `json:"incentive_master_value"`
	IncentiveMasterDescription string  `json:"incentive_master_description"`
	JobPositionId              int     `json:"job_position_id"`
	IncentiveMasterPercent     float64 `json:"incentive_master_percent"`
}

type IncentiveMasterListResponse struct {
	IsActive                   bool    `json:"is_active" parent_entity:"mtr_incentive_master"`
	IncentiveMasterId          int     `json:"incentive_master_id" parent_entity:"mtr_incentive_master" main_table:"mtr_incentive_master"`
	IncentiveMasterValue       string  `json:"incentive_master_value" parent_entity:"mtr_incentive_master"`
	IncentiveMasterLevel       string  `json:"incentive_master_level" parent_entity:"mtr_incentive_master"`
	IncentiveMasterDescription string  `json:"incentive_master_description" parent_entity:"mtr_incentive_master"`
	JobPositionId              int     `json:"job_position_id" parent_entity:"mtr_incentive_master"`
	IncentiveMasterPercent     float64 `json:"incentive_master_percent" parent_entity:"mtr_incentive_master"`
}

type IncentiveMasterRequest struct {
	IncentiveMasterId          int     `json:"incentive_master_id"`
	IncentiveMasterLevel       int     `json:"incentive_master_level"`
	IncentiveMasterValue       string  `json:"incentive_master_value"`
	IncentiveMasterDescription string  `json:"incentive_master_description"`
	JobPositionId              int     `json:"job_position_id"`
	IncentiveMasterPercent     float64 `json:"incentive_master_percent"`
}

type JobPositionResponse struct {
	JobPositionId   int    `json:"job_position_id"`
	JobPositionName string `json:"job_position_name"`
	JobPositionCode string `json:"job_position_code"`
}
