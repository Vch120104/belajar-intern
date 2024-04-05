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

func (s *ForecastMasterServiceImpl) GetForecastMasterById(id int) (masterpayloads.ForecastMasterResponse,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ForecastMasterRepo.GetForecastMasterById(tx, id)
	if err != nil {
		return results,err
	}
	return results,nil
}

func (s *ForecastMasterServiceImpl) SaveForecastMaster(req masterpayloads.ForecastMasterResponse) (bool,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.ForecastMasterId != 0 {
		_, err := s.ForecastMasterRepo.GetForecastMasterById(tx, req.ForecastMasterId)

		if err != nil {
			return false,err
		}
	}

	results, err := s.ForecastMasterRepo.SaveForecastMaster(tx, req)

	if err != nil {
		return false,err
	}
	return results,nil
}

func (s *ForecastMasterServiceImpl) ChangeStatusForecastMaster(Id int) (bool,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.ForecastMasterRepo.GetForecastMasterById(tx, Id)

	if err != nil {
		return false,err
	}

	results, err := s.ForecastMasterRepo.ChangeStatusForecastMaster(tx, Id)
	if err != nil {
		return results,nil
	}
	return true,nil
}

func (s *ForecastMasterServiceImpl) GetAllForecastMaster(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.ForecastMasterRepo.GetAllForecastMaster(tx, filterCondition, pages)
	if err != nil {
		return results,0,0,err
	}
	return results, totalPages, totalRows,nil
}
