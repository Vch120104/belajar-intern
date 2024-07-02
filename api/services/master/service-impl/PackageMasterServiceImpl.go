package masterserviceimpl

import (
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
	if err != nil {
		return nil, 0, 0, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, totalPages, totalRows, nil
}

func (s *PackageMasterServiceImpl) GetAllPackageMasterDetail(pages pagination.Pagination, id int) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, totalPages, totalRows, err := s.PackageMasterRepo.GetAllPackageMasterDetail(tx, id, pages)
	if err != nil {
		return nil, 0, 0, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, totalPages, totalRows, nil
}

func (s *PackageMasterServiceImpl) GetByIdPackageMaster(id int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.GetByIdPackageMaster(tx, id)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *PackageMasterServiceImpl) GetByIdPackageMasterDetail(id int, idhead int, LineTypeId int) (map[string]interface{}, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.GetByIdPackageMasterDetail(tx, id, idhead, LineTypeId)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *PackageMasterServiceImpl) PostPackageMaster(req masterpayloads.PackageMasterResponse) (bool, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.PostpackageMaster(tx, req)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *PackageMasterServiceImpl) PostPackageMasterDetailBodyshop(req masterpayloads.PackageMasterDetailOperationBodyshop, id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.PostPackageMasterDetailBodyshop(tx, req, id)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *PackageMasterServiceImpl) PostPackageMasterDetailWorkshop(req masterpayloads.PackageMasterDetailWorkshop) (bool, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.PostPackageMasterDetailWorkshop(tx, req)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *PackageMasterServiceImpl) ChangeStatusItemPackage(id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.ChangeStatusItemPackage(tx, id)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *PackageMasterServiceImpl) ActivateMultiIdPackageMasterDetail(ids string, idhead int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.ActivateMultiIdPackageMasterDetail(tx, ids, idhead)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *PackageMasterServiceImpl) DeactivateMultiIdPackageMasterDetail(ids string, idhead int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.DeactivateMultiIdPackageMasterDetail(tx, ids, idhead)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *PackageMasterServiceImpl) CopyToOtherModel(id int, name string, idmodel int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.PackageMasterRepo.CopyToOtherModel(tx, id, name, idmodel)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}
