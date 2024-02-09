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
	GetAllItemSubstituteDetail(*gorm.DB,[]utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetByIdItemSubstituteDetail(*gorm.DB,int) (masteritempayloads.ItemSubstituteDetailPayloads, error)
	SaveItemSubstitute(*gorm.DB,masteritempayloads.ItemSubstitutePayloads) (bool, error)
	SaveItemSubstituteDetail(*gorm.DB,masteritempayloads.ItemSubstituteDetailPayloads) (bool, error)
	ChangeStatusItemOperation(*gorm.DB,int) (bool, error)
	DeactivateItemSubstituteDetail(*gorm.DB,string) (bool, error)
	ActivateItemSubstituteDetail(*gorm.DB,string) (bool, error)
}
