package masterserviceimpl

import (
	// "after-sales/api/exceptions"
	"after-sales/api/exceptions"
	"after-sales/api/helper"

	// masterpayloads "after-sales/api/payloads/master"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type FieldActionServiceImpl struct {
	FieldActionRepo masterrepository.FieldActionRepository
	DB              *gorm.DB
}

func StartFieldActionService(FieldActionRepo masterrepository.FieldActionRepository, db *gorm.DB) masterservice.FieldActionService {
	return &FieldActionServiceImpl{
		FieldActionRepo: FieldActionRepo,
		DB:              db,
	}
}

func (s *FieldActionServiceImpl) GetAllFieldAction(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.GetAllFieldAction(tx, filterCondition, pages)
	if err != nil {
		return pages
	}
	return results
}

func (s *FieldActionServiceImpl) SaveFieldAction(req masterpayloads.FieldActionResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.SaveFieldAction(tx, req)
	if err != nil {
		return false
	}
	return results
}

func (s *FieldActionServiceImpl) GetFieldActionHeaderById(Id int) masterpayloads.FieldActionResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.GetFieldActionHeaderById(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *FieldActionServiceImpl) GetAllFieldActionVehicleDetailById(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	pages, err := s.FieldActionRepo.GetAllFieldActionVehicleDetailById(tx, Id, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return pages
}

func (s *FieldActionServiceImpl) GetFieldActionVehicleDetailById(Id int) masterpayloads.FieldActionDetailResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.GetFieldActionVehicleDetailById(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *FieldActionServiceImpl) GetAllFieldActionVehicleItemDetailById(Id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	pages, err := s.FieldActionRepo.GetAllFieldActionVehicleItemDetailById(tx, Id, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return pages
}

func (s *FieldActionServiceImpl) GetFieldActionVehicleItemDetailById(Id int) masterpayloads.FieldActionItemDetailResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.GetFieldActionVehicleItemDetailById(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *FieldActionServiceImpl) PostFieldActionVehicleItemDetail(Id int, req masterpayloads.FieldActionItemDetailResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.PostFieldActionVehicleItemDetail(tx, req, Id)
	if err != nil {
		return false
	}
	return results
}

func (s *FieldActionServiceImpl) PostFieldActionVehicleDetail(Id int, req masterpayloads.FieldActionDetailResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.PostFieldActionVehicleDetail(tx, req, Id)
	if err != nil {
		return false
	}
	return results
}

func (s *FieldActionServiceImpl) PostMultipleVehicleDetail(headerId int, companyId int, id string) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.PostMultipleVehicleDetail(tx, headerId, companyId, id)
	if err != nil {
		return false
	}
	return results
}

func (s *FieldActionServiceImpl) PostVehicleItemIntoAllVehicleDetail(headerId int, req masterpayloads.FieldActionItemDetailResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.PostVehicleItemIntoAllVehicleDetail(tx, headerId, req)
	if err != nil {
		return false
	}
	return results
}
