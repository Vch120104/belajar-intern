package masterwarehouseservice

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type WarehouseMasterService interface {
	Save(masterwarehousepayloads.GetWarehouseMasterResponse) (masterwarehouseentities.WarehouseMaster, *exceptions.BaseErrorResponse)
	Update(id int, companyId int, request masterwarehousepayloads.UpdateWarehouseMasterRequest) (masterwarehouseentities.WarehouseMaster, *exceptions.BaseErrorResponse)
	GetById(warehouseId int, pagination pagination.Pagination) (masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllIsActive() ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	DropdownWarehouse() ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetWarehouseMasterByCode(Code string) (masterwarehousepayloads.GetAllWarehouseMasterCodeResponse, *exceptions.BaseErrorResponse)
	GetWarehouseWithMultiId(MultiIds []int) ([]masterwarehousepayloads.GetAllWarehouseMasterCodeResponse, *exceptions.BaseErrorResponse)
	IsWarehouseMasterByCodeAndCompanyIdExist(int, string) bool
	GetWarehouseGroupAndMasterbyCodeandCompanyId(int, string) (int, int, *exceptions.BaseErrorResponse)
	ChangeStatus(int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	DropdownbyGroupId(int) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetAuthorizeUser(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	PostAuthorizeUser(req masterwarehousepayloads.WarehouseAuthorize) (masterwarehousepayloads.WarehouseAuthorize, *exceptions.BaseErrorResponse)
	DeleteMultiIdAuthorizeUser(id string) (bool, *exceptions.BaseErrorResponse)
	InTransitWarehouseCodeDropdown(int, int) ([]masterwarehousepayloads.DropdownWarehouseMasterByCodeResponse, *exceptions.BaseErrorResponse)
}
