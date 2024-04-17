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
	skillLevelRepo masterrepository.SkillLevelRepository
	DB             *gorm.DB
}

func StartSkillLevelService(skillLevelRepo masterrepository.SkillLevelRepository, db *gorm.DB) masterservice.SkillLevelService {
	return &SkillLevelServiceImpl{
		skillLevelRepo: skillLevelRepo,
		DB:             db,
	}
}

func (s *SkillLevelServiceImpl) GetAllSkillLevel(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.skillLevelRepo.GetAllSkilllevel(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *SkillLevelServiceImpl) GetSkillLevelById(Id int) (masterpayloads.SkillLevelResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.skillLevelRepo.GetSkillLevelById(tx, Id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *SkillLevelServiceImpl) SaveSkillLevel(req masterpayloads.SkillLevelResponse) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.SkillLevelId != 0 {
		_, err := s.skillLevelRepo.GetSkillLevelById(tx, req.SkillLevelId)
		if err != nil {
			return false, err
		}
	}

	results, err := s.skillLevelRepo.SaveSkillLevel(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *SkillLevelServiceImpl) ChangeStatusSkillLevel(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.skillLevelRepo.GetSkillLevelById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.skillLevelRepo.ChangeStatusSkillLevel(tx, Id)
	if err != nil {
		return results, err
	}
	return true, nil
}