package masterwarehouserepository

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WarehouseLocationRepository interface {
	Save(*gorm.DB, masterwarehouseentities.WarehouseLocation) (bool, *exceptions.BaseErrorResponse)
	ProcessWarehouseLocationTemplate(*gorm.DB, masterwarehousepayloads.ProcessWarehouseLocationTemplate) (bool, *exceptions.BaseErrorResponse)
	GetById(*gorm.DB, int) (masterwarehousepayloads.GetWarehouseLocationResponse, *exceptions.BaseErrorResponse)
	GetAll(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	CheckIfLocationExist(*gorm.DB, string, string, string) bool
	ChangeStatus(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
}
