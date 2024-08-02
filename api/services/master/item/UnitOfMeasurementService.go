package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type UnitOfMeasurementService interface {
	GetAllUnitOfMeasurement(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllUnitOfMeasurementIsActive() ([]masteritempayloads.UomResponse, *exceptions.BaseErrorResponse)
	GetUnitOfMeasurementById(id int) (masteritempayloads.UomIdCodeResponse, *exceptions.BaseErrorResponse)
	GetUnitOfMeasurementByCode(Code string) (masteritempayloads.UomIdCodeResponse, *exceptions.BaseErrorResponse)
	SaveUnitOfMeasurement(req masteritempayloads.UomResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusUnitOfMeasurement(Id int) (bool, *exceptions.BaseErrorResponse)
}
