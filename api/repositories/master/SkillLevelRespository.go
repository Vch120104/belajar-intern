package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type SkillLevelRepository interface {
	GetAllSkilllevel(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	SaveSkillLevel(tx *gorm.DB, request masterpayloads.SkillLevelResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	GetSkillLevelById(tx *gorm.DB, Id int) (masterpayloads.SkillLevelResponse, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusSkillLevel(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
