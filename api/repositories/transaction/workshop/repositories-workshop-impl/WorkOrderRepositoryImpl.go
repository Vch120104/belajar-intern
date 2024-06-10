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
	// Define table struct
	tableStruct := transactionworkshoppayloads.WorkOrderGetAllRequest{}

	// Define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// Apply filters
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Execute query
	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	// Define a slice to hold WorkOrderGetAllResponse
	var convertedResponses []transactionworkshoppayloads.WorkOrderGetAllResponse

	// Iterate over rows
	for rows.Next() {
		// Define variables to hold row data
		var (
			workOrderReq transactionworkshoppayloads.WorkOrderGetAllRequest
			workOrderRes transactionworkshoppayloads.WorkOrderGetAllResponse
		)

		// Scan the row into WorkOrderGetAllRequest struct
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

		// Fetch vehicle data from external service
		VehicleURL := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(workOrderReq.VehicleId)
		fmt.Println("Fetching Vehicle data from:", VehicleURL)
		var getVehicleResponse transactionworkshoppayloads.WorkOrderVehicleResponse
		if err := utils.Get(VehicleURL, &getVehicleResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch vehicle data from external service",
				Err:        err,
			}
		}

		// Fetch Customer data from external service
		CustomerURL := config.EnvConfigs.GeneralServiceUrl + "customer-detail/" + strconv.Itoa(workOrderReq.CustomerId)
		fmt.Println("Fetching Customer data from:", CustomerURL)
		var getCustomerResponse transactionworkshoppayloads.CustomerResponse
		if err := utils.Get(CustomerURL, &getCustomerResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch customer data from external service",
				Err:        err,
			}
		}

		// Fetch type of work order
		WorkOrderTypeURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-type?work_order_type_id=" + strconv.Itoa(workOrderReq.WorkOrderTypeId)
		fmt.Println("Fetching Work Order Type data from:", WorkOrderTypeURL)
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

		// fetch status of work order
		WorkOrderStatusURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-status?work_order_status_id=" + strconv.Itoa(workOrderReq.StatusId)
		fmt.Println("Fetching Work Order Status data from:", WorkOrderStatusURL)
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

		// Create WorkOrderGetAllResponse
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

		// Append WorkOrderResponse to the slice
		convertedResponses = append(convertedResponses, workOrderRes)
	}

	// Define a slice to hold map responses
	var mapResponses []map[string]interface{}

	// Iterate over convertedResponses and convert them to maps
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
			"customer_id":                 response.CustomerId,
			"work_order_status_id":        response.StatusId,
			"work_order_status_name":      response.StatusName,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) VehicleLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	// Define a slice to hold WorkOrderLookupResponse responses
	var responses []transactionworkshoppayloads.WorkOrderLookupResponse

	// Define table struct
	tableStruct := transactionworkshoppayloads.WorkOrderLookupRequest{}

	// Define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// Apply filters
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Execute query
	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	// Define a slice to hold WorkOrderLookupResponse
	var convertedResponses []transactionworkshoppayloads.WorkOrderLookupResponse

	// Iterate over rows
	for rows.Next() {
		// Define variables to hold row data
		var (
			workOrderReq transactionworkshoppayloads.WorkOrderLookupRequest
			workOrderRes transactionworkshoppayloads.WorkOrderLookupResponse
		)

		// Scan the row into WorkOrderLookupRequest struct
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

		// Fetch vehicle data from external service
		VehicleURL := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(workOrderReq.VehicleId)
		fmt.Println("Fetching Vehicle data from:", VehicleURL)
		var getVehicleResponse transactionworkshoppayloads.WorkOrderVehicleResponse
		if err := utils.Get(VehicleURL, &getVehicleResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch vehicle data from external service",
				Err:        err,
			}
		}

		// Fetch Customer data from external service
		CustomerURL := config.EnvConfigs.GeneralServiceUrl + "customer-detail/" + strconv.Itoa(workOrderReq.CustomerId)
		fmt.Println("Fetching Customer data from:", CustomerURL)
		var getCustomerResponse transactionworkshoppayloads.CustomerResponse
		if err := utils.Get(CustomerURL, &getCustomerResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch customer data from external service",
				Err:        err,
			}
		}
		// Create WorkOrderLookupResponse
		workOrderRes = transactionworkshoppayloads.WorkOrderLookupResponse{
			WorkOrderDocumentNumber: workOrderRes.WorkOrderDocumentNumber,
			WorkOrderSystemNumber:   workOrderRes.WorkOrderSystemNumber,
			VehicleId:               workOrderRes.VehicleId,
			CustomerId:              workOrderRes.CustomerId,
		}
		// Append WorkOrderResponse to the slice
		convertedResponses = append(convertedResponses, workOrderRes)
	}

	// Define a slice to hold map responses
	var mapResponses []map[string]interface{}

	// Iterate over convertedResponses and convert them to maps
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_document_number": response.WorkOrderDocumentNumber,
			"work_order_system_number":   response.WorkOrderSystemNumber,
			"vehicle_id":                 response.VehicleId,
			"customer_id":                response.CustomerId,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil

}

func (r *WorkOrderRepositoryImpl) CampaignLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var entities []mastercampaignmasterentities.CampaignMaster
	// Query to retrieve all work order entities based on the request
	query := tx.Model(&mastercampaignmasterentities.CampaignMaster{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve campaign master from the database"}
	}

	var WorkOrderCampaignResponse []map[string]interface{}

	// Loop through each entity and copy its data to the response
	for _, entity := range entities {
		campaignData := make(map[string]interface{})
		// Copy data from entity to response
		campaignData["campaign_id"] = entity.CampaignId
		campaignData["campaign_code"] = entity.CampaignCode
		campaignData["campaign_name"] = entity.CampaignName
		campaignData["campaign_period_from"] = entity.CampaignPeriodFrom
		campaignData["campaign_period_to"] = entity.CampaignPeriodTo

		WorkOrderCampaignResponse = append(WorkOrderCampaignResponse, campaignData)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(WorkOrderCampaignResponse, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptions.BaseErrorResponse) {
	// Create a new instance of WorkOrderRequest
	entities := transactionworkshopentities.WorkOrder{
		// Assign fields from request
		// Basic information
		WorkOrderSystemNumber:   request.WorkOrderSystemNumber,
		WorkOrderDocumentNumber: request.WorkOrderDocumentNumber,
		WorkOrderStatusId:       request.WorkOrderStatusId,
		WorkOrderDate:           &request.WorkOrderDate,
		WorkOrderTypeId:         request.WorkOrderTypeId,
		BrandId:                 request.BrandId,
		ServiceAdvisor:          request.ServiceAdvisorId,
		ModelId:                 request.ModelId,
		VariantId:               request.VariantId,
		VehicleId:               request.VehicleId,
		CustomerId:              request.CustomerId,
		BillableToId:            request.BilltoCustomerId,
		FromEra:                 request.FromEra,
		QueueNumber:             request.QueueSystemNumber,
		ArrivalTime:             &request.WorkOrderArrivalTime,
		ServiceMileage:          request.WorkOrderCurrentMileage,
		Storing:                 request.Storing,
		Remark:                  request.WorkOrderRemark,
		ProfitCenterId:          request.WorkOrderProfitCenter,
		CostCenterId:            request.DealerRepresentativeId,

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

	// Save the work order
	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewStatus(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse) {
	var statuses []transactionworkshopentities.WorkOrderMasterStatus

	// Apply filter to the query
	query := utils.ApplyFilter(tx, filter)

	// Execute the query and check for errors
	if err := query.Find(&statuses).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order statuses from the database"}
	}
	return statuses, nil
}

func (r *WorkOrderRepositoryImpl) NewType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterType
	// Apply filter to the query
	query := utils.ApplyFilter(tx, filter)

	if err := query.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}
	return types, nil
}

func (r *WorkOrderRepositoryImpl) NewBill(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderBillable, *exceptions.BaseErrorResponse) {
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

func (r *WorkOrderRepositoryImpl) NewDropPoint(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderDropPoint, *exceptions.BaseErrorResponse) {
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

func (r *WorkOrderRepositoryImpl) NewVehicleBrand(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderVehicleBrand, *exceptions.BaseErrorResponse) {
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

func (r *WorkOrderRepositoryImpl) NewVehicleModel(tx *gorm.DB, brandId int) ([]transactionworkshoppayloads.WorkOrderVehicleModel, *exceptions.BaseErrorResponse) {
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
		return transactionworkshoppayloads.WorkOrderRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Convert entity to payload
	payload := transactionworkshoppayloads.WorkOrderRequest{}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest, workOrderId int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Update the work order
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

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) Void(tx *gorm.DB, workOrderId int) (bool, *exceptions.BaseErrorResponse) {
	// Retrieve the work order by work_order_system_number
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, &exceptions.BaseErrorResponse{Message: "Work order not found"}
		}
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Delete the work order
	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete the work order"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) CloseOrder(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Update the work order status to 8 (Closed)
	entity.WorkOrderStatusId = 8

	// Save the updated work order
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

	// Loop through each entity and copy its data to the response
	for _, entity := range entities {
		workOrderServiceData := make(map[string]interface{})
		// Copy data from entity to response
		workOrderServiceData["work_order_service_id"] = entity.WorkOrderServiceId
		workOrderServiceData["work_order_system_number"] = entity.WorkOrderSystemNumber
		workOrderServiceData["work_order_service_remark"] = entity.WorkOrderServiceRemark
		workOrderServiceResponses = append(workOrderServiceResponses, workOrderServiceData)
	}

	// Paginate the response data
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

	// Convert entity to payload
	payload := transactionworkshoppayloads.WorkOrderServiceRequest{
		WorkOrderServiceId:     entity.WorkOrderServiceId,
		WorkOrderSystemNumber:  entity.WorkOrderSystemNumber,
		WorkOrderServiceRemark: entity.WorkOrderServiceRemark,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) UpdateRequest(tx *gorm.DB, id int, IdWorkorder int, request transactionworkshoppayloads.WorkOrderServiceRequest) *exceptions.BaseErrorResponse {
	// Retrieve the work order service request by Id
	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order service request from the database"}
	}

	// Update the work order service request
	entity.WorkOrderServiceRemark = request.WorkOrderServiceRemark

	// Save the updated work order service request
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order service request"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) AddRequest(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderServiceRequest) *exceptions.BaseErrorResponse {
	// Create a new instance of WorkOrderServiceRequest
	entities := transactionworkshopentities.WorkOrderService{
		// Assign fields from request
		WorkOrderSystemNumber:  request.WorkOrderSystemNumber,
		WorkOrderServiceRemark: request.WorkOrderServiceRemark,
	}

	// Save the work order service
	err := tx.Create(&entities).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return nil
}

func (r *WorkOrderRepositoryImpl) DeleteRequest(tx *gorm.DB, id int, IdWorkorder int) *exceptions.BaseErrorResponse {
	// Retrieve the work order service request by Id
	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		Delete(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to delete work order service request from the database"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) GetAllVehicleService(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrderServiceVehicle

	// Query to retrieve all work order service vehicle entities based on the request
	query := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}

	// Execute the query and check for errors
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order service vehicle requests from the database", Err: err}
	}

	if len(entities) == 0 {
		// Return empty data if no entities found
		return []map[string]interface{}{}, 0, 0, nil
	}

	var workOrderServiceVehicleResponses []map[string]interface{}

	// Loop through each entity and copy its data to the response
	for _, entity := range entities {
		workOrderServiceVehicleData := make(map[string]interface{})
		// Copy data from entity to response
		workOrderServiceVehicleData["work_order_service_id"] = entity.WorkOrderServiceId
		workOrderServiceVehicleData["work_order_system_number"] = entity.WorkOrderSystemNumber
		workOrderServiceVehicleData["work_order_service_date"] = entity.WorkOrderServiceDate
		workOrderServiceVehicleData["work_order_service_remark"] = entity.WorkOrderServiceRemark
		workOrderServiceVehicleResponses = append(workOrderServiceVehicleResponses, workOrderServiceVehicleData)
	}

	// Paginate the response data
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

	// Convert entity to payload
	payload := transactionworkshoppayloads.WorkOrderServiceVehicleRequest{
		WorkOrderSystemNumber:  entity.WorkOrderSystemNumber,
		WorkOrderVehicleDate:   entity.WorkOrderServiceDate,
		WorkOrderVehicleRemark: entity.WorkOrderServiceRemark,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) UpdateVehicleService(tx *gorm.DB, id int, IdWorkorder int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) *exceptions.BaseErrorResponse {
	// Retrieve the work order service request by Id
	var entity transactionworkshopentities.WorkOrderServiceVehicle
	err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order service request from the database"}
	}

	// Update the work order service request
	entity.WorkOrderServiceDate = request.WorkOrderVehicleDate
	entity.WorkOrderServiceRemark = request.WorkOrderVehicleRemark

	// Save the updated work order service request
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order service request"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) AddVehicleService(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) *exceptions.BaseErrorResponse {
	// Create a new instance of WorkOrderServiceVehicleRequest
	entities := transactionworkshopentities.WorkOrderServiceVehicle{
		// Assign fields from request
		WorkOrderSystemNumber:  request.WorkOrderSystemNumber,
		WorkOrderServiceDate:   request.WorkOrderVehicleDate,
		WorkOrderServiceRemark: request.WorkOrderVehicleRemark,
	}

	// Save the work order service
	err := tx.Create(&entities).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return nil
}

func (r *WorkOrderRepositoryImpl) DeleteVehicleService(tx *gorm.DB, id int, IdWorkorder int) *exceptions.BaseErrorResponse {
	// Retrieve the work order service request by Id
	var entity transactionworkshopentities.WorkOrderServiceVehicle
	err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		Delete(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to delete work order service request from the database"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) GenerateDocumentNumber(tx *gorm.DB, workOrderId int) (string, *exceptions.BaseErrorResponse) {
	var workOrder transactionworkshopentities.WorkOrder

	// Retrieve the work order by Id to get brand_id
	err := tx.Where("work_order_system_number = ?", workOrderId).First(&workOrder).Error
	if err != nil {
		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("failed to retrieve work order: %v", err)}
	}

	// Check if brand_id is available
	if workOrder.BrandId == 0 {
		return "", &exceptions.BaseErrorResponse{Message: "brand_id is missing in the work order"}
	}

	// Retrieve the last work order with the same brandID
	var lastWorkOrder transactionworkshopentities.WorkOrder
	err = tx.Where("brand_id = ?", workOrder.BrandId).
		Order("work_order_system_number desc").
		First(&lastWorkOrder).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("failed to retrieve last work order: %v", err)}
	}

	currentTime := time.Now()
	month := int(currentTime.Month())
	year := currentTime.Year() % 100 // Use last two digits of the year

	// Get the brand initial
	brandInitial := workOrder.BrandId

	// Check if the last work order exists
	if lastWorkOrder.WorkOrderSystemNumber == 0 {
		// If no work order exists, start the running number from 1
		return fmt.Sprintf("WSWO/%d/%02d/%02d/00001", brandInitial, month, year), nil
	}

	lastWorkOrderDate := lastWorkOrder.WorkOrderDate
	lastWorkOrderMonth := int(lastWorkOrderDate.Month())
	lastWorkOrderYear := lastWorkOrderDate.Year() % 100

	// Reset the running number if the last work order is from a different month or year
	if lastWorkOrderMonth != month || lastWorkOrderYear != year {
		return fmt.Sprintf("WSWO/%d/%02d/%02d/00001", brandInitial, month, year), nil
	}

	// Extract the running number from the last work order code and increment it by 1
	lastWorkOrderCode := lastWorkOrder.WorkOrderDocumentNumber
	codeParts := strings.Split(lastWorkOrderCode, "/")
	if len(codeParts) < 5 {
		return "", &exceptions.BaseErrorResponse{Message: "invalid last work order code format"}
	}
	lastWorkOrderNumber, err := strconv.Atoi(codeParts[4])
	if err != nil {
		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("failed to parse last work order code: %v", err)}
	}

	newWorkOrderNumber := lastWorkOrderNumber + 1

	// Format the new work order code
	newDocumentNumber := fmt.Sprintf("WSWO/%d/%02d/%02d/%05d", brandInitial, month, year, newWorkOrderNumber)
	return newDocumentNumber, nil
}

func (r *WorkOrderRepositoryImpl) Submit(tx *gorm.DB, workOrderId int) (bool, string, *exceptions.BaseErrorResponse) {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		return false, "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve work order from the database: %v", err)}
	}

	// Generate new document number
	// newDocumentNumber, genErr := r.GenerateDocumentNumber(tx, entity.WorkOrderSystemNumber)
	// if genErr != nil {
	// 	return false, "", genErr
	// }

	// // Update the work order document number
	// entity.WorkOrderDocumentNumber = newDocumentNumber

	newDocumentNumber := entity.WorkOrderDocumentNumber
	// Update the work order status to 2 (New Submitted)
	entity.WorkOrderStatusId = 2

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return false, "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to submit the work order: %v", err)}
	}

	return true, newDocumentNumber, nil
}

func (r *WorkOrderRepositoryImpl) GetAllDetailWorkOrder(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrderDetail
	// Query to retrieve all work order detail entities based on the request
	query := tx.Model(&transactionworkshopentities.WorkOrderDetail{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order detail from the database"}
	}

	var workOrderDetailResponses []map[string]interface{}

	// Loop through each entity and copy its data to the response
	for _, entity := range entities {
		workOrderDetailData := make(map[string]interface{})
		// Copy data from entity to response
		workOrderDetailData["work_order_detail_id"] = entity.WorkOrderDetailId
		workOrderDetailData["work_order_system_number"] = entity.WorkOrderSystemNumber
		workOrderDetailData["line_type_id"] = entity.LineTypeId
		workOrderDetailData["work_order_transaction_type_id"] = entity.WorkOrderTransactionTypeId
		workOrderDetailData["job_type_id"] = entity.JobTypeId
		workOrderDetailData["description"] = entity.Description
		workOrderDetailData["frt_quantity"] = entity.FrtQuantity
		workOrderDetailData["supply_quantity"] = entity.SupplyQuantity
		workOrderDetailData["price_list_id"] = entity.PriceListId

		workOrderDetailResponses = append(workOrderDetailResponses, workOrderDetailData)
	}

	// Paginate the response data
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

	// Convert entity to payload
	payload := transactionworkshoppayloads.WorkOrderDetailRequest{
		WorkOrderDetailId:          entity.WorkOrderDetailId,
		WorkOrderSystemNumber:      entity.WorkOrderSystemNumber,
		LineTypeId:                 entity.LineTypeId,
		WorkOrderTransactionTypeId: entity.WorkOrderTransactionTypeId,
		JobTypeId:                  entity.JobTypeId,
		Description:                entity.Description,
		FrtQuantity:                entity.FrtQuantity,
		SupplyQuantity:             entity.SupplyQuantity,
		PriceListId:                entity.PriceListId,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) UpdateDetailWorkOrder(tx *gorm.DB, id int, IdWorkorder int, request transactionworkshoppayloads.WorkOrderDetailRequest) *exceptions.BaseErrorResponse {
	// Retrieve the work order detail by Id
	var entity transactionworkshopentities.WorkOrderDetail
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_detail_id = ?", id, IdWorkorder).
		First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order detail from the database"}
	}

	// Update the work order detail
	entity.LineTypeId = request.LineTypeId
	entity.WorkOrderTransactionTypeId = request.WorkOrderTransactionTypeId
	entity.JobTypeId = request.JobTypeId
	entity.Description = request.Description
	entity.FrtQuantity = request.FrtQuantity
	entity.SupplyQuantity = request.SupplyQuantity
	entity.PriceListId = request.PriceListId

	// Save the updated work order detail
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order detail"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) AddDetailWorkOrder(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderDetailRequest) *exceptions.BaseErrorResponse {
	// Create a new instance of WorkOrderDetailRequest
	entities := transactionworkshopentities.WorkOrderDetail{
		// Assign fields from request
		WorkOrderSystemNumber:      request.WorkOrderSystemNumber,
		LineTypeId:                 request.LineTypeId,
		WorkOrderTransactionTypeId: request.WorkOrderTransactionTypeId,
		JobTypeId:                  request.JobTypeId,
		Description:                request.Description,
		FrtQuantity:                request.FrtQuantity,
		SupplyQuantity:             request.SupplyQuantity,
		PriceListId:                request.PriceListId,
	}

	// Save the work order detail
	err := tx.Create(&entities).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return nil
}

func (r *WorkOrderRepositoryImpl) DeleteDetailWorkOrder(tx *gorm.DB, id int, IdWorkorder int) *exceptions.BaseErrorResponse {
	// Retrieve the work order detail by Id
	var entity transactionworkshopentities.WorkOrderDetail
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_detail_id = ?", id, IdWorkorder).
		Delete(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to delete work order detail from the database"}
	}

	return nil
}
