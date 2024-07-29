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
	GetById(int) (masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllIsActive() ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	DropdownWarehouse() ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetWarehouseMasterByCode(Code string) (masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetWarehouseWithMultiId(MultiIds []string) ([]masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	ChangeStatus(int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	DropdownbyGroupId(int) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetAuthorizeUser(pages pagination.Pagination,id int)(pagination.Pagination,*exceptions.BaseErrorResponse)
	PostAuthorizeUser(req masterwarehousepayloads.WarehouseAuthorize)(masterwarehousepayloads.WarehouseAuthorize,*exceptions.BaseErrorResponse)
	DeleteMultiIdAuthorizeUser(id string)(bool,*exceptions.BaseErrorResponse)
}
