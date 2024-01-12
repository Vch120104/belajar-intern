package masterwarehouseserviceimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"

	"log"

	// "log"

	// "after-sales/api/utils"

	"gorm.io/gorm"
)

type WarehouseLocationServiceImpl struct {
	warehouseLocationRepo masterwarehouserepository.WarehouseLocationRepository
}

func OpenWarehouseLocationService(warehouseLocation masterwarehouserepository.WarehouseLocationRepository) masterwarehouseservice.WarehouseLocationService {
	return &WarehouseLocationServiceImpl{
		warehouseLocationRepo: warehouseLocation,
	}
}

func (s *WarehouseLocationServiceImpl) WithTrx(trxHandle *gorm.DB) masterwarehouseservice.WarehouseLocationService {
	s.warehouseLocationRepo = s.warehouseLocationRepo.WithTrx(trxHandle)
	return s
}

func (s *WarehouseLocationServiceImpl) Save(request masterwarehousepayloads.GetWarehouseLocationResponse) (bool, error) {
	save, err := s.warehouseLocationRepo.Save(request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *WarehouseLocationServiceImpl) GetById(warehouseLocationId int) (masterwarehousepayloads.GetWarehouseLocationResponse, error) {
	get, err := s.warehouseLocationRepo.GetById(warehouseLocationId)

	if err != nil {
		return masterwarehousepayloads.GetWarehouseLocationResponse{}, err
	}

	return get, nil
}

func (s *WarehouseLocationServiceImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseLocationRequest, pages pagination.Pagination) (pagination.Pagination, error) {
	get, err := s.warehouseLocationRepo.GetAll(request, pages)

	if err != nil {
		return pagination.Pagination{}, err
	}

	return get, nil
}

func (s *WarehouseLocationServiceImpl) ChangeStatus(warehouseLocationId int) (masterwarehousepayloads.GetWarehouseLocationResponse, error) {
	change_status, err := s.warehouseLocationRepo.ChangeStatus(warehouseLocationId)

	if err != nil {
		log.Panic(err.Error())
	}

	return change_status, nil
}