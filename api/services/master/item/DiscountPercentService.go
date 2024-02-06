package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type DiscountPercentService interface {
	GetAllDiscountPercent(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int)
	SaveDiscountPercent(req masteritempayloads.DiscountPercentResponse) bool
	GetDiscountPercentById(Id int) masteritempayloads.DiscountPercentResponse
	ChangeStatusDiscountPercent(Id int) (bool)
}
