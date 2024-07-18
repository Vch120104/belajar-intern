package transactionworkshopservice

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"
)

type VehicleHistoryService interface {
	GetAllVehicleHistory(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetVehicleHistoryById(int) (transactionworkshoppayloads.VehicleHistoryByIdResponses, *exceptions.BaseErrorResponse)
}
