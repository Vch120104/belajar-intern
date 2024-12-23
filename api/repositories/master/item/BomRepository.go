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

type BomRepository interface {
	// Parent
	GetBomList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomById(tx *gorm.DB, id int) (masteritempayloads.BomResponse, *exceptions.BaseErrorResponse)
	// Parent (unfinished)
	SaveBomMaster(*gorm.DB, masteritempayloads.BomMasterRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	UpdateBomMaster(*gorm.DB, int, masteritempayloads.BomMasterRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	ChangeStatusBomMaster(tx *gorm.DB, Id int) (masteritementities.Bom, *exceptions.BaseErrorResponse)

	// Child
	GetBomDetailByMasterId(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomDetailByMasterUn(tx *gorm.DB, id int, effectiveDate time.Time, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomDetailById(tx *gorm.DB, id int) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)
	// Child (unfinished)
	GetBomDetailList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SaveBomDetail(*gorm.DB, masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)
	UpdateBomDetail(*gorm.DB, int, masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)
	DeleteByIds(*gorm.DB, []int) (bool, *exceptions.BaseErrorResponse)

	GetBomItemList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
