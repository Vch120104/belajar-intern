package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var rdb *redis.Client

func InitDB() *gorm.DB {
	var err error = nil
	var db *gorm.DB = nil
	val := url.Values{}
	val.Add("parseTime", "True")
	val.Add("loc", "Asia/Jakarta")

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
	rdb = redis.NewClient(&redis.Options{
		Addr:     EnvConfigs.ClientRedis,
		Password: "",
		DB:       0,
	})
	return rdb
}
