package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DiscountPercentService interface {
	WithTrx(trxHandle *gorm.DB) DiscountPercentService
	GetAllDiscountPercent(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error)
	SaveDiscountPercent(req masteritempayloads.DiscountPercentResponse) (bool, error)
	GetDiscountPercentById(Id int) (masteritempayloads.DiscountPercentResponse, error)
	ChangeStatusDiscountPercent(Id int) (bool, error)
}
