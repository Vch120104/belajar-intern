package masterrepository

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemOperationRepository interface {
	GetAllItemOperation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemOperation(tx *gorm.DB, id int) (masterpayloads.ItemOperationPost, *exceptions.BaseErrorResponse)
	PostItemOperation(tx *gorm.DB, req masterpayloads.ItemOperationPost) (masterentities.MappingItemOperation, *exceptions.BaseErrorResponse)
	DeleteItemOperation(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
	UpdateItemOperation(tx *gorm.DB, id int, req masterpayloads.ItemOperationPost) (masterentities.MappingItemOperation, *exceptions.BaseErrorResponse)
}
