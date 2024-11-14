package masterserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
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

func (s *LookupServiceImpl) ItemOprCodeByCode(linetypeId int, oprItemCode string, pages pagination.Pagination, filterCondition []utils.FilterCondition) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.ItemOprCodeByCode(tx, linetypeId, oprItemCode, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) ItemOprCodeByID(linetypeId int, oprItemId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.ItemOprCodeByID(tx, linetypeId, oprItemId, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) ItemOprCodeWithPrice(linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode int, whsGroup string, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.ItemOprCodeWithPrice(tx, linetypeId, companyId, oprItemCode, brandId, modelId, jobTypeId, variantId, currencyId, billCode, whsGroup, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) GetVehicleUnitMaster(brandId int, modelId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.GetVehicleUnitMaster(tx, brandId, modelId, pages, filterCondition)
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

func (s *LookupServiceImpl) GetCampaignMaster(companyId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.GetCampaignMaster(tx, companyId, pages, filterCondition)
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

func (s *LookupServiceImpl) CustomerByTypeAndAddress(pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.CustomerByTypeAndAddress(tx, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) CustomerByTypeAndAddressByID(customerId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.CustomerByTypeAndAddressByID(tx, customerId, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) CustomerByTypeAndAddressByCode(customerCode string, pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.CustomerByTypeAndAddressByCode(tx, customerCode, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) GetOprItemPrice(linetypeId int, companyId int, oprItemCode int, brandId int, modelId int, jobTypeId int, variantId int, currencyId int, billCode int, whsGroup string) (float64, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	price, baseErr := s.LookupRepo.GetOprItemPrice(tx, linetypeId, companyId, oprItemCode, brandId, modelId, jobTypeId, variantId, currencyId, billCode, whsGroup)
	if baseErr != nil {
		return 0, baseErr
	}

	return price, nil
}

func (s *LookupServiceImpl) GetLineTypeByItemCode(itemCode string) (int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lineType, baseErr := s.LookupRepo.GetLineTypeByItemCode(tx, itemCode)
	if baseErr != nil {
		return 0, baseErr
	}

	return lineType, nil
}

func (s *LookupServiceImpl) ListItemLocation(companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	warehouse, baseErr := s.LookupRepo.ListItemLocation(tx, companyId, filterCondition, pages)
	if baseErr != nil {
		return warehouse, baseErr
	}

	return warehouse, nil
}

func (s *LookupServiceImpl) WarehouseGroupByCompany(companyId int) ([]masterpayloads.WarehouseGroupByCompanyResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	warehouse, baseErr := s.LookupRepo.WarehouseGroupByCompany(tx, companyId)
	if baseErr != nil {
		return warehouse, baseErr
	}

	return warehouse, nil
}

func (s *LookupServiceImpl) ItemListTrans(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	item, baseErr := s.LookupRepo.ItemListTrans(tx, filterCondition, pages)
	if baseErr != nil {
		return item, baseErr
	}

	return item, nil
}

func (s *LookupServiceImpl) ItemListTransPL(companyId int, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	item, baseErr := s.LookupRepo.ItemListTransPL(tx, companyId, filterCondition, pages)
	if baseErr != nil {
		return item, baseErr
	}

	return item, nil
}

func (s *LookupServiceImpl) ReferenceTypeWorkOrder(pages pagination.Pagination, filterCondition []utils.FilterCondition) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.ReferenceTypeWorkOrder(tx, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}

func (s *LookupServiceImpl) ReferenceTypeWorkOrderByID(referenceId int, pages pagination.Pagination, filterCondition []utils.FilterCondition) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx, nil)

	lookup, totalPages, totalRows, baseErr := s.LookupRepo.ReferenceTypeWorkOrderByID(tx, referenceId, pages, filterCondition)
	if baseErr != nil {
		return nil, 0, 0, baseErr
	}

	return lookup, totalPages, totalRows, nil
}
