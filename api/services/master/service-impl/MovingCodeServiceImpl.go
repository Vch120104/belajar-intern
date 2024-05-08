package masterserviceimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

	"gorm.io/gorm"
)

type MovingCodeServiceImpl struct {
	MovingCodeRepo masterrepository.MovingCodeRepository
	DB             *gorm.DB
}

// ChangeStatusMovingCode implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) ChangeStatusMovingCode(Id int) (any, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.MovingCodeRepo.ChangeStatusMovingCode(tx, Id)
	if err != nil {
		return results, err
	}
	return results, nil
}

// CreateMovingCode implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) CreateMovingCode(req masterpayloads.MovingCodeListRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.MovingCodeRepo.CreateMovingCode(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetAllMovingCode implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) GetAllMovingCode(pages pagination.Pagination) ([]map[string]any, int, int, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, totalPages, totalRows, err := s.MovingCodeRepo.GetAllMovingCode(tx, pages)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

// GetMovingCodebyId implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) GetMovingCodebyId(Id int) (any, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.MovingCodeRepo.GetMovingCodebyId(tx, Id)
	if err != nil {
		return results, err
	}
	return results, nil
}

// PushMovingCodePriority implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) PushMovingCodePriority(Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.MovingCodeRepo.PushMovingCodePriority(tx, Id)
	if err != nil {
		return results, err
	}
	return results, nil
}

// UpdateMovingCode implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) UpdateMovingCode(req masterpayloads.MovingCodeListRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.MovingCodeRepo.UpdateMovingCode(tx, req)
	if err != nil {
		return results, err
	}
	return results, nil
}

func StartMovingCodeServiceImpl(MovingCodeRepo masterrepository.MovingCodeRepository, db *gorm.DB) masterservice.MovingCodeService {
	return &MovingCodeServiceImpl{
		MovingCodeRepo: MovingCodeRepo,
		DB:             db,
	}
}
