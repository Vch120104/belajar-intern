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
	GetBomMasterList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomMasterById(tx *gorm.DB, id int, pages pagination.Pagination) (masteritempayloads.BomMasterResponseDetail, *exceptions.BaseErrorResponse)
	SaveBomMaster(*gorm.DB, masteritempayloads.BomMasterRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	UpdateBomMaster(*gorm.DB, int, masteritempayloads.BomMasterRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	ChangeStatusBomMaster(tx *gorm.DB, Id int) (masteritementities.Bom, *exceptions.BaseErrorResponse)

	GetBomDetailList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetBomDetailById(tx *gorm.DB, id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	SaveBomDetail(*gorm.DB, masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)
	UpdateBomDetail(*gorm.DB, int, masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)
	DeleteByIds(*gorm.DB, []int) (bool, *exceptions.BaseErrorResponse)

	GetBomItemList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
