package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemSubstituteService interface {
	GetAllItemSubstitute(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	GetByIdItemSubstitute(id int) masteritempayloads.ItemSubstitutePayloads
	GetAllItemSubstituteDetail(pagination.Pagination,int) pagination.Pagination
	GetByIdItemSubstituteDetail(id int) masteritempayloads.ItemSubstituteDetailGetPayloads
	SaveItemSubstitute(req masteritempayloads.ItemSubstitutePostPayloads) bool
	SaveItemSubstituteDetail(req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) bool
	ChangeStatusItemSubstitute (id int) bool
	DeactivateItemSubstituteDetail(id string)bool
	ActivateItemSubstituteDetail(id string)bool
}