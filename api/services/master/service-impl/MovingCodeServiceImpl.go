package masterserviceimpl

import (
	exceptions "after-sales/api/exceptions"
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

// ActivateMovingCode implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) ActivateMovingCode(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.MovingCodeRepo.ActivateMovingCode(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// DeactiveMovingCode implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) DeactiveMovingCode(id string) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.MovingCodeRepo.DeactiveMovingCode(tx, id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetDropdownMovingCode implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) GetDropdownMovingCode() ([]masterpayloads.MovingCodeDropDown, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.MovingCodeRepo.GetDropdownMovingCode(tx)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// ChangeStatusMovingCode implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) ChangeStatusMovingCode(Id int) (any, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.MovingCodeRepo.ChangeStatusMovingCode(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// CreateMovingCode implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) CreateMovingCode(req masterpayloads.MovingCodeListRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.MovingCodeRepo.CreateMovingCode(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// GetAllMovingCode implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) GetAllMovingCode(pages pagination.Pagination) ([]map[string]any, int, int, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, totalPages, totalRows, err := s.MovingCodeRepo.GetAllMovingCode(tx, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, totalPages, totalRows, err
	}
	return results, totalPages, totalRows, nil
}

// GetMovingCodebyId implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) GetMovingCodebyId(Id int) (masterpayloads.MovingCodeResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.MovingCodeRepo.GetMovingCodebyId(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// PushMovingCodePriority implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) PushMovingCodePriority(Id int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.MovingCodeRepo.PushMovingCodePriority(tx, Id)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return results, err
	}
	return results, nil
}

// UpdateMovingCode implements masterservice.MovingCodeService.
func (s *MovingCodeServiceImpl) UpdateMovingCode(req masterpayloads.MovingCodeListUpdate) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	results, err := s.MovingCodeRepo.UpdateMovingCode(tx, req)
	defer helper.CommitOrRollback(tx, err)
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
