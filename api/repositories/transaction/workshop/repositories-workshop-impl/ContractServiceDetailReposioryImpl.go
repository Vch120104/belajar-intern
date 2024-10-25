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
	var operationPayload transactionworkshoppayloads.Operation
	combinedPayloads := make([]map[string]interface{}, 0)

	// Initialize the query with the base table
	query := tx.Model(&transactionworkshopentities.ContractServiceDetail{}).Where("contract_service_system_number = ?", Id)

	// Apply filter conditions
	for _, condition := range filterCondition {
		query = query.Where(condition.ColumnField+" = ?", condition.ColumnValue)
	}

	// Execute the query and get the result into entities
	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	// Prepare the response payload
	for _, entity := range entities {
		var operationName, operationCode string

		// Fetch operation details using item_operation_id
		err := tx.Select("mtr_operation_code.operation_name, mtr_operation_code.operation_code").
			Joins("join mtr_item_operation on mtr_item_operation.item_operation_id = ?", entity.ItemOperationId).
			Joins("JOIN mtr_operation_model_mapping ON mtr_operation_model_mapping.operation_model_mapping_id = mtr_item_operation.item_operation_model_mapping_id").
			Joins("join mtr_operation_code on mtr_operation_code.operation_id = mtr_operation_model_mapping.operation_id").
			Table("mtr_item_operation").
			Scan(&operationPayload).Error

		if err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		operationName = operationPayload.OperationName
		operationCode = operationPayload.OperationCode

		// Construct response map
		response := map[string]interface{}{
			"contract_service_package_detail_system_number": entity.ContractServicePackageDetailSystemNumber,
			"contract_service_system_number":                entity.ContractServiceSystemNumber,
			"contract_service_line":                         entity.ContractServiceLine,
			"line_type_id":                                  entity.LineTypeId,
			"item_operation_id":                             entity.ItemOperationId,
			"frt_quantity":                                  entity.FrtQuantity,
			"item_price":                                    entity.ItemPrice,
			"item_discount_percent":                         entity.ItemDiscountPercent,
			"item_discount_amount":                          entity.ItemDiscountAmount,
			"package_id":                                    entity.PackageId,
			"total_use_frt_quantity":                        entity.TotalUseFrtQuantity,
			"operation_name":                                operationName,
			"operation_code":                                operationCode,
		}

		combinedPayloads = append(combinedPayloads, response)
	}

	// Pagination and return
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(combinedPayloads, &pages)
	return paginatedData, totalPages, totalRows, nil
}
