package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type DiscountService interface {
	GetAllDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllDiscountIsActive() ([]masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse)
	GetDiscountById(Id int) (masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse)
	GetDiscountByCode(Code string) (masterpayloads.DiscountResponse, *exceptionsss_test.BaseErrorResponse)
	SaveDiscount(req masterpayloads.DiscountResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusDiscount(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
