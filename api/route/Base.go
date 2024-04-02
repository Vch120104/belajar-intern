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
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", warehouseGroupController.GetAllWarehouseGroup)
	router.Get("/by-id/{warehouse_group_id}", warehouseGroupController.GetByIdWarehouseGroup)
	router.Post("/", warehouseGroupController.SaveWarehouseGroup)
	router.Patch("/{warehouse_group_id}", warehouseGroupController.ChangeStatusWarehouseGroup)

	//router.PanicHandler = exceptions.ErrorHandler

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
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/item-substitute/", itemSubstituteController.GetAllItemSubstitute)
	router.GET("/item-substitute/header/by-id/:item_substitute_id", itemSubstituteController.GetByIdItemSubstitute)
	router.GET("/item-substitute/detail/all/by-id/:item_substitute_id", itemSubstituteController.GetAllItemSubstituteDetail)
	router.GET("/item-substitute/detail/by-id/:item_substitute_detail_id", itemSubstituteController.GetByIdItemSubstituteDetail)
	router.POST("/item-substitute/", itemSubstituteController.SaveItemSubstitute)
	router.POST("/item-substitute/detail/:item_substitute_id", itemSubstituteController.SaveItemSubstituteDetail)
	router.PATCH("/item-substitute/header/by-id/:item_substitute_id", itemSubstituteController.ChangeStatusItemSubstitute)
	router.PATCH("/item-substitute/detail/activate/by-id/", itemSubstituteController.ActivateItemSubstituteDetail)
	router.PATCH("/item-substitute/detail/deactivate/by-id/", itemSubstituteController.DeactivateItemSubstituteDetail)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func OperationCodeRouter(
	operationCodeController masteroperationcontroller.OperationCodeController,
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/operation-code/", operationCodeController.GetAllOperationCode)
	router.GET("/operation-code/by-id/:operation_id", operationCodeController.GetByIdOperationCode)
	router.POST("/operation-code/", operationCodeController.SaveOperationCode)
	router.PATCH("/operation-code/:operation_id", operationCodeController.ChangeStatusOperationCode)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func OperationGroupRouter(
	operationGroupController masteroperationcontroller.OperationGroupController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", operationGroupController.GetAllOperationGroup)
	router.Get("/drop-down", operationGroupController.GetAllOperationGroupIsActive)
	router.Get("/by-code/:operation_group_code", operationGroupController.GetOperationGroupByCode)
	router.Post("/", operationGroupController.SaveOperationGroup)
	router.Patch("/:operation_group_id", operationGroupController.ChangeStatusOperationGroup)

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
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/item-package/", ItemPackageController.GetAllItemPackage)
	router.POST("/item-package/", ItemPackageController.SaveItemPackage)
	router.GET("/item-package/by-id/:item_package_id", ItemPackageController.GetItemPackageById)
	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ItemPackageDetailRouter(
	ItemPackageDetailController masteritemcontroller.ItemPackageDetailController,
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/item-package-detail/by-package-id/:item_package_id", ItemPackageDetailController.GetItemPackageDetailByItemPackageId)

	router.PanicHandler = exceptions.ErrorHandler

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
) *httprouter.Router {
	router := httprouter.New()

	router.GET("/operation-section/", operationSectionController.GetAllOperationSectionList)
	router.GET("/operation-section/by-id/:operation_section_id", operationSectionController.GetOperationSectionByID)
	router.GET("/operation-section/by-name", operationSectionController.GetOperationSectionName)
	router.GET("/operation-section/code-by-group-id", operationSectionController.GetSectionCodeByGroupId)
	router.PUT("/operation-section/", operationSectionController.SaveOperationSection)
	router.PATCH("/operation-section/:operation_section_id", operationSectionController.ChangeStatusOperationSection)
	router.PanicHandler = exceptions.ErrorHandler
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

) *httprouter.Router {
	router := httprouter.New()

	router.GET("/operation-key/:operation_key_id", operationKeyController.GetOperationKeyByID)
	router.GET("/operation-key/", operationKeyController.GetAllOperationKeyList)
	router.GET("/operation-key-name/", operationKeyController.GetOperationKeyName)
	router.POST("/operation-key/", operationKeyController.SaveOperationKey)
	router.PATCH("/operation-key/:operation_key_id", operationKeyController.ChangeStatusOperationKey)
	router.PanicHandler = exceptions.ErrorHandler
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
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/item/", itemController.GetAllItem)
	router.GET("/item/pop-up/", itemController.GetAllItemLookup)
	router.GET("/item/multi-id/:item_ids", itemController.GetItemWithMultiId)
	router.GET("/item/by-code/:item_code", itemController.GetItemByCode)
	router.POST("/item/", itemController.SaveItem)
	router.PATCH("/item/:item_id", itemController.ChangeStatusItem)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func PriceListRouter(
	priceListController masteritemcontroller.PriceListController,
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/price-list/", priceListController.GetPriceList)
	router.GET("/price-list/pop-up/", priceListController.GetPriceListLookup)
	router.POST("/price-list/", priceListController.SavePriceList)
	router.PATCH("/price-list/:price_list_id", priceListController.ChangeStatusPriceList)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}

func WarrantyFreeServiceRouter(
	warrantyFreeServiceController mastercontroller.WarrantyFreeServiceController,
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/warranty-free-service/", warrantyFreeServiceController.GetAllWarrantyFreeService)
	router.GET("/warranty-free-service/:warranty_free_services_id", warrantyFreeServiceController.GetWarrantyFreeServiceByID)
	router.POST("/warranty-free-service/", warrantyFreeServiceController.SaveWarrantyFreeService)
	router.PATCH("/warranty-free-service/:warranty_free_services_id", warrantyFreeServiceController.ChangeStatusWarrantyFreeService)

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
