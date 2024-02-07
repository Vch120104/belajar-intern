package main

import (
	"after-sales/api/config"
	masteritemcontroller "after-sales/api/controllers/master/item"

	// mastercontroller "after-sales/api/controllers/master"
	mastercontroller "after-sales/api/controllers/master"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	"after-sales/api/helper"

	masteritemrepositoryimpl "after-sales/api/repositories/master/item/repositories-item-impl"

	// masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"

	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	"after-sales/api/route"
	masteritemserviceimpl "after-sales/api/services/master/item/services-item-impl"

	// masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"

	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	migration "after-sales/generate/sql"
	"net/http"
	"os"

	masterserviceimpl "after-sales/api/services/master/service-impl"
)

// @title After Sales API
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
// @host localhost:8000
// @BasePath /api/aftersales
func main() {
	args := os.Args
	env := ""
	if len(args) > 1 {
		env = args[1]
	}

	if env == "migrate" {
		migration.Migrate()
	} else if env == "generate" {
		migration.Generate()
	} else if env == "migg" {
		migration.MigrateGG()
	} else if env == "debug" {
		// config.InitEnvConfigs(false, env)
		// db := config.InitDB()
		// config.InitLogger(db)
		// redis := config.InitRedis()
		// route.CreateHandler(db, env, redis)
	} else {
		config.InitEnvConfigs(false, env)
		db := config.InitDB()
		// basePath := "/api/aftersales"
		// redis := config.InitRedis()
		// route.CreateHandler(db, env, redis)

		// operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl()
		// operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository, db)
		// operationGroupController := masteroperationcontroller.NewOperationGroupController(operationGroupService)
		// OperationGroupRouter := route.OperationGroupRouter(operationGroupController, basePath)

		startFieldActionRepository := masterrepositoryimpl.StartFieldActionRepositoryImpl()
		startFieldActionService := masterserviceimpl.StartFieldActionService(startFieldActionRepository, db)
		startFieldActionController := mastercontroller.NewFieldActionController(startFieldActionService)
		startFieldActionRouter := route.FieldActionRouter(startFieldActionController)
		// operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl()
		// operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository, db)
		// operationGroupController := masteroperationcontroller.NewOperationGroupController(operationGroupService)

		discountPercentRepository := masteritemrepositoryimpl.StartDiscountPercentRepositoryImpl()
		discountPercentService := masteritemserviceimpl.StartDiscountPercentService(discountPercentRepository, db)
		discountPercentController := masteritemcontroller.NewDiscountPercentController(discountPercentService)
		operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl()
		operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository, db)
		operationGroupController := masteroperationcontroller.NewOperationGroupController(operationGroupService)

		forecastMasterRepository := masterrepositoryimpl.StartForecastMasterRepositoryImpl()
		forecastMasterService := masterserviceimpl.StartForecastMasterService(forecastMasterRepository, db)
		forecastMasterController := mastercontroller.NewForecastMasterController(forecastMasterService)

		ShiftScheduleRepository := masterrepositoryimpl.StartShiftScheduleRepositoryImpl()
		ShiftScheduleService := masterserviceimpl.StartShiftScheduleService(ShiftScheduleRepository, db)
		ShiftScheduleController := mastercontroller.NewShiftScheduleController(ShiftScheduleService)

		// OperationGroupRouter := route.OperationGroupRouter(operationGroupController)
		DiscountPercentRouter := route.DiscountPercentRouter(discountPercentController)
		OperationGroupRouter := route.OperationGroupRouter(operationGroupController)
		ForecastMasterRouter := route.ForecastMasterRouter(forecastMasterController)
		ShiftScheduleRouter := route.ShiftScheduleRouter(ShiftScheduleController)

		swaggerRouter := route.SwaggerRouter()
		mux := http.NewServeMux()

		// // mux.Handle("/api/aftersales/operation-group", OperationGroupRouter)
		mux.Handle("/field-action/", startFieldActionRouter)
		mux.Handle("/discount-percent/", DiscountPercentRouter)
		mux.Handle("/operation-group/", OperationGroupRouter)
		mux.Handle("/forecast-master/", ForecastMasterRouter)
		mux.Handle("/shift-schedule/", ShiftScheduleRouter)

		//Swagger
		mux.Handle("/swagger/", swaggerRouter)
		server := http.Server{
			Addr:    config.EnvConfigs.ClientOrigin,
			Handler: mux,
		}

		err := server.ListenAndServe()
		helper.PanicIfError(err)
	}
}
