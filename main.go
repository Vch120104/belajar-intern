package main

import (
	"after-sales/api/config"
	route "after-sales/api/route"
	migration "after-sales/generate/sql"
	"context"
	"fmt"
	"log"
	"os"
)

// test
//
//	@title			DMS After-Sales API
//	@version		v1
//	@license		AGPLv3
//	@description	This is a DMS After-Sales API Server.
//	 @basePath		/aftersales-service
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
		dbredis := config.InitRedis()

		// Set a key-value pair
		ctx := context.Background()
		err := dbredis.Set(ctx, "key", "value", 0).Err()
		if err != nil {
			log.Fatalf("could not set key: %v", err)
		}
		fmt.Println("Key set successfully")

		config.InitLogger(db)
		route.StartRouting(db, dbredis)
	}
}
