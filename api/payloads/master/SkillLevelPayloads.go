package masterpayloads

type SkillLevelResponse struct {
	IsActive                  bool   `json:"is_active"`
	SkillLevelCodeId          int    `json:"skill_level_id"`
	SkillLevelCodeValue       string `json:"skill_level_code"`
	SkillLevelCodeDescription string `json:"skill_level_description"`
}
