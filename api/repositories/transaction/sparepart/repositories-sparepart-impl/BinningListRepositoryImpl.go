package transactionsparepartrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BinningListRepositoryImpl struct {
}

func (b *BinningListRepositoryImpl) GetReferenceNumberTypoPOWithPagination(db *gorm.DB, filter []utils.FilterCondition, paginations pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//var purchaseOrderEntities []transactionsparepartentities.PurchaseOrderEntities
	var purchaseOrderEntities []transactionsparepartentities.PurchaseOrderEntities
	joinTable := db.Table(`trx_item_purchase_order A`).
		Select(`	A.purchase_order_document_date,
						A.purchase_order_document_number,
						A.supplier_id
`)
	WhereQuery := utils.ApplyFilter(joinTable, filter)
	err := WhereQuery.Scopes(pagination.Paginate(&paginations, WhereQuery)).Order("purchase_order_document_number").Scan(&purchaseOrderEntities).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return paginations, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to to paginate purchase order reference type",
			Err:        errors.New("failed To Get Paginate purchase order reference type"),
		}
	}
	if len(purchaseOrderEntities) == 0 {
		paginations.Rows = []string{}
		return paginations, nil
	}
	var responsePO []transactionsparepartpayloads.BinningListReferenceDocumentNumberTypePOResponse

	for _, item := range purchaseOrderEntities {
		supplierRes, errRes := generalserviceapiutils.GetSupplierMasterByID(item.SupplierId)
		if errRes != nil {
			return paginations, errRes
		}
		responsePO = append(responsePO, transactionsparepartpayloads.BinningListReferenceDocumentNumberTypePOResponse{
			PurchaseNumber: item.PurchaseOrderDocumentNumber,
			DocumentDate:   *item.PurchaseOrderDocumentDate,
			SupplierCode:   supplierRes.SupplierCode,
			SupplierName:   supplierRes.SupplierName,
			SupplierId:     item.SupplierId,
		})
	}
	paginations.Rows = responsePO
	return paginations, nil
}

func NewbinningListRepositoryImpl() transactionsparepartrepository.BinningListRepository {
	return &BinningListRepositoryImpl{}
}
func (b *BinningListRepositoryImpl) GetAllBinningListDetailWithPagination(db *gorm.DB, filter []utils.FilterCondition, paginations pagination.Pagination, binningListId int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
                     	AND A.purchase_order_detail_system_number = F.purchase_order_detail_system_number`).
		Joins("INNER JOIN mtr_warehouse_location WL ON A.warehouse_location_id = WL.warehouse_location_id").
		Where("A.binning_system_number = ?", binningListId)
	WhereQuery := utils.ApplyFilter(joinTable, filter)
	err := WhereQuery.Scopes(pagination.Paginate(&paginations, WhereQuery)).Order("binning_system_number").Scan(&Responses).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return paginations, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed To Get Paginate Binning List"),
		}
	}

	if len(Responses) == 0 {
		paginations.Rows = []string{}
		return paginations, nil
	}
	paginations.Rows = Responses
	return paginations, nil
}

//func GetApprovalStatusId(code string) int {
//	var DocResponse []generalservicepayloads.ApprovalStatusResponses
//
//	DocumentStatusUrl := config.EnvConfigs.GeneralServiceUrl + "approval-status-codes/" + code
//	if err := utils.GetArray(DocumentStatusUrl, &DocResponse, &DocResponse); err != nil {
//		return 0
//	}
//	return DocResponse[0].ApprovalStatusId
//}

func (b *BinningListRepositoryImpl) GetAllBinningListWithPagination(db *gorm.DB, rdb *redis.Client, filter []utils.FilterCondition, paginations pagination.Pagination, ctx context.Context) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var binningEntities []transactionsparepartentities.BinningStock
	var BinningResponses []transactionsparepartpayloads.BinningListGetPaginationResponse
	joinTable := db.Model(&binningEntities)
	WhereQuery := utils.ApplyFilter(joinTable, filter)
	err := WhereQuery.Scopes(pagination.Paginate(&paginations, WhereQuery)).Order("binning_system_number").Scan(&binningEntities).Error
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
		SupplierData, errSupplierMaster := generalserviceapiutils.GetSupplierMasterByID(binningEntity.SupplierId)
		if err != nil {
			return paginations, errSupplierMaster
		}
		exists, _ := rdb.Exists(ctx, strconv.Itoa(binningEntity.BinningDocumentStatusId)).Result()

		if exists == 0 {
			fmt.Println("Failed To Get Status on redis... queue for external service")
			DocResponse, errDocResponse := generalserviceapiutils.GetDocumentStatusById(binningEntity.BinningDocumentStatusId)
			if errDocResponse != nil {
				return paginations, errDocResponse
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
	SupplierData, supplierDataErr := generalserviceapiutils.GetSupplierMasterByID(BinningStockEntities.SupplierId)
	if supplierDataErr != nil {
		return Response, supplierDataErr
	}
	//var SupplierData generalservicepayloads.SupplierMasterCrossServicePayloads
	//SupplierDataUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(BinningStockEntities.SupplierId)
	//if err := utils.Get(SupplierDataUrl, &SupplierData, nil); err != nil {
	//	return Response, &exceptions.BaseErrorResponse{
	//		StatusCode: http.StatusInternalServerError,
	//		Message:    "Failed to fetch Supplier data from external service" + err.Error(),
	//		Err:        err,
	//	}
	//}
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
			Err:        errors.New("error On retreiving Data"),
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
				Err:        errors.New("item Has Excedeeded Reference Qty"),
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
		err = db.Model(&transactionsparepartentities.ItemClaimDetail{}).
			Where(transactionsparepartentities.ItemClaimDetail{ItemClaimDetailId: payloads.ReferenceDetailSystemNumber}).
			Update("binning_quantity", gorm.Expr("COALESCE(binning_quantity, 0) ?", payloads.DeliveryOrderQuantity)).
			Error
		if err != nil {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("failed To Update Claim Data"),
			}
		}

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
		if err != nil {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("work Order Data"),
			}
		}
		//_LINE FROM wtWorkOrder2
		//WHERE WO_SYS_NO = @Ref_Sys_No
		//AND WO_OPR_ITEM_LINE = @Ref_Line
		//AND BINNING_QTY > FRT_QTY)
		if WorkOrderDetail.BinningQuantity > WorkOrderDetail.FrtQuantity {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("binning is bigger than Work Order Qty"),
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
		err = db.Model(&transactionsparepartentities.ItemClaimDetail{}).
			Where(transactionsparepartentities.ItemClaimDetail{ItemClaimDetailId: payloads.ReferenceDetailSystemNumber}).
			Update("binning_quantity", gorm.Expr("COALESCE(binning_quantity, 0) - ? + ?", LastDoQty, payloads.DeliveryOrderQuantity)).
			Error
		if err != nil {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("failed To Update Claim Data"),
			}
		}
	}
	if BinningRefTypeCode == "WC" {
		err = db.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Where(transactionworkshopentities.WorkOrderDetail{WorkOrderDetailId: payloads.ReferenceDetailSystemNumber}).
			Update("binning_quantity", gorm.Expr("COALESCE(binning_quantity, 0) - ? + ?", LastDoQty, payloads.DeliveryOrderQuantity)).
			Error
		if err != nil {
			return BinningDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed Update Work Order Detail Data",
				Err:        err,
			}
		}
	}
	return BinningDetailEntities, nil
}

// atbinningstock0_update
// option 1
func (b *BinningListRepositoryImpl) SubmitBinningList(db *gorm.DB, BinningId int) (transactionsparepartentities.BinningStock, *exceptions.BaseErrorResponse) {
	BinningEntities := transactionsparepartentities.BinningStock{}
	var BinningDetailEntities []transactionsparepartentities.BinningStockDetail
	//get data first
	err := db.Model(&BinningEntities).
		Where("binning_system_number = ?", BinningId).
		Scan(&BinningEntities).
		Error
	if err != nil {
		return transactionsparepartentities.BinningStock{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("binning Not Found Please Check Input"),
		}
	}
	//get detail data
	err = db.Model(&BinningDetailEntities).
		Where(transactionsparepartentities.BinningStockDetail{BinningSystemNumber: BinningId}).
		Scan(&BinningDetailEntities).
		Error
	if err != nil {
		return transactionsparepartentities.BinningStock{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	//if not exists (select * from atBinningStock1 where BIN_SYS_NO = @Bin_Sys_No)
	//	begin
	//	raiserror('please add detail first',16,1)
	//	return 0
	//	end
	if len(BinningDetailEntities) == 0 {
		return transactionsparepartentities.BinningStock{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("please add detail first"),
		}
	}
	//get validation draft doc
	DocResponse, errDocResponse := generalserviceapiutils.GetDocumentStatusById(BinningEntities.BinningDocumentStatusId)
	if errDocResponse != nil {
		return transactionsparepartentities.BinningStock{}, errDocResponse
	}

	if DocResponse.DocumentStatusCode != "10" { //draft {
		return transactionsparepartentities.BinningStock{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("document Is Not Draft"),
		}
	}
	//https://testing-backendims.indomobil.co.id/general-service/v1/source-document-type-code/SPBN
	//IF @Src_Code IS NULL
	//BEGIN
	//RAISERROR('Document Source Is Not Define at Table comGenVariable',16,1)
	//RETURN 0
	//END
	SourceDocType, errSource := generalserviceapiutils.GetDocumentTypeByCode("SPBN")
	if errSource != nil {
		return BinningEntities, errSource
	}

	if SourceDocType == (generalserviceapiutils.SourceDocumentTypeMasterResponse{}) {
		return BinningEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Document Source Is Not Define at Table",
		}
	}
	var count int64
	err = db.Model(&transactionsparepartentities.BinningStockDetail{}).
		Where("binning_system_number = ?", BinningId).
		Where("COALESCE(delivery_order_quantity, 0) > COALESCE(purchase_order_quantity, 0)"). // ISNULL in SQL is equivalent to COALESCE in GORM
		Count(&count).Error
	if err != nil {
		return BinningEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to check binning stock condition"),
		}
	}
	// If any record matches the condition, raise the error

	//del comment
	if count > 0 {
		return BinningEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("binning Qty is bigger than PO Qty"),
		}
	}
	//get company code
	CompanyData, errResponse := generalserviceapiutils.GetCompanyDataById(BinningEntities.CompanyId)
	if errResponse != nil {
		return BinningEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed To Fetch Company Data From Cross service"),
		}
	}

	//if @Company_Code = '3125098'
	//and (SELECT ISNULL(WHS_GROUP,'') FROM atBinningStock0 WHERE BIN_SYS_NO = @Bin_Sys_No) = 'PG'
	//and (SELECT ISNULL(WHS_CODE,'') FROM atBinningStock0 WHERE BIN_SYS_NO = @Bin_Sys_No) = 'WPG01'
	//and exists (select 1
	//	--A.item_code,B.loc_Code
	//	from atbinningstock1 A
	//	left join amlocationitem B on B.company_code = '3125098'
	//	and B.whs_code = 'WPG01'
	//	and B.whs_group = 'PG'
	//	and A.item_code = B.item_code
	//	where bin_sys_no = @Bin_Sys_No
	//	and ISNULL(B.loc_code ,'') = '')
	//
	//	begin
	//	raiserror('There is Item Code without Master > Item Location. Please register first!',16,1)
	//	return 0
	//	end
	//select the sub query first
	//sementara hard code karna belum ada tabel comgen dan di SP masi hardcode
	var defaultWarehouseGroup = "PG"
	var defaultWarehouse = ""
	if CompanyData.CompanyCode == "1516098" { //KIA
		defaultWarehouse = "WKG01"
	} else if CompanyData.CompanyCode == "1518098" { //CITROEN
		defaultWarehouse = "WCG01"
	} else if CompanyData.CompanyCode == "140000" { //yadea
		defaultWarehouse = "WYG01"
	} else if CompanyData.CompanyCode == "3125098" { //nmdi
		defaultWarehouse = "WPG01"
	} else {
		defaultWarehouse = "WPG01"
	}
	defaultWarehouseGroupId := 0
	defaultWarehouseId := 0

	//get warehouse group id for validation
	err = db.Model(&masterwarehouseentities.WarehouseGroup{}).
		Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupCode: defaultWarehouseGroup}).
		Select("warehouse_group_id").Scan(&defaultWarehouseGroupId).Error
	if err != nil {
		return BinningEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get warehouse group",
		}
	}
	//get warehouse group id
	err = db.Model(masterwarehouseentities.WarehouseMaster{}).
		Where(masterwarehouseentities.WarehouseMaster{WarehouseCode: defaultWarehouse}).
		Select("warehouse_id").Scan(&defaultWarehouseId).Error
	if err != nil {
		return BinningEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get warehouse",
		}
	}
	var result int
	//validation for nmdi
	if CompanyData.CompanyCode == "3125098" &&
		BinningEntities.WarehouseGroupId == defaultWarehouseGroupId &&
		BinningEntities.WarehouseId == defaultWarehouseId {

		var IsLocExist bool = false
		err = db.Table("trx_binning_list_stock_detail A").
			Joins(`LEFT JOIN mtr_location_item B ON
						AND B.warehouse_id = ?
						AND B.warehouse_group_id = ?
						AND A.item_id = B.item_id
					`, defaultWarehouseId, defaultWarehouseGroupId).
			Select("1").
			Where("A.binning_system_number = ?", BinningId).
			Scan(&IsLocExist).
			Error
		if err != nil {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to get warehouse location item",
			}
		}
		if !IsLocExist {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("there is Item Code without Master > Item Location. Please register first"),
			}
		}
		//AND EXISTS (SELECT TOP 1 1 FROM atBinningStock1 where BIN_SYS_NO = @Bin_Sys_No GROUP BY ITEM_CODE HAVING COUNT(*) > 1)
		err = db.Model(&BinningDetailEntities).
			Select("1").Where("binning_system_number = ?", BinningId).
			Group("item_id").
			Having("count (*) > 1").
			Limit(1).
			Scan(&result).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed To Check Double Item",
			}
		}

		//IF @Company_Code = '1516098'
		//AND (SELECT ISNULL(WHS_GROUP,'') FROM atBinningStock0 WHERE BIN_SYS_NO = @Bin_Sys_No) = 'PG'
		//AND (SELECT ISNULL(WHS_CODE,'') FROM atBinningStock0 WHERE BIN_SYS_NO = @Bin_Sys_No) = 'WKG01'
		//AND EXISTS (SELECT TOP 1 1 FROM atBinningStock1 where BIN_SYS_NO = @Bin_Sys_No GROUP BY ITEM_CODE HAVING COUNT(*) > 1)
		//begin
		//SET @doubleItem = (SELECT TOP 1 ITEM_CODE FROM atBinningStock1 where BIN_SYS_NO = @Bin_Sys_No GROUP BY ITEM_CODE HAVING COUNT(*) > 1)
		//raiserror('Item %s multiple line! Please split Binning document!',16,1,@doubleItem)
		//return 0
		//end
		if result == 1 {
			var doubleItem string
			err = db.Table("trx_binning_list_stock_detail A").
				Joins("mtr_item B ON A.item_id = B.item_id").
				Select("B.item_code").
				Where("binning_system_number = ?", BinningId).
				Group("A.item_id").
				Having("COUNT(A.item_id) > ?", 1).
				Limit(1).
				Scan(&doubleItem).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return BinningEntities, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to Check Double Item",
				}
			}
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("item %s multiple line! Please split Binning document", doubleItem),
			}
		}
		//IF @Company_Code = '3125098'
		//AND (SELECT ISNULL(WHS_GROUP,'') FROM atBinningStock0 WHERE BIN_SYS_NO = @Bin_Sys_No) = 'PG'
		//AND (SELECT ISNULL(WHS_CODE,'') FROM atBinningStock0 WHERE BIN_SYS_NO = @Bin_Sys_No) = 'WPG01'
		//AND ISNULL((SELECT SUPPLIER_INV_NO FROM atBinningStock0 WHERE BIN_SYS_NO = @Bin_Sys_No),'') = ''
		//begin
		//raiserror('Invoice No is empty! Please type Invoice No then hit the Save button!',16,1)
		//retur
		//SupplierInvNo := ""
		//err = db.Model("")"
		if BinningEntities.SupplierInvoiceNumber == "" {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("invoice No is empty! Please type Invoice No then hit the Save button"),
			}
		}
		//AND (SELECT COUNT(DISTINCT REF_SYS_NO) FROM atBinningStock1 where BIN_SYS_NO = @Bin_Sys_No) > 1
		//begin
		//raiserror('Binning Reference Detail More Than 1 Purchase Order! Please contact your system administrator.',16,1)
		//return 0
		//end
		result := 0
		err = db.Model(&BinningDetailEntities).
			Where("binning_system_number = ?", BinningId).
			Select("COUNT(distinct reference_system_number)").
			Scan(&result).
			Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed To Check Reference Number",
			}
		}
		if result > 1 {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("binning Reference Detail More Than 1 Purchase Order! Please contact your system administrator"),
			}
		}
		//IF @Company_Code = '3125098'
		//AND (SELECT ISNULL(WHS_GROUP,'') FROM atBinningStock0 WHERE BIN_SYS_NO = @Bin_Sys_No) = 'PG'
		//AND (SELECT ISNULL(WHS_CODE,'') FROM atBinningStock0 WHERE BIN_SYS_NO = @Bin_Sys_No) = 'WPG01'
		//AND ISNULL((SELECT TRX_TYPE FROM atBinningStock0 WHERE BIN_SYS_NO = @Bin_Sys_No),'') = 'I'
		//AND NOT EXISTS (SELECT TOP 1 1 FROM atSupplierInvoice0 WHERE COMPANY_CODE = @Company_Code AND SUPP_INV_STATUS = '20' AND SUPP_INV_DOC_NO = (SELECT SUPPLIER_INV_NO FROM atBinningStock0 WHERE BIN_SYS_NO = @Bin_Sys_No))
		//begin
		//raiserror('Supplier Invoice status not Approved!',16,1)
		//return 0
	}

	if CompanyData.CompanyCode == "1516098" && BinningEntities.WarehouseGroupId == defaultWarehouseGroupId &&
		BinningEntities.WarehouseId == defaultWarehouseId {
		result = 0
		err = db.Model(&BinningDetailEntities).
			Select("1").Where("binning_system_number = ?", BinningId).
			Group("item_id").
			Having("count (*) > 1").
			Limit(1).
			Scan(&result).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed To Check Double Item",
			}
		}
		if result == 1 {
			var doubleItem string
			err = db.Table("trx_binning_list_stock_detail A").
				Joins("mtr_item B ON A.item_id = B.item_id").
				Select("B.item_code").
				Where("binning_system_number = ?", BinningId).
				Group("A.item_id").
				Having("COUNT(A.item_id) > ?", 1).
				Limit(1).
				Scan(&doubleItem).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return BinningEntities, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to Check Double Item",
				}
			}
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("item %s multiple line! Please split Binning document", doubleItem),
			}
		}
		result = 0
		err = db.Model(&BinningDetailEntities).
			Where("binning_system_number = ?", BinningId).
			Select("COUNT(distinct reference_system_number)").
			Scan(&result).
			Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed To Check Reference Number",
			}
		}
		if result > 1 {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("binning Reference Detail More Than 1 Purchase Order! Please contact your system administrator"),
			}
		}
	}
	if BinningEntities.WarehouseGroupId == 0 || BinningEntities.WarehouseId == defaultWarehouseId {
		return BinningEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("warehouse Group / warehouse Code empty! Please check your Binning Document"),
		}
	}
	var BinningRefTypeCode string
	err = db.Model(masterentities.BinningReferenceTypeMaster{}).
		Select("binning_reference_type_code").
		Where(masterentities.BinningReferenceTypeMaster{BinningReferenceTypeId: BinningEntities.BinningReferenceTypeId}).
		Scan(&BinningRefTypeCode).
		Error
	if err != nil {
		return BinningEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get bining reference type code",
		}
	}
	var vehicleBrandId int
	if BinningRefTypeCode == "PO" {
		err = db.Model(&transactionsparepartentities.PurchaseOrderEntities{}).
			Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: BinningEntities.ReferenceSystemNumber}).
			Select("brand_id").Scan(&vehicleBrandId).Error
		if err != nil {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get Brand Id From Purchase Order",
			}
		}
	}
	if BinningRefTypeCode == "CL" {
		err = db.Model(&transactionsparepartentities.ItemClaim{}).
			Where(transactionsparepartentities.ItemClaim{ClaimSystemNumber: BinningEntities.ReferenceSystemNumber}).
			Select("vehicle_brand_id").Scan(&vehicleBrandId).Error
		if err != nil {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get Brand Id From item claim",
			}
		}
	}
	if BinningRefTypeCode == "WC" {
		err = db.Model(&transactionworkshopentities.WorkOrder{}).
			Where(transactionworkshopentities.WorkOrder{WorkOrderSystemNumber: BinningEntities.ReferenceSystemNumber}).
			Select("brand_id").
			Scan(&vehicleBrandId).Error
		if err != nil {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get Brand Id Work Order",
			}
		}
	}
	//endpoint not ready
	//EXEC uspg_gmSrcDoc1_Update
	//@Option = 0 ,
	//@COMPANY_CODE = @Company_Code ,
	//@SOURCE_CODE = @Src_Code ,
	//@VEHICLE_BRAND = @VEHICLE_BRAND, --'' ,
	//@PROFIT_CENTER_CODE = '' ,
	//@TRANSACTION_CODE = '' ,
	//@BANK_ACC_CODE = '' ,
	//@TRANSACTION_DATE = @Bin_Doc_Date ,
	//@Last_Doc_No = @Bin_Doc_No OUTPUT
	//dummy for gmsercdoc 1 update

	//UPDATE atBinningStock0 SET
	//SUPPLIER_DLVR_PERSON = @Supplier_Dlvr_Person,
	//	CHANGE_USER_ID = @Change_User_Id ,
	//	CHANGE_DATETIME = @Change_Datetime
	//WHERE BIN_SYS_NO = @Bin_Sys_No AND TRX_TYPE = 'I'

	//UPDATE atBinningStock0 SET
	//WHS_GROUP = 'PG',
	//	WHS_CODE = 'WKG01'
	//WHERE BIN_SYS_NO = @Bin_Sys_No AND TRX_TYPE = 'I' AND COMPANY_CODE = 1516098 AND ISNULL(WHS_GROUP,'') = '' AND ISNULL(WHS_CODE,'') = ''

	//BinningEntities
	//get ready id
	//DocumentStatusUrl = config.EnvConfigs.GeneralServiceUrl + "document-status-by-code/20"
	//if err := utils.Get(DocumentStatusUrl, &DocResponse, nil); err != nil {
	//	return BinningEntities, &exceptions.BaseErrorResponse{
	//		StatusCode: http.StatusBadRequest,
	//		Err:        errors.New("failed To fetch Document Status From External Service"),
	//	}
	//}
	DocumentStatusUrl, errStatus := generalserviceapiutils.GetDocumentStatusByCode("20")
	if errStatus != nil {
		return BinningEntities, errStatus
	}
	//get bining reference type id first "T"
	binningTypeIdImport := 0
	err = db.Model(&masterentities.BinningTypeMaster{}).
		Where(masterentities.BinningTypeMaster{BinningTypeCode: "I"}).
		Select("binning_type_id").Scan(&binningTypeIdImport).Error
	if err != nil {
		return BinningEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get Binning Type Id",
		}
	}
	//UPDATE atBinningStock0 SET
	//BIN_STATUS = dbo.GetVariableValue('DOC_STATUS_READY') ,
	//	BIN_DOC_NO = @Bin_Doc_No ,
	//	CHANGE_NO = CHANGE_NO + 1 ,
	//	CHANGE_USER_ID = @Change_User_Id ,
	//	CHANGE_DATETIME = @Change_Datetime
	//WHERE BIN_SYS_NO = @Bin_Sys_No
	BinningEntities.BinningDocumentStatusId = DocumentStatusUrl.DocumentStatusId
	BinningEntities.ChangeNo += 1
	BinningEntities.BinningDocumentNumber = "dummy doc number waiting for gmsrcdoc1_update"
	*BinningEntities.CreatedDate = time.Now()

	err = db.Save(&BinningEntities).Error
	if err != nil {
		return BinningEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save Binning Document",
		}
	}
	if (CompanyData.CompanyCode == "1516098" || //kia
		CompanyData.CompanyCode == "1518098" || //citroen
		CompanyData.CompanyCode == "140000") && //yadea
		BinningEntities.WarehouseId == 0 &&
		BinningEntities.WarehouseGroupId == 0 &&
		BinningEntities.BinningTypeId == binningTypeIdImport {

		BinningEntities.WarehouseId = defaultWarehouseId
		BinningEntities.WarehouseGroupId = defaultWarehouseGroupId
		err = db.Save(&BinningEntities).Error
		if err != nil {
			return BinningEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to save Binning Document",
			}
		}
	}
	return BinningEntities, nil
}

func (b *BinningListRepositoryImpl) DeleteBinningList(db *gorm.DB, BinningId int) (bool, *exceptions.BaseErrorResponse) {
	//get entity first
	var binningListEntity transactionsparepartentities.BinningStock
	err := db.Model(&binningListEntity).Where(transactionsparepartentities.BinningStock{BinningSystemNumber: BinningId}).
		Preload("BinningStockDetail").
		Find(&binningListEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("binning stock is not found please check input"),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get binning stock header to delete please check input",
		}
	}
	//get status draft on document master
	ApprovalStatusResponseDraft, errs := generalserviceapiutils.GetApprovalStatusByCode("10")
	if errs != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get approval status",
		}
	}
	if binningListEntity.BinningDocumentStatusId != ApprovalStatusResponseDraft.ApprovalStatusId {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("document is not draft"),
		}
	}
	//get reference binning list to check
	var referenceTypeBinningList masterentities.BinningReferenceTypeMaster

	err = db.Model(&referenceTypeBinningList).Where(masterentities.BinningReferenceTypeMaster{BinningReferenceTypeId: binningListEntity.BinningReferenceTypeId}).
		First(&referenceTypeBinningList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("failed to get reference type please check reference type input"),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error on reference type binning list please check input",
		}
	}
	//loop through all the binning detail
	var binningStockDetailEntities []transactionsparepartentities.BinningStockDetail

	err = db.Model(&binningStockDetailEntities).Where(transactionsparepartentities.BinningStockDetail{BinningSystemNumber: BinningId}).
		Scan(&binningStockDetailEntities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get binning stock detail please check input",
		}
	}
	for _, binningDetail := range binningStockDetailEntities {
		if strings.ToUpper(referenceTypeBinningList.BinningReferenceTypeCode) == "PO" || strings.ToUpper(referenceTypeBinningList.BinningReferenceTypeCode) == "IV" {
			//cara 1
			//db.Model(&transactionsparepartentities.PurchaseOrderDetailEntities{}).
			//	Updates(map[string]interface{}{
			//		"binning_quantity": gorm.Expr(`
			//					CASE WHEN (ISNULL(binning_quantity,0) - ? < 0
			//					THEN 0
			//					ELSE
			//					(ISNULL(binning_quantity,0) - ?
			//					END
			//				`, binningDetail.DeliveryOrderQuantity, binningDetail.DeliveryOrderQuantity),
			//	}).Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderDetailSystemNumber: binningDetail.ReferenceDetailSystemNumber})

			var purchaseOrderDetail transactionsparepartentities.PurchaseOrderDetailEntities
			err = db.Model(&purchaseOrderDetail).Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderDetailSystemNumber: binningDetail.ReferenceDetailSystemNumber}).
				First(&purchaseOrderDetail).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Err:        errors.New("purchase order is not found with reference type purchase order please check input"),
					}
				}
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to update purchase order please check input",
				}
			}
			purchaseOrderDetail.BinningQuantity -= binningDetail.DeliveryOrderQuantity
			if purchaseOrderDetail.BinningQuantity <= 0 {
				purchaseOrderDetail.BinningQuantity = 0
			}
			err = db.Save(&purchaseOrderDetail).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to update purchase order please check input",
				}
			}
		} else if referenceTypeBinningList.BinningReferenceTypeCode == "CL" {
			var claimDetailEntity transactionsparepartentities.ItemClaimDetail
			err = db.Model(&claimDetailEntity).Where(transactionsparepartentities.ItemClaimDetail{ItemClaimDetailId: binningDetail.ReferenceDetailSystemNumber}).
				First(&claimDetailEntity).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Err:        errors.New("item claim detail is not found with reference type item claim please check input"),
					}
				}
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to update item claim detail please check input",
				}
			}
			claimDetailEntity.QuantityBinning -= binningDetail.DeliveryOrderQuantity
			if claimDetailEntity.QuantityBinning < 0 {
				claimDetailEntity.QuantityBinning = 0
			}
			err = db.Save(&claimDetailEntity).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to update item claim detail please check input",
				}
			}
		} else if strings.ToUpper(referenceTypeBinningList.BinningReferenceTypeCode) == "WC" {
			//get work order detail first
			var workOrderDetail transactionworkshopentities.WorkOrderDetail
			err = db.Model(&workOrderDetail).Where(transactionworkshopentities.WorkOrderDetail{WorkOrderDetailId: binningDetail.ReferenceDetailSystemNumber}).
				First(&workOrderDetail).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Err:        errors.New("work order is not found with reference type workorder please check input"),
					}
				}
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to update work order please check input",
				}
			}
			workOrderDetail.BinningQuantity -= binningDetail.DeliveryOrderQuantity
			if workOrderDetail.BinningQuantity < 0 {
				workOrderDetail.BinningQuantity = 0
			}
			err = db.Save(&workOrderDetail).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to update work order please check input",
				}
			}
		}
		//delete detail
		err = db.Delete(&binningDetail).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to delete binning detail please check input",
			}
		}
	}
	err = db.Delete(&binningListEntity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to delete binning list please check input",
		}
	}
	return true, nil
}
func (b *BinningListRepositoryImpl) DeleteBinningListDetailMultiId(db *gorm.DB, binningDetailMultiId string) (bool, *exceptions.BaseErrorResponse) {
	//var splitId []int
	splitId := strings.Split(binningDetailMultiId, ",")
	for _, item := range splitId {
		itemId, splitErr := strconv.Atoi(item)
		if splitErr != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("id is not a number please check input"),
			}
		}
		//get entity to delete
		var binningListDetailEntity transactionsparepartentities.BinningStockDetail
		err := db.Model(&binningListDetailEntity).
			Where(transactionsparepartentities.BinningStockDetail{BinningDetailSystemNumber: itemId}).
			First(&binningListDetailEntity).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        errors.New("binning list detail is not found with id please check input"),
				}
			}
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to delete binning detail please check input",
			}
		}
		//get header
		var binningListEntity transactionsparepartentities.BinningStock
		err = db.Model(&binningListEntity).
			Where(transactionsparepartentities.BinningStock{BinningSystemNumber: binningListDetailEntity.BinningSystemNumber}).
			First(&binningListEntity).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        errors.New("binning list header is not found with id please check input"),
				}
			}
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to delete binning detail please check input",
			}
		}
		docResponseDraft, DocErr := generalserviceapiutils.GetApprovalStatusByCode("10")
		if DocErr != nil {
			return false, DocErr
		}
		if binningListEntity.BinningDocumentStatusId != docResponseDraft.ApprovalStatusId {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("binning list status is not draft"),
			}
		}
		var referenceTypeBinningList masterentities.BinningReferenceTypeMaster
		err = db.Model(&referenceTypeBinningList).Where(masterentities.BinningReferenceTypeMaster{BinningReferenceTypeId: binningListEntity.BinningReferenceTypeId}).
			First(&referenceTypeBinningList).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        errors.New("failed to get reference type please check reference type input"),
				}
			}
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error on reference type binning list please check input",
			}
		}

		//get reference type binning first
		if strings.ToUpper(referenceTypeBinningList.BinningReferenceTypeCode) == "PO" {
			err = db.Model(&transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderDetailSystemNumber: binningListDetailEntity.ReferenceDetailSystemNumber}).
				Update("binning_quantity", gorm.Expr("ISNULL(binning_quantity,0) - ?", binningListDetailEntity.DeliveryOrderQuantity)).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to update purchase order please check input",
				}
			}
		}
		//type IV
		if strings.ToUpper(referenceTypeBinningList.BinningReferenceTypeCode) == "IV" {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("cannot delete ths line!. Please Void document"),
			}
		}
		if strings.ToUpper(referenceTypeBinningList.BinningReferenceTypeCode) == "CL" {
			err = db.Model(&transactionsparepartentities.ItemClaimDetail{}).
				Where(transactionsparepartentities.ItemClaimDetail{ItemClaimDetailId: binningListDetailEntity.ReferenceDetailSystemNumber}).
				Update("quantity_binning", gorm.Expr("ISNULL(quantity_binning,0) - ?", binningListDetailEntity.DeliveryOrderQuantity)).
				Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to update claim detail please check input",
				}
			}
		}
		if strings.ToUpper(referenceTypeBinningList.BinningReferenceTypeCode) == "WC" {
			err = db.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Where(transactionworkshopentities.WorkOrderDetail{WorkOrderDetailId: binningListDetailEntity.ReferenceDetailSystemNumber}).
				Update("binning_quantity", gorm.Expr("ISNULL(binning_quantity,0) - ?", binningListDetailEntity.DeliveryOrderQuantity)).
				Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to update work order detail please check input",
				}
			}
		}
		err = db.Delete(&binningListDetailEntity).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to delete binning detail please check input",
			}
		}
	}
	return true, nil
}
