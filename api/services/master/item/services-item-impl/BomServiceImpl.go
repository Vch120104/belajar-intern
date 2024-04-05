package masteritemserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type BomServiceImpl struct {
	BomRepository masteritemrepository.BomRepository
	DB            *gorm.DB
}

func StartBomService(BomRepository masteritemrepository.BomRepository, db *gorm.DB) masteritemservice.BomService {
	return &BomServiceImpl{
		BomRepository: BomRepository,
		DB:            db,
	}
}

func (s *BomServiceImpl) GetBomMasterList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	//log.Printf("Menerima kondisi filter: %+v", filterCondition) // Tambahkan log untuk menerima kondisi filter
	results, totalPages, totalRows, err := s.BomRepository.GetBomMasterList(tx, filterCondition, pages)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *BomServiceImpl) GetBomMasterById(id int) (masteritempayloads.BomMasterRequest, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.BomRepository.GetBomMasterById(tx, id)

	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *BomServiceImpl) SaveBomMaster(req masteritempayloads.BomMasterRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.BomRepository.SaveBomMaster(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *BomServiceImpl) ChangeStatusBomMaster(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.BomRepository.GetBomMasterById(tx, Id)

	if err != nil {
		return false, err
	}

	results, err := s.BomRepository.ChangeStatusBomMaster(tx, Id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *BomServiceImpl) GetBomDetailList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	//log.Printf("Menerima kondisi filter: %+v", filterCondition) // Tambahkan log untuk menerima kondisi filter
	results, totalPages, totalRows, err := s.BomRepository.GetBomDetailList(tx, filterCondition, pages)
	if err != nil {
		return results,0,0,err
	}
	return results, totalPages, totalRows,nil
}

func (s *BomServiceImpl) GetBomDetailById(id int) ([]masteritempayloads.BomDetailListResponse,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.BomRepository.GetBomDetailById(tx, id)

	if err != nil {
		return results,err
	}
	return results,nil
}

func (s *BomServiceImpl) SaveBomDetail(req masteritempayloads.BomDetailRequest) (bool,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.BomRepository.SaveBomDetail(tx, req)
	if err != nil {
		return false,err
	}
	return results,nil
}

func (s *BomServiceImpl) GetBomItemList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int,*exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	//log.Printf("Menerima kondisi filter: %+v", filterCondition) // Tambahkan log untuk menerima kondisi filter
	results, totalPages, totalRows, err := s.BomRepository.GetBomItemList(tx, filterCondition, pages)
	if err != nil {
		return results,0,0,err
	}
	return results, totalPages, totalRows,nil
}
