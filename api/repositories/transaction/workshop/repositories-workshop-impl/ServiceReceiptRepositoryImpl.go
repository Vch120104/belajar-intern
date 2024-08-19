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
	"math"
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

		// Fetch data from external APIs
		BrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(ServiceReceiptReq.BrandId)
		var brandResponses transactionworkshoppayloads.WorkOrderVehicleBrand
		if err := utils.Get(BrandUrl, &brandResponses, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve brand data from the external API",
				Err:        err,
			}
		}

		ModelUrl := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(ServiceReceiptReq.ModelId)
		var modelResponses transactionworkshoppayloads.WorkOrderVehicleModel
		if err := utils.Get(ModelUrl, &modelResponses, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve model data from the external API",
				Err:        err,
			}
		}

		VariantUrl := config.EnvConfigs.SalesServiceUrl + "unit-variant/" + strconv.Itoa(ServiceReceiptReq.VariantId)
		var variantResponses transactionworkshoppayloads.WorkOrderVehicleVariant
		if err := utils.Get(VariantUrl, &variantResponses, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve variant data from the external API",
				Err:        err,
			}
		}

		ColourUrl := config.EnvConfigs.SalesServiceUrl + "unit-color-dropdown/" + strconv.Itoa(ServiceReceiptReq.BrandId)
		var colourResponses []transactionworkshoppayloads.WorkOrderVehicleColour
		if err := utils.GetArray(ColourUrl, &colourResponses, nil); err != nil || len(colourResponses) == 0 {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve colour data from the external API",
				Err:        err,
			}
		}

		CompanyUrl := config.EnvConfigs.GeneralServiceUrl + "companies-redis?company_id=" + strconv.Itoa(ServiceReceiptReq.CompanyId)
		var companyResponses []transactionworkshoppayloads.CompanyResponse
		if err := utils.GetArray(CompanyUrl, &companyResponses, nil); err != nil || len(companyResponses) == 0 {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve company data from the external API",
				Err:        err,
			}
		}

		VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(ServiceReceiptReq.VehicleId)
		var vehicleResponses transactionworkshoppayloads.VehicleResponse
		if err := utils.Get(VehicleUrl, &vehicleResponses, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve vehicle data from the external API",
				Err:        err,
			}
		}

		ServiceRequestStatusURL := config.EnvConfigs.AfterSalesServiceUrl + "service-request/dropdown-status?service_request_status_id=" + strconv.Itoa(ServiceReceiptReq.ServiceRequestStatusId)
		var statusResponses []transactionworkshoppayloads.ServiceRequestStatusResponse
		if err := utils.GetArray(ServiceRequestStatusURL, &statusResponses, nil); err != nil || len(statusResponses) == 0 {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve status service request data from the external API",
				Err:        err,
			}
		}

		ServiceReceiptRes = transactionworkshoppayloads.ServiceReceiptResponse{
			ServiceRequestSystemNumber:   ServiceReceiptReq.ServiceRequestSystemNumber,
			ServiceRequestDocumentNumber: ServiceReceiptReq.ServiceRequestDocumentNumber,
			ServiceRequestDate:           ServiceReceiptReq.ServiceRequestDate.Format("2006-01-02 15:04:05"),
			ServiceRequestBy:             ServiceReceiptReq.ServiceRequestBy,
			ServiceRequestStatusName:     statusResponses[0].ServiceRequestStatusName,
			BrandName:                    brandResponses.BrandName,
			ModelName:                    modelResponses.ModelName,
			VariantName:                  variantResponses.VariantName,
			VariantColourName:            colourResponses[0].VariantColourName,
			VehicleCode:                  vehicleResponses.VehicleCode,
			VehicleTnkb:                  vehicleResponses.VehicleTnkb,
			CompanyName:                  companyResponses[0].CompanyName,
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
			"work_order_no":                   response.WorkOrderSystemNumber,
			"booking_no":                      response.BookingSystemNumber,
			"ref_doc_no":                      response.ReferenceDocSystemNumber,
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
	CompanyUrl := config.EnvConfigs.GeneralServiceUrl + "companies-redis?company_id=" + strconv.Itoa(entity.CompanyId)
	var companyResponses []transactionworkshoppayloads.CompanyResponse
	errCompany := utils.GetArray(CompanyUrl, &companyResponses, nil)
	if errCompany != nil || len(companyResponses) == 0 {
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
		Select("service_request_detail_id, service_request_id, service_request_system_number, line_type_id, operation_item_id, frt_quantity, reference_doc_system_number, reference_doc_id").
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
	ProfitCenterUrl := config.EnvConfigs.GeneralServiceUrl + "profit-center/" + strconv.Itoa(entity.ProfitCenterId)
	var profitCenterResponses transactionworkshoppayloads.ProfitCenter
	errProfitCenter := utils.Get(ProfitCenterUrl, &profitCenterResponses, nil)
	if errProfitCenter != nil {
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.ServiceReceiptResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve dealer representative data from the external API",
			Err:        errDealerRepresentative,
		}
	}
	totalRows := 0
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
		VehicleCode:                  vehicleResponses.VehicleCode,
		VehicleTnkb:                  vehicleResponses.VehicleTnkb,
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
		ReferenceDocNumber:           "", //entity.ReferenceDocNumber,
		ReferenceDocDate:             "", //entity.ReferenceDocDate,
		ReplyId:                      entity.ReplyId,
		ReplyBy:                      entity.ReplyBy,
		ReplyDate:                    ReplyDate,
		ReplyRemark:                  entity.ReplyRemark,
		ServiceCompanyId:             entity.ServiceCompanyId,
		ServiceCompanyName:           companyResponses[0].CompanyName,
		ServiceDate:                  serviceDate,
		ServiceRequestBy:             entity.ServiceRequestBy,
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
