package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
	"gorm.io/gorm"
)

type PurchaseOrderRepository interface {
	GetAllPurchaseOrder(db *gorm.DB, filter []utils.FilterCondition, pagination pagination.Pagination, DateParams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	//GetPurchaseOrderById(db *gorm.DB, id int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdPurchaseOrder(*gorm.DB, int) (transactionsparepartpayloads.PurchaseOrderGetByIdResponses, *exceptions.BaseErrorResponse)
	GetByIdPurchaseOrderDetail(*gorm.DB, int, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	NewPurchaseOrderHeader(*gorm.DB, transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderResponses) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse)
	UpdatePurchaseOrderHeader(*gorm.DB, transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderPayloads) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse)
}
