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
	GetByIdPurchaseRequest(*gorm.DB, int) (transactionsparepartpayloads.PurchaseRequestGetByIdNormalizeResponses, *exceptions.BaseErrorResponse)
	GetAllPurchaseRequestDetail(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdPurchaseRequestDetail(*gorm.DB, int) (transactionsparepartpayloads.PurchaseRequestDetailResponsesPayloads, *exceptions.BaseErrorResponse)
	NewPurchaseRequestHeader(*gorm.DB, transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse)
	NewPurchaseRequestDetail(*gorm.DB, transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartentities.PurchaseRequestDetail, *exceptions.BaseErrorResponse)
	SavePurchaseRequestHeader(*gorm.DB, transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, int) (transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, *exceptions.BaseErrorResponse)
	SavePurchaseRequestDetail(*gorm.DB, transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse)
	VoidPurchaseRequest(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	InsertPurchaseRequestHeader(*gorm.DB, transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, int) (transactionsparepartpayloads.PurchaseRequestGetByIdNormalizeResponses, *exceptions.BaseErrorResponse)
	InsertPurchaseRequestDetail(*gorm.DB, transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, int) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse)
}

//NewPurchaseRequestHeader(*gorm.DB, transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse)
//NewPurchaseRequestDetail(*gorm.DB, transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartentities.PurchaseRequestDetail, *exceptions.BaseErrorResponse)
//SavePurchaseRequestHeader(*gorm.DB, transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest, *exceptions.BaseErrorResponse)
//SavePurchaseRequestDetail(*gorm.DB, transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads) (transactionsparepartpayloads.PurchaseRequestSaveDetailRequestPayloads, *exceptions.BaseErrorResponse)
