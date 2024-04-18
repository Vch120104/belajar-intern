package masteroperationserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"

	"gorm.io/gorm"
)

type LabourSellingPriceServiceImpl struct {
	labourSellingPriceRepo masteroperationrepository.LabourSellingPriceRepository
	DB                     *gorm.DB
}

func StartLabourSellingPriceService(labourSellingPriceRepo masteroperationrepository.LabourSellingPriceRepository, db *gorm.DB) masteroperationservice.LabourSellingPriceService {
	return &LabourSellingPriceServiceImpl{
		labourSellingPriceRepo: labourSellingPriceRepo,
		DB:             db,
	}
}

func (s *LabourSellingPriceServiceImpl) SaveLabourSellingPrice(req masteroperationpayloads.LabourSellingPriceRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	results, err := s.labourSellingPriceRepo.SaveLabourSellingPrice(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *LabourSellingPriceServiceImpl) SaveLabourSellingPriceDetail(req masteroperationpayloads.LabourSellingPriceDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	results, err := s.labourSellingPriceRepo.SaveLabourSellingPriceDetail(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}