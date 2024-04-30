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

	transactionworksopcontroller "after-sales/api/controllers/transactions/workshop"
	transactionworkshoprepositoryimpl "after-sales/api/repositories/transaction/workshop/repositories-workshop-impl"
	transactionworkshopserviceimpl "after-sales/api/services/transaction/workshop/services-workshop-impl"

	"net/http"

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
	MovingCodeService := masterserviceimpl.StartMovingCodeService(MovingCodeRepository, db, rdb)
	MovingCodeController := mastercontroller.NewMovingCodeController(MovingCodeService)

	// ForecastMaster
	forecastMasterRepository := masterrepositoryimpl.StartForecastMasterRepositoryImpl()
	forecastMasterService := masterserviceimpl.StartForecastMasterService(forecastMasterRepository, db, rdb)
	forecastMasterController := mastercontroller.NewForecastMasterController(forecastMasterService)

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

	// Warehouse Location
	warehouseLocationRepository := masterwarehouserepositoryimpl.OpenWarehouseLocationImpl()
	warehouseLocationService := masterwarehouseserviceimpl.OpenWarehouseLocationService(warehouseLocationRepository, db, rdb)
	warehouseLocationController := masterwarehousecontroller.NewWarehouseLocationController(warehouseLocationService)

	// Warehouse Master
	warehouseMasterRepository := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl()
	warehouseMasterService := masterwarehouseserviceimpl.OpenWarehouseMasterService(warehouseMasterRepository, db, rdb)
	warehouseMasterController := masterwarehousecontroller.NewWarehouseMasterController(warehouseMasterService)

	// Bom Master
	BomRepository := masteritemrepositoryimpl.StartBomRepositoryImpl()
	BomService := masteritemserviceimpl.StartBomService(BomRepository, db, rdb)
	BomController := masteritemcontroller.NewBomController(BomService)

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

	//Field Action
	FieldActionRepository := masterrepositoryimpl.StartFieldActionRepositoryImpl()
	FieldActionService := masterserviceimpl.StartFieldActionService(FieldActionRepository, db, rdb)
	FieldActionController := mastercontroller.NewFieldActionController(FieldActionService)

	/* Transaction */
	//Work order
	WorkOrderRepository := transactionworkshoprepositoryimpl.OpenWorkOrderRepositoryImpl()
	WorkOrderService := transactionworkshopserviceimpl.OpenWorkOrderServiceImpl(WorkOrderRepository, db, rdb)
	WorkOrderController := transactionworksopcontroller.NewWorkOrderController(WorkOrderService)

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

	/* Transaction */
	WorkOrderRouter := WorkOrderRouter(WorkOrderController)

	r := chi.NewRouter()
	// Route untuk setiap versi API
	r.Route("/v1", func(r chi.Router) {
		// Tambahkan routing untuk v1 versi di sini
		/* Master */
		r.Mount("/item-class", itemClassRouter)
		r.Mount("/unit-of-measurement", unitOfMeasurementRouter)
		r.Mount("/discount-percent", DiscountPercentRouter)
		r.Mount("/markup-master", markupMasterRouter)
		r.Mount("/markup-rate", MarkupRateRouter)
		r.Mount("/item-level", itemLevelRouter)
		r.Mount("/item", itemRouter)
		r.Mount("/item-substitute", ItemSubstituteRouter)
		r.Mount("/item-location", ItemLocationRouter)
		r.Mount("/item-package", itemPackageRouter)
		r.Mount("/item-package-detail", itemPackageDetailRouter)
		r.Mount("/price-list", priceListRouter)
		r.Mount("/item-model-mapping", ItemModelMappingRouter)
		//r.Mount("/import-item", ImportItemRouter)
		r.Mount("/bom", BomRouter)
		r.Mount("/item-import", itemImportRouter)
		r.Mount("/purchase-price", PurchasePriceRouter)
		r.Mount("/landed-cost", LandedCostMasterRouter)
		//r.Mount("/import-duty", ImportDutyRouter)
		r.Mount("/operation-group", OperationGroupRouter)
		r.Mount("/operation-section", OperationSectionRouter)
		r.Mount("/operation-key", OperationKeyRouter)
		r.Mount("/operation-entries", OperationEntriesRouter)
		r.Mount("/operation-code", OperationCodeRouter)
		r.Mount("/operation-model-mapping", OperationModelMappingRouter)
		//r.Mount("/labour-selling-price", LabourSellingPriceRouter)
		r.Mount("/warehouse-group", WarehouseGroupRouter)
		r.Mount("/warehouse-master", WarehouseMaster)
		r.Mount("/warehouse-location-definition", WarehouseLocationDefinition)
		r.Mount("/warehouse-location", WarehouseLocation)
		r.Mount("/moving-code", MovingCodeRouter)
		r.Mount("/forecast-master", ForecastMasterRouter)
		//r.Mount("/agreement", AgreementRouter)
		//r.Mount("/campaign", CampaignRouter)
		//r.Mount("/package", PackageRouter)
		r.Mount("/skill-level", SkillLevelRouter)
		r.Mount("/shift-schedule", ShiftScheduleRouter)
		r.Mount("/incentive", IncentiveMasterRouter)
		//r.Mount("/work-info-massage", WorkInfoRouter)
		r.Mount("/field-action", FieldActionRouter)
		r.Mount("/warranty-free-service", warrantyFreeServiceRouter)
		r.Mount("/discount", DiscountRouter)
		r.Mount("/incentive-group", IncentiveGroupRouter)
		r.Mount("/incentive-group-detail", IncentiveGroupDetailRouter)
		r.Mount("/deduction", DeductionRouter)

		/* Transaction */
		r.Mount("/work-order", WorkOrderRouter)

	})

	// Route untuk utilities
	r.Route("/utilities", func(r chi.Router) {
		// Route untuk Swagger
		r.Mount("/swagger", SwaggerRouter())
		// Route untuk Prometheus metrics
		r.Mount("/metrics", promhttp.Handler())
	})

	server := http.Server{
		Addr:    config.EnvConfigs.ClientOrigin,
		Handler: r,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
