package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type IncentiveGroupDetailRepository interface {
	GetAllIncentiveGroupDetail(*gorm.DB, int, pagination.Pagination) (pagination.Pagination, error)
	GetIncentiveGroupDetailById(*gorm.DB, int) (masterpayloads.IncentiveGroupDetailResponse, error)
	SaveIncentiveGroupDetail(*gorm.DB, masterpayloads.IncentiveGroupDetailRequest) (bool, error)
}
