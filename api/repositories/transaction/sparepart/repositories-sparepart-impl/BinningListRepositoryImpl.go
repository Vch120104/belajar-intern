package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	generalservicepayloads "after-sales/api/payloads/crossservice/generalservice"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type BinningListRepositoryImpl struct{}

func NewbinningListRepositoryImpl() transactionsparepartrepository.BinningListRepository {
	return &BinningListRepositoryImpl{}
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
			Message:    "Failed to fetch Supplier data from external service",
			Err:        err,
		}
	}
	//get PO Data

	//Item Group Name
	ItemGroupName := ""
	err = db.Table("trx_item_purchase_order A").Select("B.item_group_name").
		Joins("JOIN mtr_item_group B ON A.item_group_id = B.item_group_id").Scan(&ItemGroupName).Error
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
