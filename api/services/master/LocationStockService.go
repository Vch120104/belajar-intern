package masterservice

import (
	"after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LocationStockService interface {
	GetViewLocationStock([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	UpdateLocationStock(payloads masterwarehousepayloads.LocationStockUpdatePayloads) (bool, *exceptions.BaseErrorResponse)
	GetAvailableQuantity(payload masterwarehousepayloads.GetAvailableQuantityPayload) (masterwarehousepayloads.GetQuantityAvailablePayload, *exceptions.BaseErrorResponse)
}
