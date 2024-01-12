package masteritemserviceimpl

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type UnitOfMeasurementServiceImpl struct {
	unitOfMeasurementRepo masteritemrepository.UnitOfMeasurementRepository
}

func StartUnitOfMeasurementService(unitOfMeasurementRepo masteritemrepository.UnitOfMeasurementRepository) masteritemservice.UnitOfMeasurementService {
	return &UnitOfMeasurementServiceImpl{
		unitOfMeasurementRepo: unitOfMeasurementRepo,
	}
}
func (s *UnitOfMeasurementServiceImpl) WithTrx(trxHandle *gorm.DB) masteritemservice.UnitOfMeasurementService {
	s.unitOfMeasurementRepo = s.unitOfMeasurementRepo.WithTrx(trxHandle)
	return s
}

func (s *UnitOfMeasurementServiceImpl) GetAllUnitOfMeasurementIsActive() ([]masteritempayloads.UomResponse, error) {
	results, err := s.unitOfMeasurementRepo.GetAllUnitOfMeasurementIsActive()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) GetUnitOfMeasurementById(id int) (masteritempayloads.UomResponse, error) {
	results, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementById(id)

	if err != nil {
		return masteritempayloads.UomResponse{}, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) GetUnitOfMeasurementByCode(Code string) (masteritempayloads.UomResponse, error) {
	results, err := s.unitOfMeasurementRepo.GetUnitOfMeasurementByCode(Code)
	if err != nil {
		return masteritempayloads.UomResponse{}, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) GetAllUnitOfMeasurement(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	results, err := s.unitOfMeasurementRepo.GetAllUnitOfMeasurement(filterCondition, pages)
	if err != nil {
		return pages, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) ChangeStatusUnitOfMeasurement(Id int) (bool, error) {
	results, err := s.unitOfMeasurementRepo.ChangeStatusUnitOfMeasurement(Id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *UnitOfMeasurementServiceImpl) SaveUnitOfMeasurement(req masteritempayloads.UomResponse) (bool, error) {
	results, err := s.unitOfMeasurementRepo.SaveUnitOfMeasurement(req)
	if err != nil {
		return false, err
	}
	return results, nil
}
