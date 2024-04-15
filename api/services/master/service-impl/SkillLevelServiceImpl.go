package masterserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

	"after-sales/api/utils"

	"gorm.io/gorm"
)

type SkillLevelServiceImpl struct {
	SkillLevelRepo masterrepository.SkillLevelRepository
	DB             *gorm.DB
}

func StartSkillLevelService(SkillLevelRepo masterrepository.SkillLevelRepository, db *gorm.DB) masterservice.SkillLevelService {
	return &SkillLevelServiceImpl{
		SkillLevelRepo: SkillLevelRepo,
		DB:             db,
	}
}

func (s *SkillLevelServiceImpl) GetSkillLevelById(id int) (masterpayloads.SkillLevelResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.SkillLevelRepo.GetSkillLevelById(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *SkillLevelServiceImpl) GetAllSkillLevel(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.SkillLevelRepo.GetAllSkillLevel(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *SkillLevelServiceImpl) ChangeStatusSkillLevel(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.SkillLevelRepo.GetSkillLevelById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.SkillLevelRepo.ChangeStatusSkillLevel(tx, Id)
	if err != nil {
		return results, err
	}
	return true, nil
}

func (s *SkillLevelServiceImpl) SaveSkillLevel(req masterpayloads.SkillLevelResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.SkillLevelCodeId != 0 {
		_, err := s.SkillLevelRepo.GetSkillLevelById(tx, req.SkillLevelCodeId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.SkillLevelRepo.SaveSkillLevel(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}
