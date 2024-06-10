package test

import (
	"after-sales/api/config"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"

	// "after-sales/api/payloads/pagination"

	masterwarehouserepo "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"
	masterwarehousegroupservice "after-sales/api/services/master/warehouse/services-warehouse-impl"
	"fmt"
	"testing"
)

func TestSaveWarehouseLocation(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	warehouseLocationRepo := masterwarehouserepo.OpenWarehouseLocationImpl()
	warehouseLocationService := masterwarehousegroupservice.OpenWarehouseLocationService(warehouseLocationRepo, db, nil)

	save, err := warehouseLocationService.Save(
		masterwarehousepayloads.GetWarehouseLocationResponse{
			IsActive:                      true,
			CompanyId:                     1,
			WarehouseGroupId:              1,
			WarehouseLocationCode:         "ADA",
			WarehouseLocationName:         "ADAADA",
			WarehouseLocationDetailName:   "ADA ADA ADA",
			WarehouseLocationPickSequence: 1,
			WarehouseLocationCapacityInM3: 1,
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(save)
}

func TestGetWarehouseLocationById(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	warehouseLocationRepo := masterwarehouserepo.OpenWarehouseLocationImpl()
	warehouseLocationService := masterwarehousegroupservice.OpenWarehouseLocationService(warehouseLocationRepo, db, nil)

	update, err := warehouseLocationService.GetById(
		2,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(update)
}

func TestGetAllWarehouseLocation(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	warehouseLocationRepo := masterwarehouserepo.OpenWarehouseLocationImpl()
	warehouseLocationService := masterwarehousegroupservice.OpenWarehouseLocationService(warehouseLocationRepo, db, nil)

	update, err := warehouseLocationService.GetAll(
		masterwarehousepayloads.GetAllWarehouseLocationRequest{}, pagination.Pagination{
			Page:  0,
			Limit: 10,
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(update)
}

func TestChangeStatusWarehouseLocationById(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	warehouseLocationRepo := masterwarehouserepo.OpenWarehouseLocationImpl()
	warehouseLocationService := masterwarehousegroupservice.OpenWarehouseLocationService(warehouseLocationRepo, db, nil)

	update, err := warehouseLocationService.ChangeStatus(
		2,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(update)
}
