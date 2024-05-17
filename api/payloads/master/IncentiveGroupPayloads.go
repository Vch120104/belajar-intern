package masterpayloads

import "time"

type IncentiveGroupResponse struct {
	IsActive           bool      `json:"is_active"`
	IncentiveGroupId   int       `json:"incentive_group_id"`
	IncentiveGroupCode string    `json:"incentive_group_code"`
	IncentiveGroupName string    `json:"incentive_group_name"`
	EffectiveDate      time.Time `json:"effective_date"`
}

type ChangeStatusIncentiveGroupRequest struct {
	IsActive bool `json:"is_active"`
}

type UpdateIncentiveGroupRequest struct {
	IncentiveGroupId   int       `json:"incentive_group_id"`
	IncentiveGroupCode string    `json:"incentive_group_code"`
	IncentiveGroupName string    `json:"incentive_group_name"`
	EffectiveDate      time.Time `json:"effective_date"`
}
