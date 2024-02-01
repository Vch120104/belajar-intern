package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type DiscountPercentRepository interface {
	GetAllDiscountPercent(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error)
	SaveDiscountPercent(tx *gorm.DB, request masteritempayloads.DiscountPercentResponse) (bool, error)
	GetDiscountPercentById(tx *gorm.DB, Id int) (masteritempayloads.DiscountPercentResponse, error)
	ChangeStatusDiscountPercent(tx *gorm.DB, Id int) (bool, error)
}
