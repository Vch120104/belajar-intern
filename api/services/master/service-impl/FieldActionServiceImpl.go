package masterserviceimpl

import (
	// "after-sales/api/exceptions"

	exceptions "after-sales/api/exceptions"
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

func (s *FieldActionServiceImpl) GetAllFieldAction(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.GetAllFieldAction(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *FieldActionServiceImpl) SaveFieldAction(req masterpayloads.FieldActionRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.SaveFieldAction(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *FieldActionServiceImpl) GetFieldActionHeaderById(Id int) (masterpayloads.FieldActionResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.GetFieldActionHeaderById(tx, Id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *FieldActionServiceImpl) GetAllFieldActionVehicleDetailById(Id int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	pages, err := s.FieldActionRepo.GetAllFieldActionVehicleDetailById(tx, Id, pages, filterCondition)
	if err != nil {
		return pages, err
	}
	return pages, nil
}

func (s *FieldActionServiceImpl) GetFieldActionVehicleDetailById(Id int) (masterpayloads.FieldActionDetailResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.GetFieldActionVehicleDetailById(tx, Id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *FieldActionServiceImpl) GetAllFieldActionVehicleItemDetailById(Id int, pages pagination.Pagination) ([]map[string]interface{},int,int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result,totalpage,totalrows, err := s.FieldActionRepo.GetAllFieldActionVehicleItemDetailById(tx, Id, pages)
	if err != nil {
		return result,totalpage,totalrows, err
	}
	return result,totalpage,totalrows, nil
}

func (s *FieldActionServiceImpl) GetFieldActionVehicleItemDetailById(Id int, linetypeid int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.GetFieldActionVehicleItemDetailById(tx, Id, linetypeid)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *FieldActionServiceImpl) PostFieldActionVehicleItemDetail(Id int, req masterpayloads.FieldActionItemDetailResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.PostFieldActionVehicleItemDetail(tx, req, Id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *FieldActionServiceImpl) PostFieldActionVehicleDetail(Id int, req masterpayloads.FieldActionDetailResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.PostFieldActionVehicleDetail(tx, req, Id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *FieldActionServiceImpl) PostMultipleVehicleDetail(headerId int, id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.PostMultipleVehicleDetail(tx, headerId, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *FieldActionServiceImpl) PostVehicleItemIntoAllVehicleDetail(headerId int, req masterpayloads.FieldActionItemDetailResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.PostVehicleItemIntoAllVehicleDetail(tx, headerId, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *FieldActionServiceImpl) ChangeStatusFieldAction(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.FieldActionRepo.ChangeStatusFieldAction(tx, id)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *FieldActionServiceImpl) ChangeStatusFieldActionVehicle(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.FieldActionRepo.ChangeStatusFieldActionVehicle(tx, id)

	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *FieldActionServiceImpl) ChangeStatusFieldActionVehicleItem(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	result, err := s.FieldActionRepo.ChangeStatusFieldActionVehicleItem(tx, id)

	if err != nil {
		return result, err
	}
	return result, nil
}
