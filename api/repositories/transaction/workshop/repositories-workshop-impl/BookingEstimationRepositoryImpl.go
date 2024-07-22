package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
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

func (r *BookingEstimationImpl) Post(tx *gorm.DB, request transactionworkshoppayloads.BookingEstimationRequest) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.BookingEstimation{
		BrandId:                    request.BrandId,
		ModelId:                    request.ModelId,
		VariantId:                  request.VariantId,
		VehicleId:                  request.VehicleId,
		ContractSystemNumber:       request.ContractSystemNumber,
		AgreementId:                request.AgreementId,
		CampaignId:                 request.CampaignId,
		CompanyId:                  request.CompanyId,
		ProfitCenterId:             request.ProfitCenterId,
		DealerRepresentativeId:     request.DealerRepresentativeId,
		CustomerId:                 request.CustomerId,
		DocumentStatusId:           request.DocumentStatusId,
		BookingEstimationBatchDate: request.BookingEstimationBatchDate,
		IsUnregistered:             request.IsUnregistered,
		InsurancePolicyNo:          request.InsurancePolicyNo,
		InsuranceExpiredDate:       request.InsuranceExpiredDate,
		InsuranceClaimNo:           request.InsuranceClaimNo,
		InsurancePic:               request.InsurancePic,
	} 	
	err := tx.Save(&entities).Error // Menggunakan method Save dari receiver saat ini, yaitu r
	if err != nil {
		return transactionworkshopentities.BookingEstimation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: err}
	}

	return entities, nil
}

func (r *BookingEstimationImpl) GetById(tx *gorm.DB, Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.BookingEstimation
	var brand masterpayloads.BrandResponse
	var model masterpayloads.GetModelResponse
	var vehicle transactionworkshoppayloads.VehicleDetailPayloads
	err := tx.Model("trx_booking_estimation.*,mtr_contract_sevice.contract_service_number,mtr_contract_service.contract_service_date,mtr_contract_service.contract_service_dealer,mtr_agreement.agreement_code,mtr_agreement.agreement_date_to,mtr_agreement.agreement_dealer").
	Joins("JOIN mtr_contract_service on mtr_contract_service.vehicle_id = trx_booking_estimation.vehicle_id").
	Joins("Join mtr_agreement on mtr_agreement.customer_id = trx_booking_estimation.customer_id").Where("batch_system_number= ?", Id).First(&entity).Error
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{Message: "Failed to retrieve booking estimation from the database"}
	}
	errUrlBrand := utils.Get(config.EnvConfigs.SalesServiceUrl+"unit-brand/"+strconv.Itoa(entity.BrandId),brand,nil)
	if errUrlBrand != nil{
		return nil,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: errors.New("brand undefined"),
		}
	}
	joinedData1:= utils.DataFrameInnerJoin([]transactionworkshopentities.BookingEstimation{entity},[]masterpayloads.BrandResponse{brand},"BrandId")
	errUrlModel := utils.Get(config.EnvConfigs.SalesServiceUrl+"unit-model/"+strconv.Itoa(entity.ModelId),model,nil)
	if errUrlModel != nil{
		return nil,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: errors.New("model undefined"),
		}
	}
	joinedData2:= utils.DataFrameInnerJoin(joinedData1,[]masterpayloads.GetModelResponse{model},"ModelId")
	errurlVehicle:= utils.Get(config.EnvConfigs.SalesServiceUrl+"vehicle-master/"+strconv.Itoa(entity.VehicleId),&vehicle,nil)
	if errurlVehicle != nil{
		return nil,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: errors.New("vehicle undefined"),
		}
	}
	joinedData3 := utils.DataFrameInnerJoin(joinedData2,[]transactionworkshoppayloads.VehicleDetailPayloads{vehicle},"VehicleId")

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

func (r *BookingEstimationImpl) Submit(tx *gorm.DB, Id int) (bool,*exceptions.BaseErrorResponse) {
	// Retrieve the booking estimation by Id
	var entity transactionworkshopentities.BookingEstimation
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return false,&exceptions.BaseErrorResponse{Message: "Failed to retrieve booking estimation from the database"}
	}

	// Perform the necessary operations to submit the booking estimation
	// ...

	// Save the updated booking estimation
	err = tx.Save(&entity).Error
	if err != nil {
		return false,&exceptions.BaseErrorResponse{Message: "Failed to save the updated booking estimation"}
	}

	return true,nil
}

func (r *BookingEstimationImpl) Void(tx *gorm.DB, Id int) (bool,*exceptions.BaseErrorResponse) {
	// Retrieve the booking estimation by Id
	var entity transactionworkshopentities.BookingEstimation
	err := tx.Model(&transactionworkshopentities.BookingEstimation{}).Where("id = ?", Id).First(&entity).Error
	if err != nil {
		return false,&exceptions.BaseErrorResponse{Message: "Failed to retrieve booking estimation from the database"}
	}

	// Perform the necessary operations to void the booking estimation
	// ...

	// Save the updated booking estimation
	err = tx.Delete(&entity).Error
	if err != nil {
		return false,&exceptions.BaseErrorResponse{Message: "Failed to save the updated booking estimation"}
	}

	return true,nil
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

func (r *BookingEstimationImpl) SaveBookEstimReq(tx *gorm.DB, req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (int, *exceptions.BaseErrorResponse) {
	var count int64
	err := tx.Select("trx_booking_estim.*").Where("booking_system_number =?", id).Count(&count).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
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
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
	return id, nil
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
            if gorm.IsRecordNotFoundError(err) {
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
	err := tx.Model(&model).Where("booking_estimation_request_id = ?", id).Scan(&payloads).Error
	if err != nil {
		return payloads, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
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
	
	err:= tx.Select(&bookestimcalc).Where("batch_system_number = ?", req.EstimationSystemNumber).Scan(&bookestimpayloads).Error
	if err != nil{
		return 0,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: errors.New("Discount can'r be created"),
		}
	}
	
	if err != nil{
		return 0,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err,
		}
	}
	return req.Book, nil
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

func (r *BookingEstimationImpl) AddPackage(tx *gorm.DB, id int, packId int) (int, *exceptions.BaseErrorResponse) {
	var modeloperation masterpackagemasterentity.PackageMasterDetailOperation
	var modelitem masterpackagemasterentity.PackageMasterDetailItem
	var operationpayloads []masterpayloads.CampaignMasterDetailOperationPayloads
	var itempayloads []masterpayloads.PackageMasterDetailItem
	err2 := tx.Model(&modelitem).Where("package_id = ?", packId).Scan(&itempayloads).Error
	if err2 != nil {
		return 0, &exceptions.BaseErrorResponse{
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
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	err3 := tx.Model(&modeloperation).Where("estimation_system_number = ?", id).Scan(&operationpayloads).Error
	if err3 != nil {
		return 0, &exceptions.BaseErrorResponse{
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
			return 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        err,
			}
		}
	}
	return id, nil
}

func (r *BookingEstimationImpl) AddContractService(tx *gorm.DB, id int, contractserviceid int) (int, *exceptions.BaseErrorResponse) {
	var model transactionworkshopentities.ContractService
	var modeloperation []transactionworkshopentities.ContractServiceOperationDetail
	var modelitem []transactionworkshopentities.ContractServiceItemDetail
	err := tx.Model(&model).Where("contract_service_system_number = ?", contractserviceid).Scan(&model).Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}
	err2 := tx.Model(&modelitem).Where("contract_service_system_number = ?", contractserviceid).Scan(&modelitem).Error
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

func (r *BookingEstimationImpl) GetByIdBookEstimDetail (tx *gorm.DB ,id int ,LineTypeID int)(map[string]interface{},*exceptions.BaseErrorResponse){
	var payloadsoperation []transactionworkshoppayloads.BookEstimDetailPayloadsOperation
	var payloadsitem []transactionworkshoppayloads.BookEstimDetailPayloadsItem
	var payloadslinetype []masterpayloads.LineTypeCode
	var payloadstransactiontype []transactionworkshoppayloads.TransactionTypePayloads

	if LineTypeID == 5{
		err := tx.Select("trx_booking_estimation_operation_detail.*,mtr_operation_code.operation_code").Where("estimation_line_id = ?",id).
		Joins("join mtr_operation_model_mapping on mtr_operation_model_mapping.operation_model_mapping_id = trx_booking_estimation_operation_detail.operation_id").
		Joins("join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id").Scan(&payloadsoperation).Error
		if err != nil{
			return nil,&exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err: err,
			}
		}
		errurllinetype := utils.Get(config.EnvConfigs.GeneralServiceUrl + "/line-type",&payloadslinetype,nil)
		if errurllinetype != nil{
			return nil,&exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err: errurllinetype,
			}
		} 
		joinedData1 := utils.DataFrameInnerJoin(payloadsoperation,payloadslinetype,"LineTypeId")
		errurltransactiontype := utils.Get(config.EnvConfigs.GeneralServiceUrl+ "transaction-type-list?page=0&limit=100000",&payloadstransactiontype,nil)
		if errurltransactiontype != nil{
			return nil,&exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err: errurltransactiontype,
			}
		}
		joineddata2 := utils.DataFrameInnerJoin(joinedData1,payloadstransactiontype,"TransactionTypeId")
		return joineddata2[0],nil

	}else{
		err:= tx.Model("tr_booking_estimation_item_detail.*,mtr_item.item_name").Where("estimation_line_id = ?",id).
		Joins("Join mtr_item on mtr_item.item_id = trx_booking_estimation_item_detail.item_id").
		Scan(&payloadsitem).Error
		if err != nil{
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err: err,
			}
		}
		errurllinetype := utils.Get(config.EnvConfigs.GeneralServiceUrl + "/line-type",&payloadslinetype,nil)
		if errurllinetype != nil{
			return nil,&exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err: errurllinetype,
			}
		} 
		joinedData1 := utils.DataFrameInnerJoin(payloadsitem,payloadslinetype,"LineTypeId")
		errurltransactiontype := utils.Get(config.EnvConfigs.GeneralServiceUrl+ "transaction-type-list?page=0&limit=100000",&payloadstransactiontype,nil)
		if errurltransactiontype != nil{
			return nil,&exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err: errurltransactiontype,
			}
		}
		joineddata2 := utils.DataFrameInnerJoin(joinedData1,payloadstransactiontype,"TransactionTypeId")
		return joineddata2[0],nil
	}
}

func (r *BookingEstimationImpl) GetAllBookEstimDetail (tx *gorm.DB,id int, pages pagination.Pagination)([]map[string]interface{},*exceptions.BaseErrorResponse){
	var operationpayloads []transactionworkshoppayloads.BookEstimDetailPayloadsOperation
	var itempayloads []transactionworkshoppayloads.BookEstimDetailPayloadsItem
	var payloadslinetype []masterpayloads.LineTypeCode
	var payloadstransactiontype []transactionworkshoppayloads.TransactionTypePayloads
	combinedpayloads := make([]map[string]interface{},0)
	
	err := tx.Select("trx_booking_estimation_operation_detail.*,mtr_operation_code.operation_code").Where("estimation_line_id = ?",id).
	Joins("join mtr_operation_model_mapping on mtr_operation_model_mapping.operation_model_mapping_id = trx_booking_estimation_operation_detail.operation_id").
	Joins("join mtr_operation_code on mtr_operation_code.operation_id=mtr_operation_model_mapping.operation_id").Scan(&operationpayloads).Error
	if err != nil{
		return nil,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: err,
		}
	}
	errurllinetype := utils.Get(config.EnvConfigs.GeneralServiceUrl + "/line-type",&payloadslinetype,nil)
	if errurllinetype != nil{
		return nil,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: errurllinetype,
		}
	} 
	joinedData1 := utils.DataFrameInnerJoin(operationpayloads,payloadslinetype,"LineTypeId")
	errurltransactiontype := utils.Get(config.EnvConfigs.GeneralServiceUrl+ "transaction-type-list?page=0&limit=100000",&payloadstransactiontype,nil)
	if errurltransactiontype != nil{
		return nil,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: errurltransactiontype,
		}
	}
	joineddata2 := utils.DataFrameInnerJoin(joinedData1,payloadstransactiontype,"TransactionTypeId")
	for _,op := range joineddata2{
		combinedpayloads = append(combinedpayloads,map[string]interface{}{
			"line_type_id": op["LineTypeid"],
			"transaction_type_id": op["TransactionTypeId"],
			"operation_id": op["OperationId"],
			"operation_name":op["OperationName"],
			"quantity": op["Quantity"],
			"price":op["Price"],
			"subtotal": op["SubTotal"],
			"original_discount":op["OriginalDiscount"],
			"proposal_discount": op["ProposalDiscount"],
			"total" :op["Total"],
		})
	}


	err2:= tx.Model("tr_booking_estimation_item_detail.*,mtr_item.item_name").Where("estimation_line_id = ?",id).
	Joins("Join mtr_item on mtr_item.item_id = trx_booking_estimation_item_detail.item_id").
	Scan(&itempayloads).Error
	if err2 != nil{
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: err,
		}
	}
	errurllinetype2 := utils.Get(config.EnvConfigs.GeneralServiceUrl + "/line-type",&payloadslinetype,nil)
	if errurllinetype2 != nil{
		return nil,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: errurllinetype2,
		}
	} 
	joinedData1_2 := utils.DataFrameInnerJoin(itempayloads,payloadslinetype,"LineTypeId")
	errurltransactiontype2 := utils.Get(config.EnvConfigs.GeneralServiceUrl+ "transaction-type-list?page=0&limit=100000",&payloadstransactiontype,nil)
	if errurltransactiontype2 != nil{
		return nil,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err: errurltransactiontype,
		}
	}
	joineddata2_2 := utils.DataFrameInnerJoin(joinedData1_2,payloadstransactiontype,"TransactionTypeId")
	for _,it := range joineddata2_2{
		combinedpayloads = append(combinedpayloads,map[string]interface{}{
			"line_type_id": it["LineTypeID"],
			"transaction_type_id": it["TransactionTypeId"],
			"item_id": it["ItemID"],
			"item_name":it["ItemName"],
			"quantity": it["Quantity"],
			"price":it["Price"],
			"subtotal": it["SubTotal"],
			"original_discount":it["OriginalDiscount"],
			"proposal_discount": it["ProposalDiscount"],
			"total" :it["Total"],
		})
	}
	return combinedpayloads,nil
}

func (r *BookingEstimationImpl) PostBookingEstimationCalculation(tx*gorm.DB,id int)(int,*exceptions.BaseErrorResponse){
	now := time.Now()
	entity := transactionworkshopentities.BookingEstimationServiceDiscount{
		BatchSystemNumber:               id,
		DocumentStatusID:                0, 
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
	if err != nil{
		return 0,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err: err,
		}
	}
	return id,nil
}

func (r *BookingEstimationImpl) PutBookingEstimationCalculation (tx *gorm.DB, id int, linetypeid int, req transactionworkshoppayloads.BookingEstimationCalculationPayloads)(int,*exceptions.BaseErrorResponse){
	var entity transactionworkshopentities.BookingEstimationServiceDiscount
	var columnName string
	var value int
	err := tx.Model(&entity).Where("batch_system_number =?",id).Scan(&entity).Error
	if err != nil{
		return 0,&exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err: err,
		}
	}

	switch linetypeid {
	case 5:
		columnName = "total_price_operation"
		value = int(req.TotalPriceOperation)
	case 6:
		columnName = "total_price_package"
		value = int(req.TotalPricePackage)
	case 1:
		columnName = "total_price_accessories"
		value = int(req.TotalPriceAccessories)
	case 2:
		columnName = "total_price_part"
		value = int(req.TotalPricePart)
	case 3:
		columnName = "total_price_material"
		value = int(req.TotalPriceMaterial)
	case 4:
		columnName = "total_price_oil"																																																										
		value = int(req.TotalPriceOil)
	case 9:
		columnName = "total_sublet"
		value = int(req.TotalSublet)



	err := tx.Model(&transactionworkshopentities.BookingEstimationOperationDetail{}).
		Where("batch_system_number = ? ", id).
		Update(columnName, value).
		Update("total", gorm.Expr("total_price_operation + total_price_package + total_price_accessories + total_price_part + total_price_material + total_price_oil + total_sublet")).
		Error
	if err != nil {
		return 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	}
	
	return id, nil
}

func (r *BookingEstimationImpl) SaveBookingEstimationFromPDI (tx *gorm.DB,id int)(transactionworkshopentities.BookingEstimation,*exceptions.BaseErrorResponse){
	entities := transactionworkshopentities.BookingEstimation{}
	
}