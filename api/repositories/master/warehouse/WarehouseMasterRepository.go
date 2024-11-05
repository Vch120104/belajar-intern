package masterwarehouserepository

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WarehouseMasterRepository interface {
	Save(*gorm.DB, masterwarehousepayloads.GetWarehouseMasterResponse) (masterwarehouseentities.WarehouseMaster, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, warehouseId int, pagination pagination.Pagination) (masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetAll(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllIsActive(*gorm.DB) ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetWarehouseMasterByCode(*gorm.DB, string) (masterwarehousepayloads.GetAllWarehouseMasterCodeResponse, *exceptions.BaseErrorResponse)
	IsWarehouseMasterByCodeAndCompanyIdExist(*gorm.DB, int, string) (bool, *exceptions.BaseErrorResponse)
	GetWarehouseGroupAndMasterbyCodeandCompanyId(*gorm.DB, int, string) (int, int, *exceptions.BaseErrorResponse)
	GetWarehouseWithMultiId(*gorm.DB, []int) ([]masterwarehousepayloads.GetAllWarehouseMasterCodeResponse, *exceptions.BaseErrorResponse)
	ChangeStatus(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	DropdownWarehouse(*gorm.DB) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	DropdownbyGroupId(*gorm.DB, int) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetAuthorizeUser(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	PostAuthorizeUser(tx *gorm.DB, req masterwarehousepayloads.WarehouseAuthorize) (masterwarehousepayloads.WarehouseAuthorize, *exceptions.BaseErrorResponse)
	DeleteMultiIdAuthorizeUser(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	InTransitWarehouseCodeDropdown(*gorm.DB, int, int) ([]masterwarehousepayloads.DropdownWarehouseMasterByCodeResponse, *exceptions.BaseErrorResponse)
}
