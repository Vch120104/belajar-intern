package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
)

type IncentiveGroupDetailService interface {
	GetAllIncentiveGroupDetail(int, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetIncentiveGroupDetailById(int) (masterpayloads.IncentiveGroupDetailResponse, *exceptionsss_test.BaseErrorResponse)
	SaveIncentiveGroupDetail(masterpayloads.IncentiveGroupDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	UpdateIncentiveGroupDetail(masterpayloads.UpdateIncentiveGroupDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
}
