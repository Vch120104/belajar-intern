package transactionworkshoprepositoryimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type PrintGatePassRepositoryImpl struct {
}

func OpenPrintGatePassRepositoryImpl() transactionworkshoprepository.PrintGatePassRepository {
	return &PrintGatePassRepositoryImpl{}
}

// GetAll implements transactionworkshoprepository.PrintGatePassRepository.
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
			"trx_gate_pass.company_id",
			"trx_work_order.work_order_system_number",
			"trx_work_order.work_order_document_number",
			"trx_work_order.work_order_date",
		).Joins("LEFT JOIN trx_work_order ON trx_gate_pass.company_id = trx_work_order.company_id")
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var selectedWOs []transactionworkshoppayloads.PrintGatePassResponse
	for _, response := range responses {
		selectedWOs = append(selectedWOs, transactionworkshoppayloads.PrintGatePassResponse{
			WorkOrderSystemNumber:   response.WorkOrderSystemNumber,
			WorkOrderDocumentNumber: response.WorkOrderDocumentNumber,
			WorkOrderDate:           response.WorkOrderDate,
			CustomerId:              response.CustomerId,
			CustomerName:            response.CustomerName,
			VehicleId:               response.VehicleId,
			VehicleBrandId:          response.VehicleBrandId,
			ModelId:                 response.ModelId,
			GatePassSystemNumber:    response.GatePassSystemNumber,
			GatePassDocumentNumber:  response.GatePassDocumentNumber,
			GatePassDate:            response.GatePassDate,
			DeliveryName:            response.DeliveryName,
			DeliveryAddress:         response.DeliveryAddress,
		})
	}

	if len(selectedWOs) > 1 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Cannot print multiple WO. Please select only one WO",
		}
	}

	gatePassDate := responses[0].GatePassDate
	currentDate := time.Now().Format("2006-01-02")

	if gatePassDate != currentDate {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Gate Pass form only can be printed sameday as gate pass date (" + gatePassDate + ")",
		}
	}
	pages.Rows = responses

	return pages, nil
}
