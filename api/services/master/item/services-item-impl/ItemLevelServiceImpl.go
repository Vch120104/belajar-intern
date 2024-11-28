package masteritemserviceimpl

import (
	exceptions "after-sales/api/exceptions"
	masteritemlevelpayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemlevelrepo "after-sales/api/repositories/master/item"
	masteritemlevelservice "after-sales/api/services/master/item"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ItemLevelServiceImpl struct {
	structItemLevelRepo masteritemlevelrepo.ItemLevelRepository
	DB                  *gorm.DB
	RedisClient         *redis.Client // Redis client
}

func StartItemLevelService(itemlevelrepo masteritemlevelrepo.ItemLevelRepository, db *gorm.DB, redisClient *redis.Client) masteritemlevelservice.ItemLevelService {
	return &ItemLevelServiceImpl{
		structItemLevelRepo: itemlevelrepo,
		DB:                  db,
		RedisClient:         redisClient,
	}
}

// GetItemLevelLookUpbyId implements masteritemservice.ItemLevelService.
func (s *ItemLevelServiceImpl) GetItemLevelLookUpbyId(filter []utils.FilterCondition, itemLevelId int) (masteritemlevelpayloads.GetItemLevelLookUp, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	get, err := s.structItemLevelRepo.GetItemLevelLookUpbyId(tx, filter, itemLevelId)

	if err != nil {
		return get, err
	}
	return get, nil
}

// GetItemLevelLookUp implements masteritemservice.ItemLevelService.
func (s *ItemLevelServiceImpl) GetItemLevelLookUp(filter []utils.FilterCondition, pages pagination.Pagination, itemClassId int) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	get, err := s.structItemLevelRepo.GetItemLevelLookUp(tx, filter, pages, itemClassId)

	if err != nil {
		return get, err
	}
	return get, nil
}

// GetItemLevelDropDown implements masteritemservice.ItemLevelService.
func (s *ItemLevelServiceImpl) GetItemLevelDropDown(itemLevel int) ([]masteritemlevelpayloads.GetItemLevelDropdownResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	get, err := s.structItemLevelRepo.GetItemLevelDropDown(tx, itemLevel)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *ItemLevelServiceImpl) Save(request masteritemlevelpayloads.SaveItemLevelRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()

	if request.ItemLevelId != 0 {
		_, err := s.structItemLevelRepo.GetById(tx, request.ItemLevel, request.ItemLevelId)

		if err != nil {
			return false, err
		}
	}

	save, err := s.structItemLevelRepo.Save(tx, request)

	if err != nil {
		return false, err
	}
	return save, nil
}

func (s *ItemLevelServiceImpl) GetById(itemLevel int, itemLevelId int) (masteritemlevelpayloads.GetItemLevelResponseById, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	get, err := s.structItemLevelRepo.GetById(tx, itemLevel, itemLevelId)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *ItemLevelServiceImpl) GetAll(filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()
	get, err := s.structItemLevelRepo.GetAll(tx, filter, pages)

	if err != nil {
		return get, err
	}
	return get, nil
}

func (s *ItemLevelServiceImpl) ChangeStatus(itemLevel int, itemLevelId int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        fmt.Errorf("panic recovered: %v", r),
			}
		} else if err != nil {
			tx.Rollback()
			logrus.Info("Transaction rollback due to error:", err)
		} else {
			tx.Commit()
			//logrus.Info("Transaction committed successfully")
		}
	}()

	_, err = s.structItemLevelRepo.GetById(tx, itemLevel, itemLevelId)

	if err != nil {
		return false, err
	}

	change_status, err := s.structItemLevelRepo.ChangeStatus(tx, itemLevel, itemLevelId)
	if err != nil {
		return change_status, err
	}

	return true, nil
}
