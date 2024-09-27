package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"time"
)

type ItemSubstituteService interface {
	GetAllItemSubstitute(filterCondition []utils.FilterCondition, pages pagination.Pagination, from time.Time, to time.Time) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetByIdItemSubstitute(id int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllItemSubstituteDetail(pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemSubstituteDetail(id int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse)
	SaveItemSubstitute(req masteritempayloads.ItemSubstitutePostPayloads) (masteritementities.ItemSubstitute, *exceptions.BaseErrorResponse)
	SaveItemSubstituteDetail(req masteritempayloads.ItemSubstituteDetailPostPayloads, id int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemSubstitute(id int) (bool, *exceptions.BaseErrorResponse)
	DeactivateItemSubstituteDetail(id string) (bool, *exceptions.BaseErrorResponse)
	ActivateItemSubstituteDetail(id string) (bool, *exceptions.BaseErrorResponse)
	GetallItemForFilter(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
