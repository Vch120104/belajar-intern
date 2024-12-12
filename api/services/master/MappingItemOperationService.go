package masterservice

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type MappingItemOperationService interface {
	GetAllItemOperation(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemOperation(id int) (masterpayloads.ItemOperationPost, *exceptions.BaseErrorResponse)
	PostItemOperation(req masterpayloads.ItemOperationPost) (masterentities.MappingItemOperation, *exceptions.BaseErrorResponse)
	DeleteItemOperation(id int) (bool, *exceptions.BaseErrorResponse)
	UpdateItemOperation(id int, req masterpayloads.ItemOperationPost) (masterentities.MappingItemOperation, *exceptions.BaseErrorResponse)
}
