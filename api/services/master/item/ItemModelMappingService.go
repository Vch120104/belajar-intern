package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
)

type ItemModelMappingService interface {
	CreateItemModelMapping(req masteritempayloads.CreateItemModelMapping) (bool, *exceptions.BaseErrorResponse)
	UpdateItemModelMapping(req masteritempayloads.CreateItemModelMapping) (bool, *exceptions.BaseErrorResponse)
	GetItemModelMappingByItemId(itemId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
