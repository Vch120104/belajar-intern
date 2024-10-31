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

type ContractServiceDetailRepositoryImpl struct {
}

func OpenContractServicelDetailRepositoryImpl() transactionworkshoprepository.ContractServiceDetailRepository {
	return &ContractServiceDetailRepositoryImpl{}
}

func (r *ContractServiceDetailRepositoryImpl) GetAllDetail(tx *gorm.DB, Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []transactionworkshopentities.ContractServiceDetail
	combinedPayloads := make([]map[string]interface{}, 0)

	// Initialize the query to filter by ID
	query := tx.Model(&transactionworkshopentities.ContractServiceDetail{}).Where("contract_service_system_number = ?", Id)

	// Apply filter conditions dynamically
	for _, condition := range filterCondition {
		query = query.Where(condition.ColumnField+" = ?", condition.ColumnValue)
	}

	// Execute the query to fetch the data into the entities slice
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	// Prepare the response payloads
	for _, entity := range entities {
		response := map[string]interface{}{
			"contract_service_package_detail_system_number": entity.ContractServicePackageDetailSystemNumber,
			"contract_service_system_number":                entity.ContractServiceSystemNumber,
			"contract_service_line":                         entity.ContractServiceLine,
			"line_type_id":                                  entity.LineTypeId,
			"item_operation_id":                             entity.ItemOperationId,
			"description":                                   entity.Description,
			"frt_quantity":                                  entity.FrtQuantity,
			"item_price":                                    entity.ItemPrice,
			"item_discount_percent":                         entity.ItemDiscountPercent,
			"item_discount_amount":                          entity.ItemDiscountAmount,
			"package_id":                                    entity.PackageId,
			"total_use_frt_quantity":                        entity.TotalUseFrtQuantity,
		}
		combinedPayloads = append(combinedPayloads, response)
	}

	// Handle pagination and return the result
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

// SaveDetail implements transactionworkshoprepository.ContractServiceDetailRepository.
// func (r *ContractServiceDetailRepositoryImpl) SaveDetail(tx *gorm.DB, req transactionworkshoppayloads.ContractServiceIdResponse, Id int) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse) {
// 	entities := transactionworkshopentities.ContractServiceDetail{
// 		ContractServicePackageDetailSystemNumber: req.ContractServicePackageDetailSystemNumber,
// 		ItemOperationId:                          req.ItemOperationId,
// 		ItemDiscountPercent:                      req.ItemDiscountPercent,
// 	}
// 	responses := transactionworkshoppayloads.ContractServiceDetailPayloads{}

// 	lineTypeResponse := transactionworkshoppayloads.LineTypeResponse{}

// 	lineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type/" + strconv.Itoa(entities.LineTypeId)

// 	if err := utils.Get(lineTypeUrl, &lineTypeResponse, nil); err != nil {
// 		return responses, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Err:        err,
// 		}
// 	}

// 	err := tx.Save(&entities).Error

// 	if err != nil {
// 		return responses, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Err:        err,
// 		}
// 	}

// 	responses.ContractServicePackageDetailSystemNumber = entities.ContractServicePackageDetailSystemNumber
// 	responses.ContractServiceSystemNumber = entities.ContractServiceSystemNumber
// 	responses.ContractServiceLine = entities.ContractServiceLine
// 	responses.LineTypeId = entities.LineTypeId
// 	responses.ItemOperationId = entities.ItemOperationId
// 	responses.Description = entities.Description
// 	responses.FrtQuantity = entities.FrtQuantity
// 	responses.ItemPrice = entities.ItemPrice
// 	responses.ItemDiscountPercent = entities.ItemDiscountPercent
// 	responses.ItemDiscountAmount = entities.ItemDiscountAmount
// 	responses.PackageId = entities.PackageId
// 	responses.TotalUseFrtQuantity = entities.TotalUseFrtQuantity

// 	return responses, nil
// }

func (r *ContractServiceDetailRepositoryImpl) SaveDetail(tx *gorm.DB, req transactionworkshoppayloads.ContractServiceIdResponse) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.ContractServiceDetail
	var existingContractService transactionworkshopentities.ContractServiceDetail

	// Validasi Duplikat
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

	// Membuat entitas baru tanpa primary key (auto-increment)
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
		tx.Rollback() // Rollback jika terjadi error
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Mengambil value_after_tax_rate dari trx_contract_service
	var taxRate float64
	err = tx.Table("trx_contract_service").
		Select("value_after_tax_rate").
		Where("contract_service_system_number = ?", req.ContractServiceSystemNumber).
		Scan(&taxRate).Error
	if err != nil {
		tx.Rollback() // Rollback jika error
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Menghitung TOTAL_PRICE berdasarkan FRT_QTY dan OPR_ITEM_PRICE
	var totalPrice float64
	err = tx.Table("trx_contract_service_detail").
		Select("SUM(frt_quantity * item_price * (1 - (item_discount_percent / 100)))").
		Where("contract_service_system_number = ?", req.ContractServiceSystemNumber).
		Scan(&totalPrice).Error
	if err != nil {
		tx.Rollback() // Rollback jika error
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Menghitung TOTAL_VALUE_AFTER_TAX (Total setelah pajak) dan VALUE_AFTER_TAX_RATE (Total pajak)
	totalValueAfterTax := totalPrice * (taxRate / 100)
	totalAfterTax := totalPrice + totalValueAfterTax

	// Update kolom terkait di tabel trx_contract_service
	err = tx.Table("trx_contract_service").
		Where("contract_service_system_number = ?", req.ContractServiceSystemNumber).
		Updates(map[string]interface{}{
			"total":                 totalPrice,
			"value_after_tax_rate":  totalValueAfterTax,
			"total_value_after_tax": totalAfterTax,
		}).Error
	if err != nil {
		tx.Rollback() // Rollback jika error
		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Membuat response untuk data yang berhasil disimpan
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

// func (r *ContractServiceDetailRepositoryImpl) SaveDetail(tx *gorm.DB, req transactionworkshoppayloads.ContractServiceIdResponse) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse) {
// 	var entity transactionworkshopentities.ContractServiceDetail
// 	var existingContractService transactionworkshopentities.ContractServiceDetail

// 	err := tx.Where("contract_service_system_number = ? AND contract_service_line = ?", req.ContractServiceSystemNumber, req.ContractServiceLine).First(&existingContractService).Error

// 	if err == nil {
// 		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusConflict,
// 			Message:    "This Contract Service already exist",
// 		}
// 	} else if err != gorm.ErrRecordNotFound {
// 		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Err:        err,
// 		}
// 	}

// 	entity = transactionworkshopentities.ContractServiceDetail{
// 		// ContractServicePackageDetailSystemNumber: req.ContractServicePackageDetailSystemNumber,
// 		ItemOperationId:                          req.ItemOperationId,
// 		ItemDiscountPercent:                      req.ItemDiscountPercent,
// 		LineTypeId:                               req.LineTypeId, // tetap disimpan tanpa logika tambahan
// 		ContractServiceSystemNumber:              req.ContractServiceSystemNumber,
// 		ContractServiceLine:                      req.ContractServiceLine,
// 		Description:                              req.Description,
// 		FrtQuantity:                              req.FrtQuantity,
// 		ItemPrice:                                req.ItemPrice, // langsung dari input
// 		ItemDiscountAmount:                       req.ItemDiscountAmount,
// 		PackageId:                                req.PackageId,
// 		TotalUseFrtQuantity:                      req.TotalUseFrtQuantity,
// 	}

// 	err = tx.Create(&entity).Error
// 	if err != nil {
// 		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Err:        err,
// 		}
// 	}

// 	var taxRate float64
// 	var totalPrice float64

// 	err = tx.Table("trx_contract_service").Select("value_after_tax_rate").Where("contract_service_system_number = ?", req.ContractServiceSystemNumber).Scan(&taxRate).Error
// 	if err != nil {
// 		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Err:        err,
// 		}
// 	}

// 	// Menghitung TOTAL_PRICE berdasarkan FRT_QTY dan OPR_ITEM_PRICE
// 	err = tx.Table("trx_contract_service_detail").
// 		Select("SUM(frt_quantity * item_price * (1 - (item_discount_percent / 100)))").
// 		Where("contract_service_system_number = ?", req.ContractServiceSystemNumber).
// 		Scan(&totalPrice).Error

// 	if err != nil {
// 		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Err:        err,
// 		}
// 	}

// 	// Menghitung TOTAL_VALUE_AFTER_TAX (Total setelah pajak) dan VALUE_AFTER_TAX_RATE (Total pajak)
// 	totalValueAfterTax := totalPrice * (taxRate / 100)
// 	totalAfterTax := totalPrice + totalValueAfterTax

// 	// Update kolom terkait di tabel trx_contract_service
// 	err = tx.Table("trx_contract_service").
// 		Where("contract_service_system_number = ?", req.ContractServiceSystemNumber).
// 		Updates(map[string]interface{}{
// 			"total":                 totalPrice,
// 			"value_after_tax_rate":  totalValueAfterTax,
// 			"total_value_after_tax": totalAfterTax,
// 		}).Error

// 	if err != nil {
// 		return transactionworkshoppayloads.ContractServiceDetailPayloads{}, &exceptions.BaseErrorResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Err:        err,
// 		}
// 	}

// 	// 4. Mengembalikan Respons
// 	response := transactionworkshoppayloads.ContractServiceDetailPayloads{
// 		ContractServicePackageDetailSystemNumber: entity.ContractServicePackageDetailSystemNumber,
// 		ItemOperationId:                          entity.ItemOperationId,
// 		ItemDiscountPercent:                      entity.ItemDiscountPercent,
// 		LineTypeId:                               entity.LineTypeId,
// 		ContractServiceSystemNumber:              entity.ContractServiceSystemNumber,
// 		ContractServiceLine:                      entity.ContractServiceLine,
// 		Description:                              entity.Description,
// 		FrtQuantity:                              entity.FrtQuantity,
// 		ItemPrice:                                entity.ItemPrice,
// 		ItemDiscountAmount:                       entity.ItemDiscountAmount,
// 		PackageId:                                entity.PackageId,
// 		TotalUseFrtQuantity:                      entity.TotalUseFrtQuantity,
// 	}

// 	return response, nil
// }
