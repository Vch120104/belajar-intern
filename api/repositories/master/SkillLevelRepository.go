package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type SkillLevelRepository interface {
	GetAllSkillLevel(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetSkillLevelById(*gorm.DB, int) (masterpayloads.SkillLevelResponse, *exceptions.BaseErrorResponse)
	GetSkillLevelByCode(tx *gorm.DB, Code string) (masterpayloads.SkillLevelResponse, *exceptions.BaseErrorResponse)
	SaveSkillLevel(tx *gorm.DB, req masterpayloads.SkillLevelResponse) (masterentities.SkillLevel, *exceptions.BaseErrorResponse)
	ChangeStatusSkillLevel(tx *gorm.DB, Id int) (masterpayloads.SkillLevelPatchResponse, *exceptions.BaseErrorResponse)
	UpdateSkillLevel(tx *gorm.DB, req masterpayloads.SkillLevelResponse, id int) (masterentities.SkillLevel, *exceptions.BaseErrorResponse)
}
