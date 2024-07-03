package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type PriceListService interface {
	GetPriceList(request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, *exceptions.BaseErrorResponse)
	GetPriceListById(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	SavePriceList(request masteritempayloads.PriceListResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusPriceList(Id int) (bool, *exceptions.BaseErrorResponse)
	GetAllPriceListNew (filterCondition []utils.FilterCondition, pages pagination.Pagination)([]map[string]interface{},int,int,*exceptions.BaseErrorResponse)
	DeactivatePriceList(id string)(bool,*exceptions.BaseErrorResponse)
	ActivatePriceList(id string)(bool,*exceptions.BaseErrorResponse)
	DeletePriceList(id string)(bool,*exceptions.BaseErrorResponse)
}
