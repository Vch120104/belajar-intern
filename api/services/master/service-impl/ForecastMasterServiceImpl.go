package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

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
