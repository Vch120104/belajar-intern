package transactionsparepartrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
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
