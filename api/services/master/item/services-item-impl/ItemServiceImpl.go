package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	"after-sales/api/helper"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"fmt"

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

func (s *ItemServiceImpl) GetAllItemSearch(filterCondition []utils.FilterCondition, itemIDs []string, supplierIDs []string, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, repoErr := s.itemRepo.GetAllItemSearch(tx, filterCondition, itemIDs, supplierIDs, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}
	return results, totalPages, totalRows, nil
}

// GetUomDropDown implements masteritemservice.ItemService.
func (s *ItemServiceImpl) GetUomDropDown(uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemRepo.GetUomDropDown(tx, uomTypeId)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

// GetUomTypeDropDown implements masteritemservice.ItemService.
func (s *ItemServiceImpl) GetUomTypeDropDown() ([]masteritempayloads.UomTypeDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemRepo.GetUomTypeDropDown(tx)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemServiceImpl) GetAllItem(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	results, totalPages, totalRows, repoErr := s.itemRepo.GetAllItem(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}
	defer helper.CommitOrRollback(tx, repoErr)
	return results, totalPages, totalRows, nil

}

func (s *ItemServiceImpl) GetAllItemLookup(filter []utils.FilterCondition) (any, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemRepo.GetAllItemLookup(tx, filter)
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemServiceImpl) GetItemById(Id int) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetItemById(tx, Id)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemServiceImpl) GetItemWithMultiId(MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetItemWithMultiId(tx, MultiIds)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemServiceImpl) GetItemCode(code string) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	results, err := s.itemRepo.GetItemCode(tx, code) // Menggunakan kode yang telah diencode
	if err != nil {
		return results, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemServiceImpl) SaveItem(req masteritempayloads.ItemRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	fmt.Print("sini?")

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
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemServiceImpl) ChangeStatusItem(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.itemRepo.GetItemById(tx, Id)
	if err != nil {
		return false, err
	}

	results, err := s.itemRepo.ChangeStatusItem(tx, Id)
	if err != nil {
		return false, err
	}
	defer helper.CommitOrRollback(tx, err)
	return results, nil
}

func (s *ItemServiceImpl) GetAllItemDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	results, totalPages, totalRows, repoErr := s.itemRepo.GetAllItemDetail(tx, filterCondition, pages)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}
	defer helper.CommitOrRollback(tx, repoErr)

	return results, totalPages, totalRows, nil

}

func (s *ItemServiceImpl) GetItemDetailById(itemID, itemDetailID int) (masteritempayloads.ItemDetailRequest, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	result, err := s.itemRepo.GetItemDetailById(tx, itemID, itemDetailID)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemServiceImpl) AddItemDetail(id int, req masteritempayloads.ItemDetailRequest) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	err := s.itemRepo.AddItemDetail(tx, id, req)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx, err)
	return nil
}

func (s *ItemServiceImpl) DeleteItemDetail(id int, itemDetailID int) *exceptions.BaseErrorResponse {
	tx := s.DB.Begin()
	err := s.itemRepo.DeleteItemDetail(tx, id, itemDetailID)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx, err)
	return nil
}

func (s *ItemServiceImpl) UpdateItem(id int, req masteritempayloads.ItemUpdateRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.UpdateItem(tx, id, req)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemServiceImpl) UpdateItemDetail(id int, req masteritempayloads.ItemDetailUpdateRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.UpdateItemDetail(tx, id, req)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemServiceImpl) GetPrincipleBrandParent(code string) ([]masteritempayloads.PrincipleBrandDropdownDescription, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetPrincipleBrandParent(tx, code)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}

func (s *ItemServiceImpl) GetPrincipleBrandDropdown() ([]masteritempayloads.PrincipleBrandDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetPrincipleBrandDropdown(tx)
	if err != nil {
		return result, err
	}
	defer helper.CommitOrRollback(tx, err)
	return result, nil
}
