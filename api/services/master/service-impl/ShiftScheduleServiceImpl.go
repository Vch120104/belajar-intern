package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ShiftScheduleServiceImpl struct {
	ShiftScheduleRepo masterrepository.ShiftScheduleRepository
	DB                *gorm.DB
}

func StartShiftScheduleService(ShiftScheduleRepo masterrepository.ShiftScheduleRepository, db *gorm.DB) masterservice.ShiftScheduleService {
	return &ShiftScheduleServiceImpl{
		ShiftScheduleRepo: ShiftScheduleRepo,
		DB:                db,
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

func (s *ShiftScheduleServiceImpl) GetShiftScheduleById(id int) masterpayloads.ShiftScheduleResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ShiftScheduleRepo.GetShiftScheduleById(tx, id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
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

func (s *ShiftScheduleServiceImpl) GetAllShiftSchedule(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ShiftScheduleRepo.GetAllShiftSchedule(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *ShiftScheduleServiceImpl) ChangeStatusShiftSchedule(oprId int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.ShiftScheduleRepo.GetShiftScheduleById(tx, oprId)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.ShiftScheduleRepo.ChangeStatusShiftSchedule(tx, oprId)
	if err != nil {
		return results
	}
	return true
}

func (s *ShiftScheduleServiceImpl) SaveShiftSchedule(req masterpayloads.ShiftScheduleResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.ShiftScheduleId != 0 {
		_, err := s.ShiftScheduleRepo.GetShiftScheduleById(tx, req.ShiftScheduleId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.ShiftScheduleRepo.SaveShiftSchedule(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
