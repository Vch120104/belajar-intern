package masteritemrepository

import (
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemPackageDetailRepository interface {
	GetItemPackageDetailByItemPackageId(tx *gorm.DB, itemPackageId int, pages pagination.Pagination) (pagination.Pagination, error)
}
