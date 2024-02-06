package main

import (
	"after-sales/api/config"
	mastercontroller "after-sales/api/controllers/master"
	"after-sales/api/helper"

	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	"after-sales/api/route"
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
// @host localhost:2000
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
		// OperationGroupRouter := route.OperationGroupRouter(operationGroupController, basePath)

		startFieldActionRepository := masterrepositoryimpl.StartFieldActionRepositoryImpl()
		startFieldActionService := masterserviceimpl.StartFieldActionService(startFieldActionRepository, db)
		startFieldActionController := mastercontroller.NewFieldActionController(startFieldActionService)
		startFieldActionRouter := route.FieldActionRouter(startFieldActionController)

		swaggerRouter := route.SwaggerRouter()
		mux := http.NewServeMux()

		// mux.Handle(basePath, OperationGroupRouter)
		mux.Handle("/field-action/", startFieldActionRouter)

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
