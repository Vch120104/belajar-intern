package masterserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LookupServiceImpl struct {
	LookupRepo  masterrepository.LookupRepository
	DB          *gorm.DB
	RedisClient *redis.Client // Redis client
}

func StartLookupService(LookupRepo masterrepository.LookupRepository, db *gorm.DB, redisClient *redis.Client) masterservice.LookupService {
	return &LookupServiceImpl{
		LookupRepo:  LookupRepo,
		DB:          db,
		RedisClient: redisClient,
	}
}

func (s *LookupServiceImpl) ItemOprCode(linetypeId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.ItemOprCode(tx, linetypeId, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) ItemOprCodeWithPrice(linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode string, whsGroup string, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.ItemOprCodeWithPrice(tx, linetypeId, companyId, oprItemCode, brandId, modelId, jobTypeId, variantId, currencyId, billCode, whsGroup, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) VehicleUnitMaster(brandId int, modelId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.VehicleUnitMaster(tx, brandId, modelId, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) GetVehicleUnitByID(vehicleID int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.GetVehicleUnitByID(tx, vehicleID, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) GetVehicleUnitByChassisNumber(chassisNumber string, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.GetVehicleUnitByChassisNumber(tx, chassisNumber, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) CampaignMaster(companyId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.CampaignMaster(tx, companyId, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) WorkOrderService(pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.WorkOrderService(tx, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}
