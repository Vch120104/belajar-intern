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
	SaveSkillLevel(tx *gorm.DB, req masterpayloads.SkillLevelResponse) (masterentities.SkillLevel, *exceptions.BaseErrorResponse)
	ChangeStatusSkillLevel(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
}
