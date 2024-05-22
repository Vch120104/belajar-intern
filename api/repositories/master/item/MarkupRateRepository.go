package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupRateRepository interface {
	GetMarkupRateById(tx *gorm.DB, Id int) (masteritempayloads.MarkupRateResponse, *exceptions.BaseErrorResponse)
	SaveMarkupRate(tx *gorm.DB, request masteritempayloads.MarkupRateRequest) (bool, *exceptions.BaseErrorResponse)
	GetAllMarkupRate(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ChangeStatusMarkupRate(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
}
