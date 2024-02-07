package migration

import (
	"after-sales/api/config"

	// masteroperationentities "after-sales/api/entities/master/operation"

	masterentities "after-sales/api/entities/master"
	// masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"time"

	"fmt"
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Migrate() {
	config.InitEnvConfigs(false, "")
	logEntry := fmt.Sprintf("Auto Migrating...")

	dsn := fmt.Sprintf(
		`%s://%s:%s@%s:%v?database=%s`,
		config.EnvConfigs.DBDriver,
		config.EnvConfigs.DBUser,
		config.EnvConfigs.DBPass,
		config.EnvConfigs.DBHost,
		config.EnvConfigs.DBPort,
		config.EnvConfigs.DBName,
	)

	//init logger
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	//constraint foreign key tidak akan ke create jika DisableForeignKeyConstraintWhenMigrating: true
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: newLogger, // Set the logger for GORM
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "dbo.", // schema name
			SingularTable: false,
		},
		DisableForeignKeyConstraintWhenMigrating: false,
	})

	// db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
	// 	NamingStrategy: schema.NamingStrategy{
	// 		//TablePrefix:   "dbo.", // schema name
	// 		SingularTable: false,
	// 	}, DisableForeignKeyConstraintWhenMigrating: false})

	db.AutoMigrate( // sesuai urutan foreign key
		// &masteroperationentities.OperationGroup{},
		// &masteroperationentities.OperationSection{},
		// &masteroperationentities.OperationKey{},
		// &masteroperationentities.OperationEntries{},
		// &masteroperationentities.OperationCode{},
		// &masteroperationentities.OperationModelMapping{},
		// &masteritementities.ItemClass{},
		// &masteritementities.Discount{},
		// &masteritementities.MarkupMaster{},
		// // &masteritementities.PrincipleBrandParent{},
		// // &masteritementities.UomType{},
		// // &masteritementities.Uom{},
		// // &masteritementities.Item{},
		// // &masteritementities.PriceList{},
		// // &masteritementities.ItemDetail{},
		// &masteritementities.DiscountPercent{},
		// &masterentities.IncentiveGroup{},
		&masterentities.ShiftSchedule{},

		// &masterwarehouseentities.WarehouseGroup{},
		&masterentities.FieldAction{},
		// &masterwarehouseentities.WarehouseMaster{},
		// &masterwarehouseentities.WarehouseLocation{},

	// &transactionentities.SupplySlip{},
	// &transactionentities.SupplySlipDetail{},
	// &transactionentities.WorkOrderItem{},
	// &transactionentities.WorkOrderOperation{},
	// &transactionentities.ServiceLog{},
	// &transactionworkshopentities.BookingEstimation{},
	)

	if db != nil && db.Error != nil {
		fmt.Sprintf("%s %s with error %s", logEntry, "Failed", db.Error)
		panic(err)
	}

	fmt.Sprintf("%s %s", logEntry, "Success")
}
