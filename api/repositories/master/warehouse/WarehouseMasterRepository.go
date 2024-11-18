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
	Save(tx *gorm.DB, request masterwarehousepayloads.GetWarehouseMasterResponse) (masterwarehouseentities.WarehouseMaster, *exceptions.BaseErrorResponse)
	Update(tx *gorm.DB, warehouseId int, companyId int, request masterwarehousepayloads.UpdateWarehouseMasterRequest) (masterwarehouseentities.WarehouseMaster, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, warehouseId int, pagination pagination.Pagination) (masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetAll(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllIsActive(tx *gorm.DB) ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetWarehouseMasterByCode(tx *gorm.DB, warehouseCode string) (masterwarehousepayloads.GetAllWarehouseMasterCodeResponse, *exceptions.BaseErrorResponse)
	IsWarehouseMasterByCodeAndCompanyIdExist(tx *gorm.DB, companyId int, warehouseCodes []string) ([]masterwarehouseentities.WarehouseMaster, *exceptions.BaseErrorResponse)
	GetWarehouseGroupAndMasterbyCodeandCompanyId(tx *gorm.DB, companyId int, warehouseCode string) (int, int, *exceptions.BaseErrorResponse)
	GetWarehouseWithMultiId(tx *gorm.DB, MultiIds []int) ([]masterwarehousepayloads.GetAllWarehouseMasterCodeResponse, *exceptions.BaseErrorResponse)
	ChangeStatus(tx *gorm.DB, warehouseId int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	DropdownWarehouse(tx *gorm.DB) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	DropdownbyGroupId(tx *gorm.DB, warehouseGroupId int, companyId int) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetAuthorizeUser(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	PostAuthorizeUser(tx *gorm.DB, req masterwarehousepayloads.WarehouseAuthorize) (masterwarehousepayloads.WarehouseAuthorize, *exceptions.BaseErrorResponse)
	DeleteMultiIdAuthorizeUser(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	InTransitWarehouseCodeDropdown(tx *gorm.DB, companyID int, warehouseGroupID int) ([]masterwarehousepayloads.DropdownWarehouseMasterByCodeResponse, *exceptions.BaseErrorResponse)
}
