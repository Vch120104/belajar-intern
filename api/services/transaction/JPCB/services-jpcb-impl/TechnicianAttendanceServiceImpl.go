package transactionjpcbserviceimpl

import (
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	"after-sales/api/payloads/pagination"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	transactionjpcbrepository "after-sales/api/repositories/transaction/JPCB"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TechnicianAttendanceImpl struct {
	TechnicianAttendanceRepository transactionjpcbrepository.TechnicianAttendanceRepository
	DB                             *gorm.DB
	RedisClient                    *redis.Client
}

func StartTechnicianAttendanceImpl(technicianAttendanceRepository transactionjpcbrepository.TechnicianAttendanceRepository, db *gorm.DB, redisClient *redis.Client) transactionjpcbservice.TechnicianAttendanceService {
	return &TechnicianAttendanceImpl{
		TechnicianAttendanceRepository: technicianAttendanceRepository,
		DB:                             db,
		RedisClient:                    redisClient,
	}
}

func (s *TechnicianAttendanceImpl) GetAllTechnicianAttendance(filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.TechnicianAttendanceRepository.GetAllTechnicianAttendance(tx, filterCondition, pages)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *TechnicianAttendanceImpl) SaveTechnicianAttendance(req transactionjpcbpayloads.TechnicianAttendanceSaveRequest) (transactionjpcbentities.TechnicianAttendance, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.TechnicianAttendanceRepository.SaveTechnicianAttendance(tx, req)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *TechnicianAttendanceImpl) ChangeStatusTechnicianAttendance(technicianAttendanceId int) (transactionjpcbentities.TechnicianAttendance, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	result, err := s.TechnicianAttendanceRepository.ChangeStatusTechnicianAttendance(tx, technicianAttendanceId)
	defer helper.CommitOrRollback(tx, err)
	if err != nil {
		return result, err
	}
	return result, nil
}
