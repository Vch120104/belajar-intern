package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemSubstituteRepository interface {
	GetAllItemSubstitute(*gorm.DB,[]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetByIdItemSubstitute(*gorm.DB,int) (masteritempayloads.ItemSubstitutePayloads, error)
	GetAllItemSubstituteDetail(*gorm.DB, pagination.Pagination,int) (pagination.Pagination, error)
	GetByIdItemSubstituteDetail(*gorm.DB,int) (masteritempayloads.ItemSubstituteDetailGetPayloads, error)
	SaveItemSubstitute(*gorm.DB,masteritempayloads.ItemSubstitutePostPayloads) (bool, error)
	SaveItemSubstituteDetail(*gorm.DB,masteritempayloads.ItemSubstituteDetailPostPayloads,int) (bool, error)
	ChangeStatusItemOperation(*gorm.DB,int) (bool, error)
	DeactivateItemSubstituteDetail(*gorm.DB,string) (bool, error)
	ActivateItemSubstituteDetail(*gorm.DB,string) (bool, error)
}
