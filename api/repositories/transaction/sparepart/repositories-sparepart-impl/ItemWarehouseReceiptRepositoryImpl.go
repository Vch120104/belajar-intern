package transactionsparepartrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	"after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func NewItemWarehouseTransferReceiptRepositoryImpl() transactionsparepartrepository.ItemWarehouseTransferReceiptRepository {
	return &ItemWarehouseReceiptRepositoryImpl{}
}

type ItemWarehouseReceiptRepositoryImpl struct {
}

// GetAll implements transactionsparepartrepository.ItemWarehouseTransferReceiptRepository.
func (*ItemWarehouseReceiptRepositoryImpl) GetAll(tx *gorm.DB, pages pagination.Pagination, filter []utils.FilterCondition, dateParams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest
	var responses []transactionsparepartpayloads.GetAllItemWarehouseTransferRequestResponse
	var statusA masteritementities.ItemTransferStatus
	var statusReady masteritementities.ItemTransferStatus
	var statusReject masteritementities.ItemTransferStatus

	errGetStatus := tx.Model(&statusA).Where("item_transfer_status_code = ?", 20).Find(&statusA)
	if errGetStatus.Error != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatus.Error,
		}
	}

	errGetStatus2 := tx.Model(&statusReady).Where("item_transfer_status_code = ?", 15).Find(&statusReady)
	if errGetStatus2.Error != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatus2.Error,
		}
	}

	errGetStatusReject := tx.Model(&statusReject).Where("item_transfer_status_code = ?", 30).Find(&statusReject)
	if errGetStatusReject.Error != nil {
		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatusReject.Error,
		}
	}

	joinTable := tx.Model(&entities).
		Select(
			"transfer_request_system_number",
			"transfer_request_document_number",
			"transfer_request_status_id",
			"item_transfer_status_code transfer_request_status_code",
			"item_transfer_status_description transfer_request_status_description",
			"transfer_request_date",
			"request_from_warehouse_id",
			"wmf.warehouse_name request_from_warehouse_name",
			"wmf.warehouse_group_id request_from_warehouse_group_id",
			"wgf.warehouse_group_name request_from_warehouse_group_name",
			"transfer_request_by_id",
			"request_to_warehouse_id",
			"wmt.warehouse_name request_to_warehouse_name",
			"wmt.warehouse_group_id request_to_warehouse_group_id",
			"wgt.warehouse_group_name request_to_warehouse_group_name",
		).
		Where("transfer_request_status_id in (?, ?, ?)", statusA.ItemTransferStatusId, statusReady.ItemTransferStatusId, statusReject.ItemTransferStatusId).
		Joins("LEFT JOIN mtr_warehouse_master wmf on wmf.warehouse_id = request_from_warehouse_id").
		Joins("LEFT JOIN mtr_item_transfer_status stat on stat.item_transfer_status_id = transfer_request_status_id").
		Joins("LEFT JOIN mtr_warehouse_master wmt on wmt.warehouse_id = request_to_warehouse_id").
		Joins("LEFT JOIN mtr_warehouse_group wgf on wgf.warehouse_group_id = wmf.warehouse_group_id").
		Joins("LEFT JOIN mtr_warehouse_group wgt on wgt.warehouse_group_id = wmt.warehouse_group_id")

	whereQuery := utils.ApplyFilter(joinTable, filter)

	var strDateFilter string
	if dateParams["transfer_request_date_from"] == "" {
		dateParams["transfer_request_date_from"] = "19000101"
	}
	if dateParams["transfer_request_date_to"] == "" {
		dateParams["transfer_request_date_to"] = "99991212"
	}
	strDateFilter = "transfer_request_date >='" + dateParams["transfer_request_date_from"] + "' AND transfer_request_date <= '" + dateParams["transfer_request_date_to"] + "'"

	whereQuery = whereQuery.Where(strDateFilter)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
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

	for i, respon := range responses {
		get, errUser := generalserviceapiutils.GetUserDetailsByID(respon.TransferRequestById)
		if errUser != nil {
			return pages, nil
		}
		responses[i].TransferRequestByName = get.EmployeeName
	}

	pages.Rows = responses
	return pages, nil
}

// Accept implements transactionsparepartrepository.ItemWarehouseTransferReceiptRepository.
func (*ItemWarehouseReceiptRepositoryImpl) Accept(tx *gorm.DB, number int, request transactionsparepartpayloads.AcceptWarehouseTransferRequestRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest

	errGetEntities := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestSystemNumber: number}).First(&entities).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item claim with that id is not found please check input",
				Err:        errGetEntities,
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get item claim please check input",
		}
	}

	var statusReady masteritementities.ItemTransferStatus
	var statusAcc masteritementities.ItemTransferStatus

	errGetStatusAcc := tx.Model(&statusAcc).Where("item_transfer_status_code = ?", 20).Find(&statusAcc)
	if errGetStatusAcc.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatusAcc.Error,
		}
	}

	errGetStatusReady := tx.Model(&statusReady).Where("item_transfer_status_code = ?", 15).Find(&statusReady)
	if errGetStatusReady.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatusReady.Error,
		}
	}

	if entities.TransferRequestStatusId != statusReady.ItemTransferStatusId {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("transfer request status is not ready"),
			Message:    "transfer request status is not ready",
		}
	}

	now := time.Now()

	entities.TransferRequestStatusId = statusAcc.ItemTransferStatusId
	entities.ApprovalById = &request.ApprovalById
	entities.ApprovalDate = &now
	entities.ApprovalRemark = request.ApprovalRemark
	entities.ModifiedById = request.ApprovalById

	err := tx.Save(&entities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to update transfer request",
		}
	}

	return entities, nil
}

// Reject implements transactionsparepartrepository.ItemWarehouseTransferReceiptRepository.
func (*ItemWarehouseReceiptRepositoryImpl) Reject(tx *gorm.DB, number int, request transactionsparepartpayloads.RejectWarehouseTransferRequestRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse) {
	var entities transactionsparepartentities.ItemWarehouseTransferRequest

	errGetEntities := tx.Model(&entities).Where(transactionsparepartentities.ItemWarehouseTransferRequest{TransferRequestSystemNumber: number}).First(&entities).Error
	if errGetEntities != nil {
		if errors.Is(errGetEntities, gorm.ErrRecordNotFound) {
			return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "item claim with that id is not found please check input",
				Err:        errGetEntities,
			}
		}
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetEntities,
			Message:    "failed to get item claim please check input",
		}
	}

	var statusReady masteritementities.ItemTransferStatus
	var statusRej masteritementities.ItemTransferStatus

	errGetStatusRej := tx.Model(&statusRej).Where("item_transfer_status_code = ?", 30).Find(&statusRej)
	if errGetStatusRej.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatusRej.Error,
		}
	}

	errGetStatusReady := tx.Model(&statusReady).Where("item_transfer_status_code = ?", 15).Find(&statusReady)
	if errGetStatusReady.Error != nil {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errGetStatusReady.Error,
		}
	}

	if entities.TransferRequestStatusId != statusReady.ItemTransferStatusId {
		return transactionsparepartentities.ItemWarehouseTransferRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("transfer request status is not ready"),
			Message:    "transfer request status is not ready",
		}
	}

	now := time.Now()

	entities.TransferRequestStatusId = statusRej.ItemTransferStatusId
	entities.ApprovalById = &request.ApprovalById
	entities.ApprovalDate = &now
	entities.ApprovalRemark = request.ApprovalRemark
	entities.ModifiedById = request.ApprovalById

	err := tx.Save(&entities).Error
	if err != nil {
		return entities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			Message:    "failed to update transfer request",
		}
	}

	return entities, nil
}
