package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	financeservice "after-sales/api/payloads/cross-service/finance-service"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
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

	if locationStockEntities.StockOpname == true && locationStockEntities.WarehouseGroupId == GoodsReceiveEntities.WarehouseGroupId {
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
		var CheckDuplicateItemClaim int = 0
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
	if goodsReceiveHeader.UseInTransitWarehouse == true {
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
	if locationStockEntities.StockOpname == true && locationStockEntities.WarehouseGroupId == goodsReceiveHeader.WarehouseGroupId {
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
				Err:        errors.New(fmt.Sprintf("goods receive with id : %d is not found", GoodsReceiveId)),
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
	err = db.Model(&transactionsparepartentities.GoodsReceiveDetail{}).
		Select("count(goods_receive_detail_system_number)").
		Where(transactionsparepartentities.GoodsReceiveDetail{WarehouseLocationId: 0, GoodsReceiveSystemNumber: GoodsReceiveId}).
		Scan(&isExist).Error
	if isExist > 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New(fmt.Sprintf("location Code must be filled", GoodsReceiveId)),
		}
	}
	//IF ((SELECT COUNT(GRPO_LINE_NO) FROM atItemGRPO1 WITH(NOLOCK) WHERE GRPO_SYS_NO = ISNULL(@Grpo_Sys_No,0) AND ISNULL(LOC_CLAIM_CODE,'') = '' AND QTY_DAMAGE + QTY_OVER + QTY_SHORT + QTY_WRONG > 0) > 0)
	//BEGIN
	//RAISERROR('Location Claim Code must be filled for Item that has Claim',16,1)
	//RETURN 0
	//END
	isExist = 0
	err = db.Model(&transactionsparepartentities.GoodsReceiveDetail{}).
		Select("count(goods_receive_detail_system_number)").
		Where("goods_receive_system_number = ? AND quantity_short+quantity_damage +quantity_over+ quantity_wrong > 0", GoodsReceiveId).
		Scan(&isExist).Error
	if isExist > 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("location Claim Code must be filled for Item that has Claim"),
		}
	}

}
