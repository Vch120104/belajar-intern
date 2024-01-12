package masteritemserviceimpl

import (
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type MarkupMasterServiceImpl struct {
	markupRepo masteritemrepository.MarkupMasterRepository
}

func StartMarkupMasterService(markupRepo masteritemrepository.MarkupMasterRepository) masteritemservice.MarkupMasterService {
	return &MarkupMasterServiceImpl{
		markupRepo: markupRepo,
	}
}

func (s *MarkupMasterServiceImpl) WithTrx(trxHandle *gorm.DB) masteritemservice.MarkupMasterService {
	s.markupRepo = s.markupRepo.WithTrx(trxHandle)
	return s
}
func (s *MarkupMasterServiceImpl) GetMarkupMasterList(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	results, err := s.markupRepo.GetMarkupMasterList(filter, pages)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (s *MarkupMasterServiceImpl) GetMarkupMasterById(id int) (masteritempayloads.MarkupMasterResponse, error) {
	results, err := s.markupRepo.GetMarkupMasterById(id)

	if err != nil {
		return masteritempayloads.MarkupMasterResponse{}, err
	}
	return results, nil
}

func (s *MarkupMasterServiceImpl) SaveMarkupMaster(req masteritempayloads.MarkupMasterResponse) (bool, error) {
	results, err := s.markupRepo.SaveMarkupMaster(req)
	if err != nil {
		return false, err
	}
	return results, nil
}
func (s *MarkupMasterServiceImpl) ChangeStatusMasterMarkupMaster(Id int) (bool, error) {
	results, err := s.markupRepo.ChangeStatusMasterMarkupMaster(Id)
	if err != nil {
		return false, err
	}
	return results, nil
}
func (s *MarkupMasterServiceImpl) GetMarkupMasterByCode(markupCode string) (masteritempayloads.MarkupMasterResponse, error) {
	result, err := s.markupRepo.GetMarkupMasterByCode(markupCode)
	if err != nil {
		return result, err
	}
	return result, nil

}
