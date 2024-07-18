package test

import (
	"after-sales/api/config"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	"after-sales/api/utils"

	// "after-sales/api/payloads/pagination"
	masterwarehouserepo "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"
	masterwarehousegroupservice "after-sales/api/services/master/warehouse/services-warehouse-impl"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveWarehouseGroup(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	warehouseGroupRepo := masterwarehouserepo.OpenWarehouseGroupImpl()
	warehouseGroupService := masterwarehousegroupservice.OpenWarehouseGroupService(warehouseGroupRepo, db, nil)

	save, err := warehouseGroupService.SaveWarehouseGroup(
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
	warehouseGroupRepo := masterwarehouserepo.OpenWarehouseGroupImpl()
	warehouseGroupService := masterwarehousegroupservice.OpenWarehouseGroupService(warehouseGroupRepo, db, nil)

	get, err := warehouseGroupService.GetByIdWarehouseGroup(
		1,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}

func TestGetAllWarehouseGroup(t *testing.T) {
	// Initialize environment configurations
	config.InitEnvConfigs(true, "")

	// Initialize the database
	db := config.InitDB()

	// Initialize the repository and service
	warehouseGroupRepo := masterwarehouserepo.OpenWarehouseGroupImpl()
	warehouseGroupService := masterwarehousegroupservice.OpenWarehouseGroupService(warehouseGroupRepo, db, nil)

	// Define filter conditions (if any)
	filterConditions := []utils.FilterCondition{}

	// Define pagination parameters
	paginationParams := pagination.Pagination{
		Page:  1,
		Limit: 10,
	}

	// Make the request to the service
	paginatedResult, err := warehouseGroupService.GetAllWarehouseGroup(filterConditions, paginationParams)

	// Handle any errors
	if err != nil {
		panic(err)
	}

	// Print the response for debugging purposes
	fmt.Println(paginatedResult)

	// Add assertions to validate the response
	assert.NotNil(t, paginatedResult)
	assert.GreaterOrEqual(t, nil, 0, "Length of warehouse groups should be >= 0")
}

func TestChangeStatusWarehouseGroup(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	warehouseGroupRepo := masterwarehouserepo.OpenWarehouseGroupImpl()
	warehouseGroupService := masterwarehousegroupservice.OpenWarehouseGroupService(warehouseGroupRepo, db, nil)

	changeStatus, err := warehouseGroupService.ChangeStatusWarehouseGroup(
		1,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(changeStatus)
}
