package masteroperationserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
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

func (s *LabourSellingPriceServiceImpl) GetLabourSellingPriceById(Id int) (map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.labourSellingPriceRepo.GetLabourSellingPriceById(tx, Id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *LabourSellingPriceServiceImpl) GetAllSellingPriceDetailByHeaderId(headerId int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.labourSellingPriceRepo.GetAllSellingPriceDetailByHeaderId(tx, headerId, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
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