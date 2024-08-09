package test

import (
	"after-sales/api/config"
	transactionjpcbentities "after-sales/api/entities/transaction/JPCB"
	"fmt"
	"testing"
)

func TestGetByIdSupplySlipDetail(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()

	// repo := transactionjpcbrepositoryimpl.OpenBayMasterRepositoryImpl()
	// service := transactionjpcbserviceimpl.StartBayService(repo, db, &redis.Client{})

	// filter := []utils.FilterCondition{}

	// filter = append(filter, utils.FilterCondition{ColumnValue: "1", ColumnField: "company_id"})

	// get, _, _, err := service.GetAllActiveBayCarWashScreen(filter, pagination.Pagination{})
	var bayEntities transactionjpcbentities.BayMaster
	var highestOrder int
	result := db.Model(&bayEntities).Select("MAX(order_number)").Scan(&highestOrder)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Print(highestOrder)
}
