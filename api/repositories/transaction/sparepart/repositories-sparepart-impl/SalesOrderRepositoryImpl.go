package transactionsparepartrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type SalesOrderRepositoryImpl struct {
}

func StartSalesOrderRepositoryImpl() transactionsparepartrepository.SalesOrderRepository {
	return &SalesOrderRepositoryImpl{}
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
	CustomerData, CustomerDataErr := generalserviceapiutils.GetCustomerMasterByID(payload.CustomerID)
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
	employee, employeeErr := generalserviceapiutils.GetEmployeeByID(payload.SalesEmployeeID)
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
		GetTaxPercent, GetTaxPercentErr := financeserviceapiutils.GetTaxPercent("PPN", "PPN", payload.SalesOrderDate)
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
	Approval, ApprovalErr := generalserviceapiutils.GetApprovalStatusByID(SoEntities.SalesOrderStatusID)
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
	AdditionalDiscountApproval, AdditionalDiscountErr := generalserviceapiutils.GetApprovalStatusByID(SoEntities.AdditionalDiscountStatusID)
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
		customerData, customerDataErr := generalserviceapiutils.GetCustomerMasterByID(item.CustomerID)
		if customerDataErr != nil {
			return pages, customerDataErr
		}
		//get approval status
		salesOrderApprovalStatus, salesOrderApprovalStatusErr := generalserviceapiutils.GetApprovalStatusByID(item.SalesOrderStatusID)
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

	approvalDraft, approvalDraftId := generalserviceapiutils.GetApprovalStatusByID(entities.SalesOrderStatusID)
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
func (r *SalesOrderRepositoryImpl) InsertSalesOrderDetail(db *gorm.DB, payload transactionsparepartpayloads.SalesOrderDetailInsertPayload) (transactionsparepartentities.SalesOrderDetail, *exceptions.BaseErrorResponse) {
	//entities
	entities := transactionsparepartentities.SalesOrderDetail{}
	isItemExist := false
	err := db.Model(&entities).
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
	return entities, nil

}
