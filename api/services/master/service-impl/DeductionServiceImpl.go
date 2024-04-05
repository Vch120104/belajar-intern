package masterserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"

	// "after-sales/api/payloads"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DeductionServiceImpl struct {
	deductionrepo masterrepository.DeductionRepository
	DB            *gorm.DB
}

func StartDeductionService(deductionRepo masterrepository.DeductionRepository, db *gorm.DB) masterservice.DeductionService {
	return &DeductionServiceImpl{
		deductionrepo: deductionRepo,
		DB:            db,
	}
}

func (s *DeductionServiceImpl) GetAllDeduction(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.deductionrepo.GetAllDeduction(tx, filterCondition, pages)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *DeductionServiceImpl) GetByIdDeductionDetail(Id int) (masterpayloads.DeductionDetailResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.deductionrepo.GetByIdDeductionDetail(tx, Id)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *DeductionServiceImpl) PostDeductionList(req masterpayloads.DeductionListResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.deductionrepo.SaveDeductionList(tx, req)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *DeductionServiceImpl) PostDeductionDetail(req masterpayloads.DeductionDetailResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.deductionrepo.SaveDeductionDetail(tx, req)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *DeductionServiceImpl) GetDeductionById(Id int) (masterpayloads.DeductionListResponse, *exceptionsss_test.BaseErrorResponse) {

	tx := s.DB.Begin()

	defer helper.CommitOrRollback(tx)

	result, err := s.deductionrepo.GetDeductionById(tx, Id)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *DeductionServiceImpl) GetAllDeductionDetail(Id int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	detail_result, detail_err := s.deductionrepo.GetAllDeductionDetail(tx, pages, Id)

	if detail_err != nil {
		return detail_result, detail_err
	}

	return detail_result, nil
}

func (s *DeductionServiceImpl) ChangeStatusDeduction(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.deductionrepo.GetDeductionById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.deductionrepo.ChangeStatusDeduction(tx, Id)
	if err != nil {
		return results, err
	}
	return true, nil
}
