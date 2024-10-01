package test

import (
	"after-sales/api/config"
	transactionjpcbpayloads "after-sales/api/payloads/transaction/JPCB"
	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	transactionjpcbrepositoryimpl "after-sales/api/repositories/transaction/JPCB/repositories-jpcb-impl"
	transactionjpcbservice "after-sales/api/services/transaction/JPCB"
	transactionjpcbserviceimpl "after-sales/api/services/transaction/JPCB/services-jpcb-impl"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func setupOutstandingJobAlloc() (*gorm.DB, *redis.Client, transactionjpcbservice.OutstandingJobAllocationService) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	OutstandingJobAllocationRepo := transactionjpcbrepositoryimpl.StartOutStandingJobAllocationRepository()
	OperationCodeRepo := masteroperationrepositoryimpl.StartOperationCodeRepositoryImpl()
	OutstandingJobAllocationService := transactionjpcbserviceimpl.StartOutstandingJobAllocationService(OutstandingJobAllocationRepo, OperationCodeRepo, db, rdb)
	return db, rdb, OutstandingJobAllocationService
}

func TestSaveOutstandingJobAlloc(t *testing.T) {
	db, rdb, service := setupOutstandingJobAlloc()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()

	db.Begin()

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Now().Location())

	req := transactionjpcbpayloads.OutstandingJobAllocationSaveRequest{
		CompanyId:      151,
		UserEmployeeId: 2,
		IsExpress:      false,
		OperationId:    3,
		ServiceDate:    date,
	}

	result, err := service.SaveOutstandingJobAllocation("BOOKING", 1, req)
	if err != nil {
		fmt.Println("test run failed because there is an error")
		fmt.Println(err)
	} else {
		fmt.Println("test run successfully")
		fmt.Println(result)
	}
}
