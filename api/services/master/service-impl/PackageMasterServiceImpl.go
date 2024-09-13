package masterserviceimpl

import (
	masterentities "after-sales/api/entities/master"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type PackageMasterServiceImpl struct {
	PackageMasterRepo masterrepository.PackageMasterRepository
	db                *gorm.DB
}

func StartPackageMasterService(PackageMasterRepo masterrepository.PackageMasterRepository, db *gorm.DB) masterservice.PackageMasterService {
	return &PackageMasterServiceImpl{
		PackageMasterRepo: PackageMasterRepo,
		db:                db,
	}
}

func (s *PackageMasterServiceImpl) GetAllPackageMaster(filtercondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, totalPages, totalRows, err := s.PackageMasterRepo.GetAllPackageMaster(tx, filtercondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, 0, 0, err
	}
	return result, totalPages, totalRows, nil
}

func (s *PackageMasterServiceImpl) GetAllPackageMasterDetail(pages pagination.Pagination, id int) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, totalPages, totalRows, err := s.PackageMasterRepo.GetAllPackageMasterDetail(tx, id, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return nil, 0, 0, err
	}
	return result, totalPages, totalRows, nil
}

func (s *PackageMasterServiceImpl) GetByIdPackageMaster(id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.GetByIdPackageMaster(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *PackageMasterServiceImpl) GetByIdPackageMasterDetail(id int, idhead int, LineTypeId int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.GetByIdPackageMasterDetail(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *PackageMasterServiceImpl) PostPackageMaster(req masterpayloads.PackageMasterResponse) (masterentities.PackageMaster, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.PostpackageMaster(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.PackageMaster{}, err
	}
	return result, nil
}

func (s *PackageMasterServiceImpl) PostPackageMasterDetail(req masterpayloads.PackageMasterDetail) (masterentities.PackageMasterDetail, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.PostPackageMasterDetail(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.PackageMasterDetail{}, err
	}
	return result, nil
}

func (s *PackageMasterServiceImpl) ChangeStatusItemPackage(id int) (masterentities.PackageMaster, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.ChangeStatusItemPackage(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masterentities.PackageMaster{}, err
	}
	return result, nil
}

func (s *PackageMasterServiceImpl) ActivateMultiIdPackageMasterDetail(ids string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.ActivateMultiIdPackageMasterDetail(tx, ids)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *PackageMasterServiceImpl) DeactivateMultiIdPackageMasterDetail(ids string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.DeactivateMultiIdPackageMasterDetail(tx, ids)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *PackageMasterServiceImpl) CopyToOtherModel(id int, name string, idmodel int) (int, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.CopyToOtherModel(tx, id, name, idmodel)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return 0, err
	}
	return result, nil
}
