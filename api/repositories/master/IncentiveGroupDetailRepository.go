package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type IncentiveGroupDetailRepository interface {
	GetAllIncentiveGroupDetail(*gorm.DB, int, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetIncentiveGroupDetailById(*gorm.DB, int) (masterpayloads.IncentiveGroupDetailResponse, *exceptionsss_test.BaseErrorResponse)
	SaveIncentiveGroupDetail(*gorm.DB, masterpayloads.IncentiveGroupDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
}
