package masteritemrepository

import (
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemPackageRepository interface {
	GetAllItemPackage(tx *gorm.DB, filterCondition []utils.FilterCondition) (pagination.Pagination, error)
}
