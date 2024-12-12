package transactionsparepartrepositoryimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type ItemLocationTransferRepositoryImpl struct {
}

func NewItemLocationTransferRepositoryImpl() transactionsparepartrepository.ItemLocationTransferRepository {
	return &ItemLocationTransferRepositoryImpl{}
}

// uspg_atTrfReq0_Select
// IF @Option = 2
func (r *ItemLocationTransferRepositoryImpl) GetAllItemLocationTransfer(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest
	var responses []transactionsparepartpayloads.GetAllItemLocationTransferResponse

	baseModelQuery := tx.Model(&entities).
		Select(
			"trx_item_warehouse_transfer_request.transfer_request_system_number",
			"trx_item_warehouse_transfer_request.transfer_request_document_number",
			"trx_item_warehouse_transfer_request.transfer_request_status_id",
			"TransferRequestStatus.item_transfer_status_code as transfer_request_status_code",
			"TransferRequestStatus.item_transfer_status_description as transfer_request_status_description",
			"trx_item_warehouse_transfer_request.transfer_request_date",
			"trx_item_warehouse_transfer_request.transfer_request_by_id",
			"trx_item_warehouse_transfer_request.request_from_warehouse_id",
			"RequestFromWarehouse.warehouse_code as request_from_warehouse_code",
			"RequestFromWarehouse.warehouse_name as request_from_warehouse_name",
			"RequestFromWarehouse.warehouse_group_id as request_from_warehouse_group_id",
			"RequestFromWarehouseGroup.warehouse_group_code as request_from_warehouse_group_code",
			"RequestFromWarehouseGroup.warehouse_group_name asrequest_from_warehouse_group_name",
			"trx_item_warehouse_transfer_request.request_to_warehouse_id",
			"RequestToWarehouse.warehouse_code as request_to_warehouse_code",
			"RequestToWarehouse.warehouse_name as request_to_warehouse_name",
			"RequestToWarehouse.warehouse_group_id as request_to_warehouse_group_id",
			"RequestToWarehouseGroup.warehouse_group_code as request_to_warehouse_group_code",
			"RequestToWarehouseGroup.warehouse_group_name as request_to_warehouse_group_name",
		).
		Joins("LEFT JOIN mtr_item_transfer_status TransferRequestStatus ON TransferRequestStatus.item_transfer_status_id = trx_item_warehouse_transfer_request.transfer_request_status_id").
		Joins("LEFT JOIN mtr_warehouse_master RequestFromWarehouse ON RequestFromWarehouse.warehouse_id = trx_item_warehouse_transfer_request.request_from_warehouse_id").
		Joins("LEFT JOIN mtr_warehouse_group RequestFromWarehouseGroup ON RequestFromWarehouseGroup.warehouse_group_id = RequestFromWarehouse.warehouse_group_id").
		Joins("LEFT JOIN mtr_warehouse_master RequestToWarehouse ON RequestToWarehouse.warehouse_id = trx_item_warehouse_transfer_request.request_to_warehouse_id").
		Joins("LEFT JOIN mtr_warehouse_group RequestToWarehouseGroup ON RequestToWarehouseGroup.warehouse_group_id = RequestToWarehouse.warehouse_group_id")
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	// -- ADD BY RAEYANS ON 05.11.2009
	// IF @TrfType = 'rec' AND ISNULL(@Trfreq_Status, '') = ''
	// BEGIN
	// 	SET  @strFilter = @strFilter + ' AND A.TRFREQ_STATUS IN (dbo.getVariableValue(''ITEM_TRF_STAT_ACCEPT''), dbo.getVariableValue(''ITEM_TRF_STAT_READY''), dbo.getVariableValue(''ITEM_TRF_STAT_REJECT''))'
	// END
	// -- END

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Scan(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var responsesCrossService []transactionsparepartpayloads.GetAllItemLocationTransferResponse
	for _, result := range responses {
		employeePayloads, employeeError := generalserviceapiutils.GetUserDetailsByID(result.TransferRequestById)
		if employeeError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		if employeePayloads.UserEmployeeId == 0 {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("employee data not found"),
			}
		}

		result.TransferRequestByName = employeePayloads.EmployeeName
		responsesCrossService = append(responsesCrossService, result)
	}

	pages.Rows = responsesCrossService

	return pages, nil
}
