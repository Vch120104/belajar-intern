package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemSubstituteRepository interface {
	GetAllItemSubstitute(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemSubstitute(*gorm.DB, int) (masteritempayloads.ItemSubstitutePayloads, *exceptions.BaseErrorResponse)
	GetAllItemSubstituteDetail(*gorm.DB, pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemSubstituteDetail(*gorm.DB, int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse)
	SaveItemSubstitute(*gorm.DB, masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptions.BaseErrorResponse)
	SaveItemSubstituteDetail(*gorm.DB, masteritempayloads.ItemSubstituteDetailPostPayloads, int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemOperation(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	DeactivateItemSubstituteDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	ActivateItemSubstituteDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
}
