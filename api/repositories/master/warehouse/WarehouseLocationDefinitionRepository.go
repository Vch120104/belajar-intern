package masterwarehouserepository

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptionsss_test "after-sales/api/expectionsss"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WarehouseLocationDefinitionRepository interface {
	GetByLevel(*gorm.DB, int, string) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptionsss_test.BaseErrorResponse)
	SaveData(*gorm.DB, masterwarehousepayloads.WarehouseLocationDefinitionResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	Save(*gorm.DB, masterwarehousepayloads.WarehouseLocationDefinitionResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	GetById(*gorm.DB, int) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptionsss_test.BaseErrorResponse)
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	ChangeStatus(tx *gorm.DB, Id int) (masterwarehouseentities.WarehouseLocationDefinition, *exceptionsss_test.BaseErrorResponse)
	PopupWarehouseLocationLevel(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
}
