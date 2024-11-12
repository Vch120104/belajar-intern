package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	financeservice "after-sales/api/payloads/cross-service/finance-service"
	generalservicepayloads "after-sales/api/payloads/cross-service/general-service"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type GoodsReceiveRepositoryImpl struct {
}

func NewGoodsReceiveRepositoryImpl() transactionsparepartrepository.GoodsReceiveRepository {
	return &GoodsReceiveRepositoryImpl{}
}

func (repository *GoodsReceiveRepositoryImpl) GetAllGoodsReceive(db *gorm.DB, filter []utils.FilterCondition, paginations pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []transactionsparepartpayloads.GoodsReceivesGetAllPayloads
	Entities := transactionsparepartentities.GoodsReceive{}
	JoinTable := db.Table("trx_goods_receive IG").
		Joins(`LEFT OUTER JOIN trx_goods_receive_detail IG1 ON IG.goods_receive_system_number = ig1.goods_receive_system_number`).
		Joins("LEFT OUTER JOIN mtr_item_group itemgroup ON IG.item_group_id = itemgroup.item_group_id").
		Joins(`INNER JOIN mtr_reference_type_goods_receive reftype ON reftype.reference_type_good_receive_id = ig.reference_type_good_receive_id`).
		Select(`
						ig.goods_receive_system_number,
						ig.goods_receive_document_number,
						itemgroup.item_group_name,
						ig.goods_receive_document_date,
						ig.reference_document_number,
						ig.supplier_id,
						ig.goods_receive_status_id,
						ig.journal_system_number,
						ig.supplier_delivery_order_number,
						SUM(ISNULL(ig1.quantity_goods_receive,0)) AS quantity_goods_receive,
						SUM(ISNULL(ig1.quantity_goods_receive,0) * ISNULL(ig1.item_price,0)) AS total_amount
					`).
		Group(`	ig.goods_receive_system_number,
						ig.goods_receive_document_number,
						itemgroup.item_group_name,
						ig.goods_receive_document_date,
						ig.reference_document_number,
						ig.supplier_id,
						ig.goods_receive_status_id,
						ig.journal_system_number,
						ig.supplier_delivery_order_number`)
	WhereQuery := utils.ApplyFilter(JoinTable, filter)
	//for i, res := range responses {
	//	var SupplierData generalservicepayloads.SupplierMasterCrossServicePayloads
	//	SupplierDataUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(res.SupplierId)
	//	if err := utils.Get(SupplierDataUrl, &SupplierData, nil); err != nil {
	//		return paginations, &exceptions.BaseErrorResponse{
	//			StatusCode: http.StatusInternalServerError,
	//			Message:    "Failed to fetch Supplier data from external service" + err.Error(),
	//			Err:        err,
	//		}
	//	}
	//responses[i].SupplierName = SupplierData.SupplierName
	//}
	err := WhereQuery.Scopes(pagination.Paginate(&Entities, &paginations, WhereQuery)).Scan(&responses).Error
	if err != nil {
		return paginations, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error On Paginate Goods Receive",
		}
	}
	paginations.Rows = responses
	return paginations, nil
}
func (repository *GoodsReceiveRepositoryImpl) GetGoodsReceiveById(db *gorm.DB, GoodsReceiveId int) (transactionsparepartpayloads.GoodsReceivesGetByIdResponses, *exceptions.BaseErrorResponse) {
	var response transactionsparepartpayloads.GoodsReceivesGetByIdResponses
	////Entitites := transactionsparepartentities.GoodsReceive{}
	//ressss := 0
	//err := db.Model(&transactionsparepartentities.GoodsReceive{}).Where(transactionsparepartentities.GoodsReceive{GoodsReceiveSystemNumber: GoodsReceiveId}).
	//	Select("1").Scan(&ressss).Error
	//
	////First(Entitites).
	////	Error
	//if err == gorm.ErrRecordNotFound {
	//	return response, nil
	//
	//}
	//return response, nil
	err := db.Table("trx_goods_receive A").
		Joins("LEFT OUTER JOIN mtr_warehouse_master D ON D.warehouse_id = A.warehouse_id AND A.company_id = D.company_id").
		Joins("LEFT OUTER JOIN mtr_warehouse_master E ON E.warehouse_id = A.warehouse_id AND E.company_id = A.company_id").
		Select(`
	 A.goods_receive_system_number,
       A.goods_receive_status_id,
       A.goods_receive_document_number,
       A.goods_receive_document_date,
       A.item_group_id,
       a.reference_type_good_receive_id,
       a.reference_system_number,
       a.reference_document_number,
       a.affiliated_purchase_order,
       a.via_binning,
       a.back_order,
       a.set_order,
       a.brand_id,
       a.cost_center_id,
       a.profit_center_id,
       a.transaction_type_id,
       a.event_id,
       a.supplier_id,
       a.supplier_delivery_order_number,
       a.supplier_invoice_number,
       a.supplier_tax_invoice_number,
       a.warehouse_id,
       D.warehouse_code,
       D.warehouse_name,
       a.warehouse_claim_id,
       E.warehouse_code AS warehouse_claim_code,
       E.warehouse_name AS warehouse_claim_name,
       a.item_class_id
		`).
		Where("A.goods_receive_system_number = ?", GoodsReceiveId).
		Scan(&response).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("header Not Found"),
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error On Get Goods Receive By Id",
		}
	}
	return response, nil
}
func (repository *GoodsReceiveRepositoryImpl) InsertGoodsReceive(db *gorm.DB, payloads transactionsparepartpayloads.GoodsReceiveInsertPayloads) (transactionsparepartentities.GoodsReceive, *exceptions.BaseErrorResponse) {
	Entities := transactionsparepartentities.GoodsReceive{}
	if payloads.CompanyId == 0 {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New(`company Is Missing Please Try Again`),
		}
	}
	if payloads.ReferenceSystemNumber == 0 {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("reference Document is missing. Please try again"),
		}
	}
	var PurchaseOrderWarehouseId int

	err := db.Model(&transactionsparepartentities.PurchaseOrderEntities{}).
		Select("warehouse_id").
		Where(transactionsparepartentities.PurchaseOrderEntities{CompanyId: payloads.CompanyId,
			PurchaseOrderDocumentNumber: payloads.ReferenceDocumentNumber,
			PurchaseOrderSystemNumber:   payloads.ReferenceSystemNumber,
			ItemGroupId:                 payloads.ItemGroupId,
		}).Scan(&PurchaseOrderWarehouseId).
		Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return Entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error On Get Purchase OrderReference number",
			}
		}
	}
	//if PurchaseOrderWarehouseId != payloads.WarehouseId {
	//	return Entities, &exceptions.BaseErrorResponse{
	//		StatusCode: http.StatusBadRequest,
	//		Err:        errors.New("the Goods Receipt Warehouse is not the same as the Purchase Order or Binning Warehouse . Please try again"),
	//	}
	//}
	Entities = transactionsparepartentities.GoodsReceive{
		CompanyId:                   payloads.CompanyId,
		GoodsReceiveDocumentNumber:  payloads.ReferenceDocumentNumber,
		GoodsReceiveDocumentDate:    payloads.GoodsReceiveDocumentDate,
		GoodsReceiveStatusId:        payloads.GoodsReceiveStatusId,
		ReferenceTypeGoodReceiveId:  payloads.ReferenceTypeGoodReceiveId,
		ReferenceSystemNumber:       payloads.ReferenceSystemNumber,
		ReferenceDocumentNumber:     payloads.ReferenceDocumentNumber,
		AffiliatedPurchaseOrder:     payloads.AffiliatedPurchaseOrder,
		ViaBinning:                  payloads.ViaBinning,
		SetOrder:                    payloads.SetOrder,
		BackOrder:                   payloads.BackOrder,
		BrandId:                     payloads.BrandId,
		CostCenterId:                payloads.CostCenterId,
		ProfitCenterId:              payloads.ProfitCenterId,
		TransactionTypeId:           payloads.TransactionTypeId,
		EventId:                     payloads.EventId,
		SupplierId:                  payloads.SupplierId,
		SupplierDeliveryOrderNumber: payloads.SupplierDeliveryOrderNumber,
		SupplierInvoiceNumber:       payloads.SupplierInvoiceNumber,
		SupplierInvoiceDate:         payloads.SupplierInvoiceDate,
		SupplierTaxInvoiceNumber:    payloads.SupplierTaxInvoiceNumber,
		SupplierTaxInvoiceDate:      payloads.SupplierTaxInvoiceDate,
		WarehouseGroupId:            payloads.WarehouseGroupId,
		WarehouseId:                 payloads.WarehouseId,
		WarehouseClaimId:            payloads.WarehouseClaimId,
		ItemGroupId:                 payloads.ItemGroupId,
		CurrencyId:                  payloads.CurrencyId,
		CurrencyExchangeRateDate:    time.Now().UTC(),
		CurrencyExchangeRate:        payloads.CurrencyExchangeRate,
		CurrencyExchangeRateTypeId:  payloads.CurrencyExchangeRateTypeId,
		UseInTransitWarehouse:       payloads.UseInTransitWarehouse,
		InTransitWarehouseId:        payloads.InTransitWarehouseId,
		ChangeNo:                    1,
		CreatedDate:                 time.Now().UTC(),
		UpdatedDate:                 time.Now().UTC(),
		CreatedByUserId:             payloads.CreatedByUserId,
		UpdatedByUserId:             payloads.UpdatedByUserId,
	}
	err = db.Create(&Entities).Error
	if err != nil {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed To Create Record Error On:" + err.Error(),
		}
	}
	return Entities, nil
}

// [uspg_atItemGRPO0_Update]
// option 0
func (repository *GoodsReceiveRepositoryImpl) UpdateGoodsReceive(db *gorm.DB, payloads transactionsparepartpayloads.GoodsReceiveUpdatePayloads, GoodsReceivesId int) (transactionsparepartentities.GoodsReceive, *exceptions.BaseErrorResponse) {
	Entities := transactionsparepartentities.GoodsReceive{}
	err := db.Model(&Entities).Where(transactionsparepartentities.GoodsReceive{GoodsReceiveSystemNumber: GoodsReceivesId}).
		First(&Entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Goods Receive Is Not Found",
				Err:        errors.New("goods Receive Is Not Found"),
			}
		}
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed To Update GoodsReceive",
		}
	}
	Entities.ReferenceSystemNumber = payloads.ReferenceSystemNumber
	Entities.ReferenceTypeGoodReceiveId = payloads.ReferenceTypeGoodReceiveId
	Entities.ReferenceDocumentNumber = payloads.ReferenceDocumentNumber
	Entities.ProfitCenterId = payloads.ProfitCenterId
	Entities.TransactionTypeId = payloads.TransactionTypeId
	Entities.EventId = payloads.EventId
	Entities.AffiliatedPurchaseOrder = payloads.AffiliatedPurchaseOrder
	Entities.ViaBinning = payloads.ViaBinning
	Entities.SupplierId = payloads.SupplierId
	Entities.SupplierDeliveryOrderNumber = payloads.SupplierDeliveryOrderNumber
	Entities.SupplierInvoiceNumber = payloads.SupplierInvoiceNumber
	Entities.SupplierTaxInvoiceNumber = payloads.SupplierTaxInvoiceNumber
	Entities.WarehouseGroupId = payloads.WarehouseGroupId
	Entities.WarehouseId = payloads.WarehouseId
	Entities.WarehouseClaimId = payloads.WarehouseClaimId
	Entities.ItemGroupId = payloads.ItemGroupId
	Entities.ChangeNo += 1
	Entities.UpdatedDate = time.Now().UTC()
	Entities.UpdatedByUserId = payloads.UpdatedByUserId
	Entities.UseInTransitWarehouse = payloads.UseInTransitWarehouse
	Entities.InTransitWarehouseId = payloads.InTransitWarehouseId
	//save data to db
	err = db.Save(&Entities).Error
	if err != nil {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed To Update GoodsReceive",
		}
	}
	return Entities, nil
}

// [uspg_atItemGRPO1_Insert]
// option 0
func (repository *GoodsReceiveRepositoryImpl) InsertGoodsReceiveDetail(db *gorm.DB, payloads transactionsparepartpayloads.GoodsReceiveDetailInsertPayloads) (transactionsparepartentities.GoodsReceiveDetail, *exceptions.BaseErrorResponse) {
	//get the header first
	GoodsReceiveId := payloads.GoodsReceiveSystemNumber
	var GoodsReceiveEntities transactionsparepartentities.GoodsReceive
	var GoodsReceiveDetail transactionsparepartentities.GoodsReceiveDetail
	err := db.Model(&GoodsReceiveEntities).Where(transactionsparepartentities.GoodsReceive{GoodsReceiveSystemNumber: GoodsReceiveId}).
		First(&GoodsReceiveEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("goods receive header is not found"),
				Message:    "GoodsReceive Is Not Found",
			}
		}
		return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed To Update GoodsReceiveDetail",
		}
	}
	//get costing type for calidasion warehouse
	var WarehouseCostingTypeHPPId int
	err = db.Model(&masterwarehouseentities.WarehouseCostingType{}).
		Select("warehouse_costing_type_id").
		Where(masterwarehouseentities.WarehouseCostingType{WarehouseCostingTypeCode: "HPP"}).
		Scan(&WarehouseCostingTypeHPPId).Error
	if err != nil {
		return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed To get costing type HPP",
		}
	}
	//@Hpp_wh_Type = COSTING_TYPE
	//FROM gmLoc1  WHERE COMPANY_CODE  = @Company_Code and WAREHOUSE_CODE=@Whs_Code2
	//var entities masteritementities.ItemLocation
	var locationStockEntities masteritementities.ItemLocation
	//get location stock for validation
	err = db.Model(&locationStockEntities).Where(masteritementities.ItemLocation{

		WarehouseLocationId: payloads.WarehouseLocationId,
		WarehouseId:         GoodsReceiveEntities.WarehouseId,
		ItemId:              payloads.ItemId}).
		First(&locationStockEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("invalid Item location"),
				Message:    "Invalid Item Location..!!",
			}
		}
		return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed To Get Location Warehouse",
		}
	}

	//			IF (@Hpp_wh_Type =(@Hpp_wh_Type_Cmp) OR @Hpp_wh_Type = @Hpp_wh_Type_Normal)
	WarehouseCostingType := 0
	err = db.Model(&masterwarehouseentities.WarehouseMaster{}).Where("warehouse_id = ?", GoodsReceiveEntities.WarehouseId).
		Select("warehouse_costing_type_id").Scan(&WarehouseCostingType).
		Error
	if err != nil || WarehouseCostingType == 0 {
		return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed To warehouse master with id = " + strconv.Itoa(GoodsReceiveEntities.WarehouseId),
		}
	}
	//IF (@Hpp_wh_Type =(@Hpp_wh_Type_Cmp) OR @Hpp_wh_Type = @Hpp_wh_Type_Normal)
	//BEGIN
	//IF (ISNULL(@Item_Price,0) <=0)
	//BEGIN
	//RAISERROR(' Price must be greater than 0',16,1)
	//RETURN 0
	//END
	//END
	if WarehouseCostingType == WarehouseCostingTypeHPPId && payloads.ItemPrice <= 0 {
		return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("price must be greater then 0"),
		}
	}
	//IF EXISTS(SELECT LOC_CODE FROM amLocationItem WHERE COMPANY_CODE = @Company_Code AND ITEM_CODE = @Item_Code AND LOC_CODE = @Loc_Code AND STOCK_OPNAME = 1 AND WHS_GROUP = @Whs_Group)

	if locationStockEntities.StockOpname && locationStockEntities.WarehouseGroupId == GoodsReceiveEntities.WarehouseGroupId {
		return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("location is under stock opname"),
		}
	}
	//var locationStockEntities masteritementities.ItemLocation
	isStockOpaname := 0
	err = db.Model(&locationStockEntities).Where(masteritementities.ItemLocation{

		WarehouseLocationId: payloads.WarehouseLocationClaimId,
		WarehouseId:         GoodsReceiveEntities.WarehouseId,
		ItemId:              payloads.ItemId,
		StockOpname:         true}).
		Select("1").
		Scan(&isStockOpaname).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Invalid Item Location",
			}
		}
		return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed To Get Claim Location Warehouse",
		}
	}
	if isStockOpaname == 1 {
		return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("location is under stock opname"),
		}
	}
	//get type item group ID
	ItemGroupTypeInventoryId := 0
	err = db.Model(&masteritementities.ItemGroup{}).Select("item_group_id").Scan(&ItemGroupTypeInventoryId).
		Error
	if err != nil || ItemGroupTypeInventoryId == 0 {
		return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed To Get Item Group Type Inventory",
		}
	}
	if payloads.WarehouseLocationClaimId != 0 {
		var CheckDuplicateItemClaim = 0
		err = db.Table("trx_goods_receive_detail GR1").
			Joins("LEFT JOIn trx_goods_receive GR0 ON GR0.goods_receive_system_number = gr1.binning_system_number").
			Where(`
						WHERE GR1.warehouse_location_id = ? AND gr1.item_id = ? and gr0.item_group_id = ?
							`, payloads.WarehouseLocationId, payloads.ItemId, ItemGroupTypeInventoryId).
			Select("1").Scan(&CheckDuplicateItemClaim).Error
		if err != nil && CheckDuplicateItemClaim == 1 {
			return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("selected location claim is already exists in Goods Receipt detail"),
			}
		}
	}
	var GoodsReceiveReferenceEntities masterentities.GoodsReceiveReferenceType
	err = db.Model(&GoodsReceiveReferenceEntities).First(&GoodsReceiveReferenceEntities).
		Where("reference_type_good_receive_id = ?", GoodsReceiveEntities.ReferenceTypeGoodReceiveId).
		Error
	if err != nil {
		return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	//for ItemPrice Update

	var ItemGoodsReceive transactionsparepartpayloads.ItemGoodsReceiveTemp
	if GoodsReceiveReferenceEntities.ReferenceTypeGoodReceiveCode == "PO" {
		//SELECT
		//@Item_Price = ISNULL(ITEM_PRICE,0) - ISNULL(ITEM_DISC_AMOUNT,0),
		//@Item_Disc_Percent = ISNULL(ITEM_DISC_PERCENT,0) ,
		//@Item_Disc_Amount = ISNULL(ITEM_DISC_AMOUNT,0)
		//FROM atItemPO1 WHERE PO_SYS_NO = @Ref_Sys_No AND PO_LINE = @Ref_Line_No
		err = db.Model(&transactionsparepartentities.PurchaseOrderDetailEntities{}).
			Select(`
						ISNULL(item_price,0) - ISNULL(item_discount_amount,0) as item_price,
						ISNULL(item_discount_percentage,0) as item_discount_percentage,
						ISNULL(item_discount_amount,0) as item_discount_amount
				`).
			Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderDetailSystemNumber: payloads.ReferenceSystemNumber}).
			Scan(&ItemGoodsReceive).
			Error
		if err != nil {
			return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed To Get Item Goods Receive from purchase order",
			}
		}
	}
	if GoodsReceiveReferenceEntities.ReferenceTypeGoodReceiveCode == "CL" {
		err = db.Model(&transactionsparepartentities.ItemClaimDetail{}).
			Select(`
						ISNULL(item_price,0) as item_price,
						ISNULL(item_discount_percentage,0) as item_discount_percentage,
						ISNULL(item_discount_amount,0) as item_discount_amount
				`).
			Where(transactionsparepartentities.ItemClaimDetail{ItemClaimDetailId: payloads.ReferenceSystemNumber}).
			Scan(&ItemGoodsReceive).Error
		if err != nil {
			return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed To Get Item Goods Receive from Item Claim",
			}
		}
	}
	if GoodsReceiveReferenceEntities.ReferenceTypeGoodReceiveCode == "WC" {
		ItemGoodsReceive.ItemDiscPercent = 0
		ItemGoodsReceive.ItemDiscAmount = 0
		ItemGoodsReceive.ItemPrice = 0
	}
	payloads.ItemPrice = ItemGoodsReceive.ItemPrice
	payloads.ItemDiscountAmount = ItemGoodsReceive.ItemDiscAmount
	payloads.ItemDiscountPercent = ItemGoodsReceive.ItemDiscPercent

	//search UomRate
	GoodsReceiveDetail = transactionsparepartentities.GoodsReceiveDetail{
		//BinningId
		GoodsReceiveSystemNumber: payloads.GoodsReceiveSystemNumber,
		GoodsReceiveLineNumber:   payloads.GoodsReceiveLineNumber,
		ItemId:                   payloads.ItemId,
		ItemUnitOfMeasurement:    payloads.ItemUnitOfMeasurement,
		UnitOfMeasurementRate:    payloads.UnitOfMeasurementRate,
		ItemPrice:                payloads.ItemPrice,
		ItemDiscountPercent:      payloads.ItemDiscountPercent,
		ItemDiscountAmount:       payloads.ItemDiscountAmount,
		QuantityReference:        payloads.QuantityReference,
		QuantityDeliveryOrder:    payloads.QuantityDeliveryOrder,
		QuantityShort:            payloads.QuantityShort,
		QuantityDamage:           payloads.QuantityDamage,
		QuantityOver:             payloads.QuantityOver,
		QuantityWrong:            payloads.QuantityWrong,
		QuantityGoodsReceive:     payloads.QuantityGoodsReceive,
		WarehouseLocationId:      payloads.WarehouseLocationId,
		WarehouseLocationClaimId: payloads.WarehouseLocationClaimId,
		CaseNumber:               payloads.CaseNumber,
		BinningId:                payloads.BinningId,
		BinningDocumentNumber:    payloads.BinningDocumentNumber,
		BinningLineNumber:        payloads.BinningLineNumber,
		ReferenceSystemNumber:    payloads.ReferenceSystemNumber,
		ReferenceLineNumber:      payloads.ReferenceLineNumber,
		ClaimSystemNumber:        payloads.ClaimSystemNumber,
		ClaimLineNumber:          payloads.ClaimLineNumber,
		ItemTotal:                (payloads.QuantityGoodsReceive + payloads.QuantityShort) * payloads.ItemPrice,
		ItemTotalBaseAmount:      (payloads.QuantityGoodsReceive + payloads.QuantityShort) * payloads.ItemPrice * GoodsReceiveEntities.CurrencyExchangeRate,
		CreatedByUserId:          payloads.CreatedByUserId,
		CreatedDate:              payloads.CreatedDate,
		UpdatedByUserId:          payloads.UpdatedByUserId,
		UpdatedDate:              payloads.UpdatedDate,
		ChangeNo:                 1,
	}
	err = db.Create(&GoodsReceiveDetail).Error
	if err != nil {
		return GoodsReceiveDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to create goods receive detail entities"),
		}
	}
	return GoodsReceiveDetail, nil
}

// [uspg_atItemGRPO1_Update] option 2
func (repository *GoodsReceiveRepositoryImpl) UpdateGoodsReceiveDetail(db *gorm.DB, payloads transactionsparepartpayloads.GoodsReceiveDetailUpdatePayloads, DetailId int) (bool, *exceptions.BaseErrorResponse) {
	var goodsReceiveHeader transactionsparepartentities.GoodsReceive
	var goodsReceiveDetail transactionsparepartentities.GoodsReceiveDetail
	var GoodsReceiveReferenceEntities masterentities.GoodsReceiveReferenceType
	//get reference typr
	//get detail to update
	err := db.Model(&goodsReceiveDetail).Where("goods_receive_detail_system_number = ?", DetailId).
		First(&goodsReceiveDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("detail data to update is not found"),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to update goods receive detail, failed to retrieves data",
		}
	}
	//get header for validation

	err = db.Model(&goodsReceiveHeader).Where("goods_receive_system_number = ?", goodsReceiveDetail.GoodsReceiveSystemNumber).
		First(&goodsReceiveHeader).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("header data is not found"),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get header data = " + err.Error(),
		}
	}

	err = db.Model(&GoodsReceiveReferenceEntities).First(&GoodsReceiveReferenceEntities).
		Where("reference_type_good_receive_id = ?", goodsReceiveHeader.ReferenceTypeGoodReceiveId).
		Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	var actualWhsId int
	if goodsReceiveHeader.UseInTransitWarehouse {
		actualWhsId = goodsReceiveHeader.InTransitWarehouseId
	} else {
		actualWhsId = goodsReceiveHeader.WarehouseId
	}
	var locationStockEntities masteritementities.ItemLocation
	err = db.Model(&locationStockEntities).Where(masteritementities.ItemLocation{

		WarehouseLocationId: payloads.WarehouseLocationId,
		WarehouseId:         actualWhsId,
		ItemId:              payloads.ItemId}).
		First(&locationStockEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("invalid Item location"),
				Message:    "Invalid Item Location..!!",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed To Get Location Warehouse",
		}
	}
	if locationStockEntities.StockOpname && locationStockEntities.WarehouseGroupId == goodsReceiveHeader.WarehouseGroupId {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("location is under stock opname"),
		}
	}

	if payloads.WarehouseLocationClaimId != 0 {
		IsOpname := 0
		err = db.Model(&locationStockEntities).
			Where(masteritementities.ItemLocation{WarehouseId: actualWhsId, ItemId: payloads.ItemId, WarehouseLocationId: payloads.WarehouseLocationClaimId, StockOpname: true}).
			Select("1").Scan(&IsOpname).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
			}
		}
		if IsOpname == 1 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("location Claim is under Stock Opname"),
			}
		}
	}
	ItemGroupTypeInventoryId := 0
	err = db.Model(&masteritementities.ItemGroup{}).Select("item_group_id").Scan(&ItemGroupTypeInventoryId).
		Error
	if err != nil || ItemGroupTypeInventoryId == 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed To Get Item Group Type Inventory",
		}
	}
	if payloads.WarehouseLocationClaimId != goodsReceiveDetail.WarehouseLocationClaimId {
		var CheckDuplicateItemClaim int = 0
		err = db.Table("trx_goods_receive_detail GR1").
			Joins("LEFT JOIN trx_goods_receive GR0 ON GR0.goods_receive_system_number = gr1.binning_system_number").
			Where(`
						WHERE GR1.warehouse_location_id = ? 
						AND GR1.item_id = ? 
						AND GR0.item_group_id = ?
						AND GR1.goods_receive_system_number = ?
							`,
				payloads.WarehouseLocationClaimId,
				payloads.ItemId,
				ItemGroupTypeInventoryId,
				goodsReceiveDetail.GoodsReceiveSystemNumber).
			Select("1").Scan(&CheckDuplicateItemClaim).Error
		if err != nil && CheckDuplicateItemClaim == 1 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("selected location claim is already exists in Goods Receipt detail"),
			}
		}
	}
	var ItemGoodsReceive transactionsparepartpayloads.ItemGoodsReceiveTemp
	if GoodsReceiveReferenceEntities.ReferenceTypeGoodReceiveCode == "PO" {
		//SELECT
		//@Item_Price = ISNULL(ITEM_PRICE,0) - ISNULL(ITEM_DISC_AMOUNT,0),
		//@Item_Disc_Percent = ISNULL(ITEM_DISC_PERCENT,0) ,
		//@Item_Disc_Amount = ISNULL(ITEM_DISC_AMOUNT,0)
		//FROM atItemPO1 WHERE PO_SYS_NO = @Ref_Sys_No AND PO_LINE = @Ref_Line_No
		err = db.Model(&transactionsparepartentities.PurchaseOrderDetailEntities{}).
			Select(`
						ISNULL(item_price,0) - ISNULL(item_discount_amount,0) as item_price,
						ISNULL(item_discount_percentage,0) as item_discount_percentage,
						ISNULL(item_discount_amount,0) as item_discount_amount
				`).
			Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderDetailSystemNumber: payloads.ReferenceSystemNumber}).
			Scan(&ItemGoodsReceive).
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed To Get Item Goods Receive from purchase order",
			}
		}
	}
	if GoodsReceiveReferenceEntities.ReferenceTypeGoodReceiveCode == "CL" {
		err = db.Model(&transactionsparepartentities.ItemClaimDetail{}).
			Select(`
						ISNULL(item_price,0) as item_price,
						ISNULL(item_discount_percentage,0) as item_discount_percentage,
						ISNULL(item_discount_amount,0) as item_discount_amount
				`).
			Where(transactionsparepartentities.ItemClaimDetail{ItemClaimDetailId: payloads.ReferenceSystemNumber}).
			Scan(&ItemGoodsReceive).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed To Get Item Goods Receive from Item Claim",
			}
		}
	}
	if GoodsReceiveReferenceEntities.ReferenceTypeGoodReceiveCode == "WC" {
		ItemGoodsReceive.ItemDiscPercent = 0
		ItemGoodsReceive.ItemDiscAmount = 0
		ItemGoodsReceive.ItemPrice = 0
	}
	payloads.ItemPrice = ItemGoodsReceive.ItemPrice
	payloads.ItemDiscountAmount = ItemGoodsReceive.ItemDiscAmount
	payloads.ItemDiscountPercent = ItemGoodsReceive.ItemDiscPercent

	goodsReceiveDetail.QuantityReference = payloads.QuantityReference
	goodsReceiveDetail.QuantityDeliveryOrder = payloads.QuantityDeliveryOrder
	goodsReceiveDetail.QuantityShort = payloads.QuantityShort
	goodsReceiveDetail.QuantityDamage = payloads.QuantityDamage
	goodsReceiveDetail.QuantityOver = payloads.QuantityOver
	goodsReceiveDetail.QuantityWrong = payloads.QuantityWrong
	goodsReceiveDetail.QuantityGoodsReceive = payloads.QuantityGoodsReceive
	goodsReceiveDetail.WarehouseLocationId = payloads.WarehouseLocationId
	goodsReceiveDetail.WarehouseLocationClaimId = payloads.WarehouseLocationClaimId
	goodsReceiveDetail.ItemPrice = payloads.ItemPrice
	goodsReceiveDetail.ItemDiscountPercent = payloads.ItemDiscountPercent
	goodsReceiveDetail.ItemDiscountAmount = payloads.ItemDiscountAmount
	goodsReceiveDetail.ItemTotal = (payloads.QuantityGoodsReceive + payloads.QuantityShort) * payloads.ItemPrice
	goodsReceiveDetail.ItemTotalBaseAmount = goodsReceiveDetail.ItemTotal * goodsReceiveHeader.CurrencyExchangeRate
	goodsReceiveDetail.CaseNumber = payloads.CaseNumber
	goodsReceiveDetail.ChangeNo += 1
	goodsReceiveDetail.UpdatedByUserId = payloads.UpdatedByUserId
	goodsReceiveDetail.UpdatedDate = payloads.UpdatedDate

	err = db.Save(&goodsReceiveDetail).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed To Update Goods Receive Error On : " + err.Error(),
		}
	}
	return true, nil
}
func (repository *GoodsReceiveRepositoryImpl) LocationItemGoodsReceive(db *gorm.DB, filter []utils.FilterCondition, PaginationParams pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	ItemLocationEntities := masteritementities.ItemLocation{}
	var Responses []transactionsparepartpayloads.GetAllLocationGRPOResponse
	joinTable := db.Table("mtr_location_item A").
		Joins(`LEFT JOIN mtr_warehouse_location B ON A.warehouse_id = B.warehouse_id AND B.warehouse_location_id = A.warehouse_location_id`).
		Joins(`LEFT JOIN mtr_warehouse_master whs on A.warehouse_id = whs.warehouse_id`).
		Joins(`INNER JOIN mtr_item item on A.item_id = item.item_id`).
		Select(`
				A.warehouse_id,
				A.item_id, 
				A.item_location_id,
				B.warehouse_location_name,
				item.item_code,
				whs.company_id,
				whs.warehouse_code
				`)
	whereQuery := utils.ApplyFilter(joinTable, filter)
	err := whereQuery.Scopes(pagination.Paginate(&ItemLocationEntities, &PaginationParams, whereQuery)).Order("warehouse_code").Scan(&Responses).Error
	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to get all location item"),
		}
	}
	if len(Responses) == 0 {
		PaginationParams.Rows = []string{}
		return PaginationParams, nil
	}
	PaginationParams.Rows = Responses
	return PaginationParams, nil
}
func (repository *GoodsReceiveRepositoryImpl) SubmitGoodsReceive(db *gorm.DB, GoodsReceiveId int) (bool, *exceptions.BaseErrorResponse) {
	//get entities first
	GoodsReceiveEntities := transactionsparepartentities.GoodsReceive{}
	err := db.Model(&GoodsReceiveEntities).Where(transactionsparepartentities.GoodsReceive{GoodsReceiveSystemNumber: GoodsReceiveId}).
		First(&GoodsReceiveEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("goods receive with id : %d is not found", GoodsReceiveId),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to retrieve goods receive",
		}
	}

	var PeriodResponse financeservice.OpenPeriodPayloadResponse
	PeriodUrl := config.EnvConfigs.FinanceServiceUrl + "closing-period-company/current-period?company_id=" + strconv.Itoa(GoodsReceiveEntities.CompanyId) + "&closing_module_detail_code=SP" //strconv.Itoa(response.ItemCode)

	if err := utils.Get(PeriodUrl, &PeriodResponse, nil); err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to Period Response data from external service",
			Err:        err,
		}
	}
	//IF ((SELECT COUNT(GRPO_LINE_NO) FROM atItemGRPO1 WITH(NOLOCK) WHERE GRPO_SYS_NO = ISNULL(@Grpo_Sys_No,0) AND ISNULL(LOC_CODE,'') = '') > 0)
	//BEGIN
	//RAISERROR('Location Code must be filled',16,1)
	//RETURN 0
	//END
	var isExist = 0
	var goodsReceiveDetailEntities transactionsparepartentities.GoodsReceiveDetail
	err = db.Model(&goodsReceiveDetailEntities).
		Select("count(goods_receive_detail_system_number)").
		Where("warehouse_location_id = 0 AND goods_receive_system_number = ?", GoodsReceiveId).
		Scan(&isExist).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get goods receive detail",
		}
	}
	if isExist > 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("location Code must be filled"),
		}
	}
	//IF ((SELECT COUNT(GRPO_LINE_NO) FROM atItemGRPO1 WITH(NOLOCK) WHERE GRPO_SYS_NO = ISNULL(@Grpo_Sys_No,0) AND ISNULL(LOC_CLAIM_CODE,'') = '' AND QTY_DAMAGE + QTY_OVER + QTY_SHORT + QTY_WRONG > 0) > 0)
	//BEGIN
	//RAISERROR('Location Claim Code must be filled for Item that has Claim',16,1)
	//RETURN 0
	//END
	isExist = 0
	err = db.Model(&goodsReceiveDetailEntities).
		Select("count(goods_receive_detail_system_number)").
		Where(`goods_receive_system_number = ? 
					AND warehouse_location_id = 0
					AND quantity_short+quantity_damage +quantity_over+ quantity_wrong > 0`, GoodsReceiveId).
		Scan(&isExist).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "cannot cek location claim code",
		}
	}
	if isExist > 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("location Claim Code must be filled for Item that has Claim"),
		}
	}
	isExist = 0
	err = db.Model(&transactionsparepartentities.GoodsReceiveDetail{}).
		Select("count(goods_receive_detail_system_number)").Where(`
		goods_receive_system_number = ? 
		AND ISNULL(quantity_delivery_order,0) <> (ISNULL(quantity_goods_receive,0) + ISNULL(quantity_short,0) + ISNULL(quantity_damage,0))
		`, GoodsReceiveId).Scan(&isExist).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	if isExist > 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("receive Qty is not valid")}
	}
	//get item group
	var GoodsReceivesItemGroupEntities masteritementities.ItemGroup
	err = db.Model(&GoodsReceivesItemGroupEntities).
		Where(masteritementities.ItemGroup{ItemGroupId: GoodsReceiveEntities.ItemGroupId}).
		First(&GoodsReceivesItemGroupEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("item group with id %d is not found", GoodsReceiveEntities.ItemGroupId),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error on fetching item group from database"),
		}
	}
	//IF @Item_Group = @ItemGroupOutsideJob
	//BEGIN
	//--IF EXISTS(SELECT RECORD_STATUS FROM atItemGRPO1 WITH (ROWLOCK) WHERE GRPO_SYS_NO = @Grpo_Sys_No AND ISNULL(QTY_DO,0) <> ISNULL(QTY_GRPO,0))
	//IF ((SELECT COUNT(RECORD_STATUS) FROM atItemGRPO1 WITH(NOLOCK) WHERE GRPO_SYS_NO = @Grpo_Sys_No AND ISNULL(QTY_DO,0) <> ISNULL(QTY_GRPO,0)) > 0)
	//BEGIN
	//RAISERROR('Qty GRPO Cannot Less Than Qty DO For Item Group OJ',16,1)
	//RETURN 0
	//END
	//END
	if GoodsReceivesItemGroupEntities.ItemGroupCode == "OJ" {
		//validate if grpo quantity les than do for item group oj
		isExist = 0
		err = db.Model(&transactionsparepartentities.GoodsReceiveDetail{}).
			Select("count(goods_receive_detail_system_number)").Scan(&isExist).
			Where("goods_receive_system_number = ? AND ISNULL(quantity_delivery_order,0) <> ISNULL(quantity_goods_receive,0)").
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to check grpo quantity with item group oj on error : " + err.Error(),
			}
		}
		if isExist > 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("quantity GRPO cannot lest than quantity DO for item group oj"),
			}
		}
	}
	//get is use dms for gm ref checking
	//hit general service
	CompanyReferenceBetByIdResponse := generalservicepayloads.CompanyReferenceBetByIdResponse{}
	CompanyReferenceUrl := fmt.Sprintf("%scompany-reference/%s", config.EnvConfigs.GeneralServiceUrl, strconv.Itoa(GoodsReceiveEntities.SupplierId))
	errFetchCompany := utils.Get(CompanyReferenceUrl, &CompanyReferenceBetByIdResponse, nil)
	if errFetchCompany != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errFetchCompany,
			Message:    errFetchCompany.Error(),
		}
	}
	if strings.ToUpper(GoodsReceivesItemGroupEntities.ItemGroupCode) == "IN" && CompanyReferenceBetByIdResponse.UseDms {
		isExist = 0
		err = db.Model(&transactionsparepartentities.GoodsReceiveDetail{}).
			Select("count(goods_receive_detail_system_number)").Scan(&isExist).
			Where("goods_receive_system_number = ? AND ISNULL(quantity_delivery_order,0) <> ISNULL(quantity_goods_receive,0)").
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to check grpo quantity with item group IN on error : " + err.Error(),
			}
		}
		if isExist > 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("there is short Qty in Goods Receive detail"),
			}
		}
		//VALIDATION BELOW IS NOT DEVELOP YET
		//WAITING FOR INVOICE FROM FINANCE READY
		//	--QTY TOTAL GRPO HARUS SAMA DENGAN QTY TOTAL PADA INVOICE SUPPLIER SESAMA PENGGUNA DMS
		//	IF	(SELECT SUM(QTY_GRPO) FROM atItemGRPO1 WHERE GRPO_SYS_NO = @Grpo_Sys_No)
		//	<>
		//	(SELECT SUM(ITEM_QTY) FROM rtInvoice1 WHERE INV_SYS_NO =	(	SELECT INV_SYS_NO
		//	FROM rtInvoice0
		//	WHERE INV_DOC_NO = @Supplier_Inv_No
		//	AND CONVERT(VARCHAR,COMPANY_CODE) = @Supplier_Code
		//	AND BILL_TO_CUST_CODE = CONVERT(VARCHAR,@Company_Code)
		//)
		//)
		//	AND @Company_Code <> '3125098' --Hanya NMDI yang bisa terima parsial atas Supplier sesama pengguna DMS
		//	BEGIN
		//	RAISERROR('Total item qty Goods Receive not match with Supplier Invoice detail!',16,1)
		//	RETURN 0
		//	END
		//	END
	}
	if GoodsReceiveEntities.ReferenceTypeGoodReceiveId == 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("reference type is empty"),
		}
	}
	//get reference type goods receive to validate
	var goodsReceiveReferenceTypeEntities masterentities.GoodsReceiveReferenceType
	err = db.Model(&goodsReceiveReferenceTypeEntities).Where(masterentities.GoodsReceiveReferenceType{ReferenceTypeGoodReceiveId: GoodsReceiveEntities.ReferenceTypeGoodReceiveId}).
		First(&goodsReceiveReferenceTypeEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("goods receive reference type is not found id : " + strconv.Itoa(GoodsReceiveEntities.ReferenceTypeGoodReceiveId)),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error on fetching reference type from database"),
		}
	}
	CostingTypeNon, errs := getCostingTypeByCode(db, "NON")
	if errs != nil {
		return false, errs
	}
	isCostingTypeNon := 0
	//get costing type non first
	err = db.Model(&masterwarehouseentities.WarehouseMaster{}).
		Where(masterwarehouseentities.WarehouseMaster{WarehouseCostingTypeId: CostingTypeNon.WarehouseCostingTypeId,
			WarehouseId: GoodsReceiveEntities.WarehouseId,
			CompanyId:   GoodsReceiveEntities.CompanyId}).
		Select("1").Scan(&isCostingTypeNon).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get warehouse costing type",
		}
	}
	if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "PO" {
		//cek purchase order table
		//select first the purchase order table

		var PurchaseOrderEntities transactionsparepartentities.PurchaseOrderEntities
		err = db.Model(&PurchaseOrderEntities).Where(transactionsparepartentities.PurchaseOrderEntities{
			PurchaseOrderSystemNumber: GoodsReceiveEntities.ReferenceSystemNumber,
		}).First(&PurchaseOrderEntities).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Err:        fmt.Errorf("purchase order with reference number : %s is not found on table purchase order", GoodsReceiveEntities.ReferenceDocumentNumber),
				}
			}
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("failed to fetch from purchase order"),
			}
		}
		if PurchaseOrderEntities.ProfitCenterId != GoodsReceiveEntities.ProfitCenterId {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("profit center with in PO With Id : %d is not  match with gr Profit Center Id %d", PurchaseOrderEntities.ProfitCenterId, GoodsReceiveEntities.ProfitCenterId),
			}
		}

		if isCostingTypeNon > 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("warehouse is not valid. Cannot use warehouse with costing type Non"),
			}
		}
		//this validtion for grpo import last. not yet dev. deve the local first
		//IF ((SELECT COUNT(a.GRPO_SYS_NO) FROM atItemGRPO1 a WITH(NOLOCK)
		//inner join atItemPO0 c WITH(NOLOCK) on c.PO_SYS_NO = a.REF_SYS_NO
		//inner join atItemPO1 d WITH(NOLOCK) on d.PO_SYS_NO = a.REF_SYS_NO and d.PO_LINE = a.REF_LINE_NO
		//inner join atItemPR0 b WITH(NOLOCK) on b.PR_SYS_NO = d.PR_SYS_NO
		//inner join atSalesOrder0 e WITH(NOLOCK) on e.SO_SYS_NO = c.SO_SYS_NO and e.PO_SYS_NO = c.PO_SYS_NO and e.PO_TYPE = @potype_inv
		//inner join atSalesOrder1 f WITH(NOLOCK) on f.SO_SYS_NO = e.SO_SYS_NO and f.PO_SYS_NO = d.PO_SYS_NO and f.PO_LINE_NO = d.PO_LINE
		//left join atPickList1 g WITH(NOLOCK) on g.SO_SYS_NO = f.SO_SYS_NO and g.SO_LINE_NO = f.SO_LINE_NO
		//left join atPackList1 h WITH(NOLOCK) on h.PICK_LIST_SYS_NO = g.PICK_LIST_SYS_NO and h.PICK_LIST_LINE = g.PICK_LIST_LINE
		//left join rtInvoice0 i WITH(NOLOCK) on i.REF2_SYS_NO = h.PACK_LIST_SYS_NO and i.REF_SYS_NO = e.SO_SYS_NO
		//where a.GRPO_SYS_NO = @Grpo_Sys_No and b.REF_DOC_TYPE = @reftype_unit and isnull(i.INV_STATUS,'') <> @approval_approved) > 0)
		//
		//begin
		//raiserror('Goods Receive Accs Inter Dealer cannot be made until Invoice SO from the opposite dealer is made',16,1)
		//RETURN 0
		//end
		if GoodsReceiveEntities.AffiliatedPurchaseOrder &&
			GoodsReceivesItemGroupEntities.ItemGroupCode != "OJ" &&
			GoodsReceivesItemGroupEntities.ItemGroupCode != "OX" &&
			GoodsReceivesItemGroupEntities.ItemGroupCode != "FA" {
			//not yet dev forinter grpo and for rtinvoice is not yet dev in finance
			//IF EXISTS(SELECT TOP 1 1
			//FROM atitemgrpo0 A
			//INNER JOIN gmRef C ON A.SUPPLIER_CODE = CAST(C.COMPANY_CODE AS VARCHAR)
			//WHERE ISNULL(C.USE_DMS,0) = 1
			//AND A.GRPO_SYS_NO = @Grpo_Sys_No)
			//BEGIN
			//IF NOT EXISTS (select TOP 1 1
			//	from atItemGRPO0 A WITH(NOLOCK)
			//	INNER JOIN atItemPO0 B WITH(NOLOCK) ON A.REF_SYS_NO = B.PO_SYS_NO
			//	INNER JOIN gmRef C WITH(NOLOCK) ON B.SUPPLIER_CODE = C.COMPANY_CODE AND ISNULL(C.USE_DMS,0) = 1
			//	INNER JOIN atSalesOrder0 D WITH(NOLOCK) ON D.PO_SYS_NO = B.PO_SYS_NO
			//	INNER JOIN rtinvoice0 E WITH(NOLOCK) ON E.REF_SYS_NO = D.SO_SYS_NO AND E.REF_TYPE = @varARinvRefTypeSO AND E.INV_TYPE = @varARInvTypeSO
			//	WHERE A.GRPO_SYS_NO = @Grpo_Sys_No
			//	AND A.REF_TYPE = @varItemGRPORefTypePO
			//	AND ISNULL(E.INV_STATUS,'') = @approval_approved)
			//	BEGIN
			//	raiserror('Goods Receive Sparepart Inter Dealer cannot be made until Invoice SO from the opposite dealer is made',16,1)
			//	RETURN 0
			//	END
			//	END
		}
	} else if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "WC" {
		//cek if profit center is match
		isExist = 0
		err = db.Model(&transactionworkshopentities.WorkOrder{}).
			Where(transactionworkshopentities.WorkOrder{WorkOrderSystemNumber: GoodsReceiveEntities.ReferenceSystemNumber, ProfitCenterId: GoodsReceiveEntities.ProfitCenterId}).
			Select("1").Scan(&isExist).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("system number %d is not found in reference work order", GoodsReceiveEntities.ReferenceSystemNumber),
			}
		}
		if isExist > 1 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("profit center not match"),
			}
		}
		if isCostingTypeNon > 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("warehouse is not valid. Cannot use warehouse with costing type Non"),
			}
		}
	} else if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "CL" {
		isExist = 0
		err = db.Model(&transactionsparepartentities.ItemClaim{}).
			Where(transactionsparepartentities.ItemClaim{ClaimSystemNumber: GoodsReceiveEntities.ReferenceSystemNumber, ProfitCenterId: GoodsReceiveEntities.ProfitCenterId}).
			Select("1").Scan(&isExist).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("system number %s is not found in reference item claim", strconv.Itoa(GoodsReceiveEntities.ReferenceSystemNumber)),
			}
		}
		if isExist > 1 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("profit center not match"),
			}
		}
	}
	var PeriodResponseSp financeservice.OpenPeriodPayloadResponse

	if GoodsReceivesItemGroupEntities.IsItemSparepart {
		PeriodUrl := config.EnvConfigs.FinanceServiceUrl + "closing-period-company/current-period?company_id=" + strconv.Itoa(GoodsReceiveEntities.CompanyId) + "&closing_module_detail_code=SP" //strconv.Itoa(response.ItemCode)
		if err := utils.Get(PeriodUrl, &PeriodResponseSp, nil); err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to Period Response data from external service",
				Err:        err,
			}
		}
	} else {
		PeriodUrl := config.EnvConfigs.FinanceServiceUrl + "closing-period-company/current-period?company_id=" + strconv.Itoa(GoodsReceiveEntities.CompanyId) + "&closing_module_detail_code=AP" //strconv.Itoa(response.ItemCode)
		if err := utils.Get(PeriodUrl, &PeriodResponseSp, nil); err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to Period Response data from external service",
				Err:        err,
			}
		}
	}
	//validate if period is open
	if strconv.Itoa(GoodsReceiveEntities.GoodsReceiveDocumentDate.Year()) != PeriodResponseSp.PeriodYear &&
		strconv.Itoa(int(GoodsReceiveEntities.GoodsReceiveDocumentDate.Month())) != PeriodResponseSp.PeriodMonth {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("period status %s for this goods receive is already closed ", GoodsReceiveEntities.GoodsReceiveDocumentDate),
		}
	}
	//generate dummy document number
	//EXEC uspg_gmSrcDoc1_Update
	//@Option = 0 ,
	//@COMPANY_CODE = @Company_Code ,
	//@SOURCE_CODE = @Src_Code ,
	//@VEHICLE_BRAND = @VEHICLE_BRAND ,
	//@PROFIT_CENTER_CODE = @PROFIT_CENTER ,
	//@TRANSACTION_CODE = '' ,
	//@BANK_ACC_CODE = '' ,
	//@TRANSACTION_DATE = @Grpo_Doc_Date ,
	//@Last_Doc_No = @Grpo_Doc_No OUTPUT
	docno, errorGenerate := GenerateDocumentNumber(db, GoodsReceiveId)
	if errorGenerate != nil {
		return false, errorGenerate
	}
	//get document status ready for grpo
	StatusIdReady := 0
	err = db.Model(masterentities.GoodsReceiveDocumentStatus{}).Where(masterentities.GoodsReceiveDocumentStatus{ItemGoodsReceiveStatusCode: "99"}).
		Select("item_goods_receive_status_id").Scan(&StatusIdReady).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error on getting status document complete",
		}
	}
	GoodsReceiveEntities.GoodsReceiveDocumentNumber = docno
	GoodsReceiveEntities.GoodsReceiveStatusId = StatusIdReady
	GoodsReceiveEntities.ChangeNo += 1
	GoodsReceiveEntities.UpdatedDate = time.Now()
	err = db.Save(&GoodsReceiveEntities).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to save goods receive header on error : " + err.Error()}
	}

	//SET @Qty_Sum_Grpo = (SELECT SUM(ISNULL(GRPO1.QTY_GRPO,0)) FROM atItemGRPO1 GRPO1 with (nolock)
	//INNER JOIN atItemGRPO0 GRPO with (nolock) ON GRPO1.GRPO_SYS_NO = GRPO.GRPO_SYS_NO
	//WHERE GRPO.GRPO_STATUS = @Grpo_Stat_Complete
	//--AND GRPO.REF_TYPE IN (@RefTypeCl,@RefTypePO)
	//AND GRPO.REF_TYPE = @Ref_Type
	//AND GRPO1.REF_SYS_NO = @Ref_Sys_No
	//AND GRPO1.REF_LINE_NO = @CSR1_Ref_Line_No)
	var quantitySumGoodsReceive float64
	err = db.Table("trx_goods_receive_detail A").Joins("INNER JOIN trx_goods_receive B ON A.goods_receive_system_number = B.goods_receive_system_number").
		Where(`
					B.goods_receive_status_id = ?
					AND B.reference_type_good_receive_id = ?
					AND A.reference_system_number = ?
			`, StatusIdReady, GoodsReceiveEntities.ReferenceTypeGoodReceiveId, GoodsReceiveEntities.ReferenceSystemNumber).
		Select("COALESCE(SUM(quantity_goods_receive), 0)").Scan(&quantitySumGoodsReceive).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error on getting quantity sum goods receive",
		}
	}

	//start looping all detailavailable
	var goodsReceivesDetailResponse []transactionsparepartpayloads.GoodsReceiveSubmitResponse
	err = db.Table("trx_goods_receive_detail A").Joins("INNER JOIN mtr_item item on A.item_id = item.item_id").
		Select(`
				A.goods_receive_detail_system_number,
				A.binning_system_number,
				A.binning_detail_id,
				A.reference_line_number,
				A.warehouse_location_id,
				A.warehouse_location_claim_id,
				A.item_id,
				a.item_unit_of_measurement,
				A.item_price,
				a.quantity_reference,
				a.quantity_goods_receive,
				A.quantity_short,
				A.quantity_damage + a.quantity_over+a.quantity_wrong as quantity_variance,
				A.item_discount_percent,
				A.quantity_delivery_order,
				item.stock_keeping,
				item.unit_of_measurement_stock_id,
				A.reference_system_number
			`).Where("goods_receive_system_number = ?", GoodsReceiveId).
		Scan(&goodsReceivesDetailResponse).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get goods receive detail",
		}
	}
	for _, resdb := range goodsReceivesDetailResponse {
		//get onhand and intransit first
		var UomRate float64
		err = db.Table("mtr_location_stock A").
			Joins("INNER JOIN mtr_warehouse_master B ON A.company_id = B.company_id AND A.warehouse_id = b.warehouse_id").
			Where(`
					A.period_year = ?
					AND period_month = ?
					AND A.company_id = ?
					AND A.item_id = ?
					AND A.warehouse_group_id  = ?
					AND B.warehouse_costing_type_id = ?
					`, PeriodResponseSp.PeriodYear,
				PeriodResponseSp.PeriodMonth,
				GoodsReceiveEntities.CompanyId,
				resdb.ItemId,
				GoodsReceiveEntities.WarehouseGroupId,
				CostingTypeNon.WarehouseCostingTypeId,
			).
			Select(`
						A.quantity_in_transit,
						ISNULL(A.quantity_claim_in, 0) + ISNULL(A.quantity_robbing_out, 0) + ISNULL(A.quantity_assembly_out, 0)
					`).Row().Scan(&resdb.QuantityInTransit, &resdb.QuantityOnHand)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				resdb.QuantityInTransit = 0
				resdb.QuantityOnHand = 0
			} else {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "error occured when fetch quantity instransit and quantity on hand",
				}
			}
		}
		if resdb.QuantityInTransit != 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("there is quantity in transit : %f. Please finish the transfer process before Goods Receive", resdb.QuantityInTransit),
			}
		}
		//IF ((SELECT COUNT(*) FROM amLocationItem  WHERE COMPANY_CODE  = @Company_Code AND LOC_CODE = @CSR1_Loc_Code AND ITEM_CODE = @CSR1_Item_Code) = 0)
		//BEGIN
		//SET @Err_Msg = 'Location Code '+@CSR1_Loc_Code+ ' must be set for this item ' + @CSR1_Item_Code
		//RAISERROR( @Err_Msg ,16,1)
		//RETURN 0
		//END
		isExist = 0
		//masteritementities.ItemLocation{}
		err = db.Model(&masteritementities.ItemLocation{}).
			Where(masteritementities.ItemLocation{ItemId: resdb.ItemId,
				WarehouseLocationId: resdb.WarehouseLocationId}).Select("1").
			Scan(&isExist).
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Sprintf("failed to cek warehouse location on error %s", err.Error()),
			}
		}
		if isExist == 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("location id %d must be set for this item %d", resdb.WarehouseLocationId, resdb.ItemId),
			}
		}
		//IF @CSR1_Qty_Grpo = 0 AND @CSR1_Qty_Variance = 0 AND @CSR1_Qty_Short = 0
		//BEGIN
		//RAISERROR('Qty is not valid',16,1)
		//RETURN 0
		//END
		if resdb.QuantityGoodsReceive == 0 && resdb.QuantityVariance == 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("quantity is not valid"),
			}
		}
		//IF @CSR1_Qty_Grpo = 0 AND @CSR1_Qty_Variance = 0 AND @CSR1_Qty_Short > 0 AND @Item_Group <> @ItemGroupInventory
		//BEGIN
		//RAISERROR('Qty is not valid',16,1)
		//RETURN 0
		//END

		var quantitySumReference = 0.0
		//IF @Ref_Type = @RefTypePO
		//BEGIN
		//SET @Qty_Sum_Ref = (SELECT SUM(ISNULL(ITEM_QTY,0)) FROM atItemPO1 WITH(NOLOCK)
		//WHERE PO_SYS_NO = @Ref_Sys_No
		//AND PO_LINE = @CSR1_Ref_Line_No)
		//END
		if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "PO" {
			err = db.Model(&transactionsparepartentities.PurchaseOrderDetailEntities{}).
				Select("SUM(ISNULL(item_quantity,0))").
				Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderSystemNumber: GoodsReceiveEntities.ReferenceSystemNumber}).
				Scan(&quantitySumReference).
				Error
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Error on get total reference quantity from purchase order",
					}
				}
			}
		}
		if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "CL" {
			err = db.Model(&transactionsparepartentities.ItemClaimDetail{}).
				Where(transactionsparepartentities.ItemClaimDetail{ClaimSystemNumber: GoodsReceiveEntities.ReferenceSystemNumber}).
				Select("COALESCE(SUM(quantity_claimed), 0)").
				Scan(&quantitySumReference).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error On Get total quantity reference item claim",
				}
			}
		}
		//IF @Ref_Type = @RefTypeCl
		//BEGIN
		//SET @Qty_Sum_Ref = (SELECT SUM(ISNULL(QTY_CLAIM,0)) FROM atItemClaim1 WITH(NOLOCK)
		//WHERE CLAIM_SYS_NO = @Ref_Sys_No
		//AND CLAIM_LINE_NO = @CSR1_Ref_Line_No)
		//END
		if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "WC" {
			quantitySumReference = resdb.QuantityReference
		}
		if quantitySumReference < quantitySumGoodsReceive {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("total quantity received : %f cannot exceed quantity reference : %f", quantitySumGoodsReceive, quantitySumReference),
			}
		}
		if GoodsReceiveEntities.ViaBinning {
			if (resdb.QuantityDeliveryOrder - resdb.QuantityReference) > 0 {
				if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "PO" {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Err:        fmt.Errorf("delivery order quantity : %f is bigger dan PO quantity :%f", resdb.QuantityDeliveryOrder, resdb.QuantityReference),
					}
				}
				if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "CL" {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Err:        fmt.Errorf("delivery order quantity : %f is bigger dan Claim quantity :%f", resdb.QuantityDeliveryOrder, resdb.QuantityReference),
					}
				}
				if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "WC" {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Err:        fmt.Errorf("delivery order quantity : %f is bigger dan Claim Warranty quantity :%f", resdb.QuantityDeliveryOrder, resdb.QuantityReference),
					}
				}
			}
		}
		if resdb.BinningId != 0 {
			//update variance
			if resdb.QuantityVariance > 0 && resdb.StockKeeping {

				//SELECT @QTY_ONHAND = ISNULL(SUM(ISNULL(VS.QTY_ON_HAND,0)),0),
				//@QTY_INTRANSIT = ISNULL(SUM(ISNULL(VS.QTY_INTRANSIT,0)),0)
				//FROM dbo.viewLocationStock VS
				//LEFT JOIN dbo.gmLoc1 L ON L.WAREHOUSE_CODE = VS.WHS_CODE AND L.COMPANY_CODE = VS.COMPANY_CODE
				//WHERE  PERIOD_YEAR = @Period_Year  AND PERIOD_MONTH = @Period_Month
				//AND VS.COMPANY_CODE = @Company_Code AND ITEM_CODE = @CSR1_Item_Code
				//AND	WHS_GROUP = @Whs_Group AND L.COSTING_TYPE <> @CostTypeNon
				//err = db.Table("mtr_location_stock A").
				//	Joins("INNER JOIN mtr_warehouse_master B ON A.company_id = B.company_id AND A.warehouse_id = b.warehouse_id").
				//	Where(`
				//	A.period_year = ?
				//	AND period_month = ?
				//	AND A.company_id = ?
				//	AND A.item_id = ?
				//	AND A.warehouse_group_id  = ?
				//	AND B.warehouse_costing_type_id = ?
				//	`, PeriodResponseSp.PeriodYear,
				//		PeriodResponseSp.PeriodMonth,
				//		GoodsReceiveEntities.CompanyId,
				//		GoodsReceiveEntities.WarehouseGroupId,
				//		CostingTypeNon.WarehouseCostingTypeId,
				//	).
				//	Select(`
				//		select A.quantity_in_transit,
				//		ISNULL(A.quantity_claim_in, 0) + ISNULL(A.quantity_robbing_out, 0) + ISNULL(A.quantity_assembly_out, 0)
				//	`).Row().Scan(&resdb.QuantityInTransit, &resdb.QuantityOnHand)
				//if err != nil {
				//	return false, &exceptions.BaseErrorResponse{
				//		StatusCode: http.StatusInternalServerError,
				//		Message:    fmt.Sprintf("error occured when fetch quantity instransit and quantity on hand"),
				//	}
				//}

				//validation if item group is inventory
				//so get item group by code first
				/*				ItemGroupInventoryId := 0
								err = db.Model(&masteritementities.ItemGroup{}).Where(masteritementities.ItemGroup{ItemGroupCode: "IN"}).
									Select("item_group_id").Scan(&ItemGroupInventoryId).Error
								if err != nil {
									return false, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "failed to get item group inventory with code IN",
									}
								}*/
				//quantityClaimIn := 1.0
				if GoodsReceivesItemGroupEntities.ItemGroupCode == "IN" {
					//not yet dev
					//SET @Qty_Claim_In = dbo.getQtyConvertion(@Source_Type, @CSR1_Item_Code, @CSR1_Qty_Variance)

					UomRate = resdb.QuantityClaimIn / (func() float64 {
						if resdb.QuantityVariance == 0 {
							return 1.0
						} else {
							return resdb.QuantityVariance
						}
					}())
					if UomRate == 0 {
						UomRate = 1
					}
					if resdb.QuantityClaimIn == 0 {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusBadRequest,
							Err:        fmt.Errorf("item Id %d does not have UOM conversion. Please define UOM Convertion", resdb.ItemId),
						}
					}
					if resdb.ItemPrice != 0 {
						resdb.PricePurchase = (resdb.QuantityVariance * resdb.ItemPrice) / resdb.QuantityClaimIn
					} else {
						resdb.ItemPrice = 0
					}
					//end of validation item group inventory
				}
				//start validation if type po or no
				if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "PO" ||
					goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "CL" {
					//hpp
					//IF EXISTS (SELECT ITEM_CODE FROM amGroupStock WHERE COMPANY_CODE = @Company_Code AND WHS_GROUP = @Whs_Group AND ITEM_CODE = @CSR1_Item_Code AND
					//PERIOD_YEAR = @Period_Year AND PERIOD_MONTH = @Period_Month)
					//BEGIN
					//SELECT @Hpp_Current = PRICE_CURRENT FROM amGroupStock WHERE COMPANY_CODE = @Company_Code AND WHS_GROUP = @Whs_Group AND ITEM_CODE = @CSR1_Item_Code AND
					//PERIOD_YEAR = @Period_Year AND PERIOD_MONTH = @Period_Month
					//END
					//ELSE
					//BEGIN
					//SET @Hpp_Current = 0
					//END
					//HppCurent := 0.0
					var groupStock masterentities.GroupStock
					err = db.Model(&groupStock).Where(masterentities.GroupStock{
						CompanyId:        GoodsReceiveEntities.CompanyId,
						WarehouseGroupId: GoodsReceiveEntities.WarehouseGroupId,
						PeriodYear:       PeriodResponseSp.PeriodYear,
						PeriodMonth:      PeriodResponseSp.PeriodMonth,
						ItemId:           resdb.ItemId,
					}).First(&groupStock).Error
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							resdb.HppCurrent = 0.0
						} else {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "failed to get group stock on error : " + err.Error(),
							}
						}
					}
					resdb.HppCurrent = groupStock.PriceCurrent
					if resdb.PricePurchase == 0 {
						resdb.PricePurchase = resdb.HppCurrent
					}
					if isCostingTypeNon == 1 {
						resdb.HppCurrent = 0
						resdb.QuantityOnHand = 0
					}
					if resdb.QuantityClaimIn+resdb.QuantityOnHand == 0 {
						resdb.HppNew = 0
					} else {
						resdb.HppNew = ((resdb.QuantityClaimIn * resdb.PricePurchase) + (resdb.QuantityOnHand * resdb.HppCurrent)) / (resdb.QuantityClaimIn + resdb.QuantityOnHand)
					}
					//end hpp
				} else {
					resdb.HppNew = 0
				}
				//execute stock transaction here
				//get by code first transaction type claim in
				//get by code first transaction type reason AP
				claimTypeInId := 0
				err = db.Model(&masterentities.StockTransactionType{}).
					Where(&masterentities.StockTransactionType{StockTransactionTypeCode: "CI"}).
					Select("stock_transaction_type_id").Scan(&claimTypeInId).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get claim type id claim in please check database",
					}
				}
				claimReasonApId := 0
				err = db.Model(&masterentities.StockTransactionReason{}).
					Where(&masterentities.StockTransactionReason{StockTransactionReasonCode: "AP"}).
					Select("stock_transaction_reason_id").Scan(&claimReasonApId).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get claim reason type AP please check database",
					}
				}
				payloadsStockTransaction := transactionsparepartpayloads.StockTransactionInsertPayloads{
					CompanyId:                    GoodsReceiveEntities.CompanyId,
					TransactionTypeId:            claimTypeInId,
					TransactionReasonId:          claimReasonApId,
					ReferenceId:                  GoodsReceiveEntities.GoodsReceiveSystemNumber,
					ReferenceDocumentNumber:      GoodsReceiveEntities.ReferenceDocumentNumber,
					ReferenceDate:                &GoodsReceiveEntities.GoodsReceiveDocumentDate,
					ReferenceWarehouseId:         GoodsReceiveEntities.WarehouseId,
					ReferenceWarehouseGroupId:    GoodsReceiveEntities.WarehouseGroupId,
					ReferenceLocationId:          resdb.WarehouseLocationClaimId,
					ReferenceItemId:              resdb.ItemId,
					ReferenceQuantity:            resdb.QuantityClaimIn,
					ReferenceUnitOfMeasurementId: resdb.UnitOfMeasurementStockId,
					ReferencePrice:               resdb.PricePurchase,
					ReferenceCurrencyId:          GoodsReceiveEntities.CurrencyId,
					TransactionCogs:              resdb.HppNew,
					ChangeNo:                     1,
					CreatedByUserId:              GoodsReceiveEntities.UpdatedByUserId,
					CreatedDate:                  time.Now(),
					UpdatedByUserId:              GoodsReceiveEntities.UpdatedByUserId,
					UpdatedDate:                  time.Now(),
				}
				stockTransactionUrl := config.EnvConfigs.AfterSalesServiceUrl + "stock-transaction"

				errStockTransaction := utils.Post(stockTransactionUrl, &payloadsStockTransaction, nil)
				if errStockTransaction != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "error on inserting stock transaction type claim in",
					}
				}
				//-- End Update Binning List
			}
			//IF @CSR1_Qty_Short > 0
			//BEGIN
			//IF @Item_Group = @ItemGroupInventory
			//BEGIN

		}
		if resdb.QuantityShort > 0 && GoodsReceivesItemGroupEntities.ItemGroupCode == "IN" {
			if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "PO" {
				//first get po data first for updating
				var PurchaseOrderDetailEntities transactionsparepartentities.PurchaseOrderDetailEntities
				err = db.Model(&PurchaseOrderDetailEntities).Where(transactionsparepartentities.PurchaseOrderDetailEntities{
					PurchaseOrderDetailSystemNumber: resdb.ReferenceSystemNumber,
				}).First(&PurchaseOrderDetailEntities).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusNotFound,
							Err:        fmt.Errorf("purchase order detail with reference number : %s is not found on table purchase order", GoodsReceiveEntities.ReferenceDocumentNumber),
						}
					}
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Err:        errors.New("failed to fetch from purchase order detail"),
					}
				}
				*PurchaseOrderDetailEntities.GoodsReceiveQuantity += resdb.QuantityShort
				PurchaseOrderDetailEntities.ChangeNo += 1
				PurchaseOrderDetailEntities.CreatedByUserId = GoodsReceiveEntities.CreatedByUserId
				PurchaseOrderDetailEntities.UpdatedByUserId = GoodsReceiveEntities.UpdatedByUserId
				err = db.Save(&PurchaseOrderDetailEntities).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "failed to save purchase order detail",
					}
				}
			}
			if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "CL" {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        errors.New("cannot have quantity short for item claim"),
				}
			}
			if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "WC" {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        errors.New("cannot have quantity short for item claim"),
				}
			}
		}
		//IF @CSR1_Qty_Grpo > 0
		//BEGIN
		//IF @Item_Group = @ItemGroupInventory

		if resdb.QuantityGoodsReceive > 0 {
			if GoodsReceivesItemGroupEntities.ItemGroupCode == "IN" {

				//SET @Qty_Purchase = dbo.getQtyConvertion(@Source_Type, @CSR1_Item_Code, @CSR1_Qty_Grpo)
				//resdb.QuantityPurchase
				UomRate = resdb.QuantityClaimIn / (func() float64 {
					if resdb.QuantityVariance == 0 {
						return 1.0
					} else {
						return resdb.QuantityVariance
					}
				}())
				if UomRate == 0 {
					UomRate = 1
				}
				if resdb.QuantityPurchase == 0 {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Err:        fmt.Errorf("item id %d doesnot have UOM conversion, please define uom conversion", resdb.ItemId),
					}
				}
				if resdb.ItemId != 0 {
					resdb.PricePurchase = (resdb.QuantityGoodsReceive * resdb.ItemPrice) / resdb.QuantityPurchase
				} else {
					resdb.PricePurchase = 0
				}
				if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "PO" ||
					goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "CL" {
					var groupStock masterentities.GroupStock
					err = db.Model(&groupStock).Where(masterentities.GroupStock{
						CompanyId:        GoodsReceiveEntities.CompanyId,
						WarehouseGroupId: GoodsReceiveEntities.WarehouseGroupId,
						PeriodYear:       PeriodResponseSp.PeriodYear,
						PeriodMonth:      PeriodResponseSp.PeriodMonth,
						ItemId:           resdb.ItemId,
					}).First(&groupStock).Error
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							resdb.HppCurrent = 0.0
						} else {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "failed to get group stock on error : " + err.Error(),
							}
						}
					}
					resdb.HppCurrent = groupStock.PriceCurrent
					if resdb.PricePurchase == 0 {
						resdb.PricePurchase = resdb.HppCurrent
					}
					if isCostingTypeNon == 1 {
						resdb.HppCurrent = 0
						resdb.QuantityOnHand = 0
					}
					if resdb.QuantityClaimIn+resdb.QuantityOnHand == 0 {
						resdb.HppNew = 0
					} else {
						resdb.HppNew = ((resdb.QuantityClaimIn * resdb.PricePurchase) + (resdb.QuantityOnHand * resdb.HppCurrent)) / (resdb.QuantityClaimIn + resdb.QuantityOnHand)
					}
					//end hpp
				} else {
					resdb.HppNew = 0
				}
				claimTypePurchaseId := 0
				err = db.Model(&masterentities.StockTransactionType{}).
					Where(masterentities.StockTransactionType{StockTransactionTypeCode: "PU"}).
					Select("stock_transaction_type_id").First(&claimTypePurchaseId).Error
				if err != nil {

					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "failed to get stock transaction type purchase on database",
					}
				}
				//cek back order on po header data
				var PurchaseOrderEntities transactionsparepartentities.PurchaseOrderEntities
				err = db.Model(&PurchaseOrderEntities).Where(transactionsparepartentities.PurchaseOrderEntities{
					PurchaseOrderSystemNumber: GoodsReceiveEntities.ReferenceSystemNumber,
				}).First(&PurchaseOrderEntities).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusNotFound,
							Err:        fmt.Errorf("purchase order with reference number : %s is not found on table purchase order", GoodsReceiveEntities.ReferenceDocumentNumber),
						}
					}
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Err:        errors.New("failed to fetch from purchase order"),
					}
				}
				stockReasonCode := ""
				if PurchaseOrderEntities.BackOrder {
					stockReasonCode = "NL"
				} else {
					stockReasonCode = "BO"
				}
				if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "WC" {
					stockReasonCode = "WP"
				}
				//get type reason code
				claimReasonId := 0
				err = db.Model(&masterentities.StockTransactionReason{}).
					Where(&masterentities.StockTransactionReason{StockTransactionReasonCode: stockReasonCode}).
					Select("stock_transaction_reason_id").First(&claimReasonId).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get claim reason type AP please check database",
					}
				}
				payloadsStockTransaction := transactionsparepartpayloads.StockTransactionInsertPayloads{
					CompanyId:                    GoodsReceiveEntities.CompanyId,
					TransactionTypeId:            claimTypePurchaseId,
					TransactionReasonId:          claimReasonId,
					ReferenceId:                  GoodsReceiveEntities.GoodsReceiveSystemNumber,
					ReferenceDocumentNumber:      GoodsReceiveEntities.ReferenceDocumentNumber,
					ReferenceDate:                &GoodsReceiveEntities.GoodsReceiveDocumentDate,
					ReferenceWarehouseId:         GoodsReceiveEntities.WarehouseId,
					ReferenceWarehouseGroupId:    GoodsReceiveEntities.WarehouseGroupId,
					ReferenceLocationId:          resdb.WarehouseLocationClaimId,
					ReferenceItemId:              resdb.ItemId,
					ReferenceQuantity:            resdb.QuantityClaimIn,
					ReferenceUnitOfMeasurementId: resdb.UnitOfMeasurementStockId,
					ReferencePrice:               resdb.PricePurchase,
					ReferenceCurrencyId:          GoodsReceiveEntities.CurrencyId,
					TransactionCogs:              resdb.HppNew,
					ChangeNo:                     1,
					CreatedByUserId:              GoodsReceiveEntities.UpdatedByUserId,
					CreatedDate:                  time.Now(),
					UpdatedByUserId:              GoodsReceiveEntities.UpdatedByUserId,
					UpdatedDate:                  time.Now(),
				}
				stockTransactionUrl := config.EnvConfigs.AfterSalesServiceUrl + "stock-transaction"

				errStockTransaction := utils.Post(stockTransactionUrl, &payloadsStockTransaction, nil)
				if errStockTransaction != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "error on inserting stock transaction type claim in",
					}
				}
				//--==Begin Update Item Claim ==-- IF REFERENCE TYPE CLAIM
				//------------------------////////////////////////////-------------
				if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "CL" {
					if resdb.QuantityDeliveryOrder-resdb.QuantityReference > 0 {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusBadRequest,
							Err:        errors.New("delivery order quantity is bigger than claim quantity"),
						}
					}
					err = db.Model(&transactionsparepartentities.ItemClaimDetail{}).
						Where(transactionsparepartentities.ItemClaimDetail{ItemClaimDetailId: resdb.ReferenceSystemNumber}).
						Updates(map[string]interface{}{
							"quantity_goods_receive": gorm.Expr("quantity_goods_receive + ?", resdb.QuantityGoodsReceive),
						}).Error
					if err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "failed to update item claim detail"}
					}
				}
			} //end 1330
			if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "PO" {
				err = db.Model(&transactionsparepartentities.PurchaseOrderDetailEntities{}).
					Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderDetailSystemNumber: resdb.ReferenceSystemNumber}).
					Updates(map[string]interface{}{
						"goods_receive_quantity": gorm.Expr("goods_receive_quantity + ?", resdb.QuantityGoodsReceive),
						"change_no":              gorm.Expr("change_no + ?", 1),
						"updated_by_user_id":     GoodsReceiveEntities.UpdatedByUserId,
						"updated_date":           time.Now(),
					}).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "failed to update purchase order detail"}
				}
				if GoodsReceivesItemGroupEntities.ItemGroupCode == "OJ" {
					//update wo but getting data from po- >pr- >wo
					var purchaseOrderDetail transactionsparepartentities.PurchaseOrderDetailEntities
					err = db.Model(&purchaseOrderDetail).
						Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderDetailSystemNumber: resdb.ReferenceSystemNumber}).
						First(&purchaseOrderDetail).Error
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusNotFound,
								Err:        errors.New("purchase order detail is not found"),
							}
						}
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "error occured when retreiving purchase order detail",
						}
					}
					//get pr
					var purchaseRequestDetail transactionsparepartentities.PurchaseRequestDetail
					err = db.Model(&purchaseRequestDetail).
						Where(transactionsparepartentities.PurchaseRequestDetail{PurchaseRequestDetailSystemNumber: purchaseOrderDetail.PurchaseRequestDetailSystemNumber}).
						First(&purchaseRequestDetail).Error
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusNotFound,
								Err:        errors.New("purchase request detail is not found in purchase request"),
							}
						}
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "error occured when retreiving purchase request detail",
						}
					}
					//get supply quantity and total cogs from wtworkorder 2
					var workOrderDetail transactionworkshopentities.WorkOrderDetail
					err = db.Model(&workOrderDetail).Where(transactionworkshopentities.WorkOrderDetail{WorkOrderDetailId: purchaseRequestDetail.ReferenceSystemNumber}).
						Select("supply_quantity,total_cost_of_goods_sold").
						First(&workOrderDetail).Error
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusNotFound,
								Err:        errors.New("work order detail is not found"),
							}
						}
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "error on getting supply quantity and total cogs from work order detail",
						}
					}
					totalCogs := (workOrderDetail.TotalCostOfGoodsSold * workOrderDetail.SupplyQuantity) + (*purchaseOrderDetail.ItemPrice*resdb.QuantityGoodsReceive)/(workOrderDetail.SupplyQuantity+resdb.QuantityGoodsReceive)

					workOrderDetail.SupplyQuantity += resdb.QuantityGoodsReceive
					workOrderDetail.TotalCostOfGoodsSold = totalCogs

					//save work order update
					err = db.Save(&workOrderDetail).Error
					if err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "error on updating work order detail",
						}
					}
				}
			} else if goodsReceiveReferenceTypeEntities.ReferenceTypeGoodReceiveCode == "WC" {
				var workOrderDetail transactionworkshopentities.WorkOrderDetail
				err = db.Model(&workOrderDetail).Where(transactionworkshopentities.WorkOrderDetail{WorkOrderDetailId: resdb.ReferenceSystemNumber}).
					First(&workOrderDetail).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusNotFound,
							Err:        errors.New("work order detail is not found type claim"),
						}
					}
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "error on getting supply quantity and total cogs from work order detail type claim",
					}
				}
				workOrderDetail.SupplyQuantity += resdb.QuantityGoodsReceive
				//save work order update
				err = db.Save(&workOrderDetail).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "error on updating work order detail",
					}
				}
			}

			//UPDATE dbo.atItemGRPO1
			//SET UOM_RATE = @Uom_Rate,
			//	ITEM_TOTAL = QTY_GRPO * ITEM_PRICE,
			//	ITEM_TOTAL_BASE_AMOUNT = QTY_GRPO * ITEM_PRICE * @Ccy_Exch_Rate,
			//	CHANGE_NO = CHANGE_NO + 1,
			//	CHANGE_DATETIME = GETDATE(),
			//	CHANGE_USER_ID = @Change_User_Id
			//WHERE GRPO_SYS_NO = @Grpo_Sys_No AND GRPO_LINE_NO = @CSR1_Grpo_Line_No
		}
		//get exchange rate for item total base amount

		//resdb.GoodsReceiveDetailSystemNumber
		err = db.Model(&goodsReceiveDetailEntities).
			Where(transactionsparepartentities.GoodsReceiveDetail{GoodsReceiveDetailSystemNumber: resdb.GoodsReceiveDetailSystemNumber}).
			Updates(map[string]interface{}{
				"unit_of_measurement_rate": UomRate,
				"item_total":               gorm.Expr("quantity_goods_receive * item_price"),
				"item_total_base_amount":   gorm.Expr("quantity_goods_receive * item_price * ?", GoodsReceiveEntities.CurrencyExchangeRate),
				"change_no":                gorm.Expr("change_no + 1"),
				"updated_by_user_id":       GoodsReceiveEntities.UpdatedByUserId,
				"updated_date":             time.Now(),
			}).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error on updating goods receive detail",
			}
		}
		//untuk else nya dms live pengecekan berulang tidak perlu dibuat --notes devin
	}
	////update purchase order
	//IF NOT EXISTS(SELECT PO_LINE FROM atItemPO1 WHERE PO_SYS_NO = @Ref_Sys_No AND ITEM_QTY <> GRPO_QTY)
	//BEGIN
	//UPDATE atItemPO0
	//SET PO_STATUS = dbo.getVariableValue('APPROVAL_CLOSED'),
	//	CHANGE_NO = CHANGE_NO + 1,
	//	CHANGE_DATETIME = @Change_Datetime,
	//	CHANGE_USER_ID = @Change_User_Id
	//WHERE PO_SYS_NO = @Ref_Sys_No
	//END
	//
	//update purchase order
	var DocResponse generalservicepayloads.ApprovalStatusResponses

	DocumentStatusUrl := config.EnvConfigs.GeneralServiceUrl + "approval-status-codes/99"
	if err := utils.Get(DocumentStatusUrl, &DocResponse, nil); err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error on getting approval status codes",
		}
	}

	err = db.Model(&transactionsparepartentities.PurchaseOrderEntities{}).
		Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: GoodsReceiveEntities.ReferenceSystemNumber}).
		Updates(map[string]interface{}{
			"purchase_order_status_id": DocResponse.ApprovalStatusId,
			"change_no":                gorm.Expr("change_no + 1"),
			"updated_date":             time.Now(),
			"updated_by_user_id":       GoodsReceiveEntities.UpdatedByUserId,
		}).
		Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error on update status purchase order",
		}
	}
	//not yet journal exec
	//IF EXISTS (SELECT GRPO_SYS_NO FROM atItemGRPO0 WHERE GRPO_SYS_NO = @Grpo_Sys_No AND REF_TYPE <> @RefTypeWcf)
	//BEGIN
	//SET @Process_Code = dbo.getVariableValue('GL_PROCESS_CODE_SP_GRPO')
	//SELECT @Process_Code,@Profit_Center,@Trx_Type,@Event_No,@Grpo_Doc_Date,@Grpo_Sys_No,@Grpo_Doc_No,@Change_User_Id
	//EXEC usp_comJournalAction
	//@Process_Code =@Process_Code,
	//@Cpc_Code = @Profit_Center,
	//@Trx_Type = @Trx_Type,
	//@Event_No = @Event_No,
	//@Journal_Sys_No = @Journal_Sys_No Output,
	//@Journal_Date = @Grpo_Doc_Date,
	//@Ref_Sys_No = @Grpo_Sys_No,
	//@Ref_Doc_No = @Grpo_Doc_No,
	//@Creation_User_Id = @Change_User_Id
	//
	//--IF ISNULL(@Journal_sys_no,0) = 0 --:GH bisa ga create journal
	//--BEGIN
	//--	RAISERROR('Journal is not created',16,1)
	//--END
	//END

	return true, nil
}

func getCostingTypeByCode(db *gorm.DB, code string) (masterwarehouseentities.WarehouseCostingType, *exceptions.BaseErrorResponse) {
	var costingType masterwarehouseentities.WarehouseCostingType
	err := db.Model(&costingType).
		Where(masterwarehouseentities.WarehouseCostingType{WarehouseCostingTypeCode: code}).
		First(&costingType).Error
	if err != nil {
		return costingType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	return costingType, nil
}
func GenerateDocumentNumber(tx *gorm.DB, id int) (string, *exceptions.BaseErrorResponse) {
	var GoodsReceivesEntities transactionsparepartentities.GoodsReceive

	err1 := tx.Model(&transactionsparepartentities.GoodsReceive{}).
		Where("goods_receive_system_number = ?", id).
		First(&GoodsReceivesEntities).
		Error
	if err1 != nil {
		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve goods receive from the database: %v", err1)}
	}

	var GoodsReceive transactionsparepartentities.GoodsReceive
	var brandResponse transactionworkshoppayloads.BrandDocResponse

	GoodsReceiveId := GoodsReceivesEntities.GoodsReceiveSystemNumber

	// Get the work order based on the work order system number
	err := tx.Model(&transactionsparepartentities.GoodsReceive{}).Where("goods_receive_system_number = ?", GoodsReceiveId).First(&GoodsReceive).Error
	if err != nil {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to goods receive order from the database: %v", err)}
	}

	if GoodsReceive.BrandId == 0 {

		return "", &exceptions.BaseErrorResponse{Message: "brand_id is missing in the work order. Please ensure the work order has a valid brand_id before generating document number."}
	}

	// Get the last work order based on the work order system number
	var LastGoodsReceive transactionsparepartentities.GoodsReceive
	err = tx.Model(&transactionsparepartentities.GoodsReceive{}).
		Where("brand_id = ?", GoodsReceive.BrandId).
		First(&LastGoodsReceive).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve last goods receive: %v", err)}
	}

	currentTime := time.Now()
	month := int(currentTime.Month())
	year := currentTime.Year() % 100 // Use last two digits of the year

	// fetch data brand from external api
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(GoodsReceive.BrandId)
	errUrl := utils.Get(brandUrl, &brandResponse, nil)
	if errUrl != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrl,
		}
	}

	// Check if BrandCode is not empty before using it
	if brandResponse.BrandCode == "" {
		return "", &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Message: "Brand code is empty"}
	}

	// Get the initial of the brand code
	brandInitial := brandResponse.BrandCode[0]

	// Handle the case when there is no last work order or the format is invalid
	newDocumentNumber := fmt.Sprintf("SPRS/%c/%02d/%02d/00001", brandInitial, month, year)
	if LastGoodsReceive.GoodsReceiveSystemNumber != 0 {
		lastWorkOrderDate := LastGoodsReceive.GoodsReceiveDocumentDate
		lastWorkOrderYear := lastWorkOrderDate.Year() % 100

		// Check if the last work order is from the same year
		if lastWorkOrderYear == year {
			lastWorkOrderCode := LastGoodsReceive.GoodsReceiveDocumentNumber
			codeParts := strings.Split(lastWorkOrderCode, "/")
			if len(codeParts) == 5 {
				lastWorkOrderNumber, err := strconv.Atoi(codeParts[4])
				if err == nil {
					newWorkOrderNumber := lastWorkOrderNumber + 1
					newDocumentNumber = fmt.Sprintf("SPRS/%c/%02d/%02d/%05d", brandInitial, month, year, newWorkOrderNumber)
				} else {
					log.Printf("failed to parse last work order code: %v\n", err)
				}
			} else {
				log.Println("Invalid last work order code format")
			}
		}
	}

	log.Printf("New document number: %s", newDocumentNumber)
	return newDocumentNumber, nil
}
