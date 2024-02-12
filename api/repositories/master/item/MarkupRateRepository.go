package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupRateRepository interface {
	GetMarkupRateById(tx *gorm.DB, Id int) (masteritempayloads.MarkupRateResponse, error)
	SaveMarkupRate(tx *gorm.DB, request masteritempayloads.MarkupRateRequest) (bool, error)
	GetAllMarkupRate(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error)
	ChangeStatusMarkupRate(tx *gorm.DB, Id int) (bool, error)
}
