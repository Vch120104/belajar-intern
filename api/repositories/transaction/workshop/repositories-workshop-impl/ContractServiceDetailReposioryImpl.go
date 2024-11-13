package transactionworkshoprepositoryimpl

import (
	masteroperationentities "after-sales/api/entities/master/operation"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"net/http"

	"gorm.io/gorm"
)

type ContractServiceDetailRepositoryImpl struct {
}

func OpenContractServicelDetailRepositoryImpl() transactionworkshoprepository.ContractServiceDetailRepository {
	return &ContractServiceDetailRepositoryImpl{}
}

func (r *ContractServiceDetailRepositoryImpl) GetAllDetail(tx *gorm.DB, Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.ContractServiceDetail
	combinedPayloads := make([]map[string]interface{}, 0)

	// Query utama untuk mengambil data ContractServiceDetail
	query := tx.Model(&transactionworkshopentities.ContractServiceDetail{}).
		Where("contract_service_system_number = ?", Id)

	// Menambahkan filter condition ke query
	for _, condition := range filterCondition {
		query = query.Where(condition.ColumnField+" = ?", condition.ColumnValue)
	}

	// Eksekusi query untuk mendapatkan ContractServiceDetail
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	// Iterasi melalui hasil query dan melakukan query tambahan untuk Line Type dan Operation Code
	for _, entity := range entities {
		// Mengambil Line Type dari API eksternal
		linetype, linetypeErr := generalserviceapiutils.GetLineTypeById(entity.LineTypeId)
		if linetypeErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        linetypeErr.Err,
			}
		}

		// Mengambil Operation Code dari tabel mtr_operation_code menggunakan GORM
		var operation masteroperationentities.OperationCode
		operationErr := tx.Model(&masteroperationentities.OperationCode{}).
			Where("operation_id = ?", entity.ItemOperationId).
			First(&operation).Error
		if operationErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        operationErr,
			}
		}

		// Membuat response untuk setiap entity detail
		response := map[string]interface{}{
			"contract_service_package_detail_system_number": entity.ContractServicePackageDetailSystemNumber,
			"contract_service_system_number":                entity.ContractServiceSystemNumber,
			"contract_service_line":                         entity.ContractServiceLine,
			"line_type_id":                                  entity.LineTypeId,
			"line_type_code":                                linetype.LineTypeCode,
			"item_operation_id":                             entity.ItemOperationId,
			"operation_code":                                operation.OperationCode,
			"operation_name":                                operation.OperationName,
			"description":                                   entity.Description,
			"frt_quantity":                                  entity.FrtQuantity,
			"item_price":                                    entity.ItemPrice,
			"item_discount_percent":                         entity.ItemDiscountPercent,
			"item_discount_amount":                          entity.ItemDiscountAmount,
			"package_id":                                    entity.PackageId,
			"total_use_frt_quantity":                        entity.TotalUseFrtQuantity,
		}

		// Menambahkan hasil ke payload kombinasi
		combinedPayloads = append(combinedPayloads, response)
	}

	// Pagination hasil data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(combinedPayloads, &pages)
	return paginatedData, totalPages, totalRows, nil
}

// GetById implements transactionworkshoprepository.ContractServiceDetailRepository.
func (r *ContractServiceDetailRepositoryImpl) GetById(tx *gorm.DB, Id int) (transactionworkshoppayloads.ContractServiceIdResponse, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.ContractServiceDetail{}
	responses := transactionworkshoppayloads.ContractServiceIdResponse{}

	err := tx.Model(&entities).Where(transactionworkshopentities.ContractServiceDetail{ContractServicePackageDetailSystemNumber: Id}).First(&responses).Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	return responses, nil
}

func (r *ContractServiceDetailRepositoryImpl) SaveDetail(tx *gorm.DB, req transactionworkshoppayloads.ContractServiceIdResponse) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ContractServiceDetail
	var existingContractService transactionworkshopentities.ContractServiceDetail

	err := tx.Where("contract_service_system_number = ? AND contract_service_line = ?", req.ContractServiceSystemNumber, req.ContractServiceLine).
		First(&existingContractService).Error
	if err == nil {
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "This Contract Service already exists",
		}
	} else if err != gorm.ErrRecordNotFound {
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	entity = transactionworkshopentities.ContractServiceDetail{
		ItemOperationId:             req.ItemOperationId,
		ItemDiscountPercent:         req.ItemDiscountPercent,
		LineTypeId:                  req.LineTypeId,
		ContractServiceSystemNumber: req.ContractServiceSystemNumber,
		ContractServiceLine:         req.ContractServiceLine,
		Description:                 req.Description,
		FrtQuantity:                 req.FrtQuantity,
		ItemPrice:                   req.ItemPrice,
		ItemDiscountAmount:          req.ItemDiscountAmount,
		PackageId:                   req.PackageId,
		TotalUseFrtQuantity:         req.TotalUseFrtQuantity,
	}

	err = tx.Create(&entity).Error
	if err != nil {
		tx.Rollback()
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var taxRate float64
	err = tx.Table("trx_contract_service").
		Select("value_after_tax_rate").
		Where("contract_service_system_number = ?", req.ContractServiceSystemNumber).
		Scan(&taxRate).Error
	if err != nil {
		tx.Rollback()
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var totalPrice float64
	err = tx.Table("trx_contract_service_detail").
		Select("SUM(frt_quantity * item_price * (1 - (item_discount_percent / 100)))").
		Where("contract_service_system_number = ?", req.ContractServiceSystemNumber).
		Scan(&totalPrice).Error
	if err != nil {
		tx.Rollback()
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	totalValueAfterTax := totalPrice * (taxRate / 100)
	totalAfterTax := totalPrice + totalValueAfterTax

	err = tx.Table("trx_contract_service").
		Where("contract_service_system_number = ?", req.ContractServiceSystemNumber).
		Updates(map[string]interface{}{
			"total":                 totalPrice,
			"value_after_tax_rate":  totalValueAfterTax,
			"total_value_after_tax": totalAfterTax,
		}).Error
	if err != nil {
		tx.Rollback()
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	response := transactionworkshoppayloads.ContractServiceDetailPayloads{
		ContractServicePackageDetailSystemNumber: entity.ContractServicePackageDetailSystemNumber,
		ItemOperationId:                          entity.ItemOperationId,
		ItemDiscountPercent:                      entity.ItemDiscountPercent,
		LineTypeId:                               entity.LineTypeId,
		ContractServiceSystemNumber:              entity.ContractServiceSystemNumber,
		ContractServiceLine:                      entity.ContractServiceLine,
		Description:                              entity.Description,
		FrtQuantity:                              entity.FrtQuantity,
		ItemPrice:                                entity.ItemPrice,
		ItemDiscountAmount:                       entity.ItemDiscountAmount,
		PackageId:                                entity.PackageId,
		TotalUseFrtQuantity:                      entity.TotalUseFrtQuantity,
	}

	return response, nil
}

// UpdateDetail implements transactionworkshoprepository.ContractServiceDetailRepository.
func (r *ContractServiceDetailRepositoryImpl) UpdateDetail(tx *gorm.DB, contractServiceSystemNumber int, contractServiceLine string, req transactionworkshoppayloads.ContractServiceDetailRequest) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse) {
	var existingDetail transactionworkshopentities.ContractServiceDetail

	// Cari data detail berdasarkan `contractServiceSystemNumber` dan `contractServiceLine`
	err := tx.Model(&transactionworkshopentities.ContractServiceDetail{}).
		Where("contract_service_system_number = ? AND contract_service_line = ?", contractServiceSystemNumber, contractServiceLine).
		First(&existingDetail).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Contract Service Detail not found",
			}
		}
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Update `item_discount_percent` dan `item_discount_amount`
	err = tx.Model(&existingDetail).Updates(map[string]interface{}{
		"item_discount_percent": req.ItemDiscountPercent,
		"item_discount_amount":  req.ItemDiscountAmount,
	}).Error
	if err != nil {
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Perhitungan `total`, `total_value_after_tax`, dan `total_after_tax`
	var taxRate float64
	err = tx.Table("trx_contract_service").
		Select("value_after_tax_rate").
		Where("contract_service_system_number = ?", contractServiceSystemNumber).
		Scan(&taxRate).Error
	if err != nil {
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var totalPrice float64
	err = tx.Table("trx_contract_service_detail").
		Select("SUM(frt_quantity * item_price * (1 - (item_discount_percent / 100)))").
		Where("contract_service_system_number = ?", contractServiceSystemNumber).
		Scan(&totalPrice).Error
	if err != nil {
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	totalValueAfterTax := totalPrice * (taxRate / 100)
	totalAfterTax := totalPrice + totalValueAfterTax

	// Update tabel `trx_contract_service` dengan hasil perhitungan
	err = tx.Model(&transactionworkshopentities.ContractService{}).
		Where("contract_service_system_number = ?", contractServiceSystemNumber).
		Updates(map[string]interface{}{
			"total":                 totalPrice,
			"value_after_tax_rate":  totalValueAfterTax,
			"total_value_after_tax": totalAfterTax,
		}).Error
	if err != nil {
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Siapkan response
	response := transactionworkshoppayloads.ContractServiceDetailPayloads{
		ContractServicePackageDetailSystemNumber: existingDetail.ContractServicePackageDetailSystemNumber,
		ItemOperationId:                          existingDetail.ItemOperationId,
		ItemDiscountPercent:                      req.ItemDiscountPercent,
		LineTypeId:                               existingDetail.LineTypeId,
		ContractServiceSystemNumber:              existingDetail.ContractServiceSystemNumber,
		ContractServiceLine:                      existingDetail.ContractServiceLine,
		Description:                              existingDetail.Description,
		FrtQuantity:                              existingDetail.FrtQuantity,
		ItemPrice:                                existingDetail.ItemPrice,
		ItemDiscountAmount:                       req.ItemDiscountAmount,
		PackageId:                                existingDetail.PackageId,
		TotalUseFrtQuantity:                      existingDetail.TotalUseFrtQuantity,
	}

	return response, nil
}
