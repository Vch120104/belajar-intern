package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemPriceCodeRepository interface {
	GetAllItemPriceCode(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetByIdItemPriceCode(tx *gorm.DB, id int) (masteritempayloads.SaveItemPriceCode, *exceptions.BaseErrorResponse)
	GetByCodeItemPriceCode(tx *gorm.DB, itemPriceCode string) (masteritempayloads.SaveItemPriceCode, *exceptions.BaseErrorResponse)
	SaveItemPriceCode(tx *gorm.DB, request masteritempayloads.SaveItemPriceCode) (masteritementities.ItemPriceCode, *exceptions.BaseErrorResponse)
	DeleteItemPriceCode(tx *gorm.DB, Id string) (bool, *exceptions.BaseErrorResponse)
	UpdateItemPriceCode(tx *gorm.DB, itemPriceId int, req masteritempayloads.UpdateItemPriceCode) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemPriceCode(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
}
