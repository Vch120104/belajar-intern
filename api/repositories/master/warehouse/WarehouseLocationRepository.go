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
	GetById(*gorm.DB, int) (masterwarehousepayloads.GetAllWarehouseLocationResponse, *exceptions.BaseErrorResponse)
	GetByCode(*gorm.DB, string) (masterwarehousepayloads.GetAllWarehouseLocationResponse, *exceptions.BaseErrorResponse)
	GetAll(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	CheckIfLocationExist(*gorm.DB, string, string, string) (bool, *exceptions.BaseErrorResponse)
	ChangeStatus(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
}
