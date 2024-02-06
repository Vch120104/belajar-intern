package masterservice

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
)

type IncentiveGroupDetailService interface {
	GetAllIncentiveGroupDetail(int, pagination.Pagination) pagination.Pagination
	GetIncentiveGroupDetailById(int) masterpayloads.IncentiveGroupDetailResponse
	SaveIncentiveGroupDetail(req masterpayloads.IncentiveGroupDetailRequest) bool
}
