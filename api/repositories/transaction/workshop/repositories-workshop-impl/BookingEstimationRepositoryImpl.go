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

func (r *BookingEstimationImpl) GetAllBookEstimReq(tx *gorm.DB, pages *pagination.Pagination, id int) ([]transactionworkshoppayloads.BookEstimRemarkRequest, *exceptions.BaseErrorResponse) {
	var payloads []transactionworkshoppayloads.BookEstimRemarkRequest
	var model []transactionworkshopentities.BookingEstimationRequest

	err := tx.Model(&model).Where("booking_system_number = ?", id).Scan(&payloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.Rows = payloads
	return payloads, nil
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

func (r *BookingEstimationImpl) SaveDetailBookEstim(tx *gorm.DB, req transactionworkshoppayloads.BookEstimDetailReq) (int, *exceptions.BaseErrorResponse) {
	var bookestimcalc transactionworkshopentities.BookingEstimationServiceDiscount
	var bookestimpayloads []transactionworkshoppayloads.BookingEstimationCalculationPayloads
	if req.LineTypeID == 5 {
		entity := transactionworkshopentities.BookingEstimationOperationDetail{
			EstimationLineID:               req.EstimationLineID,
			EstimationLineCode:             req.EstimationLineCode,
			EstimationSystemNumber:         req.EstimationSystemNumber,
			BillID:                         req.BillID,
			EstimationLineDiscountApproval: req.EstimationLineDiscountApproval,
			OperationId:                    req.OperationId,
			LineTypeID:                     req.LineTypeID,
			PackageID:                      req.PackageID,
			JobTypeID:                      req.JobTypeID,
			FieldActionSystemNumber:        req.FieldActionSystemNumber,
			ApprovalRequestNumber:          req.ApprovalRequestNumber,
			UOMID:                          req.UOMID,
			RequestDescription:             req.RequestDescription,
			FRTQuantity:                    req.FRTQuantity,
			OperationPrice:                 req.OperationItemPrice,
			DiscountOperationAmount:        req.DiscountItemAmount,
			DiscountOperationPercent:       req.DiscountItemPercent,
			DiscountRequestPercent:         req.DiscountRequestPercent,
			DiscountRequestAmount:          req.DiscountRequestAmount,
		}

		if err := tx.Save(&entity).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	} else {
		entity := transactionworkshopentities.BookingEstimationItemDetail{
			EstimationLineID:               req.EstimationLineID,
			EstimationLineCode:             req.EstimationLineCode,
			EstimationSystemNumber:         req.EstimationSystemNumber,
			BillID:                         req.BillID,
			EstimationLineDiscountApproval: req.EstimationLineDiscountApproval,
			ItemID:                         req.ItemID,
			LineTypeID:                     req.LineTypeID,
			FieldActionSystemNumber:        req.FieldActionSystemNumber,
			ApprovalRequestNumber:          req.ApprovalRequestNumber,
			UOMID:                          req.UOMID,
			RequestDescription:             req.RequestDescription,
			FRTQuantity:                    req.FRTQuantity,
			ItemPrice:                      req.OperationItemPrice,
			DiscountItemAmount:             req.DiscountItemAmount,
			DiscountItemPercent:            req.DiscountItemPercent,
			DiscountRequestPercent:         req.DiscountRequestPercent,
			DiscountRequestAmount:          req.DiscountRequestAmount,
		}

		if err := tx.Save(&entity).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}

	err := tx.Select(&bookestimcalc).Where("batch_system_number = ?", req.EstimationSystemNumber).Scan(&bookestimpayloads).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return bookestimpayloads[0].BatchSystemNumber, nil
}

func (r *BookingEstimationImpl) UpdateBookEstimDetail(tx *gorm.DB, req transactionworkshoppayloads.BookEstimDetailUpdate, id int, LineTypeId int) (bool, *exceptions.BaseErrorResponse) {
	if LineTypeId == 5 {
		var model transactionworkshopentities.BookingEstimationOperationDetail
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
	} else {
		var model transactionworkshopentities.BookingEstimationItemDetail
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
	}
	return true, nil
}

func (r *BookingEstimationImpl) DeleteBookEstimDetail(tx *gorm.DB, id int, linetypeid int) (bool, *exceptions.BaseErrorResponse) {
	if linetypeid == 5 {
		var model transactionworkshopentities.BookingEstimationOperationDetail
		err := tx.Delete(&model).Where("estimation_line_id =?", id).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
	} else {
		var model transactionworkshopentities.BookingEstimationItemDetail
		err := tx.Delete(&model).Where("estimation_line_id =?", id).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
	}
	return true, nil
}

func (r *BookingEstimationImpl) CopyFromHistory(tx *gorm.DB, id int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	var modeloperation transactionworkshopentities.BookingEstimationOperationDetail
	var modelitem transactionworkshopentities.BookingEstimationItemDetail
	var operationpayloads []transactionworkshoppayloads.BookEstimOperationPayloads
	var itempayloads []transactionworkshoppayloads.BookEstimItemPayloads

	// Slice to hold the payloads
	var payloads []map[string]interface{}

	// Query item details
	err := tx.Model(&modelitem).Where("estimation_system_number = ?", id).Scan(&itempayloads).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	// Process item details
	for _, item := range itempayloads {
		entity := transactionworkshopentities.BookingEstimationItemDetail{
			EstimationLineID:               item.EstimationLineID,
			EstimationLineCode:             item.EstimationLineCode,
			EstimationSystemNumber:         id,
			BillID:                         item.BillID,
			EstimationLineDiscountApproval: item.EstimationLineDiscountApproval,
			ItemID:                         item.ItemID,
			LineTypeID:                     item.LineTypeID,
			PackageID:                      item.PackageID,
			JobTypeID:                      item.JobTypeID,
			FieldActionSystemNumber:        item.FieldActionSystemNumber,
			ApprovalRequestNumber:          item.ApprovalRequestNumber,
			UOMID:                          item.UOMID,
			RequestDescription:             item.RequestDescription,
			FRTQuantity:                    item.FRTQuantity,
			ItemPrice:                      item.OperationItemPrice,
			DiscountItemAmount:             item.DiscountItemAmount,
			DiscountItemPercent:            item.DiscountItemPercent,
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
			"item_id":              item.ItemID,
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

	// Query operation details
	err2 := tx.Model(&modeloperation).Where("estimation_system_number = ?", id).Scan(&operationpayloads).Error
	if err2 != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err2,
		}
	}

	// Process operation details
	for _, operation := range operationpayloads {
		entity := transactionworkshopentities.BookingEstimationOperationDetail{
			EstimationLineID:               operation.EstimationLineID,
			EstimationLineCode:             operation.EstimationLineCode,
			EstimationSystemNumber:         id,
			BillID:                         operation.BillID,
			EstimationLineDiscountApproval: operation.EstimationLineDiscountApproval,
			OperationId:                    operation.OperationId,
			LineTypeID:                     operation.LineTypeID,
			FieldActionSystemNumber:        operation.FieldActionSystemNumber,
			ApprovalRequestNumber:          operation.ApprovalRequestNumber,
			UOMID:                          operation.UOMID,
			RequestDescription:             operation.RequestDescription,
			FRTQuantity:                    operation.FRTQuantity,
			OperationPrice:                 operation.OperationItemPrice,
			DiscountOperationAmount:        operation.DiscountItemAmount,
			DiscountOperationPercent:       operation.DiscountItemPercent,
			DiscountRequestPercent:         operation.DiscountRequestPercent,
			DiscountRequestAmount:          operation.DiscountRequestAmount,
		}

		if err := tx.Save(&entity).Error; err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		// Add operation payload to the payloads slice
		payload := map[string]interface{}{
			"estimation_line_id":   operation.EstimationLineID,
			"estimation_line_code": operation.EstimationLineCode,
			"operation_id":         operation.OperationId,
			"line_type_id":         operation.LineTypeID,
			"request_description":  operation.RequestDescription,
			"frt_quantity":         operation.FRTQuantity,
			"operation_price":      operation.OperationItemPrice,
			"item_name":            "", // Empty item_name
			"operation_name":       operation.RequestDescription,
		}
		payloads = append(payloads, payload)
	}

	return payloads, nil
}

func (r *BookingEstimationImpl) AddPackage(tx *gorm.DB, id int, packId int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	var model masterentities.PackageMasterDetail
	var operationpayloads []masterpayloads.CampaignMasterDetailGetPayloads
	var payloads []map[string]interface{}
	err2 := tx.Model(&model).Where("package_id = ?", packId).Scan(&operationpayloads).Error
	if err2 != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err2,
		}
	}
	for _, item := range operationpayloads {
		entity := transactionworkshopentities.BookingEstimationItemDetail{
			EstimationSystemNumber: id,
			ItemID:                 item.ItemOperationId,
			LineTypeID:             item.LineTypeId,
			PackageID:              item.PackageId,
			ItemPrice:              float64(item.PackageId),
		}
		if err := tx.Save(&entity).Error; err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		payload := map[string]interface{}{
			"item_id":      item.ItemOperationId,
			"line_type_id": item.LineTypeId,
			"package_id":   item.PackageId,
			"item_price":   float64(item.PackageId),
		}
		payloads = append(payloads, payload)
	}
	err3 := tx.Model(&model).Where("estimation_system_number = ?", id).Scan(&operationpayloads).Error
	if err3 != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err3,
		}
	}
	for _, operation := range operationpayloads {
		entity := transactionworkshopentities.BookingEstimationOperationDetail{
			EstimationSystemNumber: id,
			OperationId:            operation.ItemOperationId,
			LineTypeID:             operation.LineTypeId,
			PackageID:              operation.PackageId,
			FRTQuantity:            operation.Quantity,
			OperationPrice:         operation.Price,
		}
		if err := tx.Save(&entity).Error; err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
		payload := map[string]interface{}{
			"operation_id":    operation.ItemOperationId,
			"line_type_id":    operation.LineTypeId,
			"package_id":      operation.PackageId,
			"frt_quantity":    operation.Quantity,
			"operation_price": operation.Price,
		}
		payloads = append(payloads, payload)
	}

	return payloads, nil
}

func (r *BookingEstimationImpl) AddContractService(tx *gorm.DB, id int, contractserviceid int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	var model transactionworkshopentities.ContractService
	var modeloperation []transactionworkshopentities.ContractServiceOperationDetail
	var modelitem []transactionworkshopentities.ContractServiceItemDetail

	// Slice to hold the payloads
	var payloads []map[string]interface{}

	// Query the contract service
	err := tx.Model(&model).Where("contract_service_system_number = ?", contractserviceid).Scan(&model).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	// Query the item details
	err2 := tx.Model(&modelitem).Where("contract_service_system_number = ?", contractserviceid).Scan(&modelitem).Error
	if err2 != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err2,
		}
	}

	// Process the item details
	for _, item := range modelitem {
		entity := transactionworkshopentities.BookingEstimationItemDetail{
			EstimationSystemNumber: id,
			ItemID:                 item.ItemId,
			LineTypeID:             item.LineTypeId,
			PackageID:              item.PackageId,
			RequestDescription:     item.Description,
			FRTQuantity:            item.FrtQuantity,
			ItemPrice:              item.ItemPrice,
		}

		if err := tx.Save(&entity).Error; err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		// Add item payload to the payloads slice
		payload := map[string]interface{}{
			"item_id":      item.ItemId,
			"line_type_id": item.LineTypeId,
			"package_id":   item.PackageId,
			"item_name":    item.Description,
			"frt_quantity": item.FrtQuantity,
			"item_price":   item.ItemPrice,
		}
		payloads = append(payloads, payload)
	}

	// Query the operation details
	err3 := tx.Model(&modeloperation).Where("contract_service_system_number = ?", contractserviceid).Scan(&modeloperation).Error
	if err3 != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err3,
		}
	}

	// Process the operation details
	for _, operation := range modeloperation {
		entity := transactionworkshopentities.BookingEstimationOperationDetail{
			EstimationSystemNumber: id,
			OperationId:            operation.OperationId,
			LineTypeID:             operation.LineTypeId,
			RequestDescription:     operation.Description,
			FRTQuantity:            operation.FrtQuantity,
			OperationPrice:         operation.OperationPrice,
		}

		if err := tx.Save(&entity).Error; err != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}

		// Add operation payload to the payloads slice
		payload := map[string]interface{}{
			"operation_id":    operation.OperationId,
			"line_type_id":    operation.LineTypeId,
			"operation_name":  operation.Description,
			"frt_quantity":    operation.FrtQuantity,
			"operation_price": operation.OperationPrice,
		}
		payloads = append(payloads, payload)
	}

	return payloads, nil
}

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
	}

	for _, detail := range itemDetails {
		err := tx.Model(&transactionworkshopentities.BookingEstimationItemDetail{}).
			Where("estimation_system_number = ? AND line_type_id = ?", id, detail.LineTypeID).
			Update("discount_item_percent", detail.Value).Error
		if err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Err:        err,
			}
		}
	}

	operationDetails := []struct {
		LineTypeID int
		Value      int
	}{
		{5, req.Operation},
		{6, req.PackageDiscount},
	}

	for _, detail := range operationDetails {
		err := tx.Model(&transactionworkshopentities.BookingEstimationOperationDetail{}).
			Where("estimation_system_number = ? AND line_type_id = ?", id, detail.LineTypeID).
			Update("discount_operation_percent", detail.Value).Error
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
	var modeloperation []transactionworkshopentities.ContractServiceOperationDetail
	var modelitem []transactionworkshopentities.ContractServiceItemDetail
	err2 := tx.Model(&modelitem).Where("contract_service_system_number = ?", idrecall).Scan(&modelitem).Error
	if err2 != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err2,
		}
	}
	for _, item := range modelitem {
		entity := transactionworkshopentities.BookingEstimationItemDetail{
			EstimationSystemNumber: id,
			ItemID:                 item.ItemId,
			LineTypeID:             item.LineTypeId,
			PackageID:              item.PackageId,
			RequestDescription:     item.Description,
			FRTQuantity:            item.FrtQuantity,
			ItemPrice:              item.ItemPrice,
		}

		if err := tx.Save(&entity).Error; err != nil {
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	err3 := tx.Model(&modeloperation).Where("contract_service_system_number = ?", id).Scan(&modeloperation).Error
	if err3 !=

		nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err3,
		}
	}
	for _, operation := range modeloperation {
		entity := transactionworkshopentities.BookingEstimationOperationDetail{
			EstimationSystemNumber: id,
			OperationId:            operation.OperationId,
			LineTypeID:             operation.LineTypeId,
			RequestDescription:     operation.Description,
			FRTQuantity:            operation.FrtQuantity,
			OperationPrice:         operation.OperationPrice,
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
	entity := transactionworkshopentities.BookingEstimationServiceDiscount{
		BatchSystemNumber:                id,
		DocumentStatusID:                 0,
		EstimationDiscountApprovalStatus: 10,
		CompanyID:                        0,
		ApprovalRequestNumber:            0,
		EstimationDocumentNumber:         "",
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
	err := tx.Save(entity).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return id, nil
}

func (r *BookingEstimationImpl) PutBookingEstimationCalculationPutBookingEstimationCalculation(tx *gorm.DB, id int, linetypeid int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	const (
		LineTypePackage            = 0 // Package Bodyshop
		LineTypeOperation          = 1 // Operation
		LineTypeSparePart          = 2 // Spare Part
		LineTypeOil                = 3 // Oil
		LineTypeMaterial           = 4 // Material
		LineTypeFee                = 5 // Fee
		LineTypeAccessories        = 6 // Accessories
		LineTypeConsumableMaterial = 7 // Consumable Material
		LineTypeSublet             = 8 // Sublet
		LineTypeSouvenir           = 9 // Souvenir
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
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0), 0) ELSE 0 END) AS total_package,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_operation,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_spare_part,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_oil,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_material,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_fee,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_accessories,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_consumable_material,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_sublet,
			SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_souvenir
		FROM trx_booking_estimation_service_discount
		WHERE batch_system_number = ?`,
		LineTypePackage,
		LineTypeOperation,
		LineTypeSparePart,
		LineTypeOil,
		LineTypeMaterial,
		LineTypeFee,
		LineTypeAccessories,
		LineTypeConsumableMaterial,
		LineTypeSublet,
		LineTypeSouvenir,
		id).Scan(&result).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: fmt.Sprintf("Failed to calculate totals: %v", err)}
	}

	// Calculate grand total
	grandTotal := result.TotalPackage + result.TotalOperation + result.TotalSparePart + result.TotalOil + result.TotalMaterial + result.TotalFee + result.TotalAccessories + result.TotalConsumableMaterial + result.TotalSublet + result.TotalSouvenir

	// Update Work Order with the calculated totals
	err = tx.Model(&transactionworkshopentities.BookingEstimationServiceDiscount{}).
		Where("batch_system_number = ?", id).
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
	BookingEstimationResponse := []map[string]interface{}{
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

	return BookingEstimationResponse, nil
}

func (r *BookingEstimationImpl) PutBookingEstimationCalculation(tx *gorm.DB, id int, linetypeid int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	const (
		LineTypePackage            = 0 // Package Bodyshop
		LineTypeOperation          = 1 // Operation
		LineTypeSparePart          = 2 // Spare Part
		LineTypeOil                = 3 // Oil
		LineTypeMaterial           = 4 // Material
		LineTypeFee                = 5 // Fee
		LineTypeAccessories        = 6 // Accessories
		LineTypeConsumableMaterial = 7 // Consumable Material
		LineTypeSublet             = 8 // Sublet
		LineTypeSouvenir           = 9 // Souvenir
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
				SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0), 0) ELSE 0 END) AS total_package,
				SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_operation,
				SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_spare_part,
				SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_oil,
				SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_material,
				SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_fee,
				SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_accessories,
				SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_consumable_material,
				SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_sublet,
				SUM(CASE WHEN line_type_id = ? THEN ROUND(ISNULL(operation_item_price, 0) * ISNULL(frt_quantity, 0), 0) ELSE 0 END) AS total_souvenir
			FROM trx_work_order_detail
			WHERE work_order_system_number = ?`,
		LineTypePackage,
		LineTypeOperation,
		LineTypeSparePart,
		LineTypeOil,
		LineTypeMaterial,
		LineTypeFee,
		LineTypeAccessories,
		LineTypeConsumableMaterial,
		LineTypeSublet,
		LineTypeSouvenir,
		id).Scan(&result).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	// Calculate grand total
	grandTotal := result.TotalPackage + result.TotalOperation + result.TotalSparePart + result.TotalOil + result.TotalMaterial + result.TotalFee + result.TotalAccessories + result.TotalConsumableMaterial + result.TotalSublet + result.TotalSouvenir

	// Update Work Order with the calculated totals
	err = tx.Model(&transactionworkshopentities.BookingEstimationServiceDiscount{}).
		Where("estimation_system_number = ?", id).
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
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	// Prepare response
	estimationcalculation := []map[string]interface{}{
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

	return estimationcalculation, nil
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
	entities3 := transactionworkshopentities.BookingEstimationOperationDetail{
		EstimationSystemNumber:         entities2.EstimationSystemNumber,
		BillID:                         1, //transaction type workorder external
		EstimationLineDiscountApproval: 1, //status draft
		OperationId:                    pdidetailpayloads.OperationNumberId,
		LineTypeID:                     1, //line type id where line type description = operation
		RequestDescription:             "",
		FRTQuantity:                    pdidetailpayloads.Frt,
		OperationPrice:                 lastprice,
		DiscountOperationAmount:        0,
		DiscountOperationPercent:       0,
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
	err := tx.Select("trx_service_request.profit_center_id,trx_service_request.company_id,trx_service_request.vehicle_id,trx_service_request.service_request_document_number,trx_contraxt_service.contract_service_system_number").
		Joins("JOIN trx_contract_service on trx_contract_service.vehicle_id==trx_service_request.vehicle_id and trx_contract_service.contract_service_to < "+time.Now().Format("2006-01-02 15:04:05")+" and "+time.Now().Format("2006-01-02 15:04:05")+" > trx_contract_service.contract_service_from and trx_contract_service.contract_service_status_id = "+strconv.Itoa(20)).
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
		InsuranceExpiredDate:           time.Time{},
		InsuranceClaimNo:               "",
		InsurancePic:                   "",
		ProfitCenterId:                 initialpayloads.ProfitCenterId,
		IsUnregistered:                 false,
		BookingEstimationBatchDate:     time.Now(),
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

	now := time.Now()
	entities := transactionworkshopentities.BookingEstimationServiceDiscount{
		BatchSystemNumber:                entity.BatchSystemNumber,
		DocumentStatusID:                 10,
		EstimationDiscountApprovalStatus: 10,
		CompanyID:                        entity.CompanyId,
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
		if detail.LineTypeId != 5 { //set for line type id is not operation
			err3 := tx.Select("mtr_labour_selling_price_detail.selling_price").
				Table("mtr_labour_selling_price_detail").
				Joins("Join mtr_labour_selling_price on mtr_labour_selling_price.labour_selling_price_id = mtr_labour_selling_price_detail.labour_selling_price_id").
				Where("mtr_labour_selling_price.brand_id =?", vehiclepayloads.VehicleBrandId).
				Where("mtr_labour_selling_price.company_id = ?", initialpayloads.CompanyId).
				Where("mtr_labour_selling_price.effective_date < ?", time.Now()).
				Scan(&lastprice).Error
			if err3 != nil {
				return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusConflict,
					Err:        err3,
				}
			}
			entities3 := transactionworkshopentities.BookingEstimationOperationDetail{
				EstimationSystemNumber:         entities.EstimationSystemNumber,
				BillID:                         1, //transaction type workorder external
				EstimationLineDiscountApproval: 1, //status draft
				OperationId:                    detail.OperationItemId,
				LineTypeID:                     detail.LineTypeId, //line type id where line type description = operation
				RequestDescription:             "",
				FRTQuantity:                    detail.FrtQuantity,
				OperationPrice:                 lastprice,
				DiscountOperationAmount:        0,
				DiscountOperationPercent:       0,
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
		} else {
			err3 := tx.Select("mtr_item.last_price").Where("mtr_item.item_id = ?", detail.OperationItemId).Scan(lastprice).Error
			if err3 != nil {
				return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusConflict,
					Err:        err3,
				}
			}
			entities3 := transactionworkshopentities.BookingEstimationItemDetail{
				EstimationSystemNumber:         entities.EstimationSystemNumber,
				BillID:                         1, //transaction type workorder external
				EstimationLineDiscountApproval: 1, //status draft
				ItemID:                         detail.OperationItemId,
				LineTypeID:                     detail.LineTypeId, //line type id where line type description = operation
				RequestDescription:             "",
				FRTQuantity:                    detail.FrtQuantity,
				ItemPrice:                      lastprice,
				DiscountItemAmount:             0,
				DiscountItemPercent:            0,
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
