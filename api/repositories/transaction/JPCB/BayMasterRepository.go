package transactionjpcbrepository

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type BayMasterRepository interface {
	GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllActive(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetAllDeactive(tx *gorm.DB, filterCondition []utils.FilterCondition) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
}
