package masterrepository

import (
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"
	"time"

	"gorm.io/gorm"
)

type LookupRepository interface {
	ItemOprCode(tx *gorm.DB, linetypeId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemOprCodeByCode(tx *gorm.DB, linetypeId int, oprItemCode string, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemOprCodeByID(tx *gorm.DB, linetypeId int, oprItemId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPrice(tx *gorm.DB, linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, trxTypeId int, jobTypeId int, variantId int, currencyId int, whsGroup string, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPriceByID(tx *gorm.DB, linetypeId int, oprItemCode int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemOprCodeWithPriceByCode(tx *gorm.DB, linetypeId int, oprItemCode string, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetVehicleUnitMaster(tx *gorm.DB, brandId int, modelId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetVehicleUnitByChassisNumber(tx *gorm.DB, chassisNumber string, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetVehicleUnitByID(tx *gorm.DB, vehicleID int, paginate pagination.Pagination, filters []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetLineTypeByItemCode(tx *gorm.DB, itemCode string) (string, *exceptions.BaseErrorResponse)
	GetOprItemPrice(tx *gorm.DB, linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode int, whsGroup string) (float64, *exceptions.BaseErrorResponse)
	GetOprItemDisc(tx *gorm.DB, lineTypeStr string, billCodeId int, oprItemCode int, agreementId int, profitCenterId int, minValue float64, companyId int, brandId int, contractServSysNo int, whsGroup int, orderTypeId int) (float64, *exceptions.BaseErrorResponse)
	GetCampaignMaster(tx *gorm.DB, companyId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetWhsGroup(tx *gorm.DB, companyCode int) (int, *exceptions.BaseErrorResponse)
	GetCampaignDiscForWO(tx *gorm.DB, campaignId int, linetypeId int, oprItemId int, frtQty float64, markupAmount float64, markupPercentage float64, millage float64) (masterpayloads.CampaignDiscount, *exceptions.BaseErrorResponse)
	ListItemLocation(tx *gorm.DB, companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	WarehouseGroupByCompany(tx *gorm.DB, companyId int) ([]masterpayloads.WarehouseGroupByCompanyResponse, *exceptions.BaseErrorResponse)
	ItemListTrans(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemListTransPL(tx *gorm.DB, companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetOprItemFrt(tx *gorm.DB, oprItemId int, brandId int, modelId int, variantId int, vehicleChassisNo string) (float64, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddress(tx *gorm.DB, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddressByID(tx *gorm.DB, customerId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	CustomerByTypeAndAddressByCode(tx *gorm.DB, customerCode string, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	WorkOrderService(tx *gorm.DB, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	WorkOrderAtpmRegistration(tx *gorm.DB, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	SelectLocationStockItem(tx *gorm.DB, option int, companyId int, periodDate time.Time, whsCode int, locCode string, itemId int, whsGroup int, uomType string) (float64, *exceptions.BaseErrorResponse)
	ReferenceTypeWorkOrder(tx *gorm.DB, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ReferenceTypeWorkOrderByID(tx *gorm.DB, referenceId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ReferenceTypeSalesOrder(tx *gorm.DB, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ReferenceTypeSalesOrderByID(tx *gorm.DB, referenceId int, paginate pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetLineTypeByReferenceType(tx *gorm.DB, referenceTypeId int, paginate pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	LocationAvailable(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemDetailForItemInquiry(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemSubstituteDetailForItemInquiry(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	GetPartNumberItemImport(tx *gorm.DB, internalCondition []utils.FilterCondition, externalCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	LocationItem(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemLocUOM(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemLocUOMById(tx *gorm.DB, companyId int, itemId int) (masterpayloads.ItemLocUOMResponse, *exceptions.BaseErrorResponse)
	ItemLocUOMByCode(tx *gorm.DB, companyId int, itemCode string) (masterpayloads.ItemLocUOMResponse, *exceptions.BaseErrorResponse)
	ItemMasterForFreeAccs(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse)
	ItemMasterForFreeAccsById(tx *gorm.DB, companyId int, itemId int) (masterpayloads.ItemMasterForFreeAccsResponse, *exceptions.BaseErrorResponse)
	ItemMasterForFreeAccsByCode(tx *gorm.DB, companyId int, itemCode string) (masterpayloads.ItemMasterForFreeAccsResponse, *exceptions.BaseErrorResponse)
	ItemMasterForFreeAccsByBrand(tx *gorm.DB, companyId int, itemId int, brandId int) (masterpayloads.ItemMasterForFreeAccsBrandResponse, *exceptions.BaseErrorResponse)
}
