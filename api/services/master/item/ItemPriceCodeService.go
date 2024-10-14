package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type ItemPriceCodeService interface {
	GetAllItemPriceCode(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetByIdItemPriceCode(id int) (masteritempayloads.SaveItemPriceCode, *exceptions.BaseErrorResponse)
	GetByCodeItemPriceCode(ItemPriceCode string) (masteritempayloads.SaveItemPriceCode, *exceptions.BaseErrorResponse)
	SaveItemPriceCode(request masteritempayloads.SaveItemPriceCode) (masteritementities.ItemPriceCode, *exceptions.BaseErrorResponse)
	DeleteItemPriceCode(Id string) (bool, *exceptions.BaseErrorResponse)
	UpdateItemPriceCode(ItemPriceId int, req masteritempayloads.UpdateItemPriceCode) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemPriceCode(Id int) (bool, *exceptions.BaseErrorResponse)
}
