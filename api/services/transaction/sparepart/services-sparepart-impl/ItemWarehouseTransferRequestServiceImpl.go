package transactionsparepartserviceimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	masteritemrepository "after-sales/api/repositories/master/item"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	excelize "github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func NewWhTransferRequestImpl(transferRequestRepo transactionsparepartrepository.ItemWarehouseTransferRequestRepository, db *gorm.DB, redis *redis.Client, itemRepository masteritemrepository.ItemRepository, unitOfMeasurement masteritemrepository.UnitOfMeasurementRepository, transferReceiptRepo transactionsparepartrepository.ItemWarehouseTransferReceiptRepository) transactionsparepartservice.ItemWarehouseTransferRequestService {
	return &WhTransferRequestServiceImpl{
		TransferRequestRepo: transferRequestRepo,
		TransferReceiptRepo: transferReceiptRepo,
		DB:                  db,
		RedisClient:         redis,
		ItemRepository:      itemRepository,
		UnitOfMeasurement:   unitOfMeasurement,
	}
}

type WhTransferRequestServiceImpl struct {
	TransferRequestRepo transactionsparepartrepository.ItemWarehouseTransferRequestRepository
	TransferReceiptRepo transactionsparepartrepository.ItemWarehouseTransferReceiptRepository
	ItemRepository      masteritemrepository.ItemRepository
	UnitOfMeasurement   masteritemrepository.UnitOfMeasurementRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client
}

// AcceptTransferReceipt implements transactionsparepartservice.ItemWarehouseTransferRequestService.
func (s *WhTransferRequestServiceImpl) AcceptTransferReceipt(number int, request transactionsparepartpayloads.AcceptWarehouseTransferRequestRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferReceiptRepo.Accept(tx, number, request)
	if err != nil {
		return result, err
	}
	return result, nil
}

// RejectTransferReceipt implements transactionsparepartservice.ItemWarehouseTransferRequestService.
func (s *WhTransferRequestServiceImpl) RejectTransferReceipt(number int, request transactionsparepartpayloads.RejectWarehouseTransferRequestRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferReceiptRepo.Reject(tx, number, request)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GenerateTemplateFile implements transactionsparepartservice.ItemWarehouseTransferRequestService.
func (s *WhTransferRequestServiceImpl) GenerateTemplateFile() (*excelize.File, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	f := excelize.NewFile()
	sheetName := "Sheet1"
	defer func() {
		if err := f.Close(); err != nil {
			log.Error(err)
		}
	}()

	index, errA := f.NewSheet(sheetName)
	if errA != nil {
		return nil, &exceptions.BaseErrorResponse{Err: errA, StatusCode: http.StatusInternalServerError}
	}

	f.SetCellValue(sheetName, "A1", "ITEM_CODE") // bom_item_id
	f.SetCellValue(sheetName, "B1", "REQ_QTY")

	style, errA := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if errA != nil {
		return nil, &exceptions.BaseErrorResponse{Err: errA, StatusCode: http.StatusInternalServerError}
	}

	for col := 'A'; col <= 'B'; col++ {
		cell := string(col) + "1"
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	f.SetActiveSheet(index)
	return f, nil
}

// ProcessUploadData implements transactionsparepartservice.ItemWarehouseTransferRequestService.
func (s *WhTransferRequestServiceImpl) ProcessUploadData(request transactionsparepartpayloads.UploadProcessItemWarehouseTransferRequestPayloads) ([]transactionsparepartentities.ItemWarehouseTransferRequestDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse
	var results []transactionsparepartentities.ItemWarehouseTransferRequestDetail

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	details := request.TransferRequestDetails

	var validation_text string

	for _, detail := range details {
		item, errItem := s.ItemRepository.GetItemCode(tx, detail.ItemCode)
		if errItem != nil {
			if err.StatusCode == http.StatusNotFound {
				validation_text += "item code does not exist in Item Master. "
			} else {
				return []transactionsparepartentities.ItemWarehouseTransferRequestDetail{}, err
			}
		}
		var get transactionsparepartentities.ItemWarehouseTransferRequestDetail
		get, err = s.TransferRequestRepo.InsertWhTransferRequestDetail(tx, transactionsparepartpayloads.InsertItemWarehouseTransferDetailRequest{
			TransferRequestSystemNumberId: request.TransferRequestSystemNumber,
			ModifiedById:                  request.ModifiedById,
			ItemId:                        &item.ItemId,
			RequestQuantity:               detail.RequestQuantity,
		})

		results = append(results, get)
	}

	if err != nil {
		return []transactionsparepartentities.ItemWarehouseTransferRequestDetail{}, err
	}

	return results, nil
}

// PreviewUploadData implements transactionsparepartservice.ItemWarehouseTransferRequestService.
func (s *WhTransferRequestServiceImpl) PreviewUploadData(rows [][]string) ([]transactionsparepartpayloads.UploadPreviewItemWarehouseTransferRequestPayloads, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse
	var results []transactionsparepartpayloads.UploadPreviewItemWarehouseTransferRequestPayloads

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) != 2 {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid row length",
			}
		}

		var validation_text string

		item, err := s.ItemRepository.GetItemCode(tx, row[0])
		if err != nil {
			if err.StatusCode == http.StatusNotFound {
				validation_text += "item code does not exist in Item Master. "
			} else {
				return []transactionsparepartpayloads.UploadPreviewItemWarehouseTransferRequestPayloads{}, err
			}
		}

		uom, errUom := s.UnitOfMeasurement.GetUnitOfMeasurementById(tx, *item.UnitOfMeasurementStockId)
		if errUom != nil {
			if err.StatusCode == http.StatusNotFound {
				validation_text += "unit of measurement does not exist in Item Master. "
			} else {
				return []transactionsparepartpayloads.UploadPreviewItemWarehouseTransferRequestPayloads{}, err
			}
		}

		row[1] = strings.ReplaceAll(row[1], ",", ".") // Replace comma with dot
		reqQuantity, errA := strconv.ParseFloat(row[1], 64)
		if errA != nil {
			return nil, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid request quantity format",
			}
		}

		results = append(results, transactionsparepartpayloads.UploadPreviewItemWarehouseTransferRequestPayloads{
			ItemCode:          row[0],
			ItemName:          item.ItemName,
			RequestQuantity:   reqQuantity,
			UnitOfMeasurement: uom.UomCode,
		})
	}

	return results, nil
}

// GetByIdTransferRequestDetail implements transactionsparepartservice.ItemWarehouseTransferRequestService.
func (s *WhTransferRequestServiceImpl) GetByIdTransferRequestDetail(number int) (transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestDetailResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferRequestRepo.GetByIdTransferRequestDetail(tx, number)
	if err != nil {
		return result, err
	}
	return result, nil
}

// UpdateWhTransferRequestDetail implements transactionsparepartservice.ItemWarehouseTransferRequestService.
func (s *WhTransferRequestServiceImpl) UpdateWhTransferRequestDetail(request transactionsparepartpayloads.UpdateItemWarehouseTransferRequestDetailRequest, number int) (transactionsparepartentities.ItemWarehouseTransferRequestDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferRequestRepo.UpdateWhTransferRequestDetail(tx, request, number)
	if err != nil {
		return result, err
	}
	return result, nil
}

// DeleteDetail implements transactionsparepartservice.WhTransferRequestService.
func (s *WhTransferRequestServiceImpl) DeleteDetail(number []int, request transactionsparepartpayloads.DeleteDetailItemWarehouseTransferRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse
	var result bool

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	for _, num := range number {
		result, err = s.TransferRequestRepo.DeleteDetail(tx, num, request)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

// DeleteHeaderTransferRequest implements transactionsparepartservice.WhTransferRequestService.
func (s *WhTransferRequestServiceImpl) DeleteHeaderTransferRequest(number int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferRequestRepo.DeleteHeaderTransferRequest(tx, number)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetAllWhTransferRequest implements transactionsparepartservice.WhTransferRequestService.
func (s *WhTransferRequestServiceImpl) GetAllWhTransferRequest(pages pagination.Pagination, filter []utils.FilterCondition, dateParams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferRequestRepo.GetAllWhTransferRequest(tx, pages, filter, dateParams)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetAllDetailTransferRequest implements transactionsparepartservice.WhTransferRequestService.
func (s *WhTransferRequestServiceImpl) GetAllDetailTransferRequest(number int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferRequestRepo.GetAllDetailTransferRequest(tx, number, pages)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetByIdTransferRequest implements transactionsparepartservice.WhTransferRequestService.
func (s *WhTransferRequestServiceImpl) GetByIdTransferRequest(number int) (transactionsparepartpayloads.GetByIdItemWarehouseTransferRequestResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferRequestRepo.GetByIdTransferRequest(tx, number)
	if err != nil {
		return result, err
	}
	return result, nil
}

// InsertWhTransferRequestDetail implements transactionsparepartservice.WhTransferRequestService.
func (s *WhTransferRequestServiceImpl) InsertWhTransferRequestDetail(request transactionsparepartpayloads.InsertItemWarehouseTransferDetailRequest) (transactionsparepartentities.ItemWarehouseTransferRequestDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferRequestRepo.InsertWhTransferRequestDetail(tx, request)
	if err != nil {
		return result, err
	}
	return result, nil
}

// InsertWhTransferRequestHeader implements transactionsparepartservice.WhTransferRequestService.
func (s *WhTransferRequestServiceImpl) InsertWhTransferRequestHeader(request transactionsparepartpayloads.InsertItemWarehouseTransferRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferRequestRepo.InsertWhTransferRequestHeader(tx, request)
	if err != nil {
		return result, err
	}
	return result, nil
}

// SubmitWhTransferRequest implements transactionsparepartservice.WhTransferRequestService.
func (s *WhTransferRequestServiceImpl) SubmitWhTransferRequest(number int, request transactionsparepartpayloads.SubmitItemWarehouseTransferRequest) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferRequestRepo.SubmitWhTransferRequest(tx, number, request)
	if err != nil {
		return result, err
	}
	return result, nil
}

// UpdateWhTransferRequest implements transactionsparepartservice.WhTransferRequestService.
func (s *WhTransferRequestServiceImpl) UpdateWhTransferRequest(request transactionsparepartpayloads.UpdateItemWarehouseTransferRequest, number int) (transactionsparepartentities.ItemWarehouseTransferRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()
	result, err := s.TransferRequestRepo.UpdateWhTransferRequest(tx, request, number)
	if err != nil {
		return result, err
	}
	return result, nil
}
