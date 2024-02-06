package masterentities

var TableIncentiveGroupDetail = "mtr_incentive_group_detail"

type IncentiveGroupDetail struct {
	IsActive               bool `gorm:"column:is_active;not null;default:true" json:"is_active"`
	IncentiveGroupDetailId int  `gorm:"column:incentive_group_detail_id;not null; primaryKey; size:30" json:"incentive_group_detail_id"`
	IncentiveGroupId       int  `gorm:"column:incentive_group_id;not null; size:30" json:"incentive_group_id"`
	IncentiveGroup         IncentiveGroup
	IncentiveLevel         float64 `gorm:"column:incentive_level;not null" json:"incentive_level"`
	TargetAmount           float64 `gorm:"column:target_amount;not null" json:"target_amount"`
	TargetPercent          float64 `gorm:"column:target_percent;not null" json:"target_percent"`
}

func (*IncentiveGroupDetail) TableName() string {
	return TableIncentiveGroupDetail
}
