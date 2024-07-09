package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"errors"
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

	// Add the additional where condition
	whereQuery = whereQuery.Where("service_request_system_number != 0 AND service_request_status_id = 2")

	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	var convertedResponses []transactionworkshoppayloads.ServiceReceiptResponse
	for rows.Next() {
		var (
			ServiceReceiptReq transactionworkshoppayloads.ServiceReceiptNew
			ServiceReceiptRes transactionworkshoppayloads.ServiceReceiptResponse
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

		ServiceReceiptRes = transactionworkshoppayloads.ServiceReceiptResponse{
			ServiceRequestSystemNumber:   ServiceReceiptRes.ServiceRequestSystemNumber,
			ServiceRequestDocumentNumber: ServiceReceiptRes.ServiceRequestDocumentNumber,
			//ServiceRequestDate:           ServiceReceiptRes.ServiceRequestDate.Format("2006-01-02 15:04:05"),
			ServiceRequestBy:         ServiceReceiptRes.ServiceRequestBy,
			ServiceRequestStatusId:   ServiceReceiptRes.ServiceRequestStatusId,
			BrandId:                  ServiceReceiptRes.BrandId,
			ModelId:                  ServiceReceiptRes.ModelId,
			VehicleId:                ServiceReceiptRes.VehicleId,
			CompanyId:                ServiceReceiptRes.CompanyId,
			DealerRepresentativeId:   ServiceReceiptRes.DealerRepresentativeId,
			ProfitCenterId:           ServiceReceiptRes.ProfitCenterId,
			WorkOrderSystemNumber:    ServiceReceiptRes.WorkOrderSystemNumber,
			BookingSystemNumber:      ServiceReceiptRes.BookingSystemNumber,
			EstimationSystemNumber:   ServiceReceiptRes.EstimationSystemNumber,
			ReferenceDocSystemNumber: ServiceReceiptRes.ReferenceDocSystemNumber,
			ReplyId:                  ServiceReceiptRes.ReplyId,
			ServiceCompanyId:         ServiceReceiptRes.ServiceCompanyId,
			//ServiceDate:                  ServiceReceiptRes.ServiceDate.Format("2006-01-02 15:04:05"),
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

func (s *ServiceReceiptRepositoryImpl) GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.ServiceReceiptResponse, *exceptions.BaseErrorResponse) {
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

	payload := transactionworkshoppayloads.ServiceReceiptResponse{
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

func (s *ServiceReceiptRepositoryImpl) Save(tx *gorm.DB, Id int, request transactionworkshoppayloads.ServiceReceiptSaveRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ServiceRequest
	currentDate := time.Now()

	// cek service request system number
	err := tx.Model(&transactionworkshopentities.ServiceRequest{}).Where("service_request_system_number = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Data not found"}
		}
		return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Err: err}
	}

	// Check current service request status
	if entity.ServiceRequestStatusId != 1 {
		return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusBadRequest, Message: "Service request status is not in draft"}
	}

	// Check if ServiceDate is less than currentDate
	if request.ServiceDate.Before(currentDate) {
		return false, &exceptions.BaseErrorResponse{
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save the service request",
			Err:        err}
	}

	return true, nil
}
