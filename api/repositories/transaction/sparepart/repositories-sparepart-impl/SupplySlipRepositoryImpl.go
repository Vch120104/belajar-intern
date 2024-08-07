package transactionsparepartrepositoryimpl

import (
	"after-sales/api/config"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"gorm.io/gorm"
)

type SupplySlipRepositoryImpl struct {
}

func StartSupplySlipRepositoryImpl() transactionsparepartrepository.SupplySlipRepository {
	return &SupplySlipRepositoryImpl{}
}

func (r *SupplySlipRepositoryImpl) GetSupplySlipById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipResponse, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlip{}
	response := transactionsparepartpayloads.SupplySlipResponse{}

	rows, err := tx.Model(&entities).
		Where("supply_system_number = ?", Id).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *SupplySlipRepositoryImpl) GetSupplySlipDetailById(tx *gorm.DB, Id int) (transactionsparepartpayloads.SupplySlipDetailResponse, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlipDetail{}
	response := transactionsparepartpayloads.SupplySlipDetailResponse{}

	rows, err := tx.Model(&entities).
		Where("supply_slip_detail_system_number = ?", Id).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *SupplySlipRepositoryImpl) SaveSupplySlip(tx *gorm.DB, request transactionsparepartentities.SupplySlip) (transactionsparepartentities.SupplySlip, *exceptions.BaseErrorResponse) {
	entities := transactionsparepartentities.SupplySlip{
		SupplyStatusId:        request.SupplyStatusId,
		SupplyDate:            request.SupplyDate,
		SupplyDocumentNumber:  " ",
		SupplyTypeId:          request.SupplyTypeId,
		CompanyId:             request.CompanyId,
		WorkOrderSystemNumber: request.WorkOrderSystemNumber,
		TechnicianId:          request.TechnicianId,
		CampaignId:            request.CampaignId,
		Remark:                request.Remark,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return transactionsparepartentities.SupplySlip{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *SupplySlipRepositoryImpl) SaveSupplySlipDetail(tx *gorm.DB, request transactionsparepartentities.SupplySlipDetail) (transactionsparepartentities.SupplySlipDetail, *exceptions.BaseErrorResponse) {
	total := request.QuantitySupply - 0
	entities := transactionsparepartentities.SupplySlipDetail{
		SupplySystemNumber:                request.SupplySystemNumber,
		WorkOrderOperationId:              request.WorkOrderOperationId,
		WorkOrderItemId:                   request.WorkOrderItemId,
		LocationId:                        request.LocationId,
		UnitOfMeasurementId:               request.UnitOfMeasurementId,
		QuantitySupply:                    request.QuantitySupply,
		QuantityReturn:                    0,
		QuantityDemand:                    request.QuantityDemand,
		CostOfGoodsSold:                   0,
		PurchaseRequestDetailSystemNumber: request.PurchaseRequestDetailSystemNumber,
		WorkOrderDetailId:                 request.WorkOrderDetailId,
		WarehouseGroupId:                  request.WarehouseGroupId,
		WarehouseId:                       request.WarehouseId,
		QuantityTotal:                     total,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return transactionsparepartentities.SupplySlipDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities, nil
}

func (r *SupplySlipRepositoryImpl) GetAllSupplySlip(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []transactionsparepartpayloads.SupplySlipSearchResponse
	var getSupplyTypeResponse transactionsparepartpayloads.SupplyTypeResponse

	var internalServiceFilter, externalServiceFilter []utils.FilterCondition
	var supplyTypeId string
	responseStruct := reflect.TypeOf(transactionsparepartpayloads.SupplySlipSearchResponse{})

	for i := 0; i < len(filterCondition); i++ {
		flag := false
		for j := 0; j < responseStruct.NumField(); j++ {
			if filterCondition[i].ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, filterCondition[i])
				flag = true
				break
			}
		}
		if !flag {
			externalServiceFilter = append(externalServiceFilter, filterCondition[i])
		}
	}

	//apply external services filter
	for i := 0; i < len(externalServiceFilter); i++ {
		supplyTypeId = externalServiceFilter[i].ColumnValue
	}

	// define table struct
	tableStruct := transactionsparepartpayloads.SupplySlipSearchResponse{}
	//define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	//apply filter
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)

	// Execute the query
	rows, err := whereQuery.Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}
	defer rows.Close()

	// Scan the results into the responses slice
	for rows.Next() {
		var response transactionsparepartpayloads.SupplySlipSearchResponse
		if err := rows.Scan(&response.SupplySystemNumber, &response.SupplyDocumentNumber, &response.SupplyDate, &response.SupplyTypeId, &response.WorkOrderSystemNumber, &response.WorkOrderDocumentNumber, &response.CustomerId, &response.SupplyStatusId); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		responses = append(responses, response)
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	// Fetch order type data
	supplyTypeUrl := config.EnvConfigs.GeneralServiceUrl + "supply-type/" + supplyTypeId
	fmt.Println(supplyTypeUrl)
	errUrlSupplyType := utils.Get(supplyTypeUrl, &getSupplyTypeResponse, nil)
	if errUrlSupplyType != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errUrlSupplyType,
		}
	}

	// Perform inner join with order type data
	joinedData := utils.DataFrameInnerJoin(responses, getSupplyTypeResponse, "SupplyTypeId")

	// Paginate the joined data
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}
