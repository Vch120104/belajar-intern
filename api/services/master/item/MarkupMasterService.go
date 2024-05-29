package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type MarkupMasterService interface {
	GetMarkupMasterList(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetMarkupMasterById(id int) (masteritempayloads.MarkupMasterResponse, *exceptions.BaseErrorResponse)
	GetAllMarkupMasterIsActive() ([]masteritempayloads.MarkupMasterDropDownResponse, *exceptions.BaseErrorResponse)
	SaveMarkupMaster(masteritempayloads.MarkupMasterResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusMasterMarkupMaster(Id int) (bool, *exceptions.BaseErrorResponse)
	GetMarkupMasterByCode(markupCode string) (masteritempayloads.MarkupMasterResponse, *exceptions.BaseErrorResponse)
}
