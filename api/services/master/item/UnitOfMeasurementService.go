package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type UnitOfMeasurementService interface {
	GetAllUnitOfMeasurement(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllUnitOfMeasurementIsActive() ([]masteritempayloads.UomResponse, *exceptionsss_test.BaseErrorResponse)
	GetUnitOfMeasurementById(id int) (masteritempayloads.UomIdCodeResponse, *exceptionsss_test.BaseErrorResponse)
	GetUnitOfMeasurementByCode(Code string) (masteritempayloads.UomIdCodeResponse, *exceptionsss_test.BaseErrorResponse)
	SaveUnitOfMeasurement(req masteritempayloads.UomResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusUnitOfMeasurement(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
