package masterservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LookupService interface {
	ItemOprCode(linetypeId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPrice(linetypeId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	VehicleUnitMaster(brandId int, modelId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CampaignMaster(companyId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
