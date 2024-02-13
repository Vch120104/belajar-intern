package main

import (
	"after-sales/api/config"
	mastercontroller "after-sales/api/controllers/master"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	"after-sales/api/helper"

	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	"after-sales/api/route"

	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	migration "after-sales/generate/sql"
	"net/http"
	"os"
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
		// redis := config.InitRedis()
		// route.CreateHandler(db, env, redis)

		deductionRepository := masterrepositoryimpl.StartDeductionRepositoryImpl()
		deductionService := masterserviceimpl.StartDeductionService(deductionRepository, db)
		deductionController := mastercontroller.NewDeductionController(deductionService)

		operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl()
		operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository, db)
		operationGroupController := masteroperationcontroller.NewOperationGroupController(operationGroupService)

		forecastMasterRepository := masterrepositoryimpl.StartForecastMasterRepositoryImpl()
		forecastMasterService := masterserviceimpl.StartForecastMasterService(forecastMasterRepository, db)
		forecastMasterController := mastercontroller.NewForecastMasterController(forecastMasterService)

		OperationGroupRouter := route.OperationGroupRouter(operationGroupController)
		ForecastMasterRouter := route.ForecastMasterRouter(forecastMasterController)
		DeductionRouter := route.DeductionRouter(deductionController)

		swaggerRouter := route.SwaggerRouter()
		mux := http.NewServeMux()

		mux.Handle("/operation-group/", OperationGroupRouter)
		mux.Handle("/forecast-master/", ForecastMasterRouter)
		mux.Handle("/deduction/", DeductionRouter)

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
