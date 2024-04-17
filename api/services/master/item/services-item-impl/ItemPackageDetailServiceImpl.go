package masteritemserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"

	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"

	"gorm.io/gorm"
)

type ItemPackageDetailServiceImpl struct {
	ItemPackageDetailRepo masteritemrepository.ItemPackageDetailRepository
	DB                    *gorm.DB
}

func StartItemPackageDetailService(ItemPackageDetailRepo masteritemrepository.ItemPackageDetailRepository, db *gorm.DB) masteritemservice.ItemPackageDetailService {
	return &ItemPackageDetailServiceImpl{
		ItemPackageDetailRepo: ItemPackageDetailRepo,
		DB:                    db,
	}
}

func (s *ItemPackageDetailServiceImpl) ChangeStatusItemPackageDetail(id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.ItemPackageDetailRepo.GetItemPackageDetailById(tx, id)

	if err != nil {
		return false, err
	}

	results, err := s.ItemPackageDetailRepo.ChangeStatusItemPackageDetail(tx, id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ItemPackageDetailServiceImpl) GetItemPackageDetailByItemPackageId(itemPackageId int, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageDetailRepo.GetItemPackageDetailByItemPackageId(tx, itemPackageId, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageDetailServiceImpl) GetItemPackageDetailById(itemPackageDetailId int) (masteritempayloads.ItemPackageDetailResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageDetailRepo.GetItemPackageDetailById(tx, itemPackageDetailId)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageDetailServiceImpl) CreateItemPackageDetailByItemPackageId(req masteritempayloads.SaveItemPackageDetail) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageDetailRepo.CreateItemPackageDetailByItemPackageId(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemPackageDetailServiceImpl) UpdateItemPackageDetailByItemPackageId(req masteritempayloads.SaveItemPackageDetail) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.ItemPackageDetailRepo.UpdateItemPackageDetailByItemPackageId(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}
