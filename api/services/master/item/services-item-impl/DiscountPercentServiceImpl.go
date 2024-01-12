package masteritemserviceimpl

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DiscountPercentServiceImpl struct {
	discountPercentRepo masteritemrepository.DiscountPercentRepository
}

func StartDiscountPercentService(discountPercentRepo masteritemrepository.DiscountPercentRepository) masteritemservice.DiscountPercentService {
	return &DiscountPercentServiceImpl{
		discountPercentRepo: discountPercentRepo,
	}
}

func (s *DiscountPercentServiceImpl) WithTrx(trxHandle *gorm.DB) masteritemservice.DiscountPercentService {
	s.discountPercentRepo = s.discountPercentRepo.WithTrx(trxHandle)
	return s
}

func (s *DiscountPercentServiceImpl) GetAllDiscountPercent(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error) {
	results, err := s.discountPercentRepo.GetAllDiscountPercent(filterCondition)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountPercentServiceImpl) GetDiscountPercentById(Id int) (masteritempayloads.DiscountPercentResponse, error) {
	result, err := s.discountPercentRepo.GetDiscountPercentById(Id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *DiscountPercentServiceImpl) SaveDiscountPercent(req masteritempayloads.DiscountPercentResponse) (bool, error) {
	results, err := s.discountPercentRepo.SaveDiscountPercent(req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountPercentServiceImpl) ChangeStatusDiscountPercent(Id int) (bool, error) {
	results, err := s.discountPercentRepo.ChangeStatusDiscountPercent(Id)
	if err != nil {
		return results, err
	}
	return results, nil
}