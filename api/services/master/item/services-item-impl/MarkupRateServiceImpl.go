package masteritemserviceimpl

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupRateServiceImpl struct {
	markupRepo masteritemrepository.MarkupRateRepository
}

func StartMarkupRateService(markupRepo masteritemrepository.MarkupRateRepository) masteritemservice.MarkupRateService {
	return &MarkupRateServiceImpl{
		markupRepo: markupRepo,
	}
}

func (s *MarkupRateServiceImpl) WithTrx(trxHandle *gorm.DB) masteritemservice.MarkupRateService {
	s.markupRepo = s.markupRepo.WithTrx(trxHandle)
	return s
}

func (s *MarkupRateServiceImpl) GetAllMarkupRate(filterCondition []utils.FilterCondition) ([]map[string]interface{}, error) {
	results, err := s.markupRepo.GetAllMarkupRate(filterCondition)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *MarkupRateServiceImpl) GetMarkupRateById(id int) (masteritempayloads.MarkupRateResponse, error) {
	results, err := s.markupRepo.GetMarkupRateById(id)

	if err != nil {
		return masteritempayloads.MarkupRateResponse{}, err
	}
	return results, nil
}

func (s *MarkupRateServiceImpl) SaveMarkupRate(req masteritempayloads.MarkupRateRequest) (bool, error) {
	results, err := s.markupRepo.SaveMarkupRate(req)
	if err != nil {
		return results, err
	}
	return results, nil
}