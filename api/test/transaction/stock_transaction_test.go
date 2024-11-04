package test

import (
	"after-sales/api/config"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepositoryimpl "after-sales/api/repositories/transaction/sparepart/repositories-sparepart-impl"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setupStockTransaction() transactionsparepartservice.StockTransactionService {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	StockTransactionRepository := transactionsparepartrepositoryimpl.StartStockTransactionRepositoryImpl()
	StockTransactionService := masterserviceimpl.StartStockTransactionServiceImpl(StockTransactionRepository, db, rdb)
	return StockTransactionService
}

func TestStockTransactionInsert(t *testing.T) {
	service := setupStockTransaction()
	currentTime := time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)
	payloads := transactionsparepartpayloads.StockTransactionInsertPayloads{
		CompanyId:                    421,
		StockTransactionId:           0,
		TransactionTypeId:            12,
		TransactionReasonId:          8,
		ReferenceId:                  0,
		ReferenceDocumentNumber:      "Dummy Reference Number",
		ReferenceDate:                &currentTime,
		ReferenceWarehouseId:         7,
		ReferenceWarehouseGroupId:    1,
		ReferenceLocationId:          1,
		ReferenceItemId:              1,
		ReferenceQuantity:            1,
		ReferenceUnitOfMeasurementId: 1,
		ReferencePrice:               1,
		ReferenceCurrencyId:          1,
		TransactionCogs:              1,
		ChangeNo:                     1,
		CreatedByUserId:              1,
		CreatedDate:                  currentTime,
		UpdatedByUserId:              1,
		UpdatedDate:                  currentTime,
		VehicleId:                    1,
		ItemClassId:                  1,
	}
	_, err := service.StockTransactionInsert(payloads)
	assert.Nil(t, err, "Error should be nil")
}
