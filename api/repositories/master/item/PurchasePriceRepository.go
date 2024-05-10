package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type PurchasePriceRepository interface {
	DeletePurchasePrice(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse
	AddPurchasePrice(*gorm.DB, masteritempayloads.PurchasePriceDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	SavePurchasePrice(*gorm.DB, masteritempayloads.PurchasePriceRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllPurchasePrice(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetPurchasePriceById(tx *gorm.DB, Id int) (masteritempayloads.PurchasePriceRequest, *exceptionsss_test.BaseErrorResponse)
	GetAllPurchasePriceDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusPurchasePrice(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
