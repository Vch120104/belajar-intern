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

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type FieldActionServiceImpl struct {
	FieldActionRepo masterrepository.FieldActionRepository
	DB              *gorm.DB
	RedisClient     *redis.Client // Redis client
}

func StartFieldActionService(FieldActionRepo masterrepository.FieldActionRepository, db *gorm.DB, redisClient *redis.Client) masterservice.FieldActionService {
	return &FieldActionServiceImpl{
		FieldActionRepo: FieldActionRepo,
		DB:              db,
		RedisClient:     redisClient,
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

func (s *FieldActionServiceImpl) GetAllFieldActionVehicleDetailById(Id int, pages pagination.Pagination, filterCondition []utils.FilterCondition) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	pages, err := s.FieldActionRepo.GetAllFieldActionVehicleDetailById(tx, Id, pages, filterCondition)
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

func (s *FieldActionServiceImpl) GetAllFieldActionVehicleItemDetailById(Id int, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	pages, err := s.FieldActionRepo.GetAllFieldActionVehicleItemDetailById(tx, Id, pages)
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

func (s *FieldActionServiceImpl) PostMultipleVehicleDetail(headerId int, id string) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.PostMultipleVehicleDetail(tx, headerId, id)
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

func (s *FieldActionServiceImpl) ChangeStatusFieldAction(id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.FieldActionRepo.ChangeStatusFieldAction(tx, id)

	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *FieldActionServiceImpl) ChangeStatusFieldActionVehicle(id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.FieldActionRepo.ChangeStatusFieldActionVehicle(tx, id)

	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}

func (s *FieldActionServiceImpl) ChangeStatusFieldActionVehicleItem(id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.FieldActionRepo.ChangeStatusFieldActionVehicleItem(tx, id)

	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}
	return result
}
