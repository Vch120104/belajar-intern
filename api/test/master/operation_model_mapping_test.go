package test

import (
	"after-sales/api/config"
	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"

	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	"fmt"
	"testing"
)

// func TestStatusMapping(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	operationMappingRepo := operationmodelmappingrepositoryimpl.StartOperationModelMappingRepositoryImpl(db)
// 	operationMappingServ := operationmodelmappingserviceimpl.StartOperationMappingService(operationMappingRepo)

// 	get, err := operationMappingServ.ChangeStatusOperationModelMapping(1)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(get)
// }

func TestGetByIdMapping(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	operationMappingRepo := masteroperationrepositoryimpl.StartOperationModelMappingRepositoryImpl(db)
	operationMappingServ := masteroperationserviceimpl.StartOperationMappingService(operationMappingRepo)

	get, err := operationMappingServ.GetOperationModelMappingById(1)

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}
