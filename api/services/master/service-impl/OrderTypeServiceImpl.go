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

type OrderTypeServiceImpl struct {
	OrderTypeRepo masterrepository.OrderTypeRepository
	DB            *gorm.DB
}

func StartOrderTypeServiceImpl(orderTypeRepo masterrepository.OrderTypeRepository, db *gorm.DB) masterservice.OrderTypeService {
	return &OrderTypeServiceImpl{
		OrderTypeRepo: orderTypeRepo,
		DB:            db,
	}
}

func (s *OrderTypeServiceImpl) GetAllOrderType() ([]masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.OrderTypeRepo.GetAllOrderType(tx)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OrderTypeServiceImpl) GetOrderTypeById(id int) (masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.OrderTypeRepo.GetOrderTypeById(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OrderTypeServiceImpl) GetOrderTypeByName(name string) ([]masterpayloads.GetOrderTypeResponse, *exceptions.BaseErrorResponse) {
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
	results, err := s.OrderTypeRepo.GetOrderTypeByName(tx, name)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OrderTypeServiceImpl) SaveOrderType(req masterpayloads.OrderTypeSaveRequest) (masterentities.OrderType, *exceptions.BaseErrorResponse) {
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
	results, err := s.OrderTypeRepo.SaveOrderType(tx, req)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OrderTypeServiceImpl) UpdateOrderType(id int, req masterpayloads.OrderTypeUpdateRequest) (masterentities.OrderType, *exceptions.BaseErrorResponse) {
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
	results, err := s.OrderTypeRepo.UpdateOrderType(tx, id, req)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OrderTypeServiceImpl) ChangeStatusOrderType(id int) (masterentities.OrderType, *exceptions.BaseErrorResponse) {
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
	results, err := s.OrderTypeRepo.ChangeStatusOrderType(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *OrderTypeServiceImpl) DeleteOrderType(id int) (bool, *exceptions.BaseErrorResponse) {
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
	results, err := s.OrderTypeRepo.DeleteOrderType(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}
