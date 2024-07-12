package masterentities

var CreateSkillLevelTable = "mtr_skill_level"

type SkillLevel struct {
	IsActive              bool   `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	SkillLevelId          int    `gorm:"column:skill_level_id;size:30;not null;primaryKey"        json:"skill_level_id"`
	SkillLevelCode        string `gorm:"column:skill_level_code;size:10;not null"        json:"skill_level_code"`
	SkillLevelDescription string `gorm:"column:skill_level_description;size:50;not null"      json:"skill_level_description"`
}

func (*SkillLevel) TableName() string {
	return CreateSkillLevelTable
}
