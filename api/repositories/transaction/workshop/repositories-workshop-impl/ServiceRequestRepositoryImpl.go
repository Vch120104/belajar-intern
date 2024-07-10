package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
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

	// Check if BrandAbbreviation is not empty before using it
	if brandResponse.BrandCode == "" {
		return "", &exceptions.BaseErrorResponse{Message: "Brand code is empty"}
	}

	// Get the initial of the brand abbreviation
	brandInitial := brandResponse.BrandCode[0]

	// Handle the case when there is no last service request or the format is invalid
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

		ServiceRequestRes = transactionworkshoppayloads.ServiceRequestResponse{
			ServiceRequestSystemNumber:   ServiceRequestReq.ServiceRequestSystemNumber,
			ServiceRequestDocumentNumber: ServiceRequestReq.ServiceRequestDocumentNumber,
			ServiceRequestDate:           ServiceRequestReq.ServiceRequestDate.Format("2006-01-02 15:04:05"),
			ServiceRequestBy:             ServiceRequestReq.ServiceRequestBy,
			ServiceRequestStatusId:       ServiceRequestReq.ServiceRequestStatusId,
			BrandId:                      ServiceRequestReq.BrandId,
			ModelId:                      ServiceRequestReq.ModelId,
			VehicleId:                    ServiceRequestReq.VehicleId,
			CompanyId:                    ServiceRequestReq.CompanyId,
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
			"service_request_system_number":   response.ServiceRequestSystemNumber,
			"service_request_document_number": response.ServiceRequestDocumentNumber,
			"service_request_date":            response.ServiceRequestDate,
			"service_request_by":              response.ServiceRequestBy,
			"company_id":                      response.CompanyId,
			"service_company_id":              response.ServiceCompanyId,
			"brand_id":                        response.BrandId,
			"model_id":                        response.ModelId,
			"vehicle_id":                      response.VehicleId,
			"service_request_status_id":       response.ServiceRequestStatusId,
			"work_order_system_number":        response.WorkOrderSystemNumber,
			"booking_system_number":           response.BookingSystemNumber,
			"reference_doc_system_number":     response.ReferenceDocSystemNumber,
		}

		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (s *ServiceRequestRepositoryImpl) GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.ServiceRequestResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequest
	err := tx.Model(&transactionworkshopentities.ServiceRequest{}).Where("service_request_system_number = ?", Id).First(&entity).Error

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

	// fetch data brand from external api
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

	// fetch data model from external api
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

	// fetch data variant from external api
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

	//fetch data company from external api
	CompanyUrl := config.EnvConfigs.GeneralServiceUrl + "company/" + strconv.Itoa(entity.CompanyId)
	var companyResponse transactionworkshoppayloads.CompanyResponse
	errCompany := utils.Get(CompanyUrl, &companyResponse, nil)
	if errCompany != nil {
		return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve company data from the external API",
			Err:        errCompany,
		}
	}

	// fetch data vehicle from external api
	// VehicleUrl := config.EnvConfigs.SalesServiceUrl + "/vehicle-master?page=0&limit=1&vehicle_id=" + strconv.Itoa(entity.VehicleId)
	// var vehicleResponse transactionworkshoppayloads.VehicleResponse
	// errVehicle := utils.Get(VehicleUrl, &vehicleResponse, nil)
	// if errVehicle != nil {
	// 	return transactionworkshoppayloads.ServiceRequestResponse{}, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    "Failed to retrieve vehicle data from the external API",
	// 		Err:        errVehicle,
	// 	}
	// }

	payload := transactionworkshoppayloads.ServiceRequestResponse{
		ServiceRequestSystemNumber: entity.ServiceRequestSystemNumber,
		ServiceRequestStatusId:     entity.ServiceRequestStatusId,
		ServiceRequestDate:         ServiceRequestDate,
		BrandId:                    entity.BrandId,
		BrandName:                  brandResponse.BrandName,
		ModelId:                    entity.ModelId,
		ModelName:                  modelResponse.ModelName,
		VariantId:                  entity.VariantId,
		VariantName:                variantResponse.VariantName,
		VehicleId:                  entity.VehicleId,
		CompanyId:                  entity.CompanyId,
		CompanyName:                companyResponse.CompanyName,
		DealerRepresentativeId:     entity.DealerRepresentativeId,
		ProfitCenterId:             entity.ProfitCenterId,
		WorkOrderSystemNumber:      entity.WorkOrderSystemNumber,
		BookingSystemNumber:        entity.BookingSystemNumber,
		EstimationSystemNumber:     entity.EstimationSystemNumber,
		ReferenceDocSystemNumber:   entity.ReferenceDocSystemNumber,
		ReplyId:                    entity.ReplyId,
		ReplyBy:                    entity.ReplyBy,
		ReplyDate:                  ReplyDate,
		ReplyRemark:                entity.ReplyRemark,
		ServiceCompanyId:           entity.ServiceCompanyId,
		ServiceDate:                serviceDate,
		ServiceRequestBy:           entity.ServiceRequestBy,
	}

	return payload, nil
}

func (s *ServiceRequestRepositoryImpl) New(tx *gorm.DB, request transactionworkshoppayloads.ServiceRequestSaveRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
	defaultServiceRequestDocumentNumber := ""
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
			ServiceRequestDocumentNumber: defaultServiceRequestDocumentNumber,
			ServiceRequestStatusId:       defaultWorkOrderStatusId,
			ServiceRequestDate:           currentDate,
			BrandId:                      request.BrandId,
			ModelId:                      request.ModelId,
			VariantId:                    request.VariantId,
			VehicleId:                    request.VehicleId,
			CompanyId:                    request.CompanyId,
			DealerRepresentativeId:       request.DealerRepresentativeId,
			ProfitCenterId:               request.ProfitCenterId,
			WorkOrderSystemNumber:        request.WorkOrderSystemNumber,
			BookingSystemNumber:          request.BookingSystemNumber,
			EstimationSystemNumber:       request.EstimationSystemNumber,
			ReferenceDocSystemNumber:     request.ReferenceDocSystemNumber,
			ReplyId:                      defaultReplyId,
			ServiceCompanyId:             request.ServiceCompanyId,
			ServiceDate:                  request.ServiceDate,
			ServiceRequestBy:             request.ServiceRequestBy,
			ReferenceTypeId:              ReferenceTypeId,
			ReferenceJobType:             jobType,
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
			ServiceRequestDocumentNumber: defaultServiceRequestDocumentNumber,
			ServiceRequestStatusId:       defaultWorkOrderStatusId,
			ServiceRequestDate:           currentDate,
			BrandId:                      request.BrandId,
			ModelId:                      request.ModelId,
			VariantId:                    request.VariantId,
			VehicleId:                    request.VehicleId,
			CompanyId:                    request.CompanyId,
			DealerRepresentativeId:       request.DealerRepresentativeId,
			ProfitCenterId:               request.ProfitCenterId,
			WorkOrderSystemNumber:        request.WorkOrderSystemNumber,
			BookingSystemNumber:          request.BookingSystemNumber,
			ReferenceDocSystemNumber:     request.ReferenceDocSystemNumber,
			ServiceRequestBy:             request.ServiceRequestBy,
			ServiceDate:                  request.ServiceDate,
			ReferenceTypeId:              ReferenceTypeId,
			ReferenceJobType:             jobType,
		}

	case "SO":
		jobType := ""
		entities = transactionworkshopentities.ServiceRequest{
			ServiceRequestDocumentNumber: defaultServiceRequestDocumentNumber,
			ServiceRequestStatusId:       defaultWorkOrderStatusId,
			ServiceRequestDate:           currentDate,
			BrandId:                      request.BrandId,
			ModelId:                      request.ModelId,
			VariantId:                    request.VariantId,
			VehicleId:                    request.VehicleId,
			CompanyId:                    request.CompanyId,
			DealerRepresentativeId:       request.DealerRepresentativeId,
			ProfitCenterId:               request.ProfitCenterId,
			EstimationSystemNumber:       request.EstimationSystemNumber,
			ReplyId:                      request.ReplyId,
			ServiceCompanyId:             request.ServiceCompanyId,
			ServiceDate:                  request.ServiceDate,
			ServiceRequestBy:             request.ServiceRequestBy,
			ReferenceTypeId:              ReferenceTypeId,
			ReferenceJobType:             jobType,
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

func (s *ServiceRequestRepositoryImpl) Save(tx *gorm.DB, Id int, request transactionworkshoppayloads.ServiceRequestSaveRequest) (transactionworkshopentities.ServiceRequest, *exceptions.BaseErrorResponse) {
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

	entity.BrandId = request.BrandId
	entity.ModelId = request.ModelId
	entity.VehicleId = request.VehicleId
	entity.CompanyId = request.CompanyId
	entity.DealerRepresentativeId = request.DealerRepresentativeId
	entity.ProfitCenterId = request.ProfitCenterId
	entity.WorkOrderSystemNumber = request.WorkOrderSystemNumber
	entity.BookingSystemNumber = request.BookingSystemNumber
	entity.EstimationSystemNumber = request.EstimationSystemNumber
	entity.ReferenceDocSystemNumber = request.ReferenceDocSystemNumber
	entity.ReplyId = request.ReplyId
	entity.ServiceCompanyId = request.ServiceCompanyId
	entity.ServiceDate = request.ServiceDate
	entity.ServiceRequestBy = request.ServiceRequestBy

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

	// Check if there are service request details with non-zero FrtQuantity
	var detailCount int64
	tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).
		Where("service_request_system_number = ? AND frt_quantity > 0", Id).
		Count(&detailCount)

	if detailCount == 0 {
		return false, "", &exceptions.BaseErrorResponse{Message: "Cannot submit service request ftr / qty must be > 0"}
	}

	if entity.ServiceRequestDocumentNumber == "" && entity.ServiceRequestStatusId == 1 {
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
	err := tx.Model(&transactionworkshopentities.ServiceRequest{}).Where("service_request_system_number = ?", Id).First(&entity).Error
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

	query := tx.Model(&transactionworkshopentities.ServiceRequestDetail{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	// fetch data Item from api external
	itemUrl := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(entities[0].OperationItemId)
	fmt.Println(entities[0].OperationItemId)
	errItem := utils.Get(itemUrl, &getItemResponse, nil)
	if errItem != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errItem,
		}
	}

	// fetch data Uom from api external
	uomUrl := config.EnvConfigs.AfterSalesServiceUrl + "unit-of-measurement/?page=0&limit=10&uom_id=" + strconv.Itoa(getItemResponse.UomId)
	errUom := utils.Get(uomUrl, &getUomItems, nil)
	if errUom != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUom,
		}
	}

	var ServiceRequestDetailResponses []map[string]interface{}

	for _, entity := range entities {
		ServiceRequestDetailResponse := map[string]interface{}{
			"uom_name":          getUomItems[0].UomName,
			"item_name":         getItemResponse.ItemName,
			"line_type_id":      entity.LineTypeId,
			"operation_item_id": entity.OperationItemId,
			"frt_quantity":      entity.FrtQuantity,
		}

		ServiceRequestDetailResponses = append(ServiceRequestDetailResponses, ServiceRequestDetailResponse)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(ServiceRequestDetailResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (s *ServiceRequestRepositoryImpl) GetServiceDetailById(tx *gorm.DB, Id int) (transactionworkshoppayloads.ServiceDetailResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequestDetail
	err := tx.Model(&transactionworkshopentities.ServiceRequestDetail{}).Where("service_request_detail_id = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.ServiceDetailResponse{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Data not found"}
		}
		return transactionworkshoppayloads.ServiceDetailResponse{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Err: err}
	}

	payload := transactionworkshoppayloads.ServiceDetailResponse{
		ServiceRequestDetailId:     entity.ServiceRequestDetailId,
		ServiceRequestId:           entity.ServiceRequestId,
		ServiceRequestSystemNumber: entity.ServiceRequestSystemNumber,
		LineTypeId:                 entity.LineTypeId,
		OperationItemId:            entity.OperationItemId,
		FrtQuantity:                entity.FrtQuantity,
	}

	return payload, nil
}

func (s *ServiceRequestRepositoryImpl) AddServiceDetail(tx *gorm.DB, id int, request transactionworkshoppayloads.ServiceDetailSaveRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse) {

	entity := transactionworkshopentities.ServiceRequestDetail{
		ServiceRequestId:           request.ServiceRequestId,
		ServiceRequestSystemNumber: request.ServiceRequestSystemNumber,
		LineTypeId:                 request.LineTypeId,
		OperationItemId:            request.OperationItemId,
		FrtQuantity:                request.FrtQuantity,
		ReferenceDocSystemNumber:   0,
		ReferenceDocId:             0,
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

func (s *ServiceRequestRepositoryImpl) UpdateServiceDetail(tx *gorm.DB, Id int, DetailId int, request transactionworkshoppayloads.ServiceDetailSaveRequest) (transactionworkshopentities.ServiceRequestDetail, *exceptions.BaseErrorResponse) {

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
