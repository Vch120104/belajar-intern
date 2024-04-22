package test

import (
	"after-sales/api/config"
	"after-sales/api/payloads/pagination"
	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
	"after-sales/api/utils"
	"fmt"
	"testing"
)

// import (
// 	"after-sales/api/config"
// 	masteroperationpayloads "after-sales/api/payloads/master/operation"
// 	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
// 	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
// 	"fmt"
// 	"testing"
// )

// func TestChangeStatusOperationGroup(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	operationGroupRepo := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl(db)
// 	operationGroupServ := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepo)

// 	get, err := operationGroupServ.ChangeStatusOperationGroup(1)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(get)
// }

// func TestGetOperationGroupById(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	operationGroupRepo := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl(db)
// 	operationGroupServ := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepo)

// 	get, err := operationGroupServ.GetOperationGroupById(2)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(get)
// }

// // func TestGetAllOperationGroup(t *testing.T) {
// // 	config.InitEnvConfigs(true, "")
// // 	db := config.InitDB()
// // 	operationGroupRepo := operationgrouprepoimpl.StartOperationGroupRepositoryImpl(db)
// // 	operationGroupServ := operationgroupserviceimpl.StartOperationGroupService(operationGroupRepo)

// // 	queryOf := []string{} // Replace with actual query values
// // 	queryBy := []string{} // Replace with actual query fields
// // 	sortOf := "operation_group_code"
// // 	sortBy := "asc"

// // 	get, err := operationGroupServ.GetAllOperationGroup(queryOf, queryBy, sortOf, sortBy)

// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	fmt.Println(get)
// // }

// func TestSaveOperationGroup(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	operationGroupRepo := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl(db)
// 	operationGroupServ := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepo)
// 	//if OperationGroupId = 0 -> insert; else update
// 	req := masteroperationpayloads.OperationGroupResponse{
// 		IsActive:                  true,
// 		OperationGroupId:          0,
// 		OperationGroupCode:        "A1",
// 		OperationGroupDescription: "test",
// 	}

// 	get, err := operationGroupServ.SaveOperationGroup(req)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(get)
// }

// func TestUpdateGroupById(t *testing.T) {
// 	config.InitEnvConfigs(true, "")
// 	db := config.InitDB()
// 	operationGroupRepo := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl(db)
// 	operationGroupServ := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepo)
// 	req := masteroperationpayloads.OperationGroupResponse{
// 		IsActive: false,
// 	}

// 	get, err := operationGroupServ.SaveOperationGroup(req)

// 	if err != nil {
// 		panic(err)
// 	}

//		fmt.Println(get)
//	}

func TestGetAllOperationGroup(t *testing.T) {
	config.InitEnvConfigs(true, "")
	// Initialize Redis client
	rdb := config.InitRedis()
	db := config.InitDB()

	filterCondition := []utils.FilterCondition{
		{
			ColumnField: "operation_group_description",
			ColumnValue: "",
		},
	}

	pages := pagination.Pagination{
		Page:  0,
		Limit: 10,
	}

	operationGroupRepo := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl()
	operationGroupServ := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepo, db, rdb)

	res, err := operationGroupServ.GetAllOperationGroup(filterCondition, pages)

	if err != nil {
		fmt.Print("err ", err)
	}

	fmt.Print("result ", res)

	// entities := masteroperationentities.OperationGroup{}

	// filterCondition := []utils.FilterCondition{
	// 	{
	// 		ColumnField: "operation_group_descriptiona",
	// 		ColumnValue: "",
	// 	},
	// }

	// pages := pagination.Pagination{
	// 	Page:  0,
	// 	Limit: 10,
	// }

	// baseModelQuery := db.Model(&entities)
	// //apply where query
	// whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)
	// //apply pagination and execute
	// _, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()

	// if err != nil {
	// 	fmt.Print("err ", err)
	// }

	// fmt.Print("result ", entities)

}
