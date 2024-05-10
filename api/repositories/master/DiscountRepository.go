package masterrepository

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	exceptionsss_test "after-sales/api/expectionsss"

	"gorm.io/gorm"
)

type DiscountRepository interface {
	GetAllDiscount(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllDiscountIsActive(*gorm.DB) ([]masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse)
	GetDiscountById(*gorm.DB, int) (masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse)
	GetDiscountByCode(*gorm.DB, string) (masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse)
	SaveDiscount(*gorm.DB, masterpayloads.DiscountResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusDiscount(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
}
