package main

import (
	"after-sales/api/config"
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	masterwarehousecontroller "after-sales/api/controllers/master/warehouse"

	// masteroperationcontroller "after-sales/api/controllers/master/operation"
	mastercontroller "after-sales/api/controllers/master"
	"after-sales/api/helper"
	masteritemrepositoryimpl "after-sales/api/repositories/master/item/repositories-item-impl"
	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	masterwarehouserepositoryimpl "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"

	// masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"

	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	masteritemserviceimpl "after-sales/api/services/master/item/services-item-impl"
	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
	masterwarehouseserviceimpl "after-sales/api/services/master/warehouse/services-warehouse-impl"

	// masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"

	route "after-sales/api/route"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	migration "after-sales/generate/sql"
	"net/http"
	"os"
)

// @title After Sales API
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
// @host localhost:8000
// @BasePath /api/aftersales
func main() {
	args := os.Args
	env := ""
	if len(args) > 1 {
		env = args[1]
	}

	if env == "migrate" {
		migration.Migrate()
	} else if env == "generate" {
		migration.Generate()
	} else if env == "migg" {
		migration.MigrateGG()
	} else if env == "debug" {
		// config.InitEnvConfigs(false, env)
		// db := config.InitDB()
		// config.InitLogger(db)
		// redis := config.InitRedis()
		// route.CreateHandler(db, env, redis)
	} else {
		config.InitEnvConfigs(false, env)
		db := config.InitDB()
		// redis := config.InitRedis()
		// route.CreateHandler(db, env, redis)
		config.InitLogger(db)

		itemClassRepository := masteritemrepositoryimpl.StartItemClassRepositoryImpl()
		itemClassService := masteritemserviceimpl.StartItemClassService(itemClassRepository, db)
		itemClassController := masteritemcontroller.NewItemClassController(itemClassService)

		operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl()
		operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository, db)
		operationGroupController := masteroperationcontroller.NewOperationGroupController(operationGroupService)

		IncentiveGroupRepository := masterrepositoryimpl.StartIncentiveGroupRepositoryImpl()
		IncentiveGroupService := masterserviceimpl.StartIncentiveGroupService(IncentiveGroupRepository, db)
		IncentiveGroupController := mastercontroller.NewIncentiveGroupController(IncentiveGroupService)

		IncentiveGroupDetailRepository := masterrepositoryimpl.StartIncentiveGroupDetailRepositoryImpl()
		IncentiveGroupDetailService := masterserviceimpl.StartIncentiveGroupDetailService(IncentiveGroupDetailRepository, db)
		IncentiveGroupDetailController := mastercontroller.NewIncentiveGroupDetailController(IncentiveGroupDetailService)

		forecastMasterRepository := masterrepositoryimpl.StartForecastMasterRepositoryImpl()
		forecastMasterService := masterserviceimpl.StartForecastMasterService(forecastMasterRepository, db)
		forecastMasterController := mastercontroller.NewForecastMasterController(forecastMasterService)

		operationSectionRepository := masteroperationrepositoryimpl.StartOperationSectionRepositoryImpl()
		operationSectionService := masteroperationserviceimpl.StartOperationSectionService(operationSectionRepository, db)
		operationSectionController := masteroperationcontroller.NewOperationSectionController(operationSectionService)

		operationEntriesRepository := masteroperationrepositoryimpl.StartOperationEntriesRepositoryImpl()
		operationEntriesService := masteroperationserviceimpl.StartOperationEntriesService(operationEntriesRepository, db)
		operationEntriesController := masteroperationcontroller.NewOperationEntriesController(operationEntriesService)

		operationKeyRepository := masteroperationrepositoryimpl.StartOperationKeyRepositoryImpl()
		operationKeyService := masteroperationserviceimpl.StartOperationKeyService(operationKeyRepository, db)
		operationKeyController := masteroperationcontroller.NewOperationKeyController(operationKeyService)

		ShiftScheduleRepository := masterrepositoryimpl.StartShiftScheduleRepositoryImpl()
		ShiftScheduleService := masterserviceimpl.StartShiftScheduleService(ShiftScheduleRepository, db)
		ShiftScheduleController := mastercontroller.NewShiftScheduleController(ShiftScheduleService)

		discountPercentRepository := masteritemrepositoryimpl.StartDiscountPercentRepositoryImpl()
		discountPercentService := masteritemserviceimpl.StartDiscountPercentService(discountPercentRepository, db)
		discountPercentController := masteritemcontroller.NewDiscountPercentController(discountPercentService)

		discountRepository := masterrepositoryimpl.StartDiscountRepositoryImpl()
		discountService := masterserviceimpl.StartDiscountService(discountRepository, db)
		discountController := mastercontroller.NewDiscountController(discountService)

		markupRateRepository := masteritemrepositoryimpl.StartMarkupRateRepositoryImpl()
		markupRateService := masteritemserviceimpl.StartMarkupRateService(markupRateRepository, db)
		markupRateController := masteritemcontroller.NewMarkupRateController(markupRateService)

		warehouseGroupRepository := masterwarehouserepositoryimpl.OpenWarehouseGroupImpl()
		warehouseGroupService := masterwarehouseserviceimpl.OpenWarehouseGroupService(warehouseGroupRepository, db)
		warehouseGroupController := masterwarehousecontroller.NewWarehouseGroupController(warehouseGroupService)

		warehouseLocationRepository := masterwarehouserepositoryimpl.OpenWarehouseLocationImpl()
		warehouseLocationService := masterwarehouseserviceimpl.OpenWarehouseLocationService(warehouseLocationRepository, db)
		warehouseLocationController := masterwarehousecontroller.NewWarehouseLocationController(warehouseLocationService)

		warehouseMasterRepository := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl()
		warehouseMasterService := masterwarehouseserviceimpl.OpenWarehouseMasterService(warehouseMasterRepository, db)
		warehouseMasterController := masterwarehousecontroller.NewWarehouseMasterController(warehouseMasterService)

		itemClassRouter := route.ItemClassRouter(itemClassController)
		OperationGroupRouter := route.OperationGroupRouter(operationGroupController)
		IncentiveGroupRouter := route.IncentiveGroupRouter(IncentiveGroupController)
		IncentiveGroupDetailRouter := route.IncentiveGroupDetailRouter(IncentiveGroupDetailController)
		OperationSectionRouter := route.OperationSectionRouter(operationSectionController)
		OperationEntriesRouter := route.OperationEntriesRouter(operationEntriesController)
		OperationKeyRouter := route.OperationKeyRouter(operationKeyController)
		ForecastMasterRouter := route.ForecastMasterRouter(forecastMasterController)
		DiscountPercentRouter := route.DiscountPercentRouter(discountPercentController)
		DiscountRouter := route.DiscountRouter(discountController)
		MarkupRateRouter := route.MarkupRateRouter(markupRateController)
		WarehouseGroup := route.WarehouseGroupRouter(warehouseGroupController)
		WarehouseLocation := route.WarehouseLocationRouter(warehouseLocationController)
		WarehouseMaster := route.WarehouseMasterRouter(warehouseMasterController)
		ShiftScheduleRouter := route.ShiftScheduleRouter(ShiftScheduleController)

		swaggerRouter := route.SwaggerRouter()
		mux := http.NewServeMux()

		mux.Handle("/item-class/", itemClassRouter)
		mux.Handle("/operation-group/", OperationGroupRouter)

		mux.Handle("/incentive-group/", IncentiveGroupRouter)
		mux.Handle("/incentive-group-detail/", IncentiveGroupDetailRouter)

		mux.Handle("/operation-section/", OperationSectionRouter)
		mux.Handle("/operation-key/", OperationKeyRouter)
		mux.Handle("/operation-entries/", OperationEntriesRouter)
		mux.Handle("/forecast-master/", ForecastMasterRouter)
		mux.Handle("/discount-percent/", DiscountPercentRouter)
		mux.Handle("/discount/", DiscountRouter)
		mux.Handle("/markup-rate/", MarkupRateRouter)
		mux.Handle("/warehouse-group/", WarehouseGroup)
		mux.Handle("/warehouse-location/", WarehouseLocation)
		mux.Handle("/warehouse-master/", WarehouseMaster)
		mux.Handle("/shift-schedule/", ShiftScheduleRouter)

		//Swagger
		mux.Handle("/swagger/", swaggerRouter)
		server := http.Server{
			Addr:    config.EnvConfigs.ClientOrigin,
			Handler: mux,
		}

		err := server.ListenAndServe()
		helper.PanicIfError(err)
	}
}
