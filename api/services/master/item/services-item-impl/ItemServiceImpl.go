package masteritemserviceimpl

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemServiceImpl struct {
	itemRepo masteritemrepository.ItemRepository
}

func StartItemService(itemRepo masteritemrepository.ItemRepository) masteritemservice.ItemService {
	return &ItemServiceImpl{
		itemRepo: itemRepo,
	}
}

func (s *ItemServiceImpl) WithTrx(trxHandle *gorm.DB) masteritemservice.ItemService {
	s.itemRepo = s.itemRepo.WithTrx(trxHandle)
	return s
}

func (s *ItemServiceImpl) GetAllItem(filterCondition []utils.FilterCondition) ([]masteritempayloads.ItemLookup, error) {
	results, err := s.itemRepo.GetAllItem(filterCondition)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) GetAllItemLookup(queryParams map[string]string) ([]map[string]interface{}, error) {
	results, err := s.itemRepo.GetAllItemLookup(queryParams)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) GetItemById(Id int) (masteritempayloads.ItemResponse, error) {
	result, err := s.itemRepo.GetItemById(Id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) GetItemWithMultiId(MultiIds []string) ([]masteritempayloads.ItemResponse, error) {
	result, err := s.itemRepo.GetItemWithMultiId(MultiIds)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemServiceImpl) GetItemCode(code string) ([]map[string]interface{}, error) {
	results, err := s.itemRepo.GetItemCode(code)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) SaveItem(req masteritempayloads.ItemResponse) (bool, error) {
	results, err := s.itemRepo.SaveItem(req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemServiceImpl) ChangeStatusItem(Id int) (bool, error) {
	results, err := s.itemRepo.ChangeStatusItem(Id)
	if err != nil {
		return results, err
	}
	return results, nil
}
