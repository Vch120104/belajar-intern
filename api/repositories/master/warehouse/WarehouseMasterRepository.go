package masterwarehouserepository

import (
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseMasterRepository interface {
	Save(*gorm.DB, masterwarehousepayloads.GetWarehouseMasterResponse) (bool, *exceptions.BaseErrorResponse)
	GetById(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetAll(*gorm.DB, masterwarehousepayloads.GetAllWarehouseMasterRequest, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllIsActive(*gorm.DB) ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	GetWarehouseMasterByCode(*gorm.DB, string) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	GetWarehouseWithMultiId(*gorm.DB, []string) ([]masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	ChangeStatus(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptions.BaseErrorResponse)
	DropdownWarehouse(*gorm.DB) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse)
}
