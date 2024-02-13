package test

// import (
// 	"after-sales/api/config"
// 	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"

// 	// "after-sales/api/payloads/pagination"

// 	masterwarehouserepo "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"
// 	masterwarehousegroupservice "after-sales/api/services/master/warehouse/services-warehouse-impl"
// 	"fmt"
// 	"testing"
// )

// func TestSaveWarehouseMaster(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	warehouseMasterRepo := masterwarehouserepo.OpenWarehouseMasterImpl(db)
// 	warehouseMasterService := masterwarehousegroupservice.OpenWarehouseMasterService(warehouseMasterRepo)

// 	save, err := warehouseMasterService.Save(
// 		masterwarehousepayloads.GetWarehouseMasterResponse{
// 			IsActive:                      true,
// 			WarehouseCostingType:          "01",
// 			WarehouseKaroseri:             true,
// 			WarehouseNegativeStock:        true,
// 			WarehouseReplishmentIndicator: true,
// 			WarehouseContact:              "01",
// 			WarehouseCode:                 "01",
// 			AddressId:                     1,
// 			BrandId:                       1,
// 			SupplierId:                    1,
// 			UserId:                        1,
// 			WarehouseSalesAllow:           true,
// 			WarehouseInTransit:            true,
// 			WarehouseName:                 "1",
// 			WarehouseDetailName:           "1",
// 			WarehouseTransitDefault:       "1",
// 		},
// 	)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(save)
// }

// func TestGetWarehouseMaster(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	warehouseMasterRepo := masterwarehouserepo.OpenWarehouseMasterImpl(db)
// 	warehouseMasterService := masterwarehousegroupservice.OpenWarehouseMasterService(warehouseMasterRepo)

// 	save, err := warehouseMasterService.GetById(
// 		1,
// 	)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(save)
// }
