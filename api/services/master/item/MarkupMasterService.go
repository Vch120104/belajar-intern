package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupMasterService interface {
	WithTrx(trxHandle *gorm.DB) MarkupMasterService
	GetMarkupMasterList([]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetMarkupMasterById(int) (masteritempayloads.MarkupMasterResponse, error)
	SaveMarkupMaster(masteritempayloads.MarkupMasterResponse) (bool, error)
	ChangeStatusMasterMarkupMaster(Id int) (bool, error)
	GetMarkupMasterByCode(string) (masteritempayloads.MarkupMasterResponse, error)
}
