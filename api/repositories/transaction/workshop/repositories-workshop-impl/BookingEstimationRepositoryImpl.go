package transactionworkshoprepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type BookingEstimationImpl struct {
	workorderRepo transactionworkshoprepository.WorkOrderRepository
}

func OpenBookingEstimationRepositoryImpl() transactionworkshoprepository.BookingEstimationRepository {
	workorderRepo := OpenWorkOrderRepositoryImpl()
	return &BookingEstimationImpl{
		workorderRepo: workorderRepo,
	}
}

// uspg_wtBookEstim0_Select
// IF @Option = 0
// --USE FOR : * SELECT DATA BY KEY
func (r *BookingEstimationImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.BookingEstimation

	baseModelQuery := tx.Model(&transactionworkshopentities.BookingEstimation{})
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve booking estimations from the database",
			Err:        err,
		}
	}

	if len(entities) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, entity := range entities {
		result := map[string]interface{}{
			"batch_system_number":               entity.BatchSystemNumber,
			"booking_system_number":             entity.BookingSystemNumber,
			"brand_id":                          entity.BrandId,
			"model_id":                          entity.ModelId,
			"variant_id":                        entity.VariantId,
			"vehicle_id":                        entity.VehicleId,
			"estimation_system_number":          entity.EstimationSystemNumber,
			"pdi_system_number":                 entity.PdiSystemNumber,
			"service_request_system_number":     entity.ServiceRequestSystemNumber,
			"contract_system_number":            entity.ContractSystemNumber,
			"agreement_id":                      entity.AgreementId,
			"campaign_id":                       entity.CampaignId,
			"company_id":                        entity.CompanyId,
			"profit_center_id":                  entity.ProfitCenterId,
			"dealer_representative_id":          entity.DealerRepresentativeId,
			"customer_id":                       entity.CustomerId,
			"document_status_id":                entity.DocumentStatusId,
			"booking_estimation_batch_date":     entity.BookingEstimationBatchDate,
			"booking_estimation_vehicle_number": entity.BookingEstimationVehicleNumber,
			"agreement_number":                  entity.AgreementNumber,
			"is_unregistered":                   entity.IsUnregistered,
			"contact_person_name":               entity.ContactPersonName,
			"contact_person_phone":              entity.ContactPersonPhone,
			"contact_person_mobile":             entity.ContactPersonMobile,
			"contact_person_via_id":             entity.ContactPersonViaId,
			"insurance_policy_no":               entity.InsurancePolicyNo,
			"insurance_expired_date":            entity.InsuranceExpiredDate,
			"insurance_claim_no":                entity.InsuranceClaimNo,
			"insurance_pic":                     entity.InsurancePic,
		}
		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

// uspg_wtBookEstim0_Insert
// IF @Option = 0
// --USE FOR : * INSERT NEW DATA
func (r *BookingEstimationImpl) New(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationNormalRequest) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse) {
	var agreementId int
	var contractSystemNumber int

	// Fetch Agreement ID based on Customer ID
	err := tx.Model(&masterentities.Agreement{}).
		Select("agreement_id").
		Where("customer_id = ?", request.CustomerId).
		Scan(&agreementId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Agreement not found for the given customer",
		}
	} else if err != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve agreement from the database",
			Err:        err,
		}
	}

	// Fetch Contract System Number based on Vehicle ID
	err = tx.Model(&transactionworkshopentities.ContractService{}).
		Select("contract_service_system_number").
		Where("vehicle_id = ?", request.VehicleId).
		Scan(&contractSystemNumber).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Contract service not found for the given vehicle",
		}
	}

	if err != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve contract service from the database",
			Err:        err,
		}
	}

	loc, _ := time.LoadLocation("Asia/Jakarta") // UTC+7
	currentDate := time.Now().In(loc).Format("2006-01-02T15:04:05Z")
	parsedTime, _ := time.Parse(time.RFC3339, currentDate)

	entity := transactionworkshopentities.BookingEstimation{
		BrandId:                    request.BrandId,
		ModelId:                    request.ModelId,
		VariantId:                  request.VariantId,
		VehicleId:                  request.VehicleId,
		ContractSystemNumber:       request.ContractSystemNumber,
		AgreementId:                agreementId,
		CampaignId:                 request.CampaignId,
		CompanyId:                  request.CompanyId,
		ProfitCenterId:             request.ProfitCenterId,
		DealerRepresentativeId:     request.DealerRepresentativeId,
		CustomerId:                 request.CustomerId,
		BookingEstimationBatchDate: parsedTime,
		IsUnregistered:             request.Unregistered,
		InsurancePolicyNo:          request.InsurancePolicyNo,
		InsuranceExpiredDate:       request.InsuranceExpired,
		InsuranceClaimNo:           request.InsuranceClaimNo,
		InsurancePic:               request.InsurancePic,
		ContactPersonName:          request.NameCust,
		ContactPersonPhone:         request.PhoneCust,
		ContactPersonMobile:        request.MobileCust,
		ContactPersonViaId:         request.CallActivityId,
	}

	if err = tx.Create(&entity).Error; err != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save the new booking estimation",
			Err:        err,
		}
	}

	return entity, nil
}

// uspg_wtBookEstim0_Select
// IF @Option = 2
// --USE FOR : * SELECT DATA FOR SEARCHING
func (r *BookingEstimationImpl) GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.BookingEstimationResponseDetail, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.BookingEstimation
	err := tx.Where("batch_system_number = ?", Id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.BookingEstimationResponseDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Booking estimation not found",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.BookingEstimationResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve booking estimation",
			Err:        err,
		}
	}

	brandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(entity.BrandId)
	if brandErr != nil {
		return transactionworkshoppayloads.BookingEstimationResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve brand data",
			Err:        brandErr.Err,
		}
	}

	modelResponse, modelErr := salesserviceapiutils.GetUnitModelById(entity.ModelId)
	if modelErr != nil {
		return transactionworkshoppayloads.BookingEstimationResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve model data",
			Err:        modelErr.Err,
		}
	}

	variantResponse, variantErr := salesserviceapiutils.GetUnitVariantById(entity.VariantId)
	if variantErr != nil {
		return transactionworkshoppayloads.BookingEstimationResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve variant data",
			Err:        variantErr.Err,
		}
	}

	var bookingEstimationCampaigns []transactionworkshoppayloads.BookingEstimationCampaignResponse
	if err := tx.Model(&masterentities.CampaignMaster{}).
		Where("campaign_id = ? and campaign_id != 0", entity.CampaignId).
		Find(&bookingEstimationCampaigns).Error; err != nil {
		return transactionworkshoppayloads.BookingEstimationResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order campaigns from the database",
			Err:        err,
		}
	}

	var bookingEstimationAgreements []transactionworkshoppayloads.BookingEstimationAgreementResponse
	if err := tx.Model(&masterentities.Agreement{}).
		Where("agreement_id = ? and agreement_id != 0", entity.AgreementId).
		Find(&bookingEstimationAgreements).Error; err != nil {
		return transactionworkshoppayloads.BookingEstimationResponseDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order agreements from the database",
			Err:        err,
		}
	}

	response := transactionworkshoppayloads.BookingEstimationResponseDetail{
		BrandId:                    entity.BrandId,
		BrandName:                  brandResponse.BrandName,
		ModelId:                    entity.ModelId,
		ModelDescription:           modelResponse.ModelName,
		VariantId:                  entity.VariantId,
		VariantDescription:         variantResponse.VariantDescription,
		VehicleId:                  entity.VehicleId,
		ContractSystemNumber:       entity.ContractSystemNumber,
		AgreementId:                entity.AgreementId,
		CampaignId:                 entity.CampaignId,
		CompanyId:                  entity.CompanyId,
		ProfitCenterId:             entity.ProfitCenterId,
		DealerRepresentativeId:     entity.DealerRepresentativeId,
		CustomerId:                 entity.CustomerId,
		NameCust:                   entity.ContactPersonName,
		PhoneCust:                  entity.ContactPersonPhone,
		MobileCust:                 entity.ContactPersonMobile,
		CallActivityId:             entity.ContactPersonViaId,
		BookingEstimationBatchDate: entity.BookingEstimationBatchDate,

		BookingEstimationCampaign: transactionworkshoppayloads.BookingEstimationCampaignDetail{
			DataCampaign: bookingEstimationCampaigns,
		},
		BookingEstimationAgreement: transactionworkshoppayloads.BookingEstimationAgreement{
			DataAgreement: bookingEstimationAgreements,
		},
	}

	return response, nil
}

// uspg_wtBookEstim0_Update
// IF @Option = 0
// --USE FOR : * UPDATE DATA BY KEY
func (r *BookingEstimationImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationSaveRequest, id int) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.BookingEstimation

	err := tx.Where("batch_system_number = ?", id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Booking estimation not found",
			}
		}
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve booking from the database",
			Err:        err,
		}
	}

	updates := make(map[string]interface{})

	if request.CompanyId != 0 {
		updates["company_id"] = request.CompanyId
	}
	if request.DealerRepresentativeId != 0 {
		updates["dealer_representative_id"] = request.DealerRepresentativeId
	}
	if request.CampaignId != 0 {
		updates["campaign_id"] = request.CampaignId
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
	if request.ContactViaId != 0 {
		updates["contact_person_via_id"] = request.ContactViaId
	}
	if request.InsurancePolicyNo != "" {
		updates["insurance_policy_no"] = request.InsurancePolicyNo
	}
	if !request.InsuranceExpired.IsZero() {
		updates["insurance_expired_date"] = request.InsuranceExpired.Truncate(24 * time.Hour)
	}
	if request.InsuranceClaimNo != "" {
		updates["insurance_claim_no"] = request.InsuranceClaimNo
	}
	if request.InsurancePic != "" {
		updates["insurance_pic"] = request.InsurancePic
	}

	if len(updates) == 0 {
		return entity, nil
	}

	err = tx.Model(&entity).Updates(updates).Error
	if err != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save the updated booking",
			Err:        err,
		}
	}

	return entity, nil
}

func (r *BookingEstimationImpl) Submit(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	// Retrieve the booking estimation by Id
	var entity transactionworkshopentities.BookingEstimation
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve booking estimation from the database"}
	}

	// Perform the necessary operations to submit the booking estimation
	// ...

	// Save the updated booking estimation
	err = tx.Save(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the updated booking estimation"}
	}

	return true, nil
}

// uspg_wtBookEstim0_Delete
// IF @Option = 0
// --USE FOR : * DELETE DATA BY KEY
func (r *BookingEstimationImpl) Void(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	type BookingEstimation struct {
		BatchSystemNumber          int `gorm:"column:batch_system_number"`
		ServiceRequestSystemNumber int `gorm:"column:service_request_system_number"`
		PdiSystemNumber            int `gorm:"column:pdi_system_number"`
		PdiLine                    int `gorm:"column:pdi_line"`
		BookingSystemNumber        int `gorm:"column:booking_system_number"`
		EstimationSystemNumber     int `gorm:"column:estimation_system_number"`
	}

	var booking BookingEstimation
	err := tx.Table("trx_booking_estimation").
		Select("service_request_system_number, pdi_system_number, pdi_line, booking_system_number, estimation_system_number").
		Where("batch_system_number = ?", Id).
		First(&booking).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Booking not found",
				Err:        err,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve booking details",
			Err:        err,
		}
	}

	var bookAllocSysNo int
	err = tx.Table("trx_booking_estimation_allocation").
		Where("booking_system_number = ?", booking.BookingSystemNumber).
		Select("COALESCE(booking_allocation_system_number, 0)").
		Scan(&bookAllocSysNo).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve booking allocation",
			Err:        err,
		}
	}

	var estimDiscStatus int
	err = tx.Table("trx_booking_estimation_service_discount").
		Where("estimation_system_number = ?", booking.EstimationSystemNumber).
		Select("estimation_discount_approval_status").
		Scan(&estimDiscStatus).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve estimation status",
			Err:        err,
		}
	}

	var existingAlloc int
	err = tx.Table("trx_work_order_allocation").
		Where("booking_system_number = ?", booking.BookingSystemNumber).
		Select("COUNT(1)").
		Scan(&existingAlloc).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check booking allocation",
			Err:        err,
		}
	}
	if existingAlloc > 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Can't void document, document already allocated",
		}
	}

	if estimDiscStatus == 2 { //"APPROVAL_WAITAPPROVED"
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Document still waiting for approval",
		}
	}

	if booking.ServiceRequestSystemNumber != 0 {
		err = tx.Table("trx_service_request").
			Where("service_request_system_number = ?", booking.ServiceRequestSystemNumber).
			Updates(map[string]interface{}{
				"booking_system_number":     0,
				"service_request_status_id": 3, // "SERV_REQ_STAT_ACCEPT"
			}).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update service request",
				Err:        err,
			}
		}
	}

	if booking.PdiSystemNumber != 0 && booking.PdiLine != 0 {
		err = tx.Table("dms_microservices_sales_dev.dbo.trx_pdi_request_detail").
			Where("pdi_request_system_number = ? AND pdi_request_detail_line_number = ?", booking.PdiSystemNumber, booking.PdiLine).
			Updates(map[string]interface{}{
				"booking_system_number":             0,
				"service_date":                      nil,
				"service_time":                      "",
				"pdi_request_detail_line_status_id": 3, // "PDI_REQ_DET_STAT_ACCEPT"
			}).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update PDI records",
				Err:        err,
			}
		}
	}

	if bookAllocSysNo != 0 {
		err = tx.Table("trx_booking_allocation").
			Where("booking_allocation_system_number = ?", bookAllocSysNo).
			Delete(nil).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to delete booking allocation",
				Err:        err,
			}
		}
	}

	err = tx.Table("trx_booking_estimation").Where("batch_system_number = ?", Id).Delete(nil).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete records from trx_booking_estimation",
			Err:        err,
		}
	}

	err = tx.Table("trx_booking_estimation_allocation").Where("booking_system_number = ?", booking.BookingSystemNumber).Delete(nil).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete records from trx_booking_estimation_allocation",
			Err:        err,
		}
	}

	err = tx.Table("trx_booking_estimation_request").Where("booking_system_number = ?", booking.BookingSystemNumber).Delete(nil).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete records from trx_booking_estimation_request",
			Err:        err,
		}
	}

	err = tx.Table("trx_booking_estimation_service_discount").Where("estimation_system_number = ?", booking.EstimationSystemNumber).Delete(nil).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete records from trx_booking_estimation_service_discount",
			Err:        err,
		}
	}

	err = tx.Table("trx_booking_estimation_detail").Where("estimation_system_number = ?", booking.EstimationSystemNumber).Delete(nil).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete records from trx_booking_estimation_detail",
			Err:        err,
		}
	}

	return true, nil
}

func (r *BookingEstimationImpl) CloseOrder(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	// Retrieve the booking estimation by Id
	var entity transactionworkshopentities.BookingEstimation
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Where("batch_system_number = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve booking estimation from the database"}
	}

	// Perform the necessary operations to close the booking estimation
	// ...

	// Save the updated booking estimation
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated booking estimation"}
	}

	return nil
}

func (r *BookingEstimationImpl) GetByIdBookEstimReq(tx *gorm.DB, booksysno int, id int) (transactionworkshoppayloads.BookEstimRemarkResponse, *exceptions.BaseErrorResponse) {
	var model transactionworkshopentities.BookingEstimationRequest

	if err := tx.Where("booking_system_number = ? AND booking_estimation_request_id = ?", booksysno, id).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.BookEstimRemarkResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Booking estimation request not found",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.BookEstimRemarkResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve booking estimation request",
			Err:        err,
		}
	}

	payloads := transactionworkshoppayloads.BookEstimRemarkResponse{
		BookingServiceRequest: model.BookingServiceRequest,
		BookingDocumentNumber: model.BookingDocumentNumber,
		BookingLine:           model.BookingLine,
		IsActive:              model.IsActive,
	}

	return payloads, nil
}

func (r *BookingEstimationImpl) SaveBookEstimReq(tx *gorm.DB, req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (transactionworkshopentities.BookingEstimationRequest, *exceptions.BaseErrorResponse) {
	var maxBookingLine *int

	err := tx.Model(&transactionworkshopentities.BookingEstimationRequest{}).
		Where("booking_system_number = ?", id).
		Select("MAX(booking_line)").
		Scan(&maxBookingLine).Error
	if err != nil {
		return transactionworkshopentities.BookingEstimationRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve max booking line",
			Err:        err,
		}
	}

	newBookingLine := 1
	if maxBookingLine != nil {
		newBookingLine = *maxBookingLine + 1
	}

	newBooking := transactionworkshopentities.BookingEstimationRequest{
		BookingLine:           newBookingLine,
		BookingSystemNumber:   id,
		BookingServiceRequest: req.BookingServiceRequest,
		IsActive:              true,
	}

	if err := tx.Create(&newBooking).Error; err != nil {
		return transactionworkshopentities.BookingEstimationRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save booking estimation request",
			Err:        err,
		}
	}

	return newBooking, nil
}

func (r *BookingEstimationImpl) UpdateBookEstimReq(tx *gorm.DB, req transactionworkshoppayloads.BookEstimRemarkRequest, booksysno int, id int) (transactionworkshopentities.BookingEstimationRequest, *exceptions.BaseErrorResponse) {
	var model transactionworkshopentities.BookingEstimationRequest

	if err := tx.Where("booking_system_number = ? AND booking_estimation_request_id = ?", booksysno, id).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshopentities.BookingEstimationRequest{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Booking estimation request not found",
				Err:        err,
			}
		}
		return transactionworkshopentities.BookingEstimationRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve booking estimation request",
			Err:        err,
		}
	}

	if err := tx.Model(&model).
		Where("booking_system_number = ? AND booking_estimation_request_id = ?", booksysno, id).
		Updates(map[string]interface{}{
			"booking_service_request": req.BookingServiceRequest,
		}).Error; err != nil {
		return transactionworkshopentities.BookingEstimationRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update booking estimation request",
			Err:        err,
		}
	}

	return model, nil
}

func (r *BookingEstimationImpl) DeleteBookEstimReq(tx *gorm.DB, booksysno int, ids []int) (bool, *exceptions.BaseErrorResponse) {
	result := tx.Where("booking_system_number = ? AND booking_estimation_request_id IN (?)", booksysno, ids).
		Delete(&transactionworkshopentities.BookingEstimationRequest{})

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete booking estimation request",
			Err:        result.Error,
		}
	}

	if result.RowsAffected == 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No matching booking estimation request found",
		}
	}

	return true, nil
}

func (r *BookingEstimationImpl) GetAllBookEstimReq(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []transactionworkshopentities.BookingEstimationRequest
	pages.Rows = []transactionworkshopentities.BookingEstimationRequest{}

	baseQuery := tx.Model(&transactionworkshopentities.BookingEstimationRequest{})
	filteredQuery := utils.ApplyFilter(baseQuery, filterCondition)

	if err := filteredQuery.Scopes(pagination.Paginate(&pages, filteredQuery)).
		Find(&responses).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve booking estimation requests from the database",
			Err:        err,
		}
	}

	pages.Rows = responses
	return pages, nil
}

func (r *BookingEstimationImpl) SaveBookEstimReminderServ(tx *gorm.DB, req transactionworkshoppayloads.ReminderServicePost, id int) (transactionworkshopentities.BookingEstimationServiceReminder, *exceptions.BaseErrorResponse) {
	var maxBookingLine int

	err := tx.Model(&transactionworkshopentities.BookingEstimationServiceReminder{}).
		Where("booking_system_number = ?", id).
		Select("COALESCE(MAX(booking_line_number), 0)").
		Scan(&maxBookingLine).Error
	if err != nil {
		return transactionworkshopentities.BookingEstimationServiceReminder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve max booking line number",
			Err:        err,
		}
	}

	newEntity := transactionworkshopentities.BookingEstimationServiceReminder{
		BookingLineNumber:      maxBookingLine + 1,
		BookingSystemNumber:    id,
		BookingServiceReminder: req.BookingServiceReminder,
	}

	if err := tx.Create(&newEntity).Error; err != nil {
		return transactionworkshopentities.BookingEstimationServiceReminder{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save booking estimation reminder service",
			Err:        err,
		}
	}

	return newEntity, nil
}

// uspg_wtBookEstim2_1_Select
// IF @Option = 0
// --USE FOR : * SELECT DATA BY KEY
func (r *BookingEstimationImpl) GetAllDetailBookingEstimation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var tableStruct []transactionworkshoppayloads.BookingEstimationDetailRequest

	baseModelQuery := tx.Model(&transactionworkshopentities.BookingEstimationDetail{})
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&tableStruct).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve booking estimation details",
			Err:        err,
		}
	}

	if len(tableStruct) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var convertedResponses []transactionworkshoppayloads.BookingEstimationDetailResponse

	for _, estimationReq := range tableStruct {
		lineTypeResponse, lineTypeErr := generalserviceapiutils.GetLineTypeById(estimationReq.LineTypeId)
		if lineTypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve line type from the external API",
				Err:        lineTypeErr.Err,
			}
		}

		jobTypeResponse, jobTypeErr := generalserviceapiutils.GetJobTransactionTypeById(estimationReq.JobTypeId)
		if jobTypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve job type from the external API",
				Err:        jobTypeErr.Err,
			}
		}

		transactionTypeResponse, transactionTypeErr := generalserviceapiutils.GetWoTransactionTypeById(estimationReq.TransactionTypeId)
		if transactionTypeErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve transaction type from the external API",
				Err:        transactionTypeErr.Err,
			}
		}

		var OperationItemCode string
		var Description string

		operationItemResponse, operationItemErr := r.workorderRepo.GetOperationItemById(estimationReq.LineTypeId, estimationReq.OperationItemId)
		if operationItemErr != nil {
			return pages, operationItemErr
		}

		OperationItemCode, Description, errResp := r.workorderRepo.HandleLineTypeResponse(estimationReq.LineTypeId, operationItemResponse)
		if errResp != nil {
			return pages, errResp
		}

		// fetch data item
		var itemResponse masteritementities.Item
		if err := tx.Where("item_id = ?", estimationReq.OperationItemId).First(&itemResponse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item not found",
					Err:        fmt.Errorf("item with ID %d not found", estimationReq.OperationItemId),
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

		estimationRes := transactionworkshoppayloads.BookingEstimationDetailResponse{
			EstimationDetailId:          estimationReq.EstimationDetailId,
			EstimationSystemNumber:      estimationReq.EstimationSystemNumber,
			TransactionTypeId:           estimationReq.TransactionTypeId,
			TransactionTypeCode:         transactionTypeResponse.WoTransactionTypeCode,
			JobTypeId:                   estimationReq.JobTypeId,
			JobTypeCode:                 jobTypeResponse.JobTypeCode,
			LineTypeId:                  estimationReq.LineTypeId,
			LineTypeCode:                lineTypeResponse.LineTypeCode,
			LineTypeName:                lineTypeResponse.LineTypeName,
			FrtQuantity:                 estimationReq.FrtQuantity,
			OperationItemId:             estimationReq.OperationItemId,
			OperationItemCode:           OperationItemCode,
			Description:                 Description,
			Uom:                         uomItems.UomDescription,
			OperationItemPrice:          estimationReq.OperationItemPrice,
			OperationItemDiscountAmount: estimationReq.OperationItemDiscountAmount,
		}

		convertedResponses = append(convertedResponses, estimationRes)
	}

	var mapResponses []map[string]interface{}
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"estimation_detail_id":           response.EstimationDetailId,
			"estimation_system_number":       response.EstimationSystemNumber,
			"line_type_id":                   response.LineTypeId,
			"line_type_code":                 response.LineTypeCode,
			"transaction_type_id":            response.TransactionTypeId,
			"transaction_type_code":          response.TransactionTypeCode,
			"job_type_id":                    response.JobTypeId,
			"frt_quantity":                   response.FrtQuantity,
			"operation_item_id":              response.OperationItemId,
			"operation_item_code":            response.OperationItemCode,
			"operation_item_price":           response.OperationItemPrice,
			"operation_item_discount_amount": response.OperationItemDiscountAmount,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	pages.Rows = mapResponses
	return pages, nil
}

// uspg_wtBookEstim2_1_Select
// IF @Option = 1
// --USE FOR : * SELECT DATA BY KEY
func (r *BookingEstimationImpl) GetByIdBookEstimDetail(tx *gorm.DB, estimsysno int, id int) (transactionworkshoppayloads.BookingEstimationDetailResponse, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.BookingEstimationDetail
	err := tx.Model(&transactionworkshopentities.BookingEstimationDetail{}).
		Where("estimation_detail_id = ? AND estimation_system_number = ?", id, estimsysno).
		First(&entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactionworkshoppayloads.BookingEstimationDetailResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Booking estimation detail not found",
				Err:        err,
			}
		}
		return transactionworkshoppayloads.BookingEstimationDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve booking estimation detail from the database",
			Err:        err,
		}
	}

	// Fetch external data
	lineTypeResponse, lineErr := generalserviceapiutils.GetLineTypeById(entity.LineTypeId)
	if lineErr != nil {
		return transactionworkshoppayloads.BookingEstimationDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve line type from the external API",
			Err:        lineErr.Err,
		}
	}

	transactionTypeResponse, transactionTypeErr := generalserviceapiutils.GetWoTransactionTypeById(entity.TransactionTypeId)
	if transactionTypeErr != nil {
		return transactionworkshoppayloads.BookingEstimationDetailResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve transaction type from the external API",
			Err:        transactionTypeErr.Err,
		}
	}

	payload := transactionworkshoppayloads.BookingEstimationDetailResponse{
		EstimationDetailId:          entity.EstimationDetailId,
		EstimationLine:              entity.EstimationLine,
		EstimationSystemNumber:      entity.EstimationSystemNumber,
		EstimationDocumentNumber:    entity.EstimationDocumentNumber,
		TransactionTypeId:           entity.TransactionTypeId,
		TransactionTypeCode:         transactionTypeResponse.WoTransactionTypeCode,
		JobTypeId:                   entity.JobTypeId,
		LineTypeId:                  entity.LineTypeId,
		LineTypeCode:                lineTypeResponse.LineTypeCode,
		LineTypeName:                lineTypeResponse.LineTypeName,
		OperationItemId:             entity.OperationItemId,
		OperationItemCode:           entity.OperationItemCode,
		FrtQuantity:                 entity.FRTQuantity,
		OperationItemPrice:          entity.OperationItemPrice,
		OperationItemDiscountAmount: entity.OperationItemDiscountAmount,
	}

	return payload, nil
}

func (r *BookingEstimationImpl) SaveDetailBookEstim(tx *gorm.DB, id int, request transactionworkshoppayloads.BookingEstimationDetailRequest) (transactionworkshopentities.BookingEstimationDetail, *exceptions.BaseErrorResponse) {
	var bookingDetail transactionworkshopentities.BookingEstimationDetail
	var estimationServiceDiscount transactionworkshopentities.BookingEstimationServiceDiscount

	if err := tx.Where("estimation_system_number = ?", id).First(&estimationServiceDiscount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			var bookingEstimation transactionworkshopentities.BookingEstimation
			if err := tx.Where("batch_system_number = ?", bookingEstimation.BatchSystemNumber).First(&bookingEstimation).Error; err != nil {
				return transactionworkshopentities.BookingEstimationDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to fetch booking estimation",
					Err:        err,
				}
			}

			// VAT Tax Rate
			taxResponse, taxErr := financeserviceapiutils.GetTaxPercent(financeserviceapiutils.TaxPercentParams{
				TaxServiceCode: "PPN",
				TaxTypeCode:    "PPN",
				EffectiveDate:  time.Now(),
			})
			if taxErr != nil {
				return transactionworkshopentities.BookingEstimationDetail{}, taxErr
			}

			estimationServiceDiscount = transactionworkshopentities.BookingEstimationServiceDiscount{
				EstimationSystemNumber: id,
				BatchSystemNumber:      bookingEstimation.BatchSystemNumber,
				VATTaxRate:             taxResponse.TaxPercent,
				CompanyID:              bookingEstimation.CompanyId,
			}
			if err := tx.Create(&estimationServiceDiscount).Error; err != nil {
				return transactionworkshopentities.BookingEstimationDetail{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    "Failed to insert booking estimation service discount",
					Err:        err,
				}
			}
		} else {
			return transactionworkshopentities.BookingEstimationDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Database error",
				Err:        err,
			}
		}
	}

	var existingDetail transactionworkshopentities.BookingEstimationDetail
	if err := tx.Where("estimation_system_number = ? AND operation_item_id = ?", id, request.OperationItemId).First(&existingDetail).Error; err == nil {
		return transactionworkshopentities.BookingEstimationDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Item already exists",
		}
	}

	bookingDetail = transactionworkshopentities.BookingEstimationDetail{
		EstimationSystemNumber:              id,
		OperationItemCode:                   request.OperationItemCode,
		LineTypeId:                          request.LineTypeId,
		FRTQuantity:                         request.FrtQuantity,
		OperationItemPrice:                  request.OperationItemPrice,
		OperationItemDiscountAmount:         request.OperationItemDiscountAmount,
		OperationItemDiscountRequestAmount:  request.OperationItemDiscountRequestAmount,
		OperationItemDiscountPercent:        request.OperationItemDiscountPercent,
		OperationItemDiscountRequestPercent: request.OperationItemDiscountRequestPercent,
	}

	if err := tx.Create(&bookingDetail).Error; err != nil {
		return transactionworkshopentities.BookingEstimationDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create booking estimation detail",
			Err:        err,
		}
	}

	return bookingDetail, nil
}

func (r *BookingEstimationImpl) UpdateBookEstimDetail(tx *gorm.DB, req transactionworkshoppayloads.BookEstimDetailUpdate, id int, LineTypeId int) (bool, *exceptions.BaseErrorResponse) {
	return true, nil
}

func (r *BookingEstimationImpl) DeleteBookEstimDetail(tx *gorm.DB, id, linetypeid int) (bool, *exceptions.BaseErrorResponse) {

	return true, nil
}

func (r *BookingEstimationImpl) CopyFromHistory(tx *gorm.DB, batchid int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	return nil, nil
}

func (r *BookingEstimationImpl) InputDiscount(tx *gorm.DB, id int, req transactionworkshoppayloads.BookEstimationPayloadsDiscount) (int, *exceptions.BaseErrorResponse) {
	return 0, nil
}

func (r *BookingEstimationImpl) AddFieldAction(tx *gorm.DB, id int, idrecall int) (int, *exceptions.BaseErrorResponse) {
	return 0, nil
}

func (r *BookingEstimationImpl) PostBookingEstimationCalculation(tx *gorm.DB, id int) (int, *exceptions.BaseErrorResponse) {
	var estimationSystemNumber int

	err := tx.Table("trx_booking_estimation_service_discount").
		Select("estimation_system_number").
		Where("batch_system_number = ?", id).
		Scan(&estimationSystemNumber).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch estimation system number",
			Err:        err,
		}
	}

	if estimationSystemNumber != 0 {
		return estimationSystemNumber, nil
	}

	now := time.Now()
	entity := transactionworkshopentities.BookingEstimationServiceDiscount{
		BatchSystemNumber:    id,
		EstimationDate:       &now,
		DiscountApprovalDate: &now,
	}

	if err := tx.Create(&entity).Error; err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create estimation system number",
			Err:        err,
		}
	}

	return entity.EstimationSystemNumber, nil
}

func (r *BookingEstimationImpl) PutBookingEstimationCalculation(tx *gorm.DB, id int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {

	currentTime := time.Now()
	taxResponse, errTax := financeserviceapiutils.GetTaxPercent(financeserviceapiutils.TaxPercentParams{
		TaxServiceCode: "PPN",
		TaxTypeCode:    "PPN",
		EffectiveDate:  currentTime,
	})
	if errTax != nil {
		return nil, errTax
	}
	taxfare := taxResponse.TaxPercent

	const (
		LineTypePackage            = 0
		LineTypeOperation          = 1
		LineTypeSparePart          = 2
		LineTypeOil                = 3
		LineTypeMaterial           = 4
		LineTypeAccessories        = 6
		LineTypeConsumableMaterial = 7
		LineTypeSublet             = 8
	)
	type Result struct {
		TotalPackage            float64
		TotalOperation          float64
		TotalSparePart          float64
		TotalOil                float64
		TotalMaterial           float64
		TotalAccessories        float64
		TotalConsumableMaterial float64
		TotalSublet             float64
	}

	var result Result

	err := tx.Model(&transactionworkshopentities.BookingEstimationServiceDiscount{}).
		Select(`SUM(CASE WHEN line_type_id = ? THEN ROUND(COALESCE(trx_booking_estimation_detail.operation_item_price, 0) * COALESCE(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_package,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(COALESCE(trx_booking_estimation_detail.operation_item_price, 0) * COALESCE(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_operation,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(COALESCE(trx_booking_estimation_detail.operation_item_price, 0) * COALESCE(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_spare_part,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(COALESCE(trx_booking_estimation_detail.operation_item_price, 0) * COALESCE(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_oil,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(COALESCE(trx_booking_estimation_detail.operation_item_price, 0) * COALESCE(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_material,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(COALESCE(trx_booking_estimation_detail.operation_item_price, 0) * COALESCE(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_accessories,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(COALESCE(trx_booking_estimation_detail.operation_item_price, 0) * COALESCE(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_consumable_material,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(COALESCE(trx_booking_estimation_detail.operation_item_price, 0) * COALESCE(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_sublet`,
			LineTypePackage,
			LineTypeOperation,
			LineTypeSparePart,
			LineTypeOil,
			LineTypeMaterial,
			LineTypeAccessories,
			LineTypeConsumableMaterial,
			LineTypeSublet).
		Joins("JOIN trx_booking_estimation_detail ON trx_booking_estimation_detail.estimation_system_number = trx_booking_estimation_service_discount.estimation_system_number").
		Where("batch_system_number = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to calculate totals",
			Err:        err,
		}
	}

	// Calculate grand total
	total := result.TotalPackage + result.TotalOperation + result.TotalSparePart + result.TotalOil + result.TotalMaterial + result.TotalAccessories + result.TotalConsumableMaterial + result.TotalSublet
	tax := total * (taxfare / 100)
	grandTotal := total + tax

	err = tx.Model(&transactionworkshopentities.BookingEstimationServiceDiscount{}).
		Where("batch_system_number = ?", id).
		Updates(map[string]interface{}{
			"total_price_package":             result.TotalPackage,
			"total_price_operation":           result.TotalOperation,
			"total_price_part":                result.TotalSparePart,
			"total_price_oil":                 result.TotalOil,
			"total_price_material":            result.TotalMaterial,
			"total_price_accessories":         result.TotalAccessories,
			"total_price_consumable_material": result.TotalConsumableMaterial,
			"total_sublet":                    result.TotalSublet,
			"total_vat":                       tax,
			"total_after_vat":                 grandTotal,
		}).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update work order",
			Err:        err,
		}
	}

	bookingEstimationResponse := []map[string]interface{}{
		{"total_package": result.TotalPackage},
		{"total_operation": result.TotalOperation},
		{"total_spare_part": result.TotalSparePart},
		{"total_oil": result.TotalOil},
		{"total_material": result.TotalMaterial},
		{"total_accessories": result.TotalAccessories},
		{"total_consumable_material": result.TotalConsumableMaterial},
		{"total_sublet": result.TotalSublet},
		{"total_vat": tax},
		{"total_after_vat": grandTotal},
	}

	return bookingEstimationResponse, nil
}

func (r *BookingEstimationImpl) SaveBookingEstimationAllocation(tx *gorm.DB, id int, req transactionworkshoppayloads.BookEstimationAllocation) (transactionworkshopentities.BookingEstimationAllocation, *exceptions.BaseErrorResponse) {

	if id <= 0 {
		return transactionworkshopentities.BookingEstimationAllocation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid batch system number",
		}
	}

	if req.BookingStatusID == 0 || req.CompanyID == 0 || req.BookingDocumentNumber == "" {
		return transactionworkshopentities.BookingEstimationAllocation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Missing required fields: BookingStatusID, CompanyID, or BookingDocumentNumber",
		}
	}

	entities := transactionworkshopentities.BookingEstimationAllocation{
		BookingStatusID:       req.BookingStatusID,
		BatchSystemNumber:     id,
		CompanyID:             req.CompanyID,
		PdiSystemNumber:       req.PdiSystemNumber,
		BookingDocumentNumber: req.BookingDocumentNumber,
		BookingDate:           req.BookingDate,
		BookingStall:          req.BookingStall,
		BookingReminderDate:   req.BookingReminderDate,
		BookingServiceDate:    req.BookingServiceDate,
		BookingServiceTime:    req.BookingServiceTime,
		BookingEstimationTime: req.BookingEstimationTime,
	}

	err := tx.Create(&entities).Error
	if err != nil {
		return transactionworkshopentities.BookingEstimationAllocation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Failed to save booking estimation allocation",
			Err:        err,
		}
	}
	return entities, nil
}

func (r *BookingEstimationImpl) AddContractService(tx *gorm.DB, idheader int, Idcontract int) (bool, *exceptions.BaseErrorResponse) {
	return true, nil
}

func (r *BookingEstimationImpl) AddPackage(tx *gorm.DB, idhead int, idpackage int) (bool, *exceptions.BaseErrorResponse) {
	return true, nil
}

func (r *BookingEstimationImpl) SaveBookingEstimationFromServiceRequest(tx *gorm.DB, idservreq int, req transactionworkshoppayloads.PdiServiceRequest) (bool, *exceptions.BaseErrorResponse) {
	return true, nil
}

func (r *BookingEstimationImpl) SaveBookingEstimationFromPDI(tx *gorm.DB, idpdi int, req transactionworkshoppayloads.PdiServiceRequest) (bool, *exceptions.BaseErrorResponse) {

	return true, nil
}
