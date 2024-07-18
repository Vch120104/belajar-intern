package masterwarehouserepository

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WarehouseLocationDefinitionRepository interface {
	GetByLevel(*gorm.DB, int, string) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptions.BaseErrorResponse)
	SaveData(*gorm.DB, masterwarehousepayloads.WarehouseLocationDefinitionResponse) (masterwarehouseentities.WarehouseLocationDefinition, *exceptions.BaseErrorResponse)
	Save(*gorm.DB, masterwarehousepayloads.WarehouseLocationDefinitionResponse) (masterwarehouseentities.WarehouseLocationDefinition, *exceptions.BaseErrorResponse)
	GetById(*gorm.DB, int) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptions.BaseErrorResponse)
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ChangeStatus(tx *gorm.DB, Id int) (masterwarehouseentities.WarehouseLocationDefinition, *exceptions.BaseErrorResponse)
	PopupWarehouseLocationLevel(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
