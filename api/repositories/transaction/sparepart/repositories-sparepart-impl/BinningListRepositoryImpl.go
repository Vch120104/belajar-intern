package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	generalservicepayloads "after-sales/api/payloads/crossservice/generalservice"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type BinningListRepositoryImpl struct {
}

func NewbinningListRepositoryImpl() transactionsparepartrepository.BinningListRepository {
	return &BinningListRepositoryImpl{}
}
func (b *BinningListRepositoryImpl) GetAllBinningListDetailWithPagination(db *gorm.DB, filter []utils.FilterCondition, paginations pagination.Pagination, binningListId int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var binningDetailEntities []transactionsparepartentities.BinningStockDetail
	var Responses []transactionsparepartpayloads.BinningListGetByIdResponses
	joinTable := db.Table("trx_binning_list_stock_detail A").
		Select(`
				A.binning_system_number,
				A.binning_line_number,
				C.item_code,
				C.item_name,
				A.uom_id,
				WL.warehouse_location_code,
				A.item_price,
				A.purchase_order_quantity,
				A.delivery_order_quantity,
				A.reference_system_number,
				A.reference_line_number,
				D.item_code original_item_code,
				D.item_name original_item_code,
				B.purchase_order_document_number
			`).
		Joins("LEFT OUTER JOIN mtr_item C ON C.item_id = A.item_id").
		Joins("LEFT OUTER JOIN mtr_item D ON D.item_id = A.original_item_id").
		Joins("LEFT OUTER JOIN trx_item_purchase_order B ON A.reference_system_number = B.purchase_order_system_number").
		Joins(`	LEFT OUTER JOIN trx_item_purchase_order_detail F ON A.reference_system_number = b.purchase_order_system_number
                     	AND A.item_purchase_order_detail_id = F.purchase_order_detail_system_number`).
		Joins("INNER JOIN mtr_warehouse_location WL ON A.warehouse_location_id = WL.warehouse_location_id").
		Where("A.binning_system_number = ?", binningListId)
	WhereQuery := utils.ApplyFilter(joinTable, filter)
	err := WhereQuery.Scopes(pagination.Paginate(&binningDetailEntities, &paginations, WhereQuery)).Order("binning_system_number").Scan(&Responses).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return paginations, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed To Get Paginate Binning List"),
		}
	}

	if len(Responses) == 0 {
		paginations.Rows = []string{}
		return paginations, &exceptions.BaseErrorResponse{}
	}
	paginations.Rows = Responses
	return paginations, nil
}

func (b *BinningListRepositoryImpl) GetAllBinningListWithPagination(db *gorm.DB, rdb *redis.Client, filter []utils.FilterCondition, paginations pagination.Pagination, ctx context.Context) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var binningEntities []transactionsparepartentities.BinningStock
	var BinningResponses []transactionsparepartpayloads.BinningListGetPaginationResponse
	joinTable := db.Model(&binningEntities)
	WhereQuery := utils.ApplyFilter(joinTable, filter)
	err := WhereQuery.Scopes(pagination.Paginate(&binningEntities, &paginations, WhereQuery)).Order("binning_system_number").Scan(&binningEntities).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return paginations, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed To Get Paginate Binning List"),
		}
	}
	if len(binningEntities) == 0 {
		paginations.Rows = []string{}
		return paginations, nil
	}

	for _, binningEntity := range binningEntities {
		var SupplierData generalservicepayloads.SupplierMasterCrossServicePayloads
		SupplierDataUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(binningEntity.SupplierId)
		if err := utils.Get(SupplierDataUrl, &SupplierData, nil); err != nil {
			return paginations, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Supplier data from external service" + err.Error(),
				Err:        err,
			}
		}
		exists, _ := rdb.Exists(ctx, strconv.Itoa(binningEntity.BinningDocumentStatusId)).Result()

		if exists == 0 {
			fmt.Println("Failed To Get Status on redis... queue for external service")
			var DocResponse generalservicepayloads.DocumentStatusPayloads
			DocumentStatusUrl := config.EnvConfigs.GeneralServiceUrl + "document-status/" + strconv.Itoa(binningEntity.BinningDocumentStatusId)
			if err := utils.Get(DocumentStatusUrl, &DocResponse, nil); err != nil {
				return paginations, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errors.New("failed To fetch Document Status From External Service"),
				}
			}

			rdb.Set(ctx, strconv.Itoa(binningEntity.BinningDocumentStatusId), DocResponse.DocumentStatusCode, 1*time.Hour)
		}
		StatusCode, _ := rdb.Get(ctx, strconv.Itoa(binningEntity.BinningDocumentStatusId)).Result()

		//return DocResponse.ApprovalStatusId
		BinningResponse := transactionsparepartpayloads.BinningListGetPaginationResponse{
			BinningSystemNumber:         binningEntity.BinningSystemNumber,
			BinningDocumentStatusId:     binningEntity.BinningDocumentStatusId,
			BinningDocumentNumber:       binningEntity.BinningDocumentNumber,
			BinningDocumentDate:         binningEntity.BinningDocumentDate,
			ReferenceDocumentNumber:     binningEntity.ReferenceDocumentNumber,
			SupplierInvoiceNumber:       binningEntity.SupplierInvoiceNumber,
			SupplierName:                SupplierData.SupplierName,
			SupplierCaseNumber:          binningEntity.SupplierCaseNumber,
			Status:                      StatusCode,
			SupplierDeliveryOrderNumber: binningEntity.SupplierDeliveryOrderNumber,
		}
		BinningResponses = append(BinningResponses, BinningResponse)
	}
	paginations.Rows = BinningResponses
	return paginations, nil
}
func (b *BinningListRepositoryImpl) GetBinningListById(db *gorm.DB, BinningStockId int) (transactionsparepartpayloads.BinningListGetByIdResponse, *exceptions.BaseErrorResponse) {
	BinningStockEntities := transactionsparepartentities.BinningStock{}
	Response := transactionsparepartpayloads.BinningListGetByIdResponse{}

	err := db.Model(&BinningStockEntities).
		First(&BinningStockEntities, transactionsparepartentities.BinningStock{BinningSystemNumber: BinningStockId}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Response, &exceptions.BaseErrorResponse{

				StatusCode: http.StatusNotFound,
				Message:    "Binning Stock Not Found",
				Err:        gorm.ErrRecordNotFound}
		}
		return Response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Get Data Failed",
			Err:        err,
		}
	}
	BinningReferenceTypeEntities := masterentities.BinningReferenceTypeMaster{}
	//get binning reference type data
	err = db.Model(&BinningReferenceTypeEntities).
		First(&BinningReferenceTypeEntities, masterentities.BinningReferenceTypeMaster{BinningReferenceTypeId: BinningStockEntities.BinningReferenceTypeId}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			BinningReferenceTypeEntities.BinningReferenceTypeCode = ""

		} else {
			return Response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Get Data Failed",
				Err:        errors.New("error on getting binning reference type entities"),
			}
		}
	}
	//get warehouse group code
	warehousegroup := masterwarehouseentities.WarehouseGroup{}

	err = db.Model(&warehousegroup).
		Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupId: BinningStockEntities.WarehouseGroupId}).
		First(&warehousegroup).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			warehousegroup.WarehouseGroupCode = ""
		} else {
			return Response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Get Data Warehouse Failed",
				Err:        err,
			}
		}
	}
	//get warehouse code
	WarehouseMasterEntities := masterwarehouseentities.WarehouseMaster{}
	err = db.Model(&WarehouseMasterEntities).
		Where(masterwarehouseentities.WarehouseMaster{WarehouseGroupId: BinningStockEntities.WarehouseGroupId}).
		First(&WarehouseMasterEntities).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			WarehouseMasterEntities.WarehouseCode = ""
		} else {
			return Response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Get Data Warehouse Failed",
				Err:        err,
			}
		}
	}
	//get supplier data
	var SupplierData generalservicepayloads.SupplierMasterCrossServicePayloads
	SupplierDataUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(BinningStockEntities.SupplierId)
	if err := utils.Get(SupplierDataUrl, &SupplierData, nil); err != nil {
		return Response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Supplier data from external service" + err.Error(),
			Err:        err,
		}
	}
	//get PO Data

	//Item Group Name
	ItemGroupName := ""
	err = db.Table("trx_item_purchase_order A").Select("B.item_group_name").
		Joins("JOIN mtr_item_group B ON A.item_group_id = B.item_group_id").
		Where("A.purchase_order_system_number = ?", BinningStockEntities.ReferenceSystemNumber).
		Scan(&ItemGroupName).Error
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return Response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to Get Group Data",
			Err:        err,
		}
	}

	Response = transactionsparepartpayloads.BinningListGetByIdResponse{
		CompanyId:                   BinningStockEntities.CompanyId,
		BinningSystemNumber:         BinningStockEntities.BinningSystemNumber,
		BinningDocumentStatusId:     BinningStockEntities.BinningDocumentStatusId,
		BinningDocumentNumber:       BinningStockEntities.BinningDocumentNumber,
		BinningDocumentDate:         BinningStockEntities.BinningDocumentDate,
		BinningReferenceType:        BinningReferenceTypeEntities.BinningReferenceTypeCode,
		ReferenceSystemNumber:       BinningStockEntities.ReferenceSystemNumber,
		ReferenceDocumentNumber:     BinningStockEntities.ReferenceDocumentNumber,
		WarehouseGroupCode:          warehousegroup.WarehouseGroupCode,
		WarehouseCode:               WarehouseMasterEntities.WarehouseCode,
		SupplierCode:                SupplierData.SupplierCode,
		SupplierDeliveryOrderNumber: BinningStockEntities.SupplierDeliveryOrderNumber,
		SupplierInvoiceNumber:       BinningStockEntities.SupplierInvoiceNumber,
		SupplierInvoiceDate:         BinningStockEntities.SupplierInvoiceDate,
		SupplierFakturPajakNumber:   BinningStockEntities.SupplierFakturPajakNumber,
		SupplierFakturPajakDate:     BinningStockEntities.SupplierFakturPajakDate,
		SupplierDeliveryPerson:      BinningStockEntities.SupplierDeliveryPerson,
		SupplierCaseNumber:          BinningStockEntities.SupplierCaseNumber,
		ItemGroup:                   ItemGroupName, //ambil dari po
		CreatedByUserId:             BinningStockEntities.CreatedByUserId,
		CreatedDate:                 BinningStockEntities.CreatedDate,
		UpdatedByUserId:             BinningStockEntities.UpdatedByUserId,
		UpdatedDate:                 BinningStockEntities.UpdatedDate,
		ChangeNo:                    BinningStockEntities.ChangeNo,
	}
	return Response, nil
}

func (b *BinningListRepositoryImpl) InsertBinningListHeader(db *gorm.DB, payloads transactionsparepartpayloads.BinningListInsertPayloads) (transactionsparepartentities.BinningStock, *exceptions.BaseErrorResponse) {
	Entities := transactionsparepartentities.BinningStock{
		CompanyId:                   payloads.CompanyId,
		BinningDocumentStatusId:     payloads.BinningDocumentStatusId,
		BinningDocumentNumber:       payloads.BinningDocumentNumber,
		BinningDocumentDate:         payloads.BinningDocumentDate,
		BinningReferenceTypeId:      payloads.BinningReferenceTypeId,
		ReferenceSystemNumber:       payloads.ReferenceSystemNumber,
		ReferenceDocumentNumber:     payloads.ReferenceDocumentNumber,
		WarehouseGroupId:            payloads.WarehouseGroupId,
		WarehouseId:                 payloads.WarehouseId,
		SupplierId:                  payloads.SupplierId,
		SupplierDeliveryOrderNumber: payloads.SupplierDeliveryOrderNumber,
		SupplierInvoiceNumber:       payloads.SupplierInvoiceNumber,
		SupplierInvoiceDate:         payloads.SupplierInvoiceDate,
		SupplierFakturPajakNumber:   payloads.SupplierFakturPajakNumber,
		SupplierFakturPajakDate:     payloads.SupplierFakturPajakDate,
		SupplierDeliveryPerson:      payloads.SupplierDeliveryPerson,
		SupplierCaseNumber:          payloads.SupplierCaseNumber,
		BinningTypeId:               payloads.BinningTypeId,
		CurrencyId:                  payloads.CurrencyId,
		ExchangeId:                  payloads.ExchangeId,
		CreatedByUserId:             payloads.CreatedByUserId,
		CreatedDate:                 payloads.CreatedDate,
		UpdatedByUserId:             payloads.UpdatedByUserId,
		UpdatedDate:                 payloads.UpdatedDate,
		ChangeNo:                    0,
	}
	err := db.Create(&Entities).Scan(&Entities).Error
	if err != nil {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed To Insert",
			Err:        err,
		}
	}
	return Entities, nil
}
func (b *BinningListRepositoryImpl) UpdateBinningListHeader(db *gorm.DB, payloads transactionsparepartpayloads.BinningListSavePayload) (transactionsparepartentities.BinningStock, *exceptions.BaseErrorResponse) {
	Entities := transactionsparepartentities.BinningStock{}
	err := db.Model(&Entities).Where(transactionsparepartentities.BinningStock{BinningSystemNumber: payloads.BinningSystemNumber}).
		Scan(&Entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Binning Is Not Found Please Check Input",
			}
		}
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("Error On retreiving Data"),
		}
	}
	Entities.SupplierDeliveryOrderNumber = payloads.SupplierDeliveryOrderNumber
	Entities.SupplierInvoiceNumber = payloads.SupplierInvoiceNumber
	Entities.SupplierInvoiceDate = payloads.SupplierInvoiceDate
	Entities.SupplierFakturPajakNumber = payloads.SupplierFakturPajakNumber
	Entities.SupplierFakturPajakDate = payloads.SupplierFakturPajakDate
	Entities.SupplierDeliveryOrderNumber = payloads.SupplierDeliveryOrderNumber
	Entities.ChangeNo += 1
	Entities.UpdatedDate = payloads.UpdatedDate
	Entities.UpdatedByUserId = payloads.UpdatedByUserId
	Entities.CurrencyId = payloads.CurrencyId
	Entities.ExchangeId = payloads.ExchangeId
	Entities.SupplierId = payloads.SupplierId

	err = db.Save(&Entities).Error
	if err != nil {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed To Update Data : " + err.Error(),
			Err:        err,
		}
	}
	return Entities, nil
}
func (b *BinningListRepositoryImpl) GetBinningListDetailById(db *gorm.DB, BinningDetailId int) (transactionsparepartpayloads.BinningListGetByIdResponses, *exceptions.BaseErrorResponse) {
	//var BinningDetailEntities transactionsparepartentities.BinningStockDetail
	var Responses transactionsparepartpayloads.BinningListGetByIdResponses
	err := db.Table("trx_binning_list_stock_detail A").
		Select(`
				A.binning_system_number,
				A.binning_line_number,
				C.item_code,
				C.item_name,
				H.uom_code,
				WL.warehouse_location_code,
				A.item_price,
				A.purchase_order_quantity,
				A.delivery_order_quantity,
				A.reference_system_number,
				A.reference_line_number,
				D.item_code original_item_code,
				D.item_name original_item_code,
				B.purchase_order_document_number
			`).
		Joins("LEFT OUTER JOIN mtr_item C ON C.item_id = A.item_id").
		Joins("LEFT OUTER JOIN mtr_item D ON D.item_id = A.original_item_id").
		Joins("LEFT OUTER JOIN trx_item_purchase_order B ON A.reference_system_number = B.purchase_order_system_number").
		Joins(`	LEFT OUTER JOIN trx_item_purchase_order_detail F ON A.reference_system_number = b.purchase_order_system_number
                     	AND A.item_purchase_order_detail_id = F.purchase_order_detail_system_number`).
		Joins("INNER JOIN mtr_warehouse_location WL ON A.warehouse_location_id = WL.warehouse_location_id").
		Joins("LEFT OUTER JOIN mtr_uom H ON H.uom_id = A.uom_id").
		Where("A.binning_detail_system_number = ?", BinningDetailId).
		Scan(&Responses).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Responses, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Binning Detail Is Not Found Please Check Input",
				Err:        err,
			}
		}
		return Responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed To Get Binning List Detail error on : " + err.Error()),
		}
	}
	return Responses, nil
}

// IF @Option = 0
// [uspg_atBinningStock1_Insert]
func (b *BinningListRepositoryImpl) InsertBinningListDetail(db *gorm.DB, payloads transactionsparepartpayloads.BinningListDetailPayloads) (transactionsparepartentities.BinningStockDetail, *exceptions.BaseErrorResponse) {
	BinningRefTypeCode := ""
	BinningEntities := transactionsparepartentities.BinningStock{}
	BinningDetailEntities := transactionsparepartentities.BinningStockDetail{}
	PurchaseOrderDetailEntities := transactionsparepartentities.PurchaseOrderDetailEntities{}
	err := db.Model(&BinningEntities).
		Where(transactionsparepartentities.BinningStock{BinningSystemNumber: payloads.BinningSystemNumber}).
		Scan(&BinningEntities).Error
	//err := db.Model(&BinningEntities).
	//	Preload("BinningReferenceType"). // Preload the associated BinningReferenceType
	//	Where("binning_system_number = ?", payloads.BinningSystemNumber).
	//	Scan(&BinningEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "Binning Header Is Not Found Please Check Input",
			}
		}
		return BinningDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed To Get Binning error on" + err.Error(),
			Err:        errors.New("failed To Get Binning error on : " + err.Error()),
		}
	}
	err = db.Model(masterentities.BinningReferenceTypeMaster{}).
		Select("binning_reference_type_code").
		Where(masterentities.BinningReferenceTypeMaster{BinningReferenceTypeId: BinningEntities.BinningReferenceTypeId}).
		Scan(&BinningRefTypeCode).
		Error
	if err != nil {
		return BinningDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed To Get Binning Reference Type Code",
			Err:        errors.New("failed To Get Binning Reference Type Code"),
		}
	}
	if BinningRefTypeCode == "PO" {
		//IF @Ref_Type = @Ref_Type_PO AND EXISTS(SELECT PO_LINE FROM atItemPO1 WHERE PO_SYS_NO = @Ref_Sys_No
		//AND PO_LINE = @Ref_Line AND ISNULL(BINNING_QTY,0) + ISNULL(@Do_Qty,0) > ISNULL(ITEM_QTY,0))
		//BEGIN
		//RAISERROR('Item has exceeded Reference Qty',16,1)
		//RETURN 0
		//END
		err = db.Model(&PurchaseOrderDetailEntities).
			Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderDetailSystemNumber: payloads.ReferenceDetailSystemNumber}).
			Scan(&PurchaseOrderDetailEntities).Error
		if err != nil {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "purchase Order Detail Is Not Found",
				Err:        errors.New("purchase Order Detail Is Not Found"),
			}
		}
		if PurchaseOrderDetailEntities.BinningQuantity+payloads.DeliveryOrderQuantity > *PurchaseOrderDetailEntities.ItemQuantity {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Item Has Excedeeded Reference Qty",
				Err:        errors.New("Item Has Excedeeded Reference Qty"),
			}
		}
	}
	//insert process
	BinningDetailEntities = transactionsparepartentities.BinningStockDetail{
		BinningSystemNumber:         payloads.BinningSystemNumber,
		ReferenceDetailSystemNumber: payloads.ReferenceDetailSystemNumber,
		OriginalItemId:              payloads.OriginalItemId,
		ItemId:                      payloads.ItemId,
		ItemPrice:                   payloads.ItemPrice,
		UomId:                       payloads.ItemId,
		WarehouseLocationId:         payloads.WarehouseLocationId,
		PurchaseOrderQuantity:       payloads.PurchaseOrderQuantity,
		DeliveryOrderQuantity:       payloads.DeliveryOrderQuantity,
		ReferenceSystemNumber:       payloads.ReferenceSystemNumber,
		ReferenceLineNumber:         payloads.ReferenceLineNumber,
		GoodsReceiveSystemNumber:    payloads.GoodsReceiveSystemNumber,
		GoodsReceiveLineNumber:      payloads.GoodsReceiveLineNumber,
		CreatedByUserId:             payloads.CreatedByUserId,
		CreatedDate:                 payloads.CreatedDate,
		UpdatedByUserId:             payloads.UpdatedByUserId,
		UpdatedDate:                 payloads.UpdatedDate,
		ChangeNo:                    payloads.ChangeNo,
	}
	err = db.Create(&BinningDetailEntities).Scan(&BinningDetailEntities).Error
	if err != nil {
		return BinningDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed To Create Binning Detail Data Please Check Input"),
		}
	}
	if BinningRefTypeCode == "PO" {
		//UPDATE atItemPO1
		//SET BINNING_QTY = ISNULL(BINNING_QTY,0) + @Do_Qty
		//WHERE PO_SYS_NO = @Ref_Sys_No AND PO_LINE = @Ref_Line

		PurchaseOrderDetailEntities.BinningQuantity += payloads.DeliveryOrderQuantity
		err = db.Save(&PurchaseOrderDetailEntities).Error
		if err != nil {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("failed To Update Purchase Order Detail Data Please Check Input"),
			}
		}
	}
	if BinningRefTypeCode == "CL" {
		//item claim entity is not ready yet
		//UPDATE atItemClaim1
		//SET QTY_BINNING = ISNULL(QTY_BINNING,0) + @Do_Qty
		//WHERE CLAIM_SYS_NO = @Ref_Sys_No AND CLAIM_LINE_NO = @Ref_Line

	}
	if BinningRefTypeCode == "WC" {
		//UPDATE wtWorkOrder2
		WorkOrderDetail := transactionworkshopentities.WorkOrderDetail{}
		//SET BINNING_QTY = ISNULL(BINNING_QTY,0) + @Do_Qty
		//WHERE WO_SYS_NO = @Ref_Sys_No AND WO_OPR_ITEM_LINE = @Ref_Line
		err = db.Model(&WorkOrderDetail).
			Where(transactionworkshopentities.WorkOrderDetail{WorkOrderDetailId: payloads.ReferenceDetailSystemNumber}).
			UpdateColumn("binning_quantity", gorm.Expr("binning_quantity + ?", payloads.DeliveryOrderQuantity)).
			Scan(&WorkOrderDetail).Error

		//_LINE FROM wtWorkOrder2
		//WHERE WO_SYS_NO = @Ref_Sys_No
		//AND WO_OPR_ITEM_LINE = @Ref_Line
		//AND BINNING_QTY > FRT_QTY)
		if WorkOrderDetail.BinningQuantity > WorkOrderDetail.FrtQuantity {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("Binning is bigger than Work Order Qty"),
				Message:    "Binning is bigger than Work Order Qty",
			}
		}
	}
	return BinningDetailEntities, nil

}
func (b *BinningListRepositoryImpl) UpdateBinningListDetail(db *gorm.DB, payloads transactionsparepartpayloads.BinningListDetailUpdatePayloads) (transactionsparepartentities.BinningStockDetail, *exceptions.BaseErrorResponse) {
	BinningRefTypeCode := ""
	//LastDoQty := ""
	BinningEntities := transactionsparepartentities.BinningStock{}
	BinningDetailEntities := transactionsparepartentities.BinningStockDetail{}

	err := db.Model(&BinningEntities).
		Where(transactionsparepartentities.BinningStock{BinningSystemNumber: payloads.BinningSystemNumber}).
		Scan(&BinningEntities).Error
	//err := db.Model(&BinningEntities).
	//	Preload("BinningReferenceType"). // Preload the associated BinningReferenceType
	//	Where("binning_system_number = ?", payloads.BinningSystemNumber).
	//	Scan(&BinningEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "Binning Header Is Not Found Please Check Input",
			}
		}
		return BinningDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed To Get Binning error on : " + err.Error()),
		}
	}

	//get detail first
	err = db.Model(&BinningDetailEntities).
		Where(transactionsparepartentities.BinningStockDetail{BinningDetailSystemNumber: payloads.BinningDetailSystemNumber}).
		Scan(&BinningDetailEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
				Message:    "Binning Detail Is Not Found Please Check Input",
			}
		}
		return BinningDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "unexpected Error Occur when querying binning detail",
			//Err:        errors.New("unexpected Error Occur when querying binning detail"),
		}
	}

	LastDoQty := BinningDetailEntities.DeliveryOrderQuantity
	err = db.Model(masterentities.BinningReferenceTypeMaster{}).
		Select("binning_reference_type_code").
		Where(masterentities.BinningReferenceTypeMaster{BinningReferenceTypeId: BinningEntities.BinningReferenceTypeId}).
		Scan(&BinningRefTypeCode).
		Error
	if err != nil {
		return BinningDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed To Get Binning Reference Type Code",
			Err:        errors.New("failed To Get Binning Reference Type Code"),
		}
	}
	BinningDetailEntities.DeliveryOrderQuantity = payloads.DeliveryOrderQuantity
	BinningDetailEntities.ItemId = payloads.ItemId
	BinningDetailEntities.WarehouseLocationId = 0
	BinningDetailEntities.OriginalItemId = payloads.OriginalItemId
	BinningDetailEntities.ChangeNo += 1
	BinningDetailEntities.UpdatedByUserId = payloads.UpdatedByUserId
	BinningDetailEntities.UpdatedDate = payloads.UpdatedDate

	err = db.Save(&BinningDetailEntities).Error
	if err != nil {
		return BinningDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed To Save Binning Stock Detail"),
		}
	}
	if BinningRefTypeCode == "PO" {
		err = db.Model(&transactionsparepartentities.PurchaseOrderDetailEntities{}).
			Where("purchase_order_detail_system_number = ?", BinningDetailEntities.ReferenceSystemNumber).
			Update("binning_quantity", gorm.Expr("COALESCE(binning_quantity, 0) - ? + ?", LastDoQty, payloads.DeliveryOrderQuantity)).Error
		if err != nil {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}
	if BinningRefTypeCode == "CL" {

	}
	if BinningRefTypeCode == "WC" {
		err = db.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Where(transactionworkshopentities.WorkOrderDetail{WorkOrderDetailId: payloads.ReferenceDetailSystemNumber}).
			Update("binning_quantity", gorm.Expr("COALESCE(binning_quantity, 0) - ? + ?", LastDoQty, payloads.DeliveryOrderQuantity)).
			Error
		if err != nil {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Faild Update Work Order Detail Dataa",
				Err:        err,
			}
		}
	}
	return BinningDetailEntities, nil

}
