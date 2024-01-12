package test

// import (
// 	"after-sales/api/config"
// 	operationcoderepoimpl "after-sales/api/repositories/master/operation-code/operation-code-impl"
// 	operationcodeserviceimpl "after-sales/api/services/master/operation-code/operation-code-impl"
// 	"fmt"
// 	"testing"
// )

// func TestGetOperationCodeById(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	operationCodeRepo := operationcoderepoimpl.StartOperationCodeImpl(db)
// 	operationCodeServ := operationcodeserviceimpl.StartOperationCodeService(operationCodeRepo)
// 	get, err := operationCodeServ.GetOperationCodeById(1)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(get)
// }
