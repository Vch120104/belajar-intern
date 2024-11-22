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

// GetSellingPriceDetailById implements masteroperationservice.LabourSellingPriceService.
func (s *LabourSellingPriceServiceImpl) GetSellingPriceDetailById(detailId int) (masteroperationpayloads.LabourSellingPriceDetailbyIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.labourSellingPriceRepo.GetSellingPriceDetailById(tx, detailId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

// SaveDuplicate implements masteroperationservice.LabourSellingPriceService.
func (s *LabourSellingPriceServiceImpl) SaveDuplicate(req masteroperationpayloads.SaveDuplicateLabourSellingPrice) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	var err *exceptions.BaseErrorResponse
	defer func() {
		helper.CommitOrRollback(tx, err)
	}()

	_, err = s.labourSellingPriceRepo.SaveLabourSellingPrice(tx, req.Header)

	if err != nil {
		return false, err
	}

	_, err = s.labourSellingPriceRepo.SaveMultipleDetail(tx, req.Detail)

	if err != nil {
		return false, err
	}

	return true, nil

}

// Duplicate implements masteroperationservice.LabourSellingPriceService.
func (s *LabourSellingPriceServiceImpl) Duplicate(headerId int) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.labourSellingPriceRepo.GetAllDetailbyHeaderId(tx, headerId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

// GetAllSellingPrice implements masteroperationservice.LabourSellingPriceService.
func (s *LabourSellingPriceServiceImpl) GetAllSellingPrice(filter []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.labourSellingPriceRepo.GetAllSellingPrice(tx, filter, pages)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
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

func (s *LabourSellingPriceServiceImpl) SaveLabourSellingPrice(req masteroperationpayloads.LabourSellingPriceRequest) (int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.labourSellingPriceRepo.SaveLabourSellingPrice(tx, req)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *LabourSellingPriceServiceImpl) SaveLabourSellingPriceDetail(req masteroperationpayloads.LabourSellingPriceDetailRequest) (int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.labourSellingPriceRepo.SaveLabourSellingPriceDetail(tx, req)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *LabourSellingPriceServiceImpl) DeleteLabourSellingPriceDetail(iddet []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	deletemultiid, err := s.labourSellingPriceRepo.DeleteLabourSellingPriceDetail(tx, iddet)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return deletemultiid, nil
}