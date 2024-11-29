package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GmmPriceCodeServiceImpl struct {
	GmmPriceCodeRepo masterrepository.GmmPriceCodeRepository
	DB               *gorm.DB
}

func StartGmmPriceCodeServiceImpl(gmmPriceCodeRepo masterrepository.GmmPriceCodeRepository, db *gorm.DB) masterservice.GmmPriceCodeService {
	return &GmmPriceCodeServiceImpl{
		GmmPriceCodeRepo: gmmPriceCodeRepo,
		DB:               db,
	}
}

func (s *GmmPriceCodeServiceImpl) GetAllGmmPriceCode() ([]masterpayloads.GmmPriceCodeResponse, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := s.GmmPriceCodeRepo.GetAllGmmPriceCode(tx)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) GetGmmPriceCodeById(id int) (masterpayloads.GmmPriceCodeResponse, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := s.GmmPriceCodeRepo.GetGmmPriceCodeById(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) GetGmmPriceCodeDropdown() ([]masterpayloads.GmmPriceCodeDropdownResponse, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := s.GmmPriceCodeRepo.GetGmmPriceCodeDropdown(tx)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) SaveGmmPriceCode(req masterpayloads.GmmPriceCodeSaveRequest) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := s.GmmPriceCodeRepo.SaveGmmPriceCode(tx, req)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) UpdateGmmPriceCode(id int, req masterpayloads.GmmPriceCodeUpdateRequest) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := s.GmmPriceCodeRepo.UpdateGmmPriceCode(tx, id, req)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) ChangeStatusGmmPriceCode(id int) (masterentities.GmmPriceCode, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := s.GmmPriceCodeRepo.ChangeStatusGmmPriceCode(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *GmmPriceCodeServiceImpl) DeleteGmmPriceCode(id int) (bool, *exceptions.BaseErrorResponse) {
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
			} else {
				logrus.Info("Transaction committed successfully")
			}
		}
	}()
	results, err := s.GmmPriceCodeRepo.DeleteGmmPriceCode(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}
