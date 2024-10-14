package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
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
		//rdb.FlushDB(ctx).Err()
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
