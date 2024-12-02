package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type MarkupRateService interface {
	GetMarkupRateById(id int) (masteritempayloads.MarkupRateResponse, *exceptions.BaseErrorResponse)
	SaveMarkupRate(req masteritempayloads.MarkupRateRequest) (bool, *exceptions.BaseErrorResponse)
	GetAllMarkupRate(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ChangeStatusMarkupRate(Id int) (bool, *exceptions.BaseErrorResponse)
	GetMarkupRateByMarkupMasterAndOrderType(MarkupMasterId int, OrderTypeId int) ([]masteritempayloads.MarkupRateResponse, *exceptions.BaseErrorResponse)
}
