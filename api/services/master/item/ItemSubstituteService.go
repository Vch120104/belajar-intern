package masteritemservice

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemSubstituteService interface {
	GetAllItemSubstitute(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetByIdItemSubstitute(id int) (masteritempayloads.ItemSubstitutePayloads, *exceptionsss_test.BaseErrorResponse)
	GetAllItemSubstituteDetail(pagination.Pagination, int) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetByIdItemSubstituteDetail(id int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptionsss_test.BaseErrorResponse)
	SaveItemSubstitute(req masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveItemSubstituteDetail(req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusItemOperation(id int) (bool, *exceptionsss_test.BaseErrorResponse)
	DeactivateItemSubstituteDetail(id string) (bool, *exceptionsss_test.BaseErrorResponse)
	ActivateItemSubstituteDetail(id string) (bool, *exceptionsss_test.BaseErrorResponse)
}
