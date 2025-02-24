package transactionsparepartrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	masterrepository "after-sales/api/repositories/master"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SalesOrderRepositoryImpl struct {
	locationStockImpl masterrepository.LocationStockRepository
}

func StartSalesOrderRepositoryImpl() transactionsparepartrepository.SalesOrderRepository {
	locationStockRepoImpl := masterwarehouserepository.NewLocationStockRepositoryImpl()
	return &SalesOrderRepositoryImpl{locationStockImpl: locationStockRepoImpl}
}

// [dbo].[uspg_atSalesOrder0_Insert] option = 0
func (r *SalesOrderRepositoryImpl) InsertSalesOrderHeader(db *gorm.DB, payload transactionsparepartpayloads.SalesOrderInsertHeaderPayload) (transactionsparepartentities.SalesOrder, *exceptions.BaseErrorResponse) {
	//IF @Cpc_Code <> '' AND NOT EXISTS (SELECT CPC_CODE FROM gmCCPC0 WHERE CPC_TYPE = 'P' AND CPC_CODE = @Cpc_Code)
	//BEGIN
	//SET @Error = 'Profit Center Code is not valid'
	//END
	Entities := transactionsparepartentities.SalesOrder{}
	ProfitCenter, ProfitCenterErr := generalserviceapiutils.GetProfitCenterById(payload.ProfitCenterID)
	if ProfitCenterErr != nil {
		return Entities, ProfitCenterErr
	}
	if ProfitCenter.ProfitCenterId == 0 {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("profit Center Code is not valid"),
			Message:    "Profit Center Code is not valid",
		}
	}
	Event, EventErr := financeserviceapiutils.GetEventById(payload.EventNumberID)
	if EventErr != nil {
		return Entities, EventErr
	}
	if Event.EventId == 0 && Event.EventNo == "" {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("event number is not valid"),
			Message:    "event number is not valid",
		}
	}
	//IF @Cust_Code <> '' AND NOT EXISTS (SELECT CUSTOMER_CODE FROM GMCUST0 A INNER JOIN GMCUSTTYPE B ON A.CUSTOMER_TYPE = B.CUSTOMER_TYPE WHERE A.CUSTOMER_CODE = @Cust_Code)
	//BEGIN
	//SET @Error = @Error + ', Customer Code is not valid'
	//END
	CustomerData, CustomerDataErr := generalserviceapiutils.GetCustomerMasterById(payload.CustomerID)
	if CustomerDataErr != nil {
		return Entities, CustomerDataErr
	}
	if CustomerData.CustomerId == 0 {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("customer is not valid"),
			Message:    "customer is not valid",
		}
	}
	TermOfPayment, TermOfPaymentErr := generalserviceapiutils.GetTermOfPaymentById(payload.TermOfPaymentID)
	if TermOfPaymentErr != nil {
		return Entities, TermOfPaymentErr
	}
	if TermOfPayment.TermOfPaymentId == 0 {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("term of payment is not valid"),
			Message:    "term of payment is not valid",
		}
	}
	//IF @Whs_Group <> '' AND NOT EXISTS (SELECT WAREHOUSE_GROUP FROM GMLOC0 WHERE WAREHOUSE_GROUP = @Whs_Group)
	//BEGIN
	//IF @Error <> ''
	//BEGIN
	//SET @Error = @Error + ', Warehouse Group is not valid'
	//END
	warehouseGroup := masterwarehouseentities.WarehouseGroup{}
	err := db.Model(&warehouseGroup).Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupId: payload.WarehouseGroupID}).
		First(&warehouseGroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadGateway,
				Err:        err,
				Message:    "warehouse group is not valid please check input",
			}
		}
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get warehouse group please check input",
		}
	}
	//	IF @Sales_Emp_No <> '' AND NOT EXISTS (SELECT EMPLOYEE_NO FROM gmEmp1 WHERE EMPLOYEE_NO = @Sales_Emp_No AND COMPANY_CODE = @Company_Code AND RECORD_STATUS = @Record_Active) --:GH diganti gmEmp1 req: pak tepen
	//	BEGIN
	//	IF @Error <> ''
	//	BEGIN
	//	SET @Error = @Error + ', Sales Person Code is not valid'
	//	END
	employee, employeeErr := generalserviceapiutils.GetEmployeeById(payload.SalesEmployeeID)
	if employeeErr != nil {
		return Entities, employeeErr
	}
	if employee.UserEmployeeId == 0 && employee.EmployeeName == "" {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("sales person is not valid"),
			Message:    "sales person is not valid",
		}
	}
	//IF @Cust_Code = CAST(@Company_Code AS VARCHAR)
	//get type
	//PurchaseOrderTypePoSo
	PotypeSo := masterentities.PurchaseOrderTypeSalesOrderEntity{}
	err = db.Model(&PotypeSo).Where(masterentities.PurchaseOrderTypeSalesOrderEntity{PurchaseOrderTypeSalesOrderId: payload.PurchaseOrderTypeID}).
		First(&PotypeSo).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
			Message:    "purchase order type id is not found please check input",
		}
	}
	if PotypeSo.PurchaseOrderTypeSalesOrderCode == "SO2" {
		//IF ISNULL(@Po_Sys_No,0) <> 0 AND NOT EXISTS (SELECT PO_SYSTEM_NO FROM utPO0 WHERE PO_SYSTEM_NO = ISNULL(@Po_Sys_No,0) AND ISNULL(SO_SYS_NO,0)=0)
		//BEGIN
		//IF @Error <> ''
		//BEGIN
		//SET @Error = @Error + ', Purchase Order No. is not valid'
		//END
		POExpedition, POExpeditionErr := salesserviceapiutils.GetPurchaseOrderExpeditionById(payload.PurchaseOrderSystemNumber)
		if POExpeditionErr != nil {
			return Entities, POExpeditionErr
		}
		if POExpedition.UnitPurchaseOrderExpeditionSystemNumber == 0 {
			return Entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("purchase order number is not valid"),
				Message:    "purchase order number is not valid",
			}
		}
		//payload.VehicleSalesOrderSystemNumber = POExpedition.system
	} else {
		if payload.PurchaseOrderSystemNumber != 0 {

			//IF ISNULL(@Po_Sys_No,0) <> 0 AND NOT EXISTS (SELECT PO_SYS_NO FROM atItemPO0 WHERE PO_SYS_NO = ISNULL(@Po_Sys_No,0) AND ISNULL(SO_SYS_NO,0)=0)
			isExist := 0
			err = db.Model(&transactionsparepartentities.PurchaseOrderEntities{}).
				Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: payload.PurchaseOrderSystemNumber}).
				Select("1").Scan(&isExist).Error
			if err != nil {
				return Entities, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
					Message:    "failed to check purchase order please check input",
				}
			}
			if isExist == 0 {
				return Entities, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Err:        err,
					Message:    "Purchase Order No. is not valid",
				}
			}
		}
	}

	//PurchaseOrderEntities := transactionsparepartentities.PurchaseOrderEntities{}
	if payload.CustomerID != payload.CompanyID {
		//SELECT @Remark = PO_REMARK FROM atItemPO0 WHERE PO_SYS_NO = ISNULL(@Po_Sys_No,0)
		err = db.Model(&transactionsparepartentities.PurchaseOrderEntities{}).
			Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: payload.PurchaseOrderSystemNumber}).
			Select("purchase_order_remark").Scan(&payload.Remark).Error
		if err != nil {
			return Entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
				Message:    "failed to get remark from purchase order please check input",
			}
		}
	}
	//get transaction type sales order first
	TransactionTypeSalesOrder := masterentities.TransactionTypeSalesOrder{}
	err = db.Model(&TransactionTypeSalesOrder).
		Where(masterentities.TransactionTypeSalesOrder{TransactionTypeSalesOrderId: payload.TransactionTypeID}).
		First(&TransactionTypeSalesOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
				Message:    "transaction type sales order is not found in master table please check input",
			}
		}
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
			Message:    "failed to get transaction type sales order please check input",
		}
	}
	//BEGIN
	//DECLARE @bolInternalUsePPN bit = dbo.getVariableValue('PPN_INTERNAL_IS_ACTIVE_FLAG')
	//IF @Trx_Type = @trxtypecentralize OR (@Trx_Type = @trxtypeinternal AND @bolInternaluseppn = 0)
	//BEGIN
	//SET @Vat_Tax_Percent = 0
	//END
	//ELSE
	//BEGIN
	//SET @Vat_Tax_Percent  =	dbo.getTaxPercent(dbo.getVariableValue('TAX_TYPE_PPN'),dbo.getVariableValue('TAX_SERV_CODE_PPN'),ISNULL(@So_Date,GETDATE()))
	//END
	//END
	//code centralized = SU06
	//code internal = SU05

	if TransactionTypeSalesOrder.TransactionTypeSalesOrderCode == "SU05" ||
		TransactionTypeSalesOrder.TransactionTypeSalesOrderCode == "SU06" {
		payload.VATTaxPercentage = 0
	} else {
		GetTaxPercent, GetTaxPercentErr := financeserviceapiutils.GetTaxPercent(financeserviceapiutils.TaxPercentParams{
			TaxServiceCode: "PPN",
			TaxTypeCode:    "PPN",
			EffectiveDate:  payload.SalesOrderDate,
		})
		if GetTaxPercentErr != nil {
			return Entities, GetTaxPercentErr
		}
		payload.VATTaxPercentage = GetTaxPercent.TaxPercent
	}

	//SET @Src_Code  = dbo.getVariableValue('SRC_DOC_SP_SE')
	//
	//EXEC uspg_gmSrcDoc1_Update
	//@Option = 0 ,
	//@COMPANY_CODE = @Company_Code ,
	//@SOURCE_CODE = @Src_Code ,
	//@VEHICLE_BRAND = @Vehicle_Brand ,
	//@PROFIT_CENTER_CODE = @Cpc_Code ,
	//@TRANSACTION_CODE = '' ,
	//@BANK_ACC_CODE = '' ,
	//@TRANSACTION_DATE = @SO_DATE ,
	//@Last_Doc_No =  @Se_Doc_No OUTPUT
	//get apporval status draft
	ApprovalStatusDraft, ApprovalStatusErr := generalserviceapiutils.GetApprovalStatusByCode(utils.ApprovalDraftCode)
	if ApprovalStatusErr != nil {
		return Entities, ApprovalStatusErr
	}
	nowTime := time.Now()
	Entities = transactionsparepartentities.SalesOrder{
		CompanyID:                     payload.CompanyID,
		SalesOrderStatusID:            ApprovalStatusDraft.ApprovalStatusId,
		SalesOrderDate:                payload.SalesOrderDate,
		SalesEstimationDocumentNumber: payload.SalesEstimationDocumentNumber,
		BrandID:                       payload.BrandID,
		CostCenterID:                  payload.CostCenterID,
		ProfitCenterID:                payload.ProfitCenterID,
		EventNumberID:                 payload.EventNumberID,
		SalesOrderIsAffiliated:        payload.SalesOrderIsAffiliated,
		TransactionTypeID:             payload.TransactionTypeID,
		SalesOrderIsOneTimeCustomer:   payload.SalesOrderIsOneTimeCustomer,
		CustomerID:                    payload.CustomerID,
		TermOfPaymentID:               payload.TermOfPaymentID,
		SameTaxArea:                   payload.SameTaxArea,
		ETDTime:                       payload.ETDTime,
		DeliveryAddressSameAsInvoice:  payload.DeliveryAddressSameAsInvoice,
		DeliveryContactPerson:         payload.DeliveryContactPerson,
		DeliveryAddress:               payload.DeliveryAddress,
		DeliveryAddressLine1:          payload.DeliveryAddressLine1,
		DeliveryAddressLine2:          payload.DeliveryAddressLine2,
		DeliveryPhoneNumber:           payload.DeliveryPhoneNumber,
		PurchaseOrderSystemNumber:     payload.PurchaseOrderSystemNumber,
		OrderTypeID:                   payload.OrderTypeID,
		DeliveryViaID:                 payload.DeliveryViaID,
		PayType:                       payload.PayType,
		WarehouseGroupID:              payload.WarehouseGroupID,
		BackOrder:                     payload.BackOrder,
		SetOrder:                      payload.SetOrder,
		DownPaymentAmount:             payload.DownPaymentAmount,
		Remark:                        payload.Remark,
		AgreementID:                   payload.AgreementID,
		SalesEmployeeID:               payload.SalesEmployeeID,
		CurrencyID:                    payload.CurrencyID,
		CurrencyExchangeID:            payload.CurrencyExchangeID,
		Total:                         payload.Total,
		TotalDiscount:                 payload.TotalDiscount,
		AdditionalDiscountPercentage:  payload.AdditionalDiscountPercentage,
		AdditionalDiscountStatusID:    ApprovalStatusDraft.ApprovalStatusId,
		VATTaxID:                      payload.VATTaxId,
		TotalVAT:                      payload.TotalVAT,
		VATTaxPercentage:              payload.VATTaxPercentage,
		TotalAfterVAT:                 payload.TotalAfterVAT,
		VehicleSalesOrderSystemNumber: payload.VehicleSalesOrderSystemNumber,
		VehicleSalesOrderDetailID:     payload.VehicleSalesOrderDetailID,
		PurchaseOrderTypeID:           payload.PurchaseOrderTypeID,
		CreatedByUserId:               payload.CreatedByUserId,
		CreatedDate:                   &nowTime,
		UpdatedByUserId:               payload.CreatedByUserId,
		UpdatedDate:                   &nowTime,
	}
	err = db.Create(&Entities).First(&Entities).Error
	if err != nil {
		return Entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to save sales order please check input error",
		}
	}
	return Entities, nil
}

// USPG_ATSALESORDER0_SELECT option = 0
func (r *SalesOrderRepositoryImpl) GetSalesOrderByID(db *gorm.DB, Id int) (transactionsparepartpayloads.SalesOrderEstimationGetByIdResponse, *exceptions.BaseErrorResponse) {
	var response transactionsparepartpayloads.SalesOrderEstimationGetByIdResponse

	var SoEntities transactionsparepartentities.SalesOrder

	err := db.Model(&SoEntities).Where(transactionsparepartentities.SalesOrder{SalesOrderSystemNumber: Id}).
		First(&SoEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "sales order with that is is not found",
			}
		}
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "error occur when get sales order please contact administrator",
		}
	}

	//approval
	Approval, ApprovalErr := generalserviceapiutils.GetApprovalStatusById(SoEntities.SalesOrderStatusID)
	if ApprovalErr != nil {
		return response, ApprovalErr
	}
	//vehicle brand
	unitBrand, unitBrandErr := salesserviceapiutils.GetUnitBrandById(SoEntities.BrandID)
	if unitBrandErr != nil {
		return response, unitBrandErr
	}

	TransactionTypeEntities := masterentities.TransactionTypeSalesOrder{}

	err = db.Model(&TransactionTypeEntities).Where(masterentities.TransactionTypeSalesOrder{TransactionTypeSalesOrderId: SoEntities.TransactionTypeID}).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "error occur when get sales order transaction type please contact administrator",
		}
	}
	//term of payment
	termOfPayment, termOfPaymentErr := generalserviceapiutils.GetTermOfPaymentById(SoEntities.TermOfPaymentID)
	if termOfPaymentErr != nil {
		return response, termOfPaymentErr
	}
	//get purchase order entities
	PoEntities := transactionsparepartentities.PurchaseOrderEntities{}
	err = db.Model(&PoEntities).Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: SoEntities.PurchaseOrderSystemNumber}).
		Scan(&PoEntities).Error
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get purchase order entities please check input",
		}
	}
	AdditionalDiscountApproval, AdditionalDiscountErr := generalserviceapiutils.GetApprovalStatusById(SoEntities.AdditionalDiscountStatusID)
	if AdditionalDiscountErr != nil {
		return response, AdditionalDiscountErr
	}
	poCompany, poCompanyErr := generalserviceapiutils.GetCompanyDataById(SoEntities.CompanyID)
	if poCompanyErr != nil {
		return response, poCompanyErr
	}
	WarehouseGroup := masterwarehouseentities.WarehouseGroup{}
	err = db.Model(&WarehouseGroup).Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupId: SoEntities.WarehouseGroupID}).
		First(&WarehouseGroup).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get warehouse group please contact your administrator",
		}
	}
	//SoCompanyDetail, SoCompanyDetailErr := generalserviceapiutils.GetCompanyDataById(SoEntities.CompanyID)
	//if SoCompanyDetailErr != nil {
	//	return response, SoCompanyDetailErr
	//}
	response = transactionsparepartpayloads.SalesOrderEstimationGetByIdResponse{
		CompanyID:                           SoEntities.CompanyID,
		SalesOrderSystemNumber:              SoEntities.SalesOrderSystemNumber,
		SalesOrderDocumentNumber:            SoEntities.SalesOrderDocumentNumber,
		SalesOrderStatusID:                  SoEntities.SalesOrderStatusID,
		SalesOrderStatusDescription:         Approval.ApprovalStatusDescription,
		SalesOrderDate:                      SoEntities.SalesOrderDate,
		SalesEstimationDocumentNumber:       SoEntities.SalesEstimationDocumentNumber,
		BrandID:                             SoEntities.BrandID,
		VehicleBrandCode:                    unitBrand.BrandCode,
		VehicleBrandDescription:             unitBrand.BrandName,
		ProfitCenterID:                      SoEntities.ProfitCenterID,
		EventNumberID:                       SoEntities.EventNumberID,
		SalesOrderIsAffiliated:              SoEntities.SalesOrderIsAffiliated,
		TransactionTypeID:                   SoEntities.TransactionTypeID,
		TransactionTypeDescription:          TransactionTypeEntities.TransactionTypeSalesOrderDescription,
		SalesOrderIsOneTimeCustomer:         SoEntities.SalesOrderIsOneTimeCustomer,
		CustomerID:                          SoEntities.CustomerID,
		TermOfPaymentID:                     SoEntities.TermOfPaymentID,
		TermOfPaymentDescription:            termOfPayment.TermOfPaymentCode,
		SameTaxArea:                         SoEntities.SameTaxArea,
		EstimatedTimeOfDelivery:             SoEntities.EstimatedTimeOfDelivery,
		DeliveryAddressSameAsInvoice:        SoEntities.DeliveryAddressSameAsInvoice,
		DeliveryContactPerson:               SoEntities.DeliveryContactPerson,
		DeliveryAddress:                     SoEntities.DeliveryAddress,
		DeliveryAddressLine1:                SoEntities.DeliveryAddressLine1,
		DeliveryAddressLine2:                SoEntities.DeliveryAddressLine2,
		PurchaseOrderSystemNumber:           SoEntities.PurchaseOrderSystemNumber,
		PurchaseOrdeCompanyId:               SoEntities.CompanyID,
		PurchaseOrderCompanyCode:            poCompany.CompanyCode,
		PurchaseOrderDocumentNumber:         PoEntities.PurchaseOrderDocumentNumber,
		OrderTypeID:                         SoEntities.OrderTypeID,
		DeliveryViaID:                       SoEntities.DeliveryViaID,
		DeliveryViaDescription:              "",
		PayType:                             SoEntities.PayType,
		WarehouseGroupID:                    SoEntities.WarehouseGroupID,
		WarehouseGroupDescription:           WarehouseGroup.WarehouseGroupName,
		BackOrder:                           SoEntities.BackOrder,
		SetOrder:                            SoEntities.SetOrder,
		DownPaymentAmount:                   SoEntities.DownPaymentAmount,
		DownPaymentPaidAmount:               SoEntities.DownPaymentPaidAmount,
		DownPaymentPaymentAllocated:         SoEntities.DownPaymentPaymentAllocated,
		DownPaymentPaymentVAT:               SoEntities.DownPaymentPaymentVAT,
		DownPaymentAllocatedToInvoice:       SoEntities.DownPaymentAllocatedToInvoice,
		DownPaymentVATAllocatedToInvoice:    SoEntities.DownPaymentVATAllocatedToInvoice,
		Remark:                              SoEntities.Remark,
		AgreementID:                         SoEntities.AgreementID,
		SalesEmployeeID:                     SoEntities.SalesEmployeeID,
		CurrencyID:                          SoEntities.CurrencyID,
		CurrencyExchangeID:                  SoEntities.CurrencyExchangeID,
		CurrencyRateTypeID:                  SoEntities.CurrencyRateTypeID,
		Total:                               SoEntities.Total,
		TotalDiscount:                       SoEntities.TotalDiscount,
		AdditionalDiscountPercentage:        SoEntities.AdditionalDiscountPercentage,
		AdditionalDiscountAmount:            SoEntities.AdditionalDiscountAmount,
		AdditionalDiscountStatusID:          SoEntities.AdditionalDiscountStatusID,
		AdditionalDiscountStatusDescription: AdditionalDiscountApproval.ApprovalStatusDescription,
		VATTaxID:                            SoEntities.VATTaxID,
		VATTaxPercentage:                    SoEntities.VATTaxPercentage,
		TotalVAT:                            SoEntities.TotalVAT,
		TotalAfterVAT:                       SoEntities.TotalAfterVAT,
		ApprovalRequestNumber:               SoEntities.ApprovalRequestNumber,
		ApprovalRemark:                      SoEntities.ApprovalRemark,
		VehicleSalesOrderSystemNumber:       SoEntities.VehicleSalesOrderSystemNumber,
		VehicleSalesOrderDetailID:           SoEntities.VehicleSalesOrderDetailID,
		PurchaseOrderTypeID:                 SoEntities.PurchaseOrderTypeID,
		CostCenterID:                        SoEntities.CostCenterID,
		IsAtpm:                              false,
		AtpmInternalPurpose:                 SoEntities.AtpmInternalPurpose,
	}

	return response, nil
}

// @Option = 2 [uspg_atSalesOrder0_Select]
func (r *SalesOrderRepositoryImpl) GetAllSalesOrder(db *gorm.DB, pages pagination.Pagination, condition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	Response := []transactionsparepartpayloads.GetAllSalesOrderResponse{}
	entities := []transactionsparepartentities.SalesOrder{}
	joinTable := db.Model(&transactionsparepartentities.SalesOrder{}).Select("*")
	WhereQuery := utils.ApplyFilter(joinTable, condition)
	err := WhereQuery.Scopes(pagination.Paginate(&pages, WhereQuery)).Scan(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get sales order",
			Err:        err,
		}
	}
	if len(entities) == 0 {
		pages.Rows = []string{}
		return pages, nil
	}
	for _, item := range entities {
		//reference document number
		var referenceDocumentNumber string
		if item.PurchaseOrderSystemNumber != 0 {
			purchaseOrderEntities := transactionsparepartentities.PurchaseOrderEntities{}
			err = db.Model(&purchaseOrderEntities).Where(transactionsparepartentities.PurchaseOrderEntities{PurchaseOrderSystemNumber: item.PurchaseOrderSystemNumber}).
				Scan(&purchaseOrderEntities).Error

			if err != nil {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to get purchase order",
					Err:        err,
				}
			}
			referenceDocumentNumber = purchaseOrderEntities.PurchaseOrderDocumentNumber
		}
		if item.VehicleSalesOrderSystemNumber != 0 {
			vehicleSalesOrder, vehicleSalesOrderSys := salesserviceapiutils.GetVehicleSalesOrderById(item.VehicleSalesOrderSystemNumber)
			if vehicleSalesOrderSys != nil {
				return pages, vehicleSalesOrderSys
			}
			referenceDocumentNumber = vehicleSalesOrder.VehicleSalesOrderDocumentNumber
		}
		//get customer name
		customerData, customerDataErr := generalserviceapiutils.GetCustomerMasterById(item.CustomerID)
		if customerDataErr != nil {
			return pages, customerDataErr
		}
		//get approval status
		salesOrderApprovalStatus, salesOrderApprovalStatusErr := generalserviceapiutils.GetApprovalStatusById(item.SalesOrderStatusID)
		if salesOrderApprovalStatusErr != nil {
			return pages, salesOrderApprovalStatusErr
		}
		//get transaction type

		transactionTypeEntities := masterentities.TransactionTypeSalesOrder{}
		err = db.Model(&transactionTypeEntities).Where(masterentities.TransactionTypeSalesOrder{TransactionTypeSalesOrderId: item.TransactionTypeID}).
			First(&transactionTypeEntities).Error
		if err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to get transaction type",
				Err:        err,
			}
		}
		//created by user id
		userDetail, userDetailErr := generalserviceapiutils.GetUserDetailsByID(item.CreatedByUserId)
		if userDetailErr != nil {
			return pages, userDetailErr
		}
		Response = append(Response, transactionsparepartpayloads.GetAllSalesOrderResponse{
			SalesOrderSystemNumber:        item.SalesOrderSystemNumber,
			SalesOrderDocumentNumber:      item.SalesOrderDocumentNumber,
			SalesEstimationDocumentNumber: item.SalesEstimationDocumentNumber,
			SalesOrderDate:                item.SalesOrderDate,
			ReferenceDocumentNumber:       referenceDocumentNumber,
			CustomerID:                    item.CustomerID,
			CustomerName:                  customerData.CustomerName,
			SalesOrderStatusID:            item.SalesOrderStatusID,
			SalesOrderStatusDescription:   salesOrderApprovalStatus.ApprovalStatusDescription,
			TransactionTypeID:             item.TransactionTypeID,
			TransactionTypeDescription:    transactionTypeEntities.TransactionTypeSalesOrderDescription,
			CreatedByUserId:               item.CreatedByUserId,
			CreatedByUserName:             userDetail.EmployeeName,
			PurchaseOrderSystemNumber:     item.PurchaseOrderSystemNumber,
			VehicleSalesOrderSystemNumber: item.VehicleSalesOrderDetailID,
		})
	}
	pages.Rows = Response
	return pages, nil
}
func (r *SalesOrderRepositoryImpl) VoidSalesOrder(db *gorm.DB, salesOrderId int) (bool, *exceptions.BaseErrorResponse) {
	//get data entities first
	entities := transactionsparepartentities.SalesOrder{}
	err := db.Model(&entities).Where(transactionsparepartentities.SalesOrder{SalesOrderSystemNumber: salesOrderId}).
		First(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "sales order is not found with that id please check input",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to retreives sales order",
		}
	}
	//get id for sales order approval draft

	approvalDraft, approvalDraftId := generalserviceapiutils.GetApprovalStatusById(entities.SalesOrderStatusID)
	if approvalDraftId != nil {
		return false, approvalDraftId
	}
	if entities.SalesOrderStatusID != approvalDraft.ApprovalStatusId {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "document is not draft",
			Err:        err,
		}
	}
	if entities.PurchaseOrderTypeID != 0 {
		//update purchaseorder entities
		err = db.Model(&transactionsparepartentities.PurchaseOrderEntities{}).
			Where("purchase_order_system_number = ?", entities.PurchaseOrderSystemNumber).
			Update("sales_order_system_number", 0).
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed to update purchase order system number",
			}
		}
		//update utpo0

		if entities.VehicleSalesOrderSystemNumber != 0 {
			//update vehicle sales order waiting for endpoint
		}

		//get appoval status canceled
		approvalCancel, approvalCancelErr := generalserviceapiutils.GetApprovalStatusByCode(utils.ApprovalCancelledCode)
		if approvalCancelErr != nil {
			return false, approvalCancelErr
		}

		entities.SalesOrderStatusID = approvalCancel.ApprovalStatusId

		//get detail first
		//update status sales order
		err = db.Model(&transactionsparepartentities.SalesOrderDetail{}).
			Where(transactionsparepartentities.SalesOrderDetail{SalesOrderSystemNumber: entities.SalesOrderSystemNumber}).
			Update("sales_order_line_status_id", approvalCancel.ApprovalStatusId).
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed to update sales order detail",
			}
		}
	}
	err = db.Save(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to update sales order",
			Err:        err,
		}
	}
	return true, nil
}

// [uspg_atSalesOrder1_Insert] option 1
func (r *SalesOrderRepositoryImpl) InsertSalesOrderDetail(db *gorm.DB, payload transactionsparepartpayloads.SalesOrderDetailInsertPayload) (transactionsparepartentities.SalesOrderDetail, *exceptions.BaseErrorResponse) {
	//get item exist first
	//entities
	entities := transactionsparepartentities.SalesOrderDetail{}
	soEntities := transactionsparepartentities.SalesOrder{}
	var SalesOrderSubTotal float64
	var SalesOrderTotalDiscount float64
	var SalesOrderTotal float64
	var SalesOrderAdditionalDiscountAmount float64
	var SalesOrderTotalVat float64
	var SalesOrderTotalAfterVat float64
	////var SalesOrder
	var availableQuantity float64 = 0
	err := db.Model(&soEntities).Where(transactionsparepartentities.SalesOrder{SalesOrderSystemNumber: payload.SalesOrderSystemNumber}).
		First(&soEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "failed to get sales order header",
				Err:        err,
			}
		}
	}
	isItemExist := false
	err = db.Model(&entities).
		Where(transactionsparepartentities.SalesOrderDetail{ItemId: payload.ItemId, SalesOrderSystemNumber: payload.SalesOrderSystemNumber}).
		Select("1").
		Scan(&isItemExist).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to check sales order detail occurance please check input",
		}
	}
	if isItemExist {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
			Message:    "Data item is already exist",
		}
	}
	//check technical defect for company : 3125098 NMDI
	NmdiCompany, getCompanyErr := generalserviceapiutils.GetCompanyDataByCode("3125098")
	if getCompanyErr != nil {
		return entities, getCompanyErr
	}
	fmt.Println(NmdiCompany)
	//cek technical defect
	isTechincalDefect := false
	err = db.Model(&masteritementities.Item{}).
		Where(masteritementities.Item{ItemId: payload.ItemId}).
		Select("is_technical_defect").
		Scan(&isTechincalDefect).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to check technical defect on item",
		}
	}
	if isTechincalDefect {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("item has technical defect"),
			Message:    "Item has technical defect",
		}
	}
	var itemTypeResponse masteritementities.ItemType
	if err := db.Where("item_type_code = ?", "G").First(&itemTypeResponse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item type not found",
				Err:        fmt.Errorf("item type with code %s not found", "G"),
			}
		}
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item type code",
			Err:        err,
		}
	}

	//check if exist with item goods
	isItemExist = false
	err = db.Model(&masteritementities.Item{}).
		Where(masteritementities.Item{ItemId: payload.ItemId, ItemTypeId: itemTypeResponse.ItemTypeId}).
		Select("1").
		Scan(&isItemExist).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to check if item exist in item master",
		}
	}

	if isItemExist {
		//EXEC dbo.uspg_amLocationStockItem_Select
		//@Option = 1,
		//@Company_Code = @Company_code,
		//@Period_Date = @Creation_Datetime ,
		//@Whs_Code = '' ,--@Whs_Code ,
		//@Loc_Code = '' ,--@Loc_Code ,
		//@Item_Code = @ITEM_CODE ,
		//@Whs_Group = @WHS_GROUP ,
		//@UoM_Type = @UoM_Type ,
		//@QtyResult = @QTY_AVAIL OUTPUT
		avaibleItem, errStockLocations := r.locationStockImpl.GetAvailableQuantity(db, masterwarehousepayloads.GetAvailableQuantityPayload{
			CompanyId:        soEntities.CompanyID,
			PeriodDate:       soEntities.SalesOrderDate,
			WarehouseId:      169,
			LocationId:       0,
			ItemId:           payload.ItemId,
			WarehouseGroupId: payload.ItemId,
			UomTypeId:        1,
		})
		//availableItem, errStockLocation := aftersalesserviceapiutils.GetAvailableItemLocationStock(masterwarehousepayloads.GetAvailableQuantityPayload{
		//	CompanyId:        soEntities.CompanyID,
		//	PeriodDate:       soEntities.SalesOrderDate,
		//	WarehouseId:      169,
		//	LocationId:       0,
		//	ItemId:           payload.ItemId,
		//	WarehouseGroupId: soEntities.WarehouseGroupID,
		//	UomTypeId:        1,
		//})
		if errStockLocations != nil {
			return entities, errStockLocations
		}
		availableQuantity = avaibleItem.QuantityAvailable
	} else {
		availableQuantity = 1
	}
	//
	//IF @Company_Code IN ('3125098','1516098','200000')
	//BEGIN
	//SET @QTY_AVAIL = 1
	//END
	//checking IF @Company_Code IN ('3125098','1516098','200000') NMDI,KIA,GMM
	KiaCompany, KiaCompanyErr := generalserviceapiutils.GetCompanyDataByCode("1516098")
	if KiaCompanyErr != nil {
		return entities, KiaCompanyErr
	}
	GmmCompany, GmmCompanyErr := generalserviceapiutils.GetCompanyDataByCode("200000")
	if GmmCompanyErr != nil {
		return entities, GmmCompanyErr
	}
	if soEntities.CompanyID == NmdiCompany.CompanyId ||
		soEntities.CompanyID == KiaCompany.CompanyId ||
		soEntities.CompanyID == GmmCompany.CompanyId {
		availableQuantity = 1
	}
	//get approval draft status
	approvalDraft, approvalDraftErr := generalserviceapiutils.GetApprovalStatusByCode(utils.ApprovalDraftCode)
	if approvalDraftErr != nil {
		return entities, approvalDraftErr
	}
	itemSubstituteType := masteritementities.ItemSubstituteType{}
	err = db.Model(&itemSubstituteType).Where(masteritementities.ItemSubstituteType{ItemSubstituteTypeCode: "S"}).
		Scan(&itemSubstituteType).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get item substitute type",
		}
	}
	if availableQuantity > 0 {

		//insert
		entities = transactionsparepartentities.SalesOrderDetail{
			SalesOrderSystemNumber:              payload.SalesOrderSystemNumber,
			SalesOrderLineStatusId:              &approvalDraft.ApprovalStatusId,
			ItemSubstituteId:                    payload.ItemSubstituteId,
			ItemId:                              payload.ItemId,
			QuantityDemand:                      payload.QuantityDemand,
			IsAvailable:                         true,
			QuantitySupply:                      payload.QuantitySupply,
			QuantityPick:                        payload.QuantityPick,
			UomId:                               payload.UomId,
			Price:                               payload.Price,
			PriceEffectiveDate:                  payload.PriceEffectiveDate,
			DiscountPercent:                     payload.DiscountPercent,
			DiscountAmount:                      payload.Price * (payload.DiscountPercent / 100),
			DiscountRequestPercent:              payload.DiscountRequestPercent,
			DiscountRequestAmount:               payload.Price * (payload.DiscountRequestPercent / 100),
			Remark:                              payload.Remark,
			ApprovalRequestNumber:               payload.ApprovalRequestNumber,
			ApprovalRemark:                      payload.ApprovalRemark,
			VehicleSalesOrderSystemNumber:       payload.VehicleSalesOrderSystemNumber,
			VehicleSalesOrderDetailSystemNumber: payload.VehicleSalesOrderDetailSystemNumber,
			PriceListId:                         payload.PriceListId,
			ItemSubstituteTypeId:                itemSubstituteType.ItemSubstituteTypeId, //question item substitute type

		}
		err = db.Create(&entities).First(&entities).Error
		if err != nil {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to save sales order detail",
				Err:        err,
			}
		}
	} else {
		//process insert substitute item jka available item 0 -> not yet develop
		entities = transactionsparepartentities.SalesOrderDetail{
			SalesOrderSystemNumber:              payload.SalesOrderSystemNumber,
			SalesOrderLineStatusId:              &approvalDraft.ApprovalStatusId,
			ItemSubstituteId:                    payload.ItemSubstituteId,
			ItemId:                              payload.ItemId,
			QuantityDemand:                      payload.QuantityDemand,
			IsAvailable:                         true,
			QuantitySupply:                      payload.QuantitySupply,
			QuantityPick:                        payload.QuantityPick,
			UomId:                               payload.UomId,
			Price:                               payload.Price,
			PriceEffectiveDate:                  payload.PriceEffectiveDate,
			DiscountPercent:                     payload.DiscountPercent,
			DiscountAmount:                      payload.Price * (payload.DiscountPercent / 100),
			DiscountRequestPercent:              payload.DiscountRequestPercent,
			DiscountRequestAmount:               payload.Price * (payload.DiscountRequestPercent / 100),
			Remark:                              payload.Remark,
			ApprovalRequestNumber:               payload.ApprovalRequestNumber,
			ApprovalRemark:                      payload.ApprovalRemark,
			VehicleSalesOrderSystemNumber:       payload.VehicleSalesOrderSystemNumber,
			VehicleSalesOrderDetailSystemNumber: payload.VehicleSalesOrderDetailSystemNumber,
			PriceListId:                         payload.PriceListId,
			ItemSubstituteTypeId:                1,
		}
		err = db.Create(&entities).First(&entities).Error
		if err != nil {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to save sales order detail",
				Err:        err,
			}
		}
	}

	//process recalculate header
	//get first subsitute type SUBTITUTE_ITEM [S]

	//caclculate subtotal
	err = db.Model(&entities).
		Where(transactionsparepartentities.SalesOrderDetail{SalesOrderSystemNumber: payload.SalesOrderSystemNumber}).
		Select(`
				COALESCE(SUM(CASE WHEN item_substitute_type_id = ? THEN 0 ELSE COALESCE(price, 0) * COALESCE(quantity_demand, 0) END), 0)
				`, itemSubstituteType.ItemSubstituteTypeId).
		Scan(&SalesOrderSubTotal).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get sub total from sales order detail",
			Err:        err,
		}
	}
	//round to nearest integer
	SalesOrderSubTotal = math.Round(SalesOrderSubTotal)

	//calculate total discount
	err = db.Model(&entities).
		Where(transactionsparepartentities.SalesOrderDetail{SalesOrderSystemNumber: payload.SalesOrderSystemNumber}).
		Select(`
				COALESCE(SUM(CASE WHEN COALESCE(discount_request_amount,0) > 0 THEN 
				COALESCE(discount_request_amount, 0) * COALESCE(quantity_demand, 0) 
				ELSE 
				COALESCE(discount_amount,0) * COALESCE(quantity_demand,0)
				END), 0)
				`).
		Scan(&SalesOrderTotalDiscount).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get sub total from sales order detail",
			Err:        err,
		}
	}
	//rounding
	SalesOrderTotalDiscount = math.Round(SalesOrderTotalDiscount)
	//SET @Total = @Sub_Total - @Total_Disc
	SalesOrderTotal = SalesOrderSubTotal - SalesOrderTotalDiscount
	SalesOrderAdditionalDiscountAmount = SalesOrderTotal * (payload.AdditionalDiscountPercentage / 100)
	//rounding
	SalesOrderAdditionalDiscountAmount = math.Round(SalesOrderAdditionalDiscountAmount)
	//notes  math.Round(value/10) * 10 -> sama dengan round(...,-1)
	SalesOrderTotalVat = (SalesOrderTotal - SalesOrderAdditionalDiscountAmount) * (math.Round(SalesOrderTotalVat))

	SalesOrderTotalAfterVat = (SalesOrderTotal - SalesOrderAdditionalDiscountAmount) + SalesOrderTotalVat

	//update header
	soEntities.Total = SalesOrderTotal
	soEntities.TotalDiscount = SalesOrderTotalDiscount
	soEntities.AdditionalDiscountPercentage = payload.AdditionalDiscountPercentage
	soEntities.AdditionalDiscountAmount = SalesOrderAdditionalDiscountAmount
	soEntities.TotalVAT = SalesOrderTotalVat
	soEntities.TotalAfterVAT = SalesOrderTotalAfterVat
	soEntities.ChangeNo += 1
	soEntities.Remark = payload.HeaderRemark

	//save
	err = db.Save(&soEntities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to save sales order entities",
			Err:        err,
		}
	}
	//fmt.Println(availableQuantity)
	return entities, nil
}

// type option 1  interface{}
func (r *SalesOrderRepositoryImpl) DeleteSalesOrderDetail(db *gorm.DB, salesOrderDetailId int) (bool, *exceptions.BaseErrorResponse) {
	//get and check entities first
	var SalesOrderSubTotal float64
	var SalesOrderTotalDiscount float64
	var SalesOrderTotal float64
	var SalesOrderAdditionalDiscountAmount float64
	var SalesOrderTotalVat float64
	var SalesOrderTotalAfterVat float64

	detailEntities := transactionsparepartentities.SalesOrderDetail{}
	err := db.Model(&detailEntities).
		Where(transactionsparepartentities.SalesOrderDetail{SalesOrderDetailSystemNumber: salesOrderDetailId}).
		First(&detailEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "sales order detail to deleted is not found please check input",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get sales order detail to deleted please check input",
			Err:        err,
		}
	}

	//get sales order header
	soEntities := transactionsparepartentities.SalesOrder{}
	err = db.Model(&soEntities).
		Where(transactionsparepartentities.SalesOrder{SalesOrderSystemNumber: detailEntities.SalesOrderSystemNumber}).
		First(&soEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "sales order header is not found",
			}
		}
	}
	//validasi if already have pick
	if detailEntities.QuantityPick > 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("item already picked"),
			Message:    "item alread picked",
		}
	}
	if detailEntities.QuantitySupply > 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("item already supplied"),
			Message:    "item already supplied",
		}
	}
	//cek approved closed id
	approvalClosed, approvalClosedErr := generalserviceapiutils.GetApprovalStatusByCode(utils.ApprovalClosedCode)
	if approvalClosedErr != nil {
		return false, approvalClosedErr
	}
	if soEntities.SalesOrderStatusID == approvalClosed.ApprovalStatusId {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("sales order is already closed"),
			Message:    "sales order is already closed",
		}
	}
	if soEntities.PurchaseOrderSystemNumber != 0 || soEntities.VehicleSalesOrderSystemNumber != 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("item cannot be deleted there is purchase order and vehicle sales order number binding"),
			Message:    "item cannot be deleted",
		}
	}

	//delete sales order detail
	err = db.Delete(&detailEntities, transactionsparepartentities.SalesOrderDetail{SalesOrderDetailSystemNumber: salesOrderDetailId}).
		Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to delete sales order entity",
		}
	}
	//get item substiute type
	itemSubstituteType := masteritementities.ItemSubstituteType{}
	err = db.Model(&itemSubstituteType).Where(masteritementities.ItemSubstituteType{ItemSubstituteTypeCode: "S"}).
		Scan(&itemSubstituteType).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get item substitute type",
		}
	}

	//////////////////////////////
	//recalculated for header data
	err = db.Model(&detailEntities).
		Where(transactionsparepartentities.SalesOrderDetail{SalesOrderSystemNumber: detailEntities.SalesOrderSystemNumber}).
		Select(`
				COALESCE(SUM(CASE WHEN item_substitute_type_id = ? THEN 0 ELSE COALESCE(price, 0) * COALESCE(quantity_demand, 0) END), 0)
				`, itemSubstituteType.ItemSubstituteTypeId).
		Scan(&SalesOrderSubTotal).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get sub total from sales order detail",
			Err:        err,
		}
	}

	//round to nearest integer
	SalesOrderSubTotal = math.Round(SalesOrderSubTotal)

	//calculate total discount
	err = db.Model(&detailEntities).
		Where(transactionsparepartentities.SalesOrderDetail{SalesOrderSystemNumber: detailEntities.SalesOrderSystemNumber}).
		Select(`
				COALESCE(SUM(CASE WHEN COALESCE(discount_request_amount,0) > 0 THEN 
				COALESCE(discount_request_amount, 0) * COALESCE(quantity_demand, 0) 
				ELSE 
				COALESCE(discount_amount,0) * COALESCE(quantity_demand,0)
				END), 0)
				`).
		Scan(&SalesOrderTotalDiscount).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get sub total from sales order detail",
			Err:        err,
		}
	}
	//rounding
	SalesOrderTotalDiscount = math.Round(SalesOrderTotalDiscount)
	//SET @Total = @Sub_Total - @Total_Disc
	SalesOrderTotal = SalesOrderSubTotal - SalesOrderTotalDiscount
	//SET @Add_Disc_Amount = (@Total * ISNULL(@Add_Disc_Percent,0)) / 100

	SalesOrderAdditionalDiscountAmount = SalesOrderTotal * (soEntities.AdditionalDiscountPercentage / 100)
	//rounding
	SalesOrderAdditionalDiscountAmount = math.Round(SalesOrderAdditionalDiscountAmount)
	//notes  math.Round(value/10) * 10 -> sama dengan round(...,-1)
	SalesOrderTotalVat = (SalesOrderTotal - SalesOrderAdditionalDiscountAmount) * (math.Round(SalesOrderTotalVat))

	SalesOrderTotalAfterVat = (SalesOrderTotal - SalesOrderAdditionalDiscountAmount) + SalesOrderTotalVat

	//update header data
	soEntities.Total = SalesOrderTotal
	soEntities.TotalDiscount = SalesOrderTotalDiscount
	soEntities.AdditionalDiscountAmount = SalesOrderAdditionalDiscountAmount
	soEntities.TotalVAT = SalesOrderTotalVat
	soEntities.TotalAfterVAT = SalesOrderTotalAfterVat
	//soEntities.Remark =

	err = db.Save(&soEntities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to save sales order header",
		}
	}
	return true, nil

}

// propose disc USPG_ATSALESORDER1_UPDATE option 7
func (r *SalesOrderRepositoryImpl) SalesOrderProposedDiscountMultiId(db *gorm.DB, multiId string, proposedDiscountPercentage float64) (bool, *exceptions.BaseErrorResponse) {
	var SalesOrderSubTotal float64
	var SalesOrderTotalDiscount float64
	var SalesOrderTotal float64
	var SalesOrderAdditionalDiscountAmount float64
	var SalesOrderTotalVat float64
	var SalesOrderTotalAfterVat float64
	var soEntities transactionsparepartentities.SalesOrder

	multiIds := strings.Split(multiId, ",")
	for _, Id := range multiIds {
		SoDetailId, errConvert := strconv.Atoi(Id)
		if errConvert != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errConvert,
				Message:    "failed to read sales order id",
			}
		}

		//select the sales order Detail entities

		salesOrderDetailEntities := transactionsparepartentities.SalesOrderDetail{}
		err := db.Model(&salesOrderDetailEntities).
			Where(transactionsparepartentities.SalesOrderDetail{SalesOrderDetailSystemNumber: SoDetailId}).
			First(&salesOrderDetailEntities).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Err:        fmt.Errorf("sales order detail with id : %d is not exist", SoDetailId),
					Message:    fmt.Sprintf("sales order detail with id : %d is not exist", SoDetailId),
				}
			}
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "there is an error when getting sales order detail",
			}
		}
		//check if sales order header is exist
		err = db.Model(&soEntities).
			Where(transactionsparepartentities.SalesOrderDetail{
				SalesOrderSystemNumber: salesOrderDetailEntities.SalesOrderSystemNumber,
			}).
			First(&soEntities).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Err:        fmt.Errorf("sales order  with id : %d is not exist", salesOrderDetailEntities.SalesOrderSystemNumber),
					Message:    fmt.Sprintf("sales order  with id : %d is not exist", salesOrderDetailEntities.SalesOrderSystemNumber),
				}
			}
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed to get sales order header",
			}
		}
		//update the sales order detail
		salesOrderDetailEntities.DiscountRequestPercent = proposedDiscountPercentage
		salesOrderDetailEntities.DiscountRequestAmount = math.Round(salesOrderDetailEntities.Price * (proposedDiscountPercentage / 100))
		//save sales order

		err = db.Save(&salesOrderDetailEntities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed to save sales order detail please contact your administrator",
			}
		}
	}
	//updating header
	itemSubstituteType := masteritementities.ItemSubstituteType{}
	err := db.Model(&itemSubstituteType).Where(masteritementities.ItemSubstituteType{ItemSubstituteTypeCode: "S"}).
		Scan(&itemSubstituteType).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get item substitute type",
		}
	}

	detailEntities := transactionsparepartentities.SalesOrderDetail{}
	err = db.Model(&detailEntities).
		Where(transactionsparepartentities.SalesOrderDetail{SalesOrderSystemNumber: detailEntities.SalesOrderSystemNumber}).
		Select(`
				COALESCE(SUM(CASE WHEN item_substitute_type_id = ? THEN 0 ELSE COALESCE(price, 0) * COALESCE(quantity_demand, 0) END), 0)
				`, itemSubstituteType.ItemSubstituteTypeId).
		Scan(&SalesOrderSubTotal).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get sub total from sales order detail",
			Err:        err,
		}
	}

	//round to nearest integer
	SalesOrderSubTotal = math.Round(SalesOrderSubTotal)

	//calculate total discount
	err = db.Model(&detailEntities).
		Where(transactionsparepartentities.SalesOrderDetail{SalesOrderSystemNumber: detailEntities.SalesOrderSystemNumber}).
		Select(`
				COALESCE(SUM(CASE WHEN COALESCE(discount_request_amount,0) > 0 THEN 
				COALESCE(discount_request_amount, 0) * COALESCE(quantity_demand, 0) 
				ELSE 
				COALESCE(discount_amount,0) * COALESCE(quantity_demand,0)
				END), 0)
				`).
		Scan(&SalesOrderTotalDiscount).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get sub total from sales order detail",
			Err:        err,
		}
	}
	//rounding
	SalesOrderTotalDiscount = math.Round(SalesOrderTotalDiscount)
	//SET @Total = @Sub_Total - @Total_Disc
	SalesOrderTotal = SalesOrderSubTotal - SalesOrderTotalDiscount
	//SET @Add_Disc_Amount = (@Total * ISNULL(@Add_Disc_Percent,0)) / 100

	SalesOrderAdditionalDiscountAmount = SalesOrderTotal * (soEntities.AdditionalDiscountPercentage / 100)
	//rounding
	SalesOrderAdditionalDiscountAmount = math.Round(SalesOrderAdditionalDiscountAmount)
	//notes  math.Round(value/10) * 10 -> sama dengan round(...,-1)
	SalesOrderTotalVat = (SalesOrderTotal - SalesOrderAdditionalDiscountAmount) * (math.Round(SalesOrderTotalVat))

	SalesOrderTotalAfterVat = (SalesOrderTotal - SalesOrderAdditionalDiscountAmount) + SalesOrderTotalVat

	//process update header
	soEntities.TotalDiscount = SalesOrderTotal
	soEntities.TotalDiscount = SalesOrderTotalDiscount
	soEntities.AdditionalDiscountPercentage = SalesOrderAdditionalDiscountAmount
	soEntities.AdditionalDiscountAmount = SalesOrderAdditionalDiscountAmount
	soEntities.TotalVAT = SalesOrderTotalVat
	soEntities.TotalAfterVAT = SalesOrderTotalAfterVat
	err = db.Save(&soEntities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to save sales order entities with error : " + err.Error(),
		}
	}
	return true, nil
}
func (r *SalesOrderRepositoryImpl) UpdateSalesOrderHeader(db *gorm.DB, payload transactionsparepartpayloads.SalesOrderUpdatePayload, SalesOrderId int) (transactionsparepartentities.SalesOrder, *exceptions.BaseErrorResponse) {
	//cek first if id is exist
	soEntity := transactionsparepartentities.SalesOrder{}
	err := db.Model(&soEntity).Where(transactionsparepartentities.SalesOrder{SalesOrderSystemNumber: SalesOrderId}).
		First(&soEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return soEntity, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "sales order with that id is not found",
			}
		}
		return soEntity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get sales order entity",
		}
	}
	if payload.SalesOrderRemark != "" {
		soEntity.Remark = payload.SalesOrderRemark
	}
	if payload.DownPaymentAmount != 0 {
		soEntity.DownPaymentAmount = payload.DownPaymentAmount
	}
	//save header
	err = db.Save(&soEntity).Error
	if err != nil {
		return soEntity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to update sales order",
		}
	}
	return soEntity, nil
}

// uspg atsalesorder0 Update option 1
func (r *SalesOrderRepositoryImpl) SubmitSalesOrderHeader(db *gorm.DB, salesOrderId int) (bool, *exceptions.BaseErrorResponse) {
	//BEGIN
	//IF NOT EXISTS(SELECT TOP 1 1
	//FROM atSalesOrderLog_KIANONIMG A
	//WHERE A.SO_SYS_NO = @So_Sys_No
	//AND A.PROCESS_CODE = '02')--SO_PROCESS_CODE_DEALERSUBMIT
	//BEGIN
	//RAISERROR('Submit Failed, Non IMG dealers haven''t submitted this SPSE yet.',16,1)

	//SET @Src_Doc  = dbo.getVariableValue('SRC_DOC_SP_SO')
	//not yet dev
	//var srcDoc = "SPSO"
	var discountOriginal float64
	var discountRequest float64

	//SELECT @Disc_Ori = SUM(ISNULL(DISC_AMOUNT,0)),
	//@Disc_Req = SUM(ISNULL(DISC_REQ_AMOUNT,0))
	//FROM atSalesOrder1
	//get discountOriginal and discountRequest
	err := db.Model(&transactionsparepartentities.SalesOrderDetail{}).
		Where(transactionsparepartentities.SalesOrderDetail{SalesOrderSystemNumber: salesOrderId}).
		Select(`SUM(ISNULL(discount_amount,0)),SUM(ISNULL(discount_request_amount,0))`).
		Row().
		Scan(&discountOriginal, &discountRequest)
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get discount amount and discount request amount",
		}
	}
	if discountOriginal != discountRequest {
		//SET @Disc_Req = (SELECT ROUND(ISNULL(SUM(CASE WHEN DISC_REQ_AMOUNT > 0
		//THEN ISNULL(DISC_REQ_AMOUNT,0) * ISNULL(QTY_DEMAND,0)
		//ELSE ISNULL(DISC_AMOUNT,0) * ISNULL(QTY_DEMAND,0)
		//END ),0),0)
		//FROM atSalesOrder1 WHERE SO_SYS_NO = @So_Sys_No)
		err = db.Model(&transactionsparepartentities.SalesOrderDetail{}).
			Select(`ROUND(ISNULL(SUM(CASE WHEN discount_request_amount > 0
											   THEN ISNULL(discount_request_amount,0) * ISNULL(quantity_demand,0)
											   ELSE ISNULL(discount_amount,0) * ISNULL(quantity_demand,0)
											   END ),0),0)`).
			Where(transactionsparepartentities.SalesOrderDetail{SalesOrderSystemNumber: salesOrderId}).
			Scan(&discountRequest).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed to get discount request amount",
			}
		}
	} else {
		discountRequest = 0
	}
	//salesOrderEntities
	salesOrderEntities := transactionsparepartentities.SalesOrder{}
	err = db.Model(&salesOrderEntities).
		Where(transactionsparepartentities.SalesOrder{SalesOrderSystemNumber: salesOrderId}).
		First(&salesOrderEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "sales order data with that id is not found please check input",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get sales order data",
		}
	}
	//@Disc
	//@Disc = ISNULL(SO.ADD_DISC_AMOUNT,0) + ISNULL(@Disc_Req,0) ,
	discountApproval := salesOrderEntities.AdditionalDiscountAmount + discountRequest

	//check empty item code
	isExist := false
	err = db.Model(&transactionsparepartentities.SalesOrderDetail{}).
		Select("1").
		Where("ISNULL(item_id,0) = 0").
		Scan(&isExist).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to check item code empty",
		}
	}
	if isExist {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("please insert all item code first"),
			Message:    "Please insert all item code first",
		}
	}
	//--==SET VALUE FOR CALCULATE==--
	var SalesOrderSubTotal float64
	var SalesOrderTotalDiscount float64
	var SalesOrderTotal float64
	var SalesOrderAdditionalDiscountAmount float64
	var SalesOrderTotalVat float64
	var SalesOrderTotalAfterVat float64
	//var moduleCode = "SP"
	//var ApproveCode = "SO"
	var ApprovalRequestNumber int
	itemSubstituteType := masteritementities.ItemSubstituteType{}
	err = db.Model(&itemSubstituteType).Where(masteritementities.ItemSubstituteType{ItemSubstituteTypeCode: "S"}).
		Scan(&itemSubstituteType).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get item substitute type",
		}
	}

	err = db.Model(&transactionsparepartentities.SalesOrderDetail{}).
		Where(transactionsparepartentities.SalesOrderDetail{SalesOrderSystemNumber: salesOrderId}).
		Select(`
				COALESCE(SUM(CASE WHEN item_substitute_type_id = ? THEN 0 ELSE COALESCE(price, 0) * COALESCE(quantity_demand, 0) END), 0)
				`, itemSubstituteType.ItemSubstituteTypeId).
		Scan(&SalesOrderSubTotal).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get sub total from sales order detail",
			Err:        err,
		}
	}

	//round to nearest integer
	SalesOrderSubTotal = math.Round(SalesOrderSubTotal)

	//calculate total discount
	err = db.Model(&transactionsparepartentities.SalesOrderDetail{}).
		Where(transactionsparepartentities.SalesOrderDetail{SalesOrderSystemNumber: salesOrderId}).
		Select(`
				COALESCE(SUM(CASE WHEN COALESCE(discount_request_amount,0) > 0 THEN 
				COALESCE(discount_request_amount, 0) * COALESCE(quantity_demand, 0) 
				ELSE 
				COALESCE(discount_amount,0) * COALESCE(quantity_demand,0)
				END), 0)
				`).
		Scan(&SalesOrderTotalDiscount).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get sub total from sales order detail",
			Err:        err,
		}
	}
	//rounding
	SalesOrderTotalDiscount = math.Round(SalesOrderTotalDiscount)
	//SET @Total = @Sub_Total - @Total_Disc
	SalesOrderTotal = SalesOrderSubTotal - SalesOrderTotalDiscount
	//SET @Add_Disc_Amount = (@Total * ISNULL(@Add_Disc_Percent,0)) / 100

	SalesOrderAdditionalDiscountAmount = SalesOrderTotal * (salesOrderEntities.AdditionalDiscountPercentage / 100)
	//rounding
	SalesOrderAdditionalDiscountAmount = math.Round(SalesOrderAdditionalDiscountAmount)
	//notes  math.Round(value/10) * 10 -> sama dengan round(...,-1)
	SalesOrderTotalVat = (SalesOrderTotal - SalesOrderAdditionalDiscountAmount) * (math.Round(SalesOrderTotalVat))

	SalesOrderTotalAfterVat = (SalesOrderTotal - SalesOrderAdditionalDiscountAmount) + SalesOrderTotalVat

	if salesOrderEntities.DownPaymentPaidAmount > SalesOrderTotalAfterVat {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "DP Amount cannot larger than total sales order",
		}
	}
	companyCodeNmdi, companyCodeNmdiErr := generalserviceapiutils.GetCompanyDataByCode("3125098")
	if companyCodeNmdiErr != nil {
		return false, companyCodeNmdiErr
	}

	if companyCodeNmdi.CompanyId == salesOrderEntities.CompanyID {
		//check defect item
		itemMaster := masteritementities.Item{}
		err = db.Table("trx_sales_order_detail A").
			Joins("INNER JOIN mtr_item B ON A.item_id = B.item_id").
			Select("B.*").
			Where("A.sales_order_system_number = ? AND B.is_technical_defect = 1", salesOrderId).Scan(&itemMaster).
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed to cek ",
			}
		}
		if itemMaster.IsTechnicalDefect {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("this item has technical defect"),
				Message:    fmt.Sprintf("item : %s had technical defect", itemMaster.ItemCode),
			}
		}
		var exists bool
		err = db.Table("trx_sales_order_detail A").
			Joins("INNER JOIN trx_sales_order B ON A.sales_order_system_number = B.sales_order_system_number").
			Select("1").
			Where("A.sales_order_system_number = ?", salesOrderId).
			Where("COALESCE(A.price, 0) = 0").
			Limit(1).
			Scan(&exists).
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed on getting item price",
			}
		}
		if exists {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("item price must not be 0"),
				Message:    "item price must not be 0",
			}
		}
	}

	//updating Header
	salesOrderEntities.Total = SalesOrderTotal
	salesOrderEntities.TotalDiscount = SalesOrderTotalDiscount
	salesOrderEntities.AdditionalDiscountPercentage = SalesOrderAdditionalDiscountAmount
	salesOrderEntities.AdditionalDiscountAmount = SalesOrderAdditionalDiscountAmount
	salesOrderEntities.TotalVAT = SalesOrderTotalVat
	salesOrderEntities.TotalAfterVAT = SalesOrderTotalAfterVat

	//save
	err = db.Save(&salesOrderEntities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to save sales order header data",
		}
	}
	//checking transaction type SU05
	transactionTypeInternal := masterentities.TransactionTypeSalesOrder{}
	err = db.Model(&transactionTypeInternal).Where(masterentities.TransactionTypeSalesOrder{TransactionTypeSalesOrderCode: "SU05"}).
		First(&transactionTypeInternal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
				Message:    "transaction type sales order with code SU05 is not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get transaction type internal",
		}
	}
	//get warehouse PG
	WarehouseGroupPG := masterwarehouseentities.WarehouseGroup{}
	err = db.Model(&WarehouseGroupPG).Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupCode: "PG"}).
		First(&WarehouseGroupPG).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "warehouse group with code PG is not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get warehouse group with code pg",
		}
	}
	salesOrderStatusWaitApproved, ApprovalErr := generalserviceapiutils.GetApprovalStatusByCode(utils.ApprovalWaitApproveCode)
	if ApprovalErr != nil {
		return false, ApprovalErr
	}
	//IF @Vehicle_Brand IN ('NISSAN', 'DATSUN', 'KIA', 'CITROEN') AND @Company_Code <> '3125098'
	brandNissan, brandNissanErr := salesserviceapiutils.GetUnitBrandByCode("NISSAN")
	if brandNissanErr != nil {
		return false, brandNissanErr
	}
	brandDatsun, brandDatsunErr := salesserviceapiutils.GetUnitBrandByCode("DATSUN")
	if brandDatsunErr != nil {
		return false, brandDatsunErr
	}
	brandKIA, brandKIAErr := salesserviceapiutils.GetUnitBrandByCode("KIA")
	if brandKIAErr != nil {
		return false, brandKIAErr
	}
	brandCitroen, brandCITROENErr := salesserviceapiutils.GetUnitBrandByCode("CITROEN")
	if brandCITROENErr != nil {
		return false, brandCITROENErr
	}

	var cekBrandId []int
	cekBrandId = append(cekBrandId, brandNissan.BrandId)
	cekBrandId = append(cekBrandId, brandDatsun.BrandId)
	cekBrandId = append(cekBrandId, brandKIA.BrandId)
	cekBrandId = append(cekBrandId, brandCitroen.BrandId)

	if companyCodeNmdi.CompanyId != salesOrderEntities.CompanyID && !utils.NotInList(cekBrandId, salesOrderEntities.BrandID) {
		ApprovalRequestNumber = 12345678
		salesOrderEntities.ApprovalRequestNumber = ApprovalRequestNumber
		salesOrderEntities.SalesOrderStatusID = salesOrderStatusWaitApproved.ApprovalStatusId //wait approved
		err = db.Save(&salesOrderEntities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed to save sales order entities update approval request number",
			}
		}
		//update sales order detail status to

		err = db.Model(&transactionsparepartentities.SalesOrderDetail{}).
			Update("sales_order_line_status_id", salesOrderStatusWaitApproved.ApprovalStatusId).
			Where("sales_order_system_number = ?", salesOrderId).
			Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed to save sales order entities update line status",
			}
		}
	} else if companyCodeNmdi.CompanyId == salesOrderEntities.CompanyID &&
		salesOrderEntities.TransactionTypeID == transactionTypeInternal.TransactionTypeSalesOrderId &&
		salesOrderEntities.WarehouseGroupID == WarehouseGroupPG.WarehouseGroupId {

		//IF ISNULL((SELECT ATPM_INTERNAL_PURPOSE FROM atSalesOrder0 WHERE SO_SYS_NO = @So_Sys_No),'') = ''
		if salesOrderEntities.AtpmInternalPurpose == "" {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("please choose ATPM Internal Purpose Event"),
				Message:    "Please Choose ATPM Internal Purpose Event!",
			}
		}
		//ApproveCode = "SOIT"
		//EXEC dbo.usp_comApprovalReq_Insert
		//@Company_Code = @Company_Code ,
		//@Approval_Code = @Approv_Code ,
		//@Module_Code = @Module_Code ,
		//@Src_Doc_Type = @Src_Doc ,
		//@Src_Sys_No = @So_Sys_No ,
		//@Src_Doc_No = '' ,
		//@Src_Doc_Line = 0 ,
		//@Src_Doc_Date  = @So_Date ,
		//@Src_Doc_Amount = @Disc ,
		//@Change_No = 0 ,
		//@Creation_User_Id = @Change_User_Id ,
		//@Req_No = @Approval_Req_No OUTPUT
		ApprovalRequestNumber = 123456789
		salesOrderEntities.SalesOrderStatusID = salesOrderStatusWaitApproved.ApprovalStatusId
		salesOrderEntities.ApprovalRequestNumber = ApprovalRequestNumber
		err = db.Save(&salesOrderEntities).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
				Message:    "failed to save sales order entities",
			}
		}
		//NOT YET DEVELOP WAITING FOR APPROVAL REQUEST INSERT
	} else {
		if discountApproval != 0 {
			documentNumberSalesOrder, DocumentNumberError := r.GenerateDocumentNumber(db, salesOrderEntities.SalesOrderSystemNumber)
			if DocumentNumberError != nil {
				return false, DocumentNumberError
			}
			//SET @So_Status = dbo.getVariableValue('APPROVAL_WAITAPPROVED')
			//--SET @Approve_Code = dbo.getVariableValue('APV_CODE_SO')
			//SET @Module_Code = dbo.getVariableValue('MODULE_SP')
			//SET @Approve_Code =  dbo.getApprovalCodeBrand(@Company_Code,@Src_Doc,'',@Vehicle_Brand,@Cpc_Code,'')
			//ApproveCode = "SO"
			//waiting for com approval request from general
			//EXEC dbo.usp_comApprovalReq_Insert
			ApprovalRequestNumber = 12345
			salesOrderEntities.SalesOrderStatusID = salesOrderStatusWaitApproved.ApprovalStatusId
			salesOrderEntities.ApprovalRequestNumber = ApprovalRequestNumber
			salesOrderEntities.SalesOrderDocumentNumber = documentNumberSalesOrder
			err = db.Save(&salesOrderEntities).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
					Message:    "failed to save sales order entities",
				}
			}
		} else {

			approvalApproved, approvalApprovedErr := generalserviceapiutils.GetApprovalStatusByCode(utils.ApprovalApprovedCode)
			if approvalApprovedErr != nil {
				return false, approvalApprovedErr
			}

			//-----------
			//EXEC uspg_gmSrcDoc1_Update

			//-----------------
			salesOrderEntities.SalesOrderStatusID = approvalApproved.ApprovalStatusId
			salesOrderEntities.ApprovalRequestNumber = ApprovalRequestNumber
			err = db.Save(&salesOrderEntities).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
					Message:    "failed to save sales order entities",
				}
			}
		}
	}
	return true, nil
}
func (r *SalesOrderRepositoryImpl) GenerateDocumentNumber(tx *gorm.DB, id int) (string, *exceptions.BaseErrorResponse) {
	var salesOrder transactionsparepartentities.SalesOrder

	// Get the work order based on the work order system number
	err := tx.Model(&transactionsparepartentities.SalesOrder{}).Where("sales_order_system_number = ?", id).First(&salesOrder).Error
	if err != nil {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve sales order from the database: %v", err)}
	}

	if salesOrder.BrandID == 0 {

		return "", &exceptions.BaseErrorResponse{Message: "brand_id is missing in the sales order. Please ensure the sales order has a valid brand_id before generating document number."}
	}

	// Get the last work order based on the work order system number
	var lastWorkOrder transactionsparepartentities.SalesOrder
	err = tx.Model(&transactionsparepartentities.SalesOrder{}).
		Where("brand_id = ?", salesOrder.BrandID).
		Order("sales_order_system_number desc").
		Limit(1).
		Scan(&lastWorkOrder).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve last sales order : %v", err)}
	}

	currentTime := time.Now()
	month := int(currentTime.Month())
	year := currentTime.Year() % 100 // Use last two digits of the year

	// fetch data brand from external api
	brandResponse, brandErr := generalserviceapiutils.GetBrandGenerateDoc(salesOrder.BrandID)
	if brandErr != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch brand data from external service",
			Err:        brandErr.Err,
		}
	}

	// Check if BrandCode is not empty before using it
	if brandResponse.BrandCode == "" {
		return "", &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Message: "Brand code is empty"}
	}

	// Get the initial of the brand code
	brandInitial := brandResponse.BrandCode[0]

	// Handle the case when there is no last work order or the format is invalid
	newDocumentNumber := fmt.Sprintf("SPSO/%c/%02d/%02d/00001", brandInitial, month, year)
	if lastWorkOrder.SalesOrderStatusID != 0 {
		lastWorkOrderDate := lastWorkOrder.SalesOrderDate
		lastWorkOrderYear := lastWorkOrderDate.Year() % 100

		// Check if the last work order is from the same year
		if lastWorkOrderYear == year {
			lastWorkOrderCode := lastWorkOrder.SalesOrderDocumentNumber
			codeParts := strings.Split(lastWorkOrderCode, "/")
			if len(codeParts) == 5 {
				lastWorkOrderNumber, err := strconv.Atoi(codeParts[4])
				if err == nil {
					newWorkOrderNumber := lastWorkOrderNumber + 1
					newDocumentNumber = fmt.Sprintf("SPPR/%c/%02d/%02d/%05d", brandInitial, month, year, newWorkOrderNumber)
				} else {
					log.Printf("Failed to parse last work order code: %v", err)
				}
			} else {
				log.Println("Invalid last work order code format")
			}
		}
	}

	log.Printf("New document number: %s", newDocumentNumber)
	return newDocumentNumber, nil
}

func (r *SalesOrderRepositoryImpl) GetSalesOrderTransactionType(tx *gorm.DB) ([]transactionsparepartpayloads.GetSalesOrderTransactionType, *exceptions.BaseErrorResponse) {

	//get all sales order transaction type
	var SalesOrderTransactionType []transactionsparepartpayloads.GetSalesOrderTransactionType
	err := tx.Model(&masterentities.TransactionTypeSalesOrder{}).
		Scan(&SalesOrderTransactionType).Error
	if err != nil {
		return SalesOrderTransactionType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to get sales order transaction type please check input",
		}
	}
	return SalesOrderTransactionType, nil

}
