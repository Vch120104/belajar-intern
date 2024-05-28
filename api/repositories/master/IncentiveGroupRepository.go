package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveGroupRepository interface {
	GetAllIncentiveGroup(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllIncentiveGroupIsActive(*gorm.DB) ([]masterpayloads.IncentiveGroupResponse, *exceptions.BaseErrorResponse)
	GetIncentiveGroupById(*gorm.DB, int) (masterpayloads.IncentiveGroupResponse, *exceptions.BaseErrorResponse)
	SaveIncentiveGroup(*gorm.DB, masterpayloads.IncentiveGroupResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusIncentiveGroup(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	UpdateIncentiveGroup(*gorm.DB, int, masterpayloads.UpdateIncentiveGroupRequest) (bool, *exceptions.BaseErrorResponse)
	GetAllIncentiveGroupDropDown(tx *gorm.DB) ([]masterpayloads.IncentiveGroupDropDown, *exceptions.BaseErrorResponse)
}
