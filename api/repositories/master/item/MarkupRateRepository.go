package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupRateRepository interface {
	GetMarkupRateById(tx *gorm.DB, Id int) (masteritempayloads.MarkupRateResponse, *exceptionsss_test.BaseErrorResponse)
	SaveMarkupRate(tx *gorm.DB, request masteritempayloads.MarkupRateRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetAllMarkupRate(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusMarkupRate(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetMarkupRateByMarkupMasterAndOrderType(tx *gorm.DB, MarkupMasterId int, OrderTypeId int) ([]masteritempayloads.MarkupRateResponse, *exceptionsss_test.BaseErrorResponse)
}
