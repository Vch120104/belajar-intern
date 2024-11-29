package masteritemserviceimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	"after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ItemGroupServiceImpl struct {
	repository masteritemrepository.ItemGroupRepository
	DB         *gorm.DB
	rdb        *redis.Client
}

func (i *ItemGroupServiceImpl) GetAllItemGroupWithPagination(internalFilter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := i.repository.GetAllItemGroupWithPagination(tx, internalFilter, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) GetAllItemGroup(code string) ([]masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()

	results, err := i.repository.GetAllItemGroup(tx, code)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) GetItemGroupById(id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := i.repository.GetItemGroupById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) DeleteItemGroupById(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := i.repository.DeleteItemGroupById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) UpdateItemGroupById(payload masteritempayloads.ItemGroupUpdatePayload, id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := i.repository.UpdateItemGroupById(tx, payload, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) UpdateStatusItemGroupById(id int) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := i.repository.UpdateStatusItemGroupById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) GetItemGroupByMultiId(multiId string) ([]masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := i.repository.GetItemGroupByMultiId(tx, multiId)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (i *ItemGroupServiceImpl) NewItemGroup(payload masteritempayloads.NewItemGroupPayload) (masteritementities.ItemGroup, *exceptions.BaseErrorResponse) {
	tx := i.DB.Begin()
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := i.repository.NewItemGroup(tx, payload)
	if err != nil {
		return results, err
	}
	return results, nil
}

func NewItemGroupServiceImpl(repo masteritemrepository.ItemGroupRepository, DB *gorm.DB, rdb *redis.Client) masteritemservice.ItemGroupService {
	return &ItemGroupServiceImpl{
		repository: repo,
		DB:         DB,
		rdb:        rdb,
	}
}
