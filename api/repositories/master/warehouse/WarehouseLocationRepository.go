package masterwarehouserepository

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseLocationRepository interface {
	Save(*gorm.DB, masterwarehousepayloads.GetWarehouseLocationResponse) (bool, error)
	GetById(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseLocationResponse, error)
	GetAll(*gorm.DB, masterwarehousepayloads.GetAllWarehouseLocationRequest, pagination.Pagination) (pagination.Pagination, error)
	ChangeStatus(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseLocationResponse, error)
}
