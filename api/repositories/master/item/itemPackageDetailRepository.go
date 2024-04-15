package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemPackageDetailRepository interface {
	GetItemPackageDetailByItemPackageId(tx *gorm.DB, itemPackageId int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	SaveItemPackageDetailByItemPackageId(tx *gorm.DB, itemPackageId int, req masteritempayloads.ItemPackageDetailPayload) (bool, *exceptionsss_test.BaseErrorResponse)
}
