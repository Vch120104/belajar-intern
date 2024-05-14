package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type SkillLevelRepository interface {
	GetAllSkillLevel(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetSkillLevelById(*gorm.DB, int) (masterpayloads.SkillLevelResponse, *exceptionsss_test.BaseErrorResponse)
	SaveSkillLevel(*gorm.DB, masterpayloads.SkillLevelResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusSkillLevel(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
}
