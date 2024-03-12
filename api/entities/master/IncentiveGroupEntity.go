package masterentities

import "time"

var CreateIncentiveGroupTable = "mtr_incentive_group"

type IncentiveGroup struct {
	IsActive           bool      `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	IncentiveGroupId   int     `gorm:"column:incentive_group_id;not null;primaryKey; size:30"        json:"incentive_group_id"`
	IncentiveGroupCode string    `gorm:"column:incentive_group_code;size:50;not null"        json:"incentive_group_code"`
	IncentiveGroupName string    `gorm:"column:incentive_group_name;size:100;not null"        json:"incentive_group_name"`
	EffectiveDate      time.Time `gorm:"column:effective_date;not null"        json:"effective_date"`
}

func (*IncentiveGroup) TableName() string {
	return CreateIncentiveGroupTable
}
