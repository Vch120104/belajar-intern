package masterwarehouseserviceimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"

	"log"

	"gorm.io/gorm"
	// "log"
	// "after-sales/api/utils"
)

type WarehouseLocationServiceImpl struct {
	warehouseLocationRepo masterwarehouserepository.WarehouseLocationRepository
	DB                    *gorm.DB
}

func OpenWarehouseLocationService(warehouseLocation masterwarehouserepository.WarehouseLocationRepository, db *gorm.DB) masterwarehouseservice.WarehouseLocationService {
	return &WarehouseLocationServiceImpl{
		warehouseLocationRepo: warehouseLocation,
		DB:                    db,
	}
}

func (s *WarehouseLocationServiceImpl) Save(request masterwarehousepayloads.GetWarehouseLocationResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	save, err := s.warehouseLocationRepo.Save(tx, request)

	if err != nil {
		return false
	}

	return save
}

func (s *WarehouseLocationServiceImpl) GetById(warehouseLocationId int) masterwarehousepayloads.GetWarehouseLocationResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseLocationRepo.GetById(tx, warehouseLocationId)

	if err != nil {
		return masterwarehousepayloads.GetWarehouseLocationResponse{}
	}

	return get
}

func (s *WarehouseLocationServiceImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseLocationRequest, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseLocationRepo.GetAll(tx, request, pages)

	if err != nil {
		return pagination.Pagination{}
	}

	return get
}

func (s *WarehouseLocationServiceImpl) ChangeStatus(warehouseLocationId int) masterwarehousepayloads.GetWarehouseLocationResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	change_status, err := s.warehouseLocationRepo.ChangeStatus(tx, warehouseLocationId)

	if err != nil {
		log.Panic(err.Error())
	}

	return change_status
}
