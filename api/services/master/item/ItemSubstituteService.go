package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemSubstituteService interface {
	GetAllItemSubstitute(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemSubstitute(id int) (masteritempayloads.ItemSubstitutePayloads, *exceptions.BaseErrorResponse)
	GetAllItemSubstituteDetail(pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemSubstituteDetail(id int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse)
	SaveItemSubstitute(req masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptions.BaseErrorResponse)
	SaveItemSubstituteDetail(req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemOperation(id int) (bool, *exceptions.BaseErrorResponse)
	DeactivateItemSubstituteDetail(id string) (bool, *exceptions.BaseErrorResponse)
	ActivateItemSubstituteDetail(id string) (bool, *exceptions.BaseErrorResponse)
}
