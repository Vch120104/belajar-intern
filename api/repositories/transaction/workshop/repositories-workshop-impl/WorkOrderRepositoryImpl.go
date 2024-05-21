package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	mastercampaignmasterentities "after-sales/api/entities/master/campaign_master"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"fmt"
	"net/http"
	"strconv"

	exceptions "after-sales/api/exceptions"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderRepositoryImpl struct {
}

func OpenWorkOrderRepositoryImpl() transactionworkshoprepository.WorkOrderRepository {
	return &WorkOrderRepositoryImpl{}
}

func (r *WorkOrderRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrder
	// Query to retrieve all work order entities based on the request
	query := tx.Model(&transactionworkshopentities.WorkOrder{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work orders from the database"}
	}

	var workOrderResponses []map[string]interface{}

	// Loop through each entity and copy its data to the response
	for _, entity := range entities {
		workOrderData := make(map[string]interface{})
		// Copy data from entity to response
		workOrderData["last_approval_by_id"] = entity.LastApprovalBy
		workOrderData["last_approval_date"] = entity.LastApprovalDate
		workOrderData["work_order_system_number"] = entity.WorkOrderSystemNumber
		workOrderData["company_id"] = entity.CompanyId
		workOrderData["work_order_document_number"] = entity.WorkOrderDocumentNumber
		workOrderData["work_order_status_id"] = entity.WorkOrderStatusId
		workOrderData["work_order_date"] = entity.WorkOrderDate
		workOrderData["work_order_close_date"] = entity.WorkOrderCloseDate
		workOrderData["work_order_type_id"] = entity.WorkOrderTypeId
		workOrderData["work_order_repeated_system_number"] = entity.WorkOrderRepeatedSystemNumber
		workOrderData["work_order_repeated_document_number"] = entity.WorkOrderRepeatedDocumentNumber
		workOrderData["afiliated_company"] = entity.AffiliatedCompany
		workOrderData["profit_center_id"] = entity.ProfitCenterId
		workOrderData["brand_id"] = entity.BrandId
		workOrderData["model_id"] = entity.ModelId
		workOrderData["variant_id"] = entity.VariantId
		workOrderData["vehicle_chassis_number"] = entity.VehicleChassisNumber
		workOrderData["billable_to_id"] = entity.BillableToId
		workOrderData["customer_id"] = entity.CustomerId
		workOrderData["pay_type"] = entity.PayType
		workOrderData["from_era"] = entity.FromEra
		workOrderData["queue_number"] = entity.QueueNumber
		workOrderData["arrival_time"] = entity.ArrivalTime
		workOrderData["service_mileage"] = entity.ServiceMileage
		workOrderData["leave_car"] = entity.LeaveCar
		workOrderData["storing"] = entity.Storing
		workOrderData["era_number"] = entity.EraNumber
		workOrderData["era_expired_date"] = entity.EraExpiredDate
		workOrderData["unregister"] = entity.Unregister
		workOrderData["contact_person_name"] = entity.ContactPersonName
		workOrderData["contact_person_phone"] = entity.ContactPersonPhone
		workOrderData["contact_person_mobile"] = entity.ContactPersonMobile
		workOrderData["contact_person_contact_via"] = entity.ContactPersonContactVia
		workOrderData["contract_service_system_number"] = entity.ContractServiceSystemNumber
		workOrderData["agrement_general_repair_id"] = entity.AgreementGeneralRepairId
		workOrderData["agreement_body_repair_id"] = entity.AgreementBodyRepairId
		workOrderData["booking_system_number"] = entity.BookingSystemNumber
		workOrderData["estimation_system_number"] = entity.EstimationSystemNumber
		workOrderData["pdi_system_number"] = entity.PDISystemNumber
		workOrderData["pdi_document_number"] = entity.PDIDocumentNumber
		workOrderData["pdi_line_number"] = entity.PDILineNumber
		workOrderData["service_request_system_number"] = entity.ServiceRequestSystemNumber
		workOrderData["campaign_id"] = entity.CampaignId
		workOrderData["campaign_code"] = entity.CampaignCode
		workOrderData["insurance_policy_number"] = entity.InsurancePolicyNumber
		workOrderData["insurance_expired_date"] = entity.InsuranceExpiredDate
		workOrderData["insurance_claim_number"] = entity.InsuranceClaimNumber
		workOrderData["insurance_person_in_charge"] = entity.InsurancePersonInCharge
		workOrderData["insurance_own_risk"] = entity.InsuranceOwnRisk
		workOrderData["insurance_work_order_number"] = entity.InsuranceWorkOrderNumber
		workOrderData["total_package"] = entity.TotalPackage
		workOrderData["total_operation"] = entity.TotalOperation
		workOrderData["total_part"] = entity.TotalPart
		workOrderData["total_oil"] = entity.TotalOil
		workOrderData["total_material"] = entity.TotalMaterial
		workOrderData["total_consumable_material"] = entity.TotalConsumableMaterial
		workOrderData["total_sublet"] = entity.TotalSublet
		workOrderData["total_price_accessories"] = entity.TotalPriceAccessories
		workOrderData["total_discount"] = entity.TotalDiscount
		workOrderData["total"] = entity.Total
		workOrderData["total_vat"] = entity.TotalVAT
		workOrderData["total_after_vat"] = entity.TotalAfterVAT
		workOrderData["total_pph"] = entity.TotalPPH
		workOrderData["discount_request_percent"] = entity.DiscountRequestPercent
		workOrderData["discount_request_amount"] = entity.DiscountRequestAmount
		workOrderData["tax_id"] = entity.TaxId
		workOrderData["vat_tax_rate"] = entity.VATTaxRate
		workOrderData["additional_discount_status_approval_id"] = entity.AdditionalDiscountStatusApprovalId
		workOrderData["remark"] = entity.Remark
		workOrderData["foreman_id"] = entity.Foreman
		workOrderData["production_head_id"] = entity.ProductionHead
		workOrderData["estimate_time"] = entity.EstTime
		workOrderData["notes"] = entity.Notes
		workOrderData["suggestion"] = entity.Suggestion
		workOrderData["fs_coupon_number"] = entity.FSCouponNo
		workOrderData["service_advisor_id"] = entity.ServiceAdvisor
		workOrderData["incentive_date"] = entity.IncentiveDate
		workOrderData["work_order_cancel_reason"] = entity.WOCancelReason
		workOrderData["invoice_system_number"] = entity.InvoiceSystemNumber
		workOrderData["currency_id"] = entity.CurrencyId
		workOrderData["ATPM_warranty_claim_form_document_number"] = entity.ATPMWCFDocNo
		workOrderData["ATPM_warranty_claim_form_date"] = entity.ATPMWCFDate
		workOrderData["ATPM_free_service_document_number"] = entity.ATPMFSDocNo
		workOrderData["ATPM_free_service_date"] = entity.ATPMFSDate
		workOrderData["total_after_discount"] = entity.TotalAfterDisc
		workOrderData["approval_request_number"] = entity.ApprovalReqNo
		workOrderData["journal_system_number"] = entity.JournalSysNo
		workOrderData["approval_gatepass_request_number"] = entity.ApprovalGatepassReqNo
		workOrderData["downpayment_amount"] = entity.DPAmount
		workOrderData["downpayment_payment"] = entity.DPPayment
		workOrderData["downpayment_payment_allocated"] = entity.DPPaymentAllocated
		workOrderData["downpayment_payment_vat"] = entity.DPPaymentVAT
		workOrderData["downpayment_payment_to_invoice"] = entity.DPAllocToInv
		workOrderData["downpayment_payment_vat_to_invoice"] = entity.DPVATAllocToInv
		workOrderData["journal_overpay_system_number"] = entity.JournalOverpaySysNo
		workOrderData["downpayment_overpay"] = entity.DPOverpay
		workOrderData["work_order_site_type_id"] = entity.SiteTypeId
		workOrderData["cost_center_id"] = entity.CostCenterId
		workOrderData["promise_date"] = entity.PromiseDate
		workOrderData["promise_time"] = entity.PromiseTime
		workOrderData["car_wash"] = entity.CarWash
		workOrderData["job_on_hold_reason"] = entity.JobOnHoldReason
		workOrderData["customer_express"] = entity.CustomerExpress
		workOrderData["contact_person_title_prefix"] = entity.CPTitlePrefix

		workOrderResponses = append(workOrderResponses, workOrderData)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(workOrderResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) VehicleLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	// Define a slice to hold WorkOrderLookupResponse responses
	var responses []transactionworkshoppayloads.WorkOrderLookupResponse

	// Define table struct
	tableStruct := transactionworkshoppayloads.WorkOrderLookupRequest{}

	// Define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// Apply filters
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Execute query
	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	// Define a slice to hold WorkOrderLookupResponse
	var convertedResponses []transactionworkshoppayloads.WorkOrderLookupResponse

	// Iterate over rows
	for rows.Next() {
		// Define variables to hold row data
		var (
			workOrderReq transactionworkshoppayloads.WorkOrderLookupRequest
			workOrderRes transactionworkshoppayloads.WorkOrderLookupResponse
		)

		// Scan the row into WorkOrderLookupRequest struct
		if err := rows.Scan(
			&workOrderReq.WorkOrderSystemNumber,
			&workOrderReq.WorkOrderDocumentNumber,
			&workOrderReq.VehicleId,
			&workOrderReq.CustomerId,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch vehicle data from external service
		VehicleURL := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(workOrderReq.VehicleId)
		fmt.Println("Fetching Vehicle data from:", VehicleURL)
		var getVehicleResponse transactionworkshoppayloads.WorkOrderVehicleResponse
		if err := utils.Get(VehicleURL, &getVehicleResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch vehicle data from external service",
				Err:        err,
			}
		}

		// Fetch Customer data from external service
		CustomerURL := config.EnvConfigs.GeneralServiceUrl + "customer-detail/" + strconv.Itoa(workOrderReq.CustomerId)
		fmt.Println("Fetching Customer data from:", CustomerURL)
		var getCustomerResponse transactionworkshoppayloads.CustomerResponse
		if err := utils.Get(CustomerURL, &getCustomerResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch customer data from external service",
				Err:        err,
			}
		}
		// Create WorkOrderLookupResponse
		workOrderRes = transactionworkshoppayloads.WorkOrderLookupResponse{
			WorkOrderDocumentNumber: workOrderRes.WorkOrderDocumentNumber,
			WorkOrderSystemNumber:   workOrderRes.WorkOrderSystemNumber,
			VehicleId:               workOrderRes.VehicleId,
			CustomerId:              workOrderRes.CustomerId,
		}
		// Append WorkOrderResponse to the slice
		convertedResponses = append(convertedResponses, workOrderRes)
	}

	// Define a slice to hold map responses
	var mapResponses []map[string]interface{}

	// Iterate over convertedResponses and convert them to maps
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_document_number": response.WorkOrderDocumentNumber,
			"work_order_system_number":   response.WorkOrderSystemNumber,
			"vehicle_id":                 response.VehicleId,
			"customer_id":                response.CustomerId,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil

}

func (r *WorkOrderRepositoryImpl) CampaignLookup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var entities []mastercampaignmasterentities.CampaignMaster
	// Query to retrieve all work order entities based on the request
	query := tx.Model(&mastercampaignmasterentities.CampaignMaster{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve campaign master from the database"}
	}

	var WorkOrderCampaignResponse []map[string]interface{}

	// Loop through each entity and copy its data to the response
	for _, entity := range entities {
		campaignData := make(map[string]interface{})
		// Copy data from entity to response
		campaignData["campaign_id"] = entity.CampaignId
		campaignData["campaign_code"] = entity.CampaignCode
		campaignData["campaign_name"] = entity.CampaignName
		campaignData["campaign_period_from"] = entity.CampaignPeriodFrom
		campaignData["campaign_period_to"] = entity.CampaignPeriodTo

		WorkOrderCampaignResponse = append(WorkOrderCampaignResponse, campaignData)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(WorkOrderCampaignResponse, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderRepositoryImpl) New(tx *gorm.DB) (transactionworkshoppayloads.WorkOrderRequest, *exceptions.BaseErrorResponse) {
	// Create a new instance of WorkOrderRequest
	var workOrderRequest transactionworkshoppayloads.WorkOrderRequest

	// Save the work order
	err := tx.Create(&workOrderRequest).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderRequest{}, &exceptions.BaseErrorResponse{
			Message: "Failed to save work order",
			Err:     err,
		}
	}

	return workOrderRequest, nil
}

func (r *WorkOrderRepositoryImpl) NewStatus(tx *gorm.DB) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptions.BaseErrorResponse) {
	var statuses []transactionworkshopentities.WorkOrderMasterStatus
	if err := tx.Find(&statuses).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order statuses from the database"}
	}
	return statuses, nil
}

func (r *WorkOrderRepositoryImpl) NewType(tx *gorm.DB) ([]transactionworkshopentities.WorkOrderMasterType, *exceptions.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterType
	if err := tx.Find(&types).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}
	return types, nil
}

func (r *WorkOrderRepositoryImpl) GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Convert entity to payload
	payload := transactionworkshoppayloads.WorkOrderRequest{}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptions.BaseErrorResponse) {
	var workOrderEntities = transactionworkshopentities.WorkOrder{
		// Assign fields from request
	}

	// Create a new record
	err := tx.Create(&workOrderEntities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) Submit(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to submit the work order
	// ...

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) Void(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to void the work order
	// ...

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) CloseOrder(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to close the work order
	// ...

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return nil
}
