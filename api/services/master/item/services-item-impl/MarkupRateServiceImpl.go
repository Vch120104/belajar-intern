package masteritemserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupRateServiceImpl struct {
	markupRepo masteritemrepository.MarkupRateRepository
	DB         *gorm.DB
}

func StartMarkupRateService(markupRepo masteritemrepository.MarkupRateRepository, db *gorm.DB) masteritemservice.MarkupRateService {
	return &MarkupRateServiceImpl{
		markupRepo: markupRepo,
		DB:         db,
	}
}

func (s *MarkupRateServiceImpl) GetAllMarkupRate(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.markupRepo.GetAllMarkupRate(tx, filterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *MarkupRateServiceImpl) GetMarkupRateById(id int) (masteritempayloads.MarkupRateResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.markupRepo.GetMarkupRateById(tx, id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *MarkupRateServiceImpl) SaveMarkupRate(req masteritempayloads.MarkupRateRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.MarkupRateId != 0 {
		_, err := s.markupRepo.GetMarkupRateById(tx, req.MarkupRateId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.markupRepo.SaveMarkupRate(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *MarkupRateServiceImpl) ChangeStatusMarkupRate(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.markupRepo.GetMarkupRateById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.markupRepo.ChangeStatusMarkupRate(tx, Id)
	if err != nil {
		return results, err
	}
	return true, nil
}
