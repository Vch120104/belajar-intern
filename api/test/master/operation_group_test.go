package test

import (
	"after-sales/api/config"
	masteroperationpayloads "after-sales/api/payloads/master/operation"
	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
	"fmt"
	"testing"
)

func TestChangeStatusOperationGroup(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	operationGroupRepo := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl(db)
	operationGroupServ := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepo)

	get, err := operationGroupServ.ChangeStatusOperationGroup(1)

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}

func TestGetOperationGroupById(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	operationGroupRepo := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl(db)
	operationGroupServ := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepo)

	get, err := operationGroupServ.GetOperationGroupById(2)

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}

// func TestGetAllOperationGroup(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	operationGroupRepo := operationgrouprepoimpl.StartOperationGroupRepositoryImpl(db)
// 	operationGroupServ := operationgroupserviceimpl.StartOperationGroupService(operationGroupRepo)

// 	queryOf := []string{} // Replace with actual query values
// 	queryBy := []string{} // Replace with actual query fields
// 	sortOf := "operation_group_code"
// 	sortBy := "asc"

// 	get, err := operationGroupServ.GetAllOperationGroup(queryOf, queryBy, sortOf, sortBy)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(get)
// }

func TestSaveOperationGroup(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	operationGroupRepo := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl(db)
	operationGroupServ := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepo)
	//if OperationGroupId = 0 -> insert; else update
	req := masteroperationpayloads.OperationGroupResponse{
		IsActive:                  true,
		OperationGroupId:          0,
		OperationGroupCode:        "A1",
		OperationGroupDescription: "test",
	}

	get, err := operationGroupServ.SaveOperationGroup(req)

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}

func TestUpdateGroupById(t *testing.T) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	operationGroupRepo := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl(db)
	operationGroupServ := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepo)
	req := masteroperationpayloads.OperationGroupResponse{
		IsActive: false,
	}

	get, err := operationGroupServ.SaveOperationGroup(req)

	if err != nil {
		panic(err)
	}

	fmt.Println(get)
}
