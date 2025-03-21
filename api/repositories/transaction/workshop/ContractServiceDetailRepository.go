package transactionworkshoprepository

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ContractServiceDetailRepository interface {
	GetAllDetail(tx *gorm.DB, Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.ContractServiceIdResponse, *exceptions.BaseErrorResponse)
	SaveDetail(tx *gorm.DB, req transactionworkshoppayloads.ContractServiceIdResponse) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse)
	UpdateDetail(tx *gorm.DB, contractServiceSystemNumber int, contractServiceLine string, req transactionworkshoppayloads.ContractServiceDetailRequest) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse)
	DeleteDetail(tx *gorm.DB, contractServiceSystemNumber int, packageCode string) (bool, *exceptions.BaseErrorResponse)
}
