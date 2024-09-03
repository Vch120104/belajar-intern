package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type LookupRepository interface {
	ItemOprCode(tx *gorm.DB, linetypeId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPrice(tx *gorm.DB, linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode string, whsGroup string, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	VehicleUnitMaster(tx *gorm.DB, brandId int, modelId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CampaignMaster(tx *gorm.DB, companyId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
