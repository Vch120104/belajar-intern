package masterservice

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
)

type LookupService interface {
	ItemOprCode(linetypeId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemOprCodeByCode(linetypeId int, oprItemCode string, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemOprCodeByID(linetypeId int, oprItemId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPrice(linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, trxTypeId int, jobTypeId int, variantId int, currencyId int, whsGroup string, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPriceByID(linetypeId int, oprItemId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPriceByCode(linetypeId int, oprItemCode string, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetVehicleUnitMaster(brandId int, modelId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetVehicleUnitByID(vehicleID int) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetVehicleUnitByChassisNumber(chassisNumber string) (map[string]interface{}, *exceptions.BaseErrorResponse)
	GetCampaignMaster(companyId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetLineTypeByItemCode(itemCode string) (string, *exceptions.BaseErrorResponse)
	GetOprItemPrice(linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode int, whsGroup string) (float64, *exceptions.BaseErrorResponse)
	ListItemLocation(companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	WarehouseGroupByCompany(companyId int) ([]masterpayloads.WarehouseGroupByCompanyResponse, *exceptions.BaseErrorResponse)
	ItemListTrans(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemListTransPL(companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddress(pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddressByID(customerId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddressByCode(customerCode string, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	WorkOrderService(pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	WorkOrderAtpmRegistration(pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ReferenceTypeWorkOrder(pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ReferenceTypeWorkOrderByID(referenceId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ReferenceTypeSalesOrder(pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ReferenceTypeSalesOrderByID(referenceId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetLineTypeByReferenceType(referenceTypeId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	LocationAvailable(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemDetailForItemInquiry(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemSubstituteDetailForItemInquiry(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetPartNumberItemImport(internalCondition []utils.FilterCondition, externalCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	LocationItem(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	LocationItemGoodsReceive(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemLocUOM(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemLocUOMById(companyId int, itemId int) (masterpayloads.ItemLocUOMResponse, *exceptions.BaseErrorResponse)
	ItemLocUOMByCode(companyId int, itemCode string) (masterpayloads.ItemLocUOMResponse, *exceptions.BaseErrorResponse)
	ItemMasterForFreeAccs(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemMasterForFreeAccsById(companyId int, itemId int) (masterpayloads.ItemMasterForFreeAccsResponse, *exceptions.BaseErrorResponse)
	ItemMasterForFreeAccsByCode(companyId int, itemCode string) (masterpayloads.ItemMasterForFreeAccsResponse, *exceptions.BaseErrorResponse)
	ItemMasterForFreeAccsByBrand(companyId int, itemId int, brandId int) (masterpayloads.ItemMasterForFreeAccsBrandResponse, *exceptions.BaseErrorResponse)
}
