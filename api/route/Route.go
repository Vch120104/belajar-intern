package route

import (
	"after-sales/api/config"
	"after-sales/api/helper"
	masteritemrepositoryimpl "after-sales/api/repositories/master/item/repositories-item-impl"
	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	masterwarehouserepositoryimpl "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"
	masteritemserviceimpl "after-sales/api/services/master/item/services-item-impl"
	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	masterwarehouseserviceimpl "after-sales/api/services/master/warehouse/services-warehouse-impl"

	mastercontroller "after-sales/api/controllers/master"
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	masterwarehousecontroller "after-sales/api/controllers/master/warehouse"

	transactionsparepartcontroller "after-sales/api/controllers/transactions/sparepart"
	transactionworksopcontroller "after-sales/api/controllers/transactions/workshop"
	transactionsparepartrepositoryimpl "after-sales/api/repositories/transaction/sparepart/repositories-sparepart-impl"
	transactionworkshoprepositoryimpl "after-sales/api/repositories/transaction/workshop/repositories-workshop-impl"
	transactionsparepartserviceimpl "after-sales/api/services/transaction/sparepart/services-sparepart-impl"
	transactionworkshopserviceimpl "after-sales/api/services/transaction/workshop/services-workshop-impl"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

func StartRouting(db *gorm.DB) {
	// Initialize Redis client
	rdb := config.InitRedis()

	/* Master */
	// Unit Measurement
	unitOfMeasurementRepository := masteritemrepositoryimpl.StartUnitOfMeasurementRepositoryImpl()
	unitOfMeasurementService := masteritemserviceimpl.StartUnitOfMeasurementService(unitOfMeasurementRepository, db, rdb)
	unitOfMeasurementController := masteritemcontroller.NewUnitOfMeasurementController(unitOfMeasurementService)

	// Markup Master
	markupMasterRepository := masteritemrepositoryimpl.StartMarkupMasterRepositoryImpl()
	markupMasterService := masteritemserviceimpl.StartMarkupMasterService(markupMasterRepository, db, rdb)
	markupMasterController := masteritemcontroller.NewMarkupMasterController(markupMasterService)

	// Item Level
	itemLevelRepository := masteritemrepositoryimpl.StartItemLevelRepositoryImpl()
	itemLevelService := masteritemserviceimpl.StartItemLevelService(itemLevelRepository, db, rdb)
	itemLevelController := masteritemcontroller.NewItemLevelController(itemLevelService)

	// Item
	itemRepository := masteritemrepositoryimpl.StartItemRepositoryImpl()
	itemService := masteritemserviceimpl.StartItemService(itemRepository, db, rdb)
	itemController := masteritemcontroller.NewItemController(itemService)

	// Item Model Mapping
	ItemModelMappingRepository := masteritemrepositoryimpl.StartItemModelMappingRepositoryImpl()
	ItemModelMappingService := masteritemserviceimpl.StartItemModelMappingService(ItemModelMappingRepository, db)
	ItemModelMappingController := masteritemcontroller.NewItemModelMappingController(ItemModelMappingService)

	// PriceList
	priceListRepository := masteritemrepositoryimpl.StartPriceListRepositoryImpl()
	priceListService := masteritemserviceimpl.StartPriceListService(priceListRepository, db, rdb)
	priceListController := masteritemcontroller.NewPriceListController(priceListService)

	// Item Class
	itemClassRepository := masteritemrepositoryimpl.StartItemClassRepositoryImpl()
	itemClassService := masteritemserviceimpl.StartItemClassService(itemClassRepository, db, rdb)
	itemClassController := masteritemcontroller.NewItemClassController(itemClassService)

	// Item Location
	ItemLocationRepository := masteritemrepositoryimpl.StartItemLocationRepositoryImpl()
	ItemLocationService := masteritemserviceimpl.StartItemLocationService(ItemLocationRepository, db, rdb)
	ItemLocationController := masteritemcontroller.NewItemLocationController(ItemLocationService)

	// Item Substitute
	itemSubstituteRepository := masteritemrepositoryimpl.StartItemSubstituteRepositoryImpl()
	itemSubstituteService := masteritemserviceimpl.StartItemSubstituteService(itemSubstituteRepository, db, rdb)
	itemSubstituteController := masteritemcontroller.NewItemSubstituteController(itemSubstituteService)

	// Item Package
	itemPackageRepository := masteritemrepositoryimpl.StartItemPackageRepositoryImpl()
	itemPackageService := masteritemserviceimpl.StartItemPackageService(itemPackageRepository, db, rdb)
	itemPackageController := masteritemcontroller.NewItemPackageController(itemPackageService)

	// Item Package Detail
	itemPackageDetailRepository := masteritemrepositoryimpl.StartItemPackageDetailRepositoryImpl()
	itemPackageDetailService := masteritemserviceimpl.StartItemPackageDetailService(itemPackageDetailRepository, db, rdb)
	itemPackageDetailController := masteritemcontroller.NewItemPackageDetailController(itemPackageDetailService)

	// Item Import
	ItemImportRepository := masteritemrepositoryimpl.StartItemImportRepositoryImpl()
	ItemImportService := masteritemserviceimpl.StartItemImportService(ItemImportRepository, db)
	ItemImportController := masteritemcontroller.NewItemImportController(ItemImportService)

	// Purchase Price
	PurchasePriceRepository := masteritemrepositoryimpl.StartPurchasePriceRepositoryImpl()
	PurchasePriceService := masteritemserviceimpl.StartPurchasePriceService(PurchasePriceRepository, db, rdb)
	PurchasePriceController := masteritemcontroller.NewPurchasePriceController(PurchasePriceService)

	// // Landed Cost
	LandedCostRepository := masteritemrepositoryimpl.StartLandedCostMasterRepositoryImpl()
	LandedCostService := masteritemserviceimpl.StartLandedCostMasterService(LandedCostRepository, db, rdb)
	LandedCostController := masteritemcontroller.NewLandedCostMasterController(LandedCostService)

	// Operation Group
	operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl()
	operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository, db, rdb)
	operationGroupController := masteroperationcontroller.NewOperationGroupController(operationGroupService)

	// Incentive Group
	IncentiveGroupRepository := masterrepositoryimpl.StartIncentiveGroupRepositoryImpl()
	IncentiveGroupService := masterserviceimpl.StartIncentiveGroupService(IncentiveGroupRepository, db, rdb)
	IncentiveGroupController := mastercontroller.NewIncentiveGroupController(IncentiveGroupService)

	// IncentiveGroupDetail
	IncentiveGroupDetailRepository := masterrepositoryimpl.StartIncentiveGroupDetailRepositoryImpl()
	IncentiveGroupDetailService := masterserviceimpl.StartIncentiveGroupDetailService(IncentiveGroupDetailRepository, db, rdb)
	IncentiveGroupDetailController := mastercontroller.NewIncentiveGroupDetailController(IncentiveGroupDetailService)

	// MovingCode
	MovingCodeRepository := masterrepositoryimpl.StartMovingCodeRepositoryImpl()
	MovingCodeService := masterserviceimpl.StartMovingCodeServiceImpl(MovingCodeRepository, db)
	MovingCodeController := mastercontroller.NewMovingCodeController(MovingCodeService)

	// ForecastMaster
	forecastMasterRepository := masterrepositoryimpl.StartForecastMasterRepositoryImpl()
	forecastMasterService := masterserviceimpl.StartForecastMasterService(forecastMasterRepository, db, rdb)
	forecastMasterController := mastercontroller.NewForecastMasterController(forecastMasterService)

	// Agreement
	AgreementRepository := masterrepositoryimpl.StartAgreementRepositoryImpl()
	AgreementService := masterserviceimpl.StartAgreementService(AgreementRepository, db, rdb)
	AgreementController := mastercontroller.NewAgreementController(AgreementService)

	// operation code
	operationCodeRepository := masteroperationrepositoryimpl.StartOperationCodeRepositoryImpl()
	operationCodeService := masteroperationserviceimpl.StartOperationCodeService(operationCodeRepository, db, rdb)
	operationCodeController := masteroperationcontroller.NewOperationCodeController(operationCodeService)

	// Operation Section
	operationSectionRepository := masteroperationrepositoryimpl.StartOperationSectionRepositoryImpl()
	operationSectionService := masteroperationserviceimpl.StartOperationSectionService(operationSectionRepository, db, rdb)
	operationSectionController := masteroperationcontroller.NewOperationSectionController(operationSectionService)

	//OperationEntries
	operationEntriesRepository := masteroperationrepositoryimpl.StartOperationEntriesRepositoryImpl()
	operationEntriesService := masteroperationserviceimpl.StartOperationEntriesService(operationEntriesRepository, db, rdb)
	operationEntriesController := masteroperationcontroller.NewOperationEntriesController(operationEntriesService)

	// Operation Key
	operationKeyRepository := masteroperationrepositoryimpl.StartOperationKeyRepositoryImpl()
	operationKeyService := masteroperationserviceimpl.StartOperationKeyService(operationKeyRepository, db, rdb)
	operationKeyController := masteroperationcontroller.NewOperationKeyController(operationKeyService)

	// operation model mapping
	operationModelMappingRepository := masteroperationrepositoryimpl.StartOperationModelMappingRepositoryImpl()
	operationModelMappingService := masteroperationserviceimpl.StartOperationModelMappingService(operationModelMappingRepository, db, rdb)
	operationModelMappingController := masteroperationcontroller.NewOperationModelMappingController(operationModelMappingService)

	// Skill Level
	SkillLevelRepository := masterrepositoryimpl.StartSkillLevelRepositoryImpl()
	SkillLevelService := masterserviceimpl.StartSkillLevelService(SkillLevelRepository, db, rdb)
	SkillLevelController := mastercontroller.NewSkillLevelController(SkillLevelService)

	// Shift Schedule
	ShiftScheduleRepository := masterrepositoryimpl.StartShiftScheduleRepositoryImpl()
	ShiftScheduleService := masterserviceimpl.StartShiftScheduleService(ShiftScheduleRepository, db, rdb)
	ShiftScheduleController := mastercontroller.NewShiftScheduleController(ShiftScheduleService)

	// Discount Percent
	discountPercentRepository := masteritemrepositoryimpl.StartDiscountPercentRepositoryImpl()
	discountPercentService := masteritemserviceimpl.StartDiscountPercentService(discountPercentRepository, db, rdb)
	discountPercentController := masteritemcontroller.NewDiscountPercentController(discountPercentService)

	// Discount
	discountRepository := masterrepositoryimpl.StartDiscountRepositoryImpl()
	discountService := masterserviceimpl.StartDiscountService(discountRepository, db, rdb)
	discountController := mastercontroller.NewDiscountController(discountService)

	// Markup Rate
	markupRateRepository := masteritemrepositoryimpl.StartMarkupRateRepositoryImpl()
	markupRateService := masteritemserviceimpl.StartMarkupRateService(markupRateRepository, db, rdb)
	markupRateController := masteritemcontroller.NewMarkupRateController(markupRateService)

	// Warehouse Group
	warehouseGroupRepository := masterwarehouserepositoryimpl.OpenWarehouseGroupImpl()
	warehouseGroupService := masterwarehouseserviceimpl.OpenWarehouseGroupService(warehouseGroupRepository, db, rdb)
	warehouseGroupController := masterwarehousecontroller.NewWarehouseGroupController(warehouseGroupService)

	// Warehouse Location
	WarehouseLocationDefinitionRepository := masterwarehouserepositoryimpl.OpenWarehouseLocationDefinitionImpl()
	WarehouseLocationDefinitionService := masterwarehouseserviceimpl.OpenWarehouseLocationDefinitionService(WarehouseLocationDefinitionRepository, db, rdb)
	WarehouseLocationDefinitionController := masterwarehousecontroller.NewWarehouseLocationDefinitionController(WarehouseLocationDefinitionService)

	// Warehouse Master
	warehouseMasterRepository := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl()
	warehouseMasterService := masterwarehouseserviceimpl.OpenWarehouseMasterService(warehouseMasterRepository, db, rdb)
	warehouseMasterController := masterwarehousecontroller.NewWarehouseMasterController(warehouseMasterService)

	// Warehouse Location
	warehouseLocationRepository := masterwarehouserepositoryimpl.OpenWarehouseLocationImpl()
	warehouseLocationService := masterwarehouseserviceimpl.OpenWarehouseLocationService(warehouseLocationRepository, warehouseMasterService, db, rdb)
	warehouseLocationController := masterwarehousecontroller.NewWarehouseLocationController(warehouseLocationService)

	// Bom Master
	BomRepository := masteritemrepositoryimpl.StartBomRepositoryImpl()
	BomService := masteritemserviceimpl.StartBomService(BomRepository, db, rdb)
	BomController := masteritemcontroller.NewBomController(BomService)

	//package master
	PackageMasterRepository := masterrepositoryimpl.StartPackageMasterRepositoryImpl()
	PackageMasterService := masterserviceimpl.StartPackageMasterService(PackageMasterRepository, db)
	PackageMasterController := mastercontroller.NewPackageMasterController(PackageMasterService)

	// Deduction
	DeductionRepository := masterrepositoryimpl.StartDeductionRepositoryImpl()
	DeductionService := masterserviceimpl.StartDeductionService(DeductionRepository, db, rdb)
	DeductionController := mastercontroller.NewDeductionController(DeductionService)

	// Warranty Free Service
	WarrantyFreeServiceRepository := masterrepositoryimpl.StartWarrantyFreeServiceRepositoryImpl()
	WarrantyFreeServiceService := masterserviceimpl.StartWarrantyFreeServiceService(WarrantyFreeServiceRepository, db, rdb)
	WarrantyFreeServiceController := mastercontroller.NewWarrantyFreeServiceController(WarrantyFreeServiceService)

	// Incentive Master
	IncentiveMasterRepository := masterrepositoryimpl.StartIncentiveMasterRepositoryImpl()
	IncentiveMasterService := masterserviceimpl.StartIncentiveMasterService(IncentiveMasterRepository, db, rdb)
	IncentiveMasterController := mastercontroller.NewIncentiveMasterController(IncentiveMasterService)

	// Skill Level
	SkillLevelRepository := masterrepositoryimpl.StartSkillLevelRepositoryImpl()
	SkillLevelService := masterserviceimpl.StartSkillLevelService(SkillLevelRepository, db)
	SkillLevelController := mastercontroller.NewSkillLevelController(SkillLevelService)

	// Master
	//Field Action
	FieldActionRepository := masterrepositoryimpl.StartFieldActionRepositoryImpl()
	FieldActionService := masterserviceimpl.StartFieldActionService(FieldActionRepository, db, rdb)
	FieldActionController := mastercontroller.NewFieldActionController(FieldActionService)

	/* Transaction */
	//Supply Slip
	SupplySlipRepository := transactionsparepartrepositoryimpl.StartSupplySlipRepositoryImpl()
	SupplySlipService := transactionsparepartserviceimpl.StartSupplySlipService(SupplySlipRepository, db, rdb)
	SupplySlipController := transactionsparepartcontroller.NewSupplySlipController(SupplySlipService)

	//Booking Estimation
	BookingEstimationRepository := transactionworkshoprepositoryimpl.OpenBookingEstimationRepositoryImpl()
	BookingEstimationService := transactionworkshopserviceimpl.OpenBookingEstimationServiceImpl(BookingEstimationRepository, db, rdb)
	BookingEstimationController := transactionworksopcontroller.NewBookingEstimationController(BookingEstimationService)

	//Work order
	WorkOrderRepository := transactionworkshoprepositoryimpl.OpenWorkOrderRepositoryImpl()
	WorkOrderService := transactionworkshopserviceimpl.OpenWorkOrderServiceImpl(WorkOrderRepository, db, rdb)
	WorkOrderController := transactionworksopcontroller.NewWorkOrderController(WorkOrderService)

	//Sales Order
	SalesOrderRepository := transactionsparepartrepositoryimpl.StartSalesOrderRepositoryImpl()
	SalesOrderService := transactionsparepartserviceimpl.StartSalesOrderService(SalesOrderRepository, db, rdb)
	SalesOrderController := transactionsparepartcontroller.NewSalesOrderController(SalesOrderService)

	//Service Request
	ServiceRequestRepository := transactionworkshoprepositoryimpl.OpenServiceRequestRepositoryImpl()
	ServiceRequestService := transactionworkshopserviceimpl.OpenServiceRequestServiceImpl(ServiceRequestRepository, db, rdb)
	ServiceRequestController := transactionworksopcontroller.NewServiceRequestController(ServiceRequestService)

	//vehicle history
	VehicleHistoryRepository := transactionworkshoprepositoryimpl.NewVehicleHistoryImpl()
	VehicleHistoryServices := transactionworkshopserviceimpl.NewVehicleHistoryServiceImpl(VehicleHistoryRepository, db, rdb)
	VehicleHistoryController := transactionworksopcontroller.NewVehicleHistoryController(VehicleHistoryServices)

	//Service Receipt
	ServiceReceiptRepository := transactionworkshoprepositoryimpl.OpenServiceReceiptRepositoryImpl()
	ServiceReceiptService := transactionworkshopserviceimpl.OpenServiceReceiptServiceImpl(ServiceReceiptRepository, db, rdb)
	ServiceReceiptController := transactionworksopcontroller.NewServiceReceiptController(ServiceReceiptService)

	//Work order bypass
	WorkOrderBypassRepository := transactionworkshoprepositoryimpl.OpenWorkOrderBypassRepositoryImpl()
	WorkOrderBypassService := transactionworkshopserviceimpl.OpenWorkOrderBypassServiceImpl(WorkOrderBypassRepository, db, rdb)
	WorkOrderBypassController := transactionworksopcontroller.NewWorkOrderBypassController(WorkOrderBypassService)

	/* Master */
	itemClassRouter := ItemClassRouter(itemClassController)
	itemPackageRouter := ItemPackageRouter(itemPackageController)
	ItemModelMappingRouter := ItemModelMappingRouter(ItemModelMappingController)
	itemPackageDetailRouter := ItemPackageDetailRouter(itemPackageDetailController)
	itemImportRouter := ItemImportRouter(ItemImportController)
	OperationGroupRouter := OperationGroupRouter(operationGroupController)
	PurchasePriceRouter := PurchasePriceRouter(PurchasePriceController)
	LandedCostMasterRouter := LandedCostMasterRouter(LandedCostController)
	IncentiveGroupRouter := IncentiveGroupRouter(IncentiveGroupController)
	IncentiveGroupDetailRouter := IncentiveGroupDetailRouter(IncentiveGroupDetailController)
	IncentiveMasterRouter := IncentiveMasterRouter(IncentiveMasterController)
	OperationCodeRouter := OperationCodeRouter(operationCodeController)
	OperationSectionRouter := OperationSectionRouter(operationSectionController)
	OperationEntriesRouter := OperationEntriesRouter(operationEntriesController)
	OperationKeyRouter := OperationKeyRouter(operationKeyController)
	OperationModelMappingRouter := OperationModelMappingRouter(operationModelMappingController)
	MovingCodeRouter := MovingCodeRouter(MovingCodeController)
	ForecastMasterRouter := ForecastMasterRouter(forecastMasterController)
	AgreementRouter := AgreementRouter(AgreementController)
	DiscountPercentRouter := DiscountPercentRouter(discountPercentController)
	DiscountRouter := DiscountRouter(discountController)
	MarkupRateRouter := MarkupRateRouter(markupRateController)
	ItemSubstituteRouter := ItemSubstituteRouter(itemSubstituteController)
	ItemLocationRouter := ItemLocationRouter(ItemLocationController)
	WarehouseGroupRouter := WarehouseGroupRouter(warehouseGroupController)
	WarehouseLocation := WarehouseLocationRouter(warehouseLocationController)
	WarehouseLocationDefinition := WarehouseLocationDefinitionRouter(WarehouseLocationDefinitionController)
	WarehouseMaster := WarehouseMasterRouter(warehouseMasterController)
	SkillLevelRouter := SkillLevelRouter(SkillLevelController)
	ShiftScheduleRouter := ShiftScheduleRouter(ShiftScheduleController)
	unitOfMeasurementRouter := UnitOfMeasurementRouter(unitOfMeasurementController)
	markupMasterRouter := MarkupMasterRouter(markupMasterController)
	itemLevelRouter := ItemLevelRouter(itemLevelController)
	itemRouter := ItemRouter(itemController)
	priceListRouter := PriceListRouter(priceListController)
	FieldActionRouter := FieldActionRouter(FieldActionController)
	warrantyFreeServiceRouter := WarrantyFreeServiceRouter(WarrantyFreeServiceController)
	BomRouter := BomRouter(BomController)
	DeductionRouter := DeductionRouter(DeductionController)
	SkillLevelRouter := SkillLevelRouter(SkillLevelController)

	r := chi.NewRouter()
	r.Mount("/item-class", itemClassRouter)
	r.Mount("/unit-of-measurement", unitOfMeasurementRouter)
	r.Mount("/markup-master", markupMasterRouter)
	r.Mount("/item-level", itemLevelRouter)
	// mux.Handle("/operation-group/", OperationGroupRouter)
	r.Mount("/operation-group", OperationGroupRouter)
	r.Mount("/incentive", IncentiveMasterRouter)
	r.Mount("/bom", BomRouter)
	r.Mount("/deduction", DeductionRouter)

	r.Mount("/item-package", itemPackageRouter) //null value
	// r.Mount("/item-package-detail", itemPackageDetailRouter) //notfound
	r.Mount("/item", itemRouter) //error mssql: The correlation name 'mtr_item_class' is specified multiple times in a FROM clause.
	r.Mount("/item-substitute", ItemSubstituteRouter)

	r.Mount("/incentive-group", IncentiveGroupRouter)
	r.Mount("/incentive-group-detail", IncentiveGroupDetailRouter) //method notalowed

	r.Mount("/operation-code", OperationCodeRouter)
	r.Mount("/operation-section", OperationSectionRouter)
	r.Mount("/operation-key", OperationKeyRouter)
	r.Mount("/operation-entries", OperationEntriesRouter)

	r.Mount("/discount-percent", DiscountPercentRouter) //error Could not get response
	r.Mount("/discount", DiscountRouter)

	r.Mount("/markup-rate", MarkupRateRouter) //error Could not get response

	r.Mount("/warehouse-group", WarehouseGroupRouter) //null value
	r.Mount("/warehouse-location", WarehouseLocation)
	r.Mount("/warehouse-master", WarehouseMaster)
	r.Mount("/warehouse-free-service", warrantyFreeServiceRouter)

	r.Mount("/forecast-master", ForecastMasterRouter) //error Could not get response
	r.Mount("/shift-schedule", ShiftScheduleRouter)
	r.Mount("/price-list", priceListRouter) //null value
	r.Mount("/warranty-free-service", warrantyFreeServiceRouter)
	r.Mount("/skill-level", SkillLevelRouter)

	server := http.Server{
		Addr:    config.EnvConfigs.ClientOrigin,
		Handler: r,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
