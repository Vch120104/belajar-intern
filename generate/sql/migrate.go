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

	// Initialize logger
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// Disable foreign key constraints when migrating
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: newLogger, // Set the logger for GORM
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "dbo.", // schema name
			SingularTable: false,
		},
		DisableForeignKeyConstraintWhenMigrating: false,
	})

	if err != nil {
		log.Printf("%s Failed to connect to database with error: %s", logEntry, err)
		panic(err)
	}

	// AutoMigrate models
	err = db.AutoMigrate( // according to foreign key order
	//&masteroperationentities.OperationModelMapping{},
	//&masteroperationentities.OperationFrt{},
	//&masteroperationentities.OperationGroup{},
	//&masteroperationentities.OperationSection{},
	//&masteroperationentities.OperationKey{},
	//&masteroperationentities.OperationEntries{},
	//&masteroperationentities.OperationCode{},
	//
	//&masterwarehouseentities.WarehouseGroup{},
	//&masterwarehouseentities.WarehouseMaster{},
	//&masterwarehouseentities.WarehouseLocation{},
	//&masterwarehouseentities.WarehouseLocationDefinition{},
	//&masterwarehouseentities.WarehouseLocationDefinitionLevel{},
	//
	//&masteritementities.LandedCost{},
	//&masteritementities.PurchasePrice{},
	//&masteritementities.PurchasePriceDetail{},
	//&masteritementities.UomType{},
	//&masteritementities.Uom{},
	//&masteritementities.UomItem{},
	//&masteritementities.MarkupRate{},
	//&masteritementities.PrincipleBrandParent{},
	//&masteritementities.MarkupMaster{},
	//&masteritementities.ItemLevel{},
	//&masteritementities.PriceList{},
	//&masteritementities.ItemSubstituteDetail{},
	//&masteritementities.ItemSubstitute{},
	//&masteritementities.ItemPackage{},
	//&masteritementities.ItemPackageDetail{},
	//&masteritementities.ItemImport{},
	//&masteritementities.ItemDetail{},
	//&masteritementities.ItemLocationSource{},
	//&masteritementities.Item{},
	//&masteritementities.ItemLocation{},
	//&masteritementities.ItemLocationDetail{},
	//&masteritementities.ItemClass{},
	//&masteritementities.Bom{},
	//&masteritementities.BomDetail{},
	//&masteritementities.Discount{},
	//&masteritementities.DiscountPercent{},
	//
	//&masterentities.ForecastMaster{},
	//&masterentities.MovingCode{},
	//&masterentities.IncentiveGroup{},
	//&masterentities.PackageMaster{},
	//&masterentities.ShiftSchedule{},
	//&masterentities.IncentiveMaster{},
	//&masterentities.IncentiveGroupDetail{},
	//&masterentities.SkillLevel{},
	//&masterentities.WarrantyFreeService{},
	//&masterentities.DeductionList{},
	//&masterentities.DeductionDetail{},
	//&masterentities.FieldAction{},
	//&masterentities.FieldActionEligibleVehicleItem{},
	//&masterentities.FieldActionEligibleVehicleOperation{},
	//&masterentities.FieldActionEligibleVehicle{},
	//&masteritementities.DiscountPercent{},
	//&masterentities.Agreement{},
	//&masterentities.AgreementDiscount{},
	//&masterentities.AgreementDiscountGroupDetail{},
	//&masterentities.AgreementItemDetail{},
	//
	//&masterentities.FieldAction{},
	//&masterentities.FieldActionEligibleVehicleItem{},
	//&masterentities.FieldActionEligibleVehicle{},
	//&masterentities.Agreement{},
	//&masterentities.AgreementDiscount{},
	//&masterentities.AgreementDiscountGroupDetail{},
	//&masterentities.AgreementItemDetail{},
	//
	//&mastercampaignmasterentities.CampaignMaster{},
	//&mastercampaignmasterentities.CampaignMasterDetailItem{},
	//&mastercampaignmasterentities.CampaignMasterOperationDetail{},
	//
	//&transactionsparepartentities.SupplySlip{},
	//&transactionsparepartentities.SupplySlipDetail{},
	//
	//&transactionworkshopentities.WorkOrderMaster{},
	//&transactionworkshopentities.WorkOrderMasterStatus{},
	//&transactionworkshopentities.WorkOrderMasterType{},
	//&transactionworkshopentities.WorkOrderMasterBillAbleto{},
	//&transactionworkshopentities.WorkOrder{},
	//&transactionworkshopentities.WorkOrderRequestDescription{},
	//&transactionworkshopentities.WorkOrderDetail{},
	//&transactionworkshopentities.WorkOrderHistory{},
	//&transactionworkshopentities.WorkOrderHistoryRequest{},
	//&transactionworkshopentities.WorkOrderHistoryDetail{},
	//&transactionworkshopentities.WorkOrderService{},
	//&transactionworkshopentities.WorkOrderServiceVehicle{},
	//&transactionworkshopentities.ServiceRequest{},
	//&transactionworkshopentities.ServiceRequestDetail{},
	//&transactionworkshopentities.ServiceRequestMasterStatus{},
	//&transactionworkshopentities.BookingEstimation{},
	//&transactionworkshopentities.BookingEstimationAllocation{},
	//&transactionworkshopentities.BookingEstimationRequest{},
	//&transactionworkshopentities.BookingEstimationServiceReminder{},
	//&transactionworkshopentities.BookingEstimationServiceDiscount{},
	//&transactionworkshopentities.BookingEstimationItemDetail{},
	//&transactionworkshopentities.BookingEstimationOperationDetail{},
	//&transactionworkshopentities.BookingEstimationRequest{},
	//
	//&transactionworkshopentities.BookingEstimation{},
	//&transactionworkshopentities.BookingEstimationAllocation{},
	//&transactionworkshopentities.BookingEstimationRequest{},
	//&transactionworkshopentities.BookingEstimationServiceReminder{},
	//&transactionworkshopentities.BookingEstimationServiceDiscount{},
	//&transactionworkshopentities.BookingEstimationDetail{},

	// &transactionsparepartentities.PurchaseRequestEntities{},
	// &transactionsparepartentities.PurchaseRequestDetail{},

	// &masterentities.LocationStock{},
	//&transactionsparepartentities.PurchaseRequestReferenceType{},
	)

	if err != nil {
		log.Printf("%s Failed with error: %s", logEntry, err)
		panic(err)
	}

	log.Printf("%s Success", logEntry)
}
