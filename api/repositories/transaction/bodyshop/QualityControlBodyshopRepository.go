package transactionbodyshoprepository

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionbodyshoppayloads "after-sales/api/payloads/transaction/bodyshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type QualityControlBodyshopRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionbodyshoppayloads.QualityControlIdResponse, *exceptions.BaseErrorResponse)
	Qcpass(tx *gorm.DB, id int, iddet int) (transactionbodyshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse)
	Reorder(tx *gorm.DB, id int, iddet int, payload transactionbodyshoppayloads.QualityControlReorder) (transactionbodyshoppayloads.QualityControlUpdateResponse, *exceptions.BaseErrorResponse)
}
