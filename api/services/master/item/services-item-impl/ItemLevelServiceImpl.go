package masteritemserviceimpl

import (
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	masteritemlevelservice "after-sales/api/services/master/item"
	"log"

	"gorm.io/gorm"
)

type ItemLevelServiceImpl struct {
	structItemLevelRepo masteritemlevelrepo.ItemLevelRepository
}

func StartItemLevelService(itemlevelrepo masteritemlevelrepo.ItemLevelRepository) masteritemlevelservice.ItemLevelService {
	return &ItemLevelServiceImpl{
		structItemLevelRepo: itemlevelrepo,
	}
}

func (s *ItemLevelServiceImpl) WithTrx(Trxhandle *gorm.DB) masteritemlevelservice.ItemLevelService {
	s.structItemLevelRepo = s.structItemLevelRepo.WithTrx(Trxhandle)
	return s
}

func (s *ItemLevelServiceImpl) Save(request masteritemlevelpayloads.SaveItemLevelRequest) (bool, error) {
	save, err := s.structItemLevelRepo.Save(request)

	if err != nil {
		return false, err
	}

	return save, nil
}

func (s *ItemLevelServiceImpl) Update(request masteritemlevelpayloads.SaveItemLevelRequest) (bool, error) {
	update, err := s.structItemLevelRepo.Update(request)

	if err != nil {
		return false, err
	}

	return update, nil
}

func (s *ItemLevelServiceImpl) GetById(itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponse, error) {
	get, err := s.structItemLevelRepo.GetById(itemLevelId)

	if err != nil {
		return masteritemlevelpayloads.GetItemLevelResponse{}, err
	}

	return get, nil
}

func (s *ItemLevelServiceImpl) GetAll(request masteritemlevelpayloads.GetAllItemLevelResponse, pages pagination.Pagination) (pagination.Pagination, error) {
	get, err := s.structItemLevelRepo.GetAll(request, pages)

	if err != nil {
		return pagination.Pagination{}, err
	}

	return get, nil
}

func (s *ItemLevelServiceImpl) ChangeStatus(itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponse, error) {
	change_status, err := s.structItemLevelRepo.ChangeStatus(itemLevelId)

	if err != nil {
		log.Panic(err.Error())
	}

	return change_status, nil
}
