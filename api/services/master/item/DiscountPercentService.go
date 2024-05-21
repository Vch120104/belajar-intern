package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type DiscountPercentService interface {
	GetAllDiscountPercent(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	SaveDiscountPercent(req masteritempayloads.DiscountPercentResponse) (bool, *exceptions.BaseErrorResponse)
	GetDiscountPercentById(Id int) (masteritempayloads.DiscountPercentResponse, *exceptions.BaseErrorResponse)
	ChangeStatusDiscountPercent(Id int) (bool, *exceptions.BaseErrorResponse)
}
