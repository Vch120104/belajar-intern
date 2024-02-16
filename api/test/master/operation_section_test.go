package test

// import (
// 	"after-sales/api/config"
// 	masteroperationpayloads "after-sales/api/payloads/master/operation"
// 	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
// 	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
// 	"fmt"
// 	"testing"
// )

// func TestStatusSection(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	operationSectionRepo := masteroperationrepositoryimpl.StartOperationSectionRepositoryImpl(db)
// 	operationSectionServ := masteroperationserviceimpl.StartOperationSectionService(operationSectionRepo)

// 	get, err := operationSectionServ.ChangeStatusOperationSection(1)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(get)
// }

// func TestGetAllSection(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	operationSectionRepo := masteroperationrepositoryimpl.StartOperationSectionRepositoryImpl(db)
// 	operationSectionServ := masteroperationserviceimpl.StartOperationSectionService(operationSectionRepo)

// 	get, err := operationSectionServ.GetAllOperationSection()

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(get)
// }

// // func TestGetAllSectionList(t *testing.T) {
// // 	config.InitEnvConfigs(true, "")
// // 	db := config.InitDB()
// // 	operationSectionRepo := masteroperationrepositoryimpl.StartOperationSectionRepositoryImpl(db)
// // 	operationSectionServ := masteroperationserviceimpl.StartOperationSectionService(operationSectionRepo)

// // 	get, err := operationSectionServ.GetAllOperationSectionList()

// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	fmt.Println(get)
// // }

// func TestGetSectionDescription(t *testing.T) {
// 	// config.InitEnvConfigs(true, "")
// 	// db := config.InitDB()

// 	// operationSectionRepo := masteroperationrepositoryimpl.StartOperationSectionRepositoryImpl(db)
// 	// operationSectionServ := masteroperationserviceimpl.StartOperationSectionService(operationSectionRepo)

// 	// // get, err := operationSectionServ.GetOperationSectionDescription("as", "ad")

// 	// // if err != nil {
// 	// // 	panic(err)
// 	// // }

// 	// // fmt.Println(get)
// }

// func TestSaveOperationSection(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()

// 	operationSectionRepo := masteroperationrepositoryimpl.StartOperationSectionRepositoryImpl(db)
// 	operationSectionServ := masteroperationserviceimpl.StartOperationSectionService(operationSectionRepo)
// 	req := masteroperationpayloads.OperationSectionRequest{
// 		IsActive:                    true,
// 		OperationSectionCode:        "gac",
// 		OperationGroupId:            1,
// 		OperationSectionDescription: "kochai bau kambing",
// 	}

// 	get, err := operationSectionServ.SaveOperationSection(req)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(get)
// }
