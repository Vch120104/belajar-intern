package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	exceptions "after-sales/api/exceptions"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ServiceRequestRepositoryImpl struct {
	workorderRepo transactionworkshoprepository.WorkOrderRepository
}

func OpenServiceRequestRepositoryImpl() transactionworkshoprepository.ServiceRequestRepository {
	workorderRepo := OpenWorkOrderRepositoryImpl()
	return &ServiceRequestRepositoryImpl{
		workorderRepo: workorderRepo,
	}
}

func (s *ServiceRequestRepositoryImpl) GenerateDocumentNumberServiceRequest(tx *gorm.DB, ServiceRequestId int) (string, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequest

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
	brandResponse, brandErr := generalserviceapiutils.GetBrandGenerateDoc(entity.BrandId)
	if brandErr != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch brand data from external service",
			Err:        brandErr.Err,
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

	// Apply filters to the query
	query := utils.ApplyFilter(tx, filter)

	// Fetch records that match the filter
	if err := query.Find(&statuses).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			Message:    "Failed to retrieve service request statuses from the database",
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return statuses, nil
}

// uspg_atServiceReq0_Select
// IF @Option = 0
// --USE IN MODUL :
func (s *ServiceRequestRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var payloads []transactionworkshoppayloads.ServiceRequestNew

	joinTable := utils.CreateJoinSelectStatement(tx, transactionworkshoppayloads.ServiceRequestNew{})
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

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
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch variant data from external service",
				Err:        variantErr.Err,
			}
		}

		// Fetch color data from external service
		colourResponses, colourErr := salesserviceapiutils.GetUnitColorById(payload.ColourId)
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
		// vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(payload.VehicleId)
		// if vehicleErr != nil {
		// 	return pages, &exceptions.BaseErrorResponse{
		// 		StatusCode: vehicleErr.StatusCode,
		// 		Message:    "Failed to retrieve vehicle data from the external API",
		// 		Err:        vehicleErr.Err,
		// 	}
		// }

		// Fetch service request status from external service
		StatusResponses, statusErr := generalserviceapiutils.GetServiceRequestStatusById(payload.ServiceRequestStatusId)
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
			"service_request_date":            payload.ServiceRequestDate,
			"service_request_by":              payload.ServiceRequestBy,
			"service_company_name":            companyResponses.CompanyName,
			"brand_name":                      brandResponses.BrandName,
			"model_description":               modelResponses.ModelName,
			"variant_description":             variantResponses.VariantDescription,
			"colour_name":                     colourResponses.ColourCommercialName,
			"chassis_no":                      "vehicleResponses.Data.Master.VehicleChassisNumber",
			"no_polisi":                       "vehicleResponses.Data.STNK.VehicleRegistrationCertificateTNKB",
			"status":                          StatusResponses.ServiceRequestStatusDescription,
			"work_order_system_number":        payload.WorkOrderSystemNumber,
			"booking_system_number":           payload.BookingSystemNumber,
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
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

	// Fetch data brand from external API
	brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
	if brandErr != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch brand data from external service",
			Err:        brandErr.Err,
		}
	}

	// Fetch data model from external API
	modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
	if modelErr != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from the external API",
			Err:        modelErr.Err,
		}
	}

	// Fetch data variant from external API
	variantResponse, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
	if variantErr != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        variantErr.Err,
		}
	}

	// Fetch data colour from external API
	colourResponses, colourErr := salesserviceapiutils.GetUnitColorById(entity.ColourId)
	if colourErr != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve colour data from the external API",
			Err:        colourErr.Err,
		}
	}

	// Fetch data vehicle from external
	// vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
	// if vehicleErr != nil {
	// 	return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    "Failed to retrieve vehicle data from the external API",
	// 		Err:        vehicleErr.Err,
	// 	}
	// }

	// Fetch data Service Request Status from external API
	StatusResponses, statusErr := generalserviceapiutils.GetServiceRequestStatusById(entity.ServiceRequestStatusId)
	if statusErr != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve service request status data from the external API",
			Err:        statusErr.Err,
		}
	}

	// Fetch data company from external API
	companyResponses, companyErr := generalserviceapiutils.GetCompanyDataById(entity.CompanyId)
	if companyErr != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch company data from internal service",
			Err:        companyErr.Err,
		}
	}

	// Fetch data company from external API
	servicecompanyResponses, servicecompanyErr := generalserviceapiutils.GetCompanyDataById(entity.ServiceCompanyId)
	if servicecompanyErr != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch company data from internal service",
			Err:        servicecompanyErr.Err,
		}
	}

	// Fetch service details with pagination
	var count int64
	var serviceDetails []transactionworkshoppayloads.ServiceDetailResponse
	totalRows := tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
		Where("service_request_system_number = ?", Id).
		Count(&count).Error

	if totalRows != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count service details",
			Err:        totalRows,
		}
	}

	query := tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
		Select("service_request_detail_id, service_request_line_number, service_request_system_number, line_type_id, operation_item_id, frt_quantity").
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

	// Fetch item and UOM details for each service detail
	for i, detail := range serviceDetails {
		// Fetch data Item from external API
		var OperationItemCode string
		var Description string

		lineTypeResponse, lineErr := generalserviceapiutils.GetLineTypeById(detail.LineTypeId)
		if lineErr != nil {
			return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve line type from the external API",
				Err:        lineErr.Err,
			}
		}

		operationItemResponse, operationItemErr := s.workorderRepo.GetOperationItemById(detail.LineTypeId, detail.OperationItemId)
		if operationItemErr != nil {
			return transactionworkshoppayloads.ServiceRequestResponse{}, operationItemErr
		}

		OperationItemCode, Description, errResp := s.workorderRepo.HandleLineTypeResponse(detail.LineTypeId, operationItemResponse)
		if errResp != nil {
			return transactionworkshoppayloads.ServiceRequestResponse{}, errResp
		}

		// fetch data item
		var itemResponse masteritementities.Item
		if err := tx.Where("item_id = ?", detail.OperationItemId).First(&itemResponse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item not found",
					Err:        fmt.Errorf("item with ID %d not found", detail.OperationItemId),
				}
			}
			return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item",
				Err:        err,
			}
		}

		// Fetch data UOM from external API
		var uomItems masteritementities.Uom
		if err := tx.Where("uom_id = ?", itemResponse.UnitOfMeasurementStockId).First(&uomItems).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "UOM not found",
					Err:        fmt.Errorf("uom with ID %d not found", itemResponse.UnitOfMeasurementStockId),
				}

			}
			return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch UOM",
				Err:        err,
			}
		}

		// Update service detail with item and UOM data
		serviceDetails[i].OperationItemId = detail.OperationItemId
		serviceDetails[i].OperationItemCode = OperationItemCode
		serviceDetails[i].OperationItemName = Description
		serviceDetails[i].LineTypeCode = lineTypeResponse.LineTypeCode
		serviceDetails[i].UomName = uomItems.UomDescription
	}

	// fetch profit center from external API
	serviceprofitCenterResponses, profitCenterErr := generalserviceapiutils.GetServiceProfitCenterById(entity.ProfitCenterId)
	if profitCenterErr != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve profit center data from the external API",
			Err:        profitCenterErr.Err,
		}
	}

	// fetch dealer representative from external API
	dealerRepresentativeResponses, dealerRepresentativeErr := generalserviceapiutils.GetDealerRepresentativeById(entity.DealerRepresentativeId)
	if dealerRepresentativeErr != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve dealer representative data from the external API",
			Err:        dealerRepresentativeErr.Err,
		}
	}

	// fetch reference type from external API
	referenceTypeResponses, referenceTypeErr := generalserviceapiutils.GetReferenceTypeById(entity.ReferenceTypeId)
	if referenceTypeErr != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve reference type data from the external API",
			Err:        referenceTypeErr.Err,
		}
	}

	// fetch reference document from external API
	ReferenceDocUrl := config.EnvConfigs.GeneralServiceUrl + "service-request-reference-type/" + strconv.Itoa(entity.ReferenceSystemNumber)
	var referenceDocResponses transactionworkshoppayloads.ReferenceDoc
	errReferenceDoc := utils.Get(ReferenceDocUrl, &referenceDocResponses, nil)
	if errReferenceDoc != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve reference document data from the external API",
			Err:        errReferenceDoc,
		}
	}
	// Construct the payload with pagination information
	payload := transactionworkshoppayloads.ServiceRequestResponse{
		ServiceRequestSystemNumber:   entity.ServiceRequestSystemNumber,
		ServiceRequestStatusName:     StatusResponses.ServiceRequestStatusDescription,
		ServiceRequestDocumentNumber: entity.ServiceRequestDocumentNumber,
		ServiceRequestDate:           entity.ServiceRequestDate,
		ServiceTypeId:                entity.ServiceTypeId,
		ProfitCenterId:               entity.ProfitCenterId,
		ProfitCenterName:             serviceprofitCenterResponses.ServiceProfitCenterName,
		DealerRepresentativeName:     dealerRepresentativeResponses.DealerRepresentativeName,
		CompanyId:                    entity.CompanyId,
		CompanyName:                  companyResponses.CompanyName,
		ServiceRequestBy:             entity.ServiceRequestBy,
		ReferenceTypeId:              entity.ReferenceTypeId,
		ReferenceTypeName:            referenceTypeResponses.ReferenceTypeName,
		ReferenceDocId:               referenceDocResponses.ReferenceDocSystemNumber,
		ReferenceDocNumber:           referenceDocResponses.ReferenceDocNumber,
		ReferenceDocDate:             referenceDocResponses.ReferenceDocDate,
		BrandName:                    brandResponse.BrandName,
		ModelName:                    modelResponse.ModelName,
		VariantDescription:           variantResponse.VariantDescription,
		VariantColourName:            colourResponses.ColourCommercialName,
		VehicleId:                    entity.VehicleId,
		VehicleCode:                  "vehicleResponses.Data.Master.VehicleChassisNumber",
		VehicleTnkb:                  "vehicleResponses.Data.STNK.VehicleRegistrationCertificateTNKB",
		ServiceRemark:                entity.ServiceRemark,
		ServiceCompanyId:             entity.ServiceCompanyId,
		ServiceCompanyName:           servicecompanyResponses.CompanyName,
		ServiceDate:                  entity.ServiceDate,
		ReplyBy:                      entity.ReplyBy,
		ReplyDate:                    entity.ReplyDate,
		ReplyRemark:                  entity.ReplyRemark,
		BookingSystemNumber:          entity.BookingSystemNumber,
		EstimationSystemNumber:       entity.EstimationSystemNumber,
		ServiceDetails: transactionworkshoppayloads.ServiceRequestDetailsResponse{
			Page:       pagination.GetPage(),
			Limit:      pagination.GetLimit(),
			TotalPages: int(math.Ceil(float64(count) / float64(pagination.GetLimit()))),
			TotalRows:  int(count),
			Data:       serviceDetails,
		},
	}

	return payload, nil
}

// uspg_atServiceReq0_Insert
// IF @Option = 0
// --USE IN MODUL :
func (s *ServiceRequestRepositoryImpl) New(tx *gorm.DB, request transactionworkshoppayloads.ServiceRequestSaveRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
	defaultWorkOrderStatusId := 1 // Default status ID
	defaultReplyId := 0

	var refType string
	var ReferenceTypeId int

	loc, _ := time.LoadLocation("Asia/Jakarta") // UTC+7
	currentDate := time.Now().In(loc).Format("2006-01-02T15:04:05Z")
	parsedTime, _ := time.Parse(time.RFC3339, currentDate)

	// Check if ServiceDate is less than currentDate
	if request.ServiceDate.Before(parsedTime) {
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Service date cannot be before the current date",
			Err:        errors.New("service date cannot be before the current date"),
		}
	}

	if request.ReferenceSystemNumber == 0 && request.EstimationSystemNumber == 0 {
		refType = "SR" // New Service Request
		ReferenceTypeId = 1
	} else if request.ReferenceSystemNumber != 0 {
		refType = "WO" // Work Order
		ReferenceTypeId = 2
	} else if request.EstimationSystemNumber != 0 {
		refType = "SO" // Sales Order
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
			ServiceRequestStatusId: defaultWorkOrderStatusId,
			ServiceRequestDate:     parsedTime,
			BrandId:                request.BrandId,
			ModelId:                request.ModelId,
			VariantId:              request.VariantId,
			ColourId:               request.ColourId,
			VehicleId:              request.VehicleId,
			CompanyId:              request.CompanyId,
			DealerRepresentativeId: request.DealerRepresentativeId,
			ProfitCenterId:         request.ProfitCenterId,
			WorkOrderSystemNumber:  request.WorkOrderSystemNumber,
			BookingSystemNumber:    request.BookingSystemNumber,
			EstimationSystemNumber: request.EstimationSystemNumber,
			ReferenceSystemNumber:  request.ReferenceSystemNumber,
			ReplyId:                defaultReplyId,
			ServiceCompanyId:       request.ServiceCompanyId,
			ServiceDate:            request.ServiceDate,
			ServiceRequestBy:       request.ServiceRequestBy,
			ReferenceTypeId:        ReferenceTypeId,
			ReferenceJobType:       jobType,
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
			ServiceRequestStatusId: defaultWorkOrderStatusId,
			ServiceRequestDate:     parsedTime,
			BrandId:                request.BrandId,
			ModelId:                request.ModelId,
			VariantId:              request.VariantId,
			ColourId:               request.ColourId,
			VehicleId:              request.VehicleId,
			CompanyId:              request.CompanyId,
			DealerRepresentativeId: request.DealerRepresentativeId,
			ProfitCenterId:         request.ProfitCenterId,
			WorkOrderSystemNumber:  request.WorkOrderSystemNumber,
			BookingSystemNumber:    request.BookingSystemNumber,
			ReferenceSystemNumber:  request.ReferenceSystemNumber,
			ServiceRequestBy:       request.ServiceRequestBy,
			ServiceDate:            request.ServiceDate,
			ReferenceTypeId:        ReferenceTypeId,
			ReferenceJobType:       jobType,
		}

	case "SO":
		jobType := ""
		entities = transactionworkshopentities.ServiceRequest{
			ServiceRequestStatusId: defaultWorkOrderStatusId,
			ServiceRequestDate:     parsedTime,
			BrandId:                request.BrandId,
			ModelId:                request.ModelId,
			VariantId:              request.VariantId,
			ColourId:               request.ColourId,
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
	currentDate := time.Now().Truncate(24 * time.Hour)
	err := tx.Model(&transactionworkshopentities.ServiceRequest{}).
		Where("service_request_system_number = ?", Id).
		First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
			}
		}
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve service request",
			Err:        err,
		}
	}

	if entity.ServiceRequestStatusId != 1 {
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Service request status is not in draft",
		}
	}

	serviceDate := request.ServiceDate.Truncate(24 * time.Hour)
	if serviceDate.Before(currentDate) {
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Service date cannot be before the current date",
		}
	}

	updates := make(map[string]interface{})
	if request.ServiceTypeId != 0 {
		updates["service_type_id"] = request.ServiceTypeId
	}
	if request.ServiceCompanyId != 0 {
		updates["service_company_id"] = request.ServiceCompanyId
	}
	if !request.ServiceDate.IsZero() {
		updates["service_date"] = request.ServiceDate
	}
	if request.ServiceRemark != "" {
		updates["service_remark"] = request.ServiceRemark
	}

	if len(updates) == 0 {
		return entity, nil
	}

	err = tx.Model(&entity).Updates(updates).Error
	if err != nil {
		return transactionworkshopentities.ServiceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update service request",
			Err:        err,
		}
	}

	return entity, nil
}

// uspg_atServiceReq0_Update
// IF @Option = 2
// --USE IN MODUL :
func (s *ServiceRequestRepositoryImpl) Submit(tx *gorm.DB, Id int) (bool, string, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequest

	err := tx.Where("service_request_system_number = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", &exceptions.BaseErrorResponse{
				Message:    "Data not found",
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		return false, "", &exceptions.BaseErrorResponse{
			Message:    "Failed to retrieve service request from the database",
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if entity.ServiceRequestStatusId != 1 {
		return false, "", &exceptions.BaseErrorResponse{
			Message:    "Service Request cannot be submitted",
			StatusCode: http.StatusConflict,
			Err:        errors.New("service request cannot be submitted"),
		}
	}

	if entity.BrandId == 0 {
		return false, "", &exceptions.BaseErrorResponse{
			Message:    "Brand must be filled",
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("brand must be filled"),
		}
	}
	if entity.ModelId == 0 {
		return false, "", &exceptions.BaseErrorResponse{
			Message:    "Model must be filled",
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("model must be filled"),
		}
	}
	if entity.VariantId == 0 {
		return false, "", &exceptions.BaseErrorResponse{
			Message:    "Variant must be filled",
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("variant must be filled"),
		}
	}
	if entity.VehicleId == 0 {
		return false, "", &exceptions.BaseErrorResponse{
			Message:    "Vehicle must be filled",
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("vehicle must be filled"),
		}
	}

	if entity.ServiceRequestDocumentNumber == "" && entity.ServiceRequestStatusId == 1 {
		var serviceItemCount int64
		err = tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
			Joins("JOIN mtr_item IT ON IT.item_id = trx_service_request_detail.operation_item_id").
			Where("service_request_system_number = ? AND IT.item_type_id = ?", Id, 2).
			Where("trx_service_request_detail.line_type_id IS NULL OR trx_service_request_detail.line_type_id = ''").
			Count(&serviceItemCount).Error
		if err != nil {
			return false, "", &exceptions.BaseErrorResponse{
				Message:    "Failed to validate service items",
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		if serviceItemCount > 0 {
			return false, "", &exceptions.BaseErrorResponse{
				Message:    "Service Request has Item Type Service in details",
				StatusCode: http.StatusConflict,
				Err:        errors.New("service request has item type service in details"),
			}
		}

		var detailCount int64
		err = tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
			Where("service_request_system_number = ?", Id).
			Where("frt_quantity <= ?", 0).
			Count(&detailCount).Error
		if err != nil {
			return false, "", &exceptions.BaseErrorResponse{
				Message:    "Failed to count service request details",
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		if detailCount > 0 {
			return false, "", &exceptions.BaseErrorResponse{
				Message:    "Cannot submit service request; FRT / qty must be bigger than 0",
				StatusCode: http.StatusConflict,
				Err:        errors.New("cannot submit service request; FRT / qty must be bigger than 0"),
			}
		}

		docRequest := generalserviceapiutils.DocumentMasterRequest{
			CompanyId:       entity.CompanyId,
			TransactionDate: time.Now(),
			DocumentTypeId:  283, // document type ID
			BrandId:         entity.BrandId,
		}
		docResponse, docErr := generalserviceapiutils.GetDocumentNumber(docRequest)
		if docErr != nil {
			return false, "", docErr
		}

		entity.ServiceRequestDocumentNumber = docResponse.GeneratedDocumentNumber
		entity.ServiceRequestStatusId = 2 // Ready

		err = tx.Save(&entity).Error
		if err != nil {
			return false, "", &exceptions.BaseErrorResponse{
				Message:    "Failed to submit the service request",
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		return true, entity.ServiceRequestDocumentNumber, nil
	}

	return false, entity.ServiceRequestDocumentNumber, &exceptions.BaseErrorResponse{Message: "Service request has been submitted or the document number is already generated", StatusCode: http.StatusConflict}
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
	companyResponse, companyErr := generalserviceapiutils.GetCompanyDataById(entity.CompanyId)
	if companyErr != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve company data from the external API",
			Err:        companyErr.Err,
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

// uspg_atServiceReq1_Select
// IF @Option = 0
// --USE IN MODUL :
func (s *ServiceRequestRepositoryImpl) GetAllServiceDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.ServiceRequestDetail

	query := tx.Model(&transactionworkshopentities.ServiceRequestDetail{})
	for _, condition := range filterCondition {
		if condition.ColumnField == "service_request_system_number" {
			query = query.Where("service_request_system_number = ?", condition.ColumnValue)
		} else {
			query = query.Where(condition.ColumnField+" = ?", condition.ColumnValue)
		}
	}

	err := query.Scopes(pagination.Paginate(&pages, query)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var serviceRequestDetailResponses []map[string]interface{}

	for _, entity := range entities {
		// Fetch data Item from external API
		itemUrl := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(entity.OperationItemId)
		var itemResponse transactionworkshoppayloads.ItemServiceRequestDetail
		errItem := utils.Get(itemUrl, &itemResponse, nil)
		if errItem != nil {
			return pages, &exceptions.BaseErrorResponse{
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

		// Fetch line type from external API
		linetype, linetypeErr := generalserviceapiutils.GetLineTypeById(entity.LineTypeId)
		if linetypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: linetypeErr.StatusCode,
				Err:        linetypeErr.Err,
			}
		}

		serviceRequestDetailResponse := map[string]interface{}{
			"service_request_detail_id":     entity.ServiceRequestDetailId,
			"service_request_system_number": entity.ServiceRequestSystemNumber,
			"uom_name":                      uomItems[0].UomName,
			"item_code":                     itemResponse.ItemCode,
			"item_name":                     itemResponse.ItemName,
			"line_type_code":                linetype.LineTypeCode,
			"operation_item_id":             entity.OperationItemId,
			"reference_item_code":           itemResponse.ItemCode,
			"reference_item_name":           itemResponse.ItemName,
			"frt_quantity":                  entity.FrtQuantity,
			"reference_system_number":       entity.ReferenceSystemNumber,
			"reference_line_number":         entity.ReferenceLineNumber,
			"reference_qty":                 entity.FrtQuantity,
		}

		serviceRequestDetailResponses = append(serviceRequestDetailResponses, serviceRequestDetailResponse)
	}

	pages.Rows = serviceRequestDetailResponses
	return pages, nil
}

// uspg_atServiceReq1_Select
// IF @Option = 1
// --USE IN MODUL :
func (s *ServiceRequestRepositoryImpl) GetServiceDetailById(tx *gorm.DB, Id int) (transactionworkshoppayloads.ServiceDetailResponse, *exceptions.BaseErrorResponse) {
	var detail transactionworkshopentities.ServiceRequestDetail
	var getItemResponse transactionworkshoppayloads.ItemServiceRequestDetail
	var getUomItems []transactionworkshoppayloads.UomItemServiceRequestDetail

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

	// Fetch data Item from external API
	itemUrl := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(detail.OperationItemId)
	errItem := utils.Get(itemUrl, &getItemResponse, nil)
	if errItem != nil {
		return transactionworkshoppayloads.ServiceDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errItem,
		}
	}

	// Fetch data Uom from external API
	uomUrl := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/?page=0&limit=10&uom_id=" + strconv.Itoa(getItemResponse.UomId)
	errUom := utils.Get(uomUrl, &getUomItems, nil)
	if errUom != nil {
		return transactionworkshoppayloads.ServiceDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUom,
		}
	}

	getLineTypeResponse, getLineTypeErr := generalserviceapiutils.GetLineTypeById(detail.LineTypeId)
	if getLineTypeErr != nil {
		return transactionworkshoppayloads.ServiceDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve line type data from the external API",
			Err:        getLineTypeErr.Err,
		}
	}

	serviceDetail := transactionworkshoppayloads.ServiceDetailResponse{
		ServiceRequestDetailId:     detail.ServiceRequestDetailId,
		ServiceRequestSystemNumber: detail.ServiceRequestSystemNumber,
		LineTypeId:                 detail.LineTypeId,
		LineTypeCode:               getLineTypeResponse.LineTypeCode,
		OperationItemId:            detail.OperationItemId,
		OperationItemCode:          getItemResponse.ItemCode,
		OperationItemName:          getItemResponse.ItemName,
		UomName:                    getUomItems[0].UomName,
		FrtQuantity:                detail.FrtQuantity,
		ReferenceSystemNumber:      detail.ReferenceSystemNumber,
	}

	return serviceDetail, nil

}

// uspg_atServiceReq1_Insert
// IF @Option = 50
// --USE IN MODUL :
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

// uspg_atServiceReq1_Update
// IF @Option = 1
// --USE IN MODUL :
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

// uspg_atServiceReq1_Delete
// IF @Option = 0
// --USE IN MODUL :
func (s *ServiceRequestRepositoryImpl) DeleteServiceDetailMultiId(tx *gorm.DB, Id int, DetailIds []int) (bool, *exceptions.BaseErrorResponse) {
	if err := tx.Where("service_request_system_number = ? AND service_request_detail_id IN (?)", Id, DetailIds).
		Delete(&transactionworkshopentities.ServiceRequestDetail{}).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete the service request detail",
			Err:        err,
		}
	}

	if tx.RowsAffected == 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No records found to delete",
		}
	}

	return true, nil
}

func (s *ServiceRequestRepositoryImpl) NewServiceType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.ServiceRequestMasterServiceType, *exceptions.BaseErrorResponse) {
	var statuses []transactionworkshopentities.ServiceRequestMasterServiceType

	// Apply filters to the query
	query := utils.ApplyFilter(tx, filter)

	// Fetch records that match the filter
	if err := query.Find(&statuses).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			Message:    "Failed to retrieve service request statuses from the database",
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return statuses, nil
}
