package masterwarehouseserviceimpl

import (
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
}

func OpenWarehouseMasterService(warehouseMaster masterwarehouserepository.WarehouseMasterRepository) masterwarehouseservice.WarehouseMasterService {
	return &WarehouseMasterServiceImpl{
		warehouseMasterRepo: warehouseMaster,
	}
}

func (s *WarehouseMasterServiceImpl) WithTrx(trxHandle *gorm.DB) masterwarehouseservice.WarehouseMasterService {
	s.warehouseMasterRepo = s.warehouseMasterRepo.WithTrx(trxHandle)
	return s
}

func (s *WarehouseMasterServiceImpl) Save(request masterwarehousepayloads.GetWarehouseMasterResponse) (bool, error) {
	save, err := s.warehouseMasterRepo.Save(request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *WarehouseMasterServiceImpl) GetById(warehouseId int) (masterwarehousepayloads.GetWarehouseMasterResponse, error) {
	get, err := s.warehouseMasterRepo.GetById(warehouseId)

	if err != nil {
		return masterwarehousepayloads.GetWarehouseMasterResponse{}, err
	}

	return get, nil
}

func (s *WarehouseMasterServiceImpl) GetAllIsActive() ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, error) {
	get, err := s.warehouseMasterRepo.GetAllIsActive()

	if err != nil {
		return nil, err
	}

	return get, nil
}

func (s *WarehouseMasterServiceImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseMasterRequest, pages pagination.Pagination) (pagination.Pagination, error) {
	get, err := s.warehouseMasterRepo.GetAll(request, pages)

	if err != nil {
		return pagination.Pagination{}, err
	}

	return get, nil
}

func (s *WarehouseMasterServiceImpl) GetWarehouseMasterByCode(Code string) ([]map[string]interface{}, error) {
	get, err := s.warehouseMasterRepo.GetWarehouseMasterByCode(Code)

	if err != nil {
		return nil, err
	}

	return get, nil
}

func (s *WarehouseMasterServiceImpl) ChangeStatus(warehouseId int) (masterwarehousepayloads.GetWarehouseMasterResponse, error) {
	change_status, err := s.warehouseMasterRepo.ChangeStatus(warehouseId)

	if err != nil {
		log.Panic(err.Error())
	}

	return change_status, nil
}
