package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	"after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemGroupService interface {
	GetAllItemGroupWithPagination(internalFilter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetAllItemGroup(code string) ([]masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
	GetItemGroupById(id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
	DeleteItemGroupById(id int) (bool, *exceptions.BaseErrorResponse)
	UpdateItemGroupById(payload masteritempayloads.ItemGroupUpdatePayload, id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
	UpdateStatusItemGroupById(id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
	GetItemGroupByMultiId(multiId string) ([]masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
	NewItemGroup(payload masteritempayloads.NewItemGroupPayload) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse)
}
