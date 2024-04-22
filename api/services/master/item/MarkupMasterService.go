package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type MarkupMasterService interface {
	GetMarkupMasterList(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetMarkupMasterById(id int) (masteritempayloads.MarkupMasterResponse, *exceptionsss_test.BaseErrorResponse)
	SaveMarkupMaster(masteritempayloads.MarkupMasterResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusMasterMarkupMaster(Id int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetMarkupMasterByCode(markupCode string) (masteritempayloads.MarkupMasterResponse, *exceptionsss_test.BaseErrorResponse)
}
