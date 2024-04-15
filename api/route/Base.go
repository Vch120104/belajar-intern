package route

import (
	mastercontroller "after-sales/api/controllers/master"
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	masterwarehousecontroller "after-sales/api/controllers/master/warehouse"
	"after-sales/api/exceptions"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Master */
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

func DiscountRouter(
	discountController mastercontroller.DiscountController,
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/discount/", discountController.GetAllDiscount)
	router.GET("/discount/drop-down/", discountController.GetAllDiscountIsActive)
	router.GET("/discount/by-code/:discount_code", discountController.GetDiscountByCode)
	router.POST("/discount/", discountController.SaveDiscount)
	router.PATCH("/discount/:discount_code_id", discountController.ChangeStatusDiscount)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func WarehouseMasterRouter(
	warehouseMasterController masterwarehousecontroller.WarehouseMasterController,
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/warehouse-master/", warehouseMasterController.GetAll)
	router.GET("/warehouse-master/by-id/:warehouse_id", warehouseMasterController.GetById)
	router.GET("/warehouse-master/by-code/:warehouse_code", warehouseMasterController.GetByCode)
	router.GET("/warehouse-master/multi-id/:warehouse_ids", warehouseMasterController.GetWarehouseWithMultiId)
	router.GET("/warehouse-master/drop-down", warehouseMasterController.GetAllIsActive)
	router.POST("/warehouse-master/", warehouseMasterController.Save)
	router.PATCH("/warehouse-master/:warehouse_id", warehouseMasterController.ChangeStatus)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func WarehouseGroupRouter(
	warehouseGroupController masterwarehousecontroller.WarehouseGroupController,
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/warehouse-group/", warehouseGroupController.GetAll)
	router.GET("/warehouse-group/:warehouse_group_id", warehouseGroupController.GetById)
	router.POST("/warehouse-group/", warehouseGroupController.Save)
	router.PATCH("/warehouse-group/:warehouse_group_id", warehouseGroupController.ChangeStatus)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func WarehouseLocationRouter(
	warehouseLocationController masterwarehousecontroller.WarehouseLocationController,
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/warehouse-location/", warehouseLocationController.GetAll)
	router.GET("/warehouse-location/:warehouse_location_id", warehouseLocationController.GetById)
	router.POST("/warehouse-location/", warehouseLocationController.Save)
	router.PATCH("/warehouse-location/:warehouse_location_id", warehouseLocationController.ChangeStatus)

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
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", itemSubstituteController.GetAllItemSubstitute)
	router.Get("/header/by-id/{item_substitute_id}", itemSubstituteController.GetByIdItemSubstitute)
	router.Get("/detail/all/by-id/{item_substitute_id}", itemSubstituteController.GetAllItemSubstituteDetail)
	router.Get("/detail/by-id/{item_substitute_detail_id}", itemSubstituteController.GetByIdItemSubstituteDetail)
	router.Post("/", itemSubstituteController.SaveItemSubstitute)
	router.Post("/detail/{item_substitute_id}", itemSubstituteController.SaveItemSubstituteDetail)
	router.Patch("/header/by-id/{item_substitute_id}", itemSubstituteController.ChangeStatusItemSubstitute)
	router.Patch("/detail/activate/by-id/{item_substitute_detail_id}", itemSubstituteController.ActivateItemSubstituteDetail)
	router.Patch("/detail/deactivate/by-id/{item_substitute_detail_id}", itemSubstituteController.DeactivateItemSubstituteDetail)

	return router
}

func OperationGroupRouter(
	operationGroupController masteroperationcontroller.OperationGroupController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", operationGroupController.GetAllOperationGroup)
	router.Get("/drop-down", operationGroupController.GetAllOperationGroupIsActive)
	router.Get("/by-code/{operation_group_code}", operationGroupController.GetOperationGroupByCode)
	router.Post("/", operationGroupController.SaveOperationGroup)
	router.Patch("/{operation_group_id}", operationGroupController.ChangeStatusOperationGroup)

	// router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ItemClassRouter(
	itemClassController masteritemcontroller.ItemClassController,
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/item-class/", itemClassController.GetAllItemClass)
	router.GET("/item-class/pop-up/", itemClassController.GetAllItemClassLookup)
	router.POST("/item-class/", itemClassController.SaveItemClass)
	router.PATCH("/item-class/:item_class_id", itemClassController.ChangeStatusItemClass)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ItemPackageRouter(
	ItemPackageController masteritemcontroller.ItemPackageController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", ItemPackageController.GetAllItemPackage)
	router.Post("/", ItemPackageController.SaveItemPackage)
	router.Get("/by-id/{item_package_id}", ItemPackageController.GetItemPackageById)
	router.Patch("/{item_package_id}", ItemPackageController.ChangeStatusItemPackage)

	return router
}

func ItemPackageDetailRouter(
	ItemPackageDetailController masteritemcontroller.ItemPackageDetailController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/by-package-id/{item_package_id}", ItemPackageDetailController.GetItemPackageDetailByItemPackageId)

	return router
}

func IncentiveGroupRouter(
	incentiveGroupController mastercontroller.IncentiveGroupController,
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/incentive-group/", incentiveGroupController.GetAllIncentiveGroup)
	router.GET("/incentive-group/drop-down/", incentiveGroupController.GetAllIncentiveGroupIsActive)
	router.GET("/incentive-group/by-id/:incentive_group_id", incentiveGroupController.GetIncentiveGroupById)
	router.POST("/incentive-group/", incentiveGroupController.SaveIncentiveGroup)
	router.PATCH("/incentive-group/:incentive_group_id", incentiveGroupController.ChangeStatusIncentiveGroup)

	router.PanicHandler = exceptions.ErrorHandler
	return router
}

func IncentiveGroupDetailRouter(
	incentiveGroupDetailController mastercontroller.IncentiveGroupDetailController,
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/incentive-group-detail/by-header-id/", incentiveGroupDetailController.GetAllIncentiveGroupDetail)
	router.GET("/incentive-group-detail/by-detail-id/:incentive_group_detail_id", incentiveGroupDetailController.GetIncentiveGroupDetailById)
	router.POST("/incentive-group-detail/", incentiveGroupDetailController.SaveIncentiveGroupDetail)

	router.PanicHandler = exceptions.ErrorHandler
	return router
}

func ShiftScheduleRouter(
	ShiftScheduleController mastercontroller.ShiftScheduleController,
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/shift-schedule/", ShiftScheduleController.GetAllShiftSchedule)
	// router.GET("/shift-schedule/drop-down", ShiftScheduleController.GetAllShiftScheduleIsActive)
	// router.GET("/shift-schedule/by-code/:operation_group_code", ShiftScheduleController.GetShiftScheduleByCode)
	router.POST("/shift-schedule/", ShiftScheduleController.SaveShiftSchedule)
	router.GET("/shift-schedule/:shift_schedule_id", ShiftScheduleController.GetShiftScheduleById)
	router.PATCH("/shift-schedule/:shift_schedule_id", ShiftScheduleController.ChangeStatusShiftSchedule)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func OperationSectionRouter(
	operationSectionController masteroperationcontroller.OperationSectionController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", operationSectionController.GetAllOperationSectionList)
	router.Get("/{operation_section_id}", operationSectionController.GetOperationSectionByID)
	router.Get("/by-name", operationSectionController.GetOperationSectionName)
	router.Get("/code-by-group-id/{operation_group_id}", operationSectionController.GetSectionCodeByGroupId)
	router.Post("/", operationSectionController.SaveOperationSection)
	router.Patch("/{operation_section_id}", operationSectionController.ChangeStatusOperationSection)
	return router
}

func OperationEntriesRouter(
	operationEntriesController masteroperationcontroller.OperationEntriesController,
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/operation-entries/", operationEntriesController.GetAllOperationEntries)
	router.GET("/operation-entries/:operation_entries_id", operationEntriesController.GetOperationEntriesByID)
	router.GET("/operation-entries-by-name/", operationEntriesController.GetOperationEntriesName)
	router.POST("/operation-entries/", operationEntriesController.SaveOperationEntries)
	router.PATCH("/operation-entries/:operation_entries_id", operationEntriesController.ChangeStatusOperationEntries)

	router.PanicHandler = exceptions.ErrorHandler
	return router
}

func OperationKeyRouter(
	operationKeyController masteroperationcontroller.OperationKeyController,

) chi.Router {
	router := chi.NewRouter()

	router.Get("/{operation_key_id}", operationKeyController.GetOperationKeyByID)
	router.Get("/", operationKeyController.GetAllOperationKeyList)
	router.Get("/operation-key-name", operationKeyController.GetOperationKeyName)
	router.Post("/", operationKeyController.SaveOperationKey)
	router.Patch("/{operation_key_id}", operationKeyController.ChangeStatusOperationKey)

	return router
}

func ForecastMasterRouter(
	forecastMasterController mastercontroller.ForecastMasterController,
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/forecast-master/", forecastMasterController.GetAllForecastMaster)
	router.GET("/forecast-master/:forecast_master_id", forecastMasterController.GetForecastMasterById)
	router.POST("/forecast-master/", forecastMasterController.SaveForecastMaster)
	router.PATCH("/forecast-master/:forecast_master_id", forecastMasterController.ChangeStatusForecastMaster)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func UnitOfMeasurementRouter(
	unitOfMeasurementController masteritemcontroller.UnitOfMeasurementController,
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/unit-of-measurement/", unitOfMeasurementController.GetAllUnitOfMeasurement)
	router.GET("/unit-of-measurement/drop-down", unitOfMeasurementController.GetAllUnitOfMeasurementIsActive)
	router.GET("/unit-of-measurement/by-code/:uom_code", unitOfMeasurementController.GetUnitOfMeasurementByCode)
	router.POST("/unit-of-measurement/", unitOfMeasurementController.SaveUnitOfMeasurement)
	router.PATCH("/unit-of-measurement/:uom_id", unitOfMeasurementController.ChangeStatusUnitOfMeasurement)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func MarkupMasterRouter(
	markupMasterController masteritemcontroller.MarkupMasterController,
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/markup-master/", markupMasterController.GetMarkupMasterList)
	router.GET("/markup-master/by-code/:markup_master_code", markupMasterController.GetMarkupMasterByCode)
	router.POST("/markup-master/", markupMasterController.SaveMarkupMaster)
	router.PATCH("/markup-master/:markup_master_id", markupMasterController.ChangeStatusMarkupMaster)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ItemLevelRouter(
	itemLevelController masteritemcontroller.ItemLevelController,
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/item-level/", itemLevelController.GetAll)
	router.GET("/item-level/:item_level_id", itemLevelController.GetById)
	router.POST("/item-level/", itemLevelController.Save)
	router.PATCH("/item-level/:item_level_id", itemLevelController.ChangeStatus)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ItemRouter(
	itemController masteritemcontroller.ItemController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", itemController.GetAllItem)
	router.Get("/pop-up/", itemController.GetAllItemLookup)
	router.Get("/multi-id/{item_id}", itemController.GetItemWithMultiId)
	router.Get("/by-code/{item_code}", itemController.GetItemByCode)
	router.Post("/", itemController.SaveItem)
	router.Patch("/{item_id}", itemController.ChangeStatusItem)

	return router
}

func PriceListRouter(
	priceListController masteritemcontroller.PriceListController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", priceListController.GetPriceList)
	router.Get("/pop-up", priceListController.GetPriceListLookup)
	router.Post("/", priceListController.SavePriceList)
	router.Patch("/{price_list_id}", priceListController.ChangeStatusPriceList)

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
