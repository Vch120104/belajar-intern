package transactionworkshoprepositoryimpl

////  NOTES  ////
import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masteroperationentities "after-sales/api/entities/master/operation"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	transactionunitentities "after-sales/api/entities/transaction/Unit"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	masterrepository "after-sales/api/repositories/master"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
func (r *WorkOrderRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	tableStruct := transactionworkshoppayloads.WorkOrderGetAllRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	whereQuery = whereQuery.Where("service_request_system_number = 0 AND booking_system_number = 0 AND estimation_system_number = 0 and cpc_code = '00002'")

	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Failed to retrieve work order data from the database",
			Err:        err,
		}
	}
	defer rows.Close()

	var convertedResponses []transactionworkshoppayloads.WorkOrderGetAllResponse

	for rows.Next() {
		var workOrderReq transactionworkshoppayloads.WorkOrderGetAllRequest
		var workOrderRes transactionworkshoppayloads.WorkOrderGetAllResponse

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
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check work order data from the database",
				Err:        err,
			}
		}

		// Fetch external data for brand, model, vehicle, work order type, and status
		getBrandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(workOrderReq.BrandId)
		if brandErr != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: brandErr.StatusCode,
				Message:    "Failed to fetch brand data from external service",
				Err:        brandErr.Err,
			}
		}

		getModelResponse, modelErr := salesserviceapiutils.GetUnitModelById(workOrderReq.ModelId)
		if modelErr != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: modelErr.StatusCode,
				Message:    "Failed to fetch model data from external service",
				Err:        modelErr.Err,
			}
		}

		// vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(workOrderReq.VehicleId)
		// if vehicleErr != nil {
		// 	return pagination.Pagination{}, &exceptions.BaseErrorResponse{
		// 		StatusCode: vehicleErr.StatusCode,
		// 		Message:    "Failed to retrieve vehicle data from the external API",
		// 		Err:        vehicleErr.Err,
		// 	}
		// }

		getWorkOrderTypeResponses, workOrderTypeErr := generalserviceapiutils.GetWorkOrderTypeByID(workOrderReq.WorkOrderTypeId)
		if workOrderTypeErr != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: workOrderTypeErr.StatusCode,
				Message:    "Failed to fetch work order type data from external service",
				Err:        workOrderTypeErr.Err,
			}
		}

		getWorkOrderStatusResponses, workOrderStatusErr := generalserviceapiutils.GetWorkOrderStatusById(workOrderReq.StatusId)
		if workOrderStatusErr != nil {
			return pagination.Pagination{}, &exceptions.BaseErrorResponse{
				StatusCode: workOrderStatusErr.StatusCode,
				Message:    "Failed to fetch work order status data from external service",
				Err:        workOrderStatusErr.Err,
			}
		}

		workOrderRes = transactionworkshoppayloads.WorkOrderGetAllResponse{
			WorkOrderDocumentNumber: workOrderReq.WorkOrderDocumentNumber,
			WorkOrderSystemNumber:   workOrderReq.WorkOrderSystemNumber,
			WorkOrderDate:           workOrderReq.WorkOrderDate,
			FormattedWorkOrderDate:  utils.FormatRFC3339(workOrderReq.WorkOrderDate), // Use RFC3339 format
			WorkOrderTypeId:         workOrderReq.WorkOrderTypeId,
			WorkOrderTypeName:       getWorkOrderTypeResponses.WorkOrderTypeName,
			BrandId:                 workOrderReq.BrandId,
			BrandName:               getBrandResponse.BrandName,
			VehicleCode:             "", //vehicleResponses.VehicleChassisNumber,
			VehicleTnkb:             "", //vehicleResponses.VehicleRegistrationCertificateTNKB,
			ModelId:                 workOrderReq.ModelId,
			ModelName:               getModelResponse.ModelName,
			VehicleId:               workOrderReq.VehicleId,
			CustomerId:              workOrderReq.CustomerId,
			StatusId:                workOrderReq.StatusId,
			StatusName:              getWorkOrderStatusResponses.WorkOrderStatusName,
			RepeatedJob:             workOrderReq.RepeatedJob,
		}

		convertedResponses = append(convertedResponses, workOrderRes)
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_document_number":  response.WorkOrderDocumentNumber,
			"work_order_system_number":    response.WorkOrderSystemNumber,
			"work_order_date":             response.FormattedWorkOrderDate,
			"work_order_type_id":          response.WorkOrderTypeId,
			"work_order_type_description": response.WorkOrderTypeName,
			"brand_id":                    response.BrandId,
			"brand_name":                  response.BrandName,
			"model_id":                    response.ModelId,
			"model_name":                  response.ModelName,
			"vehicle_id":                  response.VehicleId,
			//"vehicle_chassis_number":      response.VehicleCode,
			//"vehicle_tnkb":                response.VehicleTnkb,
			"work_order_status_id":   response.StatusId,
			"work_order_status_name": response.StatusName,
			"repeated_system_number": response.RepeatedJob,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	pages.Rows = paginatedData
	pages.TotalRows = int64(totalRows)
	pages.TotalPages = totalPages

	return pages, nil
}

// uspg_wtWorkOrder0_Insert
// IF @Option = 0
// --USE FOR : * INSERT NEW DATA
func (r *WorkOrderRepositoryImpl) New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderNormalRequest) (transactionworkshopentities.WorkOrder, *exceptions.BaseErrorResponse) {

	// Default values
	// 1:Draft, 2:New, 3:Ready, 4:On Going, 5:Stop, 6:QC Pass, 7:Cancel, 8:Closed
	// 1:Normal, 2:Campaign, 3:Affiliated, 4:Repeat Job
	defaultWorkOrderDocumentNumber := ""
	defaultCPCcode := "00002" // Default CPC code 00002 for workshop
	workOrderTypeId := 1      // Default work order type ID 1 for normal

	// fetch currency by code
	getCurrencyId, currencyErr := financeserviceapiutils.GetCurrencyByCode("IDR")
	if currencyErr != nil {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: currencyErr.StatusCode,
			Message:    "Failed to fetch currency data from external service",
			Err:        currencyErr.Err,
		}
	}

	// Validation: request date
	loc, _ := time.LoadLocation("Asia/Jakarta") // UTC+7
	currentDate := time.Now().In(loc).Format("2006-01-02T15:04:05Z")
	parsedTime, _ := time.Parse(time.RFC3339, currentDate)

	requestDate := request.WorkOrderArrivalTime.In(loc).Truncate(24 * time.Hour)

	if requestDate.Before(parsedTime.Truncate(24*time.Hour)) || requestDate.After(parsedTime.Truncate(24*time.Hour)) {
		request.WorkOrderArrivalTime = parsedTime
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
	// vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(request.VehicleId)
	// if vehicleErr != nil {
	// 	return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    "Failed to retrieve vehicle data from the external API",
	// 		Err:        vehicleErr.Err,
	// 	}
	// }

	// fetch approval
	approvalStatus, approvalErr := generalserviceapiutils.GetApprovalStatusByCode("10")
	if approvalErr != nil {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: approvalErr.StatusCode,
			Message:    "Failed to fetch approval status data from external service",
			Err:        approvalErr.Err,
		}
	}

	// Create WorkOrder entity
	entitieswo := transactionworkshopentities.WorkOrder{
		// page 1
		WorkOrderDocumentNumber:    defaultWorkOrderDocumentNumber,
		WorkOrderStatusId:          utils.WoStatDraft,
		WorkOrderDate:              parsedTime,
		CPCcode:                    defaultCPCcode,
		ServiceAdvisorId:           request.ServiceAdvisorId,
		WorkOrderTypeId:            workOrderTypeId,
		BookingSystemNumber:        request.BookingSystemNumber,
		EstimationSystemNumber:     request.EstimationSystemNumber,
		ServiceRequestSystemNumber: request.ServiceRequestSystemNumber,
		PDISystemNumber:            request.PDISystemNumber,
		RepeatedSystemNumber:       request.RepeatedSystemNumber,
		ServiceSite:                "OD - Service On Dealer",
		VehicleChassisNumber:       "", //vehicleResponses.VehicleChassisNumber,

		// Page 1
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
		EraExpiredDate:           utils.ToTimePointer(request.WorkOrderEraExpiredDate),
		InsuranceCheck:           request.WorkOrderInsuranceCheck,
		InsurancePolicyNumber:    &request.WorkOrderInsurancePolicyNo,
		InsuranceExpiredDate:     utils.ToTimePointer(request.WorkOrderInsuranceExpiredDate),
		InsuranceClaimNumber:     &request.WorkOrderInsuranceClaimNo,
		InsurancePersonInCharge:  &request.WorkOrderInsurancePic,
		InsuranceOwnRisk:         &request.WorkOrderInsuranceOwnRisk,
		InsuranceWorkOrderNumber: &request.WorkOrderInsuranceWONumber,

		// page 2
		AdditionalDiscountStatusApprovalId: approvalStatus.ApprovalStatusId,
		EstTime:                            request.EstimationDuration,
		CustomerExpress:                    request.CustomerExpress,
		LeaveCar:                           request.LeaveCar,
		CarWash:                            request.CarWash,
		PromiseDate:                        request.PromiseDate,
		PromiseTime:                        request.PromiseTime,
		FSCouponNo:                         request.FSCouponNo,
		Notes:                              request.Notes,
		Suggestion:                         request.Suggestion,
		DPAmount:                           request.DownpaymentAmount,
		CurrencyId:                         getCurrencyId.CurrencyId,
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
	brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
	if brandErr != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve brand data from the external API",
			Err:        brandErr.Err,
		}
	}

	// Fetch data model from external API
	modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
	if modelErr != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from the external API",
			Err:        modelErr.Err,
		}
	}

	// Fetch data variant from external API
	variantResponse, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
	if variantErr != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        variantErr.Err,
		}
	}

	// Fetch data vehicle from external API
	// vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
	// if vehicleErr != nil {
	// 	return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    "Failed to retrieve vehicle data from the external API",
	// 		Err:        vehicleErr.Err,
	// 	}
	// }

	// Fetch workorder details count
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
		Select("trx_work_order_detail.work_order_detail_id, trx_work_order_detail.work_order_system_number, trx_work_order_detail.line_type_id, lt.line_type_code as line_type_code, lt.line_type_name as line_type_name,trx_work_order_detail.transaction_type_id, tt.work_order_transaction_type_code AS transaction_type_code, trx_work_order_detail.job_type_id, tc.job_type_code AS job_type_code, trx_work_order_detail.warehouse_group_id, trx_work_order_detail.frt_quantity, trx_work_order_detail.supply_quantity, trx_work_order_detail.operation_item_price, trx_work_order_detail.operation_item_discount_amount, trx_work_order_detail.operation_item_discount_request_amount").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_line_type AS lt ON lt.line_type_id = trx_work_order_detail.line_type_id").
		Joins("INNER JOIN dms_microservices_general_dev.dbo.mtr_work_order_transaction_type AS tt ON tt.work_order_transaction_type_id = trx_work_order_detail.transaction_type_id").
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
	getWorkOrderStatusResponses, workOrderStatusErr := generalserviceapiutils.GetWorkOrderStatusById(entity.WorkOrderStatusId)
	if workOrderStatusErr != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order status data from external service",
			Err:        workOrderStatusErr.Err,
		}
	}

	// fetch data type work order
	getWorkOrderTypeResponses, workOrderTypeErr := generalserviceapiutils.GetWorkOrderTypeByID(entity.WorkOrderTypeId)
	if workOrderTypeErr != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order type data from external service",
			Err:        workOrderTypeErr.Err,
		}
	}

	// fetch approval status
	getApprovalStatusResponses, approvalStatusErr := generalserviceapiutils.GetApprovalStatusById(entity.AdditionalDiscountStatusApprovalId)
	if approvalStatusErr != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch approval status data from external service",
			Err:        approvalStatusErr.Err,
		}
	}

	// fecth data currency
	getCurrencyResponses, currencyErr := financeserviceapiutils.GetCurrencyId(entity.CurrencyId)
	if currencyErr != nil {
		return transactionworkshoppayloads.WorkOrderResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch currency data from external service",
			Err:        currencyErr.Err,
		}
	}

	payload := transactionworkshoppayloads.WorkOrderResponseDetail{
		WorkOrderSystemNumber:              entity.WorkOrderSystemNumber,
		WorkOrderDate:                      entity.WorkOrderDate,
		WorkOrderDocumentNumber:            entity.WorkOrderDocumentNumber,
		WorkOrderTypeId:                    entity.WorkOrderTypeId,
		WorkOrderTypeName:                  getWorkOrderTypeResponses.WorkOrderTypeName,
		WorkOrderStatusId:                  entity.WorkOrderStatusId,
		WorkOrderStatusName:                getWorkOrderStatusResponses.WorkOrderStatusName,
		ServiceAdvisorId:                   entity.ServiceAdvisorId,
		BrandId:                            entity.BrandId,
		BrandName:                          brandResponse.BrandName,
		ModelId:                            entity.ModelId,
		ModelName:                          modelResponse.ModelName,
		VariantId:                          entity.VariantId,
		VariantDescription:                 variantResponse.VariantDescription,
		VehicleId:                          entity.VehicleId,
		VehicleCode:                        "", //vehicleResponses.VehicleChassisNumber,
		VehicleTnkb:                        "", //vehicleResponses.VehicleRegistrationCertificateTNKB,
		CustomerId:                         entity.CustomerId,
		ServiceSite:                        entity.ServiceSite,
		BilltoCustomerId:                   entity.BillableToId,
		CampaignId:                         entity.CampaignId,
		FromEra:                            entity.FromEra,
		WorkOrderEraNo:                     entity.EraNumber,
		Storing:                            entity.Storing,
		WorkOrderCurrentMileage:            entity.ServiceMileage,
		WorkOrderProfitCenterId:            entity.ProfitCenterId,
		AgreementId:                        entity.AgreementBodyRepairId,
		BoookingId:                         entity.BookingSystemNumber,
		EstimationId:                       entity.EstimationSystemNumber,
		ContractSystemNumber:               entity.ContractServiceSystemNumber,
		QueueSystemNumber:                  entity.QueueNumber,
		WorkOrderArrivalTime:               entity.ArrivalTime,
		WorkOrderRemark:                    entity.Remark,
		DealerRepresentativeId:             entity.CostCenterId,
		CompanyId:                          entity.CompanyId,
		Titleprefix:                        entity.CPTitlePrefix,
		NameCust:                           entity.ContactPersonName,
		PhoneCust:                          entity.ContactPersonPhone,
		MobileCust:                         entity.ContactPersonMobile,
		MobileCustAlternative:              entity.ContactPersonMobileAlternative,
		MobileCustDriver:                   entity.ContactPersonMobileDriver,
		ContactVia:                         entity.ContactPersonContactVia,
		WorkOrderInsurancePolicyNo:         entity.InsurancePolicyNumber,
		WorkOrderInsuranceClaimNo:          entity.InsuranceClaimNumber,
		WorkOrderInsuranceExpiredDate:      entity.InsuranceExpiredDate,
		WorkOrderEraExpiredDate:            entity.EraExpiredDate,
		PromiseDate:                        entity.PromiseDate,
		PromiseTime:                        entity.PromiseTime,
		EstimationDuration:                 entity.EstTime,
		WorkOrderInsuranceOwnRisk:          entity.InsuranceOwnRisk,
		WorkOrderInsurancePic:              entity.InsurancePersonInCharge,
		WorkOrderInsuranceWONumber:         entity.InsuranceWorkOrderNumber,
		WorkOrderInsuranceCheck:            entity.InsuranceCheck,
		AdditionalDiscountStatusApprovalId: entity.AdditionalDiscountStatusApprovalId,
		AdditionalDiscountStatusApproval:   getApprovalStatusResponses.ApprovalStatusDescription,
		CustomerExpress:                    entity.CustomerExpress,
		LeaveCar:                           entity.LeaveCar,
		CarWash:                            entity.CarWash,
		FSCouponNo:                         entity.FSCouponNo,
		Notes:                              entity.Notes,
		Suggestion:                         entity.Suggestion,
		DPAmount:                           entity.DPAmount,
		DPPayment:                          entity.DPPayment,
		DPPaymentAllocated:                 entity.DPPaymentAllocated,
		DPPaymentVAT:                       entity.DPPaymentVAT,
		DPAllocToInv:                       entity.DPAllocToInv,
		InvoiceSystemNumber:                entity.InvoiceSystemNumber,
		CurrencyId:                         entity.CurrencyId,
		CurrencyCode:                       getCurrencyResponses.CurrencyCode,

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
			TotalRows:  int(pagination.TotalRows),
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

	updates := make(map[string]interface{})

	if request.BilltoCustomerId != 0 {
		updates["billable_to_id"] = request.BilltoCustomerId
	}
	if request.FromEra {
		updates["from_era"] = request.FromEra
	}
	if request.QueueSystemNumber != 0 {
		updates["queue_number"] = request.QueueSystemNumber
	}
	if !request.WorkOrderArrivalTime.IsZero() {
		updates["arrival_time"] = utils.FormatTimeForJSON(request.WorkOrderArrivalTime)
	}
	if request.WorkOrderCurrentMileage != 0 {
		updates["service_mileage"] = request.WorkOrderCurrentMileage
	}
	if request.Storing {
		updates["storing"] = request.Storing
	}
	if request.WorkOrderRemark != "" {
		updates["remark"] = request.WorkOrderRemark
	}
	if request.Unregistered {
		updates["unregister"] = request.Unregistered
	}
	if request.WorkOrderProfitCenter != 0 {
		updates["profit_center_id"] = request.WorkOrderProfitCenter
	}
	if request.DealerRepresentativeId != 0 {
		updates["cost_center_id"] = request.DealerRepresentativeId
	}
	if request.CompanyId != 0 {
		updates["company_id"] = request.CompanyId
	}

	// Contact person details
	if request.Titleprefix != "" {
		updates["contact_person_title_prefix"] = request.Titleprefix
	}
	if request.NameCust != "" {
		updates["contact_person_name"] = request.NameCust
	}
	if request.PhoneCust != "" {
		updates["contact_person_phone"] = request.PhoneCust
	}
	if request.MobileCust != "" {
		updates["contact_person_mobile"] = request.MobileCust
	}
	if request.MobileCustAlternative != "" {
		updates["contact_person_mobile_alternative"] = request.MobileCustAlternative
	}
	if request.MobileCustDriver != "" {
		updates["contact_person_mobile_driver"] = request.MobileCustDriver
	}
	if request.ContactVia != "" {
		updates["contact_person_contact_via"] = request.ContactVia
	}

	// Insurance details
	if request.WorkOrderInsuranceCheck {
		updates["insurance_check"] = request.WorkOrderInsuranceCheck
	}
	if request.WorkOrderInsurancePolicyNo != "" {
		updates["insurance_policy_number"] = request.WorkOrderInsurancePolicyNo
	}
	if !request.WorkOrderInsuranceExpiredDate.IsZero() {
		updates["insurance_expired_date"] = utils.FormatTimeForJSON(request.WorkOrderInsuranceExpiredDate)
	}
	if request.WorkOrderInsuranceClaimNo != "" {
		updates["insurance_claim_number"] = request.WorkOrderInsuranceClaimNo
	}
	if request.WorkOrderInsurancePic != "" {
		updates["insurance_person_in_charge"] = request.WorkOrderInsurancePic
	}
	if request.WorkOrderInsuranceOwnRisk != 0 {
		updates["insurance_own_risk"] = request.WorkOrderInsuranceOwnRisk
	}
	if request.WorkOrderInsuranceWONumber != "" {
		updates["insurance_work_order_number"] = request.WorkOrderInsuranceWONumber
	}

	// Other work order details (Page 2 fields)
	if request.EstimationDuration != 0 {
		updates["estimate_time"] = request.EstimationDuration
	}
	if request.CustomerExpress {
		updates["customer_express"] = request.CustomerExpress
	}
	if request.LeaveCar {
		updates["leave_car"] = request.LeaveCar
	}
	if request.CarWash {
		updates["car_wash"] = request.CarWash
	}
	if !request.PromiseDate.IsZero() {
		updates["promise_date"] = utils.FormatTimeForJSON(request.PromiseDate)
	}
	if !request.PromiseTime.IsZero() { //
		updates["promise_time"] = utils.FormatTimeForJSON(request.PromiseTime)
	}
	if request.FSCouponNo != "" {
		updates["fs_coupon_number"] = request.FSCouponNo
	}
	if request.Notes != "" {
		updates["notes"] = request.Notes
	}
	if request.Suggestion != "" {
		updates["suggestion"] = request.Suggestion
	}
	if request.DownpaymentAmount != 0 {
		updates["downpayment_amount"] = request.DownpaymentAmount
	}

	// Handling VAT Tax Rate
	if generalserviceapiutils.IsFTZCompany(request.CompanyId) {
		vatZero := 0.0
		updates["vat_tax_rate"] = vatZero
	} else {
		// Call getTaxPercent
		vatTaxRate, err := getTaxPercent(tx, 10, 11, time.Now()) // 10.PPN 11.PPN
		if err != nil {
			return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate VAT tax rate",
				Err:        err,
			}
		}
		updates["vat_tax_rate"] = vatTaxRate
	}

	// Jika tidak ada perubahan, langsung return entity
	if len(updates) == 0 {
		return entity, nil
	}

	// Perbaikan: Gunakan Updates agar perubahan benar-benar tersimpan
	err = tx.Model(&entity).Updates(updates).Error
	if err != nil {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save the updated work order",
			Err:        err,
		}
	}

	return entity, nil
}

func getTaxPercent(tx *gorm.DB, taxTypeId int, taxServCode int, effDate time.Time) (float64, error) {
	var taxPercent float64

	// effective date (EFF_DATE)
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

	//--Tambahan validasi pada saat Status = New dan semua Item sudah memiliki Invoice
	//--Karena Status = Ready terjadi saat Operation di Alokasi
	if entity.WorkOrderStatusId == utils.WoStatNew {

		var allInvoice int
		// fetch wo transaction type api no charge
		getWorkOrderTrxTypeNoCharge, workOrderTypeErr := generalserviceapiutils.GetWoTransactionTypeByCode("N")
		if workOrderTypeErr != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: workOrderTypeErr.StatusCode,
				Message:    "Failed to fetch work order type data from external service",
				Err:        workOrderTypeErr.Err,
			}
		}

		var existsPrimary bool
		err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("work_order_operation_item_line").
			Where(
				"work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id <> ? AND substitute_type_id <> ?",
				Id,
				utils.WoStatClosed,
				getWorkOrderTrxTypeNoCharge.WoTransactionTypeId,
				0,
			).
			Limit(1).
			Find(&existsPrimary).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check primary condition",
				Err:        err,
			}
		}

		// set `AllInvoice = 0`
		if existsPrimary {
			allInvoice = 0
		} else {
			var existsSecondary bool
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select("work_order_operation_item_line").
				Where("work_order_system_number = ?", Id).
				Limit(1).
				Find(&existsSecondary).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to check secondary condition",
					Err:        err,
				}
			}

			if existsSecondary {
				// set `AllInvoice = 0`
				allInvoice = 0
			} else {
				// set `AllInvoice = 1`
				allInvoice = 1
			}
		}

		if entity.WorkOrderStatusId == utils.WoStatNew && allInvoice == 0 {
			// fetch wo transaction type api no charge
			getWorkOrderTrxTypeNoCharge, workOrderTypeErr := generalserviceapiutils.GetWoTransactionTypeByCode("N")
			if workOrderTypeErr != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: workOrderTypeErr.StatusCode,
					Message:    "Failed to fetch work order type data from external service",
					Err:        workOrderTypeErr.Err,
				}
			}

			var count int64
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id <> ? AND substitute_type_id <> ?",
					Id, utils.WoStatClosed, getWorkOrderTrxTypeNoCharge.WoTransactionTypeId, 0).
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
					StatusCode: http.StatusConflict,
					Message:    "Detail Work Order with out Invoice No must be deleted",
					Err:        err,
				}
			}

			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				UpdateColumns(map[string]interface{}{
					"work_order_status_id": utils.WoStatClosed,
				}).
				Where("work_order_system_number = ?", Id).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to update status work order, status service work order must be closed",
					Err:        err,
				}
			}

			//--DELETE WO FROM CAR_WASH
			err = tx.Model(&transactionjpcbentities.CarWash{}).
				Where("work_order_system_number = ?", Id).
				Delete(&transactionjpcbentities.CarWash{}).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to delete work order from car wash",
					Err:        err,
				}
			}

			// update service request
			if entity.ServiceRequestSystemNumber != 0 {
				// fetch service request status api closed
				getServiceRequestStatusClosed, serviceRequestStatusErr := generalserviceapiutils.GetServiceRequestStatusByCode("C")
				if serviceRequestStatusErr != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: serviceRequestStatusErr.StatusCode,
						Message:    "Failed to fetch service request status data from external service",
						Err:        serviceRequestStatusErr.Err,
					}
				}

				err = tx.Model(&transactionworkshopentities.ServiceRequest{}).
					Where("service_request_system_number = ?", entity.ServiceRequestSystemNumber).
					Update("service_request_status_id", getServiceRequestStatusClosed.ServiceRequestStatusID).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update service request status",
						Err:        err,
					}
				}

			} else if entity.PDISystemNumber != 0 && entity.PDILineNumber != 0 {
				var bookingSystemNo int64
				err := tx.Model(&transactionunitentities.PdiRequestDetail{}).
					Select("COALESCE(booking_system_number, 0)").
					Where("pdi_request_system_number = ? AND pdi_request_detail_line_number = ?", entity.PDISystemNumber, entity.PDILineNumber).
					Scan(&bookingSystemNo).Error

				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to retrieve booking system number from the database",
						Err:        err,
					}
				}

				var PDIStatusId int
				if bookingSystemNo == 0 {
					// fetch pdi status api accept if booking system number is 0
					getPDIStatusAccept, pdiStatusErr := generalserviceapiutils.GetPDIStatusByCode("20")
					if pdiStatusErr != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: pdiStatusErr.StatusCode,
							Message:    "Failed to fetch pdi status data from external service",
							Err:        pdiStatusErr.Err,
						}
					}

					PDIStatusId = getPDIStatusAccept.PDIStatusId

				} else {
					// fetch pdi status api booking
					getPDIStatusBooking, pdiStatusErr := generalserviceapiutils.GetPDIStatusByCode("40")
					if pdiStatusErr != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: pdiStatusErr.StatusCode,
							Message:    "Failed to fetch pdi status data from external service",
							Err:        pdiStatusErr.Err,
						}
					}

					PDIStatusId = getPDIStatusBooking.PDIStatusId
				}

				err = tx.Model(&transactionunitentities.PdiRequestDetail{}).
					Where("pdi_request_system_number = ? AND pdi_request_detail_line_number = ?", entity.PDISystemNumber, entity.PDILineNumber).
					Update("pdi_status_id", PDIStatusId).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to update pdi status",
						Err:        err,
					}
				}

			} else if entity.BookingSystemNumber != 0 && entity.EstimationSystemNumber != 0 && entity.PDISystemNumber != 0 && entity.PDILineNumber != 0 && entity.ServiceRequestSystemNumber != 0 {
				var BatchSystemNumber int64

				if entity.BookingSystemNumber != 0 {
					err = tx.Model(&transactionworkshopentities.BookingEstimation{}).
						Select("batch_system_number").
						Where("booking_system_number = ? ", entity.BookingSystemNumber).
						Scan(&BatchSystemNumber).Error
					if err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to retrieve batch system number from the database",
							Err:        err,
						}
					}
				} else {

					err = tx.Model(&transactionworkshopentities.BookingEstimation{}).
						Select("batch_system_number").
						Where("estimation_system_number = ? ", entity.EstimationSystemNumber).
						Scan(&BatchSystemNumber).Error
					if err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to retrieve batch system number from the database",
							Err:        err,
						}
					}
				}

				if BatchSystemNumber != 0 {

					err = tx.Model(&transactionworkshopentities.BookingEstimation{}).
						UpdateColumns(map[string]interface{}{
							"batch_status_id": 2,
						}).
						Where("batch_system_number = ?", BatchSystemNumber).Error
					if err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to update batch status",
							Err:        err,
						}
					}

				}
				//-- RPS/06/17/00677
				//-- mohon bantuan untuk perbaikan program saat WO sudah ada pembayaran kemudian WO diclose dari menu WO, pembayaran menjadi DPOT dan membentuk jurnal seperti WO Overpay.kalau perlu diskusi, nanti bisa dengan saya dan pak Santo juga.
				//-- If Work Order still has DP Payment not allocated for Invoice
				var exists bool
				err = tx.Model(&transactionworkshopentities.WorkOrder{}).
					Select("work_order_system_number").
					Where("work_order_system_number = ? AND work_order_status_id = ? AND (downpayment_payment - downpayment_payment_to_invoice) > 0", Id, utils.WoStatClosed).
					Find(&exists).Error
				if err != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to check if Work Order still has DP Payment not allocated for Invoice",
						Err:        err,
					}
				}
				if exists {

					// --Generate DP Other Payment
					//var RefType = "DPI08"
					var Total_Dp_Overpay float64

					err = tx.Model(&transactionworkshopentities.WorkOrder{}).
						Select("COALESCE(downpayment_payment - downpayment_payment_to_invoice, 0) as total_dp_overpay").
						Where("work_order_system_number = ? AND work_order_status_id = ? AND (downpayment_payment - downpayment_payment_to_invoice) > 0", Id, utils.WoStatClosed).
						Scan(&Total_Dp_Overpay).Error
					if err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to retrieve total DP Overpay from the database",
							Err:        err,
						}
					}

					// Generate ctDPIn (DP Other)
					// Call dbo.uspg_ctDPIn_Insert here
					// TODO: Implement logic for dbo.uspg_ctDPIn_Insert
					// @Option = 5, -- int
					// @Company_Code = @Company_Code, -- numeric
					// @Ref_Type = @Ref_Type, -- varchar(10)
					// @Ref_Sys_No = @Wo_Sys_No, -- numeric
					// @Ref_Doc_No = @Wo_Doc_No,
					// @Creation_User_Id = @Change_User_Id, -- varchar(10)
					// @Change_User_Id = @Change_User_Id, -- varchar(10)
					// @Total_Diff = @Total_Dp_Overpay, -- numeric
					// @Di_Date = @Change_Datetime ,
					// @Vehicle_Brand = @Vehicle_Brand
					// --End Generate DP Other--

					// update work order
					err = tx.Model(&transactionworkshopentities.WorkOrder{}).
						UpdateColumns(map[string]interface{}{
							"downpayment_payment_to_invoice": gorm.Expr("downpayment_payment"),
							"downpayment_overpay":            Total_Dp_Overpay,
						}).
						Where("work_order_system_number = ?", Id).Error
					if err != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to update work order DP Payment to Invoice",
							Err:        err,
						}
					}

					// case "dealer", "imsi":
					// 	eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO_D"
					// case "atpm", "salim", "maintained":
					// 	eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO_A"
					// default:
					// 	eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO"
					// }

					// Determine customer type and set event number
					getCustomerInfo, custErr := generalserviceapiutils.GetCustomerMasterDetailById(entity.CustomerId)
					if custErr != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: custErr.StatusCode,
							Message:    "Failed to fetch customer data from external service",
							Err:        custErr.Err,
						}
					}

					custTypeId := getCustomerInfo.ClientTypeId

					eventNo := ""
					switch custTypeId {
					case 1:
						eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO"
					case 2:
						eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO_D"
					case 3:
						eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO_A"
					default:
						eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO"
					}

					if eventNo == "" {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Event for Returning DP Customer to DP Other is not exists",
							Err:        err,
						}
					}

					// Generate Journal (DP Customer -> DP Other)
					// Call usp_comJournalAction here
					// TODO: Implement logic for usp_comJournalAction

					// Update JOURNAL_SYS_NO on DPOT
					// TODO: Implement logic for updating JOURNAL_SYS_NO on DPOT
					//}
				}

			}
		} else {
			// fetch wo transaction type api warranty
			getWorkOrderTrxTypeWarranty, workOrderTypeErr := generalserviceapiutils.GetWoTransactionTypeByCode("W")
			if workOrderTypeErr != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: workOrderTypeErr.StatusCode,
					Message:    "Failed to fetch work order type data from external service",
					Err:        workOrderTypeErr.Err,
				}
			}
			var allPtpSupply bool
			var count int64
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id = ? AND substitute_type_id <> ?",
					Id, utils.WoStatClosed, getWorkOrderTrxTypeWarranty.WoTransactionTypeId, 0).
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

				// --Jika ada bill code warranty,
				// 	--1.Validasi part to part harus tersupply semua
				// 	--2.Validasi part to money dan operation harus status close
				// fetch wo transaction type api warranty
				getWorkOrderTrxTypeWarranty, workOrderTypeErr = generalserviceapiutils.GetWoTransactionTypeByCode("W")
				if workOrderTypeErr != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: workOrderTypeErr.StatusCode,
						Message:    "Failed to fetch work order type data from external service",
						Err:        workOrderTypeErr.Err,
					}
				}

				// fetch warranty claim type api part to part
				getWarrantyClaimTypeResponses, warrantyClaimTypeErr := generalserviceapiutils.GetWarrantyClaimTypeByCode("PP")
				if warrantyClaimTypeErr != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: warrantyClaimTypeErr.StatusCode,
						Message:    "Failed to fetch warranty claim type data from external service",
						Err:        warrantyClaimTypeErr.Err,
					}
				}
				// --Validasi 1--
				// Validate part-to-part supply //cek statusid <> 8(closed), billcode <> warranty (6), substituteid , warrantyclaim_type = 0 (part), frt_qty > supply_qty
				err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id = ? AND substitute_type_id <> ? AND warranty_claim_type_id = ? AND frt_qty > supply_qty",
						Id, utils.WoStatClosed, getWorkOrderTrxTypeWarranty.WoTransactionTypeId, 0, getWarrantyClaimTypeResponses.WarrantyClaimTypeId).
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

				// fetch wo transaction type api warranty
				getWorkOrderTrxTypeWarranty, workOrderTypeErr = generalserviceapiutils.GetWoTransactionTypeByCode("W")
				if workOrderTypeErr != nil {
					return false, &exceptions.BaseErrorResponse{
						StatusCode: workOrderTypeErr.StatusCode,
						Message:    "Failed to fetch work order type data from external service",
						Err:        workOrderTypeErr.Err,
					}
				}
				// --Validasi 2--
				// Validate part-to-money and operation status //cek statusid <> 8(closed), billcode <> warranty (6), substituteid , warrantyclaim_type = 0 (part)
				err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
					Where("work_order_system_number = ? AND work_order_status_id <> ? AND transaction_type_id = ? AND substitute_type_id <> ? AND warranty_claim_type_id <> ?",
						Id, utils.WoStatClosed, getWorkOrderTrxTypeWarranty.WoTransactionTypeId, 0, 0).
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
			// Validate mileage and update vehicle master
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

			// --== UPDATE STATUS DI HEADER DAN DETAIL MENJADI CLOSED ==--
			// UPDATE wtWorkOrder0
			// SET	WO_STATUS = @WoStatClose ,
			// 	WO_CLOSE_DATE = @Change_Datetime ,
			// 	CHANGE_NO = CHANGE_NO + 1 ,
			// 	CHANGE_USER_ID = @Change_User_Id  ,
			// 	CHANGE_DATETIME = @Change_Datetime
			// WHERE WO_SYS_NO = @WO_SYS_NO

			// EXEC dbo.uspg_wtWorkOrder0_Update
			// @Option = 9,
			// @Wo_Sys_No = @Wo_Sys_No,
			// @Change_User_Id = @Change_User_Id

			// EXEC dbo.uspg_wtWorkOrder0_Update
			// @Option = 10,
			// @Wo_Sys_No = @Wo_Sys_No,
			// @Change_User_Id = @Change_User_Id

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

			// Generate ctDPIn (DP Other)
			// Call dbo.uspg_ctDPIn_Insert here
			// TODO: Implement logic for dbo.uspg_ctDPIn_Insert
			// @Option = 5, -- int
			// @Company_Code = @Company_Code, -- numeric
			// @Ref_Type = @Ref_Type, -- varchar(10)
			// @Ref_Sys_No = @Wo_Sys_No, -- numeric
			// @Ref_Doc_No = @Wo_Doc_No,
			// @Creation_User_Id = @Change_User_Id, -- varchar(10)
			// @Change_User_Id = @Change_User_Id, -- varchar(10)
			// @Total_Diff = @Total_Dp_Overpay, -- numeric
			// @Di_Date = @Change_Datetime ,
			// @Vehicle_Brand = @Vehicle_Brand
			// --End Generate DP Other--

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

			// case "dealer", "imsi":
			// 	eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO_D"
			// case "atpm", "salim", "maintained":
			// 	eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO_A"
			// default:
			// 	eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO"
			// }

			// Determine customer type and set event number
			getCustomerInfo, custErr := generalserviceapiutils.GetCustomerMasterDetailById(entity.CustomerId)
			if custErr != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: custErr.StatusCode,
					Message:    "Failed to fetch customer data from external service",
					Err:        custErr.Err,
				}
			}

			custTypeId := getCustomerInfo.ClientTypeId

			eventNo := ""
			switch custTypeId {
			case 1:
				eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO"
			case 2:
				eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO_D"
			case 3:
				eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO_A"
			default:
				eventNo = "GL_EVENT_NO_CLOSE_ORDER_WO"
			}

			if eventNo == "" {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Event for Returning DP Customer to DP Other is not exists",
					Err:        err,
				}
			}

			// Generate Journal (DP Customer -> DP Other)
			// Call usp_comJournalAction here
			// TODO: Implement logic for usp_comJournalAction

			// Update JOURNAL_SYS_NO on DPOT
			// TODO: Implement logic for updating JOURNAL_SYS_NO on DPOT
			//}

			entity.WorkOrderStatusId = utils.WoStatClosed
			err = tx.Save(&entity).Error
			if err != nil {
				return false, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to close the work order",
					Err:        err,
				}
			}
		}
	}

	return true, nil
}

// uspg_wtWorkOrder1_Insert
// IF @Option = 0
// --USE FOR : * INSERT NEW DATA
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (r *WorkOrderRepositoryImpl) GetAllRequest(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrderService

	baseModelQuery := tx.Model(&transactionworkshopentities.WorkOrderService{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order service requests",
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, entity := range entities {
		workOrderServiceData := map[string]interface{}{
			"work_order_service_id":     entity.WorkOrderServiceId,
			"work_order_system_number":  entity.WorkOrderSystemNumber,
			"work_order_service_remark": entity.WorkOrderServiceRemark,
		}
		results = append(results, workOrderServiceData)
	}

	pages.Rows = results

	return pages, nil
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

	var maxWonextLine int
	err = tx.Model(&transactionworkshopentities.WorkOrderService{}).
		Select("COALESCE(MAX(work_order_service_request_line), 0)").
		Where("work_order_system_number = ?", request.WorkOrderSystemNumber).
		Scan(&maxWonextLine).Error

	if err != nil {
		return transactionworkshopentities.WorkOrderService{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve maximum work order service item line",
			Err:        err,
		}
	}

	if maxWonextLine == 0 {
		maxWonextLine = 1
	} else {
		maxWonextLine++
	}

	loc, _ := time.LoadLocation("Asia/Jakarta") // UTC+7
	currentDate := time.Now().In(loc).Format("2006-01-02T15:04:05Z")
	parsedTime, _ := time.Parse(time.RFC3339, currentDate)

	entities := transactionworkshopentities.WorkOrderService{
		WorkOrderSystemNumber:       request.WorkOrderSystemNumber,
		WorkOrderServiceRemark:      request.WorkOrderServiceRemark,
		WorkOrderServiceDate:        parsedTime,
		WorkOrderServiceRequestLine: maxWonextLine,
	}

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
	if len(requests) > 5 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "A maximum of 5 requests can be added at once",
		}
	}

	var entities []transactionworkshopentities.WorkOrderService

	for _, request := range requests {
		var lastService transactionworkshopentities.WorkOrderService

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

		var maxWonextLine int
		err = tx.Model(&transactionworkshopentities.WorkOrderService{}).
			Select("COALESCE(MAX(work_order_service_request_line), 0)").
			Where("work_order_system_number = ?", request.WorkOrderSystemNumber).
			Scan(&maxWonextLine).Error

		if err != nil {
			return entities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve maximum work order service item line",
				Err:        err,
			}
		}

		if maxWonextLine == 0 {
			maxWonextLine = 1
		} else {
			maxWonextLine++
		}

		entity := transactionworkshopentities.WorkOrderService{
			WorkOrderSystemNumber:       request.WorkOrderSystemNumber,
			WorkOrderServiceRemark:      request.WorkOrderServiceRemark,
			WorkOrderServiceDate:        time.Now(),
			WorkOrderServiceRequestLine: maxWonextLine,
		}

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
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (r *WorkOrderRepositoryImpl) GetAllVehicleService(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrderServiceVehicle

	baseModelQuery := tx.Model(&transactionworkshopentities.WorkOrderServiceVehicle{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order service vehicle requests",
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, entity := range entities {
		workOrderServiceVehicleData := map[string]interface{}{
			"work_order_service_vehicle_id": entity.WorkOrderServiceVehicleId,
			"work_order_system_number":      entity.WorkOrderSystemNumber,
			"work_order_vehicle_date":       entity.WorkOrderVehicleDate,
			"work_order_vehicle_remark":     entity.WorkOrderVehicleRemark,
		}
		results = append(results, workOrderServiceVehicleData)
	}

	pages.Rows = results

	return pages, nil
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

	updates := make(map[string]interface{})

	if request.WorkOrderVehicleRemark != "" {
		updates["work_order_vehicle_remark"] = request.WorkOrderVehicleRemark
	}
	if !request.WorkOrderVehicleDate.IsZero() {
		updates["work_order_vehicle_date"] = request.WorkOrderVehicleDate
	}

	if len(updates) == 0 {
		return entity, nil
	}

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

	loc, _ := time.LoadLocation("Asia/Jakarta") // UTC+7
	currentDate := time.Now().In(loc).Format("2006-01-02T15:04:05Z")
	parsedTime, _ := time.Parse(time.RFC3339, currentDate)

	entities := transactionworkshopentities.WorkOrderServiceVehicle{

		WorkOrderSystemNumber:  request.WorkOrderSystemNumber,
		WorkOrderVehicleDate:   parsedTime,
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
	brandResponse, brandErr := generalserviceapiutils.GetBrandGenerateDoc(workOrder.BrandId)
	if brandErr != nil {
		return "", &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch brand data from external service",
			Err:        brandErr.Err,
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

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (r *WorkOrderRepositoryImpl) GetAllDetailWorkOrder(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var tableStruct []transactionworkshoppayloads.WorkOrderDetailRequest

	baseModelQuery := tx.Model(&transactionworkshopentities.WorkOrderDetail{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&tableStruct).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order details",
			Err:        err,
		}
	}

	if len(tableStruct) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var convertedResponses []transactionworkshoppayloads.WorkOrderDetailResponse

	for _, workOrderReq := range tableStruct {

		lineTypeResponse, linetypeErr := generalserviceapiutils.GetLineTypeById(workOrderReq.LineTypeId)
		if linetypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve line type from the external API",
				Err:        linetypeErr.Err,
			}
		}

		transactionTypeResponse, transactionTypeErr := generalserviceapiutils.GetWoTransactionTypeById(workOrderReq.TransactionTypeId)
		if transactionTypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve transaction type from the external API",
				Err:        transactionTypeErr.Err,
			}
		}

		jobTypeResponse, jobTypeErr := generalserviceapiutils.GetJobTransactionTypeById(workOrderReq.JobTypeId)
		if jobTypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve job type from the external API",
				Err:        jobTypeErr.Err,
			}
		}

		wcfTypeResponse, wcfTypeErr := generalserviceapiutils.GetWarrantyClaimTypeById(workOrderReq.AtpmWCFTypeId)
		if wcfTypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: wcfTypeErr.StatusCode,
				Message:    "Failed to retrieve wcf type from the external API",
				Err:        wcfTypeErr.Err,
			}
		}

		woStatus, woStatusErr := generalserviceapiutils.GetWorkOrderStatusById(workOrderReq.WorkorderStatusId)
		if woStatusErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: woStatusErr.StatusCode,
				Message:    "Failed to retrieve work order status from the external API",
				Err:        woStatusErr.Err,
			}
		}

		var whsGroup masterwarehouseentities.WarehouseGroup
		tx.Where("warehouse_group_id = ?", workOrderReq.WarehouseGroupId).Find(&whsGroup)

		var subsTypeName string
		if workOrderReq.SubstituteTypeId == 0 {
			subsTypeName = ""
		} else {
			subsType, subsTypeErr := generalserviceapiutils.GetSubstituteTypeById(workOrderReq.SubstituteTypeId)
			if subsTypeErr != nil {
				if subsTypeErr.StatusCode == http.StatusNotFound {
					subsTypeName = ""
				} else {
					return pages, &exceptions.BaseErrorResponse{
						StatusCode: subsTypeErr.StatusCode,
						Message:    "Failed to retrieve substitute type from the external API",
						Err:        subsTypeErr.Err,
					}
				}
			} else {
				subsTypeName = subsType.SubstituteTypeName
			}
		}

		var OperationItemCode string
		var Description string

		operationItemResponse, operationItemErr := r.GetOperationItemById(workOrderReq.LineTypeId, workOrderReq.OperationItemId)
		if operationItemErr != nil {
			return pages, operationItemErr
		}

		OperationItemCode, Description, errResp := r.HandleLineTypeResponse(workOrderReq.LineTypeId, operationItemResponse)
		if errResp != nil {
			return pages, errResp
		}

		// fetch data item
		var itemResponse masteritementities.Item
		if err := tx.Where("item_id = ?", workOrderReq.OperationItemId).First(&itemResponse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item not found",
					Err:        fmt.Errorf("item with ID %d not found", workOrderReq.OperationItemId),
				}
			}
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item",
				Err:        err,
			}
		}

		// Fetch data UOM from external API
		var uomItems masteritementities.Uom
		if err := tx.Where("uom_id = ?", itemResponse.UnitOfMeasurementStockId).First(&uomItems).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "UOM not found",
					Err:        fmt.Errorf("uom with ID %d not found", itemResponse.UnitOfMeasurementStockId),
				}

			}
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch UOM",
				Err:        err,
			}
		}

		workOrderRes := transactionworkshoppayloads.WorkOrderDetailResponse{
			WorkOrderDetailId:                   workOrderReq.WorkOrderDetailId,
			WorkOrderSystemNumber:               workOrderReq.WorkOrderSystemNumber,
			LineTypeId:                          workOrderReq.LineTypeId,
			LineTypeCode:                        lineTypeResponse.LineTypeCode,
			LineTypeName:                        lineTypeResponse.LineTypeName,
			TransactionTypeId:                   workOrderReq.TransactionTypeId,
			TransactionTypeCode:                 transactionTypeResponse.WoTransactionTypeCode,
			JobTypeId:                           workOrderReq.JobTypeId,
			JobTypeCode:                         jobTypeResponse.JobTypeCode,
			FrtQuantity:                         workOrderReq.FrtQuantity,
			SupplyQuantity:                      workOrderReq.SupplyQuantity,
			OperationItemId:                     workOrderReq.OperationItemId,
			OperationItemCode:                   OperationItemCode,
			Description:                         Description,
			Uom:                                 uomItems.UomDescription,
			OperationItemPrice:                  workOrderReq.OperationItemPrice,
			OperationItemDiscountAmount:         workOrderReq.OperationItemDiscountAmount,
			OperationItemDiscountRequestAmount:  workOrderReq.OperationItemDiscountRequestAmount,
			OperationItemDiscountPercent:        workOrderReq.OperationItemDiscountPercent,
			OperationItemDiscountRequestPercent: workOrderReq.OperationItemDiscountRequestPercent,
			WorkOrderStatusName:                 woStatus.WorkOrderStatusName,
			InvoiceSystemNumber:                 workOrderReq.InvoiceSystemNumber,
			TechnicianId:                        workOrderReq.TechnicianId,
			TechnicianName:                      "",
			ForemanName:                         "",
			SubstituteTypeId:                    workOrderReq.SubstituteTypeId,
			SubstituteTypeName:                  subsTypeName,
			AtpmWCFTypeId:                       workOrderReq.AtpmWCFTypeId,
			WarrantyClaimTypeDescription:        wcfTypeResponse.WarrantyClaimTypeDescription,
			QualityControlPassDatetime:          nil,
			InvoiceDate:                         nil,
			ClaimNumber:                         "",
			Package:                             "",
			WarehouseGroupId:                    workOrderReq.WarehouseGroupId,
			WarehouseGroupName:                  whsGroup.WarehouseGroupName,
			PendingReason:                       "",
			RecSystemNumber:                     0,
		}

		convertedResponses = append(convertedResponses, workOrderRes)
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {

		totalDetails, totalErr := r.CalculateWorkOrderTotal(tx, response.WorkOrderSystemNumber)
		if totalErr != nil {
			return pages, totalErr
		}

		var totalValue float64
		for _, total := range totalDetails {
			if val, ok := total["total"].(float64); ok {
				totalValue = val
				break
			}
		}

		responseMap := map[string]interface{}{
			"work_order_detail_id":                   response.WorkOrderDetailId,
			"work_order_system_number":               response.WorkOrderSystemNumber,
			"service_status":                         response.WorkOrderStatusName,
			"line_type_id":                           response.LineTypeId,
			"line_type_code":                         response.LineTypeCode,
			"line_type_name":                         response.LineTypeName,
			"transaction_type_id":                    response.TransactionTypeId,
			"transaction_type_code":                  response.TransactionTypeCode,
			"job_type_id":                            response.JobTypeId,
			"job_type_code":                          response.JobTypeCode,
			"frt_quantity":                           response.FrtQuantity,
			"supply_quantity":                        response.SupplyQuantity,
			"operation_item_id":                      response.OperationItemId,
			"operation_item_code":                    response.OperationItemCode,
			"description":                            response.Description,
			"uom":                                    response.Uom,
			"operation_item_price":                   response.OperationItemPrice,
			"operation_item_discount_amount":         response.OperationItemDiscountAmount,
			"operation_item_discount_request_amount": response.OperationItemDiscountRequestAmount,
			"operation_item_discount_percent":        response.OperationItemDiscountPercent,
			"total":                                  totalValue,
			"sub_total":                              totalValue,
			"invoice_system_number":                  response.InvoiceSystemNumber,
			"technician_id":                          response.TechnicianId,
			"technician_name":                        response.TechnicianName,
			"foreman_id":                             response.ForemanId,
			"foreman_name":                           response.ForemanName,
			"substitute_type_id":                     response.SubstituteTypeId,
			"substitute_type_name":                   response.SubstituteTypeName,
			"atpm_wcf_type_id":                       response.AtpmWCFTypeId,
			"atpm_wcf_type_description":              response.WarrantyClaimTypeDescription,
			"warehouse_group_id":                     response.WarehouseGroupId,
			"warehouse_group_name":                   response.WarehouseGroupName,
			"pending_reason":                         response.PendingReason,
			"recall_system_number":                   response.RecSystemNumber,
			"quality_control_pass_datetime":          response.QualityControlPassDatetime,
			"invoice_date":                           response.InvoiceDate,
			"claim_number":                           response.ClaimNumber,
			"package":                                response.Package,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	pages.Rows = mapResponses

	return pages, nil
}

// uspg_wtWorkOrder2_Insert
// IF @Option = 0
// --USE FOR : * INSERT NEW DATA
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

	// Fetch external data

	lineTypeResponse, lineErr := generalserviceapiutils.GetLineTypeById(entity.LineTypeId)
	if lineErr != nil {
		return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve line type from the external API",
			Err:        lineErr.Err,
		}
	}

	transactionTypeResponse, transactionTypeErr := generalserviceapiutils.GetWoTransactionTypeById(entity.TransactionTypeId)
	if transactionTypeErr != nil {
		return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve transaction type from the external API",
			Err:        transactionTypeErr.Err,
		}
	}

	jobTypeResponse, jobTypeErr := generalserviceapiutils.GetJobTransactionTypeById(entity.JobTypeId)
	if jobTypeErr != nil {
		return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve job type from the external API",
			Err:        jobTypeErr.Err,
		}
	}

	wcfTypeResponse, wcfTypeErr := generalserviceapiutils.GetWarrantyClaimTypeById(entity.AtpmWCFTypeId)
	if wcfTypeErr != nil {
		return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: wcfTypeErr.StatusCode,
			Message:    "Failed to retrieve wcf type from the external API",
			Err:        wcfTypeErr.Err,
		}
	}

	woStatus, woStatusErr := generalserviceapiutils.GetWorkOrderStatusById(entity.WorkorderStatusId)
	if woStatusErr != nil {
		return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: woStatusErr.StatusCode,
			Message:    "Failed to retrieve work order status from the external API",
			Err:        woStatusErr.Err,
		}
	}

	var whsGroup masterwarehouseentities.WarehouseGroup
	if err := tx.Where("warehouse_group_id = ?", entity.WarehouseGroupId).First(&whsGroup).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Warehouse group not found",
				Err:        fmt.Errorf("warehouse group with ID %d not found", entity.WarehouseGroupId),
			}
		}
		return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Warehouse group",
			Err:        err,
		}
	}

	var subsTypeName string
	if entity.SubstituteTypeId == 0 {
		subsTypeName = ""
	} else {
		subsType, subsTypeErr := generalserviceapiutils.GetSubstituteTypeById(entity.SubstituteTypeId)
		if subsTypeErr != nil {
			if subsTypeErr.StatusCode == http.StatusNotFound {
				subsTypeName = ""
			} else {
				return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
					StatusCode: subsTypeErr.StatusCode,
					Message:    "Failed to retrieve substitute type from the external API",
					Err:        subsTypeErr.Err,
				}
			}
		} else {
			subsTypeName = subsType.SubstituteTypeName
		}
	}

	var OperationItemCode string
	var Description string

	operationItemResponse, operationItemErr := r.GetOperationItemById(entity.LineTypeId, entity.OperationItemId)
	if operationItemErr != nil {
		return transactionworkshoppayloads.WorkOrderDetailResponse{}, operationItemErr
	}

	OperationItemCode, Description, errResp := r.HandleLineTypeResponse(entity.LineTypeId, operationItemResponse)
	if errResp != nil {
		return transactionworkshoppayloads.WorkOrderDetailResponse{}, errResp
	}

	// fetch data item
	var itemResponse masteritementities.Item
	if err := tx.Where("item_id = ?", entity.OperationItemId).First(&itemResponse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item not found",
				Err:        fmt.Errorf("item with ID %d not found", entity.OperationItemId),
			}
		}
		return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item",
			Err:        err,
		}
	}

	// Fetch data UOM from external API
	var uomItems masteritementities.Uom
	if err := tx.Where("uom_id = ?", itemResponse.UnitOfMeasurementStockId).First(&uomItems).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "UOM not found",
				Err:        fmt.Errorf("uom with ID %d not found", itemResponse.UnitOfMeasurementStockId),
			}

		}
		return transactionworkshoppayloads.WorkOrderDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch UOM",
			Err:        err,
		}
	}

	payload := transactionworkshoppayloads.WorkOrderDetailResponse{
		WorkOrderDetailId:                   entity.WorkOrderDetailId,
		WorkOrderSystemNumber:               entity.WorkOrderSystemNumber,
		LineTypeId:                          entity.LineTypeId,
		LineTypeCode:                        lineTypeResponse.LineTypeCode,
		LineTypeName:                        lineTypeResponse.LineTypeName,
		TransactionTypeId:                   entity.TransactionTypeId,
		TransactionTypeCode:                 transactionTypeResponse.WoTransactionTypeCode,
		JobTypeId:                           entity.JobTypeId,
		JobTypeCode:                         jobTypeResponse.JobTypeCode,
		FrtQuantity:                         entity.FrtQuantity,
		SupplyQuantity:                      entity.SupplyQuantity,
		PriceListId:                         entity.PriceListId,
		WarehouseGroupId:                    entity.WarehouseGroupId,
		WarehouseGroupName:                  whsGroup.WarehouseGroupName,
		OperationItemId:                     entity.OperationItemId,
		OperationItemCode:                   OperationItemCode,
		Description:                         Description,
		Uom:                                 uomItems.UomDescription,
		OperationItemPrice:                  entity.OperationItemPrice,
		OperationItemDiscountAmount:         entity.OperationItemDiscountAmount,
		OperationItemDiscountRequestAmount:  entity.OperationItemDiscountRequestAmount,
		OperationItemDiscountPercent:        entity.OperationItemDiscountPercent,
		OperationItemDiscountRequestPercent: entity.OperationItemDiscountRequestPercent,
		WorkOrderStatusName:                 woStatus.WorkOrderStatusName,
		AtpmWCFTypeId:                       entity.AtpmWCFTypeId,
		WarrantyClaimTypeDescription:        wcfTypeResponse.WarrantyClaimTypeDescription,
		InvoiceSystemNumber:                 entity.InvoiceSystemNumber,
		TechnicianId:                        entity.TechnicianId,
		SubstituteTypeId:                    entity.SubstituteTypeId,
		SubstituteTypeName:                  subsTypeName,
		QualityControlPassDatetime:          entity.QualityControlPassDatetime,
	}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) CalculateWorkOrderTotal(tx *gorm.DB, workOrderSystemNumber int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {

	type Result struct {
		TotalPackage            float64
		TotalOperation          float64
		TotalSparePart          float64
		TotalOil                float64
		TotalMaterial           float64
		TotalAccessories        float64
		TotalConsumableMaterial float64
		TotalSublet             float64
		TotalSouvenir           float64
	}

	var result Result

	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select(`
			SUM(CASE WHEN line_type_id = 1 THEN ROUND(COALESCE(operation_item_price, 0), 0) ELSE 0 END) AS total_package,
			SUM(CASE WHEN line_type_id = 2 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_operation,
			SUM(CASE WHEN line_type_id = 3 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_spare_part,
			SUM(CASE WHEN line_type_id = 4 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_oil,
			SUM(CASE WHEN line_type_id = 5 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_material,
			SUM(CASE WHEN line_type_id = 6 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_sublet,
			SUM(CASE WHEN line_type_id = 7 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_accessories,
			SUM(CASE WHEN line_type_id = 8 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_consumable_material,
			SUM(CASE WHEN line_type_id = 9 THEN ROUND(COALESCE(operation_item_price, 0) * COALESCE(frt_quantity, 0), 0) ELSE 0 END) AS total_souvenir
		`).
		Where("work_order_system_number = ?", workOrderSystemNumber).
		Scan(&result).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate work order total",
			Err:        err,
		}
	}

	// Calculate grand total
	grandTotal := result.TotalPackage + result.TotalOperation + result.TotalSparePart + result.TotalOil + result.TotalMaterial + result.TotalAccessories + result.TotalConsumableMaterial + result.TotalSublet + result.TotalSouvenir

	// Update Work Order with the calculated totals
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).
		Where("work_order_system_number = ?", workOrderSystemNumber).
		Updates(map[string]interface{}{
			"total_package":             result.TotalPackage,
			"total_operation":           result.TotalOperation,
			"total_part":                result.TotalSparePart,
			"total_oil":                 result.TotalOil,
			"total_material":            result.TotalMaterial,
			"total_sublet":              result.TotalSublet,
			"total_accessories":         result.TotalAccessories,
			"total_consumable_material": result.TotalConsumableMaterial,
			"total_souvenir":            result.TotalSouvenir,
			"total":                     grandTotal,
		}).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to update work order: %v", err)}
	}

	workOrderDetailResponses := []map[string]interface{}{
		{"total_package": result.TotalPackage},
		{"total_operation": result.TotalOperation},
		{"total_part": result.TotalSparePart},
		{"total_oil": result.TotalOil},
		{"total_material": result.TotalMaterial},
		{"total_sublet": result.TotalSublet},
		{"total_accessories": result.TotalAccessories},
		{"total_consumable_material": result.TotalConsumableMaterial},
		{"total_souvenir": result.TotalSouvenir},
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

			if maxWoOprItemLine == 0 {
				maxWoOprItemLine = 1
			} else {
				maxWoOprItemLine++
			}

			var bookingEstim transactionworkshopentities.BookingEstimation
			var bookingEstimDetail transactionworkshopentities.BookingEstimationDetail

			err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
				Select("trx_booking_estimation.*, trx_booking_estimation_detail.line_type_id, trx_booking_estimation_detail.transaction_type_id, trx_booking_estimation_detail.job_type_id, trx_booking_estimation_detail.operation_item_code, trx_booking_estimation_detail.frt_quantity, trx_booking_estimation_detail.operation_item_price, trx_booking_estimation_detail.operation_item_discount_amount, trx_booking_estimation_detail.operation_item_discount_request_amount, trx_booking_estimation_detail.operation_item_discount_percent, trx_booking_estimation_detail.operation_item_discount_request_percent, trx_booking_estimation_detail.estimation_line").
				Joins("LEFT JOIN trx_booking_estimation_detail ON trx_booking_estimation.estimation_system_number = trx_booking_estimation_detail.estimation_system_number").
				Where("trx_booking_estimation.estimation_system_number = ?", estimSystemNo).
				First(&bookingEstim).Error

			if err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to retrieve booking estimation data",
					Err:        err,
				}
			}

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:               id,
				LineTypeId:                          bookingEstimDetail.LineTypeId,        // BE0.LineTypeId,
				TransactionTypeId:                   bookingEstimDetail.TransactionTypeId, // utils.TrxTypeWoExternal,
				JobTypeId:                           bookingEstimDetail.JobTypeId,         // CASE WHEN BE0.CPC_CODE = @Profit_Center_BR THEN @JobTypeBR ELSE @JobTypePM END,
				OperationItemId:                     bookingEstimDetail.OperationItemId,   // BE0.OperationItemCode,
				OperationItemCode:                   bookingEstimDetail.OperationItemCode, //BE.OPR_ITEM_CODE,
				WarehouseGroupId:                    0,                                    //Whs_Group_Sp
				FrtQuantity:                         bookingEstimDetail.FRTQuantity,       //BE.FrtQuantity,
				SupplyQuantity:                      0,                                    //CASE WHEN BE.LINE_TYPE = @LINETYPE_OPR OR BE.LINE_TYPE = @LINETYPE_PACKAGE THEN BE.FRT_QTY ELSE CASE WHEN I.ITEM_TYPE = @ItemTypeService AND I.ITEM_GROUP <> @ItemGrpOJ THEN BE.FRT_QTY ELSE 0 END END
				WorkorderStatusId:                   utils.WoStatDraft,
				OperationItemDiscountAmount:         bookingEstimDetail.OperationItemDiscountAmount,         //BE.OPR_ITEM_DISC_AMOUNT,
				OperationItemDiscountRequestAmount:  bookingEstimDetail.OperationItemDiscountRequestAmount,  //BE.OPR_ITEM_DISC_REQ_AMOUNT,
				OperationItemDiscountPercent:        bookingEstimDetail.OperationItemDiscountPercent,        //BE.OPR_ITEM_DISC_PERCENT,
				OperationItemDiscountRequestPercent: bookingEstimDetail.OperationItemDiscountRequestPercent, //BE.OPR_ITEM_DISC_REQ_PERCENT,
				OperationItemPrice:                  bookingEstimDetail.OperationItemPrice,                  //BE.OPR_ITEM_PRICE,
				PphAmount:                           0,                                                      //BE.PPH_AMOUNT,
				PphTaxRate:                          0,                                                      //BE.PPH_TAX_RATE,
				AtpmWCFTypeId:                       0,                                                      //CASE WHEN BE.LINE_TYPE = @LINETYPE_OPR OR BE.LINE_TYPE = @LINETYPE_PACKAGE THEN '' ELSE ATPM_WCF_TYPE END
				WorkOrderOperationItemLine:          maxWoOprItemLine,                                       //BE.ESTIM_LINE,
				ServiceStatusId:                     utils.SrvStatDraft,
			}

			if err := tx.Create(&workOrderDetail).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to create work order detail",
					Err:        err,
				}
			}

			if _, err := r.CalculateWorkOrderTotal(tx, id); err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, err
			}

		} else {

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

			if maxWoOprItemLine == 0 {
				maxWoOprItemLine = 1
			} else {
				maxWoOprItemLine++
			}

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:               id,
				LineTypeId:                          request.LineTypeId,
				TransactionTypeId:                   request.TransactionTypeId,
				JobTypeId:                           request.JobTypeId,
				OperationItemId:                     request.OperationItemId,
				WarehouseGroupId:                    request.WarehouseGroupId, // Whs_Group_Sp
				FrtQuantity:                         request.FrtQuantity,      // BE.FrtQuantity,
				SupplyQuantity:                      request.SupplyQuantity,   // CASE WHEN BE.LINE_TYPE = @LINETYPE_OPR OR BE.LINE_TYPE = @LINETYPE_PACKAGE THEN BE.FRT_QTY ELSE CASE WHEN I.ITEM_TYPE = @ItemTypeService AND I.ITEM_GROUP <> @ItemGrpOJ THEN BE.FRT_QTY ELSE 0 END END
				WorkorderStatusId:                   utils.WoStatDraft,
				OperationItemDiscountAmount:         0,                          // BE.OPR_ITEM_DISC_AMOUNT,
				OperationItemDiscountRequestAmount:  0,                          // BE.OPR_ITEM_DISC_REQ_AMOUNT,
				OperationItemDiscountPercent:        0,                          // BE.OPR_ITEM_DISC_PERCENT,
				OperationItemDiscountRequestPercent: 0,                          // BE.OPR_ITEM_DISC_REQ_PERCENT,
				OperationItemPrice:                  request.OperationItemPrice, // BE.OPR_ITEM_PRICE,
				PphAmount:                           0,                          // BE.PPH_AMOUNT,
				PphTaxRate:                          0,                          // BE.PPH_TAX_RATE,
				AtpmWCFTypeId:                       request.AtpmWCFTypeId,      // CASE WHEN BE.LINE_TYPE = @LINETYPE_OPR OR BE.LINE_TYPE = @LINETYPE_PACKAGE THEN '' ELSE ATPM_WCF_TYPE END
				WorkOrderOperationItemLine:          maxWoOprItemLine,
				ServiceStatusId:                     utils.SrvStatDraft,
			}

			if err := tx.Create(&workOrderDetail).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to create work order detail",
					Err:        err,
				}
			}

			if _, err := r.CalculateWorkOrderTotal(tx, id); err != nil {
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

			if maxWoOprItemLine == 0 {
				maxWoOprItemLine = 1
			} else {
				maxWoOprItemLine++
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

			// fetch linetype from campaign items
			linetypeCheck, LinetypeErr := generalserviceapiutils.GetLineTypeByCode(campaignItems[0].LineTypeCode)
			if LinetypeErr != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: LinetypeErr.StatusCode,
					Message:    "Failed to retrieve line type from the external API",
					Err:        LinetypeErr.Err,
				}
			}

			if len(campaignItems) > 0 {
				workOrderDetail = transactionworkshopentities.WorkOrderDetail{
					WorkOrderSystemNumber:               id,
					LineTypeId:                          linetypeCheck.LineTypeId,                       // C1.LINE_TYPE,
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
					WorkOrderOperationItemLine:          maxWoOprItemLine, // 0
					ServiceStatusId:                     utils.SrvStatDraft,
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

			if _, err := r.CalculateWorkOrderTotal(tx, id); err != nil {
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

			if maxWoOprItemLine == 0 {
				maxWoOprItemLine = 1
			} else {
				maxWoOprItemLine++
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

			// fetch linetype
			linetypeCheck, LinetypeErr := generalserviceapiutils.GetLineTypeByCode(utils.LinetypeOperation)
			if LinetypeErr != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: LinetypeErr.StatusCode,
					Message:    "Failed to retrieve line type from the external API",
					Err:        LinetypeErr.Err,
				}
			}

			workOrderDetail = transactionworkshopentities.WorkOrderDetail{
				WorkOrderSystemNumber:               id,
				LineTypeId:                          linetypeCheck.LineTypeId, // LINETYPE_OPR,
				TransactionTypeId:                   0,                        // dbo.FCT_getBillCode(@COMPANY_CODE ,CAST(P1.COMPANY_CODE AS VARCHAR(10)),'W'),
				JobTypeId:                           8,                        // dbo.getVariableValue('JOBTYPE_PDI'),
				OperationItemCode:                   "",                       // P1.OPERATION_NO,
				WarehouseGroupId:                    38,                       // Whs_Group_Sp
				FrtQuantity:                         0,                        // P1.FRT,
				SupplyQuantity:                      0,                        // P1.FRT
				WorkorderStatusId:                   0,                        // ""
				OperationItemDiscountAmount:         0,                        // 0,
				OperationItemDiscountRequestAmount:  0,                        // 0,
				OperationItemDiscountPercent:        0,                        // 0,
				OperationItemDiscountRequestPercent: 0,                        // 0,
				OperationItemPrice:                  0,                        // LSP1.SELLING_PRICE,
				PphAmount:                           0,                        // 0,
				PphTaxRate:                          0,                        // OP.TAX_CODE,
				AtpmWCFTypeId:                       0,                        // 0
				WorkOrderOperationItemLine:          maxWoOprItemLine,
				ServiceStatusId:                     utils.SrvStatDraft,
			}

			if err := tx.Create(&workOrderDetail).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to create work order detail",
					Err:        err,
				}
			}

			if _, err := r.CalculateWorkOrderTotal(tx, id); err != nil {
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

			if maxWoOprItemLine == 0 {
				maxWoOprItemLine = 1
			} else {
				maxWoOprItemLine++
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
				WorkOrderOperationItemLine:          maxWoOprItemLine,
				ServiceStatusId:                     utils.SrvStatDraft,
			}

			if err := tx.Create(&workOrderDetail).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to create work order detail",
					Err:        err,
				}
			}

			if _, err := r.CalculateWorkOrderTotal(tx, id); err != nil {
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

			if maxWoOprItemLine == 0 {
				maxWoOprItemLine = 1
			} else {
				maxWoOprItemLine++
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
				WorkOrderOperationItemLine:          maxWoOprItemLine,
				ServiceStatusId:                     utils.SrvStatDraft,
			}

			if err := tx.Create(&workOrderDetail).Error; err != nil {
				return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to create work order detail",
					Err:        err,
				}
			}

			if _, err := r.CalculateWorkOrderTotal(tx, id); err != nil {
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

	updates := make(map[string]interface{})
	if request.LineTypeId != 0 {
		updates["line_type_id"] = request.LineTypeId
	}
	if request.TransactionTypeId != 0 {
		updates["transaction_type_id"] = request.TransactionTypeId
	}
	if request.JobTypeId != 0 {
		updates["job_type_id"] = request.JobTypeId
	}
	if request.WarehouseGroupId != 0 {
		updates["warehouse_group_id"] = request.WarehouseGroupId
	}
	if request.OperationItemId != 0 {
		updates["operation_item_id"] = request.OperationItemId
	}
	if request.FrtQuantity != 0 {
		updates["frt_quantity"] = request.FrtQuantity
	}
	if request.SupplyQuantity != 0 {
		updates["supply_quantity"] = request.SupplyQuantity
	}
	if request.PriceListId != 0 {
		updates["price_list_id"] = request.PriceListId
	}
	if request.OperationItemDiscountRequestAmount != 0 {
		updates["operation_item_discount_request_amount"] = request.OperationItemDiscountRequestAmount
	}
	if request.OperationItemPrice != 0 {
		updates["operation_item_price"] = request.OperationItemPrice
	}

	if len(updates) == 0 {
		return entity, nil
	}

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
	_, calcErr := r.CalculateWorkOrderTotal(tx, id)
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
	vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(request.VehicleId)
	if vehicleErr != nil {
		return transactionworkshopentities.WorkOrder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        vehicleErr.Err,
		}
	}

	// Create WorkOrder entity
	entitieswo := transactionworkshopentities.WorkOrder{
		// Default values
		WorkOrderDocumentNumber: defaultWorkOrderDocumentNumber,
		WorkOrderStatusId:       utils.WoStatDraft,
		WorkOrderDate:           currentDate,
		CPCcode:                 defaultCPCcode,
		ServiceAdvisorId:        defaultServiceAdvisorId,
		WorkOrderTypeId:         workOrderTypeId,
		BookingSystemNumber:     request.BookingSystemNumber,
		EstimationSystemNumber:  request.EstimationSystemNumber,
		ServiceSite:             "OD - Service On Dealer",
		VehicleChassisNumber:    vehicleResponses.Data.Master.VehicleChassisNumber,

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

func (r *WorkOrderRepositoryImpl) GetAllBooking(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var tableStruct []transactionworkshoppayloads.WorkOrderBooking

	baseModelQuery := tx.Model(&transactionworkshoppayloads.WorkOrderBooking{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	whereQuery = whereQuery.Where("booking_system_number != 0 OR estimation_system_number != 0")

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&tableStruct).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order bookings",
			Err:        err,
		}
	}

	if len(tableStruct) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var convertedResponses []transactionworkshoppayloads.WorkOrderBookingResponse

	for _, workOrderReq := range tableStruct {
		var (
			workOrderRes transactionworkshoppayloads.WorkOrderBookingResponse
		)

		getBrandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(workOrderReq.BrandId)
		if brandErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: brandErr.StatusCode,
				Message:    "Failed to fetch brand data from external service",
				Err:        brandErr.Err,
			}
		}

		// Fetch external data for model
		getModelResponse, modelErr := salesserviceapiutils.GetUnitModelById(workOrderReq.ModelId)
		if modelErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: modelErr.StatusCode,
				Message:    "Failed to fetch model data from external service",
				Err:        modelErr.Err,
			}
		}

		// Fetch external data for vehicle
		vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(workOrderReq.VehicleId)
		if vehicleErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: vehicleErr.StatusCode,
				Message:    "Failed to retrieve vehicle data from external API",
				Err:        vehicleErr.Err,
			}
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
			VehicleCode:                vehicleResponses.Data.Master.VehicleChassisNumber,
			VehicleTnkb:                vehicleResponses.Data.STNK.VehicleRegistrationCertificateTNKB,
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
			"batch_system_number":           "",
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

	pages.Rows = mapResponses

	return pages, nil
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
	brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
	if brandErr != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve brand data from the external API",
			Err:        brandErr.Err,
		}
	}

	// Fetch data model from external API
	modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
	if modelErr != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from the external API",
			Err:        modelErr.Err,
		}
	}

	// Fetch data variant from external API
	variantResponse, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
	if variantErr != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        variantErr.Err,
		}
	}

	// Fetch data vehicle from external API
	vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
	if vehicleErr != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        vehicleErr.Err,
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
	getWorkOrderStatusResponses, workOrderStatusErr := generalserviceapiutils.GetWorkOrderStatusById(entity.WorkOrderStatusId)
	if workOrderStatusErr != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order status data from external service",
			Err:        workOrderStatusErr.Err,
		}
	}

	// fetch data type work order
	getWorkOrderTypeResponses, workOrderTypeErr := generalserviceapiutils.GetWorkOrderTypeByID(entity.WorkOrderTypeId)
	if workOrderTypeErr != nil {
		return transactionworkshoppayloads.WorkOrderBookingResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order type data from external service",
			Err:        workOrderTypeErr.Err,
		}
	}

	payload := transactionworkshoppayloads.WorkOrderBookingResponse{
		WorkOrderSystemNumber:         entity.WorkOrderSystemNumber,
		WorkOrderDate:                 entity.WorkOrderDate.Format("2006-01-02"),
		WorkOrderDocumentNumber:       entity.WorkOrderDocumentNumber,
		WorkOrderTypeId:               entity.WorkOrderTypeId,
		WorkOrderTypeName:             getWorkOrderTypeResponses.WorkOrderTypeName,
		WorkOrderStatusId:             entity.WorkOrderStatusId,
		WorkOrderStatusName:           getWorkOrderStatusResponses.WorkOrderStatusName,
		ServiceAdvisorId:              entity.ServiceAdvisorId,
		BrandId:                       entity.BrandId,
		BrandName:                     brandResponse.BrandName,
		ModelId:                       entity.ModelId,
		ModelName:                     modelResponse.ModelName,
		VariantId:                     entity.VariantId,
		VariantDescription:            variantResponse.VariantDescription,
		VehicleId:                     entity.VehicleId,
		VehicleCode:                   vehicleResponses.Data.Master.VehicleChassisNumber,
		VehicleTnkb:                   vehicleResponses.Data.STNK.VehicleRegistrationCertificateTNKB,
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

func (r *WorkOrderRepositoryImpl) GetAllAffiliated(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var tableStruct []transactionworkshoppayloads.WorkOrderAffiliate

	baseModelQuery := tx.Model(&transactionworkshoppayloads.WorkOrderAffiliate{})

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	whereQuery = whereQuery.Where("service_request_system_number != 0")

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&tableStruct).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order affiliated data",
			Err:        err,
		}
	}

	if len(tableStruct) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var convertedResponses []transactionworkshoppayloads.WorkOrderAffiliateGetResponse

	for _, workOrderReq := range tableStruct {
		var workOrderRes transactionworkshoppayloads.WorkOrderAffiliateGetResponse

		// Fetch external data for brand
		getBrandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(workOrderReq.BrandId)
		if brandErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: brandErr.StatusCode,
				Message:    "Failed to fetch brand data from external service",
				Err:        brandErr.Err,
			}
		}

		// Fetch external data for model
		getModelResponse, modelErr := salesserviceapiutils.GetUnitModelById(workOrderReq.ModelId)
		if modelErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: modelErr.StatusCode,
				Message:    "Failed to fetch model data from external service",
				Err:        modelErr.Err,
			}
		}

		// Fetch external data for vehicle
		getVehicleResponse, vehicleErr := salesserviceapiutils.GetVehicleById(workOrderReq.VehicleId)
		if vehicleErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: vehicleErr.StatusCode,
				Message:    "Failed to retrieve vehicle data from external API",
				Err:        vehicleErr.Err,
			}
		}

		// Fetch service request data from internal service
		serviceRequestURL := config.EnvConfigs.AfterSalesServiceUrl + "service-request/" + strconv.Itoa(workOrderReq.ServiceRequestSystemNumber)
		var getServiceRequestResponse transactionworkshoppayloads.ServiceRequestResponse
		if err := utils.Get(serviceRequestURL, &getServiceRequestResponse, nil); err != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch service request data from internal service",
				Err:        err,
			}
		}

		// Fetch company data from external service
		getCompanyResponse, companyErr := generalserviceapiutils.GetCompanyDataById(workOrderReq.CompanyId)
		if companyErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: companyErr.StatusCode,
				Message:    "Failed to fetch company data from external service",
				Err:        companyErr.Err,
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
			VehicleCode:                  getVehicleResponse.Data.Master.VehicleChassisNumber,
			VehicleTnkb:                  getVehicleResponse.Data.STNK.VehicleRegistrationCertificateTNKB,
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

	pages.Rows = mapResponses

	return pages, nil
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
	brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
	if brandErr != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve brand data from the external API",
			Err:        brandErr.Err,
		}
	}

	// Fetch data model from external API
	modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
	if modelErr != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data from the external API",
			Err:        modelErr.Err,
		}
	}

	// Fetch data variant from external API
	variantResponse, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
	if variantErr != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data from the external API",
			Err:        variantErr.Err,
		}
	}

	// Fetch data vehicle from external API
	vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
	if vehicleErr != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        vehicleErr.Err,
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
	getWorkOrderStatusResponses, workOrderStatusErr := generalserviceapiutils.GetWorkOrderStatusById(entity.WorkOrderStatusId)
	if workOrderStatusErr != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order status data from external service",
			Err:        workOrderStatusErr.Err,
		}
	}

	// fetch data type work order
	getWorkOrderTypeResponses, workOrderTypeErr := generalserviceapiutils.GetWorkOrderTypeByID(entity.WorkOrderTypeId)
	if workOrderTypeErr != nil {
		return transactionworkshoppayloads.WorkOrderAffiliateResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch work order type data from external service",
			Err:        workOrderTypeErr.Err,
		}
	}

	payload := transactionworkshoppayloads.WorkOrderAffiliateResponse{
		WorkOrderSystemNumber:         entity.WorkOrderSystemNumber,
		WorkOrderDate:                 entity.WorkOrderDate.Format("2006-01-02"),
		WorkOrderDocumentNumber:       entity.WorkOrderDocumentNumber,
		WorkOrderTypeId:               entity.WorkOrderTypeId,
		WorkOrderTypeName:             getWorkOrderTypeResponses.WorkOrderTypeName,
		WorkOrderStatusId:             entity.WorkOrderStatusId,
		WorkOrderStatusName:           getWorkOrderStatusResponses.WorkOrderStatusName,
		ServiceAdvisorId:              entity.ServiceAdvisorId,
		ServiceSite:                   entity.ServiceSite,
		BrandId:                       entity.BrandId,
		BrandName:                     brandResponse.BrandName,
		ModelId:                       entity.ModelId,
		ModelName:                     modelResponse.ModelName,
		VariantId:                     entity.VariantId,
		VariantDescription:            variantResponse.VariantDescription,
		VehicleId:                     entity.VehicleId,
		VehicleCode:                   vehicleResponses.Data.Master.VehicleChassisNumber,
		VehicleTnkb:                   vehicleResponses.Data.STNK.VehicleRegistrationCertificateTNKB,
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

func (r *WorkOrderRepositoryImpl) DeleteVehicleServiceMultiId(tx *gorm.DB, Id int, DetailIds []int) (bool, *exceptions.BaseErrorResponse) {
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

func (r *WorkOrderRepositoryImpl) DeleteDetailWorkOrderMultiId(tx *gorm.DB, Id int, DetailIds []int) (bool, *exceptions.BaseErrorResponse) {
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

func (r *WorkOrderRepositoryImpl) DeleteRequestMultiId(tx *gorm.DB, Id int, DetailIds []int) (bool, *exceptions.BaseErrorResponse) {
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

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// usp_comLookUp
// IF @strEntity = 'CustomerByTypeAndAddress'--CUSTOMER MASTER
// uspg_wtWorkOrder0_Update
// IF @Option = 8
// --USE FOR : * WORK ORDER CHANGE BILL TO
func (r *WorkOrderRepositoryImpl) ChangeBillTo(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.ChangeBillToRequest) (transactionworkshoppayloads.ChangeBillToResponse, *exceptions.BaseErrorResponse) {
	var existingWorkOrder struct {
		WorkOrderOperationItemLine int
	}

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

	var entity transactionworkshopentities.WorkOrder
	err = tx.Model(&transactionworkshopentities.WorkOrder{}).Where("work_order_system_number = ?", workOrderId).First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.ChangeBillToResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Failed to retrieve work order from the database",
			Err:        err,
		}
	}

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

	return transactionworkshoppayloads.ChangeBillToResponse{
		WorkOrderSystemNumber: workOrderId,
		BillToCustomerId:      entity.CustomerId,
		BillableToId:          entity.BillableToId,
	}, nil
}

// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// uspg_wtWorkOrder0_Update
// IF @Option = 13
//
//	--USE FOR : * WORK ORDER CHANGE PHONE NO
func (r *WorkOrderRepositoryImpl) ChangePhoneNo(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.ChangePhoneNoRequest) (*transactionworkshoppayloads.ChangePhoneNoResponse, *exceptions.BaseErrorResponse) {
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

// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// uspg_wtWorkOrder2_Update
// IF @Option = 14
// --USE FOR : CONFIRM PRICE
func (r *WorkOrderRepositoryImpl) ConfirmPrice(tx *gorm.DB, workOrderId int, idwos []int, request transactionworkshoppayloads.WorkOrderConfirmPriceRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	var response transactionworkshopentities.WorkOrderDetail
	var markupPercentage, totalPackage, totalOpr, totalPart, totalOil, totalMaterial, totalConsumableMaterial, totalSublet, totalAccs, totalSouvenir float64
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
	vehicleResponses, vehicleErr := salesserviceapiutils.GetVehicleById(entity.VehicleId)
	if vehicleErr != nil {
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve vehicle data from the external API",
			Err:        vehicleErr.Err,
		}
	}

	vehicleChassisNo = vehicleResponses.Data.Master.VehicleChassisNumber

	// Check if the vehicle is in the grey market by looking up the vehicle chassis number
	err = tx.Table("dms_microservices_sales_dev.dbo.mtr_vehicle").
		Select("ISNULL(vehicle_is_grey_market, 0)").
		Where("vehicle_chassis_number = ?", vehicleChassisNo).
		Scan(&greyMarket).Error

	if err != nil {
		return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check vehicle grey market status",
			Err:        err,
		}
	}

	// Check for grey market markup
	if greyMarketMarkupPercentageExists && greyMarket {
		if invSysNo == 0 {
			markupPercentage = 40
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

		err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
			Select("COALESCE(SUM(ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0, 0)), 0)").
			Where("work_order_system_number = ? AND line_type_id = ?", workOrderId, utils.LinetypeSouvenir).
			Scan(&totalSouvenir).Error

		if err != nil {
			return transactionworkshopentities.WorkOrderDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to calculate total souvenir",
				Err:        err,
			}
		}

		// Calculate total and discounts
		total = totalPackage + totalOpr + totalPart + totalOil + totalMaterial + totalSublet + totalAccs + totalConsumableMaterial + totalSouvenir

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

// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// uspg_wtWorkOrder2_Update
// IF @Option = 6
// --USE FOR : CHECK DETAIL
func (r *WorkOrderRepositoryImpl) CheckDetail(tx *gorm.DB, workOrderId int, idwos []int) (bool, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	var detailentity transactionworkshopentities.WorkOrderDetail

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

	// Fetch wo transaction type campaign
	wotransactionType, fetchErr := generalserviceapiutils.GetWoTransactionTypeByCode("G")
	if fetchErr != nil {
		return false, fetchErr
	}

	// Check transaction type bukan campaign
	// logic Validasi Jenis Transaksi ,Cek Ketersediaan Stok, Proses Substitusi,Penghitungan Total dan Diskon, Update Work Order
	if detailentity.TransactionTypeId != wotransactionType.WoTransactionTypeId {
		if detailentity.InvoiceSystemNumber == 0 {

			// fetch linetype operation or package
			linetypeOperation, LinetypeErr := generalserviceapiutils.GetLineTypeByCode(utils.LinetypeOperation)
			if LinetypeErr != nil {
				return false, LinetypeErr
			}

			linetypePackage, LinetypeErr := generalserviceapiutils.GetLineTypeByCode(utils.LinetypePackage)
			if LinetypeErr != nil {
				return false, LinetypeErr
			}

			// Check if the line type is an operation or package
			if detailentity.LineTypeId == linetypeOperation.LineTypeId || detailentity.LineTypeId == linetypePackage.LineTypeId {
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

					qtyAvail, errResponse := r.lookupRepo.SelectLocationStockItem(tx, 1, entity.CompanyId, time.Now(), 0, "", detailentity.OperationItemId, detailentity.WarehouseGroupId, "S")
					if errResponse != nil {
						return false, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to check quantity available",
							Err:        err,
						}
					}

					if qtyAvail == 0 {
						// Step 1: Define substitute table structure
						type Substitute struct {
							SubsItemCode string  `gorm:"column:SUBS_ITEM_CODE"`
							ItemName     string  `gorm:"column:ITEM_NAME"`
							SupplyQty    float64 `gorm:"column:SUPPLY_QTY"`
							SubsType     string  `gorm:"column:SUBS_TYPE"`
						}

						// Step 2: Insert data into the temporary table using the stored procedure
						// err = tx.Exec(`
						// 	INSERT INTO #SUBS
						// 	EXEC dbo.uspg_smSubstitute0_Select @Option = 1, @Company_Code = ?, @Item_Code = ?, @Whs_Group = ?, @UoM_Type = ?
						// `, entity.CompanyId, detailentity.OperationItemCode, detailentity.WarehouseGroupId, "S").Error
						// if err != nil {
						// 	return false, &exceptions.BaseErrorResponse{
						// 		StatusCode: http.StatusInternalServerError,
						// 		Message:    "Failed to insert data into temporary table",
						// 		Err:        err,
						// 	}
						// }

						// Step 2: Create a temporary substitute table
						if err := tx.Migrator().CreateTable(&Substitute{}); err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to create substitute table",
								Err:        err,
							}
						}

						// Step 3: Insert substitute data
						var substituteItems []Substitute

						err := tx.Table("mtr_item_substitute_detail").
							Select("mtr_item_substitute_detail.item_id AS subs_item_code, mtr_item_substitute.description AS item_name, mtr_item_substitute_detail.quantity AS supply_qty, mtr_item_substitute.substitute_type_id AS subs_type").
							Joins("INNER JOIN mtr_item_substitute ON mtr_item_substitute_detail.item_substitute_id = mtr_item_substitute.item_substitute_id").
							Where("mtr_item_substitute_detail.is_active = ? AND mtr_item_substitute.item_id = ?", true, detailentity.OperationItemId).
							Scan(&substituteItems).Error

						if err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to fetch substitute items with JOIN",
								Err:        err,
							}
						}

						// Log the fetched data (debug)
						// for _, s := range substituteItems {
						// 	fmt.Printf("Substitute Item: %+v\n", s)
						// }

						// Create temporary table
						if err := tx.Migrator().CreateTable(&Substitute{}); err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to create temporary substitute table",
								Err:        err,
							}
						}

						// Insert mapped data into the temporary substitute table
						if err := tx.Create(&substituteItems).Error; err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to insert data into the substitute table",
								Err:        err,
							}
						}

						// Step 4: Check if the original item exists in the substitution table
						var exists bool
						if err := tx.Model(&Substitute{}).
							Select("COUNT(1) > 0").
							Where("SUBS_ITEM_CODE = ?", detailentity.OperationItemCode).
							Find(&exists).Error; err != nil {
							return false, &exceptions.BaseErrorResponse{
								StatusCode: http.StatusInternalServerError,
								Message:    "Failed to check if the original item is in the substitution table",
								Err:        err,
							}
						}

						if exists {
							return true, nil
						}

						if !exists {
							var substitutes []Substitute
							if err := tx.Find(&substitutes).Error; err != nil {
								return false, &exceptions.BaseErrorResponse{
									StatusCode: http.StatusInternalServerError,
									Message:    "Failed to fetch substitute items",
									Err:        err,
								}
							}

							// Step 5: Process each substitute item
							for _, substituteItem := range substitutes {

								// Step 6: Check and update the original item if not substituted
								if substituteItem.SubsType != "" {
									if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
										Where("work_order_system_number = ? AND operation_item_code = ? AND work_order_operation_item_line = ? AND substitute_type_id IS NULL",
											workOrderId, detailentity.OperationItemCode, idwos).
										Updates(map[string]interface{}{
											"substitute_type_id": 1,
											"substitute_type":    "SUBSTITUTE_ITEM",
											"warehouse_group_id": detailentity.WarehouseGroupId,
										}).Error; err != nil {
										return false, &exceptions.BaseErrorResponse{
											StatusCode: http.StatusInternalServerError,
											Message:    "Failed to update original item",
											Err:        err,
										}
									}
								}

								// Step 7: Get vehicle brand and currency code
								var vehicleInfo struct {
									BrandId    string `gorm:"column:vehicle_brand"`
									CurrencyId string `gorm:"column:currency"`
								}
								if err := tx.Model(&transactionworkshopentities.WorkOrder{}).
									Select("brand_id AS vehicle_brand, currency_id AS currency").
									Where("work_order_system_number = ?", workOrderId).
									Scan(&vehicleInfo).Error; err != nil {
									return false, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Failed to get vehicle brand and currency code",
										Err:        err,
									}
								}

								// Step 8: Calculate the next line number for the work order
								var nextLineNumber int
								if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
									Select("COALESCE(MAX(work_order_operation_item_line), 0) + 1").
									Where("work_order_system_number = ?", workOrderId).
									Scan(&nextLineNumber).Error; err != nil {
									return false, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Failed to calculate next line number",
										Err:        err,
									}
								}

								// Step 9: Insert the substitute item into wtWorkOrder2
								// Check if the substitute item already exists
								var count int64
								if err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
									Where("work_order_system_number = ? AND operation_item_code = ? AND substitute_item_code = ?",
										workOrderId, detailentity.OperationItemCode, substituteItem.SubsItemCode).
									Count(&count).Error; err != nil {
									return false, &exceptions.BaseErrorResponse{
										StatusCode: http.StatusInternalServerError,
										Message:    "Failed to check existence of substitute item",
										Err:        err,
									}
								}

								// check if the substitute item does not exist
								if count == 0 {
									// Get operation item price from the lookup repository
									var oprItemPrice float64

									// fetch linetype
									// linetypechecks, LinetypeErr := generalserviceapiutils.GetLineTypeById(detailentity.LineTypeId)
									// if LinetypeErr != nil {
									// 	return false, LinetypeErr
									// }

									// Fetch Opr_Item_Price
									oprItemPrice, _ = r.lookupRepo.GetOprItemPrice(tx, detailentity.LineTypeId, entity.CompanyId, detailentity.OperationItemId, entity.BrandId, entity.ModelId, detailentity.JobTypeId, entity.VariantId, entity.CurrencyId, utils.TrxTypeWoWarranty.ID, "1")

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
								// cleanup of temporary table
								defer func() {
									_ = tx.Migrator().DropTable(&Substitute{})
								}()

							}
						}
					}
				}
			}

			//-- By default price of all items will be replaced with the price from ItemPriceList
			//-- Exclude Fee from the replacing process
			var oprItemPrice, oprItemPriceDisc, discountPercent float64
			var markupAmount, markupPercentage float64
			var warrantyClaimType string

			if detailentity.LineTypeId == 6 { //utils.LinetypeSublet
				// fetch linetype
				linetypechecks, LinetypeErr := generalserviceapiutils.GetLineTypeById(detailentity.LineTypeId)
				if LinetypeErr != nil {
					return false, LinetypeErr
				}

				// Fetch Opr_Item_Price
				oprItemPrice, _ = r.lookupRepo.GetOprItemPrice(tx, detailentity.LineTypeId, entity.CompanyId, detailentity.OperationItemId, entity.BrandId, entity.ModelId, detailentity.JobTypeId, entity.VariantId, entity.CurrencyId, utils.TrxTypeWoWarranty.ID, "1")

				// Set markup percentage based on company ID
				if entity.CompanyId == 139 {
					markupPercentage = 11.00
				}

				// // Apply markup amount and percentage
				oprItemPrice = oprItemPrice + markupAmount + (oprItemPrice * (markupPercentage / 100))

				// Fetch Opr_Item_Disc_Percent
				oprItemPriceDisc, _ = r.lookupRepo.GetOprItemDisc(tx, linetypechecks.LineTypeCode, 6, detailentity.OperationItemId, entity.AgreementGeneralRepairId, entity.ProfitCenterId, detailentity.FrtQuantity*detailentity.OperationItemPrice, entity.CompanyId, entity.BrandId, entity.ContractServiceSystemNumber, utils.TrxTypeWoWarranty.ID, utils.EstWoOrderTypeId)

			} else {

				// fetch linetype
				linetypechecks, LinetypeErr := generalserviceapiutils.GetLineTypeById(detailentity.LineTypeId)
				if LinetypeErr != nil {
					return false, LinetypeErr
				}

				// Fetch Opr_Item_Price
				oprItemPrice, _ = r.lookupRepo.GetOprItemPrice(tx, detailentity.LineTypeId, entity.CompanyId, detailentity.OperationItemId, entity.BrandId, entity.ModelId, detailentity.JobTypeId, entity.VariantId, entity.CurrencyId, utils.TrxTypeWoWarranty.ID, "1")

				// Set markup percentage based on company ID
				if entity.CompanyId == 139 {
					markupPercentage = 11.00
				}

				// // Apply markup amount and percentage
				oprItemPrice = oprItemPrice + markupAmount + (oprItemPrice * (markupPercentage / 100))

				// Fetch Opr_Item_Disc_Percent
				oprItemPriceDisc, _ = r.lookupRepo.GetOprItemDisc(tx, linetypechecks.LineTypeCode, 6, detailentity.OperationItemId, entity.AgreementGeneralRepairId, entity.ProfitCenterId, detailentity.FrtQuantity*detailentity.OperationItemPrice, entity.CompanyId, entity.BrandId, entity.ContractServiceSystemNumber, utils.TrxTypeWoWarranty.ID, utils.EstWoOrderTypeId)

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

			// fetch linetype
			linetypechecks, LinetypeErr := generalserviceapiutils.GetLineTypeById(detailentity.LineTypeId)
			if LinetypeErr != nil {
				return false, LinetypeErr
			}

			if linetypechecks.LineTypeCode != utils.LinetypeOperation && linetypechecks.LineTypeCode != utils.LinetypePackage {
				warrantyClaimType = entity.ATPMWCFDocNo
			} else {
				warrantyClaimType = ""
			}

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

			GetApprovalStatus, GetApprovalStatusErr := generalserviceapiutils.GetApprovalStatusByCode("20")
			if GetApprovalStatusErr != nil {
				return false, GetApprovalStatusErr
			}

			// Calculate discounts and VAT
			err = tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
				Select(`
					SUM(CASE
						WHEN line_type_id = ? THEN
							CASE WHEN approval_id = ? AND ISNULL(operation_item_discount_request_amount, 0) > 0
							THEN ISNULL(operation_item_discount_request_amount, 0)
							ELSE ISNULL(operation_item_discount_amount, 0)
							END
						ELSE
							CASE WHEN approval_id = ? AND ISNULL(operation_item_discount_request_amount, 0) > 0
							THEN ISNULL(operation_item_discount_request_amount, 0)
							ELSE ISNULL(operation_item_discount_amount, 0)
							END

							*

							CASE WHEN LINE_TYPE <> ? THEN ISNULL(frt_quantity, 0)
							ELSE CASE WHEN ISNULL(supply_quantity, 0) > 0
							THEN ISNULL(supply_quantity, 0) ELSE ISNULL(frt_quantity, 0) END
							END
					END)`, utils.LinetypePackage, GetApprovalStatus.ApprovalStatusId, GetApprovalStatus.ApprovalStatusId, utils.LinetypeOperation).
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

// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
func checkCampaignExistence(tx *gorm.DB, workOrderId int) (int, error) {
	var campaignId int
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).
		Select("ISNULL(campaign_id, 0) AS campaign_id").
		Where("work_order_system_number = ?", workOrderId).
		Scan(&campaignId).Error
	return campaignId, err
}

func checkOperationAllocation(tx *gorm.DB, workOrderId int, woTrxTypeId int) (bool, error) {
	var exists bool
	// Check if operation is already allocated
	err := tx.Model(&transactionworkshopentities.WorkOrderDetail{}).
		Select("1").
		Where("work_order_system_number = ? AND (line_type_id = ? OR line_type_id = ?) AND transaction_type_id = ? AND ISNULL(service_status_id, '') <> ''",
			workOrderId, utils.LinetypeOperation, utils.LinetypePackage, woTrxTypeId).
		Scan(&exists).Error

	if err != nil {
		return false, fmt.Errorf("failed to check operation allocation: %w", err)
	}

	return exists, nil
}

// uspg_wtWorkOrder0_Update
// IF @Option = 7
func (r *WorkOrderRepositoryImpl) DeleteCampaign(tx *gorm.DB, workOrderId int) (transactionworkshoppayloads.DeleteCampaignPayload, *exceptions.BaseErrorResponse) {
	campaignId, err := checkCampaignExistence(tx, workOrderId)
	if err != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check campaign data from work order, campaign is empty",
			Err:        err,
		}
	}

	woTrxType, wotrxtypeErr := generalserviceapiutils.GetWoTransactionTypeByCode("Campaign")
	if wotrxtypeErr != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: wotrxtypeErr.StatusCode,
			Message:    "Failed to get transaction type",
			Err:        wotrxtypeErr.Err,
		}
	}

	exists, err := checkOperationAllocation(tx, workOrderId, woTrxType.WoTransactionTypeId)
	if err != nil {
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
		//fmt.Println("Invalid supply quantity:", supplyQuantity) // Debug
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

	// fetch subtitute type item
	getsubtitutetype, subsErr := generalserviceapiutils.GetSubstituteTypeByCode("S")
	if subsErr != nil {
		return transactionworkshoppayloads.DeleteCampaignPayload{}, &exceptions.BaseErrorResponse{
			StatusCode: subsErr.StatusCode,
			Message:    "Failed to get substitute type",
			Err:        subsErr.Err,
		}
	}

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
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeSparepart, getsubtitutetype.SubstituteTypeId).
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
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeOil, getsubtitutetype.SubstituteTypeId).
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
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeMaterial, getsubtitutetype.SubstituteTypeId).
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
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeConsumableMaterial, getsubtitutetype.SubstituteTypeId).
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
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeSublet, getsubtitutetype.SubstituteTypeId).
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
		Where("work_order_system_number = ? AND line_type_id = ? AND subtitute_type_id <> ?", workOrderId, utils.LinetypeAccesories, getsubtitutetype.SubstituteTypeId).
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

// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// uspg_wtWorkOrder2_Insert
// IF @Option = 1
// --USE FOR : * INSERT NEW DATA FROM PACKAGE IN CONTRACT SERVICE
func (r *WorkOrderRepositoryImpl) AddContractService(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.WorkOrderContractServiceRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {

	type ContractServiceData struct {
		ContractServSysNo float64 `gorm:"column:contract_service_system_number"`
		AddDiscStat       int     `gorm:"column:additional_discount_status_approval_id"`
		WhsGroupId        int     `gorm:"column:warehouse_group_id"`
		TaxFree           int     `gorm:"column:tax_free"`
	}

	var woentities transactionworkshopentities.WorkOrder
	var response transactionworkshopentities.WorkOrderDetail
	var contractServiceData ContractServiceData

	var (
		csrDescription, pphTaxCode                                                                     string
		csrFrtQty, csrPrice, csrDiscPercent, addDiscReqAmount, newFrtQty, supplyQty, oprItemDiscAmount float64
		csrOprItemCode, wcfTypeMoney, csrLineType, woOprItemLine, atpmWcfType, addDiscStat, itemTypeId int
	)

	// Set default WCF type Money
	atpmWcfTypeName, wcfErr := generalserviceapiutils.GetWarrantyClaimTypeByCode("PM")
	if wcfErr != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: wcfErr.StatusCode,
			Message:    "Failed to get ATPM WCF Type",
			Err:        wcfErr.Err,
		}
	}
	wcfTypeMoney = atpmWcfTypeName.WarrantyClaimTypeId

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

	//fmt.Println("Contract Service Data: ", contractServiceData)

	// Initialize new freight quantity
	newFrtQty = 0

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

	//fmt.Println("Contract Service Items: ", contractServiceItems)

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

		// If Atpm_Wcf_Type is empty, set it to wcfTypeMoney
		if atpmWcfType == 0 {
			atpmWcfType = wcfTypeMoney
		}

		switch csrLineType {
		case 1: //utils.LinetypePackage
			csrFrtQty = 1
			supplyQty = 1
			atpmWcfType = 0
		case 2: //utils.LinetypeOperation
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
			supplyQty = csrFrtQty
			atpmWcfType = 0
		default:
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
			WarehouseGroupId:             contractServiceData.WhsGroupId,
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
		// utils.LinetypeOperation
		// utils.LinetypePackage
		if csrLineType == 2 || csrLineType == 1 {
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
			"total_accessories":         totalAccs,
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

// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// uspg_wtWorkOrder2_Insert
// IF @Option = 2
// --USE FOR : * INSERT NEW DATA FROM PACKAGE MASTER
func (r *WorkOrderRepositoryImpl) AddGeneralRepairPackage(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.WorkOrderGeneralRepairPackageRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
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
	const profitCenterGR = "00002"

	// Step 1: Fetch data from work order and related tables
	if err := tx.Table("trx_work_order AS wo").
		Select("wo.company_id, wo.work_order_document_number,trx_work_order_detail.line_type_id, trx_work_order_detail.job_type_id, wo.agreement_general_repair_id, trx_work_order_detail.transaction_type_id, wo.brand_id, wo.campaign_id, wo.variant_id").
		Joins("INNER JOIN dms_microservices_aftersales_dev.dbo.trx_work_order_detail ON wo.work_order_system_number = trx_work_order_detail.work_order_system_number").
		Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_customer ON wo.customer_id = mtr_customer.customer_id").
		Where("wo.work_order_system_number = ?", workOrderId).
		Scan(&result).Error; err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch data from work order",
			Err:        err,
		}
	}

	// fetch wo job type from external api
	woJobTypePM, err := generalserviceapiutils.GetJobTransactionTypeByCode("PM")
	if err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: err.StatusCode,
			Message:    "Failed to fetch job type PM",
			Err:        err.Err,
		}
	}

	woJobTypeTG, err := generalserviceapiutils.GetJobTransactionTypeByCode("TG")
	if err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: err.StatusCode,
			Message:    "Failed to fetch job type TG",
			Err:        err.Err,
		}
	}

	// Step 2: Determine job type based on Profit Center GR
	if request.CPCCode == profitCenterGR {
		result.JobTypeId = woJobTypePM.JobTypeId // "PM" - Job Type for Periodical Maintenance
	} else {
		result.JobTypeId = woJobTypeTG.JobTypeId // "TG" - Job Type for Transfer to General Repair
	}

	// Step 3: Fetch Whs_Group based on company code
	whsGroupValue, err := r.lookupRepo.GetWhsGroup(tx, result.CompanyCode)
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
		LineTypeCode int
		OprItemId    int
		OprItemCode  string
		FrtQty       float64
		JobTypeId    int
		TrxTypeId    int
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

	woJobTypePDI, err := generalserviceapiutils.GetJobTransactionTypeByCode("PDI")
	if err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: err.StatusCode,
			Message:    "Failed to fetch job type PDI",
			Err:        err.Err,
		}
	}

	woJobTypeFSI, err := generalserviceapiutils.GetJobTransactionTypeByCode("FSI")
	if err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: err.StatusCode,
			Message:    "Failed to fetch job type FSI",
			Err:        err.Err,
		}
	}

	woJobTypeWR, err := generalserviceapiutils.GetJobTransactionTypeByCode("WR")
	if err != nil {
		return entity, &exceptions.BaseErrorResponse{
			StatusCode: err.StatusCode,
			Message:    "Failed to fetch job type FSI",
			Err:        err.Err,
		}
	}

	// Step 6: Process each package and apply business logic
	for _, pkg := range packages {
		csrLineTypeId := pkg.LineTypeCode
		csrOprItemId := pkg.OprItemId
		csrOprItemCode := pkg.OprItemCode
		csrFrtQty := pkg.FrtQty
		csrJobTypeId := pkg.JobTypeId
		csrTrxTypeId := pkg.TrxTypeId

		// Update FRT_QTY if LineType is a package
		if csrLineTypeId == 1 { //utils.LinetypePackage
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
		if csrJobTypeId == woJobTypePDI.JobTypeId || csrJobTypeId == woJobTypeFSI.JobTypeId || csrJobTypeId == woJobTypeWR.JobTypeId {
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

		// fetch line type
		// linetypecheck, linetypeErr := generalserviceapiutils.GetLineTypeById(csrLineTypeId)
		// if linetypeErr != nil {
		// 	return entity, &exceptions.BaseErrorResponse{
		// 		StatusCode: linetypeErr.StatusCode,
		// 		Message:    "Failed to fetch line type",
		// 		Err:        linetypeErr.Err,
		// 	}
		// }

		// Fetch operation item price and discount percent
		oprItemPrice, err := r.lookupRepo.GetOprItemPrice(tx, csrLineTypeId, result.CompanyCode, result.BrandId, csrOprItemId, result.AgreementNo, result.JobTypeId, csrTrxTypeId, int(csrFrtQty), whsGroup, strconv.Itoa(result.VariantId))
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
		oprItemDiscPercent, err := r.lookupRepo.GetOprItemDisc(tx, strconv.Itoa(result.LinetypeId), result.BillCodeExt, csrOprItemId, result.AgreementNo, result.CpcCode, oprItemPrice*csrFrtQty, result.CompanyCode, result.BrandId, 0, whsGroup, 0)
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
		if csrLineTypeId == 2 || csrLineTypeId == 1 {
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
			getLineTypeByItemCode, err := r.lookupRepo.GetLineTypeByItemCode(tx, csrOprItemCode)
			if err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch line type by item code",
					Err:        errors.New("failed to fetch line type by item code"),
				}
			}

			// fetch line type
			linetypecheck, linetypeErr := generalserviceapiutils.GetLineTypeByCode(getLineTypeByItemCode)
			if linetypeErr != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: linetypeErr.StatusCode,
					Message:    "Failed to fetch line type",
					Err:        linetypeErr.Err,
				}
			}

			// Check if the line type matches the one retrieved by item code
			if csrLineTypeId != linetypecheck.LineTypeId {
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
			qtyAvailable, err := r.lookupRepo.SelectLocationStockItem(tx, 1, result.CompanyCode, currentDate, 0, "", csrOprItemId, whsGroup, uomType)
			if err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to execute stock item location query",
					Err:        errors.New("failed to execute stock item location query"),
				}
			}

			// Check if the item is available or the line type is Sublet
			if qtyAvailable > 0 || csrLineTypeId == 6 { //utils.LinetypeSublet
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

					// fetch linetype
					// linetypecheck, linetypeErr := generalserviceapiutils.GetLineTypeById(csrLineTypeId)
					// if linetypeErr != nil {
					// 	return entity, &exceptions.BaseErrorResponse{
					// 		StatusCode: linetypeErr.StatusCode,
					// 		Message:    "Failed to fetch line type",
					// 		Err:        linetypeErr.Err,
					// 	}
					// }

					// Fetch operation item price and discount percent
					oprItemPrice, err := r.lookupRepo.GetOprItemPrice(tx, csrLineTypeId, result.CompanyCode, result.BrandId, csrOprItemId, result.AgreementNo, result.JobTypeId, csrTrxTypeId, int(csrFrtQty), newWhsGroup, strconv.Itoa(result.VariantId))
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
					oprItemDiscPercent, err := r.lookupRepo.GetOprItemDisc(tx, strconv.Itoa(result.LinetypeId), result.BillCodeExt, csrOprItemId, result.AgreementNo, result.CpcCode, oprItemPrice*newFrtQty, result.CompanyCode, result.BrandId, 0, whsGroup, 0)
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
					SubstituteItemCode string  `gorm:"column:SUBS_ITEM_CODE"`
					ItemName           string  `gorm:"column:ITEM_NAME"`
					SupplyQty          float64 `gorm:"column:SUPPLY_QTY"`
					SubstituteType     string  `gorm:"column:SUBS_TYPE"`
				}

				// Step 2: Create a temporary substitute table
				if err := tx.Migrator().CreateTable(&Substitute{}); err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to create substitute table",
						Err:        err,
					}
				}

				// Step 2: Execute stored procedure to populate the temporary table
				var substitutes []Substitute

				err := tx.Table("mtr_item_substitute_detail").
					Select("mtr_item_substitute_detail.item_id AS subs_item_code, mtr_item_substitute.description AS item_name, mtr_item_substitute_detail.quantity AS supply_qty, mtr_item_substitute.substitute_type_id AS subs_type").
					Joins("INNER JOIN mtr_item_substitute ON mtr_item_substitute_detail.item_substitute_id = mtr_item_substitute.item_substitute_id").
					Where("mtr_item_substitute_detail.is_active = ? AND mtr_item_substitute.item_id = ?", true, entity.OperationItemId).
					Scan(&substitutes).Error

				if err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to fetch substitute items with JOIN",
						Err:        err,
					}
				}

				// Log the fetched data (debug)
				// for _, s := range substituteItems {
				// 	fmt.Printf("Substitute Item: %+v\n", s)
				// }

				// Step 3: Fetch data from the temporary table

				// Create temporary table
				if err := tx.Migrator().CreateTable(&Substitute{}); err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to create temporary substitute table",
						Err:        err,
					}
				}

				// Insert mapped data into the temporary substitute table
				if err := tx.Create(&substitutes).Error; err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to insert data into the substitute table",
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

						// fetch linetype
						// linetypecheck, linetypeErr := generalserviceapiutils.GetLineTypeById(csrLineTypeId)
						// if linetypeErr != nil {
						// 	return entity, &exceptions.BaseErrorResponse{
						// 		StatusCode: linetypeErr.StatusCode,
						// 		Message:    "Failed to fetch line type",
						// 		Err:        linetypeErr.Err,
						// 	}
						// }

						// Fetch operation item price and discount percent
						oprItemPrice, err := r.lookupRepo.GetOprItemPrice(tx, csrLineTypeId, result.CompanyCode, result.BrandId, csrOprItemId, result.AgreementNo, result.JobTypeId, csrTrxTypeId, int(csrFrtQty), whsGroup, strconv.Itoa(result.VariantId))
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
						oprItemDiscPercent, err := r.lookupRepo.GetOprItemDisc(tx, strconv.Itoa(result.LinetypeId), result.BillCodeExt, csrOprItemId, result.AgreementNo, result.CpcCode, oprItemPrice*substitute.SupplyQty, result.CompanyCode, result.BrandId, 0, whsGroup, 0)
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

				// cleanup of temporary table
				defer func() {
					_ = tx.Migrator().DropTable(&Substitute{})
				}()
			}
		}
	}

	return entity, nil
}

// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////
// uspg_wtWorkOrder2_Insert
// IF @Option = 3
func (r *WorkOrderRepositoryImpl) AddFieldAction(tx *gorm.DB, workOrderId int, request transactionworkshoppayloads.WorkOrderFieldActionRequest) (transactionworkshopentities.WorkOrderDetail, *exceptions.BaseErrorResponse) {
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
			WhsGroup, err := r.lookupRepo.GetWhsGroup(tx, companyCode)
			if err != nil {
				return entity, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Error fetching Warehouse Group",
					Err:        errors.New("error fetching Warehouse Group"),
				}
			}

			var agreementNo int
			if recallRecord.LineTypeId == 2 { //utils.LinetypeOperation
				agreementNo = agreementNoBR
			} else {
				agreementNo = agreementNoGR
			}

			// fetch linetype
			// linetypecheck, linetypeErr := generalserviceapiutils.GetLineTypeById(recallRecord.LineTypeId)
			// if linetypeErr != nil {
			// 	return entity, linetypeErr
			// }

			oprItemPrice, err := r.lookupRepo.GetOprItemPrice(
				tx, recallRecord.LineTypeId, companyCode, vehicleBrand, recallRecord.OprItemId, agreementNo,
				entity.JobTypeId, utils.TrxTypeWoWarranty.ID,
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
			case 2: //utils.LinetypeOperation
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

				frtQty, err := r.lookupRepo.GetOprItemFrt(tx, recallRecord.OprItemId, vehicleBrand, modelCode, variantCode, vehicleChassisNo)
				if err != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Error fetching Frt_Qty and Supply_Qty",
						Err:        errors.New("error fetching Frt_Qty and Supply_Qty"),
					}
				}
				entity.FrtQuantity = frtQty

			case 1: //utils.LinetypePackage
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

				lineTypeStr, err := r.lookupRepo.GetLineTypeByItemCode(tx, recallRecord.OprItemCode)
				if err != nil {
					return entity, err
				}

				linetypeCheck, linetypeErr := generalserviceapiutils.GetLineTypeByCode(lineTypeStr)
				if linetypeErr != nil {
					return entity, linetypeErr
				}

				if linetypeCheck.LineTypeId != recallRecord.LineTypeId {
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

			if recallRecord.LineTypeId == 1 { //utils.LinetypePackage
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
			//utils.LinetypeOperation
			//utils.LinetypePackage
			if recallRecord.LineTypeId == 2 || recallRecord.LineTypeId == 1 {

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
						if recallRecord.LineTypeId == 1 { //utils.LinetypePackage
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
				linetypecode, err := r.lookupRepo.GetLineTypeByItemCode(tx, recallRecord.OprItemCode)
				if err != nil {
					return entity, err
				}

				// fetch line type from external service
				linetypeCheck, LinetypeErr := generalserviceapiutils.GetLineTypeByCode(linetypecode)
				if LinetypeErr != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to get line type by code",
						Err:        LinetypeErr,
					}
				}

				if recallRecord.LineTypeId != linetypeCheck.LineTypeId {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusBadRequest,
						Message:    "Item Code belong to Line Type",
						Err:        errors.New("item code belong to line type"),
					}
				}

				UomType := utils.UomTypeService

				var qtyAvail float64
				qtyAvail, errResponse := r.lookupRepo.SelectLocationStockItem(tx, 1, companyCode, time.Now(), 0, "", recallRecord.OprItemId, WhsGroup, UomType)
				if errResponse != nil {
					return entity, &exceptions.BaseErrorResponse{
						StatusCode: http.StatusInternalServerError,
						Message:    "Failed to check quantity available",
						Err:        errResponse.Err,
					}
				}

				// Assuming necessary imports and definitions are in place
				if qtyAvail > 0 || recallRecord.LineTypeId == 6 { //utils.LinetypeSublet
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
								LineTypeId:                          6, //utils.LinetypeSublet
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
								lineTypeId string
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
						SubstituteItemCode string  `gorm:"column:SUBS_ITEM_CODE"`
						ItemName           string  `gorm:"column:ITEM_NAME"`
						SupplyQty          float64 `gorm:"column:SUPPLY_QTY"`
						SubstituteType     string  `gorm:"column:SUBS_TYPE"`
					}

					// Step 1: Create a temporary table with GORM
					if err := tx.Migrator().CreateTable(&Substitute{}); err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to create substitute table",
							Err:        err,
						}
					}

					// Step 2: Execute stored procedure to populate the temporary table
					var substitutes []Substitute
					err := tx.Table("mtr_item_substitute_detail").
						Select("mtr_item_substitute_detail.item_id AS subs_item_code, mtr_item_substitute.description AS item_name, mtr_item_substitute_detail.quantity AS supply_qty, mtr_item_substitute.substitute_type_id AS subs_type").
						Joins("INNER JOIN mtr_item_substitute ON mtr_item_substitute_detail.item_substitute_id = mtr_item_substitute.item_substitute_id").
						Where("mtr_item_substitute_detail.is_active = ? AND mtr_item_substitute.item_id = ?", true, entity.OperationItemId).
						Scan(&substitutes).Error

					if err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to fetch substitute items with JOIN",
							Err:        err,
						}
					}

					// Log the fetched data (debug)
					// for _, s := range substituteItems {
					// 	fmt.Printf("Substitute Item: %+v\n", s)
					// }

					// Step 3: Fetch data from the temporary table
					if err := tx.Migrator().CreateTable(&Substitute{}); err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to create temporary substitute table",
							Err:        err,
						}
					}

					// Insert mapped data into the temporary substitute table
					if err := tx.Create(&substitutes).Error; err != nil {
						return entity, &exceptions.BaseErrorResponse{
							StatusCode: http.StatusInternalServerError,
							Message:    "Failed to insert data into the substitute table",
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
									LineTypeId:                   6, //utils.LinetypeSublet
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

							// Fetch operation item price and discount percent  utils.LinetypeSublet
							oprItemPrice, err := r.lookupRepo.GetOprItemPrice(tx, 9, result.CompanyId, result.BrandId, recallRecord.OprItemId, agreementNo, utils.TrxTypeWoFreeService.ID, utils.TrxTypeWoFreeService.ID, int(recallRecord.FrtQty), WhsGroup, strconv.Itoa(result.VariantId))
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
							oprItemDiscPercent, err := r.lookupRepo.GetOprItemDisc(tx, utils.LinetypeSublet, utils.TrxTypeWoFreeService.ID, recallRecord.OprItemId, agreementNo, 00002, oprItemPrice*substitute.SupplyQty, result.CompanyId, result.BrandId, 0, WhsGroup, 0)
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
								LineTypeId:                   6, //utils.LinetypeSublet
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

					// cleanup of temporary table
					defer func() {
						_ = tx.Migrator().DropTable(&Substitute{})
					}()
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

func (r *WorkOrderRepositoryImpl) GetOperationItemById(LineTypeId int, OperationItemId int) (interface{}, *exceptions.BaseErrorResponse) {
	// Validate LineType
	if LineTypeId < 0 || LineTypeId > 9 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid LineType",
			Err:        fmt.Errorf("invalid LineType: %d", LineTypeId),
		}
	}

	// URL for the request
	url := fmt.Sprintf("%slookup/item-opr-code/%d/by-id/%d", config.EnvConfigs.AfterSalesServiceUrl, LineTypeId, OperationItemId)
	log.Printf("Requesting URL: %s", url)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create request",
			Err:        err,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to make request",
			Err:        err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    "Failed to get operation item",
			Err:        errors.New("failed to get operation item"),
		}
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to read response body",
			Err:        err,
		}
	}

	//log.Printf("Raw Response: %s", string(body))

	var apiResponse transactionworkshoppayloads.ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to decode API response",
			Err:        err,
		}
	}

	if apiResponse.StatusCode != http.StatusOK {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: apiResponse.StatusCode,
			Message:    apiResponse.Message,
			Err:        errors.New(apiResponse.Message),
		}
	}

	// Handle data according to LineType
	var responseData interface{}
	switch LineTypeId {
	case 2, 3, 4, 5, 6, 7, 8, 9:
		var response transactionworkshoppayloads.LineType2To9Response
		if err := r.mapToStruct(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error unmarshaling 'data' into LineType2To9Response",
				Err:        err,
			}
		}
		responseData = response
	case 1:
		var response transactionworkshoppayloads.LineType1Response
		if err := r.mapToStruct(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error unmarshaling 'data' into LineType1Response",
				Err:        err,
			}
		}
		responseData = response
	case 0:
		var response transactionworkshoppayloads.LineType0Response
		if err := r.mapToStruct(apiResponse.Data, &response); err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error unmarshaling 'data' into LineType0Response",
				Err:        err,
			}
		}
		responseData = response
	default:
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Unknown line type in operation item response",
			Err:        fmt.Errorf("unexpected line type %d", LineTypeId),
		}
	}

	return responseData, nil
}

// Helper function to map data into the correct struct
func (r *WorkOrderRepositoryImpl) mapToStruct(data map[string]interface{}, result interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	if err := json.Unmarshal(dataBytes, result); err != nil {
		return fmt.Errorf("error unmarshaling into struct: %v", err)
	}
	return nil
}

// aftersalesserviceapiutils/response_utils.go
func (r *WorkOrderRepositoryImpl) HandleLineTypeResponse(lineTypeId int, operationItemResponse interface{}) (string, string, *exceptions.BaseErrorResponse) {
	var OperationItemCode, Description string

	switch lineTypeId {
	case 1:
		if response, ok := operationItemResponse.(transactionworkshoppayloads.LineType0Response); ok {
			OperationItemCode = response.PackageCode
			Description = response.PackageName
		} else {
			return "", "", &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to decode LineType0 response",
				Err:        errors.New("failed to decode LineType0 response"),
			}
		}

	case 2:
		if response, ok := operationItemResponse.(transactionworkshoppayloads.LineType1Response); ok {
			OperationItemCode = response.OperationCode
			Description = response.OperationName
		} else {
			return "", "", &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to decode LineType1 response",
				Err:        errors.New("failed to decode LineType1 response"),
			}
		}

	default:
		if response, ok := operationItemResponse.(transactionworkshoppayloads.LineType2To9Response); ok {
			OperationItemCode = response.ItemCode
			Description = response.ItemName
		} else {
			return "", "", &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to decode LineType2-9 response",
				Err:        errors.New("failed to decode LineType2-9 response"),
			}
		}
	}

	return OperationItemCode, Description, nil
}
