package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type SkillLevelService interface {
	GetAllSkillLevel(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetSkillLevelById(Id int) (masterpayloads.SkillLevelResponse, *exceptionsss_test.BaseErrorResponse)
	SaveSkillLevel(req masterpayloads.SkillLevelResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusSkillLevel(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
