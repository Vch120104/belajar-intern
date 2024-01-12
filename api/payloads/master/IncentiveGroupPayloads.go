package masterpayloads

import "time"

type IncentiveGroupResponse struct {
	IncentiveGroupId   int32     `json:"incentive_group_id"`
	IncentiveGroupCode string    `json:"incentive_group_code"`
	IncentiveGroupName string    `json:"incentive_group_name"`
	EffectiveDate      time.Time `json:"effective_date"`
}
