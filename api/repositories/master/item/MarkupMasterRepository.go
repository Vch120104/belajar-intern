package masteritemrepository

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupMasterRepository interface {
	GetMarkupMasterList(tx *gorm.DB,filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse)
	GetMarkupMasterById(tx *gorm.DB,Id int) (masteritempayloads.MarkupMasterResponse, *exceptionsss_test.BaseErrorResponse)
	SaveMarkupMaster(tx *gorm.DB,req masteritempayloads.MarkupMasterResponse) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusMasterMarkupMaster(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse)
	GetMarkupMasterByCode(*gorm.DB, string) (masteritempayloads.MarkupMasterResponse, *exceptionsss_test.BaseErrorResponse)
}
