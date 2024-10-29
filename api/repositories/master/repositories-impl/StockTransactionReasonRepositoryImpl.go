package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterpayloads "after-sales/api/payloads/master"
	"after-sales/api/payloads/pagination"
	masterrepository "after-sales/api/repositories/master"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

type StockTransactionReasonRepositoryImpl struct {
}

func StartStockTraansactionReasonRepositoryImpl() masterrepository.StockTransactionReasonRepository {

	return &StockTransactionReasonRepositoryImpl{}
}

func (s *StockTransactionReasonRepositoryImpl) InsertStockTransactionReason(db *gorm.DB, payloads masterpayloads.StockTransactionReasonInsertPayloads) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse) {
	StockTransactionEntities := masterentities.StockTransactionReason{
		StockTransactionReasonCode:  payloads.StockTransactionReasonCode,
		StockTransactionDescription: payloads.StockTransactionDescription}
	err := db.Create(&StockTransactionEntities).Error
	if err != nil {
		return StockTransactionEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("error on insert stock transaction : %s", err.Error()),
		}
	}
	return StockTransactionEntities, nil
}

func (s *StockTransactionReasonRepositoryImpl) GetStockTransactionReasonById(db *gorm.DB, id int) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse) {
	StockTransactionReasonEntities := masterentities.StockTransactionReason{}
	err := db.Model(&StockTransactionReasonEntities).Where("stock_transaction_reason_id =?", id).
		First(&StockTransactionReasonEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return StockTransactionReasonEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
			}
		}
		return StockTransactionReasonEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error on stock transaction by id",
		}
	}
	return StockTransactionReasonEntities, nil
}

func (s *StockTransactionReasonRepositoryImpl) GetAllStockTransactionReason(db *gorm.DB, conditions []utils.FilterCondition, paginationParams pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var StockTransactionReason masterentities.StockTransactionReason
	var Responses []masterentities.StockTransactionReason

	Jointable := db.Model(&StockTransactionReason)
	WhereQuery := utils.ApplyFilter(Jointable, conditions)

	err := WhereQuery.Scopes(pagination.Paginate(&StockTransactionReason, &paginationParams, WhereQuery)).Order("stock_transaction_reason_id").Scan(&Responses).Error
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

func (s *StockTransactionReasonRepositoryImpl) GetStockTransactionReasonByCode(db *gorm.DB, Code string) (masterentities.StockTransactionReason, *exceptions.BaseErrorResponse) {
	StockTransactionReasonEntities := masterentities.StockTransactionReason{}
	err := db.Model(&StockTransactionReasonEntities).Where("stock_transaction_reason_code =?", Code).
		First(&StockTransactionReasonEntities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return StockTransactionReasonEntities, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
			}
		}
		return StockTransactionReasonEntities, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error on getting code",
		}
	}
	return StockTransactionReasonEntities, nil
}
