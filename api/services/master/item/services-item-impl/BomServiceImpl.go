package masteritemserviceimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BomServiceImpl struct {
	BomRepository masteritemrepository.BomRepository
	DB            *gorm.DB
	RedisClient   *redis.Client // Redis client
}

func StartBomService(BomRepository masteritemrepository.BomRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.BomService {
	return &BomServiceImpl{
		BomRepository: BomRepository,
		DB:            db,
		RedisClient:   redisClient,
	}
}

func (s *BomServiceImpl) GetBomMasterList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, totalPages, totalRows, err := s.BomRepository.GetBomMasterList(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *BomServiceImpl) GetBomMasterById(id int) (masteritempayloads.BomMasterRequest, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.BomRepository.GetBomMasterById(tx, id)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return results, err
	}

	return results, nil
}

func (s *BomServiceImpl) SaveBomMaster(req masteritempayloads.BomMasterRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.BomRepository.SaveBomMaster(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *BomServiceImpl) UpdateBomMaster(id int, req masteritempayloads.BomMasterRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.BomRepository.UpdateBomMaster(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}

	return results, nil
}

func (s *BomServiceImpl) ChangeStatusBomMaster(Id int) (masteritementities.Bom, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	// Ubah status
	entity, err := s.BomRepository.ChangeStatusBomMaster(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritementities.Bom{}, err
	}

	return entity, nil
}

func (s *BomServiceImpl) GetBomDetailList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	//log.Printf("Menerima kondisi filter: %+v", filterCondition) // Tambahkan log untuk menerima kondisi filter
	results, totalPages, totalRows, err := s.BomRepository.GetBomDetailList(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}

	return results, totalPages, totalRows, nil
}

func (s *BomServiceImpl) GetBomDetailById(id int, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.BomRepository.GetBomDetailById(tx, id, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *BomServiceImpl) SaveBomDetail(req masteritempayloads.BomDetailRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.BomRepository.SaveBomDetail(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *BomServiceImpl) UpdateBomDetail(id int, req masteritempayloads.BomDetailRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	results, err := s.BomRepository.UpdateBomDetail(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *BomServiceImpl) GetBomItemList(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	//log.Printf("Menerima kondisi filter: %+v", filterCondition) // Tambahkan log untuk menerima kondisi filter
	results, totalPages, totalRows, err := s.BomRepository.GetBomItemList(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *BomServiceImpl) DeleteByIds(ids []int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	delete, err := s.BomRepository.DeleteByIds(tx, ids)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return false, err
	}
	return delete, nil
}
