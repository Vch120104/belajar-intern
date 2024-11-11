package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ServiceReceiptRepositoryImpl struct {
}

func OpenServiceReceiptRepositoryImpl() transactionworkshoprepository.ServiceReceiptRepository {
	return &ServiceReceiptRepositoryImpl{}
}

func (s *ServiceReceiptRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tableStruct := transactionworkshoppayloads.ServiceReceiptNew{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	whereQuery = whereQuery.Where("service_request_system_number != 0 AND service_request_status_id = 2")

	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	var convertedResponses []transactionworkshoppayloads.ServiceReceiptGetAllResponse
	for rows.Next() {
		var (
			ServiceReceiptReq transactionworkshoppayloads.ServiceReceiptNew
			ServiceReceiptRes transactionworkshoppayloads.ServiceReceiptGetAllResponse
		)

		if err := rows.Scan(
			&ServiceReceiptReq.ServiceRequestSystemNumber,
			&ServiceReceiptReq.ServiceRequestDocumentNumber,
			&ServiceReceiptReq.ServiceRequestDate,
			&ServiceReceiptReq.ServiceRequestBy,
			&ServiceReceiptReq.ServiceRequestStatusId,
			&ServiceReceiptReq.BrandId,
			&ServiceReceiptReq.ModelId,
			&ServiceReceiptReq.VariantId,
			&ServiceReceiptReq.VehicleId,
			&ServiceReceiptReq.BookingSystemNumber,
			&ServiceReceiptReq.EstimationSystemNumber,
			&ServiceReceiptReq.WorkOrderSystemNumber,
			&ServiceReceiptReq.ReferenceDocSystemNumber,
			&ServiceReceiptReq.ProfitCenterId,
			&ServiceReceiptReq.CompanyId,
			&ServiceReceiptReq.DealerRepresentativeId,
			&ServiceReceiptReq.ServiceTypeId,
			&ServiceReceiptReq.ReferenceTypeId,
			&ServiceReceiptReq.ServiceRemark,
			&ServiceReceiptReq.ServiceCompanyId,
			&ServiceReceiptReq.ServiceDate,
			&ServiceReceiptReq.ReplyId,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch data brand from external services
		brandResponses, brandErr := salesserviceapiutils.GetUnitBrandById(ServiceReceiptReq.BrandId)
		if brandErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch brand data from external service",
				Err:        brandErr.Err,
			}
		}

		// Fetch data model from external services
		modelResponses, modelErr := salesserviceapiutils.GetUnitModelById(ServiceReceiptReq.ModelId)
		if modelErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch model data from external service",
				Err:        modelErr.Err,
			}
		}

		// Fetch data variant from external services
		variantResponses, variantErr := salesserviceapiutils.GetUnitVariantById(ServiceReceiptReq.VariantId)
		if variantErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch variant data from external service",
				Err:        variantErr.Err,
			}
		}

		colourResponses, colourErr := salesserviceapiutils.GetUnitColourByBrandId(ServiceReceiptReq.BrandId)
		if colourErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve colour data from the external API",
				Err:        colourErr.Err,
			}
		}

		companyResponses, companyErr := generalserviceapiutils.GetCompanyDataById(ServiceReceiptReq.CompanyId)
		if companyErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve company data from the external API",
				Err:        companyErr.Err,
			}
		}

		vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(ServiceReceiptReq.VehicleId)
		if vehicleErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve vehicle data from the external API",
				Err:        vehicleErr.Err,
			}
		}

		statusResponses, statusErr := generalserviceapiutils.GetServiceRequestStatusById(ServiceReceiptReq.ServiceRequestStatusId)
		if statusErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve status service request data from the external API",
				Err:        statusErr.Err,
			}
		}

		ServiceReceiptRes = transactionworkshoppayloads.ServiceReceiptGetAllResponse{
			ServiceRequestSystemNumber:   ServiceReceiptReq.ServiceRequestSystemNumber,
			ServiceRequestDocumentNumber: ServiceReceiptReq.ServiceRequestDocumentNumber,
			ServiceRequestDate:           ServiceReceiptReq.ServiceRequestDate.Format("2006-01-02 15:04:05"),
			ServiceRequestBy:             ServiceReceiptReq.ServiceRequestBy,
			ServiceRequestStatusName:     statusResponses.ServiceRequestStatusDescription,
			BrandName:                    brandResponses.BrandName,
			ModelName:                    modelResponses.ModelName,
			VariantName:                  variantResponses.VariantName,
			VariantColourName:            colourResponses[0].VariantColourName,
			VehicleCode:                  vehicleResponses.VehicleChassisNumber,
			VehicleTnkb:                  vehicleResponses.VehicleRegistrationCertificateTNKB,
			CompanyName:                  companyResponses.CompanyName,
			WorkOrderSystemNumber:        ServiceReceiptReq.WorkOrderSystemNumber,
			BookingSystemNumber:          ServiceReceiptReq.BookingSystemNumber,
			EstimationSystemNumber:       ServiceReceiptReq.EstimationSystemNumber,
			ReferenceDocSystemNumber:     ServiceReceiptReq.ReferenceDocSystemNumber,
			ServiceDate:                  ServiceReceiptReq.ServiceDate.Format("2006-01-02 15:04:05"),
		}

		convertedResponses = append(convertedResponses, ServiceReceiptRes)
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"service_request_system_number":   response.ServiceRequestSystemNumber,
			"service_request_document_number": response.ServiceRequestDocumentNumber,
			"service_request_date":            response.ServiceRequestDate,
			"service_request_by":              response.ServiceRequestBy,
			"service_company_name":            response.CompanyName,
			"brand_name":                      response.BrandName,
			"model_code_description":          response.ModelName,
			"variant_code_description":        response.VariantName,
			"colour_name":                     response.VariantColourName,
			"chassis_no":                      response.VehicleCode,
			"no_polisi":                       response.VehicleTnkb,
			"status":                          response.ServiceRequestStatusName,
			"work_order_system_number":        response.WorkOrderSystemNumber,
			"work_order_no":                   response.WorkOrderDocumentNumber,
			"booking_system_number":           response.BookingSystemNumber,
			"booking_no":                      response.BookingDocumentNumber,
			"reference_doc_system_number":     response.ReferenceDocSystemNumber,
			"ref_doc_no":                      response.ReferenceDocDocumentNumber,
		}

		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)
	return paginatedData, totalPages, totalRows, nil
}

func (s *ServiceReceiptRepositoryImpl) GetById(tx *gorm.DB, Id int, pagination pagination.Pagination) (transactionworkshoppayloads.ServiceReceiptResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequest
	err := tx.Model(&transactionworkshopentities.ServiceRequest{}).Where("service_request_system_number = ? AND service_request_status_id = 2", Id).First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
			}
		}

		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Convert service date format
	serviceDate := utils.SafeConvertDateFormat(entity.ServiceDate)
	ServiceRequestDate := utils.SafeConvertDateFormat(entity.ServiceRequestDate)
	ReplyDate := utils.SafeConvertDateFormat(entity.ReplyDate)

	// fetch data brand from external api
	brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
	if brandErr != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch brand data from external service",
			Err:        brandErr.Err,
		}
	}

	// fetch data model from external api
	modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
	if modelErr != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from the external API",
			Err:        modelErr.Err,
		}
	}

	// fetch data variant from external api
	variantResponse, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
	if variantErr != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        variantErr.Err,
		}
	}

	// Fetch data colour from external API
	colourResponses, colourErr := salesserviceapiutils.GetUnitColourByBrandId(entity.BrandId)
	if colourErr != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve colour data from the external API",
			Err:        colourErr.Err,
		}
	}

	//fetch data company from external api
	companyResponses, companyErr := generalserviceapiutils.GetCompanyDataById(entity.CompanyId)
	if companyErr != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve company data from the external API",
			Err:        companyErr.Err,
		}
	}

	// fetch data vehicle from external api
	vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
	if vehicleErr != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        vehicleErr.Err,
		}
	}

	// Fetch data Service Request Status from external API
	StatusResponses, statusErr := generalserviceapiutils.GetServiceRequestStatusById(entity.ServiceRequestStatusId)
	if statusErr != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve status service request data from the external API",
			Err:        statusErr.Err,
		}
	}

	// Fetch work order from external API
	// WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order?work_order_system_number=" + strconv.Itoa(entity.WorkOrderSystemNumber)
	// var WorkOrderResponses []transactionworkshoppayloads.WorkOrderRequestResponse
	// errWorkOrder := utils.GetArray(WorkOrderUrl, &WorkOrderResponses, nil)

	// // Check for error and assign blank value if a 404 error
	// workOrderDocumentNumber := ""
	// if errWorkOrder != nil {
	// 	if strings.Contains(errWorkOrder.Error(), "404") {
	// 		workOrderDocumentNumber = ""
	// 	} else {
	// 		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
	// 			StatusCode: http.StatusInternalServerError,
	// 			Message:    "Failed to retrieve work order data from the external API",
	// 			Err:        errWorkOrder,
	// 		}
	// 	}
	// } else if len(WorkOrderResponses) > 0 {
	// 	workOrderDocumentNumber = WorkOrderResponses[0].WorkOrderDocumentNumber
	// }

	// Fetch service details with pagination
	var serviceDetails []transactionworkshoppayloads.ServiceReceiptDetailResponse
	totalRowsQuery := tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
		Where("service_request_system_number = ?", Id).
		Count(new(int64)).Error

	if totalRowsQuery != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count service details",
			Err:        totalRowsQuery,
		}
	}

	query := tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
		Select("service_request_detail_id, service_request_line_number, service_request_system_number, line_type_id, operation_item_id, frt_quantity, reference_doc_system_number, reference_doc_id").
		Where("service_request_system_number = ?", Id).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit())
	errServiceDetails := query.Find(&serviceDetails).Error
	if errServiceDetails != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve service details from the database",
			Err:        errServiceDetails,
		}
	}

	// Fetch item and UOM details for each service detail
	for i, detail := range serviceDetails {
		// Fetch data Item from external API
		itemUrl := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(detail.OperationItemId)
		var itemResponse transactionworkshoppayloads.ItemServiceRequestDetail
		errItem := utils.Get(itemUrl, &itemResponse, nil)
		if errItem != nil {
			return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errItem,
			}
		}

		// Fetch data UOM from external API
		uomUrl := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/?page=0&limit=10&uom_id=" + strconv.Itoa(itemResponse.UomId)
		var uomItems []transactionworkshoppayloads.UomItemServiceRequestDetail
		errUom := utils.Get(uomUrl, &uomItems, nil)
		if errUom != nil || len(uomItems) == 0 {
			uomItems = []transactionworkshoppayloads.UomItemServiceRequestDetail{
				{UomName: "N/A"},
			}
		}

		// Update service detail with item and UOM data
		serviceDetails[i].OperationItemCode = itemResponse.ItemCode
		serviceDetails[i].OperationItemName = itemResponse.ItemName
		serviceDetails[i].UomName = uomItems[0].UomName
	}

	// fetch profit center from external API
	profitCenterResponses, profitCenterErr := generalserviceapiutils.GetProfitCenterById(entity.ProfitCenterId)
	if profitCenterErr != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve profit center data from the external API",
			Err:        profitCenterErr.Err,
		}
	}

	// fetch dealer representative from external API
	dealerRepresentativeResponses, dealerRepresentativeErr := generalserviceapiutils.GetDealerRepresentativeById(entity.DealerRepresentativeId)
	if dealerRepresentativeErr != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve dealer representative data from the external API",
			Err:        dealerRepresentativeErr.Err,
		}
	}

	totalRows := 0
	payload := transactionworkshoppayloads.ServiceReceiptResponse{
		ServiceRequestSystemNumber:   entity.ServiceRequestSystemNumber,
		ServiceRequestStatusId:       entity.ServiceRequestStatusId,
		ServiceRequestStatusName:     StatusResponses.ServiceRequestStatusDescription,
		ServiceRequestDocumentNumber: entity.ServiceRequestDocumentNumber,
		ServiceRequestDate:           ServiceRequestDate,
		BrandId:                      entity.BrandId,
		BrandName:                    brandResponse.BrandName,
		ModelId:                      entity.ModelId,
		ModelName:                    modelResponse.ModelName,
		VariantId:                    entity.VariantId,
		VariantName:                  variantResponse.VariantName,
		VariantColourName:            colourResponses[0].VariantColourName,
		VehicleId:                    entity.VehicleId,
		VehicleCode:                  vehicleResponses.VehicleChassisNumber,
		VehicleTnkb:                  vehicleResponses.VehicleRegistrationCertificateTNKB,
		CompanyId:                    entity.CompanyId,
		CompanyName:                  companyResponses.CompanyName,
		DealerRepresentativeId:       entity.DealerRepresentativeId,
		DealerRepresentativeName:     dealerRepresentativeResponses.DealerRepresentativeName,
		ProfitCenterId:               entity.ProfitCenterId,
		ProfitCenterName:             profitCenterResponses.ProfitCenterName,
		WorkOrderSystemNumber:        entity.WorkOrderSystemNumber,
		//WorkOrderDocumentNumber:      workOrderDocumentNumber,
		BookingSystemNumber:      entity.BookingSystemNumber,
		EstimationSystemNumber:   entity.EstimationSystemNumber,
		ReferenceDocSystemNumber: entity.ReferenceDocSystemNumber,
		ReferenceDocNumber:       "", //entity.ReferenceDocNumber,
		ReferenceDocDate:         "", //entity.ReferenceDocDate,
		ReplyId:                  entity.ReplyId,
		ReplyBy:                  entity.ReplyBy,
		ReplyDate:                ReplyDate,
		ReplyRemark:              entity.ReplyRemark,
		ServiceCompanyId:         entity.ServiceCompanyId,
		ServiceCompanyName:       companyResponses.CompanyName,
		ServiceDate:              serviceDate,
		ServiceRequestBy:         entity.ServiceRequestBy,
		ServiceDetails: transactionworkshoppayloads.ServiceReceiptDetailsResponse{
			Page:       pagination.GetPage(),
			Limit:      pagination.GetLimit(),
			TotalPages: int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit()))),
			TotalRows:  totalRows,
			Data:       serviceDetails,
		},
	}

	return payload, nil
}

func (s *ServiceReceiptRepositoryImpl) Save(tx *gorm.DB, Id int, request transactionworkshoppayloads.ServiceReceiptSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequest
	currentDate := time.Now()

	// Check if the service request exists
	err := tx.Model(&transactionworkshopentities.ServiceRequest{}).Where("service_request_system_number = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Data not found", Err: err}
		}
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Message: "Query error", Err: err}
	}

	// Check if ServiceRequestStatusId is in draft (1) status
	if entity.ServiceRequestStatusId == 1 {
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Service request is in draft status",
			Err:        errors.New("service request is in draft status"),
		}
	}

	// Check if ServiceRequestStatusId is in ready (2) status
	if entity.ServiceRequestStatusId == 2 {
		// Check if ServiceDate is before the current date
		if entity.ServiceDate.Before(currentDate) {
			return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Service date cannot be before the current date",
				Err:        errors.New("service date cannot be before the current date"),
			}
		}

		// Check if status is transitioning to 3 (Accept) or 6 (Reject) from a state other than 2 (Ready)
		if (request.ServiceRequestStatusId == 3 || request.ServiceRequestStatusId == 6) && entity.ServiceRequestStatusId != 2 {
			return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Message:    "Cannot change status to Accept or Reject; request is no longer in Ready state",
				Err:        errors.New("status transition conflict"),
			}
		}

		// Update entity fields
		entity.ReplyRemark = request.ReplyRemark
		entity.ReplyBy = "Admin"
		entity.ReplyDate = currentDate

		// Set ServiceRequestStatusId based on request.ServiceRequestStatusId
		switch request.ServiceRequestStatusId {
		case 3:
			entity.ServiceRequestStatusId = 3 // Accept
		case 6:
			entity.ServiceRequestStatusId = 6 // Reject
		default:
			return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid ServiceRequestStatusId provided",
				Err:        errors.New("invalid ServiceRequestStatusId"),
			}
		}

		err = tx.Save(&entity).Error
		if err != nil {
			return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to save the service request",
				Err:        err,
			}
		}

		return entity, nil
	}

	return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
		StatusCode: http.StatusBadRequest,
		Message:    "Service request status is not in a valid state for saving data",
		Err:        errors.New("service request status is not valid for saving data"),
	}
}
