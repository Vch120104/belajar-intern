package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type IncentiveGroupDetailService interface {
	GetAllIncentiveGroupDetail([]utils.FilterCondition, pagination.Pagination) pagination.Pagination
	GetIncentiveGroupDetailById(int) masterpayloads.IncentiveGroupDetailResponse
	SaveIncentiveGroupDetail(int, masterpayloads.IncentiveGroupDetailResponse) bool
}
