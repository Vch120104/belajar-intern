package transactionsparepartrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"gorm.io/gorm"
	"net/http"
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
		ClaimSystemNumber:        payloads.ClaimSystemNumber,
		LocationItemId:           payloads.LocationId,
		ItemIds:                  goodsReceiveDetailEntities.ItemId,
		ItemUnitOfMeasurement:    goodsReceiveDetailEntities.ItemUnitOfMeasurement,
		ItemPrice:                goodsReceiveDetailEntities.ItemPrice,
		ItemDiscountPercentage:   goodsReceiveDetailEntities.ItemPrice,
		ItemDiscountAmount:       goodsReceiveDetailEntities.ItemDiscountAmount,
		QuantityDeliveryOrder:    goodsReceiveDetailEntities.QuantityDeliveryOrder,
		QuantityShort:            goodsReceiveDetailEntities.QuantityShort,
		QuantityDamaged:          goodsReceiveDetailEntities.QuantityDamage,
		QuantityOver:             goodsReceiveDetailEntities.QuantityOver,
		QuantityWrong:            goodsReceiveDetailEntities.QuantityWrong,
		QuantityClaimed:          goodsReceiveDetailEntities.QuantityShort + goodsReceiveDetailEntities.QuantityDamage,
		QuantityGoodsReceive:     goodsReceiveDetailEntities.QuantityGoodsReceive,
		Remark:                   payloads.Remark,
		CaseNumber:               goodsReceiveDetailEntities.CaseNumber,
		GoodsReceiveSystemNumber: payloads.GoodsReceiveSystemNumber,
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
