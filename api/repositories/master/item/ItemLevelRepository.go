package masteritemrepository

import (
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemLevelRepository interface {
	WithTrx(trxHandle *gorm.DB) ItemLevelRepository
	Save(masteritemlevelpayloads.SaveItemLevelRequest) (bool, error)
	Update(masteritemlevelpayloads.SaveItemLevelRequest) (bool, error)
	GetById(int) (masteritemlevelpayloads.GetItemLevelResponse, error)
	GetAll(request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) (pagination.Pagination, error)
	ChangeStatus(int) (masteritemlevelpayloads.GetItemLevelResponse, error)
}
