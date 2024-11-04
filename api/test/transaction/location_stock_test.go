package test

import (
	"after-sales/api/config"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterservice "after-sales/api/services/master"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setupLocationStock() masterservice.LocationStockService {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	LocationStockRepository := masterwarehouserepository.NewLocationStockRepositoryImpl()
	LocationStockService := masterserviceimpl.NewLocationStockServiceImpl(LocationStockRepository, db, rdb)
	return LocationStockService
}
func TestUpdateLocationStock(t *testing.T) {
	service := setupLocationStock()
	//defer func() {
	//	sqlDB, _ := db.DB()
	//	sqlDB.Close()
	//	rdb.Close()
	//}()
	timeNow := time.Now().UTC()
	payloads := masterwarehousepayloads.LocationStockUpdatePayloads{
		CompanyId:                3,
		PeriodYear:               "2013",
		PeriodMonth:              "02",
		WarehouseId:              1,
		LocationId:               1,
		ItemId:                   293773,
		WarehouseGroupId:         1,
		QuantityBegin:            1,
		QuantitySales:            1,
		QuantitySalesReturn:      1,
		QuantityPurchase:         100,
		QuantityPurchaseReturn:   1,
		QuantityTransferIn:       1,
		QuantityTransferOut:      1,
		QuantityClaimIn:          1,
		QuantityClaimOut:         1,
		QuantityAdjustment:       5,
		QuantityAllocated:        2,
		QuantityInTransit:        2,
		QuantityEnding:           5,
		QuantityRobbingIn:        3,
		QuantityRobbingOut:       2,
		QuantityAssemblyIn:       5,
		QuantityAssemblyOut:      5,
		StockTransactionTypeId:   6,
		StockTransactionReasonId: 3,
		CreatedByUserId:          1,
		CreatedDate:              &timeNow,
		UpdatedByUserId:          1,
		UpdatedDate:              &timeNow,
	}
	_, err := service.UpdateLocationStock(payloads)
	assert.Nil(t, err)
}
func TestGetLocationStock(t *testing.T) {
	//db, rdb, service := setupLocationStock()
	//defer func() {
	//	sqlDB, _ := db.DB()
	//	sqlDB.Close()
	//	rdb.Close()s
	//}()
}
