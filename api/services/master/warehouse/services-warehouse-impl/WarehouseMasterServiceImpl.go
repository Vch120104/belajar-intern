package masterwarehouseserviceimpl

import (
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"log"

	// "log"

	// "after-sales/api/utils"

	"gorm.io/gorm"
)

type WarehouseMasterServiceImpl struct {
	warehouseMasterRepo masterwarehouserepository.WarehouseMasterRepository
	DB                  *gorm.DB
}

func OpenWarehouseMasterService(warehouseMaster masterwarehouserepository.WarehouseMasterRepository, db *gorm.DB) masterwarehouseservice.WarehouseMasterService {
	return &WarehouseMasterServiceImpl{
		warehouseMasterRepo: warehouseMaster,
		DB:                  db,
	}
}

func (s *WarehouseMasterServiceImpl) Save(request masterwarehousepayloads.GetWarehouseMasterResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	save, err := s.warehouseMasterRepo.Save(tx, request)

	if err != nil {
		return false
	}

	return save
}

func (s *WarehouseMasterServiceImpl) GetById(warehouseId int) masterwarehousepayloads.GetWarehouseMasterResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseMasterRepo.GetById(tx, warehouseId)

	if err != nil {
		return masterwarehousepayloads.GetWarehouseMasterResponse{}
	}

	return get
}

func (s *WarehouseMasterServiceImpl) GetAllIsActive() []masterwarehousepayloads.IsActiveWarehouseMasterResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseMasterRepo.GetAllIsActive(tx)

	if err != nil {
		return nil
	}

	return get
}

func (s *WarehouseMasterServiceImpl) GetWarehouseWithMultiId(MultiIds []string) []masterwarehousepayloads.GetAllWarehouseMasterResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseMasterRepo.GetWarehouseWithMultiId(tx, MultiIds)

	if err != nil {
		return nil
	}

	return get
}

func (s *WarehouseMasterServiceImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseMasterRequest, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseMasterRepo.GetAll(tx, request, pages)

	if err != nil {
		return pagination.Pagination{}
	}

	return get
}

func (s *WarehouseMasterServiceImpl) GetWarehouseMasterByCode(Code string) []map[string]interface{} {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseMasterRepo.GetWarehouseMasterByCode(tx, Code)

	if err != nil {
		return nil
	}

	return get
}

func (s *WarehouseMasterServiceImpl) ChangeStatus(warehouseId int) masterwarehousepayloads.GetWarehouseMasterResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	change_status, err := s.warehouseMasterRepo.ChangeStatus(tx, warehouseId)

	if err != nil {
		log.Panic(err.Error())
	}

	return change_status
}
