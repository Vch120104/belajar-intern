package transactionsparepartservice

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type PurchaseOrderService interface {
	//GetSupplySlipById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipResponse, *exceptions.BaseErrorResponse)
	GetAllPurchaseOrder(filter []utils.FilterCondition, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
