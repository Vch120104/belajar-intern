package masterpayloads

type SkillLevelResponse struct {
	IsActive              bool   `json:"is_active"`
	SkillLevelId          int    `json:"skill_level_id"`
	SkillLevelCode        string `json:"skill_level_code"`
	SkillLevelDescription string `json:"skill_level_description"`
}

type SkillLevelRequest struct {
	SkillLevelId          int    `json:"skill_level_id"`
	SkillLevelCode        string `json:"skill_level_code"`
	SkillLevelDescription string `json:"skill_level_description"`
}
