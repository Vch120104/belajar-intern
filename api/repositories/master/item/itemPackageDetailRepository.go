package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemPackageDetailRepository interface {
	GetItemPackageDetailByItemPackageId(tx *gorm.DB, itemPackageId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	CreateItemPackageDetailByItemPackageId(tx *gorm.DB, req masteritempayloads.SaveItemPackageDetail) (bool, *exceptions.BaseErrorResponse)
	UpdateItemPackageDetail(tx *gorm.DB, req masteritempayloads.SaveItemPackageDetail) (bool, *exceptions.BaseErrorResponse)
	GetItemPackageDetailById(tx *gorm.DB, itemPackageDetailId int) (masteritempayloads.ItemPackageDetailResponse, *exceptions.BaseErrorResponse)
	ChangeStatusItemPackageDetail(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
	DeactiveItemPackageDetail(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
	ActivateItemPackageDetail(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse)
}
