package masteritemservice

import (
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemLevelService interface {
	WithTrx(Trxhandle *gorm.DB) ItemLevelService
	Save(masteritemlevelpayloads.SaveItemLevelRequest) (bool, error)
	GetById(int) (masteritemlevelpayloads.GetItemLevelResponseById, error)
	GetAll(request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) (pagination.Pagination, error)
	ChangeStatus(int) (bool, error)
}
