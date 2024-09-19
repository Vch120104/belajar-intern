package transactionworkshoprepositoryimpl

////  NOTES  ////
import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	masterrepository "after-sales/api/repositories/master"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
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

type WorkOrderRepositoryImpl struct {
	lookupRepo masterrepository.LookupRepository
}

func OpenWorkOrderRepositoryImpl() transactionworkshoprepository.WorkOrderRepository {
	lookupRepo := masterrepositoryimpl.StartLookupRepositoryImpl()
	return &WorkOrderRepositoryImpl{
		lookupRepo: lookupRepo,
	}
}

// uspg_wtWorkOrder0_Insert
// IF @Option = 0
// --USE FOR : * INSERT NEW DATA
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (r *WorkOrderRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tableStruct := transactionworkshoppayloads.WorkOrderGetAllRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Add the additional where condition
	whereQuery = whereQuery.Where("service_request_system_number = 0 AND booking_system_number = 0 AND estimation_system_number = 0")

	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	var convertedResponses []transactionworkshoppayloads.WorkOrderGetAllResponse

	for rows.Next() {

		var (
			workOrderReq transactionworkshoppayloads.WorkOrderGetAllRequest
			workOrderRes transactionworkshoppayloads.WorkOrderGetAllResponse
		)

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
			&workOrderReq.RepeatedJob,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch data brand from external services
		BrandURL := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(workOrderReq.BrandId)
		var getBrandResponse transactionworkshoppayloads.WorkOrderVehicleBrand
		if err := utils.Get(BrandURL, &getBrandResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch brand data from external service",
				Err:        err,
			}
		}

		// Fetch data model from external services
		ModelURL := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(workOrderReq.ModelId)
		var getModelResponse transactionworkshoppayloads.WorkOrderVehicleModel
		if err := utils.Get(ModelURL, &getModelResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch model data from external service",
				Err:        err,
			}
		}

		VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100&vehicle_id=" + strconv.Itoa(workOrderReq.VehicleId)
		var vehicleResponses []transactionworkshoppayloads.VehicleResponse
		errVehicle := utils.GetArray(VehicleUrl, &vehicleResponses, nil)
		if errVehicle != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve vehicle data from the external API",
				Err:        errVehicle,
			}
		}

		if len(vehicleResponses) == 0 {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "No vehicle data found",
			}
		}

		CustomerURL := config.EnvConfigs.GeneralServiceUrl + "customer-detail/" + strconv.Itoa(workOrderReq.CustomerId)
		var getCustomerResponse transactionworkshoppayloads.CustomerResponse
		if err := utils.Get(CustomerURL, &getCustomerResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch customer data from external service",
				Err:        err,
			}
		}

		WorkOrderTypeURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-type?work_order_type_id=" + strconv.Itoa(workOrderReq.WorkOrderTypeId)
		var getWorkOrderTypeResponses []transactionworkshoppayloads.WorkOrderTypeResponse
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

		WorkOrderStatusURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-status?work_order_status_id=" + strconv.Itoa(workOrderReq.StatusId)
		var getWorkOrderStatusResponses []transactionworkshoppayloads.WorkOrderStatusResponse
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

		workOrderRes = transactionworkshoppayloads.WorkOrderGetAllResponse{
			WorkOrderDocumentNumber: workOrderReq.WorkOrderDocumentNumber,
			WorkOrderSystemNumber:   workOrderReq.WorkOrderSystemNumber,
			WorkOrderDate:           workOrderReq.WorkOrderDate,
			FormattedWorkOrderDate:  utils.FormatRFC3339(workOrderReq.WorkOrderDate), // Use RFC3339 format
			WorkOrderTypeId:         workOrderReq.WorkOrderTypeId,
			WorkOrderTypeName:       workOrderTypeName,
			BrandId:                 workOrderReq.BrandId,
			BrandName:               getBrandResponse.BrandName,
			VehicleCode:             vehicleResponses[0].VehicleCode,
			VehicleTnkb:             vehicleResponses[0].VehicleTnkb,
			ModelId:                 workOrderReq.ModelId,
			ModelName:               getModelResponse.ModelName,
			VehicleId:               workOrderReq.VehicleId,
			CustomerId:              workOrderReq.CustomerId,
			StatusId:                workOrderReq.StatusId,
			StatusName:              workOrderStatusName,
			RepeatedJob:             workOrderReq.RepeatedJob,
		}

		convertedResponses = append(convertedResponses, workOrderRes)
	}

	var mapResponses []map[string]interface{}

	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_document_number":  response.WorkOrderDocumentNumber,
			"work_order_system_number":    response.WorkOrderSystemNumber,
			"work_order_date":             response.FormattedWorkOrderDate, // Use formatted date
			"work_order_type_id":          response.WorkOrderTypeId,
			"work_order_type_description": response.WorkOrderTypeName,
			"brand_id":                    response.BrandId,
			"brand_name":                  response.BrandName,
			"model_id":                    response.ModelId,
			"model_name":                  response.ModelName,
			"vehicle_id":                  response.VehicleId,
			"vehicle_chassis_number":      response.VehicleCode,
			"vehicle_tnkb":                response.VehicleTnkb,
			"work_order_status_id":        response.StatusId,
			"work_order_status_name":      response.StatusName,
			"repeated_system_number":      response.RepeatedJob,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

// uspg_wtWorkOrder0_Insert
// IF @Option = 0
// --USE FOR : * INSERT NEW DATA
func (r *WorkOrderRepositoryImpl) New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderNormalRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse) {

	// Default values
	// 1:Draft, 2:New, 3:Ready, 4:On Going, 5:Stop, 6:QC Pass, 7:Cancel, 8:Closed
	// 1:Normal, 2:Campaign, 3:Affiliated, 4:Repeat Job
	defaultWorkOrderDocumentNumber := ""
	defaultServiceAdvisorId := 1 // Default advisor ID
	defaultCPCcode := "00002"    // Default CPC code 00002 for workshop
	workOrderTypeId := 1         // Default work order type ID 1 for normal

	// Validation: request date
	currentDate := time.Now()
	requestDate := request.WorkOrderArrivalTime.Truncate(24 * time.Hour)
	if requestDate.Before(currentDate) || requestDate.After(currentDate) {
		request.WorkOrderArrivalTime = currentDate
	}

	// Check if the CompanyId is provided
	if request.CompanyId == 0 {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Company ID is required",
			Err:        errors.New("parameter has lost session, please refresh the data"),
		}
	}

	// Validation: WorkOrderTypeId
	if request.CampaignId > 0 {
		// Campaign
		workOrderTypeId = 2
	} else if (request.PDISystemNumber > 0 && request.ServiceRequestSystemNumber == 0) ||
		(request.PDISystemNumber == 0 && request.ServiceRequestSystemNumber > 0) {
		// Affiliated (PDI/ServiceRequestSystemNumber)
		workOrderTypeId = 3
	} else if request.RepeatedSystemNumber > 0 {
		// Repeat Job
		workOrderTypeId = 4
	}

	// fetch vehicle
	vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100&vehicle_id=" + strconv.Itoa(request.VehicleId)
	var vehicleResponses []transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.GetArray(vehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	// Create WorkOrder entity
	entitieswo := transactionworkshopentities.WorkOrder{
		// Default values
		WorkOrderDocumentNumber:    defaultWorkOrderDocumentNumber,
		WorkOrderStatusId:          utils.WoStatDraft,
		WorkOrderDate:              currentDate,
		CPCcode:                    defaultCPCcode,
		ServiceAdvisor:             defaultServiceAdvisorId,
		WorkOrderTypeId:            workOrderTypeId,
		BookingSystemNumber:        request.BookingSystemNumber,
		EstimationSystemNumber:     request.EstimationSystemNumber,
		ServiceRequestSystemNumber: request.ServiceRequestSystemNumber,
		PDISystemNumber:            request.PDISystemNumber,
		RepeatedSystemNumber:       request.RepeatedSystemNumber,
		ServiceSite:                "OD - Service On Dealer",
		VehicleChassisNumber:       vehicleResponses[0].VehicleCode,

		// Provided values
		BrandId:                  request.BrandId,
		ModelId:                  request.ModelId,
		VariantId:                request.VariantId,
		VehicleId:                request.VehicleId,
		CustomerId:               request.CustomerId,
		BillableToId:             request.BilltoCustomerId,
		FromEra:                  request.FromEra,
		QueueNumber:              request.QueueSystemNumber,
		ArrivalTime:              request.WorkOrderArrivalTime,
		ServiceMileage:           request.WorkOrderCurrentMileage,
		Storing:                  request.Storing,
		Remark:                   request.WorkOrderRemark,
		ProfitCenterId:           request.WorkOrderProfitCenter,
		CostCenterId:             request.DealerRepresentativeId,
		CampaignId:               request.CampaignId,
		CompanyId:                request.CompanyId,
		CPTitlePrefix:            request.Titleprefix,
		ContactPersonName:        request.NameCust,
		ContactPersonPhone:       request.PhoneCust,
		ContactPersonMobile:      request.MobileCust,
		ContactPersonContactVia:  request.ContactVia,
		EraNumber:                request.WorkOrderEraNo,
		EraExpiredDate:           request.WorkOrderEraExpiredDate,
		InsurancePolicyNumber:    request.WorkOrderInsurancePolicyNo,
		InsuranceExpiredDate:     request.WorkOrderInsuranceExpiredDate,
		InsuranceClaimNumber:     request.WorkOrderInsuranceClaimNo,
		InsurancePersonInCharge:  request.WorkOrderInsurancePic,
		InsuranceOwnRisk:         request.WorkOrderInsuranceOwnRisk,
		InsuranceWorkOrderNumber: request.WorkOrderInsuranceWONumber,
		EstTime:                  request.EstimationDuration,
		CustomerExpress:          request.CustomerExpress,
		LeaveCar:                 request.LeaveCar,
		CarWash:                  request.CarWash,
		PromiseDate:              request.PromiseDate,
		PromiseTime:              request.PromiseTime,
		FSCouponNo:               request.FSCouponNo,
		Notes:                    request.Notes,
		Suggestion:               request.Suggestion,
		DPAmount:                 request.DownpaymentAmount,
	}

	if err := tx.Create(&entitieswo).Error; err != nil {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create work order",
			Err:        err,
		}
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Memperbarui status pemesanan dan estimasi jika Booking_System_No atau Estim_System_No tidak nol

	var (
		batchSystemNo       int
		bookingStatusClosed int
		bookingSystemNo     = request.BookingSystemNumber
		estimationSystemNo  = request.EstimationSystemNumber
	)
	bookingStatusClosed = 8

	// Update related records if necessary
	if bookingSystemNo != 0 {
		if batchSystemNo == 0 {
			var batchSystemNoResult struct {
				BatchSystemNo int
			}
			if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Select("batch_system_number").
				Where("booking_system_number = ?", bookingSystemNo).
				Scan(&batchSystemNoResult).Error; err != nil {
				return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve batch system number from the database",
					Err:        err,
				}
			}
			batchSystemNo = batchSystemNoResult.BatchSystemNo
		}

		// Update BOOKING_STATUS
		if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
			Where("booking_system_number = ?", bookingSystemNo).
			Update("booking_status_id", bookingStatusClosed).Error; err != nil {
			return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update booking status",
				Err:        err,
			}
		}
	}

	if estimationSystemNo != 0 {
		if batchSystemNo == 0 {
			var batchSystemNoResult struct {
				BatchSystemNo int
			}
			if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Select("batch_system_number").
				Where("estimation_system_number = ?", estimationSystemNo).
				Scan(&batchSystemNoResult).Error; err != nil {
				return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve batch system number from the database",
					Err:        err,
				}
			}
			batchSystemNo = batchSystemNoResult.BatchSystemNo
		}

		// Update ESTIM_STATUS
		if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
			Where("estimation_system_number = ?", estimationSystemNo).
			Update("estimation_status_id", bookingStatusClosed).Error; err != nil {
			return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update estimation status",
				Err:        err,
			}
		}
	}

	// Update BATCH_STATUS
	if batchSystemNo != 0 {
		if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
			Where("batch_system_number = ?", batchSystemNo).
			Update("batch_status_id", bookingStatusClosed).Error; err != nil {
			return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update batch status",
				Err:        err,
			}
		}
	}

	return entitieswo, nil
}

func (r *WorkOrderRepositoryImpl) GetById(tx *gorm.DB, Id int, pagination pagination.Pagination) (transactionworkshoppayloads.WorkOrderResponseDetail, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
			}
		}
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order from the database",
			Err:        err,
		}
	}

	// Fetch data brand from external API
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entity.BrandId)
	var brandResponse transactionworkshoppayloads.WorkOrderVehicleBrand
	errBrand := utils.Get(brandUrl, &brandResponse, nil)
	if errBrand != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        errVariant,
		}
	}

	// Fetch data colour from external API
	colourUrl := config.EnvConfigs.SalesServiceUrl + "unit-color-dropdown/" + strconv.Itoa(entity.BrandId)
	var colourResponses []transactionworkshoppayloads.WorkOrderVehicleColour
	errColour := utils.GetArray(colourUrl, &colourResponses, nil)
	if errColour != nil || len(colourResponses) == 0 {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve colour data from the external API",
			Err:        errColour,
		}
	}

	// Fetch data vehicle from external API
	vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100&vehicle_id=" + strconv.Itoa(entity.VehicleId)
	var vehicleResponses []transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.GetArray(vehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	// Fetch workorder details with pagination
	var workorderDetails []transactionworkshoppayloads.WorkOrderDetailResponse
	query := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ?", Id).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit())
	errWorkOrderDetails := query.Find(&workorderDetails).Error
	if errWorkOrderDetails != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order details from the database",
			Err:        errWorkOrderDetails,
		}
	}

	// Fetch work order services
	var workorderServices []transactionworkshoppayloads.WorkOrderServiceResponse
	if err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ?", Id).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderServices).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order services from the database",
			Err:        err,
		}
	}

	// Fetch work order vehicles
	var workorderVehicle []transactionworkshoppayloads.WorkOrderServiceVehicleResponse
	if err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ?", Id).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderVehicle).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order vehicles from the database",
			Err:        err,
		}
	}

	// Fetch work order campaigns
	var workorderCampaigns []transactionworkshoppayloads.WorkOrderCampaignResponse
	if err := tx.Model(&masterentities.CampaignMaster{}).
		Where("campaign_id = ? and campaign_id != 0", entity.CampaignId).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderCampaigns).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order campaigns from the database",
			Err:        err,
		}
	}

	// Fetch work order agreements
	var workorderAgreements []transactionworkshoppayloads.WorkOrderGeneralRepairAgreementResponse
	if err := tx.Model(&masterentities.Agreement{}).
		Where("agreement_id = ? and agreement_id != 0", entity.AgreementGeneralRepairId).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderAgreements).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order agreements from the database",
			Err:        err,
		}
	}

	// Fetch work order bookings
	var workorderBookings []transactionworkshoppayloads.WorkOrderBookingsResponse
	if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
		Where("booking_system_number = ? and booking_system_number != 0", entity.BookingSystemNumber).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderBookings).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order bookings from the database",
			Err:        err,
		}
	}

	// // Fetch work order estimations
	var workorderEstimations []transactionworkshoppayloads.WorkOrderEstimationsResponse
	if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
		Where("estimation_system_number = ? and estimation_system_number != 0", entity.EstimationSystemNumber).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderEstimations).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order estimations from the database",
			Err:        err,
		}
	}

	// // Fetch work order contracts
	var workorderContracts []transactionworkshoppayloads.WorkOrderContractsResponse
	if err := tx.Model(&transactionworkshopentities.ContractService{}).
		Where("contract_service_system_number = ? and contract_service_system_number != 0", entity.ContractServiceSystemNumber).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderContracts).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order contracts from the database",
			Err:        err,
		}
	}

	// Fetch work order users
	var workorderUsers []transactionworkshoppayloads.WorkOrderCurrentUserResponse
	if err := tx.Table("dms_microservices_general_dev.dbo.mtr_customer AS c").
		Select(`
		c.customer_id AS customer_id,
		c.customer_name AS customer_name,
		c.customer_code AS customer_code,
		c.id_address_id AS address_id,
		a.address_street_1 AS address_street_1,
		a.address_street_2 AS address_street_2,
		a.address_street_3 AS address_street_3,
		a.village_id AS village_id,
		v.village_name AS village_name,
		v.district_id AS district_id,
		d.district_name AS district_name,
		d.city_id AS city_id,
		ct.city_name AS city_name,
		ct.province_id AS province_id,
		p.province_name AS province_name,
		v.village_zip_code AS zip_code,
		td.npwp_no AS current_user_npwp
	`).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_address AS a ON c.id_address_id = a.address_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_village AS v ON a.village_id = v.village_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_district AS d ON v.district_id = d.district_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_city AS ct ON d.city_id = ct.city_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_province AS p ON ct.province_id = p.province_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_tax_data AS td ON c.tax_customer_id = td.tax_id").
		Where("c.customer_id = ?", entity.CustomerId).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderUsers).Error; err != nil {
		fmt.Println("Error executing query:", err)
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order users from the database",
			Err:        err,
		}
	}

	// Fetch work order detail vehicles
	var workorderVehicleDetails []transactionworkshoppayloads.WorkOrderVehicleDetailResponse
	if err := tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle AS v").
		Select(`
		v.vehicle_id AS vehicle_id,
        v.vehicle_chassis_number AS vehicle_chassis_number,
		vrc.vehicle_registration_certificate_tnkb AS vehicle_registration_certificate_tnkb,
		vrc.vehicle_registration_certificate_owner_name AS vehicle_registration_certificate_owner_name,
		v.vehicle_production_year AS vehicle_production_year,
		CONCAT(vv.variant_code , ' - ', vv.variant_description) AS vehicle_variant,
		v.option_id AS vehicle_option,
		CONCAT(vm.colour_code , ' - ', vm.colour_commercial_name) AS vehicle_colour,
		v.vehicle_sj_date AS vehicle_sj_date,
        v.vehicle_last_service_date AS vehicle_last_service_date,
        v.vehicle_last_km AS vehicle_last_km
		`).
		Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_vehicle_registration_certificate AS vrc ON v.vehicle_id = vrc.vehicle_id").
		Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_unit_variant AS vv ON v.vehicle_variant_id = vv.variant_id").
		Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_colour AS vm ON v.vehicle_colour_id = vm.colour_id").
		Where("v.vehicle_id = ? AND v.vehicle_brand_id = ? and v.vehicle_variant_id = ?", entity.VehicleId, entity.BrandId, entity.VariantId).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderVehicleDetails).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order vehicles from the database",
			Err:        err,
		}
	}

	// Fetch work order stnk
	var workorderStnk []transactionworkshoppayloads.WorkOrderStnkResponse
	if err := tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle_registration_certificate").
		Select(`
		vehicle_registration_certificate_id AS stnk_id,
		vehicle_registration_certificate_owner_name AS stnk_name
		`).
		Where("vehicle_id = ? ", entity.VehicleId).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderStnk).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order stnk from the database",
			Err:        err,
		}
	}

	// Fetch work order billings
	var workorderBillings []transactionworkshoppayloads.WorkOrderBillingResponse
	if err := tx.Table("dms_microservices_general_dev.dbo.mtr_customer AS c").
		Select(`
		c.customer_id AS bill_to_id,
		c.customer_name AS bill_to_name,
		c.customer_code AS bill_to_code,
		c.id_address_id AS address_id,
		a.address_street_1 AS address_street_1,
		a.address_street_2 AS address_street_2,
		a.address_street_3 AS address_street_3,
		v.village_name AS bill_to_village,
		d.district_name AS bill_to_district,
		ct.city_name AS bill_to_city,
		p.province_name AS bill_to_province,
		v.village_zip_code AS bill_to_zip_code,
		c.customer_mobile_phone AS bill_to_phone,
		c.home_fax_no AS bill_to_fax,
		td.npwp_no AS bill_to_npwp
	`).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_address AS a ON c.id_address_id = a.address_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_village AS v ON a.village_id = v.village_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_district AS d ON v.district_id = d.district_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_city AS ct ON d.city_id = ct.city_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_province AS p ON ct.province_id = p.province_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_tax_data AS td ON c.tax_customer_id = td.tax_id").
		Where("c.customer_id = ?", entity.CustomerId).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderBillings).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order billings from the database",
			Err:        err,
		}
	}

	// fetch data status work order
	WorkOrderStatusURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-status?work_order_status_id=" + strconv.Itoa(entity.WorkOrderStatusId)
	var getWorkOrderStatusResponses []transactionworkshoppayloads.WorkOrderStatusResponse // Use slice of WorkOrderStatusResponse
	if err := utils.Get(WorkOrderStatusURL, &getWorkOrderStatusResponses, nil); err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order status data from external service",
			Err:        err,
		}
	}

	// fetch data type work order
	WorkOrderTypeURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-type?work_order_type_id=" + strconv.Itoa(entity.WorkOrderTypeId)
	//fmt.Println("Fetching Work Order Type data from:", WorkOrderTypeURL)
	var getWorkOrderTypeResponses []transactionworkshoppayloads.WorkOrderTypeResponse // Use slice of WorkOrderTypeResponse
	if err := utils.Get(WorkOrderTypeURL, &getWorkOrderTypeResponses, nil); err != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order type data from external service",
			Err:        err,
		}
	}

	var workOrderTypeName string
	if len(getWorkOrderTypeResponses) > 0 {
		workOrderTypeName = getWorkOrderTypeResponses[0].WorkOrderTypeName
	}

	payload := transactionworkshoppayloads.WorkOrderResponseDetail{
		WorkOrderSystemNumber:         entity.WorkOrderSystemNumber,
		WorkOrderDate:                 entity.WorkOrderDate,
		WorkOrderDocumentNumber:       entity.WorkOrderDocumentNumber,
		WorkOrderTypeId:               entity.WorkOrderTypeId,
		WorkOrderTypeName:             workOrderTypeName,
		WorkOrderStatusId:             entity.WorkOrderStatusId,
		WorkOrderStatusName:           getWorkOrderStatusResponses[0].WorkOrderStatusName,
		ServiceAdvisorId:              entity.ServiceAdvisor,
		BrandId:                       entity.BrandId,
		BrandName:                     brandResponse.BrandName,
		ModelId:                       entity.ModelId,
		ModelName:                     modelResponse.ModelName,
		VariantId:                     entity.VariantId,
		VariantName:                   variantResponse.VariantName,
		VehicleId:                     entity.VehicleId,
		VehicleCode:                   vehicleResponses[0].VehicleCode,
		VehicleTnkb:                   vehicleResponses[0].VehicleTnkb,
		CustomerId:                    entity.CustomerId,
		ServiceSite:                   entity.ServiceSite,
		BilltoCustomerId:              entity.BillableToId,
		CampaignId:                    entity.CampaignId,
		FromEra:                       entity.FromEra,
		WorkOrderEraNo:                entity.EraNumber,
		Storing:                       entity.Storing,
		WorkOrderCurrentMileage:       entity.ServiceMileage,
		WorkOrderProfitCenterId:       entity.ProfitCenterId,
		AgreementId:                   entity.AgreementBodyRepairId,
		BoookingId:                    entity.BookingSystemNumber,
		EstimationId:                  entity.EstimationSystemNumber,
		ContractSystemNumber:          entity.ContractServiceSystemNumber,
		QueueSystemNumber:             entity.QueueNumber,
		WorkOrderArrivalTime:          entity.ArrivalTime,
		WorkOrderRemark:               entity.Remark,
		DealerRepresentativeId:        entity.CostCenterId,
		CompanyId:                     entity.CompanyId,
		Titleprefix:                   entity.CPTitlePrefix,
		NameCust:                      entity.ContactPersonName,
		PhoneCust:                     entity.ContactPersonPhone,
		MobileCust:                    entity.ContactPersonMobile,
		MobileCustAlternative:         entity.ContactPersonMobileAlternative,
		MobileCustDriver:              entity.ContactPersonMobileDriver,
		ContactVia:                    entity.ContactPersonContactVia,
		WorkOrderInsurancePolicyNo:    entity.InsurancePolicyNumber,
		WorkOrderInsuranceClaimNo:     entity.InsuranceClaimNumber,
		WorkOrderInsuranceExpiredDate: entity.InsuranceExpiredDate,
		WorkOrderEraExpiredDate:       entity.EraExpiredDate,
		PromiseDate:                   entity.PromiseDate,
		PromiseTime:                   entity.PromiseTime,
		EstimationDuration:            entity.EstTime,
		WorkOrderInsuranceOwnRisk:     entity.InsuranceOwnRisk,
		WorkOrderInsurancePic:         entity.InsurancePersonInCharge,
		WorkOrderInsuranceWONumber:    entity.InsuranceWorkOrderNumber,
		CustomerExpress:               entity.CustomerExpress,
		LeaveCar:                      entity.LeaveCar,
		CarWash:                       entity.CarWash,
		FSCouponNo:                    entity.FSCouponNo,
		Notes:                         entity.Notes,
		Suggestion:                    entity.Suggestion,
		DownpaymentAmount:             entity.DPAmount,
		WorkOrderCampaign: transactionworkshoppayloads.WorkOrderCampaignDetail{
			DataCampaign: workorderCampaigns,
		},
		WorkOrderGeneralRepairAgreement: transactionworkshoppayloads.WorkOrderGeneralRepairAgreement{
			DataAgreement: workorderAgreements,
		},
		WorkOrderBooking: transactionworkshoppayloads.WorkOrderBookingDetail{
			DataBooking: workorderBookings,
		},
		WorkOrderEstimation: transactionworkshoppayloads.WorkOrderEstimationDetail{
			DataEstimation: workorderEstimations,
		},
		WorkOrderContract: transactionworkshoppayloads.WorkOrderContractDetail{
			DataContract: workorderContracts,
		},
		WorkOrderCurrentUserDetail: transactionworkshoppayloads.WorkOrderCurrentUserDetail{
			DataCurrentUser: workorderUsers,
		},
		WorkOrderVehicleDetail: transactionworkshoppayloads.WorkOrderVehicleDetail{
			DataVehicle: workorderVehicleDetails,
		},
		WorkOrderStnkDetail: transactionworkshoppayloads.WorkOrderStnkDetail{
			DataStnk: workorderStnk,
		},
		WorkOrderBillingDetail: transactionworkshoppayloads.WorkOrderBillingDetail{
			DataBilling: workorderBillings,
		},
		WorkOrderDetailService: transactionworkshoppayloads.WorkOrderDetailsResponseRequest{
			DataRequest: workorderServices,
		},
		WorkOrderDetailVehicle: transactionworkshoppayloads.WorkOrderDetailsResponseVehicle{
			DataVehicle: workorderVehicle,
		},
		WorkOrderDetails: transactionworkshoppayloads.WorkOrderDetailsResponse{
			Page:       pagination.GetPage(),
			Limit:      pagination.GetLimit(),
			TotalPages: pagination.TotalPages,
			TotalRows:  int(pagination.TotalRows), // Convert int64 to int
			Data:       workorderDetails,
		},
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderNormalSaveRequest, workOrderId int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	entity.BillableToId = request.BilltoCustomerId
	entity.FromEra = request.FromEra
	entity.QueueNumber = request.QueueSystemNumber
	entity.ArrivalTime = request.WorkOrderArrivalTime
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
	entity.InsuranceExpiredDate = request.WorkOrderInsuranceExpiredDate
	entity.InsuranceClaimNumber = request.WorkOrderInsuranceClaimNo
	entity.InsurancePersonInCharge = request.WorkOrderInsurancePic
	entity.InsuranceOwnRisk = request.WorkOrderInsuranceOwnRisk
	entity.InsuranceWorkOrderNumber = request.WorkOrderInsuranceWONumber

	//page2
	entity.EstTime = request.EstimationDuration
	entity.CustomerExpress = request.CustomerExpress
	entity.LeaveCar = request.LeaveCar
	entity.CarWash = request.CarWash
	entity.PromiseDate = request.PromiseDate
	entity.PromiseTime = request.PromiseTime
	entity.FSCouponNo = request.FSCouponNo
	entity.Notes = request.Notes
	entity.Suggestion = request.Suggestion
	entity.DPAmount = request.DownpaymentAmount

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) Void(tx *gorm.DB, workOrderId int) (bool, *exceptions.BaseErrorResponse) {
	// Check if the work order exists
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", workOrderId).
		First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Work order not found"}
		}
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database", Err: err}
	}

	// Delete the work order
	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete the work order", Err: err}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) CloseOrder(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	// uspg_wtWorkOrder0_Update
	// IF @Option = 2

	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Work order not found"}
		}
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database", Err: err}
	}

	// Check if WorkOrderStatusId is equal to 1 (Draft)
	if entity.WorkOrderStatusId == utils.WoStatDraft {
		return false, &exceptions.BaseErrorResponse{Message: "Work order cannot be closed because status is draft"}
	}

	// Check if there is still DP payment that has not been settled
	var dpPaymentAllocated float64
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).Select("COALESCE(downpayment_payment_allocated, 0) as downpayment_payment_allocated").Scan(&dpPaymentAllocated).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve DP payment allocated from the database", Err: err}
	}
	if dpPaymentAllocated > 0 {
		return false, &exceptions.BaseErrorResponse{Message: "There is still DP payment that has not been settled"}
	}

	// Check if there are any work order items without invoices
	var count int64 //cek statusid <> 8(closed), billcode <> no_charge (5), substituteid
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id <> ? AND substitute_id <> ?",
			Id, 8, 5, 0).
		Count(&count).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order items from the database", Err: err}
	}
	if count > 0 {
		return false, &exceptions.BaseErrorResponse{Message: "Detail Work Order without Invoice No must be deleted"}
	}

	// Check for warranty items
	var allPtpSupply bool //cek statusid <> 8(closed), billcode <> warranty (6), substituteid
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id = ? AND substitute_id <> ?",
			Id, 8, 6, 0).
		Count(&count).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve warranty items from the database", Err: err}
	}
	if count == 0 {
		allPtpSupply = true
	} else {
		// Validate part-to-part supply //cek statusid <> 8(closed), billcode <> warranty (6), substituteid , warrantyclaim_type = 0 (part), frt_qty > supply_qty
		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id = ? AND substitute_id <> ? AND warranty_claim_type_id = ? AND frt_qty > supply_qty",
				Id, 8, 6, 0, 0).
			Count(&count).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{Message: "Failed to validate part-to-part supply", Err: err}
		}
		if count > 0 {
			return false, &exceptions.BaseErrorResponse{Message: "Warranty Item (PTP) must be supplied"}
		}

		// Validate part-to-money and operation status //cek statusid <> 8(closed), billcode <> warranty (6), substituteid , warrantyclaim_type = 0 (part)
		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id = ? AND substitute_id <> ? AND warranty_claim_type_id <> ?",
				Id, 8, 6, 0, 0).
			Count(&count).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{Message: "Failed to validate part-to-money and operation status", Err: err}
		}
		if count > 0 {
			return false, &exceptions.BaseErrorResponse{Message: "Warranty Item (PTM)/Operation must be Invoiced"}
		}

		allPtpSupply = true
	}

	// Check if all items/operations/packages other than warranty are closed
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_status_id <> ? AND substitute_id <> ? AND transaction_type_id NOT IN (?, ?)",
			Id, 8, 0, 6, 5).
		Count(&count).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to check if all items/operations/packages are closed", Err: err}
	}
	if allPtpSupply && count > 0 {
		return false, &exceptions.BaseErrorResponse{Message: "There is Work Order detail that has not been Invoiced"}
	}

	// Validate mileage and update vehicle master if necessary
	var servMileage, lastKm int
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).Select("service_mileage").Scan(&servMileage).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve service mileage", Err: err}
	}
	// err = tx.Table("umVehicle0").Where("vehicle_chassis_number = ?", entity.VehicleChassisNumber).Select("last_km").Scan(&lastKm).Error
	// if err != nil {
	// 	return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve last mileage", Err: err}
	// }
	if servMileage <= lastKm {
		return false, &exceptions.BaseErrorResponse{Message: "Service Mileage must be larger than Last Mileage."}
	}

	// Update vehicle master
	// err = tx.Table("umVehicle0").Where("vehicle_chassis_number = ?", entity.VehicleChassisNumber).
	// 	Updates(map[string]interface{}{
	// 		"last_km":         servMileage,
	// 		"last_serv_date":  entity.WorkOrderDate,
	// 	}).Error
	// if err != nil {
	// 	return false, &exceptions.BaseErrorResponse{Message: "Failed to update vehicle master", Err: err}
	// }

	// If Work Order still has DP Payment not allocated for Invoice
	type DPPaymentDetails struct {
		DPPayment    float64 `gorm:"column:downpayment_payment"`
		DPAllocToInv float64 `gorm:"column:downpayment_payment_to_invoice"`
	}

	var details DPPaymentDetails
	var dpOverpay float64

	err = tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).
		Select("downpayment_payment, downpayment_payment_to_invoice").Scan(&details).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve DP payment details", Err: err}
	}

	if details.DPPayment-details.DPAllocToInv > 0 {
		dpOverpay = details.DPPayment - details.DPAllocToInv
	}

	// Generate DP Other
	// Call dbo.uspg_ctDPIn_Insert and generate journal here
	// TODO: Implement logic for uspg_ctDPIn_Insert and journal generation

	err = tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).
		Updates(map[string]interface{}{
			"downpayment_payment_allocated": details.DPPayment,
			"downpayment_overpay":           dpOverpay,
		}).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to update DP payment details", Err: err}
	}

	// // Determine customer type and set event number
	// var custType string
	// err = tx.Table("gmCust0").Select("customer_type").
	// 	Joins("LEFT JOIN wtWorkOrder0 ON gmCust0.customer_code = wtWorkOrder0.bill_cust_code").
	// 	Where("wtWorkOrder0.work_order_system_number = ?", Id).Scan(&custType).Error
	// if err != nil {
	// 	return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve customer type", Err: err}
	// }
	// var eventNo string
	// switch custType {
	// case "dealer", "imsi":
	// 	eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO_D"
	// case "atpm", "salim", "maintained":
	// 	eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO_A"
	// default:
	// 	eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO"
	// }
	// if eventNo == "" {
	// 	return false, &exceptions.BaseErrorResponse{Message: "Event for Returning DP Customer to DP Other is not exists"}
	// }

	// Generate Journal (DP Customer -> DP Other)
	// Call usp_comJournalAction here
	// TODO: Implement logic for usp_comJournalAction

	// Update JOURNAL_SYS_NO on DPOT
	// TODO: Implement logic for updating JOURNAL_SYS_NO on DPOT
	//}

	// Update the work order status to 8 (Closed)
	entity.WorkOrderStatusId = utils.WoStatClosed
	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to close the work order"}
	}

	return true, nil
}

// uspg_wtWorkOrder1_Insert
// IF @Option = 0
// --USE FOR : * INSERT NEW DATA
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

	for _, entity := range entities {
		workOrderServiceData := make(map[string]interface{})

		workOrderServiceData["work_order_service_id"] = entity.WorkOrderServiceId
		workOrderServiceData["work_order_system_number"] = entity.WorkOrderSystemNumber
		workOrderServiceData["work_order_service_remark"] = entity.WorkOrderServiceRemark
		workOrderServiceResponses = append(workOrderServiceResponses, workOrderServiceData)
	}

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

	payload := transactionworkshoppayloads.WorkOrderServiceRequest{
		WorkOrderServiceId:     entity.WorkOrderServiceId,
		WorkOrderSystemNumber:  entity.WorkOrderSystemNumber,
		WorkOrderServiceRemark: entity.WorkOrderServiceRemark,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) UpdateRequest(tx *gorm.DB, id int, IdWorkorder int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderRequestDescription, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		First(&entity).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderRequestDescription{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "failed to retrieve work order service request from the database",
			Err:        errors.New("failed to retrieve work order service request from the database"),
		}
	}

	entity.WorkOrderServiceRemark = request.WorkOrderServiceRemark

	err = tx.Model(&entity).Updates(map[string]interface{}{
		"work_order_service_remark": entity.WorkOrderServiceRemark,
	}).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderRequestDescription{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to update the work order service request",
			Err:        errors.New("failed to update the work order service request"),
		}
	}

	workOrderRequestDescription := transactionworkshopentities.WorkOrderRequestDescription{}

	return workOrderRequestDescription, nil
}

func (r *WorkOrderRepositoryImpl) AddRequest(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderRequestDescription, *exceptions.BaseErrorResponse) {

	entities := transactionworkshopentities.WorkOrderRequestDescription{
		WorkOrderSystemNumber:   request.WorkOrderSystemNumber,
		WorkOrderServiceRequest: request.WorkOrderServiceRemark,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderRequestDescription{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return entities, nil
}

func (r *WorkOrderRepositoryImpl) DeleteRequest(tx *gorm.DB, id int, IdWorkorder int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete work order service request from the database"}
	}

	return true, nil
}

// uspg_wtWorkOrder1_Insert
// IF @Option = 0
// --USE FOR : * INSERT NEW DATA
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (r *WorkOrderRepositoryImpl) GetAllVehicleService(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrderServiceVehicle

	query := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order service vehicle requests from the database", Err: err}
	}

	if len(entities) == 0 {

		return []map[string]interface{}{}, 0, 0, nil
	}

	var workOrderServiceVehicleResponses []map[string]interface{}

	for _, entity := range entities {
		workOrderServiceVehicleData := make(map[string]interface{})
		workOrderServiceVehicleData["work_order_service_vehicle_id"] = entity.WorkOrderServiceVehicleId
		workOrderServiceVehicleData["work_order_system_number"] = entity.WorkOrderSystemNumber
		workOrderServiceVehicleData["work_order_vehicle_date"] = entity.WorkOrderVehicleDate
		workOrderServiceVehicleData["work_order_vehicle_remark"] = entity.WorkOrderVehicleRemark
		workOrderServiceVehicleResponses = append(workOrderServiceVehicleResponses, workOrderServiceVehicleData)
	}

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

	payload := transactionworkshoppayloads.WorkOrderServiceVehicleRequest{
		WorkOrderSystemNumber:  entity.WorkOrderSystemNumber,
		WorkOrderVehicleDate:   entity.WorkOrderVehicleDate,
		WorkOrderVehicleRemark: entity.WorkOrderVehicleRemark,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) UpdateVehicleService(tx *gorm.DB, id int, IdWorkorder int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderServiceVehicle
	err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		First(&entity).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderServiceVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "failed to retrieve work order service vehicle from the database",
			Err:        errors.New("failed to retrieve work order service vehicle from the database"),
		}
	}

	// Update the fields
	entity.WorkOrderVehicleDate = request.WorkOrderVehicleDate
	entity.WorkOrderVehicleRemark = request.WorkOrderVehicleRemark

	err = tx.Save(&entity).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderServiceVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to save the updated work order service vehicle",
			Err:        errors.New("failed to save the updated work order service vehicle"),
		}
	}

	return entity, nil
}

func (r *WorkOrderRepositoryImpl) AddVehicleService(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse) {

	CurrentDate := time.Now()

	entities := transactionworkshopentities.WorkOrderServiceVehicle{

		WorkOrderSystemNumber:  request.WorkOrderSystemNumber,
		WorkOrderVehicleDate:   CurrentDate,
		WorkOrderVehicleRemark: request.WorkOrderVehicleRemark,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderServiceVehicle{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return entities, nil
}

func (r *WorkOrderRepositoryImpl) DeleteVehicleService(tx *gorm.DB, id int, IdWorkorder int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderServiceVehicle
	err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", id, IdWorkorder).
		Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete work order service request from the database"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) GenerateDocumentNumber(tx *gorm.DB, workOrderId int) (string, *exceptions.BaseErrorResponse) {
	var workOrder transactionworkshopentities.WorkOrder
	var brandResponse transactionworkshoppayloads.BrandDocResponse

	// Get the work order based on the work order system number
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&workOrder).Error
	if err != nil {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve work order from the database: %v", err)}
	}

	if workOrder.BrandId == 0 {

		return "", &exceptions.BaseErrorResponse{Message: "brand_id is missing in the work order. Please ensure the work order has a valid brand_id before generating document number."}
	}

	// Get the last work order based on the work order system number
	var lastWorkOrder transactionworkshopentities.WorkOrder
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("brand_id = ?", workOrder.BrandId).
		Order("work_order_document_number desc").
		First(&lastWorkOrder).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		return "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve last work order: %v", err)}
	}

	currentTime := time.Now()
	month := int(currentTime.Month())
	year := currentTime.Year() % 100 // Use last two digits of the year

	// fetch data brand from external api
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(workOrder.BrandId)
	errUrl := utils.Get(brandUrl, &brandResponse, nil)
	if errUrl != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrl,
		}
	}

	// Check if BrandCode is not empty before using it
	if brandResponse.BrandCode == "" {
		return "", &exceptions.BaseErrorResponse{StatusCode: http.StatusInternalServerError, Message: "Brand code is empty"}
	}

	// Get the initial of the brand code
	brandInitial := brandResponse.BrandCode[0]

	// Handle the case when there is no last work order or the format is invalid
	newDocumentNumber := fmt.Sprintf("WSWO/%c/%02d/%02d/00001", brandInitial, month, year)
	if lastWorkOrder.WorkOrderSystemNumber != 0 {
		lastWorkOrderDate := lastWorkOrder.WorkOrderDate
		lastWorkOrderYear := lastWorkOrderDate.Year() % 100

		// Check if the last work order is from the same year
		if lastWorkOrderYear == year {
			lastWorkOrderCode := lastWorkOrder.WorkOrderDocumentNumber
			codeParts := strings.Split(lastWorkOrderCode, "/")
			if len(codeParts) == 5 {
				lastWorkOrderNumber, err := strconv.Atoi(codeParts[4])
				if err == nil {
					newWorkOrderNumber := lastWorkOrderNumber + 1
					newDocumentNumber = fmt.Sprintf("WSWO/%c/%02d/%02d/%05d", brandInitial, month, year, newWorkOrderNumber)
				} else {
					log.Printf("Failed to parse last work order code: %v", err)
				}
			} else {
				log.Println("Invalid last work order code format")
			}
		}
	}

	log.Printf("New document number: %s", newDocumentNumber)
	return newDocumentNumber, nil
}

func (r *WorkOrderRepositoryImpl) Submit(tx *gorm.DB, workOrderId int) (bool, string, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "", &exceptions.BaseErrorResponse{Message: "No work order data found"}
		}
		return false, "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to retrieve work order from the database: %v", err)}
	}

	if entity.WorkOrderDocumentNumber == "" && entity.WorkOrderStatusId == utils.WoStatDraft {
		//Generate new document number
		newDocumentNumber, genErr := r.GenerateDocumentNumber(tx, entity.WorkOrderSystemNumber)
		if genErr != nil {
			return false, "", genErr
		}
		//newDocumentNumber := "WSWO/1/21/21/00001"

		entity.WorkOrderDocumentNumber = newDocumentNumber

		// Update work order status to 2 (New Submitted)
		entity.WorkOrderStatusId = utils.WoStatNew

		err = tx.Save(&entity).Error
		if err != nil {
			return false, "", &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to submit the work order: %v", err)}
		}

		return true, newDocumentNumber, nil
	} else {

		return false, "", &exceptions.BaseErrorResponse{Message: "Document number has already been generated"}
	}
}

// uspg_wtWorkOrder2_Insert
// IF @Option = 0
// --USE FOR : * INSERT NEW DATA
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (r *WorkOrderRepositoryImpl) GetAllDetailWorkOrder(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tableStruct := transactionworkshoppayloads.WorkOrderDetailRequest{}

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

	var convertedResponses []transactionworkshoppayloads.WorkOrderDetailResponse

	for rows.Next() {

		var (
			workOrderReq transactionworkshoppayloads.WorkOrderDetailRequest
			workOrderRes transactionworkshoppayloads.WorkOrderDetailResponse
		)

		if err := rows.Scan(
			&workOrderReq.WorkOrderDetailId,
			&workOrderReq.WorkOrderSystemNumber,
			&workOrderReq.LineTypeId,
			&workOrderReq.TransactionTypeId,
			&workOrderReq.JobTypeId,
			&workOrderReq.FrtQuantity,
			&workOrderReq.SupplyQuantity,
			&workOrderReq.PriceListId,
			&workOrderReq.WarehouseId,
			&workOrderReq.ItemId,
			&workOrderReq.OperationId,
			&workOrderReq.OperationItemCode,
			&workOrderReq.ProposedPrice,
			&workOrderReq.OperationItemPrice,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		workOrderRes = transactionworkshoppayloads.WorkOrderDetailResponse{
			WorkOrderDetailId:     workOrderReq.WorkOrderDetailId,
			WorkOrderSystemNumber: workOrderReq.WorkOrderSystemNumber,
			LineTypeId:            workOrderReq.LineTypeId,
			TransactionTypeId:     workOrderReq.TransactionTypeId,
			JobTypeId:             workOrderReq.JobTypeId,
			FrtQuantity:           workOrderReq.FrtQuantity,
			SupplyQuantity:        workOrderReq.SupplyQuantity,
		}

		convertedResponses = append(convertedResponses, workOrderRes)
	}

	var mapResponses []map[string]interface{}

	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_detail_id":     response.WorkOrderDetailId,
			"work_order_system_number": response.WorkOrderSystemNumber,
			"line_type_id":             response.LineTypeId,
			"transaction_type_id":      response.TransactionTypeId,
			"job_type_id":              response.JobTypeId,
			"frt_quantity":             response.FrtQuantity,
			"supply_quantity":          response.SupplyQuantity,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil

}

func (r *WorkOrderRepositoryImpl) GetDetailByIdWorkOrder(tx *gorm.DB, id int, IdWorkorder int) (transactionworkshoppayloads.WorkOrderDetailRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderDetail
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_detail_id = ?", id, IdWorkorder).
		First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.WorkOrderDetailRequest{}, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Message: "Work order detail not found"}
		}
		return transactionworkshoppayloads.WorkOrderDetailRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order detail from the database", Err: err}
	}

	payload := transactionworkshoppayloads.WorkOrderDetailRequest{
		WorkOrderDetailId:     entity.WorkOrderDetailId,
		WorkOrderSystemNumber: entity.WorkOrderSystemNumber,
		LineTypeId:            entity.LineTypeId,
		TransactionTypeId:     entity.TransactionTypeId,
		JobTypeId:             entity.JobTypeId,
		FrtQuantity:           entity.FrtQuantity,
		SupplyQuantity:        entity.SupplyQuantity,
		PriceListId:           entity.PriceListId,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) CalculateWorkOrderTotal(tx *gorm.DB, workOrderSystemNumber int, lineTypeId int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {

	type Result struct {
		TotalPackage            float64
		TotalOperation          float64
		TotalSparePart          float64
		TotalOil                float64
		TotalMaterial           float64
		TotalFee                float64
		TotalAccessories        float64
		TotalConsumableMaterial float64
		TotalSublet             float64
		TotalSouvenir           float64
	}

	var result Result

	// Aggregate data using GORM
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select(`
			SUM(CASE WHEN line_type_id = 0 THEN ROUND(COALESCE(operation_item_price, 0), 0) ELSE 0 END) AS total_package,
			SUM(CASE WHEN line_type_id = 1 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_operation,
			SUM(CASE WHEN line_type_id = 2 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_spare_part,
			SUM(CASE WHEN line_type_id = 3 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_oil,
			SUM(CASE WHEN line_type_id = 4 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_material,
			SUM(CASE WHEN line_type_id = 5 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_fee,
			SUM(CASE WHEN line_type_id = 6 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_accessories,
			SUM(CASE WHEN line_type_id = 7 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_consumable_material,
			SUM(CASE WHEN line_type_id = 8 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_sublet,
			SUM(CASE WHEN line_type_id = 9 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_souvenir
		`).
		Where("work_order_system_number = ?", workOrderSystemNumber).
		Where("line_type_id = ?", lineTypeId).
		Scan(&result).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to calculate totals: %v", err)}
	}

	// Calculate grand total
	grandTotal := result.TotalPackage + result.TotalOperation + result.TotalSparePart + result.TotalOil + result.TotalMaterial + result.TotalFee + result.TotalAccessories + result.TotalConsumableMaterial + result.TotalSublet + result.TotalSouvenir

	// Update Work Order with the calculated totals
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", workOrderSystemNumber).
		Updates(map[string]interface{}{
			"total_package":             result.TotalPackage,
			"total_operation":           result.TotalOperation,
			"total_part":                result.TotalSparePart,
			"total_oil":                 result.TotalOil,
			"total_material":            result.TotalMaterial,
			"total_price_accessories":   result.TotalAccessories,
			"total_consumable_material": result.TotalConsumableMaterial,
			"total_sublet":              result.TotalSublet,
			"total":                     grandTotal,
			//"total_fee":                 result.TotalFee,
			//"total_souvenir":            result.TotalSouvenir,
		}).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to update work order: %v", err)}
	}

	// Prepare response
	workOrderDetailResponses := []map[string]interface{}{
		{"total_package": result.TotalPackage},
		{"total_operation": result.TotalOperation},
		{"total_part": result.TotalSparePart},
		{"total_oil": result.TotalOil},
		{"total_material": result.TotalMaterial},
		{"total_price_accessories": result.TotalAccessories},
		{"total_consumable_material": result.TotalConsumableMaterial},
		{"total_sublet": result.TotalSublet},
		{"total": grandTotal},
	}

	return workOrderDetailResponses, nil
}

// uspg_wtWorkOrder2_Insert
// IF @Option = 0
// --USE FOR : * INSERT NEW DATA DETAIL
func (r *WorkOrderRepositoryImpl) AddDetailWorkOrder(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	workOrderDetail := transactionworkshopentities.WorkOrderDetail{}
	currentDate := time.Now()

	var workOrderTypeId int
	if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Select("work_order_type_id").
		Where("work_order_system_number = ?", id).
		Scan(&workOrderTypeId).Error; err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order type id",
			Err:        err,
		}
	}

	// Validate if the work order is still draft

	// Validasi untuk chassis yang sudah pernah PDI,FSI,WR

	// Validate Line Type Item must be inside item master

	// Validate if Warranty to Vehicle Age

	// LINE TYPE <> 1 , NEED SUBSTITUTE

	///////////////////////////////////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Insert detil work order (WO2) berdasarkan tipe work order (Normal, Campaign, Affiliated, Repeat Job):

	switch workOrderTypeId {
	case 1: // Normal Work Order
		var estimSystemNo int
		if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
			Where("work_order_system_number = ?", id).
			Select("estimation_system_number").
			Scan(&estimSystemNo).Error; err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve estimation system number",
				Err:        err,
			}
		}

		if estimSystemNo == 0 {
			var maxWoOprItemLine int
			if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("ISNULL(MAX(work_order_operation_item_line), 0)").
				Where("work_order_system_number = ?", id).
				Scan(&maxWoOprItemLine).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve maximum work order operation item line",
					Err:        err,
				}
			}

			// var bookingEstim21 []transactionworkshopentities.BookingEstimation
			// tx.Model(&bookingEstim21).
			// 	Select("BE.ESTIM_LINE, BE.LINE_TYPE, BE.OPR_ITEM_CODE, BE.DESCRIPTION, I.SELLING_UOM, BE.FRT_QTY, BE.OPR_ITEM_PRICE, BE.OPR_ITEM_DISC_AMOUNT, BE.OPR_ITEM_DISC_REQ_AMOUNT, BE.OPR_ITEM_DISC_PERCENT, BE.OPR_ITEM_DISC_REQ_PERCENT, BE.PPH_AMOUNT, BE.PPH_TAX_CODE, BE.PPH_TAX_RATE").
			// 	Joins("LEFT OUTER JOIN wtBookEstim0 BE0 ON BE0.ESTIM_SYSTEM_NO = BE.ESTIM_SYSTEM_NO").
			// 	Joins("LEFT OUTER JOIN gmItem0 I ON I.ITEM_CODE = BE.OPR_ITEM_CODE").
			// 	Where("BE.ESTIM_SYSTEM_NO = ?", estimSystemNo).
			// 	Find(&bookingEstim21)

			workOrderDetail.WorkOrderOperationItemLine = maxWoOprItemLine + 1

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:              id,
				LineTypeId:                         request.LineTypeId,
				TransactionTypeId:                  request.TransactionTypeId,
				JobTypeId:                          request.JobTypeId,
				WarehouseId:                        request.WarehouseId,
				FrtQuantity:                        request.FrtQuantity,
				SupplyQuantity:                     request.SupplyQuantity,
				WorkorderStatusId:                  utils.WoStatDraft,
				PriceListId:                        request.PriceListId,
				OperationItemDiscountRequestAmount: request.ProposedPrice,
				OperationItemPrice:                 request.OperationItemPrice,
			}

			if request.LineTypeId == 1 {
				workOrderDetail.OperationId = request.OperationId
				workOrderDetail.ItemId = 0
			} else {
				workOrderDetail.ItemId = request.ItemId
				workOrderDetail.OperationId = 0
			}

			if err := tx.Create(&workOrderDetail).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to create work order detail",
					Err:        err,
				}
			}

			if _, err := r.CalculateWorkOrderTotal(tx, id, request.LineTypeId); err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, err
			}
		}
	case 2: // Campaign Work Order
		var campaignId int
		if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
			Where("work_order_system_number = ?", id).
			Select("campaign_id").
			Scan(&campaignId).Error; err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve campaign id",
				Err:        err,
			}
		}

		if campaignId > 0 {
			var campaignMaster masterentities.CampaignMaster

			if err := tx.Model(&masterentities.CampaignMaster{}).
				Where("campaign_id = ? AND ? BETWEEN campaign_period_from AND campaign_period_to", campaignId, currentDate).
				First(&campaignMaster).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {

					return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "Campaign Code is not valid",
						Err:        errors.New("campaign Code is not valid"),
					}
				}

				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check if campaign exists",
					Err:        err,
				}
			}

			var maxWoOprItemLine int
			if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("ISNULL(MAX(work_order_operation_item_line), 0)").
				Where("work_order_system_number = ?", id).
				Scan(&maxWoOprItemLine).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve maximum work order operation item line",
					Err:        err,
				}
			}

			// var campaignItems []mastercampaignmasterentities.CampaignMaster
			// tx.Model(&campaignItems).Select("C1.LINE_TYPE, C1.OPR_ITEM_CODE, C1.DESCRIPTION, I.SELLING_UOM, C1.FRT_QTY, C1.OPR_ITEM_PRICE, C1.OPR_ITEM_DISC_PERCENT").
			// 	Joins("INNER JOIN amCampaign0 C ON C1.CAMPAIGN_CODE = C.CAMPAIGN_CODE").
			// 	Joins("LEFT JOIN gmItem0 I ON I.ITEM_CODE = C1.OPR_ITEM_CODE").
			// 	Where("C1.campaign_id = ?", campaignId).
			// 	Find(&campaignItems)

			workOrderDetail.WorkOrderOperationItemLine = maxWoOprItemLine + 1

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:              id,
				LineTypeId:                         request.LineTypeId,
				TransactionTypeId:                  request.TransactionTypeId,
				JobTypeId:                          request.JobTypeId,
				WarehouseId:                        request.WarehouseId,
				ItemId:                             request.ItemId,
				FrtQuantity:                        request.FrtQuantity,
				SupplyQuantity:                     request.SupplyQuantity,
				WorkorderStatusId:                  utils.WoStatDraft,
				PriceListId:                        request.PriceListId,
				OperationItemDiscountRequestAmount: request.ProposedPrice,
				OperationItemPrice:                 request.OperationItemPrice,
			}

			if request.LineTypeId == 1 {
				workOrderDetail.OperationId = request.OperationId
				workOrderDetail.ItemId = 0
			} else {
				workOrderDetail.ItemId = request.ItemId
				workOrderDetail.OperationId = 0
			}

			if err := tx.Create(&workOrderDetail).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to create work order detail",
					Err:        err,
				}
			}

			if _, err := r.CalculateWorkOrderTotal(tx, id, request.LineTypeId); err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, err
			}
		}
	case 3: // Affiliated Work Order
		var results struct {
			PDISystemNo  int `gorm:"column:pdi_system_number"`
			PDILineNo    int `gorm:"column:pdi_line_number"`
			ServReqSysNo int `gorm:"column:service_request_system_number"`
		}

		if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
			Select("ISNULL(pdi_system_number, 0) AS pdi_system_number, ISNULL(pdi_line_number, 0) AS pdi_line_number, ISNULL(service_request_system_number, 0) AS service_request_system_number").
			Where("work_order_system_number = ?", id).
			Scan(&results).Error; err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve PDI system number, PDI line number, and service request system number",
				Err:        err,
			}
		}

		pdiSystemNo := results.PDISystemNo
		servReqSysNo := results.ServReqSysNo

		fmt.Printf("PDI System No: %d, Service Request Sys No: %d\n", pdiSystemNo, servReqSysNo)

		if pdiSystemNo != 0 && servReqSysNo == 0 {
			var maxWoOprItemLine int
			if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("ISNULL(MAX(work_order_operation_item_line), 0)").
				Where("work_order_system_number = ?", id).
				Scan(&maxWoOprItemLine).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve maximum work order operation item line",
					Err:        err,
				}
			}

			// Insert into wtWorkOrder2
			// 	tx.Exec(`
			// 	INSERT INTO wtWorkOrder2 (
			// 		RECORD_STATUS, WO_SYS_NO, WO_DOC_NO, WO_OPR_ITEM_LINE, WO_LINE_STAT, LINE_TYPE, WO_OPR_STATUS, BILL_CODE, JOB_TYPE,
			// 		WO_LINE_DISC_STAT, OPR_ITEM_CODE, DESCRIPTION, ITEM_UOM, FRT_QTY, OPR_ITEM_PRICE, OPR_ITEM_DISC_AMOUNT, OPR_ITEM_DISC_REQ_AMOUNT,
			// 		OPR_ITEM_DISC_PERCENT, OPR_ITEM_DISC_REQ_PERCENT, PPH_AMOUNT, PPH_TAX_CODE, PPH_TAX_RATE, SUPPLY_QTY, WHS_GROUP, CHANGE_NO,
			// 		CREATION_USER_ID, CREATION_DATETIME, CHANGE_USER_ID, CHANGE_DATETIME
			// 	)
			// 	SELECT ?, ?, ?, ? + ROW_NUMBER() OVER (ORDER BY P1.PDI_LINE), ?, ?, '',
			// 		dbo.FCT_getBillCode(?, CAST(P1.COMPANY_CODE AS VARCHAR(10)), 'W'),
			// 		dbo.getVariableValue('JOBTYPE_PDI'), ?, P1.OPERATION_NO, OPR.OPERATION_NAME, '', P1.FRT, LSP1.SELLING_PRICE, 0, 0, 0, 0, 0,
			// 		0, OPR.TAX_CODE, 0, P1.FRT, ?, 0, ?, ?, ?, ?
			// 	FROM dbo.atPDI1 P1
			// 	INNER JOIN dbo.atPDI0 P ON P1.PDI_SYS_NO = P.PDI_SYS_NO
			// 	LEFT OUTER JOIN dbo.amOperationCode OPR ON OPR.OPERATION_CODE = P1.OPERATION_NO
			// 	LEFT OUTER JOIN dbo.amOperation0 OP ON OP.OPERATION_CODE = P1.OPERATION_NO AND OP.VEHICLE_BRAND = P.VEHICLE_BRAND AND OP.MODEL_CODE = P1.MODEL_CODE
			// 	LEFT OUTER JOIN dbo.amLabSellPrice1 LSP1 ON LSP1.EFFECTIVE_DATE = (
			// 		SELECT TOP 1 EFFECTIVE_DATE
			// 		FROM dbo.amLabSellPrice1
			// 		WHERE EFFECTIVE_DATE <= CONVERT(VARCHAR, GETDATE(), 106) AND MODEL_CODE = P1.MODEL_CODE AND VEHICLE_BRAND = P.VEHICLE_BRAND
			// 		AND JOB_TYPE = dbo.getVariableValue('JOBTYPE_PDI') AND RECORD_STATUS = ? AND COMPANY_CODE = P1.COMPANY_CODE
			// 		ORDER BY EFFECTIVE_DATE DESC
			// 	) AND LSP1.MODEL_CODE = P1.MODEL_CODE AND LSP1.VEHICLE_BRAND = P.VEHICLE_BRAND AND LSP1.JOB_TYPE = dbo.getVariableValue('JOBTYPE_PDI')
			// 	AND LSP1.RECORD_STATUS = ? AND LSP1.COMPANY_CODE = P1.COMPANY_CODE
			// 	WHERE P.PDI_SYS_NO = ? AND P1.PDI_LINE = ?
			// `,
			// 		request.RecordStatus, request.WorkOrderSystemNo, request.WorkOrderDocumentNo, maxWoOprItemLine, request.WorkOrderStatusNew,
			// 		request.LineTypeOpr, request.CompanyCode, request.ApprovalDraft, request.WhsGroupSp, request.CreationUserId, time.Now(),
			// 		request.ChangeUserId, time.Now(), request.RecordStatus, request.RecordStatus, pdiSystemNo, pdiLineNo,
			// 	)

			// Update atPDI1
			// tx.Model(&AtPDI1{}).
			// 	Where("PDI_SYS_NO = ? AND PDI_LINE = ?", pdiSystemNo, pdiLineNo).
			// 	Updates(map[string]interface{}{
			// 		"WO_SYS_NO":   request.WorkOrderSystemNo,
			// 		"WO_DOC_NO":   request.WorkOrderDocumentNo,
			// 		"WO_DATE":     request.WorkOrderDate,
			// 		"LINE_STATUS": request.PDIStatusWO,
			// 	})

			workOrderDetail.WorkOrderOperationItemLine = maxWoOprItemLine + 1

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:              id,
				LineTypeId:                         request.LineTypeId,
				TransactionTypeId:                  request.TransactionTypeId,
				JobTypeId:                          request.JobTypeId,
				WarehouseId:                        request.WarehouseId,
				ItemId:                             request.ItemId,
				FrtQuantity:                        request.FrtQuantity,
				SupplyQuantity:                     request.SupplyQuantity,
				WorkorderStatusId:                  utils.WoStatDraft,
				PriceListId:                        request.PriceListId,
				OperationItemDiscountRequestAmount: request.ProposedPrice,
				OperationItemPrice:                 request.OperationItemPrice,
			}

			if request.LineTypeId == 1 {
				workOrderDetail.OperationId = request.OperationId
				workOrderDetail.ItemId = 0
			} else {
				workOrderDetail.ItemId = request.ItemId
				workOrderDetail.OperationId = 0
			}

			if err := tx.Create(&workOrderDetail).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to create work order detail",
					Err:        err,
				}
			}

			if _, err := r.CalculateWorkOrderTotal(tx, id, request.LineTypeId); err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, err
			}

		} else if pdiSystemNo == 0 && servReqSysNo != 0 {

			fmt.Printf("2. PDI System No: %d, Service Request Sys No: %d\n", pdiSystemNo, servReqSysNo)

			var maxWoOprItemLine int
			if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("ISNULL(MAX(work_order_operation_item_line), 0)").
				Where("work_order_system_number = ?", id).
				Scan(&maxWoOprItemLine).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve maximum work order operation item line",
					Err:        err,
				}
			}

			// Insert into wtWorkOrder2
			// 	tx.Exec(`
			// 	INSERT INTO wtWorkOrder2 (
			// 		RECORD_STATUS, WO_SYS_NO, WO_DOC_NO, WO_OPR_ITEM_LINE, WO_LINE_STAT, LINE_TYPE, WO_OPR_STATUS, BILL_CODE, JOB_TYPE,
			// 		WO_LINE_DISC_STAT, OPR_ITEM_CODE, DESCRIPTION, ITEM_UOM, FRT_QTY, OPR_ITEM_PRICE, OPR_ITEM_DISC_AMOUNT, OPR_ITEM_DISC_REQ_AMOUNT,
			// 		OPR_ITEM_DISC_PERCENT, OPR_ITEM_DISC_REQ_PERCENT, PPH_AMOUNT, PPH_TAX_CODE, PPH_TAX_RATE, SUPPLY_QTY, WHS_GROUP, CHANGE_NO,
			// 		CREATION_USER_ID, CREATION_DATETIME, CHANGE_USER_ID, CHANGE_DATETIME
			// 	)
			// 	SELECT ?, ?, ?, ? + ROW_NUMBER() OVER (ORDER BY SR1.SERV_REQ_LINE_NO), ?, SR1.LINE_TYPE, '',
			// 		dbo.FCT_getBillCode(?, SR.COMPANY_CODE, 'W'),
			// 		CASE WHEN SR.SERV_PROFIT_CENTER = ? THEN ? ELSE ? END, ?,
			// 		SR1.OPR_ITEM_CODE, CASE
			// 			WHEN SR1.LINE_TYPE = ? THEN OP.OPERATION_NAME
			// 			WHEN SR1.LINE_TYPE = ? THEN PCK.PACKAGE_NAME
			// 			ELSE I.ITEM_NAME
			// 		END, I.SELLING_UOM, SR1.FRT_QTY, dbo.getOprItemPrice(SR1.LINE_TYPE, ?, dbo.FCT_getBillCode(?, SR.COMPANY_CODE, 'W'), ?, ?, ?, ?, ?, ?, ?, 0, DEFAULT, ?), 0, 0, 0, 0, 0,
			// 		CASE
			// 			WHEN SR1.LINE_TYPE = ? THEN OPR.TAX_CODE
			// 			WHEN SR1.LINE_TYPE = ? THEN PCK.PPH_TAX_CODE
			// 			ELSE ''
			// 		END, 0, SR1.FRT_QTY, ?, 0, ?, ?, ?, ?
			// 	FROM dbo.atServiceReq1 SR1
			// 	INNER JOIN dbo.atServiceReq0 SR ON SR1.SERV_REQ_SYS_NO = SR.SERV_REQ_SYS_NO
			// 	LEFT OUTER JOIN dbo.amOperationCode OPR ON OPR.OPERATION_CODE = SR1.OPR_ITEM_CODE
			// 	LEFT OUTER JOIN dbo.amOperation0 OP ON OP.OPERATION_CODE = SR1.OPR_ITEM_CODE AND OP.VEHICLE_BRAND = SR.VEHICLE_BRAND AND OP.MODEL_CODE = SR.MODEL_CODE
			// 	LEFT OUTER JOIN dbo.amPackage0 PCK ON PCK.PACKAGE_CODE = SR1.OPR_ITEM_CODE AND PCK.VEHICLE_BRAND = SR.VEHICLE_BRAND AND PCK.MODEL_CODE = SR.MODEL_CODE
			// 	LEFT OUTER JOIN dbo.amItem0 I ON I.ITEM_CODE = SR1.OPR_ITEM_CODE AND I.VEHICLE_BRAND = SR.VEHICLE_BRAND AND I.MODEL_CODE = SR.MODEL_CODE
			// 	WHERE SR.SERV_REQ_SYS_NO = ?
			// `,
			// 		request.RecordStatus, request.WorkOrderSystemNo, request.WorkOrderDocumentNo, maxWoOprItemLine, request.WorkOrderStatusNew,
			// 		request.LineTypeOpr, request.CompanyCode, request.JobTypePDI, request.ApprovalDraft, request.JobTypePDI, request.JobTypePDI,
			// 		request.LineTypeOpr, request.JobTypePDI, request.CompanyCode, request.JobTypePDI, request.CompanyCode, request.JobTypePDI,
			// 		request.ModelCode, request.LineTypeOpr, request.PDIStatusWO, request.WhsGroupSp, request.CreationUserId, time.Now(),
			// 		request.ChangeUserId, time.Now(), request.RecordStatus, request.RecordStatus, servReqSysNo,
			// 	)

			// Update atServiceReq0
			// tx.Model(&transactionworkshopentities.ServiceRequest{}).
			// 	Where("service_request_system_number = ?", servReqSysNo).
			// 	Updates(map[string]interface{}{
			// 		"work_order_system_number":   request.WorkOrderSystemNo,
			// 		"work_order_document_number": request.WorkOrderDocumentNo,
			// 		"WO_DATE":                    request.WorkOrderDate,
			// 		"LINE_STATUS":                request.PDIStatusWO,
			// })

			workOrderDetail.WorkOrderOperationItemLine = maxWoOprItemLine + 1

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:              id,
				LineTypeId:                         request.LineTypeId,
				TransactionTypeId:                  request.TransactionTypeId,
				JobTypeId:                          request.JobTypeId,
				WarehouseId:                        request.WarehouseId,
				ItemId:                             request.ItemId,
				FrtQuantity:                        request.FrtQuantity,
				SupplyQuantity:                     request.SupplyQuantity,
				WorkorderStatusId:                  utils.WoStatDraft,
				PriceListId:                        request.PriceListId,
				OperationItemDiscountRequestAmount: request.ProposedPrice,
				OperationItemPrice:                 request.OperationItemPrice,
			}

			if request.LineTypeId == 1 {
				workOrderDetail.OperationId = request.OperationId
				workOrderDetail.ItemId = 0
			} else {
				workOrderDetail.ItemId = request.ItemId
				workOrderDetail.OperationId = 0
			}

			if err := tx.Create(&workOrderDetail).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to create work order detail",
					Err:        err,
				}
			}

			if _, err := r.CalculateWorkOrderTotal(tx, id, request.LineTypeId); err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, err
			}

		}
	case 4: // Repeat Job Work Order
		var jobId int
		if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
			Where("work_order_system_number = ?", id).
			Select("repeated_system_number").
			Scan(&jobId).Error; err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve repeated system number",
				Err:        err,
			}
		}

		if jobId != 0 {
			var maxWoOprItemLine int
			if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("ISNULL(MAX(work_order_operation_item_line), 0)").
				Where("work_order_system_number = ?", id).
				Scan(&maxWoOprItemLine).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve maximum work order operation item line",
					Err:        err,
				}
			}

			// Insert into wtWorkOrder2
			// tx.Exec(`
			// 	INSERT INTO wtWorkOrder2 (
			// 		RECORD_STATUS, WO_SYS_NO, WO_DOC_NO, WO_OPR_ITEM_LINE, WO_LINE_STAT, LINE_TYPE, WO_OPR_STATUS, BILL_CODE, JOB_TYPE,
			// 		WO_LINE_DISC_STAT, OPR_ITEM_CODE, DESCRIPTION, ITEM_UOM, FRT_QTY, OPR_ITEM_PRICE, OPR_ITEM_DISC_AMOUNT, OPR_ITEM_DISC_REQ_AMOUNT,
			// 		OPR_ITEM_DISC_PERCENT, OPR_ITEM_DISC_REQ_PERCENT, PPH_AMOUNT, PPH_TAX_CODE, PPH_TAX_RATE, SUPPLY_QTY, WHS_GROUP, CHANGE_NO,
			// 		CREATION_USER_ID, CREATION_DATETIME, CHANGE_USER_ID, CHANGE_DATETIME
			// 	)
			// 	SELECT
			// 		?, ?, ?, ? + ROW_NUMBER() OVER (ORDER BY RW1.WO_OPR_ITEM_LINE), ?, RW1.LINE_TYPE, '',
			// 		dbo.getVariableValue('TRXTYPE_WO_NOCHARGE'), RW1.JOB_TYPE, ?, RW1.OPR_ITEM_CODE, RW1.DESCRIPTION, RW1.ITEM_UOM,
			// 		RW1.FRT_QTY, RW1.OPR_ITEM_PRICE, 0, 0, 0, 0, 0, '', 0, RW1.SUPPLY_QTY, RW1.WHS_GROUP, 0, ?, ?, ?, ?
			// 	FROM dbo.wtWorkOrder0 RW
			// 	INNER JOIN dbo.wtWorkOrder2 RW1 ON RW1.WO_SYS_NO = RW.WO_SYS_NO
			// 	WHERE RW.WO_SYS_NO = ? AND RW1.LINE_TYPE = ?
			// `,
			// 	request.RecordStatus, request.WorkOrderSystemNo, request.WorkOrderDocumentNo, maxWoOprItemLine, request.WorkOrderStatusNew,
			// 	request.ApprovalDraft, request.CreationUserId, time.Now(), request.ChangeUserId, time.Now(),
			// 	repeatedWoSysNo, request.LineTypeOpr,
			// )

			workOrderDetail.WorkOrderOperationItemLine = maxWoOprItemLine + 1

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:              id,
				LineTypeId:                         request.LineTypeId,
				TransactionTypeId:                  request.TransactionTypeId,
				JobTypeId:                          request.JobTypeId,
				WarehouseId:                        request.WarehouseId,
				ItemId:                             request.ItemId,
				FrtQuantity:                        request.FrtQuantity,
				SupplyQuantity:                     request.SupplyQuantity,
				WorkorderStatusId:                  utils.WoStatDraft,
				PriceListId:                        request.PriceListId,
				OperationItemDiscountRequestAmount: request.ProposedPrice,
				OperationItemPrice:                 request.OperationItemPrice,
			}

			if request.LineTypeId == 1 {
				workOrderDetail.OperationId = request.OperationId
				workOrderDetail.ItemId = 0
			} else {
				workOrderDetail.ItemId = request.ItemId
				workOrderDetail.OperationId = 0
			}

			if err := tx.Create(&workOrderDetail).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to create work order detail",
					Err:        err,
				}
			}

			if _, err := r.CalculateWorkOrderTotal(tx, id, request.LineTypeId); err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, err
			}
		}
	default:
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid work order type",
			Err:        errors.New("invalid work order type"),
		}
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////////////////////////////
	// substitusi item jika stok tidak mencukupi
	// tx = tx.Begin()
	// if tx.Error != nil {
	// 	return false, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Err:        tx.Error,
	// 	}
	// }

	// var workOrderItems []transactionworkshopentities.WorkOrderDetail
	// result := tx.Raw(`
	// 	SELECT WO_OPR_ITEM_LINE, OPR_ITEM_CODE, WHS_GROUP, ISNULL(FRT_QTY, 0) AS FRT_QTY
	// 	FROM dbo.wtWorkOrder2
	// 	WHERE WO_SYS_NO = ? AND LINE_TYPE <> ? AND LINE_TYPE <> ?
	// `, request.WorkOrderSystemNo, request.LineTypeOpr, request.LineTypePackage).Scan(&workOrderItems)

	// if result.Error != nil {
	// 	tx.Rollback()
	// 	return false, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Err:        result.Error,
	// 	}
	// }

	// // Iterate over cursor results
	// for _, item := range workOrderItems {
	// 	uomType := getVariableValue("UOM_TYPE_SELL")

	// 	// get from table amLocationStockItem currently not developed
	// 	qtyAvail, err := uspg_amLocationStockItem_Select(
	// 		request.CompanyCode, request.CreationDatetime, item.OPR_ITEM_CODE, item.WHS_GROUP, uomType)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return false, &exceptions.BaseErrorResponse{
	// 			StatusCode: http.StatusInternalServerError,
	// 			Err:        err,
	// 		}
	// 	}

	// 	if qtyAvail == 0 {
	// 		// get from table uspg_smSubstitute currently not developed
	// 		substituteItems, err := uspg_smSubstitute0_Select(
	// 			request.CompanyCode, item.OPR_ITEM_CODE, item.FRT_QTY)
	// 		if err != nil {
	// 			tx.Rollback()
	// 			return false, &exceptions.BaseErrorResponse{
	// 				StatusCode: http.StatusInternalServerError,
	// 				Err:        err,
	// 			}
	// 		}

	// 		for _, subsItem := range substituteItems {
	// 			if subsItem.SUBS_TYPE != "" {
	// 				// Update original item code record before substituted
	// 				if err := updateOriginalItemCodeRecord(tx, request, item, subsItem); err != nil {
	// 					tx.Rollback()
	// 					return false, &exceptions.BaseErrorResponse{
	// 						StatusCode: http.StatusInternalServerError,
	// 						Err:        err,
	// 					}
	// 				}

	// 				// Insert new substituted item
	// 				if err := insertSubstitutedItem(tx, request, item, subsItem); err != nil {
	// 					tx.Rollback()
	// 					return false, &exceptions.BaseErrorResponse{
	// 						StatusCode: http.StatusInternalServerError,
	// 						Err:        err,
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// Menghitung total biaya work order berdasarkan tipe line item dan melakukan update pada tabel

	// Menghitung total diskon dan PPN, serta memperbarui work order

	return workOrderDetail, nil
}

func (r *WorkOrderRepositoryImpl) UpdateDetailWorkOrder(tx *gorm.DB, IdWorkorder int, id int, request transactionworkshoppayloads.WorkOrderDetailRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderDetail
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_detail_id = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order detail from the database"}
	}

	entity.LineTypeId = request.LineTypeId
	entity.TransactionTypeId = request.TransactionTypeId
	entity.JobTypeId = request.JobTypeId
	entity.WarehouseId = request.WarehouseId
	entity.ItemId = request.ItemId
	entity.FrtQuantity = request.FrtQuantity
	entity.SupplyQuantity = request.SupplyQuantity
	entity.PriceListId = request.PriceListId
	entity.OperationItemDiscountRequestAmount = request.ProposedPrice
	entity.OperationItemPrice = request.OperationItemPrice

	if request.LineTypeId == 1 {
		entity.OperationId = request.OperationId
	} else {
		entity.ItemId = request.ItemId
	}

	err = tx.Save(&entity).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order detail"}
	}

	// Call CalculateWorkOrderTotal to update the totals in trx_work_order
	_, calcErr := r.CalculateWorkOrderTotal(tx, id, request.LineTypeId)
	if calcErr != nil {
		return transactionworkshopentities.WorkOrderDetail{}, calcErr
	}

	return entity, nil
}

func (r *WorkOrderRepositoryImpl) DeleteDetailWorkOrder(tx *gorm.DB, id int, IdWorkorder int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderDetail
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_detail_id = ?", id, IdWorkorder).
		Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete work order detail from the database"}
	}

	return true, nil
}

// uspg_wtWorkOrder0_Insert
// IF @Option = 1
// --USE FOR : * INSERT NEW DATA FROM BOOKING AND ESTIMATION
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (r *WorkOrderRepositoryImpl) NewBooking(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderBookingRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse) {

	// Default values
	defaultWorkOrderDocumentNumber := ""
	//defaultWorkOrderTypeId := 1   // 1:Normal, 2:Campaign, 3:Affiliated, 4:Repeat Job
	defaultServiceAdvisorId := 1 // Default advisor ID
	defaultCPCcode := "00002"    // Default CPC code 00002 for workshop
	workOrderTypeId := 1         // Default work order type ID 1 for normal

	// Validate request date
	currentDate := time.Now()
	requestDate := request.WorkOrderArrivalTime.Truncate(24 * time.Hour)
	if requestDate.Before(currentDate) || requestDate.After(currentDate) {
		request.WorkOrderArrivalTime = currentDate
	}

	// Check if the CompanyId is provided
	if request.CompanyId == 0 {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Company ID is required",
			Err:        errors.New("parameter has lost session, please refresh the data"),
		}
	}

	// fetch vehicle
	vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100&vehicle_id=" + strconv.Itoa(request.VehicleId)
	var vehicleResponses []transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.GetArray(vehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	// Create WorkOrder entity
	entitieswo := transactionworkshopentities.WorkOrder{
		// Default values
		WorkOrderDocumentNumber: defaultWorkOrderDocumentNumber,
		WorkOrderStatusId:       utils.WoStatDraft,
		WorkOrderDate:           currentDate,
		CPCcode:                 defaultCPCcode,
		ServiceAdvisor:          defaultServiceAdvisorId,
		WorkOrderTypeId:         workOrderTypeId,
		BookingSystemNumber:     request.BookingSystemNumber,
		EstimationSystemNumber:  request.EstimationSystemNumber,
		ServiceSite:             "OD - Service On Dealer",
		VehicleChassisNumber:    vehicleResponses[0].VehicleCode,

		// Provided values
		BrandId:                  request.BrandId,
		ModelId:                  request.ModelId,
		VariantId:                request.VariantId,
		VehicleId:                request.VehicleId,
		CustomerId:               request.CustomerId,
		BillableToId:             request.BilltoCustomerId,
		FromEra:                  request.FromEra,
		QueueNumber:              request.QueueSystemNumber,
		ArrivalTime:              request.WorkOrderArrivalTime,
		ServiceMileage:           request.WorkOrderCurrentMileage,
		Storing:                  request.Storing,
		Remark:                   request.WorkOrderRemark,
		ProfitCenterId:           request.WorkOrderProfitCenterId,
		CostCenterId:             request.DealerRepresentativeId,
		CampaignId:               request.CampaignId,
		CompanyId:                request.CompanyId,
		CPTitlePrefix:            request.Titleprefix,
		ContactPersonName:        request.NameCust,
		ContactPersonPhone:       request.PhoneCust,
		ContactPersonMobile:      request.MobileCust,
		ContactPersonContactVia:  request.ContactVia,
		EraNumber:                request.WorkOrderEraNo,
		EraExpiredDate:           request.WorkOrderEraExpiredDate,
		InsurancePolicyNumber:    request.WorkOrderInsurancePolicyNo,
		InsuranceExpiredDate:     request.WorkOrderInsuranceExpiredDate,
		InsuranceClaimNumber:     request.WorkOrderInsuranceClaimNo,
		InsurancePersonInCharge:  request.WorkOrderInsurancePic,
		InsuranceOwnRisk:         request.WorkOrderInsuranceOwnRisk,
		InsuranceWorkOrderNumber: request.WorkOrderInsuranceWONumber,
		EstTime:                  request.EstimationDuration,
		CustomerExpress:          request.CustomerExpress,
		LeaveCar:                 request.LeaveCar,
		CarWash:                  request.CarWash,
		PromiseDate:              request.PromiseDate,
		PromiseTime:              request.PromiseTime,
		FSCouponNo:               request.FSCouponNo,
		Notes:                    request.Notes,
		Suggestion:               request.Suggestion,
		DPAmount:                 request.DownpaymentAmount,
	}

	if err := tx.Create(&entitieswo).Error; err != nil {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create work order",
			Err:        err,
		}
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Memperbarui status pemesanan dan estimasi jika Booking_System_No atau Estim_System_No tidak nol

	var (
		batchSystemNo       int
		bookingStatusClosed int
		bookingSystemNo     = request.BookingSystemNumber
		estimationSystemNo  = request.EstimationSystemNumber
	)
	bookingStatusClosed = 8

	// Update related records if necessary
	if bookingSystemNo != 0 {
		if batchSystemNo == 0 {
			var batchSystemNoResult struct {
				BatchSystemNo int
			}
			if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Select("batch_system_number").
				Where("booking_system_number = ?", bookingSystemNo).
				Scan(&batchSystemNoResult).Error; err != nil {
				fmt.Println(batchSystemNoResult)
				return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve batch system number from the database",
					Err:        err,
				}
			}
			batchSystemNo = batchSystemNoResult.BatchSystemNo
		}

		// Update BOOKING_STATUS
		if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
			Where("booking_system_number = ?", bookingSystemNo).
			Update("document_status_id", bookingStatusClosed).Error; err != nil {
			return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update booking status",
				Err:        err,
			}
		}
	}

	if estimationSystemNo != 0 {
		if batchSystemNo == 0 {
			var batchSystemNoResult struct {
				BatchSystemNo int
			}
			if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Select("batch_system_number").
				Where("estimation_system_number = ?", estimationSystemNo).
				Scan(&batchSystemNoResult).Error; err != nil {
				return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve batch system number from the database",
					Err:        err,
				}
			}
			batchSystemNo = batchSystemNoResult.BatchSystemNo
		}

		// Update ESTIM_STATUS
		if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
			Where("estimation_system_number = ?", estimationSystemNo).
			Update("document_status_id", bookingStatusClosed).Error; err != nil {
			return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update estimation status",
				Err:        err,
			}
		}
	}

	// Update BATCH_STATUS
	if batchSystemNo != 0 {
		if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
			Where("batch_system_number = ?", batchSystemNo).
			Update("document_status_id", bookingStatusClosed).Error; err != nil {
			return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update batch status",
				Err:        err,
			}
		}
	}

	return entitieswo, nil
}

func (r *WorkOrderRepositoryImpl) GetAllBooking(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var tableStruct transactionworkshoppayloads.WorkOrderBooking

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	whereQuery = whereQuery.Where("booking_system_number != 0 OR estimation_system_number != 0")

	var workOrders []transactionworkshoppayloads.WorkOrderBooking
	if err := whereQuery.Find(&workOrders).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	var convertedResponses []transactionworkshoppayloads.WorkOrderBookingResponse

	for _, workOrderReq := range workOrders {
		var (
			workOrderRes transactionworkshoppayloads.WorkOrderBookingResponse
		)

		// Fetch data brand from external services
		BrandURL := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(workOrderReq.BrandId)
		var getBrandResponse transactionworkshoppayloads.WorkOrderVehicleBrand
		if err := utils.Get(BrandURL, &getBrandResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch brand data from external service",
				Err:        err,
			}
		}

		// Fetch data model from external services
		ModelURL := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(workOrderReq.ModelId)
		var getModelResponse transactionworkshoppayloads.WorkOrderVehicleModel
		if err := utils.Get(ModelURL, &getModelResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch model data from external service",
				Err:        err,
			}
		}

		// Fetch vehicle data
		VehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100&vehicle_id=" + strconv.Itoa(workOrderReq.VehicleId)
		var vehicleResponses []transactionworkshoppayloads.VehicleResponse
		errVehicle := utils.GetArray(VehicleUrl, &vehicleResponses, nil)

		if errVehicle != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve vehicle data from the external API",
				Err:        errVehicle,
			}
		}

		if len(vehicleResponses) == 0 {
			log.Printf("No vehicle data found for vehicle_id %d at URL %s", workOrderReq.VehicleId, VehicleUrl)
			continue
		}

		workOrderRes = transactionworkshoppayloads.WorkOrderBookingResponse{
			WorkOrderDocumentNumber:    workOrderReq.WorkOrderDocumentNumber,
			WorkOrderSystemNumber:      workOrderReq.WorkOrderSystemNumber,
			BookingSystemNumber:        workOrderReq.BookingSystemNumber,
			EstimationSystemNumber:     workOrderReq.EstimationSystemNumber,
			ServiceRequestSystemNumber: workOrderReq.ServiceRequestSystemNumber,
			WorkOrderTypeId:            workOrderReq.WorkOrderTypeId,
			BrandId:                    workOrderReq.BrandId,
			BrandName:                  getBrandResponse.BrandName,
			VehicleCode:                vehicleResponses[0].VehicleCode,
			VehicleTnkb:                vehicleResponses[0].VehicleTnkb,
			ModelId:                    workOrderReq.ModelId,
			ModelName:                  getModelResponse.ModelName,
			VehicleId:                  workOrderReq.VehicleId,
			CustomerId:                 workOrderReq.CustomerId,
			WorkOrderStatusId:          workOrderReq.StatusId,
		}

		convertedResponses = append(convertedResponses, workOrderRes)
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"batch_system_number":           "", // Adjust if needed
			"booking_system_number":         response.BookingSystemNumber,
			"estimation_system_number":      response.EstimationSystemNumber,
			"service_request_system_number": response.ServiceRequestSystemNumber,
			"brand_id":                      response.BrandId,
			"brand_name":                    response.BrandName,
			"model_id":                      response.ModelId,
			"model_name":                    response.ModelName,
			"vehicle_id":                    response.VehicleId,
			"vehicle_chassis_number":        response.VehicleCode,
			"vehicle_tnkb":                  response.VehicleTnkb,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) GetBookingById(tx *gorm.DB, IdWorkorder int, id int, pagination pagination.Pagination) (transactionworkshoppayloads.WorkOrderBookingResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ? AND booking_system_number = ?", IdWorkorder, id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
			}
		}
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order from the database",
			Err:        err,
		}
	}

	// Fetch data brand from external API
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entity.BrandId)
	var brandResponse transactionworkshoppayloads.WorkOrderVehicleBrand
	errBrand := utils.Get(brandUrl, &brandResponse, nil)
	if errBrand != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        errVariant,
		}
	}

	// Fetch data colour from external API
	colourUrl := config.EnvConfigs.SalesServiceUrl + "unit-color-dropdown/" + strconv.Itoa(entity.BrandId)
	var colourResponses []transactionworkshoppayloads.WorkOrderVehicleColour
	errColour := utils.GetArray(colourUrl, &colourResponses, nil)
	if errColour != nil || len(colourResponses) == 0 {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve colour data from the external API",
			Err:        errColour,
		}
	}

	// Fetch data vehicle from external API
	vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100&vehicle_id=" + strconv.Itoa(entity.VehicleId)
	var vehicleResponses []transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.GetArray(vehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	// Fetch workorder details with pagination
	var workorderDetails []transactionworkshoppayloads.WorkOrderDetailResponse
	query := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ?", IdWorkorder).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit())
	errWorkOrderDetails := query.Find(&workorderDetails).Error
	if errWorkOrderDetails != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order details from the database",
			Err:        errWorkOrderDetails,
		}
	}

	// Fetch work order services
	var workorderServices []transactionworkshoppayloads.WorkOrderServiceResponse
	if err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ?", IdWorkorder).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderServices).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order services from the database",
			Err:        err,
		}
	}

	// Fetch work order vehicles
	var workorderVehicles []transactionworkshoppayloads.WorkOrderServiceVehicleResponse
	if err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ?", IdWorkorder).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderVehicles).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order vehicles from the database",
			Err:        err,
		}
	}

	// fetch data status work order
	WorkOrderStatusURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-status?work_order_status_id=" + strconv.Itoa(entity.WorkOrderStatusId)
	var getWorkOrderStatusResponses []transactionworkshoppayloads.WorkOrderStatusResponse // Use slice of WorkOrderStatusResponse
	if err := utils.Get(WorkOrderStatusURL, &getWorkOrderStatusResponses, nil); err != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order status data from external service",
			Err:        err,
		}
	}

	// fetch data type work order
	WorkOrderTypeURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-type?work_order_type_id=" + strconv.Itoa(entity.WorkOrderTypeId)
	//fmt.Println("Fetching Work Order Type data from:", WorkOrderTypeURL)
	var getWorkOrderTypeResponses []transactionworkshoppayloads.WorkOrderTypeResponse // Use slice of WorkOrderTypeResponse
	if err := utils.Get(WorkOrderTypeURL, &getWorkOrderTypeResponses, nil); err != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order type data from external service",
			Err:        err,
		}
	}

	var workOrderTypeName string
	if len(getWorkOrderTypeResponses) > 0 {
		workOrderTypeName = getWorkOrderTypeResponses[0].WorkOrderTypeName
	}

	payload := transactionworkshoppayloads.WorkOrderBookingResponse{
		WorkOrderSystemNumber:         entity.WorkOrderSystemNumber,
		WorkOrderDate:                 entity.WorkOrderDate.Format("2006-01-02"),
		WorkOrderDocumentNumber:       entity.WorkOrderDocumentNumber,
		WorkOrderTypeId:               entity.WorkOrderTypeId,
		WorkOrderTypeName:             workOrderTypeName,
		WorkOrderStatusId:             entity.WorkOrderStatusId,
		WorkOrderStatusName:           getWorkOrderStatusResponses[0].WorkOrderStatusName,
		ServiceAdvisorId:              entity.ServiceAdvisor,
		BrandId:                       entity.BrandId,
		BrandName:                     brandResponse.BrandName,
		ModelId:                       entity.ModelId,
		ModelName:                     modelResponse.ModelName,
		VariantId:                     entity.VariantId,
		VariantName:                   variantResponse.VariantName,
		VehicleId:                     entity.VehicleId,
		VehicleCode:                   vehicleResponses[0].VehicleCode,
		VehicleTnkb:                   vehicleResponses[0].VehicleTnkb,
		CustomerId:                    entity.CustomerId,
		BilltoCustomerId:              entity.BillableToId,
		CampaignId:                    entity.CampaignId,
		FromEra:                       entity.FromEra,
		WorkOrderProfitCenterId:       entity.ProfitCenterId,
		AgreementId:                   entity.AgreementBodyRepairId,
		BookingSystemNumber:           entity.BookingSystemNumber,
		EstimationSystemNumber:        entity.EstimationSystemNumber,
		ContractSystemNumber:          entity.ContractServiceSystemNumber,
		QueueSystemNumber:             entity.QueueNumber,
		WorkOrderArrivalTime:          entity.ArrivalTime,
		WorkOrderRemark:               entity.Remark,
		DealerRepresentativeId:        entity.CostCenterId,
		CompanyId:                     entity.CompanyId,
		Titleprefix:                   entity.CPTitlePrefix,
		NameCust:                      entity.ContactPersonName,
		PhoneCust:                     entity.ContactPersonPhone,
		MobileCust:                    entity.ContactPersonMobile,
		MobileCustAlternative:         entity.ContactPersonMobileAlternative,
		MobileCustDriver:              entity.ContactPersonMobileDriver,
		ContactVia:                    entity.ContactPersonContactVia,
		WorkOrderInsurancePolicyNo:    entity.InsurancePolicyNumber,
		WorkOrderInsuranceClaimNo:     entity.InsuranceClaimNumber,
		WorkOrderInsuranceExpiredDate: entity.InsuranceExpiredDate,
		WorkOrderEraExpiredDate:       entity.EraExpiredDate,
		PromiseDate:                   entity.PromiseDate,
		PromiseTime:                   entity.PromiseTime,
		EstimationDuration:            entity.EstTime,
		WorkOrderInsuranceOwnRisk:     entity.InsuranceOwnRisk,
		WorkOrderInsurancePic:         entity.InsurancePersonInCharge,
		WorkOrderInsuranceWONumber:    entity.InsuranceWorkOrderNumber,
		CustomerExpress:               entity.CustomerExpress,
		LeaveCar:                      entity.LeaveCar,
		CarWash:                       entity.CarWash,
		FSCouponNo:                    entity.FSCouponNo,
		Notes:                         entity.Notes,
		Suggestion:                    entity.Suggestion,
		DownpaymentAmount:             entity.DPAmount,
		WorkOrderDetailService: transactionworkshoppayloads.WorkOrderDetailsResponseRequest{
			DataRequest: workorderServices,
		},
		WorkOrderDetailVehicle: transactionworkshoppayloads.WorkOrderDetailsResponseVehicle{
			DataVehicle: workorderVehicles,
		},
		WorkOrderDetails: transactionworkshoppayloads.WorkOrderDetailsResponse{
			Page:       pagination.GetPage(),
			Limit:      pagination.GetLimit(),
			TotalPages: pagination.TotalPages,
			TotalRows:  int(pagination.TotalRows), // Convert int64 to int
			Data:       workorderDetails,
		},
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) SaveBooking(tx *gorm.DB, IdWorkorder int, id int, request transactionworkshoppayloads.WorkOrderBookingRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND booking_system_number = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order booking from the database"}
	}

	entity.WorkOrderSystemNumber = request.WorkOrderSystemNumber
	entity.BookingSystemNumber = request.BookingSystemNumber
	entity.BillableToId = request.BilltoCustomerId
	entity.FromEra = request.FromEra
	entity.QueueNumber = request.QueueSystemNumber
	entity.ArrivalTime = request.WorkOrderArrivalTime
	entity.ServiceMileage = request.WorkOrderCurrentMileage
	entity.Storing = request.Storing
	entity.Remark = request.WorkOrderRemark
	entity.Unregister = request.Unregistered
	entity.ProfitCenterId = request.WorkOrderProfitCenterId
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
	entity.InsuranceExpiredDate = request.WorkOrderInsuranceExpiredDate
	entity.InsuranceClaimNumber = request.WorkOrderInsuranceClaimNo
	entity.InsurancePersonInCharge = request.WorkOrderInsurancePic
	entity.InsuranceOwnRisk = request.WorkOrderInsuranceOwnRisk
	entity.InsuranceWorkOrderNumber = request.WorkOrderInsuranceWONumber

	//page2
	entity.EstTime = request.EstimationDuration
	entity.CustomerExpress = request.CustomerExpress
	entity.LeaveCar = request.LeaveCar
	entity.CarWash = request.CarWash
	entity.PromiseDate = request.PromiseDate
	entity.PromiseTime = request.PromiseTime
	entity.FSCouponNo = request.FSCouponNo
	entity.Notes = request.Notes
	entity.Suggestion = request.Suggestion
	entity.DPAmount = request.DownpaymentAmount

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order booking"}
	}

	return true, nil
}

// uspg_wtWorkOrder0_Insert
// IF @Option = 50
// --USE FOR : * INSERT NEW DATA FROM AFFILIATED
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (r *WorkOrderRepositoryImpl) GetAllAffiliated(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tableStruct := transactionworkshoppayloads.WorkOrderAffiliate{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Add the additional where condition
	whereQuery = whereQuery.Where("service_request_system_number != 0")

	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	var convertedResponses []transactionworkshoppayloads.WorkOrderAffiliateGetResponse

	for rows.Next() {

		var (
			workOrderReq transactionworkshoppayloads.WorkOrderAffiliate
			workOrderRes transactionworkshoppayloads.WorkOrderAffiliateGetResponse
		)

		if err := rows.Scan(
			&workOrderReq.WorkOrderSystemNumber,
			&workOrderReq.WorkOrderDocumentNumber,
			&workOrderReq.ServiceRequestSystemNumber,
			&workOrderReq.BrandId,
			&workOrderReq.ModelId,
			&workOrderReq.VehicleId,
			&workOrderReq.CompanyId,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to scan the work order affiliated",
				Err:        err,
			}
		}

		// Fetch data brand from external services
		BrandURL := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(workOrderReq.BrandId)
		//fmt.Println("Fetching Brand data from:", BrandURL)
		var getBrandResponse transactionworkshoppayloads.WorkOrderVehicleBrand
		if err := utils.Get(BrandURL, &getBrandResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch brand data from external service",
				Err:        err,
			}
		}

		// Fetch data model from external services
		ModelURL := config.EnvConfigs.SalesServiceUrl + "unit-model/" + strconv.Itoa(workOrderReq.ModelId)
		//fmt.Println("Fetching Model data from:", ModelURL)
		var getModelResponse transactionworkshoppayloads.WorkOrderVehicleModel
		if err := utils.Get(ModelURL, &getModelResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch model data from external service",
				Err:        err,
			}
		}

		VehicleURL := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100&vehicle_id=" + strconv.Itoa(workOrderReq.VehicleId)
		//fmt.Println("Fetching Vehicle data from:", VehicleURL)
		var getVehicleResponse []transactionworkshoppayloads.VehicleResponse
		if err := utils.GetArray(VehicleURL, &getVehicleResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch vehicle data from external service",
				Err:        err,
			}
		}

		// fetch data service request from internal services
		ServiceRequestURL := config.EnvConfigs.AfterSalesServiceUrl + "service-request/" + strconv.Itoa(workOrderReq.ServiceRequestSystemNumber)
		fmt.Println("Fetching Service Request data from:", ServiceRequestURL)
		var getServiceRequestResponse transactionworkshoppayloads.ServiceRequestResponse
		if err := utils.Get(ServiceRequestURL, &getServiceRequestResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch service request data from internal service",
				Err:        err,
			}
		}

		// fetch data company from internal services
		CompanyURL := config.EnvConfigs.GeneralServiceUrl + "company/" + strconv.Itoa(workOrderReq.CompanyId)
		fmt.Println("Fetching Company data from:", CompanyURL)
		var getCompanyResponse transactionworkshoppayloads.CompanyResponse
		if err := utils.Get(CompanyURL, &getCompanyResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch company data from internal service",
				Err:        err,
			}
		}

		workOrderRes = transactionworkshoppayloads.WorkOrderAffiliateGetResponse{
			WorkOrderDocumentNumber:      workOrderReq.WorkOrderDocumentNumber,
			WorkOrderSystemNumber:        workOrderReq.WorkOrderSystemNumber,
			ServiceRequestSystemNumber:   workOrderReq.ServiceRequestSystemNumber,
			ServiceRequestDate:           getServiceRequestResponse.ServiceRequestDate,
			ServiceRequestDocumentNumber: getServiceRequestResponse.ServiceRequestDocumentNumber,
			BrandId:                      workOrderReq.BrandId,
			BrandName:                    getBrandResponse.BrandName,
			VehicleCode:                  getVehicleResponse[0].VehicleCode,
			VehicleTnkb:                  getVehicleResponse[0].VehicleTnkb,
			ModelId:                      workOrderReq.ModelId,
			ModelName:                    getModelResponse.ModelName,
			VehicleId:                    workOrderReq.VehicleId,
			CompanyId:                    workOrderReq.CompanyId,
			CompanyName:                  getCompanyResponse.CompanyName,
		}

		convertedResponses = append(convertedResponses, workOrderRes)
	}

	var mapResponses []map[string]interface{}

	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"service_request_system_number":   response.ServiceRequestSystemNumber,
			"service_request_date":            response.ServiceRequestDate,
			"service_request_document_number": response.ServiceRequestDocumentNumber,
			"brand_id":                        response.BrandId,
			"brand_name":                      response.BrandName,
			"model_id":                        response.ModelId,
			"model_name":                      response.ModelName,
			"vehicle_id":                      response.VehicleId,
			"vehicle_chassis_number":          response.VehicleCode,
			"vehicle_tnkb":                    response.VehicleTnkb,
			"company_id":                      response.CompanyId,
			"company_name":                    response.CompanyName,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil

}

func (r *WorkOrderRepositoryImpl) GetAffiliatedById(tx *gorm.DB, IdWorkorder int, id int, pagination pagination.Pagination) (transactionworkshoppayloads.WorkOrderAffiliateResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ? AND service_request_system_number = ?", IdWorkorder, id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
			}
		}
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order from the database",
			Err:        err,
		}
	}

	// Fetch data brand from external API
	brandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entity.BrandId)
	var brandResponse transactionworkshoppayloads.WorkOrderVehicleBrand
	errBrand := utils.Get(brandUrl, &brandResponse, nil)
	if errBrand != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        errVariant,
		}
	}

	// Fetch data colour from external API
	colourUrl := config.EnvConfigs.SalesServiceUrl + "unit-color-dropdown/" + strconv.Itoa(entity.BrandId)
	var colourResponses []transactionworkshoppayloads.WorkOrderVehicleColour
	errColour := utils.GetArray(colourUrl, &colourResponses, nil)
	if errColour != nil || len(colourResponses) == 0 {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve colour data from the external API",
			Err:        errColour,
		}
	}

	// Fetch data vehicle from external API
	vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100&vehicle_id=" + strconv.Itoa(entity.VehicleId)
	var vehicleResponses []transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.GetArray(vehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	// Fetch workorder details with pagination
	var workorderDetails []transactionworkshoppayloads.WorkOrderDetailResponse
	query := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ?", IdWorkorder).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit())
	errWorkOrderDetails := query.Find(&workorderDetails).Error
	if errWorkOrderDetails != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order details from the database",
			Err:        errWorkOrderDetails,
		}
	}

	// Fetch work order services
	var workorderServices []transactionworkshoppayloads.WorkOrderServiceResponse
	if err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ?", IdWorkorder).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderServices).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order services from the database",
			Err:        err,
		}
	}

	// Fetch work order vehicles
	var workorderVehicles []transactionworkshoppayloads.WorkOrderServiceVehicleResponse
	if err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ?", IdWorkorder).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderVehicles).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order vehicles from the database",
			Err:        err,
		}
	}

	// fetch data status work order
	WorkOrderStatusURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-status?work_order_status_id=" + strconv.Itoa(entity.WorkOrderStatusId)
	var getWorkOrderStatusResponses []transactionworkshoppayloads.WorkOrderStatusResponse // Use slice of WorkOrderStatusResponse
	if err := utils.Get(WorkOrderStatusURL, &getWorkOrderStatusResponses, nil); err != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order status data from external service",
			Err:        err,
		}
	}

	// fetch data type work order
	WorkOrderTypeURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-type?work_order_type_id=" + strconv.Itoa(entity.WorkOrderTypeId)
	//fmt.Println("Fetching Work Order Type data from:", WorkOrderTypeURL)
	var getWorkOrderTypeResponses []transactionworkshoppayloads.WorkOrderTypeResponse // Use slice of WorkOrderTypeResponse
	if err := utils.Get(WorkOrderTypeURL, &getWorkOrderTypeResponses, nil); err != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order type data from external service",
			Err:        err,
		}
	}

	var workOrderTypeName string
	if len(getWorkOrderTypeResponses) > 0 {
		workOrderTypeName = getWorkOrderTypeResponses[0].WorkOrderTypeName
	}

	payload := transactionworkshoppayloads.WorkOrderAffiliateResponse{
		WorkOrderSystemNumber:         entity.WorkOrderSystemNumber,
		WorkOrderDate:                 entity.WorkOrderDate.Format("2006-01-02"),
		WorkOrderDocumentNumber:       entity.WorkOrderDocumentNumber,
		WorkOrderTypeId:               entity.WorkOrderTypeId,
		WorkOrderTypeName:             workOrderTypeName,
		WorkOrderStatusId:             entity.WorkOrderStatusId,
		WorkOrderStatusName:           getWorkOrderStatusResponses[0].WorkOrderStatusName,
		ServiceAdvisorId:              entity.ServiceAdvisor,
		ServiceSite:                   entity.ServiceSite,
		BrandId:                       entity.BrandId,
		BrandName:                     brandResponse.BrandName,
		ModelId:                       entity.ModelId,
		ModelName:                     modelResponse.ModelName,
		VariantId:                     entity.VariantId,
		VariantName:                   variantResponse.VariantName,
		VehicleId:                     entity.VehicleId,
		VehicleCode:                   vehicleResponses[0].VehicleCode,
		VehicleTnkb:                   vehicleResponses[0].VehicleTnkb,
		CustomerId:                    entity.CustomerId,
		BilltoCustomerId:              entity.BillableToId,
		CampaignId:                    entity.CampaignId,
		FromEra:                       entity.FromEra,
		WorkOrderProfitCenterId:       entity.ProfitCenterId,
		AgreementId:                   entity.AgreementBodyRepairId,
		BookingSystemNumber:           entity.BookingSystemNumber,
		EstimationSystemNumber:        entity.EstimationSystemNumber,
		ContractSystemNumber:          entity.ContractServiceSystemNumber,
		QueueSystemNumber:             entity.QueueNumber,
		WorkOrderArrivalTime:          entity.ArrivalTime,
		WorkOrderRemark:               entity.Remark,
		DealerRepresentativeId:        entity.CostCenterId,
		CompanyId:                     entity.CompanyId,
		Titleprefix:                   entity.CPTitlePrefix,
		NameCust:                      entity.ContactPersonName,
		PhoneCust:                     entity.ContactPersonPhone,
		MobileCust:                    entity.ContactPersonMobile,
		MobileCustAlternative:         entity.ContactPersonMobileAlternative,
		MobileCustDriver:              entity.ContactPersonMobileDriver,
		ContactVia:                    entity.ContactPersonContactVia,
		WorkOrderInsurancePolicyNo:    entity.InsurancePolicyNumber,
		WorkOrderInsuranceClaimNo:     entity.InsuranceClaimNumber,
		WorkOrderInsuranceExpiredDate: entity.InsuranceExpiredDate,
		WorkOrderEraExpiredDate:       entity.EraExpiredDate,
		PromiseDate:                   entity.PromiseDate,
		PromiseTime:                   entity.PromiseTime,
		EstimationDuration:            entity.EstTime,
		WorkOrderInsuranceOwnRisk:     entity.InsuranceOwnRisk,
		WorkOrderInsurancePic:         entity.InsurancePersonInCharge,
		WorkOrderInsuranceWONumber:    entity.InsuranceWorkOrderNumber,
		CustomerExpress:               entity.CustomerExpress,
		LeaveCar:                      entity.LeaveCar,
		CarWash:                       entity.CarWash,
		FSCouponNo:                    entity.FSCouponNo,
		Notes:                         entity.Notes,
		Suggestion:                    entity.Suggestion,
		DownpaymentAmount:             entity.DPAmount,
		WorkOrderDetailService: transactionworkshoppayloads.WorkOrderDetailsResponseRequest{
			DataRequest: workorderServices,
		},
		WorkOrderDetailVehicle: transactionworkshoppayloads.WorkOrderDetailsResponseVehicle{
			DataVehicle: workorderVehicles,
		},
		WorkOrderDetails: transactionworkshoppayloads.WorkOrderDetailsResponse{
			Page:       pagination.GetPage(),
			Limit:      pagination.GetLimit(),
			TotalPages: pagination.TotalPages,
			TotalRows:  int(pagination.TotalRows), // Convert int64 to int
			Data:       workorderDetails,
		},
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) NewAffiliated(tx *gorm.DB, IdWorkorder int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrder{

		WorkOrderSystemNumber: request.WorkOrderSystemNumber,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the work order affiliate"}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) SaveAffiliated(tx *gorm.DB, IdWorkorder int, id int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND affiliate_id = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order affiliate from the database"}
	}

	entity.WorkOrderSystemNumber = request.WorkOrderSystemNumber

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order affiliate"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewStatus(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse) {
	var statuses []transactionworkshopentities.WorkOrderMasterStatus

	// Apply filters
	query := utils.ApplyFilter(tx, filter)

	// Execute the query
	if err := query.Find(&statuses).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order statuses from the database",
			Err:        err,
		}
	}

	return statuses, nil
}

func (r *WorkOrderRepositoryImpl) AddStatus(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderStatusRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderMasterStatus{
		WorkOrderStatusCode:        request.WorkOrderStatusCode,
		WorkOrderStatusDescription: request.WorkOrderStatusName,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateStatus(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderStatusRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterStatus
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterStatus{}).Where("work_order_status_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order status from the database"}
	}

	entity.WorkOrderStatusCode = request.WorkOrderStatusCode
	entity.WorkOrderStatusDescription = request.WorkOrderStatusName

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to update work order status"}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteStatus(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterStatus
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterStatus{}).Where("work_order_status_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order status from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete work order status"}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterType

	query := utils.ApplyFilter(tx, filter)

	if err := query.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}
	return types, nil
}

func (r *WorkOrderRepositoryImpl) AddType(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderTypeRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderMasterType{
		WorkOrderTypeCode:        request.WorkOrderTypeCode,
		WorkOrderTypeDescription: request.WorkOrderTypeName,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateType(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderTypeRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterType{}).Where("work_order_type_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}

	entity.WorkOrderTypeCode = request.WorkOrderTypeCode
	entity.WorkOrderTypeDescription = request.WorkOrderTypeName

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to update work order type"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteType(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterType{}).Where("work_order_type_id = ?", id).First(&entity).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete work order type"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewBill(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderBillable, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterBillAbleto

	if err := tx.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}

	var getBillables []transactionworkshoppayloads.WorkOrderBillable
	for _, t := range types {
		getBillables = append(getBillables, transactionworkshoppayloads.WorkOrderBillable{
			BillableToID:   t.WorkOrderBillabletoId,
			BillableToCode: t.WorkOrderBillabletoCode,
			BillableToName: t.WorkOrderBillabletoName,
		})
	}

	return getBillables, nil
}

func (r *WorkOrderRepositoryImpl) AddBill(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderBillableRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderMasterBillAbleto{
		WorkOrderBillabletoName: request.BillableToName,
		WorkOrderBillabletoCode: request.BillableToCode,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateBill(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderBillableRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterBillAbleto
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterBillAbleto{}).Where("billable_to_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve billable data from the database"}
	}

	entity.WorkOrderBillabletoName = request.BillableToName
	entity.WorkOrderBillabletoCode = request.BillableToCode

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to update billable data"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteBill(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterBillAbleto
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterBillAbleto{}).Where("billable_to_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve billable data from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete billable data"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewTrxType(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterTrxType

	if err := tx.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}

	var payloadTypes []transactionworkshoppayloads.WorkOrderTransactionType
	for _, t := range types {
		payloadTypes = append(payloadTypes, transactionworkshoppayloads.WorkOrderTransactionType{
			TransactionTypeCode: t.WorkOrderTrxTypeCode,
			TransactionTypeName: t.WorkOrderTrxTypeDescription,
		})
	}

	return payloadTypes, nil
}

func (r *WorkOrderRepositoryImpl) AddTrxType(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderMasterTrxType{
		WorkOrderTrxTypeDescription: request.TransactionTypeName,
		WorkOrderTrxTypeCode:        request.TransactionTypeCode,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateTrxType(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterTrxType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterTrxType{}).Where("work_order_transaction_type_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve billable data from the database"}
	}

	entity.WorkOrderTrxTypeDescription = request.TransactionTypeName
	entity.WorkOrderTrxTypeCode = request.TransactionTypeCode

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to update billable data"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteTrxType(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterTrxType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterTrxType{}).Where("work_order_transaction_type_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve billable data from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete billable data"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewTrxTypeSo(tx *gorm.DB) ([]transactionworkshoppayloads.WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterTrxSoType

	if err := tx.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}

	var payloadTypes []transactionworkshoppayloads.WorkOrderTransactionType
	for _, t := range types {
		payloadTypes = append(payloadTypes, transactionworkshoppayloads.WorkOrderTransactionType{
			TransactionTypeCode: t.WorkOrderTrxTypeSoCode,
			TransactionTypeName: t.WorkOrderTrxTypeSoDescription,
		})
	}

	return payloadTypes, nil

}

func (r *WorkOrderRepositoryImpl) AddTrxTypeSo(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderMasterTrxSoType{
		WorkOrderTrxTypeSoDescription: request.TransactionTypeName,
		WorkOrderTrxTypeSoCode:        request.TransactionTypeCode,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateTrxTypeSo(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterTrxSoType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterTrxSoType{}).Where("work_order_transaction_type_so_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve billable data from the database"}
	}

	entity.WorkOrderTrxTypeSoDescription = request.TransactionTypeName
	entity.WorkOrderTrxTypeSoCode = request.TransactionTypeCode

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to update billable data"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteTrxTypeSo(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterTrxSoType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterTrxSoType{}).Where("work_order_transaction_type_so_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve billable data from the database"}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to delete billable data"}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewDropPoint(*gorm.DB) ([]transactionworkshoppayloads.WorkOrderDropPoint, *exceptions.BaseErrorResponse) {
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

func (r *WorkOrderRepositoryImpl) NewVehicleBrand(*gorm.DB) ([]transactionworkshoppayloads.WorkOrderVehicleBrand, *exceptions.BaseErrorResponse) {
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

func (r *WorkOrderRepositoryImpl) NewVehicleModel(_ *gorm.DB, brandId int) ([]transactionworkshoppayloads.WorkOrderVehicleModel, *exceptions.BaseErrorResponse) {
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

func (s *WorkOrderRepositoryImpl) DeleteVehicleServiceMultiId(tx *gorm.DB, Id int, DetailIds []int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderServiceVehicle
	err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).Where("work_order_system_number = ? AND work_order_service_id IN (?)", Id, DetailIds).First(&entity).Error
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

func (s *WorkOrderRepositoryImpl) DeleteDetailWorkOrderMultiId(tx *gorm.DB, Id int, DetailIds []int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderDetail
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).Where("work_order_system_number = ? AND work_order_detail_id = ? IN (?)", Id, DetailIds).First(&entity).Error
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

func (s *WorkOrderRepositoryImpl) DeleteRequestMultiId(tx *gorm.DB, Id int, DetailIds []int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).Where("work_order_system_number = ? AND work_order_service_id IN (?)", Id, DetailIds).First(&entity).Error
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

// usp_comLookUp
// IF @strEntity = 'CustomerByTypeAndAddress'--CUSTOMER MASTER
func (s *WorkOrderRepositoryImpl) ChangeBillTo(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.ChangeBillToRequest) (bool, *exceptions.BaseErrorResponse) {
	var existingWorkOrder struct {
		WorkOrderOperationItemLine int
	}

	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("work_order_operation_item_line").
		Where("work_order_system_number = ? AND ISNULL(invoice_system_number, 0) <> 0 AND transaction_type_id = 'E' ", workOrderId).
		First(&existingWorkOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Failed to retrieve work order from the database",
				Err:        err,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "External detail has already been invoiced",
			Err:        err}
	}

	var entity transactionworkshopentities.WorkOrder
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Failed to retrieve work order from the database"}
	}

	entity.CustomerId = request.BillToCustomerId

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update the billable to data"}
	}

	// Update data IF @strEntity = 'CustomerByTypeAndAddressFilter'--CUSTOMER MASTER

	return true, nil
}

func (s *WorkOrderRepositoryImpl) ChangePhoneNo(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.ChangePhoneNoRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Failed to retrieve work order from the database"}
	}

	entity.ContactPersonPhone = request.PhoneNo

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update the phone number data"}
	}

	return true, nil
}

// uspg_wtWorkOrder2_Update
// IF @Option = 14
// --USE FOR : CONFIRM PRICE
func (s *WorkOrderRepositoryImpl) ConfirmPrice(tx *gorm.DB, workOrderId int, idwos []int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	var markupPercentage, totalPackage, totalOpr, totalPart, totalOil, totalMaterial, totalConsumableMaterial, totalSublet, totalAccs float64
	var invSysNo, woOprItemLine, greyMarket int
	var oprItemCode, vehicleChassisNo string
	var total, totalDisc, totalAfterDisc, totalVat, totalNonVat float64
	var greyMarketMarkupPercentageExists bool

	// Fetch WorkOrder data
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Select(`
			company_id, ISNULL(contract_service_system_number, 0) AS ContractServSysNo, additional_discount_status_approval, 
			vehicle_id, model_id, currency_id, ISNULL(vat_tax_rate, 0) AS VatTaxRate, 
			ISNULL(work_order_site_type_id, '') AS SiteCode, ISNULL(campaign_id, '') AS CampaignCode
		`).
		Where("work_order_system_number = ?", workOrderId).
		Scan(&entity).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Failed to retrieve work order from the database",
			Err:        err,
		}
	}

	type WorkOrderDetailResult struct {
		LineType     int
		InvoiceSysNo int
		WhsGroup     int
	}

	var detailResult WorkOrderDetailResult

	// Fetch line type, invoice system number, etc.
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("ISNULL(line_type_id, 0) AS LineType, ISNULL(invoice_system_number, 0) AS InvoiceSysNo, ISNULL(warehouse_id, 0) AS WhsGroup").
		Where("work_order_system_number = ? AND work_order_operation_item_line IN (?)", workOrderId, idwos).
		Scan(&detailResult).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Failed to retrieve work order detail from the database",
			Err:        err,
		}
	}

	invSysNo = detailResult.InvoiceSysNo

	// Check if the grey market markup percentage exists in the comGenVariable table
	err = tx.Table("dms_microservices_general_dev.dbo.mtr_company").
		Select("1").
		Where("company_code = 120027"). // Indomobil Trada Nasional - TB Simatupang
		Limit(1).
		Find(&greyMarketMarkupPercentageExists).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check grey market markup percentage",
			Err:        err,
		}
	}

	// fetch vehicle from external API
	vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100&vehicle_id=" + strconv.Itoa(entity.VehicleId)
	var vehicleResponses []transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.GetArray(vehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	vehicleChassisNo = vehicleResponses[0].VehicleCode

	// Check if the vehicle is in the grey market by looking up the vehicle chassis number
	err = tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle").
		Select("ISNULL(vehicle_is_grey_market, 0)").
		Where("vehicle_chassis_no = ?", vehicleChassisNo).
		Scan(&greyMarket).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check vehicle grey market status",
			Err:        err,
		}
	}

	// Check for grey market markup and update
	if greyMarketMarkupPercentageExists && greyMarket == 1 {
		if invSysNo == 0 {
			markupPercentage = 40 // Set the markup percentage based on logic
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Where("work_order_system_number = ? AND work_order_operation_item_line = ? AND operation_item_code = ?", workOrderId, woOprItemLine, oprItemCode).
				Updates(map[string]interface{}{
					"operation_item_price":                   gorm.Expr("operation_item_price + (operation_item_price * ? / 100)", markupPercentage),
					"operation_item_discount_amount":         gorm.Expr("(operation_item_price + (operation_item_price * ? / 100)) * operation_item_discount_percent / 100", markupPercentage),
					"operation_item_discount_request_amount": gorm.Expr("(operation_item_price + (operation_item_price * ? / 100)) * operation_item_discount_request_amount / 100", markupPercentage),
				}).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update work order details with grey market markup",
					Err:        err,
				}
			}
		}

		// Compute totals
		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("SUM(ROUND(ISNULL(operation_item_price, 0), 0, 0))").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypePackage).
			Scan(&totalPackage).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total package",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOperation).
			Scan(&totalOpr).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total operation",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSparepart).
			Scan(&totalPart).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total sparepart",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOil).
			Scan(&totalOil).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total oil",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeMaterial).
			Scan(&totalMaterial).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total material",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LineTypeSublet).
			Scan(&totalSublet).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total sublet",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeAccesories).
			Scan(&totalAccs).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total accessories",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeConsumableMaterial).
			Scan(&totalConsumableMaterial).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total consumable material",
				Err:        err,
			}
		}

		// Calculate total and discounts
		total = totalPackage + totalOpr + totalPart + totalOil + totalMaterial + totalSublet + totalAccs + totalConsumableMaterial

		// Calculate discounts and VAT
		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select(`
				SUM(CASE 
					WHEN line_type_id = ? THEN 
						CASE WHEN approval_id = "20" AND ISNULL(operation_item_discount_request_amount, 0) > 0 
						THEN ISNULL(operation_item_discount_request_amount, 0) 
						ELSE ISNULL(operation_item_discount_amount, 0) 
						END 
					ELSE 
						CASE WHEN approval_id = "20" AND ISNULL(operation_item_discount_request_amount, 0) > 0 
						THEN ISNULL(operation_item_discount_request_amount, 0) 
						ELSE ISNULL(operation_item_discount_amount, 0) 
						END 
						* 
						CASE WHEN LINE_TYPE <> ? THEN ISNULL(frt_quantity, 0) 
						ELSE CASE WHEN ISNULL(supply_quantity, 0) > 0 
						THEN ISNULL(supply_quantity, 0) ELSE ISNULL(frt_quantity, 0) END 
						END 
				END)`, utils.LinetypePackage, utils.LinetypeOperation).
			Scan(&totalDisc).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total discount",
				Err:        err,
			}
		}

		totalDisc = math.Round(totalDisc)
		totalAfterDisc = total - totalDisc
		totalVat = (totalAfterDisc * *entity.VATTaxRate) / 100
		totalNonVat = totalAfterDisc + totalVat

		// Update totals in the WorkOrder table
		err = tx.Model(&transactionworkshopentities.WorkOrder{}).
			Where("work_order_system_number = ?", workOrderId).
			Updates(map[string]interface{}{
				"total_package":             totalPackage,
				"total_operation":           totalOpr,
				"total_part":                totalPart,
				"total_oil":                 totalOil,
				"total_material":            totalMaterial,
				"total_consumable_material": totalConsumableMaterial,
				"total_sublet":              totalSublet,
				"total_price_accessories":   totalAccs,
				"total_discount":            totalDisc,
				"total_after_discount":      totalAfterDisc,
				"total_vat":                 totalVat,
				"total_after_vat":           totalNonVat,
			}).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update work order totals",
				Err:        err,
			}
		}
	}

	return true, nil
}

// uspg_wtWorkOrder2_Update
// IF @Option = 6
// --USE FOR : CHECK DETAIL
func (s *WorkOrderRepositoryImpl) CheckDetail(tx *gorm.DB, workOrderId int, idwos []int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	var detailentity transactionworkshopentities.WorkOrderDetail

	// Fetch WorkOrder Detail data
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select(`
			line_type_id, ISNULL(invoice_system_number, 0) AS invoice_system_number,
			ISNULL(warehouse_id, 0) AS warehouse_group, 
			operation_item_code, transaction_type_id
		`).
		Where("work_order_system_number = ? AND work_order_operation_item_line IN (?)", workOrderId, idwos).
		First(&detailentity).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Failed to retrieve work order detail from the database",
			Err:        err,
		}
	}

	// Check transaction type
	if detailentity.TransactionTypeId != 2 {
		if detailentity.InvoiceSystemNumber == 0 {
			if detailentity.LineTypeId == utils.LinetypeOperation || detailentity.LineTypeId == utils.LinetypePackage {
				var supplyQty float64
				err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Select("ISNULL(supply_quantity, 0) AS supply_quantity").
					Where("work_order_system_number = ? AND work_order_operation_item_line = ? AND line_type_id <> ? AND line_type_id <> ?", workOrderId, idwos, utils.LinetypeOperation, utils.LinetypePackage).
					Scan(&supplyQty).Error

				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to check supply quantity",
						Err:        err,
					}
				}

				if supplyQty == 0 {
					var qtyAvail float64
					err = tx.Raw("EXEC dbo.uspg_amLocationStockItem_Select @Option = 1, @Company_Code = ?, @Period_Date = ?, @Whs_Code = '', @Loc_Code = '', @Item_Code = ?, @Whs_Group = ?, @UoM_Type = ?, @QtyResult = ? OUTPUT",
						entity.CompanyId, time.Now(), detailentity.OperationItemCode, detailentity.WarehouseId, "S", &qtyAvail).Error

					if err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to check quantity available",
							Err:        err,
						}
					}

					if qtyAvail == 0 {
						// Step 1: Create a temporary table
						err = tx.Exec(`
							CREATE TABLE #SUBS (
								SUBS_ITEM_CODE VARCHAR(15),
								ITEM_NAME VARCHAR(40),
								SUPPLY_QTY NUMERIC(7,2),
								SUBS_TYPE CHAR(2)
							)
						`).Error
						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to create temporary table",
								Err:        err,
							}
						}

						// Step 2: Insert data into the temporary table using the stored procedure
						err = tx.Exec(`
							INSERT INTO #SUBS 
							EXEC dbo.uspg_amSubstituteItem_Select @Option = 1, @Company_Code = ?, @Item_Code = ?, @Whs_Group = ?, @UoM_Type = ?
						`, entity.CompanyId, detailentity.OperationItemCode, detailentity.WarehouseId, "S").Error
						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to insert data into temporary table",
								Err:        err,
							}
						}

						// Step 3: Check if the original item is in the substitution table
						var exists bool
						err = tx.Raw(`
							SELECT CASE WHEN EXISTS (SELECT 1 FROM #SUBS WHERE SUBS_ITEM_CODE = ?) THEN 1 ELSE 0 END
						`, detailentity.OperationItemCode).Scan(&exists).Error
						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to check if the original item is in the substitution table",
								Err:        err,
							}
						}

						if !exists {
							// Step 4: Fetch the substitute items from the temporary table
							rows, err := tx.Raw(`
								SELECT SUBS_ITEM_CODE, ITEM_NAME, SUPPLY_QTY, SUBS_TYPE FROM #SUBS
							`).Rows()
							if err != nil {
								return false, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to fetch substitute items from the temporary table",
									Err:        err,
								}
							}
							defer rows.Close()

							// Step 5: Process each substitute item
							for rows.Next() {
								var substituteItem struct {
									SubsItemCode string
									ItemName     string
									SupplyQty    float64
									SubsType     string
								}
								err := rows.Scan(&substituteItem.SubsItemCode, &substituteItem.ItemName, &substituteItem.SupplyQty, &substituteItem.SubsType)
								if err != nil {
									return false, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Failed to process substitute items",
										Err:        err,
									}
								}

								// Step 6: Check and update the original item if not substituted
								if substituteItem.SubsType != "" {
									err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
										Where("work_order_system_number = ? AND operation_item_code = ? AND work_order_operation_item_line = ? AND substitute_id IS NULL",
											workOrderId, detailentity.OperationItemCode, idwos).
										Updates(map[string]interface{}{
											"substitute_id":   1,
											"substitute_type": "SUBSTITUTE_ITEM",
											"warehouse_id":    detailentity.WarehouseId,
										}).Error

									if err != nil {
										return false, &exceptions.BaseErrorResponse{
											StatusCode: http.StatusInternalServerError,
											Message:    "Failed to check and update the original item if not substituted",
											Err:        err,
										}
									}
								}

								// Step 7: Get vehicle brand and currency code
								err = tx.Model(&transactionworkshopentities.WorkOrder{}).
									Select("ISNULL(brand_id, '') AS vehicle_brand, ISNULL(currency_id, '') AS currency").
									Where("work_order_system_number = ?", workOrderId).
									Scan(&entity).Error
								if err != nil {
									return false, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Failed to get vehicle brand and currency code",
										Err:        err,
									}
								}

								// Step 8: Calculate the next line number for the work order
								var nextLineNumber int
								err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
									Select("COALESCE(MAX(work_order_operation_item_line), 0) + 1").
									Where("work_order_system_number = ?", workOrderId).
									Pluck("work_order_operation_item_line", &nextLineNumber).Error

								if err != nil {
									return false, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Failed to calculate the next line number for the work order",
										Err:        err,
									}
								}

								// Step 9: Insert the substitute item into wtWorkOrder2
								// Check if the substitute item already exists
								var count int64
								err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
									Where("work_order_system_number = ? AND operation_item_code = ? AND substritute_item_code = ?",
										workOrderId, detailentity.OperationItemCode, substituteItem.SubsItemCode).
									Count(&count).Error
								if err != nil {
									return false, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Failed to check the existence of the substitute item",
										Err:        err,
									}
								}

								// If the substitute item does not exist, insert it
								if count == 0 {
									// Get operation item price using custom logic
									var oprItemPrice float64

									// Fetch Opr_Item_Price
									oprItemPrice, _ = s.lookupRepo.GetOprItemPrice(tx, detailentity.LineTypeId, entity.CompanyId, detailentity.OperationId, entity.BrandId, entity.ModelId, detailentity.JobTypeId, entity.VariantId, entity.CurrencyId, "W", "1")

									// Apply markup to the item price
									oprItemPrice = oprItemPrice + 10.00 + (oprItemPrice * (5.00 / 100))

									// Insert the substitute item into wtWorkOrder2
									err = tx.Create(&transactionworkshopentities.WorkOrderDetail{
										WorkOrderSystemNumber:      workOrderId,
										WorkOrderOperationItemLine: nextLineNumber,
										LineType:                   "Substitute",
										WorkorderStatusId:          detailentity.WorkorderStatusId,
										OperationItemCode:          detailentity.OperationItemCode,
										OperationItemPrice:         oprItemPrice,
										SubstrituteItemCode:        substituteItem.SubsItemCode,
										SupplyQuantity:             substituteItem.SupplyQty,
										WarehouseId:                detailentity.WarehouseId,
									}).Error
									if err != nil {
										return false, &exceptions.BaseErrorResponse{
											StatusCode: http.StatusInternalServerError,
											Message:    "Failed to insert the substitute item into wtWorkOrder2",
											Err:        err,
										}
									}
								}
							}

							// Drop the temporary table
							err = tx.Exec(`DROP TABLE #SUBS`).Error
							if err != nil {
								return false, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to drop the temporary table",
									Err:        err,
								}
							}
						}

						// If substitute not found, return true with no errors
						return true, nil
					}
				}
			}

			//-- By default price of all items will be replaced with the price from PriceList
			//-- Exclude Fee from the replacing process
			var oprItemPrice, oprItemPriceDisc, discountPercent float64
			var markupAmount, markupPercentage float64
			var warrantyClaimType string

			if detailentity.LineTypeId == utils.LineTypeSublet {

				// Fetch Opr_Item_Price
				oprItemPrice, _ = s.lookupRepo.GetOprItemPrice(tx, detailentity.LineTypeId, entity.CompanyId, detailentity.OperationId, entity.BrandId, entity.ModelId, detailentity.JobTypeId, entity.VariantId, entity.CurrencyId, "W", "1")

				// Set markup percentage based on company ID
				if entity.CompanyId == 139 {
					markupPercentage = 11.00
				}

				// // Apply markup amount and percentage
				oprItemPrice = oprItemPrice + markupAmount + (oprItemPrice * (markupPercentage / 100))

				// Fetch Opr_Item_Disc_Percent
				oprItemPriceDisc, _ = s.lookupRepo.GetOprItemDisc(tx, detailentity.LineTypeId, utils.TrxTypeSoDeCentralize, detailentity.ItemId, entity.AgreementGeneralRepairId, entity.ProfitCenterId, detailentity.FrtQuantity*detailentity.OperationItemPrice, entity.CompanyId, entity.BrandId, entity.ContractServiceSystemNumber, "W", utils.EstWoOrderTypeId)

			} else {
				// Fetch Opr_Item_Price
				oprItemPrice, _ = s.lookupRepo.GetOprItemPrice(tx, detailentity.LineTypeId, entity.CompanyId, detailentity.OperationId, entity.BrandId, entity.ModelId, detailentity.JobTypeId, entity.VariantId, entity.CurrencyId, "W", "1")

				// Set markup percentage based on company ID
				if entity.CompanyId == 139 {
					markupPercentage = 11.00
				}

				// // Apply markup amount and percentage
				oprItemPrice = oprItemPrice + markupAmount + (oprItemPrice * (markupPercentage / 100))

				// Fetch Opr_Item_Disc_Percent
				oprItemPriceDisc, _ = s.lookupRepo.GetOprItemDisc(tx, detailentity.LineTypeId, utils.TrxTypeSoDeCentralize, detailentity.ItemId, entity.AgreementGeneralRepairId, entity.ProfitCenterId, detailentity.FrtQuantity*detailentity.OperationItemPrice, entity.CompanyId, entity.BrandId, entity.ContractServiceSystemNumber, "W", utils.EstWoOrderTypeId)

			}

			err = tx.Model(&masteritementities.Item{}).
				Where("item_code = ?", detailentity.OperationItemCode).
				Select("atpm_warranty_claim_type_id").Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch ATPM warranty claim type ID",
					Err:        err,
				}
			}

			// nilai warranty_claim_type_id berdasarkan kondisi
			if detailentity.LineTypeId != utils.LinetypeOperation && detailentity.LineTypeId != utils.LinetypePackage {
				warrantyClaimType = entity.ATPMWCFDocNo
			} else {
				warrantyClaimType = ""
			}

			// nilai discountPercent berdasarkan kondisi
			discountPercent = oprItemPriceDisc
			if oprItemPriceDisc == 0 {
				discountPercent = 0
			}

			// update work order detail
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Where("work_order_system_number = ? AND work_order_operation_item_line = ?", workOrderId, idwos).
				Updates(map[string]interface{}{
					"operation_item_price":                   oprItemPrice,
					"operation_item_discount_amount_percent": oprItemPriceDisc,
					"operation_item_discount_amount":         math.Round(oprItemPrice * (discountPercent / 100)),
					"operation_item_discount_request_amount": math.Round(oprItemPrice * (discountPercent / 100)),
					"warranty_claim_type_id":                 warrantyClaimType,
					"price_list_id":                          "",
				}).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update work order detail",
					Err:        err,
				}
			}

			var totals struct {
				TotalOpr                float64
				TotalPart               float64
				TotalOil                float64
				TotalMaterial           float64
				TotalConsumableMaterial float64
				TotalSublet             float64
				TotalAccs               float64
				Total                   float64
				TotalDisc               float64
				TotalAfterDisc          float64
				TotalVat                float64
				TotalAfterVat           float64
				TotalNonVat             float64
				TaxFree                 int
			}

			// Calculate totals
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
				Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOperation).
				Scan(&totals.TotalOpr).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to calculate total operation",
					Err:        err,
				}
			}

			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
				Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSparepart).
				Scan(&totals.TotalPart).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to calculate total sparepart",
					Err:        err,
				}
			}

			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
				Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOil).
				Scan(&totals.TotalOil).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to calculate total oil",
					Err:        err,
				}
			}

			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
				Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeMaterial).
				Scan(&totals.TotalMaterial).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to calculate total material",
					Err:        err,
				}
			}

			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
				Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeConsumableMaterial).
				Scan(&totals.TotalConsumableMaterial).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to calculate total consumable material",
					Err:        err,
				}
			}

			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
				Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LineTypeSublet).
				Scan(&totals.TotalSublet).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to calculate total sublet",
					Err:        err,
				}
			}

			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
				Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeAccesories).
				Scan(&totals.TotalAccs).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to calculate total accessories",
					Err:        err,
				}
			}

			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
				Where("work_order_system_number = ? AND line_type_id IN (?, ?, ?, ?, ?, ?, ?)",
					workOrderId, utils.LinetypePackage, utils.LinetypeOperation, utils.LinetypeSparepart, utils.LinetypeOil, utils.LinetypeMaterial, utils.LinetypeConsumableMaterial, utils.LineTypeSublet).
				Scan(&totals.Total).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to calculate total",
					Err:        err,
				}
			}

			// Calculate discounts and VAT
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select(`
					SUM(CASE
						WHEN line_type_id = ? THEN
							CASE WHEN approval_id = "20" AND ISNULL(operation_item_discount_request_amount, 0) > 0
							THEN ISNULL(operation_item_discount_request_amount, 0)
							ELSE ISNULL(operation_item_discount_amount, 0)
							END
						ELSE
							CASE WHEN approval_id = "20" AND ISNULL(operation_item_discount_request_amount, 0) > 0
							THEN ISNULL(operation_item_discount_request_amount, 0)
							ELSE ISNULL(operation_item_discount_amount, 0)
							END 
							
							*

							CASE WHEN LINE_TYPE <> ? THEN ISNULL(frt_quantity, 0)
							ELSE CASE WHEN ISNULL(supply_quantity, 0) > 0
							THEN ISNULL(supply_quantity, 0) ELSE ISNULL(frt_quantity, 0) END
							END
					END)`, utils.LinetypePackage, utils.LinetypeOperation).
				Scan(&totals.TotalDisc).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to calculate total discount",
					Err:        err,
				}
			}

			// Calculate TotalDisc and TotalAfterDisc
			totals.TotalDisc = math.Round(totals.TotalDisc)
			totals.TotalAfterDisc = math.Round(totals.Total - totals.TotalDisc)

			// Calculate totalNonVat
			var totalNonVat float64
			err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("ROUND(SUM(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0)), 0)").
				Where("work_order_system_number = ? AND transaction_type_id = 'E'", workOrderId).
				Pluck("ROUND(SUM(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0)), 0)", &totalNonVat).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to calculate total non-VAT",
					Err:        err,
				}
			}

			var taxFree int
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("ISNULL(tax_id, 0)").
				Where("work_order_system_number = ?", workOrderId).
				Pluck("ISNULL(tax_id, 0)", &taxFree).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to get tax free status",
					Err:        err,
				}
			}

			// Calculate totalVat based on your logic
			var totalVat float64
			if taxFree == 0 {
				err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Select("FLOOR((@Total_After_Disc - @Total_Non_Vat) * ISNULL(vat_tax_rate, 0) / 100)").
					Where("work_order_system_number = ?", workOrderId).
					Pluck("FLOOR((@Total_After_Disc - @Total_Non_Vat) * ISNULL(vat_tax_rate, 0) / 100)", &totalVat).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to calculate total VAT",
						Err:        err,
					}
				}
			} else {
				totalVat = 0
			}

			// Calculate TotalAfterVat
			totals.TotalAfterVat = math.Round(totals.TotalAfterDisc + totalVat)

			// Update totals in the WorkOrder table
			err = tx.Model(&transactionworkshopentities.WorkOrder{}).
				Where("work_order_system_number = ?", workOrderId).
				Updates(map[string]interface{}{
					"total_operation":           totals.TotalOpr,
					"total_part":                totals.TotalPart,
					"total_oil":                 totals.TotalOil,
					"total_material":            totals.TotalMaterial,
					"total_consumable_material": totals.TotalConsumableMaterial,
					"total_sublet":              totals.TotalSublet,
					"total_price_accessories":   totals.TotalAccs,
					"total_price":               totals.Total,
					"total_discount":            totals.TotalDisc,
					"total_after_discount":      totals.TotalAfterDisc,
					"total_vat":                 totalVat,
					"total_after_vat":           totals.TotalAfterVat,
					"total_after_vat_non_vat":   totalNonVat,
				}).Error

			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update work order totals",
					Err:        err,
				}
			}

		}
	}

	return true, nil
}
