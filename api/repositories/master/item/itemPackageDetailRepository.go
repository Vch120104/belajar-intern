package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemPackageDetailRepository interface {
	GetItemPackageDetailByItemPackageId(tx *gorm.DB, itemPackageId int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	CreateItemPackageDetailByItemPackageId(tx *gorm.DB, req masteritempayloads.SaveItemPackageDetail) (bool, *exceptionsss_test.BaseErrorResponse)
	UpdateItemPackageDetailByItemPackageId(tx *gorm.DB, req masteritempayloads.SaveItemPackageDetail) (bool, *exceptionsss_test.BaseErrorResponse)
	GetItemPackageDetailById(tx *gorm.DB, itemPackageDetailId int) (masteritempayloads.ItemPackageDetailResponse, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusItemPackageDetail(tx *gorm.DB, id int) (bool, *exceptionsss_test.BaseErrorResponse)
}
