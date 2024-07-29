package masterwarehouseserviceimpl

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"after-sales/api/utils"

	// "log"

	// "after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WarehouseMasterServiceImpl struct {
	warehouseMasterRepo masterwarehouserepository.WarehouseMasterRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client // Redis client
}

func OpenWarehouseMasterService(warehouseMaster masterwarehouserepository.WarehouseMasterRepository, db *gorm.DB, redisClient *redis.Client) masterwarehouseservice.WarehouseMasterService {
	return &WarehouseMasterServiceImpl{
		warehouseMasterRepo: warehouseMaster,
		DB:                  db,
		RedisClient:         redisClient,
	}
}

// GetWarehouseGroupbyCodeandCompanyId implements masterwarehouseservice.WarehouseMasterService.
func (s *WarehouseMasterServiceImpl) GetWarehouseGroupAndMasterbyCodeandCompanyId(companyId int, warehouseCode string) (int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	groupId, warehouseId, err := s.warehouseMasterRepo.GetWarehouseGroupAndMasterbyCodeandCompanyId(tx, companyId, warehouseCode)

	defer helper.CommitOrRollback(tx, err)

	return groupId, warehouseId, nil
}

// IsWarehouseMasterByCodeAndCompanyIdExist implements masterwarehouseservice.WarehouseMasterService.
func (s *WarehouseMasterServiceImpl) IsWarehouseMasterByCodeAndCompanyIdExist(companyId int, warehouseCode string) bool {
	tx := s.DB.Begin()
	isExist, err := s.warehouseMasterRepo.IsWarehouseMasterByCodeAndCompanyIdExist(tx, companyId, warehouseCode)

	defer helper.CommitOrRollback(tx, err)

	return isExist
}

// DropdownbyGroupId implements masterwarehouseservice.WarehouseMasterService.
func (s *WarehouseMasterServiceImpl) DropdownbyGroupId(warehouseGroupId int) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.warehouseMasterRepo.DropdownbyGroupId(tx, warehouseGroupId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseMasterServiceImpl) Save(request masterwarehousepayloads.GetWarehouseMasterResponse) (masterwarehouseentities.WarehouseMaster, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	save, err := s.warehouseMasterRepo.Save(tx, request)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return masterwarehouseentities.WarehouseMaster{}, err
	}
	return save, nil
}

func (s *WarehouseMasterServiceImpl) GetById(warehouseId int) (masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.warehouseMasterRepo.GetById(tx, warehouseId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseMasterServiceImpl) DropdownWarehouse() ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.warehouseMasterRepo.DropdownWarehouse(tx)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseMasterServiceImpl) GetAllIsActive() ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.warehouseMasterRepo.GetAllIsActive(tx)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseMasterServiceImpl) GetWarehouseWithMultiId(MultiIds []string) ([]masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.warehouseMasterRepo.GetWarehouseWithMultiId(tx, MultiIds)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseMasterServiceImpl) GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.warehouseMasterRepo.GetAll(tx, filter, pages)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseMasterServiceImpl) GetWarehouseMasterByCode(Code string) (masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	get, err := s.warehouseMasterRepo.GetWarehouseMasterByCode(tx, Code)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *WarehouseMasterServiceImpl) ChangeStatus(warehouseId int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	change_status, err := s.warehouseMasterRepo.ChangeStatus(tx, warehouseId)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return change_status, err
	}
	return change_status, nil
}

func (s *WarehouseMasterServiceImpl) GetAuthorizeUser(pages pagination.Pagination,id int)(pagination.Pagination,*exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.warehouseMasterRepo.GetAuthorizeUser(tx,pages,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return pages,err
	}
	return result,nil
}

func (s *WarehouseMasterServiceImpl) PostAuthorizeUser(req masterwarehousepayloads.WarehouseAuthorize)(masterwarehousepayloads.WarehouseAuthorize,*exceptions.BaseErrorResponse){
	tx := s.DB.Begin()
	result,err := s.warehouseMasterRepo.PostAuthorizeUser(tx,req)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return masterwarehousepayloads.WarehouseAuthorize{},err
	}
	return result,err
}

func (s *WarehouseMasterServiceImpl) DeleteMultiIdAuthorizeUser(id string)(bool,*exceptions.BaseErrorResponse){
	tx :=s.DB.Begin()
	result,err := s.warehouseMasterRepo.DeleteMultiIdAuthorizeUser(tx,id)
	defer helper.CommitOrRollback(tx,err)
	if err != nil{
		return false,err
	}
	return result,nil
}