package transactionsparepartserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"after-sales/api/utils"

	"github.com/redis/go-redis/v9"
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

func (s *StockOpnameServiceImpl) GetAllStockOpname(
	filteredCondition []utils.FilterCondition, pages pagination.Pagination,
	dateParams map[string]interface{}) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: 500,
				Err:        r.(error),
			}
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = &exceptions.BaseErrorResponse{
					StatusCode: 500,
					Err:        commitErr,
				}
			}
		}
	}()

	res, err := s.Repository.GetAllStockOpname(tx, filteredCondition, pages, dateParams)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *StockOpnameServiceImpl) GetAllStockOpnameDetail(pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: 500,
				Err:        r.(error),
			}
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = &exceptions.BaseErrorResponse{
					StatusCode: 500,
					Err:        commitErr,
				}
			}
		}
	}()
	res, err := s.Repository.GetAllStockOpnameDetail(tx, pages)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *StockOpnameServiceImpl) GetStockOpnameByStockOpnameSystemNumber(stockOpnameSystemNumber int) (
	[]transactionsparepartpayloads.GetStockOpnameByStockOpnameSystemNumberResponse, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: 500,
				Err:        r.(error),
			}
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = &exceptions.BaseErrorResponse{
					StatusCode: 500,
					Err:        commitErr,
				}
			}
		}
	}()

	res, err := s.Repository.GetStockOpnameByStockOpnameSystemNumber(tx, stockOpnameSystemNumber)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *StockOpnameServiceImpl) GetStockOpnameAllDetailByStockOpnameSystemNumber(stockOpnameSystemNumber int, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: 500,
				Err:        r.(error),
			}
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = &exceptions.BaseErrorResponse{
					StatusCode: 500,
					Err:        commitErr,
				}
			}
		}
	}()

	res, err := s.Repository.GetStockOpnameAllDetailByStockOpnameSystemNumber(tx, stockOpnameSystemNumber, pages)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *StockOpnameServiceImpl) InsertStockOpname(request transactionsparepartpayloads.StockOpnameInsertRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: 500,
				Err:        r.(error),
			}
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = &exceptions.BaseErrorResponse{
					StatusCode: 500,
					Err:        commitErr,
				}
			}
		}
	}()

	_, err = s.Repository.InsertStockOpname(tx, request)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StockOpnameServiceImpl) SubmitStockOpname(systemNumber int, request transactionsparepartpayloads.StockOpnameSubmitRequest) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: 500,
				Err:        r.(error),
			}
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = &exceptions.BaseErrorResponse{
					StatusCode: 500,
					Err:        commitErr,
				}
			}
		}
	}()

	_, err = s.Repository.SubmitStockOpname(tx, systemNumber, request)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StockOpnameServiceImpl) InsertStockOpnameDetail(request transactionsparepartpayloads.StockOpnameInsertDetailRequest,
	systemNumber int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: 500,
				Err:        r.(error),
			}
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = &exceptions.BaseErrorResponse{
					StatusCode: 500,
					Err:        commitErr,
				}
			}
		}
	}()

	_, err = s.Repository.InsertStockOpnameDetail(tx, request, systemNumber)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StockOpnameServiceImpl) UpdateStockOpname(request transactionsparepartpayloads.StockOpnameInsertRequest,
	systemNumber int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: 500,
				Err:        r.(error),
			}
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = &exceptions.BaseErrorResponse{
					StatusCode: 500,
					Err:        commitErr,
				}
			}
		}
	}()

	_, err = s.Repository.UpdateStockOpname(tx, request, systemNumber)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *StockOpnameServiceImpl) UpdateStockOpnameDetail(request transactionsparepartpayloads.StockOpnameUpdateDetailRequest,
	systemNumber int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: 500,
				Err:        r.(error),
			}
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = &exceptions.BaseErrorResponse{
					StatusCode: 500,
					Err:        commitErr,
				}
			}
		}
	}()

	_, err = s.Repository.UpdateStockOpnameDetail(tx, request, systemNumber)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StockOpnameServiceImpl) DeleteStockOpname(systemNumber int) (bool, *exceptions.BaseErrorResponse) {
	tx := s.DB.Begin()
	var err *exceptions.BaseErrorResponse

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = &exceptions.BaseErrorResponse{
				StatusCode: 500,
				Err:        r.(error),
			}
		} else if err != nil {
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = &exceptions.BaseErrorResponse{
					StatusCode: 500,
					Err:        commitErr,
				}
			}
		}
	}()
	_, err = s.Repository.DeleteStockOpname(tx, systemNumber)
	if err != nil {
		return false, err
	}
	return true, nil
}

// func (s *StockOpnameServiceImpl) GetAllStockOpname(filteredCondition []utils.FilterCondition, pages pagination.Pagination, companyCode float64, dateParams map[string]interface{}) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

// 	pages, err = s.Repository.GetAllStockOpname(tx, filteredCondition, pages, companyCode, dateParams)
// 	if err != nil {
// 		return pages, err
// 	}
// 	return pages, nil
// }

// func (s *StockOpnameServiceImpl) GetLocationList(filteredCondition []utils.FilterCondition, pages pagination.Pagination,
// 	companyCode float64, warehouseGroup string, warehouseCode string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

// 	pages, err = s.Repository.GetLocationList(tx, filteredCondition, pages, companyCode, warehouseGroup, warehouseCode)
// 	if err != nil {
// 		return pages, err
// 	}
// 	return pages, nil
// }

// func (s *StockOpnameServiceImpl) GetPersonInChargeList(filteredCondition []utils.FilterCondition, pages pagination.Pagination, companyCode float64) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

// 	pages, err = s.Repository.GetPersonInChargeList(tx, filteredCondition, pages, companyCode)
// 	if err != nil {
// 		return pages, err
// 	}
// 	return pages, nil
// }

// func (s *StockOpnameServiceImpl) GetItemList(pages pagination.Pagination, whsCode string,
// 	itemGroup string) (pagination.Pagination, *exceptions.BaseErrorResponse) {
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

// 	pages, err = s.Repository.GetItemList(tx, pages, whsCode, itemGroup)
// 	if err != nil {
// 		return pages, err
// 	}
// 	return pages, nil
// }

// func (s *StockOpnameServiceImpl) GetOnGoingStockOpname(companyCode float64, sysNo float64) ([]transactionsparepartpayloads.GetOnGoingStockOpnameResponse, *exceptions.BaseErrorResponse) {
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

// 	datas, err := s.Repository.GetOnGoingStockOpname(tx, companyCode, sysNo)
// 	if err != nil {
// 		return datas, err
// 	}
// 	return datas, nil
// }

// func (s *StockOpnameServiceImpl) InsertNewStockOpname(request transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse) {
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

// 	_, err = s.Repository.InsertNewStockOpname(tx, request)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func (s *StockOpnameServiceImpl) UpdateOnGoingStockOpname(sysNo float64, request transactionsparepartpayloads.InsertNewStockOpnameRequest) (bool, *exceptions.BaseErrorResponse) {
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

// 	_, err = s.Repository.UpdateOnGoingStockOpname(tx, sysNo, request)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
