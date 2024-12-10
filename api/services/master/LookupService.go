package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LookupService interface {
	ItemOprCode(linetypeId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeByCode(linetypeId int, oprItemCode string, pages pagination.Pagination, filterCondition []utils.FilterCondition) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeByID(linetypeId int, oprItemId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPrice(linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode int, whsGroup string, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleUnitMaster(brandId int, modelId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleUnitByID(vehicleID int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleUnitByChassisNumber(chassisNumber string, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetCampaignMaster(companyId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetLineTypeByItemCode(itemCode string) (int, *exceptions.BaseErrorResponse)
	GetOprItemPrice(linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode int, whsGroup string) (float64, *exceptions.BaseErrorResponse)
	ListItemLocation(companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	WarehouseGroupByCompany(companyId int) ([]masterpayloads.WarehouseGroupByCompanyResponse, *exceptions.BaseErrorResponse)
	ItemListTrans(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemListTransPL(companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddress(pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddressByID(customerId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddressByCode(customerCode string, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	WorkOrderService(pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ReferenceTypeWorkOrder(pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ReferenceTypeWorkOrderByID(referenceId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ReferenceTypeSalesOrder(pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ReferenceTypeSalesOrderByID(referenceId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetLineTypeByReferenceType(referenceTypeId int) ([]map[string]interface{}, *exceptions.BaseErrorResponse)
	LocationAvailable(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemDetailForItemInquiry(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemSubstituteDetailForItemInquiry(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetPartNumberItemImport(internalCondition []utils.FilterCondition, externalCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
}
