package main

import (
	"after-sales/api/config"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	"after-sales/api/helper"
	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	"after-sales/api/route"
	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
	migration "after-sales/generate/sql"
	"net/http"
	"os"
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

		operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl()
		operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository, db)
		operationGroupController := masteroperationcontroller.NewOperationGroupController(operationGroupService)
		OperationGroupRouter := route.OperationGroupRouter(operationGroupController)

		swaggerRouter := route.SwaggerRouter()
		mux := http.NewServeMux()

		mux.Handle("/test/", OperationGroupRouter)

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
