package masterwarehouserepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseMasterRepository interface {
	Save(*gorm.DB, masterwarehousepayloads.GetWarehouseMasterResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	GetById(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse)
	GetAll(*gorm.DB, masterwarehousepayloads.GetAllWarehouseMasterRequest, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllIsActive(*gorm.DB) ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse)
	GetWarehouseMasterByCode(*gorm.DB, string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	GetWarehouseWithMultiId(*gorm.DB, []string) ([]masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse)
	ChangeStatus(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse)
	DropdownWarehouse(*gorm.DB) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse)
}
