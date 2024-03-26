package masteritemrepository

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type BomRepository interface {
	GetBomMasterList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error)
	GetBomMasterById(*gorm.DB, int) (masteritempayloads.BomMasterRequest, error)
	SaveBomMaster(*gorm.DB, masteritempayloads.BomMasterRequest) (bool, error)
	ChangeStatusBomMaster(tx *gorm.DB, Id int) (bool, error)
	GetBomDetailList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, error)
}
