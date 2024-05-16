package masteritemrepository

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type BomRepository interface {
	GetBomMasterList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetBomMasterById(*gorm.DB, int) (masteritempayloads.BomMasterRequest, *exceptionsss_test.BaseErrorResponse)
	SaveBomMaster(*gorm.DB, masteritempayloads.BomMasterRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	ChangeStatusBomMaster(tx *gorm.DB, Id int) (masteritementities.Bom, *exceptionsss_test.BaseErrorResponse)
	GetBomDetailList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	GetBomDetailById(*gorm.DB, int) ([]masteritempayloads.BomDetailListResponse, *exceptionsss_test.BaseErrorResponse)
	GetBomDetailByIds(*gorm.DB, int) ([]masteritempayloads.BomDetailListResponse, *exceptionsss_test.BaseErrorResponse)
	SaveBomDetail(*gorm.DB, masteritempayloads.BomDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse)
	GetBomItemList(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse)
	DeleteByIds(*gorm.DB, []int) (bool, *exceptionsss_test.BaseErrorResponse)
}
