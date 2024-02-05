package main

import (
	"after-sales/api/config"
	mastercontroller "after-sales/api/controllers/master"
	"after-sales/api/helper"

	// masteritemrepositoryimpl "after-sales/api/repositories/master/item/repositories-item-impl"

	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	"after-sales/api/route"

	// masteritemserviceimpl "after-sales/api/services/master/item/services-item-impl"

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
		config.InitLogger(db)
		// basePath := "/api/aftersales/discount-percent"
		// redis := config.InitRedis()
		// route.CreateHandler(db, env, redis)

		// operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl()
		// operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository, db)
		// operationGroupController := masteroperationcontroller.NewOperationGroupController(operationGroupService)

		forecastMasterRepository := masterrepositoryimpl.StartForecastMasterRepositoryImpl()
		forecastMasterService := masterserviceimpl.StartForecastMasterService(forecastMasterRepository, db)
		forecastMasterController := mastercontroller.NewForecastMasterController(forecastMasterService)

		DeductionRepository := masterrepositoryimpl.StartDeductionRepositoryImpl()
		DeductionService := masterserviceimpl.StartDeductionService(DeductionRepository, db)
		DeductionController := mastercontroller.NewDeductionController(DeductionService)

		// OperationGroupRouter := route.OperationGroupRouter(operationGroupController )
		ForecastMasterRouter := route.ForecastMasterRouter(forecastMasterController)
		DeductionRouter := route.DeductionRouter(DeductionController)

		swaggerRouter := route.SwaggerRouter()
		mux := http.NewServeMux()

		// mux.Handle("/operation-group/",OperationGroupRouter)
		mux.Handle("/forecast-master/", ForecastMasterRouter)
		mux.Handle("/deduction-master/", DeductionRouter)

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
