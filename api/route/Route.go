package route

import (
	"after-sales/api/config"
	"after-sales/api/helper"
	masteritemrepositoryimpl "after-sales/api/repositories/master/item/repositories-item-impl"
	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	masterwarehouserepositoryimpl "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"

	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"

	masteritemserviceimpl "after-sales/api/services/master/item/services-item-impl"
	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	masterwarehouseserviceimpl "after-sales/api/services/master/warehouse/services-warehouse-impl"

	mastercontroller "after-sales/api/controllers/master"
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	masterwarehousecontroller "after-sales/api/controllers/master/warehouse"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func StartRouting(db *gorm.DB) {
	// Unit Measurement
	unitOfMeasurementRepository := masteritemrepositoryimpl.StartUnitOfMeasurementRepositoryImpl()
	unitOfMeasurementService := masteritemserviceimpl.StartUnitOfMeasurementService(unitOfMeasurementRepository, db)
	unitOfMeasurementController := masteritemcontroller.NewUnitOfMeasurementController(unitOfMeasurementService)

	// Markup Master
	markupMasterRepository := masteritemrepositoryimpl.StartMarkupMasterRepositoryImpl()
	markupMasterService := masteritemserviceimpl.StartMarkupMasterService(markupMasterRepository, db)
	markupMasterController := masteritemcontroller.NewMarkupMasterController(markupMasterService)

	// Item Level
	itemLevelRepository := masteritemrepositoryimpl.StartItemLevelRepositoryImpl()
	itemLevelService := masteritemserviceimpl.StartItemLevelService(itemLevelRepository, db)
	itemLevelController := masteritemcontroller.NewItemLevelController(itemLevelService)

	// Item
	itemRepository := masteritemrepositoryimpl.StartItemRepositoryImpl()
	itemService := masteritemserviceimpl.StartItemService(itemRepository, db)
	itemController := masteritemcontroller.NewItemController(itemService)

	// PriceList
	priceListRepository := masteritemrepositoryimpl.StartPriceListRepositoryImpl()
	priceListService := masteritemserviceimpl.StartPriceListService(priceListRepository, db)
	priceListController := masteritemcontroller.NewPriceListController(priceListService)

	// Item Class
	itemClassRepository := masteritemrepositoryimpl.StartItemClassRepositoryImpl()
	itemClassService := masteritemserviceimpl.StartItemClassService(itemClassRepository, db)
	itemClassController := masteritemcontroller.NewItemClassController(itemClassService)

	// Item Substitute
	itemSubstituteRepository := masteritemrepositoryimpl.StartItemSubstituteRepositoryImpl()
	itemSubstituteService := masteritemserviceimpl.StartItemSubstituteService(itemSubstituteRepository, db)
	itemSubstituteController := masteritemcontroller.NewItemSubstituteController(itemSubstituteService)

	// Item Package
	itemPackageRepository := masteritemrepositoryimpl.StartItemPackageRepositoryImpl()
	itemPackageService := masteritemserviceimpl.StartItemPackageService(itemPackageRepository, db)
	itemPackageController := masteritemcontroller.NewItemPackageController(itemPackageService)

	// Item Package Detail
	itemPackageDetailRepository := masteritemrepositoryimpl.StartItemPackageDetailRepositoryImpl()
	itemPackageDetailService := masteritemserviceimpl.StartItemPackageDetailService(itemPackageDetailRepository, db)
	itemPackageDetailController := masteritemcontroller.NewItemPackageDetailController(itemPackageDetailService)

	// Operation Group
	operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl()
	operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository, db)
	operationGroupController := masteroperationcontroller.NewOperationGroupController(operationGroupService)

	// Incentive Group
	IncentiveGroupRepository := masterrepositoryimpl.StartIncentiveGroupRepositoryImpl()
	IncentiveGroupService := masterserviceimpl.StartIncentiveGroupService(IncentiveGroupRepository, db)
	IncentiveGroupController := mastercontroller.NewIncentiveGroupController(IncentiveGroupService)

	// IncentiveGroupDetail
	IncentiveGroupDetailRepository := masterrepositoryimpl.StartIncentiveGroupDetailRepositoryImpl()
	IncentiveGroupDetailService := masterserviceimpl.StartIncentiveGroupDetailService(IncentiveGroupDetailRepository, db)
	IncentiveGroupDetailController := mastercontroller.NewIncentiveGroupDetailController(IncentiveGroupDetailService)

	// ForecastMaster
	forecastMasterRepository := masterrepositoryimpl.StartForecastMasterRepositoryImpl()
	forecastMasterService := masterserviceimpl.StartForecastMasterService(forecastMasterRepository, db)
	forecastMasterController := mastercontroller.NewForecastMasterController(forecastMasterService)

	// operation code
	operationCodeRepository := masteroperationrepositoryimpl.StartOperationCodeRepositoryImpl()
	operationCodeService := masteroperationserviceimpl.StartOperationCodeService(operationCodeRepository, db)
	operationCodeController := masteroperationcontroller.NewOperationCodeController(operationCodeService)

	// Operation Section
	operationSectionRepository := masteroperationrepositoryimpl.StartOperationSectionRepositoryImpl()
	operationSectionService := masteroperationserviceimpl.StartOperationSectionService(operationSectionRepository, db)
	operationSectionController := masteroperationcontroller.NewOperationSectionController(operationSectionService)

	//OperationEntries
	operationEntriesRepository := masteroperationrepositoryimpl.StartOperationEntriesRepositoryImpl()
	operationEntriesService := masteroperationserviceimpl.StartOperationEntriesService(operationEntriesRepository, db)
	operationEntriesController := masteroperationcontroller.NewOperationEntriesController(operationEntriesService)

	// Operation Key
	operationKeyRepository := masteroperationrepositoryimpl.StartOperationKeyRepositoryImpl()
	operationKeyService := masteroperationserviceimpl.StartOperationKeyService(operationKeyRepository, db)
	operationKeyController := masteroperationcontroller.NewOperationKeyController(operationKeyService)

	// Shift Schedule
	ShiftScheduleRepository := masterrepositoryimpl.StartShiftScheduleRepositoryImpl()
	ShiftScheduleService := masterserviceimpl.StartShiftScheduleService(ShiftScheduleRepository, db)
	ShiftScheduleController := mastercontroller.NewShiftScheduleController(ShiftScheduleService)

	// Discount Percent
	discountPercentRepository := masteritemrepositoryimpl.StartDiscountPercentRepositoryImpl()
	discountPercentService := masteritemserviceimpl.StartDiscountPercentService(discountPercentRepository, db)
	discountPercentController := masteritemcontroller.NewDiscountPercentController(discountPercentService)

	// Discount
	discountRepository := masterrepositoryimpl.StartDiscountRepositoryImpl()
	discountService := masterserviceimpl.StartDiscountService(discountRepository, db)
	discountController := mastercontroller.NewDiscountController(discountService)

	// Markup Rate
	markupRateRepository := masteritemrepositoryimpl.StartMarkupRateRepositoryImpl()
	markupRateService := masteritemserviceimpl.StartMarkupRateService(markupRateRepository, db)
	markupRateController := masteritemcontroller.NewMarkupRateController(markupRateService)

	// Warehouse Group
	warehouseGroupRepository := masterwarehouserepositoryimpl.OpenWarehouseGroupImpl()
	warehouseGroupService := masterwarehouseserviceimpl.OpenWarehouseGroupService(warehouseGroupRepository, db)
	warehouseGroupController := masterwarehousecontroller.NewWarehouseGroupController(warehouseGroupService)

	// Warehouse Location
	warehouseLocationRepository := masterwarehouserepositoryimpl.OpenWarehouseLocationImpl()
	warehouseLocationService := masterwarehouseserviceimpl.OpenWarehouseLocationService(warehouseLocationRepository, db)
	warehouseLocationController := masterwarehousecontroller.NewWarehouseLocationController(warehouseLocationService)

	// Warehouse Master
	warehouseMasterRepository := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl()
	warehouseMasterService := masterwarehouseserviceimpl.OpenWarehouseMasterService(warehouseMasterRepository, db)
	warehouseMasterController := masterwarehousecontroller.NewWarehouseMasterController(warehouseMasterService)

	// Bom Master
	BomRepository := masteritemrepositoryimpl.StartBomRepositoryImpl()
	BomService := masteritemserviceimpl.StartBomService(BomRepository, db)
	BomController := masteritemcontroller.NewBomController(BomService)

	// Deduction
	DeductionRepository := masterrepositoryimpl.StartDeductionRepositoryImpl()
	DeductionService := masterserviceimpl.StartDeductionService(DeductionRepository, db)
	DeductionController := mastercontroller.NewDeductionController(DeductionService)

	// Warranty Free Service
	WarrantyFreeServiceRepository := masterrepositoryimpl.StartWarrantyFreeServiceRepositoryImpl()
	WarrantyFreeServiceService := masterserviceimpl.StartWarrantyFreeServiceService(WarrantyFreeServiceRepository, db)
	WarrantyFreeServiceController := mastercontroller.NewWarrantyFreeServiceController(WarrantyFreeServiceService)

	// Incentive Master
	IncentiveMasterRepository := masterrepositoryimpl.StartIncentiveMasterRepositoryImpl()
	IncentiveMasterService := masterserviceimpl.StartIncentiveMasterService(IncentiveMasterRepository, db)
	IncentiveMasterController := mastercontroller.NewIncentiveMasterController(IncentiveMasterService)

	// Field Action
	FieldActionRepository := masterrepositoryimpl.StartFieldActionRepositoryImpl()
	FieldActionService := masterserviceimpl.StartFieldActionService(FieldActionRepository, db)
	FieldActionController := mastercontroller.NewFieldActionController(FieldActionService)

	// Master
	itemClassRouter := ItemClassRouter(itemClassController)
	itemPackageRouter := ItemPackageRouter(itemPackageController)
	itemPackageDetailRouter := ItemPackageDetailRouter(itemPackageDetailController)
	OperationGroupRouter := OperationGroupRouter(operationGroupController)
	IncentiveGroupRouter := IncentiveGroupRouter(IncentiveGroupController)
	IncentiveGroupDetailRouter := IncentiveGroupDetailRouter(IncentiveGroupDetailController)
	IncentiveMasterRouter := IncentiveMasterRouter(IncentiveMasterController)
	OperationCodeRouter := OperationCodeRouter(operationCodeController)
	OperationSectionRouter := OperationSectionRouter(operationSectionController)
	OperationEntriesRouter := OperationEntriesRouter(operationEntriesController)
	OperationKeyRouter := OperationKeyRouter(operationKeyController)
	ForecastMasterRouter := ForecastMasterRouter(forecastMasterController)
	DiscountPercentRouter := DiscountPercentRouter(discountPercentController)
	DiscountRouter := DiscountRouter(discountController)
	MarkupRateRouter := MarkupRateRouter(markupRateController)
	ItemSubstituteRouter := ItemSubstituteRouter(itemSubstituteController)
	WarehouseGroupRouter := WarehouseGroupRouter(warehouseGroupController)
	WarehouseLocation := WarehouseLocationRouter(warehouseLocationController)
	WarehouseMaster := WarehouseMasterRouter(warehouseMasterController)
	ShiftScheduleRouter := ShiftScheduleRouter(ShiftScheduleController)
	unitOfMeasurementRouter := UnitOfMeasurementRouter(unitOfMeasurementController)
	markupMasterRouter := MarkupMasterRouter(markupMasterController)
	itemLevelRouter := ItemLevelRouter(itemLevelController)
	itemRouter := ItemRouter(itemController)
	priceListRouter := PriceListRouter(priceListController)
	warrantyFreeServiceRouter := WarrantyFreeServiceRouter(WarrantyFreeServiceController)
	BomRouter := BomRouter(BomController)
	DeductionRouter := DeductionRouter(DeductionController)
	FieldActionRouter := FieldActionRouter(FieldActionController)

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

	r.Mount("/item-package", itemPackageRouter)              //null value
	r.Mount("/item-package-detail", itemPackageDetailRouter) //notfound
	r.Mount("/item", itemRouter)                             //error mssql: The correlation name 'mtr_item_class' is specified multiple times in a FROM clause.
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
	r.Mount("/field-action", FieldActionRouter)

	r.Mount("/forecast-master", ForecastMasterRouter) //error Could not get response
	r.Mount("/shift-schedule", ShiftScheduleRouter)
	r.Mount("/price-list", priceListRouter) //null value
	r.Mount("/warranty-free-service", warrantyFreeServiceRouter)

	server := http.Server{
		Addr:    config.EnvConfigs.ClientOrigin,
		Handler: r,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
