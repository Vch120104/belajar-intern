package masterwarehouserepository

import (
	masterentities "after-sales/api/entities/master"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type LocationStockRepositoryImpl struct {
}

func NewLocationStockRepositoryImpl() masterrepository.LocationStockRepository {
	return &LocationStockRepositoryImpl{}
}

func (repo *LocationStockRepositoryImpl) GetAllStock(db *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	//var
	var response []masterwarehousepayloads.LocationStockDBResponse
	entities := masterentities.LocationStock{}
	Jointable := db.Table("mtr_location_stock a").Select("a.company_id," +
		"a.period_year," +
		"a.period_month," +
		"a.warehouse_id," +
		"a.location_id," +
		"a.item_id," +
		"a.warehouse_group," +
		"a.quantity_begin," +
		"a.quantity_sales," +
		"a.quantity_sales_return," +
		"a.quantity_purchase," +
		"a.quantity_purchase_return," +
		"a.quantity_transfer_in," +
		"a.quantity_transfer_out," +
		"a.quantity_claim_in," +
		"a.quantity_claim_out," +
		"a.quantity_robbing_in," +
		"a.quantity_robbing_out," +
		"a.quantity_adjustment," +
		"a.quantity_allocated," +
		"a.quantity_in_transit," +
		"a.quantity_ending," +
		"b.warehouse_costing_type_id," +
		"b.brand_id," +
		"(ISNULL(a.quantity_begin,0) + ISNULL(a.quantity_purchase,0)-ISNULL(a.quantity_purchase_return,0)" +
		"+ ISNULL(A.quantity_transfer_in, 0) + ISNULL(A.quantity_claim_in, 0) + ISNULL(A.quantity_robbing_in, 0) +" +
		"ISNULL(A.quantity_adjustment, 0) + ISNULL(A.quantity_sales_return, 0)" +
		"+ ISNULL(A.quantity_assembly_in, 0)) - (ISNULL(A.quantity_sales, 0) + ISNULL(A.quantity_transfer_out, 0)" +
		"+ ISNULL(A.quantity_claim_in, 0) + ISNULL(A.quantity_robbing_out, 0) + ISNULL(A.quantity_assembly_out, 0))" +
		"  AS quantity_on_hand," +
		" (ISNULL(A.quantity_begin, 0) + ISNULL(A.quantity_purchase, 0) - ISNULL(A.quantity_purchase_return, 0) +" +
		"ISNULL(A.quantity_transfer_in, 0) + ISNULL(A.quantity_robbing_in, 0) + ISNULL(A.quantity_adjustment, 0)" +
		" + ISNULL(A.quantity_sales_return, 0) + ISNULL(A.quantity_assembly_in, 0)) - (ISNULL(A.quantity_sales, 0) " +
		" + ISNULL(A.quantity_transfer_out, 0) + ISNULL(A.quantity_robbing_out, 0) + ISNULL(A.quantity_assembly_out, 0) + " +
		"ISNULL(A.quantity_allocated, 0))" +
		"AS quantity_available").Joins("left outer join mtr_warehouse_master b ON a.company_id = b.company_id AND a.warehouse_id = b.warehouse_id")
	whereQuaery := utils.ApplyFilter(Jointable, filter)
	err := whereQuaery.Scopes(pagination.Paginate(&entities, &pages, whereQuaery)).Scan(&response).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch",
			Data:       nil,
			Err:        errors.New("Failed to fetch"),
		}
	}
	pages.Rows = response
	//page := pagination.Pagination{}

	return pages, nil
}
func periodMonthConcatinatedString(periodMonth string) string {
	if len(periodMonth) == 1 {
		return "0" + periodMonth
	}
	return periodMonth
}
func FetchMaxYear(db *gorm.DB, payloads masterwarehousepayloads.LocationStockUpdatePayloads) (int, error) {
	var lastPeriodYear int
	err := db.Model(&masterentities.LocationStock{}).
		Select("MAX(CAST(period_year AS INTEGER))").
		Where(masterentities.LocationStock{
			CompanyId:  payloads.CompanyId,
			LocationId: payloads.LocationId,
			ItemId:     payloads.ItemId,
		}).
		Scan(&lastPeriodYear).Error
	return lastPeriodYear, err
}

func FetchMaxMonth(db *gorm.DB, lastPeriodYear int, payloads masterwarehousepayloads.LocationStockUpdatePayloads) (int, error) {
	var lastPeriodMonth int
	err := db.Model(&masterentities.LocationStock{}).
		Select("MAX(CAST(period_month AS INTEGER))").
		Where(masterentities.LocationStock{
			PeriodYear: strconv.Itoa(lastPeriodYear),
			CompanyId:  payloads.CompanyId,
			LocationId: payloads.LocationId,
			ItemId:     payloads.ItemId,
		}).
		Scan(&lastPeriodMonth).Error
	return lastPeriodMonth, err
}

func FetchLocationStockData(db *gorm.DB, payloads masterwarehousepayloads.LocationStockUpdatePayloads, periodMonth, periodYear string) error {
	return db.Model(&masterentities.LocationStock{}).
		Where(masterentities.LocationStock{
			PeriodMonth: periodMonth,
			PeriodYear:  periodYear,
			CompanyId:   payloads.CompanyId,
			LocationId:  payloads.LocationId,
			ItemId:      payloads.ItemId,
		}).
		Select("quantity_begin, quantity_allocated").
		Row().
		Scan(&payloads.QuantityBegin, &payloads.QuantityAllocated)
}

func CreateNewLocationStockRecord(db *gorm.DB, payloads masterwarehousepayloads.LocationStockUpdatePayloads) error {
	newLocationStock := masterentities.LocationStock{
		CompanyId:              payloads.CompanyId,
		PeriodMonth:            payloads.PeriodMonth,
		PeriodYear:             payloads.PeriodYear,
		ItemId:                 payloads.ItemId,
		WarehouseId:            payloads.WarehouseId,
		WarehouseGroupId:       payloads.WarehouseGroupId,
		LocationId:             payloads.LocationId,
		QuantityBegin:          payloads.QuantityBegin,
		QuantityAdjustment:     0,
		QuantityAllocated:      payloads.QuantityAllocated,
		QuantityPurchase:       0,
		QuantityPurchaseReturn: 0,
		QuantitySales:          0,
		QuantitySalesReturn:    0,
		QuantityClaimOut:       0,
		QuantityClaimIn:        0,
		QuantityTransferOut:    0,
		QuantityTransferIn:     0,
		QuantityInTransit:      0,
		QuantityEnding:         0,
		ChangeNo:               0,
		CreatedByUserId:        payloads.CreatedByUserId,
		CreatedDate:            payloads.CreatedDate,
		UpdatedByUserId:        payloads.UpdatedByUserId,
		UpdatedDate:            payloads.UpdatedDate,
		QuantityRobbingIn:      0,
		QuantityRobbingOut:     0,
		QuantityAssemblyOut:    0,
		QuantityAssemblyIn:     0,
	}

	return db.Create(&newLocationStock).Error
}

// ALTER PROCEDURE [dbo].[uspg_amLocationStock_Update]
func (repo *LocationStockRepositoryImpl) UpdateLocationStock(db *gorm.DB, payloads masterwarehousepayloads.LocationStockUpdatePayloads) (bool, *exceptions.BaseErrorResponse) {
	//SELECT @Negative_Stock = NEGATIVE_STOCK FROM gmLoc1 WHERE COMPANY_CODE = @Company_Code AND WAREHOUSE_CODE = @Whs_Code
	var WarehouseMasterEntity masterwarehouseentities.WarehouseMaster

	//TS(SELECT TOP 1 1
	//FROM amLocationStock WITH(UPDLOCK)
	//WHERE COMPANY_CODE = @Company_Code
	//AND WHS_CODE = @Whs_Code
	//AND LOC_CODE = @Loc_Code
	//AND ITEM_CODE = @Item_Code
	//AND PERIOD_YEAR = @Period_Year
	//AND PERIOD_MONTH = @Period_Month)
	//BEGIN

	//begin refactor
	var LocationStockEntity masterentities.LocationStock
	err := db.Model(&masterentities.LocationStock{}).Where(masterentities.LocationStock{
		CompanyId:   payloads.CompanyId,
		WarehouseId: payloads.WarehouseId,
		LocationId:  payloads.LocationId,
		ItemId:      payloads.ItemId,
		PeriodMonth: payloads.PeriodMonth,
		PeriodYear:  payloads.PeriodYear,
	}).First(&LocationStockEntity).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch Location Stock Data",
		}
	}
	//if location stock belum ada dari period diatas
	if errors.Is(err, gorm.ErrRecordNotFound) {
		currentPeriodString := fmt.Sprintf("%s-%s-%s", payloads.PeriodYear, payloads.PeriodMonth, "01")
		ParsePeriodDatetime, errParse := time.Parse("2006-01-02", currentPeriodString)
		if errParse != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to parse current period",
			}
		}
		LastParsePeriodDatetime := ParsePeriodDatetime.AddDate(0, -1, 0)
		LastPeriodMonth := int(LastParsePeriodDatetime.Month())
		LastPeriodYear := LastParsePeriodDatetime.Year()
		//DECLARE @Qty_Allocated_LastPeriod numeric(10,2) = 0
		//SET @Last_Period =  (DATEADD(month, -1, CAST(@Period_Year + @Period_Month + '01' As DateTime)))
		//SELECT	@Qty_Begin	= ISNULL(QTY_ENDING,0),
		//@Qty_Allocated_LastPeriod = ISNULL(QTY_ALLOCATED,0)
		//FROM	amLocationStock
		//WHERE	COMPANY_CODE = @Company_Code AND WHS_CODE = @Whs_Code AND LOC_CODE = @Loc_Code AND  ITEM_CODE = @Item_Code AND
		//PERIOD_YEAR = CAST(YEAR(@Last_Period) As VarChar(4)) AND PERIOD_MONTH = (SELECT RIGHT('0' + CAST(MONTH(@Last_Period) As VarChar(2)), 2))
		err = db.Model(&LocationStockEntity).Where(masterentities.LocationStock{
			PeriodMonth: periodMonthConcatinatedString(strconv.Itoa(LastPeriodMonth)),
			PeriodYear:  strconv.Itoa(LastPeriodYear),
			CompanyId:   payloads.CompanyId,
			LocationId:  payloads.LocationId,
			ItemId:      payloads.ItemId,
		}).Select("quantity_begin,quantity_allocated").Row().Scan(&payloads.QuantityBegin, &payloads.QuantityAllocated)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				err = db.Model(&LocationStockEntity).Select("COALESCE(MAX(CAST(period_year AS INTEGER)),?)", LastPeriodYear).
					Where(masterentities.LocationStock{
						CompanyId:  payloads.CompanyId,
						LocationId: payloads.LocationId,
						ItemId:     payloads.ItemId,
					}).
					Scan(&LastPeriodYear).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "failed to fetch LastPeriodYear",
					}
				}
				err = db.Model(&LocationStockEntity).Select("COALESCE(MAX(CAST(period_month AS INTEGER)),?)", LastPeriodMonth).
					Where(masterentities.LocationStock{
						PeriodYear: strconv.Itoa(LastPeriodYear),
						CompanyId:  payloads.CompanyId,
						LocationId: payloads.LocationId,
						ItemId:     payloads.ItemId,
					}).
					Scan(&LastPeriodMonth).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "failed to fetch Last Period Month",
					}
				}
				err = db.Model(&LocationStockEntity).Where(masterentities.LocationStock{
					PeriodMonth: periodMonthConcatinatedString(strconv.Itoa(LastPeriodMonth)),
					PeriodYear:  strconv.Itoa(LastPeriodYear),
					CompanyId:   payloads.CompanyId,
					LocationId:  payloads.LocationId,
					ItemId:      payloads.ItemId,
				}).Select("quantity_begin,quantity_allocated").Row().Scan(&payloads.QuantityBegin, &payloads.QuantityAllocated)
				if err != nil && err.Error() != "sql: no rows in result set" {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "failed to fetch Last Location Stock Data",
					}
				}
			} else {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to fetch Last Location Stock Data",
				}
			}
		}
		var NewLocationStockEntity masterentities.LocationStock
		NewLocationStockEntity.CompanyId = payloads.CompanyId
		NewLocationStockEntity.PeriodMonth = payloads.PeriodMonth
		NewLocationStockEntity.PeriodYear = payloads.PeriodYear
		NewLocationStockEntity.ItemId = payloads.ItemId
		NewLocationStockEntity.WarehouseId = payloads.WarehouseId
		NewLocationStockEntity.WarehouseGroupId = payloads.WarehouseGroupId
		NewLocationStockEntity.LocationId = payloads.LocationId
		NewLocationStockEntity.QuantityBegin = payloads.QuantityBegin
		NewLocationStockEntity.QuantityAdjustment = 0
		NewLocationStockEntity.QuantityAllocated = payloads.QuantityAllocated
		NewLocationStockEntity.QuantityPurchase = 0
		NewLocationStockEntity.QuantityPurchaseReturn = 0
		NewLocationStockEntity.QuantitySales = 0
		NewLocationStockEntity.QuantitySalesReturn = 0
		NewLocationStockEntity.QuantityClaimOut = 0
		NewLocationStockEntity.QuantityClaimIn = 0
		NewLocationStockEntity.QuantityTransferOut = 0
		NewLocationStockEntity.QuantityTransferIn = 0
		NewLocationStockEntity.QuantityInTransit = 0
		NewLocationStockEntity.QuantityEnding = 0
		NewLocationStockEntity.ChangeNo = 0
		NewLocationStockEntity.CreatedByUserId = payloads.CreatedByUserId
		NewLocationStockEntity.CreatedDate = payloads.CreatedDate
		NewLocationStockEntity.UpdatedByUserId = payloads.UpdatedByUserId
		NewLocationStockEntity.UpdatedDate = payloads.UpdatedDate
		NewLocationStockEntity.QuantityRobbingIn = 0
		NewLocationStockEntity.QuantityRobbingOut = 0
		NewLocationStockEntity.QuantityAssemblyOut = 0
		NewLocationStockEntity.QuantityAssemblyIn = 0
		err = db.Create(&NewLocationStockEntity).Scan(&NewLocationStockEntity).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to create Location Stock Data",
			}
		}
		return true, nil
	}

	//endd refactor
	TransactionTypeEntity := masterentities.StockTransactionType{}
	err = db.Model(&TransactionTypeEntity).Where(masterentities.StockTransactionType{StockTransactionTypeId: payloads.StockTransactionTypeId}).
		First(&TransactionTypeEntity).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch Transaction Type Data",
		}
	}
	StockTransactionReason := masterentities.StockTransactionReason{}
	err = db.Model(&StockTransactionReason).Where(masterentities.StockTransactionReason{StockTransactionReasonId: payloads.StockTransactionReasonId}).
		First(&StockTransactionReason).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch Transaction Reason Data",
		}
	}
	var NegativeStock bool = false

	err = db.Model(&WarehouseMasterEntity).
		Where(masterwarehouseentities.WarehouseMaster{CompanyId: payloads.ItemId, WarehouseId: payloads.WarehouseId}).
		Select("warehouse_negative_stock").Scan(&NegativeStock).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch warehouse master data",
		}
	}
	var CurrentLocationStock masterentities.LocationStock
	err = db.Model(&CurrentLocationStock).Where(masterentities.LocationStock{
		LocationId:  payloads.LocationId,
		CompanyId:   payloads.CompanyId,
		ItemId:      payloads.ItemId,
		WarehouseId: payloads.WarehouseId,
		PeriodMonth: payloads.PeriodMonth,
		PeriodYear:  payloads.PeriodYear,
	}).First(&CurrentLocationStock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New(fmt.Sprintf("location stock is not found from this period year %s and period month %s", payloads.PeriodYear, payloads.PeriodMonth)),
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch Location Stock Data on " + err.Error(),
		}
	}
	//SET	QTY_ADJUSTMENT		= ISNULL(QTY_ADJUSTMENT,0) + ISNULL(@Qty_Adjustment,0) ,
	//QTY_ALLOCATED		= ISNULL(QTY_ALLOCATED,0) + ISNULL(@Qty_Allocated,0) ,
	//QTY_PURCHASE		= ISNULL(QTY_PURCHASE,0) + ISNULL(@Qty_Purchase,0)	,
	//QTY_PURCHASE_RETURN	= ISNULL(QTY_PURCHASE_RETURN,0) + ISNULL(@Qty_Purchase_Return,0)	,
	//QTY_SALES			= ISNULL(QTY_SALES,0) + ISNULL(@Qty_Sales,0)	,
	//QTY_SALES_RETURN	= ISNULL(QTY_SALES_RETURN,0) + ISNULL(@Qty_Sales_Return,0)	,
	//QTY_CLAIM_IN		= ISNULL(QTY_CLAIM_IN,0) + ISNULL(@Qty_Claim_In,0),
	//QTY_CLAIM_OUT		= ISNULL(QTY_CLAIM_OUT,0) + ISNULL(@Qty_Claim_Out,0) ,
	//QTY_TRANSFER_IN		= ISNULL(QTY_TRANSFER_IN,0) + ISNULL(@Qty_Transfer_In,0),
	//QTY_TRANSFER_OUT	= ISNULL(QTY_TRANSFER_OUT,0) + ISNULL(@Qty_Transfer_Out,0),
	//QTY_INTRANSIT		= ISNULL(QTY_INTRANSIT,0) + ISNULL(@Qty_Intransit,0),
	//--CHANGE_NO = CHANGE_NO + 1 ,
	//	CHANGE_USER_ID = @Change_User_Id ,
	//	CHANGE_DATETIME = @Change_Datetime,
	//	QTY_ROBBING_IN		= ISNULL(QTY_ROBBING_IN,0) + ISNULL(@Qty_Robbing_In,0),
	//QTY_ROBBING_OUT		= ISNULL(QTY_ROBBING_OUT,0) + ISNULL(@Qty_Robbing_Out,0),
	//QTY_ASSY_IN			= ISNULL(QTY_ASSY_IN,0) + ISNULL(@Qty_Assy_In,0),
	//QTY_ASSY_OUT		= ISNULL(QTY_ASSY_OUT,0) + ISNULL(@Qty_ASsy_Out,0)
	CurrentLocationStock.QuantityAdjustment += payloads.QuantityAdjustment
	CurrentLocationStock.QuantityAllocated += payloads.QuantityAllocated
	CurrentLocationStock.QuantityPurchase += payloads.QuantityPurchaseReturn
	CurrentLocationStock.QuantityPurchaseReturn += payloads.QuantityPurchaseReturn
	CurrentLocationStock.QuantitySales += payloads.QuantitySales
	CurrentLocationStock.QuantitySalesReturn += payloads.QuantitySalesReturn
	CurrentLocationStock.QuantityClaimIn += payloads.QuantityClaimIn
	CurrentLocationStock.QuantityClaimOut += payloads.QuantityClaimOut
	CurrentLocationStock.QuantityTransferIn += payloads.QuantityTransferIn
	CurrentLocationStock.QuantityTransferOut += payloads.QuantityTransferOut
	CurrentLocationStock.QuantityInTransit += payloads.QuantityInTransit
	CurrentLocationStock.ChangeNo += 1
	CurrentLocationStock.UpdatedDate = payloads.UpdatedDate
	CurrentLocationStock.UpdatedByUserId = payloads.UpdatedByUserId
	CurrentLocationStock.QuantityRobbingOut += payloads.QuantityRobbingOut
	CurrentLocationStock.QuantityRobbingIn += payloads.QuantityRobbingIn
	CurrentLocationStock.QuantityAssemblyIn += payloads.QuantityAssemblyIn
	CurrentLocationStock.QuantityAssemblyOut += payloads.QuantityAssemblyOut

	if (TransactionTypeEntity.StockTransactionTypeCode == "PU" &&
		(StockTransactionReason.StockTransactionReasonCode == "NL" ||
			StockTransactionReason.StockTransactionReasonCode == "BO" ||
			StockTransactionReason.StockTransactionReasonCode == "WP")) ||
		(TransactionTypeEntity.StockTransactionTypeCode == "AD" &&
			(StockTransactionReason.StockTransactionReasonCode == "NL" ||
				StockTransactionReason.StockTransactionReasonCode == "SF")) ||
		(TransactionTypeEntity.StockTransactionTypeCode == "TI" &&
			StockTransactionReason.StockTransactionReasonCode == "NL") ||
		(TransactionTypeEntity.StockTransactionTypeCode == "CI" &&
			StockTransactionReason.StockTransactionReasonCode == "AP") ||
		TransactionTypeEntity.StockTransactionTypeCode == "SR" ||
		TransactionTypeEntity.StockTransactionTypeCode == "RV" ||
		(TransactionTypeEntity.StockTransactionTypeCode == "SS" && NegativeStock == true) {
		//COMPANY_CODE = @Company_Code AND WHS_CODE = @Whs_Code AND
		//LOC_CODE = @Loc_Code AND ITEM_CODE = @Item_Code AND
		//PERIOD_YEAR = @Period_Year AND PERIOD_MONTH = @Period_Month
		//get current location stock from this month

		err = db.Save(&CurrentLocationStock).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error on saving location stock"}
		}
	} else {
		//SET @Qty_Ending = ISNULL((SELECT ISNULL(QTY_BEGIN,0)
		//+ ISNULL(QTY_PURCHASE,0)
		//- ISNULL(QTY_PURCHASE_RETURN,0)
		//- (ISNULL(QTY_SALES,0) - ISNULL(QTY_SALES_RETURN,0))
		//+ (ISNULL(QTY_TRANSFER_IN,0) - ISNULL(QTY_TRANSFER_OUT,0))
		//+ ISNULL(QTY_ADJUSTMENT,0)
		//+ (ISNULL(QTY_CLAIM_IN,0) - ISNULL(QTY_CLAIM_OUT,0))
		//+ (ISNULL(QTY_ROBBING_IN,0) - ISNULL(QTY_ROBBING_OUT,0))
		//+ (ISNULL(QTY_ASSY_IN,0) - ISNULL(QTY_ASSY_OUT,0))
		//FROM amLocationStock
		//WHERE	COMPANY_CODE = @Company_Code AND WHS_CODE = @Whs_Code AND
		//LOC_CODE = @Loc_Code AND ITEM_CODE = @Item_Code AND
		//PERIOD_YEAR = @Period_Year AND PERIOD_MONTH = @Period_Month	),0)
		var QuantityEnding float64
		QuantityEnding = CurrentLocationStock.QuantityBegin +
			CurrentLocationStock.QuantityPurchase -
			CurrentLocationStock.QuantityPurchaseReturn -
			(CurrentLocationStock.QuantitySales - CurrentLocationStock.QuantitySalesReturn) +
			(CurrentLocationStock.QuantityTransferIn - CurrentLocationStock.QuantityTransferOut) +
			CurrentLocationStock.QuantityAdjustment +
			(CurrentLocationStock.QuantityClaimIn - CurrentLocationStock.QuantityClaimOut) +
			(CurrentLocationStock.QuantityRobbingIn - CurrentLocationStock.QuantityRobbingOut) +
			(CurrentLocationStock.QuantityAssemblyIn - CurrentLocationStock.QuantityAssemblyOut)

		if !NegativeStock && QuantityEnding < 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("quantity ending cannot be negative"),
			}
		}
		CurrentLocationStock.QuantityEnding = QuantityEnding
		err = db.Save(&CurrentLocationStock).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "error on saving location stock"}
		}
	}
	return true, nil
}
