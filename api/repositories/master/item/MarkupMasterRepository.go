package masteritemrepository

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupMasterRepository interface {
	GetMarkupMasterList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetMarkupMasterById(tx *gorm.DB, Id int) (masteritempayloads.MarkupMasterResponse, *exceptions.BaseErrorResponse)
	GetAllMarkupMasterIsActive(tx *gorm.DB) ([]masteritempayloads.MarkupMasterDropDownResponse, *exceptions.BaseErrorResponse)
	SaveMarkupMaster(tx *gorm.DB, req masteritempayloads.MarkupMasterResponse) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusMasterMarkupMaster(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse)
	GetMarkupMasterByCode(*gorm.DB, string) (masteritempayloads.MarkupMasterResponse, *exceptions.BaseErrorResponse)
}
