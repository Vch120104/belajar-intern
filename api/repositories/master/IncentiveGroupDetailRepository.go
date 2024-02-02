package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type IncentiveGroupDetailRepository interface {
	GetAllIncentiveGroupDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination)
	GetAlIncentiveGroupDetailIsActive() ([]masterpayloads.IncentiveGroupDetailResponse, error)
	GetIncentiveGroupDetailById(Id int) (masterpayloads.IncentiveGroupDetailResponse, error)
	GetIncentiveGroupDetailByCode(Code string) (masterpayloads.IncentiveGroupDetailResponse, error)
	SaveIncentiveGroupDetail(req masterpayloads.IncentiveGroupDetailResponse) (bool, error)
	ChangeStatusIncentiveGroupDetail(Id int) (bool, error)
}
