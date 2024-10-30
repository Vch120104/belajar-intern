package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"net/http"
	"strconv"

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
func (r *ContractServiceDetailRepositoryImpl) SaveDetail(tx *gorm.DB, req transactionworkshoppayloads.ContractServiceIdResponse) (transactionworkshoppayloads.ContractServiceDetailPayloads, *exceptions.BaseErrorResponse) {
	entities := transactionworkshopentities.ContractServiceDetail{
		ContractServicePackageDetailSystemNumber: req.ContractServicePackageDetailSystemNumber,
		ItemOperationId:                          req.ItemOperationId,
		ItemDiscountPercent:                      req.ItemDiscountPercent,
	}
	responses := transactionworkshoppayloads.ContractServiceDetailPayloads{}

	lineTypeResponse := transactionworkshoppayloads.LineTypeResponse{}

	lineTypeUrl := config.EnvConfigs.GeneralServiceUrl + "line-type/" + strconv.Itoa(entities.LineTypeId)

	if err := utils.Get(lineTypeUrl, &lineTypeResponse, nil); err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return responses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	responses.ContractServicePackageDetailSystemNumber = entities.ContractServicePackageDetailSystemNumber
	responses.ContractServiceSystemNumber = entities.ContractServiceSystemNumber
	responses.ContractServiceLine = entities.ContractServiceLine
	responses.LineTypeId = entities.LineTypeId
	responses.ItemOperationId = entities.ItemOperationId
	responses.Description = entities.Description
	responses.FrtQuantity = entities.FrtQuantity
	responses.ItemPrice = entities.ItemPrice
	responses.ItemDiscountPercent = entities.ItemDiscountPercent
	responses.ItemDiscountAmount = entities.ItemDiscountAmount
	responses.PackageId = entities.PackageId
	responses.TotalUseFrtQuantity = entities.TotalUseFrtQuantity

	return responses, nil
}
