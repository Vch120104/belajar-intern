package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	mastercampaignmasterentities "after-sales/api/entities/master/campaign_master"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"fmt"
	"net/http"
	"strconv"

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
	var entities []transactionworkshopentities.WorkOrder
	// Query to retrieve all work order entities based on the request
	query := tx.Model(&transactionworkshopentities.WorkOrder{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work orders from the database"}
	}

	var workOrderResponses []map[string]interface{}

	// Loop through each entity and copy its data to the response
	for _, entity := range entities {
		workOrderData := make(map[string]interface{})
		// Copy data from entity to response
		workOrderData["work_order_document_number"] = entity.WorkOrderDocumentNumber
		workOrderData["work_order_status_id"] = entity.WorkOrderStatusId
		workOrderData["work_order_date"] = entity.WorkOrderDate.Format("2006-01-02")
		workOrderData["work_order_type_id"] = entity.WorkOrderTypeId
		workOrderData["work_order_repeated_system_number"] = entity.WorkOrderRepeatedSystemNumber
		workOrderData["brand_id"] = entity.BrandId
		workOrderData["model_id"] = entity.ModelId
		workOrderData["vehicle_id"] = entity.VehicleId
		workOrderResponses = append(workOrderResponses, workOrderData)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(workOrderResponses, &pages)

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

func (r *WorkOrderRepositoryImpl) NewStatus(tx *gorm.DB) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse) {
	var statuses []transactionworkshopentities.WorkOrderMasterStatus
	if err := tx.Find(&statuses).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order statuses from the database"}
	}
	return statuses, nil
}

func (r *WorkOrderRepositoryImpl) NewType(tx *gorm.DB) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterType
	if err := tx.Find(&types).Error; err != nil {
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
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Convert entity to payload
	payload := transactionworkshoppayloads.WorkOrderRequest{}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptions.BaseErrorResponse) {
	var workOrderEntities = transactionworkshopentities.WorkOrder{
		// Assign fields from request
	}

	// Create a new record
	err := tx.Create(&workOrderEntities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) Submit(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to submit the work order
	// ...

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) Void(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to void the work order
	// ...

	// Save the updated work order
	err = tx.Delete(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to delete the work order"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) CloseOrder(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to close the work order
	// ...

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order"}
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
