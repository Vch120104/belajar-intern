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
	"errors"
	"fmt"
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
	return &BookingEstimationImpl{
	}
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
		bookingEstimationData["contact_person_via_id"] = entity.ContactPersonViaId
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

func (r *BookingEstimationImpl) Post(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse) {
	var agreement masterentities.Agreement
	var contractService transactionworkshopentities.ContractService
	var id int
	var id2 int
	err := tx.Select("mtr_agreement.agreement_id").Model(agreement).Where("customer_id = ?", strconv.Itoa(request.CustomerId)).Scan(&id).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	err2 := tx.Select("trx_contract_service.contract_service_system_number").Model(contractService).Where("vehicle_id = ?", strconv.Itoa(request.VehicleId)).Scan(&id2).Error
	if err2 != nil && !errors.Is(err2, gorm.ErrRecordNotFound) {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err2,
		}
	}

	entities := transactionworkshopentities.BookingEstimation{
		BrandId:                    request.BrandId,
		ModelId:                    request.ModelId,
		VariantId:                  request.VariantId,
		VehicleId:                  request.VehicleId,
		ContractSystemNumber:       id2,
		AgreementId:                id,
		CampaignId:                 request.CampaignId,
		CompanyId:                  request.CompanyId,
		ProfitCenterId:             request.ProfitCenterId,
		DealerRepresentativeId:     request.DealerRepresentativeId,
		CustomerId:                 request.CustomerId,
		DocumentStatusId:           request.DocumentStatusId,
		BookingEstimationBatchDate: time.Now(),
		IsUnregistered:             request.IsUnregistered,
		InsurancePolicyNo:          request.InsurancePolicyNo,
		InsuranceExpiredDate:       request.InsuranceExpiredDate,
		InsuranceClaimNo:           request.InsuranceClaimNo,
		InsurancePic:               request.InsurancePic,
		ContactPersonName:          request.ContactPersonName,
		ContactPersonPhone:         request.ContactPersonPhone,
		ContactPersonMobile:        request.ContactPersonMobile,
		ContactPersonViaId:         request.ContactPersonViaId,
	}
	err3 := tx.Save(&entities).Error // Menggunakan method Save dari receiver saat ini, yaitu r
	if err != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err3}
	}

	return entities, nil
}

func (r *BookingEstimationImpl) GetById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	var payloads transactionworkshoppayloads.GetBookingById
	var brand masterpayloads.BrandResponse
	var model masterpayloads.GetModelResponse
	var vehicle transactionworkshoppayloads.VehicleDetailPayloads
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).
		Select("trx_booking_estimation.batch_system_number, trx_booking_estimation.brand_id, trx_contract_service.contract_service_system_number, trx_contract_service.contract_service_date, mtr_agreement.agreement_code, mtr_agreement.agreement_date_to, mtr_agreement.company_id").
		Joins("JOIN trx_contract_service on trx_contract_service.vehicle_id = trx_booking_estimation.vehicle_id").
		Joins("JOIN mtr_agreement on mtr_agreement.customer_id = trx_booking_estimation.customer_id").
		Where("batch_system_number = ?", Id).
		First(&payloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{StatusCode: http.StatusNotFound, Err: err}
	}
	errUrlBrand := utils.Get(config.EnvConfigs.SalesServiceUrl+"unit-brand/"+strconv.Itoa(payloads.BrandId), brand, nil)
	if errUrlBrand != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("brand undefined"),
		}
	}
	joinedData1, errdf := utils.DataFrameInnerJoin([]transactionworkshoppayloads.GetBookingById{payloads}, []masterpayloads.BrandResponse{brand}, "BrandId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}
	errUrlModel := utils.Get(config.EnvConfigs.SalesServiceUrl+"unit-model/"+strconv.Itoa(payloads.ModelId), model, nil)
	if errUrlModel != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("model undefined"),
		}
	}
	joinedData2, errdf := utils.DataFrameInnerJoin(joinedData1, []masterpayloads.GetModelResponse{model}, "ModelId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}
	errurlVehicle := utils.Get(config.EnvConfigs.SalesServiceUrl+"vehicle-master/"+strconv.Itoa(payloads.VehicleId), &vehicle, nil)
	if errurlVehicle != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("vehicle undefined"),
		}
	}
	joinedData3, errdf := utils.DataFrameInnerJoin(joinedData2, []transactionworkshoppayloads.VehicleDetailPayloads{vehicle}, "VehicleId")

	if errdf != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errdf,
		}
	}

	return joinedData3[0], nil
}

func (r *BookingEstimationImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptions.BaseErrorResponse) {
	var bookingEstimationEntities = transactionworkshopentities.BookingEstimation{
		BrandId:                request.BrandId,
		ModelId:                request.ModelId,
		VehicleId:              request.VehicleId,
		DealerRepresentativeId: request.DealerRepresentativeId,
	}

	// Create a new record
	err := tx.Create(&bookingEstimationEntities).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return true, nil
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

func (r *BookingEstimationImpl) Void(tx *gorm.DB, Id int) (bool, *exceptions.BaseErrorResponse) {
	// Retrieve the booking estimation by Id
	var entity transactionworkshopentities.BookingEstimation
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Where("batch_system_number = ?", Id).First(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to retrieve booking estimation from the database"}
	}

	// Perform the necessary operations to void the booking estimation
	// ...

	// Save the updated booking estimation
	err = tx.Delete(&entity).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save the updated booking estimation"}
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
	var count int64
	err := tx.Table("trx_booking_estimation_request").Where("booking_system_number =?", id).Count(&count).Error
	if err != nil {
		return transactionworkshopentities.BookingEstimationRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	total := count + 1
	entities := transactionworkshopentities.BookingEstimationRequest{
		BookingEstimationRequestCode: int(total),
		BookingSystemNumber:          id,
		BookingServiceRequest:        req.BookingServiceRequest,
	}
	err2 := tx.Save(&entities).Error
	if err2 != nil {
		return transactionworkshopentities.BookingEstimationRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entities, nil
}

func (r *BookingEstimationImpl) UpdateBookEstimReq(tx *gorm.DB, req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (int, *exceptions.BaseErrorResponse) {
	model := transactionworkshopentities.BookingEstimationRequest{}
	result := tx.Model(&model).Where(transactionworkshopentities.BookingEstimationRequest{BookingEstimationRequestID: id}).First(&model).Updates(req)
	if result.Error != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if model == (transactionworkshopentities.BookingEstimationRequest{}) {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	return id, nil
}

func (r *BookingEstimationImpl) DeleteBookEstimReq(tx *gorm.DB, ids string) ([]string, *exceptions.BaseErrorResponse) {
	// Split the ids string into a slice of individual IDs
	idSlice := strings.Split(ids, ",")

	// Prepare a slice to hold successfully deleted IDs
	deletedIds := []string{}

	// Iterate over each ID and delete the corresponding record
	for _, id := range idSlice {
		model := transactionworkshopentities.BookingEstimationRequest{}

		// Retrieve the record to ensure it exists
		if err := tx.First(&model, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue // If the record is not found, skip to the next ID
			}
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Delete the record
		if err := tx.Delete(&model).Error; err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Append the deleted ID to the deletedIds slice
		deletedIds = append(deletedIds, id)
	}

	// Update line numbers for the remaining records
	if err := updateLineNumbers(tx); err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return deletedIds, nil
}

func updateLineNumbers(tx *gorm.DB) error {
	var records []transactionworkshopentities.BookingEstimationRequest

	// Retrieve all remaining records, ordered by their current line numbers
	if err := tx.Order("booking_estimation_request_code asc").Find(&records).Error; err != nil {
		return err
	}

	// Update line numbers sequentially
	for i, record := range records {
		record.BookingEstimationRequestCode = i + 1
		if err := tx.Save(&record).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *BookingEstimationImpl) GetByIdBookEstimReq(tx *gorm.DB, id int) (transactionworkshoppayloads.BookEstimRemarkRequest, *exceptions.BaseErrorResponse) {
	var model transactionworkshopentities.BookingEstimationRequest
	var payloads transactionworkshoppayloads.BookEstimRemarkRequest

	// Query the database
	err := tx.Model(&model).Where("booking_estimation_request_id = ?", id).Scan(&payloads).Error

	// Check if there was an error during the query
	if err != nil {
		return payloads, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Check if the payload is empty
	if (payloads == transactionworkshoppayloads.BookEstimRemarkRequest{}) {
		return payloads, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("data not found"),
		}
	}

	return payloads, nil
}

func (r *BookingEstimationImpl) GetAllBookEstimReq(tx *gorm.DB, pages *pagination.Pagination, id int) (*pagination.Pagination, *exceptions.BaseErrorResponse) {
	var model []transactionworkshopentities.BookingEstimationRequest

	err := tx.Model(&model).Where("booking_system_number = ?", id).Scan(&model).Scopes(pagination.Paginate(model, pages, tx)).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	pages.Rows = model
	return pages, nil
}

func (r *BookingEstimationImpl) SaveBookEstimReminderServ(tx *gorm.DB, req transactionworkshoppayloads.ReminderServicePost, id int) (int, *exceptions.BaseErrorResponse) {
	var count int64
	err := tx.Select("trx_booking_estimation_request.*").Where("booking_system_number =?", id).Count(&count).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	count32 := int32(count)
	total := count32 + 1
	entities := transactionworkshopentities.BookingEstimationServiceReminder{
		BookingLineNumber:      int(total),
		BookingSystemNumber:    id,
		BookingServiceReminder: req.BookingServiceReminder,
	}
	err2 := tx.Save(&entities).Error
	if err2 != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return entities.BookingSystemNumber, nil
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
		err := tx.Select("mtr_price_list.price_list_amount").Table("mtr_price_list").
			Joins("JOIN mtr_item on mtr_item.item_id=mtr_price_list.item_id").
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
	err := tx.Model(&model).Where("estimation_line_id =?", id).Updates(map[string]interface{}{
		"frt_quantity":             req.FRTQuantity,
		"discount_request_percent": req.DiscountRequestPercent,
	}).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	_, err2 := r.PutBookingEstimationCalculation(tx, model.EstimationSystemNumber)
	if err2 != nil {
		return false, err2
	}

	return true, nil
}

func (r *BookingEstimationImpl) DeleteBookEstimDetail(tx *gorm.DB, id int, linetypeid int) (bool, *exceptions.BaseErrorResponse) {
	var model transactionworkshopentities.BookingEstimationDetail
	err := tx.Delete(&model).Where("estimation_line_id =?", id).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
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
			"item_operation_id":   item.ItemOperationID,
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
	_,err3:= r.PutBookingEstimationCalculation(tx,batchid)
	if err3 != nil{
		return nil,err3
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
// 			err := tx.Select("mtr_price_list.price_list_amount,mtr_item.item_name,mtr_item.item_code").Table("mtr_price_list").
// 				Joins("JOIN mtr_item on mtr_item.item_id=mtr_price_list.item_id").
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
	var taxfare float64
	time := time.Now()
	errUrlGetTax := utils.Get(config.EnvConfigs.FinanceServiceUrl+"tax-fare/detail/tax-percent?tax_service_code=PPN&tax_type_code=PPN&effective_date="+time.String(), taxfare, nil)
	if errUrlGetTax != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errUrlGetTax,
		}
	}
	const (
		LineTypePackage            = 0 // Package Bodyshop
		LineTypeOperation          = 1 // Operation
		LineTypeSparePart          = 2 // Spare Part
		LineTypeOil                = 3 // Oil
		LineTypeMaterial           = 4 // Material
		LineTypeAccessories        = 6 // Accessories
		LineTypeConsumableMaterial = 7 // Consumable Material
		LineTypeSublet             = 8 // Sublet
	)

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

	// Calculate totals for each line type
	err := tx.Raw(`
		SELECT
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(trx_booking_estimation_detail.operation_item_price, 0) * ISNULL(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_price_package,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(trx_booking_estimation_detail.operation_item_price, 0) * ISNULL(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_price_operation,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(trx_booking_estimation_detail.operation_item_price, 0) * ISNULL(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_price_spare_part,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(trx_booking_estimation_detail.operation_item_price, 0) * ISNULL(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_price_oil,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(trx_booking_estimation_detail.operation_item_price, 0) * ISNULL(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_price_material,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(trx_booking_estimation_detail.operation_item_price, 0) * ISNULL(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_price_accessories,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(trx_booking_estimation_detail.operation_item_price, 0) * ISNULL(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_price_consumable_material,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(trx_booking_estimation_detail.operation_item_price, 0) * ISNULL(trx_booking_estimation_detail.frt_quantity, 0), 0) ELSE 0 END) AS total_sublet
		FROM trx_booking_estimation_service_discount
		JOIN trx_booking_estimation_detail on trx_booking_estimation_detail.estimation_system_number = trx_booking_estimation_service_discount.estimation_system_number 
		WHERE batch_system_number = ?`,
		LineTypePackage,
		LineTypeOperation,
		LineTypeSparePart,
		LineTypeOil,
		LineTypeMaterial,
		LineTypeAccessories,
		LineTypeConsumableMaterial,
		LineTypeSublet,
		id).Scan(&result).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to calculate totals: %v", err)}
	}

	// Calculate grand total
	Total := result.TotalPackage + result.TotalOperation + result.TotalSparePart + result.TotalOil + result.TotalMaterial + result.TotalFee + result.TotalAccessories + result.TotalConsumableMaterial + result.TotalSublet + result.TotalSouvenir
	tax := Total * (taxfare) / 100
	grandTotal := Total + tax
	// Update Work Order with the calculated totals
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
		return nil, &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to update work order: %v", err)}
	}

	// Prepare response
	BookingEstimationResponse := []map[string]interface{}{
		{"total_package": result.TotalPackage},
		{"total_operation": result.TotalOperation},
		{"total_part": result.TotalSparePart},
		{"total_oil": result.TotalOil},
		{"total_material": result.TotalMaterial},
		{"total_price_accessories": result.TotalAccessories},
		{"total_consumable_material": result.TotalConsumableMaterial},
		{"total_sublet": result.TotalSublet},
		{"total_after_vat": grandTotal},
	}

	return BookingEstimationResponse, nil
}

func (r *BookingEstimationImpl) SaveBookingEstimationFromPDI(tx *gorm.DB, id int) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse) {
	var pdipayload transactionunitpayloads.PdiRequest
	var pdidetailpayloads transactionunitpayloads.PdiRequestDetail
	var agreement masterpayloads.AgreementResponse
	var agreementdocno string
	var lastprice float64
	var operationtotal float64
	var profitcenter transactionunitpayloads.ProfitCenterResponse
	var approvalstatus transactionunitpayloads.ApprovalStatus
	var contractservice transactionunitpayloads.ContractService
	errUrlPdiRequest := utils.Get(config.EnvConfigs.SalesServiceUrl+"pdi-request/"+strconv.Itoa(id), pdipayload, nil)
	if errUrlPdiRequest != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlPdiRequest,
		}
	}
	errUrlPdiDetailrequest := utils.Get(config.EnvConfigs.SalesServiceUrl+"pdi-request-full/"+strconv.Itoa(id), pdidetailpayloads, nil)
	if errUrlPdiDetailrequest != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlPdiDetailrequest,
		}
	}
	errUrlProfitCenter := utils.Get(config.EnvConfigs.GeneralServiceUrl+"cost-profit-map?page=0&limit=1000000&profit_center_code=profit_center_gr", profitcenter, nil)
	if errUrlProfitCenter != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlProfitCenter,
		}
	}
	erragreementcompare := tx.Select("mtr_agreement.agreement_document_number").Table("mtr_agreement").
		Where("mtr_agreement.customer_id =?", pdipayload.CompanyID).
		Where("mtr_agreement.profit_center_code=?", profitcenter.ProfitCenterId).
		Where("mtr_agreement.agreement_date_from < ?", time.Now()).
		Where(time.Now(), "?<mtr_agreement.agreement_date_to").Scan(agreementdocno)
	if erragreementcompare != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errUrlProfitCenter,
		}
	}
	erragreement := tx.Select("mtr_agreement.agreement_document_number").Table("mtr_agreement").
		Where("mtr_agreement.agreement_document_number=?", agreementdocno).Scan(agreement).Error
	if erragreement != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        erragreement,
		}
	}
	errapprovalstatuscontractservice := utils.Get(config.EnvConfigs.GeneralServiceUrl+"approval-status-by-code/25", approvalstatus, nil)
	if errapprovalstatuscontractservice != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        erragreement,
		}
	}
	errcontractservice := tx.Select("trx_contract_service.contract_service_system_number").Table("trx_contract_service").
		Where("trx_contract_service.contract_service_status_id=?", approvalstatus.ApprovalStatusId).
		Where("trx_contract_service.contract_service_from < ?", time.Now()).
		Where(time.Now(), "?<trx_contract_service.contract_service_to").
		Where("trx_contract_service.vehicle_id=?", pdidetailpayloads.VehicleId).Scan(contractservice).Error
	if errcontractservice != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errcontractservice,
		}
	}
	entities := transactionworkshopentities.BookingEstimation{
		BrandId:                        pdipayload.BrandID,
		ModelId:                        pdipayload.ModelID,
		VariantId:                      pdipayload.VariantID,
		VehicleId:                      pdidetailpayloads.VehicleId,
		ContractSystemNumber:           contractservice.ContractServiceId,
		CompanyId:                      pdipayload.CompanyID,
		BookingSystemNumber:            0,
		ServiceRequestSystemNumber:     0,
		EstimationSystemNumber:         0,
		AgreementNumberBr:              "",
		AgreementId:                    0,
		ContactPersonName:              "",
		ContactPersonPhone:             "",
		ContactPersonViaId:             0,
		ContactPersonMobile:            "",
		InsurancePolicyNo:              "",
		InsuranceExpiredDate:           time.Time{},
		InsuranceClaimNo:               "",
		InsurancePic:                   "",
		ProfitCenterId:                 profitcenter.ProfitCenterId,
		IsUnregistered:                 false,
		BookingEstimationBatchDate:     time.Now(),
		BookingEstimationVehicleNumber: pdidetailpayloads.VehicleRegistrationCertificateTnkb,
	}
	err := tx.Save(entities).Error
	if err != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}
	entities8 := transactionworkshopentities.BookingEstimationAllocation{
		DocumentStatusID:      15, //document status new
		BatchSystemNumber:     entities.BatchSystemNumber,
		CompanyID:             pdipayload.CompanyID,
		PdiSystemNumber:       id,
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
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
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
	err2 := tx.Save(entities2).Error
	if err2 != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
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
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err3,
		}
	}
	entities3 := transactionworkshopentities.BookingEstimationDetail{
		EstimationSystemNumber:         entities2.EstimationSystemNumber,
		BillID:                         1, //transaction type workorder external
		EstimationLineDiscountApproval: 1, //status draft
		ItemOperationID:                pdidetailpayloads.OperationNumberId,
		LineTypeID:                     1, //line type id where line type description = operation
		RequestDescription:             "",
		FRTQuantity:                    pdidetailpayloads.Frt,
		ItemOperationPrice:             lastprice,
		DiscountItemOperationAmount:    0,
		DiscountItemOperationPercent:   0,
		DiscountRequestPercent:         0,
		DiscountRequestAmount:          0,
	}
	err4 := tx.Save(entities3).Error
	if err4 != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err4,
		}
	}
	err5 := tx.Select("trx_booking_estimation_operation_detail.operation_price").Where("estimation_system_number=?", entities2.EstimationSystemNumber).Scan(&operationtotal)
	if err5 != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err4,
		}
	}
	return entities, nil
}

func (r *BookingEstimationImpl) SaveBookingEstimationFromServiceRequest(tx *gorm.DB, id int) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse) {
	var initialpayloads transactionworkshoppayloads.ServiceRequestBookingEstimation
	var vehiclepayloads transactionworkshoppayloads.VehicleTnkb
	var lastprice float64
	var servicerequestdetail []transactionworkshoppayloads.ServiceRequestDetailBookingPayloads
	time := time.Now()
	err := tx.Select("trx_service_request.profit_center_id,trx_service_request.company_id,trx_service_request.vehicle_id,trx_service_request.service_request_document_number,trx_contraxt_service.contract_service_system_number").
		Joins("JOIN trx_contract_service on trx_contract_service.vehicle_id==trx_service_request.vehicle_id and trx_contract_service.contract_service_to < "+time.String()+" and "+time.String()+" > trx_contract_service.contract_service_from and trx_contract_service.contract_service_status_id = "+strconv.Itoa(20)).
		Where("trx_service_request.service_request_system_number=?", id).Scan(initialpayloads).Error
	if err != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}
	errUrlVehicle := utils.Get(config.EnvConfigs.SalesServiceUrl+"vehicle-master/"+strconv.Itoa(initialpayloads.VehicleId), vehiclepayloads, nil)
	if errUrlVehicle != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
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
		ContactPersonName:              "",
		ContactPersonPhone:             "",
		ContactPersonViaId:             0,
		ContactPersonMobile:            "",
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
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err1,
		}
	}

	entities8 := transactionworkshopentities.BookingEstimationAllocation{
		DocumentStatusID:      15, //document status new
		BatchSystemNumber:     entity.BatchSystemNumber,
		CompanyID:             initialpayloads.CompanyId,
		PdiSystemNumber:       id,
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
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err8,
		}
	}

	entities := transactionworkshopentities.BookingEstimationServiceDiscount{
		BatchSystemNumber:                entity.BatchSystemNumber,
		DocumentStatusID:                 10,
		EstimationDiscountApprovalStatus: 10,
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
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err2,
		}
	}

	err4 := tx.Model(transactionworkshopentities.ServiceRequestDetail{}).Where("service_request_system_number = ?", id).Scan(&servicerequestdetail).Error
	if err4 != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err4,
		}
	}
	for _, detail := range servicerequestdetail {
		err3 := tx.Select("mtr_labour_selling_price_detail.selling_price").
			Table("mtr_labour_selling_price_detail").
			Joins("Join mtr_labour_selling_price on mtr_labour_selling_price.labour_selling_price_id = mtr_labour_selling_price_detail.labour_selling_price_id").
			Where("mtr_labour_selling_price.brand_id =?", vehiclepayloads.VehicleBrandId).
			Where("mtr_labour_selling_price.company_id = ?", initialpayloads.CompanyId).
			Where("mtr_labour_selling_price.effective_date < ?", time).
			Scan(&lastprice).Error
		if err3 != nil {
			return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err3,
			}
		}
		entities3 := transactionworkshopentities.BookingEstimationDetail{
			EstimationSystemNumber:         entities.EstimationSystemNumber,
			BillID:                         1, //transaction type workorder external
			EstimationLineDiscountApproval: 1, //status draft
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
			return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err4,
			}
		}
	}
	return entity, nil
}

func (r *BookingEstimationImpl) SaveBookingEstimationAllocation(tx *gorm.DB, id int, req transactionworkshoppayloads.BookEstimationAllocation) (transactionworkshopentities.BookingEstimationAllocation, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.BookingEstimationAllocation{
		DocumentStatusID:      req.DocumentStatusID,
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
	err := tx.Save(&entities).Error
	if err != nil {
		return transactionworkshopentities.BookingEstimationAllocation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}
	return entities, nil
}


func (r *BookingEstimationImpl) AddContractService(tx *gorm.DB,idheader int, Idcontract int)(bool,*exceptions.BaseErrorResponse){
	var firststruct transactionworkshoppayloads.BookingEstimationFirstContractService
	var contractservice transactionworkshoppayloads.ContractService
	var taxfare float64
	now := time.Now()
	var count int64
	err := tx.Raw(`
		SELECT BE.contract_service_system_number, BE2.estimation_discount_approval_status, BE.booking_system_number, BE.estimation_system_number, 
		       BE.brand_id, BE.profit_center_id, BE.model_id, BE.Company_id
		FROM trx_booking_estimation BE
		LEFT JOIN trx_booking_estimation_service_discount BE2 ON BE.batch_system_number = BE2.batch_system_number
		WHERE BE.batch_system_number = ?
	`, idheader).Scan(&firststruct).Error
	if err != nil {
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err,
		}
	}

	errUrlGetTax:= utils.Get(config.EnvConfigs.FinanceServiceUrl+"tax-fare/detail/tax-percent?tax_service_code=PPN&tax_type_code=PPN&effective_date="+time.Now().String(), taxfare, nil)
	if errUrlGetTax != nil {
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errUrlGetTax,
		}
	}
	entities := transactionworkshopentities.BookingEstimationServiceDiscount{
		DocumentStatusID: 10,//status new
		BatchSystemNumber: idheader,
		EstimationDate: &now,
		EstimationDiscountApprovalStatus: 10,//status draft
		CompanyID: firststruct.CompanyId,
		VATTaxRate: taxfare,
	}
	err3:=tx.Save(&entities).Error
	if err3 != nil{
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err3,
		}
	}

	err4:= tx.Select("trx_booking_estimation_detail.estimation_line_code").Table("trx_booking_estimation_detail").Where("estimation_system_number=?",entities.EstimationSystemNumber).Count(&count).Error
	if err4 != nil{
		return false,&exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err: err4,
		}
	}
	err5:= tx.Select("CSD.line_type_id,CSD.item_operation_id,CSD.description,CSD.frt_quantity,CSD.item_price,CSD.item_discount_percent").
	Table("trx_contract_service CS").
	Joins("Join trx_contract_service_detail CSD on CSD.contract_service_system_number = CS.contract_service_system_number").
	Where("CS.contract_service_system_number=?",Idcontract).Scan(&contractservice).Error
	if err5 != nil{
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err5,
		}
	}

	entities2:= transactionworkshopentities.BookingEstimationDetail{
		EstimationSystemNumber: entities.EstimationSystemNumber,
		EstimationLineCode: int(count)+1,
		LineTypeID: contractservice.LineTypeId,
		BillID: 1,//id of TRXTYPE_WO_CONTRACT_SERVICE
		EstimationLineDiscountApproval: 10,//id of approval status draft
		ItemOperationID: contractservice.ItemOperationId,
		RequestDescription: contractservice.Description,
		FRTQuantity: float64(contractservice.FrtQuantity),
		ItemOperationPrice: contractservice.ItemPrice,
		DiscountItemOperationAmount: math.Round(contractservice.ItemPrice*contractservice.ItemDiscountPercent/100),
		DiscountRequestAmount: 0,
		DiscountRequestPercent:0,
		DiscountItemOperationPercent:  contractservice.ItemDiscountPercent,
		DiscountApprovalBy: "",
		DiscountApprovalDate: nil,
	}
	err6:= tx.Save(&entities2).Error
	if err6 != nil{
			return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err6,
		}
	}
	_,err7:= r.PutBookingEstimationCalculation(tx,idheader)
	if err7 != nil{
		return false,err7
	}
	if contractservice.LineTypeId == 0{
		entities3:= transactionworkshopentities.BookingEstimationAllocation{
			BookingStall: "",
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
	return true,nil
}

func (r *BookingEstimationImpl) AddPackage(tx *gorm.DB, idhead int, idpackage int)(bool,*exceptions.BaseErrorResponse){
	var headerdata transactionworkshoppayloads.BookingEstimationFirstContractService
	var taxfare float64
	var count int64
	var uom int
	var price float64
	var discpercent float64
	var data transactionworkshoppayloads.PackageForDetail
	time := time.Now()

	err := tx.Raw(`
		SELECT BE.contract_service_system_number, BE2.estimation_discount_approval_status, BE.booking_system_number, BE.estimation_system_number, 
		       BE.brand_id, BE.profit_center_id, BE.model_id, BE.Company_id
		FROM trx_booking_estimation BE
		LEFT JOIN trx_booking_estimation_service_discount BE2 ON BE.batch_system_number = BE2.batch_system_number
		WHERE BE.batch_system_number = ?
	`, idhead).Scan(&headerdata).Error
	if err != nil {
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err,
		}
	}

	errUrlGetTax:= utils.Get(config.EnvConfigs.FinanceServiceUrl+"tax-fare/detail/tax-percent?tax_service_code=PPN&tax_type_code=PPN&effective_date="+time.String(), taxfare, nil)
	if errUrlGetTax != nil {
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errUrlGetTax,
		}
	}

	entities := transactionworkshopentities.BookingEstimationServiceDiscount{
		EstimationDocumentNumber: headerdata.EstimationDocumentNumber,
		EstimationDate: &time,
		DocumentStatusID: 5,//status new
		EstimationDiscountApprovalStatus: 10,//approval draft
		BatchSystemNumber: idhead,
		CompanyID: headerdata.CompanyId,
		VATTaxRate:  taxfare,
	}
	err3:=tx.Save(&entities).Error
	if err3 != nil{
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err3,
		}
	}

	err4 := tx.Table("trx_booking_estimation_detail AS tb").
    Select("tb.estimation_line_id").
    Where("tb.estimation_system_number = ?", entities.EstimationSystemNumber).
    Count(&count).Error	
	if err4 != nil{
		return false,&exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err: err4,
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
	if err5 != nil{
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err5,
		}
	}
	err6 := tx.Select("mtr_item.unit_of_measurement_type_id").Table("mtr_item").Where("mtr_item.item_name=?",data.ItemOrOperationName).Scan(&uom).Error
	if err6 != nil{
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err6,
		}
	}
	err7:= tx.Select("mtr_price_list.price_list_amount").Table("mtr_price_list").
	Joins("join mtr_item on mtr_price_list.price_list_id = mtr_item.price_list_item").
	Joins("join mtr_item_operation on mtr_item_operation.item_operation_model_mapping_id=mtr_item.item_id").Where("mtr_item_operation.item_operation_id=?",data.ItemOperationId).
	Scan(&price).Error
	if err7 != nil{
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err7,
		}
	}
	err8:= tx.Select("trx_contract_service_detail.item_discount_percent").
	Table("trx_contract_service_detail").
	Joins("Join trx_contract_service on trx_contract_service.contract_service_system_number=trx_contract_service_detail.contract_service_system_number").
	Where("trx_contract_service.contract_service_system_number=0").
	Where("trx_contract_service_detail.item_operation_id=?",data.ItemOperationId).
	Where("trx_contract_service_detail.line_type_id=?",data.LineTypeId).Scan(&discpercent).Error
	if err8 != nil{
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err8,
		}
	}
	operationdiscount:=math.Round(price*data.FrtQuantity*discpercent/100)
	entities2 := transactionworkshopentities.BookingEstimationDetail{
		EstimationSystemNumber: entities.EstimationSystemNumber,
		EstimationLineCode: int(count+1),
		LineTypeID: data.LineTypeId,
		JobTypeID: data.JobTypeId,
		BillID: data.BillId,
		EstimationLineDiscountApproval: 10,//approval draft id
		ItemOperationID: data.ItemOperationId,
		RequestDescription: data.ItemOrOperationName,
		PackageID: idpackage,
		UOMID: uom,
		FRTQuantity: data.FrtQuantity,
		ItemOperationPrice: price,
		DiscountItemOperationAmount: operationdiscount,
		DiscountItemOperationPercent: discpercent,
	}
	err9 := tx.Save(&entities2).Error
	if err9 != nil{
		return false,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err9,
		}
	}
	_,err0:=r.PutBookingEstimationCalculation(tx,idhead)
	if err0 != nil{
		return false,err0
	}
	return true,nil
}