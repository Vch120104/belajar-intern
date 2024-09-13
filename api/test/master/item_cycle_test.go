package test

import (
	"after-sales/api/config"
	masterpayloads "after-sales/api/payloads/master"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	masterservice "after-sales/api/services/master"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func setupItemCycle() (*gorm.DB, *redis.Client, masterservice.ItemCycleService) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	ItemCycleRepository := masterrepositoryimpl.NewItemCycleRepositoryImpl()
	ItemCycleService := masterserviceimpl.NewItemCycleServiceImpl(ItemCycleRepository, db, rdb)
	return db, rdb, ItemCycleService
}
func TestItemCycle(t *testing.T) {
	db, rdb, ItemCycleService := setupItemCycle()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()

	payloads := masterpayloads.ItemCycleInsertPayloads{
		CompanyId:         1,
		PeriodYear:        "2023",
		PeriodMonth:       "11",
		ItemId:            1,
		OrderCycle:        12,
		QuantityOnOrder:   13,
		QuantityBackOrder: 14,
	}
	_, err := ItemCycleService.ItemCycleInsert(payloads)
	assert.Nil(t, err)
}
