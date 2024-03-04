package masterwarehouseservice

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
)

type WarehouseMasterService interface {
	Save(masterwarehousepayloads.GetWarehouseMasterResponse) bool
	GetById(int) masterwarehousepayloads.GetWarehouseMasterResponse
	GetAll(request masterwarehousepayloads.GetAllWarehouseMasterRequest, pages pagination.Pagination) pagination.Pagination
	GetAllIsActive() []masterwarehousepayloads.IsActiveWarehouseMasterResponse
	GetWarehouseMasterByCode(Code string) []map[string]interface{}
	GetWarehouseWithMultiId(MultiIds []string) []masterwarehousepayloads.GetAllWarehouseMasterResponse
	ChangeStatus(int) masterwarehousepayloads.GetWarehouseMasterResponse
}
