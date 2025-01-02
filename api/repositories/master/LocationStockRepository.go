package masterrepository

import (
	"after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"gorm.io/gorm"
)

type LocationStockRepository interface {
	GetViewLocationStock(db *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	UpdateLocationStock(db *gorm.DB, payloads masterwarehousepayloads.LocationStockUpdatePayloads) (bool, *exceptions.BaseErrorResponse)
	GetAvailableQuantity(db *gorm.DB, payload masterwarehousepayloads.GetAvailableQuantityPayload) (masterwarehousepayloads.GetQuantityAvailablePayload, *exceptions.BaseErrorResponse)
}
