package masteritemservice


import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type UnitOfMeasurementService interface {
	WithTrx(trxHandle *gorm.DB) UnitOfMeasurementService
	GetAllUnitOfMeasurement(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error)
	GetAllUnitOfMeasurementIsActive() ([]masteritempayloads.UomResponse, error)
	GetUnitOfMeasurementById(Id int) (masteritempayloads.UomResponse, error)
	GetUnitOfMeasurementByCode(Code string) (masteritempayloads.UomResponse, error)
	SaveUnitOfMeasurement(req masteritempayloads.UomResponse) (bool, error)
	ChangeStatusUnitOfMeasurement(Id int) (bool, error)
}
