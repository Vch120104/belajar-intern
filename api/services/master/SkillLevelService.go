package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type SkillLevelService interface {
	GetAllSkillLevel(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetSkillLevelById(Id int) (masterpayloads.SkillLevelResponse, *exceptions.BaseErrorResponse)
	SaveSkillLevel(req masterpayloads.SkillLevelResponse) (masterentities.SkillLevel, *exceptions.BaseErrorResponse)
	ChangeStatusSkillLevel(Id int) (bool, *exceptions.BaseErrorResponse)
}
