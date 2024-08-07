package masteroperationserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	"after-sales/api/payloads/pagination"
	masteroperationrepository "after-sales/api/repositories/master/operation"
	masteroperationservice "after-sales/api/services/master/operation"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type LabourSellingPriceServiceImpl struct {
	labourSellingPriceRepo masteroperationrepository.LabourSellingPriceRepository
	DB                     *gorm.DB
}

func StartLabourSellingPriceService(labourSellingPriceRepo masteroperationrepository.LabourSellingPriceRepository, db *gorm.DB) masteroperationservice.LabourSellingPriceService {
	return &LabourSellingPriceServiceImpl{
		labourSellingPriceRepo: labourSellingPriceRepo,
		DB:                     db,
	}
}

// GetAllSellingPrice implements masteroperationservice.LabourSellingPriceService.
func (s *LabourSellingPriceServiceImpl) GetAllSellingPrice(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.labourSellingPriceRepo.GetAllSellingPrice(tx, filter, pages)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *LabourSellingPriceServiceImpl) GetLabourSellingPriceById(Id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.labourSellingPriceRepo.GetLabourSellingPriceById(tx, Id)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *LabourSellingPriceServiceImpl) GetAllSellingPriceDetailByHeaderId(headerId int, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.labourSellingPriceRepo.GetAllSellingPriceDetailByHeaderId(tx, headerId, pages)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *LabourSellingPriceServiceImpl) SaveLabourSellingPrice(req masteroperationpayloads.LabourSellingPriceRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.labourSellingPriceRepo.SaveLabourSellingPrice(tx, req)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *LabourSellingPriceServiceImpl) SaveLabourSellingPriceDetail(req masteroperationpayloads.LabourSellingPriceDetailRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.labourSellingPriceRepo.SaveLabourSellingPriceDetail(tx, req)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return false, err
	}
	return results, nil
}
