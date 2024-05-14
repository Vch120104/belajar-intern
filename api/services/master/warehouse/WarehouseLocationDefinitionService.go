package masterwarehouseservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type WarehouseLocationDefinitionService interface {
	GetByLevel(idlevel int, idwhl string) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptionsss_test.BaseErrorResponse)
	SaveData(masterwarehousepayloads.WarehouseLocationDefinitionResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	Save(masterwarehousepayloads.WarehouseLocationDefinitionResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	GetById(Id int) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptionsss_test.BaseErrorResponse)
	GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	ChangeStatus(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
	PopupWarehouseLocationLevel(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
}
