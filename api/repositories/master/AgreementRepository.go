package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	exceptionsss_test "after-sales/api/expectionsss"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type AgreementRepository interface {
	GetAgreementById(*gorm.DB, int) (masterpayloads.AgreementRequest, *exceptionsss_test.BaseErrorResponse)
	SaveAgreement(*gorm.DB, masterpayloads.AgreementResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusAgreement(*gorm.DB, int) (masterentities.Agreement, *exceptionsss_test.BaseErrorResponse)
	GetAllAgreement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	AddDiscountGroup(*gorm.DB, int, masterpayloads.DiscountGroupRequest) *exceptionsss_test.BaseErrorResponse
	DeleteDiscountGroup(*gorm.DB, int, int) *exceptionsss_test.BaseErrorResponse
	AddItemDiscount(*gorm.DB, int, masterpayloads.ItemDiscountRequest) *exceptionsss_test.BaseErrorResponse
	DeleteItemDiscount(*gorm.DB, int, int) *exceptionsss_test.BaseErrorResponse
	AddDiscountValue(*gorm.DB, int, masterpayloads.DiscountValueRequest) *exceptionsss_test.BaseErrorResponse
	DeleteDiscountValue(*gorm.DB, int, int) *exceptionsss_test.BaseErrorResponse
	GetDiscountGroupAgreementById(*gorm.DB, int, int) (masterpayloads.DiscountGroupRequest, *exceptionsss_test.BaseErrorResponse)
	GetDiscountItemAgreementById(*gorm.DB, int, int) (masterpayloads.ItemDiscountRequest, *exceptionsss_test.BaseErrorResponse)
	GetDiscountValueAgreementById(*gorm.DB, int, int) (masterpayloads.DiscountValueRequest, *exceptionsss_test.BaseErrorResponse)
	GetAllDiscountGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetAllItemDiscount(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetAllDiscountValue(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
}
