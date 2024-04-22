package masterentities

var CreateIncentiveMasterTable = "mtr_aftersales_incentive"

type IncentiveMaster struct {
	IncentiveLevelId      int     `gorm:"column:incentive_level_id;size:30;not null;primaryKey"        json:"incentive_level_id"`
	IncentiveLevelCode    int     `gorm:"column:incentive_level_code;size:30;not null"        json:"incentive_level_code"`
	JobPositionId         int     `gorm:"column:job_position_id;size:30;not null"        json:"job_position_id"` //fk with job_position_id in general service
	IncentiveLevelPercent float64 `gorm:"column:incentive_level_percent;not null"        json:"incentive_level_percent"`
	IsActive              bool    `gorm:"column:is_active;not null;default:true"        json:"is_active"`
}

func (*IncentiveMaster) TableName() string {
	return CreateIncentiveMasterTable
}
