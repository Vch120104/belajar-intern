package masteritemservice

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
)

type ItemSubstituteService interface {
	GetByIdItemSubstitute(id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllItemSubstituteDetail(pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemSubstituteDetail(id int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse)
	SaveItemSubstitute(req masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptions.BaseErrorResponse)
	SaveItemSubstituteDetail(req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemOperation(id int) (bool, *exceptions.BaseErrorResponse)
	DeactivateItemSubstituteDetail(id string) (bool, *exceptions.BaseErrorResponse)
	ActivateItemSubstituteDetail(id string) (bool, *exceptions.BaseErrorResponse)
	GetAllItemSubstitute(filterCondition map[string]string, pages pagination.Pagination) ([]map[string]interface{},int,int, *exceptions.BaseErrorResponse)
}