package masterserviceimpl

import (
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

type ShiftScheduleServiceImpl struct {
	ShiftScheduleRepo masterrepository.ShiftScheduleRepository
	DB                *gorm.DB
	RedisClient       *redis.Client // Redis client
}

func StartShiftScheduleService(ShiftScheduleRepo masterrepository.ShiftScheduleRepository, db *gorm.DB, redisClient *redis.Client) masterservice.ShiftScheduleService {
	return &ShiftScheduleServiceImpl{
		ShiftScheduleRepo: ShiftScheduleRepo,
		DB:                db,
		RedisClient:       redisClient,
	}
}

// func (s *ShiftScheduleServiceImpl) GetAllShiftScheduleIsActive() []masterpayloads.ShiftScheduleResponse {
// 	tx := s.DB.Begin()
// 	defer helper.CommitOrRollback(tx)
// 	get, err := s.ShiftScheduleRepo.GetAllShiftScheduleIsActive(tx)

// 	if err != nil {
// 		panic(exceptions.NewAppExceptionError(err.Error()))
// 	}

// 	return get
// }

func (s *ShiftScheduleServiceImpl) GetShiftScheduleById(id int) (masterpayloads.ShiftScheduleResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.ShiftScheduleRepo.GetShiftScheduleById(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// func (s *ShiftScheduleServiceImpl) GetShiftScheduleByCode(Code string) masterpayloads.ShiftScheduleResponse {
// 	tx := s.DB.Begin()
// 	defer helper.CommitOrRollback(tx)
// 	results, err := s.ShiftScheduleRepo.GetShiftScheduleByCode(tx, Code)
// 	if err != nil {
// 		panic(exceptions.NewNotFoundError(err.Error()))
// 	}
// 	return results
// }

func (s *ShiftScheduleServiceImpl) GetAllShiftSchedule(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.ShiftScheduleRepo.GetAllShiftSchedule(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ShiftScheduleServiceImpl) ChangeStatusShiftSchedule(oprId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.ShiftScheduleRepo.GetShiftScheduleById(tx, oprId)

	if err != nil {
		return false, err
	}

	results, err := s.ShiftScheduleRepo.ChangeStatusShiftSchedule(tx, oprId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ShiftScheduleServiceImpl) SaveShiftSchedule(req masterpayloads.ShiftScheduleResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if req.ShiftScheduleId != 0 {
		_, err := s.ShiftScheduleRepo.GetShiftScheduleById(tx, req.ShiftScheduleId)

		if err != nil {
			return false, err
		}
	}

	results, err := s.ShiftScheduleRepo.SaveShiftSchedule(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ShiftScheduleServiceImpl) GetShiftScheduleDropDown() ([]masterpayloads.ShiftScheduleDropDownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.ShiftScheduleRepo.GetShiftScheduleDropDown(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, err
	}
	return results, nil
}
