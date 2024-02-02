package masterentities

var TableIncentiveGroupDetail = "mtr_incentive_group_detail"

type IncentiveGroupDetail struct {
	IncentiveGroupDetailId int     `gorm:"column:incentive_group_detail_id;not null; primaryKey; size:30"        json:"incentive_group_detail_id"`
	IncentiveGroupId       int     `gorm:"column:incentive_group_id;not null; size:30"        json:"incentive_group_id"`
	IncentiveGroupCode     string  `gorm:"column:incentive_group_code;not null"        json:"incentive_group_code"`
	IncentiveLevel         float64 `gorm:"column:incentive_level;not null"        json:"incentive_level"`
	TargetAmount           float64 `gorm:"column:target_amount;not null"        json:"target_amount"`
	TargetPercent          float64 `gorm:"column:target_percent;not null"        json:"target_percent"`
}

func (*IncentiveGroupDetail) TableName() string {
	return TableIncentiveGroupDetail
}
