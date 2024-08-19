package transactionworkshoprepository

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ServiceWorkshopRepository interface {
	GetAllByTechnicianWO(*gorm.DB, int, int, []utils.FilterCondition, pagination.Pagination) (transactionworkshoppayloads.ServiceWorkshopDetailResponse, *exceptions.BaseErrorResponse)
	StartService(*gorm.DB, int, int, int) (bool, *exceptions.BaseErrorResponse)
	PendingService(*gorm.DB, int, int, int) (bool, *exceptions.BaseErrorResponse)
}
