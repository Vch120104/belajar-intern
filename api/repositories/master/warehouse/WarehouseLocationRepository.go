package masterwarehouserepository

import (
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseLocationRepository interface {
	Save(*gorm.DB, masterwarehousepayloads.GetWarehouseLocationResponse) (bool, *exceptions.BaseErrorResponse)
	GetById(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseLocationResponse, *exceptions.BaseErrorResponse)
	GetAll(*gorm.DB, masterwarehousepayloads.GetAllWarehouseLocationRequest, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ChangeStatus(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
}
