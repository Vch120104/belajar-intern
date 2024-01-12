package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupMasterRepository interface {
	WithTrx(trxHandle *gorm.DB) MarkupMasterRepository
	GetMarkupMasterList([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetMarkupMasterById(int) (masteritempayloads.MarkupMasterResponse, error)
	SaveMarkupMaster(masteritempayloads.MarkupMasterResponse) (bool, error)
	ChangeStatusMasterMarkupMaster(Id int) (bool, error)
	GetMarkupMasterByCode(string) (masteritempayloads.MarkupMasterResponse, error)
}
