package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type UnitOfMeasurementRepository interface {
	GetAllUnitOfMeasurement(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetAllUnitOfMeasurementIsActive(tx *gorm.DB) ([]masteritempayloads.UomResponse, *exceptionsss_test.BaseErrorResponse)
	GetUnitOfMeasurementById(tx *gorm.DB,Id int) (masteritempayloads.UomIdCodeResponse, *exceptionsss_test.BaseErrorResponse)
	GetUnitOfMeasurementByCode(tx *gorm.DB,Code string) (masteritempayloads.UomIdCodeResponse, *exceptionsss_test.BaseErrorResponse)
	SaveUnitOfMeasurement(tx *gorm.DB,req masteritempayloads.UomResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusUnitOfMeasurement(tx *gorm.DB,Id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
