package masterwarehouseservice

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseLocationService interface {
	WithTrx(Trxhandle *gorm.DB) WarehouseLocationService
	Save(masterwarehousepayloads.GetWarehouseLocationResponse) (bool, error)
	GetById(int) (masterwarehousepayloads.GetWarehouseLocationResponse, error)
	GetAll(request masterwarehousepayloads.GetAllWarehouseLocationRequest, pages pagination.Pagination) (pagination.Pagination, error)
	ChangeStatus(int) (masterwarehousepayloads.GetWarehouseLocationResponse, error)
}
