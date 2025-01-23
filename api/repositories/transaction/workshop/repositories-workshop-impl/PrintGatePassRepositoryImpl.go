package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type PrintGatePassRepositoryImpl struct {
}

func OpenPrintGatePassRepositoryImpl() transactionworkshoprepository.PrintGatePassRepository {
	return &PrintGatePassRepositoryImpl{}
}

// GetAll implements transactionworkshoprepository.PrintGatePassRepository.
func (p *PrintGatePassRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var gatePasses []transactionworkshopentities.PrintGatePass

	baseQuery := tx.Model(&transactionworkshopentities.PrintGatePass{})
	whereQuery := utils.ApplyFilter(baseQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&gatePasses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(gatePasses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, gatePass := range gatePasses {
		workOrderResponse, errWorkOrder := p.getWorkOrderData(gatePass.CompanyId, gatePass.BpkSystemNumber)
		if errWorkOrder != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: errWorkOrder.StatusCode,
				Err:        errWorkOrder.Err,
			}
		}

		if workOrderResponse == nil {
			continue
		}

		result := map[string]interface{}{
			"wo_sys_no":        gatePass.BpkSystemNumber,
			"wo_doc_no":        workOrderResponse.WorkOrderDocumentNumber,
			"wo_date":          workOrderResponse.WorkOrderDate,
			"customer_name":    workOrderResponse.NameCust,
			"vehicle_id":       gatePass.VehicleId,
			"vehicle_brand_id": gatePass.VehicleBrandId,
			"model_id":         gatePass.ModelId,
			"gate_pass_sys_no": gatePass.GatePassSystemNumber,
			"gate_pass_doc_no": gatePass.GatePassDocumentNumber,
			"gate_pass_date":   gatePass.GatePassDate,
			"delivery_name":    gatePass.DeliveryName,
			"delivery_address": gatePass.DeliveryAddress,
		}
		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
}

func (p *PrintGatePassRepositoryImpl) getWorkOrderData(companyId int, workOrderSystemNumber int) (*transactionworkshoppayloads.WorkOrderResponse, *exceptions.BaseErrorResponse) {
	workOrderURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(workOrderSystemNumber)
	var workOrderResponse transactionworkshoppayloads.WorkOrderResponse

	err := utils.Get(workOrderURL, &workOrderResponse, nil)
	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve work order data from the external API",
			Err:        err,
		}
	}

	if workOrderResponse.WorkOrderDocumentNumber == "" {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Work Order document number not found",
		}
	}

	return &workOrderResponse, nil
}
