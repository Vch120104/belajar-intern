package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type IncentiveGroupDetailRepository interface {
	GetAllIncentiveGroupDetail(*gorm.DB, int, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetIncentiveGroupDetailById(*gorm.DB, int) (masterpayloads.IncentiveGroupDetailResponse, *exceptions.BaseErrorResponse)
	SaveIncentiveGroupDetail(*gorm.DB, masterpayloads.IncentiveGroupDetailRequest) (bool, *exceptions.BaseErrorResponse)
	UpdateIncentiveGroupDetail(*gorm.DB, int, masterpayloads.UpdateIncentiveGroupDetailRequest) (bool, *exceptions.BaseErrorResponse)
}
