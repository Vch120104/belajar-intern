package route

import (
	mastercontroller "after-sales/api/controllers/master"
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	"after-sales/api/exceptions"

	_ "after-sales/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"net/http"

	"github.com/julienschmidt/httprouter"
)

func DiscountPercentRouter(
	discountPercentController masteritemcontroller.DiscountPercentController,
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/discount-percent/", discountPercentController.GetAllDiscountPercent)
	router.GET("/discount-percent/:discount_percent_id", discountPercentController.GetDiscountPercentByID)
	router.POST("/discount-percent/", discountPercentController.SaveDiscountPercent)
	router.PATCH("/discount-percent/:discount_percent_id", discountPercentController.ChangeStatusDiscountPercent)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func MarkupRateRouter(
	markupRateController masteritemcontroller.MarkupRateController,
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/markup-rate/", markupRateController.GetAllMarkupRate)
	router.GET("/markup-rate/:markup_rate_id", markupRateController.GetMarkupRateByID)
	router.POST("/markup-rate/", markupRateController.SaveMarkupRate)
	router.PATCH("/markup-rate/:markup_rate_id", markupRateController.ChangeStatusMarkupRate)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ItemSubstituteRouter(
	itemSubstituteController masteritemcontroller.ItemSubstituteController,
) *httprouter.Router{
	router := httprouter.New()
	router.GET("/item-substitute/",itemSubstituteController.GetAllItemSubstitute)
	router.GET("/item-substitute/header/by-id/:item_substitute_id",itemSubstituteController.GetByIdItemSubstitute)
	router.GET("/item-substitute/detail/all/by-id/:item_substitute_id",itemSubstituteController.GetAllItemSubstituteDetail)
	router.GET("/item-substitute/detail/by-id/:item_substitute_detail_id",itemSubstituteController.GetByIdItemSubstituteDetail)
	router.POST("/item-substitute/",itemSubstituteController.SaveItemSubstitute)
	router.POST("/item-substitute/detail/:item_substitute_id",itemSubstituteController.SaveItemSubstituteDetail)
	router.PATCH("/item-substitute/header/by-id/:item_substitute_id",itemSubstituteController.ChangeStatusItemSubstitute)
	router.PATCH("/item-substitute/detail/activate/by-id/",itemSubstituteController.ActivateItemSubstituteDetail)
	router.PATCH("/item-substitute/detail/deactivate/by-id/",itemSubstituteController.DeactivateItemSubstituteDetail)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func OperationGroupRouter(
	operationGroupController masteroperationcontroller.OperationGroupController,
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/operation-group/", operationGroupController.GetAllOperationGroup)
	router.GET("/operation-group/drop-down", operationGroupController.GetAllOperationGroupIsActive)
	router.GET("/operation-group/by-code/:operation_group_code", operationGroupController.GetOperationGroupByCode)
	router.POST("/operation-group/", operationGroupController.SaveOperationGroup)
	router.PATCH("/operation-group/:operation_group_id", operationGroupController.ChangeStatusOperationGroup)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ForecastMasterRouter(
	forecastMasterController mastercontroller.ForecastMasterController,
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/forecast-master", forecastMasterController.GetForecastMasterById)
	router.PanicHandler = exceptions.ErrorHandler

	return router
}
func SwaggerRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/swagger/*any", adaptHandler(swaggerHandler()))
	return router
}

func adaptHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}

func swaggerHandler() http.HandlerFunc {
	return httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json"))
}

// type OperationGroupController interface {
// 	GetAllOperationGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
// 	GetAllOperationGroupIsActive(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
// 	GetOperationGroupByCode(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
// 	SaveOperationGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
// 	ChangeStatusOperationGroup(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
// }

// func CreateHandler(db *gorm.DB, env string, redis *redis.Client) {
// 	r := gin.New()
// 	//mtr_operation_group
// 	operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl(db)
// 	operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository)
// 	//mtr_operation_section
// 	operationSectionRepository := masteroperationrepositoryimpl.StartOperationSectionRepositoryImpl(db)
// 	operationSectionService := masteroperationserviceimpl.StartOperationSectionService(operationSectionRepository)
// 	//mtr_operation_key
// 	operationKeyRepository := masteroperationrepositoryimpl.StartOperationKeyRepositoryImpl(db)
// 	operationKeyService := masteroperationserviceimpl.StartOperationKeyService(operationKeyRepository)
// 	// //mtr_operation_code
// 	operationCodeRepository := masteroperationrepositoryimpl.StartOperationCodeRepositoryImpl(db)
// 	operationCodeService := masteroperationserviceimpl.StartOperationCodeService(operationCodeRepository)
// 	// //mtr_operation_entries
// 	operationEntriesRepository := masteroperationrepositoryimpl.StartOperationEntriesRepositoryImpl(db)
// 	operationEntriesService := masteroperationserviceimpl.StartOperationEntriesService(operationEntriesRepository)
// 	// //mtr_operation_model_mapping
// 	operationModelMappingRepository := masteroperationrepositoryimpl.StartOperationModelMappingRepositoryImpl(db)
// 	operationModelMappingService := masteroperationserviceimpl.StartOperationMappingService(operationModelMappingRepository)

// 	//mtr_markup_master
// 	markupMasterRepository := masteritemrepositoryimpl.StartMarkupMasterRepositoryImpl(db)
// 	markupMasterService := masteritemserviceimpl.StartMarkupMasterService(markupMasterRepository)

// 	//mtr_uom
// 	UnitOfMeasurementRepository := masteritemrepositoryimpl.StartUnitOfMeasurementRepositoryImpl(db)
// 	UnitOfMeasurementService := masteritemserviceimpl.StartUnitOfMeasurementService(UnitOfMeasurementRepository)
// 	//mtr_item
// 	itemRepository := masteritemrepositoryimpl.StartItemRepositoryImpl(db, redis)
// 	itemService := masteritemserviceimpl.StartItemService(itemRepository)
// 	//mtr_discount
// 	discountRepository := masterrepositoryimpl.StartDiscountRepositoryImpl(db)
// 	discountService := masterserviceimpl.StartDiscountService(discountRepository)
// 	//mtr_incentive_group
// 	incentiveGroupRepository := masterrepositoryimpl.StartIncentiveGroupImpl(db)
// 	incentiveGroupService := masterserviceimpl.StartIncentiveGroup(incentiveGroupRepository)
// 	//mtr_item_class
// 	itemClassRepository := masteritemrepositoryimpl.StartItemClassRepositoryImpl(db)
// 	itemClassService := masteritemserviceimpl.StartItemClassService(itemClassRepository)
// 	//mtr_price_list
// 	priceListRepository := masteritemrepositoryimpl.StartPriceListRepositoryImpl(db)
// 	priceListService := masteritemserviceimpl.StartPriceListService(priceListRepository)

// 	//mtr_discount_percent
// 	discountPercentRepository := masteritemrepositoryimpl.StartDiscountPercentRepositoryImpl(db)
// 	discountPercentService := masteritemserviceimpl.StartDiscountPercentService(discountPercentRepository)

// 	//mtr_item_level
// 	itemLevelRepository := masteritemrepositoryimpl.StartItemLevelRepositoryImpl(db)
// 	itemLevelService := masteritemserviceimpl.StartItemLevelService(itemLevelRepository)

// 	warehouseGroupRepository := masterwarehouserepositoryimpl.OpenWarehouseGroupImpl(db)
// 	warehouseGroupService := masterwarehouseserviceimpl.OpenWarehouseGroupService(warehouseGroupRepository)

// 	warehouseMasterRepository := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl(db)
// 	warehouseMasterService := masterwarehouseserviceimpl.OpenWarehouseMasterService(warehouseMasterRepository)

// 	warehouseLocationRepository := masterwarehouserepositoryimpl.OpenWarehouseLocationImpl(db)
// 	warehouseLocationService := masterwarehouseserviceimpl.OpenWarehouseLocationService(warehouseLocationRepository)

// 	bookingEstimationRepository := transactionworkshoprepositoryimpl.OpenBookingEstimationImpl(db)
// 	bookingEstimationService := transactionworkshopserviceimpl.OpenBookingEstimationServiceImpl(bookingEstimationRepository)

// 	r.Use(middlewares.SetupCorsMiddleware())

// 	if env != "prod" {
// 		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
// 	}

// 	api := r.Group("aftersales-service/api/aftersales")
// 	//master
// 	mastercontroller.StartDiscountRoutes(db, api, discountService)
// 	mastercontroller.StartIncentiveGroupRoutes(db, api, incentiveGroupService)
// 	//operation
// 	masteroperationcontroller.StartOperationGroupRoutes(db, api, operationGroupService)
// 	masteroperationcontroller.StartOperationSectionRoutes(db, api, operationSectionService)
// 	masteroperationcontroller.StartOperationKeyRoutes(db, api, operationKeyService)
// 	masteroperationcontroller.StartOperationCodeRoutes(db, api, operationCodeService)
// 	masteroperationcontroller.StartOperationEntriesRoutes(db, api, operationEntriesService)
// 	masteroperationcontroller.StartOperationModelMappingRoutes(db, api, operationModelMappingService)
// 	masteritemcontroller.StartUnitOfMeasurementRoutes(db, api, UnitOfMeasurementService)
// 	masteritemcontroller.StartItemRoutes(db, api, itemService)
// 	masteritemcontroller.StartItemLevelRoutes(db, api, itemLevelService)
// 	masteritemcontroller.StartItemClassRoutes(db, api, itemClassService)
// 	masteritemcontroller.StartPriceListRoutes(db, api, priceListService)
// 	masteritemcontroller.StartMarkupMasterRoutes(db, api, markupMasterService)
// 	masteritemcontroller.StartDiscountPercentRoutes(db, api, discountPercentService)

// 	masterwarehousecontroller.OpenWarehouseGroupRoutes(db, api, warehouseGroupService)
// 	masterwarehousecontroller.OpenWarehouseMasterRoutes(db, api, warehouseMasterService)
// 	masterwarehousecontroller.OpenWarehouseLocationRoutes(db, api, warehouseLocationService)
// 	// //transaction
// 	// transactioncontroller.StartSupplySlipRoutes(db, api.Group("/transaction"), supplySlipService)

// 	transactionworkshopcontroller.OpenBookingEstimationRoutes(db, api, bookingEstimationService)

// 	server := &http.Server{Handler: r}
// 	l, err := net.Listen("tcp4", fmt.Sprintf(":%v", config.EnvConfigs.Port))
// 	err = server.Serve(l)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// }
