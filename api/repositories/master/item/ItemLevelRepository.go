package masteritemrepository

import (
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemLevelRepository interface {
	WithTrx(trxHandle *gorm.DB) ItemLevelRepository
	Save(masteritemlevelpayloads.SaveItemLevelRequest) (bool, error)
	GetById(int) (masteritemlevelpayloads.GetItemLevelResponseById, error)
	GetAll(request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) (pagination.Pagination, error)
	ChangeStatus(int) (bool, error)
}
