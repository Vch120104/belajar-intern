package masterrepositoryimpl

import (
	masterentities "after-sales/api/entities/master"
	"after-sales/api/exceptions"
	masterrepository "after-sales/api/repositories/master"
	"errors"
	"gorm.io/gorm"
	"net/http"
)

type StockTransactionReasonRepositoryImpl struct {
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

func StartStockTraansactionReasonRepositoryImpl() masterrepository.StockTransactionReasonRepository {

	return &StockTransactionReasonRepositoryImpl{}
}
