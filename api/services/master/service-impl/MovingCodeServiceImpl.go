package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	masterservice "after-sales/api/services/master"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MovingCodeServiceImpl struct {
	MovingCodeRepository masterrepository.MovingCodeRepository
	DB                   *gorm.DB
	RedisClient          *redis.Client // Redis client
}

func StartMovingCodeService(MovingCodeRepository masterrepository.MovingCodeRepository, db *gorm.DB, redisClient *redis.Client) masterservice.MovingCodeService {
	return &MovingCodeServiceImpl{
		MovingCodeRepository: MovingCodeRepository,
		DB:                   db,
		RedisClient:          redisClient,
	}
}

func (s *MovingCodeServiceImpl) GetAllMovingCode(pages pagination.Pagination) pagination.Pagination {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)
	results, err := s.MovingCodeRepository.GetAllMovingCode(tx, pages)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}
	return results
}

func (s *MovingCodeServiceImpl) SaveMovingCode(req masterpayloads.MovingCodeRequest) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if req.MovingCodeId != 0 {
		_, err := s.MovingCodeRepository.GetMovingCodeById(tx, req.MovingCodeId)

		if err != nil {
			panic(exceptions.NewNotFoundError(err.Error()))
		}
	}

	results, err := s.MovingCodeRepository.SaveMovingCode(tx, req)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	return results
}

func (s *MovingCodeServiceImpl) ChangePriorityMovingCode(id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// Check if the moving code with the given ID exists
	decreasedMovingCode, err := s.MovingCodeRepository.GetMovingCodeById(tx, id)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	priority := decreasedMovingCode.Priority - 1
	// Find the moving code with the increased priority
	increasedMovingCode, err := s.MovingCodeRepository.GetMovingCodeByPriority(tx, priority)
	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	// Decrease the priority of the specified moving code
	movingCodeDecreased, err := s.MovingCodeRepository.DecreasePriorityMovingCode(tx, id)
	if err != nil {
		return movingCodeDecreased
	}

	// Increase the priority of the moving code with increased priority
	_, err = s.MovingCodeRepository.IncreasePriorityMovingCode(tx, increasedMovingCode.MovingCodeId)
	if err != nil {
		panic(exceptions.NewAppExceptionError(err.Error()))
	}

	return true
}

func (s *MovingCodeServiceImpl) ChangeStatusMovingCode(id int) bool {
	tx := s.DB.Begin()
	defer helper.CommitOrRollback(tx)

	_, err := s.MovingCodeRepository.GetMovingCodeById(tx, id)

	if err != nil {
		panic(exceptions.NewNotFoundError(err.Error()))
	}

	results, err := s.MovingCodeRepository.ChangeStatusMovingCode(tx, id)
	if err != nil {
		return results
	}
	return true
}
