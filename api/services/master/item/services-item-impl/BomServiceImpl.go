package masteritemserviceimpl

import (
	"after-sales/api/exceptions"
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

func (s *BomServiceImpl) GetBomMasterList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	//log.Printf("Menerima kondisi filter: %+v", filterCondition) // Tambahkan log untuk menerima kondisi filter
	results, totalPages, totalRows, err := s.BomRepository.GetBomMasterList(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results, totalPages, totalRows
}

func (s *BomServiceImpl) GetBomMasterById(id int) masteritempayloads.BomMasterRequest {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.BomRepository.GetBomMasterById(tx, id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *BomServiceImpl) SaveBomMaster(req masteritempayloads.BomMasterRequest) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.BomRepository.SaveBomMaster(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *BomServiceImpl) ChangeStatusBomMaster(Id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.BomRepository.GetBomMasterById(tx, Id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.BomRepository.ChangeStatusBomMaster(tx, Id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *BomServiceImpl) GetBomDetailList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	//log.Printf("Menerima kondisi filter: %+v", filterCondition) // Tambahkan log untuk menerima kondisi filter
	results, totalPages, totalRows, err := s.BomRepository.GetBomDetailList(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results, totalPages, totalRows
}

func (s *BomServiceImpl) GetBomDetailById(id int) []masteritempayloads.BomDetailRequest {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.BomRepository.GetBomDetailById(tx, id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}
