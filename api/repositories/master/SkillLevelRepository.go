package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type SkillLevelRepository interface {
	GetAllSkillLevel(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetSkillLevelById(*gorm.DB, int) (masterpayloads.SkillLevelResponse, *exceptions.BaseErrorResponse)
	SaveSkillLevel(*gorm.DB, masterpayloads.SkillLevelResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusSkillLevel(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
}
