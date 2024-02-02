package masterserviceimpl

import (
	// "after-sales/api/exceptions"
	"after-sales/api/helper"
	// masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"
	"after-sales/api/utils"

	"gorm.io/gorm"
)

type FieldActionServiceImpl struct {
	FieldActionRepo masterrepository.FieldActionRepository
	DB              *gorm.DB
}

func StartFieldActionService(FieldActionRepo masterrepository.FieldActionRepository, db *gorm.DB) masterservice.FieldActionService {
	return &FieldActionServiceImpl{
		FieldActionRepo: FieldActionRepo,
		DB:              db,
	}
}

func (s *FieldActionServiceImpl) GetAllFieldAction(filterCondition []utils.FilterCondition, pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.FieldActionRepo.GetAllFieldAction(tx, filterCondition, pages)
	if err != nil {
		return pages
	}
	return results
}
