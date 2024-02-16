package masteritemserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type UnitOfMeasurementServiceImpl struct {
	unitOfMeasurementRepo masteritemrepository.UnitOfMeasurementRepository
	DB                    *gorm.DB
}

func StartUnitOfMeasurementService(unitOfMeasurementRepo masteritemrepository.UnitOfMeasurementRepository, db *gorm.DB) masteritemservice.UnitOfMeasurementService {
	return &UnitOfMeasurementServiceImpl{
		unitOfMeasurementRepo: unitOfMeasurementRepo,
		DB:                    db,
	}
}

func (s *UnitOfMeasurementServiceImpl) GetAllUnitOfMeasurementIsActive() []masteritempayloads.UomResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.unitOfMeasurementRepo.GetAllUnitOfMeasurementIsActive(tx)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *UnitOfMeasurementServiceImpl) GetUnitOfMeasurementById(id int) masteritempayloads.UomIdCodeResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementById(tx, id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *UnitOfMeasurementServiceImpl) GetUnitOfMeasurementByCode(Code string) masteritempayloads.UomIdCodeResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementByCode(tx, Code)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *UnitOfMeasurementServiceImpl) GetAllUnitOfMeasurement(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.unitOfMeasurementRepo.GetAllUnitOfMeasurement(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *UnitOfMeasurementServiceImpl) ChangeStatusUnitOfMeasurement(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementById(tx, Id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.unitOfMeasurementRepo.ChangeStatusUnitOfMeasurement(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *UnitOfMeasurementServiceImpl) SaveUnitOfMeasurement(req masteritempayloads.UomResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.UomId != 0 {
		_, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementById(tx, req.UomId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.unitOfMeasurementRepo.SaveUnitOfMeasurement(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
