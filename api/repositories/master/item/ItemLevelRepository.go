package masteritemrepository

import (
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemLevelRepository interface {
	Save(*gorm.DB, masteritemlevelpayloads.SaveItemLevelRequest) (bool, error)
	GetById(*gorm.DB, int) (masteritemlevelpayloads.GetItemLevelResponseById, error)
	GetAll(tx *gorm.DB, request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) (pagination.Pagination, error)
	ChangeStatus(*gorm.DB, int) (bool, error)
}
