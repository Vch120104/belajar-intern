package test

import (
	"after-sales/api/config"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"

	// "after-sales/api/payloads/pagination"
	masterwarehouserepo "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"
	masterwarehousegroupservice "after-sales/api/services/master/warehouse/services-warehouse-impl"
	"fmt"
	"testing"
)

func TestSaveWarehouseGroup(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	warehouseGroupRepo := masterwarehouserepo.OpenWarehouseGroupImpl(db)
	warehouseGroupService := masterwarehousegroupservice.OpenWarehouseGroupService(warehouseGroupRepo)

	save, err := warehouseGroupService.Save(
		masterwarehousepayloads.GetWarehouseGroupResponse{
			IsActive:           true,
			WarehouseGroupCode: "01",
			WarehouseGroupName: "01",
			ProfitCenterId:     01,
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(save)
}

func TestGetWarehouseGroup(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	warehouseGroupRepo := masterwarehouserepo.OpenWarehouseGroupImpl(db)
	warehouseGroupService := masterwarehousegroupservice.OpenWarehouseGroupService(warehouseGroupRepo)

	get, err := warehouseGroupService.GetById(
		1,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}

func TestGetAllWarehouseGroup(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	warehouseGroupRepo := masterwarehouserepo.OpenWarehouseGroupImpl(db)
	warehouseGroupService := masterwarehousegroupservice.OpenWarehouseGroupService(warehouseGroupRepo)

	get, err := warehouseGroupService.GetAll(
		masterwarehousepayloads.GetAllWarehouseGroupRequest{},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}

func TestChangeStatusWarehouseGroup(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	warehouseGroupRepo := masterwarehouserepo.OpenWarehouseGroupImpl(db)
	warehouseGroupService := masterwarehousegroupservice.OpenWarehouseGroupService(warehouseGroupRepo)

	changeStatus, err := warehouseGroupService.ChangeStatus(
		1,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(changeStatus)
}
