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

type ForecastMasterServiceImpl struct {
	ForecastMasterRepo masterrepository.ForecastMasterRepository
	DB                 *gorm.DB
}

func StartForecastMasterService(ForecastMasterRepo masterrepository.ForecastMasterRepository, db *gorm.DB) masterservice.ForecastMasterService {
	return &ForecastMasterServiceImpl{
		ForecastMasterRepo: ForecastMasterRepo,
		DB:                 db,
	}
}

func (s *ForecastMasterServiceImpl) GetForecastMasterById(id int) masterpayloads.ForecastMasterResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ForecastMasterRepo.GetForecastMasterById(tx, id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *ForecastMasterServiceImpl) SaveForecastMaster(req masterpayloads.ForecastMasterResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.ForecastMasterId != 0 {
		_, err := s.ForecastMasterRepo.GetForecastMasterById(tx, req.ForecastMasterId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.ForecastMasterRepo.SaveForecastMaster(tx, req)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *ForecastMasterServiceImpl) ChangeStatusForecastMaster(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.ForecastMasterRepo.GetForecastMasterById(tx, Id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.ForecastMasterRepo.ChangeStatusForecastMaster(tx, Id)
	if err != nil {
		return results
	}
	return true
}

func (s *ForecastMasterServiceImpl) GetAllForecastMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.ForecastMasterRepo.GetAllForecastMaster(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results, totalPages, totalRows
}
