package test

import (
	"after-sales/api/config"
	"after-sales/api/payloads/pagination"
	transactionjpcbrepositoryimpl "after-sales/api/repositories/transaction/JPCB/repositories-jpcb-impl"
	transactionjpcbserviceimpl "after-sales/api/services/transaction/JPCB/services-jpcb-impl"
	"after-sales/api/utils"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestGetByIdSupplySlipDetail(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()

	repo := transactionjpcbrepositoryimpl.OpenBayMasterRepositoryImpl()
	service := transactionjpcbserviceimpl.StartBayService(repo, db, &redis.Client{})

	filter := []utils.FilterCondition{}

	filter = append(filter, utils.FilterCondition{ColumnValue: "1", ColumnField: "company_id"})

	get, _, _, err := service.GetAllActiveBayCarWashScreen(filter, pagination.Pagination{})

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}
