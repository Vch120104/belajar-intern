package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	mastercampaignmasterentities "after-sales/api/entities/master/campaign_master"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	exceptions "after-sales/api/exceptions"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderRepositoryImpl struct {
}

func OpenWorkOrderRepositoryImpl() transactionworkshoprepository.WorkOrderRepository {
	return &WorkOrderRepositoryImpl{}
}

func (r *WorkOrderRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tableStruct := transactionworkshoppayloads.WorkOrderGetAllRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	var convertedResponses []transactionworkshoppayloads.WorkOrderGetAllResponse

	for rows.Next() {

		var (
			workOrderReq transactionworkshoppayloads.WorkOrderGetAllRequest
			workOrderRes transactionworkshoppayloads.WorkOrderGetAllResponse
		)

		if err := rows.Scan(
			&workOrderReq.WorkOrderSystemNumber,
			&workOrderReq.WorkOrderDocumentNumber,
			&workOrderReq.WorkOrderDate,
			&workOrderReq.WorkOrderTypeId,
			&workOrderReq.ServiceAdvisorId,
			&workOrderReq.BrandId,
			&workOrderReq.ModelId,
			&workOrderReq.VariantId,
			&workOrderReq.ServiceSite,
			&workOrderReq.VehicleId,
			&workOrderReq.CustomerId,
			&workOrderReq.BilltoCustomerId,
			&workOrderReq.StatusId,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		VehicleURL := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(workOrderReq.VehicleId)
		//fmt.Println("Fetching Vehicle data from:", VehicleURL)
		var getVehicleResponse transactionworkshoppayloads.WorkOrderVehicleResponse
		if err := utils.Get(VehicleURL, &getVehicleResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch vehicle data from external service",
				Err:        err,
			}
		}

		CustomerURL := config.EnvConfigs.GeneralServiceUrl + "customer-detail/" + strconv.Itoa(workOrderReq.CustomerId)

		var getCustomerResponse transactionworkshoppayloads.CustomerResponse
		if err := utils.Get(CustomerURL, &getCustomerResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch customer data from external service",
				Err:        err,
			}
		}

		WorkOrderTypeURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-type?work_order_type_id=" + strconv.Itoa(workOrderReq.WorkOrderTypeId)
		//fmt.Println("Fetching Work Order Type data from:", WorkOrderTypeURL)
		var getWorkOrderTypeResponses []transactionworkshoppayloads.WorkOrderTypeResponse // Use slice of WorkOrderTypeResponse
		if err := utils.Get(WorkOrderTypeURL, &getWorkOrderTypeResponses, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch work order type data from external service",
				Err:        err,
			}
		}

		var workOrderTypeName string
		if len(getWorkOrderTypeResponses) > 0 {
			workOrderTypeName = getWorkOrderTypeResponses[0].WorkOrderTypeName
		}

		WorkOrderStatusURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-status?work_order_status_id=" + strconv.Itoa(workOrderReq.StatusId)
		//fmt.Println("Fetching Work Order Status data from:", WorkOrderStatusURL)
		var getWorkOrderStatusResponses []transactionworkshoppayloads.WorkOrderStatusResponse // Use slice of WorkOrderStatusResponse
		if err := utils.Get(WorkOrderStatusURL, &getWorkOrderStatusResponses, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch work order status data from external service",
				Err:        err,
			}
		}
		var workOrderStatusName string
		if len(getWorkOrderStatusResponses) > 0 {
			workOrderStatusName = getWorkOrderStatusResponses[0].WorkOrderStatusName
		}

		workOrderRes = transactionworkshoppayloads.WorkOrderGetAllResponse{
			WorkOrderDocumentNumber: workOrderReq.WorkOrderDocumentNumber,
			WorkOrderSystemNumber:   workOrderReq.WorkOrderSystemNumber,
			WorkOrderDate:           workOrderReq.WorkOrderDate,
			FormattedWorkOrderDate:  workOrderReq.WorkOrderDate.Format("2006-01-02"), // Set formatted date
			WorkOrderTypeId:         workOrderReq.WorkOrderTypeId,
			WorkOrderTypeName:       workOrderTypeName,
			BrandId:                 workOrderReq.BrandId,
			VehicleCode:             getVehicleResponse.VehicleCode,
			VehicleTnkb:             getVehicleResponse.VehicleTnkb,
			ModelId:                 workOrderReq.ModelId,
			VehicleId:               workOrderReq.VehicleId,
			CustomerId:              workOrderReq.CustomerId,
			StatusId:                workOrderReq.StatusId,
			StatusName:              workOrderStatusName,
		}

		convertedResponses = append(convertedResponses, workOrderRes)
	}

	var mapResponses []map[string]interface{}

	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_document_number":  response.WorkOrderDocumentNumber,
			"work_order_system_number":    response.WorkOrderSystemNumber,
			"work_order_date":             response.FormattedWorkOrderDate, // Use formatted date
			"work_order_type_id":          response.WorkOrderTypeId,
			"work_order_type_description": response.WorkOrderTypeName,
			"brand_id":                    response.BrandId,
			"vehicle_id":                  response.VehicleId,
			"vehicle_chassis_number":      response.VehicleCode,
			"vehicle_tnkb":                response.VehicleTnkb,
			"work_order_status_id":        response.StatusId,
			"work_order_status_name":      response.StatusName,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) VehicleLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var responses []transactionworkshoppayloads.WorkOrderLookupResponse

	tableStruct := transactionworkshoppayloads.WorkOrderLookupRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	var convertedResponses []transactionworkshoppayloads.WorkOrderLookupResponse

	for rows.Next() {

		var (
			workOrderReq transactionworkshoppayloads.WorkOrderLookupRequest
			workOrderRes transactionworkshoppayloads.WorkOrderLookupResponse
		)

		if err := rows.Scan(
			&workOrderReq.WorkOrderSystemNumber,
			&workOrderReq.WorkOrderDocumentNumber,
			&workOrderReq.VehicleId,
			&workOrderReq.CustomerId,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		VehicleURL := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(workOrderReq.VehicleId)
		//fmt.Println("Fetching Vehicle data from:", VehicleURL)
		var getVehicleResponse transactionworkshoppayloads.WorkOrderVehicleResponse
		if err := utils.Get(VehicleURL, &getVehicleResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch vehicle data from external service",
				Err:        err,
			}
		}

		CustomerURL := config.EnvConfigs.GeneralServiceUrl + "customer-detail/" + strconv.Itoa(workOrderReq.CustomerId)
		//fmt.Println("Fetching Customer data from:", CustomerURL)
		var getCustomerResponse transactionworkshoppayloads.CustomerResponse
		if err := utils.Get(CustomerURL, &getCustomerResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch customer data from external service",
				Err:        err,
			}
		}

		workOrderRes = transactionworkshoppayloads.WorkOrderLookupResponse{
			WorkOrderDocumentNumber: workOrderRes.WorkOrderDocumentNumber,
			WorkOrderSystemNumber:   workOrderRes.WorkOrderSystemNumber,
			VehicleId:               workOrderRes.VehicleId,
			CustomerId:              workOrderRes.CustomerId,
		}

		convertedResponses = append(convertedResponses, workOrderRes)
	}

	var mapResponses []map[string]interface{}

	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_document_number": response.WorkOrderDocumentNumber,
			"work_order_system_number":   response.WorkOrderSystemNumber,
			"vehicle_id":                 response.VehicleId,
			"customer_id":                response.CustomerId,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil

}

func (r *WorkOrderRepositoryImpl) CampaignLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var entities []mastercampaignmasterentities.CampaignMaster

	query := tx.Model(&mastercampaignmasterentities.CampaignMaster{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve campaign master from the database"}
	}

	var WorkOrderCampaignResponse []map[string]interface{}

	for _, entity := range entities {
		campaignData := make(map[string]interface{})

		campaignData["campaign_id"] = entity.CampaignId
		campaignData["campaign_code"] = entity.CampaignCode
		campaignData["campaign_name"] = entity.CampaignName
		campaignData["campaign_period_from"] = entity.CampaignPeriodFrom
		campaignData["campaign_period_to"] = entity.CampaignPeriodTo

		WorkOrderCampaignResponse = append(WorkOrderCampaignResponse, campaignData)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(WorkOrderCampaignResponse, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderNormalRequest) (bool, *exceptions.BaseErrorResponse) {

	// uspg_wtWorkOrder0_Insert  @Option = 0
	// Default values
	defaultWorkOrderDocumentNumber := ""
	defaultWorkOrderStatusId := 1 // 1:Draft, 2:New, 3:Ready, 4:On Going, 5:Stop, 6:QC Pass, 7:Cancel, 8:Closed
	defaultWorkOrderTypeId := 1   // 1:Normal, 2:Campaign, 3:Affiliated, 4:Repeat Job
	defaultServiceAdvisorId := 1  // set default 1 nanti pass from session FE

	// SET @Price_Code = dbo.getVariableValue('DEFAULT_PRICECODE')

	// Menentukan tipe dokumen (SRC_DOC_TYPE) berdasarkan (PROFIT_CENTER): IF @PROFIT_CENTER = @Profit_Center_GR
	// IF @PROFIT_CENTER = @Profit_Center_GR BEGIN SET @SRC_DOC_TYPE = @SrcDocWs END ELSE IF @PROFIT_CENTER = @Profit_Center_BR BEGIN SET @SRC_DOC_TYPE = @SrcDocBs END
	// if request.WorkOrderProfitCenter == 1 {
	// 	request.WorkOrderProfitCenter = 1
	// } else if request.WorkOrderProfitCenter == 2 {
	// 	request.WorkOrderProfitCenter = 2
	// }

	// Get Kode Mata Uang (CCY_CODE) dan Status Bebas Pajak (TAX_FREE)
	// SELECT @CCY_CODE = CCY_CODE FROM gmRef WHERE COMPANY_CODE = @COMPANY_CODE (ccycode)
	// SELECT @Tax_Free = ISNULL(CT.TAX_FREE, 0) FROM dbo.gmCust0 C LEFT JOIN dbo.gmCustType CT ON C.CUSTOMER_TYPE = CT.CUSTOMER_TYPE WHERE CUSTOMER_CODE = @Bill_Cust_Code

	// Menentukan tarif pajak (VAT_TAX_RATE) berdasarkan apakah perusahaan termasuk dalam Free Trade Zone (FTZ) atau tidak

	// Validate current date
	currentDate := time.Now()
	requestDate := request.WorkOrderArrivalTime.Truncate(24 * time.Hour)

	// Check if the WorkOrderDate is backdate or future date
	if requestDate.Before(currentDate) || requestDate.After(currentDate) {
		request.WorkOrderArrivalTime = currentDate
	}

	// check company session
	if request.CompanyId == 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("parameter has lost session, please refresh the data"),
		}
	}

	// Pengecekan apakah Work Order sudah ada, jika belum maka insert data ke tabel work order

	entities := transactionworkshopentities.WorkOrder{
		// Basic information

		// Default values
		WorkOrderDocumentNumber: defaultWorkOrderDocumentNumber,
		WorkOrderStatusId:       defaultWorkOrderStatusId,
		WorkOrderDate:           &currentDate,
		WorkOrderTypeId:         defaultWorkOrderTypeId,
		ServiceAdvisor:          defaultServiceAdvisorId,

		BrandId:        request.BrandId,
		ModelId:        request.ModelId,
		VariantId:      request.VariantId,
		VehicleId:      request.VehicleId,
		CustomerId:     request.CustomerId,
		BillableToId:   request.BilltoCustomerId,
		FromEra:        request.FromEra,
		QueueNumber:    request.QueueSystemNumber,
		ArrivalTime:    &request.WorkOrderArrivalTime,
		ServiceMileage: request.WorkOrderCurrentMileage,
		Storing:        request.Storing,
		Remark:         request.WorkOrderRemark,
		ProfitCenterId: request.WorkOrderProfitCenter,
		CostCenterId:   request.DealerRepresentativeId,

		//general information
		CampaignId: request.CampaignId,
		CompanyId:  request.CompanyId,

		// Customer contact information
		CPTitlePrefix:           request.Titleprefix,
		ContactPersonName:       request.NameCust,
		ContactPersonPhone:      request.PhoneCust,
		ContactPersonMobile:     request.MobileCust,
		ContactPersonContactVia: request.ContactVia,

		// Work order status and details
		EraNumber:      request.WorkOrderEraNo,
		EraExpiredDate: &request.WorkOrderEraExpiredDate,

		// Insurance details
		InsurancePolicyNumber:    request.WorkOrderInsurancePolicyNo,
		InsuranceExpiredDate:     &request.WorkOrderInsuranceExpiredDate,
		InsuranceClaimNumber:     request.WorkOrderInsuranceClaimNo,
		InsurancePersonInCharge:  request.WorkOrderInsurancePic,
		InsuranceOwnRisk:         &request.WorkOrderInsuranceOwnRisk,
		InsuranceWorkOrderNumber: request.WorkOrderInsuranceWONumber,

		// Estimation and service details
		EstTime:         &request.EstimationDuration,
		CustomerExpress: request.CustomerExpress,
		LeaveCar:        request.LeaveCar,
		CarWash:         request.CarWash,
		PromiseDate:     &request.PromiseDate,
		PromiseTime:     &request.PromiseTime,

		// Additional information
		FSCouponNo: request.FSCouponNo,
		Notes:      request.Notes,
		Suggestion: request.Suggestion,
		DPAmount:   &request.DownpaymentAmount,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}

	}

	return true, nil

	// Memperbarui status pemesanan dan estimasi jika Booking_System_No atau Estim_System_No tidak nol

	// Insert detil work order (WO1) berdasarkan pemesanan (wtBookEstim1_1):

	// Insert detil work order (WO2) berdasarkan tipe work order (Normal, Campaign, Affiliated, Repeat Job):

	// substitusi item jika stok tidak mencukupi

	// Menghitung total biaya work order berdasarkan tipe line item dan melakukan update pada tabel

	// Menghitung total diskon dan PPN, serta memperbarui work order
}

func (r *WorkOrderRepositoryImpl) NewStatus(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse) {
	var statuses []transactionworkshopentities.WorkOrderMasterStatus

	query := utils.ApplyFilter(tx, filter)

	if err := query.Find(&statuses).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order statuses from the database"}
	}
	return statuses, nil
}

func (r *WorkOrderRepositoryImpl) AddStatus(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderStatusRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderMasterStatus{
		WorkOrderStatusCode:        request.WorkOrderStatusCode,
		WorkOrderStatusDescription: request.WorkOrderStatusName,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateStatus(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderStatusRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterStatus
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterStatus{}).Where("work_order_status_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order status from the database"}
	}

	entity.WorkOrderStatusCode = request.WorkOrderStatusCode
	entity.WorkOrderStatusDescription = request.WorkOrderStatusName

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to update work order status"}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteStatus(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterStatus
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterStatus{}).Where("work_order_status_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order status from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete work order status"}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterType

	query := utils.ApplyFilter(tx, filter)

	if err := query.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}
	return types, nil
}

func (r *WorkOrderRepositoryImpl) AddType(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderTypeRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderMasterType{
		WorkOrderTypeCode:        request.WorkOrderTypeCode,
		WorkOrderTypeDescription: request.WorkOrderTypeName,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateType(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderTypeRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterType{}).Where("work_order_type_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}

	entity.WorkOrderTypeCode = request.WorkOrderTypeCode
	entity.WorkOrderTypeDescription = request.WorkOrderTypeName

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to update work order type"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteType(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterType{}).Where("work_order_type_id = ?", id).First(&entity).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete work order type"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewBill(*gorm.DB) ([]transactionworkshoppayloads.WorkOrderBillable, *exceptions.BaseErrorResponse) {
	BillableURL := config.EnvConfigs.GeneralServiceUrl + "billable-to"
	fmt.Println("Fetching Billable data from:", BillableURL)

	var getBillables []transactionworkshoppayloads.WorkOrderBillable
	if err := utils.Get(BillableURL, &getBillables, nil); err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch billable data from external service",
			Err:        err,
		}
	}

	return getBillables, nil
}

func (r *WorkOrderRepositoryImpl) AddBill(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderBillableRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderMasterBillAbleto{
		WorkOrderBillabletoName: request.BillableToName,
		WorkOrderBillabletoCode: request.BillableToCode,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateBill(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderBillableRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterBillAbleto
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterBillAbleto{}).Where("billable_to_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve billable data from the database"}
	}

	entity.WorkOrderBillabletoName = request.BillableToName
	entity.WorkOrderBillabletoCode = request.BillableToCode

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to update billable data"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteBill(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterBillAbleto
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterBillAbleto{}).Where("billable_to_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve billable data from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete billable data"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewDropPoint(*gorm.DB) ([]transactionworkshoppayloads.WorkOrderDropPoint, *exceptions.BaseErrorResponse) {
	DropPointURL := config.EnvConfigs.GeneralServiceUrl + "company-selection-dropdown"
	fmt.Println("Fetching Drop Point data from:", DropPointURL)

	var getDropPoints []transactionworkshoppayloads.WorkOrderDropPoint
	if err := utils.Get(DropPointURL, &getDropPoints, nil); err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch drop point data from external service",
			Err:        err,
		}
	}

	return getDropPoints, nil
}

// func (r *WorkOrderRepositoryImpl) AddDropPoint(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderDropPointRequest) (bool, *exceptions.BaseErrorResponse) {
// 	entities := transactionworkshopentities.WorkOrderDropPoint{
// 		CompanySelectionName:        request.CompanySelectionName,
// 		CompanySelectionDescription: request.CompanySelectionDescription,
// 	}

// 	err := tx.Create(&entities).Error
// 	if err != nil {
// 		return false, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Err:        err,
// 		}
// 	}

// 	return true, nil
// }

// func (r *WorkOrderRepositoryImpl) UpdateDropPoint(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderDropPointRequest) (bool, *exceptions.BaseErrorResponse) {
// 	var entity transactionworkshopentities.WorkOrderDropPoint
// 	err := tx.Model(&transactionworkshopentities.WorkOrderDropPoint{}).Where("company_selection_id = ?", id).First(&entity).Error
// 	if err != nil {
// 		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve drop point data from the database"}
// 	}

// 	entity.CompanySelectionName = request.CompanySelectionName
// 	entity.CompanySelectionDescription = request.CompanySelectionDescription

// 	err = tx.Save(&entity).Error
// 	if err != nil {
// 		return false, &exceptions.BaseErrorResponse{Message: "Failed to update drop point data"}
// 	}

// 	return true, nil
// }

// func (r *WorkOrderRepositoryImpl) DeleteDropPoint(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
// 	var entity transactionworkshopentities.WorkOrderDropPoint
// 	err := tx.Model(&transactionworkshopentities.WorkOrderDropPoint{}).Where("company_selection_id = ?", id).First(&entity).Error
// 	if err != nil {
// 		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve drop point data from the database"}
// 	}

// 	err = tx.Delete(&entity).Error
// 	if err != nil {
// 		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete drop point data"}
// 	}

// 	return true, nil
// }

func (r *WorkOrderRepositoryImpl) NewVehicleBrand(*gorm.DB) ([]transactionworkshoppayloads.WorkOrderVehicleBrand, *exceptions.BaseErrorResponse) {
	VehicleBrandURL := config.EnvConfigs.SalesServiceUrl + "unit-brand-dropdown"
	fmt.Println("Fetching Vehicle Brand data from:", VehicleBrandURL)

	var getVehicleBrands []transactionworkshoppayloads.WorkOrderVehicleBrand
	if err := utils.Get(VehicleBrandURL, &getVehicleBrands, nil); err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch vehicle brand data from external service",
			Err:        err,
		}
	}

	return getVehicleBrands, nil
}

func (r *WorkOrderRepositoryImpl) NewVehicleModel(_ *gorm.DB, brandId int) ([]transactionworkshoppayloads.WorkOrderVehicleModel, *exceptions.BaseErrorResponse) {
	VehicleModelURL := config.EnvConfigs.SalesServiceUrl + "unit-model-dropdown/" + strconv.Itoa(brandId)
	fmt.Println("Fetching Vehicle Model data from:", VehicleModelURL)

	var getVehicleModels []transactionworkshoppayloads.WorkOrderVehicleModel
	if err := utils.Get(VehicleModelURL, &getVehicleModels, nil); err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch vehicle model data from external service",
			Err:        err,
		}
	}

	return getVehicleModels, nil
}

func (r *WorkOrderRepositoryImpl) GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.WorkOrderRequest{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Work order not found"}
		}
		return transactionworkshoppayloads.WorkOrderRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database", Err: err}
	}

	payload := transactionworkshoppayloads.WorkOrderRequest{
		WorkOrderSystemNumber:      entity.WorkOrderSystemNumber,
		WorkOrderDocumentNumber:    entity.WorkOrderDocumentNumber,
		WorkOrderTypeId:            entity.WorkOrderTypeId,
		ServiceAdvisorId:           entity.ServiceAdvisor,
		BrandId:                    entity.BrandId,
		ModelId:                    entity.ModelId,
		ServiceSite:                entity.ServiceSite,
		VehicleId:                  entity.VehicleId,
		CustomerId:                 entity.CustomerId,
		BilltoCustomerId:           entity.BillableToId,
		CampaignId:                 entity.CampaignId,
		AgreementId:                entity.AgreementBodyRepairId,
		BoookingId:                 entity.BookingSystemNumber,
		EstimationId:               entity.EstimationSystemNumber,
		ContractSystemNumber:       entity.ContractServiceSystemNumber,
		QueueSystemNumber:          entity.QueueNumber,
		WorkOrderRemark:            entity.Remark,
		DealerRepresentativeId:     entity.CostCenterId,
		CompanyId:                  entity.CompanyId,
		Titleprefix:                entity.CPTitlePrefix,
		NameCust:                   entity.ContactPersonName,
		PhoneCust:                  entity.ContactPersonPhone,
		MobileCust:                 entity.ContactPersonMobile,
		MobileCustAlternative:      entity.ContactPersonMobileAlternative,
		MobileCustDriver:           entity.ContactPersonMobileDriver,
		ContactVia:                 entity.ContactPersonContactVia,
		WorkOrderInsurancePolicyNo: entity.InsurancePolicyNumber,
		WorkOrderInsuranceClaimNo:  entity.InsuranceClaimNumber,
		WorkOrderInsurancePic:      entity.InsurancePersonInCharge,
		WorkOrderInsuranceWONumber: entity.InsuranceWorkOrderNumber,
	}

	if entity.WorkOrderDate != nil {
		payload.WorkOrderDate = *entity.WorkOrderDate
	}

	if entity.ArrivalTime != nil {
		payload.WorkOrderArrivalTime = *entity.ArrivalTime
	}

	if entity.ServiceMileage != 0 {
		payload.WorkOrderCurrentMileage = entity.ServiceMileage
	}

	if entity.EraExpiredDate != nil {
		payload.WorkOrderEraExpiredDate = *entity.EraExpiredDate
	}

	if entity.InsuranceExpiredDate != nil {
		payload.WorkOrderInsuranceExpiredDate = *entity.InsuranceExpiredDate
	}

	if entity.PromiseDate != nil {
		payload.PromiseDate = *entity.PromiseDate
	}

	if entity.PromiseTime != nil {
		payload.PromiseTime = *entity.PromiseTime
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderNormalSaveRequest, workOrderId int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	entity.BillableToId = request.BilltoCustomerId
	entity.FromEra = request.FromEra
	entity.QueueNumber = request.QueueSystemNumber
	entity.ArrivalTime = &request.WorkOrderArrivalTime
	entity.ServiceMileage = request.WorkOrderCurrentMileage
	entity.Storing = request.Storing
	entity.Remark = request.WorkOrderRemark
	entity.Unregister = request.Unregistered
	entity.ProfitCenterId = request.WorkOrderProfitCenter
	entity.CostCenterId = request.DealerRepresentativeId
	entity.CompanyId = request.CompanyId

	entity.CPTitlePrefix = request.Titleprefix
	entity.ContactPersonName = request.NameCust
	entity.ContactPersonPhone = request.PhoneCust
	entity.ContactPersonMobile = request.MobileCust
	entity.ContactPersonMobileAlternative = request.MobileCustAlternative
	entity.ContactPersonMobileDriver = request.MobileCustDriver
	entity.ContactPersonContactVia = request.ContactVia

	entity.InsuranceCheck = request.WorkOrderInsuranceCheck
	entity.InsurancePolicyNumber = request.WorkOrderInsurancePolicyNo
	entity.InsuranceExpiredDate = &request.WorkOrderInsuranceExpiredDate
	entity.InsuranceClaimNumber = request.WorkOrderInsuranceClaimNo
	entity.InsurancePersonInCharge = request.WorkOrderInsurancePic
	entity.InsuranceOwnRisk = &request.WorkOrderInsuranceOwnRisk
	entity.InsuranceWorkOrderNumber = request.WorkOrderInsuranceWONumber

	//page2
	entity.EstTime = &request.EstimationDuration
	entity.CustomerExpress = request.CustomerExpress
	entity.LeaveCar = request.LeaveCar
	entity.CarWash = request.CarWash
	entity.PromiseDate = &request.PromiseDate
	entity.PromiseTime = &request.PromiseTime
	entity.FSCouponNo = request.FSCouponNo
	entity.Notes = request.Notes
	entity.Suggestion = request.Suggestion
	entity.DPAmount = &request.DownpaymentAmount

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) Void(tx *gorm.DB, workOrderId int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, &exceptions.BaseErrorResponse{Message: "Work order not found"}
		}
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete the work order"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) CloseOrder(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Check if WorkOrderStatusId is equal to 1 (Draft)
	if entity.WorkOrderStatusId == 1 {
		return false, &exceptions.BaseErrorResponse{Message: "Work order cannot be closed because status is draft"}
	}

	// Update the work order status to 8 (Closed)
	entity.WorkOrderStatusId = 8

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to close the work order"}
	}

	return true, nil

}

func (r *WorkOrderRepositoryImpl) GetAllRequest(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrderService
	// Query to retrieve all work order service entities based on the request
	query := tx.Model(&transactionworkshopentities.WorkOrderService{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order service requests from the database"}
	}

	var workOrderServiceResponses []map[string]interface{}

	for _, entity := range entities {
		workOrderServiceData := make(map[string]interface{})

		workOrderServiceData["work_order_service_id"] = entity.WorkOrderServiceId
		workOrderServiceData["work_order_system_number"] = entity.WorkOrderSystemNumber
		workOrderServiceData["work_order_service_remark"] = entity.WorkOrderServiceRemark
		workOrderServiceResponses = append(workOrderServiceResponses, workOrderServiceData)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(workOrderServiceResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) GetRequestById(tx *gorm.DB, id int, IdWorkorder int) (transactionworkshoppayloads.WorkOrderServiceRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderServiceRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order service request from the database"}
	}

	payload := transactionworkshoppayloads.WorkOrderServiceRequest{
		WorkOrderServiceId:     entity.WorkOrderServiceId,
		WorkOrderSystemNumber:  entity.WorkOrderSystemNumber,
		WorkOrderServiceRemark: entity.WorkOrderServiceRemark,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) UpdateRequest(tx *gorm.DB, id int, IdWorkorder int, request transactionworkshoppayloads.WorkOrderServiceRequest) *exceptions.BaseErrorResponse {

	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order service request from the database"}
	}

	entity.WorkOrderServiceRemark = request.WorkOrderServiceRemark

	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order service request"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) AddRequest(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderServiceRequest) (bool, *exceptions.BaseErrorResponse) {

	CurrentTime := time.Now()

	entities := transactionworkshopentities.WorkOrderService{

		WorkOrderSystemNumber:  request.WorkOrderSystemNumber,
		WorkOrderServiceRemark: request.WorkOrderServiceRemark,
		WorkOrderServiceDate:   CurrentTime,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteRequest(tx *gorm.DB, id int, IdWorkorder int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete work order service request from the database"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) GetAllVehicleService(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrderServiceVehicle

	query := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order service vehicle requests from the database", Err: err}
	}

	if len(entities) == 0 {

		return []map[string]interface{}{}, 0, 0, nil
	}

	var workOrderServiceVehicleResponses []map[string]interface{}

	for _, entity := range entities {
		workOrderServiceVehicleData := make(map[string]interface{})

		workOrderServiceVehicleData["work_order_system_number"] = entity.WorkOrderSystemNumber
		workOrderServiceVehicleData["work_order_vehicle_date"] = entity.WorkOrderVehicleDate
		workOrderServiceVehicleData["work_order_vehicle_remark"] = entity.WorkOrderVehicleRemark
		workOrderServiceVehicleResponses = append(workOrderServiceVehicleResponses, workOrderServiceVehicleData)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(workOrderServiceVehicleResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) GetVehicleServiceById(tx *gorm.DB, id int, IdWorkorder int) (transactionworkshoppayloads.WorkOrderServiceVehicleRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderServiceVehicle
	err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderServiceVehicleRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order service vehicle request from the database"}
	}

	payload := transactionworkshoppayloads.WorkOrderServiceVehicleRequest{
		WorkOrderSystemNumber:  entity.WorkOrderSystemNumber,
		WorkOrderVehicleDate:   entity.WorkOrderVehicleDate,
		WorkOrderVehicleRemark: entity.WorkOrderVehicleRemark,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) UpdateVehicleService(tx *gorm.DB, id int, IdWorkorder int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) *exceptions.BaseErrorResponse {

	var entity transactionworkshopentities.WorkOrderServiceVehicle
	err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order service request from the database"}
	}

	entity.WorkOrderVehicleDate = request.WorkOrderVehicleDate
	entity.WorkOrderVehicleRemark = request.WorkOrderVehicleRemark

	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order service request"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) AddVehicleService(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (bool, *exceptions.BaseErrorResponse) {

	CurrentDate := time.Now()

	entities := transactionworkshopentities.WorkOrderServiceVehicle{

		WorkOrderSystemNumber:  request.WorkOrderSystemNumber,
		WorkOrderVehicleDate:   CurrentDate,
		WorkOrderVehicleRemark: request.WorkOrderVehicleRemark,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteVehicleService(tx *gorm.DB, id int, IdWorkorder int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderServiceVehicle
	err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete work order service request from the database"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) GenerateDocumentNumber(tx *gorm.DB, workOrderId int) (string, *exceptions.BaseErrorResponse) {
	var workOrder transactionworkshopentities.WorkOrder

	// Get the work order based on the work order system number
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&workOrder).Error
	if err != nil {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve work order from the database: %v", err)}
	}

	if workOrder.BrandId == 0 {

		return "", &exceptions.BaseErrorResponse{Message: "brand_id is missing in the work order. Please ensure the work order has a valid brand_id before generating document number."}
	}

	// Get the last work order based on the work order system number
	var lastWorkOrder transactionworkshopentities.WorkOrder
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("brand_id = ?", workOrder.BrandId).
		Order("work_order_document_number desc").
		First(&lastWorkOrder).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve last work order: %v", err)}
	}

	currentTime := time.Now()
	month := int(currentTime.Month())
	year := currentTime.Year() % 100 // Use last two digits of the year

	brandInitial := workOrder.BrandId

	// Handle the case when there is no last work order or the format is invalid
	newDocumentNumber := fmt.Sprintf("WSWO/%d/%02d/%02d/00001", brandInitial, month, year)
	if lastWorkOrder.WorkOrderSystemNumber != 0 {
		lastWorkOrderDate := lastWorkOrder.WorkOrderDate
		lastWorkOrderYear := lastWorkOrderDate.Year() % 100

		// Check if the last work order is from the same year
		if lastWorkOrderYear == year {
			lastWorkOrderCode := lastWorkOrder.WorkOrderDocumentNumber
			codeParts := strings.Split(lastWorkOrderCode, "/")
			if len(codeParts) == 5 {
				lastWorkOrderNumber, err := strconv.Atoi(codeParts[4])
				if err == nil {
					newWorkOrderNumber := lastWorkOrderNumber + 1
					newDocumentNumber = fmt.Sprintf("WSWO/%d/%02d/%02d/%05d", brandInitial, month, year, newWorkOrderNumber)
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

func (r *WorkOrderRepositoryImpl) Submit(tx *gorm.DB, workOrderId int) (bool, string, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", &exceptions.BaseErrorResponse{Message: "No work order data found"}
		}
		return false, "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve work order from the database: %v", err)}
	}

	if entity.WorkOrderDocumentNumber == "" && entity.WorkOrderStatusId == 1 {
		//Generate new document number
		newDocumentNumber, genErr := r.GenerateDocumentNumber(tx, entity.WorkOrderSystemNumber)
		if genErr != nil {
			return false, "", genErr
		}
		//newDocumentNumber := "WSWO/1/21/21/00001"

		entity.WorkOrderDocumentNumber = newDocumentNumber

		// Update work order status to 2 (New Submitted)
		entity.WorkOrderStatusId = 2

		err = tx.Save(&entity).Error
		if err != nil {
			return false, "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to submit the work order: %v", err)}
		}

		return true, newDocumentNumber, nil
	} else {

		return false, "", &exceptions.BaseErrorResponse{Message: "Document number has already been generated"}
	}
}

func (r *WorkOrderRepositoryImpl) GetAllDetailWorkOrder(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrderDetail

	query := tx.Model(&transactionworkshopentities.WorkOrderDetail{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order detail from the database"}
	}

	var workOrderDetailResponses []map[string]interface{}

	for _, entity := range entities {
		workOrderDetailData := make(map[string]interface{})

		workOrderDetailData["work_order_detail_id"] = entity.WorkOrderDetailId
		workOrderDetailData["work_order_system_number"] = entity.WorkOrderSystemNumber
		workOrderDetailData["line_type_id"] = entity.LineTypeId
		workOrderDetailData["transaction_type_id"] = entity.TransactionTypeId
		workOrderDetailData["job_type_id"] = entity.JobTypeId
		workOrderDetailData["description"] = entity.Description
		workOrderDetailData["frt_quantity"] = entity.FrtQuantity
		workOrderDetailData["supply_quantity"] = entity.SupplyQuantity
		workOrderDetailData["price_list_id"] = entity.PriceListId

		workOrderDetailResponses = append(workOrderDetailResponses, workOrderDetailData)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(workOrderDetailResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) GetDetailByIdWorkOrder(tx *gorm.DB, id int, IdWorkorder int) (transactionworkshoppayloads.WorkOrderDetailRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderDetail
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_detail_id = ?", id, IdWorkorder).
		First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderDetailRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order detail from the database"}
	}

	payload := transactionworkshoppayloads.WorkOrderDetailRequest{
		WorkOrderDetailId:     entity.WorkOrderDetailId,
		WorkOrderSystemNumber: entity.WorkOrderSystemNumber,
		LineTypeId:            entity.LineTypeId,
		TransactionTypeId:     entity.TransactionTypeId,
		JobTypeId:             entity.JobTypeId,
		FrtQuantity:           entity.FrtQuantity,
		SupplyQuantity:        entity.SupplyQuantity,
		PriceListId:           entity.PriceListId,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) CalculateWorkOrderTotal(tx *gorm.DB, workOrderSystemNumber int, lineTypeId int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	const (
		LineTypePackage            = 0 // Package Bodyshop
		LineTypeOperation          = 1 // Operation
		LineTypeSparePart          = 2 // Spare Part
		LineTypeOil                = 3 // Oil
		LineTypeMaterial           = 4 // Material
		LineTypeFee                = 5 // Fee
		LineTypeAccessories        = 6 // Accessories
		LineTypeConsumableMaterial = 7 // Consumable Material
		LineTypeSublet             = 8 // Sublet
		LineTypeSouvenir           = 9 // Souvenir
	)

	type Result struct {
		TotalPackage            float32
		TotalOperation          float32
		TotalSparePart          float32
		TotalOil                float32
		TotalMaterial           float32
		TotalFee                float32
		TotalAccessories        float32
		TotalConsumableMaterial float32
		TotalSublet             float32
		TotalSouvenir           float32
	}

	var result Result

	// Calculate totals for each line type
	err := tx.Raw(`
		SELECT
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0), 0) ELSE 0 END) AS total_package,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_operation,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_spare_part,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_oil,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_material,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_fee,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_accessories,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_consumable_material,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_sublet,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_souvenir
		FROM trx_work_order_detail
		WHERE work_order_system_number = ?`,
		LineTypePackage,
		LineTypeOperation,
		LineTypeSparePart,
		LineTypeOil,
		LineTypeMaterial,
		LineTypeFee,
		LineTypeAccessories,
		LineTypeConsumableMaterial,
		LineTypeSublet,
		LineTypeSouvenir,
		workOrderSystemNumber).Scan(&result).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to calculate totals: %v", err)}
	}

	// Calculate grand total
	grandTotal := result.TotalPackage + result.TotalOperation + result.TotalSparePart + result.TotalOil + result.TotalMaterial + result.TotalFee + result.TotalAccessories + result.TotalConsumableMaterial + result.TotalSublet + result.TotalSouvenir

	// Update Work Order with the calculated totals
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", workOrderSystemNumber).
		Updates(map[string]interface{}{
			"total_package":             result.TotalPackage,
			"total_operation":           result.TotalOperation,
			"total_part":                result.TotalSparePart,
			"total_oil":                 result.TotalOil,
			"total_material":            result.TotalMaterial,
			"total_price_accessories":   result.TotalAccessories,
			"total_consumable_material": result.TotalConsumableMaterial,
			"total_sublet":              result.TotalSublet,
			"total":                     grandTotal,
			//"total_fee":                 result.TotalFee,
			//"total_souvenir":            result.TotalSouvenir,
		}).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to update work order: %v", err)}
	}

	// Prepare response
	workOrderDetailResponses := []map[string]interface{}{
		{"total_package": result.TotalPackage},
		{"total_operation": result.TotalOperation},
		{"total_part": result.TotalSparePart},
		{"total_oil": result.TotalOil},
		{"total_material": result.TotalMaterial},
		{"total_price_accessories": result.TotalAccessories},
		{"total_consumable_material": result.TotalConsumableMaterial},
		{"total_sublet": result.TotalSublet},
		{"total": grandTotal},
	}

	return workOrderDetailResponses, nil
}

func (r *WorkOrderRepositoryImpl) AddDetailWorkOrder(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderDetailRequest) (bool, *exceptions.BaseErrorResponse) {

	entities := transactionworkshopentities.WorkOrderDetail{

		WorkOrderSystemNumber:              request.WorkOrderSystemNumber,
		LineTypeId:                         request.LineTypeId,
		TransactionTypeId:                  request.TransactionTypeId,
		JobTypeId:                          request.JobTypeId,
		WarehouseId:                        request.WarehouseId,
		ItemId:                             request.ItemId,
		FrtQuantity:                        request.FrtQuantity,
		SupplyQuantity:                     request.SupplyQuantity,
		PriceListId:                        request.PriceListId,
		OperationItemDiscountRequestAmount: request.ProposedPrice,
		OperationItemPrice:                 request.OperationItemPrice,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Call CalculateWorkOrderTotal to update the totals in trx_work_order
	_, calcErr := r.CalculateWorkOrderTotal(tx, id, request.LineTypeId)
	if calcErr != nil {
		return false, calcErr
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateDetailWorkOrder(tx *gorm.DB, IdWorkorder int, id int, request transactionworkshoppayloads.WorkOrderDetailRequest) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderDetail
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_detail_id = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order detail from the database"}
	}

	entity.LineTypeId = request.LineTypeId
	entity.TransactionTypeId = request.TransactionTypeId
	entity.JobTypeId = request.JobTypeId
	entity.WarehouseId = request.WarehouseId
	entity.ItemId = request.ItemId
	entity.FrtQuantity = request.FrtQuantity
	entity.SupplyQuantity = request.SupplyQuantity
	entity.PriceListId = request.PriceListId
	entity.OperationItemDiscountRequestAmount = request.ProposedPrice
	entity.OperationItemPrice = request.OperationItemPrice

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order detail"}
	}

	// Call CalculateWorkOrderTotal to update the totals in trx_work_order
	_, calcErr := r.CalculateWorkOrderTotal(tx, id, request.LineTypeId)
	if calcErr != nil {
		return false, calcErr
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteDetailWorkOrder(tx *gorm.DB, id int, IdWorkorder int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderDetail
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_detail_id = ?", id, IdWorkorder).
		Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete work order detail from the database"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewBooking(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse) {

	entities := transactionworkshopentities.WorkOrder{

		WorkOrderSystemNumber: request.WorkOrderSystemNumber,
		BookingSystemNumber:   request.BoookingId,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the work order booking"}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) GetAllBooking(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrder

	query := tx.Model(&transactionworkshopentities.WorkOrder{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order booking from the database"}
	}

	var workOrderBookingResponses []map[string]interface{}

	for _, entity := range entities {
		workOrderBookingData := make(map[string]interface{})

		workOrderBookingData["work_order_system_number"] = entity.WorkOrderSystemNumber
		workOrderBookingData["booking_system_number"] = entity.BookingSystemNumber
		workOrderBookingData["service_request_system_number"] = entity.ServiceRequestSystemNumber
		workOrderBookingData["brand_id"] = entity.BrandId
		workOrderBookingData["model_id"] = entity.ModelId
		workOrderBookingData["vehicle_id"] = entity.VehicleId

		workOrderBookingResponses = append(workOrderBookingResponses, workOrderBookingData)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(workOrderBookingResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) GetBookingById(tx *gorm.DB, IdWorkorder int, id int) (transactionworkshoppayloads.WorkOrderBookingRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND booking_system_number = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderBookingRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order booking from the database"}
	}

	payload := transactionworkshoppayloads.WorkOrderBookingRequest{
		WorkOrderSystemNumber: entity.WorkOrderSystemNumber,
		BoookingId:            entity.BookingSystemNumber,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) SaveBooking(tx *gorm.DB, IdWorkorder int, id int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND booking_system_number = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order booking from the database"}
	}

	entity.WorkOrderSystemNumber = request.WorkOrderSystemNumber
	entity.BookingSystemNumber = request.BoookingId

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order booking"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) SubmitBooking(tx *gorm.DB, IdWorkorder int, id int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND booking_system_number = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order booking from the database"}
	}

	// Update the work order booking status to 2 (Submitted)
	entity.WorkOrderStatusId = 2

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to submit the work order booking"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) VoidBooking(tx *gorm.DB, IdWorkorder int, id int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND booking_system_number = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order booking from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete the work order booking"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) CloseBooking(tx *gorm.DB, IdWorkorder int, id int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND booking_system_number = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order booking from the database"}
	}

	// Update the work order booking status to 3 (Closed)
	entity.WorkOrderStatusId = 3

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to close the work order booking"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) GetAllAffiliated(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrder

	query := tx.Model(&transactionworkshopentities.WorkOrder{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order affiliate from the database"}
	}

	var workOrderAffiliateResponses []map[string]interface{}

	for _, entity := range entities {
		workOrderAffiliateData := make(map[string]interface{})

		workOrderAffiliateData["work_order_system_number"] = entity.WorkOrderSystemNumber

		workOrderAffiliateResponses = append(workOrderAffiliateResponses, workOrderAffiliateData)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(workOrderAffiliateResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) GetAffiliatedById(tx *gorm.DB, IdWorkorder int, id int) (transactionworkshoppayloads.WorkOrderAffiliatedRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND affiliate_id = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderAffiliatedRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order affiliate from the database"}
	}

	payload := transactionworkshoppayloads.WorkOrderAffiliatedRequest{
		WorkOrderSystemNumber: entity.WorkOrderSystemNumber,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) NewAffiliated(tx *gorm.DB, IdWorkorder int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrder{

		WorkOrderSystemNumber: request.WorkOrderSystemNumber,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the work order affiliate"}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) SaveAffiliated(tx *gorm.DB, IdWorkorder int, id int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND affiliate_id = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order affiliate from the database"}
	}

	entity.WorkOrderSystemNumber = request.WorkOrderSystemNumber

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order affiliate"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) VoidAffiliated(tx *gorm.DB, IdWorkorder int, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND affiliate_id = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order affiliate from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete the work order affiliate"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) CloseAffiliated(tx *gorm.DB, IdWorkorder int, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND affiliate_id = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order affiliate from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to close the work order affiliate"}
	}

	return true, nil
}
