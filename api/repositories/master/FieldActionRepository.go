package masterrepository

import (
	// masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type FieldActionRepository interface {
	GetAllFieldAction(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
}
