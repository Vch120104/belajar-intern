package transactionworkshoprepositoryimpl

////  NOTES  ////
import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masteroperationentities "after-sales/api/entities/master/operation"
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
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
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (r *WorkOrderRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tableStruct := transactionworkshoppayloads.WorkOrderGetAllRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilterSearch(joinTable, filterCondition)

	// Add the additional where condition
	whereQuery = whereQuery.Where("service_request_system_number = 0 AND booking_system_number = 0 AND estimation_system_number = 0 and cpc_code = '00002'")

	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Failed to retrieve work order data from the database",
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
				Err:        errVehicle,
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
		InsuranceCheck:           request.WorkOrderInsuranceCheck,
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
	// Update related statuses if necessary

	if err := r.UpdateStatusBookEstim(tx, request); err != nil {
		return transactionworkshopentities.WorkOrder{}, err
	}

	return entitieswo, nil
}

func (r *WorkOrderRepositoryImpl) UpdateStatusBookEstim(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderNormalRequest) *exceptions.BaseErrorResponse {
	var (
		batchSystemNo       int
		bookingStatusClosed = 8
		bookingSystemNo     = request.BookingSystemNumber
		estimationSystemNo  = request.EstimationSystemNumber
	)

	// Update booking status if necessary
	if bookingSystemNo != 0 {
		if batchSystemNo == 0 {
			var batchSystemNoResult struct {
				BatchSystemNo int
			}
			if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
				Select("batch_system_number").
				Where("booking_system_number = ?", bookingSystemNo).
				Scan(&batchSystemNoResult).Error; err != nil {
				return &exceptions.BaseErrorResponse{
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
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update booking status",
				Err:        err,
			}
		}
	}

	// Update estimation status if necessary
	if estimationSystemNo != 0 {
		if batchSystemNo == 0 {
			var batchSystemNoResult struct {
				BatchSystemNo int
			}
			if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
				Select("batch_system_number").
				Where("estimation_system_number = ?", estimationSystemNo).
				Scan(&batchSystemNoResult).Error; err != nil {
				return &exceptions.BaseErrorResponse{
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
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update estimation status",
				Err:        err,
			}
		}
	}

	// Update batch status if necessary
	if batchSystemNo != 0 {
		if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
			Where("batch_system_number = ?", batchSystemNo).
			Update("batch_status_id", bookingStatusClosed).Error; err != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update batch status",
				Err:        err,
			}
		}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) GetById(tx *gorm.DB, Id int, pagination pagination.Pagination) (transactionworkshoppayloads.WorkOrderResponseDetail, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
				Err:        err,
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
	vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100000000&vehicle_id=" + strconv.Itoa(entity.VehicleId)
	var vehicleResponses []transactionworkshoppayloads.VehicleResponse
	errVehicle := utils.GetArray(vehicleUrl, &vehicleResponses, nil)
	if errVehicle != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	// Fetch workorder details without pagination to get total count
	var totalRows int64
	errCount := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ?", Id).
		Count(&totalRows).Error
	if errCount != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count work order details",
			Err:        errCount,
		}
	}

	// Fetch workorder details with pagination
	var workorderDetails []transactionworkshoppayloads.WorkOrderDetailResponse
	errWorkOrderDetails := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("trx_work_order_detail.work_order_detail_id, trx_work_order_detail.work_order_system_number, trx_work_order_detail.line_type_id, lt.line_type_code, trx_work_order_detail.transaction_type_id, tt.transaction_type_code AS transaction_type_code, trx_work_order_detail.job_type_id, tc.job_type_code AS job_type_code, trx_work_order_detail.warehouse_group_id, trx_work_order_detail.frt_quantity, trx_work_order_detail.supply_quantity, trx_work_order_detail.operation_item_price, trx_work_order_detail.operation_item_discount_amount, trx_work_order_detail.operation_item_discount_request_amount").
		Joins("INNER JOIN mtr_work_order_line_type AS lt ON lt.line_type_code = trx_work_order_detail.line_type_id").
		Joins("INNER JOIN mtr_work_order_transaction_type AS tt ON tt.transaction_type_id = trx_work_order_detail.transaction_type_id").
		Joins("INNER JOIN mtr_work_order_job_type AS tc ON tc.job_type_id = trx_work_order_detail.job_type_id").
		Where("work_order_system_number = ?", Id).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderDetails).Error
	if errWorkOrderDetails != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order details from the database",
			Err:        errWorkOrderDetails,
		}
	}

	// Calculate total pages
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))

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
	if err := tx.Model(&transactionworkshopentities.BookingEstimationAllocation{}).
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

	// Fetch work order estimations
	var workorderEstimations []transactionworkshoppayloads.WorkOrderEstimationsResponse
	if err := tx.Model(&transactionworkshopentities.BookingEstimationServiceDiscount{}).
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

	// Fetch work order contracts
	var workorderContracts []transactionworkshoppayloads.WorkOrderContractsResponse
	if err := tx.Model(&transactionworkshopentities.ContractService{}).
		Select("contract_service_system_number, contract_service_document_number, contract_service_date, company_id").
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

// uspg_wtWorkOrder0_Update
// IF @Option = 0
func (r *WorkOrderRepositoryImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderNormalSaveRequest, workOrderId int) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", workOrderId).
		First(&entity).Error
	if err != nil {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order from the database",
			Err:        err,
		}
	}

	// Mapping request fields to entity fields
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

	// Contact person details
	entity.CPTitlePrefix = request.Titleprefix
	entity.ContactPersonName = request.NameCust
	entity.ContactPersonPhone = request.PhoneCust
	entity.ContactPersonMobile = request.MobileCust
	entity.ContactPersonMobileAlternative = request.MobileCustAlternative
	entity.ContactPersonMobileDriver = request.MobileCustDriver
	entity.ContactPersonContactVia = request.ContactVia

	// Insurance details
	entity.InsuranceCheck = request.WorkOrderInsuranceCheck
	entity.InsurancePolicyNumber = request.WorkOrderInsurancePolicyNo
	entity.InsuranceExpiredDate = request.WorkOrderInsuranceExpiredDate
	entity.InsuranceClaimNumber = request.WorkOrderInsuranceClaimNo
	entity.InsurancePersonInCharge = request.WorkOrderInsurancePic
	entity.InsuranceOwnRisk = request.WorkOrderInsuranceOwnRisk
	entity.InsuranceWorkOrderNumber = request.WorkOrderInsuranceWONumber

	// Other work order details (Page 2 fields)
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

	// Handling VAT Tax Rate as per SQL logic
	if isFTZCompany(request.CompanyId) {

		entity.VATTaxRate = &[]float64{0}[0]
	} else {

		// Call getTaxPercent with the correct arguments and handle both return values
		vatTaxRate, err := getTaxPercent(tx, 10, 11, time.Now()) // 10.PPN 11.PPN
		if err != nil {
			return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate VAT tax rate",
				Err:        err,
			}
		}
		entity.VATTaxRate = &vatTaxRate
	}

	err = tx.Save(&entity).Error
	if err != nil {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save the updated work order",
			Err:        err,
		}
	}

	return entity, nil
}

func isFTZCompany(companyId int) bool {
	return companyId == 139 //1520098 - Rodamas Makmur Motor
}

func getTaxPercent(tx *gorm.DB, taxTypeId int, taxServCode int, effDate time.Time) (float64, error) {
	var taxPercent float64

	// Subquery for the effective date (EFF_DATE)
	var effDateSubquery time.Time
	subquery := tx.Table("dms_microservices_finance_dev.dbo.mtr_tax_fare").
		Select("mtr_tax_fare.effective_date").
		Where("mtr_tax_fare.tax_type_id = ? AND mtr_tax_fare.effective_date <= ?", taxTypeId, effDate).
		Order("mtr_tax_fare.effective_date DESC").
		Limit(1).
		Find(&effDateSubquery)
	if subquery.Error != nil {
		return 0, subquery.Error
	}

	// Main query to get the tax percent
	err := tx.Table("dms_microservices_finance_dev.dbo.mtr_tax_fare").
		Select("CASE WHEN mtr_tax_fare_detail.is_use_net = 0 THEN mtr_tax_fare_detail.tax_percent ELSE (mtr_tax_fare_detail.tax_percent * (COALESCE(mtr_tax_fare_detail.net_percent, 0) / 100)) END").
		Joins("LEFT JOIN dms_microservices_finance_dev.dbo.mtr_tax_fare_detail ON mtr_tax_fare.tax_fare_id = mtr_tax_fare_detail.tax_fare_id").
		Where("mtr_tax_fare.tax_type_id = ? AND mtr_tax_fare_detail.tax_service_id = ? AND mtr_tax_fare.effective_date = ?", taxTypeId, taxServCode, effDateSubquery).
		Scan(&taxPercent).Error
	if err != nil {
		return 0, err
	}

	return taxPercent, nil
}

func (r *WorkOrderRepositoryImpl) Void(tx *gorm.DB, workOrderId int) (bool, *exceptions.BaseErrorResponse) {
	// Check if the work order exists
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", workOrderId).
		First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
				Err:        err,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order from the database",
			Err:        err,
		}
	}

	// Delete the work order
	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete the work order",
			Err:        err,
		}
	}

	return true, nil
}

// uspg_wtWorkOrder0_Update
// IF @Option = 2
func (r *WorkOrderRepositoryImpl) CloseOrder(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order not found",
				Err:        err,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order from the database",
			Err:        err,
		}
	}

	// Check if WorkOrderStatusId is equal to 1 (Draft)
	if entity.WorkOrderStatusId == utils.WoStatDraft {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Work order cannot be closed because status is draft",
			Err:        err,
		}
	}

	// Check if there is still DP payment that has not been settled
	var dpPaymentAllocated float64
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", Id).
		Select("COALESCE(downpayment_payment_allocated, 0) as downpayment_payment_allocated").
		Scan(&dpPaymentAllocated).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve DP payment allocated from the database",
			Err:        err,
		}
	}
	if dpPaymentAllocated > 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "There is still DP payment that has not been settled",
			Err:        err,
		}
	}

	// Check if there are any work order items without invoices
	var count int64 //cek statusid <> 8(closed), billcode <> no_charge (5), substituteid
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id <> ? AND substitute_type_id <> ?",
			Id, utils.WoStatClosed, 5, 0).
		Count(&count).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order items from the database",
			Err:        err,
		}
	}
	if count > 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Detail Work Order without Invoice No must be deleted",
			Err:        err,
		}
	}

	// Check for warranty items
	var allPtpSupply bool //cek statusid <> 8(closed), billcode <> warranty (6), substituteid
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id = ? AND substitute_type_id <> ?",
			Id, utils.WoStatClosed, 6, 0).
		Count(&count).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve warranty items from the database",
			Err:        err,
		}
	}
	if count == 0 {
		allPtpSupply = true
	} else {
		// Validate part-to-part supply //cek statusid <> 8(closed), billcode <> warranty (6), substituteid , warrantyclaim_type = 0 (part), frt_qty > supply_qty
		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id = ? AND substitute_type_id <> ? AND warranty_claim_type_id = ? AND frt_qty > supply_qty",
				Id, utils.WoStatClosed, 6, 0, 0).
			Count(&count).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to validate part-to-part supply",
				Err:        err,
			}
		}
		if count > 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Warranty Item (PTP) must be supplied",
				Err:        err,
			}
		}

		// Validate part-to-money and operation status //cek statusid <> 8(closed), billcode <> warranty (6), substituteid , warrantyclaim_type = 0 (part)
		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id = ? AND substitute_type_id <> ? AND warranty_claim_type_id <> ?",
				Id, utils.WoStatClosed, 6, 0, 0).
			Count(&count).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to validate part-to-money and operation status",
				Err:        err,
			}
		}
		if count > 0 {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Warranty Item (PTM)/Operation must be Invoiced",
				Err:        err,
			}
		}

		allPtpSupply = true
	}

	// Check if all items/operations/packages other than warranty are closed
	err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_status_id <> ? AND substitute_type_id <> ? AND transaction_type_id NOT IN (?, ?)",
			Id, utils.WoStatClosed, 0, 6, 5).
		Count(&count).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check if all items/operations/packages are closed",
			Err:        err,
		}
	}
	if allPtpSupply && count > 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "There is Work Order detail that has not been Invoiced",
			Err:        err,
		}
	}

	// Validate mileage and update vehicle master if necessary
	var servMileage, lastKm int
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", Id).
		Select("service_mileage").Scan(&servMileage).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve service mileage",
			Err:        err,
		}
	}
	err = tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle").
		Where("vehicle_chassis_number = ?", entity.VehicleChassisNumber).
		Select("vehicle_last_km").
		Scan(&lastKm).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve last mileage",
			Err:        err,
		}
	}
	if servMileage <= lastKm {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Service Mileage must be larger than Last Mileage.",
			Err:        err,
		}
	}

	// Update vehicle master
	err = tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle").Where("vehicle_chassis_number = ?", entity.VehicleChassisNumber).
		Updates(map[string]interface{}{
			"vehicle_last_km":           servMileage,
			"vehicle_last_service_date": entity.WorkOrderDate,
		}).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update vehicle master",
			Err:        err,
		}
	}

	// If Work Order still has DP Payment not allocated for Invoice
	type DPPaymentDetails struct {
		DPPayment    float64 `gorm:"column:downpayment_payment"`
		DPAllocToInv float64 `gorm:"column:downpayment_payment_to_invoice"`
	}

	var details DPPaymentDetails
	var dpOverpay float64

	err = tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", Id).
		Select("downpayment_payment, downpayment_payment_to_invoice").
		Scan(&details).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve DP payment details",
			Err:        err,
		}
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update DP payment details",
			Err:        err,
		}
	}

	// Determine customer type and set event number
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to close the work order",
			Err:        err,
		}
	}

	return true, nil
}

// uspg_wtWorkOrder1_Insert
// IF @Option = 0
// --USE FOR : * INSERT NEW DATA
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (r *WorkOrderRepositoryImpl) GetAllRequest(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrderService

	query := tx.Model(&transactionworkshopentities.WorkOrderService{})

	if len(filterCondition) > 0 {
		query = utils.ApplyFilterSearch(query, filterCondition)
	}

	query = query.Scopes(pagination.Paginate(&entities, &pages, tx))

	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order service requests",
			Err:        err,
		}
	}

	var workOrderServiceResponses []map[string]interface{}
	for _, entity := range entities {
		workOrderServiceData := map[string]interface{}{
			"work_order_service_id":     entity.WorkOrderServiceId,
			"work_order_system_number":  entity.WorkOrderSystemNumber,
			"work_order_service_remark": entity.WorkOrderServiceRemark,
		}
		workOrderServiceResponses = append(workOrderServiceResponses, workOrderServiceData)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(workOrderServiceResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) GetRequestById(tx *gorm.DB, workorderID int, detailID int) (transactionworkshoppayloads.WorkOrderServiceResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", workorderID, detailID).
		First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderServiceResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order service request from the database",
			Err:        err,
		}
	}

	payload := transactionworkshoppayloads.WorkOrderServiceResponse{
		WorkOrderServiceId:     entity.WorkOrderServiceId,
		WorkOrderSystemNumber:  entity.WorkOrderSystemNumber,
		WorkOrderServiceRemark: entity.WorkOrderServiceRemark,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) UpdateRequest(tx *gorm.DB, workorderID int, detailID int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderService, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", workorderID, detailID).
		First(&entity).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderService{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "failed to retrieve work order service request from the database",
			Err:        err,
		}
	}

	entity.WorkOrderServiceRemark = request.WorkOrderServiceRemark

	err = tx.Model(&entity).Updates(map[string]interface{}{
		"work_order_service_remark": entity.WorkOrderServiceRemark,
	}).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderService{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to update the work order service request",
			Err:        err,
		}
	}

	return entity, nil
}

func (r *WorkOrderRepositoryImpl) AddRequest(tx *gorm.DB, workorderID int, request transactionworkshoppayloads.WorkOrderServiceRequest) (transactionworkshopentities.WorkOrderService, *exceptions.BaseErrorResponse) {

	var lastService transactionworkshopentities.WorkOrderService
	// Find the latest service request for the same WorkOrderSystemNumber to get the latest line number
	err := tx.Where("work_order_system_number = ?", request.WorkOrderSystemNumber).
		Order("work_order_service_request_line DESC").
		First(&lastService).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return transactionworkshopentities.WorkOrderService{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get the last work order service request",
			Err:        err,
		}
	}

	// Increment the line number, starting from 1 if no previous records are found
	nextLine := lastService.WorkOrderServiceRequestLine + 1

	// Create new WorkOrderService entity with the incremented line number
	entities := transactionworkshopentities.WorkOrderService{
		WorkOrderSystemNumber:       request.WorkOrderSystemNumber,
		WorkOrderServiceRemark:      request.WorkOrderServiceRemark,
		WorkOrderServiceDate:        time.Now(),
		WorkOrderServiceRequestLine: nextLine,
	}

	// Save the new work order service request
	err = tx.Create(&entities).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderService{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save the work order service request",
			Err:        err,
		}
	}

	return entities, nil
}

func (r *WorkOrderRepositoryImpl) AddRequestMultiId(tx *gorm.DB, workorderID int, requests []transactionworkshoppayloads.WorkOrderServiceRequest) ([]transactionworkshopentities.WorkOrderService, *exceptions.BaseErrorResponse) {
	// Limit the number of requests to 5
	if len(requests) > 5 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "A maximum of 5 requests can be added at once",
		}
	}

	var entities []transactionworkshopentities.WorkOrderService

	for _, request := range requests {
		var lastService transactionworkshopentities.WorkOrderService

		// Find the latest service request for the same WorkOrderSystemNumber to get the latest line number
		err := tx.Where("work_order_system_number = ?", request.WorkOrderSystemNumber).
			Order("work_order_service_request_line DESC").
			First(&lastService).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get the last work order service request",
				Err:        err,
			}
		}

		// Increment the line number, starting from 1 if no previous records are found
		nextLine := lastService.WorkOrderServiceRequestLine + 1

		// Create new WorkOrderService entity with the incremented line number
		entity := transactionworkshopentities.WorkOrderService{
			WorkOrderSystemNumber:       request.WorkOrderSystemNumber,
			WorkOrderServiceRemark:      request.WorkOrderServiceRemark,
			WorkOrderServiceDate:        time.Now(),
			WorkOrderServiceRequestLine: nextLine,
		}

		// Save the new work order service request
		err = tx.Create(&entity).Error
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to save the work order service request",
				Err:        err,
			}
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *WorkOrderRepositoryImpl) DeleteRequest(tx *gorm.DB, workorderID int, detailID int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", workorderID, detailID).
		Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete work order service request from the database",
			Err:        err,
		}
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
		query = utils.ApplyFilterSearch(query, filterCondition)
	}

	query = query.Scopes(pagination.Paginate(&entities, &pages, tx))

	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order service vehicle requests",
			Err:        err,
		}
	}

	var workOrderServiceVehicleResponses []map[string]interface{}
	for _, entity := range entities {
		workOrderServiceVehicleData := map[string]interface{}{
			"work_order_service_vehicle_id": entity.WorkOrderServiceVehicleId,
			"work_order_system_number":      entity.WorkOrderSystemNumber,
			"work_order_vehicle_date":       entity.WorkOrderVehicleDate,
			"work_order_vehicle_remark":     entity.WorkOrderVehicleRemark,
		}
		workOrderServiceVehicleResponses = append(workOrderServiceVehicleResponses, workOrderServiceVehicleData)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(workOrderServiceVehicleResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) GetVehicleServiceById(tx *gorm.DB, workorderID int, detailID int) (transactionworkshoppayloads.WorkOrderServiceVehicleResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderServiceVehicle
	err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ? AND work_order_service_vehicle_id = ?", workorderID, detailID).
		First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderServiceVehicleResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order service vehicle request from the database",
			Err:        err,
		}
	}

	payload := transactionworkshoppayloads.WorkOrderServiceVehicleResponse{
		WorkOrderServiceVehicleId: entity.WorkOrderServiceVehicleId,
		WorkOrderSystemNumber:     entity.WorkOrderSystemNumber,
		WorkOrderVehicleDate:      entity.WorkOrderVehicleDate,
		WorkOrderVehicleRemark:    entity.WorkOrderVehicleRemark,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) UpdateVehicleService(tx *gorm.DB, workorderID int, detailID int, request transactionworkshoppayloads.WorkOrderServiceVehicleRequest) (transactionworkshopentities.WorkOrderServiceVehicle, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderServiceVehicle
	err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ? AND work_order_service_vehicle_id = ?", workorderID, detailID).
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

func (r *WorkOrderRepositoryImpl) DeleteVehicleService(tx *gorm.DB, workorderID int, detailID int) (bool, *exceptions.BaseErrorResponse) {

	var entity transactionworkshopentities.WorkOrderServiceVehicle
	err := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{}).
		Where("work_order_system_number = ? AND work_order_service_id = ?", workorderID, detailID).
		Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete work order service request from the database",
			Err:        err,
		}
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
			&workOrderReq.WarehouseGroupId,
			&workOrderReq.OperationItemId,
			&workOrderReq.OperationItemCode,
			&workOrderReq.OperationItemPrice,
			&workOrderReq.OperationItemDiscountAmount,
			&workOrderReq.OperationItemDiscountRequestAmount,
			&workOrderReq.OperationItemDiscountPercent,
			&workOrderReq.OperationItemDiscountRequestPercent,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// fetch line type from external api
		lineTypeUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-line-type?line_type_code=" + strconv.Itoa(workOrderReq.LineTypeId)
		var lineTypeResponse []transactionworkshoppayloads.Linetype
		err := utils.GetArray(lineTypeUrl, &lineTypeResponse, nil)
		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve line type from the external API",
				Err:        err,
			}
		}

		// fetch transaction type from external api
		transactionTypeUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-transaction-type?transaction_type_id=" + strconv.Itoa(workOrderReq.TransactionTypeId)
		var transactionTypeResponse []transactionworkshoppayloads.WorkOrderTransactionType
		err = utils.GetArray(transactionTypeUrl, &transactionTypeResponse, nil)
		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve transaction type from the external API",
				Err:        err,
			}
		}

		// fetch job type from external api
		jobTypeUrl := config.EnvConfigs.AfterSalesServiceUrl + "work-order/dropdown-job-type?job_type_id=" + strconv.Itoa(workOrderReq.JobTypeId)
		var jobTypeResponse []transactionworkshoppayloads.WorkOrderJobType
		err = utils.GetArray(jobTypeUrl, &jobTypeResponse, nil)
		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve job type from the external API",
				Err:        err,
			}
		}

		workOrderRes = transactionworkshoppayloads.WorkOrderDetailResponse{
			WorkOrderDetailId:                   workOrderReq.WorkOrderDetailId,
			WorkOrderSystemNumber:               workOrderReq.WorkOrderSystemNumber,
			LineTypeId:                          workOrderReq.LineTypeId,
			LineTypeCode:                        lineTypeResponse[0].LineTypeCode,
			TransactionTypeId:                   workOrderReq.TransactionTypeId,
			TransactionTypeCode:                 transactionTypeResponse[0].TransactionTypeCode,
			JobTypeId:                           workOrderReq.JobTypeId,
			JobTypeCode:                         jobTypeResponse[0].JobTypeCode,
			OperationItemId:                     workOrderReq.OperationItemId,
			FrtQuantity:                         workOrderReq.FrtQuantity,
			SupplyQuantity:                      workOrderReq.SupplyQuantity,
			OperationItemPrice:                  workOrderReq.OperationItemPrice,
			OperationItemDiscountAmount:         workOrderReq.OperationItemDiscountAmount,
			OperationItemDiscountRequestAmount:  workOrderReq.OperationItemDiscountRequestAmount,
			OperationItemDiscountPercent:        workOrderReq.OperationItemDiscountPercent,
			OperationItemDiscountRequestPercent: workOrderReq.OperationItemDiscountRequestPercent,
		}

		convertedResponses = append(convertedResponses, workOrderRes)
	}

	var mapResponses []map[string]interface{}

	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_detail_id":                   response.WorkOrderDetailId,
			"work_order_system_number":               response.WorkOrderSystemNumber,
			"line_type_id":                           response.LineTypeId,
			"line_type_code":                         response.LineTypeCode,
			"transaction_type_id":                    response.TransactionTypeId,
			"transaction_type_code":                  response.TransactionTypeCode,
			"job_type_id":                            response.JobTypeId,
			"job_type_code":                          response.JobTypeCode,
			"frt_quantity":                           response.FrtQuantity,
			"supply_quantity":                        response.SupplyQuantity,
			"operation_item_id":                      response.OperationItemId,
			"operation_item_price":                   response.OperationItemPrice,
			"operation_item_discount_amount":         response.OperationItemDiscountAmount,
			"operation_item_discount_request_amount": response.OperationItemDiscountRequestAmount,
			"operation_item_discount_percent":        response.OperationItemDiscountPercent,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil

}

func (r *WorkOrderRepositoryImpl) GetDetailByIdWorkOrder(tx *gorm.DB, workorderID int, detailID int) (transactionworkshoppayloads.WorkOrderDetailResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderDetail
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_detail_id = ?", workorderID, detailID).
		First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Work order detail not found",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order detail from the database",
			Err:        err}
	}

	payload := transactionworkshoppayloads.WorkOrderDetailResponse{
		WorkOrderDetailId:     entity.WorkOrderDetailId,
		WorkOrderSystemNumber: entity.WorkOrderSystemNumber,
		LineTypeId:            entity.LineTypeId,
		TransactionTypeId:     entity.TransactionTypeId,
		JobTypeId:             entity.JobTypeId,
		FrtQuantity:           entity.FrtQuantity,
		SupplyQuantity:        entity.SupplyQuantity,
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

		if estimSystemNo != 0 {
			var maxWoOprItemLine int
			if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("COALESCE(MAX(work_order_operation_item_line), 0)").
				Where("work_order_system_number = ?", id).
				Scan(&maxWoOprItemLine).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve maximum work order operation item line",
					Err:        err,
				}
			}

			// var bookingEstim21 []transactionworkshopentities.BookingEstimation
			// err := tx.Model(&bookingEstim21).
			// 	Select("BE.ESTIM_LINE, BE.LINE_TYPE, BE.OPR_ITEM_CODE, BE.DESCRIPTION, I.SELLING_UOM, BE.FRT_QTY, BE.OPR_ITEM_PRICE, BE.OPR_ITEM_DISC_AMOUNT, BE.OPR_ITEM_DISC_REQ_AMOUNT, BE.OPR_ITEM_DISC_PERCENT, BE.OPR_ITEM_DISC_REQ_PERCENT, BE.PPH_AMOUNT, BE.PPH_TAX_CODE, BE.PPH_TAX_RATE").
			// 	Joins("LEFT OUTER JOIN wtBookEstim0 BE0 ON BE0.ESTIM_SYSTEM_NO = BE.ESTIM_SYSTEM_NO").
			// 	Joins("LEFT OUTER JOIN gmItem0 I ON I.ITEM_CODE = BE.OPR_ITEM_CODE").
			// 	Where("BE.ESTIM_SYSTEM_NO = ?", estimSystemNo).
			// 	Find(&bookingEstim21).Error

			// if err != nil {
			// 	return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
			// 		StatusCode: http.StatusInternalServerError,
			// 		Message:    "Failed to retrieve booking estimation data",
			// 		Err:        err,
			// 	}
			// }

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:               id,
				LineTypeId:                          0,  // BE0.LineTypeId,
				TransactionTypeId:                   0,  // utils.TrxTypeWoExternal,
				JobTypeId:                           0,  // CASE WHEN BE0.CPC_CODE = @Profit_Center_BR THEN @JobTypeBR ELSE @JobTypePM END,
				OperationItemCode:                   "", //BE.OPR_ITEM_CODE,
				WarehouseGroupId:                    0,  //Whs_Group_Sp
				FrtQuantity:                         0,  //BE.FrtQuantity,
				SupplyQuantity:                      0,  //CASE WHEN BE.LINE_TYPE = @LINETYPE_OPR OR BE.LINE_TYPE = @LINETYPE_PACKAGE THEN BE.FRT_QTY ELSE CASE WHEN I.ITEM_TYPE = @ItemTypeService AND I.ITEM_GROUP <> @ItemGrpOJ THEN BE.FRT_QTY ELSE 0 END END
				WorkorderStatusId:                   utils.WoStatDraft,
				OperationItemDiscountAmount:         0,                    //BE.OPR_ITEM_DISC_AMOUNT,
				OperationItemDiscountRequestAmount:  0,                    //BE.OPR_ITEM_DISC_REQ_AMOUNT,
				OperationItemDiscountPercent:        0,                    //BE.OPR_ITEM_DISC_PERCENT,
				OperationItemDiscountRequestPercent: 0,                    //BE.OPR_ITEM_DISC_REQ_PERCENT,
				OperationItemPrice:                  0,                    //BE.OPR_ITEM_PRICE,
				PphAmount:                           0,                    //BE.PPH_AMOUNT,
				PphTaxRate:                          0,                    //BE.PPH_TAX_RATE,
				AtpmWCFTypeId:                       0,                    //CASE WHEN BE.LINE_TYPE = @LINETYPE_OPR OR BE.LINE_TYPE = @LINETYPE_PACKAGE THEN '' ELSE ATPM_WCF_TYPE END
				WorkOrderOperationItemLine:          maxWoOprItemLine + 1, //BE.ESTIM_LINE,
			}

			if request.LineTypeId == 1 {
				workOrderDetail.OperationItemId = request.OperationItemId
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

		} else {

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

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:               id,
				LineTypeId:                          request.LineTypeId,        // BE0.LineTypeId,
				TransactionTypeId:                   request.TransactionTypeId, // utils.TrxTypeWoExternal,
				JobTypeId:                           request.JobTypeId,         // CASE WHEN BE0.CPC_CODE = @Profit_Center_BR THEN @JobTypeBR ELSE @JobTypePM END,
				OperationItemCode:                   request.OperationItemCode, // BE.OPR_ITEM_CODE,
				WarehouseGroupId:                    request.WarehouseGroupId,  // Whs_Group_Sp
				FrtQuantity:                         request.FrtQuantity,       // BE.FrtQuantity,
				SupplyQuantity:                      request.SupplyQuantity,    // CASE WHEN BE.LINE_TYPE = @LINETYPE_OPR OR BE.LINE_TYPE = @LINETYPE_PACKAGE THEN BE.FRT_QTY ELSE CASE WHEN I.ITEM_TYPE = @ItemTypeService AND I.ITEM_GROUP <> @ItemGrpOJ THEN BE.FRT_QTY ELSE 0 END END
				WorkorderStatusId:                   utils.WoStatDraft,
				OperationItemDiscountAmount:         0,                          // BE.OPR_ITEM_DISC_AMOUNT,
				OperationItemDiscountRequestAmount:  0,                          // BE.OPR_ITEM_DISC_REQ_AMOUNT,
				OperationItemDiscountPercent:        0,                          // BE.OPR_ITEM_DISC_PERCENT,
				OperationItemDiscountRequestPercent: 0,                          // BE.OPR_ITEM_DISC_REQ_PERCENT,
				OperationItemPrice:                  request.OperationItemPrice, // BE.OPR_ITEM_PRICE,
				PphAmount:                           0,                          // BE.PPH_AMOUNT,
				PphTaxRate:                          0,                          // BE.PPH_TAX_RATE,
				AtpmWCFTypeId:                       0,                          // CASE WHEN BE.LINE_TYPE = @LINETYPE_OPR OR BE.LINE_TYPE = @LINETYPE_PACKAGE THEN '' ELSE ATPM_WCF_TYPE END
				WorkOrderOperationItemLine:          maxWoOprItemLine + 1,
			}

			if request.LineTypeId == 1 {
				workOrderDetail.OperationItemId = request.OperationItemId
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

			var campaignItems []masterentities.CampaignMasterDetail
			err := tx.Model(&campaignItems).
				Select("line_type_id, item_operation_id, quantity, price, discount_percent").
				Joins("INNER JOIN mtr_campaign C ON campaign_id = C.campaign_id").
				Joins("LEFT JOIN mtr_item I ON I.item_id = C1.item_operation_id").
				Where("campaign_id = ?", campaignId).
				Find(&campaignItems).Error

			if err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve campaign items",
					Err:        err,
				}
			}

			if len(campaignItems) > 0 {
				workOrderDetail = transactionworkshopentities.WorkOrderDetail{
					WorkOrderSystemNumber:               id,
					LineTypeId:                          campaignItems[0].LineTypeId,                    // C1.LINE_TYPE,
					TransactionTypeId:                   3,                                              // utils.TrxTypeWoExternal,
					JobTypeId:                           2,                                              // JobTypeCampaign,
					OperationItemCode:                   strconv.Itoa(campaignItems[0].ItemOperationId), // C1.OPR_ITEM_CODE,
					WarehouseGroupId:                    1,                                              // Whs_Group_Campaign
					FrtQuantity:                         campaignItems[0].Quantity,                      // C1.FRT_QTY,
					SupplyQuantity:                      0,                                              // CASE WHEN C1.LINE_TYPE = @LINETYPE_OPR THEN C1.FRT_QTY ELSE CASE WHEN I.ITEM_TYPE = @ItemTypeService AND I.ITEM_GROUP <> @ItemGrpOJ THEN C1.FRT_QTY ELSE 0 END END
					WorkorderStatusId:                   utils.WoStatDraft,
					OperationItemDiscountAmount:         math.Round(campaignItems[0].Price * campaignItems[0].DiscountPercent / 100), // ROUND((C1.OPR_ITEM_PRICE * C1.OPR_ITEM_DISC_PERCENT /100),0,0),
					OperationItemDiscountRequestAmount:  0,                                                                           // 0,
					OperationItemDiscountPercent:        campaignItems[0].DiscountPercent,                                            // C1.OPR_ITEM_DISC_PERCENT,
					OperationItemDiscountRequestPercent: 0,                                                                           // 0,
					OperationItemPrice:                  campaignItems[0].Price,                                                      // C1.OPR_ITEM_PRICE,
					PphAmount:                           0,                                                                           // 0,
					PphTaxRate:                          0,                                                                           // CASE WHEN C1.LINE_TYPE = @LINETYPE_OPR THEN OPR.TAX_CODE	ELSE ''	END,
					AtpmWCFTypeId:                       0,
					WorkOrderOperationItemLine:          maxWoOprItemLine + 1, // 0
				}
			} else {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "No campaign items found",
					Err:        errors.New("campaign items not found"),
				}
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

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:               id,
				LineTypeId:                          utils.LinetypeOperation, // LINETYPE_OPR,
				TransactionTypeId:                   0,                       // dbo.FCT_getBillCode(@COMPANY_CODE ,CAST(P1.COMPANY_CODE AS VARCHAR(10)),'W'),
				JobTypeId:                           8,                       // dbo.getVariableValue('JOBTYPE_PDI'),
				OperationItemCode:                   "",                      // P1.OPERATION_NO,
				WarehouseGroupId:                    38,                      // Whs_Group_Sp
				FrtQuantity:                         0,                       // P1.FRT,
				SupplyQuantity:                      0,                       // P1.FRT
				WorkorderStatusId:                   0,                       // ""
				OperationItemDiscountAmount:         0,                       // 0,
				OperationItemDiscountRequestAmount:  0,                       // 0,
				OperationItemDiscountPercent:        0,                       // 0,
				OperationItemDiscountRequestPercent: 0,                       // 0,
				OperationItemPrice:                  0,                       // LSP1.SELLING_PRICE,
				PphAmount:                           0,                       // 0,
				PphTaxRate:                          0,                       // OP.TAX_CODE,
				AtpmWCFTypeId:                       0,                       // 0
				WorkOrderOperationItemLine:          maxWoOprItemLine + 1,
			}

			if request.LineTypeId == 1 {
				workOrderDetail.OperationItemId = request.OperationItemId
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
				WorkOrderSystemNumber:               id,
				LineTypeId:                          0,  // SR1.LINE_TYPE,
				TransactionTypeId:                   0,  // dbo.FCT_getBillCode(@COMPANY_CODE ,SR.COMPANY_CODE,'W'),
				JobTypeId:                           0,  // CASE WHEN SR.SERV_PROFIT_CENTER = @Profit_Center_BR	THEN @JobTypeBR	ELSE @JobTypePM	END,
				OperationItemCode:                   "", // SR1.OPR_ITEM_CODE,
				WarehouseGroupId:                    38, // Whs_Group_Sp
				FrtQuantity:                         0,  // SR1.FRT_QTY,
				SupplyQuantity:                      0,  // //CASE WHEN SR1.LINE_TYPE = @LINETYPE_OPR OR SR1.LINE_TYPE = @LINETYPE_PACKAGE THEN SR1.FRT_QTY ELSE CASE WHEN I.ITEM_TYPE = @ItemTypeService AND I.ITEM_GROUP <> @ItemGrpOJ THEN SR1.FRT_QTY ELSE 0 END END
				WorkorderStatusId:                   utils.WoStatDraft,
				OperationItemDiscountAmount:         0, // 0,
				OperationItemDiscountRequestAmount:  0, // 0,
				OperationItemDiscountPercent:        0, // 0,
				OperationItemDiscountRequestPercent: 0, // 0,
				OperationItemPrice:                  0, // dbo.getOprItemPrice(SR1.LINE_TYPE ,	@Whs_Group_Sp, --@Whs_Group, dbo.FCT_getBillCode(@COMPANY_CODE ,SR.COMPANY_CODE,'W') , @COMPANY_CODE,@VEHICLE_BRAND , @JobTypeGR , --==TEMPORARY UNTIL REVISION ON TABEL SERV REQ DONE (CASE WHEN SR1.CPC_CODE = dbo.getVariableValue('PROFIT_CENTER_BR') THEN dbo.getVariableValue('JOBTYPE_BR') ELSE dbo.getVariableValue('JOBTYPE_GR') END ), @Model_Code, SR1.OPR_ITEM_CODE , @CCY_CODE,'','',0,default,@Price_Code),,
				PphAmount:                           0, // 0,
				PphTaxRate:                          0, // CASE WHEN SR1.LINE_TYPE = @LINETYPE_OPR THEN OPR.TAX_CODE WHEN SR1.LINE_TYPE = @LINETYPE_PACKAGE	THEN PCK.PPH_TAX_CODE ELSE '' END,
				AtpmWCFTypeId:                       0, // 0
				WorkOrderOperationItemLine:          maxWoOprItemLine + 1,
			}

			if request.LineTypeId == 1 {
				workOrderDetail.OperationItemId = request.OperationItemId
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

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:               id,
				LineTypeId:                          0,  // RW1.LINE_TYPE,
				TransactionTypeId:                   0,  // dbo.getVariableValue('TRXTYPE_WO_NOCHARGE'),
				JobTypeId:                           0,  // RW1.JOB_TYPE,
				OperationItemCode:                   "", // RW1.OPR_ITEM_CODE,
				WarehouseGroupId:                    0,  // RW1.WHS_GROUP
				FrtQuantity:                         0,  // RW1.FRT_QTY,
				SupplyQuantity:                      0,  // RW1.SUPPLY_QTY
				WorkorderStatusId:                   utils.WoStatDraft,
				OperationItemDiscountAmount:         0, // 0,
				OperationItemDiscountRequestAmount:  0, // 0,
				OperationItemDiscountPercent:        0, // 0,
				OperationItemDiscountRequestPercent: 0, // 0,
				OperationItemPrice:                  0, // RW1.OPR_ITEM_PRICE
				PphAmount:                           0, // 0,
				PphTaxRate:                          0, // 0
				AtpmWCFTypeId:                       0, // 0
				WorkOrderOperationItemLine:          maxWoOprItemLine + 1,
			}

			if request.LineTypeId == 1 {
				workOrderDetail.OperationItemId = request.OperationItemId
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

	// Validate if the work order is still draft

	// Validasi untuk chassis yang sudah pernah PDI,FSI,WR

	// Validate Line Type Item must be inside item master

	// Validate if Warranty to Vehicle Age

	// LINE TYPE <> 1 , NEED SUBSTITUTE

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
	// result := tx.Model(`
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
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order detail from the database",
			Err:        err,
		}
	}

	entity.LineTypeId = request.LineTypeId
	entity.TransactionTypeId = request.TransactionTypeId
	entity.JobTypeId = request.JobTypeId
	entity.WarehouseGroupId = request.WarehouseGroupId
	entity.OperationItemId = request.OperationItemId
	entity.FrtQuantity = request.FrtQuantity
	entity.SupplyQuantity = request.SupplyQuantity
	entity.PriceListId = request.PriceListId
	entity.OperationItemDiscountRequestAmount = request.OperationItemDiscountRequestAmount
	entity.OperationItemPrice = request.OperationItemPrice

	if request.LineTypeId == 1 {
		entity.OperationItemId = request.OperationItemId
	}

	err = tx.Save(&entity).Error
	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save the updated work order detail",
			Err:        err,
		}
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete work order detail from the database",
			Err:        err,
		}
	}

	return true, nil
}

// uspg_wtWorkOrder0_Insert
// IF @Option = 1
// --USE FOR : * INSERT NEW DATA FROM BOOKING AND ESTIMATION
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
	if err := r.UpdateStatusBookEstimNewBooking(tx, request); err != nil {
		return transactionworkshopentities.WorkOrder{}, err
	}

	return entitieswo, nil
}

func (r *WorkOrderRepositoryImpl) UpdateStatusBookEstimNewBooking(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderBookingRequest) *exceptions.BaseErrorResponse {
	var (
		batchSystemNo       int
		bookingStatusClosed = 8
		bookingSystemNo     = request.BookingSystemNumber
		estimationSystemNo  = request.EstimationSystemNumber
	)

	// Update booking status if necessary
	if bookingSystemNo != 0 {
		if batchSystemNo == 0 {
			var batchSystemNoResult struct {
				BatchSystemNo int
			}
			if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
				Select("batch_system_number").
				Where("booking_system_number = ?", bookingSystemNo).
				Scan(&batchSystemNoResult).Error; err != nil {
				return &exceptions.BaseErrorResponse{
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
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update booking status",
				Err:        err,
			}
		}
	}

	// Update estimation status if necessary
	if estimationSystemNo != 0 {
		if batchSystemNo == 0 {
			var batchSystemNoResult struct {
				BatchSystemNo int
			}
			if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
				Select("batch_system_number").
				Where("estimation_system_number = ?", estimationSystemNo).
				Scan(&batchSystemNoResult).Error; err != nil {
				return &exceptions.BaseErrorResponse{
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
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update estimation status",
				Err:        err,
			}
		}
	}

	// Update batch status if necessary
	if batchSystemNo != 0 {
		if err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
			Where("batch_system_number = ?", batchSystemNo).
			Update("batch_status_id", bookingStatusClosed).Error; err != nil {
			return &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update batch status",
				Err:        err,
			}
		}
	}

	return nil
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
				Err:        err,
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
	vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=100000000&vehicle_id=" + strconv.Itoa(entity.VehicleId)
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
		Select("trx_work_order_detail.work_order_detail_id, trx_work_order_detail.work_order_system_number, trx_work_order_detail.line_type_id,lt.line_type_code, trx_work_order_detail.transaction_type_id, tt.transaction_type_code AS transaction_type_code, trx_work_order_detail.job_type_id, tc.job_type_code AS job_type_code, trx_work_order_detail.warehouse_group_id, trx_work_order_detail.frt_quantity, trx_work_order_detail.supply_quantity, trx_work_order_detail.operation_item_price, trx_work_order_detail.operation_item_discount_amount, trx_work_order_detail.operation_item_discount_request_amount").
		Joins("INNER JOIN mtr_work_order_line_type AS lt ON lt.line_type_code = trx_work_order_detail.line_type_id").
		Joins("INNER JOIN mtr_work_order_transaction_type AS tt ON tt.transaction_type_id = trx_work_order_detail.transaction_type_id").
		Joins("INNER JOIN mtr_work_order_job_type AS tc ON tc.job_type_id = trx_work_order_detail.job_type_id").
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

	// Fetch work order campaigns
	var workorderCampaigns []transactionworkshoppayloads.WorkOrderCampaignResponse
	if err := tx.Model(&masterentities.CampaignMaster{}).
		Where("campaign_id = ? and campaign_id != 0", entity.CampaignId).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderCampaigns).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order agreements from the database",
			Err:        err,
		}
	}

	// Fetch work order bookings
	var workorderBookings []transactionworkshoppayloads.WorkOrderBookingsResponse
	if err := tx.Model(&transactionworkshopentities.BookingEstimationAllocation{}).
		Where("booking_system_number = ? and booking_system_number != 0", entity.BookingSystemNumber).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderBookings).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order bookings from the database",
			Err:        err,
		}
	}

	// Fetch work order estimations
	var workorderEstimations []transactionworkshoppayloads.WorkOrderEstimationsResponse
	if err := tx.Model(&transactionworkshopentities.BookingEstimationServiceDiscount{}).
		Where("estimation_system_number = ? and estimation_system_number != 0", entity.EstimationSystemNumber).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderEstimations).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order estimations from the database",
			Err:        err,
		}
	}

	// Fetch work order contracts
	var workorderContracts []transactionworkshoppayloads.WorkOrderContractsResponse
	if err := tx.Model(&transactionworkshopentities.ContractService{}).
		Where("contract_service_system_number = ? and contract_service_system_number != 0", entity.ContractServiceSystemNumber).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Find(&workorderContracts).Error; err != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order billings from the database",
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order booking from the database",
			Err:        err,
		}
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save the updated work order booking",
			Err:        err,
		}
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
				Err:        err,
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
	vehicleUrl := config.EnvConfigs.SalesServiceUrl + "vehicle-master?page=0&limit=1000000&vehicle_id=" + strconv.Itoa(entity.VehicleId)
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
		Select("trx_work_order_detail.work_order_detail_id, trx_work_order_detail.work_order_system_number, trx_work_order_detail.line_type_id,lt.line_type_code, trx_work_order_detail.transaction_type_id, tt.transaction_type_code AS transaction_type_code, trx_work_order_detail.job_type_id, tc.job_type_code AS job_type_code, trx_work_order_detail.warehouse_group_id, trx_work_order_detail.frt_quantity, trx_work_order_detail.supply_quantity, trx_work_order_detail.operation_item_price, trx_work_order_detail.operation_item_discount_amount, trx_work_order_detail.operation_item_discount_request_amount").
		Joins("INNER JOIN mtr_work_order_line_type AS lt ON lt.line_type_code = trx_work_order_detail.line_type_id").
		Joins("INNER JOIN mtr_work_order_transaction_type AS tt ON tt.transaction_type_id = trx_work_order_detail.transaction_type_id").
		Joins("INNER JOIN mtr_work_order_job_type AS tc ON tc.job_type_id = trx_work_order_detail.job_type_id").
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save the work order affiliate",
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) SaveAffiliated(tx *gorm.DB, IdWorkorder int, id int, request transactionworkshoppayloads.WorkOrderAffiliatedRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ? AND affiliate_id = ?", IdWorkorder, id).
		First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order affiliate from the database",
			Err:        err,
		}
	}

	entity.WorkOrderSystemNumber = request.WorkOrderSystemNumber

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save the updated work order affiliate",
			Err:        err,
		}
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order status from the database",
			Err:        err,
		}
	}

	entity.WorkOrderStatusCode = request.WorkOrderStatusCode
	entity.WorkOrderStatusDescription = request.WorkOrderStatusName

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update work order status",
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteStatus(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterStatus
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterStatus{}).Where("work_order_status_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order status from the database",
			Err:        err,
		}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete work order status",
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterType

	query := utils.ApplyFilter(tx, filter)

	if err := query.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order type from the database",
			Err:        err,
		}
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
			Message:    "Failed to save work order type",
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateType(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderTypeRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterType{}).Where("work_order_type_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order type from the database",
			Err:        err,
		}
	}

	entity.WorkOrderTypeCode = request.WorkOrderTypeCode
	entity.WorkOrderTypeDescription = request.WorkOrderTypeName

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update work order type",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteType(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterType{}).Where("work_order_type_id = ?", id).First(&entity).Error

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order type from the database",
			Err:        err,
		}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete work order type",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewLineType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshoppayloads.Linetype, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterLineType

	query := utils.ApplyFilter(tx, filter)

	if err := query.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order line type from the database",
			Err:        err,
		}
	}

	var getlinetype []transactionworkshoppayloads.Linetype
	for _, t := range types {
		getlinetype = append(getlinetype, transactionworkshoppayloads.Linetype{
			LineTypeId:   t.WorkOrderLineTypeId,
			LineTypeCode: t.WorkOrderLineTypeCode,
			LineTypeName: t.WorkOrderLineTypeDescription,
		})
	}

	return getlinetype, nil
}

func (r *WorkOrderRepositoryImpl) AddLineType(tx *gorm.DB, request transactionworkshoppayloads.Linetype) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderMasterLineType{
		WorkOrderLineTypeCode:        request.LineTypeCode,
		WorkOrderLineTypeDescription: request.LineTypeName,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save linetype data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateLineType(tx *gorm.DB, id int, request transactionworkshoppayloads.Linetype) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterLineType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterLineType{}).Where("line_type_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve linetype data from the database",
			Err:        err,
		}
	}

	entity.WorkOrderLineTypeCode = request.LineTypeCode
	entity.WorkOrderLineTypeDescription = request.LineTypeName

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update linetype data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteLineType(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterLineType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterLineType{}).Where("line_type_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve linetype data from the database",
			Err:        err,
		}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete linetype data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewBill(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshoppayloads.WorkOrderBillable, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterBillAbleto

	query := utils.ApplyFilter(tx, filter)

	if err := query.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order type from the database",
			Err:        err,
		}
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
			Message:    "Failed to save billable data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateBill(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderBillableRequest) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterBillAbleto
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterBillAbleto{}).Where("billable_to_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve billable data from the database",
			Err:        err,
		}
	}

	entity.WorkOrderBillabletoName = request.BillableToName
	entity.WorkOrderBillabletoCode = request.BillableToCode

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update billable data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteBill(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterBillAbleto
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterBillAbleto{}).Where("billable_to_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve billable data from the database",
			Err:        err,
		}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete billable data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewTrxType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshoppayloads.WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterTrxType

	query := utils.ApplyFilter(tx, filter)

	if err := query.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order type from the database",
			Err:        err,
		}
	}

	var payloadTypes []transactionworkshoppayloads.WorkOrderTransactionType
	for _, t := range types {
		payloadTypes = append(payloadTypes, transactionworkshoppayloads.WorkOrderTransactionType{
			TransactionTypeId:   t.WorkOrderTrxTypeId,
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
			Message:    "Failed to save transaction type data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateTrxType(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterTrxType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterTrxType{}).Where("work_order_transaction_type_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve transaction type data from the database",
			Err:        err,
		}
	}

	entity.WorkOrderTrxTypeDescription = request.TransactionTypeName
	entity.WorkOrderTrxTypeCode = request.TransactionTypeCode

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update transaction type data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteTrxType(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterTrxType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterTrxType{}).Where("work_order_transaction_type_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve transaction type data from the database",
			Err:        err,
		}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete transaction type data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewTrxTypeSo(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshoppayloads.WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterTrxSoType

	query := utils.ApplyFilter(tx, filter)

	if err := query.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order type from the database",
			Err:        err,
		}
	}

	var payloadsoTypes []transactionworkshoppayloads.WorkOrderTransactionType
	for _, t := range types {
		payloadsoTypes = append(payloadsoTypes, transactionworkshoppayloads.WorkOrderTransactionType{
			TransactionTypeId:   t.WorkOrderTrxTypeSoId,
			TransactionTypeCode: t.WorkOrderTrxTypeSoCode,
			TransactionTypeName: t.WorkOrderTrxTypeSoDescription,
		})
	}

	return payloadsoTypes, nil

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
			Message:    "Failed to save transaction type so data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateTrxTypeSo(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderTransactionType) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterTrxSoType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterTrxSoType{}).Where("work_order_transaction_type_so_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve transaction type so data from the database",
			Err:        err,
		}
	}

	entity.WorkOrderTrxTypeSoDescription = request.TransactionTypeName
	entity.WorkOrderTrxTypeSoCode = request.TransactionTypeCode

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update transaction type so data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteTrxTypeSo(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterTrxSoType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterTrxSoType{}).Where("work_order_transaction_type_so_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve transaction type so data from the database",
			Err:        err,
		}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete transaction type so data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) NewJobType(tx *gorm.DB, filter []utils.FilterCondition) ([]transactionworkshoppayloads.WorkOrderJobType, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterJobType

	query := utils.ApplyFilter(tx, filter)

	if err := query.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order job type from the database",
			Err:        err,
		}
	}

	var payloadJobTypes []transactionworkshoppayloads.WorkOrderJobType
	for _, t := range types {
		payloadJobTypes = append(payloadJobTypes, transactionworkshoppayloads.WorkOrderJobType{
			JobTypeId:   t.WorkOrderJobTypeId,
			JobTypeCode: t.WorkOrderJobTypeCode,
			JobTypeName: t.WorkOrderJobTypeDescription,
		})
	}

	return payloadJobTypes, nil
}

func (r *WorkOrderRepositoryImpl) AddJobType(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderJobType) (bool, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.WorkOrderMasterJobType{
		WorkOrderJobTypeCode:        request.JobTypeCode,
		WorkOrderJobTypeDescription: request.JobTypeName,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save job type data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) UpdateJobType(tx *gorm.DB, id int, request transactionworkshoppayloads.WorkOrderJobType) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterJobType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterJobType{}).Where("job_type_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve job type data from the database",
			Err:        err,
		}
	}

	entity.WorkOrderJobTypeCode = request.JobTypeCode
	entity.WorkOrderJobTypeDescription = request.JobTypeName

	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update job type data",
			Err:        err,
		}
	}

	return true, nil
}

func (r *WorkOrderRepositoryImpl) DeleteJobType(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderMasterJobType
	err := tx.Model(&transactionworkshopentities.WorkOrderMasterJobType{}).Where("job_type_id = ?", id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve job type data from the database",
			Err:        err,
		}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete job type data",
			Err:        err,
		}
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
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order service vehicle from the database",
			Err:        err,
		}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete work order service vehicle",
			Err:        err}
	}

	return true, nil
}

func (s *WorkOrderRepositoryImpl) DeleteDetailWorkOrderMultiId(tx *gorm.DB, Id int, DetailIds []int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderDetail
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).Where("work_order_system_number = ? AND work_order_detail_id = ? IN (?)", Id, DetailIds).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order detail from the database",
			Err:        err}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete work order detail",
			Err:        err}
	}

	return true, nil
}

func (s *WorkOrderRepositoryImpl) DeleteRequestMultiId(tx *gorm.DB, Id int, DetailIds []int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderService
	err := tx.Model(&transactionworkshopentities.WorkOrderService{}).Where("work_order_system_number = ? AND work_order_service_id IN (?)", Id, DetailIds).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        err,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order service from the database",
			Err:        err,
		}
	}

	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete work order service",
			Err:        err}
	}

	return true, nil
}

// usp_comLookUp
// IF @strEntity = 'CustomerByTypeAndAddress'--CUSTOMER MASTER
// uspg_wtWorkOrder0_Update
// IF @Option = 8
// --USE FOR : * WORK ORDER CHANGE BILL TO
func (s *WorkOrderRepositoryImpl) ChangeBillTo(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.ChangeBillToRequest) (transactionworkshoppayloads.ChangeBillToResponse, *exceptions.BaseErrorResponse) {
	var existingWorkOrder struct {
		WorkOrderOperationItemLine int
	}

	// Retrieve work order item line
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("work_order_operation_item_line").
		Where("work_order_system_number = ? AND transaction_type_id = 3", workOrderId). // 3 = External, and invoice_system_number should not be null
		First(&existingWorkOrder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.ChangeBillToResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Failed to retrieve work order item line from the database",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.ChangeBillToResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "External detail has already been invoiced",
			Err:        err,
		}
	}

	// Retrieve the work order entity
	var entity transactionworkshopentities.WorkOrder
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.ChangeBillToResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Failed to retrieve work order from the database",
			Err:        err,
		}
	}

	// Update customer ID
	entity.CustomerId = request.BillToCustomerId
	entity.BillableToId = request.BillableToId

	// Save the updated entity
	err = tx.Save(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.ChangeBillToResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update the billable to data",
			Err:        err,
		}
	}

	// Return the response as a value, not a pointer
	return transactionworkshoppayloads.ChangeBillToResponse{
		WorkOrderSystemNumber: workOrderId,
		BillToCustomerId:      entity.CustomerId,
		BillableToId:          entity.BillableToId,
	}, nil
}

// uspg_wtWorkOrder0_Update
// IF @Option = 13
//
//	--USE FOR : * WORK ORDER CHANGE PHONE NO
func (s *WorkOrderRepositoryImpl) ChangePhoneNo(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.ChangePhoneNoRequest) (*transactionworkshoppayloads.ChangePhoneNoResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Failed to retrieve work order from the database",
			Err:        err,
		}
	}

	entity.ContactPersonPhone = request.PhoneNo

	err = tx.Save(&entity).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update the phone number data",
			Err:        err,
		}
	}

	return &transactionworkshoppayloads.ChangePhoneNoResponse{
		WorkOrderSystemNumber: workOrderId,
		CustomerId:            entity.CustomerId,
		BillableToId:          entity.BillableToId,
		PhoneNo:               entity.ContactPersonPhone,
	}, nil
}

// uspg_wtWorkOrder2_Update
// IF @Option = 14
// --USE FOR : CONFIRM PRICE
func (s *WorkOrderRepositoryImpl) ConfirmPrice(tx *gorm.DB, workOrderId int, idwos []int, request transactionworkshoppayloads.WorkOrderConfirmPriceRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	var response transactionworkshopentities.WorkOrderDetail
	var markupPercentage, totalPackage, totalOpr, totalPart, totalOil, totalMaterial, totalConsumableMaterial, totalSublet, totalAccs float64
	var invSysNo, woOprItemLine int
	var oprItemCode, vehicleChassisNo string
	var total, totalDisc, totalAfterDisc, totalNonVat float64
	var greyMarket bool
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
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
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
		Select("ISNULL(line_type_id, 0) AS LineType, ISNULL(invoice_system_number, 0) AS InvoiceSysNo, ISNULL(warehouse_group_id, 0) AS WhsGroup").
		Where("work_order_system_number = ? AND work_order_operation_item_line IN (?)", workOrderId, idwos).
		Scan(&detailResult).Error

	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
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
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        errVehicle,
		}
	}

	vehicleChassisNo = vehicleResponses[0].VehicleCode

	// Check if the vehicle is in the grey market by looking up the vehicle chassis number
	err = tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle").
		Select("ISNULL(vehicle_is_grey_market, 0)"). // Mengubah tipe hasil query ke boolean
		Where("vehicle_chassis_number = ?", vehicleChassisNo).
		Scan(&greyMarket).Error // Scan ke tipe boolean

	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check vehicle grey market status",
			Err:        err,
		}
	}

	// Check for grey market markup and update
	if greyMarketMarkupPercentageExists && greyMarket {
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
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update work order details with grey market markup",
					Err:        err,
				}
			}
		}

		// Compute totals
		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0), 0, 0)), 0)").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypePackage).
			Scan(&totalPackage).Error

		if err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total package",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOperation).
			Scan(&totalOpr).Error

		if err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total operation",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSparepart).
			Scan(&totalPart).Error

		if err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total sparepart",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOil).
			Scan(&totalOil).Error

		if err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total oil",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeMaterial).
			Scan(&totalMaterial).Error

		if err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total material",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSublet).
			Scan(&totalSublet).Error

		if err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total sublet",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeAccesories).
			Scan(&totalAccs).Error

		if err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total accessories",
				Err:        err,
			}
		}

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeConsumableMaterial).
			Scan(&totalConsumableMaterial).Error

		if err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total consumable material",
				Err:        err,
			}
		}

		// Calculate total and discounts
		total = totalPackage + totalOpr + totalPart + totalOil + totalMaterial + totalSublet + totalAccs + totalConsumableMaterial

		// Calculate discounts and VAT
		if err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select(`
				SUM(CASE 
					WHEN line_type_id = ? THEN 
						CASE WHEN approval_id = 20 AND ISNULL(operation_item_discount_request_amount, 0) > 0 
						THEN ISNULL(operation_item_discount_request_amount, 0) 
						ELSE ISNULL(operation_item_discount_amount, 0) 
						END 
					ELSE 
						CASE WHEN approval_id = 20 AND ISNULL(operation_item_discount_request_amount, 0) > 0 
						THEN ISNULL(operation_item_discount_request_amount, 0) 
						ELSE ISNULL(operation_item_discount_amount, 0) 
						END 
						* 
						CASE WHEN line_type_id <> ? THEN ISNULL(frt_quantity, 0) 
						ELSE CASE WHEN ISNULL(supply_quantity, 0) > 0 
						THEN ISNULL(supply_quantity, 0) ELSE ISNULL(frt_quantity, 0) END 
						END 
				END)`, utils.LinetypePackage, utils.LinetypeOperation).
			Scan(&totalDisc).Error; err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total discount",
				Err:        err,
			}
		}

		totalDisc = math.Round(totalDisc)
		totalAfterDisc = total - totalDisc

		var totalVat float64
		if entity.VATTaxRate != nil {
			totalVat = (totalAfterDisc * (*entity.VATTaxRate)) / 100
		} else {
			totalVat = 0 // Set to 0 if VATTaxRate is nil
		}

		totalNonVat = totalAfterDisc + totalVat

		// Update totals in the WorkOrder table
		if err = tx.Model(&transactionworkshopentities.WorkOrder{}).
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
			}).Error; err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update work order totals",
				Err:        err,
			}
		}
	}

	return response, nil
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

					qtyAvail, errResponse := s.lookupRepo.SelectLocationStockItem(tx, 1, entity.CompanyId, time.Now(), 0, "", detailentity.OperationItemId, detailentity.WarehouseGroupId, "S")
					if errResponse != nil {
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
						`, entity.CompanyId, detailentity.OperationItemCode, detailentity.WarehouseGroupId, "S").Error
						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to insert data into temporary table",
								Err:        err,
							}
						}

						// Step 3: Check if the original item is in the substitution table
						var exists bool

						type Substitution struct {
							SubsItemCode string `gorm:"column:SUBS_ITEM_CODE"`
						}

						err := tx.Model(&Substitution{}).
							Select("CASE WHEN EXISTS (SELECT 1 FROM #SUBS WHERE SUBS_ITEM_CODE = ?) THEN 1 ELSE 0 END", detailentity.OperationItemCode).
							Scan(&exists).Error

						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to check if the original item is in the substitution table",
								Err:        err,
							}
						}

						if !exists {
							// Step 4: Fetch the substitute items from the temporary table
							type SubstituteItem struct {
								SubsItemCode string  `gorm:"column:SUBS_ITEM_CODE"`
								ItemName     string  `gorm:"column:ITEM_NAME"`
								SupplyQty    float64 `gorm:"column:SUPPLY_QTY"`
								SubsType     string  `gorm:"column:SUBS_TYPE"`
							}

							var substituteItems []SubstituteItem

							err := tx.Table("#SUBS").
								Select("SUBS_ITEM_CODE, ITEM_NAME, SUPPLY_QTY, SUBS_TYPE").
								Find(&substituteItems).Error
							if err != nil {
								return false, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to fetch substitute items from the temporary table",
									Err:        err,
								}
							}

							// Step 5: Process each substitute item
							for _, substituteItem := range substituteItems {

								// Step 6: Check and update the original item if not substituted
								if substituteItem.SubsType != "" {
									err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
										Where("work_order_system_number = ? AND operation_item_code = ? AND work_order_operation_item_line = ? AND substitute_type_id IS NULL",
											workOrderId, detailentity.OperationItemCode, idwos).
										Updates(map[string]interface{}{
											"substitute_type_id": 1,
											"substitute_type":    "SUBSTITUTE_ITEM",
											"warehouse_group_id": detailentity.WarehouseGroupId,
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
									oprItemPrice, _ = s.lookupRepo.GetOprItemPrice(tx, detailentity.LineTypeId, entity.CompanyId, detailentity.OperationItemId, entity.BrandId, entity.ModelId, detailentity.JobTypeId, entity.VariantId, entity.CurrencyId, utils.TrxTypeWoWarranty.ID, "1")

									// Apply markup to the item price
									oprItemPrice = oprItemPrice + 10.00 + (oprItemPrice * (5.00 / 100))

									// Insert the substitute item into wtWorkOrder2
									err = tx.Create(&transactionworkshopentities.WorkOrderDetail{
										WorkOrderSystemNumber:      workOrderId,
										WorkOrderOperationItemLine: nextLineNumber,
										LineTypeId:                 detailentity.LineTypeId,
										WorkorderStatusId:          detailentity.WorkorderStatusId,
										OperationItemCode:          detailentity.OperationItemCode,
										OperationItemPrice:         oprItemPrice,
										SubstituteItemCode:         substituteItem.SubsItemCode,
										SupplyQuantity:             substituteItem.SupplyQty,
										WarehouseGroupId:           detailentity.WarehouseGroupId,
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

			//-- By default price of all items will be replaced with the price from ItemPriceList
			//-- Exclude Fee from the replacing process
			var oprItemPrice, oprItemPriceDisc, discountPercent float64
			var markupAmount, markupPercentage float64
			var warrantyClaimType string

			if detailentity.LineTypeId == utils.LinetypeSublet {

				// Fetch Opr_Item_Price
				oprItemPrice, _ = s.lookupRepo.GetOprItemPrice(tx, detailentity.LineTypeId, entity.CompanyId, detailentity.OperationItemId, entity.BrandId, entity.ModelId, detailentity.JobTypeId, entity.VariantId, entity.CurrencyId, utils.TrxTypeWoWarranty.ID, "1")

				// Set markup percentage based on company ID
				if entity.CompanyId == 139 {
					markupPercentage = 11.00
				}

				// // Apply markup amount and percentage
				oprItemPrice = oprItemPrice + markupAmount + (oprItemPrice * (markupPercentage / 100))

				// Fetch Opr_Item_Disc_Percent
				oprItemPriceDisc, _ = s.lookupRepo.GetOprItemDisc(tx, detailentity.LineTypeId, 6, detailentity.OperationItemId, entity.AgreementGeneralRepairId, entity.ProfitCenterId, detailentity.FrtQuantity*detailentity.OperationItemPrice, entity.CompanyId, entity.BrandId, entity.ContractServiceSystemNumber, utils.TrxTypeWoWarranty.ID, utils.EstWoOrderTypeId)

			} else {
				// Fetch Opr_Item_Price
				oprItemPrice, _ = s.lookupRepo.GetOprItemPrice(tx, detailentity.LineTypeId, entity.CompanyId, detailentity.OperationItemId, entity.BrandId, entity.ModelId, detailentity.JobTypeId, entity.VariantId, entity.CurrencyId, utils.TrxTypeWoWarranty.ID, "1")

				// Set markup percentage based on company ID
				if entity.CompanyId == 139 {
					markupPercentage = 11.00
				}

				// // Apply markup amount and percentage
				oprItemPrice = oprItemPrice + markupAmount + (oprItemPrice * (markupPercentage / 100))

				// Fetch Opr_Item_Disc_Percent
				oprItemPriceDisc, _ = s.lookupRepo.GetOprItemDisc(tx, detailentity.LineTypeId, 6, detailentity.OperationItemId, entity.AgreementGeneralRepairId, entity.ProfitCenterId, detailentity.FrtQuantity*detailentity.OperationItemPrice, entity.CompanyId, entity.BrandId, entity.ContractServiceSystemNumber, utils.TrxTypeWoWarranty.ID, utils.EstWoOrderTypeId)

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
				Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSublet).
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
					workOrderId, utils.LinetypePackage, utils.LinetypeOperation, utils.LinetypeSparepart, utils.LinetypeOil, utils.LinetypeMaterial, utils.LinetypeConsumableMaterial, utils.LinetypeSublet).
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

// uspg_wtWorkOrder0_Update
// IF @Option = 7
func (s *WorkOrderRepositoryImpl) DeleteCampaign(tx *gorm.DB, workOrderId int) (transactionworkshoppayloads.DeleteCampaignPayload, *exceptions.BaseErrorResponse) {
	// Scan campaign data from work order
	var campaignId int
	if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Select("ISNULL(campaign_id, 0) AS campaign_id").
		Where("work_order_system_number = ?", workOrderId).
		Scan(&campaignId).Error; err != nil {

		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check campaign data from work order, campaign is empty",
			Err:        err,
		}
	}

	// Check if operation is already allocated
	var exists bool
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("1").
		Where("work_order_system_number = ? AND (line_type_id = ? OR line_type_id = ?) AND transaction_type_id = ? AND ISNULL(service_status_id, '') <> ''", workOrderId, utils.LinetypeOperation, utils.LinetypePackage, 5).
		Scan(&exists).Error; err != nil {

		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check operation allocation",
			Err:        err,
		}
	}

	if exists {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Operation already allocated",
			Err:        errors.New("operation already allocated"),
		}
	}

	// Check if SUPPLY_QTY is valid
	var supplyQuantity float64
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("ISNULL(supply_quantity, 0) AS supply_quantity").
		Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOperation).
		Scan(&supplyQuantity).Error; err != nil {

		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check supply quantity",
			Err:        err,
		}
	}

	if supplyQuantity <= 0 {
		fmt.Println("Invalid supply quantity:", supplyQuantity) // Tambahkan log
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "SUPPLY QTY is not Valid",
			Err:        errors.New("supply qty is not Valid"),
		}
	}

	// Delete work order lines
	if err := tx.Where("work_order_system_number = ? AND transaction_type_id = ?", workOrderId, 5).
		Delete(&transactionworkshopentities.WorkOrderDetail{}).Error; err != nil {

		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete work order lines",
			Err:        err,
		}
	}

	// Declare totals
	var totalPackage, totalOpr, totalPart, totalOil, totalMaterial, totalConsumableMaterial, totalSublet, totalAccs, totalDisc, totalAfterDisc, totalNonVat, totalVat, totalAfterVat float64

	substituteTypeItem := "S"

	// Calculate total package
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("SUM(ROUND(ISNULL(operation_item_price, 0), 0, 0))").
		Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypePackage).
		Scan(&totalPackage).Error; err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total package",
			Err:        err,
		}
	}

	// Calculate total operation
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
		Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOperation).
		Scan(&totalOpr).Error; err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total operation",
			Err:        err,
		}
	}

	// Calculate total sparepart
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeSparepart, substituteTypeItem).
		Scan(&totalPart).Error; err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total sparepart",
			Err:        err,
		}
	}

	// Calculate total oil
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeOil, substituteTypeItem).
		Scan(&totalOil).Error; err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total oil",
			Err:        err,
		}
	}

	// Calculate total material
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeMaterial, substituteTypeItem).
		Scan(&totalMaterial).Error; err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total material",
			Err:        err,
		}
	}

	// Calculate total consumable material
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeConsumableMaterial, substituteTypeItem).
		Scan(&totalConsumableMaterial).Error; err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total consumable material",
			Err:        err,
		}
	}

	// Calculate total sublet
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeSublet, substituteTypeItem).
		Scan(&totalSublet).Error; err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total sublet",
			Err:        err,
		}
	}

	// Calculate total accessories
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeAccesories, substituteTypeItem).
		Scan(&totalAccs).Error; err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total accessories",
			Err:        err,
		}
	}

	// Calculate overall total
	total := totalPackage + totalOpr + totalPart + totalOil + totalMaterial + totalConsumableMaterial + totalSublet + totalAccs

	// Calculate total discount
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("SUM(CASE WHEN line_type_id = ? THEN ISNULL(operation_item_discount_amount, 0) ELSE ISNULL(operation_item_discount_amount, 0) * ISNULL(frt_quantity, 0) END)", utils.LinetypePackage).
		Where("work_order_system_number = ?", workOrderId).
		Scan(&totalDisc).Error; err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total discount",
			Err:        err,
		}
	}

	// Calculate total after discount
	totalAfterDisc = total - totalDisc

	// Calculate total non-VAT
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0))").
		Where("work_order_system_number = ? AND line_type_id NOT IN (?, ?)", workOrderId, utils.LinetypePackage, utils.LinetypeOperation).
		Scan(&totalNonVat).Error; err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total non-VAT",
			Err:        err,
		}
	}

	// Calculate total VAT
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0) * ISNULL(vat_rate, 0) / 100, 0, 0))").
		Where("work_order_system_number = ?", workOrderId).
		Scan(&totalVat).Error; err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total VAT",
			Err:        err,
		}
	}

	// Calculate total after VAT
	totalAfterVat = totalAfterDisc + totalVat

	// Create payload
	payload := transactionworkshoppayloads.DeleteCampaignPayload{
		TotalPackage:            totalPackage,
		TotalOpr:                totalOpr,
		TotalPart:               totalPart,
		TotalOil:                totalOil,
		TotalMaterial:           totalMaterial,
		TotalConsumableMaterial: totalConsumableMaterial,
		TotalSublet:             totalSublet,
		TotalAccs:               totalAccs,
		TotalDisc:               totalDisc,
		TotalAfterDisc:          totalAfterDisc,
		TotalNonVat:             totalNonVat,
		TotalVat:                totalVat,
		TotalAfterVat:           totalAfterVat,
		AddDiscStat:             "APPROVED",
		WorkOrderSystemNumber:   workOrderId,
		CampaignId:              campaignId,
	}

	return payload, nil
}

// uspg_wtWorkOrder2_Insert
// IF @Option = 1
// --USE FOR : * INSERT NEW DATA FROM PACKAGE IN CONTRACT SERVICE
func (s *WorkOrderRepositoryImpl) AddContractService(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.WorkOrderContractServiceRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {

	type ContractServiceData struct {
		ContractServSysNo float64 `gorm:"column:contract_service_system_number"`
		AddDiscStat       int     `gorm:"column:additional_discount_status_approval_id"`
		WhsGroup          float64 `gorm:"column:whs_group"`
		TaxFree           int     `gorm:"column:tax_free"`
	}

	var woentities transactionworkshopentities.WorkOrder
	var response transactionworkshopentities.WorkOrderDetail
	var contractServiceData ContractServiceData

	// Initialize variables
	var (
		csrDescription, pphTaxCode                                                                     string
		csrFrtQty, csrPrice, csrDiscPercent, addDiscReqAmount, newFrtQty, supplyQty, oprItemDiscAmount float64
		csrOprItemCode, wcfTypeMoney, woOprItemLine, csrLineType, atpmWcfType, addDiscStat, itemTypeId int
	)

	// Set default WCF type
	wcfTypeMoney = 1

	// Fetch contract service data
	if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_customer ON trx_work_order.customer_id = mtr_customer.customer_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_client_type ON mtr_customer.client_type_id = mtr_client_type.client_type_id").
		Select("contract_service_system_number, additional_discount_status_approval").
		Where("work_order_system_number = ?", workOrderId).
		Scan(&contractServiceData).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch contract service data",
			Err:        err,
		}
	}

	fmt.Println("Contract Service Data: ", contractServiceData)

	// Initialize new freight quantity
	newFrtQty = 0

	// Fetch contract service items
	type ContractServiceItem struct {
		LineType       int     `gorm:"column:line_type_id"`
		OprItemCode    int     `gorm:"column:operation_id"`
		Description    string  `gorm:"column:description"`
		FrtQty         float64 `gorm:"column:frt_quantity"`
		OprItemPrice   float64 `gorm:"column:operation_price"`
		OprItemDiscPct float64 `gorm:"column:operation_discount_percent"`
	}

	var contractServiceItems []ContractServiceItem
	if err := tx.Model(&transactionworkshopentities.ContractServiceOperationDetail{}).
		Select("line_type_id, operation_id, description, frt_quantity, operation_price, operation_discount_percent").
		Where("contract_service_system_number = ? AND package_id = ?", contractServiceData.ContractServSysNo, request.PackageCodeId).
		Find(&contractServiceItems).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch contract service items",
			Err:        err,
		}
	}

	fmt.Println("Contract Service Items: ", contractServiceItems)

	// Process each contract service item
	for _, item := range contractServiceItems {
		csrLineType = item.LineType
		csrOprItemCode = item.OprItemCode
		csrDescription = item.Description
		csrFrtQty = item.FrtQty
		csrPrice = item.OprItemPrice
		csrDiscPercent = item.OprItemDiscPct

		// Get the next available work order operation item line
		if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("ISNULL(MAX(work_order_operation_item_line), 0) + 1").
			Where("work_order_system_number = ?", workOrderId).
			Scan(&woOprItemLine).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get the next work order operation item line",
				Err:        err,
			}
		}

		fmt.Println("Work Order Operation Item Line: ", woOprItemLine)

		// Set Atpm_Wcf_Type based on conditions
		if err := tx.Model(&masteritementities.Item{}).
			Select("atpm_warranty_claim_type_id").
			Where("item_code = ?", csrOprItemCode).
			Scan(&atpmWcfType).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Atpm Wcf Type",
				Err:        err,
			}
		}

		fmt.Println("Atpm Wcf Type: ", atpmWcfType)

		// If Atpm_Wcf_Type is empty, set it to wcfTypeMoney
		if atpmWcfType == 0 {
			atpmWcfType = wcfTypeMoney
		}

		// Handle logic based on LineType
		switch csrLineType {
		case utils.LinetypePackage:
			csrFrtQty = 1
			supplyQty = 1
			atpmWcfType = 0
		case utils.LinetypeOperation:
			// Fetch PPH tax code for operations
			if err := tx.Model(&masteroperationentities.OperationModelMapping{}).
				Select("tax_code").
				Where("operation_id = ?", csrOprItemCode).
				Scan(&pphTaxCode).Error; err != nil {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch PPH tax code",
					Err:        err,
				}
			}
			fmt.Println("PPH Tax Code: ", pphTaxCode)
			supplyQty = csrFrtQty
			atpmWcfType = 0
		default:
			// Fetch item UOM and type for other items
			type ItemUOMType struct {
				ItemUom    string `gorm:"column:unit_of_measurement_selling_id"`
				ItemTypeId int    `gorm:"column:item_type_id"`
			}

			var itemDetails ItemUOMType

			if err := tx.Model(&masteritementities.Item{}).
				Select("unit_of_measurement_selling_id, item_type").
				Where("item_code = ?", csrOprItemCode).
				Scan(&itemDetails).Error; err != nil {
				return response, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch item UOM and type",
					Err:        err,
				}
			}

			itemTypeId = itemDetails.ItemTypeId

			supplyQty = 0
			if itemTypeId == 2 {
				supplyQty = csrFrtQty
			}
		}

		oprItemDiscAmount = math.Round(csrPrice * csrDiscPercent / 100)

		workOrderLine := transactionworkshopentities.WorkOrderDetail{
			WorkOrderSystemNumber:        workOrderId,
			WorkOrderOperationItemLine:   woOprItemLine,
			LineTypeId:                   csrLineType,
			OperationItemId:              csrOprItemCode,
			Description:                  csrDescription,
			FrtQuantity:                  csrFrtQty,
			OperationItemPrice:           csrPrice,
			OperationItemDiscountAmount:  oprItemDiscAmount,
			OperationItemDiscountPercent: csrDiscPercent,
			SupplyQuantity:               supplyQty,
			WarehouseGroupId:             int(contractServiceData.WhsGroup),
			AtpmWCFTypeId:                atpmWcfType,
		}

		// Insert into work order line table
		if err := tx.Create(&workOrderLine).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert work order line",
				Err:        err,
			}
		}

		// Update new freight quantity if LineType is "Operation" or "Package"
		if csrLineType == utils.LinetypeOperation || csrLineType == utils.LinetypePackage {
			newFrtQty += csrFrtQty
		}
	}

	var estTime, sumFrtQty, totalPackage, totalOpr, totalPart, totalOil, totalMaterial, totalConsumableMaterial, totalSublet, totalAccs, totalNonVat, totalVat, totalAfterDisc, totalAfterVat, totalPph, totalDisc float64
	var timeTolerance float64 = 0.25

	// Calculate EST_TIME
	if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Select("ISNULL(estimate_time, 0)").
		Where("work_order_system_number = ?", workOrderId).
		Scan(&estTime).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate EST_TIME",
			Err:        err,
		}
	}

	fmt.Println("EST_TIME: ", estTime)

	if estTime == 0 {
		if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("SUM(ISNULL(frt_quantity, 0))").
			Where("work_order_system_number = ? AND (line_type_id = ? OR line_type_id = ?)", workOrderId, utils.LinetypeOperation, utils.LinetypePackage).
			Scan(&sumFrtQty).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate sum of frt_qty",
				Err:        err,
			}
		}

		estTime = sumFrtQty * timeTolerance
	} else {
		estTime += newFrtQty
	}

	// Calculate totals for various line types
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0), 0, 0)), 0)").
		Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypePackage).
		Scan(&totalPackage).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total package",
			Err:        err,
		}
	}

	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
		Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOperation).
		Scan(&totalOpr).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total operation",
			Err:        err,
		}
	}

	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
		Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSparepart).
		Scan(&totalPart).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total sparepart",
			Err:        err,
		}
	}

	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
		Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOil).
		Scan(&totalOil).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total oil",
			Err:        err,
		}
	}

	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
		Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeMaterial).
		Scan(&totalMaterial).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total material",
			Err:        err,
		}
	}

	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
		Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeConsumableMaterial).
		Scan(&totalConsumableMaterial).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total consumable material",
			Err:        err,
		}
	}

	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
		Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSublet).
		Scan(&totalSublet).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total sublet",
			Err:        err,
		}
	}

	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
		Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeAccesories).
		Scan(&totalAccs).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total accessories",
			Err:        err,
		}
	}

	if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
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
		}).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update work order totals",
			Err:        err,
		}
	}

	//--DELETE WO FROM CAR_WASH
	if err := tx.Where("work_order_system_number = ?", workOrderId).
		Delete(&transactionjpcbentities.CarWash{}).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete work order from car wash",
			Err:        err,
		}
	}

	//--==UPDATE TOTAL WORK ORDER==--
	// Query to calculate TOTAL_DISC
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select(`
	SUM(
		CASE
			WHEN line_type_id = ? THEN
				CASE
					WHEN approval_id = ? THEN COALESCE(operation_item_discount_request_amount, 0)
					ELSE COALESCE(operation_item_discount_amount, 0)
				END
			ELSE
				CASE
					WHEN approval_id = ? THEN COALESCE(operation_item_discount_request_amount, 0)
					ELSE COALESCE(operation_item_discount_amount, 0)
				END *
				CASE
					WHEN line_type_id <> ? THEN COALESCE(frt_quantity, 0)
					ELSE
						CASE
							WHEN COALESCE(supply_quantity, 0) > 0 THEN COALESCE(supply_quantity, 0)
							ELSE COALESCE(frt_quantity, 0)
						END
				END
		END
	)`,
			utils.LinetypePackage, 20, 20, utils.LinetypeOperation).
		Where("work_order_system_number = ?", workOrderId).
		Scan(&totalDisc).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total discount",
			Err:        err,
		}
	}

	// Calculate TOTAL
	total := totalPackage + totalOpr + totalPart + totalOil + totalMaterial + totalConsumableMaterial + totalSublet + totalAccs

	// Calculate AddDiscStat
	if woentities.AdditionalDiscountStatusApprovalId == 20 {
		if woentities.DiscountRequestPercent != nil && *woentities.DiscountRequestPercent > 0 {
			woentities.AdditionalDiscountStatusApprovalId = 30 // Use assignment
		}
	}

	addDiscStat = woentities.AdditionalDiscountStatusApprovalId

	// Safely dereference the pointer
	if woentities.DiscountRequestPercent != nil {
		addDiscReqAmount = *woentities.DiscountRequestPercent
	} else {
		addDiscReqAmount = 0
	}

	// Rounding TOTAL_DISC
	totalDisc = math.Round(totalDisc)

	// Total After Discount
	totalAfterDisc = math.Round(total - totalDisc)

	// Calculate TOTAL_NON_VAT
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
		Where("work_order_system_number = ? AND transaction_type_id = ?", workOrderId, 6). //TrxTypeWOInternal
		Scan(&totalNonVat).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total non-VAT",
			Err:        err,
		}
	}

	// VAT Calculation
	if contractServiceData.TaxFree == 0 {
		var vatRate float64
		if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
			Select("COALESCE(vat_tax_rate, 0.0)").
			Where("work_order_system_number = ?", workOrderId).
			Scan(&vatRate).Error; err != nil {
			return response, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch VAT tax rate",
				Err:        err,
			}
		}

		// VAT Amount Calculation
		totalVat = math.Floor((totalAfterDisc - totalNonVat) * vatRate / 100)

	} else {

		totalVat = 0
	}

	// Total After VAT
	totalAfterVat = math.Round(totalAfterDisc + totalVat)

	// PPH Calculation
	if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND (line_type_id = ? OR line_type_id = ?)", workOrderId, utils.LinetypePackage, utils.LinetypeOperation).
		Pluck("FLOOR(SUM(COALESCE(pph_amount, 0)))", &totalPph).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate total PPH",
			Err:        err,
		}
	}

	//update work order
	if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", workOrderId).
		Updates(map[string]interface{}{
			"total":                               total,
			"total_pph":                           totalPph,
			"total_discount":                      totalDisc,
			"total_after_discount":                totalAfterDisc,
			"total_vat":                           totalVat,
			"total_after_vat":                     totalAfterVat,
			"additional_discount_status_approval": addDiscStat,
			"discount_request_amount":             addDiscReqAmount,
		}).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update work order",
			Err:        err,
		}
	}

	return response, nil
}

// uspg_wtWorkOrder2_Insert
// IF @Option = 2
// --USE FOR : * INSERT NEW DATA FROM PACKAGE MASTER
func (s *WorkOrderRepositoryImpl) AddGeneralRepairPackage(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.WorkOrderGeneralRepairPackageRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderDetail
	var result struct {
		CompanyCode int
		WoDocNo     string
		LinetypeId  int
		JobTypeId   int
		AgreementNo int
		BillCodeExt int
		BrandId     int
		CampaignId  int
		Mileage     int
		VariantId   int
		CpcCode     int
	}
	const profitCenterGR = 2

	// Step 1: Fetch data from work order and related tables
	if err := tx.Table("trx_work_order AS wo").
		Select("wo.company_id, wo.work_order_document_number,trx_work_order_detail.line_type_id, trx_work_order_detail.job_type_id, wo.agreement_general_repair_id, trx_work_order_detail.transaction_type_id, wo.brand_id, wo.campaign_id, wo.variant_id").
		Joins("INNER JOIN dms_microservices_aftersales_dev.dbo.trx_work_order_detail ON wo.work_order_system_number = trx_work_order_detail.work_order_system_number").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_customer ON wo.customer_id = mtr_customer.customer_id").
		Where("wo.work_order_system_number = ?", workOrderId).
		Scan(&result).Error; err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch data from work order",
			Err:        err,
		}
	}

	fmt.Println("Work Order Data: ", result)

	// Step 2: Determine job type based on Profit Center GR
	if request.CPCCode == profitCenterGR {
		result.JobTypeId = 9 // "PM" - Job Type for Periodical Maintenance
	} else {
		result.JobTypeId = 13 // "TG" - Job Type for Transfer to General Repair
	}

	// Step 3: Fetch Whs_Group based on company code
	whsGroupValue, err := s.lookupRepo.GetWhsGroup(tx, result.CompanyCode)
	if err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch warehouse group",
			Err:        errors.New("failed to fetch warehouse group"),
		}
	}

	// Step 4: Set WhsGroup if conditions are met
	var whsGroup int
	if result.BrandId == 23 && whsGroupValue != 38 {
		whsGroup = 38
	}

	// Step 5: Fetch package details from the database
	var packages []struct {
		LineTypeId  int
		OprItemId   int
		OprItemCode string
		FrtQty      float64
		JobTypeId   int
		TrxTypeId   int
	}

	if err := tx.Table("mtr_package AS ap").
		Select("apd.line_type_id, apd.item_operation_id,ap.package_code, COALESCE(apd.frt_quantity, 0) AS frt_qty, apd.job_type_id, apd.workorder_transaction_type_id").
		Joins("INNER JOIN mtr_package_detail AS apd ON ap.package_id = apd.package_id").
		Where("ap.package_code = ?", request.PackageId).
		Scan(&packages).Error; err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch package data",
			Err:        err,
		}
	}

	// Step 6: Process each package and apply business logic
	for _, pkg := range packages {
		csrLineTypeId := pkg.LineTypeId
		csrOprItemId := pkg.OprItemId
		csrOprItemCode := pkg.OprItemCode
		csrFrtQty := pkg.FrtQty
		csrJobTypeId := pkg.JobTypeId
		csrTrxTypeId := pkg.TrxTypeId

		// Update FRT_QTY if LineType is a package
		if csrLineTypeId == utils.LinetypePackage {
			csrFrtQty = 1
		}

		// Update job type and transaction type if present
		if csrJobTypeId != 0 {
			result.JobTypeId = csrJobTypeId
		}
		if csrTrxTypeId != 0 {
			result.BillCodeExt = csrTrxTypeId
		}

		// Validation for chassis that have already undergone PDI, FSI, WR
		if csrJobTypeId == 8 || csrJobTypeId == 15 || csrJobTypeId == 4 {
			var blockingExists int64
			if err := tx.Model(&transactionworkshopentities.WorkOrderMasterBlockingChassis{}).
				Where("vehicle_id = ?", request.VehicleId).
				Count(&blockingExists).Error; err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check blocking status",
					Err:        err,
				}
			}

			if blockingExists > 0 {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "This vehicle has already been blocked for Free service Inspection, PDI Service or Warranty Claim",
					Err:        errors.New("this vehicle has already been blocked for Free service Inspection, PDI Service or Warranty Claim"),
				}
			}
		}

		// Ambil markup berdasarkan Company dan Vehicle Brand
		// OLDSITE
		// if err := tx.Model(&GmSiteMarkup{}).
		// 	Where("company_id = ? AND brand_id = ? AND site_code = ? AND trx_type_id = ?", result.CompanyCode, result.BrandId, siteCode, billCodeExt).
		// 	Select("markup_amount, markup_percentage").
		// 	Scan(&markupAmount, &markupPercentage).Error; err != nil {
		// 	return entity, &exceptions.BaseErrorResponse{
		// 		StatusCode: http.StatusInternalServerError,
		// 		Message:    "Failed to fetch markup",
		// 		Err:        err,
		// 	}
		// }
		// NEW SITE
		// Fetch markup based on Company and Vehicle Brand
		markupAmount := 0.0
		markupPercentage := 0.0

		// Handle GET DISC FOR CAMPAIGN
		if result.CampaignId != 0 {
			result.BillCodeExt = utils.TrxTypeWoCampaign.ID // utils.TrxTypeWoCampaign
		}

		// Fetch operation item price and discount percent
		oprItemPrice, err := s.lookupRepo.GetOprItemPrice(tx, result.CompanyCode, result.BrandId, csrOprItemId, result.AgreementNo, result.JobTypeId, csrLineTypeId, csrTrxTypeId, int(csrFrtQty), whsGroup, strconv.Itoa(result.VariantId))
		if err != nil {
			return entity, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch operation item price",
				Err:        errors.New("failed to fetch operation item price"),
			}
		}

		// Apply markup and percentage increase
		oprItemPrice += markupAmount + (oprItemPrice * (markupPercentage / 100))

		// Fetch item discount percent (You need to define or fix GetOprItemDiscPercent method)
		oprItemDiscPercent, err := s.lookupRepo.GetOprItemDisc(tx, result.LinetypeId, result.BillCodeExt, csrOprItemId, result.AgreementNo, result.CpcCode, oprItemPrice*csrFrtQty, result.CompanyCode, result.BrandId, 0, whsGroup, 0)
		if err != nil {
			return entity, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch operation item discount percent",
				Err:        errors.New("failed to fetch operation item discount percent"),
			}
		}

		// Calculate discount amount
		oprItemDiscAmount := math.Round((oprItemPrice * oprItemDiscPercent / 100))

		// Fetch Pph_Tax_Code and ATPM_WCF_Type
		var pphTaxCode, atpmWcfType int

		if err := tx.Model(&masteroperationentities.OperationModelMapping{}).
			Select("tax_code").
			Where("operation_id = ?", csrOprItemId).
			Scan(&pphTaxCode).Error; err != nil {
			return entity, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Pph_Tax_Code",
				Err:        err,
			}
		}

		if err := tx.Model(&masteritementities.Item{}).
			Select("atpm_warranty_claim_type_id").
			Where("item_id = ?", csrOprItemId).
			Scan(&atpmWcfType).Error; err != nil {
			return entity, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch ATPM_WCF_Type",
				Err:        err,
			}
		}

		// Handle different line types (Operation or Package)
		if csrLineTypeId == utils.LinetypeOperation || csrLineTypeId == utils.LinetypePackage {
			// Check if this operation item already exists
			var workOrder2 transactionworkshopentities.WorkOrderDetail

			if err := tx.Where("work_order_system_number = ? AND operation_item_id = ?", workOrderId, csrOprItemId).
				First(&workOrder2).Error; err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check if operation item exists",
					Err:        err,
				}
			}

			if workOrder2.WorkOrderOperationItemLine == 0 {
				// Insert new item if it doesn't exist
				var woOprItemLine int

				if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Select("ISNULL(MAX(work_order_operation_item_line), 0) + 1").
					Where("work_order_system_number = ?", workOrderId).
					Scan(&woOprItemLine).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get next work order operation item line",
						Err:        err,
					}
				}

				workOrder2 := transactionworkshopentities.WorkOrderDetail{
					WorkOrderSystemNumber:               workOrderId,
					WorkOrderOperationItemLine:          woOprItemLine,
					WorkorderStatusId:                   utils.WoStatNew,
					LineTypeId:                          csrLineTypeId,
					OperationItemId:                     csrOprItemId,
					TransactionTypeId:                   csrTrxTypeId,
					JobTypeId:                           csrJobTypeId,
					FrtQuantity:                         csrFrtQty,
					OperationItemPrice:                  oprItemPrice,
					OperationItemDiscountAmount:         oprItemDiscAmount,
					OperationItemDiscountPercent:        oprItemDiscPercent,
					OperationItemDiscountRequestAmount:  0,
					OperationItemDiscountRequestPercent: 0,
					PphAmount:                           0,
					SupplyQuantity:                      csrFrtQty,
					WarehouseGroupId:                    whsGroup,
					AtpmWCFTypeId:                       0,
				}

				if err := tx.Create(&workOrder2).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to insert new operation item",
						Err:        err,
					}
				}

				// Update estimation time if needed
				var estTime float64
				var timeTolerance float64 = 0.25

				// Step 1: Check if estTime is zero
				if estTime == 0 {
					// Calculate estimation time using FRT quantities plus the time tolerance
					if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
						Select("SUM(COALESCE(frt_quantity, 0))").
						Where("work_order_system_number = ? AND (line_type_id = ? OR line_type_id = ?)", workOrderId, utils.LinetypeOperation, utils.LinetypePackage).
						Scan(&estTime).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to calculate estimation time",
							Err:        err,
						}
					}
					// Add timeTolerance after summing the frt_quantity
					estTime += timeTolerance
				} else {
					// Step 2: If estTime is not zero, just add csrFrtQty
					estTime += csrFrtQty
				}

				if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
					Where("work_order_system_number = ?", workOrderId).
					Updates(map[string]interface{}{
						"estimate_time": estTime,
					}).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update estimation time",
						Err:        err,
					}
				}
			}

		} else {
			// Handle other line types
			uomType := "UOM_TYPE_SELL" // get variable value for UOM_TYPE_SELL

			// Fetch line type by item code
			getLineTypeByItemCode, err := s.lookupRepo.GetLineTypeByItemCode(tx, csrOprItemCode)
			if err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch line type by item code",
					Err:        errors.New("failed to fetch line type by item code"),
				}
			}

			// Check if the line type matches the one retrieved by item code
			if csrLineTypeId != getLineTypeByItemCode {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "Item code does not belong to the provided line type",
					Err:        errors.New("item code does not belong to line type"),
				}
			}

			// Fetch warehouse group if not provided
			if whsGroup == 0 {
				if err := tx.Model(&masteritementities.Item{}).
					Select("whs_group").
					Where("item_code = ? AND company_code = ?", csrOprItemId, result.CompanyCode).
					Limit(1).
					Scan(&whsGroup).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to fetch warehouse group",
						Err:        err,
					}
				}
			}

			var currentDate = time.Now()

			// Execute the stored procedure `SelectLocationStockItem`
			qtyAvailable, err := s.lookupRepo.SelectLocationStockItem(tx, 1, result.CompanyCode, currentDate, 0, "", csrOprItemId, whsGroup, uomType)
			if err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to execute stock item location query",
					Err:        errors.New("failed to execute stock item location query"),
				}
			}

			// Check if the item is available or the line type is Sublet
			if qtyAvailable > 0 || csrLineTypeId == utils.LinetypeSublet {
				// Fetch the next work order operation item line
				var woOprItemLine int
				if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Select("COALESCE(MAX(work_order_operation_item_line), 0) + 1").
					Where("work_order_system_number = ?", workOrderId).
					Scan(&woOprItemLine).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get the next work order operation item line",
						Err:        err,
					}
				}

				// Check if the work order line exists
				var exists bool

				// Check if the work order line already exists using Exists
				if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Select("1").
					Where("work_order_system_number = ? AND work_order_operation_item_line = ?", workOrderId, woOprItemLine).
					Scan(&exists).Error; err != nil {
					// Handle error if the query fails
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to check if work order line exists",
						Err:        err,
					}
				}

				if !exists {
					// If the record doesn't exist, proceed with insertion
					workOrderDetail := transactionworkshopentities.WorkOrderDetail{
						WorkOrderSystemNumber:        workOrderId,
						WorkOrderOperationItemLine:   woOprItemLine,
						LineTypeId:                   csrLineTypeId,
						OperationItemId:              csrOprItemId,
						FrtQuantity:                  csrFrtQty,
						OperationItemPrice:           oprItemPrice,
						OperationItemDiscountAmount:  oprItemDiscAmount,
						OperationItemDiscountPercent: oprItemDiscPercent,
						SupplyQuantity:               csrFrtQty,
						WarehouseGroupId:             whsGroup,
						AtpmWCFTypeId:                atpmWcfType,
					}

					// Insert the new operation item into the database
					if err := tx.Create(&workOrderDetail).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to insert new operation item",
							Err:        err,
						}
					}
				} else {

					var existingWorkOrderLine transactionworkshopentities.WorkOrderDetail
					if err := tx.Where("work_order_system_number = ? AND work_order_operation_item_line = ?", workOrderId, woOprItemLine).
						First(&existingWorkOrderLine).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to fetch existing work order line",
							Err:        err,
						}
					}
					// Update the existing work order line
					newFrtQty := existingWorkOrderLine.FrtQuantity + csrFrtQty
					newWhsGroup := existingWorkOrderLine.WarehouseGroupId

					// Ambil markup berdasarkan Company dan Vehicle Brand
					// OLDSITE
					// if err := tx.Model(&GmSiteMarkup{}).
					// 	Where("company_id = ? AND brand_id = ? AND site_code = ? AND trx_type_id = ?", result.CompanyCode, result.BrandId, siteCode, billCodeExt).
					// 	Select("markup_amount, markup_percentage").
					// 	Scan(&markupAmount, &markupPercentage).Error; err != nil {
					// 	return entity, &exceptions.BaseErrorResponse{
					// 		StatusCode: http.StatusInternalServerError,
					// 		Message:    "Failed to fetch markup",
					// 		Err:        err,
					// 	}
					// }
					// NEW SITE
					// Fetch markup based on Company and Vehicle Brand
					markupAmount := 0.0
					markupPercentage := 0.0

					// Fetch operation item price and discount percent
					oprItemPrice, err := s.lookupRepo.GetOprItemPrice(tx, result.CompanyCode, result.BrandId, csrOprItemId, result.AgreementNo, result.JobTypeId, csrLineTypeId, csrTrxTypeId, int(csrFrtQty), newWhsGroup, strconv.Itoa(result.VariantId))
					if err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to fetch operation item price",
							Err:        errors.New("failed to fetch operation item price"),
						}
					}

					// Apply markup and percentage increase
					oprItemPrice += markupAmount + (oprItemPrice * (markupPercentage / 100))

					// Fetch item discount percent (You need to define or fix GetOprItemDiscPercent method)
					oprItemDiscPercent, err := s.lookupRepo.GetOprItemDisc(tx, result.LinetypeId, result.BillCodeExt, csrOprItemId, result.AgreementNo, result.CpcCode, oprItemPrice*newFrtQty, result.CompanyCode, result.BrandId, 0, whsGroup, 0)
					if err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to fetch operation item discount percent",
							Err:        errors.New("failed to fetch operation item discount percent"),
						}
					}

					// Calculate discount amount
					oprItemDiscAmount := math.Round((oprItemPrice * oprItemDiscPercent / 100))

					// Update the work order line
					if err := tx.Model(&existingWorkOrderLine).
						Updates(map[string]interface{}{
							"frt_quantity":                    newFrtQty,
							"operation_item_price":            oprItemPrice,
							"operation_item_discount_amount":  oprItemDiscAmount,
							"operation_item_discount_percent": oprItemDiscPercent,
							"approval_id":                     10,
						}).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to update existing operation item",
							Err:        err,
						}
					}
				}

			} else {

				type Substitute struct {
					SubstituteItemCode string  // Substitute item code
					ItemName           string  // Item name
					SupplyQty          float64 // Supply quantity
					SubstituteType     string  // Substitute type
				}

				// Step 1: Create a temporary table with GORM
				if err := tx.Exec("CREATE TABLE #SUBSTITUTE2 (SUBS_ITEM_CODE VARCHAR(15), ITEM_NAME VARCHAR(40), SUPPLY_QTY NUMERIC(7,2), SUBS_TYPE CHAR(2))").Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to create temporary table",
						Err:        err,
					}
				}

				// Step 2: Execute stored procedure to populate the temporary table
				if err := tx.Exec("INSERT INTO #SUBSTITUTE2 EXEC dbo.uspg_smSubstitute0_Select @OPTION = ?, @COMPANY_CODE = ?, @ITEM_CODE = ?, @ITEM_QTY = ?",
					2, result.CompanyCode, csrOprItemId, csrFrtQty).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to execute stored procedure",
						Err:        err,
					}
				}

				// Step 3: Fetch data from the temporary table
				var substitutes []Substitute
				if err := tx.Table("#SUBSTITUTE2").Find(&substitutes).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to fetch data from temporary table",
						Err:        err,
					}
				}

				// Step 4: Process substitution items
				for _, substitute := range substitutes {
					if substitute.SubstituteType != "" {
						var woOprItemLinesub int

						// Check if item already exists
						var existingItemCount int64
						if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
							Where("work_order_system_number = ? AND transaction_type_id = ? AND operation_item_id = ? AND substitute_type_id = ?", workOrderId, csrTrxTypeId, substitute.SubstituteItemCode, substitute.SubstituteType).
							Count(&existingItemCount).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to check if item exists",
								Err:        err,
							}
						}

						if existingItemCount == 0 {
							// Get the next line number
							if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
								Select("ISNULL(MAX(work_order_operation_item_line), 0) + 1").
								Where("work_order_system_number = ?", workOrderId).
								Scan(&woOprItemLinesub).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to get the next work order operation item line",
									Err:        err,
								}
							}

							// Step 5: Insert new data into the work order detail table
							workOrderDetail := transactionworkshopentities.WorkOrderDetail{
								WorkOrderSystemNumber:        workOrderId,
								WorkOrderOperationItemLine:   woOprItemLinesub,
								LineTypeId:                   csrLineTypeId,
								FrtQuantity:                  substitute.SupplyQty,
								OperationItemPrice:           0,
								OperationItemDiscountAmount:  0,
								OperationItemDiscountPercent: 0,
								SupplyQuantity:               substitute.SupplyQty,
								WarehouseGroupId:             whsGroup,
								AtpmWCFTypeId:                0,
							}

							if err := tx.Create(&workOrderDetail).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to insert new operation item",
									Err:        err,
								}
							}
						}

						// Step 6: Fetch markup based on Company and Vehicle Brand
						var markupAmount, markupPercentage float64

						// Fetch operation item price and discount percent
						oprItemPrice, err := s.lookupRepo.GetOprItemPrice(tx, result.CompanyCode, result.BrandId, csrOprItemId, result.AgreementNo, result.JobTypeId, csrLineTypeId, csrTrxTypeId, int(csrFrtQty), whsGroup, strconv.Itoa(result.VariantId))
						if err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to fetch operation item price",
								Err:        errors.New("failed to fetch operation item price"),
							}
						}

						// Apply markup and percentage increase
						oprItemPrice += markupAmount + (oprItemPrice * (markupPercentage / 100))

						// Fetch item discount percent
						oprItemDiscPercent, err := s.lookupRepo.GetOprItemDisc(tx, result.LinetypeId, result.BillCodeExt, csrOprItemId, result.AgreementNo, result.CpcCode, oprItemPrice*substitute.SupplyQty, result.CompanyCode, result.BrandId, 0, whsGroup, 0)
						if err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to fetch operation item discount percent",
								Err:        errors.New("failed to fetch operation item discount percent"),
							}
						}

						// Calculate discount amount
						oprItemDiscAmount := math.Round((oprItemPrice * oprItemDiscPercent / 100))

						// Get the next work order operation item line
						if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
							Select("COALESCE(MAX(work_order_operation_item_line), 0) + 1").
							Where("work_order_system_number = ?", workOrderId).
							Scan(&woOprItemLinesub).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to get the next work order operation item line",
								Err:        err,
							}
						}

						// Step 7: Process supply quantity based on substitute type
						var supplyQty float64
						if substitute.SubstituteType == "I" {
							// Fetch quantity from smSubsStockInterchange1
							if err := tx.Model(&masteritementities.Item{}).
								Select("COALESCE(SUM(COALESCE(qty, 0)), 0)").
								Where("item_code = ?", substitute.SubstituteItemCode).
								Scan(&supplyQty).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to fetch quantity from smSubsStockInterchange1",
									Err:        err,
								}
							}
						} else {
							supplyQty = substitute.SupplyQty
						}

						// Check item existence in item master
						var itemExists int64

						if err := tx.Model(&masteritementities.Item{}).
							Where("item_code = ? and item_group_id <> ? and item_type_id = ?", csrOprItemCode, 1, 2).
							Count(&itemExists).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to check if item exists in item master",
								Err:        err,
							}
						}

						if itemExists == 0 {
							supplyQty = csrFrtQty
						} else {
							supplyQty = 0
						}

						if substitute.SubstituteType == "" {
							substitute.SupplyQty = csrFrtQty
						}

						// Set ATPM warranty claim type ID
						var atpmWcfTypeId int

						if err := tx.Model(&masteritementities.Item{}).
							Select("atpm_warranty_claim_type_id").
							Where("item_code = ?", substitute.SubstituteItemCode).
							Scan(&atpmWcfTypeId).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to fetch ATPM_WCF_Type",
								Err:        err,
							}
						}

						// Step 8: Insert or update the final operation item
						finalWorkOrderDetail := transactionworkshopentities.WorkOrderDetail{
							WorkOrderSystemNumber:        workOrderId,
							WorkOrderOperationItemLine:   woOprItemLinesub,
							LineTypeId:                   csrLineTypeId,
							FrtQuantity:                  supplyQty,
							OperationItemPrice:           oprItemPrice,
							OperationItemDiscountAmount:  oprItemDiscAmount,
							OperationItemDiscountPercent: oprItemDiscPercent,
							SupplyQuantity:               supplyQty,
							WarehouseGroupId:             whsGroup,
							AtpmWCFTypeId:                atpmWcfTypeId,
						}

						if err := tx.Create(&finalWorkOrderDetail).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to insert or update operation item",
								Err:        err,
							}
						}

						// Step 9: Calculate totals for various line types and update work order
						var (
							totalPart               float64
							totalOil                float64
							totalMaterial           float64
							totalConsumableMaterial float64
							totalSublet             float64
							totalAccs               float64
						)

						if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
							Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
							Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSparepart).
							Scan(&totalPart).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to calculate total part",
								Err:        err,
							}
						}

						if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
							Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
							Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOil).
							Scan(&totalOil).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to calculate total oil",
								Err:        err,
							}
						}

						if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
							Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
							Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeMaterial).
							Scan(&totalMaterial).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to calculate total material",
								Err:        err,
							}
						}

						if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
							Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
							Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeConsumableMaterial).
							Scan(&totalConsumableMaterial).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to calculate total consumable material",
								Err:        err,
							}
						}

						if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
							Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
							Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSublet).
							Scan(&totalSublet).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to calculate total sublet",
								Err:        err,
							}
						}

						if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
							Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
							Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeAccesories).
							Scan(&totalAccs).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to calculate total accessories",
								Err:        err,
							}
						}

						// Update the work order with the new totals
						if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
							Where("work_order_system_number = ?", workOrderId).
							Updates(map[string]interface{}{
								"total_part":                totalPart,
								"total_oil":                 totalOil,
								"total_material":            totalMaterial,
								"total_consumable_material": totalConsumableMaterial,
								"total_sublet":              totalSublet,
								"total_accessories":         totalAccs,
							}).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to update work order with new totals",
								Err:        err,
							}
						}
					}
				}

				if err := tx.Exec("DROP TABLE #SUBSTITUTE2").Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to drop temporary table",
						Err:        err,
					}
				}
			}
		}
	}

	return entity, nil
}

// uspg_wtWorkOrder2_Insert
// IF @Option = 3
func (s *WorkOrderRepositoryImpl) AddFieldAction(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.WorkOrderFieldActionRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderDetail
	var result transactionworkshopentities.WorkOrder
	var recdettype transactionworkshopentities.RecallDetailType

	var companyCode, servMileage, vehicleBrand, modelCode, variantCode, woStatus, agreementNoBR, agreementNoGR, currencyCode int
	var vehicleChassisNo, cpcCode string

	if err := tx.Table("trx_work_order W").
		Select(`
            W.company_id, W.vehicle_chassis_number, W.brand_id, 
            W.model_id, W.variant_id, W.service_mileage, W.cpc_code,
            W.work_order_status_id, W.agreement_body_repair_id, W.agreement_general_repair_id,
            W.currency_id`).
		Where("work_order_system_number = ?", workOrderId).
		Row().Scan(&companyCode, &vehicleChassisNo, &vehicleBrand, &modelCode, &variantCode, &servMileage, &cpcCode, &woStatus, &agreementNoBR, &agreementNoGR, &currencyCode); err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error fetching Work Order data",
			Err:        err,
		}
	}

	// Recall validation
	var recallExists bool
	if err := tx.Table("trx_recall_detail A").
		Joins("INNER JOIN trx_recall B ON A.recall_system_number = B.recall_system_number").
		Where("A.has_recall = 0 AND A.vechicle_chassis_number = ? AND B.recall_document_number = ?", vehicleChassisNo, request.RecallNo).
		Scan(&recallExists).Error; err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error checking recall eligibility",
			Err:        err,
		}
	}

	// If recall is not eligible
	if !recallExists {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Vehicle is not eligible for Field Action",
			Err:        errors.New("vehicle is not eligible for Field Action"),
		}
	} else {

		// Step 1: Query the count of recall data
		var count int
		if err := tx.Table("trx_recall_detail A").
			Select("COUNT(*)").
			Joins("INNER JOIN trx_recall B ON B.recall_system_number = A.recall_system_number").
			Joins("INNER JOIN trx_recall_detail_type C ON C.recall_system_number = A.recall_system_number AND C.recall_line_number = A.recall_line_number").
			Where("A.vechicle_chassis_number = ? AND B.recall_document_number = ? AND C.has_recall = 0", vehicleChassisNo, request.RecallNo).
			Scan(&count).Error; err != nil {
			return entity, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error counting recall data",
				Err:        err,
			}
		}

		// Step 2: Fetch recall data and process each record
		var recallRecords []struct {
			LineTypeId  int
			OprItemId   int
			OprItemCode string
			Description string
			FrtQty      float64
		}

		if err := tx.Table("trx_recall_detail A").
			Select(`
			C.recall_type_id AS recall_type_id,
			C.operation_item_code AS operation_item_code,
			CASE C.recall_type_id WHEN '1' THEN E.operation_name ELSE D.item_name END AS item_name,
			C.frt_qty AS frt_qty`).
			Joins("INNER JOIN trx_recall B ON B.recall_system_number = A.recall_system_number").
			Joins("INNER JOIN trx_recall_detail_type C ON C.recall_system_number = A.recall_system_number AND C.recall_line_number = A.recall_line_number").
			Joins("LEFT JOIN mtr_item D ON C.operation_item_id = D.item_id").
			Joins("LEFT JOIN mtr_operation_code E ON C.operation_item_id = E.operation_id").
			Where("A.vechicle_chassis_number = ? AND B.recall_document_number = ? AND C.has_recall = 0", vehicleChassisNo, request.RecallNo).
			Find(&recallRecords).Error; err != nil {
			return entity, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error fetching recall data",
				Err:        err,
			}
		}

		// Step 3: Process each recall record
		for _, recallRecord := range recallRecords {

			billCode := utils.TrxTypeWoWarranty.Code
			jobTypeId := utils.JobTypeWarranty.ID
			WhsGroup, err := s.lookupRepo.GetWhsGroup(tx, companyCode)
			if err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error fetching Warehouse Group",
					Err:        errors.New("error fetching Warehouse Group"),
				}
			}

			var agreementNo int
			if recallRecord.LineTypeId == utils.LinetypeOperation {
				agreementNo = agreementNoBR
			} else {
				agreementNo = agreementNoGR
			}

			oprItemPrice, err := s.lookupRepo.GetOprItemPrice(
				tx, companyCode, vehicleBrand, recallRecord.OprItemId, agreementNo,
				entity.JobTypeId, recallRecord.LineTypeId, utils.TrxTypeWoWarranty.ID,
				int(recallRecord.FrtQty), WhsGroup, strconv.Itoa(variantCode),
			)
			if err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch operation item price",
					Err:        errors.New("failed to fetch operation item price"),
				}
			}

			oprItemDiscPercent := 0.0
			oprItemDiscReqPercent := 0.0
			oprItemDiscAmount := 0.0
			oprItemDiscReqAmount := 0.0
			substituteItemCode := ""

			var itemUom string
			if err := tx.Table("mtr_item").
				Select("uom_type").
				Where("item_id = ?", recallRecord.OprItemId).
				Row().Scan(&itemUom); err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error fetching item UOM",
					Err:        err,
				}
			}

			switch recallRecord.LineTypeId {
			case utils.LinetypeOperation:
				var operationExists bool
				if err := tx.Table("mtr_operation_code").
					Select("1").
					Where("operation_id = ?", recallRecord.OprItemId).
					Scan(&operationExists).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Error checking Operation existence in Master",
						Err:        err,
					}
				} else if !operationExists {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "Operation does not exist in Master",
						Err:        errors.New("operation does not exist in Master"),
					}
				}

				frtQty, err := s.lookupRepo.GetOprItemFrt(tx, recallRecord.OprItemId, vehicleBrand, modelCode, variantCode, vehicleChassisNo)
				if err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Error fetching Frt_Qty and Supply_Qty",
						Err:        errors.New("error fetching Frt_Qty and Supply_Qty"),
					}
				}
				entity.FrtQuantity = frtQty

			case utils.LinetypePackage:
				entity.FrtQuantity = 1
				entity.SupplyQuantity = 1
				entity.WarehouseGroupId = 0

				var packageData struct {
					PackagePrice float64
					PphTaxCode   float64
				}

				if err := tx.Table("mtr_package").
					Select("package_price, pph_tax_code").
					Where("package_id = ?", recallRecord.OprItemId).
					Scan(&packageData).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Error fetching Package data",
						Err:        err,
					}
				}

				entity.OperationItemPrice = packageData.PackagePrice
				entity.PphTaxRate = packageData.PphTaxCode

			default:
				var itemCodeExists bool
				if err := tx.Table("mtr_item").
					Select("1").
					Where("item_id = ?", recallRecord.OprItemId).
					Scan(&itemCodeExists).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Error checking Item existence in Master",
						Err:        err,
					}
				} else if !itemCodeExists {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "Item does not exist in Master",
						Err:        errors.New("item does not exist in Master"),
					}
				}

				lineTypeId, err := s.lookupRepo.GetLineTypeByItemCode(tx, recallRecord.OprItemCode)
				if err != nil {
					return entity, err
				}

				if lineTypeId != recallRecord.LineTypeId {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "Line Type does not match for Item",
						Err:        errors.New("line Type does not match for Item"),
					}
				}

				var itemExists int64
				if err := tx.Table("mtr_item_detail").
					Where("item_code = ? AND brand_id = ? AND model_id = ?", recallRecord.OprItemCode, vehicleBrand, modelCode).
					Count(&itemExists).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Error checking Item existence in Detail",
						Err:        err,
					}
				} else if itemExists == 0 {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "Item cannot be used for current Model and Variant",
						Err:        errors.New("item cannot be used for current Model and Variant"),
					}
				}
			}

			var vehicleAgeMnth int
			var warrantyExp struct {
				ExpireMonth   int
				ExpireMileage int
			}
			var fsExp struct {
				ExpireMonth   int
				ExpireMileage int
			}
			var frtPackage float64
			var atpmWcfType string

			// Calculate Vehicle Age in months
			if err := tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle_registration_certificate").
				Joins("INNER JOIN dms_microservices_sales_dev.dbo.mtr_vehicle ON mtr_vehicle_registration_certificate.vehicle_id = mtr_vehicle.vehicle_id").
				Select("DATEDIFF(month, vehicle_bpk_date, GETDATE()) AS vehicle_age_month").
				Where("mtr_vehicle.vehicle_chassis_number = ?", vehicleChassisNo).
				Scan(&vehicleAgeMnth).Error; err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error calculating Vehicle Age",
					Err:        err,
				}
			}

			if err := tx.Table("mtr_warranty_free_service B").
				Select("B.expire_month, B.expire_mileage").
				Where("B.model_id = ? AND B.brand_id = ? AND B.warranty_free_service_type_id = 'W'", modelCode, vehicleBrand).
				Where("B.effective_date = (SELECT TOP 1 A.effective_date FROM mtr_warranty_free_service A WHERE A.model_id = B.model_id AND A.brand_id = B.brand_id AND A.warranty_free_service_type_id = B.warranty_free_service_type_id ORDER BY A.effective_date DESC)").
				Scan(&warrantyExp).Error; err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error fetching Warranty Expiry details",
					Err:        err,
				}
			}

			if err := tx.Table("mtr_warranty_free_service B").
				Select("B.expire_month, B.expire_mileage").
				Where("B.model_id = ? AND B.brand_id = ? AND B.warranty_free_service_type_id = 'F'", modelCode, vehicleBrand).
				Where("B.effective_date = (SELECT TOP 1 A.effective_date FROM mtr_warranty_free_service A WHERE A.model_id = B.model_id AND A.brand_id = B.brand_id AND A.warranty_free_service_type_id = B.warranty_free_service_type_id ORDER BY A.effective_date DESC)").
				Scan(&fsExp).Error; err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error fetching Free Service Expiry details",
					Err:        err,
				}
			}

			if recallRecord.LineTypeId == utils.LinetypePackage {
				if err := tx.Model(&masterentities.PackageMasterDetail{}).
					Select("SUM(frt_quantity)").
					Where("package_id = ?", recallRecord.OprItemId).
					Scan(&frtPackage).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Error fetching FRT Package",
						Err:        err,
					}
				}
			} else {
				frtPackage = 0
			}

			if billCode == utils.TrxTypeWoFreeService.Code {
				if vehicleAgeMnth > fsExp.ExpireMonth {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "Free Service for this vehicle is already expired",
						Err:        errors.New("vehicle's Free Service expired"),
					}
				}
				if servMileage > fsExp.ExpireMileage {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "Free Service for this vehicle is already expired",
						Err:        errors.New("vehicle's Free Service expired"),
					}
				}
			}

			if err := tx.Model(&masteritementities.Item{}).
				Select("atpm_warranty_claim_type_id").
				Where("item_id = ?", recallRecord.OprItemId).
				Scan(&atpmWcfType).Error; err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error fetching ATPM WCF Type",
					Err:        err,
				}
			}

			var woOprItemLine int
			var totalPackage, totalOpr, estTime, totalFrtPackage, totalFrt float64

			if recallRecord.LineTypeId == utils.LinetypeOperation || recallRecord.LineTypeId == utils.LinetypePackage {

				var exists int64
				if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Select("work_order_operation_item_line").
					Where("work_order_system_number = ? AND operation_item_code = ?", workOrderId, recallRecord.OprItemCode).
					Count(&exists).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Error checking existing operation/package",
						Err:        err,
					}
				}

				if exists == 0 {
					if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
						Select("ISNULL(MAX(work_order_operation_item_line), 0) + 1").
						Where("work_order_system_number = ?", workOrderId).
						Scan(&woOprItemLine).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Error calculating WO_OPR_ITEM_LINE",
							Err:        err,
						}
					}

					if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{
						WorkOrderSystemNumber:               workOrderId,
						WorkOrderOperationItemLine:          woOprItemLine,
						LineTypeId:                          recallRecord.LineTypeId,
						JobTypeId:                           jobTypeId,
						OperationItemId:                     recallRecord.OprItemId,
						OperationItemCode:                   recallRecord.OprItemCode,
						Description:                         recallRecord.Description,
						FrtQuantity:                         recallRecord.FrtQty,
						OperationItemPrice:                  oprItemPrice,
						OperationItemDiscountPercent:        oprItemDiscPercent,
						OperationItemDiscountAmount:         oprItemDiscAmount,
						OperationItemDiscountRequestAmount:  oprItemDiscReqAmount,
						OperationItemDiscountRequestPercent: oprItemDiscReqPercent,
						SupplyQuantity:                      entity.SupplyQuantity,
						WarehouseGroupId:                    entity.WarehouseGroupId,
						SubstituteTypeId:                    0,
						SubstituteItemCode:                  substituteItemCode,
					}).Create(&entity).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Error inserting into wtWorkOrder2",
							Err:        err,
						}
					}

					if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
						Select("ISNULL(SUM(ROUND(ISNULL(operation_item_price, 0), 0, 0)), 0)").
						Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypePackage).
						Scan(&totalPackage).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Error calculating total package",
							Err:        err,
						}
					}

					if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
						Select("ISNULL(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
						Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOperation).
						Scan(&totalOpr).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Error calculating total operation",
							Err:        err,
						}
					}

					var currentEstTime float64
					if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
						Select("ISNULL(estimate_time, 0)").
						Where("work_order_system_number = ?", workOrderId).
						Scan(&currentEstTime).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Error fetching EST_TIME",
							Err:        err,
						}
					}

					if currentEstTime == 0 {
						estTime = 0.25
						if err := tx.Table("trx_work_order_detail WO2").
							Joins("INNER JOIN mtr_package_detail P1 ON WO2.operation_item_code = P1.package_code").
							Select("ISNULL(SUM(ISNULL(P1.frt_quantity, 0)), 0)").
							Where("WO2.work_order_system_number = ? AND WO2.line_type_id = ?", workOrderId, utils.LinetypePackage).
							Scan(&totalFrtPackage).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Error calculating total FRT package",
								Err:        err,
							}
						}

						if err := tx.Table("trx_work_order_detail").
							Select("ISNULL(SUM(ISNULL(frt_quantity, 0)), 0)").
							Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOperation).
							Scan(&totalFrt).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Error calculating total FRT",
								Err:        err,
							}
						}

						estTime += totalFrt + totalFrtPackage
					} else {
						if recallRecord.LineTypeId == utils.LinetypePackage {
							estTime = currentEstTime + frtPackage
						} else {
							estTime = currentEstTime + recallRecord.FrtQty
						}
					}

					// Check if Wo_Stat3 equals Wo_StatQcPass3 and if any records exist for the given conditions
					var exists int64
					var woStat3, productionHead int
					var notes, suggestion, fsCouponNo string
					var promiseDate time.Time
					var promiseTime int

					if woStat3 == utils.WoStatQC {
						if err := tx.Table("trx_work_order_detail").
							Where("work_order_system_number = ? AND (line_type_id = ? OR line_type_id = ?)", workOrderId, utils.LinetypeOperation, utils.LinetypePackage).
							Count(&exists).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Error checking record existence in wtWorkOrder2",
								Err:        err,
							}
						}

						if exists > 0 {
							if cpcCode == "00002" {
								var countOprStatus int64
								if err := tx.Table("trx_work_order_detail").
									Where("work_order_system_number = ? AND line_type_id = ? AND ISNULL(work_order_status_id, '') <> ?", workOrderId, utils.LinetypeOperation, utils.SrvStatQcPass).
									Count(&countOprStatus).Error; err != nil {
									return entity, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Error checking WO_OPR_STATUS in wtWorkOrder2",
										Err:        err,
									}
								}

								if countOprStatus == 0 {
									woStat3 = utils.WoStatQC
								} else {
									woStat3 = utils.WoStatOngoing
								}
							} else {
								var countPkgStatus int64
								if err := tx.Table("trx_work_order_detail").
									Where("work_order_system_number = ? AND line_type_id = ? AND ISNULL(work_order_status_id, '') <> ?", workOrderId, utils.LinetypePackage, utils.SrvStatQcPass).
									Count(&countPkgStatus).Error; err != nil {
									return entity, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Error checking WO_OPR_STATUS in wtWorkOrder2 for Package",
										Err:        err,
									}
								}

								if countPkgStatus == 0 {
									woStat3 = utils.WoStatQC
								} else {
									woStat3 = utils.WoStatOngoing
								}
							}
						}
					}

					if err := tx.Table("trx_work_order").
						Where("work_order_system_number = ?", workOrderId).
						Updates(map[string]interface{}{
							"work_order_status_id": woStat3,
							"notes":                notes,
							"suggestion":           suggestion,
							"fs_coupon_number":     fsCouponNo,
							"production_head_id":   productionHead,
							"total_operation":      totalOpr,
							"total_package":        totalPackage,
							"estimate_time":        estTime,
							"PROMISE_DATE":         gorm.Expr("CASE WHEN ISNULL(?, '') = '' THEN PROMISE_DATE ELSE ? END", promiseDate, promiseDate),
							"PROMISE_TIME":         gorm.Expr("CASE WHEN ISNULL(?, 0) = 0 THEN PROMISE_TIME ELSE ? END", promiseTime, promiseTime),
						}).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Error updating wtWorkOrder0",
							Err:        err,
						}
					}
				} else {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "Operation or Package already exists",
						Err:        errors.New("operation or package already exists"),
					}
				}
			} else {
				//-- LINE TYPE <> 1 , NEED SUBSTITUTE
				linetypecode, err := s.lookupRepo.GetLineTypeByItemCode(tx, recallRecord.OprItemCode)
				if err != nil {
					return entity, err
				}

				if recallRecord.LineTypeId != linetypecode {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "Item Code belong to Line Type",
						Err:        errors.New("item code belong to line type"),
					}
				}

				UomType := utils.UomTypeService

				var qtyAvail float64
				qtyAvail, errResponse := s.lookupRepo.SelectLocationStockItem(tx, 1, companyCode, time.Now(), 0, "", recallRecord.OprItemId, WhsGroup, UomType)
				if errResponse != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to check quantity available",
						Err:        errResponse.Err,
					}
				}

				// Assuming necessary imports and definitions are in place
				if qtyAvail > 0 || recallRecord.LineTypeId == utils.LinetypeSublet {
					// Get the next available WO_OPR_ITEM_LINE
					var woOprItemLine int
					if err := tx.Table("trx_work_order_detail").
						Select("ISNULL(MAX(work_order_operation_item_line), 0) + 1").
						Where("work_order_system_number = ?", workOrderId).Scan(&woOprItemLine).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Error fetching WO_OPR_ITEM_LINE",
							Err:        err,
						}
					}

					// Check if the WO_OPR_ITEM_LINE already exists
					var exists bool
					if err := tx.Table("trx_work_order_detail").
						Select("EXISTS(SELECT 1 FROM trx_work_order_detail WHERE work_order_system_number = ? AND work_order_operation_item_line = ?) AS exists", workOrderId, woOprItemLine).
						Scan(&exists).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Error checking WO_OPR_ITEM_LINE existence",
							Err:        err,
						}
					}

					if !exists {
						// Check if the OPR_ITEM_CODE and BILL_CODE already exist
						if err := tx.Table("trx_work_order_detail").
							Select("EXISTS(SELECT 1 FROM trx_work_order_detail WHERE work_order_system_number = ? AND operation_item_code = ? AND transaction_type_id = ?) AS exists", workOrderId, recallRecord.OprItemCode, billCode).Scan(&exists).Error; err != nil {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Error checking OPR_ITEM_CODE existence",
								Err:        err,
							}
						}

						if !exists {
							// Check if ITEM_CODE exists in gmItem0
							var itemCodeExists bool
							itemGrpOJ := 6 // OJ
							if err := tx.Table("mtr_item").
								Select("EXISTS(SELECT 1 FROM mtr_item WHERE item_group_id <> ? AND item_code = ? AND item_type_id = ?) AS exists", itemGrpOJ, recallRecord.OprItemCode, 2).
								Scan(&itemCodeExists).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Error checking ITEM_CODE in gmItem0",
									Err:        err,
								}
							}

							var supplyQty float64
							if itemCodeExists {
								supplyQty = entity.FrtQuantity
							} else {
								supplyQty = 0
							}

							// Prepare the new record to insert into wtWorkOrder2
							newRecord := &transactionworkshopentities.WorkOrderDetail{
								WorkOrderSystemNumber:               workOrderId,
								WorkOrderOperationItemLine:          woOprItemLine,
								WorkorderStatusId:                   utils.WoStatNew,
								LineTypeId:                          utils.LinetypeSublet,
								ServiceStatusId:                     0,
								OperationItemCode:                   recallRecord.OprItemCode,
								Description:                         recallRecord.Description,
								FrtQuantity:                         recallRecord.FrtQty,
								OperationItemPrice:                  entity.OperationItemPrice,
								OperationItemDiscountPercent:        oprItemDiscPercent,
								OperationItemDiscountAmount:         oprItemDiscAmount,
								OperationItemDiscountRequestPercent: oprItemDiscReqPercent,
								OperationItemDiscountRequestAmount:  oprItemDiscReqAmount,
								PphAmount:                           0,
								PphTaxRate:                          0,
								SupplyQuantity:                      supplyQty,
								SubstituteItemCode:                  substituteItemCode,
								WarehouseGroupId:                    WhsGroup,
								RecSystemNumber:                     entity.RecSystemNumber,
								ServiceCategoryId:                   entity.ServiceCategoryId,
							}

							// Insert the new record into wtWorkOrder2
							if err := tx.Table("trx_work_order_detail").
								Create(newRecord).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Error inserting into wtWorkOrder2",
									Err:        err,
								}
							}

							// Calculate totals
							var totalPart, totalOil, totalMaterial, totalConsumableMaterial, totalSublet, totalAccs float64
							var productionHead int
							var notes, suggestion, fsCouponNo string
							var promiseDate time.Time
							var promiseTime int

							totals := []struct {
								total      *float64
								lineTypeId int
							}{
								{&totalPart, utils.LinetypeSparepart},
								{&totalOil, utils.LinetypeOil},
								{&totalMaterial, utils.LinetypeMaterial},
								{&totalConsumableMaterial, utils.LinetypeConsumableMaterial},
								{&totalSublet, utils.LinetypeSublet},
								{&totalAccs, utils.LinetypeAccesories},
							}

							for _, t := range totals {
								if err := tx.Table("trx_work_order_detail").
									Select("ISNULL(SUM(ROUND(ISNULL(operation_item_price,0) * ISNULL(frt_quantity,0),0,0)), 0)").
									Where("work_order_system_number = ? AND line_type_id = ? AND substitute_type_id <> ?", workOrderId, t.lineTypeId, substituteItemCode).
									Scan(t.total).Error; err != nil {
									return entity, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Error calculating totals",
										Err:        err,
									}
								}
							}

							// Update wtWorkOrder0
							if err := tx.Table("trx_work_order").
								Where("work_order_system_number = ?", workOrderId).
								Updates(map[string]interface{}{
									"notes":                     notes,
									"suggestion":                suggestion,
									"fs_coupon_number":          fsCouponNo,
									"production_head_id":        productionHead,
									"total_part":                totalPart,
									"total_oil":                 totalOil,
									"total_material":            totalMaterial,
									"total_consumable_material": totalConsumableMaterial,
									"total_sublet":              totalSublet,
									"total_price_accessories":   totalAccs,
									"promise_date":              gorm.Expr("CASE WHEN ISNULL(?, '') = '' THEN PROMISE_DATE ELSE ? END", promiseDate, promiseDate),
									"promise_time":              gorm.Expr("CASE WHEN ISNULL(?, 0) = 0 THEN PROMISE_TIME ELSE ? END", promiseTime, promiseTime),
								}).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Error updating wtWorkOrder0",
									Err:        err,
								}
							}

						} else {
							return entity, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusConflict,
								Message:    "Data Item already exists",
								Err:        errors.New("data item already exists"),
							}
						}
					} else {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusConflict,
							Message:    "Data is already exists",
							Err:        errors.New("data is already exists"),
						}
					}
				} else { // --IF @QTY_AVAIL = 0
					type Substitute struct {
						SubstituteItemCode string  // Substitute item code
						ItemName           string  // Item name
						SupplyQty          float64 // Supply quantity
						SubstituteType     string  // Substitute type
					}

					// Step 1: Create a temporary table with GORM
					if err := tx.Exec("CREATE TABLE #SUBSTITUTE4 (SUBS_ITEM_CODE VARCHAR(15), ITEM_NAME VARCHAR(40), SUPPLY_QTY NUMERIC(7,2), SUBS_TYPE CHAR(2))").Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to create temporary table",
							Err:        err,
						}
					}

					// Step 2: Execute stored procedure to populate the temporary table
					if err := tx.Exec("INSERT INTO #SUBSTITUTE4 EXEC dbo.uspg_smSubstitute0_Select @OPTION = ?, @COMPANY_CODE = ?, @ITEM_CODE = ?, @ITEM_QTY = ?",
						2, result.CompanyId, recallRecord.OprItemCode, recallRecord.FrtQty).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to execute stored procedure",
							Err:        err,
						}
					}

					// Step 3: Fetch data from the temporary table
					var substitutes []Substitute
					if err := tx.Table("#SUBSTITUTE4").Find(&substitutes).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to fetch data from temporary table",
							Err:        err,
						}
					}

					// Step 4: Process substitution items
					for _, substitute := range substitutes {
						if substitute.SubstituteType != "" {
							var woOprItemLinesub int

							// Check if item already exists
							var existingItemCount int64
							if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
								Where("work_order_system_number = ? AND transaction_type_id = ? AND operation_item_id = ? AND substitute_type_id = ?", workOrderId, utils.TrxTypeWoFreeService.ID, substitute.SubstituteItemCode, substitute.SubstituteType).
								Count(&existingItemCount).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to check if item exists",
									Err:        err,
								}
							}

							if existingItemCount == 0 {
								// Get the next line number
								if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
									Select("ISNULL(MAX(work_order_operation_item_line), 0) + 1").
									Where("work_order_system_number = ?", workOrderId).
									Scan(&woOprItemLinesub).Error; err != nil {
									return entity, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Failed to get the next work order operation item line",
										Err:        err,
									}
								}

								// Step 5: Insert new data into the work order detail table
								workOrderDetail := transactionworkshopentities.WorkOrderDetail{
									WorkOrderSystemNumber:        workOrderId,
									WorkOrderOperationItemLine:   woOprItemLinesub,
									LineTypeId:                   utils.LinetypeSublet,
									FrtQuantity:                  substitute.SupplyQty,
									OperationItemPrice:           0,
									OperationItemDiscountAmount:  0,
									OperationItemDiscountPercent: 0,
									SupplyQuantity:               substitute.SupplyQty,
									WarehouseGroupId:             WhsGroup,
									AtpmWCFTypeId:                0,
								}

								if err := tx.Create(&workOrderDetail).Error; err != nil {
									return entity, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Failed to insert new operation item",
										Err:        err,
									}
								}
							}

							// Step 6: Fetch markup based on Company and Vehicle Brand
							var markupAmount, markupPercentage float64

							// Fetch operation item price and discount percent
							oprItemPrice, err := s.lookupRepo.GetOprItemPrice(tx, result.CompanyId, result.BrandId, recallRecord.OprItemId, agreementNo, utils.TrxTypeWoFreeService.ID, utils.LinetypeSublet, utils.TrxTypeWoFreeService.ID, int(recallRecord.FrtQty), WhsGroup, strconv.Itoa(result.VariantId))
							if err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to fetch operation item price",
									Err:        errors.New("failed to fetch operation item price"),
								}
							}

							// Apply markup and percentage increase
							oprItemPrice += markupAmount + (oprItemPrice * (markupPercentage / 100))

							// Fetch item discount percent
							oprItemDiscPercent, err := s.lookupRepo.GetOprItemDisc(tx, utils.LinetypeSublet, utils.TrxTypeWoFreeService.ID, recallRecord.OprItemId, agreementNo, 00002, oprItemPrice*substitute.SupplyQty, result.CompanyId, result.BrandId, 0, WhsGroup, 0)
							if err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to fetch operation item discount percent",
									Err:        errors.New("failed to fetch operation item discount percent"),
								}
							}

							// Calculate discount amount
							oprItemDiscAmount := math.Round((oprItemPrice * oprItemDiscPercent / 100))

							// Get the next work order operation item line
							if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
								Select("COALESCE(MAX(work_order_operation_item_line), 0) + 1").
								Where("work_order_system_number = ?", workOrderId).
								Scan(&woOprItemLinesub).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to get the next work order operation item line",
									Err:        err,
								}
							}

							// Step 7: Process supply quantity based on substitute type
							var supplyQty float64
							if substitute.SubstituteType == "I" {
								// Fetch quantity from smSubsStockInterchange1
								if err := tx.Model(&masteritementities.Item{}).
									Select("COALESCE(SUM(COALESCE(qty, 0)), 0)").
									Where("item_code = ?", substitute.SubstituteItemCode).
									Scan(&supplyQty).Error; err != nil {
									return entity, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Failed to fetch quantity from smSubsStockInterchange1",
										Err:        err,
									}
								}
							} else {
								supplyQty = substitute.SupplyQty
							}

							// Check item existence in item master
							var itemExists int64

							if err := tx.Model(&masteritementities.Item{}).
								Where("item_code = ? and item_group_id <> ? and item_type_id = ?", recallRecord.OprItemCode, 1, 2).
								Count(&itemExists).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to check if item exists in item master",
									Err:        err,
								}
							}

							if itemExists == 0 {
								supplyQty = entity.FrtQuantity
							} else {
								supplyQty = 0
							}

							if substitute.SubstituteType == "" {
								substitute.SupplyQty = entity.FrtQuantity
							}

							// Set ATPM warranty claim type ID
							var atpmWcfTypeId int

							if err := tx.Model(&masteritementities.Item{}).
								Select("atpm_warranty_claim_type_id").
								Where("item_code = ?", substitute.SubstituteItemCode).
								Scan(&atpmWcfTypeId).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to fetch ATPM_WCF_Type",
									Err:        err,
								}
							}

							// Step 8: Insert or update the final operation item
							finalWorkOrderDetail := transactionworkshopentities.WorkOrderDetail{
								WorkOrderSystemNumber:        workOrderId,
								WorkOrderOperationItemLine:   woOprItemLinesub,
								LineTypeId:                   utils.LinetypeSublet,
								FrtQuantity:                  supplyQty,
								OperationItemPrice:           oprItemPrice,
								OperationItemDiscountAmount:  oprItemDiscAmount,
								OperationItemDiscountPercent: oprItemDiscPercent,
								SupplyQuantity:               supplyQty,
								WarehouseGroupId:             WhsGroup,
								AtpmWCFTypeId:                atpmWcfTypeId,
							}

							if err := tx.Create(&finalWorkOrderDetail).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to insert or update operation item",
									Err:        err,
								}
							}

							// Step 9: Calculate totals for various line types and update work order
							var (
								totalPart               float64
								totalOil                float64
								totalMaterial           float64
								totalConsumableMaterial float64
								totalSublet             float64
								totalAccs               float64
							)

							if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
								Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
								Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSparepart).
								Scan(&totalPart).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to calculate total part",
									Err:        err,
								}
							}

							if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
								Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
								Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeOil).
								Scan(&totalOil).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to calculate total oil",
									Err:        err,
								}
							}

							if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
								Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
								Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeMaterial).
								Scan(&totalMaterial).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to calculate total material",
									Err:        err,
								}
							}

							if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
								Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
								Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeConsumableMaterial).
								Scan(&totalConsumableMaterial).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to calculate total consumable material",
									Err:        err,
								}
							}

							if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
								Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
								Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSublet).
								Scan(&totalSublet).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to calculate total sublet",
									Err:        err,
								}
							}

							if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
								Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
								Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeAccesories).
								Scan(&totalAccs).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to calculate total accessories",
									Err:        err,
								}
							}

							// Update the work order with the new totals
							if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
								Where("work_order_system_number = ?", workOrderId).
								Updates(map[string]interface{}{
									"total_part":                totalPart,
									"total_oil":                 totalOil,
									"total_material":            totalMaterial,
									"total_consumable_material": totalConsumableMaterial,
									"total_sublet":              totalSublet,
									"total_accessories":         totalAccs,
								}).Error; err != nil {
								return entity, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to update work order with new totals",
									Err:        err,
								}
							}
						}
					}

					if err := tx.Exec("DROP TABLE #SUBSTITUTE4").Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to drop temporary table",
							Err:        err,
						}
					}
				}

				//--==UPDATE TOTAL WORK ORDER==--
				// Query to calculate TOTAL_DISC
				var totalPackage, totalOpr, totalPart, totalOil, totalMaterial, totalConsumableMaterial, totalSublet, totalAccs, totalNonVat, totalVat, totalAfterDisc, totalAfterVat, totalPph, totalDisc float64
				var addDiscStat, TaxFree int

				if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Select(`
				SUM(
					CASE
						WHEN line_type_id = ? THEN
							CASE
								WHEN approval_id = ? THEN COALESCE(operation_item_discount_request_amount, 0)
								ELSE COALESCE(operation_item_discount_amount, 0)
							END
						ELSE
							CASE
								WHEN approval_id = ? THEN COALESCE(operation_item_discount_request_amount, 0)
								ELSE COALESCE(operation_item_discount_amount, 0)
							END *
							CASE
								WHEN line_type_id <> ? THEN COALESCE(frt_quantity, 0)
								ELSE
									CASE
										WHEN COALESCE(supply_quantity, 0) > 0 THEN COALESCE(supply_quantity, 0)
										ELSE COALESCE(frt_quantity, 0)
									END
							END
					END
				)`,
						utils.LinetypePackage, 20, 20, utils.LinetypeOperation).
					Where("work_order_system_number = ?", workOrderId).
					Scan(&totalDisc).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to calculate total discount",
						Err:        err,
					}
				}

				// Calculate TOTAL
				total := totalPackage + totalOpr + totalPart + totalOil + totalMaterial + totalConsumableMaterial + totalSublet + totalAccs

				// Calculate AddDiscStat
				if result.AdditionalDiscountStatusApprovalId == 20 {
					if result.DiscountRequestPercent != nil && *result.DiscountRequestPercent > 0 {
						result.AdditionalDiscountStatusApprovalId = 30 // Use assignment
					}
				}

				addDiscStat = result.AdditionalDiscountStatusApprovalId

				// Safely dereference the pointer
				var addDiscReqAmount *float64
				if result.DiscountRequestPercent != nil {
					addDiscReqAmount = result.DiscountRequestPercent
				} else {
					addDiscReqAmount = nil
				}

				// Rounding TOTAL_DISC
				totalDisc = math.Round(totalDisc)

				// Total After Discount
				totalAfterDisc = math.Round(total - totalDisc)

				// Calculate TOTAL_NON_VAT
				if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Select("COALESCE(SUM(ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0)), 0)").
					Where("work_order_system_number = ? AND transaction_type_id = ?", workOrderId, 6). //TrxTypeWOInternal
					Scan(&totalNonVat).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to calculate total non-VAT",
						Err:        err,
					}
				}

				// VAT Calculation
				if TaxFree == 0 {
					var vatRate float64
					if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
						Select("COALESCE(vat_tax_rate, 0.0)").
						Where("work_order_system_number = ?", workOrderId).
						Scan(&vatRate).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to fetch VAT tax rate",
							Err:        err,
						}
					}

					// VAT Amount Calculation
					totalVat = math.Floor((totalAfterDisc - totalNonVat) * vatRate / 100)

				} else {

					totalVat = 0
				}

				// Total After VAT
				totalAfterVat = math.Round(totalAfterDisc + totalVat)

				// PPH Calculation
				if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Where("work_order_system_number = ? AND (line_type_id = ? OR line_type_id = ?)", workOrderId, utils.LinetypePackage, utils.LinetypeOperation).
					Pluck("FLOOR(SUM(COALESCE(pph_amount, 0)))", &totalPph).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to calculate total PPH",
						Err:        err,
					}
				}

				//update work order
				if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
					Where("work_order_system_number = ?", workOrderId).
					Updates(map[string]interface{}{
						"total":                               total,
						"total_pph":                           totalPph,
						"total_discount":                      totalDisc,
						"total_after_discount":                totalAfterDisc,
						"total_vat":                           totalVat,
						"total_after_vat":                     totalAfterVat,
						"additional_discount_status_approval": addDiscStat,
						"discount_request_amount":             addDiscReqAmount,
					}).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update work order",
						Err:        err,
					}
				}

				// update recall
				if err := tx.Model(&transactionworkshopentities.RecallDetailType{}).
					Where("rec_system_number = ? and recall_line_number = ? and operation_item_code = ?", entity.RecSystemNumber, recdettype.RecallLineNumber, recallRecord.OprItemCode).
					Updates(map[string]interface{}{
						"has_recall": true,
					}).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update recall",
						Err:        err,
					}
				}

			}

		}

		// Step 1: Check if there are any records in atRecall2 where HAS_RECALL is 0
		var recallCount int64
		if err := tx.Model(&transactionworkshopentities.RecallDetailType{}).
			Where("recall_system_number = ? AND recall_line_number = ? AND has_recall = 0", recdettype.RecallSystemNumber, recdettype.RecallLineNumber).
			Count(&recallCount).Error; err != nil {
			return entity, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check atRecall2 record count",
				Err:        err,
			}
		}

		// Step 2: If no records exist, update atRecall1 to set HAS_RECALL = 1
		if recallCount == 0 {
			if err := tx.Model(&transactionworkshopentities.RecallDetail{}).
				Where("vechicle_chassis_number = ? AND recall_system_number = ? AND recall_line_number = ?", vehicleChassisNo, recdettype.RecallSystemNumber, recdettype.RecallLineNumber).
				Updates(map[string]interface{}{
					"has_recall":  1,
					"recall_by":   companyCode,
					"recall_date": time.Now(),
				}).Error; err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update atRecall1",
					Err:        err,
				}
			}
		}

		// Step 3: Check if there is any record in wtWorkOrder1 for the given WO_SYS_NO
		var woServReqLine3 int64
		if err := tx.Model(&transactionworkshopentities.WorkOrderRequestDescription{}).
			Where("work_order_system_number = ?", workOrderId).
			Count(&woServReqLine3).Error; err != nil {
			return entity, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check wtWorkOrder1 for WO_SYS_NO",
				Err:        err,
			}
		}

		// Step 4: Determine the value for @Wo_Serv_Req_Line3
		if woServReqLine3 == 0 {
			woServReqLine3 = 1
		} else {
			if err := tx.Model(&transactionworkshopentities.WorkOrderRequestDescription{}).
				Where("work_order_system_number = ?", workOrderId).
				Select("MAX(work_order_service_request_line) + 1").
				Scan(&woServReqLine3).Error; err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to calculate WO_SERV_REQ_LINE",
					Err:        err,
				}
			}

		}

		// Step 5: Insert new record into wtWorkOrder1
		var recallData struct {
			WoServReq string
		}

		if err := tx.Model(&transactionworkshopentities.Recall{}).
			Select("is_active, ISNULL(recall_document_number, '') + ' - ' + ISNULL(recall_name, '') AS WoServReq").
			Where("recall_system_number = ?", entity.RecSystemNumber).
			Scan(&recallData).Error; err != nil {
			return entity, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch data from atRecall0",
				Err:        err,
			}
		}

		// Step 2: Insert the queried data into wtWorkOrder1
		insertData := map[string]interface{}{
			"work_order_system_number":        workOrderId,
			"work_order_service_request_line": woServReqLine3,
			"work_order_service_request":      recallData.WoServReq,
		}

		if err := tx.Model(&transactionworkshopentities.WorkOrderRequestDescription{}).
			Create(insertData).Error; err != nil {
			return entity, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to insert into wtWorkOrder1",
				Err:        err,
			}
		}
	}

	return entity, nil
}

func (r *WorkOrderRepositoryImpl) GetServiceRequestByWO(tx *gorm.DB, workOrderId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	// Initialize the service request structure with the entity
	var tableStruct transactionworkshopentities.ServiceRequest

	// Create base query to apply filters and conditions
	query := tx.Model(&tableStruct).Where("work_order_system_number = ? AND work_order_system_number != 0", workOrderId)

	// Apply filters to the query
	query = utils.ApplyFilterSearch(query, filterCondition)

	// Get total rows for pagination
	var totalRows int64
	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count service requests",
			Err:        err,
		}
	}

	// Apply pagination
	paginatedQuery := pagination.Paginate(&tableStruct, &pages, query)

	// Execute the query
	var results []transactionworkshopentities.ServiceRequest
	if err := paginatedQuery(tx).Find(&results).Error; err != nil { // Call paginatedQuery with tx
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve service requests",
			Err:        err,
		}
	}

	// Prepare the response data
	response := make([]map[string]interface{}, len(results))
	for i, request := range results {
		response[i] = map[string]interface{}{
			"service_request_system_number":   request.ServiceRequestSystemNumber,
			"service_request_document_number": request.ServiceRequestDocumentNumber,
			"work_order_system_number":        request.WorkOrderSystemNumber,
		}
	}

	// Set pagination metadata
	pages.TotalRows = totalRows
	pages.TotalPages = int(math.Ceil(float64(totalRows) / float64(pages.Limit)))

	return response, pages.TotalPages, int(totalRows), nil
}

func (r *WorkOrderRepositoryImpl) GetClaimByWO(tx *gorm.DB, workOrderId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	entities := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Where("work_order_system_number = ? AND work_order_system_number != 0 AND line_type_id IN ('0','1') AND work_order_status_id != ?", workOrderId, utils.WoStatOngoing)

	entities = utils.ApplyFilter(entities, filterCondition)

	var totalRows int64
	if err := entities.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count filtered claims",
			Err:        err,
		}
	}

	var results []transactionworkshopentities.WorkOrderDetail

	if err := entities.Order("work_order_detail_id").Offset(pages.GetOffset()).Limit(pages.GetLimit()).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve filtered claims",
			Err:        err,
		}
	}

	response := make([]map[string]interface{}, len(results))
	for i, claim := range results {
		response[i] = map[string]interface{}{
			"work_order_system_number": claim.WorkOrderSystemNumber,
		}
	}

	pages.TotalRows = totalRows
	pages.TotalPages = int(math.Ceil(float64(totalRows) / float64(pages.GetLimit())))

	return response, pages.TotalPages, int(totalRows), nil
}

func (r *WorkOrderRepositoryImpl) GetClaimItemByWO(tx *gorm.DB, workOrderId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	entities := tx.Table("trx_work_order_detail").
		Joins("INNER JOIN mtr_item ON trx_work_order_detail.operation_item_id = mtr_item.item_id").
		Where("trx_work_order_detail.work_order_system_number = ? AND trx_work_order_detail.work_order_system_number != 0", workOrderId).
		Where("trx_work_order_detail.work_order_status_id = ?", utils.WoStatOngoing).
		Where("trx_work_order_detail.supply_quantity > 0 AND trx_work_order_detail.invoice_system_number = 0").
		Where("trx_work_order_detail.warranty_claim_type_id IN (?, ?)", 3, "")

	entities = utils.ApplyFilter(entities, filterCondition)

	var totalRows int64
	if err := entities.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count filtered claims",
			Err:        err,
		}
	}

	if totalRows == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No claims found for the given work order.",
			Err:        nil,
		}
	}

	var results []transactionworkshopentities.WorkOrderDetail
	paginatedQuery := pagination.Paginate(&transactionworkshopentities.WorkOrderDetail{}, &pages, entities)

	if err := paginatedQuery(tx).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve filtered claims",
			Err:        err,
		}
	}

	response := []map[string]interface{}{
		{
			"work_order_system_number": results[0].WorkOrderSystemNumber,
		},
	}

	pages.TotalRows = totalRows
	pages.TotalPages = int(math.Ceil(float64(totalRows) / float64(pages.GetLimit())))

	return response, pages.TotalPages, int(totalRows), nil
}

func (r *WorkOrderRepositoryImpl) GetWOByBillCode(tx *gorm.DB, workOrderId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	query := tx.Table("trx_work_order_detail").
		Where("trx_work_order_detail.work_order_system_number = ? AND trx_work_order_detail.work_order_system_number != 0", workOrderId).
		Where("trx_work_order_detail.frt_quantity > trx_work_order_detail.supply_quantity AND trx_work_order_detail.invoice_system_number = 0").
		Where("trx_work_order_detail.line_type_id NOT IN ('0', '1')")

	query = utils.ApplyFilter(query, filterCondition)

	var totalRows int64
	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count filtered work orders",
			Err:        err,
		}
	}

	if totalRows == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No work orders found for the given work order ID.",
			Err:        nil,
		}
	}

	paginatedQuery := pagination.Paginate(&transactionworkshopentities.WorkOrderDetail{}, &pages, query)

	var results []transactionworkshopentities.WorkOrderDetail
	if err := paginatedQuery(tx).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve filtered work orders",
			Err:        err,
		}
	}

	response := []map[string]interface{}{
		{
			"work_order_system_number": results[0].WorkOrderSystemNumber,
		},
	}

	pages.TotalRows = totalRows
	pages.TotalPages = int(math.Ceil(float64(totalRows) / float64(pages.GetLimit())))

	return response, pages.TotalPages, int(totalRows), nil
}

func (r *WorkOrderRepositoryImpl) GetDetailWOByClaimBillCode(tx *gorm.DB, workOrderId int, transactionTypeId int, atpmClaimNumber string, pages pagination.Pagination) ([]transactionworkshoppayloads.GetClaimResponsePayload, *exceptions.BaseErrorResponse) {
	var responsePayload []transactionworkshoppayloads.GetClaimResponsePayload

	query := tx.Table("trx_work_order_detail AS A").
		Joins("INNER JOIN trx_work_order AS B ON B.work_order_system_number = A.work_order_system_number").
		Joins("LEFT JOIN mtr_item AS C ON C.item_id = A.operation_item_id").
		Where("A.work_order_system_number = ? AND A.transaction_type_id = ? AND A.invoice_system_number = 0", workOrderId, transactionTypeId).
		Where("A.atpm_claim_number != ''").
		Where("A.warranty_claim_type_id = (CASE WHEN A.line_type_id = '1' THEN A.warranty_claim_type_id WHEN A.line_type_id = '0' THEN A.warranty_claim_type_id ELSE 'PM' END)").
		Where("A.atpm_claim_number = ?", atpmClaimNumber).
		Where("A.work_order_status_id = (CASE WHEN A.line_type_id = '1' THEN ? WHEN A.line_type_id = '0' THEN ? ELSE 0 END)", utils.WoStatOngoing, utils.WoStatOngoing).
		Where("A.supply_quantity = A.frt_quantity")

	err := query.Select(
		"A.work_order_system_number, B.work_order_document_number, A.work_order_operation_item_line, B.vehicle_chassis_number, B.brand_id, B.model_id, B.variant_id, " +
			"C.item_group_id, A.line_type_id, A.operation_item_code, A.frt_quantity, A.supply_quantity, A.approval_id, A.operation_item_price, " +
			"A.operation_item_discount_request_percent, A.operation_item_discount_percent, A.total_cost_of_goods_sold, A.job_type_id, " +
			"A.purchase_order_system_number, A.purchase_order_line, A.description").
		Limit(pages.GetLimit()).
		Offset(pages.GetOffset()).
		Find(&responsePayload).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch claim details",
			Err:        err,
		}
	}

	// Check if responsePayload is empty and return a Not Found error if so
	if len(responsePayload) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No claim details found for the specified work order.",
			Err:        nil,
		}
	}

	return responsePayload, nil
}

func (r *WorkOrderRepositoryImpl) GetDetailWOByBillCode(tx *gorm.DB, workOrderId int, transactionTypeId int, pages pagination.Pagination) ([]transactionworkshoppayloads.GetClaimResponsePayload, *exceptions.BaseErrorResponse) {

	var responsePayload []transactionworkshoppayloads.GetClaimResponsePayload

	query := tx.Table("trx_work_order_detail AS A").
		Joins("INNER JOIN trx_work_order AS B ON B.work_order_system_number = A.work_order_system_number").
		Joins("LEFT JOIN mtr_item AS C ON C.item_id = A.operation_item_id").
		Where("A.work_order_system_number = ? AND A.transaction_type_id = ? AND A.invoice_system_number = 0", workOrderId, transactionTypeId).
		Where("A.work_order_status_id = CASE WHEN A.line_type_id = '1' THEN ? WHEN A.line_type_id = '0' THEN ? ELSE 0 END", utils.WoStatOngoing, utils.WoStatOngoing).
		Where("A.supply_quantity = A.frt_quantity")

	err := query.Select(
		"A.work_order_system_number, B.work_order_document_number, A.work_order_operation_item_line, B.vehicle_chassis_number, B.brand_id, B.model_id, B.variant_id, " +
			"C.item_group_id, A.line_type_id, A.operation_item_code, A.frt_quantity, A.supply_quantity, A.approval_id, A.operation_item_price, " +
			"A.operation_item_discount_request_percent, A.operation_item_discount_percent, A.total_cost_of_goods_sold, A.job_type_id, " +
			"A.purchase_order_system_number, A.purchase_order_line, A.description").
		Limit(pages.GetLimit()).
		Offset(pages.GetOffset()).
		Find(&responsePayload).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch claim details",
			Err:        err,
		}
	}

	// Check if responsePayload is empty and return a Not Found error if so
	if len(responsePayload) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No claim details found for the specified work order.",
			Err:        nil,
		}
	}

	return responsePayload, nil
}

func (r *WorkOrderRepositoryImpl) GetDetailWOByATPMBillCode(tx *gorm.DB, workOrderId int, transactionTypeId int, pages pagination.Pagination) ([]transactionworkshoppayloads.GetClaimResponsePayload, *exceptions.BaseErrorResponse) {

	var responsePayload []transactionworkshoppayloads.GetClaimResponsePayload

	query := tx.Table("trx_work_order_detail AS A").
		Joins("INNER JOIN trx_work_order AS B ON B.work_order_system_number = A.work_order_system_number").
		Joins("LEFT JOIN mtr_item AS C ON C.item_id = A.operation_item_id").
		Where("A.work_order_system_number = ? AND A.transaction_type_id = ? AND A.invoice_system_number = 0", workOrderId, transactionTypeId).
		Where("A.warranty_claim_type_id = CASE WHEN A.line_type_id = '1' THEN A.warranty_claim_type_id WHEN A.line_type_id = '0' THEN A.warranty_claim_type_id ELSE 'PM' END").
		Where("A.work_order_status_id = CASE WHEN A.line_type_id = '1' THEN ? WHEN A.line_type_id = '0' THEN ? ELSE 0 END", utils.WoStatOngoing, utils.WoStatOngoing).
		Where("A.supply_quantity = A.frt_quantity")

	err := query.Select(
		"A.work_order_system_number, B.work_order_document_number, A.work_order_operation_item_line, B.vehicle_chassis_number, B.brand_id, B.model_id, B.variant_id, " +
			"C.item_group_id, A.line_type_id, A.operation_item_code, A.frt_quantity, A.supply_quantity, A.approval_id, A.operation_item_price, " +
			"A.operation_item_discount_request_percent, A.operation_item_discount_percent, A.total_cost_of_goods_sold, A.job_type_id, " +
			"A.purchase_order_system_number, A.purchase_order_line, A.description").
		Limit(pages.GetLimit()).
		Offset(pages.GetOffset()).
		Find(&responsePayload).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch claim details",
			Err:        err,
		}
	}

	// Check if responsePayload is empty and return a Not Found error if so
	if len(responsePayload) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No claim details found for the specified work order.",
			Err:        nil,
		}
	}

	return responsePayload, nil
}

func (r *WorkOrderRepositoryImpl) GetSupplyByWO(tx *gorm.DB, workOrderId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	query := tx.Table("trx_supply_slip AS ss0").
		Joins("INNER JOIN trx_supply_slip_return_detail AS ss1 ON ss1.supply_system_number = ss0.supply_system_number").
		Joins("LEFT OUTER JOIN trx_work_order_detail AS w2 ON w2.work_order_system_number = ss1.work_order_system_number AND w2.operation_item_code = ss1.item_code").
		Where("ss0.work_order_system_number = ? AND ss0.supply_type = 'B'", workOrderId).
		Where("(ss1.quantity_supply - ss1.quantity_return) > 0").
		Where("ISNULL(w2.work_order_system_number, 0) = 0")

	query = utils.ApplyFilter(query, filterCondition)

	var totalRows int64
	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count filtered supply orders",
			Err:        err,
		}
	}

	if totalRows == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No supply orders found for the given work order ID.",
			Err:        nil,
		}
	}

	paginatedQuery := pagination.Paginate(&transactionsparepartentities.SupplySlip{}, &pages, query)

	var results []transactionsparepartentities.SupplySlip
	if err := paginatedQuery(tx).Find(&results).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve filtered supply orders",
			Err:        err,
		}
	}

	response := []map[string]interface{}{}
	for _, result := range results {
		response = append(response, map[string]interface{}{
			"supply_system_number":   result.SupplySystemNumber,
			"supply_document_number": result.SupplyDocumentNumber,
			// Add more fields as required
		})
	}

	pages.TotalRows = totalRows
	pages.TotalPages = int(math.Ceil(float64(totalRows) / float64(pages.GetLimit())))

	return response, pages.TotalPages, int(totalRows), nil
}
