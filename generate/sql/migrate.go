package migration

import (
	"after-sales/api/config"

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

	db.AutoMigrate( // sesuai urutan foreign key
	// &masteroperationentities.OperationModelMapping{},
	// &masteroperationentities.OperationFrt{},
	// &masteroperationentities.OperationGroup{},
	// &masteroperationentities.OperationSection{},
	// &masteroperationentities.OperationKey{},
	// &masteroperationentities.OperationEntries{},
	// &masteroperationentities.OperationCode{},

	// &masterwarehouseentities.WarehouseGroup{},
	// &masterwarehouseentities.WarehouseMaster{},
	// &masterwarehouseentities.WarehouseLocation{},
	// &masterwarehouseentities.WarehouseLocationDefinition{},
	// &masterwarehouseentities.WarehouseLocationDefinitionLevel{},

	// &masteritementities.ItemLocation{},
	// &masteritementities.ItemLocationSource{},
	// &masteritementities.ItemLocationDetail{},
	// &masteritementities.PurchasePrice{},
	// &masteritementities.PurchasePriceDetail{},
	// &masteritementities.UomType{},
	// &masteritementities.Uom{},
	// &masteritementities.Bom{},
	// &masteritementities.BomDetail{},
	// &masteritementities.MarkupRate{},
	// &masteritementities.PrincipleBrandParent{},
	// &masteritementities.DiscountPercent{},
	// &masteritementities.MarkupMaster{},
	// &masteritementities.ItemLevel{},
	// &masteritementities.ItemClass{},
	// &masteritementities.PriceList{},
	// &masteritementities.ItemSubstituteDetail{},
	// &masteritementities.ItemSubstitute{},
	// &masteritementities.ItemPackage{},
	// &masteritementities.ItemPackageDetail{},
	// &masteritementities.ItemDetail{},
	// &masteritementities.ItemImport{},
	// &masteritementities.Item{},

	// &masterentities.ForecastMaster{},
	// &masterentities.MovingCode{},
	// &masterentities.IncentiveGroup{},
	// &masterentities.PackageMaster{},
	// &masterentities.ShiftSchedule{},
	// &masterentities.IncentiveMaster{},
	// &masterentities.IncentiveGroupDetail{},
	// &masterentities.SkillLevel{},
	// &masterentities.WarrantyFreeService{},
	// &masterentities.DeductionList{},
	// &masterentities.DeductionDetail{},
	// &masterentities.FieldActionEligibleVehicleItem{},
	// &masterentities.FieldActionEligibleVehicle{},
	// &masterentities.FieldAction{},
	// &masterentities.Discount{},
	// &masterentities.Agreement{},
	// &masterentities.AgreementDiscount{},
	// &masterentities.AgreementDiscountGroupDetail{},
	// &masterentities.AgreementItemDetail{},

	// &transactionsparepartentities.SupplySlip{},
	// &transactionsparepartentities.SupplySlipDetail{},
	// &transactionworkshopentities.WorkOrderMaster{},
	// &transactionworkshopentities.WorkOrderMasterStatus{},
	// &transactionworkshopentities.WorkOrderMasterType{},
	// &transactionworkshopentities.WorkOrder{},
	// &transactionworkshopentities.WorkOrderRequestDescription{},
	// &transactionworkshopentities.WorkOrderDetail{},
	// &transactionworkshopentities.WorkOrderHistory{},
	// &transactionworkshopentities.WorkOrderHistoryRequest{},
	// &transactionworkshopentities.WorkOrderHistoryDetail{},

	// &transactionworkshopentities.BookingEstimation{},
	// &transactionworkshopentities.BookingEstimationAllocation{},
	// &transactionworkshopentities.BookingEstimationRequest{},
	// &transactionworkshopentities.BookingEstimationServiceReminder{},
	// &transactionworkshopentities.BookingEstimationServiceDiscount{},
	// &transactionworkshopentities.BookingEstimationDetail{},
	)

	if db != nil && db.Error != nil {
		log.Printf("%s Failed with error %s", logEntry, db.Error)
		panic(err)
	}

	log.Printf("%s Success", logEntry)
}
