package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type MarkupMasterService interface {
	GetMarkupMasterList([]utils.FilterCondition, pagination.Pagination) pagination.Pagination
	GetMarkupMasterById(int) masteritempayloads.MarkupMasterResponse
	SaveMarkupMaster(masteritempayloads.MarkupMasterResponse) bool
	ChangeStatusMasterMarkupMaster(Id int) bool
	GetMarkupMasterByCode(string) masteritempayloads.MarkupMasterResponse
}
