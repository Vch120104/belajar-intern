package transactionsparepartrepository

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	"after-sales/api/utils"
	"gorm.io/gorm"
)

type SalesOrderRepository interface {
	InsertSalesOrderHeader(db *gorm.DB, payload transactionsparepartpayloads.SalesOrderInsertHeaderPayload) (transactionsparepartentities.SalesOrder, *exceptions.BaseErrorResponse)
	GetSalesOrderByID(tx *gorm.DB, Id int) (transactionsparepartpayloads.SalesOrderEstimationGetByIdResponse, *exceptions.BaseErrorResponse)
	GetAllSalesOrder(db *gorm.DB, pages pagination.Pagination, condition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	VoidSalesOrder(db *gorm.DB, salesOrderId int) (bool, *exceptions.BaseErrorResponse)
	InsertSalesOrderDetail(db *gorm.DB, payload transactionsparepartpayloads.SalesOrderDetailInsertPayload) (transactionsparepartentities.SalesOrderDetail, *exceptions.BaseErrorResponse)
	DeleteSalesOrderDetail(tx *gorm.DB, salesOrderDetailId int) (bool, *exceptions.BaseErrorResponse)
	SalesOrderProposedDiscountMultiId(tx *gorm.DB, multiId string, proposedDiscountPercentage float64) (bool, *exceptions.BaseErrorResponse)
	UpdateSalesOrderHeader(tx *gorm.DB, payload transactionsparepartpayloads.SalesOrderUpdatePayload, SalesOrderId int) (transactionsparepartentities.SalesOrder, *exceptions.BaseErrorResponse)
	SubmitSalesOrderHeader(tx *gorm.DB, salesOrderId int) (bool, *exceptions.BaseErrorResponse)
	GenerateDocumentNumber(tx *gorm.DB, id int) (string, *exceptions.BaseErrorResponse)
	GetSalesOrderTransactionType(tx *gorm.DB) ([]transactionsparepartpayloads.GetSalesOrderTransactionType, *exceptions.BaseErrorResponse)
}
