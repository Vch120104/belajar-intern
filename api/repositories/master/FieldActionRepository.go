package masterrepository

import (
	// masterpayloads "after-sales/api/payloads/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type FieldActionRepository interface {
	GetAllFieldAction(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	SaveFieldAction(tx *gorm.DB, req masterpayloads.FieldActionResponse) (bool, error)

	GetFieldActionHeaderById(tx *gorm.DB, Id int) (masterpayloads.FieldActionResponse, error)
	GetFieldActionDetailById(tx *gorm.DB, Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error)
}

