package transactionworkshoprepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	masterpackagemasterentity "after-sales/api/entities/master/package-master"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"errors"
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

func (r *BookingEstimationImpl) New(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptions.BaseErrorResponse) {
	// Create a new instance of WorkOrderRepositoryImpl
	// Save the booking estimation
	success, err := r.Save(tx, request) // Menggunakan method Save dari receiver saat ini, yaitu r
	if err != nil {
		return false, &exceptions.BaseErrorResponse{Message: "Failed to save booking estimation"}
	}

	return success, nil
}

func (r *BookingEstimationImpl) GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.BookingEstimationRequest, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.BookingEstimation
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Where("batch_system_number = ?", Id).First(&entity).Error
	if err != nil {
		return transactionworkshoppayloads.BookingEstimationRequest{}, &exceptions.BaseErrorResponse{Message: "Failed to retrieve booking estimation from the database"}
	}

	// Convert entity to payload
	payload := transactionworkshoppayloads.BookingEstimationRequest{}

	return payload, nil
}

func (r *BookingEstimationImpl) Save(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (bool, *exceptions.BaseErrorResponse) {
	var bookingEstimationEntities = transactionworkshopentities.BookingEstimation{
		BrandId: request.BrandId,
		ModelId: request.ModelId,
		VehicleId: request.VehicleId,
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

func (r *BookingEstimationImpl) SaveBookEstimReq(tx *gorm.DB, req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (bool, *exceptions.BaseErrorResponse) {
	var count int64
	err := tx.Select("trx_booking_estim.*").Where("booking_system_number =?", id).Count(&count).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	count32 := int32(count)
	total := count32 + 1
	entities := transactionworkshopentities.BookingEstimationRequest{
		BookingEstimationRequestCode: int(total),
		BookingSystemNumber:          id,
		BookingServiceRequest:        req.BookingServiceRequest,
	}
	err2 := tx.Save(&entities).Error
	if err2 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return true, nil
}

func (r *BookingEstimationImpl) UpdateBookEstimReq(tx *gorm.DB, req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (bool, *exceptions.BaseErrorResponse) {
	model := transactionworkshopentities.BookingEstimationRequest{}
	result := tx.Model(&model).Where(transactionworkshopentities.BookingEstimationRequest{BookingEstimationRequestID: id}).First(&model).Updates(req)
	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if model == (transactionworkshopentities.BookingEstimationRequest{}) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	return true, nil
}

func (r *BookingEstimationImpl) DeleteBookEstimReq(tx *gorm.DB, ids string) (bool, *exceptions.BaseErrorResponse) {
	model := transactionworkshopentities.BookingEstimationRequest{}
	if err := tx.Delete(&model, ids).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *BookingEstimationImpl) GetByIdBookEstim(tx *gorm.DB, id int) (transactionworkshoppayloads.BookEstimRemarkRequest, *exceptions.BaseErrorResponse) {
	var model transactionworkshopentities.BookingEstimationRequest
	var payloads transactionworkshoppayloads.BookEstimRemarkRequest
	err := tx.Model(&model).Where("booking_estimation_request_id = ?", id).Scan(&payloads).Error
	if err != nil {
		return payloads, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	return payloads, nil
}

func (r *BookingEstimationImpl) GetAllBookEstimReq(tx *gorm.DB, pages pagination.Pagination, id int) ([]transactionworkshoppayloads.BookEstimRemarkRequest, *exceptions.BaseErrorResponse) {
	var payloads []transactionworkshoppayloads.BookEstimRemarkRequest
	var model []transactionworkshopentities.BookingEstimationRequest
	err := tx.Model(&model).Scan(&payloads).Where("booking_system_number = ?", id).Error

	if len(payloads) == 0 {
		return payloads, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if err != nil {

		return payloads, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	pages.Rows = payloads
	return payloads, nil
}

func (r *BookingEstimationImpl) SaveBookEstimReminderServ(tx *gorm.DB, req transactionworkshoppayloads.ReminderServicePost, id int) (bool, *exceptions.BaseErrorResponse) {
	var count int64
	err := tx.Select("trx_booking_estimation_request.*").Where("booking_system_number =?", id).Count(&count).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return true, nil
}

func (r *BookingEstimationImpl) SaveDetailBookEstim(tx *gorm.DB, req transactionworkshoppayloads.BookEstimDetailReq) (bool, *exceptions.BaseErrorResponse) {
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
			return false, &exceptions.BaseErrorResponse{
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
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	return true, nil
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

func (r *BookingEstimationImpl) CopyFromHistory(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	var modeloperation transactionworkshopentities.BookingEstimationOperationDetail
	var modelitem transactionworkshopentities.BookingEstimationItemDetail
	var operationpayloads []transactionworkshoppayloads.BookEstimOperationPayloads
	var itempayloads []transactionworkshoppayloads.BookEstimItemPayloads
	err := tx.Model(&modelitem).Where("estimation_system_number = ?", id).Scan(&itempayloads).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}
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
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	err2 := tx.Model(&modeloperation).Where("estimation_system_number = ?", id).Scan(&operationpayloads).Error
	if err2 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}
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
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	return true, nil
}

func (r *BookingEstimationImpl) AddPackage(tx *gorm.DB, id int, packId int) (bool, *exceptions.BaseErrorResponse) {
	var modeloperation masterpackagemasterentity.PackageMasterDetailOperation
	var modelitem masterpackagemasterentity.PackageMasterDetailItem
	var operationpayloads []masterpayloads.CampaignMasterDetailOperationPayloads
	var itempayloads []masterpayloads.PackageMasterDetailItem
	err2 := tx.Model(&modelitem).Where("package_id = ?", packId).Scan(&itempayloads).Error
	if err2 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err2,
		}
	}
	for _, item := range itempayloads {
		entity := transactionworkshopentities.BookingEstimationItemDetail{
			EstimationSystemNumber: id,
			ItemID:                 item.ItemId,
			LineTypeID:             item.LineTypeId,
			PackageID:              item.PackageId,
			RequestDescription:     item.ItemName,
			FRTQuantity:            item.FrtQuantity,
			ItemPrice:              float64(item.PackageId),
		}

		if err := tx.Save(&entity).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	err3 := tx.Model(&modeloperation).Where("estimation_system_number = ?", id).Scan(&operationpayloads).Error
	if err3 != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err3,
		}
	}
	for _, operation := range operationpayloads {
		entity := transactionworkshopentities.BookingEstimationOperationDetail{
			EstimationSystemNumber: id,
			OperationId:            operation.OperationId,
			LineTypeID:             operation.LineTypeId,
			PackageID:              operation.PackageId,
			RequestDescription:     operation.OperationName,
			FRTQuantity:            operation.Quantity,
			OperationPrice:         operation.Price,
		}

		if err := tx.Save(&entity).Error; err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	return true, nil
}

func (r *BookingEstimationImpl) AddContractService(tx *gorm.DB, id int, contractserviceid int) (bool, *exceptions.BaseErrorResponse) {
	var model transactionworkshopentities.ContractService
	var modeloperation []transactionworkshopentities.ContractServiceOperationDetail
	var modelitem []transactionworkshopentities.ContractServiceItemDetail
	err := tx.Model(&model).Where("contract_service_system_number = ?", contractserviceid).Scan(&model).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}
	err2 := tx.Model(&modelitem).Where("contract_service_system_number = ?", contractserviceid).Scan(&modelitem).Error
	if err2 != nil {
		return false, &exceptions.BaseErrorResponse{
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
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	err3 := tx.Model(&modeloperation).Where("contract_service_system_number = ?", id).Scan(&modeloperation).Error
	if err3 !=
	
	nil {
		return false, &exceptions.BaseErrorResponse{
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
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	return true, nil
}

func (r *BookingEstimationImpl) InputDiscount(tx *gorm.DB, id int, req transactionworkshoppayloads.BookEstimationPayloadsDiscount) (bool, *exceptions.BaseErrorResponse) {
	// Update discount_item_percent for different LineTypeIDs in BookingEstimationItemDetail
	err := tx.Model(&transactionworkshopentities.BookingEstimationItemDetail{}).
		Where("estimation_system_number = ? AND line_type_id = ?", id, 1).
		Update("discount_item_percent", req.Accessories).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	err = tx.Model(&transactionworkshopentities.BookingEstimationItemDetail{}).
		Where("estimation_system_number = ? AND line_type_id = ?", id, 3).
		Update("discount_item_percent", req.Material).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	err = tx.Model(&transactionworkshopentities.BookingEstimationItemDetail{}).
		Where("estimation_system_number = ? AND line_type_id = ?", id, 4).
		Update("discount_item_percent", req.Oil).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	err = tx.Model(&transactionworkshopentities.BookingEstimationItemDetail{}).
		Where("estimation_system_number = ? AND line_type_id = ?", id, 7).
		Update("discount_item_percent", req.Souvenir).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	err = tx.Model(&transactionworkshopentities.BookingEstimationItemDetail{}).
		Where("estimation_system_number = ? AND line_type_id = ?", id, 8).
		Update("discount_item_percent", req.Sparepart).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	err = tx.Model(&transactionworkshopentities.BookingEstimationItemDetail{}).
		Where("estimation_system_number = ? AND line_type_id = ?", id, 9).
		Update("discount_item_percent", req.Fee).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	// Update discount_operation_percent for different LineTypeIDs in BookingEstimationOperationDetail
	err = tx.Model(&transactionworkshopentities.BookingEstimationOperationDetail{}).
		Where("estimation_system_number = ? AND line_type_id = ?", id, 5).
		Update("discount_operation_percent", req.Operation).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	err = tx.Model(&transactionworkshopentities.BookingEstimationOperationDetail{}).
		Where("estimation_system_number = ? AND line_type_id = ?", id, 6).
		Update("discount_operation_percent", req.PackageDiscount).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	return true, nil
}

func (r *BookingEstimationImpl) AddFieldAction(tx *gorm.DB, id int, idrecall int)(bool, *exceptions.BaseErrorResponse){
	var model masterentities.FieldAction
	var modeloperation []transactionworkshopentities.ContractServiceOperationDetail
	var modelitem []transactionworkshopentities.ContractServiceItemDetail
	err := tx.Model(&model).Where("field_action_system_number = ?", idrecall).Scan(&model).Error
	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}
	err2 := tx.Model(&modelitem).Where("contract_service_system_number = ?", idrecall).Scan(&modelitem).Error
	if err2 != nil {
		return false, &exceptions.BaseErrorResponse{
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
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	err3 := tx.Model(&modeloperation).Where("contract_service_system_number = ?", id).Scan(&modeloperation).Error
	if err3 !=
	
	nil {
		return false, &exceptions.BaseErrorResponse{
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
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	return true, nil
}