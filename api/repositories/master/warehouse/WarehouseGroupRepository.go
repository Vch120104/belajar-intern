package masterwarehouserepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WarehouseGroupRepository interface {
	SaveWarehouseGroup(*gorm.DB, masterwarehousepayloads.GetWarehouseGroupResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	GetByIdWarehouseGroup(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseGroupResponse, *exceptionsss_test.BaseErrorResponse)
	GetAllWarehouseGroup(*gorm.DB, []utils.FilterCondition,pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusWarehouseGroup(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
}
