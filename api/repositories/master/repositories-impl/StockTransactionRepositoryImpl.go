package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"errors"
	"gorm.io/gorm"
	"net/http"
)

type StockTransactionRepositoryImpl struct {
}

func NewStockTransactionRepositoryImpl() masterrepository.StockTransactionTypeRepository {
	return &StockTransactionRepositoryImpl{}
}
func (s *StockTransactionRepositoryImpl) GetStockTransactionTypeByCode(db *gorm.DB, Code string) (masterentities.StockTransactionType, *exceptions.BaseErrorResponse) {
	var StockTransaction masterentities.StockTransactionType
	err := db.Model(&StockTransaction).Where(masterentities.StockTransactionType{StockTransactionTypeCode: Code}).First(&StockTransaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return StockTransaction, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        errors.New("stock transaction type with code :" + Code + "not found"),
			}
		}
		return StockTransaction, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	return StockTransaction, nil
}

func (s *StockTransactionRepositoryImpl) GetAllStockTransactionType(db *gorm.DB, conditions []utils.FilterCondition, paginationParams pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var StockTransaction masterentities.StockTransactionType
	var Responses []masterentities.StockTransactionType

	Jointable := db.Model(&StockTransaction)
	WhereQuery := utils.ApplyFilter(Jointable, conditions)

	err := WhereQuery.Scopes(pagination.Paginate(&StockTransaction, &paginationParams, WhereQuery)).Order("stock_transaction_type_id").Scan(&Responses).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return paginationParams, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	if len(Responses) == 0 {
		paginationParams.Rows = []string{}
		return paginationParams, nil
	}
	paginationParams.Rows = Responses
	return paginationParams, nil
}
