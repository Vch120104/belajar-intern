package masterservice

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LookupService interface {
	ItemOprCode(linetypeId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPrice(linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode string, whsGroup string, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	VehicleUnitMaster(brandId int, modelId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleUnitByID(vehicleID int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleUnitByChassisNumber(chassisNumber string, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CampaignMaster(companyId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	WorkOrderService(pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
}
