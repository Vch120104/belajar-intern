package masteritemserviceimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	masteritemservice "after-sales/api/services/master/item"
	"after-sales/api/utils"

	"gorm.io/gen/helper"
	"gorm.io/gorm"
)

type ItemSubstituteServiceImpl struct {
	itemSubstituteRepo masteritemrepository.ItemSubstituteRepository
	Db                 *gorm.DB
}

func StartItemSubstitute(itemSubstituteRepo masteritemrepository.ItemSubstituteRepository, db *gorm.DB) masteritemservice.ItemSubstituteService {
	return &ItemSubstituteServiceImpl{
		itemSubstituteRepo: itemSubstituteRepo,
		Db:                 db,
	}
}

func (s *ItemSubstituteServiceImpl) GetAllItemSubstitute(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, error) {
	tx := s.Db.Begin()

	defer helper.CommitOrRollback(tx)

	results, err := s.itemSubstituteRepo.GetAllItemSubstitute(tx, filterCondition, pages)
	if err != nil {
		panic(exceptions.newApp)
	}

}
