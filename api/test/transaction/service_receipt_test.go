package test

import (
	"after-sales/api/config"
	transactionworkshopcontroller "after-sales/api/controllers/transactions/workshop"
	"after-sales/api/payloads/pagination"
	transactionworkshoprepositoryimpl "after-sales/api/repositories/transaction/workshop/repositories-workshop-impl"
	transactionworkshopservice "after-sales/api/services/transaction/workshop"
	transactionworkshopserviceimpl "after-sales/api/services/transaction/workshop/services-workshop-impl"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func setupServiceReceipt() (*gorm.DB, *redis.Client, transactionworkshopservice.ServiceReceiptService) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	ServiceReceiptRepository := transactionworkshoprepositoryimpl.OpenServiceReceiptRepositoryImpl()
	ServiceReceiptService := transactionworkshopserviceimpl.OpenServiceReceiptServiceImpl(ServiceReceiptRepository, db, rdb)
	return db, rdb, ServiceReceiptService
}

func TestGetServiceReceiptById_Success(t *testing.T) {
	db, rdb, ServiceReceiptService := setup()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()

	req, _ := http.NewRequest("GET", "/v1/service-receipt/2", nil)
	rr := httptest.NewRecorder()

	controller := transactionworkshopcontroller.NewServiceRequestController(ServiceReceiptService)
	controller.GetById(rr, req)

	pagination := pagination.Pagination{
		Limit: 10,
		Page:  0,
	}

	result, _ := ServiceReceiptService.GetById(1, pagination)

	fmt.Println(result)

}

func BenchmarkGetByIdServiceReceipt(b *testing.B) {
	db, rdb, ServiceReceiptService := setupServiceReceipt()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()

	paginate := pagination.Pagination{
		Limit:  10,
		Page:   0,
		SortOf: "",
		SortBy: "",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ServiceReceiptService.GetById(2, paginate)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}
