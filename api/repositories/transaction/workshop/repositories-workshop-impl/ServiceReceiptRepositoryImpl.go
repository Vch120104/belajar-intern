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

func (s *ServiceReceiptRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var payloads []transactionworkshoppayloads.ServiceReceiptNew

	joinTable := utils.CreateJoinSelectStatement(tx, transactionworkshoppayloads.ServiceReceiptNew{})
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	whereQuery = whereQuery.Where("service_request_system_number != 0 AND service_request_status_id = 2")

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&payloads).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(payloads) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, payload := range payloads {
		// Fetch brand data from external service
		brandResponses, brandErr := salesserviceapiutils.GetUnitBrandById(payload.BrandId)
		if brandErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: brandErr.StatusCode,
				Message:    "Failed to fetch brand data from external service",
				Err:        brandErr.Err,
			}
		}

		// Fetch model data from external service
		modelResponses, modelErr := salesserviceapiutils.GetUnitModelById(payload.ModelId)
		if modelErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: modelErr.StatusCode,
				Message:    "Failed to fetch model data from external service",
				Err:        modelErr.Err,
			}
		}

		// Fetch variant data from external service
		variantResponses, variantErr := salesserviceapiutils.GetUnitVariantById(payload.VariantId)
		if variantErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: variantErr.StatusCode,
				Message:    "Failed to fetch variant data from external service",
				Err:        variantErr.Err,
			}
		}

		// Fetch color data from external service
		colourResponses, colourErr := salesserviceapiutils.GetUnitColourByBrandId(payload.BrandId)
		if colourErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: colourErr.StatusCode,
				Message:    "Failed to fetch color data from external service",
				Err:        colourErr.Err,
			}
		}

		// Fetch company data from external service
		companyResponses, companyErr := generalserviceapiutils.GetCompanyDataById(payload.CompanyId)
		if companyErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: companyErr.StatusCode,
				Message:    "Failed to fetch company data from internal service",
				Err:        companyErr.Err,
			}
		}

		// Fetch vehicle data from external service
		vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(payload.VehicleId)
		if vehicleErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: vehicleErr.StatusCode,
				Message:    "Failed to retrieve vehicle data from the external API",
				Err:        vehicleErr.Err,
			}
		}

		// Fetch service request status from external service
		statusResponses, statusErr := generalserviceapiutils.GetServiceRequestStatusById(payload.ServiceRequestStatusId)
		if statusErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: statusErr.StatusCode,
				Message:    "Failed to retrieve service request status data from the external API",
				Err:        statusErr.Err,
			}
		}

		result := map[string]interface{}{
			"service_request_system_number":   payload.ServiceRequestSystemNumber,
			"service_request_document_number": payload.ServiceRequestDocumentNumber,
			"service_request_date":            payload.ServiceRequestDate.Format("2006-01-02 15:04:05"),
			"service_request_by":              payload.ServiceRequestBy,
			"service_company_name":            companyResponses.CompanyName,
			"brand_name":                      brandResponses.BrandName,
			"model_description":               modelResponses.ModelName,
			"variant_description":             variantResponses.VariantDescription,
			"colour_name":                     colourResponses.Data[0].ColourCommercialName,
			"chassis_no":                      vehicleResponses.Data.Master.VehicleChassisNumber,
			"no_polisi":                       vehicleResponses.Data.STNK.VehicleRegistrationCertificateTNKB,
			"status":                          statusResponses.ServiceRequestStatusDescription,
			"work_order_system_number":        payload.WorkOrderSystemNumber,
			"booking_system_number":           payload.BookingSystemNumber,
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
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
		Select("service_request_detail_id, service_request_line_number, service_request_system_number, line_type_id, operation_item_id, frt_quantity").
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
		VariantDescription:           variantResponse.VariantDescription,
		VariantColourName:            colourResponses.Data[0].ColourCommercialName,
		VehicleId:                    entity.VehicleId,
		VehicleCode:                  vehicleResponses.Data.Master.VehicleChassisNumber,
		VehicleTnkb:                  vehicleResponses.Data.STNK.VehicleRegistrationCertificateTNKB,
		CompanyId:                    entity.CompanyId,
		CompanyName:                  companyResponses.CompanyName,
		DealerRepresentativeId:       entity.DealerRepresentativeId,
		DealerRepresentativeName:     dealerRepresentativeResponses.DealerRepresentativeName,
		ProfitCenterId:               entity.ProfitCenterId,
		ProfitCenterName:             profitCenterResponses.ProfitCenterName,
		WorkOrderSystemNumber:        entity.WorkOrderSystemNumber,
		//WorkOrderDocumentNumber:      workOrderDocumentNumber,
		BookingSystemNumber:    entity.BookingSystemNumber,
		EstimationSystemNumber: entity.EstimationSystemNumber,
		ReferenceSystemNumber:  entity.ReferenceSystemNumber,
		ReferenceDocNumber:     "", //entity.ReferenceDocNumber,
		ReferenceDocDate:       "", //entity.ReferenceDocDate,
		ReplyId:                entity.ReplyId,
		ReplyBy:                entity.ReplyBy,
		ReplyDate:              ReplyDate,
		ReplyRemark:            entity.ReplyRemark,
		ServiceCompanyId:       entity.ServiceCompanyId,
		ServiceCompanyName:     companyResponses.CompanyName,
		ServiceDate:            serviceDate,
		ServiceRequestBy:       entity.ServiceRequestBy,
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
