package transactionworkshopserviceimpl

import (
	transactionworkshopentities "after-sales/api/entities/transaction/workshop"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionworkshoppayloads "after-sales/api/payloads/transaction/workshop"
	transactionworkshoprepository "after-sales/api/repositories/transaction/workshop"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookingEstimationServiceImpl struct {
	structBookingEstimationRepo transactionworkshoprepository.BookingEstimationRepository
	DB                          *gorm.DB
	RedisClient                 *redis.Client // Redis client
}

func OpenBookingEstimationServiceImpl(bookingEstimationRepo transactionworkshoprepository.BookingEstimationRepository, db *gorm.DB, redisClient *redis.Client) transactionworkshopservice.BookingEstimationService {
	return &BookingEstimationServiceImpl{
		structBookingEstimationRepo: bookingEstimationRepo,
		DB:                          db,
		RedisClient:                 redisClient,
	}
}

func (s *BookingEstimationServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	results, err := s.structBookingEstimationRepo.GetAll(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *BookingEstimationServiceImpl) New(request transactionworkshoppayloads.BookingEstimationNormalRequest) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse) {
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

	results, err := s.structBookingEstimationRepo.New(tx, request)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *BookingEstimationServiceImpl) GetById(id int) (transactionworkshoppayloads.BookingEstimationResponseDetail, *exceptions.BaseErrorResponse) {
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
	results, err := s.structBookingEstimationRepo.GetById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil

}

func (s *BookingEstimationServiceImpl) Save(request transactionworkshoppayloads.BookingEstimationSaveRequest, id int) (transactionworkshopentities.BookingEstimation, *exceptions.BaseErrorResponse) {
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
	save, err := s.structBookingEstimationRepo.Save(tx, request, id)
	if err != nil {
		return transactionworkshopentities.BookingEstimation{}, err
	}
	return save, nil
}

func (s *BookingEstimationServiceImpl) Submit(tx *gorm.DB, id int) (bool, *exceptions.BaseErrorResponse) {
	tx = s.DB.Begin()
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
	result, err := s.structBookingEstimationRepo.Submit(tx, id)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) Void(id int) (bool, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.Void(tx, id)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) CloseOrder(tx *gorm.DB, id int) *exceptions.BaseErrorResponse {
	tx = s.DB.Begin()
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
	err = s.structBookingEstimationRepo.CloseOrder(tx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *BookingEstimationServiceImpl) SaveBookEstimReq(req transactionworkshoppayloads.BookEstimRemarkRequest, id int) (transactionworkshopentities.BookingEstimationRequest, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.SaveBookEstimReq(tx, req, id)
	if err != nil {
		return transactionworkshopentities.BookingEstimationRequest{}, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) UpdateBookEstimReq(req transactionworkshoppayloads.BookEstimRemarkRequest, booksysno int, id int) (transactionworkshopentities.BookingEstimationRequest, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.UpdateBookEstimReq(tx, req, booksysno, id)
	if err != nil {
		return transactionworkshopentities.BookingEstimationRequest{}, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) GetByIdBookEstimReq(booksysno int, id int) (transactionworkshoppayloads.BookEstimRemarkResponse, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.GetByIdBookEstimReq(tx, booksysno, id)
	if err != nil {
		return transactionworkshoppayloads.BookEstimRemarkResponse{}, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) GetAllBookEstimReq(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.GetAllBookEstimReq(tx, filterCondition, pages)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) SaveBookEstimReminderServ(req transactionworkshoppayloads.ReminderServicePost, id int) (transactionworkshopentities.BookingEstimationServiceReminder, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.SaveBookEstimReminderServ(tx, req, id)
	if err != nil {
		return transactionworkshopentities.BookingEstimationServiceReminder{}, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) SaveDetailBookEstim(id int, req transactionworkshoppayloads.BookingEstimationDetailRequest) (transactionworkshopentities.BookingEstimationDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var errResponse *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			errResponse = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if errResponse != nil {
			tx.Rollback()
			logrus.WithError(errResponse.Err).Info("Transaction rollback due to error")
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				errResponse = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	// Simpan detail booking estimation
	result, err := s.structBookingEstimationRepo.SaveDetailBookEstim(tx, id, req)
	if err != nil {
		errResponse = &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to save booking estimation detail: %w", err),
		}
		return transactionworkshopentities.BookingEstimationDetail{}, errResponse
	}

	// Update perhitungan estimasi
	if _, err = s.structBookingEstimationRepo.PutBookingEstimationCalculation(tx, result.EstimationSystemNumber); err != nil {
		errResponse = &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to update booking estimation calculation: %w", err),
		}
		return transactionworkshopentities.BookingEstimationDetail{}, errResponse
	}

	return result, nil
}

func (s *BookingEstimationServiceImpl) AddPackage(id int, packId int) (bool, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.AddPackage(tx, id, packId)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) AddContractService(id int, contractserviceid int) (bool, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.AddContractService(tx, id, contractserviceid)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) CopyFromHistory(batchid int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.CopyFromHistory(tx, batchid)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) InputDiscount(id int, req transactionworkshoppayloads.BookEstimationPayloadsDiscount) (int, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.InputDiscount(tx, id, req)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) AddFieldAction(id int, idrecall int) (int, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.AddFieldAction(tx, id, idrecall)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) GetByIdBookEstimDetail(estimsysno int, id int) (transactionworkshoppayloads.BookingEstimationDetailResponse, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.GetByIdBookEstimDetail(tx, estimsysno, id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) PostBookingEstimationCalculation(id int) (int, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.PostBookingEstimationCalculation(tx, id)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) SaveBookingEstimationFromPDI(id int, req transactionworkshoppayloads.PdiServiceRequest) (bool, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.SaveBookingEstimationFromPDI(tx, id, req)
	if err != nil || !result {
		return false, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) SaveBookingEstimationFromServiceRequest(id int, req transactionworkshoppayloads.PdiServiceRequest) (bool, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.SaveBookingEstimationFromServiceRequest(tx, id, req)
	if err != nil || !result {
		return false, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) SaveBookingEstimationAllocation(id int, req transactionworkshoppayloads.BookEstimationAllocation) (transactionworkshopentities.BookingEstimationAllocation, *exceptions.BaseErrorResponse) {
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
	result, err := s.structBookingEstimationRepo.SaveBookingEstimationAllocation(tx, id, req)
	if err != nil {
		return transactionworkshopentities.BookingEstimationAllocation{}, err
	}
	return result, nil
}

func (s *BookingEstimationServiceImpl) GetAllDetailBookingEstimation(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	results, err := s.structBookingEstimationRepo.GetAllDetailBookingEstimation(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *BookingEstimationServiceImpl) DeleteBookEstimReq(booksysno int, id []int) (bool, *exceptions.BaseErrorResponse) {
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

	result, err := s.structBookingEstimationRepo.DeleteBookEstimReq(tx, booksysno, id)
	if err != nil {
		return false, err
	}
	return result, nil
}
