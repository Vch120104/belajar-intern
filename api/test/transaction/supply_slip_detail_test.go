package test

import (
	"after-sales/api/config"
	transactionsparepartrepositoryimpl "after-sales/api/repositories/transaction/sparepart/repositories-sparepart-impl"
	transactionsparepartserviceimpl "after-sales/api/services/transaction/sparepart/services-sparepart-impl"
	"fmt"
	"testing"
)

func TestGetByIdSupplySlipDetail(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	supplySlipDetailRepo := transactionsparepartrepositoryimpl.StartSupplySlipRepositoryImpl(db)
	supplySlipDetailServ := transactionsparepartserviceimpl.StartSupplySlipService(supplySlipDetailRepo)

	get, err := supplySlipDetailServ.GetSupplySlipDetailById(6)

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}
