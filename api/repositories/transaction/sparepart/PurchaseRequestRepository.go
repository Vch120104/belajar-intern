package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type PurchaseRequestRepository interface {
	GetAllPurchaseRequest(*gorm.DB, []utils.FilterCondition, pagination.Pagination, map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdPurchaseRequest(*gorm.DB, int) (transactionsparepartpayloads.PurchaseRequestGetByIdResponses, *exceptions.BaseErrorResponse)
	GetAllPurchaseRequestDetail(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdPurchaseRequestDetail(*gorm.DB, int) (transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads, *exceptions.BaseErrorResponse)
	NewPurchaseRequestHeader(*gorm.DB, transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse)
	NewPurchaseRequestDetail(*gorm.DB, transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartentities.PurchaseRequestDetail, *exceptions.BaseErrorResponse)
	SavePurchaseRequestHeader(*gorm.DB, transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, int) (transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, *exceptions.BaseErrorResponse)
	SavePurchaseRequestDetail(*gorm.DB, transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse)
	VoidPurchaseRequest(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	SubmitPurchaseRequest(*gorm.DB, transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, int) (transactionsparepartpayloads.PurchaseRequestGetByIdResponses, *exceptions.BaseErrorResponse)
	InsertPurchaseRequestDetail(*gorm.DB, transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse)
	GetAllItemTypePrRequest(*gorm.DB, []utils.FilterCondition, pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdPurchaseRequestItemPr(*gorm.DB, int, int) (transactionsparepartpayloads.PurchaseRequestItemGetAll, *exceptions.BaseErrorResponse)
	GetByCodePurchaseRequestItemPr(*gorm.DB, int, string) (transactionsparepartpayloads.PurchaseRequestItemGetAll, *exceptions.BaseErrorResponse)
	VoidPurchaseRequestDetailMultiId(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	GenerateDocumentNumber(tx *gorm.DB, id int) (string, *exceptions.BaseErrorResponse)
}

//NewPurchaseRequestHeader(*gorm.DB, transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse)
//NewPurchaseRequestDetail(*gorm.DB, transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartentities.PurchaseRequestDetail, *exceptions.BaseErrorResponse)
//SavePurchaseRequestHeader(*gorm.DB, transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, *exceptions.BaseErrorResponse)
//SavePurchaseRequestDetail(*gorm.DB, transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse)
