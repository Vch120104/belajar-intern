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
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

type SalesOrderServiceImpl struct {
	salesOrderRepo transactionsparepartrepository.SalesOrderRepository
	DB             *gorm.DB
	RedisClient    *redis.Client // Redis client
}

func StartSalesOrderService(salesOrderRepo transactionsparepartrepository.SalesOrderRepository, db *gorm.DB, redisClient *redis.Client) transactionsparepartservice.SalesOrderServiceInterface {
	return &SalesOrderServiceImpl{
		salesOrderRepo: salesOrderRepo,
		DB:             db,
		RedisClient:    redisClient,
	}
}

func (s *SalesOrderServiceImpl) InsertSalesOrderHeader(payload transactionsparepartpayloads.SalesOrderInsertHeaderPayload) (transactionsparepartentities.SalesOrder, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	//ini repo kedua
	result, err := s.salesOrderRepo.InsertSalesOrderHeader(tx, payload)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *SalesOrderServiceImpl) GetSalesOrderByID(Id int) (transactionsparepartpayloads.SalesOrderEstimationGetByIdResponse, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	result, err := s.salesOrderRepo.GetSalesOrderByID(tx, Id)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (s *SalesOrderServiceImpl) GetAllSalesOrder(pages pagination.Pagination, condition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	result, err := s.salesOrderRepo.GetAllSalesOrder(tx, pages, condition)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *SalesOrderServiceImpl) VoidSalesOrder(salesOrderId int) (bool, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	result, err := s.salesOrderRepo.VoidSalesOrder(tx, salesOrderId)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (s *SalesOrderServiceImpl) InsertSalesOrderDetail(payload transactionsparepartpayloads.SalesOrderDetailInsertPayload) (transactionsparepartentities.SalesOrderDetail, *exceptions.BaseErrorResponse) {
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
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	result, err := s.salesOrderRepo.InsertSalesOrderDetail(tx, payload)
	if err != nil {
		return result, err
	}
	return result, nil
}
