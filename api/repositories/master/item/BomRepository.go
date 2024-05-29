package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type BomRepository interface {
	GetBomMasterList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetBomMasterById(*gorm.DB, int) (masteritempayloads.BomMasterRequest, *exceptions.BaseErrorResponse)
	SaveBomMaster(*gorm.DB, masteritempayloads.BomMasterRequest) (bool, *exceptions.BaseErrorResponse)
	ChangeStatusBomMaster(tx *gorm.DB, Id int) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	GetBomDetailList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetBomDetailById(*gorm.DB, int) ([]masteritempayloads.BomDetailListResponse, *exceptions.BaseErrorResponse)
	GetBomDetailByIds(*gorm.DB, int) ([]masteritempayloads.BomDetailListResponse, *exceptions.BaseErrorResponse)
	SaveBomDetail(*gorm.DB, masteritempayloads.BomDetailRequest) (bool, *exceptions.BaseErrorResponse)
	GetBomItemList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	DeleteByIds(*gorm.DB, []int) (bool, *exceptions.BaseErrorResponse)
}
