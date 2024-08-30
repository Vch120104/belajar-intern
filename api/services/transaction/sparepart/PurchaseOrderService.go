package transactionsparepartservice

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
)

type PurchaseOrderService interface {
	//GetSupplySlipById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipResponse, *exceptions.BaseErrorResponse)
	GetAllPurchaseOrder(filter []utils.FilterCondition, page pagination.Pagination, DateParams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdPurchaseOrder(int) (transactionsparepartpayloads.PurchaseOrderGetByIdResponses, *exceptions.BaseErrorResponse)
	GetByIdPurchaseOrderDetail(id int, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	NewPurchaseOrderHeader(responses transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderResponses) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse)
	UpdatePurchaseOrderHeader(int, transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderPayloads) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse)
	GetPurchaseOrderDetailById(int) (transactionsparepartpayloads.PurchaseOrderGetDetail, *exceptions.BaseErrorResponse)
	NewPurchaseOrderDetail(payloads transactionsparepartpayloads.PurchaseOrderDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse)
	DeletePurchaseOrderDetailMultiId(string) (bool, *exceptions.BaseErrorResponse)
	SavePurchaseOrderDetail(transactionsparepartpayloads.PurchaseOrderSaveDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse)
}

//NewPurchaseOrderDetail(*gorm.DB, transactionsparepartpayloads.PurchaseOrderDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse)
//DeletePurchaseOrderDetailMultiId(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
//SavePurchaseOrderDetail(*gorm.DB, transactionsparepartpayloads.PurchaseOrderSaveDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse)
