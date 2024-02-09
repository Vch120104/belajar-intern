package masteritemservice

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemSubstituteService interface {
	GetAllItemSubstitute(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination
	GetByIdItemSubstitute(id int) masteritempayloads.ItemSubstitutePayloads
	GetAllItemSubstituteDetail([]utils.FilterCondition, pagination.Pagination,int) pagination.Pagination
	GetByIdItemSubstituteDetail(id int) masteritempayloads.ItemSubstituteDetailPayloads
	SaveItemSubstitute(req masteritempayloads.ItemSubstitutePayloads) bool
	SaveItemSubstituteDetail(req masteritempayloads.ItemSubstituteDetailPayloads) bool
	ChangeStatusItemOperation(id int) bool
	DeactivateItemSubstituteDetail(id string)bool
	ActivateItemSubstituteDetail(id string)bool
}