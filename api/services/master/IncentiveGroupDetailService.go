package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
)

type IncentiveGroupDetailService interface {
	GetAllIncentiveGroupDetail(int, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetIncentiveGroupDetailById(int) (masterpayloads.IncentiveGroupDetailResponse, *exceptions.BaseErrorResponse)
	SaveIncentiveGroupDetail(masterpayloads.IncentiveGroupDetailRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateIncentiveGroupDetail(int, masterpayloads.UpdateIncentiveGroupDetailRequest) (bool, *exceptions.BaseErrorResponse)
}
