package transactionsparepartserviceimpl

import (
	transactionsparepartentities "after-sales/api/entities/transaction/sparepart"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PurchaseOrderServiceImpl struct {
	PurchaseOrderRepo transactionsparepartrepository.PurchaseOrderRepository
	DB                *gorm.DB
	redis             *redis.Client
}

func NewPurchaseOrderService(PuchaseOrderRepo transactionsparepartrepository.PurchaseOrderRepository, db *gorm.DB, redisclient *redis.Client) transactionsparepartservice.PurchaseOrderService {
	return &PurchaseOrderServiceImpl{
		PurchaseOrderRepo: PuchaseOrderRepo,
		DB:                db,
		redis:             redisclient,
	}
}

func (s *PurchaseOrderServiceImpl) GetAllPurchaseOrder(filter []utils.FilterCondition, page pagination.Pagination, DateParams map[string]string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	result, err := s.PurchaseOrderRepo.GetAllPurchaseOrder(tx, filter, page, DateParams)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *PurchaseOrderServiceImpl) GetByIdPurchaseOrder(i int) (transactionsparepartpayloads.PurchaseOrderGetByIdResponses, *exceptions.BaseErrorResponse) {
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
	result, err := s.PurchaseOrderRepo.GetByIdPurchaseOrder(tx, i)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *PurchaseOrderServiceImpl) GetByIdPurchaseOrderDetail(id int, page pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	result, err := s.PurchaseOrderRepo.GetAllDetailByHeaderId(tx, id, page)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *PurchaseOrderServiceImpl) NewPurchaseOrderHeader(responses transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderResponses) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse) {
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
	result, err := s.PurchaseOrderRepo.NewPurchaseOrderHeader(tx, responses)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *PurchaseOrderServiceImpl) UpdatePurchaseOrderHeader(id int, responses transactionsparepartpayloads.PurchaseOrderNewPurchaseOrderPayloads) (transactionsparepartentities.PurchaseOrderEntities, *exceptions.BaseErrorResponse) {
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
	result, err := s.PurchaseOrderRepo.UpdatePurchaseOrderHeader(tx, id, responses)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *PurchaseOrderServiceImpl) GetPurchaseOrderDetailById(id int) (transactionsparepartpayloads.PurchaseOrderGetDetail, *exceptions.BaseErrorResponse) {
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
	result, err := s.PurchaseOrderRepo.GetPurchaseOrderDetailById(tx, id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *PurchaseOrderServiceImpl) NewPurchaseOrderDetail(payloads transactionsparepartpayloads.PurchaseOrderDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse) {
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
	result, err := s.PurchaseOrderRepo.NewPurchaseOrderDetail(tx, payloads)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *PurchaseOrderServiceImpl) DeletePurchaseOrderDetailMultiId(multiid string) (bool, *exceptions.BaseErrorResponse) {
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
	result, err := s.PurchaseOrderRepo.DeletePurchaseOrderDetailMultiId(tx, multiid)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *PurchaseOrderServiceImpl) SavePurchaseOrderDetail(payloads transactionsparepartpayloads.PurchaseOrderSaveDetailPayloads) (transactionsparepartentities.PurchaseOrderDetailEntities, *exceptions.BaseErrorResponse) {
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
	result, err := s.PurchaseOrderRepo.SavePurchaseOrderDetail(tx, payloads)
	if err != nil {
		return result, err
	}
	return result, nil
}
