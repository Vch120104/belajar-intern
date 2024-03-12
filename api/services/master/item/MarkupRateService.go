package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type MarkupRateService interface {
	GetMarkupRateById(id int) masteritempayloads.MarkupRateResponse
	SaveMarkupRate(req masteritempayloads.MarkupRateRequest) bool
	GetAllMarkupRate(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int)
	ChangeStatusMarkupRate(Id int) bool
}
