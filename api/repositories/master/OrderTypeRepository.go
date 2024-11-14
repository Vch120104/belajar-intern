package masterrepository

import (
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"

	"gorm.io/gorm"
)

type OrderTypeRepository interface {
	GetAllOrderType(tx *gorm.DB) ([]masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse)
}
