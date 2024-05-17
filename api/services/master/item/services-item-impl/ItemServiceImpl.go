package masteritemserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ItemServiceImpl struct {
	itemRepo    masteritemrepository.ItemRepository
	DB          *gorm.DB
	RedisClient *redis.Client // Redis client
}

func StartItemService(itemRepo masteritemrepository.ItemRepository, db *gorm.DB, redisClient *redis.Client) masteritemservice.ItemService {
	return &ItemServiceImpl{
		itemRepo:    itemRepo,
		DB:          db,
		RedisClient: redisClient,
	}
}

// GetUomDropDown implements masteritemservice.ItemService.
func (s *ItemServiceImpl) GetUomDropDown(uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemRepo.GetUomDropDown(tx, uomTypeId)
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetUomTypeDropDown implements masteritemservice.ItemService.
func (s *ItemServiceImpl) GetUomTypeDropDown() ([]masteritempayloads.UomTypeDropdownResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemRepo.GetUomTypeDropDown(tx)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) GetAllItem(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemRepo.GetAllItem(tx, filterCondition, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) GetAllItemLookup(internalFilterCondition []utils.FilterCondition, externalFilterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.itemRepo.GetAllItemLookup(tx, internalFilterCondition, externalFilterCondition, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

func (s *ItemServiceImpl) GetItemById(Id int) (map[string]any, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemRepo.GetItemById(tx, Id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) GetItemWithMultiId(MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemRepo.GetItemWithMultiId(tx, MultiIds)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) GetItemCode(code string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {
	// Melakukan URL encoding pada parameter code
	// encodedCode := url.PathEscape(code)

	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.itemRepo.GetItemCode(tx, code) // Menggunakan kode yang telah diencode
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) SaveItem(req masteritempayloads.ItemRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.ItemId != 0 {
		_, err := s.itemRepo.GetItemById(tx, req.ItemId)
		if err != nil {
			return false, err

		}
	}

	results, err := s.itemRepo.SaveItem(tx, req)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ItemServiceImpl) ChangeStatusItem(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.itemRepo.GetItemById(tx, Id)
	if err != nil {
		return false, err
	}

	results, err := s.itemRepo.ChangeStatusItem(tx, Id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ItemServiceImpl) GetAllItemDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.itemRepo.GetAllItemDetail(tx, filterCondition, pages)
	if err != nil {
		return results, 0, 0, err
	}
	return results, totalPages, totalRows, nil
}

func (s *ItemServiceImpl) GetItemDetailById(itemID, itemDetailID int) (masteritempayloads.ItemDetailRequest, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	result, err := s.itemRepo.GetItemDetailById(tx, itemID, itemDetailID)
	if err != nil {
		return masteritempayloads.ItemDetailRequest{}, err
	}
	return result, nil
}

func (s *ItemServiceImpl) AddItemDetail(id int, req masteritempayloads.ItemDetailRequest) *exceptionsss_test.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.itemRepo.AddItemDetail(tx, id, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *ItemServiceImpl) DeleteItemDetail(id int, itemDetailID int) *exceptionsss_test.BaseErrorResponse {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	err := s.itemRepo.DeleteItemDetail(tx, id, itemDetailID)
	if err != nil {
		return err
	}
	return nil
}
