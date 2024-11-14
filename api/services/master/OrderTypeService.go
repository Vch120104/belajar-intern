package masterservice

import (
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
)

type OrderTypeService interface {
	GetAllOrderType() ([]masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse)
}
