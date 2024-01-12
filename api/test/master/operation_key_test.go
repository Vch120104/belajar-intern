package test

import (
	"after-sales/api/config"
	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"

	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	"fmt"
	"testing"
)

func TestGetByIdKey(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	operationKeyRepo := masteroperationrepositoryimpl.StartOperationKeyRepositoryImpl(db)
	operationKeyServ := masteroperationserviceimpl.StartOperationKeyService(operationKeyRepo)

	get, err := operationKeyServ.GetOperationKeyById(1)

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}
