package masterwarehouseserviceimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	// "log"
	// "after-sales/api/utils"
)

type WarehouseLocationDefinitionServiceImpl struct {
	WarehouseLocationDefinitionRepo masterwarehouserepository.WarehouseLocationDefinitionRepository
	DB                              *gorm.DB
	RedisClient                     *redis.Client // Redis client
}

func OpenWarehouseLocationDefinitionService(WarehouseLocationDefinition masterwarehouserepository.WarehouseLocationDefinitionRepository, db *gorm.DB, redisClient *redis.Client) masterwarehouseservice.WarehouseLocationDefinitionService {
	return &WarehouseLocationDefinitionServiceImpl{
		WarehouseLocationDefinitionRepo: WarehouseLocationDefinition,
		DB:                              db,
		RedisClient:                     redisClient,
	}
}

func (s *WarehouseLocationDefinitionServiceImpl) Save(request masterwarehousepayloads.WarehouseLocationDefinitionResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if request.WarehouseLocationDefinitionId != 0 {
		_, err := s.WarehouseLocationDefinitionRepo.GetById(tx, request.WarehouseLocationDefinitionId)

		if err != nil {
			return false, err
		}
	}

	save, err := s.WarehouseLocationDefinitionRepo.Save(tx, request)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return false, err
	}
	return save, err
}

func (s *WarehouseLocationDefinitionServiceImpl) SaveData(request masterwarehousepayloads.WarehouseLocationDefinitionResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	if request.WarehouseLocationDefinitionId != 0 {
		_, err := s.WarehouseLocationDefinitionRepo.GetById(tx, request.WarehouseLocationDefinitionId)

		if err != nil {
			return false, err
		}
	}

	save, err := s.WarehouseLocationDefinitionRepo.SaveData(tx, request)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return false, err
	}
	return save, err
}

func (s *WarehouseLocationDefinitionServiceImpl) GetById(Id int) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.WarehouseLocationDefinitionRepo.GetById(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *WarehouseLocationDefinitionServiceImpl) GetByLevel(idlevel int, idwhl string) (masterwarehousepayloads.WarehouseLocationDefinitionResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.WarehouseLocationDefinitionRepo.GetByLevel(tx, idlevel, idwhl)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *WarehouseLocationDefinitionServiceImpl) GetAll(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.WarehouseLocationDefinitionRepo.GetAll(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *WarehouseLocationDefinitionServiceImpl) ChangeStatus(Id int) (masterwarehouseentities.WarehouseLocationDefinition, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	// Ubah status
	entity, err := s.WarehouseLocationDefinitionRepo.ChangeStatus(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterwarehouseentities.WarehouseLocationDefinition{}, err
	}
	return entity, nil
}

func (s *WarehouseLocationDefinitionServiceImpl) PopupWarehouseLocationLevel(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.WarehouseLocationDefinitionRepo.PopupWarehouseLocationLevel(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}
