package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type UnitOfMeasurementService interface {
	GetAllUnitOfMeasurement(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	GetAllUnitOfMeasurementIsActive() []masteritempayloads.UomResponse
	GetUnitOfMeasurementById(Id int) masteritempayloads.UomIdCodeResponse
	GetUnitOfMeasurementByCode(Code string) masteritempayloads.UomIdCodeResponse
	SaveUnitOfMeasurement(req masteritempayloads.UomResponse) bool
	ChangeStatusUnitOfMeasurement(Id int) bool
}
