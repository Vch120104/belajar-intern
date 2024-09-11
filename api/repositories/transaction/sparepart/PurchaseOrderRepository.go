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
	GetAllDetailByHeaderId(*gorm.DB, int, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	NewPurchaseOrderHeader(*gorm.DB, transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderResponses) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse)
	UpdatePurchaseOrderHeader(*gorm.DB, int, transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderPayloads) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse)
	GetPurchaseOrderDetailById(*gorm.DB, int) (transactionsparepartpayloads.PurchaseOrderGetDetail, *exceptions.BaseErrorResponse)
	NewPurchaseOrderDetail(*gorm.DB, transactionsparepartpayloads.PurchaseOrderDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse)
	DeletePurchaseOrderDetailMultiId(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	SavePurchaseOrderDetail(*gorm.DB, transactionsparepartpayloads.PurchaseOrderSaveDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse)
	DeleteDocument(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)

	//-=-=-=-=-=
	GetFromPurchaseRequest(db *gorm.DB, filter []utils.FilterCondition, pagination pagination.Pagination)
	SubmitPurchaseOrderRequest(db *gorm.DB, payloads transactionsparepartpayloads.PurchaseOrderHeaderSubmitRequest) (bool, *exceptions.BaseErrorResponse)
	CloseOrderPurchaseOrder(db *gorm.DB, payloads transactionsparepartpayloads.PurchaseOrderCloseOrderPayloads) (bool, *exceptions.BaseErrorResponse)
}
