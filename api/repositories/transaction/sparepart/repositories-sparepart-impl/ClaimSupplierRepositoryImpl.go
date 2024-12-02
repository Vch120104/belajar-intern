package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ClaimSupplierRepositoryImpl struct {
}

func NewClaimSupplierRepositoryImpl() transactionsparepartrepository.ClaimSupplierRepository {
	return &ClaimSupplierRepositoryImpl{}
}

// InsertItemClaim insert item claim header
// uspg_atItemClaim1_Update option 1
func (c *ClaimSupplierRepositoryImpl) InsertItemClaim(db *gorm.DB, payloads transactionsparepartpayloads.ClaimSupplierInsertPayload) (transactionsparepartentities.ItemClaim, *exceptions.BaseErrorResponse) {
	//get goods receive header for insert item claim
	var goodsReceiveEntities transactionsparepartentities.GoodsReceive
	var itemClaimEntities transactionsparepartentities.ItemClaim
	err := db.Model(&goodsReceiveEntities).
		Where(transactionsparepartentities.GoodsReceive{GoodsReceiveSystemNumber: payloads.GoodsReceiveSystemNumber}).
		First(&goodsReceiveEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return itemClaimEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "goods receive please check input",
			}
		}
		return itemClaimEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get goods receive header please check input",
		}
	}
	//get document ready
	docStatusReady, errDoc := generalserviceapiutils.GetDocumentStatusByCode("20")
	if errDoc != nil {
		return itemClaimEntities, errDoc
	}
	//begin insert item claim entities
	itemClaimEntities.ClaimTypeId = payloads.ClaimTypeId
	itemClaimEntities.CompanyId = payloads.CompanyId
	itemClaimEntities.GoodsReceiveSystemNumber = payloads.GoodsReceiveSystemNumber
	itemClaimEntities.GoodsReceiveDocumentNumber = goodsReceiveEntities.GoodsReceiveDocumentNumber
	itemClaimEntities.VehicleBrandId = goodsReceiveEntities.BrandId
	itemClaimEntities.ProfitCenterId = goodsReceiveEntities.ProfitCenterId
	itemClaimEntities.TransactionTypeId = payloads.TransactionTypeId
	itemClaimEntities.EventId = payloads.EventId
	itemClaimEntities.ReferenceTypeGoodReceiveId = goodsReceiveEntities.ReferenceTypeGoodReceiveId
	itemClaimEntities.ReferenceSystemNumber = goodsReceiveEntities.ReferenceSystemNumber
	itemClaimEntities.ReferenceDocumentNumber = goodsReceiveEntities.ReferenceDocumentNumber
	itemClaimEntities.SupplierId = goodsReceiveEntities.SupplierId
	itemClaimEntities.SuppplierDoNumber = goodsReceiveEntities.SupplierDeliveryOrderNumber
	itemClaimEntities.WarehouseGroupId = goodsReceiveEntities.WarehouseGroupId
	itemClaimEntities.WarehouseId = goodsReceiveEntities.WarehouseId
	itemClaimEntities.ViaBinning = goodsReceiveEntities.ViaBinning
	itemClaimEntities.CurrencyId = goodsReceiveEntities.CurrencyId
	itemClaimEntities.CurrencyExchangeRate = goodsReceiveEntities.CurrencyExchangeRate
	itemClaimEntities.CurrencyExchangeRateDate = time.Now()
	itemClaimEntities.CurrencyRateType = goodsReceiveEntities.CurrencyExchangeRateTypeId
	itemClaimEntities.ItemGroupId = goodsReceiveEntities.ItemGroupId
	//itemClaimEntities.ClaimStatus = "20"
	itemClaimEntities.ClaimStatusId = docStatusReady.DocumentStatusId
	err = db.Create(&itemClaimEntities).Error
	if err != nil {
		return itemClaimEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to save item claim supplier please check input",
		}
	}
	return itemClaimEntities, nil
}
func (c *ClaimSupplierRepositoryImpl) InsertItemClaimDetail(db *gorm.DB, payloads transactionsparepartpayloads.ClaimSupplierInsertDetailPayload) (transactionsparepartentities.ItemClaimDetail, *exceptions.BaseErrorResponse) {
	//get header goods receive for input data
	var goodsReceiveEntities transactionsparepartentities.GoodsReceive
	var goodsReceiveDetailEntities transactionsparepartentities.GoodsReceiveDetail
	var itemClaimDetailEntities transactionsparepartentities.ItemClaimDetail
	err := db.Model(&goodsReceiveEntities).Where(transactionsparepartentities.GoodsReceive{GoodsReceiveSystemNumber: payloads.GoodsReceiveSystemNumber}).
		First(&goodsReceiveEntities).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return itemClaimDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "goods receive is not found please check input",
			}
		}
		return itemClaimDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get goods receive header please check input",
		}
	}
	//get detail
	err = db.Model(&goodsReceiveDetailEntities).Where(transactionsparepartentities.GoodsReceiveDetail{GoodsReceiveSystemNumber: payloads.GoodsReceiveSystemNumber}).
		First(&goodsReceiveDetailEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return itemClaimDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "goods receive detail is not found please check input",
			}
		}
		return itemClaimDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get goods receive detail please check input",
		}
	}
	itemClaimDetailEntities = transactionsparepartentities.ItemClaimDetail{
		ClaimSystemNumber:              payloads.ClaimSystemNumber,
		LocationItemId:                 payloads.LocationId,
		ItemIds:                        goodsReceiveDetailEntities.ItemId,
		ItemUnitOfMeasurement:          goodsReceiveDetailEntities.ItemUnitOfMeasurement,
		ItemPrice:                      goodsReceiveDetailEntities.ItemPrice,
		ItemDiscountPercentage:         goodsReceiveDetailEntities.ItemDiscountPercent,
		ItemDiscountAmount:             goodsReceiveDetailEntities.ItemDiscountAmount,
		QuantityDeliveryOrder:          goodsReceiveDetailEntities.QuantityDeliveryOrder,
		QuantityShort:                  goodsReceiveDetailEntities.QuantityShort,
		QuantityDamaged:                goodsReceiveDetailEntities.QuantityDamage,
		QuantityOver:                   goodsReceiveDetailEntities.QuantityOver,
		QuantityWrong:                  goodsReceiveDetailEntities.QuantityWrong,
		QuantityClaimed:                goodsReceiveDetailEntities.QuantityShort + goodsReceiveDetailEntities.QuantityDamage,
		QuantityGoodsReceive:           goodsReceiveDetailEntities.QuantityGoodsReceive,
		Remark:                         payloads.Remark,
		CaseNumber:                     goodsReceiveDetailEntities.CaseNumber,
		GoodsReceiveDetailSystemNumber: payloads.GoodsReceiveDetailSystemNumber,
		//GoodsReceiveLineNumber:   0,
		//QuantityBinning:          0,
		ItemTotal:           (goodsReceiveDetailEntities.QuantityDamage + goodsReceiveDetailEntities.QuantityOver + goodsReceiveDetailEntities.QuantityWrong) * goodsReceiveDetailEntities.ItemPrice,
		ItemTotalBaseAmount: (goodsReceiveDetailEntities.QuantityDamage + goodsReceiveDetailEntities.QuantityOver + goodsReceiveDetailEntities.QuantityWrong) * goodsReceiveDetailEntities.ItemPrice * goodsReceiveEntities.CurrencyExchangeRate,
	}
	err = db.Create(&itemClaimDetailEntities).First(&itemClaimDetailEntities).Error
	if err != nil {
		return itemClaimDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to create item detail please check input",
		}
	}
	return itemClaimDetailEntities, nil
}
func (c *ClaimSupplierRepositoryImpl) GetItemClaimById(db *gorm.DB, itemClaimId int) (transactionsparepartpayloads.ClaimSupplierGetByIdResponse, *exceptions.BaseErrorResponse) {
	//select
	var itemClaim transactionsparepartentities.ItemClaim
	var itemClaimReponse transactionsparepartpayloads.ClaimSupplierGetByIdResponse
	err := db.Model(&itemClaim).Where(transactionsparepartentities.ItemClaim{ClaimSystemNumber: itemClaimId}).First(&itemClaim).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return itemClaimReponse, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item claim with that id is not found please check input",
				Err:        err,
			}
		}
		return itemClaimReponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get item claim please check input",
		}
	}
	//get status type first
	docStatus, errDocStatus := generalserviceapiutils.GetDocumentStatusById(itemClaim.ClaimStatusId)
	if errDocStatus != nil {
		return itemClaimReponse, errDocStatus
	}
	//get supplier
	supplierData, errSupplier := generalserviceapiutils.GetSupplierMasterByID(itemClaim.SupplierId)
	if errSupplier != nil {
		return itemClaimReponse, errSupplier
	}
	//get warehouse group data first
	var warehouseGroupData masterwarehouseentities.WarehouseGroup
	err = db.Model(&warehouseGroupData).Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupId: itemClaim.WarehouseGroupId}).
		First(&warehouseGroupData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return itemClaimReponse, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse group with that id is not found please check input",
				Err:        err,
			}
		}
		return itemClaimReponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "Failed to get warehouse group please check input",
		}
	}
	//get warehouse masrter
	var warehouseMasterData masterwarehouseentities.WarehouseMaster
	err = db.Model(&warehouseMasterData).
		Where(masterwarehouseentities.WarehouseMaster{WarehouseGroupId: itemClaim.WarehouseGroupId}).
		First(&warehouseMasterData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return itemClaimReponse, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse with that id is not found please check input",
			}
		}
		return itemClaimReponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "Failed to get warehouse id please check input",
		}
	}
	//get company
	companyData, errCompanyData := generalserviceapiutils.GetCompanyDataById(itemClaim.CompanyId)
	if errCompanyData != nil {
		return itemClaimReponse, errCompanyData
	}

	//vehicle
	vehicleMaster, errVehicleId := salesserviceapiutils.GetUnitBrandById(itemClaim.VehicleBrandId)
	if errVehicleId != nil {
		return itemClaimReponse, errVehicleId
	}
	//get goods receive reference type status
	var goodsReceiveReferenceType masterentities.GoodsReceiveReferenceType
	err = db.Model(&goodsReceiveReferenceType).Where(masterentities.GoodsReceiveReferenceType{ReferenceTypeGoodReceiveId: itemClaim.ReferenceTypeGoodReceiveId}).
		First(&goodsReceiveReferenceType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return itemClaimReponse, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "goods receive reference type is not found please check input",
				Err:        err,
			}
		}
		return itemClaimReponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "Failed to get goods receive reference type please check input",
		}
	}
	//get item claim response
	claimType, claimTypeErr := generalserviceapiutils.GetItemClaimTypeMasterById(itemClaim.ClaimTypeId)
	if claimTypeErr != nil {
		return itemClaimReponse, claimTypeErr
	}
	//get transaction type from general
	transType, transTypeErr := generalserviceapiutils.GetTransactionTypeById(itemClaim.TransactionTypeId)
	if transTypeErr != nil {
		return itemClaimReponse, transTypeErr
	}
	itemClaimReponse = transactionsparepartpayloads.ClaimSupplierGetByIdResponse{
		CompanyId:                     companyData.CompanyId,
		CompanyCode:                   companyData.CompanyCode,
		CompanyName:                   companyData.CompanyName,
		ClaimSystemNumber:             itemClaim.ClaimSystemNumber,
		ClaimStatusId:                 itemClaim.ClaimStatusId,
		ClaimStatusDesc:               docStatus.DocumentStatusDescription,
		ClaimDocumentNumber:           itemClaim.ClaimDocumentNumber,
		ClaimDate:                     itemClaim.ClaimDate,
		ClaimType:                     claimType.ItemClaimTypeDescription,
		ClaimTypeId:                   itemClaim.ClaimTypeId,
		GoodsReceiveSystemNumber:      itemClaim.GoodsReceiveSystemNumber,
		GoodsReceiveDocumentNumber:    itemClaim.GoodsReceiveDocumentNumber,
		VehicleId:                     itemClaim.VehicleBrandId,
		VehicleBrand:                  vehicleMaster.BrandName,
		CostCenterId:                  itemClaim.CostCenterId,
		ProfitCenterId:                itemClaim.ProfitCenterId,
		TransactionTypeId:             itemClaim.TransactionTypeId,
		TransactionTypeDesc:           transType.TransactionTypeName,
		EventId:                       itemClaim.EventId,
		SupplierId:                    itemClaim.SupplierId,
		SupplierName:                  supplierData.SupplierName,
		SupplierDoNo:                  itemClaim.SuppplierDoNumber,
		ReferenceTypeGoodsReceiveId:   itemClaim.ReferenceTypeGoodReceiveId,
		ReferenceTypeGoodsReceiveDesc: goodsReceiveReferenceType.ReferenceTypeGoodsReceiveDescription,
		ReferenceSystemNumber:         itemClaim.ReferenceSystemNumber,
		ReferenceDocumentNumber:       itemClaim.ReferenceDocumentNumber,
		WarehouseGroupId:              itemClaim.WarehouseGroupId,
		WarehouseGroupCode:            warehouseGroupData.WarehouseGroupCode,
		WarehouseId:                   warehouseMasterData.WarehouseId,
		WarehouseCode:                 warehouseMasterData.WarehouseCode,
		WarehouseName:                 warehouseMasterData.WarehouseName,
	}
	return itemClaimReponse, nil
}
func (c *ClaimSupplierRepositoryImpl) GetItemClaimDetailByHeaderId(db *gorm.DB, Paginations pagination.Pagination, filter []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//return Paginations, nil
	entities := transactionsparepartentities.ItemClaimDetail{}
	var response []transactionsparepartpayloads.ClaimSupplierGetAllDetailResponse
	JoinTable := db.Model(&entities).Select("*")
	whereQuery := utils.ApplyFilter(JoinTable, filter)
	err := whereQuery.Scopes(pagination.Paginate(&entities, &Paginations, whereQuery)).Scan(&response).Error
	if err != nil {
		return Paginations, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get item claim detail by header id",
			Err:        err,
		}
	}
	if len(response) == 0 {
		Paginations.Rows = []string{}
		return Paginations, nil
	}
	for i, item := range response {
		itemEntities := masteritementities.Item{}
		err = db.Model(&itemEntities).Where(masteritementities.Item{ItemId: item.ItemId}).Scan(&itemEntities).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return Paginations, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to get item data",
				Err:        err,
			}
		}
		response[i].ItemName = itemEntities.ItemName
	}
	Paginations.Rows = response
	return Paginations, nil
}
func (c *ClaimSupplierRepositoryImpl) SubmitItemClaim(db *gorm.DB, claimId int) (bool, *exceptions.BaseErrorResponse) {
	//get claim header
	itemClaimEntities := transactionsparepartentities.ItemClaim{}

	err := db.Model(&itemClaimEntities).Where(transactionsparepartentities.ItemClaim{ClaimSystemNumber: claimId}).First(&itemClaimEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "item claim header is not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get item claim",
			Err:        err,
		}
	}
	periodCompany, periodCompanyErr := financeserviceapiutils.GetOpenPeriodByCompany(itemClaimEntities.CompanyId, "SP")
	if periodCompanyErr != nil {
		return false, periodCompanyErr
	}
	docMonth := strconv.Itoa(int(itemClaimEntities.ClaimDate.Month()))
	docYear := strconv.Itoa(itemClaimEntities.ClaimDate.Year())

	if periodCompany.PeriodMonth != docMonth && periodCompany.PeriodYear != docYear {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("period Status For This Claim Supplier Is Already Close"),
			Message:    "Period Status For This Claim Supplier Is Already Close",
		}
	}
	docNumber, docNumberErr := GenerateDocumentNumber(db, claimId)
	if docNumberErr != nil {
		return false, docNumberErr
	}
	itemClaimEntities.ClaimDocumentNumber = docNumber
	var itemClaimDetailResponds []transactionsparepartpayloads.ClaimSupplierDetailSubmitCursor
	err = db.Table("trx_item_claim").
		Select(`A.location_item_id,
					a.item_id,
					A.item_unit_of_measurement,
					A.item_total_base_amount,
					A.quantity_over + a.quantity_damaged + a.quantity_wrong AS quantity_variance,
					A.goods_receive_detail_system_number`).Where("claim_system_number = ?", claimId).
		Scan(&itemClaimDetailResponds).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "error when getting entities detail",
		}
	}
	for _, item := range itemClaimDetailResponds {
		if item.QuantityVariance != 0 {
			//<<localhostp8000>>unit-of-measurement/get_quantity_conversion?source_type=S&item_id=893891&quantity=1.0
			sourceTypeConversionResponse := masteritempayloads.SourceTypeConversionResponse{}
			sourceTypeConversionUrl := config.EnvConfigs.AfterSalesServiceUrl + fmt.Sprintf("unit-of-measurement/get_quantity_conversion?source_type=%s&item_id=%s&quantity=%s")
			err = utils.CallAPI("GET", sourceTypeConversionUrl, nil, &sourceTypeConversionResponse)
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
					Message:    "failed to get source type conversion",
				}
			}
			if sourceTypeConversionResponse.QuantityConversion == 0 {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        err,
					Message:    "UOM does not have conversion, please define UOM Conversion",
				}
			}

			//IF EXISTS (SELECT * FROM	amLocationStock WHERE	PERIOD_YEAR = @Period_Year AND PERIOD_MONTH = @Period_Month AND
			//COMPANY_CODE = @Company_Code AND ITEM_CODE = @CSR_Item_Code AND
			//WHS_CODE = @Whs_Code AND LOC_CODE = @CSR_Loc_Code )
			isExist := false
			err = db.Model(&masterentities.LocationStock{}).Where(
				masterentities.LocationStock{
					PeriodYear:  docYear,
					PeriodMonth: docMonth,
					CompanyId:   itemClaimEntities.CompanyId,
					ItemId:      item.ItemId,
					WarehouseId: itemClaimEntities.WarehouseId,
					LocationId:  item.LocationItemId,
				}).Select("1").Scan(&isExist).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
					Message:    "failed on getting location stock please check input",
				}
			}
			if isExist {
				//SET @Qty_OnHand = ISNULL((SELECT QTY_ON_HAND FROM dbo.viewLocationStock WHERE  PERIOD_YEAR = @Period_Year  AND PERIOD_MONTH = @Period_Month
				//AND COMPANY_CODE = @Company_Code AND ITEM_CODE = @CSR_Item_Code
				//AND	WHS_CODE = @Whs_Code AND LOC_CODE = @CSR_Loc_Code ),0)
				quantityOnHand := 0.0
				err = db.Table("mtr_location_stock").
					Select(`
			item_inquiry_id,
			period_year,
			period_month,
			company_id,
			warehouse_group_id,
			warehouse_id,
			location_id,
			(
				ISNULL(quantity_sales, 0) +
				ISNULL(quantity_transfer_out, 0) +
				ISNULL(quantity_claim_out, 0) +
				ISNULL(quantity_robbing_out, 0) +
				ISNULL(quantity_assembly_out, 0)
			) AS quantity_on_hand
		`).
					Joins("LEFT JOIN mtr_warehouse_master mwm ON mwm.company_id = mls.company_id AND mwm.warehouse_id = mls.warehouse_id").
					Where(`period_year = ? AND
								 period_month = ? AND
								 company_id = ? AND
								 item_id = ? AND
								warehouse_id = ? AND
								location_id = ?
							`, docYear, docMonth, itemClaimEntities.CompanyId, item.ItemId, itemClaimEntities.WarehouseId, item.LocationItemId).
					Scan(&quantityOnHand).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Err:        err,
						Message:    "failed to get quantity on hand from view location stock",
					}
				}
				var hppCurent float64
				if quantityOnHand != 0 && !itemClaimEntities.ViaBinning {
					isExist = false
					err = db.Model(&masterentities.GroupStock{}).Where(&masterentities.GroupStock{CompanyId: itemClaimEntities.CompanyId,
						PeriodMonth:      docMonth,
						PeriodYear:       docYear,
						WarehouseGroupId: itemClaimEntities.WarehouseGroupId,
						ItemId:           item.ItemId,
					}).Select("1").Scan(&isExist).Error
					if err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Err:        err,
							Message:    "failed to check group stock",
						}
					}
					if isExist {
						err = db.Model(&masterentities.GroupStock{}).Where(&masterentities.GroupStock{CompanyId: itemClaimEntities.CompanyId,
							PeriodMonth:      docMonth,
							PeriodYear:       docYear,
							WarehouseGroupId: itemClaimEntities.WarehouseGroupId,
							ItemId:           item.ItemId,
						}).Select("price_current").Scan(&hppCurent).Error
					} else {
						hppCurent = 0
					}
				} else {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Err:        errors.New("data stock is not valid"),
						Message:    "data stock is not valid",
					}
				}
			} else {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        errors.New("data stock is not valid"),
					Message:    "data stock is not valid",
				}
			}
		}
	}
	docComplete, docCompleteErr := generalserviceapiutils.GetDocumentStatusByCode("99")
	if docCompleteErr != nil {
		return false, docCompleteErr
	}
	itemClaimEntities.ClaimStatusId = docComplete.DocumentStatusId
	itemClaimEntities.ClaimDocumentNumber = docNumber
	err = db.Save(&itemClaimEntities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to save claim document",
		}
	}
	return true, nil

}
func (p *ClaimSupplierRepositoryImpl) GenerateDocumentNumber(tx *gorm.DB, id int) (string, *exceptions.BaseErrorResponse) {
	var workOrder transactionsparepartentities.PurchaseRequestEntities

	// Get the work order based on the work order system number
	err := tx.Model(&transactionsparepartentities.PurchaseRequestEntities{}).Where("purchase_request_system_number = ?", id).First(&workOrder).Error
	if err != nil {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve work order from the database: %v", err)}
	}

	if workOrder.BrandId == 0 {

		return "", &exceptions.BaseErrorResponse{Message: "brand_id is missing in the work order. Please ensure the work order has a valid brand_id before generating document number."}
	}

	// Get the last work order based on the work order system number
	var lastWorkOrder transactionsparepartentities.PurchaseRequestEntities
	err = tx.Model(&transactionsparepartentities.PurchaseRequestEntities{}).
		Where("brand_id = ?", workOrder.BrandId).
		Order("purchase_request_document_number desc").
		First(&lastWorkOrder).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve last work order: %v", err)}
	}

	currentTime := time.Now()
	month := int(currentTime.Month())
	year := currentTime.Year() % 100 // Use last two digits of the year

	// fetch data brand from external api
	brandResponse, brandErr := generalserviceapiutils.GetBrandGenerateDoc(workOrder.BrandId)
	if brandErr != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch brand data from external service",
			Err:        brandErr.Err,
		}
	}

	// Check if BrandCode is not empty before using it
	if brandResponse.BrandCode == "" {
		return "", &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Message: "Brand code is empty"}
	}

	// Get the initial of the brand code
	brandInitial := brandResponse.BrandCode[0]

	// Handle the case when there is no last work order or the format is invalid
	newDocumentNumber := fmt.Sprintf("SPPR/%c/%02d/%02d/00001", brandInitial, month, year)
	if lastWorkOrder.PurchaseRequestSystemNumber != 0 {
		lastWorkOrderDate := lastWorkOrder.PurchaseRequestDocumentDate
		lastWorkOrderYear := lastWorkOrderDate.Year() % 100

		// Check if the last work order is from the same year
		if lastWorkOrderYear == year {
			lastWorkOrderCode := lastWorkOrder.PurchaseRequestDocumentNumber
			codeParts := strings.Split(lastWorkOrderCode, "/")
			if len(codeParts) == 5 {
				lastWorkOrderNumber, err := strconv.Atoi(codeParts[4])
				if err == nil {
					newWorkOrderNumber := lastWorkOrderNumber + 1
					newDocumentNumber = fmt.Sprintf("SPPR/%c/%02d/%02d/%05d", brandInitial, month, year, newWorkOrderNumber)
				} else {
					log.Printf("Failed to parse last work order code: %v", err)
				}
			} else {
				log.Println("Invalid last work order code format")
			}
		}
	}

	log.Printf("New document number: %s", newDocumentNumber)
	return newDocumentNumber, nil
}
func (c *ClaimSupplierRepositoryImpl) GetAllItemClaim(db *gorm.DB, page pagination.Pagination, filter []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var response []transactionsparepartpayloads.ClaimSupplierGetAllResponds
	entities := transactionsparepartentities.ItemClaim{}

	joinTable := db.Model(&entities).Select(`
				claim_system_number,
				claim_document_number,
				goods_receive_document_number,
				reference_document_number,
				supplier_id,
				claim_status_id`)
	//WhereQuery := utils.ApplyFilter(joinTable, filter)
	//err:=WhereQuery.Scan(pagination.Paginate())
}
