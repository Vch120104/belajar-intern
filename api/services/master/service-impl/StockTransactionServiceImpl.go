package masterserviceimpl

import (
	"after-sales/api/exceptions"
	"after-sales/api/helper"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepository "after-sales/api/repositories/transaction/sparepart"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type StockTransactionServiceImpl struct {
	StockTransactionRepository transactionsparepartrepository.StockTransactionRepository
	db                         *gorm.DB
	rdb                        *redis.Client
}

func StartStockTransactionServiceImpl(StockTransactionRepository transactionsparepartrepository.StockTransactionRepository, db *gorm.DB, rdb *redis.Client) transactionsparepartservice.StockTransactionService {
	return &StockTransactionServiceImpl{
		StockTransactionRepository: StockTransactionRepository,
		db:                         db,
		rdb:                        rdb,
	}
}
func (s *StockTransactionServiceImpl) StockTransactionInsert(payloads transactionsparepartpayloads.StockTransactionInsertPayloads) (bool, *exceptions.BaseErrorResponse) {
	tx := s.db.Begin()
	result, err := s.StockTransactionRepository.StockTransactionInsert(tx, payloads)
	defer helper.CommitOrRollbackTrx(tx)
	if err != nil {
		return result, err
	}
	return result, nil
}
