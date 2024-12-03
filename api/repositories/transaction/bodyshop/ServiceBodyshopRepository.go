package transactionbodyshoprepository

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ServiceBodyshopRepository interface {
	GetAllByTechnicianWOBodyshop(*gorm.DB, int, int, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	StartService(*gorm.DB, int, int, int) (bool, *exceptions.BaseErrorResponse)
	PendingService(*gorm.DB, int, int, int) (bool, *exceptions.BaseErrorResponse)
	TransferService(*gorm.DB, int, int, int) (bool, *exceptions.BaseErrorResponse)
	StopService(*gorm.DB, int, int, int) (bool, *exceptions.BaseErrorResponse)
}
