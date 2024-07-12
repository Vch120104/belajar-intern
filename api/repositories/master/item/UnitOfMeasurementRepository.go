package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type UnitOfMeasurementRepository interface {
	GetAllUnitOfMeasurement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllUnitOfMeasurementIsActive(tx *gorm.DB) ([]masteritempayloads.UomResponse, *exceptions.BaseErrorResponse)
	GetUnitOfMeasurementById(tx *gorm.DB, Id int) (masteritempayloads.UomIdCodeResponse, *exceptions.BaseErrorResponse)
	GetUnitOfMeasurementByCode(tx *gorm.DB, Code string) (masteritempayloads.UomIdCodeResponse, *exceptions.BaseErrorResponse)
	SaveUnitOfMeasurement(tx *gorm.DB, req masteritempayloads.UomResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusUnitOfMeasurement(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
}
