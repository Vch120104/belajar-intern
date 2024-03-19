package masterentities

var CreateIncentiveMasterTable = "mtr_incentive_master"

type IncentiveMaster struct {
	IsActive                   bool    `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	IncentiveMasterId          int     `gorm:"column:incentive_master_id;size:30;not null;primaryKey"        json:"incentive_master_id"`
	IncentiveMasterLevel       int     `gorm:"column:incentive_master_level;size:30;not null"        json:"incentive_master_level"`
	IncentiveMasterValue       string  `gorm:"column:incentive_master_value;size:30;not null"        json:"incentive_master_value"`
	IncentiveMasterDescription string  `gorm:"column:incentive_master_description;size:30;not null"        json:"incentive_master_description"`
	JobPositionId              int     `gorm:"column:job_position_id;size:30;not null"        json:"job_position_id"` //fk with job_position_id in general service
	IncentiveMasterPercent     float64 `gorm:"column:incentive_master_percent;not null"        json:"incentive_master_percent"`
}

func (*IncentiveMaster) TableName() string {
	return CreateIncentiveMasterTable
}
