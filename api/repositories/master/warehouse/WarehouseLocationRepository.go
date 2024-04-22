package masterwarehouserepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseLocationRepository interface {
	Save(*gorm.DB, masterwarehousepayloads.GetWarehouseLocationResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	GetById(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseLocationResponse, *exceptionsss_test.BaseErrorResponse)
	GetAll(*gorm.DB, masterwarehousepayloads.GetAllWarehouseLocationRequest, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	ChangeStatus(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
}
