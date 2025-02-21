package transactionsparepartrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouserepositoryimpl "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func NewItemWarehouseTransferOutRepositoryImpl() transactionsparepartrepository.ItemWarehouseTransferOutRepository {
	return &ItemWarehouseTransferOutRepositoryImpl{}
}

type ItemWarehouseTransferOutRepositoryImpl struct {
}

// UpdateTransferOutDetail implements transactionsparepartrepository.ItemWarehouseTransferOutRepository.
func (*ItemWarehouseTransferOutRepositoryImpl) UpdateTransferOutDetail(tx *gorm.DB, request transactionsparepartpayloads.UpdateItemWarehouseTransferOutDetailRequest, number int) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse) {
	var entitiesDetail transactionsparepartentities.ItemWarehouseTransferOutDetail
	errDet := tx.Model(&entitiesDetail).Find(&entitiesDetail, number)
	if errDet.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errDet.Error,
		}
	}

	if request.LocationId != 0 {
		entitiesDetail.LocationIdFrom = &request.LocationId
		entitiesDetail.LocationIdTo = &request.LocationToId
		entitiesDetail.QuantityOut = request.QuatityOut
	}

	errSave := tx.Save(&entitiesDetail).Error
	if errSave != nil {
		return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errSave,
		}
	}

	return entitiesDetail, nil
}

// DeleteTransferOut implements transactionsparepartrepository.ItemWarehouseTransferOutRepository.
func (*ItemWarehouseTransferOutRepositoryImpl) DeleteTransferOut(tx *gorm.DB, number int) (bool, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferOut
	var entitiesDetail transactionsparepartentities.ItemWarehouseTransferOutDetail

	errGet := tx.Model(&entities).
		Where("transfer_out_system_number = ?", number).First(&entities).Error
	if errGet != nil {
		if errors.Is(errGet, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGet,
				Message:    "transfer out with that id is not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGet,
			Message:    "failed to get transfer out entity",
		}
	}

	err := tx.Model(&entitiesDetail).Where("transfer_out_system_number = ?", number).Delete(&entitiesDetail).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete the transfer out detai",
			Err:        err,
		}
	}

	errDel := tx.Delete(&entities, number).Error
	if errDel != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errDel,
		}
	}

	return true, nil
}

// DeleteTransferOutDetail implements transactionsparepartrepository.ItemWarehouseTransferOutRepository.
func (*ItemWarehouseTransferOutRepositoryImpl) DeleteTransferOutDetail(tx *gorm.DB, number int) (bool, *exceptions.BaseErrorResponse) {
	var entitiesDetail transactionsparepartentities.ItemWarehouseTransferOutDetail

	errGet := tx.Model(&entitiesDetail).
		Where("transfer_out_detail_system_number = ?", number).First(&entitiesDetail).Error
	if errGet != nil {
		if errors.Is(errGet, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGet,
				Message:    "transfer out detail with that id is not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGet,
			Message:    "failed to get transfer out detail entity",
		}
	}

	if err := tx.Delete(&entitiesDetail).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete the transfer out detai",
			Err:        err,
		}
	}

	return true, nil
}

// SubmitTransferOut implements transactionsparepartrepository.ItemWarehouseTransferOutRepository.
func (*ItemWarehouseTransferOutRepositoryImpl) SubmitTransferOut(tx *gorm.DB, number int) (transactionsparepartentities.ItemWarehouseTransferOut, *exceptions.BaseErrorResponse) {
	var entitiesRequest transactionsparepartentities.ItemWarehouseTransferRequest
	var entities transactionsparepartentities.ItemWarehouseTransferOut
	var entitiesDetailOut []transactionsparepartentities.ItemWarehouseTransferOutDetail
	var status masteritementities.ItemTransferStatus

	transType, errTrans := masterrepositoryimpl.NewStockTransactionRepositoryImpl().GetStockTransactionTypeByCode(tx, "TO")
	if errTrans != nil {
		return transactionsparepartentities.ItemWarehouseTransferOut{}, errTrans
	}

	transReason, errReason := masterrepositoryimpl.StartStockTraansactionReasonRepositoryImpl().GetStockTransactionReasonByCode(tx, "NL")
	if errReason != nil {
		return transactionsparepartentities.ItemWarehouseTransferOut{}, errReason
	}

	errGetStatus := tx.Model(&status).Where("item_transfer_status_code = ?", 40).Find(&status)
	if errGetStatus.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatus.Error,
		}
	}

	errGetEntities := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferOut{TransferOutSystemNumber: number}).First(&entities).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntities,
				Message:    "transfer out with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get transfer out entity",
		}
	}

	errGetEntitiesDetail := tx.Model(&entitiesDetailOut).Where(transactionsparepartentities.ItemWarehouseTransferOutDetail{TransferOutSystemNumber: number}).Find(&entitiesDetailOut).Error
	if errGetEntitiesDetail != nil {
		if errors.Is(errGetEntitiesDetail, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntitiesDetail,
				Message:    "transfer out detail with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntitiesDetail,
			Message:    "failed to get transfer out detail entity",
		}
	}

	errGetEntitiesRequest := tx.Model(&entitiesRequest).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestSystemNumber: entities.TransferRequestSystemNumbers}).First(&entitiesRequest).Error
	if errGetEntitiesRequest != nil {
		if errors.Is(errGetEntitiesRequest, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntities,
				Message:    "transfer request with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntitiesRequest,
			Message:    "failed to get transfer request entity",
		}
	}

	// openPeriodResponse, openPeriodErr := financeserviceapiutils.GetOpenPeriodByCompany(entities.CompanyId, "SP")
	// if openPeriodErr != nil {
	// 	return transactionsparepartentities.ItemWarehouseTransferOut{}, openPeriodErr
	// }

	periodYear := entities.TransferOutDate.Format("2006")
	periodMonth := entities.TransferOutDate.Format("01")

	// if openPeriodResponse.PeriodYear != periodYear || openPeriodResponse.PeriodMonth != periodMonth {
	// 	return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    "Period is closed",
	// 		Err:        errors.New("period is closed"),
	// 	}
	// }

	for _, detail := range entitiesDetailOut {
		if *detail.LocationIdFrom == 0 || detail.LocationIdFrom == nil {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("location from is null"),
				Message:    "location origin are empty!",
			}
		} else if *detail.LocationIdTo == 0 || detail.LocationIdTo == nil {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("location to is null"),
				Message:    "location destination are empty!",
			}
		}
	}

	// IF @Src_Code IS NULL
	// 	BEGIN
	// 		RAISERROR('Please setting Transfer Out Document No at Document Master',16,1)
	// 	END
	// 	EXEC dbo.getCurrentPeriod 0, @Company_Code, @Module_Code, @Period_Year OUTPUT ,@Period_Month OUTPUT

	// 	EXEC uspg_gmSrcDoc1_Update
	// 		@Option = 0 ,
	// 		@COMPANY_CODE = @Company_Code ,
	// 		@TRANSACTION_DATE = @Trfout_Date ,
	// 		@SOURCE_CODE = @Src_Code ,
	// 		@VEHICLE_BRAND = @Whs_Brand ,
	// 		@PROFIT_CENTER_CODE = '' ,
	// 		@TRANSACTION_CODE = '' ,
	// 		@BANK_ACC_CODE =  '' ,
	// 		@Change_User_Id = @Change_User_Id,
	// 		@Last_Doc_No = @Trfout_Doc_No OUTPUT

	entities.TransferOutStatusId = status.ItemTransferStatusId
	entities.TransferOutDate = time.Now()

	errUpdateTransfer := tx.Save(&entities).Error
	if errUpdateTransfer != nil {
		return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUpdateTransfer,
		}
	}

	entitiesRequest.TransferRequestStatusId = status.ItemTransferStatusId
	errUpdateRequest := tx.Save(&entitiesRequest).Error
	if errUpdateRequest != nil {
		return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUpdateRequest,
		}
	}

	getWarehouse, errWarehouse := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl().GetById(tx, *entities.WarehouseId)
	if errWarehouse != nil {
		return transactionsparepartentities.ItemWarehouseTransferOut{}, errWarehouse
	}

	for _, detail := range entitiesDetailOut {
		var entitiesDetailRequest transactionsparepartentities.ItemWarehouseTransferRequestDetail
		updates := map[string]interface{}{
			"location_id_from": detail.LocationIdFrom,
		}
		errUpdate := tx.Model(&transactionsparepartentities.ItemWarehouseTransferRequestDetail{}).
			Where("transfer_request_detail_system_number = ?", detail.TransferRequestDetailSystemNumber).
			Updates(updates).Scan(&entitiesDetailRequest).Error

		if errUpdate != nil {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUpdate,
			}
		}

		// -- Begin Insert amStockTransaction
		// 	SET @Trans_Type = dbo.getVariableValue('STK_TRXTYPE_TRANSFER_OUT')
		// 	SET @Trans_Reason_Code = dbo.getVariableValue('STK_REASON_NORMAL')

		var hpp float64
		var gse masterentities.GroupStock
		errGse := tx.Model(&gse).
			Where("company_id = ?", entities.CompanyId).
			Where("warehouse_group_id = ?", getWarehouse.WarehouseGroupId).
			Where("item_id = ?", detail.ItemId).
			Where("period_year = ?", periodYear).
			Where("period_month = ?", periodMonth).
			Find(&gse)

		if errGse.Error != nil {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errGse.Error,
			}
		}

		if gse.GroupStockId != 0 {
			hpp = gse.PriceCurrent
		} else {
			hpp = 0
		}

		detail.CostOfGoodsSold = hpp
		detail.TotalTransferCost = hpp * detail.QuantityOut

		errSaveDetail := tx.Save(&detail).Error
		if errSaveDetail != nil {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errSaveDetail,
			}
		}

		get, errQuant := masterwarehouserepository.NewLocationStockRepositoryImpl().GetAvailableQuantity(tx, masterwarehousepayloads.GetAvailableQuantityPayload{
			CompanyId:        entities.CompanyId,
			ItemId:           *detail.ItemId,
			WarehouseId:      getWarehouse.WarehouseId,
			WarehouseGroupId: getWarehouse.WarehouseGroupId,
			LocationId:       *detail.LocationIdFrom,
			PeriodDate:       entities.TransferOutDate,
		})

		if errQuant != nil {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, errQuant
		}

		if get.QuantityAvailable-detail.QuantityOut < 0 {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("quantity is not available"),
			}
		}

		_, err := StartStockTransactionRepositoryImpl().StockTransactionInsert(tx, transactionsparepartpayloads.StockTransactionInsertPayloads{
			CompanyId:                 entities.CompanyId,
			ReferenceId:               entities.TransferOutSystemNumber,
			TransactionTypeId:         transType.StockTransactionTypeId,
			TransactionReasonId:       transReason.StockTransactionReasonId,
			ReferenceDate:             &entities.TransferOutDate,
			ReferenceQuantity:         detail.QuantityOut,
			TransactionCogs:           hpp,
			ReferenceItemId:           *detail.ItemId,
			ReferenceLocationId:       *detail.LocationIdFrom,
			ReferenceWarehouseId:      *entities.WarehouseId,
			ReferenceWarehouseGroupId: getWarehouse.WarehouseGroupId,
			ReferencePrice:            0,
			ReferenceCurrencyId:       0,
		})

		if err != nil {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, err
		}

		// 	EXEC [dbo].[uspg_atStockTransaction_Insert]
		// 	@Option    = 0,
		// 	@Company_Code  = @Company_Code,
		// 	@Trans_Line   = @Csr_Trfout_Line_No,
		// 	@Trans_Type   = @Trans_Type,
		// 	@Trans_Reason_Code = @Trans_Reason_Code,
		// 	@Ref_Sys_No   = @Trfout_Sys_No,
		// 	@Ref_Doc_No   = @Trfout_Doc_No,
		// 	@Ref_Date   = @Trfout_Date,
		// 	@Ref_Whs_Code  = @Whs_Code,
		// 	@Ref_Whs_Group  = @Whs_Group,
		// 	@Ref_Loc_Code  = @Csr_Loc_Code,
		// 	@Ref_Item_Code  = @Csr_Item_Code,
		// 	@Ref_Qty   = @Csr_Qty_Supply,
		// 	@Ref_Uom   = @Csr_Uom,
		// 	@Ref_Price   = 0,
		// 	@Ref_Ccy_Code  = '',
		// 	@Cogs    = @Hpp,
		// 	@Creation_User_Id = @Change_User_Id,
		// 	@Creation_Datetime = @Change_Datetime
		// 	-- End Insert amStockTransaction
		// FETCH NEXT FROM Csr_TO INTO
		// 	@Csr_Trfout_Line_No ,
		// 	@Csr_Loc_Code ,
		// 	@Csr_Item_Code ,
		// 	@Csr_Qty_Supply	,
		// 	@Csr_Uom
	}

	getEvent, errEvent := financeserviceapiutils.GetEventByCode("1601000000", "TRI", "SPTRIN")
	if errEvent != nil {
		return transactionsparepartentities.ItemWarehouseTransferOut{}, errEvent
	}

	transIn, errIn := NewItemWarehouseTransferInRepositoryImpl().InsertDetail(tx, transactionsparepartpayloads.InsertItemWarehouseHeaderTransferInRequest{
		EventId:                 getEvent.EventId,
		TransferInDate:          entities.TransferOutDate,
		TransferOutSystemNumber: entities.TransferOutSystemNumber,
		CompanyId:               entities.CompanyId,
		WarehouseId:             &entitiesRequest.RequestToWarehouseId,
	})

	if errIn != nil {
		return transactionsparepartentities.ItemWarehouseTransferOut{}, errIn
	}

	for _, out := range entitiesDetailOut {
		updates := map[string]interface{}{
			"quantity_received": out.QuantityOut,
			"location_id":       out.LocationIdTo,
		}
		errUpdate := tx.Model(&transactionsparepartentities.ItemWarehouseTransferInDetail{}).
			Where("transfer_in_system_number = ?", transIn.TransferInSystemNumber).
			Where("transfer_out_detail_system_number = ?", out.TransferOutDetailSystemNumber).
			Where("item_id = ?", out.ItemId).
			Updates(updates).Error

		if errUpdate != nil {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUpdate,
			}
		}
	}

	_, errTransInSubmit := NewItemWarehouseTransferInRepositoryImpl().Submit(tx, transIn.TransferInSystemNumber)
	if errTransInSubmit != nil {
		return transactionsparepartentities.ItemWarehouseTransferOut{}, errTransInSubmit
	}

	return entities, nil
}

// InsertDetailFromReceipt implements transactionsparepartrepository.ItemWarehouseTransferOutRepository.
func (*ItemWarehouseTransferOutRepositoryImpl) InsertDetailFromReceipt(tx *gorm.DB, request transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailCopyReceiptRequest) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest
	var entitiesOut transactionsparepartentities.ItemWarehouseTransferOut
	var entitiesDetail []transactionsparepartentities.ItemWarehouseTransferRequestDetail
	var entitiesDetailOut transactionsparepartentities.ItemWarehouseTransferOutDetail

	errEntities := tx.Model(&entitiesOut).Where(transactionsparepartentities.ItemWarehouseTransferOut{TransferOutSystemNumber: request.TransferOutSystemNumber}).First(&entitiesOut).Error
	if errEntities != nil {
		if errors.Is(errEntities, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errEntities,
				Message:    "transfer out with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errEntities,
			Message:    "failed to get transfer request entity",
		}
	}

	errGetEntities := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestSystemNumber: entitiesOut.TransferRequestSystemNumbers}).First(&entities).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntities,
				Message:    "transfer request with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get transfer request entity",
		}
	}

	errGetEntitiesDetail := tx.Model(&entitiesDetail).Where(transactionsparepartentities.ItemWarehouseTransferRequestDetail{TransferRequestSystemNumberId: entitiesOut.TransferRequestSystemNumbers}).Find(&entitiesDetail).Error
	if errGetEntitiesDetail != nil {
		if errors.Is(errGetEntitiesDetail, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntitiesDetail,
				Message:    "transfer request detail with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntitiesDetail,
			Message:    "failed to get transfer request entity",
		}
	}

	if len(entitiesDetail) == 0 {
		return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("this transfer request have no detail"),
		}
	}

	var locationItemFrom masteritementities.ItemLocation
	var locationItemTo masteritementities.ItemLocation

	for _, detail := range entitiesDetail {
		errGet := tx.Model(&locationItemFrom).
			Select(
				"TOP 1 item_location_id",
			).
			Joins("LEFT JOIN mtr_warehouse_master wr on wr.warehouse_id = ?", entities.RequestFromWarehouseId).
			// Where("company_id = ?").
			Where("mtr_location_item.warehouse_id = ?", entities.RequestFromWarehouseId).
			Where("mtr_location_item.warehouse_group_id = wr.warehouse_group_id"). //warehoouse can be deleted?
			Where("item_id = ?", detail.ItemId).
			Find(&locationItemFrom)
		if errGet.Error != nil {
			return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errGet.Error,
			}
		}

		errGetTo := tx.Model(&locationItemTo).
			Select(
				"TOP 1 item_location_id",
			).
			Joins("LEFT JOIN mtr_warehouse_master wr on wr.warehouse_id = ?", entities.RequestToWarehouseId).
			// Where("company_id = ?").
			Where("mtr_location_item.warehouse_id = ?", entities.RequestToWarehouseId).
			Where("mtr_location_item.warehouse_group_id = wr.warehouse_group_id").
			Where("item_id = ?", detail.ItemId).
			Find(&locationItemTo)
		if errGetTo.Error != nil {
			return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errGetTo.Error,
			}
		}

		entitiesDetailOut.TransferOutDetailSystemNumber = 0
		entitiesDetailOut.TransferOutSystemNumber = request.TransferOutSystemNumber
		entitiesDetailOut.TransferRequestDetailSystemNumber = detail.TransferRequestDetailSystemNumber
		entitiesDetailOut.LocationIdFrom = &locationItemFrom.ItemLocationId
		entitiesDetailOut.LocationIdTo = &locationItemTo.ItemLocationId
		if locationItemFrom.ItemLocationId == 0 {
			entitiesDetailOut.LocationIdFrom = nil
		}
		if locationItemTo.ItemLocationId == 0 {
			entitiesDetailOut.LocationIdTo = nil
		}
		entitiesDetailOut.ItemId = detail.ItemId
		entitiesDetailOut.QuantityOut = detail.RequestQuantity

		errCreate := tx.Create(&entitiesDetailOut).Error
		if errCreate != nil {
			return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errCreate,
			}
		}
	}

	return entitiesDetailOut, nil
}

// GetAllTransferOutDetail implements transactionsparepartrepository.ItemWarehouseTransferOutRepository.
func (*ItemWarehouseTransferOutRepositoryImpl) GetAllTransferOutDetail(tx *gorm.DB, number int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferOutDetail
	var entitiesHeader transactionsparepartentities.ItemWarehouseTransferOut
	var responses []transactionsparepartpayloads.GetAllDetailTransferOutResponse

	errGetEntities := tx.Model(&entitiesHeader).Where(transactionsparepartentities.ItemWarehouseTransferOut{TransferOutSystemNumber: number}).First(&entitiesHeader).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntities,
				Message:    "transfer out with that id is not found",
			}
		}
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get transfer request entity",
		}
	}

	joinTable := tx.Model(&entities).
		Select(
			"transfer_out_system_number",
			"transfer_out_detail_system_number",
			"trx_item_warehouse_transfer_out_detail.item_id",
			"trx_item_warehouse_transfer_out_detail.location_id_from location_id_from",
			"locF.warehouse_location_code location_code_from",
			"trx_item_warehouse_transfer_out_detail.location_id_to location_id_to",
			"locT.warehouse_location_code location_code_to",
			"it.item_name",
			"uom.uom_code unit_of_measurement",
			// "quantity_available",
			"det.request_quantity request_quantity",
			"quantity_out",
			"cost_of_goods_sold",
			"wmf.warehouse_group_id",
		).
		Joins("LEFT JOIN mtr_item it on it.item_id = trx_item_warehouse_transfer_out_detail.item_id").
		Joins("LEFT JOIN mtr_uom uom on uom.uom_id = it.unit_of_measurement_stock_id").
		Joins("LEFT JOIN trx_item_warehouse_transfer_request_detail det on det.transfer_request_detail_system_number = trx_item_warehouse_transfer_out_detail.transfer_request_detail_system_number").
		Joins("LEFT JOIN mtr_warehouse_master wmf on wmf.warehouse_id = ?", entitiesHeader.WarehouseId).
		Joins("LEFT JOIN mtr_warehouse_location locF on trx_item_warehouse_transfer_out_detail.location_id_from = locF.warehouse_location_id").
		Joins("LEFT JOIN mtr_warehouse_location locT on trx_item_warehouse_transfer_out_detail.location_id_to = locT.warehouse_location_id").
		Where("transfer_out_system_number = ?", number)

	err := joinTable.Scopes(pagination.Paginate(&pages, joinTable)).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	for i := range responses {

		get, err := masterwarehouserepository.NewLocationStockRepositoryImpl().GetAvailableQuantity(tx, masterwarehousepayloads.GetAvailableQuantityPayload{
			CompanyId:        entitiesHeader.CompanyId,
			PeriodDate:       entitiesHeader.TransferOutDate,
			WarehouseId:      *entitiesHeader.WarehouseId,
			LocationId:       responses[i].LocationIdFrom,
			ItemId:           responses[i].ItemId,
			WarehouseGroupId: responses[i].WarehouseGroupId,
		})
		responses[i].QuantityAvailable = get.QuantityAvailable
		if err != nil {
			return pagination.Pagination{}, err
		}

	}

	pages.Rows = responses
	return pages, nil
}

// GetAllTransferOut implements transactionsparepartrepository.ItemWarehouseTransferOutRepository.
func (*ItemWarehouseTransferOutRepositoryImpl) GetAllTransferOut(tx *gorm.DB, filter []utils.FilterCondition, dateParams map[string]string, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferOut
	var responses []transactionsparepartpayloads.GetAllTransferOutResponse

	joinTable := tx.Model(&entities).
		Select(
			"trx_item_warehouse_transfer_out.company_id",
			"transfer_out_system_number",
			"transfer_out_document_number",
			"transfer_out_status_id",
			"stat.item_transfer_status_code transfer_out_status_code",
			"stat.item_transfer_status_description transfer_out_status_description",
			"transfer_out_date",
			"transfer_request_system_number",
			"wmf.warehouse_id",
			"wmf.warehouse_name",
			"wgf.warehouse_group_id",
			"wgf.warehouse_group_name",
			"trx_item_warehouse_transfer_out.profit_center_id",
		).
		Joins("LEFT JOIN mtr_warehouse_master wmf on wmf.warehouse_id = trx_item_warehouse_transfer_out.warehouse_id").
		Joins("LEFT JOIN mtr_item_transfer_status stat on stat.item_transfer_status_id = transfer_out_status_id").
		Joins("LEFT JOIN mtr_warehouse_group wgf on wgf.warehouse_group_id = wmf.warehouse_group_id")

	whereQuery := utils.ApplyFilter(joinTable, filter)

	var strDateFilter string
	if dateParams["transfer_out_date_from"] == "" {
		dateParams["transfer_out_date_from"] = "19000101"
	}
	if dateParams["transfer_out_date_to"] == "" {
		dateParams["transfer_out_date_to"] = "99991212"
	}
	strDateFilter = "transfer_out_date >='" + dateParams["transfer_out_date_from"] + "' AND transfer_out_date <= '" + dateParams["transfer_out_date_to"] + "'"

	whereQuery = whereQuery.Where(strDateFilter)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	pages.Rows = responses
	return pages, nil
}

// GetTransferOutById implements transactionsparepartrepository.ItemWarehouseTransferOutRepository.
func (*ItemWarehouseTransferOutRepositoryImpl) GetTransferOutById(tx *gorm.DB, number int) (transactionsparepartpayloads.GetTransferOutByIdResponse, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferOut
	var warehouseEntities masterwarehouseentities.WarehouseMaster
	var responses transactionsparepartpayloads.GetTransferOutByIdResponse
	errGetEntities := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferOut{TransferOutSystemNumber: number}).First(&entities).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return transactionsparepartpayloads.GetTransferOutByIdResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntities,
				Message:    "transfer request with that id is not found",
			}
		}
		return transactionsparepartpayloads.GetTransferOutByIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get transfer request entity",
		}
	}

	var transferRequestStatus masteritementities.ItemTransferStatus
	errGetTransferStatus := tx.Model(&transferRequestStatus).Where(masteritementities.ItemTransferStatus{ItemTransferStatusId: entities.TransferOutStatusId}).First(&transferRequestStatus).Error
	if errGetTransferStatus != nil {
		if errors.Is(errGetTransferStatus, gorm.ErrRecordNotFound) {
			return transactionsparepartpayloads.GetTransferOutByIdResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "transfer status with that id is not found please check input",
				Err:        errGetTransferStatus,
			}
		}
		return transactionsparepartpayloads.GetTransferOutByIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetTransferStatus,
			Message:    "failed to get transfer status please check input",
		}
	}

	errWarehouse := tx.Model(&warehouseEntities).Where(masterwarehouseentities.WarehouseMaster{WarehouseId: *entities.WarehouseId}).Find(&warehouseEntities)
	if errWarehouse.Error != nil {
		return transactionsparepartpayloads.GetTransferOutByIdResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errWarehouse.Error,
			Message:    "failed to get warehouse",
		}
	}

	responses.CompanyId = entities.CompanyId
	responses.TransferOutSystemNumber = entities.TransferOutSystemNumber
	responses.TransferOutDate = entities.TransferOutDate
	responses.TransferOutDocumentSystemNumber = entities.TransferOutDocumentNumber
	responses.TransferRequestSystemNumber = entities.TransferOutSystemNumber
	responses.ProfitCenterId = entities.ProfitCenterId
	responses.WarehouseId = *entities.WarehouseId
	responses.WarehouseCode = warehouseEntities.WarehouseCode
	responses.WarehouseName = warehouseEntities.WarehouseName
	responses.WarehouseGroupId = warehouseEntities.WarehouseGroupId
	responses.TransferStatusId = entities.TransferOutStatusId

	return responses, nil
}

// InsertDetail implements transactionsparepartrepository.ItemWarehouseTransferOutRepository.
func (*ItemWarehouseTransferOutRepositoryImpl) InsertDetail(tx *gorm.DB, request transactionsparepartpayloads.InsertItemWarehouseTransferOutDetailRequest) (transactionsparepartentities.ItemWarehouseTransferOutDetail, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferOut
	var entitiesDetail transactionsparepartentities.ItemWarehouseTransferOutDetail
	var entitiesTransferRequest transactionsparepartentities.ItemWarehouseTransferRequestDetail

	fmt.Println(request.TransferOutSystemNumber)

	errGetEntities := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferOut{TransferOutSystemNumber: request.TransferOutSystemNumber}).First(&entities).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return entitiesDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntities,
				Message:    "transfer out with that id is not found",
			}
		}
		return entitiesDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get transfer out entity",
		}
	}

	errGetEntitiesReq := tx.Model(&entitiesTransferRequest).Where(transactionsparepartentities.ItemWarehouseTransferRequestDetail{TransferRequestDetailSystemNumber: request.TransferRequestDetailSystemNumber}).First(&entitiesTransferRequest).Error
	if errGetEntitiesReq != nil {
		if errors.Is(errGetEntitiesReq, gorm.ErrRecordNotFound) {
			return entitiesDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntitiesReq,
				Message:    "transfer request with that id is not found",
			}
		}
		return entitiesDetail, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntitiesReq,
			Message:    "failed to get transfer request entity",
		}
	}

	errGetEntitiesDetail := tx.Model(&entitiesDetail).Where("transfer_out_system_number = ?", request.TransferOutSystemNumber).
		Where(transactionsparepartentities.ItemWarehouseTransferOutDetail{ItemId: entitiesTransferRequest.ItemId}).First(&entitiesDetail).Error
	if errGetEntitiesDetail != nil {
		if !errors.Is(errGetEntitiesDetail, gorm.ErrRecordNotFound) {
			return entitiesDetail, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errGetEntitiesDetail,
				Message:    "failed to get transfer out detail entity",
			}
		}
	}

	fmt.Println(entitiesDetail.TransferOutDetailSystemNumber)

	if entitiesDetail.TransferOutDetailSystemNumber != 0 {
		return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("item already exist in this transfer out"),
		}
	}

	if request.QuantityOut > entitiesTransferRequest.RequestQuantity {
		return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("not enough request quantity for this request"),
		}
	}

	entitiesDetail.ItemId = entitiesTransferRequest.ItemId
	entitiesDetail.TransferOutSystemNumber = request.TransferOutSystemNumber
	entitiesDetail.TransferRequestDetailSystemNumber = request.TransferRequestDetailSystemNumber
	entitiesDetail.QuantityOut = request.QuantityOut
	entitiesDetail.LocationIdFrom = &request.LocationIdFrom
	entitiesDetail.LocationIdTo = &request.LocationIdTo

	errCreate := tx.Create(&entitiesDetail).Scan(&entitiesDetail).Error
	if errCreate != nil {
		return transactionsparepartentities.ItemWarehouseTransferOutDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errCreate,
		}
	}

	return entitiesDetail, nil
}

// InsertHeader implements transactionsparepartrepository.ItemWarehouseTransferOutRepository.
func (*ItemWarehouseTransferOutRepositoryImpl) InsertHeader(tx *gorm.DB, request transactionsparepartpayloads.InsertItemWarehouseHeaderTransferOutRequest) (transactionsparepartentities.ItemWarehouseTransferOut, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferOut
	var entitiesRequest transactionsparepartentities.ItemWarehouseTransferRequest
	var status masteritementities.ItemTransferStatus

	errGetEntities := tx.Model(&entitiesRequest).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestSystemNumber: request.TransferRequestSystemNumber}).First(&entitiesRequest).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntities,
				Message:    "transfer request with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get transfer request entity",
		}
	}

	errGetStatus := tx.Model(&status).Where("item_transfer_status_code = ?", 10).Find(&status)
	if errGetStatus.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatus.Error,
		}
	}

	entities.CompanyId = request.CompanyId
	entities.TransferOutDate = request.TransferOutDate
	entities.TransferRequestSystemNumbers = request.TransferRequestSystemNumber
	entities.WarehouseId = &entitiesRequest.RequestFromWarehouseId
	entities.TransferOutStatusId = status.ItemTransferStatusId

	errCreate := tx.Create(&entities).Scan(&entities).Error
	if errCreate != nil {
		return transactionsparepartentities.ItemWarehouseTransferOut{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errCreate,
		}
	}

	return entities, nil
}
