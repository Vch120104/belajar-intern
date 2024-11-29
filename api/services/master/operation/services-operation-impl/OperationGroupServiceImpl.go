package masteroperationserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OperationGroupServiceImpl struct {
	operationGroupRepo masteroperationrepository.OperationGroupRepository
	DB                 *gorm.DB
	RedisClient        *redis.Client // Redis client
}

func StartOperationGroupService(operationGroupRepo masteroperationrepository.OperationGroupRepository, db *gorm.DB, redisClient *redis.Client) masteroperationservice.OperationGroupService {
	return &OperationGroupServiceImpl{
		operationGroupRepo: operationGroupRepo,
		DB:                 db,
		RedisClient:        redisClient,
	}
}

func (s *OperationGroupServiceImpl) GetOperationGroupDropDown() ([]masteroperationpayloads.OperationGroupDropDownResponse, *exceptions.BaseErrorResponse) {
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
	get, err := s.operationGroupRepo.GetOperationGroupDropDown(tx)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *OperationGroupServiceImpl) GetOperationGroupById(id int) (masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.operationGroupRepo.GetOperationGroupById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationGroupServiceImpl) GetOperationGroupByCode(Code string) (masteroperationpayloads.OperationGroupResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.operationGroupRepo.GetOperationGroupByCode(tx, Code)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OperationGroupServiceImpl) GetAllOperationGroup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	// tx := s.DB.Begin()
	// defer helper.CommitOrRollback(tx)
	// results, err := s.operationGroupRepo.GetAllOperationGroup(tx, filterCondition, pages)
	// if err != nil {
	// 	panic(exceptions.NewNotFoundError(err.Error()))
	// }
	// return results
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
	get, err := s.operationGroupRepo.GetAllOperationGroup(tx, filterCondition, pages)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *OperationGroupServiceImpl) ChangeStatusOperationGroup(oprId int) (bool, *exceptions.BaseErrorResponse) {
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

	_, err = s.operationGroupRepo.GetOperationGroupById(tx, oprId)

	if err != nil {
		return false, err
	}

	results, err := s.operationGroupRepo.ChangeStatusOperationGroup(tx, oprId)

	if err != nil {
		return results, err
	}
	return true, nil
}

func (s *OperationGroupServiceImpl) SaveOperationGroup(req masteroperationpayloads.OperationGroupResponse) (bool, *exceptions.BaseErrorResponse) {
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

	if req.OperationGroupId != 0 {
		_, err := s.operationGroupRepo.GetOperationGroupById(tx, req.OperationGroupId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.operationGroupRepo.SaveOperationGroup(tx, req)

	if err != nil {
		return false, err
	}
	return results, nil
}
