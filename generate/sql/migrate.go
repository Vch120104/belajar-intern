package migration

import (
	"after-sales/api/config"

	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masteroperationentities "after-sales/api/entities/master/operation"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"

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
	logEntry := "Auto Migrating..."

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
		&masteroperationentities.OperationModelMapping{},
		&masteroperationentities.OperationFrt{},
		&masteroperationentities.OperationCode{},
		&masteroperationentities.OperationGroup{},
		&masteroperationentities.OperationSection{},
		&masteroperationentities.OperationKey{},
		&masteroperationentities.OperationEntries{},

		&masterwarehouseentities.WarehouseGroup{},
		&masterwarehouseentities.WarehouseMaster{},
		&masterwarehouseentities.WarehouseLocation{},
		&masterwarehouseentities.WarehouseLocationDefinition{},
		&masterwarehouseentities.WarehouseLocationDefinitionLevel{},

		&masteritementities.MarkupMaster{},
		&masteritementities.PrincipleBrandParent{},
		&masteritementities.UomType{},
		&masteritementities.PriceList{},
		&masteritementities.ItemLocationSource{},
		&masteritementities.ItemLocationDetail{},
		&masteritementities.ItemLocation{},
		&masteritementities.PurchasePrice{},
		&masteritementities.PurchasePriceDetail{},
		&masteritementities.ItemDetail{},
		&masteritementities.ItemImport{},
		&masteritementities.DiscountPercent{},
		&masteritementities.ItemSubstitute{},
		&masteritementities.ItemSubstituteDetail{},
		&masteritementities.ItemPackage{},
		&masteritementities.ItemPackageDetail{},
		&masteritementities.PrincipleBrandParent{},
		&masteritementities.UomType{},
		&masteritementities.Uom{},
		&masteritementities.Bom{},
		&masteritementities.BomDetail{},
		&masteritementities.ItemClass{},
		&masteritementities.Item{},
		&masteritementities.MarkupRate{},
		&masteritementities.ItemLevel{},

		&masterentities.IncentiveGroup{},
		&masterentities.ForecastMaster{},
		&masterentities.MovingCode{},
		&masterentities.ShiftSchedule{},
		&masterentities.IncentiveGroup{},
		&masterentities.ForecastMaster{},
		&masterentities.MovingCode{},
		&masterentities.ShiftSchedule{},
		&masterentities.IncentiveMaster{},
		&masterentities.IncentiveGroupDetail{},
		&masterentities.SkillLevel{},
		&masterentities.WarrantyFreeService{},
		&masterentities.DeductionList{},
		&masterentities.DeductionDetail{},
		&masterentities.FieldActionEligibleVehicleItem{},
		&masterentities.FieldActionEligibleVehicle{},
		&masterentities.FieldAction{},

		&masterentities.Discount{},
	// &transactionentities.SupplySlip{},
	// &transactionentities.SupplySlipDetail{},
	// &transactionentities.WorkOrderItem{},
	// &transactionentities.WorkOrderOperation{},
	// &transactionentities.ServiceLog{},
	// &transactionworkshopentities.BookingEstimation{},
	)

	if db != nil && db.Error != nil {
		log.Printf("%s Failed with error %s", logEntry, db.Error)
		panic(err)
	}

	log.Printf("%s Success", logEntry)
}
