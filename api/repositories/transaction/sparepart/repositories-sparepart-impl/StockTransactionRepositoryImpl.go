package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
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

type StockTransactionRepositoryImpl struct {
}

func StartStockTransactionRepositoryImpl() transactionsparepartrepository.StockTransactionRepository {
	return &StockTransactionRepositoryImpl{}
}
func ConvertMonth(month string) string {
	if len(month) == 0 {
		return "0" + month
	}
	return month
}
func (s *StockTransactionRepositoryImpl) StockTransactionInsert(db *gorm.DB, payloads transactionsparepartpayloads.StockTransactionInsertPayloads) (bool, *exceptions.BaseErrorResponse) {
	//select type and trans reason for option
	//this repo is make for rewrite uspg_atStockTransaction_Insert.sql Devloper that needed the endpoint
	//for special trans type or trans reason please add condition in this endpoint

	//select trans type
	var stockTransactionType masterentities.StockTransactionType
	err := db.Model(&stockTransactionType).
		Where(masterentities.StockTransactionType{StockTransactionTypeId: payloads.TransactionTypeId}).
		First(&stockTransactionType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("stock transaction type not found"),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("failed to get stock transaction on error : %s", err.Error()),
		}
	}
	//get transaction reason
	var stockTransactionReason masterentities.StockTransactionReason
	err = db.Model(&stockTransactionReason).Where(masterentities.StockTransactionReason{StockTransactionReasonId: payloads.TransactionReasonId}).
		First(&stockTransactionReason).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("stock transaction reason not found"),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("failed to get stock transaction reason on error : %s", err.Error()),
		}
	}
	//DECLARE @Item_Class Varchar(10)
	//SELECT @Item_Class = ITEM_CLASS FROM gmItem0 WHERE ITEM_CODE = @Ref_Item_Code
	//getting item master from reference item
	var ItemEntities masteritementities.Item
	err = db.Model(&ItemEntities).Where(masteritementities.Item{ItemId: payloads.ReferenceItemId}).
		First(&ItemEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New(fmt.Sprintf("failed to get item entity with Item Id : %d", payloads.ReferenceItemId)),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("failed to get item entity on error : %s", err.Error()),
		}
	}
	if payloads.ReferenceDate == nil {
		payloads.ReferenceDate = new(time.Time)
	}
	var newStockTransaction transactionsparepartentities.StockTransaction
	newStockTransaction.CompanyId = payloads.CompanyId
	newStockTransaction.IsActive = true
	newStockTransaction.TransactionTypeId = payloads.TransactionTypeId
	newStockTransaction.TransactionReasonId = payloads.TransactionReasonId
	newStockTransaction.ReferenceId = payloads.ReferenceId
	newStockTransaction.ReferenceDocumentNumber = payloads.ReferenceDocumentNumber
	newStockTransaction.ReferenceDate = *payloads.ReferenceDate
	newStockTransaction.ReferenceWarehouseId = payloads.ReferenceWarehouseId
	newStockTransaction.ReferenceWarehouseGroupId = payloads.ReferenceWarehouseGroupId
	newStockTransaction.ReferenceLocationId = payloads.ReferenceLocationId
	newStockTransaction.ReferenceItemId = payloads.ReferenceItemId
	newStockTransaction.ReferenceQuantity = payloads.ReferenceQuantity
	newStockTransaction.ReferenceUnitOfMeasurementId = payloads.ReferenceUnitOfMeasurementId
	newStockTransaction.ReferencePrice = payloads.ReferencePrice
	newStockTransaction.ReferenceCurrencyId = payloads.ReferenceCurrencyId
	newStockTransaction.TransactionCogs = payloads.TransactionCogs
	newStockTransaction.ChangeNo += 1
	newStockTransaction.CreatedByUserId = payloads.CreatedByUserId
	newStockTransaction.CreatedDate = payloads.CreatedDate
	newStockTransaction.UpdatedDate = payloads.UpdatedDate
	newStockTransaction.UpdatedByUserId = payloads.UpdatedByUserId
	newStockTransaction.VehicleId = payloads.VehicleId
	newStockTransaction.ItemClassId = ItemEntities.ItemClassId
	//insert to stock transaction
	err = db.Create(&newStockTransaction).Scan(&newStockTransaction).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("failed to insert stock transaction on error : %s", err.Error()),
		}
	}
	//DECLARE @Period_Year varchar(4), @Period_Month Varchar(2), @Ref_Qty_Negatif numeric(10,2)
	var periodYear string
	var periodMonth string
	var referenceQuantityNegative float64
	//
	if payloads.ReferenceDate.IsZero() {
		periodYear = strconv.Itoa(time.Now().UTC().Year())
		periodMonth = ConvertMonth(strconv.Itoa(int(time.Now().UTC().Month())))
	} else {
		periodYear = payloads.ReferenceDate.Format("2006")
		periodMonth = ConvertMonth(payloads.ReferenceDate.Format("01"))
	}
	//validate item is sellable
	//SET @Sellable = (SELECT SALES_ALLOW FROM gmLoc1 WHERE WAREHOUSE_CODE = @Ref_Whs_Code AND COMPANY_CODE = @Company_Code)
	//change to true
	var sellable bool = false
	err = db.Model(&masterwarehouseentities.WarehouseMaster{}).Select("warehouse_sales_allow").
		Where(&masterwarehouseentities.WarehouseMaster{WarehouseId: payloads.ReferenceWarehouseId, CompanyId: payloads.CompanyId}).
		Scan(&sellable).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("failed to get warehouse master on error : %s", err.Error()),
		}
	}
	if !sellable {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("item on this warehouse is not sellable : warehouse id : %d", payloads.ReferenceWarehouseId),
		}
	}
	//SET @Ref_Qty_Negatif = (-1 * @Ref_Qty)
	referenceQuantityNegative = payloads.ReferenceQuantity * -1
	//SET @TransType_PU = dbo.getVariableValue('STK_TRXTYPE_PURCHASE_RECEIPT')
	//@Reason_NL
	//CEK LINE 377 TEMP DLS
	//Making For Trans Type PU and trans reason Normal or Back Order
	//For another trans type add below
	if stockTransactionType.StockTransactionTypeCode == "PU" &&
		(stockTransactionReason.StockTransactionReasonCode == "NL" ||
			stockTransactionReason.StockTransactionReasonCode == "BO" ||
			stockTransactionReason.StockTransactionReasonCode == "WP") {
		RequestBodyLocationStock := masterwarehousepayloads.LocationStockUpdatePayloads{}
		//Exec uspg_amLocationStock_Update
		//@Option			= 0 ,
		//@Company_Code	= @Company_Code ,
		//@Period_Year	= @Period_Year ,
		//@Period_Month	= @Period_Month	 ,
		//@Whs_Code		= @Ref_Whs_Code ,
		//@Whs_Group		= @Ref_Whs_Group	,
		//@Loc_Code		= @Ref_Loc_Code ,
		//@Item_Code		= @Ref_Item_Code ,
		//@Qty_Purchase	= @Ref_Qty	,
		//@Trans_Type		= @Trans_Type ,
		//@Trans_Reason_Code = @Trans_Reason_Code ,
		//@Change_User_Id	= @Creation_User_Id ,
		//@Change_Datetime = @Creation_Datetime
		RequestBodyLocationStock.CompanyId = payloads.CompanyId
		RequestBodyLocationStock.PeriodYear = periodYear
		RequestBodyLocationStock.PeriodMonth = periodMonth
		RequestBodyLocationStock.WarehouseId = payloads.ReferenceWarehouseId
		RequestBodyLocationStock.WarehouseGroupId = payloads.ReferenceWarehouseGroupId
		RequestBodyLocationStock.LocationId = payloads.ReferenceLocationId
		RequestBodyLocationStock.ItemId = payloads.ReferenceItemId
		RequestBodyLocationStock.QuantityPurchase = payloads.ReferenceQuantity
		RequestBodyLocationStock.StockTransactionTypeId = payloads.TransactionTypeId
		RequestBodyLocationStock.StockTransactionReasonId = payloads.TransactionReasonId
		RequestBodyLocationStock.UpdatedDate = &payloads.UpdatedDate
		RequestBodyLocationStock.UpdatedByUserId = payloads.UpdatedByUserId

		//{
		//	"company_id": 4,
		//	"period_year": "2013",
		//	"period_month": "02",
		//	"warehouse_id": 1,
		//	"location_id": 1,
		//	"item_id": 293773,
		//	"warehouse_group_id": 1,
		//	"quantity_begin": 1,
		//	"quantity_sales": 1,
		//	"quantity_sales_return": 1,
		//	"quantity_purchase": 100,
		//	"quantity_purchase_return": 1,
		//	"quantity_transfer_in": 1,
		//	"quantity_transfer_out": 1,
		//	"quantity_claim_in": 1,
		//	"quantity_claim_out": 1,
		//	"quantity_adjustment": 5,
		//	"quantity_allocated": 2,
		//	"quantity_in_transit": 2,
		//	"quantity_ending": 5,
		//	"quantity_robbing_in": 3,
		//	"quantity_robbing_out": 2,
		//	"quantity_assembly_in": 5,
		//	"quantity_assembly_out": 5,
		//	"stock_transaction_type_id": 6,
		//	"stock_transaction_reason_id": 3,
		//	"created_by_user_id": 1,
		//	"created_date": "2023-10-31T12:00:00Z",
		//	"updated_by_user_id": 1,
		//	"updated_date": "2023-10-31T12:00:00Z"
		//}

		//var responseApi transactionsparepartpayloads.LocationUpdateResponse
		urlLocationStock := config.EnvConfigs.AfterSalesServiceUrl + "location-stock"
		errCrossService := utils.Put(urlLocationStock, &RequestBodyLocationStock, nil)
		if errCrossService != nil {
			fmt.Println("cross service pertama gagal")
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    errCrossService.Error(),
			}
		}
		//if responseApi.StatusCode != 200 {
		//	fmt.Println("cross service pertama gagal failed false" + responseApi.Message)
		//
		//	return false, &exceptions.BaseErrorResponse{
		//		StatusCode: http.StatusBadRequest,
		//		Err:        errors.New(fmt.Sprintf("failed to hit location update error with status code %d and message %s", responseApi.StatusCode, responseApi.Message)),
		//	}
		//}
	}
	//SELECT @Bin_Sys_No = ISNULL(BIN_SYS_NO,0) , @Bin_Line_No = ISNULL(BIN_LINE_NO,0)  FROM atItemGRPO1 WHERE GRPO_SYS_NO = @Ref_Sys_No AND ITEM_CODE = @Ref_Item_Code
	//getting binning for type PU
	//
	var goodsReceiveDetailEntities transactionsparepartentities.GoodsReceiveDetail
	err = db.Model(&goodsReceiveDetailEntities).First(&goodsReceiveDetailEntities).
		Where(transactionsparepartentities.GoodsReceiveDetail{GoodsReceiveSystemNumber: payloads.ReferenceId, ItemId: payloads.ReferenceItemId}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("reference system number is not found on binning list detail please check input"),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("failed to get goods receive detail on error : %s", err.Error()),
		}
	}
	OriginalItemId := payloads.ReferenceItemId
	if goodsReceiveDetailEntities.BinningDetailId != 0 {
		err = db.Model(&transactionsparepartentities.BinningStockDetail{}).
			Where(transactionsparepartentities.BinningStockDetail{BinningDetailSystemNumber: goodsReceiveDetailEntities.BinningDetailId}).
			Select("original_item_id").Scan(&OriginalItemId).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to get goods receive detail entities on : " + err.Error(),
			}
		}
	}
	if OriginalItemId == payloads.ReferenceItemId {
		//IF @Trans_Reason_Code = @Reason_NL
		//BEGIN
		//EXEC [dbo].[uspg_amItemCycle_Insert]
		//@Option			= 0,
		//@Company_Code	= @Company_Code,
		//@Period_Year	= @Period_Year,
		//@Period_Month	= @Period_Month,
		//@Item_Code		= @Ref_Item_Code,
		//@Order_Cycle	= 0,
		//@Qty_On_Order	= @Ref_Qty_Negatif ,
		//@Qty_Back_Order = 0
		//END
		if stockTransactionReason.StockTransactionReasonCode == "NL" {
			ItemCyclePayloads := masterpayloads.ItemCycleInsertPayloads{
				CompanyId:         payloads.CompanyId,
				PeriodYear:        periodYear,
				PeriodMonth:       periodMonth,
				ItemId:            OriginalItemId,
				OrderCycle:        0,
				QuantityOnOrder:   referenceQuantityNegative,
				QuantityBackOrder: 0,
			}
			var responseApi transactionsparepartpayloads.LocationUpdateResponse
			ItemCycleUrl := config.EnvConfigs.AfterSalesServiceUrl + "item-cycle"
			errCrossService := utils.Post(ItemCycleUrl, &ItemCyclePayloads, &responseApi)
			if errCrossService != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    errCrossService.Error(),
				}
			}
		}
		if stockTransactionReason.StockTransactionReasonCode != "WP" {
			
		}

	}

	return false, nil
}
