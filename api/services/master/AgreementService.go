package masterservice

import (
	masterentities "after-sales/api/entities/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type AgreementService interface {
	GetAgreementById(int) (masterpayloads.AgreementRequest, *exceptionsss_test.BaseErrorResponse)
	SaveAgreement(masterpayloads.AgreementResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusAgreement(int) (masterentities.Agreement, *exceptionsss_test.BaseErrorResponse)
	GetAllAgreement(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	AddDiscountGroup(int, masterpayloads.DiscountGroupRequest) *exceptionsss_test.BaseErrorResponse
	DeleteDiscountGroup(int, int) *exceptionsss_test.BaseErrorResponse
	AddItemDiscount(int, masterpayloads.ItemDiscountRequest) *exceptionsss_test.BaseErrorResponse
	DeleteItemDiscount(int, int) *exceptionsss_test.BaseErrorResponse
	AddDiscountValue(int, masterpayloads.DiscountValueRequest) *exceptionsss_test.BaseErrorResponse
	DeleteDiscountValue(int, int) *exceptionsss_test.BaseErrorResponse
	GetAllDiscountGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetAllItemDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetAllDiscountValue(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetDiscountGroupAgreementById(int, int) (masterpayloads.DiscountGroupRequest, *exceptionsss_test.BaseErrorResponse)
	GetDiscountItemAgreementById(int, int) (masterpayloads.ItemDiscountRequest, *exceptionsss_test.BaseErrorResponse)
	GetDiscountValueAgreementById(int, int) (masterpayloads.DiscountValueRequest, *exceptionsss_test.BaseErrorResponse)
}
