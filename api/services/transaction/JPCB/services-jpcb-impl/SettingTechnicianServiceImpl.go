package transactionjpcbserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SettingTechnicianServiceImpl struct {
	SettingTechnicianRepository transactionjpcbrepository.SettingTechnicianRepository
	DB                          *gorm.DB
	RedisClient                 *redis.Client
}

func StartServiceTechnicianService(SettingTechnicianRepo transactionjpcbrepository.SettingTechnicianRepository, db *gorm.DB, redisClient *redis.Client) transactionjpcbservice.SettingTechnicianService {
	return &SettingTechnicianServiceImpl{
		SettingTechnicianRepository: SettingTechnicianRepo,
		DB:                          db,
		RedisClient:                 redisClient,
	}
}

func (s *SettingTechnicianServiceImpl) GetAllSettingTechnician(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	pages, err := s.SettingTechnicianRepository.GetAllSettingTechnician(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return pages, err
	}
	return pages, nil
}

func (s *SettingTechnicianServiceImpl) GetAllSettingTechnicianDetail(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	pages, err := s.SettingTechnicianRepository.GetAllSettingTechnicianDetail(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return pages, err
	}
	return pages, nil
}

func (s *SettingTechnicianServiceImpl) GetSettingTechnicianById(settingTechnicianId int) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.SettingTechnicianRepository.GetSettingTechnicianById(tx, settingTechnicianId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *SettingTechnicianServiceImpl) GetSettingTechnicianDetailById(settingTechnicianDetailId int) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.SettingTechnicianRepository.GetSettingTechnicianDetailById(tx, settingTechnicianDetailId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *SettingTechnicianServiceImpl) GetSettingTechnicianByCompanyDate(companyId int, effectiveDate time.Time) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.SettingTechnicianRepository.GetSettingTechnicianByCompanyDate(tx, companyId, effectiveDate)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *SettingTechnicianServiceImpl) SaveSettingTechnician(CompanyId int) (transactionjpcbpayloads.SettingTechnicianGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.SettingTechnicianRepository.SaveSettingTechnician(tx, CompanyId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *SettingTechnicianServiceImpl) SaveSettingTechnicianDetail(req transactionjpcbpayloads.SettingTechnicianDetailSaveRequest) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse
	defer helper.CommitOrRollback(tx, err)

	if req.SettingTechnicianSystemNumber == 0 {
		headerData, err := s.SettingTechnicianRepository.SaveSettingTechnician(tx, req.CompanyId)
		if err != nil {
			return transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse{}, err
		}
		req.SettingTechnicianSystemNumber = headerData.SettingTechnicianSystemNumber
	}

	result, err := s.SettingTechnicianRepository.SaveSettingTechnicianDetail(tx, req)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *SettingTechnicianServiceImpl) UpdateSettingTechnicianDetail(settingTechnicianDetailId int, req transactionjpcbpayloads.SettingTechnicianDetailUpdateRequest) (transactionjpcbpayloads.SettingTechnicianDetailGetByIdResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.SettingTechnicianRepository.UpdateSettingTechnicianDetail(tx, settingTechnicianDetailId, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}
