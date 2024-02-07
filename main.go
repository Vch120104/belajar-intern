package main

import (
	"after-sales/api/config"
	masteritemcontroller "after-sales/api/controllers/master/item"

	// masteroperationcontroller "after-sales/api/controllers/master/operation"
	mastercontroller "after-sales/api/controllers/master"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	"after-sales/api/helper"
	masteritemrepositoryimpl "after-sales/api/repositories/master/item/repositories-item-impl"

	// masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"

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

		// operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl()
		// operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository, db)
		// operationGroupController := masteroperationcontroller.NewOperationGroupController(operationGroupService)

		forecastMasterRepository := masterrepositoryimpl.StartForecastMasterRepositoryImpl()
		forecastMasterService := masterserviceimpl.StartForecastMasterService(forecastMasterRepository, db)
		forecastMasterController := mastercontroller.NewForecastMasterController(forecastMasterService)
		discountPercentRepository := masteritemrepositoryimpl.StartDiscountPercentRepositoryImpl()
		discountPercentService := masteritemserviceimpl.StartDiscountPercentService(discountPercentRepository, db)
		discountPercentController := masteritemcontroller.NewDiscountPercentController(discountPercentService)

		OperationGroupRouter := route.OperationGroupRouter(operationGroupController)
		ForecastMasterRouter := route.ForecastMasterRouter(forecastMasterController)
		markupRateRepository := masteritemrepositoryimpl.StartMarkupRateRepositoryImpl()
		markupRateService := masteritemserviceimpl.StartMarkupRateService(markupRateRepository, db)
		markupRateController := masteritemcontroller.NewMarkupRateController(markupRateService)

		// OperationGroupRouter := route.OperationGroupRouter(operationGroupController)
		DiscountPercentRouter := route.DiscountPercentRouter(discountPercentController)
		MarkupRateRouter := route.MarkupRateRouter(markupRateController)

		swaggerRouter := route.SwaggerRouter()
		mux := http.NewServeMux()

		mux.Handle("/operation-group/", OperationGroupRouter)
		mux.Handle("/forecast-master/", ForecastMasterRouter)
		// mux.Handle("/api/aftersales/operation-group", OperationGroupRouter)
		mux.Handle("/discount-percent/", DiscountPercentRouter)

		mux.Handle("/markup-rate/", MarkupRateRouter)

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
