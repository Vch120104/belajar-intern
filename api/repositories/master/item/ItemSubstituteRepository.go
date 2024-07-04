package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type ItemSubstituteRepository interface {
	GetByIdItemSubstitute(*gorm.DB, int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllItemSubstituteDetail(*gorm.DB, pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemSubstituteDetail(*gorm.DB, int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse)
	SaveItemSubstitute(*gorm.DB, masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptions.BaseErrorResponse)
	SaveItemSubstituteDetail(*gorm.DB, masteritempayloads.ItemSubstituteDetailPostPayloads, int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemOperation(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	DeactivateItemSubstituteDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	ActivateItemSubstituteDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	GetAllItemSubstitute(tx *gorm.DB, filterCondition map[string]string, pages pagination.Pagination) ([]map[string]interface{},int,int, *exceptions.BaseErrorResponse)
}