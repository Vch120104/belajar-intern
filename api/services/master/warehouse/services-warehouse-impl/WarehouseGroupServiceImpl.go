package masterwarehouseserviceimpl

import (
	"after-sales/api/helper"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouseservice "after-sales/api/services/master/warehouse"
	"log"

	"gorm.io/gorm"
	// "after-sales/api/utils"
)

type WarehouseGroupServiceImpl struct {
	warehouseGroupRepo masterwarehouserepository.WarehouseGroupRepository
	DB                 *gorm.DB
}

func OpenWarehouseGroupService(warehouseGroup masterwarehouserepository.WarehouseGroupRepository, db *gorm.DB) masterwarehouseservice.WarehouseGroupService {
	return &WarehouseGroupServiceImpl{
		warehouseGroupRepo: warehouseGroup,
		DB:                 db,
	}
}

func (s *WarehouseGroupServiceImpl) Save(request masterwarehousepayloads.GetWarehouseGroupResponse) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	save, err := s.warehouseGroupRepo.Save(tx, request)

	if err != nil {
		return false
	}

	return save
}

func (s *WarehouseGroupServiceImpl) GetById(warehouseGroupId int) masterwarehousepayloads.GetWarehouseGroupResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseGroupRepo.GetById(tx, warehouseGroupId)

	if err != nil {
		return masterwarehousepayloads.GetWarehouseGroupResponse{}
	}

	return get
}

func (s *WarehouseGroupServiceImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseGroupRequest) []masterwarehousepayloads.GetWarehouseGroupResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	get, err := s.warehouseGroupRepo.GetAll(tx, request)

	if err != nil {
		return nil
	}

	return get
}

func (s *WarehouseGroupServiceImpl) ChangeStatus(warehouseGroupId int) masterwarehousepayloads.GetWarehouseGroupResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	change_status, err := s.warehouseGroupRepo.ChangeStatus(tx, warehouseGroupId)

	if err != nil {
		log.Panic(err.Error())
	}

	return change_status
}
