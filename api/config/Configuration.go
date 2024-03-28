package config

import (
	"os"
	"time"

	logging "log"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Configuration struct {
	Database setupDatabase
}

type setupDatabase struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPass            string `mapstructure:"DB_PASS"`
	DBName            string `mapstructure:"DB_NAME"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            int    `mapstructure:"DB_PORT"`
	Port              string `mapstructure:"SERVER_PORT"`
	JWTKey            string `mapstructure:"JWT_KEY"`
	ClientOrigin      string `mapstructure:"CLIENT_ORIGIN"`
	ClientRedis       string `mapstructure:"CLIENT_REDIS"`
	GeneralServiceUrl string `mapstructure:"GENERAL_SERVICE"`
	SalesServiceUrl   string `mapstructure:"SALES_SERVICE"`
	ReportServiceUrl  string `mapstructure:"REPORT_SERVICE"`
	FinanceServiceUrl string `mapstructure:"FINANCE_SERVICE"`
	CommonServiceUrl  string `mapstructure:"COMMON_SERVICE"`
}

var EnvConfigs *setupDatabase

func InitEnvConfigs(gen bool, env string) {
	EnvConfigs = SetupConfiguration(gen, env)
}

func SetupConfiguration(gen bool, env string) (config *setupDatabase) {
	// Tell viper the path/location of your env file. If it is root just add "."
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
		if gen {
			viper.AddConfigPath("../../.production")
		} else {
			viper.AddConfigPath(".production")
		}
	} else {
		if gen {
			viper.AddConfigPath("../../../.development")
		} else {
			viper.AddConfigPath("./.development")
		}
	}

	// Tell viper the name of your file
	viper.SetConfigName("app")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	// Viper unmarshals the loaded env varialbes into the struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	return
}

func InitLogger(db *gorm.DB) {
	newLogger := logger.New(
		logging.New(logging.Writer(), "\r\n", logging.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db.Logger = newLogger
}
