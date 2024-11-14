package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

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
	results, err := s.OrderTypeRepo.GetAllOrderType(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}
