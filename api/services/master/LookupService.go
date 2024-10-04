package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LookupService interface {
	GetLineTypeByItemCode(itemCode string) (int, *exceptions.BaseErrorResponse)
	GetOprItemPrice(linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode string, whsGroup string) (float64, *exceptions.BaseErrorResponse)
	ItemOprCode(linetypeId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeByCode(linetypeId int, oprItemCode string, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeByID(linetypeId int, oprItemId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPrice(linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode string, whsGroup string, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	VehicleUnitMaster(brandId int, modelId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleUnitByID(vehicleID int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleUnitByChassisNumber(chassisNumber string, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CampaignMaster(companyId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddress(pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddressByID(customerId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddressByCode(customerCode string, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	WorkOrderService(pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetItemLocationWarehouse(companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetWarehouseGroupByCompany(companyId int) ([]masterpayloads.WarehouseGroupByCompanyResponse, *exceptions.BaseErrorResponse)
}
