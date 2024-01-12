package masteritemserviceimpl

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type ItemClassServiceImpl struct {
	itemRepo masteritemrepository.ItemClassRepository
}

func StartItemClassService(itemRepo masteritemrepository.ItemClassRepository) masteritemservice.ItemClassService {
	return &ItemClassServiceImpl{
		itemRepo: itemRepo,
	}
}

func (s *ItemClassServiceImpl) WithTrx(trxHandle *gorm.DB) masteritemservice.ItemClassService {
	s.itemRepo = s.itemRepo.WithTrx(trxHandle)
	return s
}

func (s *ItemClassServiceImpl) GetAllItemClass(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error) {
	results, err := s.itemRepo.GetAllItemClass(filterCondition)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemClassServiceImpl) GetItemClassById(Id int) (masteritempayloads.ItemClassResponse, error) {
	result, err := s.itemRepo.GetItemClassById(Id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *ItemClassServiceImpl) SaveItemClass(req masteritempayloads.ItemClassResponse) (bool, error) {
	results, err := s.itemRepo.SaveItemClass(req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *ItemClassServiceImpl) ChangeStatusItemClass(Id int) (bool, error) {
	results, err := s.itemRepo.ChangeStatusItemClass(Id)
	if err != nil {
		return results, err
	}
	return results, nil
}
