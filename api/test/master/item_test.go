package test

// import (
// 	"after-sales/api/config"
// 	operationitemrepoimpl "after-sales/api/repositories/master/item/item-impl"
// 	operationitemserviceimpl "after-sales/api/services/master/item/item-impl"
// 	"fmt"
// 	"testing"
// )

// func TestGetOperationItemById(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	operationItemRepo := operationitemrepoimpl.StartItemImpl(db)
// 	operationItemServ := operationitemserviceimpl.StartItemService(operationItemRepo)
// 	get, err := operationItemServ.GetItemById(2)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(get)
// }
