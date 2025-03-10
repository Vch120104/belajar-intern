package transactionsparepartserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StockOpnameServiceImpl struct {
	Repository transactionsparepartrepository.StockOpnameRepository
	DB         *gorm.DB
	Redis      *redis.Client
}

func NewStockOpnameServiceImpl(
	repository transactionsparepartrepository.StockOpnameRepository,
	db *gorm.DB,
	redis *redis.Client,
) transactionsparepartservice.StockOpnameService {
	return &StockOpnameServiceImpl{
		Repository: repository,
		DB:         db,
		Redis:      redis,
	}
}

// func (s *StockOpnameServiceImpl) GetAllStockOpname(filteredCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
// 	tx := s.DB.Begin()
// 	var err *exceptions.BaseErrorResponse

// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			err = &exceptions.BaseErrorResponse{
// 				StatusCode: http.StatusInternalServerError,
// 				Err:        fmt.Errorf("panic recovered: %v", r),
// 			}
// 		} else if err != nil {
// 			tx.Rollback()
// 			logrus.Info("Transaction rollback due to error:", err)
// 		} else {
// 			if commitErr := tx.Commit().Error; commitErr != nil {
// 				logrus.WithError(commitErr).Error("Transaction commit failed")
// 				err = &exceptions.BaseErrorResponse{
// 					StatusCode: http.StatusInternalServerError,
// 					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
// 				}
// 			}
// 		}
// 	}()

// 	pages, err = s.Repository.GetAllStockOpname(tx, filteredCondition, pages)
// 	if err != nil {
// 		return pages, err
// 	}
// 	return pages, nil
// }

func (s *StockOpnameServiceImpl) GetAllStockOpname(filteredCondition []utils.FilterCondition, pages pagination.Pagination, companyCode float64, dateParams map[string]interface{}) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	pages, err = s.Repository.GetAllStockOpname(tx, filteredCondition, pages, companyCode, dateParams)
	if err != nil {
		return pages, err
	}
	return pages, nil
}

func (s *StockOpnameServiceImpl) GetLocationList(filteredCondition []utils.FilterCondition, pages pagination.Pagination,
	companyCode float64, warehouseGroup string, warehouseCode string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	pages, err = s.Repository.GetLocationList(tx, filteredCondition, pages, companyCode, warehouseGroup, warehouseCode)
	if err != nil {
		return pages, err
	}
	return pages, nil
}

func (s *StockOpnameServiceImpl) GetPersonInChargeList(filteredCondition []utils.FilterCondition, pages pagination.Pagination, companyCode float64) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	pages, err = s.Repository.GetPersonInChargeList(tx, filteredCondition, pages, companyCode)
	if err != nil {
		return pages, err
	}
	return pages, nil
}

func (s *StockOpnameServiceImpl) GetItemList(pages pagination.Pagination, whsCode string,
	itemGroup string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	pages, err = s.Repository.GetItemList(tx, pages, whsCode, itemGroup)
	if err != nil {
		return pages, err
	}
	return pages, nil
}

// func (s *StockOpnameServiceImpl) GetListForOnGoing(sysNo string) *exceptions.BaseErrorResponse {
// 	tx := s.DB.Begin()
// 	var err *exceptions.BaseErrorResponse

// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			err = &exceptions.BaseErrorResponse{
// 				StatusCode: http.StatusInternalServerError,
// 				Err:        fmt.Errorf("panic recovered: %v", r),
// 			}
// 		} else if err != nil {
// 			tx.Rollback()
// 			logrus.Info("Transaction rollback due to error:", err)
// 		} else {
// 			if commitErr := tx.Commit().Error; commitErr != nil {
// 				logrus.WithError(commitErr).Error("Transaction commit failed")
// 				err = &exceptions.BaseErrorResponse{
// 					StatusCode: http.StatusInternalServerError,
// 					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
// 				}
// 			}
// 		}
// 	}()

// 	err = s.Repository.GetListForOnGoing(tx, sysNo)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (s *StockOpnameServiceImpl) GetOnGoingStockOpname(companyCode float64, sysNo float64) ([]transactionsparepartpayloads.GetOnGoingStockOpnameResponse, *exceptions.BaseErrorResponse) {
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	datas, err := s.Repository.GetOnGoingStockOpname(tx, companyCode, sysNo)
	if err != nil {
		return datas, err
	}
	return datas, nil
}

func (s *StockOpnameServiceImpl) InsertNewStockOpname(request transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse) {
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	_, err = s.Repository.InsertNewStockOpname(tx, request)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StockOpnameServiceImpl) UpdateOnGoingStockOpname(sysNo float64, request transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse) {
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
			if commitErr := tx.Commit().Error; commitErr != nil {
				logrus.WithError(commitErr).Error("Transaction commit failed")
				err = &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        fmt.Errorf("failed to commit transaction: %w", commitErr),
				}
			}
		}
	}()

	_, err = s.Repository.UpdateOnGoingStockOpname(tx, sysNo, request)
	if err != nil {
		return false, err
	}
	return true, nil
}