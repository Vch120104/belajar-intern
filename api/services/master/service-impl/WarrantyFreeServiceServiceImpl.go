package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

	"gorm.io/gorm"
)

type WarrantyFreeServiceServiceImpl struct {
	warrantyFreeServiceRepo masterrepository.WarrantyFreeServiceRepository
	DB                      *gorm.DB
}

func StartWarrantyFreeServiceService(warrantyFreeServiceRepo masterrepository.WarrantyFreeServiceRepository, db *gorm.DB) masterservice.WarrantyFreeServiceService {
	return &WarrantyFreeServiceServiceImpl{
		warrantyFreeServiceRepo: warrantyFreeServiceRepo,
		DB:                  db,
	}
}

func (s *WarrantyFreeServiceServiceImpl) GetWarrantyFreeServiceById(Id int) map[string]interface{} {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.warrantyFreeServiceRepo.GetWarrantyFreeServiceById(tx, Id)
	
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *WarrantyFreeServiceServiceImpl) SaveWarrantyFreeService(req masterpayloads.WarrantyFreeServiceRequest) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.WarrantyFreeServicesId != 0 {
		_, err := s.warrantyFreeServiceRepo.GetWarrantyFreeServiceById(tx, req.WarrantyFreeServicesId)
		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.warrantyFreeServiceRepo.SaveWarrantyFreeService(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}