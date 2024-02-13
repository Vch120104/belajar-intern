package masterwarehouseserviceimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"

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

	if request.WarehouseLocationId != 0 {
		_, err := s.warehouseLocationRepo.GetById(tx, request.WarehouseLocationId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	save, err := s.warehouseLocationRepo.Save(tx, request)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return save
}

func (s *WarehouseLocationServiceImpl) GetById(warehouseLocationId int) masterwarehousepayloads.GetWarehouseLocationResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseLocationRepo.GetById(tx, warehouseLocationId)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return get
}

func (s *WarehouseLocationServiceImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseLocationRequest, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseLocationRepo.GetAll(tx, request, pages)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return get
}

func (s *WarehouseLocationServiceImpl) ChangeStatus(warehouseLocationId int) masterwarehousepayloads.GetWarehouseLocationResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.warehouseLocationRepo.GetById(tx, warehouseLocationId)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	change_status, err := s.warehouseLocationRepo.ChangeStatus(tx, warehouseLocationId)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return change_status
}
