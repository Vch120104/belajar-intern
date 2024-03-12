package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupMasterRepository interface {
	GetMarkupMasterList(*gorm.DB, []utils.FilterCondition, pagination.Pagination) (pagination.Pagination, error)
	GetMarkupMasterById(*gorm.DB, int) (masteritempayloads.MarkupMasterResponse, error)
	SaveMarkupMaster(*gorm.DB, masteritempayloads.MarkupMasterResponse) (bool, error)
	ChangeStatusMasterMarkupMaster(tx *gorm.DB, Id int) (bool, error)
	GetMarkupMasterByCode(*gorm.DB, string) (masteritempayloads.MarkupMasterResponse, error)
}
