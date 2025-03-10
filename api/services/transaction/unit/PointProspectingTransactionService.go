package transactionunitservice

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionunitpayloads "after-sales/api/payloads/transaction/unit"
	"after-sales/api/utils"
)

type PointProspectingTransactionService interface {
	GetAllCompanyData([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllSalesRepresentative([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetSalesByCompanyCode(float64, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	Process(transactionunitpayloads.ProcessRequest) (bool, *exceptions.BaseErrorResponse)
}
