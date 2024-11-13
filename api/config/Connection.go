package config

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitDB() *gorm.DB {
	var db *gorm.DB = nil
	val := url.Values{}
	val.Add("parseTime", "True")
	val.Add("loc", "Asia/Jakarta")

	var err error

	if strings.Contains(EnvConfigs.DBDriver, "postgre") {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", EnvConfigs.DBHost, EnvConfigs.DBUser, EnvConfigs.DBPass, EnvConfigs.DBName, EnvConfigs.DBPort)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		dsn := fmt.Sprintf(`%s://%s:%s@%s:%v?database=%s`, EnvConfigs.DBDriver, EnvConfigs.DBUser, EnvConfigs.DBPass, EnvConfigs.DBHost, EnvConfigs.DBPort, EnvConfigs.DBName) //SQLSEVER
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				//TablePrefix:   "dbo.", // schema name
				SingularTable: false,
			}})
	}

	if err != nil {
		log.Fatal("Cannot connected database ", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error connecting to database ", err)
		return nil
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Request Timeout ", err)
		return nil
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(time.Minute * 3)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 3)

	log.Info("Connected Database " + EnvConfigs.DBDriver + " -- running in -- " + EnvConfigs.ClientOrigin)

	return db
}

func InitRedis() *redis.Client {
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     EnvConfigs.ClientRedis + ":" + EnvConfigs.PortRedis,
	//	Password: "OOg6hZ7KvrU4aAIhmhq2cNfhgUjMYlif",
	//	DB:       0,
	//})
	const (
		REDIS_HOST     = "redis-16833.c334.asia-southeast2-1.gce.redns.redis-cloud.com"
		REDIS_PASSWORD = "bSJgYT9L4UiMJr4CCSismkxcHZAtnOci"
		REDIS_PORT     = 16833
		REDIS_DB       = 0 // Replace "Devin-free-db" with the database number if needed; typically, DB is an integer.
	)
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", REDIS_HOST, REDIS_PORT), // Combine host and port
		Password: REDIS_PASSWORD,                               // Password if set
		DB:       REDIS_DB,                                     // Use the default DB
	})
	// Menguji koneksi Redis
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Info("Connected Redis -- running in -- " + REDIS_HOST + ":" + strconv.Itoa(REDIS_PORT))

	return rdb
}
