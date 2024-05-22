package masterwarehouserepository

import (
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WarehouseGroupRepository interface {
	SaveWarehouseGroup(*gorm.DB, masterwarehousepayloads.GetWarehouseGroupResponse) (bool, *exceptions.BaseErrorResponse)
	GetByIdWarehouseGroup(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseGroupResponse, *exceptions.BaseErrorResponse)
	GetAllWarehouseGroup(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ChangeStatusWarehouseGroup(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
}
