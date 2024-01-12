package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DiscountPercentRepository interface {
	WithTrx(trxHandle *gorm.DB) DiscountPercentRepository
	GetAllDiscountPercent(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error)
	SaveDiscountPercent(request masteritempayloads.DiscountPercentResponse) (bool, error)
	GetDiscountPercentById(Id int) (masteritempayloads.DiscountPercentResponse, error)
	ChangeStatusDiscountPercent(Id int) (bool, error)
}
