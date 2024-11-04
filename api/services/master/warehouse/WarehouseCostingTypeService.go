package masterwarehouseservice

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"after-sales/api/exceptions"
)

type WarehouseCostingTypeService interface {
	GetByCodeWarehouseCostingType(string) (masterwarehouseentities.WarehouseCostingType, *exceptions.BaseErrorResponse)
}
