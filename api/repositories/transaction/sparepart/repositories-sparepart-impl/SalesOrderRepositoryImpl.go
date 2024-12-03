package transactionsparepartrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	exceptions "after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type SalesOrderRepositoryImpl struct {
}

func StartSalesOrderRepositoryImpl() transactionsparepartrepository.SalesOrderRepository {
	return &SalesOrderRepositoryImpl{}
}

func (r *SalesOrderRepositoryImpl) GetSalesOrderByID(tx *gorm.DB, Id int) (transactionsparepartpayloads.SalesOrderResponse, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SalesOrder{}
	response := transactionsparepartpayloads.SalesOrderResponse{}

	rows, err := tx.Model(&entities).
		Where("sales_order_system_number = ?", Id).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
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
		return Entities, nil
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
		Where(masterentities.TransactionTypeSalesOrder{TransactionTypeTypeSalesOrderId: payload.TransactionTypeID}).
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

	if TransactionTypeSalesOrder.TransactionTypeTypeSalesOrderCode == "SU05" ||
		TransactionTypeSalesOrder.TransactionTypeTypeSalesOrderCode == "SU06" {
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
	
	return Entities, nil
}
