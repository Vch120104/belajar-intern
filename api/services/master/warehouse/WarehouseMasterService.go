package masterwarehouseservice

import (
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type WarehouseMasterService interface {
	Save(masterwarehousepayloads.GetWarehouseMasterResponse) (bool, *exceptions.BaseErrorResponse)
	GetById(int) (masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllIsActive() ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	DropdownWarehouse() ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetWarehouseMasterByCode(Code string) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	GetWarehouseWithMultiId(MultiIds []string) ([]masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	IsWarehouseMasterByCodeAndCompanyIdExist(int, string) bool
	GetWarehouseGroupAndMasterbyCodeandCompanyId(int, string) (int, int, *exceptions.BaseErrorResponse)
	ChangeStatus(int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	DropdownbyGroupId(int) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse)
}
