package transactionsparepartservice

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
	"gorm.io/gorm"
)

// service sales order service
type SalesOrderServiceInterface interface {
	InsertSalesOrderHeader(payload transactionsparepartpayloads.SalesOrderInsertHeaderPayload) (transactionsparepartentities.SalesOrder, *exceptions.BaseErrorResponse)
	GetSalesOrderByID(Id int) (transactionsparepartpayloads.SalesOrderEstimationGetByIdResponse, *exceptions.BaseErrorResponse)
	GetAllSalesOrder(pages pagination.Pagination, condition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	VoidSalesOrder(salesOrderId int) (bool, *exceptions.BaseErrorResponse)
	InsertSalesOrderDetail(payload transactionsparepartpayloads.SalesOrderDetailInsertPayload) (transactionsparepartentities.SalesOrderDetail, *exceptions.BaseErrorResponse)
	DeleteSalesOrderDetail(salesOrderDetailId int) (bool, *exceptions.BaseErrorResponse)
	SalesOrderProposedDiscountMultiId(db *gorm.DB, multiId string, proposedDiscountPercentage float64) (bool, *exceptions.BaseErrorResponse)
}
