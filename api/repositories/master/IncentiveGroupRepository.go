package masterrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type IncentiveGroupRepository interface {
	GetAllIncentiveGroup(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllIncentiveGroupIsActive(*gorm.DB) ([]masterpayloads.IncentiveGroupResponse, *exceptionsss_test.BaseErrorResponse)
	GetIncentiveGroupById(*gorm.DB, int) (masterpayloads.IncentiveGroupResponse, *exceptionsss_test.BaseErrorResponse)
	SaveIncentiveGroup(*gorm.DB, masterpayloads.IncentiveGroupResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusIncentiveGroup(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
	UpdateIncentiveGroup(*gorm.DB, masterpayloads.UpdateIncentiveGroupRequest) (bool, *exceptionsss_test.BaseErrorResponse)
}
