package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveGroupRepository interface {
	GetAllIncentiveGroup(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetAllIncentiveGroupIsActive(*gorm.DB) ([]masterpayloads.IncentiveGroupResponse, error)
	GetIncentiveGroupById(*gorm.DB, int) (masterpayloads.IncentiveGroupResponse, error)
	SaveIncentiveGroup(*gorm.DB, masterpayloads.IncentiveGroupResponse) (bool, error)
	ChangeStatusIncentiveGroup(*gorm.DB, int) (bool, error)
}
