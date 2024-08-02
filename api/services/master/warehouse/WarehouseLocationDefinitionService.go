package masterwarehouseservice

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type WarehouseLocationDefinitionService interface {
	GetByLevel(idlevel int, idwhl string) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptions.BaseErrorResponse)
	SaveData(masterwarehousepayloads.WarehouseLocationDefinitionResponse) (masterwarehouseentities.WarehouseLocationDefinition, *exceptions.BaseErrorResponse)
	Save(masterwarehousepayloads.WarehouseLocationDefinitionResponse) (masterwarehouseentities.WarehouseLocationDefinition, *exceptions.BaseErrorResponse)
	GetById(Id int) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptions.BaseErrorResponse)
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ChangeStatus(Id int) (masterwarehouseentities.WarehouseLocationDefinition, *exceptions.BaseErrorResponse)
	PopupWarehouseLocationLevel(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
