package route

import (
	mastercontroller "after-sales/api/controllers/master"
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	masterwarehousecontroller "after-sales/api/controllers/master/warehouse"
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

func DiscountRouter(

	discountController mastercontroller.DiscountController,
) *httprouter.Router {
	router := httprouter.New()
	router.GET("/discount/", discountController.GetAllDiscount)
	router.GET("/discount-drop-down/", discountController.GetAllDiscountIsActive)
	router.GET("/discount-by-code/:discount_code", discountController.GetDiscountByCode)
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
