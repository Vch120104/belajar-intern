package masteritemservice

import (
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
)

type ItemLevelService interface {
	Save(masteritemlevelpayloads.SaveItemLevelRequest) bool
	GetById(int) masteritemlevelpayloads.GetItemLevelResponseById
	GetAll(request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) pagination.Pagination
	ChangeStatus(int) bool
}
