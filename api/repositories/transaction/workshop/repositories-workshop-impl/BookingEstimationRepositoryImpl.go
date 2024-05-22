package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"net/http"

	"gorm.io/gorm"
)

type BookingEstimationImpl struct {
}

func OpenBookingEstimationRepositoryImpl() transactionworkshoprepository.BookingEstimationRepository {
	return &BookingEstimationImpl{}
}

func (r *BookingEstimationImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.BookingEstimation
	// Query to retrieve all booking estimation entities based on the request
	query := tx.Model(&transactionworkshopentities.BookingEstimation{})
	if len(filterCondition) > 0 {
		query = query.Where(filterCondition)
	}
	err := query.Find(&entities).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{Message: "Failed to retrieve booking estimations from the database"}
	}

	var bookingEstimationResponses []map[string]interface{}

	// Loop through each entity and copy its data to the response
	for _, entity := range entities {
		bookingEstimationData := make(map[string]interface{})
		// Copy data from entity to response
		bookingEstimationData["batch_system_number"] = entity.BatchSystemNumber
		bookingEstimationData["booking_system_number"] = entity.BookingSystemNumber
		bookingEstimationData["brand_id"] = entity.BrandId
		bookingEstimationData["model_id"] = entity.ModelId
		bookingEstimationData["variant_id"] = entity.VariantId
		bookingEstimationData["vehicle_id"] = entity.VehicleId
		bookingEstimationData["estimation_system_number"] = entity.EstimationSystemNumber
		bookingEstimationData["pdi_system_number"] = entity.PdiSystemNumber
		bookingEstimationData["service_request_system_number"] = entity.ServiceRequestSystemNumber
		bookingEstimationData["contract_system_number"] = entity.ContractSystemNumber
		bookingEstimationData["agreement_id"] = entity.AgreementId
		bookingEstimationData["campaign_id"] = entity.CampaignId
		bookingEstimationData["company_id"] = entity.CompanyId
		bookingEstimationData["profit_center_id"] = entity.ProfitCenterId
		bookingEstimationData["dealer_representative_id"] = entity.DealerRepresentativeId
		bookingEstimationData["customer_id"] = entity.CustomerId
		bookingEstimationData["document_status_id"] = entity.DocumentStatusId
		bookingEstimationData["booking_estimation_batch_date"] = entity.BookingEstimationBatchDate
		bookingEstimationData["booking_estimation_vehicle_number"] = entity.BookingEstimationVehicleNumber
		bookingEstimationData["agreement_number_br"] = entity.AgreementNumberBr
		bookingEstimationData["is_unregistered"] = entity.IsUnregistered
		bookingEstimationData["contact_person_name"] = entity.ContactPersonName
		bookingEstimationData["contact_person_phone"] = entity.ContactPersonPhone
		bookingEstimationData["contact_person_mobile"] = entity.ContactPersonMobile
		bookingEstimationData["contact_person_via"] = entity.ContactPersonVia
		bookingEstimationData["insurance_policy_no"] = entity.InsurancePolicyNo
		bookingEstimationData["insurance_expired_date"] = entity.InsuranceExpiredDate
		bookingEstimationData["insurance_claim_no"] = entity.InsuranceClaimNo
		bookingEstimationData["insurance_pic"] = entity.InsurancePic

		// Append the response data to the array
		bookingEstimationResponses = append(bookingEstimationResponses, bookingEstimationData)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(bookingEstimationResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *BookingEstimationImpl) New(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (*exceptions.BaseErrorResponse, error) {
	// Create a new instance of WorkOrderRepositoryImpl
	// Save the booking estimation
	success, err := r.Save(tx, request) // Menggunakan method Save dari receiver saat ini, yaitu r
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save booking estimation"}, err
	}

	return success, nil
}

func (r *BookingEstimationImpl) GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.BookingEstimationRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.BookingEstimation
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.BookingEstimationRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve booking estimation from the database"}
	}

	// Convert entity to payload
	payload := transactionworkshoppayloads.BookingEstimationRequest{}

	return payload, nil
}

func (r *BookingEstimationImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (*exceptions.BaseErrorResponse, error) {
	var bookingEstimationEntities = transactionworkshopentities.BookingEstimation{
		// Assign fields from request
	}

	// Create a new record
	err := tx.Create(&bookingEstimationEntities).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}, err
	}
	return nil, nil
}

func (r *BookingEstimationImpl) Submit(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	// Retrieve the booking estimation by Id
	var entity transactionworkshopentities.BookingEstimation
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve booking estimation from the database"}
	}

	// Perform the necessary operations to submit the booking estimation
	// ...

	// Save the updated booking estimation
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated booking estimation"}
	}

	return nil
}

func (r *BookingEstimationImpl) Void(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	// Retrieve the booking estimation by Id
	var entity transactionworkshopentities.BookingEstimation
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to retrieve booking estimation from the database"}
	}

	// Perform the necessary operations to void the booking estimation
	// ...

	// Save the updated booking estimation
	err = tx.Save(&entity).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{Message: "Failed to save the updated booking estimation"}
	}

	return nil
}

func (r *BookingEstimationImpl) CloseOrder(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	// Retrieve the booking estimation by Id
	var entity transactionworkshopentities.BookingEstimation
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Where("id = ?", Id).First(&entity).Error
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
