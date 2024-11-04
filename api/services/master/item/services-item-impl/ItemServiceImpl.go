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

// CheckItemCodeExist implements masteritemservice.ItemService.
func (s *ItemServiceImpl) CheckItemCodeExist(itemCode string, itemGroupId int, commonPriceList bool, brandId int) (bool, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, itemId, itemClassId, repoErr := s.itemRepo.CheckItemCodeExist(tx, itemCode, itemGroupId, commonPriceList, brandId)
	defer helper.CommitOrRollback(tx, repoErr)

	if repoErr != nil {
		return results, itemId, itemClassId, repoErr
	}
	return results, itemId, itemClassId, nil
}

func (s *ItemServiceImpl) GetAllItemSearch(filterCondition []utils.FilterCondition, itemIDs []string, supplierIDs []string, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, repoErr := s.itemRepo.GetAllItemSearch(tx, filterCondition, itemIDs, supplierIDs, pages)
	defer helper.CommitOrRollback(tx, repoErr)

	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}
	return results, totalPages, totalRows, nil
}

// GetUomDropDown implements masteritemservice.ItemService.
func (s *ItemServiceImpl) GetUomDropDown(uomTypeId int) ([]masteritempayloads.UomDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemRepo.GetUomDropDown(tx, uomTypeId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetUomTypeDropDown implements masteritemservice.ItemService.
func (s *ItemServiceImpl) GetUomTypeDropDown() ([]masteritempayloads.UomTypeDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemRepo.GetUomTypeDropDown(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) GetAllItem(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	results, totalPages, totalRows, repoErr := s.itemRepo.GetAllItem(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, repoErr)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}
	return results, totalPages, totalRows, nil

}

func (s *ItemServiceImpl) GetAllItemListTransLookup(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	pages, err := s.itemRepo.GetAllItemListTransLookup(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return pages, err
	}
	return pages, nil
}

func (s *ItemServiceImpl) GetAllItemLookup(filter []utils.FilterCondition) (any, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.itemRepo.GetAllItemLookup(tx, filter)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) GetItemById(Id int) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetItemById(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) GetCatalogCode() ([]masteritempayloads.GetCatalogCode, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetCatalogCode(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) GetItemWithMultiId(MultiIds []string) ([]masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetItemWithMultiId(tx, MultiIds)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) GetItemCode(code string) (masteritempayloads.ItemResponse, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	results, err := s.itemRepo.GetItemCode(tx, code) // Menggunakan kode yang telah diencode
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) SaveItem(req masteritempayloads.ItemRequest) (masteritempayloads.ItemSaveResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result := masteritempayloads.ItemSaveResponse{}

	if req.ItemId != 0 {
		_, err := s.itemRepo.GetItemById(tx, req.ItemId)
		if err != nil {
			return result, err
		}
	}

	results, err := s.itemRepo.SaveItem(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return results, nil
}

func (s *ItemServiceImpl) ChangeStatusItem(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()

	_, err := s.itemRepo.GetItemById(tx, Id)
	if err != nil {
		return false, err
	}

	results, err := s.itemRepo.ChangeStatusItem(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *ItemServiceImpl) GetAllItemDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	results, totalPages, totalRows, repoErr := s.itemRepo.GetAllItemDetail(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, repoErr)
	if repoErr != nil {
		return results, totalPages, totalRows, repoErr
	}

	return results, totalPages, totalRows, nil

}

func (s *ItemServiceImpl) GetItemDetailById(itemID, itemDetailID int) (masteritempayloads.ItemDetailRequest, *exceptions.BaseErrorResponse) {

	tx := s.DB.Begin()
	result, err := s.itemRepo.GetItemDetailById(tx, itemID, itemDetailID)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) AddItemDetail(id int, req masteritempayloads.ItemDetailRequest) (masteritementities.ItemDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	entities, err := s.itemRepo.AddItemDetail(tx, id, req)
	defer helper.CommitOrRollback(tx, err)

	if err != nil {
		return masteritementities.ItemDetail{}, err
	}
	return entities, nil
}

func (s *ItemServiceImpl) DeleteItemDetails(id int, itemDetailIDs []int) (masteritempayloads.DeleteItemResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	resp, err := s.itemRepo.DeleteItemDetails(tx, id, itemDetailIDs)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return masteritempayloads.DeleteItemResponse{}, err
	}
	return resp, nil
}

func (s *ItemServiceImpl) UpdateItem(id int, req masteritempayloads.ItemUpdateRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.UpdateItem(tx, id, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) UpdateItemDetail(id int, itemDetail int, req masteritempayloads.ItemDetailUpdateRequest) (masteritementities.ItemDetail, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, errResp := s.itemRepo.UpdateItemDetail(tx, id, itemDetail, req)
	defer helper.CommitOrRollback(tx, errResp)
	if errResp != nil {
		return masteritementities.ItemDetail{}, errResp
	}
	return result, nil
}

func (s *ItemServiceImpl) GetPrincipleBrandParent(id int) ([]masteritempayloads.PrincipleBrandDropdownDescription, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetPrincipleBrandParent(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) GetPrincipleBrandDropdown() ([]masteritempayloads.PrincipleBrandDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.GetPrincipleBrandDropdown(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) AddItemDetailByBrand(id string, itemId int) ([]masteritempayloads.ItemDetailResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.itemRepo.AddItemDetailByBrand(tx, id, itemId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}
