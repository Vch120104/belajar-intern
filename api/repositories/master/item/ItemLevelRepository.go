package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemLevelRepository interface {
	Save(*gorm.DB, masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptions.BaseErrorResponse)
	GetById(tx *gorm.DB, itemLevel int, itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptions.BaseErrorResponse)
	GetAll(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemLevelDropDown(tx *gorm.DB, itemLevel string) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptions.BaseErrorResponse)
	GetItemLevelLookUp(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination, itemClassId int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetItemLevelLookUpbyId(tx *gorm.DB, itemLevelId int) (masteritemlevelpayloads.GetItemLevelLookUp, *exceptions.BaseErrorResponse)
	ChangeStatus(tx *gorm.DB, itemLevel int, itemLevelId int) (bool, *exceptions.BaseErrorResponse)
}
