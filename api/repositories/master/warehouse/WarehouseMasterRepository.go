package masterwarehouserepository

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseMasterRepository interface {
	Save(*gorm.DB, masterwarehousepayloads.GetWarehouseMasterResponse) (bool, error)
	GetById(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseMasterResponse, error)
	GetAll(*gorm.DB, masterwarehousepayloads.GetAllWarehouseMasterRequest, pagination.Pagination) (pagination.Pagination, error)
	GetAllIsActive(*gorm.DB) ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, error)
	GetWarehouseMasterByCode(*gorm.DB, string) ([]map[string]interface{}, error)
	GetWarehouseWithMultiId(*gorm.DB, []string) ([]masterwarehousepayloads.GetAllWarehouseMasterResponse, error)
	ChangeStatus(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseMasterResponse, error)
}
