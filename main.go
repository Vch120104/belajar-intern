package main

import (
	"after-sales/api/config"
	"after-sales/api/route"
	migration "after-sales/generate/sql"
	"os"
)

// @title After Sales API
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
// @host 10.1.32.26:2000
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
		config.InitEnvConfigs(false, env)
		db := config.InitDB()
		config.InitLogger(db)
		redis := config.InitRedis()
		route.CreateHandler(db, env, redis)
	} else {
		config.InitEnvConfigs(false, env)
		db := config.InitDB()
		redis := config.InitRedis()
		route.CreateHandler(db, env, redis)
	}
}
