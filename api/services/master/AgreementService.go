package masterservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type AgreementService interface {
	GetAgreementById(int) (masterpayloads.AgreementRequest, *exceptionsss_test.BaseErrorResponse)
	SaveAgreement(masterpayloads.AgreementResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusAgreement(int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllAgreement(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	AddDiscountGroup(int, masterpayloads.DiscountGroupRequest) *exceptionsss_test.BaseErrorResponse
	DeleteDiscountGroup(int, int) *exceptionsss_test.BaseErrorResponse
	AddItemDiscount(int, masterpayloads.ItemDiscountRequest) *exceptionsss_test.BaseErrorResponse
	DeleteItemDiscount(int, int) *exceptionsss_test.BaseErrorResponse
	AddDiscountValue(int, masterpayloads.DiscountValueRequest) *exceptionsss_test.BaseErrorResponse
	DeleteDiscountValue(int, int) *exceptionsss_test.BaseErrorResponse
}
