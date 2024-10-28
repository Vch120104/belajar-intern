package masterwarehouserepository

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"after-sales/api/exceptions"
	"gorm.io/gorm"
)

type WarehouseCostingTypeRepository interface {
	GetByCodeWarehouseCostingType(*gorm.DB, string) (masterwarehouseentities.WarehouseCostingType, *exceptions.BaseErrorResponse)
}
