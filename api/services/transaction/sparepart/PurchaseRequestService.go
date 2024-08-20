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
	GetByIdPurchaseRequest(id int) (transactionsparepartpayloads.PurchaseRequestGetByIdResponses, *exceptions.BaseErrorResponse)
	GetAllPurchaseRequestDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdPurchaseRequestDetail(id int) (transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads, *exceptions.BaseErrorResponse)
	NewPurchaseRequestHeader(transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse)
	NewPurchaseRequestDetail(payloads transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartentities.PurchaseRequestDetail, *exceptions.BaseErrorResponse)
	SavePurchaseRequestUpdateHeader(transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, int) (transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, *exceptions.BaseErrorResponse)
	SavePurchaseRequestUpdateDetail(transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse)
	VoidPurchaseRequest(int) (bool, *exceptions.BaseErrorResponse)
	InsertPurchaseRequestUpdateHeader(transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, int) (transactionsparepartpayloads.PurchaseRequestGetByIdResponses, *exceptions.BaseErrorResponse)
	InsertPurchaseRequestUpdateDetail(transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse)
	GetAllItemTypePurchaseRequest(filterCondition []utils.FilterCondition, pages pagination.Pagination, companyId int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemTypePurchaseRequest(companyId int, id int) (transactionsparepartpayloads.PurchaseRequestItemGetAll, *exceptions.BaseErrorResponse)
	GetByCodeItemTypePurchaseRequest(companyId int, stingcode string) (transactionsparepartpayloads.PurchaseRequestItemGetAll, *exceptions.BaseErrorResponse)
	VoidPurchaseRequestDetail(string2 string) (bool, *exceptions.BaseErrorResponse)
}
