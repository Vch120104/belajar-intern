package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	"after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemGroupRepository interface {
	GetAllItemGroupWithPagination(db *gorm.DB, internalFilter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllItemGroup(db *gorm.DB, code string) ([]masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
	GetItemGroupById(db *gorm.DB, id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
	GetItemGroupByCode(db *gorm.DB, code string) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
	DeleteItemGroupById(db *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse)
	UpdateItemGroupById(tx *gorm.DB, payload masteritempayloads.ItemGroupUpdatePayload, id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
	UpdateStatusItemGroupById(tx *gorm.DB, id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
	GetItemGroupByMultiId(db *gorm.DB, multiId string) ([]masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
	NewItemGroup(db *gorm.DB, payload masteritempayloads.NewItemGroupPayload) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
}
