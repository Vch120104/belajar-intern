package masterwarehouseservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
)

type WarehouseMasterService interface {
	Save(masterwarehousepayloads.GetWarehouseMasterResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	GetById(int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse)
	GetAll(request masterwarehousepayloads.GetAllWarehouseMasterRequest, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllIsActive() ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse)
	DropdownWarehouse() ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse)
	GetWarehouseMasterByCode(Code string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse)
	GetWarehouseWithMultiId(MultiIds []string) ([]masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse)
	ChangeStatus(int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse)
}
