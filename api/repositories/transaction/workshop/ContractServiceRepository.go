package transactionworkshoprepository

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ContractServiceRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (transactionworkshoppayloads.ContractServiceResponseId, *exceptions.BaseErrorResponse)
	Save(tx *gorm.DB, payload transactionworkshoppayloads.ContractServiceInsert) (transactionworkshoppayloads.ContractServiceInsert, *exceptions.BaseErrorResponse)
	Void(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
	Submit(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
}
