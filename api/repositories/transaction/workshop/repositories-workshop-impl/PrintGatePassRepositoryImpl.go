package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"net/http"

	"gorm.io/gorm"
)

type PrintGatePassRepositoryImpl struct {
}

func OpenPrintGatePassRepositoryImpl() transactionworkshoprepository.PrintGatePassRepository {
	return &PrintGatePassRepositoryImpl{}
}

func (p *PrintGatePassRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities transactionworkshopentities.PrintGatePass
	var responses []transactionworkshoppayloads.PrintGatePassResponse

	baseModelQuery := tx.Model(&entities).
		Select(
			"trx_gate_pass.gate_pass_system_number",
			"trx_gate_pass.gate_pass_document_number",
			"trx_gate_pass.gate_pass_date",
			"trx_gate_pass.delivery_name",
			"trx_gate_pass.delivery_address",
			"trx_gate_pass.customer_id",
			"trx_gate_pass.vehicle_id",
			"trx_gate_pass.vehicle_brand_id",
			"trx_gate_pass.model_id",
			"trx_work_order.work_order_system_number",
			"trx_work_order.work_order_document_number",
			"trx_work_order.work_order_date",
		).Joins("LEFT JOIN trx_work_order ON trx_gate_pass.company_id = trx_work_order.company_id")

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch gate pass data",
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var enrichedResponses []map[string]interface{}
	for _, response := range responses {
		// Get customer data from external API
		customerData, customerErr := generalserviceapiutils.GetCustomerMasterById(response.CustomerId)
		if customerErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve customer data from external API",
				Err:        customerErr.Err,
			}
		}

		// Map the data including customer information
		responseMap := map[string]interface{}{
			"work_order_system_number":   response.WorkOrderSystemNumber,
			"work_order_document_number": response.WorkOrderDocumentNumber,
			"work_order_date":            response.WorkOrderDate,
			"customer_id":                response.CustomerId,
			"customer_name":              customerData.CustomerName,
			"customer_code":              customerData.CustomerCode,
			"vehicle_id":                 response.VehicleId,
			"vehicle_brand_id":           response.VehicleBrandId,
			"model_id":                   response.ModelId,
			"gate_pass_system_number":    response.GatePassSystemNumber,
			"gate_pass_document_number":  response.GatePassDocumentNumber,
			"gate_pass_date":             response.GatePassDate,
			"delivery_name":              response.DeliveryName,
			"delivery_address":           response.DeliveryAddress,
		}
		enrichedResponses = append(enrichedResponses, responseMap)
	}

	// Handle pagination
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(enrichedResponses, &pages)

	pages.Rows = paginatedData
	pages.TotalRows = int64(totalRows)
	pages.TotalPages = totalPages

	return pages, nil
}
