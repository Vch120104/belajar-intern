package transactionsparepartservice

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
)

type PurchaseRequestService interface {
	//tes
	GetAllPurchaseRequest(filterCondition []utils.FilterCondition, pages pagination.Pagination, Dateparams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdPurchaseRequest(id int) (transactionsparepartpayloads.PurchaseRequestGetByIdNormalizeResponses, *exceptions.BaseErrorResponse)
	GetAllPurchaseRequestDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdPurchaseRequestDetail(id int) (transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads, *exceptions.BaseErrorResponse)
	PurchaseRequestSaveNewHeader(transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse)
	PurchaseRequestSaveNewDetail(payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartentities.PurchaseRequestDetail, *exceptions.BaseErrorResponse)
}
