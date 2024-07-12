package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SkillLevelServiceImpl struct {
	SkillLevelRepo masterrepository.SkillLevelRepository
	DB             *gorm.DB
	RedisClient    *redis.Client // Redis client
}

func StartSkillLevelService(SkillLevelRepo masterrepository.SkillLevelRepository, db *gorm.DB, redisClient *redis.Client) masterservice.SkillLevelService {
	return &SkillLevelServiceImpl{
		SkillLevelRepo: SkillLevelRepo,
		DB:             db,
		RedisClient:    redisClient,
	}
}

func (s *SkillLevelServiceImpl) GetSkillLevelById(id int) (masterpayloads.SkillLevelResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.SkillLevelRepo.GetSkillLevelById(tx, id)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *SkillLevelServiceImpl) GetAllSkillLevel(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.SkillLevelRepo.GetAllSkillLevel(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *SkillLevelServiceImpl) ChangeStatusSkillLevel(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.SkillLevelRepo.GetSkillLevelById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.SkillLevelRepo.ChangeStatusSkillLevel(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return true, nil
}

func (s *SkillLevelServiceImpl) SaveSkillLevel(req masterpayloads.SkillLevelResponse) (masterentities.SkillLevel, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.SkillLevelRepo.SaveSkillLevel(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.SkillLevel{}, err
	}
	return results, nil
}
