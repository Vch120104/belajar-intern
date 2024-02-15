package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"gorm.io/gorm"
)

type UnitOfMeasurementRepository interface {
	GetAllUnitOfMeasurement(tx *gorm.DB,filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error)
	GetAllUnitOfMeasurementIsActive(tx *gorm.DB,) ([]masteritempayloads.UomResponse, error)
	GetUnitOfMeasurementById(tx *gorm.DB,Id int) (masteritempayloads.UomResponse, error)
	GetUnitOfMeasurementByCode(tx *gorm.DB,Code string) (masteritempayloads.UomResponse, error)
	SaveUnitOfMeasurement(tx *gorm.DB,req masteritempayloads.UomResponse) (bool, error)
	ChangeStatusUnitOfMeasurement(tx *gorm.DB,Id int) (bool, error)
}
