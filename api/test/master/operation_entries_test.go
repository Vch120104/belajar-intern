package test

// import (
// 	"after-sales/api/config"
// 	operationentriesrepoimpl "after-sales/api/repositories/master/operation-entries/operation-entries-impl"
// 	operationentriesserviceimpl "after-sales/api/services/master/operation-entries/operation-entries-impl"
// 	"fmt"
// 	"testing"
// )

// func TestGetOperationEntriesById(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	operationEntriesRepo := operationentriesrepoimpl.StartOperationEntriesImpl(db)
// 	operationEntriesServ := operationentriesserviceimpl.StartOperationEntriesService(operationEntriesRepo)
// 	get, err := operationEntriesServ.GetOperationEntriesById(1)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(get)
// }
