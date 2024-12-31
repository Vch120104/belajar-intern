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
	SaveBomMaster(*gorm.DB, masteritempayloads.BomMasterNewRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	FirstOrCreateBom(*gorm.DB, masteritempayloads.BomMasterNewRequest) (int, *exceptions.BaseErrorResponse)
	UpdateBomMaster(tx *gorm.DB, id int, qty float64) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	ChangeStatusBomMaster(tx *gorm.DB, Id int) (masteritementities.Bom, *exceptions.BaseErrorResponse)

	// Child
	GetBomDetailByMasterId(tx *gorm.DB, id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomDetailByMasterUn(tx *gorm.DB, id int, effectiveDate time.Time, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomDetailById(tx *gorm.DB, id int) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)
	GetBomDetailTemplate(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]masteritempayloads.BomDetailTemplate, *exceptions.BaseErrorResponse)
	GetBomDetailMaxSeq(tx *gorm.DB, id int) (int, *exceptions.BaseErrorResponse)
	SaveBomDetail(*gorm.DB, masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)
	DeleteByIds(*gorm.DB, []int) (bool, *exceptions.BaseErrorResponse)
}
