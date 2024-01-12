package masteritemserviceimpl

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"

	"gorm.io/gorm"
)

type PriceListServiceImpl struct {
	priceListRepo masteritemrepository.PriceListRepository
}

func StartPriceListService(priceListRepo masteritemrepository.PriceListRepository) masteritemservice.PriceListService {
	return &PriceListServiceImpl{
		priceListRepo: priceListRepo,
	}
}

func (s *PriceListServiceImpl) WithTrx(trxHandle *gorm.DB) masteritemservice.PriceListService {
	s.priceListRepo = s.priceListRepo.WithTrx(trxHandle)
	return s
}

func (s *PriceListServiceImpl) GetPriceList(request masteritempayloads.PriceListGetAllRequest) ([]masteritempayloads.PriceListResponse, error) {
	results, err := s.priceListRepo.GetPriceList(request)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *PriceListServiceImpl) GetPriceListById(Id int) (masteritempayloads.PriceListResponse, error) {
	results, err := s.priceListRepo.GetPriceListById(Id)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *PriceListServiceImpl) SavePriceList(request masteritempayloads.PriceListResponse) (bool, error) {
	result, err := s.priceListRepo.SavePriceList(request)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *PriceListServiceImpl) ChangeStatusPriceList(Id int) (bool, error){
	result, err := s.priceListRepo.ChangeStatusPriceList(Id)
	if err != nil {
		return result, err
	}
	return result, nil
}