package transactionsparepartrepositoryimpl

import "C"
import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/crossservice/generalservice"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PurchaseOrderRepositoryImpl struct {
}

func NewPurchaseOrderRepositoryImpl() transactionsparepartrepository.PurchaseOrderRepository {
	return &PurchaseOrderRepositoryImpl{}
}

func (repo *PurchaseOrderRepositoryImpl) GetAllPurchaseOrder(db *gorm.DB, filter []utils.FilterCondition, page pagination.Pagination, DateParams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var payloadsresdb []transactionsparepartpayloads.GetAllDBResponses
	var Result []transactionsparepartpayloads.GetAllPurchaseOrderResponses
	entities := transactionsparepartentities.PurchaseOrderEntities{}
	var strfilter string
	if DateParams["PurchaseRequestDocNo"] == "" {
		strfilter = "1=1"
	} else {
		strfilter = DateParams["PurchaseRequestDocNo"]
		fmt.Println(strfilter)
	}
	JoinTable := db.Table("trx_item_purchase_order as A").
		Select("*").
		//Select("A.purchase_order_system_number,A.purchase_order_document_number,A.purchase_order_document_date,A.purchase_order_status_id,A.purchase_order_type_id,A.warehouse_id,A.supplier_id,C.purchase_request_document_number").
		Joins("LEFT JOIN trx_item_purchase_order_detail B ON A.purchase_order_system_number = B.purchase_order_system_number LEFT JOIN trx_purchase_request C ON B.purchase_request_system_number = C.purchase_request_system_number").
		Where(strfilter)
	whereQuery := utils.ApplyFilter(JoinTable, filter)
	var strDateFilter string
	if DateParams["purchase_order_date_from"] == "" {
		DateParams["purchase_order_date_from"] = "19000101"
	}
	if DateParams["purchase_order_date_to"] == "" {
		DateParams["purchase_order_date_to"] = "99991212"
	}
	strDateFilter = "purchase_order_document_date >='" + DateParams["purchase_order_date_from"] + "' AND purchase_order_document_date <= '" + DateParams["purchase_order_date_to"] + "'"

	err := whereQuery.Scopes(pagination.Paginate(&entities, &page, JoinTable)).Order("A.purchase_order_document_date desc").Where(strDateFilter).Scan(&payloadsresdb).Error
	if err != nil {
		return page, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	if len(payloadsresdb) == 0 {
		return page, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for _, i := range payloadsresdb {
		var purchaseRequestStatusDesc transactionsparepartpayloads.PurchaseRequestStatusResponse
		StatusURL := config.EnvConfigs.GeneralServiceUrl + "document-status/" + strconv.Itoa(i.PurchaseOrderStatusId)
		if err := utils.Get(StatusURL, &purchaseRequestStatusDesc, nil); err != nil {
			return page, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Status data from external service",
				Err:        err,
			}
		}
		var OrderType transactionsparepartpayloads.OrderTypeStatusResponse
		OrderTypeUrl := config.EnvConfigs.GeneralServiceUrl + "order-type/" + strconv.Itoa(i.OrderTypeId)
		if err := utils.Get(OrderTypeUrl, &OrderType, nil); err != nil {
			return page, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Status data from external service",
				Err:        err,
			}
		}
		var WhsEntities masterwarehouseentities.WarehouseMaster
		err = db.Model(&WhsEntities).Where(masterwarehouseentities.WarehouseMaster{WarehouseId: i.WarehouseId}).First(&WhsEntities).Error
		if err != nil {
			return page, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch warehouse data from database",
				Err:        err,
			}
		}
		var SupplierResponse transactionsparepartpayloads.SupplierResponsesAPI
		SupplierByIdUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(i.SupplierId)
		if err := utils.Get(SupplierByIdUrl, &SupplierResponse, nil); err != nil {
			return page, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Status data from external service",
				Err:        err,
			}
		}
		var prEntities transactionsparepartentities.PurchaseRequestEntities
		err = db.Model(&prEntities).Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: i.PurchaseRequestSystemNumber}).First(&prEntities).Error
		if err != nil {
			return page, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Purchase request data from database",
				Err:        err,
			}
		}
		tempres := transactionsparepartpayloads.GetAllPurchaseOrderResponses{
			PurchaseOrderSystemNumber:   i.PurchaseOrderSystemNumber,
			PurchaseOrderDocumentNumber: i.PurchaseOrderDocumentNumber,
			PurchaseOrderDocumentDate:   i.PurchaseOrderDocumentDate,
			PurchaseOrderStatus:         purchaseRequestStatusDesc.PurchaseRequestStatusDescription,
			OrderType:                   OrderType.OrderTypeName,
			WarehouseName:               WhsEntities.WarehouseName,
			SupplierName:                SupplierResponse.SupplierName,
			PurchaseRequestDocNo:        prEntities.PurchaseRequestDocumentNumber,
		}
		Result = append(Result, tempres)
	}
	page.Rows = Result
	return page, nil
}
func (repo *PurchaseOrderRepositoryImpl) GetByIdPurchaseOrder(db *gorm.DB, id int) (transactionsparepartpayloads.PurchaseOrderGetByIdResponses, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.PurchaseOrderEntities
	response := transactionsparepartpayloads.PurchaseOrderGetByIdResponses{}
	//response := transactionsparepartpayloads.PurchaseOrderGetByIdResponses{}
	err := db.Model(&entities).Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: id}).First(&response).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Purchase Order Not Found",
			Err:        errors.New("data not found"),
		}
	}
	return response, nil
}
func (repo *PurchaseOrderRepositoryImpl) GetAllDetailByHeaderId(db *gorm.DB, i int, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var payloadsresdb []transactionsparepartentities.PurchaseOrderDetailEntities
	var Result []transactionsparepartpayloads.PurchaseOrderGetDetail
	entities := transactionsparepartentities.PurchaseOrderDetailEntities{}
	JoinTable := db.Model(&entities).Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderSystemNumber: i})

	err := JoinTable.Scopes(pagination.Paginate(&entities, &page, JoinTable)).Scan(&payloadsresdb).Error
	if err != nil {
		return page, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	for _, data := range payloadsresdb {
		var ItemEntities masteritementities.Item
		err = db.Model(ItemEntities).Where(masteritementities.Item{ItemId: data.ItemId}).First(&ItemEntities).Error
		if err != nil {
			return page, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New("failed to get item name"),
			}
		}
		var PurchaseRequestEntities transactionsparepartentities.PurchaseRequestEntities
		err = db.Model(PurchaseRequestEntities).Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: data.PurchaseRequestSystemNumber}).First(&PurchaseRequestEntities).Error
		if err != nil {
			return page, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Purchase Order Not Found",
				Err:        errors.New("failed to get item name"),
			}
		}
		ResultPerId := transactionsparepartpayloads.PurchaseOrderGetDetail{
			PurchaseOrderDetailSystemNumber: data.PurchaseOrderDetailSystemNumber,
			Snp:                             data.Snp,
			ItemDiscountAmount:              data.ItemDiscountAmount,
			ItemPrice:                       data.ItemPrice,
			ItemQuantity:                    data.ItemQuantity,
			ItemUnitOfMeasurement:           data.ItemUnitOfMeasurement,
			UnitOfMeasurementRate:           data.UnitOfMeasurementRate,
			ItemCode:                        ItemEntities.ItemCode,
			ItemName:                        ItemEntities.ItemName,
			PurchaseOrderSystemNumber:       data.PurchaseOrderSystemNumber,
			PurchaseOrderLineNumber:         data.PurchaseOrderLineNumber,
			ItemTotal:                       data.ItemTotal,
			PurchaseRequestSystemNumber:     data.PurchaseRequestSystemNumber,
			PurchaseRequestLineNumber:       data.PurchaseRequestLineNumber,
			PurchaseRequestDocumentNumber:   PurchaseRequestEntities.PurchaseRequestDocumentNumber,
			StockOnHand:                     data.StockOnHand,
			ItemRemark:                      data.ItemRemark,
		}
		Result = append(Result, ResultPerId)
	}
	page.Rows = Result
	return page, nil
}

func (repo *PurchaseOrderRepositoryImpl) NewPurchaseOrderHeader(db *gorm.DB, request transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderResponses) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse) {
	var res transactionsparepartentities.PurchaseOrderEntities
	if request.CompanyId == 0 {
		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Company Code is missing. Please try again",
			Err:        errors.New("company Code is missing. Please try again"),
		}
	}
	if request.BrandId == 0 {
		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Vehicle Brand is missing is missing. Please try again",
			Err:        errors.New("vehicle Brand is missing. Please try again"),
		}
	}
	if request.ItemGroupId == 0 {
		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Item Group is missing. Please try again",
			Err:        errors.New("item Group is missing. Please try again"),
		}
	}
	if request.SupplierId == 0 {
		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "supplier id is missing. Please try again",
			Err:        errors.New("supplier id is missing. Please try again"),
		}
	}
	if request.WarehouseGroupId == 0 {
		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "warehouse group id is missing. Please try again",
			Err:        errors.New("warehouse group id is missing. Please try again"),
		}
	}
	if request.WarehouseId == 0 {
		return res, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "warehouse id is missing. Please try again",
			Err:        errors.New("warehouse id is missing. Please try again"),
		}
	}

	entities := transactionsparepartentities.PurchaseOrderEntities{
		CompanyId:                 request.CompanyId,
		PurchaseOrderSystemNumber: 0,
		//PurchaseOrderSystemNumber:           0,
		PurchaseOrderDocumentNumber:         request.PurchaseOrderDocumentNumber,
		PurchaseOrderDocumentDate:           request.PurchaseOrderDocumentDate,
		PurchaseOrderStatusId:               request.PurchaseOrderStatusId,
		BrandId:                             request.BrandId,
		ItemGroupId:                         request.ItemGroupId,
		OrderTypeId:                         request.PurchaseOrderTypeId,
		SupplierId:                          request.SupplierId,
		SupplierPicId:                       request.SupplierPicId,
		WarehouseId:                         request.WarehouseId,
		WarehouseGroupId:                    request.WarehouseGroupId,
		CostCenterId:                        request.CostCenterId,
		ProfitType:                          "P",
		ProfitCenterId:                      request.ProfitCenterId,
		AffiliatedPurchaseOrder:             request.AffiliatedPurchaseOrder,
		CurrencyId:                          request.CurrencyId,
		BackOrder:                           request.BackOrder,
		SetOrder:                            request.SetOrder,
		ViaBinning:                          request.ViaBinning,
		VatCode:                             "",
		PphCode:                             "",
		TotalDiscount:                       request.TotalDiscount,
		TotalAmount:                         request.TotalAmount,
		TotalVat:                            request.TotalVat,
		TotalAfterVat:                       request.TotalAfterVat,
		LastTotalDiscount:                   nil,
		LastTotalAmount:                     nil,
		LastTotalVat:                        nil,
		LastTotalAfterVat:                   nil,
		TotalAmountConfirm:                  nil,
		PurchaseOrderRemark:                 request.PurchaseOrderRemark,
		DpRequest:                           request.DpRequest,
		DpPayment:                           nil,
		DpPaymentAllocated:                  nil,
		DpPaymentAllocatedInvoice:           nil,
		DpPaymentAllocatedPpn:               nil,
		DpPaymentAllocatedRequestForPayment: nil,
		DeliveryId:                          request.DeliveryId,
		ExpectedDeliveryDate:                request.ExpectedDeliveryDate,
		ExpectedArrivalDate:                 request.ExpectedArrivalDate,
		EstimatedDeliveryDate:               nil,
		EstimatedDeliveryTime:               "",
		SalesOrderSystemNumber:              0,
		SalesOrderDocumentNumber:            "",
		LastPrintById:                       0,
		ApprovalRequestById:                 0,
		ApprovalRequestNumber:               0,
		ApprovalRequestDate:                 nil,
		ApprovalRemark:                      "",
		ApprovalLastById:                    0,
		ApprovalLastDate:                    nil,
		TotalInvoiceDownPayment:             nil,
		TotalInvoiceDownPaymentVat:          nil,
		TotalInvoiceDownPaymentAfterVat:     nil,
		DownPaymentReturn:                   nil,
		JournalSystemNumber:                 0,
		EventNumber:                         "",
		ItemClassId:                         0,
		APMIsDirectShipment:                 "",
		DirectShipmentCustomerId:            0,
		ExternalPurchaseOrderNumber:         request.ExternalPurchaseOrderNumber,
		PurchaseOrderTypeId:                 request.PurchaseOrderTypeId,
		CurrencyExchangeRate:                nil,
		PurchaseOrderDetail:                 nil,
		CreatedByUserId:                     request.CreatedByUserId,
		CreatedDate:                         request.CreatedDate,
		UpdatedByUserId:                     request.UpdatedByUserId,
		UpdatedDate:                         request.UpdatedDate,
		ChangeNo:                            1,
	}
	err := db.Create(&entities).Scan(&entities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			//Data:       err,
			Err: err,
		}
	}
	return entities, nil
}

func (repo *PurchaseOrderRepositoryImpl) UpdatePurchaseOrderHeader(db *gorm.DB, id int, payloads transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderPayloads) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse) {
	var taxRate float64
	var totalDiscount float64
	var totalAmount float64
	var pkptype bool
	var CompanyVatPkpNo string
	var SupplierVatPkpNo string
	var SupplierId int
	var totalVat float64
	var EntitiesPurchaseOrder transactionsparepartentities.PurchaseOrderEntities
	currentTime := time.Now().UTC()
	timeString := currentTime.Format("2006-01-02T15:04:05.000Z")
	var TaxRateResponse transactionsparepartpayloads.TaxRateResponseApi
	//		SET @TAX_RATE = dbo.getTaxPercent(dbo.getVariableValue('TAX_TYPE_PPN'),dbo.getVariableValue('TAX_SERV_CODE_PPN'),@Change_Datetime)
	TaxRateUrl := config.EnvConfigs.FinanceServiceUrl + "tax-fare/detail/tax-percent?tax_service_code=PPN&tax_type_code=PPN&effective_date=" + timeString
	if err := utils.Get(TaxRateUrl, &TaxRateResponse, nil); err != nil {
		fmt.Println(err.Error())
		return EntitiesPurchaseOrder, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Status data from external service",
			Err:        err,
		}
	}
	if TaxRateResponse.TaxPercent != nil {
		taxRate = *TaxRateResponse.TaxPercent

	} else {
		taxRate = 0
	}
	//taxRate = *TaxRateResponse.TaxPercent
	//SET @TOTAL_DISCOUNT = (SELECT SUM(ITEM_QTY * ITEM_DISC_AMOUNT) FROM atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No)

	err := db.Table("trx_item_purchase_order_detail A").Select("ISNULL(SUM(item_quantity * item_discount_amount),0)").
		Where("A.purchase_order_system_number =?", id).Scan(&totalDiscount).Error
	if err != nil {
		return EntitiesPurchaseOrder, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Total Discount",
			Err:        err,
		}
	}
	//SET @TOTAL_AMOUNT = (SELECT SUM(ITEM_QTY * ITEM_PRICE) FROM atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No)
	err = db.Table("trx_item_purchase_order_detail A").Select("isnull(SUM(item_quantity * item_price),0)").
		Where("A.purchase_order_system_number =?", id).Scan(&totalAmount).Error
	if err != nil {
		return EntitiesPurchaseOrder, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Total Amount",
			Err:        err,
		}
	}
	//SET @SUPPLIER = (SELECT SUPPLIER_CODE FROM atItemPO0 WHERE PO_SYS_NO = @Po_Sys_No)
	err = db.Model(&transactionsparepartentities.PurchaseOrderEntities{}).Select("supplier_id").
		Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: id}).
		Scan(&SupplierId).Error
	if err != nil {
		return EntitiesPurchaseOrder, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Get Supplier ID from purchase order",
			Err:        err,
		}
	}
	//SET @Supplier_Vat_Pkp_No = (SELECT ISNULL(VAT_PKP_NO, '') FROM gmSupplier0 WHERE SUPPLIER_CODE = @SUPPLIER)
	var SupplierResponse transactionsparepartpayloads.SupplierResponsesAPI
	SupplierByIdUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(SupplierId)
	if err := utils.Get(SupplierByIdUrl, &SupplierResponse, nil); err != nil {
		fmt.Println(err.Error())
		return EntitiesPurchaseOrder, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed to fetch Supplier Id from external service",
			Err:        err,
		}
	}
	pkptype = SupplierResponse.TaxSupplier.PkpType
	SupplierVatPkpNo = SupplierResponse.TaxSupplier.PkpNo
	//SET @Company_Vat_Pkp_No = (SELECT ISNULL(VAT_PKP_NO, '') FROM gmComp0 WHERE COMPANY_CODE = @Company_Code)
	var CompanyDetailResponse transactionsparepartpayloads.CompanyDetailResponses
	CompanyDetailUrl := config.EnvConfigs.GeneralServiceUrl + "company-detail/" + strconv.Itoa(payloads.CompanyId)
	if err := utils.Get(CompanyDetailUrl, &CompanyDetailResponse, nil); err != nil {
		return EntitiesPurchaseOrder, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed to fetch Company Id Detail from external service",
			Err:        err,
		}
	}
	//IF @PKP='Y'
	//BEGIN
	//IF LTRIM(RTRIM(ISNULL(@Supplier_Vat_Pkp_No,'')))<>LTRIM(RTRIM(ISNULL(@Company_Vat_Pkp_No,'')))
	//BEGIN
	//SET @TOTAL_VAT = (@TOTAL_AMOUNT - @TOTAL_DISCOUNT) * (@TAX_RATE / 100)
	//END
	//ELSE
	//BEGIN
	//SET @TOTAL_VAT = 0
	//END
	//END
	//ELSE
	//BEGIN
	//SET @TOTAL_VAT = 0
	//END

	if pkptype {
		if SupplierVatPkpNo != CompanyVatPkpNo {
			totalVat = (totalAmount - totalDiscount) * (taxRate / 100)
		} else {
			totalVat = 0
		}

	} else {
		totalVat = 0
	}
	//SET @TOTAL_AFTER_VAT = @TOTAL_AMOUNT - @TOTAL_DISCOUNT + @TOTAL_VAT

	var totalAfterVat = totalAmount - totalDiscount + totalVat
	err = db.Model(&EntitiesPurchaseOrder).Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: id}).
		Scan(&EntitiesPurchaseOrder).Error
	//process update
	//UPDATE atItemPO0 SET
	//SUPPLIER_CODE = @Supplier_Code ,
	//	SUPPLIER_PIC_CODE = @Supplier_Pic_Code ,
	//	WHS_GROUP = @Whs_Group ,
	//	WHS_CODE = @Whs_Code ,
	//	WHS_PIC_CODE = @Whs_Pic_Code ,
	//	COST_TYPE = @Cost_Type ,
	//	COST_CENTER = @Cost_Center ,
	//	PROFIT_TYPE = @Profit_Type ,
	//	PROFIT_CENTER = @Profit_Center ,
	//	AFFILIATED_PO = @Affiliated_Po ,
	//	CCY_CODE = @Ccy_Code ,
	//	VIA_BINNING = @Via_Binning ,
	//--VAT_CODE = @Vat_Code ,
	//	PO_REMARK = @Po_Remark ,
	//	DP_REQUEST = @Dp_Request ,
	//	DELIVERY_VIA = @Delivery_Via ,
	//	EXP_ARRIVAL_DATE = @Exp_Arrival_Date ,
	//	EXP_DELIVERY_DATE = @Exp_Delivery_Date ,
	//	TOTAL_VAT = @TOTAL_VAT ,
	//	TOTAL_AFTER_VAT = @TOTAL_AFTER_VAT ,
	//	IS_DIRECT_SHIPMENT = @Is_Direct_Shipment ,
	//	CUSTOMER_CODE = @Customer_Code ,
	//--CHANGE_NO = CHANGE_NO + 1 ,
	//	CHANGE_USER_ID = @Change_User_Id ,
	//	CHANGE_DATETIME = @Change_Datetime,
	//	EXTERNAL_PO_NO = @External_Po_No
	//WHERE PO_SYS_NO = @Po_Sys_No
	EntitiesPurchaseOrder.SupplierId = payloads.SupplierId
	EntitiesPurchaseOrder.SupplierPicId = payloads.SupplierPicId
	EntitiesPurchaseOrder.WarehouseId = payloads.WarehouseId
	EntitiesPurchaseOrder.WarehouseGroupId = payloads.WarehouseGroupId
	EntitiesPurchaseOrder.CostCenterId = payloads.CostCenterId
	EntitiesPurchaseOrder.ProfitCenterId = payloads.ProfitCenterId
	EntitiesPurchaseOrder.AffiliatedPurchaseOrder = payloads.AffiliatedPurchaseOrder
	EntitiesPurchaseOrder.CurrencyId = payloads.CurrencyId
	EntitiesPurchaseOrder.ViaBinning = payloads.ViaBinning
	EntitiesPurchaseOrder.PurchaseOrderRemark = payloads.PurchaseOrderRemark
	EntitiesPurchaseOrder.DpRequest = payloads.DpRequest
	EntitiesPurchaseOrder.DeliveryId = payloads.DeliveryId
	EntitiesPurchaseOrder.ExpectedArrivalDate = payloads.ExpectedArrivalDate
	EntitiesPurchaseOrder.ExpectedDeliveryDate = payloads.ExpectedDeliveryDate
	EntitiesPurchaseOrder.TotalVat = &totalVat
	EntitiesPurchaseOrder.TotalAfterVat = &totalAfterVat
	EntitiesPurchaseOrder.APMIsDirectShipment = payloads.APMIsDirectShipment
	EntitiesPurchaseOrder.DirectShipmentCustomerId = payloads.CustomerId
	EntitiesPurchaseOrder.ExternalPurchaseOrderNumber = payloads.ExternalPurchaseOrderNumber
	EntitiesPurchaseOrder.UpdatedByUserId = payloads.UpdatedByUserId
	EntitiesPurchaseOrder.UpdatedDate = payloads.UpdatedDate
	err = db.Model(&EntitiesPurchaseOrder).Save(&EntitiesPurchaseOrder).Error
	if err != nil {
		return EntitiesPurchaseOrder, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Err:        err,
		}
	}
	return EntitiesPurchaseOrder, nil
}
func (repo *PurchaseOrderRepositoryImpl) GetPurchaseOrderDetailById(db *gorm.DB, i int) (transactionsparepartpayloads.PurchaseOrderGetDetail, *exceptions.BaseErrorResponse) {
	//var payloadsresdb []transactionsparepartentities.PurchaseOrderDetailEntities
	entities := transactionsparepartentities.PurchaseOrderDetailEntities{}
	var ResultPerId transactionsparepartpayloads.PurchaseOrderGetDetail
	err := db.Model(&entities).Where(&transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderDetailSystemNumber: i}).First(&entities).Error

	var ItemEntities masteritementities.Item
	err = db.Model(ItemEntities).Where(masteritementities.Item{ItemId: entities.ItemId}).First(&ItemEntities).Error
	if err != nil {
		return ResultPerId, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("failed to get item name"),
		}
	}
	var PurchaseRequestEntities transactionsparepartentities.PurchaseRequestEntities
	err = db.Model(PurchaseRequestEntities).Where(transactionsparepartentities.PurchaseRequestEntities{PurchaseRequestSystemNumber: entities.PurchaseRequestSystemNumber}).First(&PurchaseRequestEntities).Error
	if err != nil {
		return ResultPerId, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Purchase Order Not Found",
			Err:        errors.New("failed to get item name"),
		}
	}
	ResultPerId = transactionsparepartpayloads.PurchaseOrderGetDetail{
		PurchaseOrderDetailSystemNumber: entities.PurchaseOrderDetailSystemNumber,
		Snp:                             entities.Snp,
		ItemDiscountAmount:              entities.ItemDiscountAmount,
		ItemPrice:                       entities.ItemPrice,
		ItemQuantity:                    entities.ItemQuantity,
		ItemUnitOfMeasurement:           entities.ItemUnitOfMeasurement,
		UnitOfMeasurementRate:           entities.UnitOfMeasurementRate,
		ItemCode:                        ItemEntities.ItemCode,
		ItemName:                        ItemEntities.ItemName,
		PurchaseOrderSystemNumber:       entities.PurchaseOrderSystemNumber,
		PurchaseOrderLineNumber:         entities.PurchaseOrderLineNumber,
		ItemTotal:                       entities.ItemTotal,
		PurchaseRequestSystemNumber:     entities.PurchaseRequestSystemNumber,
		PurchaseRequestLineNumber:       entities.PurchaseRequestLineNumber,
		PurchaseRequestDocumentNumber:   PurchaseRequestEntities.PurchaseRequestDocumentNumber,
		StockOnHand:                     entities.StockOnHand,
		ItemRemark:                      entities.ItemRemark,
	}
	return ResultPerId, nil
}

type PurchaseDatas struct {
	VehicleChassisNo        string
	PurchaseRequestQuantity float64
	PurchaseOrderQuantity   float64
}

func (repo *PurchaseOrderRepositoryImpl) NewPurchaseOrderDetail(db *gorm.DB, payloads transactionsparepartpayloads.PurchaseOrderDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse) {
	//SET @Po_Line = (SELECT ISNULL(MAX(PO_LINE), 0) + 1 FROM atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No)
	//if payloads.Snp == nil || *payloads.Snp > 1 {
	//	if payloads.purchaseorder {
	//
	//	}
	//}
	//	IF NOT EXISTS (SELECT * FROM atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No AND PO_LINE = @Po_Line)
	//replace with sysno
	count := 0
	var entities transactionsparepartentities.PurchaseOrderDetailEntities
	err := db.Model(&entities).
		Select("count(*)").
		Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseRequestDetailSystemNumber: payloads.PurchaseRequestDetailSystemNumber}).
		Scan(&count).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Data Is Already Exist",
			Err:        errors.New("data Is Already Exist"),
		}
	}
	if count == 0 {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Data is already exists",
			Err:        errors.New("data is already exists"),
		}
	}
	var VehicleId int
	var PurchaseRequestQuantity float64
	var PurchaseOrderQuantity float64
	//SELECT
	//@Vehicle_Chassis_No = VEHICLE_CHASSIS_NO
	//,@PR_QTY  = ITEM_QTY
	//,@PO_QTY = ISNULL(PO_QTY,0)-- For 1 PR , 1< PO
	//FROM atItemPR1 WHERE PR_SYS_NO = @Pr_Sys_No AND PR_LINE = @Pr_Line
	var prentities transactionsparepartentities.PurchaseRequestDetail
	err = db.Model(&prentities).Select("vehicle_id,item_quantity,item_quantity").
		Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseRequestDetailSystemNumber: payloads.PurchaseRequestDetailSystemNumber}).
		First(&prentities).
		Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Data Is Already Exist",
			Err:        errors.New("data Is Already Exist"),
		}
	}
	err = db.Table("trx_item_purchase_order_detail A").Select("ISNULL(SUM(A.item_quantity),0)").
		Joins("JOIN trx_item_purchase_order B ON A.purchase_order_system_number = b.purchase_order_system_number").
		Where("B.purchase_order_status_id NOT IN (?,?) AND A.purchase_request_system_number = 1", 80, 90).
		Scan(&PurchaseOrderQuantity).Error
	//dimana ga closed dan ga canceled
	VehicleId = prentities.VehicleId
	PurchaseRequestQuantity = *prentities.ItemQuantity

	//PurchaseOrderQuantity = *prentities.ItemPrice
	fmt.Println(VehicleId)
	//IF EXISTS (SELECT *
	//	FROM atItemPO1 A
	//INNER JOIN gmItem0 B ON A.ITEM_CODE=B.ITEM_CODE AND B.ITEM_GROUP IN (@Item_Group_OJ,@Item_Group_IN)
	//WHERE PO_SYS_NO = @Po_Sys_No AND A.ITEM_CODE  = @Item_Code)
	//BEGIN
	//RAISERROR('Cannot insert duplicate Item Code...',16,1)
	//RETURN 0
	//END
	var exists bool
	err = db.Model(&entities). // Assuming AtItemPO1 is your struct
					Select("1").
					Joins("INNER JOIN mtr_item B ON A.item_id = B.item_id AND B.item_group_id IN (?,?)", 15, 23). //in and oj
					Where("purchase_order_system_number = ? A.item_id = ?", payloads.PurchaseOrderSystemNumber, payloads.ItemId).
					Limit(1).
					Find(&exists).Error
	if exists {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Cannot insert duplicate Item Code...",
			Err:        err,
		}
	}
	//			SET @PO_QTY = @PO_QTY + @Item_Qty
	PurchaseOrderQuantity = PurchaseRequestQuantity + PurchaseOrderQuantity
	if PurchaseRequestQuantity < PurchaseOrderQuantity {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Total Quantity Purhase Order cannot exceed Quantity Purchase Request",
			Err:        err,
		}
	}
	//INSERT INTO atItemPO1 (
	//	RECORD_STATUS ,
	//	PO_SYS_NO ,
	//	PO_LINE ,
	//	ITEM_CODE ,
	//	ITEM_UOM ,
	//	UOM_RATE ,
	//	ITEM_QTY ,
	//	ITEM_PRICE ,
	//	ITEM_DISC_PERCENT ,
	//	ITEM_DISC_AMOUNT ,
	//	ITEM_TOTAL ,
	//	SUBSTITUTE_TYPE,
	//	PR_SYS_NO ,
	//	PR_LINE ,
	//	OLD_PO_SYS_NO ,
	//	OLD_PO_LINE ,
	//	VEHICLE_CHASSIS_NO,
	//	STOCK_ON_HAND,
	//	ITEM_REMARK,
	//	SNP,
	//	CHANGE_NO ,
	//	CREATION_USER_ID ,
	//	CREATION_DATETIME ,
	//	CHANGE_USER_ID ,
	//	CHANGE_DATETIME )
	//VALUES (
	//@Record_Status ,
	//@Po_Sys_No ,
	//@Po_Line ,
	//@Item_Code ,
	//@Item_Uom ,
	//@Uom_Rate ,
	//@Item_Qty ,
	//@Item_Price ,
	//@Item_Disc_Percent ,
	//@Item_Disc_Amount ,
	//@Item_Total ,
	//@Substitute_Type ,
	//@Pr_Sys_No ,
	//@Pr_Line ,
	//@Old_Po_Sys_No ,
	//@Old_Po_Line ,
	//@Vehicle_Chassis_No,
	//@Stock_On_Hand,
	//@Item_Remark,
	//@SNP,
	//	0 ,
	//@Creation_User_Id ,
	//@Creation_Datetime ,
	//@Change_User_Id ,
	//@Change_Datetime )

	entities = transactionsparepartentities.PurchaseOrderDetailEntities{
		//PurchaseOrderDetailSystemNumber:   0,
		PurchaseOrderSystemNumber: payloads.PurchaseOrderSystemNumber,
		PurchaseOrderLineNumber:   payloads.PurchaseOrderLine,
		ItemId:                    payloads.ItemId,
		ItemUnitOfMeasurement:     payloads.ItemUnitOfMeasurement,
		UnitOfMeasurementRate:     payloads.UnitOfMeasurementRate,
		ItemQuantity:              payloads.ItemQuantity,
		ItemPrice:                 payloads.ItemPrice,
		ItemDiscountPercentage:    payloads.ItemDiscountPercentage,
		ItemDiscountAmount:        payloads.ItemDiscountAmount,
		ItemTotal:                 payloads.ItemTotal,
		//SubstituteTypeId: 0
		PurchaseRequestSystemNumber:                payloads.PurchaseRequestSystemNumber,
		PurchaseRequestLineNumber:                  payloads.PurchaseRequestLineNumber,
		PurchaseRequestDetailSystemNumber:          payloads.PurchaseRequestDetailSystemNumber,
		ChangedItemPurchaseOrderLineNumber:         0,
		ChangedItemPurchaseOrderSystemNumber:       0,
		ChangedItemPurchaseOrderDetailSystemNumber: 0,
		StockOnHand:     payloads.StockOnHand,
		ItemRemark:      payloads.ItemRemark,
		Snp:             payloads.Snp,
		ChangeNo:        1,
		CreatedByUserId: payloads.CreatedByUserId,
		UpdatedDate:     payloads.UpdatedDate,
		CreatedDate:     payloads.CreatedDate,
		UpdatedByUserId: payloads.UpdatedByUserId,
		//VehicleChassisNumber: payloads.vehi
	}
	err = db.Create(&entities).Scan(&entities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Creating purchase request detail failed",
			Err:        err,
		}
	}
	//-=-==-=-=-=-=-=------=-=-=-=-=-=-=-=-=-=-
	var taxRate float64
	var totalDiscount float64
	var totalAmount float64
	var pkptype bool
	var CompanyVatPkpNo string
	var SupplierVatPkpNo string
	var SupplierId int
	var totalVat float64
	currentTime := time.Now().UTC()
	timeString := currentTime.Format("2006-01-02T15:04:05.000Z")
	var TaxRateResponse transactionsparepartpayloads.TaxRateResponseApi
	//		SET @TAX_RATE = dbo.getTaxPercent(dbo.getVariableValue('TAX_TYPE_PPN'),dbo.getVariableValue('TAX_SERV_CODE_PPN'),@Change_Datetime)
	TaxRateUrl := config.EnvConfigs.FinanceServiceUrl + "tax-fare/detail/tax-percent?tax_service_code=PPN&tax_type_code=PPN&effective_date=" + timeString
	if err := utils.Get(TaxRateUrl, &TaxRateResponse, nil); err != nil {
		fmt.Println(err.Error())
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Status data from external service",
			Err:        err,
		}
	}
	if TaxRateResponse.TaxPercent != nil {
		taxRate = *TaxRateResponse.TaxPercent

	} else {
		taxRate = 0
	}
	//SET @TOTAL_DISCOUNT = (SELECT SUM(ITEM_QTY * ITEM_DISC_AMOUNT) FROM atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No)

	err = db.Table("trx_item_purchase_order_detail A").Select("ISNULL(SUM(item_quantity * item_discount_amount),0)").
		Where("A.purchase_order_system_number =?", payloads.PurchaseOrderSystemNumber).Scan(&totalDiscount).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Total Discount",
			Err:        err,
		}
	}
	//SET @TOTAL_AMOUNT = (SELECT SUM(ITEM_QTY * ITEM_PRICE) FROM atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No)
	err = db.Table("trx_item_purchase_order_detail A").Select("isnull(SUM(item_quantity * item_price),0)").
		Where("A.purchase_order_system_number =?", payloads.PurchaseOrderSystemNumber).Scan(&totalAmount).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Total Amount",
			Err:        err,
		}
	}
	//			SET @SUPPLIER = (SELECT SUPPLIER_CODE FROM atItemPO0 WHERE PO_SYS_NO = @Po_Sys_No)
	err = db.Model(&transactionsparepartentities.PurchaseOrderEntities{}).Select("supplier_id").
		Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: payloads.PurchaseOrderSystemNumber}).
		Scan(&SupplierId).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Get Supplier ID from purchase order",
			Err:        err,
		}
	}
	//SET @Company_Code = (SELECT COMPANY_CODE FROM atItemPO0 WHERE PO_SYS_NO = @Po_Sys_No)
	//
	//var CompanyId int
	var poEntities transactionsparepartentities.PurchaseOrderEntities
	err = db.Model(&poEntities).Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: payloads.PurchaseOrderSystemNumber}).
		First(&poEntities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Header Not Found",
			Err:        err,
		}
	}
	//SET @Supplier_Vat_Pkp_No = (SELECT ISNULL(VAT_PKP_NO, '') FROM gmSupplier0 WHERE SUPPLIER_CODE = @SUPPLIER)
	var SupplierResponse transactionsparepartpayloads.SupplierResponsesAPI
	SupplierByIdUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(SupplierId)
	if err := utils.Get(SupplierByIdUrl, &SupplierResponse, nil); err != nil {
		fmt.Println(err.Error())
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed to fetch Supplier Id from external service",
			Err:        err,
		}
	}
	pkptype = SupplierResponse.TaxSupplier.PkpType
	SupplierVatPkpNo = SupplierResponse.TaxSupplier.PkpNo
	//SET @Company_Vat_Pkp_No = (SELECT ISNULL(VAT_PKP_NO, '') FROM gmComp0 WHERE COMPANY_CODE = @Company_Code)
	var CompanyDetailResponse transactionsparepartpayloads.CompanyDetailResponses
	CompanyDetailUrl := config.EnvConfigs.GeneralServiceUrl + "company-detail/" + strconv.Itoa(poEntities.CompanyId)
	if err := utils.Get(CompanyDetailUrl, &CompanyDetailResponse, nil); err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed to fetch Company Id Detail from external service",
			Err:        err,
		}
	}
	//IF @PKP='Y'
	//BEGIN
	//IF LTRIM(RTRIM(ISNULL(@Supplier_Vat_Pkp_No,'')))<>LTRIM(RTRIM(ISNULL(@Company_Vat_Pkp_No,'')))
	//BEGIN
	//SET @TOTAL_VAT = (@TOTAL_AMOUNT - @TOTAL_DISCOUNT) * (@TAX_RATE / 100)
	//END
	//ELSE
	//BEGIN
	//SET @TOTAL_VAT = 0
	//END
	//END
	//ELSE
	//BEGIN
	//SET @TOTAL_VAT = 0
	//END

	if pkptype {
		if SupplierVatPkpNo != CompanyVatPkpNo {
			totalVat = (totalAmount - totalDiscount) * (taxRate / 100)
		} else {
			totalVat = 0
		}

	} else {
		totalVat = 0
	}
	//SET @TOTAL_AFTER_VAT = @TOTAL_AMOUNT - @TOTAL_DISCOUNT + @TOTAL_VAT

	totalAfterVat := totalAmount - totalDiscount + totalVat
	//err = db.Model(&poEntities).Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: payloads.PurchaseOrderSystemNumber}).
	//	Scan(&poEntities).Error

	// AMBIL DP REQ DARI SUPPLIER

	//var SupplierResponses transactionsparepartpayloads.SupplierResponsesAPI
	//		SET @TAX_RATE = dbo.getTaxPercent(dbo.getVariableValue('TAX_TYPE_PPN'),dbo.getVariableValue('TAX_SERV_CODE_PPN'),@Change_Datetime)
	//SupplierResponseUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(poEntities.SupplierId)
	//if err := utils.Get(SupplierResponseUrl, &SupplierResponses, nil); err != nil {
	//	fmt.Println(err.Error())
	//	return entities, &exceptions.BaseErrorResponse{
	//		StatusCode: http.StatusInternalServerError,
	//		Message:    "Failed to fetch Status data from external service",
	//		Err:        err,
	//	}
	//}
	//dpRequest = SupplierResponse.MinimumDownPayment
	var dpRequest float64
	//SupplierResponse.MinimumDownPayment = new(float64)
	//*SupplierResponse.MinimumDownPayment = 2

	if SupplierResponse.MinimumDownPayment != nil && *SupplierResponse.MinimumDownPayment != 0 {
		dpRequest = totalAfterVat * (*SupplierResponse.MinimumDownPayment / 100)
	} else {
		dpRequest = *poEntities.DpRequest
	}
	*poEntities.DpRequest = dpRequest
	*poEntities.TotalDiscount = totalDiscount
	*poEntities.TotalAmount = totalAmount
	*poEntities.TotalVat = totalVat
	*poEntities.TotalAfterVat = totalAfterVat
	poEntities.ChangeNo += 1
	poEntities.UpdatedByUserId = payloads.UpdatedByUserId
	*poEntities.UpdatedDate = currentTime

	err = db.Save(&poEntities).Error

	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Update Header Failed",
			Data:       nil,
			Err:        err,
		}
	}
	//UPDATE atItemPR1
	//SET PO_SYS_NO = @Po_Sys_No,
	//	PO_LINE = @Po_Line,
	//	PO_QTY = ISNULL(PO_QTY,0) + @Item_Qty ,
	//	CHANGE_NO = CHANGE_NO + 1 ,
	//	CHANGE_USER_ID = @Change_User_Id ,
	//	CHANGE_DATETIME = @Change_Datetime
	//WHERE PR_SYS_NO = @Pr_Sys_No AND PR_LINE = @Pr_Line
	err = db.Model(&prentities).Where(transactionsparepartentities.PurchaseRequestDetail{PurchaseOrderDetailSystemNumber: payloads.PurchaseRequestDetailSystemNumber}).
		First(&prentities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed To Update Purchaes Request Because Purchase Request Data Not Found",
				Data:       nil,
				Err:        err,
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Update PR Failed",
			Data:       nil,
			Err:        err,
		}
	}
	prentities.PurchaseOrderDetailSystemNumber = entities.PurchaseOrderDetailSystemNumber
	prentities.PurchaseOrderSystemNumber = payloads.PurchaseOrderSystemNumber
	prentities.ChangeNo += 1
	prentities.UpdatedByUserId = payloads.UpdatedByUserId
	prentities.UpdatedDate = &currentTime

	err = db.Save(&prentities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Update PR Failed",
			Data:       nil,
			Err:        err,
		}
	}
	return entities, nil
}
func (repo *PurchaseOrderRepositoryImpl) DeletePurchaseOrderDetailMultiId(db *gorm.DB, s string) (bool, *exceptions.BaseErrorResponse) {
	ids := strings.Split(s, ",")
	for _, i2 := range ids {
		entities := transactionsparepartentities.PurchaseOrderDetailEntities{}
		converted, _ := strconv.Atoi(i2)
		//SELECT @POQty =ITEM_QTY  FROM atItemPO1
		//WHERE PO_SYS_NO = @Po_Sys_No AND PO_LINE = @Po_Line
		err := db.Model(&entities).Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderDetailSystemNumber: converted}).First(&entities).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Purchase Order Detail Is Not Found" + err.Error()}
			}
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusUnprocessableEntity, Message: err.Error()}
		}
		//PurchaseOrderDetailSystemNumber := entities.PurchaseOrderDetailSystemNumber
		//PurchaseOrderSystemNumber := entities.PurchaseOrderSystemNumber
		//get HeaderPO Fiest
		poEntities := transactionsparepartentities.PurchaseOrderEntities{}
		err = db.Model(&poEntities).Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: entities.PurchaseOrderSystemNumber}).First(&poEntities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: "Failed to get PO Header" + err.Error()}
		}
		//DELETE atItemPO1
		//WHERE PO_SYS_NO = @Po_Sys_No AND PO_LINE = @Po_Line
		err = db.Model(&entities).Where(transactionsparepartentities.PurchaseOrderDetailEntities{
			PurchaseOrderDetailSystemNumber: converted}).Delete(&entities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: "Failed To Delete = " + err.Error()}
		}
		PrEntities := transactionsparepartentities.PurchaseRequestDetail{}
		err = db.Model(&PrEntities).Where(transactionsparepartentities.PurchaseRequestDetail{PurchaseOrderDetailSystemNumber: converted}).
			First(&PrEntities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest,
				Message: "Purchase Request Is Not Found On Requested Payloads" + err.Error()}
		}
		//UPDATE atItemPR1 SET
		//PO_SYS_NO = 0,
		//	PO_LINE = 0,
		//	PO_QTY = ISNULL(PO_QTY,0) - ISNULL(@POQty,0),
		//CHANGE_NO = CHANGE_NO + 1 ,
		//	CHANGE_USER_ID = @Change_User_Id ,
		//	CHANGE_DATETIME = @Change_Datetime
		//WHERE PO_SYS_NO = @Po_Sys_No AND PO_LINE = @Po_Line
		PrEntities.PurchaseOrderDetailSystemNumber = 0
		PrEntities.PurchaseOrderSystemNumber = 0
		PrEntities.PurchaseOrderLine = 0
		PrEntities.ChangeNo += 1
		err = db.Save(&PrEntities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
		}
		//=-=-=-=-=-=-=-=-=-=
		var taxRate float64
		var totalDiscount float64
		var totalAmount float64
		var pkptype bool
		var CompanyVatPkpNo string
		var SupplierVatPkpNo string
		var SupplierId int
		var totalVat float64
		currentTime := time.Now().UTC()
		timeString := currentTime.Format("2006-01-02T15:04:05.000Z")
		var TaxRateResponse transactionsparepartpayloads.TaxRateResponseApi
		//		SET @TAX_RATE = dbo.getTaxPercent(dbo.getVariableValue('TAX_TYPE_PPN'),dbo.getVariableValue('TAX_SERV_CODE_PPN'),@Change_Datetime)
		TaxRateUrl := config.EnvConfigs.FinanceServiceUrl + "tax-fare/detail/tax-percent?tax_service_code=PPN&tax_type_code=PPN&effective_date=" + timeString
		if err := utils.Get(TaxRateUrl, &TaxRateResponse, nil); err != nil {
			fmt.Println(err.Error())
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Status data from external service",
				Err:        err,
			}
		}
		if TaxRateResponse.TaxPercent != nil {
			taxRate = *TaxRateResponse.TaxPercent

		} else {
			taxRate = 0
		}
		//SET @TOTAL_DISCOUNT = (SELECT SUM(ITEM_QTY * ITEM_DISC_AMOUNT) FROM atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No)

		err = db.Table("trx_item_purchase_order_detail A").Select("ISNULL(SUM(item_quantity * item_discount_amount),0)").
			Where("A.purchase_order_system_number =?", poEntities.PurchaseOrderSystemNumber).Scan(&totalDiscount).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed To Total Discount",
				Err:        err,
			}
		}
		//SET @TOTAL_AMOUNT = (SELECT SUM(ITEM_QTY * ITEM_PRICE) FROM atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No)
		err = db.Table("trx_item_purchase_order_detail A").Select("isnull(SUM(item_quantity * item_price),0)").
			Where("A.purchase_order_system_number =?", poEntities.PurchaseOrderSystemNumber).Scan(&totalAmount).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed To Total Amount",
				Err:        err,
			}
		}
		//			SET @SUPPLIER = (SELECT SUPPLIER_CODE FROM atItemPO0 WHERE PO_SYS_NO = @Po_Sys_No)
		err = db.Model(&transactionsparepartentities.PurchaseOrderEntities{}).Select("supplier_id").
			Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: poEntities.PurchaseOrderSystemNumber}).
			Scan(&SupplierId).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed To Get Supplier ID from purchase order",
				Err:        err,
			}
		}
		//SET @Company_Code = (SELECT COMPANY_CODE FROM atItemPO0 WHERE PO_SYS_NO = @Po_Sys_No)
		//

		//SET @Supplier_Vat_Pkp_No = (SELECT ISNULL(VAT_PKP_NO, '') FROM gmSupplier0 WHERE SUPPLIER_CODE = @SUPPLIER)
		var SupplierResponse transactionsparepartpayloads.SupplierResponsesAPI
		SupplierByIdUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(SupplierId)
		if err := utils.Get(SupplierByIdUrl, &SupplierResponse, nil); err != nil {
			fmt.Println(err.Error())
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed to fetch Supplier Id from external service",
				Err:        err,
			}
		}
		pkptype = SupplierResponse.TaxSupplier.PkpType

		if pkptype {
			if SupplierVatPkpNo != CompanyVatPkpNo {
				totalVat = (totalAmount - totalDiscount) * (taxRate / 100)
			} else {
				totalVat = 0
			}

		} else {
			totalVat = 0
		}
		//SET @TOTAL_AFTER_VAT = @TOTAL_AMOUNT - @TOTAL_DISCOUNT + @TOTAL_VAT

		var totalAfterVat = totalAmount - totalDiscount + totalVat
		var dpRequest float64
		if SupplierResponse.MinimumDownPayment != nil && *SupplierResponse.MinimumDownPayment != 0 {
			dpRequest = totalAfterVat * (*SupplierResponse.MinimumDownPayment / 100)
		} else {
			dpRequest = *poEntities.DpRequest
		}
		*poEntities.DpRequest = dpRequest
		*poEntities.TotalDiscount = totalDiscount
		*poEntities.TotalAmount = totalAmount
		*poEntities.TotalVat = totalVat
		*poEntities.TotalAfterVat = totalAfterVat
		poEntities.ChangeNo += 1
		err = db.Save(&poEntities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: err.Error()}
		}
	}
	return true, nil
}
func (repo *PurchaseOrderRepositoryImpl) SavePurchaseOrderDetail(db *gorm.DB, payloads transactionsparepartpayloads.PurchaseOrderSaveDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse) {
	//SELECT
	//@Pr_Po_Qty = ITEM_QTY ,-- get Old PO Qty
	//@Pr_Sys_No = PR_SYS_NO,
	//@Pr_Line = PR_LINE
	//FROM atItemPO1 where PO_SYS_NO = @Po_Sys_No AND PO_LINE = @Po_Line
	var poDetailEntities transactionsparepartentities.PurchaseOrderDetailEntities
	//var poChangedItemDetailEntities transactionsparepartentities.PurchaseOrderDetailChangedItem
	err := db.Model(&poDetailEntities).Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderDetailSystemNumber: payloads.PurchaseOrderDetailSystemNumber}).
		First(&poDetailEntities).Error
	if err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
		}
	}

	//SELECT
	//@Pr_Qty = ITEM_QTY ,
	//@PO_QTY = ISNULL(PO_QTY,0) - ISNULL(@Pr_Po_Qty,0)-- For 1 PR , 1< PO
	//FROM atItemPR1 WHERE PR_SYS_NO = @Pr_Sys_No AND PR_LINE = @Pr_Line
	var prDetailEntities transactionsparepartentities.PurchaseRequestDetail

	err = db.Model(&prDetailEntities).Where(transactionsparepartentities.PurchaseRequestDetail{PurchaseOrderDetailSystemNumber: payloads.PurchaseOrderDetailSystemNumber}).
		First(&prDetailEntities).Error
	if err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed On Get Data Purchase Request Detail" + err.Error(),
		}
	}
	var PurchaseOrderQuantity *float64

	err = db.Table("trx_item_purchase_order_detail A").Select("ISNULL(SUM(A.item_quantity),0)").
		Joins("JOIN trx_item_purchase_order B ON A.purchase_order_system_number = b.purchase_order_system_number").
		Where("B.purchase_order_status_id NOT IN (?,?) AND A.purchase_request_system_number = 1", 80, 90).
		Scan(&PurchaseOrderQuantity).Error
	if err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed Get Purchase Order Quantity",
		}
	}
	//@PO_QTY = ISNULL(PO_QTY,0) - ISNULL(@Pr_Po_Qty,0)-- For 1 PR , 1< PO
	if PurchaseOrderQuantity == nil {
		PurchaseOrderQuantity = new(float64)
		*PurchaseOrderQuantity = 0
	}
	*PurchaseOrderQuantity = *PurchaseOrderQuantity - *poDetailEntities.ItemQuantity + *payloads.ItemQuantity
	//IF @SNP>0
	//BEGIN
	//IF (@PO_QTY%@SNP)<>0
	//BEGIN
	//RAISERROR('Total Quantity does not equal multiples SNP ',16,1)
	//RETURN 0
	//END
	//END
	if payloads.Snp != nil && *payloads.Snp > 0 {
		if math.Mod(*PurchaseOrderQuantity, *payloads.Snp) != 0 {
			return poDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Total Quantity Does Not Equal Multiples SNP",
			}
		}
	}
	//IF @Pr_Qty < @PO_QTY
	//BEGIN
	//RAISERROR('Total Quantity Purchase Order cannot exceed Quantity Purchase Request',16,1)
	//RETURN 0
	//END
	if *prDetailEntities.ItemQuantity < *PurchaseOrderQuantity {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Total Quantity Purchase Order cannot exceed Quantity Purchase Request",
		}
	}
	//IF EXISTS (SELECT PO_SYS_NO FROM atItemPO1 WHERE OLD_PO_SYS_NO = @Po_Sys_No AND OLD_PO_LINE = @Old_Po_Line
	//AND ISNULL(SUBSTITUTE_TYPE,'') = dbo.getVariableValue('SUBSTITUTE_NON'))
	var exist bool
	err = db.Table("trx_item_purchase_order_detail A").Select("1").
		Joins("JOIN trx_item_purchase_order_detail_changed_item B ON A.purchase_order_detail_system_number = B.changed_item_purchase_order_detail_system_number").
		Where(transactionsparepartentities.PurchaseOrderDetailEntities{
			ChangedItemPurchaseOrderDetailSystemNumber: payloads.PurchaseOrderDetailSystemNumber,
			SubstituteTypeId: 4,
		}).Scan(&exist).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Total Quantity Purchase Order cannot exceed Quantity Purchase Request",
			Err:        err,
		}
	}

	if exist {
		poDetailEntities.ItemQuantity = payloads.ItemQuantity
		poDetailEntities.ItemPrice = payloads.ItemPrice
		poDetailEntities.ItemDiscountAmount = payloads.ItemDiscountAmount
		poDetailEntities.ItemDiscountPercentage = payloads.ItemDiscountPercentage
		*poDetailEntities.ItemTotal = *payloads.ItemTotal * (*payloads.ItemPrice * *payloads.ItemDiscountAmount)
		err = db.Save(&poDetailEntities).Error
		if err != nil {
			return poDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Total Quantity Purchase Order cannot exceed Quantity Purchase Request",
				Err:        err,
			}
		}
		//SELECT
		//@Old_Po_Sys_No = P.OLD_PO_SYS_NO,
		//@Old_Po_Line = P.OLD_PO_LINE ,
		//@Substitute_Type = P2.SUBSTITUTE_TYPE
		//FROM atItemPO1 P
		//INNER JOIN atItemPO2 P2 ON P.OLD_PO_SYS_NO = P2.OLD_PO_SYS_NO AND P.OLD_PO_LINE = P2.OLD_PO_LINE
		//WHERE PO_SYS_NO = @Po_Sys_No
		//AND PO_LINE = @Old_Po_Line
		var ChangedItemPurchaseOrderDetailSystemNumber int
		err = db.Table("trx_item_purchase_order_detail A").Select("A.changed_item_purchase_order_detail_system_number").
			Joins("trx_item_purchase_order_detail_changed_item B ON A.purchase_order_detail_system_number = B.changed_item_purchase_order_detail_system_number").
			Where(transactionsparepartentities.PurchaseOrderDetailEntities{
				ChangedItemPurchaseOrderDetailSystemNumber: payloads.PurchaseOrderDetailSystemNumber,
			}).Scan(&ChangedItemPurchaseOrderDetailSystemNumber).Error

		if err != nil {
			return poDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed To Get changes purchase order detail system number",
				Err:        err,
			}
		}
		var doQty int
		var poQty int

		//SET @Do_Qty = ISNULL((SELECT
		//SUM(DO_QTY)
		//FROM atBinningStock1
		//WHERE REF_SYS_NO = @Po_Sys_No
		//AND REF_LINE = @Old_Po_Line),0)
		err = db.Model(transactionsparepartentities.BinningStockDetail{}).
			Select("ISNULL(SUM(delivery_order_quantity),0)").
			Where(transactionsparepartentities.BinningStockDetail{PurchaseOrderDetailSystemNumber: payloads.PurchaseRequestDetailSystemNumber}).
			Scan(&doQty).Error
		if err != nil {
			return poDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed To Get Delivery Order Quantity",
				Err:        err,
			}
		}

		//SET @Po_Qty = ISNULL((SELECT
		//SUM(ITEM_QTY)
		//FROM atItemPO1
		//WHERE OLD_PO_SYS_NO = @Po_Sys_No
		//AND OLD_PO_LINE = @Old_Po_Line),0)
		err = db.Model(&poDetailEntities).Select("ISNULL(SUM(item_quantity),0)").Where(&transactionsparepartentities.PurchaseOrderDetailEntities{ChangedItemPurchaseOrderDetailSystemNumber: ChangedItemPurchaseOrderDetailSystemNumber}).
			Scan(&poQty).Error
		if err != nil {
			return poDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed To Get Purchase Order Quantity",
				Err:        err,
			}
		}
		if poQty < doQty {
			return poDetailEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Item cannot be less than " + strconv.Itoa(doQty),
				Err:        err,
			}
		}
	}
	//prDetailEntities.
	//UPDATE atItemPR1
	//SET
	//PO_QTY = ISNULL(PO_QTY,0)- ISNULL(@Pr_Po_Qty,0) + ISNUll(@Item_Qty,0) ,
	//--CHANGE_NO = CHANGE_NO + 1 ,
	//	CHANGE_USER_ID = @Change_User_Id ,
	//	CHANGE_DATETIME = @Change_Datetime
	//WHERE PR_SYS_NO = @Pr_Sys_No AND PR_LINE  = @Pr_Line
	currentTime := time.Now().UTC()
	//timeString := currentTime.Format("2006-01-02T15:04:05.000Z")
	prDetailEntities.ChangeNo += 1
	prDetailEntities.UpdatedDate = &currentTime
	prDetailEntities.UpdatedByUserId = payloads.UpdatedByUserId
	err = db.Save(&prDetailEntities).Error
	if err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Save Purchase Request Detail =" + err.Error()}
	}

	//UPDATE atItemPO1
	//SET ITEM_PRICE = @Item_Price ,
	//	ITEM_DISC_PERCENT = @Item_Disc_Percent ,
	//	ITEM_DISC_AMOUNT = @Item_Disc_Amount ,
	//	ITEM_TOTAL = @Item_Total ,
	//	ITEM_REMARK = @Item_Remark,
	//	ITEM_QTY = @Item_Qty,
	//--CHANGE_NO = CHANGE_NO + 1 ,
	//	CHANGE_USER_ID = @Change_User_Id ,
	//	CHANGE_DATETIME = @Change_Datetime
	//WHERE PO_SYS_NO = @Po_Sys_No AND PO_LINE = @Po_Line

	poDetailEntities.ItemPrice = payloads.ItemPrice
	poDetailEntities.ItemDiscountPercentage = payloads.ItemDiscountPercentage
	poDetailEntities.ItemDiscountAmount = payloads.ItemDiscountAmount
	poDetailEntities.ItemTotal = payloads.ItemTotal
	poDetailEntities.ItemRemark = payloads.ItemRemark
	poDetailEntities.ItemQuantity = payloads.ItemQuantity
	*poDetailEntities.UpdatedDate = currentTime
	poDetailEntities.UpdatedByUserId = payloads.UpdatedByUserId
	err = db.Save(&poDetailEntities).Error
	if err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Save Purchase Order Detail =" + err.Error()}
	}

	var taxRate float64
	var totalDiscount float64
	var totalAmount float64
	var pkptype bool
	var CompanyVatPkpNo string
	var SupplierVatPkpNo string
	var SupplierId int
	var totalVat float64
	//currentTime := time.Now().UTC()
	timeString := currentTime.Format("2006-01-02T15:04:05.000Z")
	var TaxRateResponse transactionsparepartpayloads.TaxRateResponseApi
	//		SET @TAX_RATE = dbo.getTaxPercent(dbo.getVariableValue('TAX_TYPE_PPN'),dbo.getVariableValue('TAX_SERV_CODE_PPN'),@Change_Datetime)
	TaxRateUrl := config.EnvConfigs.FinanceServiceUrl + "tax-fare/detail/tax-percent?tax_service_code=PPN&tax_type_code=PPN&effective_date=" + timeString
	if err := utils.Get(TaxRateUrl, &TaxRateResponse, nil); err != nil {
		fmt.Println(err.Error())
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Status data from external service",
			Err:        err,
		}
	}
	if TaxRateResponse.TaxPercent != nil {
		taxRate = *TaxRateResponse.TaxPercent

	} else {
		taxRate = 0
	}
	//SET @TOTAL_DISCOUNT = (SELECT SUM(ITEM_QTY * ITEM_DISC_AMOUNT) FROM atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No)

	err = db.Table("trx_item_purchase_order_detail A").Select("ISNULL(SUM(item_quantity * item_discount_amount),0)").
		Where("A.purchase_order_system_number =?", payloads.PurchaseOrderSystemNumber).Scan(&totalDiscount).Error
	if err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Total Discount",
			Err:        err,
		}
	}
	//SET @TOTAL_AMOUNT = (SELECT SUM(ITEM_QTY * ITEM_PRICE) FROM atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No)
	err = db.Table("trx_item_purchase_order_detail A").Select("isnull(SUM(item_quantity * item_price),0)").
		Where("A.purchase_order_system_number =?", payloads.PurchaseOrderSystemNumber).Scan(&totalAmount).Error
	if err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Total Amount",
			Err:        err,
		}
	}
	//			SET @SUPPLIER = (SELECT SUPPLIER_CODE FROM atItemPO0 WHERE PO_SYS_NO = @Po_Sys_No)
	err = db.Model(&transactionsparepartentities.PurchaseOrderEntities{}).Select("supplier_id").
		Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: payloads.PurchaseOrderSystemNumber}).
		Scan(&SupplierId).Error
	if err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Get Supplier ID from purchase order",
			Err:        err,
		}
	}
	//SET @Company_Code = (SELECT COMPANY_CODE FROM atItemPO0 WHERE PO_SYS_NO = @Po_Sys_No)
	//
	//var CompanyId int
	var poEntities transactionsparepartentities.PurchaseOrderEntities
	err = db.Model(&poEntities).Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: payloads.PurchaseOrderSystemNumber}).
		First(&poEntities).Error
	if err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed To Get Customer from purchase order",
			Err:        err,
		}
	}
	//SET @Supplier_Vat_Pkp_No = (SELECT ISNULL(VAT_PKP_NO, '') FROM gmSupplier0 WHERE SUPPLIER_CODE = @SUPPLIER)
	var SupplierResponse transactionsparepartpayloads.SupplierResponsesAPI
	SupplierByIdUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(SupplierId)
	if err := utils.Get(SupplierByIdUrl, &SupplierResponse, nil); err != nil {
		fmt.Println(err.Error())
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed to fetch Supplier Id from external service",
			Err:        err,
		}
	}
	pkptype = SupplierResponse.TaxSupplier.PkpType
	SupplierVatPkpNo = SupplierResponse.TaxSupplier.PkpNo
	//SET @Company_Vat_Pkp_No = (SELECT ISNULL(VAT_PKP_NO, '') FROM gmComp0 WHERE COMPANY_CODE = @Company_Code)
	var CompanyDetailResponse transactionsparepartpayloads.CompanyDetailResponses
	CompanyDetailUrl := config.EnvConfigs.GeneralServiceUrl + "company-detail/" + strconv.Itoa(poEntities.CompanyId)
	if err := utils.Get(CompanyDetailUrl, &CompanyDetailResponse, nil); err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed to fetch Company Id Detail from external service",
			Err:        err,
		}
	}
	//IF @PKP='Y'
	//BEGIN
	//IF LTRIM(RTRIM(ISNULL(@Supplier_Vat_Pkp_No,'')))<>LTRIM(RTRIM(ISNULL(@Company_Vat_Pkp_No,'')))
	//BEGIN
	//SET @TOTAL_VAT = (@TOTAL_AMOUNT - @TOTAL_DISCOUNT) * (@TAX_RATE / 100)
	//END
	//ELSE
	//BEGIN
	//SET @TOTAL_VAT = 0
	//END
	//END
	//ELSE
	//BEGIN
	//SET @TOTAL_VAT = 0
	//END

	if pkptype {
		if SupplierVatPkpNo != CompanyVatPkpNo {
			totalVat = (totalAmount - totalDiscount) * (taxRate / 100)
		} else {
			totalVat = 0
		}

	} else {
		totalVat = 0
	}
	//SET @TOTAL_AFTER_VAT = @TOTAL_AMOUNT - @TOTAL_DISCOUNT + @TOTAL_VAT

	totalAfterVat := totalAmount - totalDiscount + totalVat
	err = db.Model(&poEntities).Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: payloads.PurchaseOrderSystemNumber}).
		Scan(&poEntities).Error
	if err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Header Not Found",
			Err:        err,
		}
	}
	// AMBIL DP REQ DARI SUPPLIER

	//var SupplierResponses transactionsparepartpayloads.SupplierResponsesAPI
	//		SET @TAX_RATE = dbo.getTaxPercent(dbo.getVariableValue('TAX_TYPE_PPN'),dbo.getVariableValue('TAX_SERV_CODE_PPN'),@Change_Datetime)
	//SupplierResponseUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(poEntities.SupplierId)
	//if err := utils.Get(SupplierResponseUrl, &SupplierResponses, nil); err != nil {
	//	fmt.Println(err.Error())
	//	return entities, &exceptions.BaseErrorResponse{
	//		StatusCode: http.StatusInternalServerError,
	//		Message:    "Failed to fetch Status data from external service",
	//		Err:        err,
	//	}
	//}
	//dpRequest = SupplierResponse.MinimumDownPayment
	var dpRequest float64
	if SupplierResponse.MinimumDownPayment != nil && *SupplierResponse.MinimumDownPayment != 0 {
		dpRequest = totalAfterVat * (*SupplierResponse.MinimumDownPayment / 100)
	} else {
		dpRequest = *poEntities.DpRequest
	}
	*poEntities.DpRequest = dpRequest
	*poEntities.TotalDiscount = totalDiscount
	*poEntities.TotalAmount = totalAmount
	*poEntities.TotalVat = totalVat
	*poEntities.TotalAfterVat = totalAfterVat
	poEntities.ChangeNo += 1
	poEntities.UpdatedByUserId = payloads.UpdatedByUserId
	*poEntities.UpdatedDate = currentTime

	err = db.Save(&poEntities).Error

	if err != nil {
		return poDetailEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Update Header Failed",
			Data:       nil,
			Err:        err,
		}
	}
	return poDetailEntities, nil
}
func (repo *PurchaseOrderRepositoryImpl) DeleteDocument(db *gorm.DB, i int) (bool, *exceptions.BaseErrorResponse) {
	var PrDetailSysNo []int
	var poDetailEntities transactionsparepartentities.PurchaseOrderDetailEntities
	err := db.Model(&poDetailEntities).Select("purchase_request_detail_system_number").Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderSystemNumber: i}).
		Scan(&PrDetailSysNo).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Update Header Failed",
				Data:       nil,
				Err:        err,
			}
		}
	}
	for _, prDetail := range PrDetailSysNo {
		var prDetailEntities transactionsparepartentities.PurchaseRequestDetail
		err = db.Model(&transactionsparepartentities.PurchaseRequestDetail{}).
			Where(transactionsparepartentities.PurchaseRequestDetail{PurchaseOrderDetailSystemNumber: prDetail}).
			First(&prDetailEntities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Purchase Request Is Not Found",
				Data:       nil,
				Err:        err,
			}
		}
		//update pr
		prDetailEntities.PurchaseOrderDetailSystemNumber = 0
		prDetailEntities.PurchaseOrderLine = 0
		prDetailEntities.ChangeNo += 1
		*prDetailEntities.UpdatedDate = time.Now()
		err = db.Save(&prDetailEntities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Purchase Request Update Failed, " + err.Error(),
				Data:       nil,
				Err:        err,
			}
		}
		err = db.Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: i}).
			Delete(&transactionsparepartentities.PurchaseOrderEntities{}).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed To Delete Header, " + err.Error(),
				Data:       nil,
				Err:        err,
			}
		}
		err = db.Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderSystemNumber: i}).
			Delete(&transactionsparepartentities.PurchaseOrderDetailEntities{}).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed To Delete Detail, " + err.Error(),
				Data:       nil,
				Err:        err,
			}
		}
	}
	return true, nil
}

func (repo *PurchaseOrderRepositoryImpl) GetFromPurchaseRequest(db *gorm.DB, filter []utils.FilterCondition, pagination pagination.Pagination) {
	//TODO implement me
	panic("implement me")
}

func (repo *PurchaseOrderRepositoryImpl) SubmitPurchaseOrderRequest(db *gorm.DB, payloads transactionsparepartpayloads.PurchaseOrderHeaderSubmitRequest) (bool, *exceptions.BaseErrorResponse) {
	var poEntities transactionsparepartentities.PurchaseOrderEntities
	//SELECT @Company_Code = COMPANY_CODE,
	//@Supplier_Code = SUPPLIER_CODE,
	//@Order_Type = ORDER_TYPE
	//FROM atitempo0
	//WHERE PO_SYS_NO = @Po_Sys_No
	whereQuery := transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: payloads.PurchaseOrderSystemNumber}
	err := db.Model(&poEntities).Where(whereQuery).Scan(&poEntities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Purchase Order Header Is Not Found" + err.Error(),
			Data:       nil,
			Err:        err,
		}
	}
	// cek apakah supplier codenya ini itu atpm vw atau bukan
	//123 itu id nya gmm(garuda mataram motor)
	if poEntities.SupplierId == 123 {
		totalOrder := 0
		totalOrderLimit := 0
		//SET @TotalOrder = (SELECT COUNT(PO_SYS_NO)
		//FROM atitempo0
		//WHERE COMPANY_CODE = @Company_Code
		//AND YEAR(@Po_Doc_Date) = YEAR(PO_DOC_DATE)
		//AND MONTH(@Po_Doc_Date) = MONTH(PO_DOC_DATE))
		whereQuery = transactionsparepartentities.PurchaseOrderEntities{CompanyId: payloads.CompanyId}
		err = db.Model(&poEntities).Select("count(purchase_order_system_number)").
			Where(whereQuery).
			Where("YEAR(purchase_order_date) = ? "+
				"AND MONTH(purchase_order_date) = ?",
				payloads.PurchaseOrderDocumentDate.Year(),
				payloads.PurchaseOrderDocumentDate.Month()).
			Scan(&totalOrder).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Failed To Get Total Order" + err.Error(),
				Data:       nil,
				Err:        err,
			}
		}
		var OrderLimitEntities transactionsparepartentities.PurchaseOrderLimit
		//SET @TotalOrderLimit = (SELECT CAST(ISNULL(TABLE_KEY1,0) AS INT)
		//FROM comgentable2
		//WHERE TABLE_CODE = 'GMMPOLimit'
		//AND TABLE_KEY0 = @ORDER_TYPE)
		err = db.Model(&OrderLimitEntities).Select("ISNULL(order_limit)").
			Where(transactionsparepartentities.PurchaseOrderLimit{CompanyId: payloads.CompanyId,
				OrderTypeId: payloads.PurchaseOrderTypeId}).
			Scan(&totalOrder).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Failed To Get Purchase Order Limit" + err.Error(),
				Data:       nil,
				Err:        err,
			}
		}
		//IF (@TotalOrder > @TotalOrderLimit)
		//BEGIN
		//SET @ErrMsg = 'This dealer has exceed''s GMM monthly order limit.\nGMM Order Limit : '+CAST(@TotalOrderLimit AS VARCHAR)+'\nCurrent Dealer Order : '+CAST(@TotalOrder AS VARCHAR)
		//RAISERROR(@ErrMsg,16,1)
		//RETURN 0
		//END
		if totalOrder > totalOrderLimit {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message: "This dealer has exceed''s GMM monthly order limit. \nGMM Order Limit" +
					"\n GMM Order Limit : " + strconv.Itoa(totalOrderLimit) +
					"\n Current Dealer Order : " + strconv.Itoa(totalOrder),
				Data: nil,
				Err:  err,
			}
		}
	}
	//SET @Supplier_Default = dbo.getVariableValue('FREE_ACCS_SUPPLIER_DEFAULT')
	//62
	//DECLARE @Is_SP_PO integer = 1
	//,@tmpPo_Date DATETIME = GETDATE()
	//select top 1 @Is_SP_PO = ISNULL(V.ORDER_NO,1), @tmpPo_Date = ISNULL(P0.PO_DOC_DATE,GETDATE())
	//	from atItemPO0 P0 LEFT JOIN comGenVariable V ON V.VALUE = P0.ITEM_GROUP and VARIABLE like 'ITMGRP_%'
	//	where P0.PO_SYS_NO = @Po_Sys_No

	//IsSparepartPo := 1
	//tmpPoDate := time.Now()

	//CEK APAKAH GRUP ITU IN ATAU OJ?
	var ItemGroup transactionsparepartpayloads.PurchaseOrderItemGroupResponse
	ItemGroupURL := config.EnvConfigs.GeneralServiceUrl + "item-group/" + strconv.Itoa(payloads.ItemGroupId)
	if err := utils.Get(ItemGroupURL, &ItemGroup, nil); err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item Group data from external service",
			Err:        err,
		}
	}
	//		SET @TAX_RATE = dbo.getTaxPercent(dbo.getVariableValue('TAX_TYPE_PPN'),dbo.getVariableValue('TAX_SERV_CODE_PPN'),@Change_Datetime)
	var PeriodResponse masterpayloads.OpenPeriodPayloadResponse
	PeriodUrl := config.EnvConfigs.FinanceServiceUrl + "closing-period-company/current-period?company_id=" + strconv.Itoa(payloads.CompanyId) + "&closing_module_detail_code=" //strconv.Itoa(response.ItemCode)
	//UomItem := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/" + res.ItemCode + "/P" //strconv.Itoa(response.ItemCode)
	//IF @Is_SP_PO = 1
	//if purchase order is sparepart po
	if ItemGroup.ItemGroupCode == "IN" || ItemGroup.ItemGroupCode == "OJ" {
		//BEGIN
		//IF dbo.getPeriodStatus(@Company_Code,@Po_Doc_Date,dbo.getVariableValue('MODULE_SP')) <> 'O'
		//BEGIN
		//RAISERROR('Period is already Closed',16,1)
		//RETURN 0
		//END
		//SET @Po_Doc_Date = GETDATE()
		PeriodUrl = PeriodUrl + "SP"
		if err := utils.Get(PeriodUrl, &PeriodResponse, nil); err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed to Period Response data from external service",
				Err:        err,
			}
		}
		if PeriodResponse.PeriodYear != strconv.Itoa(payloads.PurchaseOrderDocumentDate.Year()) && PeriodResponse.PeriodMonth != strconv.Itoa(payloads.PurchaseOrderDocumentDate.Year()) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Period is already Closed",
				Err:        err,
			}
		}
		//SET @Po_Doc_Date = GETDATE()

		*payloads.PurchaseOrderDocumentDate = time.Now()
		//get Current Open Period
	} else {
		//if purchase order is ga po
		PeriodUrl = PeriodUrl + "FXA"
		if err := utils.Get(PeriodUrl, &PeriodResponse, nil); err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed to Get Period Response data from external service",
				Err:        err,
			}
		}
		if PeriodResponse.PeriodYear != strconv.Itoa(payloads.PurchaseOrderDocumentDate.Year()) && PeriodResponse.PeriodMonth != strconv.Itoa(payloads.PurchaseOrderDocumentDate.Year()) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Period is already Closed",
				Err:        err,
			}
		}
		//SET @Po_Doc_Date = @tmpPo_Date

		if poEntities.PurchaseOrderDocumentDate == nil {
			poEntities.PurchaseOrderDocumentDate = new(time.Time)
			*poEntities.PurchaseOrderDocumentDate = time.Now()
		} else {
			*payloads.PurchaseOrderDocumentDate = *poEntities.PurchaseOrderDocumentDate
		}
	}

	//--tidak boleh campur item gmitem0.item_regulation
	//IF (SELECT COUNT(DISTINCT ISNULL(B.ITEM_REGULATION,'N'))
	//FROM dbo.atItemPO1 A
	//INNER JOIN gmitem0 B ON A.ITEM_CODE = B.ITEM_CODE
	//INNER JOIN ATITEMPO0 C ON C.PO_SYS_NO = A.PO_SYS_NO
	//WHERE A.PO_SYS_NO = @Po_Sys_No
	//AND C.COMPANY_CODE = 3125098
	//AND ISNULL(B.ITEM_REGULATION,'N') IN ('Y','N')) > 1
	//BEGIN
	//RAISERROR('Regulated item can''t be mixed with non regulated item on same PO.',16,1)
	//RETURN 0
	//END
	var CompanyResponses generalservicepayloads.GetCompanyByIdResponses
	CompanyUrl := config.EnvConfigs.GeneralServiceUrl + "company/" + strconv.Itoa(payloads.CompanyId)
	if err := utils.Get(CompanyUrl, &CompanyResponses, nil); err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Failed to Get Company data from external service",
			Err:        err,
		}
	}
	//checking regulated item but based on sp only checking for nmdi
	if CompanyResponses.CompanyCode == "3125098" {
		count := 0
		//Untuk hitung apaakah ada item regulation yang dicampur atau aman
		err = db.Table("trx_item_purchase_order_detail A").Select("distinct COUNT(ISNULL(B.is_item_regulation,'N')").
			Joins("INNER JOIN mtr_item B ON A.item_id = B.item_id").
			Joins("INNER JOIN trx_item_purchase_order C ON A.old_purchase_order_system_no = C.purchase_order_system_number").
			Where("A.purchase_order_system_number = ? AND ISNULL(B.is_item_regulation,'N') IN ('Y','N') ", payloads.PurchaseOrderSystemNumber).
			Scan(&count).
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed To Checked Regulated item ",
				Err:        err,
			}
		}
		if count != 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Regulated item can''t be mixed with non regulated item on same PO.",
				Err:        err,
			}
		}

		//check of defect item
		//DECLARE @ListDefectItem varchar(500) = (SELECT STUFF((SELECT DISTINCT ',' + A.ITEM_CODE +'<br>'
		//FROM dbo.atItemPO1 A
		//INNER JOIN gmitem0 B ON A.ITEM_CODE = B.ITEM_CODE
		//INNER JOIN ATITEMPO0 C ON C.PO_SYS_NO = A.PO_SYS_NO
		//WHERE A.PO_SYS_NO = @Po_Sys_No
		//AND C.SUPPLIER_CODE = '3125098'
		//AND ISNULL(B.TECHNICAL_DEFECT,'NO') IN ('YES','Y')
		//FOR XML PATH ('')),1,1,''))
		var listOfDefectItem []string
		err = db.Table("trx_item_purchase_order_detail A").Select("distinct COUNT(ISNULL(B.item_code,'N')").
			Joins("INNER JOIN mtr_item B ON A.item_id = B.item_id").
			Joins("INNER JOIN trx_item_purchase_order C ON A.old_purchase_order_system_no = C.purchase_order_system_number").
			Where("A.purchase_order_system_number = ? AND ISNULL(B.is_technical_defect,'NO') IN ('YES','Y') ", payloads.PurchaseOrderSystemNumber).
			Scan(&listOfDefectItem).
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Failed To technical defect item ",
				Err:        err,
			}
		}
		if len(listOfDefectItem) > 0 {
			DefectedItem := "There is technical defect Item in Purchase Order : \n"
			for _, defected := range listOfDefectItem {
				DefectedItem = DefectedItem + defected + "\n"
			}
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    DefectedItem,
				Err:        err,
			}
		}

	}
	//IF EXISTS (SELECT PO_LINE FROM dbo.atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No AND ITEM_QTY <= 0)
	//BEGIN
	//RAISERROR('Qty PO must be bigger than 0',16,1)
	//RETURN 0
	//END
	var exist bool
	err = db.Model(transactionsparepartentities.PurchaseOrderDetailEntities{}).Select("1").
		Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderSystemNumber: payloads.PurchaseOrderSystemNumber}).
		Where("item_quantity <= 0").Scan(&exist).Error
	if err != nil || !exist {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Qty PO must be bigger than 0",
			Err:        err,
		}
	}

	//IF EXISTS (SELECT PO_LINE FROM dbo.atItemPO1 WHERE PO_SYS_NO = @Po_Sys_No AND ITEM_PRICE <= 0)
	//BEGIN
	//RAISERROR('Price PO must be bigger than 0',16,1)
	//RETURN 0
	//END
	err = db.Model(transactionsparepartentities.PurchaseOrderDetailEntities{}).Select("1").
		Where(transactionsparepartentities.PurchaseOrderDetailEntities{PurchaseOrderSystemNumber: payloads.PurchaseOrderSystemNumber}).
		Where("item_price <= 0").Scan(&exist).Error
	if err != nil || !exist {
		return false, &exceptions.BaseErrorResponse{

			StatusCode: http.StatusBadRequest,
			Message:    "Price PO must be bigger than 0",
			Err:        err,
		}
	}
	//CHECK IF PO FOR NEW ACCESCORIES MUST TO DEFAULT FREE ACCESORIES SUPPLIER
	//IF EXISTS (SELECT PO_LINE FROM atItemPO1 P1
	//INNER JOIN atItemPO0 P ON P.PO_SYS_NO = P1.PO_SYS_NO
	//WHERE P1.PO_SYS_NO = @Po_Sys_No AND ISNULL(VEHICLE_CHASSIS_NO,'') <> ''
	//AND P.SUPPLIER_CODE <> @Supplier_Default)
	//BEGIN
	//SET @ErrMsg = 'PO for Free Accessories Standard must be directed to Supplier ' + @Supplier_Default
	//RAISERROR(@ErrMsg,16,1)
	//RETURN 0
	//END
	exist = false
	err = db.Table("trx_item_purchase_order A").Select("1").
		Joins("trx_item_purchase_order_detail B ON A.").
		Where("A.purchase_order_system_number = ? AND ISNULL(B.vehicle_id,'')<>''", payloads.PurchaseOrderSystemNumber).
		Scan(&exist).
		Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Err:        err,
		}
	}
	var supplierDefault string
	if !exist {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "PO for Free Accessories Standard must be directed to Supplier " + supplierDefault,
			Err:        err,
		}
	}
	//err = db.Table()
	return true, nil
}
