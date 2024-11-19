package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"time"

	"gorm.io/gorm"
)

type ItemSubstituteRepository interface {
	GetAllItemSubstitute(*gorm.DB, []utils.FilterCondition, pagination.Pagination, time.Time, time.Time) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetByIdItemSubstitute(*gorm.DB, int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetAllItemSubstituteDetail(*gorm.DB, pagination.Pagination, int) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetByIdItemSubstituteDetail(*gorm.DB, int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptions.BaseErrorResponse)
	SaveItemSubstitute(*gorm.DB, masteritempayloads.ItemSubstitutePostPayloads) (masteritementities.ItemSubstitute, *exceptions.BaseErrorResponse)
	SaveItemSubstituteDetail(*gorm.DB, masteritempayloads.ItemSubstituteDetailPostPayloads, int) (masteritementities.ItemSubstituteDetail, *exceptions.BaseErrorResponse)
	ChangeStatusItemSubstitute(*gorm.DB, int) (bool, *exceptions.BaseErrorResponse)
	DeactivateItemSubstituteDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	ActivateItemSubstituteDetail(*gorm.DB, string) (bool, *exceptions.BaseErrorResponse)
	GetallItemForFilter(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
