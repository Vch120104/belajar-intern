package test

import (
	"after-sales/api/config"
	transactionworkshoprepositoryimpl "after-sales/api/repositories/transaction/workshop/repositories-workshop-impl"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	transactionworkshopserviceimpl "after-sales/api/services/transaction/workshop/services-workshop-impl"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"testing"
)

func setupPurchaseOrder() (*gorm.DB, *redis.Client, transactionworkshopservice.ServiceReceiptService) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	ServiceReceiptRepository := transactionworkshoprepositoryimpl.OpenServiceReceiptRepositoryImpl()
	ServiceReceiptService := transactionworkshopserviceimpl.OpenServiceReceiptServiceImpl(ServiceReceiptRepository, db, rdb)
	return db, rdb, ServiceReceiptService
}
func TestGetAllPurchaseOrder(t *testing.T) {
	//t.Error("das")
	db, rdb, _ := setupPurchaseOrder()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()
	db.Begin()
	//err := db.s

}
