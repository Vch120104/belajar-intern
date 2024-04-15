package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"

	"gorm.io/gorm"
)

type PriceListRepository interface {
	GetPriceList(*gorm.DB, masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptionsss_test.BaseErrorResponse)
	SavePriceList(tx *gorm.DB, request masteritempayloads.PriceListResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	GetPriceListById(tx *gorm.DB, Id int) (masteritempayloads.PriceListResponse, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusPriceList(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
