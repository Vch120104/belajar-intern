package transactionworkshoprepositoryimpl

import (
	"after-sales/api/config"
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	"after-sales/api/utils"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type WorkOrderBypassRepositoryImpl struct {
}

func OpenWorkOrderBypassRepositoryImpl() transactionworkshoprepository.WorkOrderBypassRepository {
	return &WorkOrderBypassRepositoryImpl{}
}

func (r *WorkOrderBypassRepositoryImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tableStruct := transactionworkshoppayloads.WorkOrderDetailRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Add the additional where condition
	whereQuery = whereQuery.Where("work_order_system_number > 0 and line_type_id = 1")

	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	var convertedResponses []transactionworkshoppayloads.WorkOrderDetailResponse

	for rows.Next() {

		var (
			workOrderReq transactionworkshoppayloads.WorkOrderDetailRequest
			workOrderRes transactionworkshoppayloads.WorkOrderDetailResponse
		)

		if err := rows.Scan(
			&workOrderReq.WorkOrderDetailId,
			&workOrderReq.WorkOrderSystemNumber,
			&workOrderReq.LineTypeId,
			&workOrderReq.TransactionTypeId,
			&workOrderReq.JobTypeId,
			&workOrderReq.FrtQuantity,
			&workOrderReq.SupplyQuantity,
			&workOrderReq.PriceListId,
			&workOrderReq.WarehouseId,
			&workOrderReq.ItemId,
			&workOrderReq.ProposedPrice,
			&workOrderReq.OperationItemPrice,
		); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch data work order from internal services
		ModelURL := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(workOrderReq.WorkOrderSystemNumber)
		//fmt.Println("Fetching  work order data from:", ModelURL)
		var getModelResponse transactionworkshoppayloads.WorkOrderLookupResponse
		if err := utils.Get(ModelURL, &getModelResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch model data from external service",
				Err:        err,
			}
		}

		// fetch data item from internal services
		ItemURL := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(workOrderReq.ItemId)
		//fmt.Println("Fetching  item data from:", ItemURL)
		var getItemResponse masteritempayloads.BomItemNameResponse
		if err := utils.Get(ItemURL, &getItemResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch item data from external service",
				Err:        err,
			}
		}

		workOrderRes = transactionworkshoppayloads.WorkOrderDetailResponse{
			WorkOrderDocumentNumber: getModelResponse.WorkOrderDocumentNumber,
			WorkOrderSystemNumber:   workOrderReq.WorkOrderSystemNumber,
			LineTypeId:              workOrderReq.LineTypeId,
			ItemId:                  workOrderReq.ItemId,
			ItemName:                getItemResponse.ItemName,
		}

		convertedResponses = append(convertedResponses, workOrderRes)
	}

	var mapResponses []map[string]interface{}

	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"work_order_document_number": response.WorkOrderDocumentNumber,
			"work_order_system_number":   response.WorkOrderSystemNumber,
			"line_type_id":               response.LineTypeId,
			"item_id":                    response.ItemId,
			"item_name":                  response.ItemName,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *WorkOrderBypassRepositoryImpl) GetById(tx *gorm.DB, id int) (transactionworkshoppayloads.WorkOrderBypassResponse, *exceptions.BaseErrorResponse) {
	tableStruct := transactionworkshoppayloads.WorkOrderBypassResponse{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := joinTable.Where("work_order_system_number = ?", id)

	if err := whereQuery.Find(&tableStruct).Error; err != nil {
		return tableStruct, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return tableStruct, nil
}

func (r *WorkOrderBypassRepositoryImpl) Bypass(tx *gorm.DB, request transactionworkshoppayloads.WorkOrderBypassRequestDetail) (transactionworkshopentities.WorkOrderQualityControl, *exceptions.BaseErrorResponse) {
	var entity transactionworkshopentities.WorkOrderQualityControl

	entity.WorkOrderSystemNumber = request.WorkOrderSystemNumber
	entity.WorkOrderQualityControlStatusID = request.WorkOrderQualityControlStatusID
	entity.WorkOrderStartDateTime = request.WorkOrderStartDateTime
	entity.WorkOrderEndDateTime = request.WorkOrderEndDateTime
	entity.WorkOrderActualTime = request.WorkOrderActualTime

	if err := tx.Save(&entity).Error; err != nil {
		return transactionworkshopentities.WorkOrderQualityControl{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entity, nil
}
