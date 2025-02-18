package transactionsparepartrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	masterwarehouserepositoryimpl "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	"errors"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func NewItemWarehouseTransferInRepositoryImpl() transactionsparepartrepository.ItemWarehouseTransferInRepository {
	return &ItemWarehouseTransferInRepositoryImpl{}
}

type ItemWarehouseTransferInRepositoryImpl struct {
}

// Submit implements transactionsparepartrepository.ItemWarehouseTransferInRepository.
func (*ItemWarehouseTransferInRepositoryImpl) Submit(tx *gorm.DB, number int) (transactionsparepartentities.ItemWarehouseTransferIn, *exceptions.BaseErrorResponse) {
	var entitiesIn transactionsparepartentities.ItemWarehouseTransferIn
	var entitiesOut transactionsparepartentities.ItemWarehouseTransferOut
	var entitiesRequest transactionsparepartentities.ItemWarehouseTransferRequest
	var entitiesInDetail []transactionsparepartentities.ItemWarehouseTransferInDetail
	var statusClosed masteritementities.ItemTransferStatus

	errGetStatus := tx.Model(&statusClosed).Where("item_transfer_status_code = ?", 50).Find(&statusClosed)
	if errGetStatus.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatus.Error,
		}
	}

	errGetEntities := tx.Model(&entitiesIn).Where(transactionsparepartentities.ItemWarehouseTransferIn{TransferInSystemNumber: number}).First(&entitiesIn).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntities,
				Message:    "transfer in with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get transfer in entity",
		}
	}

	errGetEntitiesOut := tx.Model(&entitiesOut).Where(transactionsparepartentities.ItemWarehouseTransferOut{TransferOutSystemNumber: entitiesIn.TransferOutSystemNumberId}).First(&entitiesOut).Error
	if errGetEntitiesOut != nil {
		if errors.Is(errGetEntitiesOut, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntities,
				Message:    "transfer in with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntitiesOut,
			Message:    "failed to get transfer in entity",
		}
	}

	errGetEntitiesReq := tx.Model(&entitiesRequest).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestSystemNumber: entitiesOut.TransferRequestSystemNumbers}).First(&entitiesRequest).Error
	if errGetEntitiesReq != nil {
		if errors.Is(errGetEntitiesReq, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntities,
				Message:    "transfer request with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntitiesReq,
			Message:    "failed to request transfer in entity",
		}
	}

	getPeriod, errPeriod := financeserviceapiutils.GetOpenPeriodByCompany(entitiesIn.CompanyId, "SP")
	if errPeriod != nil {
		return transactionsparepartentities.ItemWarehouseTransferIn{}, errPeriod
	}

	// EXEC uspg_gmSrcDoc1_Update
	// 	@Option = 0 ,
	// 	@COMPANY_CODE = @Company_Code ,
	// 	@SOURCE_CODE = @Src_Code ,
	// 	@VEHICLE_BRAND = @Whs_Brand ,
	// 	@PROFIT_CENTER_CODE = '' ,
	// 	@TRANSACTION_CODE = '' ,
	// 	@BANK_ACC_CODE = '' ,
	// 	@TRANSACTION_DATE = @Trfin_Date ,
	// 	@Last_Doc_No =  @Trfin_Doc_No OUTPUT

	// 	IF ISNULL(@Trfin_Doc_No ,'') = ''
	// 	BEGIN
	// 		RAISERROR('Document no. cannot be generated.',16,1)
	// 		RETURN 0
	// 	END

	entitiesRequest.TransferRequestStatusId = statusClosed.ItemTransferStatusId
	errSaveRequest := tx.Save(&entitiesRequest).Error
	if errSaveRequest != nil {
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errSaveRequest,
		}
	}

	entitiesOut.TransferOutStatusId = statusClosed.ItemTransferStatusId
	errSaveOut := tx.Save(&entitiesOut).Error
	if errSaveOut != nil {
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errSaveOut,
		}
	}

	entitiesIn.TransferInStatusId = statusClosed.ItemTransferStatusId
	errSaveIn := tx.Save(&entitiesIn).Error
	if errSaveIn != nil {
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errSaveIn,
		}
	}

	errGetIn := tx.Model(&entitiesInDetail).Where("transfer_in_system_number = ?", entitiesIn.TransferInSystemNumber).Find(&entitiesInDetail)
	if errGetIn.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetIn.Error,
		}
	}

	getWarehouseTo, errWare := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl().GetById(tx, entitiesRequest.RequestToWarehouseId)
	if errWare != nil {
		return transactionsparepartentities.ItemWarehouseTransferIn{}, errWare
	}

	getWarehouseFrom, errFrom := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl().GetById(tx, entitiesRequest.RequestFromWarehouseId)
	if errFrom != nil {
		return transactionsparepartentities.ItemWarehouseTransferIn{}, errFrom
	}

	var locationStock masterentities.LocationStock
	var hpp float64
	var inHpp float64
	for _, detail := range entitiesInDetail {
		updates := map[string]interface{}{
			"location_id_to": detail.LocationId,
		}
		errUpdate := tx.Model(&transactionsparepartentities.ItemWarehouseTransferRequestDetail{}).
			Where("transfer_request_system_number = ?", entitiesRequest.TransferRequestSystemNumber).
			Where("item_id = ?", detail.ItemId).
			Updates(updates).Error

		if errUpdate != nil {
			return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUpdate,
			}
		}

		// periodYear := entities.TransferOutDate.Format("2006")
		// periodMonth := entities.TransferOutDate.Format("01")

		errGetLocStock := tx.Model(&locationStock).
			Select(
				"ISNULL(SUM(quantity_ending), 0) quantity_ending",
			).
			Where("company_id = ?", entitiesIn.CompanyId).
			Where("period_month = ?", getPeriod.PeriodMonth).
			Where("period_year = ?", getPeriod.PeriodYear).
			Where("item_id = ?", detail.ItemId).
			Where("warehouse_group_id = ?", getWarehouseTo.WarehouseGroupId).
			Find(&locationStock).Error

		if errGetLocStock != nil {
			return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errGetLocStock,
			}
		}
		if locationStock.QuantityEnding == 0 {
			year, errYear := strconv.Atoi(getPeriod.PeriodYear)
			if errYear != nil {
				return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errYear,
				}
			}

			month, errMonth := strconv.Atoi(getPeriod.PeriodMonth)
			if errMonth != nil {
				return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errMonth,
				}
			}

			date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			periodYear := date.Format("2006")
			periodMonth := date.Format("01")
			errGetLocStock := tx.Model(&locationStock).
				Select(
					"ISNULL(SUM(quantity_ending), 0) quantity_ending",
				).
				Where("company_id = ?", entitiesIn.CompanyId).
				Where("period_month = ?", periodMonth).
				Where("period_year = ?", periodYear).
				Where("item_id = ?", detail.ItemId).
				Where("warehouse_group_id = ?", getWarehouseTo.WarehouseGroupId).
				Find(&locationStock).Error

			if errGetLocStock != nil {
				return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errGetLocStock,
				}
			}
		}

		var trfOutDet transactionsparepartentities.ItemWarehouseTransferOutDetail
		if getWarehouseFrom.WarehouseGroupId == getWarehouseTo.WarehouseGroupId {
			errGetOutDetail := tx.Model(&trfOutDet).
				Select(
					"trx_item_warehouse_transfer_out_detail.cost_of_goods_sold",
				).
				Joins("INNER JOIN trx_item_warehouse_transfer_request b on b.transfer_out_system_number = trx_item_warehouse_transfer_out_detail.transfer_out_system_number").
				Joins("INNER JOIN trx_item_warehouse_transfer_in_detail c on b.transfer_in_system_number = c.transfer_in_system_number and c.transfer_out_detail_system_number = trx_item_warehouse_transfer_out_detail.transfer_out_detail_system_number").
				Where("c.transfer_in_system_number = ?", detail.TransferInDetailSystemNumber).Scan(&hpp)
			if errGetOutDetail.Error != nil {
				return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        errGetOutDetail.Error,
				}
			}

			if hpp == 0 {
				tx.Model(&masterentities.GroupStock{}).
					Select(
						"price_current",
					).
					Where("company_id = ?").Where("warehouse_group_id = ?").Where("item_id = ?").
					Where("period_year = ?").Where("period_month = ?").Scan(&hpp)
			}
		} else {
			if hpp == 0 {
				errGetOutDetail := tx.Model(&trfOutDet).
					Select(
						"trx_item_warehouse_transfer_out_detail.cost_of_goods_sold",
					).
					Joins("INNER JOIN trx_item_warehouse_transfer_request b on b.transfer_out_system_number = trx_item_warehouse_transfer_out_detail.transfer_out_system_number").
					Joins("INNER JOIN trx_item_warehouse_transfer_in_detail c on b.transfer_in_system_number = c.transfer_in_system_number and c.transfer_out_detail_system_number = trx_item_warehouse_transfer_out_detail.transfer_out_detail_system_number").
					Where("c.transfer_in_system_number = ?", detail.TransferInDetailSystemNumber).Scan(&hpp)
				if errGetOutDetail.Error != nil {
					return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Err:        errGetOutDetail.Error,
					}
				}

				if hpp == 0 {
					tx.Model(&masterentities.GroupStock{}).
						Select(
							"price_current",
						).
						Where("company_id = ?", entitiesOut.CompanyId).Where("warehouse_group_id = ?", getWarehouseFrom.WarehouseGroupId).Where("item_id = ?", detail.ItemId).
						Where("period_year = ?", getPeriod.PeriodYear).Where("period_month = ?", getPeriod.PeriodMonth).Scan(&hpp)
				}
			} else {
				errGetOutDetail := tx.Model(&trfOutDet).
					Select(
						"trx_item_warehouse_transfer_out_detail.cost_of_goods_sold",
					).
					Joins("INNER JOIN trx_item_warehouse_transfer_request b on b.transfer_out_system_number = trx_item_warehouse_transfer_out_detail.transfer_out_system_number").
					Joins("INNER JOIN trx_item_warehouse_transfer_in_detail c on b.transfer_in_system_number = c.transfer_in_system_number and c.transfer_out_detail_system_number = trx_item_warehouse_transfer_out_detail.transfer_out_detail_system_number").
					Where("c.transfer_in_system_number = ?", detail.TransferInDetailSystemNumber).Scan(&hpp)
				if errGetOutDetail.Error != nil {
					return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Err:        errGetOutDetail.Error,
					}
				}

				if hpp == 0 {
					tx.Model(&masterentities.GroupStock{}).
						Select(
							"price_current",
						).
						Where("company_id = ?", entitiesOut.CompanyId).Where("warehouse_group_id = ?", getWarehouseFrom.WarehouseGroupId).Where("item_id = ?", detail.ItemId).
						Where("period_year = ?", getPeriod.PeriodYear).Where("period_month = ?", getPeriod.PeriodMonth).Scan(&hpp)
				}

				var grpStckRes transactionsparepartpayloads.SubmitItemWarehouseTransferOutGroupStock

				errgroup := tx.Table("mtr_group_stock as a").
					Select(
						"price_current",
						"SUM(ISNULL(), 0) quantity_ending",
					).
					Joins("LEFT JOIN mtr_location_stock loc on loc.company_id = a.company_id and loc.item_id = a.item_id "+
						"and a.period_month = loc.period_month and loc.period_year = a.period_year and warehouse_id in (?)").
					Where("company_id = ?", entitiesOut.CompanyId).
					Where("warehouse_group_id = ?", getWarehouseFrom.WarehouseGroupId).
					Where("item_id = ?", detail.ItemId).
					Where("period_year = ?", getPeriod.PeriodYear).
					Where("period_month = ?", getPeriod.PeriodMonth).Scan(&grpStckRes)

				if errgroup.Error != nil {
					return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
					}
				}

				inHpp = ((grpStckRes.PriceCurrent*grpStckRes.QuantityEnding)+(hpp*detail.QuantityReceived))/grpStckRes.QuantityEnding + detail.QuantityReceived
				// 	SELECT @PriceCurrent = A.PRICE_CURRENT,
				// 		   @QtyEnding = SUM(ISNULL(B.QTY_ENDING,0))		--B.QTY_ENDING		Fali 30 Apr 2014 : Hasil select bisa multiple row, qty harus di sum
				// 	FROM amGroupStock A
				// 	LEFT JOIN amLocationStock B ON A.COMPANY_CODE = B.COMPANY_CODE
				// 								and A.ITEM_CODE = B.ITEM_CODE
				// 								and A.PERIOD_MONTH = B.PERIOD_MONTH
				// 								and A.PERIOD_YEAR = B.PERIOD_YEAR
				// 								--and B.WHS_CODE = @To_Whs_Code
				// 								and B.WHS_GROUP = @To_Whs_Group
				// 								and B.WHS_CODE IN (SELECT C.WAREHOUSE_CODE
				// 												   FROM gmLoc1 C
				// 												   WHERE B.COMPANY_CODE = C.COMPANY_CODE
				// 												   AND B.WHS_CODE = C.WAREHOUSE_CODE
				// 												   AND C.COSTING_TYPE = @varHPP_WH_TYPE_NORMAL)
				// 	WHERE A.COMPANY_CODE = @Company_Code
				// 	and A.PERIOD_YEAR = @Period_Year
				// 	and A.PERIOD_MONTH = @Period_Month
				// 	and A.ITEM_CODE = @Csr_Item_Code
				// 	AND A.WHS_GROUP = @To_Whs_Group
				// 	GROUP BY A.PRICE_CURRENT,A.ITEM_CODE,B.WHS_GROUP--,B.WHS_CODE

				// 	SET @In_Hpp = ((ISNULL(@PriceCurrent,0) * ISNULL(@QtyEnding,0)) + (ISNULL(@HPP,0) * ISNULL(@CSR_QTY_RECEIVE,0)))/(ISNULL(@QtyEnding,0) + ISNULL(@Csr_Qty_Receive,0)
			}
		}
		if inHpp == 0 {
			inHpp = hpp
		}

		// UPDATE attrfIn1
		// 	SET COGS = ISNULL(@In_Hpp,0),
		// 		TRF_COST = (QTY_RECEIVE * ISNULL(@Hpp,0)),--(QTY_RECEIVE * ISNULL(@In_Hpp,0)),
		// 		REF_PRICE = ISNULL(@Hpp,0),
		// 		HPP_VARIANCE = (ISNULL(@Hpp,0) - ISNULL(@In_Hpp,0))*QTY_RECEIVE
		// 	WHERE TRFIN_SYS_NO = @Trfin_Sys_No
		// 	AND TRFIN_LINE = @Csr_Trfin_Line
		// 	AND ITEM_CODE = @Csr_Item_Code

		// 	-- HPP --
		// 	EXEC [dbo].[uspg_atStockTransaction_Insert]
		// 	@Option    = 0,
		// 	@Company_Code  = @Company_Code,
		// 	@Trans_Line   = @Csr_Trfin_Line,
		// 	@Trans_Type   = @Trans_Type,
		// 	@Trans_Reason_Code = @Trans_Reason_Code,
		// 	@Ref_Sys_No   = @Trfin_Sys_No,
		// 	@Ref_Doc_No   = @Trfin_Doc_No,
		// 	@Ref_Date   = @Trfin_Date,
		// 	@Ref_Whs_Code  = @To_Whs_Code,
		// 	@Ref_Whs_Group  = @To_Whs_Group,
		// 	@Ref_Loc_Code  = @Csr_Loc_Code,
		// 	@Ref_Item_Code  = @Csr_Item_Code,
		// 	@Ref_Qty   = @Csr_Qty_Receive,
		// 	@Ref_Uom   = @Csr_Uom,
		// 	@Ref_Price   = @Hpp,--@In_Hpp,
		// 	@Ref_Ccy_Code  = '',
		// 	@Cogs    = @In_Hpp,--@Hpp,
		// 	@Creation_User_Id = @Change_User_Id,
		// 	@Creation_Datetime = @Change_Datetime
	}

	// DECLARE @Process_Code varchar(30)= dbo.getVariableValue('GL_PROCESS_CODE_SP_TRFITEM'),
	// 			@Profit_Center varchar(5) = dbo.getVariableValue('PROFIT_CENTER_SP'),
	// 			@Trx_Type varchar(10) = dbo.getVariableValue('TRXTYPE_SP_TRFITEM_IN'),
	// 			@Journal_Sys_No numeric(15,0) = 0,
	// 			@WarehouseFrom_Brand varchar(10) = '',
	// 			@WarehouseTo_Brand varchar(10)=''

	// 	SELECT @WarehouseTo_Brand = ISNULL(B.WAREHOUSE_BRAND,''),
	// 		   @WarehouseFrom_Brand = ISNULL(C.WAREHOUSE_BRAND,'')
	// 	FROM atTrfReq0 A
	// 	INNER JOIN gmLoc1 B ON A.COMPANY_CODE=B.COMPANY_CODE AND A.REQ_TO_WHS_CODE = B.WAREHOUSE_CODE
	// 	INNER JOIN gmLoc1 C ON A.COMPANY_CODE=C.COMPANY_CODE AND A.REQ_FROM_WHS_CODE = C.WAREHOUSE_CODE
	// 	WHERE A.TRFREQ_SYS_NO = @Trfreq_Sys_No

	// 	IF @WarehouseTo_Brand <> @WarehouseFrom_Brand AND @WarehouseTo_Brand <> '' AND @WarehouseFrom_Brand <> ''
	// 	BEGIN
	// 		-- Generate Journal
	// 		EXEC usp_comJournalAction
	// 		@Process_Code = @Process_Code,
	// 		@Cpc_Code = @Profit_Center,
	// 		@Trx_Type = @Trx_Type,
	// 		@Event_No = @Event_No,
	// 		@Journal_Sys_No = @Journal_Sys_No Output,
	// 		@Ref_Sys_No = @TrfIn_Sys_No,
	// 		@Ref_Doc_No = @TrfIn_Doc_No,
	// 		@Creation_User_Id = @Change_User_Id

	// 		-- Update Journal Sys No Pada attrfin0
	// 		UPDATE atTrfIn0
	// 		SET JOURNAL_SYS_NO = @Journal_Sys_No,
	// 			CPC_CODE = @Profit_Center
	// 		WHERE TRFIN_SYS_NO = @TrfIn_Sys_No

	// 		-- Update Journal Sys No Pada attrfout0
	// 		UPDATE atTrfOut0
	// 		SET JOURNAL_SYS_NO = @Journal_Sys_No ,
	// 			CPC_CODE = @Profit_Center
	// 		WHERE TRFOUT_SYS_NO = @Trfout_Sys_No
	// 	END

	return transactionsparepartentities.ItemWarehouseTransferIn{}, nil
}

// Insert implements transactionsparepartrepository.ItemWarehouseTransferInRepository.
func (*ItemWarehouseTransferInRepositoryImpl) InsertDetail(tx *gorm.DB, request transactionsparepartpayloads.InsertItemWarehouseHeaderTransferInRequest) (transactionsparepartentities.ItemWarehouseTransferIn, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferIn
	var entitiesDetail transactionsparepartentities.ItemWarehouseTransferInDetail
	var entitiesOut transactionsparepartentities.ItemWarehouseTransferOut
	var entitiesOutDetail []transactionsparepartentities.ItemWarehouseTransferOutDetail
	var statusDraft masteritementities.ItemTransferStatus

	errGetStatus := tx.Model(&statusDraft).Where("item_transfer_status_code = ?", 10).Find(&statusDraft)
	if errGetStatus.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatus.Error,
		}
	}

	errGetEntities := tx.Model(&entitiesOut).Where(transactionsparepartentities.ItemWarehouseTransferOut{TransferOutSystemNumber: request.TransferOutSystemNumber}).First(&entitiesOut).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntities,
				Message:    "transfer out with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get transfer out entity",
		}
	}

	errGetEntitiesOutDetail := tx.Model(&entitiesOutDetail).Where(transactionsparepartentities.ItemWarehouseTransferOutDetail{TransferOutSystemNumber: request.TransferOutSystemNumber}).
		Find(&entitiesOutDetail).Error
	if errGetEntitiesOutDetail != nil {
		if errors.Is(errGetEntitiesOutDetail, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errGetEntitiesOutDetail,
				Message:    "transfer out with that id is not found",
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get transfer out entity",
		}
	}

	entities.CompanyId = request.CompanyId
	entities.EventId = request.EventId
	entities.TransferOutSystemNumberId = request.TransferOutSystemNumber
	entities.TransferInDate = request.TransferInDate
	entities.WarehouseId = request.WarehouseId
	entities.TransferInStatusId = statusDraft.ItemTransferStatusId

	errCreate := tx.Create(&entities).Scan(&entities).Error
	if errCreate != nil {
		return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errCreate,
		}
	}

	for _, detail := range entitiesOutDetail {
		entitiesDetail.TransferInDetailSystemNumber = 0
		entitiesDetail.TransferInSystemNumberId = entities.TransferInSystemNumber
		entitiesDetail.TransferOutDetailSystemNumberId = detail.TransferOutDetailSystemNumber
		entitiesDetail.QuantityReceived = 0
		entitiesDetail.ItemId = *detail.ItemId

		errCreateDetail := tx.Create(&entitiesDetail).Scan(&entitiesDetail).Error
		if errCreateDetail != nil {
			return transactionsparepartentities.ItemWarehouseTransferIn{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errCreateDetail,
			}
		}
	}

	return entities, nil
}
