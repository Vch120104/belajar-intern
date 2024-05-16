package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"net/http"

	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type WorkOrderRepositoryImpl struct {
}

func OpenWorkOrderRepositoryImpl() transactionworkshoprepository.WorkOrderRepository {
	return &WorkOrderRepositoryImpl{}
}

func (r *WorkOrderRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var entities []transactionworkshopentities.WorkOrder
	// Query to retrieve all work order entities based on the request
	query := tx.Model(&transactionworkshopentities.WorkOrder{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work orders from the database"}
	}

	var workOrderResponses []map[string]interface{}

	// Loop through each entity and copy its data to the response
	for _, entity := range entities {
		workOrderData := make(map[string]interface{})
		// Copy data from entity to response
		workOrderData["last_approval_by_id"] = entity.LastApprovalBy
		workOrderData["last_approval_date"] = entity.LastApprovalDate
		workOrderData["work_order_system_number"] = entity.WorkOrderSystemNumber
		workOrderData["company_id"] = entity.CompanyID
		workOrderData["work_order_document_number"] = entity.WorkOrderDocumentNumber
		workOrderData["work_order_status_id"] = entity.WorkOrderStatusID
		workOrderData["work_order_date"] = entity.WorkOrderDate
		workOrderData["work_order_close_date"] = entity.WorkOrderCloseDate
		workOrderData["work_order_type_id"] = entity.WorkOrderTypeID
		workOrderData["work_order_repeated_system_number"] = entity.WorkOrderRepeatedSystemNumber
		workOrderData["work_order_repeated_document_number"] = entity.WorkOrderRepeatedDocumentNumber
		workOrderData["afiliated_company"] = entity.AffiliatedCompany
		workOrderData["profit_center_id"] = entity.ProfitCenterID
		workOrderData["brand_id"] = entity.BrandID
		workOrderData["model_id"] = entity.ModelID
		workOrderData["variant_id"] = entity.VariantID
		workOrderData["vehicle_chassis_number"] = entity.VehicleChassisNumber
		workOrderData["billable_to_id"] = entity.BillableToID
		workOrderData["customer_id"] = entity.CustomerID
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
		workOrderData["agrement_general_repair_id"] = entity.AgreementGeneralRepairID
		workOrderData["agreement_body_repair_id"] = entity.AgreementBodyRepairID
		workOrderData["booking_system_number"] = entity.BookingSystemNumber
		workOrderData["estimation_system_number"] = entity.EstimationSystemNumber
		workOrderData["pdi_system_number"] = entity.PDISystemNumber
		workOrderData["pdi_document_number"] = entity.PDIDocumentNumber
		workOrderData["pdi_line_number"] = entity.PDILineNumber
		workOrderData["service_request_system_number"] = entity.ServiceRequestSystemNumber
		workOrderData["campaign_id"] = entity.CampaignID
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
		workOrderData["tax_id"] = entity.TaxID
		workOrderData["vat_tax_rate"] = entity.VATTaxRate
		workOrderData["additional_discount_status_approval_id"] = entity.AdditionalDiscountStatusApprovalID
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
		workOrderData["currency_id"] = entity.CurrencyID
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

func (r *WorkOrderRepositoryImpl) New(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	// Create a new instance of WorkOrderRepositoryImpl
	// Save the work order
	success, err := r.Save(tx, request) // Menggunakan method Save dari receiver saat ini, yaitu r
	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{Message: "Failed to save work order"}
	}

	return success, nil
}

func (r *WorkOrderRepositoryImpl) NewStatus(tx *gorm.DB) ([]transactionworkshopentities.WorkOrderMasterStatus, *exceptionsss_test.BaseErrorResponse) {
	var statuses []transactionworkshopentities.WorkOrderMasterStatus
	if err := tx.Find(&statuses).Error; err != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work order statuses from the database"}
	}
	return statuses, nil
}

func (r *WorkOrderRepositoryImpl) NewType(tx *gorm.DB) ([]transactionworkshopentities.WorkOrderMasterType, *exceptionsss_test.BaseErrorResponse) {
	var types []transactionworkshopentities.WorkOrderMasterType
	if err := tx.Find(&types).Error; err != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work order type from the database"}
	}
	return types, nil
}

func (r *WorkOrderRepositoryImpl) GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.WorkOrderRequest, *exceptionsss_test.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.WorkOrderRequest{}, &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Convert entity to payload
	payload := transactionworkshoppayloads.WorkOrderRequest{}

	return payload, nil
}

func (r *WorkOrderRepositoryImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	var workOrderEntities = transactionworkshopentities.WorkOrder{
		// Assign fields from request
	}

	// Create a new record
	err := tx.Create(&workOrderEntities).Error
	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
}

func (r *WorkOrderRepositoryImpl) Submit(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to submit the work order
	// ...

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) Void(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to void the work order
	// ...

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return nil
}

func (r *WorkOrderRepositoryImpl) CloseOrder(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse {
	// Retrieve the work order by Id
	var entity transactionworkshopentities.WorkOrder
	err := tx.Model(&transactionworkshopentities.WorkOrder{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to retrieve work order from the database"}
	}

	// Perform the necessary operations to close the work order
	// ...

	// Save the updated work order
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{Message: "Failed to save the updated work order"}
	}

	return nil
}
