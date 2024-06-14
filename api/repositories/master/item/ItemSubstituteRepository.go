package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"time"

	"gorm.io/gorm"
)

type ItemSubstituteRepository interface {
	GetAllItemSubstitute(*gorm.DB, []utils.FilterCondition, pagination.Pagination, time.Time, time.Time) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemSubstitute(*gorm.DB, int) (masteritempayloads.ItemSubstitutePayloads, *exceptions.BaseErrorResponse)
	GetAllItemSubstituteDetail(*gorm.DB, pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemSubstituteDetail(*gorm.DB, int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse)
	SaveItemSubstitute(*gorm.DB, masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptions.BaseErrorResponse)
	SaveItemSubstituteDetail(*gorm.DB, masteritempayloads.ItemSubstituteDetailPostPayloads, int) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusItemSubstitute(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	DeactivateItemSubstituteDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	ActivateItemSubstituteDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
}
