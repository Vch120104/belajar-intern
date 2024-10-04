package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type LookupRepository interface {
	GetLineTypeByItemCode(tx *gorm.DB, itemCode string) (int, *exceptions.BaseErrorResponse)
	GetOprItemPrice(tx *gorm.DB, linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode string, whsGroup string) (float64, *exceptions.BaseErrorResponse)
	GetOprItemDisc(tx *gorm.DB, lineTypeId int, billCode string, oprItemCode int, agreementId int, profitCenterId int, minValue float64, companyId int, brandId int, contractServSysNo int, whsGroup string, orderTypeId int) (float64, *exceptions.BaseErrorResponse)
	ItemOprCode(tx *gorm.DB, linetypeId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeByCode(tx *gorm.DB, linetypeId int, oprItemCode string, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeByID(tx *gorm.DB, linetypeId int, oprItemId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPrice(tx *gorm.DB, linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode string, whsGroup string, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	VehicleUnitMaster(tx *gorm.DB, brandId int, modelId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleUnitByChassisNumber(tx *gorm.DB, chassisNumber string, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetVehicleUnitByID(tx *gorm.DB, vehicleID int, paginate pagination.Pagination, filters []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CampaignMaster(tx *gorm.DB, companyId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddress(tx *gorm.DB, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddressByID(tx *gorm.DB, customerId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddressByCode(tx *gorm.DB, customerCode string, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	WorkOrderService(tx *gorm.DB, paginate pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse)
	GetWhsGroup(tx *gorm.DB, companyCode int) (string, *exceptions.BaseErrorResponse)
	GetCampaignDiscForWO(tx *gorm.DB, campaignId int, linetypeId int, oprItemCode string, frtQty float64, markupAmount float64, markupPercentage float64, millage float64) (masterpayloads.CampaignDiscount, *exceptions.BaseErrorResponse)
	GetItemLocationWarehouse(tx *gorm.DB, companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetWarehouseGroupByCompany(tx *gorm.DB, companyId int) ([]masterpayloads.WarehouseGroupByCompanyResponse, *exceptions.BaseErrorResponse)
}
