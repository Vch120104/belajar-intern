package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemSubstituteRepository interface {
	GetAllItemSubstitute(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetByIdItemSubstitute(*gorm.DB, int) (masteritempayloads.ItemSubstitutePayloads, *exceptionsss_test.BaseErrorResponse)
	GetAllItemSubstituteDetail(*gorm.DB, pagination.Pagination, int) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetByIdItemSubstituteDetail(*gorm.DB, int) (masteritempayloads.ItemSubstituteDetailGetPayloads, *exceptionsss_test.BaseErrorResponse)
	SaveItemSubstitute(*gorm.DB, masteritempayloads.ItemSubstitutePostPayloads) (bool, *exceptionsss_test.BaseErrorResponse)
	SaveItemSubstituteDetail(*gorm.DB, masteritempayloads.ItemSubstituteDetailPostPayloads, int) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusItemOperation(*gorm.DB, int) (bool, *exceptionsss_test.BaseErrorResponse)
	DeactivateItemSubstituteDetail(*gorm.DB, string) (bool, *exceptionsss_test.BaseErrorResponse)
	ActivateItemSubstituteDetail(*gorm.DB, string) (bool, *exceptionsss_test.BaseErrorResponse)
}
