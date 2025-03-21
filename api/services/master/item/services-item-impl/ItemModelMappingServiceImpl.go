package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ItemModelMappingServiceImpl struct {
	ItemModelMappingRepo masteritemrepository.ItemModelMappingRepository
	DB                   *gorm.DB
}

// GetItemModelMappingByItemId implements masteritemservice.ItemModelMappingService.
func (s *ItemModelMappingServiceImpl) GetItemModelMappingByItemId(itemId int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
	results, err := s.ItemModelMappingRepo.GetItemModelMappingByItemId(tx, itemId, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

// UpdateItemModelMapping implements masteritemservice.ItemModelMappingService.
func (s *ItemModelMappingServiceImpl) UpdateItemModelMapping(req masteritempayloads.CreateItemModelMapping) (bool, *exceptions.BaseErrorResponse) {
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
	results, err := s.ItemModelMappingRepo.UpdateItemModelMapping(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

// CreateItemModelMapping implements masteritemservice.ItemModelMappingService.
func (s *ItemModelMappingServiceImpl) CreateItemModelMapping(req masteritempayloads.CreateItemModelMapping) (bool, *exceptions.BaseErrorResponse) {
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
	results, err := s.ItemModelMappingRepo.CreateItemModelMapping(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func StartItemModelMappingService(ItemModelMappingRepo masteritemrepository.ItemModelMappingRepository, db *gorm.DB) masteritemservice.ItemModelMappingService {
	return &ItemModelMappingServiceImpl{
		ItemModelMappingRepo: ItemModelMappingRepo,
		DB:                   db,
	}
}
