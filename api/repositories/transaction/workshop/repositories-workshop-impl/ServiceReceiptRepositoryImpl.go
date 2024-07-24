package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ServiceReceiptRepositoryImpl struct {
}

func OpenServiceReceiptRepositoryImpl() transactionworkshoprepository.ServiceReceiptRepository {
	return &ServiceReceiptRepositoryImpl{}
}

func (s *ServiceReceiptRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshoppayloads.ServiceReceiptNew

	joinTable := utils.CreateJoinSelectStatement(tx, transactionworkshoppayloads.ServiceReceiptNew{})

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)
	whereQuery = whereQuery.Where("service_request_system_number != 0 AND service_request_status_id = 2")

	if err := whereQuery.Find(&entities).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Fetch data brand from external API
	BrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entities[0].BrandId)
	var brandResponses transactionworkshoppayloads.WorkOrderVehicleBrand
	errBrand := utils.Get(BrandUrl, &brandResponses, nil)
	if errBrand != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve brand data from the external API",
			Err:        errBrand,
		}
	}

	// Fetch data model from external API
	ModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(entities[0].ModelId)
	var modelResponses transactionworkshoppayloads.WorkOrderVehicleModel
	errModel := utils.Get(ModelUrl, &modelResponses, nil)
	if errModel != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from the external API",
			Err:        errModel,
		}
	}

	// fetch data variant from external API
	VariantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(entities[0].VariantId)
	var variantResponses transactionworkshoppayloads.WorkOrderVehicleVariant
	errVariant := utils.Get(VariantUrl, &variantResponses, nil)
	if errVariant != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        errVariant,
		}
	}

	// fetch data colour from external API
	ColourUrl := config.EnvConfigs.SalesServiceUrl + "unit-color-dropdown/" + strconv.Itoa(entities[0].BrandId)
	var colourResponses []transactionworkshoppayloads.WorkOrderVehicleColour
	errColour := utils.GetArray(ColourUrl, &colourResponses, nil)
	if errColour != nil || len(colourResponses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve colour data from the external API",
			Err:        errColour,
		}
	}

	// Fetch data company from external API
	// CompanyUrl := config.EnvConfigs.GeneralServiceUrl + "company-id/" + strconv.Itoa(ServiceRequestReq.CompanyId)
	// var companyResponses []transactionworkshoppayloads.CompanyResponse
	// errCompany := utils.GetArray(CompanyUrl, &companyResponses, nil)
	// if errCompany != nil || len(companyResponses) == 0 {
	// 	return nil, 0, 0, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    "Failed to retrieve company data from the external API",
	// 		Err:        errCompany,
	// 	}
	// }

	// Fetch data vehicle from external API
	VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(entities[0].VehicleId)
	fmt.Println("Fetching URL: ", VehicleUrl) // Debug: Print URL

	var vehicleResponses transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.Get(VehicleUrl, &vehicleResponses, nil)
	// Debug: Print the fetched vehicle response
	vehicleResponseJSON, _ := json.Marshal(vehicleResponses)
	fmt.Println("Vehicle Response: ", string(vehicleResponseJSON))
	if errVehicle != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	// Fetch data Service Request Status from external API
	ServiceRequestStatusURL := config.EnvConfigs.AfterSalesServiceUrl + "service-request/dropdown-status?service_request_status_id=" + strconv.Itoa(entities[0].ServiceRequestStatusId)
	//fmt.Println("Fetching Work Order Status data from:", ServiceRequestStatusURL)
	var StatusResponses []transactionworkshoppayloads.ServiceRequestStatusResponse // Use slice of ServiceRequestStatusResponse
	errStatus := utils.GetArray(ServiceRequestStatusURL, &StatusResponses, nil)
	if errStatus != nil || len(StatusResponses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve status service request data from the external API",
			Err:        errStatus,
		}
	}

	var convertedResponses []transactionworkshoppayloads.ServiceReceiptResponse
	for _, entity := range entities {
		convertedResponses = append(convertedResponses, transactionworkshoppayloads.ServiceReceiptResponse{
			ServiceRequestSystemNumber:   entity.ServiceRequestSystemNumber,
			ServiceRequestDocumentNumber: entity.ServiceRequestDocumentNumber,
			ServiceRequestDate:           entity.ServiceRequestDate.Format("2006-01-02 15:04:05"),
			ServiceRequestBy:             entity.ServiceRequestBy,
			CompanyId:                    entity.CompanyId,
			CompanyName:                  "", //entity.CompanyName,
			ServiceCompanyId:             entity.ServiceCompanyId,
			BrandId:                      entity.BrandId,
			BrandName:                    brandResponses.BrandName,
			ModelId:                      entity.ModelId,
			ModelName:                    modelResponses.ModelName,
			VariantId:                    entity.VariantId,
			VariantName:                  variantResponses.VariantName,
			VariantColourName:            colourResponses[0].VariantColourName,
			VehicleId:                    entity.VehicleId,
			VehicleCode:                  vehicleResponses.Master.VehicleCode,
			VehicleTnkb:                  vehicleResponses.Stnk.VehicleTnkb,
			ServiceRequestStatusId:       entity.ServiceRequestStatusId,
			ServiceRequestStatusName:     StatusResponses[0].ServiceRequestStatusName,
			WorkOrderSystemNumber:        entity.WorkOrderSystemNumber,
			BookingSystemNumber:          entity.BookingSystemNumber,
			EstimationSystemNumber:       entity.EstimationSystemNumber,
			ReferenceDocSystemNumber:     entity.ReferenceDocSystemNumber,
			ReplyId:                      entity.ReplyId,
			ServiceDate:                  entity.ServiceDate.Format("2006-01-02 15:04:05"),
		})
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"service_request_system_number":         response.ServiceRequestSystemNumber,
			"service_request_document_number":       response.ServiceRequestDocumentNumber,
			"service_request_date":                  response.ServiceRequestDate,
			"service_request_by":                    response.ServiceRequestBy,
			"company_name":                          response.CompanyName,
			"brand_name":                            response.BrandName,
			"model_description":                     response.ModelName,
			"variant_description":                   response.VariantName,
			"colour_name":                           response.VariantColourName,
			"vehicle_chassis_number":                response.VehicleCode,
			"vehicle_registration_certificate_tnkb": response.VehicleTnkb,
			"service_request_status_name":           response.ServiceRequestStatusName,
			"work_order_system_number":              response.WorkOrderSystemNumber,
			"booking_system_number":                 response.BookingSystemNumber,
			"reference_doc_system_number":           response.ReferenceDocSystemNumber,
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
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entity.BrandId)
	var brandResponse transactionworkshoppayloads.WorkOrderVehicleBrand
	errBrand := utils.Get(brandUrl, &brandResponse, nil)
	if errBrand != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve brand data from the external API",
			Err:        errBrand,
		}
	}

	// fetch data model from external api
	modelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(entity.ModelId)
	var modelResponse transactionworkshoppayloads.WorkOrderVehicleModel
	errModel := utils.Get(modelUrl, &modelResponse, nil)
	if errModel != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from the external API",
			Err:        errModel,
		}
	}

	// fetch data variant from external api
	variantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(entity.VariantId)
	var variantResponse transactionworkshoppayloads.WorkOrderVehicleVariant
	errVariant := utils.Get(variantUrl, &variantResponse, nil)
	if errVariant != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        errVariant,
		}
	}

	// Fetch data colour from external API
	ColourUrl := config.EnvConfigs.SalesServiceUrl + "unit-color-dropdown/" + strconv.Itoa(entity.BrandId)
	var colourResponses []transactionworkshoppayloads.WorkOrderVehicleColour
	errColour := utils.GetArray(ColourUrl, &colourResponses, nil)
	if errColour != nil || len(colourResponses) == 0 {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve colour data from the external API",
			Err:        errColour,
		}
	}

	//fetch data company from external api
	CompanyUrl := config.EnvConfigs.GeneralServiceUrl + "company/" + strconv.Itoa(entity.CompanyId)
	var companyResponse transactionworkshoppayloads.CompanyResponse
	errCompany := utils.Get(CompanyUrl, &companyResponse, nil)
	if errCompany != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve company data from the external API",
			Err:        errCompany,
		}
	}

	// fetch data vehicle from external api
	VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(entity.VehicleId)
	var vehicleResponses transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.Get(VehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	// Fetch data Service Request Status from external API
	ServiceRequestStatusURL := config.EnvConfigs.AfterSalesServiceUrl + "service-request/dropdown-status?service_request_status_id=" + strconv.Itoa(entity.ServiceRequestStatusId)
	var StatusResponses []transactionworkshoppayloads.ServiceRequestStatusResponse
	errStatus := utils.GetArray(ServiceRequestStatusURL, &StatusResponses, nil)
	if errStatus != nil || len(StatusResponses) == 0 {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve status service request data from the external API",
			Err:        errStatus,
		}
	}

	// Fetch work order from external API
	WorkOrderUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order?work_order_system_number=" + strconv.Itoa(entity.WorkOrderSystemNumber)
	var WorkOrderResponses []transactionworkshoppayloads.WorkOrderRequestResponse
	errWorkOrder := utils.GetArray(WorkOrderUrl, &WorkOrderResponses, nil)

	// Check for error and assign blank value if a 404 error
	workOrderDocumentNumber := ""
	if errWorkOrder != nil {
		if strings.Contains(errWorkOrder.Error(), "404") {
			workOrderDocumentNumber = ""
		} else {
			return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve work order data from the external API",
				Err:        errWorkOrder,
			}
		}
	} else if len(WorkOrderResponses) > 0 {
		workOrderDocumentNumber = WorkOrderResponses[0].WorkOrderDocumentNumber
	}

	// Fetch service details with pagination
	var serviceDetails []transactionworkshoppayloads.ServiceReceiptDetailResponse
	query := tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
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

	payload := transactionworkshoppayloads.ServiceReceiptResponse{
		ServiceRequestSystemNumber:   entity.ServiceRequestSystemNumber,
		ServiceRequestStatusId:       entity.ServiceRequestStatusId,
		ServiceRequestStatusName:     StatusResponses[0].ServiceRequestStatusName,
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
		VehicleCode:                  vehicleResponses.Master.VehicleCode,
		VehicleTnkb:                  vehicleResponses.Stnk.VehicleTnkb,
		CompanyId:                    entity.CompanyId,
		CompanyName:                  "", //companyResponses[0].CompanyName,
		DealerRepresentativeId:       entity.DealerRepresentativeId,
		ProfitCenterId:               entity.ProfitCenterId,
		WorkOrderSystemNumber:        entity.WorkOrderSystemNumber,
		WorkOrderDocumentNumber:      workOrderDocumentNumber,
		BookingSystemNumber:          entity.BookingSystemNumber,
		EstimationSystemNumber:       entity.EstimationSystemNumber,
		ReferenceDocSystemNumber:     entity.ReferenceDocSystemNumber,
		ReferenceDocNumber:           0,  //entity.ReferenceDocNumber,
		ReferenceDocDate:             "", //entity.ReferenceDocDate,
		ReplyId:                      entity.ReplyId,
		ReplyBy:                      entity.ReplyBy,
		ReplyDate:                    ReplyDate,
		ReplyRemark:                  entity.ReplyRemark,
		ServiceCompanyId:             entity.ServiceCompanyId,
		ServiceCompanyName:           "", //servicecompanyResponses[0].CompanyName,
		ServiceDate:                  serviceDate,
		ServiceRequestBy:             entity.ServiceRequestBy,
		ServiceDetails: transactionworkshoppayloads.ServiceReceiptDetailsResponse{
			Page:       pagination.GetPage(),
			Limit:      pagination.GetLimit(),
			TotalPages: pagination.TotalPages,
			TotalRows:  int(pagination.TotalRows), // Convert int64 to int
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

	// Check if ServiceRequestStatusId is 2 (can save data)
	if entity.ServiceRequestStatusId == 2 {
		// Check if ServiceDate is before the current date
		if entity.ServiceDate.Before(currentDate) {
			return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Service date cannot be before the current date",
				Err:        errors.New("service date cannot be before the current date"),
			}
		}

		// Update entity fields
		entity.ReplyRemark = request.ReplyRemark
		entity.ReplyBy = "Admin" // Hardcoded value; this should be passed from the session user
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

		// Save the updated entity
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
