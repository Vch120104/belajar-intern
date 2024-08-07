package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"encoding/json"
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

type ServiceRequestRepositoryImpl struct {
}

func OpenServiceRequestRepositoryImpl() transactionworkshoprepository.ServiceRequestRepository {
	return &ServiceRequestRepositoryImpl{}
}

func (s *ServiceRequestRepositoryImpl) GenerateDocumentNumberServiceRequest(tx *gorm.DB, ServiceRequestId int) (string, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequest
	var brandResponse transactionworkshoppayloads.WorkOrderVehicleBrand

	// Get the service request based on the service request system number
	err := tx.Model(&transactionworkshopentities.ServiceRequest{}).Where("service_request_system_number = ?", ServiceRequestId).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Service request not found"}
		}
		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve service request from the database: %v", err)}
	}

	if entity.BrandId == 0 {
		return "", &exceptions.BaseErrorResponse{Message: "brand_id is missing in the service request. Please ensure the service request has a valid brand_id before generating document number."}
	}

	// Get the last service request based on the service request system number
	var lastServiceRequest transactionworkshopentities.ServiceRequest
	err = tx.Model(&transactionworkshopentities.ServiceRequest{}).
		Where("brand_id = ?", entity.BrandId).
		Order("service_request_document_number desc").
		First(&lastServiceRequest).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve last service request: %v", err)}
	}

	currentTime := time.Now()
	month := int(currentTime.Month())
	year := currentTime.Year() % 100 // Use last two digits of the year

	// fetch data brand from external api
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entity.BrandId)
	errUrl := utils.Get(brandUrl, &brandResponse, nil)
	if errUrl != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrl,
		}
	}

	// Check if the brand code is empty
	if brandResponse.BrandCode == "" {
		return "", &exceptions.BaseErrorResponse{Message: "Brand code is empty"}
	}

	// Get the initial of the brand abbreviation
	brandInitial := brandResponse.BrandCode[0]

	newDocumentNumber := fmt.Sprintf("WSSQ/%c/%02d/%02d/00001", brandInitial, month, year)
	if lastServiceRequest.ServiceRequestSystemNumber != 0 {
		lastServiceRequestDate := lastServiceRequest.ServiceRequestDate
		lastServiceRequestYear := lastServiceRequestDate.Year() % 100

		// Check if the last service request is in the same month and year
		if lastServiceRequestYear == year {
			lastServiceRequestCode := lastServiceRequest.ServiceRequestDocumentNumber
			codeParts := strings.Split(lastServiceRequestCode, "/")
			if len(codeParts) == 5 {
				lastDocumentNumber, err := strconv.Atoi(codeParts[4])
				if err == nil {
					newServiceRequestNumber := lastDocumentNumber + 1
					newDocumentNumber = fmt.Sprintf("WSSQ/%c/%02d/%02d/%05d", brandInitial, month, year, newServiceRequestNumber)
				} else {
					log.Printf("Failed to parse last service request number: %v", err)
				}
			} else {
				log.Println("Invalid service request number format")
			}
		}
	}

	log.Printf("New document number: %s", newDocumentNumber)
	return newDocumentNumber, nil
}

func (r *ServiceRequestRepositoryImpl) NewStatus(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.ServiceRequestMasterStatus, *exceptions.BaseErrorResponse) {
	var statuses []transactionworkshopentities.ServiceRequestMasterStatus

	query := utils.ApplyFilter(tx, filter)

	if err := query.Find(&statuses).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			Message:    "Failed to retrieve service request statuses from the database",
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return statuses, nil
}

func (s *ServiceRequestRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tableStruct := transactionworkshoppayloads.ServiceRequestNew{}

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

	var convertedResponses []transactionworkshoppayloads.ServiceRequestResponse
	for rows.Next() {
		var (
			ServiceRequestReq transactionworkshoppayloads.ServiceRequestNew
			ServiceRequestRes transactionworkshoppayloads.ServiceRequestResponse
		)

		if err := rows.Scan(
			&ServiceRequestReq.ServiceRequestSystemNumber,
			&ServiceRequestReq.ServiceRequestDocumentNumber,
			&ServiceRequestReq.ServiceRequestDate,
			&ServiceRequestReq.ServiceRequestBy,
			&ServiceRequestReq.ServiceRequestStatusId,
			&ServiceRequestReq.BrandId,
			&ServiceRequestReq.ModelId,
			&ServiceRequestReq.VariantId,
			&ServiceRequestReq.VehicleId,
			&ServiceRequestReq.BookingSystemNumber,
			&ServiceRequestReq.EstimationSystemNumber,
			&ServiceRequestReq.WorkOrderSystemNumber,
			&ServiceRequestReq.ReferenceDocSystemNumber,
			&ServiceRequestReq.ProfitCenterId,
			&ServiceRequestReq.CompanyId,
			&ServiceRequestReq.DealerRepresentativeId,
			&ServiceRequestReq.ServiceTypeId,
			&ServiceRequestReq.ReferenceTypeId,
			&ServiceRequestReq.ServiceRemark,
			&ServiceRequestReq.ServiceCompanyId,
			&ServiceRequestReq.ServiceDate,
			&ServiceRequestReq.ReplyId,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch data brand from external API
		BrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(ServiceRequestReq.BrandId)
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
		ModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(ServiceRequestReq.ModelId)
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
		VariantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(ServiceRequestReq.VariantId)
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
		ColourUrl := config.EnvConfigs.SalesServiceUrl + "unit-color-dropdown/" + strconv.Itoa(ServiceRequestReq.BrandId)
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
		CompanyUrl := config.EnvConfigs.GeneralServiceUrl + "companies-redis?company_id=" + strconv.Itoa(ServiceRequestReq.CompanyId)
		var companyResponses []transactionworkshoppayloads.CompanyResponse
		errCompany := utils.GetArray(CompanyUrl, &companyResponses, nil)
		if errCompany != nil || len(companyResponses) == 0 {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve company data from the external API",
				Err:        errCompany,
			}
		}

		// Fetch data vehicle from external API
		VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(ServiceRequestReq.VehicleId)
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
		ServiceRequestStatusURL := config.EnvConfigs.AfterSalesServiceUrl + "service-request/dropdown-status?service_request_status_id=" + strconv.Itoa(ServiceRequestReq.ServiceRequestStatusId)
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

		ServiceRequestRes = transactionworkshoppayloads.ServiceRequestResponse{
			ServiceRequestSystemNumber:   ServiceRequestReq.ServiceRequestSystemNumber,
			ServiceRequestDocumentNumber: ServiceRequestReq.ServiceRequestDocumentNumber,
			ServiceRequestDate:           ServiceRequestReq.ServiceRequestDate.Format("2006-01-02 15:04:05"),
			ServiceRequestBy:             ServiceRequestReq.ServiceRequestBy,
			ServiceRequestStatusId:       ServiceRequestReq.ServiceRequestStatusId,
			ServiceRequestStatusName:     StatusResponses[0].ServiceRequestStatusName,
			BrandName:                    brandResponses.BrandName,
			ModelName:                    modelResponses.ModelName,
			VariantName:                  variantResponses.VariantName,
			VariantColourName:            colourResponses[0].VariantColourName,
			VehicleCode:                  vehicleResponses.Master.VehicleCode,
			VehicleTnkb:                  vehicleResponses.Stnk.VehicleTnkb,
			CompanyId:                    ServiceRequestReq.CompanyId,
			CompanyName:                  companyResponses[0].CompanyName,
			DealerRepresentativeId:       ServiceRequestReq.DealerRepresentativeId,
			ProfitCenterId:               ServiceRequestReq.ProfitCenterId,
			WorkOrderSystemNumber:        ServiceRequestReq.WorkOrderSystemNumber,
			BookingSystemNumber:          ServiceRequestReq.BookingSystemNumber,
			EstimationSystemNumber:       ServiceRequestReq.EstimationSystemNumber,
			ReferenceDocSystemNumber:     ServiceRequestReq.ReferenceDocSystemNumber,
			ReplyId:                      ServiceRequestReq.ReplyId,
			ServiceCompanyId:             ServiceRequestReq.ServiceCompanyId,
			ServiceDate:                  ServiceRequestReq.ServiceDate.Format("2006-01-02 15:04:05"),
		}

		convertedResponses = append(convertedResponses, ServiceRequestRes)
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

func (s *ServiceRequestRepositoryImpl) GetById(tx *gorm.DB, Id int, pagination pagination.Pagination) (transactionworkshoppayloads.ServiceRequestResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequest
	err := tx.Model(&transactionworkshopentities.ServiceRequest{}).
		Where("service_request_system_number = ?", Id).
		First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
			}
		}
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Convert service date format
	serviceDate := utils.SafeConvertDateFormat(entity.ServiceDate)
	ServiceRequestDate := utils.SafeConvertDateFormat(entity.ServiceRequestDate)
	ReplyDate := utils.SafeConvertDateFormat(entity.ReplyDate)

	// Fetch data brand from external API
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entity.BrandId)
	var brandResponse transactionworkshoppayloads.WorkOrderVehicleBrand
	errBrand := utils.Get(brandUrl, &brandResponse, nil)
	if errBrand != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve brand data from the external API",
			Err:        errBrand,
		}
	}

	// Fetch data model from external API
	modelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(entity.ModelId)
	var modelResponse transactionworkshoppayloads.WorkOrderVehicleModel
	errModel := utils.Get(modelUrl, &modelResponse, nil)
	if errModel != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from the external API",
			Err:        errModel,
		}
	}

	// Fetch data variant from external API
	variantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(entity.VariantId)
	var variantResponse transactionworkshoppayloads.WorkOrderVehicleVariant
	errVariant := utils.Get(variantUrl, &variantResponse, nil)
	if errVariant != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve colour data from the external API",
			Err:        errColour,
		}
	}

	// Fetch data vehicle from external API
	VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(entity.VehicleId)
	var vehicleResponses transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.Get(VehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve status service request data from the external API",
			Err:        errStatus,
		}
	}

	// Fetch data company from external API
	CompanyUrl := config.EnvConfigs.GeneralServiceUrl + "companies-redis?company_id=" + strconv.Itoa(entity.CompanyId)
	var companyResponses []transactionworkshoppayloads.CompanyResponse
	errCompany := utils.GetArray(CompanyUrl, &companyResponses, nil)
	if errCompany != nil || len(companyResponses) == 0 {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve company data from the external API",
			Err:        errCompany,
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
			return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve work order data from the external API",
				Err:        errWorkOrder,
			}
		}
	} else if len(WorkOrderResponses) > 0 {
		workOrderDocumentNumber = WorkOrderResponses[0].WorkOrderDocumentNumber
	}

	// Fetch service details with pagination
	var serviceDetails []transactionworkshoppayloads.ServiceRequestDetailResponse
	query := tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
		Where("service_request_system_number = ?", Id).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit())
	errServiceDetails := query.Find(&serviceDetails).Error
	if errServiceDetails != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve service details from the database",
			Err:        errServiceDetails,
		}
	}

	// fetch profit center from external API
	ProfitCenterUrl := config.EnvConfigs.GeneralServiceUrl + "profit-center/" + strconv.Itoa(entity.ProfitCenterId)
	var profitCenterResponses transactionworkshoppayloads.ProfitCenter
	errProfitCenter := utils.Get(ProfitCenterUrl, &profitCenterResponses, nil)
	if errProfitCenter != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve profit center data from the external API",
			Err:        errProfitCenter,
		}
	}

	// fetch dealer representative from external API
	DealerRepresentativeUrl := config.EnvConfigs.GeneralServiceUrl + "dealer-representative/" + strconv.Itoa(entity.DealerRepresentativeId)
	var dealerRepresentativeResponses transactionworkshoppayloads.DealerRepresentative
	errDealerRepresentative := utils.Get(DealerRepresentativeUrl, &dealerRepresentativeResponses, nil)
	if errDealerRepresentative != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve dealer representative data from the external API",
			Err:        errDealerRepresentative,
		}
	}

	// Construct the payload with pagination information
	payload := transactionworkshoppayloads.ServiceRequestResponse{
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
		CompanyName:                  companyResponses[0].CompanyName,
		DealerRepresentativeId:       entity.DealerRepresentativeId,
		DealerRepresentativeName:     dealerRepresentativeResponses.DealerRepresentativeName,
		ProfitCenterId:               entity.ProfitCenterId,
		ProfitCenterName:             profitCenterResponses.ProfitCenterName,
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
		ServiceCompanyName:           companyResponses[0].CompanyName,
		ServiceDate:                  serviceDate,
		ServiceRequestBy:             entity.ServiceRequestBy,
		ServiceDetails: transactionworkshoppayloads.ServiceRequestDetailsResponse{
			Page:       pagination.GetPage(),
			Limit:      pagination.GetLimit(),
			TotalPages: pagination.TotalPages,
			TotalRows:  int(pagination.TotalRows),
			Data:       serviceDetails,
		},
	}

	return payload, nil
}

func (s *ServiceRequestRepositoryImpl) New(tx *gorm.DB, request transactionworkshoppayloads.ServiceRequestSaveRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
	defaultWorkOrderStatusId := 1 // 1:Draft, 2:Ready, 3:Accept, 4:Work Order, 5:Booking, 6:Reject, 7:Cancel, 8:Closed
	currentDate := time.Now()
	defaultReplyId := 0

	var refType string
	var ReferenceTypeId int

	// Check if ServiceDate is less than currentDate
	if request.ServiceDate.Before(currentDate) {
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Service date cannot be before the current date",
			Err:        errors.New("service date cannot be before the current date"),
		}
	}

	if request.ReferenceDocSystemNumber == 0 && request.EstimationSystemNumber == 0 {
		refType = "SR" // Use "SR" for New Service Request
		ReferenceTypeId = 1
	} else if request.ReferenceDocSystemNumber != 0 {
		refType = "WO" // Use "WO" for Work Order reference type
		ReferenceTypeId = 2
	} else if request.EstimationSystemNumber != 0 {
		refType = "SO" // Use "SO" for Sales Order reference type
		ReferenceTypeId = 3
	} else {
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid reference document type",
			Err:        errors.New("invalid reference document type"),
		}
	}

	var entities transactionworkshopentities.ServiceRequest

	switch refType {
	case "SR":
		jobType := ""
		entities = transactionworkshopentities.ServiceRequest{
			ServiceRequestStatusId:   defaultWorkOrderStatusId,
			ServiceRequestDate:       currentDate,
			BrandId:                  request.BrandId,
			ModelId:                  request.ModelId,
			VariantId:                request.VariantId,
			VehicleId:                request.VehicleId,
			CompanyId:                request.CompanyId,
			DealerRepresentativeId:   request.DealerRepresentativeId,
			ProfitCenterId:           request.ProfitCenterId,
			WorkOrderSystemNumber:    request.WorkOrderSystemNumber,
			BookingSystemNumber:      request.BookingSystemNumber,
			EstimationSystemNumber:   request.EstimationSystemNumber,
			ReferenceDocSystemNumber: request.ReferenceDocSystemNumber,
			ReplyId:                  defaultReplyId,
			ServiceCompanyId:         request.ServiceCompanyId,
			ServiceDate:              request.ServiceDate,
			ServiceRequestBy:         request.ServiceRequestBy,
			ReferenceTypeId:          ReferenceTypeId,
			ReferenceJobType:         jobType,
		}

	case "WO":
		jobType := getJobType(request.ProfitCenterId, request.ServiceProfitCenterId)
		if jobType == "" {
			return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid profit center combination",
				Err:        errors.New("invalid profit center combination"),
			}
		}

		entities = transactionworkshopentities.ServiceRequest{

			ServiceRequestStatusId:   defaultWorkOrderStatusId,
			ServiceRequestDate:       currentDate,
			BrandId:                  request.BrandId,
			ModelId:                  request.ModelId,
			VariantId:                request.VariantId,
			VehicleId:                request.VehicleId,
			CompanyId:                request.CompanyId,
			DealerRepresentativeId:   request.DealerRepresentativeId,
			ProfitCenterId:           request.ProfitCenterId,
			WorkOrderSystemNumber:    request.WorkOrderSystemNumber,
			BookingSystemNumber:      request.BookingSystemNumber,
			ReferenceDocSystemNumber: request.ReferenceDocSystemNumber,
			ServiceRequestBy:         request.ServiceRequestBy,
			ServiceDate:              request.ServiceDate,
			ReferenceTypeId:          ReferenceTypeId,
			ReferenceJobType:         jobType,
		}

	case "SO":
		jobType := ""
		entities = transactionworkshopentities.ServiceRequest{
			ServiceRequestStatusId: defaultWorkOrderStatusId,
			ServiceRequestDate:     currentDate,
			BrandId:                request.BrandId,
			ModelId:                request.ModelId,
			VariantId:              request.VariantId,
			VehicleId:              request.VehicleId,
			CompanyId:              request.CompanyId,
			DealerRepresentativeId: request.DealerRepresentativeId,
			ProfitCenterId:         request.ProfitCenterId,
			EstimationSystemNumber: request.EstimationSystemNumber,
			ReplyId:                request.ReplyId,
			ServiceCompanyId:       request.ServiceCompanyId,
			ServiceDate:            request.ServiceDate,
			ServiceRequestBy:       request.ServiceRequestBy,
			ReferenceTypeId:        ReferenceTypeId,
			ReferenceJobType:       jobType,
		}

	default:
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid reference document type",
			Err:        errors.New("invalid reference document type"),
		}
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func getJobType(profitCenterId, serviceProfitCenterId int) string {
	switch {
	case profitCenterId == 1 && serviceProfitCenterId == 2:
		return "JOBTYPE_TB"
	case profitCenterId == 1 && serviceProfitCenterId == 1:
		return "JOBTYPE_TG"
	case profitCenterId == 2 && serviceProfitCenterId == 2:
		return "JOBTYPE_TB"
	case profitCenterId == 2 && serviceProfitCenterId == 1:
		return "JOBTYPE_TG"
	default:
		return ""
	}
}

func (s *ServiceRequestRepositoryImpl) Save(tx *gorm.DB, Id int, request transactionworkshoppayloads.ServiceRequestSaveDataRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequest
	currentDate := time.Now()

	// cek service request system number
	err := tx.Model(&transactionworkshopentities.ServiceRequest{}).Where("service_request_system_number = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Data not found"}
		}
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Err: err}
	}

	// Check current service request status
	if entity.ServiceRequestStatusId != 1 {
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: "Service request status is not in draft"}
	}

	// Check if ServiceDate is less than currentDate
	if request.ServiceDate.Before(currentDate) {
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Service date cannot be before the current date",
			Err:        errors.New("service date cannot be before the current date"),
		}
	}

	entity.ServiceTypeId = request.ServiceTypeId
	entity.ServiceCompanyId = request.ServiceCompanyId
	entity.ServiceDate = request.ServiceDate
	entity.ServiceRemark = request.ServiceRemark

	err = tx.Save(&entity).Error
	if err != nil {
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save the service request",
			Err:        err}
	}

	return entity, nil
}

func (s *ServiceRequestRepositoryImpl) Submit(tx *gorm.DB, Id int) (bool, string, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequest

	err := tx.Where("service_request_system_number = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", &exceptions.BaseErrorResponse{Message: "Data not found"}
		}
		return false, "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve service request from the database: %v", err)}
	}

	if entity.BrandId == 0 {
		return false, "", &exceptions.BaseErrorResponse{Message: "Brand must be filled"}
	}
	if entity.ModelId == 0 {
		return false, "", &exceptions.BaseErrorResponse{Message: "Model must be filled"}
	}

	if entity.ServiceRequestDocumentNumber == "" && entity.ServiceRequestStatusId == 1 {

		// Check if there are service request details with non-zero FrtQuantity
		var detailCount int64
		tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
			Where("service_request_system_number = ? AND frt_quantity > 0", Id).
			Count(&detailCount)

		if detailCount == 0 {
			return false, "", &exceptions.BaseErrorResponse{Message: "Cannot submit service request detail ftr / qty must be > 0"}
		}

		newDocumentNumber, genErr := s.GenerateDocumentNumberServiceRequest(tx, entity.ServiceRequestSystemNumber)
		if genErr != nil {
			return false, "", genErr
		}

		entity.ServiceRequestDocumentNumber = newDocumentNumber
		entity.ServiceRequestStatusId = 2 // Ready

		err = tx.Save(&entity).Error
		if err != nil {
			return false, "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to submit the service request: %v", err)}
		}

		return true, newDocumentNumber, nil
	} else {
		return false, entity.ServiceRequestDocumentNumber, &exceptions.BaseErrorResponse{Message: "Service request has been submitted or the document number is already generated"}
	}
}

func (s *ServiceRequestRepositoryImpl) Void(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequest

	err := tx.Model(&transactionworkshopentities.ServiceRequest{}).
		Where("service_request_system_number = ?", Id).
		First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if entity.ServiceRequestStatusId != 1 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Service request status is not in draft",
			Err:        errors.New("service request status is not in draft"),
		}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete the service request",
			Err:        err,
		}
	}

	return true, nil
}

func (s *ServiceRequestRepositoryImpl) CloseOrder(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.ServiceRequest
	err := tx.Model(&transactionworkshopentities.ServiceRequest{}).
		Select("work_order_system_number, booking_system_number, service_request_status_id").
		Where("service_request_system_number = ?", Id).
		First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Service request not found",
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve service request from the database",
			Err:        err}
	}

	if entity.ServiceRequestStatusId == 1 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Service request status cannot be closed because status is draft"}
	}

	var woSysNo, woStatClosed int
	var servReqStatusId, woStatusId, servReqStatWorkOrderId, refType, jobTypeId int

	refType = 2                // 1:Service Request 2:Work Order 3:Sales Order
	woStatClosed = 8           // 8:Closed
	servReqStatWorkOrderId = 4 // 4:Work Order
	servReqStatAccept := 3     // 3:Accept

	woSysNo = entity.WorkOrderSystemNumber
	servReqStatusId = entity.ServiceRequestStatusId

	var workOrder transactionworkshopentities.WorkOrder
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Select("work_order_status_id").
		Where("work_order_system_number = ?", woSysNo).
		First(&workOrder).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order, please check the work order is exist.",
			Err:        err,
		}
	}

	woStatusId = workOrder.WorkOrderStatusId

	// Check business category from gmComp0 table
	CompanyUrl := config.EnvConfigs.GeneralServiceUrl + "company/" + strconv.Itoa(entity.CompanyId)
	var companyResponse transactionworkshoppayloads.CompanyResponse
	errCompany := utils.Get(CompanyUrl, &companyResponse, nil)
	if errCompany != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errCompany,
		}
	}

	// check company category is not "001" uspg_atServiceReq0_Update / IMG Bina Trada - Pusat / company_id 130
	if companyResponse.CompanyId != 130 {

		if entity.WorkOrderSystemNumber == 0 && entity.BookingSystemNumber == 0 {

			if woStatusId != woStatClosed {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "Service Request cannot be closed. because Work Order is not closed. "}
			}

			if servReqStatusId != servReqStatWorkOrderId {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "Service Request cannot be closed. because Service Request status is in Work Order. "}
			} else {

				err = tx.Model(&transactionworkshopentities.ServiceRequest{}).
					Select("reference_doc_system_number, service_profit_center_id").
					Where("service_request_system_number = ?", Id).
					First(&entity).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to retrieve service reference document number, please check the service reference document number is exist.",
						Err:        err}
				}

				// cek reference type WO
				if entity.ReferenceTypeId == refType {

					jobType := getJobType(entity.ProfitCenterId, entity.ServiceProfitCenterId)
					if jobType == "JOBTYPE_TB" {
						jobTypeId = 1
					} else {
						jobTypeId = 2
					}

					srvStatQcPass := 6
					lineTypeOpr := 1
					lineTypePack := 2
					woStatQcPass := 6 // 1:Draft, 2:New, 3:Ready, 4:On Going, 5:Stop, 6:QC Pass, 7:Cancel, 8:Closed

					// update detail work order
					err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
						Where("work_order_system_number = ? AND job_type_id = ?", entity.WorkOrderSystemNumber, jobTypeId).
						Updates(map[string]interface{}{
							"work_order_status_id": gorm.Expr("CASE WHEN line_type_id = ? OR line_type_id = ? THEN ? ELSE work_order_status_id END", lineTypeOpr, lineTypePack, srvStatQcPass),
							"supply_quantity":      gorm.Expr("FRT_QTY"),
						}).Error
					if err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to update work order detail, please check the work order detail is exist.",
							Err:        err}
					}

					var count int64
					tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
						Where("work_order_system_number = ? AND work_order_status_id <> ? AND line_type_id IN (?, ?)", entity.WorkOrderSystemNumber, srvStatQcPass, lineTypeOpr, lineTypePack).
						Count(&count)
					if count == 0 {
						err := tx.Model(&transactionworkshopentities.WorkOrder{}).
							Where("work_order_system_number = ? ", entity.WorkOrderSystemNumber).
							Updates(map[string]interface{}{
								"work_order_status_id": woStatQcPass,
							}).Error
						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to update work order, please check the work order status in qc pass.",
								Err:        err}
						}
					}

					err = tx.Model(&transactionworkshopentities.ServiceRequest{}).
						Where("service_request_system_number = ?", Id).
						Update("service_request_status_id", 8).Error
					if err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to close service request status, please check the service request status is exist.",
							Err:        err}
					}
				}
			}

		} else {
			if servReqStatusId != servReqStatAccept {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "Service Request cannot be closed. because Service Request status is not in Accept."}
			} else {
				err = tx.Model(&transactionworkshopentities.ServiceRequest{}).
					Where("service_request_system_number = ?", Id).
					Update("service_request_status_id", 8).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to close service request status, please check the service request status is Accept.",
						Err:        err}
				}
			}
		}
	} else {

		err = tx.Model(&transactionworkshopentities.ServiceRequest{}).
			Where("service_request_system_number = ?", Id).
			Update("service_request_status_id", 8).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to close service request status, please check the company category .",
				Err:        err}
		}
	}

	return true, nil
}

func (s *ServiceRequestRepositoryImpl) GetAllServiceDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.ServiceRequestDetail
	var getItemResponse transactionworkshoppayloads.ItemServiceRequestDetail
	var getUomItems []transactionworkshoppayloads.UomItemServiceRequestDetail

	// Build the query with filters
	query := tx.Model(&transactionworkshopentities.ServiceRequestDetail{})
	if len(filterCondition) > 0 {
		for _, condition := range filterCondition {
			if condition.ColumnField == "service_request_system_number" {
				query = query.Where("service_request_system_number = ?", condition.ColumnValue)
			} else {
				query = query.Where(condition.ColumnField+" = ?", condition.ColumnValue)
			}
		}
	}

	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	// Fetch data Item from external API
	if len(entities) > 0 {
		itemUrl := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(entities[0].OperationItemId)
		errItem := utils.Get(itemUrl, &getItemResponse, nil)
		if errItem != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errItem,
			}
		}

		// Fetch data Uom from external API
		uomUrl := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/?page=0&limit=10&uom_id=" + strconv.Itoa(getItemResponse.UomId)
		errUom := utils.Get(uomUrl, &getUomItems, nil)
		if errUom != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        errUom,
			}
		}
	}

	var ServiceRequestDetailResponses []map[string]interface{}

	for _, entity := range entities {
		ServiceRequestDetailResponse := map[string]interface{}{
			"service_request_system_number": entity.ServiceRequestSystemNumber,
			"uom_name":                      getUomItems[0].UomName,
			"item_code":                     getItemResponse.ItemCode,
			"item_name":                     getItemResponse.ItemName,
			"line_type_id":                  entity.LineTypeId,
			"operation_item_id":             entity.OperationItemId,
			"frt_quantity":                  entity.FrtQuantity,
		}

		ServiceRequestDetailResponses = append(ServiceRequestDetailResponses, ServiceRequestDetailResponse)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(ServiceRequestDetailResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (s *ServiceRequestRepositoryImpl) GetServiceDetailById(tx *gorm.DB, Id int) (transactionworkshoppayloads.ServiceDetailResponse, *exceptions.BaseErrorResponse) {
	var detail transactionworkshopentities.ServiceRequestDetail
	err := tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
		Where("service_request_detail_id = ?", Id).
		First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.ServiceDetailResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Service detail not found",
			}
		}
		return transactionworkshoppayloads.ServiceDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	serviceDetail := transactionworkshoppayloads.ServiceDetailResponse{
		ServiceRequestDetailId:     detail.ServiceRequestDetailId,
		ServiceRequestId:           detail.ServiceRequestId,
		ServiceRequestSystemNumber: detail.ServiceRequestSystemNumber,
		LineTypeId:                 detail.LineTypeId,
		OperationItemId:            detail.OperationItemId,
		FrtQuantity:                detail.FrtQuantity,
		ReferenceDocSystemNumber:   detail.ReferenceDocSystemNumber,
		ReferenceDocId:             detail.ReferenceDocId,
	}

	return serviceDetail, nil

}

func (s *ServiceRequestRepositoryImpl) AddServiceDetail(tx *gorm.DB, id int, request transactionworkshoppayloads.ServiceDetailSaveRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse) {

	entity := transactionworkshopentities.ServiceRequestDetail{
		ServiceRequestSystemNumber: request.ServiceRequestSystemNumber,
		LineTypeId:                 request.LineTypeId,
		OperationItemId:            request.OperationItemId,
		FrtQuantity:                request.FrtQuantity,
	}

	err := tx.Create(&entity).Error
	if err != nil {
		return transactionworkshopentities.ServiceRequestDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entity, nil
}

func (s *ServiceRequestRepositoryImpl) UpdateServiceDetail(tx *gorm.DB, Id int, DetailId int, request transactionworkshoppayloads.ServiceDetailUpdateRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse) {

	var serviceRequest transactionworkshopentities.ServiceRequest
	err := tx.Where("service_request_system_number = ?", Id).First(&serviceRequest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshopentities.ServiceRequestDetail{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Data not found"}
		}
		return transactionworkshopentities.ServiceRequestDetail{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Err: err}
	}

	if serviceRequest.ServiceRequestStatusId != 1 {
		return transactionworkshopentities.ServiceRequestDetail{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: "Service request status is not in draft"}
	}

	var entity transactionworkshopentities.ServiceRequestDetail
	err = tx.Where("service_request_detail_id = ?", DetailId).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshopentities.ServiceRequestDetail{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Data not found"}
		}
		return transactionworkshopentities.ServiceRequestDetail{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Err: err}
	}

	// Validate the service request detail
	var count int64
	err = tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
		Where("service_request_system_number = ? AND service_request_detail_id = ?", Id, DetailId).
		Count(&count).Error
	if err != nil {
		return transactionworkshopentities.ServiceRequestDetail{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Message: fmt.Sprintf("Failed to validate service request detail: %v", err)}
	}
	if count == 0 {
		return transactionworkshopentities.ServiceRequestDetail{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Service request detail not found for the given service request"}
	}

	// Validate the FRT / Qty value
	if request.FrtQuantity <= 0 {
		return transactionworkshopentities.ServiceRequestDetail{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: "FRT / Qty must be bigger than 0"}
	}

	entity.FrtQuantity = request.FrtQuantity

	err = tx.Save(&entity).Error
	if err != nil {
		return transactionworkshopentities.ServiceRequestDetail{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return entity, nil
}

func (s *ServiceRequestRepositoryImpl) DeleteServiceDetail(tx *gorm.DB, Id int, DetailId int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequestDetail
	err := tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).Where("service_request_detail_id = ?", DetailId).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Data not found"}
		}
		return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Err: err}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return true, nil
}

func (s *ServiceRequestRepositoryImpl) DeleteServiceDetailMultiId(tx *gorm.DB, Id int, DetailIds []int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequestDetail
	err := tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).Where("service_request_system_number = ? AND service_request_detail_id IN (?)", Id, DetailIds).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Data not found"}
		}
		return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Err: err}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return true, nil
}
