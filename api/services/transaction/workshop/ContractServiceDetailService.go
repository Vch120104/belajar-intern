package transactionworkshopservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type ContractServiceDetailService interface {
	GetAllDetail(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetById(Id int) (transactionworkshoppayloads.ContractServiceIdResponse, *exceptions.BaseErrorResponse)
	SaveDetail(req transactionworkshoppayloads.ContractServiceIdResponse) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse)
	UpdateDetail(contractServiceSystemNumber int, contractServiceLine string, req transactionworkshoppayloads.ContractServiceDetailRequest) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse)
	DeleteDetail(contractServiceSystemNumber int, packageCode string) (bool, *exceptions.BaseErrorResponse)
}
