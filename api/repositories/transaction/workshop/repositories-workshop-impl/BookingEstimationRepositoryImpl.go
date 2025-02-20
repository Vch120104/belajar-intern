package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	masterentities "after-sales/api/entities/master"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	transactionunitpayloads "after-sales/api/payloads/transaction/unit"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BookingEstimationImpl struct {
}

func OpenBookingEstimationRepositoryImpl() transactionworkshoprepository.BookingEstimationRepository {
	return &BookingEstimationImpl{}
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
			"agreement_number_br":               entity.AgreementNumberBr,
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

	// Validation: request date
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

	// Fetch work order agreements
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Booking estimation not found",
			Err:        err,
		}
	}

	if err != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve booking from the database",
			Err:        err,
		}
	}

	entity.CompanyId = request.CompanyId
	entity.DealerRepresentativeId = request.DealerRepresentativeId
	entity.CampaignId = request.CampaignId
	entity.ContactPersonName = request.NameCust
	entity.ContactPersonPhone = request.PhoneCust
	entity.ContactPersonMobile = request.MobileCust
	entity.ContactPersonViaId = request.ContactViaId
	entity.InsurancePolicyNo = request.InsurancePolicyNo
	entity.InsuranceExpiredDate = request.InsuranceExpired
	entity.InsuranceClaimNo = request.InsuranceClaimNo
	entity.InsurancePic = request.InsurancePic

	err = tx.Save(&entity).Error
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

	// Cek apakah booking sudah dialokasikan
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

	// Cek apakah estimasi masih dalam status approval
	if estimDiscStatus == 2 { //"APPROVAL_WAITAPPROVED"
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Document still waiting for approval",
		}
	}

	// Update service request
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

	// Update PDI jika ada
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

	// Hapus booking allocation
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

func (r *BookingEstimationImpl) DeleteBookEstimReq(tx *gorm.DB, ids string) ([]string, *exceptions.BaseErrorResponse) {

	idSlice := strings.Split(ids, ",")

	var models []transactionworkshopentities.BookingEstimationRequest
	if err := tx.Where("booking_estimation_request_id IN (?)", idSlice).Find(&models).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch booking estimation requests",
			Err:        err,
		}
	}

	if len(models) == 0 {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No booking estimation requests found to delete",
			Err:        nil,
		}
	}

	var deletedIds []string
	for _, model := range models {
		deletedIds = append(deletedIds, strconv.Itoa(model.BookingEstimationRequestID))
	}

	if err := tx.Delete(&models).Error; err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete booking estimation requests",
			Err:        err,
		}
	}

	if err := updateLineNumbers(tx); err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update line numbers",
			Err:        err,
		}
	}

	return deletedIds, nil
}

func updateLineNumbers(tx *gorm.DB) error {
	var records []transactionworkshopentities.BookingEstimationRequest

	if err := tx.Model(&transactionworkshopentities.BookingEstimationRequest{}).
		Order("booking_estimation_request_code ASC").
		Find(&records).Error; err != nil {
		return err
	}

	for i := range records {
		records[i].BookingLine = i + 1
	}

	if len(records) > 0 {
		if err := tx.Save(&records).Error; err != nil {
			return err
		}
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

func (r *BookingEstimationImpl) SaveDetailBookEstim(tx *gorm.DB, req transactionworkshoppayloads.BookEstimDetailReq, id int) (transactionworkshopentities.BookingEstimationDetail, *exceptions.BaseErrorResponse) {
	var lastprice float64
	var entity transactionworkshopentities.BookingEstimation
	var count int64
	// repo := masterrepositoryimpl.LookupRepositoryImpl{}
	result, err := r.PostBookingEstimationCalculation(tx, id)
	if err != nil {
		return transactionworkshopentities.BookingEstimationDetail{}, err
	}

	if req.LineTypeID != 9 && req.LineTypeID != 0 {
		err := tx.Select("mtr_item_price_list.price_list_amount").Table("mtr_item_price_list").
			Joins("JOIN mtr_item on mtr_item.item_id=mtr_item_price_list.item_id").
			Joins("Join mtr_item_operation on mtr_item.item_id=mtr_item_operation.item_operation_model_mapping_id").
			Where("item_operation_id=?", req.ItemOperationID).
			Scan(&lastprice).Error
		if err != nil {
			return transactionworkshopentities.BookingEstimationDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
	} else {
		err := tx.Select("mtr_labour_selling_price_detail.selling_price").
			Table("mtr_labour_selling_price_detail").
			Joins("Join mtr_labour_selling_price on mtr_labour_selling_price.labour_selling_price_id = mtr_labour_selling_price_detail.labour_selling_price_id").
			Where("mtr_labour_selling_price.brand_id =?", entity.BrandId).
			Where("mtr_labour_selling_price_detail.model_id=?", entity.ModelId).
			Where("mtr_labour_selling_price.company_id = ?", entity.CompanyId).
			Where("mtr_labour_selling_price.effective_date < ?", time.Now()).
			Scan(&lastprice).Error
		if err != nil {
			return transactionworkshopentities.BookingEstimationDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
	}
	err3 := tx.Select("estimation_line_id").Table("trx_booking_estimation_detail").Where("estimation_system_number=?", result).Count(&count).Error
	if err3 != nil {
		return transactionworkshopentities.BookingEstimationDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err3,
		}
	}
	// price,err:= repo.GetOprItemPrice(tx,req.LineTypeID,entity.CompanyId,req.ItemOperationID,entity.BrandId,entity.ModelId,req.JobTypeID,entity.VariantId)
	entities := transactionworkshopentities.BookingEstimationDetail{
		EstimationLineID:               req.EstimationLineID,
		EstimationLineCode:             int(count + 1),
		EstimationSystemNumber:         result,
		BillID:                         req.BillID,
		EstimationLineDiscountApproval: req.EstimationLineDiscountApproval,
		ItemOperationID:                req.ItemOperationID,
		LineTypeID:                     req.LineTypeID,
		PackageID:                      req.PackageID,
		JobTypeID:                      req.JobTypeID,
		FieldActionSystemNumber:        req.FieldActionSystemNumber,
		ApprovalRequestNumber:          req.ApprovalRequestNumber,
		UOMID:                          req.UOMID,
		RequestDescription:             req.RequestDescription,
		FRTQuantity:                    req.FRTQuantity,
		ItemOperationPrice:             lastprice,
		DiscountItemOperationAmount:    req.DiscountItemAmount,
		DiscountItemOperationPercent:   req.DiscountItemPercent,
		DiscountRequestPercent:         req.DiscountRequestPercent,
		DiscountRequestAmount:          req.DiscountRequestAmount,
	}

	if err := tx.Save(&entities).Error; err != nil {
		return transactionworkshopentities.BookingEstimationDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	return entities, nil
}

func (r *BookingEstimationImpl) UpdateBookEstimDetail(tx *gorm.DB, req transactionworkshoppayloads.BookEstimDetailUpdate, id int, LineTypeId int) (bool, *exceptions.BaseErrorResponse) {
	var model transactionworkshopentities.BookingEstimationDetail

	if err := tx.First(&model, "estimation_line_id = ?", id).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Booking estimation detail not found",
			Err:        err,
		}
	}

	if err := tx.Model(&model).Where("estimation_line_id = ?", id).Updates(map[string]interface{}{
		"frt_quantity":             req.FRTQuantity,
		"discount_request_percent": req.DiscountRequestPercent,
	}).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	_, err := r.PutBookingEstimationCalculation(tx, model.EstimationSystemNumber)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *BookingEstimationImpl) DeleteBookEstimDetail(tx *gorm.DB, id int, linetypeid int) (bool, *exceptions.BaseErrorResponse) {
	var model transactionworkshopentities.BookingEstimationDetail

	result := tx.Where("estimation_line_id = ? AND line_type_id = ?", id, linetypeid).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Booking estimation detail not found",
				Err:        result.Error,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check booking estimation detail",
			Err:        result.Error,
		}
	}

	deleteResult := tx.Where("estimation_line_id = ? AND line_type_id = ?", id, linetypeid).Delete(&model)
	if deleteResult.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete booking estimation detail",
			Err:        deleteResult.Error,
		}
	}

	if deleteResult.RowsAffected == 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No booking estimation detail deleted",
			Err:        errors.New("no rows affected"),
		}
	}

	return true, nil
}

func (r *BookingEstimationImpl) CopyFromHistory(tx *gorm.DB, batchid int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	var modeldetail transactionworkshopentities.BookingEstimationDetail
	var detailpayloads []transactionworkshoppayloads.BookEstimDetailPayloads

	// Slice to hold the payloads
	var payloads []map[string]interface{}
	result, err2 := r.PostBookingEstimationCalculation(tx, batchid)
	if err2 != nil {
		return nil, err2
	}
	// Query item details
	err := tx.Model(&modeldetail).Where("estimation_system_number = ?", result).Scan(&detailpayloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	// Process item details
	for _, item := range detailpayloads {
		entity := transactionworkshopentities.BookingEstimationDetail{
			EstimationLineID:               item.EstimationLineID,
			EstimationLineCode:             item.EstimationLineCode,
			EstimationSystemNumber:         result,
			BillID:                         item.BillID,
			EstimationLineDiscountApproval: item.EstimationLineDiscountApproval,
			ItemOperationID:                item.ItemOperationID,
			LineTypeID:                     item.LineTypeID,
			PackageID:                      item.PackageID,
			JobTypeID:                      item.JobTypeID,
			FieldActionSystemNumber:        item.FieldActionSystemNumber,
			ApprovalRequestNumber:          item.ApprovalRequestNumber,
			UOMID:                          item.UOMID,
			RequestDescription:             item.RequestDescription,
			FRTQuantity:                    item.FRTQuantity,
			ItemOperationPrice:             item.OperationItemPrice,
			DiscountItemOperationAmount:    item.DiscountItemOperationAmount,
			DiscountItemOperationPercent:   item.DiscountItemOperationPercent,
			DiscountRequestPercent:         item.DiscountRequestPercent,
			DiscountRequestAmount:          item.DiscountRequestAmount,
		}

		if err := tx.Save(&entity).Error; err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		// Add item payload to the payloads slice
		payload := map[string]interface{}{
			"estimation_line_id":   item.EstimationLineID,
			"estimation_line_code": item.EstimationLineCode,
			"item_operation_id":    item.ItemOperationID,
			"line_type_id":         item.LineTypeID,
			"package_id":           item.PackageID,
			"job_type_id":          item.JobTypeID,
			"request_description":  item.RequestDescription,
			"frt_quantity":         item.FRTQuantity,
			"item_price":           item.OperationItemPrice,
			"item_name":            item.RequestDescription,
			"operation_name":       "", // Empty operation_name
		}
		payloads = append(payloads, payload)
	}
	_, err3 := r.PutBookingEstimationCalculation(tx, batchid)
	if err3 != nil {
		return nil, err3
	}
	return payloads, nil
}

// func (r *BookingEstimationImpl) AddPackage(tx *gorm.DB, id int, packId int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
// 	var model masterentities.PackageMasterDetail
// 	var entity transactionworkshopentities.BookingEstimation
// 	var detailpayloads []masterpayloads.PackageMasterDetail
// 	var payloads []map[string]interface{}
// 	var pricename masterpayloads.PriceCodeName
// 	var lastprice float64
// 	var operation masterpayloads.Operation
// 	result, err3 := r.PostBookingEstimationCalculation(tx, id)
// 	if err3 != nil {
// 		return nil, err3
// 	}
// 	err2 := tx.Model(&model).Where("package_id = ?", packId).Scan(&detailpayloads).Error
// 	if err2 != nil {
// 		return nil, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusConflict,
// 			Err:        err2,
// 		}
// 	}
// 	errs := tx.Model(&entity).Where("batch_system_number=?", id).Scan(&entity).Error
// 	if errs != nil {
// 		return nil, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusNotFound,
// 			Err:        errs,
// 		}
// 	}
// 	for _, item := range detailpayloads {
// 		if item.LineTypeId != 9 && item.LineTypeId != 0 {
// 			err := tx.Select("mtr_item_price_list.price_list_amount,mtr_item.item_name,mtr_item.item_code").Table("mtr_item_price_list").
// 				Joins("JOIN mtr_item on mtr_item.item_id=mtr_item_price_list.item_id").
// 				Joins("Join mtr_item_operation on mtr_item.item_id=mtr_item_operation.item_operation_model_mapping_id").
// 				Where("item_operation_id=?", item.ItemOperationId).
// 				Scan(&pricename).Error
// 			if err != nil {
// 				return nil, &exceptions.BaseErrorResponse{
// 					StatusCode: http.StatusBadRequest,
// 					Err:        err,
// 				}
// 			}
// 			lastprice = pricename.Price
// 			operation = masterpayloads.Operation{
// 				OperationCode: pricename.Code,
// 				OperationName: pricename.Name,
// 			}
// 		} else {
// 			err := tx.Select("mtr_labour_selling_price_detail.selling_price").
// 				Table("mtr_labour_selling_price_detail").
// 				Joins("Join mtr_labour_selling_price on mtr_labour_selling_price.labour_selling_price_id = mtr_labour_selling_price_detail.labour_selling_price_id").
// 				Where("mtr_labour_selling_price.brand_id =?", entity.BrandId).
// 				Where("mtr_labour_selling_price_detail.model_id=?", entity.ModelId).
// 				Where("mtr_labour_selling_price.company_id = ?", entity.CompanyId).
// 				Where("mtr_labour_selling_price.effective_date < ?", time.Now()).
// 				Scan(&lastprice).Error
// 			if err != nil {
// 				return nil, &exceptions.BaseErrorResponse{
// 					StatusCode: http.StatusBadRequest,
// 					Err:        err,
// 				}
// 			}
// 			err2 := tx.Select("mtr_operation_code.operation_code,mtr_operation_code.operation_name").Table("mtr_operation_code").
// 				Joins("mtr_item_operation on mtr_item_operation.item_operation_model_mapping_id=mtr_operation_model_mapping.operation_model_mapping_id").
// 				Joins("mtr_operation_code on mtr_operation_model_mapping.operation_code_id=mtr_operation_code.operation_code_id").Where("mtr_item_operation.item_operation_id=?", item.ItemOperationId).
// 				Scan(&operation).Error
// 			if err2 != nil {
// 				return nil, &exceptions.BaseErrorResponse{
// 					StatusCode: http.StatusBadRequest,
// 					Err:        err2,
// 				}
// 			}
// 		}
// 		// price,err:= repo.GetOprItemPrice(tx,req.LineTypeID,entity.CompanyId,req.ItemOperationID,entity.BrandId,entity.ModelId,req.JobTypeID,entity.VariantId)
// 		entity := transactionworkshopentities.BookingEstimationDetail{
// 			EstimationSystemNumber: result,
// 			ItemOperationID:        item.ItemOperationId,
// 			LineTypeID:             item.LineTypeId,
// 			PackageID:              item.PackageId,
// 			RequestDescription:     operation.OperationName,
// 			FRTQuantity:            item.FrtQuantity,
// 			ItemOperationPrice:     lastprice,
// 		}
// 		err := tx.Save(&entity).Error
// 		if err != nil {
// 			return nil, &exceptions.BaseErrorResponse{
// 				StatusCode: http.StatusBadRequest,
// 				Err:        err,
// 			}
// 		}
// 		payload := map[string]interface{}{
// 			"item_id":      item.ItemOperationId,
// 			"line_type_id": item.LineTypeId,
// 			"package_id":   item.PackageId,
// 			"item_price":   float64(item.PackageId),
// 		}
// 		payloads = append(payloads, payload)
// 	}
// 	return payloads, nil
// }

func (r *BookingEstimationImpl) InputDiscount(tx *gorm.DB, id int, req transactionworkshoppayloads.BookEstimationPayloadsDiscount) (int, *exceptions.BaseErrorResponse) {
	itemDetails := []struct {
		LineTypeID int
		Value      int
	}{
		{1, req.Accessories},
		{3, req.Material},
		{4, req.Oil},
		{7, req.Souvenir},
		{8, req.Sparepart},
		{9, req.Fee},
		{5, req.Operation},
		{6, req.PackageDiscount},
	}

	for _, detail := range itemDetails {
		err := tx.Model(&transactionworkshopentities.BookingEstimationDetail{}).
			Where("estimation_system_number = ? AND line_type_id = ?", id, detail.LineTypeID).
			Update("discount_item_percent", detail.Value).Error
		if err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
	}
	return id, nil
}

func (r *BookingEstimationImpl) AddFieldAction(tx *gorm.DB, id int, idrecall int) (int, *exceptions.BaseErrorResponse) {
	var modelitem []transactionworkshopentities.ContractServiceDetail
	err2 := tx.Model(&modelitem).Where("contract_service_system_number = ?", idrecall).Scan(&modelitem).Error
	if err2 != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err2,
		}
	}
	for _, item := range modelitem {
		entity := transactionworkshopentities.BookingEstimationDetail{
			EstimationSystemNumber: id,
			ItemOperationID:        item.ItemOperationId,
			LineTypeID:             item.LineTypeId,
			PackageID:              item.PackageId,
			RequestDescription:     item.Description,
			FRTQuantity:            item.FrtQuantity,
			ItemOperationPrice:     item.ItemPrice,
		}

		if err := tx.Save(&entity).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	return id, nil
}

func (r *BookingEstimationImpl) GetByIdBookEstimDetail(tx *gorm.DB, id int, LineTypeID int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	var payloadsoperation []transactionworkshoppayloads.BookEstimDetailPayloadsOperation
	var payloadsitem []transactionworkshoppayloads.BookEstimDetailPayloadsItem
	var payloadslinetype []masterpayloads.LineTypeCode
	var payloadstransactiontype []transactionworkshoppayloads.TransactionTypePayloads

	if LineTypeID == 5 {
		err := tx.Select("trx_booking_estimation_operation_detail.*,mtr_operation_code.operation_code").Where("estimation_line_id = ?", id).
			Joins("join mtr_operation_model_mapping on mtr_operation_model_mapping.operation_model_mapping_id = trx_booking_estimation_operation_detail.operation_id").
			Joins("join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id").Scan(&payloadsoperation).Error
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		errurllinetype := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/line-type", &payloadslinetype, nil)
		if errurllinetype != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errurllinetype,
			}
		}
		joinedData1, errdf := utils.DataFrameInnerJoin(payloadsoperation, payloadslinetype, "LineTypeId")
		if errdf != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errdf,
			}
		}
		errurltransactiontype := utils.Get(config.EnvConfigs.GeneralServiceUrl+"transaction-type-list?page=0&limit=100000", &payloadstransactiontype, nil)
		if errurltransactiontype != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errurltransactiontype,
			}
		}
		joineddata2, errdf := utils.DataFrameInnerJoin(joinedData1, payloadstransactiontype, "TransactionTypeId")
		if errdf != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errdf,
			}
		}
		return joineddata2[0], nil

	} else {
		err := tx.Model("tr_booking_estimation_item_detail.*,mtr_item.item_name").Where("estimation_line_id = ?", id).
			Joins("Join mtr_item on mtr_item.item_id = trx_booking_estimation_item_detail.item_id").
			Scan(&payloadsitem).Error
		if err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		errurllinetype := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/line-type", &payloadslinetype, nil)
		if errurllinetype != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errurllinetype,
			}
		}
		joinedData1, errdf := utils.DataFrameInnerJoin(payloadsitem, payloadslinetype, "LineTypeId")
		if errdf != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errdf,
			}
		}
		errurltransactiontype := utils.Get(config.EnvConfigs.GeneralServiceUrl+"transaction-type-list?page=0&limit=100000", &payloadstransactiontype, nil)
		if errurltransactiontype != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errurltransactiontype,
			}
		}
		joineddata2, errdf := utils.DataFrameInnerJoin(joinedData1, payloadstransactiontype, "TransactionTypeId")
		if errdf != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errdf,
			}
		}
		return joineddata2[0], nil
	}
}

func (r *BookingEstimationImpl) GetAllBookEstimDetail(tx *gorm.DB, id int, pages pagination.Pagination) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	var operationpayloads []transactionworkshoppayloads.BookEstimDetailPayloadsOperation
	var itempayloads []transactionworkshoppayloads.BookEstimDetailPayloadsItem
	var payloadslinetype []masterpayloads.LineTypeCode
	var payloadstransactiontype []transactionworkshoppayloads.TransactionTypePayloads
	combinedpayloads := make([]map[string]interface{}, 0)

	err := tx.Select("trx_booking_estimation_operation_detail.*,mtr_operation_code.operation_code").Where("estimation_line_id = ?", id).
		Joins("join mtr_operation_model_mapping on mtr_operation_model_mapping.operation_model_mapping_id = trx_booking_estimation_operation_detail.operation_id").
		Joins("join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id").Scan(&operationpayloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	errurllinetype := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/line-type", &payloadslinetype, nil)
	if errurllinetype != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errurllinetype,
		}
	}
	joinedData1, errdf := utils.DataFrameInnerJoin(operationpayloads, payloadslinetype, "LineTypeId")
	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}
	errurltransactiontype := utils.Get(config.EnvConfigs.GeneralServiceUrl+"transaction-type-list?page=0&limit=100000", &payloadstransactiontype, nil)
	if errurltransactiontype != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errurltransactiontype,
		}
	}
	joineddata2, errdf := utils.DataFrameInnerJoin(joinedData1, payloadstransactiontype, "TransactionTypeId")
	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}
	for _, op := range joineddata2 {
		combinedpayloads = append(combinedpayloads, map[string]interface{}{
			"line_type_id":        op["LineTypeid"],
			"transaction_type_id": op["TransactionTypeId"],
			"operation_id":        op["OperationId"],
			"operation_name":      op["OperationName"],
			"quantity":            op["Quantity"],
			"price":               op["Price"],
			"subtotal":            op["SubTotal"],
			"original_discount":   op["OriginalDiscount"],
			"proposal_discount":   op["ProposalDiscount"],
			"total":               op["Total"],
		})
	}

	err2 := tx.Model("tr_booking_estimation_item_detail.*,mtr_item.item_name").Where("estimation_line_id = ?", id).
		Joins("Join mtr_item on mtr_item.item_id = trx_booking_estimation_item_detail.item_id").
		Scan(&itempayloads).Error
	if err2 != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	errurllinetype2 := utils.Get(config.EnvConfigs.GeneralServiceUrl+"/line-type", &payloadslinetype, nil)
	if errurllinetype2 != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errurllinetype2,
		}
	}
	joinedData1_2, errdf := utils.DataFrameInnerJoin(itempayloads, payloadslinetype, "LineTypeId")
	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}
	errurltransactiontype2 := utils.Get(config.EnvConfigs.GeneralServiceUrl+"transaction-type-list?page=0&limit=100000", &payloadstransactiontype, nil)
	if errurltransactiontype2 != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errurltransactiontype,
		}
	}
	joineddata2_2, errdf := utils.DataFrameInnerJoin(joinedData1_2, payloadstransactiontype, "TransactionTypeId")
	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}
	for _, it := range joineddata2_2 {
		combinedpayloads = append(combinedpayloads, map[string]interface{}{
			"line_type_id":        it["LineTypeID"],
			"transaction_type_id": it["TransactionTypeId"],
			"item_id":             it["ItemID"],
			"item_name":           it["ItemName"],
			"quantity":            it["Quantity"],
			"price":               it["Price"],
			"subtotal":            it["SubTotal"],
			"original_discount":   it["OriginalDiscount"],
			"proposal_discount":   it["ProposalDiscount"],
			"total":               it["Total"],
		})
	}
	return combinedpayloads, nil
}

func (r *BookingEstimationImpl) PostBookingEstimationCalculation(tx *gorm.DB, id int) (int, *exceptions.BaseErrorResponse) {
	now := time.Now()
	var estimationSystemNumber int
	err2 := tx.Select("trx_booking_estimation_service_discount.estimation_system_number").Table("trx_booking_estimation_service_discount").Where("batch_system_number=?", id).Scan(&estimationSystemNumber).Error
	if err2 != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err2,
		}
	}
	if estimationSystemNumber != 0 {
		return estimationSystemNumber, nil
	}
	entity := &transactionworkshopentities.BookingEstimationServiceDiscount{
		BatchSystemNumber:    id,
		EstimationDate:       &now,
		DiscountApprovalDate: &now,
	}
	err := tx.Save(entity).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entity.EstimationSystemNumber, nil
}

func (r *BookingEstimationImpl) PutBookingEstimationCalculation(tx *gorm.DB, id int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {

	currentTime := time.Now()
	taxResponse, errTax := financeserviceapiutils.GetTaxPercent("PPN", "PPN", currentTime)
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
	var firststruct transactionworkshoppayloads.BookingEstimationFirstContractService
	var contractservice transactionworkshoppayloads.ContractService
	var taxfare float64
	now := time.Now()
	var count int64
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
		Select(`BE.contract_system_number, BE2.estimation_discount_approval_status, BE.booking_system_number, 
			BE.estimation_system_number, BE.brand_id, BE.profit_center_id, BE.model_id, BE.company_id`).
		Joins(`LEFT JOIN trx_booking_estimation_service_discount BE2 ON BE.batch_system_number = BE2.batch_system_number`).
		Where(`BE.batch_system_number = ?`, idheader).
		Scan(&firststruct).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	errUrlGetTax := utils.Get(config.EnvConfigs.FinanceServiceUrl+"tax-fare/detail/tax-percent?tax_service_code=PPN&tax_type_code=PPN&effective_date="+time.Now().String(), taxfare, nil)
	if errUrlGetTax != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errUrlGetTax,
		}
	}
	entities := transactionworkshopentities.BookingEstimationServiceDiscount{
		DocumentStatusID:                 10, //status new
		BatchSystemNumber:                idheader,
		EstimationDate:                   &now,
		EstimationDiscountApprovalStatus: 10, //status draft
		CompanyID:                        firststruct.CompanyId,
		VATTaxRate:                       taxfare,
	}
	err3 := tx.Save(&entities).Error
	if err3 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err3,
		}
	}

	err4 := tx.Select("trx_booking_estimation_detail.estimation_line_code").Table("trx_booking_estimation_detail").Where("estimation_system_number=?", entities.EstimationSystemNumber).Count(&count).Error
	if err4 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err4,
		}
	}
	err5 := tx.Select("CSD.line_type_id,CSD.item_operation_id,CSD.description,CSD.frt_quantity,CSD.item_price,CSD.item_discount_percent").
		Table("trx_contract_service CS").
		Joins("Join trx_contract_service_detail CSD on CSD.contract_service_system_number = CS.contract_service_system_number").
		Where("CS.contract_service_system_number=?", Idcontract).Scan(&contractservice).Error
	if err5 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err5,
		}
	}

	entities2 := transactionworkshopentities.BookingEstimationDetail{
		EstimationSystemNumber:         entities.EstimationSystemNumber,
		EstimationLineCode:             int(count) + 1,
		LineTypeID:                     contractservice.LineTypeId,
		BillID:                         1,  //id of TRXTYPE_WO_CONTRACT_SERVICE
		EstimationLineDiscountApproval: 10, //id of approval status draft
		ItemOperationID:                contractservice.ItemOperationId,
		RequestDescription:             contractservice.Description,
		FRTQuantity:                    float64(contractservice.FrtQuantity),
		ItemOperationPrice:             contractservice.ItemPrice,
		DiscountItemOperationAmount:    math.Round(contractservice.ItemPrice * contractservice.ItemDiscountPercent / 100),
		DiscountRequestAmount:          0,
		DiscountRequestPercent:         0,
		DiscountItemOperationPercent:   contractservice.ItemDiscountPercent,
		DiscountApprovalBy:             "",
		DiscountApprovalDate:           nil,
	}
	err6 := tx.Save(&entities2).Error
	if err6 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err6,
		}
	}
	_, err7 := r.PutBookingEstimationCalculation(tx, idheader)
	if err7 != nil {
		return false, err7
	}
	if contractservice.LineTypeId == 0 {
		entities3 := transactionworkshopentities.BookingEstimationAllocation{
			BookingStall:       "",
			BookingServiceTime: float32(entities2.FRTQuantity),
			BookingServiceDate: entities.EstimationDate,
		}
		err := tx.Where("booking_system_number = ?", entities3.BookingSystemNumber).
			Updates(&entities3).Error

		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}

	}
	return true, nil
}

func (r *BookingEstimationImpl) AddPackage(tx *gorm.DB, idhead int, idpackage int) (bool, *exceptions.BaseErrorResponse) {
	var headerdata transactionworkshoppayloads.BookingEstimationFirstContractService
	var taxfare float64
	var count int64
	var uom int
	var price float64
	var discpercent float64
	var data transactionworkshoppayloads.PackageForDetail
	time := time.Now()

	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
		Select(`BE.contract_system_number, BE2.estimation_discount_approval_status, BE.booking_system_number, 
			BE.estimation_system_number, BE.brand_id, BE.profit_center_id, BE.model_id, BE.company_id`).
		Joins(`LEFT JOIN trx_booking_estimation_service_discount BE2 ON BE.batch_system_number = BE2.batch_system_number`).
		Where(`BE.batch_system_number = ?`, idhead).
		Scan(&headerdata).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	errUrlGetTax := utils.Get(config.EnvConfigs.FinanceServiceUrl+"tax-fare/detail/tax-percent?tax_service_code=PPN&tax_type_code=PPN&effective_date="+time.String(), taxfare, nil)
	if errUrlGetTax != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errUrlGetTax,
		}
	}

	entities := transactionworkshopentities.BookingEstimationServiceDiscount{
		EstimationDocumentNumber:         headerdata.EstimationDocumentNumber,
		EstimationDate:                   &time,
		DocumentStatusID:                 5,  //status new
		EstimationDiscountApprovalStatus: 10, //approval draft
		BatchSystemNumber:                idhead,
		CompanyID:                        headerdata.CompanyId,
		VATTaxRate:                       taxfare,
	}
	err3 := tx.Save(&entities).Error
	if err3 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err3,
		}
	}

	err4 := tx.Table("trx_booking_estimation_detail AS tb").
		Select("tb.estimation_line_id").
		Where("tb.estimation_system_number = ?", entities.EstimationSystemNumber).
		Count(&count).Error
	if err4 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err4,
		}
	}
	err5 := tx.Select(`
    p1.item_operation_id, 
    p1.line_type_id, 
    p1.frt_quantity, 
    p0.currency_id, 
    p1.job_type_id, 
    p1.workorder_transaction_type_id, 
    CASE 
        WHEN p1.line_type_id = 9 THEN op.operation_name 
        WHEN p1.line_type_id = 0 THEN p0.package_name 
        ELSE it.item_name 
    END AS item_or_operation_name
`).
		Table("mtr_package p0").
		Joins("JOIN mtr_package_master_detail p1 ON p0.package_id = p1.package_id").
		Joins("LEFT JOIN mtr_item_operation io ON io.item_operation_id = p1.item_operation_id").
		Joins(`LEFT JOIN mtr_operation_model_mapping opr 
		ON opr.brand_id = p0.brand_id 
		AND opr.model_id = p0.model_id 
		AND opr.operation_id = CASE WHEN p1.line_type_id = 9 THEN io.item_operation_model_mapping_id ELSE NULL END`).
		Joins(`LEFT JOIN mtr_operation_code op 
		ON op.operation_id = CASE WHEN p1.line_type_id = 9 THEN io.item_operation_model_mapping_id ELSE NULL END`).
		Joins(`LEFT JOIN mtr_item it 
		ON it.item_id = CASE WHEN p1.line_type_id != 9 THEN io.item_operation_model_mapping_id ELSE NULL END`).
		Where("p0.package_id = ?", idpackage).
		Scan(&data).Error
	if err5 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err5,
		}
	}
	err6 := tx.Select("mtr_item.unit_of_measurement_type_id").Table("mtr_item").Where("mtr_item.item_name=?", data.ItemOrOperationName).Scan(&uom).Error
	if err6 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err6,
		}
	}
	err7 := tx.Select("mtr_price_list.price_list_amount").Table("mtr_price_list").
		Joins("join mtr_item on mtr_price_list.price_list_id = mtr_item.price_list_item").
		Joins("join mtr_item_operation on mtr_item_operation.item_operation_model_mapping_id=mtr_item.item_id").Where("mtr_item_operation.item_operation_id=?", data.ItemOperationId).
		Scan(&price).Error
	if err7 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err7,
		}
	}
	err8 := tx.Select("trx_contract_service_detail.item_discount_percent").
		Table("trx_contract_service_detail").
		Joins("Join trx_contract_service on trx_contract_service.contract_service_system_number=trx_contract_service_detail.contract_service_system_number").
		Where("trx_contract_service.contract_service_system_number=0").
		Where("trx_contract_service_detail.item_operation_id=?", data.ItemOperationId).
		Where("trx_contract_service_detail.line_type_id=?", data.LineTypeId).Scan(&discpercent).Error
	if err8 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err8,
		}
	}
	operationdiscount := math.Round(price * data.FrtQuantity * discpercent / 100)
	entities2 := transactionworkshopentities.BookingEstimationDetail{
		EstimationSystemNumber:         entities.EstimationSystemNumber,
		EstimationLineCode:             int(count + 1),
		LineTypeID:                     data.LineTypeId,
		JobTypeID:                      data.JobTypeId,
		BillID:                         data.BillId,
		EstimationLineDiscountApproval: 10, //approval draft id
		ItemOperationID:                data.ItemOperationId,
		RequestDescription:             data.ItemOrOperationName,
		PackageID:                      idpackage,
		UOMID:                          uom,
		FRTQuantity:                    data.FrtQuantity,
		ItemOperationPrice:             price,
		DiscountItemOperationAmount:    operationdiscount,
		DiscountItemOperationPercent:   discpercent,
	}
	err9 := tx.Save(&entities2).Error
	if err9 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err9,
		}
	}
	_, err0 := r.PutBookingEstimationCalculation(tx, idhead)
	if err0 != nil {
		return false, err0
	}
	return true, nil
}

func (r *BookingEstimationImpl) SaveBookingEstimationFromServiceRequest(tx *gorm.DB, idservreq int, req transactionworkshoppayloads.PdiServiceRequest) (bool, *exceptions.BaseErrorResponse) {
	var initialpayloads *transactionworkshoppayloads.ServiceRequestBookingEstimation
	var vehiclepayloads transactionworkshoppayloads.VehicleTnkb
	var lastprice float64
	var linetype int
	var approvalstatus transactionunitpayloads.ApprovalStatus
	var documentstatus transactionworkshoppayloads.DocumentStatus
	var workordertransaction transactionworkshoppayloads.WorkorderTransactionType
	var servicerequestdetail []transactionworkshoppayloads.ServiceRequestDetailBookingPayloads
	time := time.Now()
	formattedTime := time.Format("2006-01-02 15:04:05")

	// Build the query
	err := tx.Select("trx_service_request.profit_center_id, trx_service_request.company_id, trx_service_request.vehicle_id, trx_service_request.service_request_document_number, trx_contract_service.contract_service_system_number").Table("trx_service_request").
		Joins("JOIN trx_contract_service ON trx_contract_service.vehicle_id = trx_service_request.vehicle_id AND trx_contract_service.contract_service_to < ? AND ? > trx_contract_service.contract_service_from AND trx_contract_service.contract_service_status_id = ?", formattedTime, formattedTime, 20).
		Where("trx_service_request.service_request_system_number = ?", idservreq).
		Scan(&initialpayloads).Error

	if err != nil || initialpayloads == nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}
	errUrlVehicle := utils.Get(config.EnvConfigs.SalesServiceUrl+"vehicle-master/"+strconv.Itoa(initialpayloads.VehicleId), vehiclepayloads, nil)
	if errUrlVehicle != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlVehicle,
		}
	}
	entity := transactionworkshopentities.BookingEstimation{
		BrandId:                        vehiclepayloads.VehicleBrandId,
		ModelId:                        vehiclepayloads.VehicleBrandId,
		VariantId:                      vehiclepayloads.VehicleVariantId,
		VehicleId:                      initialpayloads.VehicleId,
		ContractSystemNumber:           initialpayloads.ContractServiceSystemNumber,
		CompanyId:                      initialpayloads.CompanyId,
		BookingSystemNumber:            0,
		ServiceRequestSystemNumber:     0,
		EstimationSystemNumber:         0,
		AgreementNumberBr:              "",
		AgreementId:                    0,
		ContactPersonName:              req.ContactPersonName,
		ContactPersonPhone:             req.ContactPersonPhone,
		ContactPersonViaId:             req.ContactPersonViaId,
		ContactPersonMobile:            req.ContactPersonMobile,
		InsurancePolicyNo:              "",
		InsuranceExpiredDate:           time,
		InsuranceClaimNo:               "",
		InsurancePic:                   "",
		ProfitCenterId:                 initialpayloads.ProfitCenterId,
		IsUnregistered:                 false,
		BookingEstimationBatchDate:     time,
		BookingEstimationVehicleNumber: vehiclepayloads.Tnkb,
	}
	err1 := tx.Save(entity).Error
	if err1 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err1,
		}
	}
	errUrlLineType := utils.Get(config.EnvConfigs.GeneralServiceUrl+"line-type-by-name/operation", &linetype, nil)
	if errUrlLineType != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errUrlLineType,
		}
	}
	errUrlApprovalStatus := utils.Get(config.EnvConfigs.GeneralServiceUrl+"approval-status-description/draft", &approvalstatus, nil)
	if errUrlApprovalStatus != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errUrlLineType,
		}
	}
	errUrlWorkorderTransactionType := utils.Get(config.EnvConfigs.GeneralServiceUrl+"work-order-transaction-type-by-code/external", &workordertransaction, nil)
	if errUrlWorkorderTransactionType != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errUrlLineType,
		}
	}
	errUrlDocumentStatus := utils.Get(config.EnvConfigs.GeneralServiceUrl+"document-status-by-description/New%20Document", &documentstatus, nil)
	if errUrlDocumentStatus != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errUrlDocumentStatus,
		}
	}
	entities8 := transactionworkshopentities.BookingEstimationAllocation{
		BookingStatusID:       documentstatus.DocumentStatusId, //document status new
		BatchSystemNumber:     entity.BatchSystemNumber,
		CompanyID:             initialpayloads.CompanyId,
		PdiSystemNumber:       idservreq,
		BookingDocumentNumber: initialpayloads.ServiceRequestDocumentNumber,
		BookingDate:           nil,
		BookingStall:          " ",
		BookingReminderDate:   nil,
		BookingServiceDate:    nil,
		BookingServiceTime:    0,
		BookingEstimationTime: 0,
	}
	err8 := tx.Save(&entities8).Error
	if err8 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err8,
		}
	}
	approvalstatusid, _ := strconv.Atoi(approvalstatus.ApprovalStatusId)
	entities := transactionworkshopentities.BookingEstimationServiceDiscount{
		BatchSystemNumber:                entity.BatchSystemNumber,
		DocumentStatusID:                 documentstatus.DocumentStatusId,
		EstimationDiscountApprovalStatus: approvalstatusid,
		CompanyID:                        entity.CompanyId,
		ApprovalRequestNumber:            0,
		EstimationDate:                   &time,
		TotalPricePackage:                0.0,
		TotalPriceOperation:              0.0,
		TotalPricePart:                   0.0,
		TotalPriceOil:                    0.0,
		TotalPriceMaterial:               0.0,
		TotalPriceConsumableMaterial:     0.0,
		TotalSublet:                      0.0,
		TotalPriceAccessories:            0.0,
		TotalDiscount:                    0.0,
		TotalVAT:                         0.0,
		TotalAfterVAT:                    0.0,
		AdditionalDiscountRequestPercent: 0.0,
		AdditionalDiscountRequestAmount:  0.0,
		VATTaxRate:                       0.0,
		DiscountApprovalBy:               "",
		DiscountApprovalDate:             &time,
		TotalAfterDiscount:               0.0,
	}
	err2 := tx.Save(entities).Error
	if err2 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err2,
		}
	}

	err4 := tx.Model(transactionworkshopentities.ServiceRequestDetail{}).Where("service_request_system_number = ?", idservreq).Scan(&servicerequestdetail).Error
	if err4 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err4,
		}
	}
	err3 := tx.Select("mtr_labour_selling_price_detail.selling_price").
		Table("mtr_labour_selling_price_detail").
		Joins("Join mtr_labour_selling_price on mtr_labour_selling_price.labour_selling_price_id = mtr_labour_selling_price_detail.labour_selling_price_id").
		Where("mtr_labour_selling_price.brand_id =?", vehiclepayloads.VehicleBrandId).
		Where("mtr_labour_selling_price.company_id = ?", initialpayloads.CompanyId).
		Where("mtr_labour_selling_price.effective_date < ?", time).
		Scan(&lastprice).Error
	if err3 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err3,
		}
	}
	approvalstatusid1, _ := strconv.Atoi(approvalstatus.ApprovalStatusId)
	for _, detail := range servicerequestdetail {
		entities3 := transactionworkshopentities.BookingEstimationDetail{
			EstimationSystemNumber:         entities.EstimationSystemNumber,
			BillID:                         workordertransaction.WorkOrderTransactionTypeId, //transaction type workorder external
			EstimationLineDiscountApproval: approvalstatusid1,                               //status draft
			ItemOperationID:                detail.OperationItemId,
			LineTypeID:                     detail.LineTypeId, //line type id where line type description = operation
			RequestDescription:             "",
			FRTQuantity:                    detail.FrtQuantity,
			ItemOperationPrice:             lastprice,
			DiscountItemOperationAmount:    0,
			DiscountItemOperationPercent:   0,
			DiscountRequestPercent:         0,
			DiscountRequestAmount:          0,
		}
		err4 := tx.Save(entities3).Error
		if err4 != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err4,
			}
		}
		_, errs := r.PutBookingEstimationCalculation(tx, entity.BatchSystemNumber)
		if errs != nil {
			return false, errs
		}
	}
	return true, nil
}

func (r *BookingEstimationImpl) SaveBookingEstimationFromPDI(tx *gorm.DB, idpdi int, req transactionworkshoppayloads.PdiServiceRequest) (bool, *exceptions.BaseErrorResponse) {
	var pdipayload transactionunitpayloads.PdiRequest
	var pdidetailpayloads []transactionunitpayloads.PdiRequestDetail
	var pdidetailbyid []transactionunitpayloads.PdiRequestDetailById
	var agreement masterpayloads.AgreementResponse
	var lastprice float64
	var linetype masterpayloads.LineTypeCode
	var vehicle transactionworkshoppayloads.VehicleTnkb
	var workordertransaction transactionworkshoppayloads.WorkorderTransactionType
	var profitcenter transactionunitpayloads.ProfitCenterResponse
	var approvalstatus []transactionunitpayloads.ApprovalStatus
	var contractservice transactionunitpayloads.ContractService
	errUrlPdiRequest := utils.Get(config.EnvConfigs.SalesServiceUrl+"pdi-request/"+strconv.Itoa(idpdi), &pdipayload, nil)
	if errUrlPdiRequest != nil || pdipayload.CompanyID == 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlPdiRequest,
		}
	}
	errUrlPdiDetailrequest := utils.Get(config.EnvConfigs.SalesServiceUrl+"pdi-request-full/"+strconv.Itoa(idpdi)+"?page=0&limit=1000000", &pdidetailpayloads, nil)
	if errUrlPdiDetailrequest != nil || len(pdidetailpayloads) == 0 {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlPdiDetailrequest,
		}
	}

	for _, detail := range pdidetailpayloads {
		var tempDetail []transactionunitpayloads.PdiRequestDetailById // This is a slice to handle the array response
		errUrlPdiDetail := utils.Get(config.EnvConfigs.SalesServiceUrl+"pdi-request-by-detail-id/"+strconv.Itoa(detail.PdiRequestDetailSystemNumber), &tempDetail, nil)
		if errUrlPdiDetail != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errUrlPdiDetail,
			}
		}

		pdidetailbyid = append(pdidetailbyid, tempDetail...) // Append the slice to the existing slice
	}
	profitcenterurl := config.EnvConfigs.GeneralServiceUrl + "profit-center-by-name/Workshop"
	errUrlProfitCenter := utils.Get(profitcenterurl, &profitcenter, nil)
	if errUrlProfitCenter != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlProfitCenter,
		}
	}
	erragreementcompare := tx.Select("mtr_agreement.*").Table("mtr_agreement").
		Where("mtr_agreement.customer_id =?", pdipayload.CompanyID).
		Where("mtr_agreement.profit_center_id=?", profitcenter.ProfitCenterId).
		Where("? between mtr_agreement.agreement_date_from and mtr_agreement.agreement_date_to", time.Now()).
		Scan(&agreement).Error
	if erragreementcompare != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        erragreementcompare,
		}
	}
	urlApprovalStatus1 := config.EnvConfigs.GeneralServiceUrl + "approval-status-description/Ready"
	errapprovalstatuscontractservice := utils.Get(urlApprovalStatus1, &approvalstatus, nil)
	if errapprovalstatuscontractservice != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errapprovalstatuscontractservice,
		}
	}
	errcontractservice := tx.Select("trx_contract_service.*").Table("trx_contract_service").
		Where("trx_contract_service.contract_service_status_id = ?", approvalstatus[0].ApprovalStatusId).
		Where("trx_contract_service.vehicle_id = ?", pdidetailbyid[0].VehicleId).
		Where("trx_contract_service.contract_service_from < ? AND trx_contract_service.contract_service_to > ?", time.Now(), time.Now()).
		Scan(&contractservice).Error

	if errcontractservice != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errcontractservice,
		}
	}
	errUrlVehicle := utils.Get(config.EnvConfigs.SalesServiceUrl+"vehicle-master/"+strconv.Itoa(pdidetailbyid[0].VehicleId), &vehicle, nil)
	if errUrlVehicle != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlVehicle,
		}
	}
	entities := transactionworkshopentities.BookingEstimation{
		BrandId:                        pdipayload.BrandID,
		ModelId:                        pdidetailbyid[0].ModelId,
		VariantId:                      pdidetailbyid[0].VariantId,
		VehicleId:                      pdidetailbyid[0].VehicleId,
		ContractSystemNumber:           contractservice.ContractServiceId,
		CompanyId:                      pdipayload.CompanyID,
		BookingSystemNumber:            0,
		ServiceRequestSystemNumber:     0,
		EstimationSystemNumber:         0,
		AgreementNumberBr:              "",
		AgreementId:                    0,
		ContactPersonName:              req.ContactPersonName,
		ContactPersonPhone:             req.ContactPersonPhone,
		ContactPersonViaId:             req.ContactPersonViaId,
		ContactPersonMobile:            req.ContactPersonMobile,
		InsurancePolicyNo:              "",
		InsuranceExpiredDate:           time.Time{},
		InsuranceClaimNo:               "",
		InsurancePic:                   "",
		ProfitCenterId:                 profitcenter.ProfitCenterId,
		IsUnregistered:                 false,
		BookingEstimationBatchDate:     time.Now(),
		BookingEstimationVehicleNumber: vehicle.Tnkb,
	}
	err := tx.Save(&entities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	entities8 := transactionworkshopentities.BookingEstimationAllocation{
		BookingStatusID:       1, //document status new
		BatchSystemNumber:     entities.BatchSystemNumber,
		CompanyID:             pdipayload.CompanyID,
		PdiSystemNumber:       idpdi,
		BookingDocumentNumber: pdipayload.PdiDocumentNumber,
		BookingDate:           nil,
		BookingStall:          " ",
		BookingReminderDate:   nil,
		BookingServiceDate:    nil,
		BookingServiceTime:    0,
		BookingEstimationTime: 0,
	}
	err8 := tx.Save(&entities8).Error
	if err8 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err8,
		}
	}

	now := time.Now()
	entities2 := transactionworkshopentities.BookingEstimationServiceDiscount{
		BatchSystemNumber:                entities.BatchSystemNumber,
		DocumentStatusID:                 10,
		EstimationDiscountApprovalStatus: 10,
		CompanyID:                        entities.CompanyId,
		ApprovalRequestNumber:            0,
		EstimationDate:                   &now,
		TotalPricePackage:                0.0,
		TotalPriceOperation:              0.0,
		TotalPricePart:                   0.0,
		TotalPriceOil:                    0.0,
		TotalPriceMaterial:               0.0,
		TotalPriceConsumableMaterial:     0.0,
		TotalSublet:                      0.0,
		TotalPriceAccessories:            0.0,
		TotalDiscount:                    0.0,
		TotalVAT:                         0.0,
		TotalAfterVAT:                    0.0,
		AdditionalDiscountRequestPercent: 0.0,
		AdditionalDiscountRequestAmount:  0.0,
		VATTaxRate:                       0.0,
		DiscountApprovalBy:               "",
		DiscountApprovalDate:             &now,
		TotalAfterDiscount:               0.0,
	}
	err2 := tx.Save(&entities2).Error
	if err2 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err2,
		}
	}
	err3 := tx.Select("mtr_labour_selling_price_detail.selling_price").
		Table("mtr_labour_selling_price_detail").
		Joins("Join mtr_labour_selling_price on mtr_labour_selling_price.labour_selling_price_id = mtr_labour_selling_price_detail.labour_selling_price_id").
		Where("mtr_labour_selling_price.brand_id =?", pdipayload.BrandID).
		Where("mtr_labour_selling_price.company_id = ?", pdipayload.CompanyID).
		Where("mtr_labour_selling_price.effective_date < ?", time.Now()).
		Scan(&lastprice).Error
	if err3 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err3,
		}
	}
	errUrlLineType := utils.Get(config.EnvConfigs.GeneralServiceUrl+"line-type-by-name/Operation", &linetype, nil)
	if errUrlLineType != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errUrlLineType,
		}
	}
	urlApprovalStatus := config.EnvConfigs.GeneralServiceUrl + "approval-status-description/Draft"
	errUrlApprovalStatus := utils.Get(urlApprovalStatus, &approvalstatus, nil)
	if errUrlApprovalStatus != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errUrlApprovalStatus,
		}
	}
	errUrlWorkorderTransactionType := utils.Get(config.EnvConfigs.GeneralServiceUrl+"work-order-transaction-type-by-code/External", &workordertransaction, nil)
	if errUrlWorkorderTransactionType != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        errUrlWorkorderTransactionType,
		}
	}
	approvalstatusid, _ := strconv.Atoi(approvalstatus[0].ApprovalStatusId)
	for _, detail := range pdidetailbyid {
		entities3 := transactionworkshopentities.BookingEstimationDetail{
			EstimationSystemNumber:         entities2.EstimationSystemNumber,
			BillID:                         workordertransaction.WorkOrderTransactionTypeId, //transaction type workorder external
			EstimationLineDiscountApproval: approvalstatusid,                                //status draft
			ItemOperationID:                detail.OperationNumberId,
			LineTypeID:                     linetype.LineTypeId, //line type id where line type description = operation
			RequestDescription:             "",
			FRTQuantity:                    detail.Frt,
			ItemOperationPrice:             lastprice,
			DiscountItemOperationAmount:    0,
			DiscountItemOperationPercent:   0,
			DiscountRequestPercent:         0,
			DiscountRequestAmount:          0,
		}
		err4 := tx.Save(&entities3).Error
		if err4 != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err4,
			}
		}
	}

	return true, nil
}
