package masterserviceimpl

import (
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services"

	"after-sales/api/utils"

	"gorm.io/gorm"
)

type DiscountServiceImpl struct {
	discountRepo masterrepository.DiscountRepository
}

func StartDiscountService(discountRepo masterrepository.DiscountRepository) masterservice.DiscountService {
	return &DiscountServiceImpl{
		discountRepo: discountRepo,
	}
}
func (s *DiscountServiceImpl) WithTrx(trxHandle *gorm.DB) masterservice.DiscountService {
	s.discountRepo = s.discountRepo.WithTrx(trxHandle)
	return s
}

func (s *DiscountServiceImpl) GetAllDiscountIsActive() ([]masterpayloads.DiscountResponse, error) {
	results, err := s.discountRepo.GetAllDiscountIsActive()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) GetDiscountById(id int) (masterpayloads.DiscountResponse, error) {
	results, err := s.discountRepo.GetDiscountById(id)

	if err != nil {
		return masterpayloads.DiscountResponse{}, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) GetDiscountByCode(Code string) (masterpayloads.DiscountResponse, error) {
	results, err := s.discountRepo.GetDiscountByCode(Code)
	if err != nil {
		return masterpayloads.DiscountResponse{}, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) GetAllDiscount(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	results, err := s.discountRepo.GetAllDiscount(filterCondition, pages)
	if err != nil {
		return pages, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) ChangeStatusDiscount(Id int) (bool, error) {
	results, err := s.discountRepo.ChangeStatusDiscount(Id)
	if err != nil {
		return false, err
	}
	return results, nil
}

func (s *DiscountServiceImpl) SaveDiscount(req masterpayloads.DiscountResponse) (bool, error) {
	results, err := s.discountRepo.SaveDiscount(req)
	if err != nil {
		return false, err
	}
	return results, nil
}
