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
	PurchaseRequestSaveHeader(*gorm.DB, transactionsparepartpayloads.PurchaseRequestHeaderSaveRequest) (transactionsparepartentities.PurchaseRequestEntities, *exceptions.BaseErrorResponse)
}
