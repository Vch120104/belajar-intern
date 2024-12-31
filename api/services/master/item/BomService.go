package masteritemservice

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"time"

	"github.com/xuri/excelize/v2"
)

type BomService interface {
	// Master
	GetBomList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomById(id int) (masteritempayloads.BomResponse, *exceptions.BaseErrorResponse)

	SaveBomMaster(request masteritempayloads.BomMasterNewRequest) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	UpdateBomMaster(id int, qty float64) (masteritementities.Bom, *exceptions.BaseErrorResponse)
	ChangeStatusBomMaster(id int) (masteritementities.Bom, *exceptions.BaseErrorResponse)

	// Detail
	GetBomDetailByMasterId(id int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomDetailByMasterUn(id int, effectiveDate time.Time, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomDetailById(id int) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)
	GetBomDetailList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetBomDetailMaxSeq(id int) (int, *exceptions.BaseErrorResponse)
	UpdateBomDetail(id int, request masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)

	SaveBomDetail(request masteritempayloads.BomDetailRequest) (masteritementities.BomDetail, *exceptions.BaseErrorResponse)
	DeleteByIds(ids []int) (bool, *exceptions.BaseErrorResponse)
	GetBomItemList(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)

	GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse)
	PreviewUploadData([][]string) ([]masteritempayloads.BomDetailTemplate, *exceptions.BaseErrorResponse)
	ProcessDataUpload(request masteritempayloads.BomDetailUpload) ([]masteritementities.BomDetail, *exceptions.BaseErrorResponse)
}
