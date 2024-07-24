package transactionworkshoprepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
	"gorm.io/gorm"
)

type VehicleHistoryRepository interface {
	GetAllVehicleHistory(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetVehicleHistoryById(*gorm.DB, int) (transactionworkshoppayloads.VehicleHistoryByIdResponses, *exceptions.BaseErrorResponse)

	GetAllVehicleHistoryChassis(*gorm.DB, transactionworkshoppayloads.VehicleHistoryChassisRequest, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
