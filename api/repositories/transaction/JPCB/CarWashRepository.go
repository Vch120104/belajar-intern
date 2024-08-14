package transactionjpcbrepository

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type CarWashRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	UpdatePriority(tx *gorm.DB, workOrderSystemNumber, carWashPriorityId int) (transactionjpcbentities.CarWash, *exceptions.BaseErrorResponse)
}
