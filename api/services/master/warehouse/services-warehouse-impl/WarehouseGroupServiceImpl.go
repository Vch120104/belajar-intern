package masterwarehouseserviceimpl

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"log"

	// "after-sales/api/utils"

	"gorm.io/gorm"
)

type WarehouseGroupServiceImpl struct {
	warehouseGroupRepo masterwarehouserepository.WarehouseGroupRepository
}

func OpenWarehouseGroupService(warehouseGroup masterwarehouserepository.WarehouseGroupRepository) masterwarehouseservice.WarehouseGroupService {
	return &WarehouseGroupServiceImpl{
		warehouseGroupRepo: warehouseGroup,
	}
}

func (s *WarehouseGroupServiceImpl) WithTrx(trxHandle *gorm.DB) masterwarehouseservice.WarehouseGroupService {
	s.warehouseGroupRepo = s.warehouseGroupRepo.WithTrx(trxHandle)
	return s
}

func (s *WarehouseGroupServiceImpl) Save(request masterwarehousepayloads.GetWarehouseGroupResponse) (bool, error) {
	save, err := s.warehouseGroupRepo.Save(request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *WarehouseGroupServiceImpl) GetById(warehouseGroupId int) (masterwarehousepayloads.GetWarehouseGroupResponse, error) {
	get, err := s.warehouseGroupRepo.GetById(warehouseGroupId)

	if err != nil {
		return masterwarehousepayloads.GetWarehouseGroupResponse{}, err
	}

	return get, nil
}

func (s *WarehouseGroupServiceImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseGroupRequest) ([]masterwarehousepayloads.GetWarehouseGroupResponse, error) {
	get, err := s.warehouseGroupRepo.GetAll(request)

	if err != nil {
		return nil, err
	}

	return get, nil
}

func (s *WarehouseGroupServiceImpl) ChangeStatus(warehouseGroupId int) (masterwarehousepayloads.GetWarehouseGroupResponse, error) {
	change_status, err := s.warehouseGroupRepo.ChangeStatus(warehouseGroupId)

	if err != nil {
		log.Panic(err.Error())
	}

	return change_status, nil
}
