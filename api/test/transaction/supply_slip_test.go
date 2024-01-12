package test

import (
	"after-sales/api/config"
	transactionsparepartrepositoryimpl "after-sales/api/repositories/transaction/sparepart/repositories-sparepart-impl"
	transactionsparepartserviceimpl "after-sales/api/services/transaction/sparepart/services-sparepart-impl"
	"fmt"
	"testing"
)

func TestGetByIdSupplySlip(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	supplySlipRepo := transactionsparepartrepositoryimpl.StartSupplySlipRepositoryImpl(db)
	supplySlipServ := transactionsparepartserviceimpl.StartSupplySlipService(supplySlipRepo)

	get, err := supplySlipServ.GetSupplySlipById(1)

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}
