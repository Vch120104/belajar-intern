package transactionsparepartrepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"gorm.io/gorm"
)

type PurchaseOrderRepository interface {
	GetAllPurchaseOrder(db *gorm.DB, filter []utils.FilterCondition, pagination pagination.Pagination, DateParams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse)
	//GetPurchaseOrderById(db *gorm.DB, id int) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
