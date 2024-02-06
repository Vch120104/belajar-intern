package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveGroupDetailRepository interface {
	GetAllIncentiveGroupDetail(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetIncentiveGroupDetailById(*gorm.DB, int) (masterpayloads.IncentiveGroupDetailResponse, error)
	SaveIncentiveGroupDetail(*gorm.DB, int, masterpayloads.IncentiveGroupDetailResponse) (bool, error)
}
