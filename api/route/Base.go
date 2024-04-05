package route

import (
	mastercontroller "after-sales/api/controllers/master"
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	masterwarehousecontroller "after-sales/api/controllers/master/warehouse"

	"github.com/go-chi/chi/v5"

	"net/http"
)

/* Master */
func DiscountPercentRouter(
	discountPercentController masteritemcontroller.DiscountPercentController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", discountPercentController.GetAllDiscountPercent)
	router.Get("/{discount_percent_id}", discountPercentController.GetDiscountPercentByID)
	router.Post("/", discountPercentController.SaveDiscountPercent)
	router.Patch("/{discount_percent_id}", discountPercentController.ChangeStatusDiscountPercent)

	//router.PanicHandler = exceptions.ErrorHandler

	return router
}

func DiscountRouter(
	discountController mastercontroller.DiscountController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", discountController.GetAllDiscount)
	router.Get("/drop-down", discountController.GetAllDiscountIsActive)
	router.Get("/by-code", discountController.GetDiscountByCode)
	router.Get("/by-id/{id}", discountController.GetDiscountById)
	router.Post("/", discountController.SaveDiscount)
	router.Patch("/{id}", discountController.ChangeStatusDiscount)

	return router
}

func WarehouseMasterRouter(
	warehouseMasterController masterwarehousecontroller.WarehouseMasterController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", warehouseMasterController.GetAll)
	router.Get("/by-id/{warehouse_id}", warehouseMasterController.GetById)
	router.Get("/by-code/{warehouse_code}", warehouseMasterController.GetByCode)
	router.Get("/multi-id/{warehouse_ids}", warehouseMasterController.GetWarehouseWithMultiId)
	router.Get("/drop-down", warehouseMasterController.GetAllIsActive)
	router.Post("/", warehouseMasterController.Save)
	router.Patch("/{warehouse_id}", warehouseMasterController.ChangeStatus)

	//router.PanicHandler = exceptions.ErrorHandler

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
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", warehouseLocationController.GetAll)
	router.Get("/by-id/{warehouse_location_id}", warehouseLocationController.GetById)
	router.Post("/", warehouseLocationController.Save)
	router.Patch("/{warehouse_location_id}", warehouseLocationController.ChangeStatus)

	// router.PanicHandler = exceptions.ErrorHandler

	return router
}

func MarkupRateRouter(
	markupRateController masteritemcontroller.MarkupRateController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", markupRateController.GetAllMarkupRate)
	router.Get("/{markup_rate_id}", markupRateController.GetMarkupRateByID)
	router.Post("/", markupRateController.SaveMarkupRate)
	router.Patch("/{markup_rate_id}", markupRateController.ChangeStatusMarkupRate)

	//router.PanicHandler = exceptions.ErrorHandler

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
	router.Patch("/detail/activate/by-id/", itemSubstituteController.ActivateItemSubstituteDetail)
	router.Patch("/detail/deactivate/by-id/", itemSubstituteController.DeactivateItemSubstituteDetail)

	//router.PanicHandler = exceptions.ErrorHandler

	return router
}

func OperationCodeRouter(
	operationCodeController masteroperationcontroller.OperationCodeController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", operationCodeController.GetAllOperationCode)
	router.Get("/by-id/{operation_id}", operationCodeController.GetByIdOperationCode)
	router.Post("/", operationCodeController.SaveOperationCode)
	router.Patch("/{operation_id}", operationCodeController.ChangeStatusOperationCode)

	// router.PanicHandler = exceptions.ErrorHandler

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

	// //router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ItemClassRouter(
	itemClassController masteritemcontroller.ItemClassController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", itemClassController.GetAllItemClass)
	router.Get("/pop-up", itemClassController.GetAllItemClassLookup)
	router.Post("/", itemClassController.SaveItemClass)
	router.Patch("/{item_class_id}", itemClassController.ChangeStatusItemClass)

	// //router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ItemPackageRouter(
	ItemPackageController masteritemcontroller.ItemPackageController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", ItemPackageController.GetAllItemPackage)
	router.Post("/", ItemPackageController.SaveItemPackage)
	router.Get("/by-id/{item_package_id}", ItemPackageController.GetItemPackageById)
	//router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ItemPackageDetailRouter(
	ItemPackageDetailController masteritemcontroller.ItemPackageDetailController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/item-package-detail/by-package-id/:item_package_id", ItemPackageDetailController.GetItemPackageDetailByItemPackageId)

	//router.PanicHandler = exceptions.ErrorHandler

	return router
}

func IncentiveGroupRouter(
	incentiveGroupController mastercontroller.IncentiveGroupController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", incentiveGroupController.GetAllIncentiveGroup)
	router.Get("/drop-down", incentiveGroupController.GetAllIncentiveGroupIsActive)
	router.Get("/by-id/{id}", incentiveGroupController.GetIncentiveGroupById)
	router.Post("/", incentiveGroupController.SaveIncentiveGroup)
	router.Patch("/{id}", incentiveGroupController.ChangeStatusIncentiveGroup)

	// router.PanicHandler = exceptions.ErrorHandler
	return router
}

func DeductionRouter(
	DeductionController mastercontroller.DeductionController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", DeductionController.GetAllDeductionList)
	router.Get("/{id}", DeductionController.GetAllDeductionDetail)
	router.Get("/by-detail-id/{id}", DeductionController.GetByIdDeductionDetail)
	router.Get("/by-header-id/{id}", DeductionController.GetDeductionById)
	router.Post("/detail", DeductionController.SaveDeductionDetail)
	router.Post("/", DeductionController.SaveDeductionList)
	router.Patch("/{id}", DeductionController.ChangeStatusDeduction)

	// router.PanicHandler = exceptions.ErrorHandler
	return router
}

func IncentiveGroupDetailRouter(
	incentiveGroupDetailController mastercontroller.IncentiveGroupDetailController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/{id}", incentiveGroupDetailController.GetAllIncentiveGroupDetail)
	router.Get("/by-id/{incentive_group_detail_id}", incentiveGroupDetailController.GetIncentiveGroupDetailById)
	router.Post("/", incentiveGroupDetailController.SaveIncentiveGroupDetail)

	// router.PanicHandler = exceptions.ErrorHandler
	return router
}

func IncentiveMasterRouter(
	IncentiveMasterController mastercontroller.IncentiveMasterController,
) chi.Router {
	router := chi.NewRouter()
	// Gunakan middleware NotFoundHandler
	// router.Use(middleware.NotFoundHandler)

	router.Get("/", http.HandlerFunc(IncentiveMasterController.GetAllIncentiveMaster))
	router.Get("/{incentive_level_id}", http.HandlerFunc(IncentiveMasterController.GetIncentiveMasterById))
	router.Post("/", http.HandlerFunc(IncentiveMasterController.SaveIncentiveMaster))
	router.Patch("/{incentive_level_id}", http.HandlerFunc(IncentiveMasterController.ChangeStatusIncentiveMaster))

	////router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ShiftScheduleRouter(
	ShiftScheduleController mastercontroller.ShiftScheduleController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", ShiftScheduleController.GetAllShiftSchedule)
	// router.Get("/drop-down", ShiftScheduleController.GetAllShiftScheduleIsActive)
	// router.Get("/by-code/{operation_group_code}", ShiftScheduleController.GetShiftScheduleByCode)
	router.Post("/", ShiftScheduleController.SaveShiftSchedule)
	router.Get("/by-id/{shift_schedule_id}", ShiftScheduleController.GetShiftScheduleById)
	router.Patch("/{shift_schedule_id}", ShiftScheduleController.ChangeStatusShiftSchedule)

	// router.PanicHandler = exceptions.ErrorHandler

	return router
}

func OperationSectionRouter(
	operationSectionController masteroperationcontroller.OperationSectionController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/", operationSectionController.GetAllOperationSectionList)
	router.Get("/by-id/{operation_section_id}", operationSectionController.GetOperationSectionByID)
	router.Get("/by-name", operationSectionController.GetOperationSectionName)
	router.Get("/code-by-group-id", operationSectionController.GetSectionCodeByGroupId)
	router.Put("/", operationSectionController.SaveOperationSection)
	router.Patch("/{operation_section_id}", operationSectionController.ChangeStatusOperationSection)
	//router.PanicHandler = exceptions.ErrorHandler
	return router
}

func OperationEntriesRouter(
	operationEntriesController masteroperationcontroller.OperationEntriesController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", operationEntriesController.GetAllOperationEntries)
	router.Get("/by-id/{operation_entries_id}", operationEntriesController.GetOperationEntriesByID)
	router.Get("/by-name", operationEntriesController.GetOperationEntriesName)
	router.Post("/", operationEntriesController.SaveOperationEntries)
	router.Patch("/{operation_entries_id}", operationEntriesController.ChangeStatusOperationEntries)

	return router
}

func OperationKeyRouter(
	operationKeyController masteroperationcontroller.OperationKeyController,

) chi.Router {
	router := chi.NewRouter()

	router.Get("/{operation_key_id}", operationKeyController.GetOperationKeyByID)
	router.Get("/", operationKeyController.GetAllOperationKeyList)
	router.Get("/operation-key-name/", operationKeyController.GetOperationKeyName)
	router.Post("/", operationKeyController.SaveOperationKey)
	router.Patch("/{operation_key_id}", operationKeyController.ChangeStatusOperationKey)
	//router.PanicHandler = exceptions.ErrorHandler
	return router
}

func ForecastMasterRouter(
	forecastMasterController mastercontroller.ForecastMasterController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", forecastMasterController.GetAllForecastMaster)
	router.Get("/by-id/{forecast_master_id}", forecastMasterController.GetForecastMasterById)
	router.Post("/", forecastMasterController.SaveForecastMaster)
	router.Patch("/{forecast_master_id}", forecastMasterController.ChangeStatusForecastMaster)

	// router.PanicHandler = exceptions.ErrorHandler

	return router
}

func UnitOfMeasurementRouter(
	unitOfMeasurementController masteritemcontroller.UnitOfMeasurementController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", unitOfMeasurementController.GetAllUnitOfMeasurement)
	router.Get("/drop-down", unitOfMeasurementController.GetAllUnitOfMeasurementIsActive)
	router.Get("/code/{uom_code}", unitOfMeasurementController.GetUnitOfMeasurementByCode)
	router.Post("/", unitOfMeasurementController.SaveUnitOfMeasurement)
	router.Patch("/{uom_id}", unitOfMeasurementController.ChangeStatusUnitOfMeasurement)

	// //router.PanicHandler = exceptions.ErrorHandler

	return router
}

func MarkupMasterRouter(
	markupMasterController masteritemcontroller.MarkupMasterController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", markupMasterController.GetMarkupMasterList)
	router.Get("/code/{markup_master_code}", markupMasterController.GetMarkupMasterByCode)
	router.Post("/", markupMasterController.SaveMarkupMaster)
	router.Patch("/{markup_master_id}", markupMasterController.ChangeStatusMarkupMaster)

	// router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ItemLevelRouter(
	itemLevelController masteritemcontroller.ItemLevelController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", itemLevelController.GetAll)
	router.Get("/{item_level_id}", itemLevelController.GetById)
	router.Post("/", itemLevelController.Save)
	router.Patch("/{item_level_id}", itemLevelController.ChangeStatus)

	// router.PanicHandler = exceptions.ErrorHandler

	return router
}

func ItemRouter(
	itemController masteritemcontroller.ItemController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", itemController.GetAllItem)
	router.Get("/pop-up/", itemController.GetAllItemLookup)
	router.Get("/multi-id/{item_ids}", itemController.GetItemWithMultiId)
	router.Get("/by-code/{item_code}", itemController.GetItemByCode)
	router.Post("/", itemController.SaveItem)
	router.Patch("/{item_id}", itemController.ChangeStatusItem)

	//router.PanicHandler = exceptions.ErrorHandler

	return router
}

func PriceListRouter(
	priceListController masteritemcontroller.PriceListController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", priceListController.GetPriceList)
	router.Get("/pop-up/", priceListController.GetPriceListLookup)
	router.Post("/", priceListController.SavePriceList)
	router.Patch("/{price_list_id}", priceListController.ChangeStatusPriceList)

	//router.PanicHandler = exceptions.ErrorHandler

	return router
}

func WarrantyFreeServiceRouter(
	warrantyFreeServiceController mastercontroller.WarrantyFreeServiceController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", warrantyFreeServiceController.GetAllWarrantyFreeService)
	router.Get("/{warranty_free_services_id}", warrantyFreeServiceController.GetWarrantyFreeServiceByID)
	router.Post("/", warrantyFreeServiceController.SaveWarrantyFreeService)
	router.Patch("/{warranty_free_services_id}", warrantyFreeServiceController.ChangeStatusWarrantyFreeService)

	// router.PanicHandler = exceptions.ErrorHandler

	return router
}

func BomRouter(
	BomController masteritemcontroller.BomController,
) chi.Router {
	router := chi.NewRouter()
	// Gunakan middleware NotFoundHandler

	//bom master
	router.Get("/", BomController.GetBomMasterList)
	router.Get("/{bom_master_id}", BomController.GetBomMasterById)
	router.Post("/", BomController.SaveBomMaster)
	router.Patch("/{bom_master_id}", BomController.ChangeStatusBomMaster)

	//bom detail
	router.Get("/all/detail", BomController.GetBomDetailList)
	router.Get("/{bom_master_id}/detail", BomController.GetBomDetailById)
	router.Post("/all/detail", BomController.SaveBomDetail)
	//router.Put("/all/detail", BomController.SubmitBomDetail)

	//router.Delete("/{bom_detail_id}/detail", BomController.SaveBomDetail)

	//bom lookup
	router.Get("/{bom_master_id}/popup-item", BomController.GetBomItemList)

	////router.PanicHandler = exceptions.ErrorHandler

	return router
}

// func SwaggerRouter() chi.Router {
// 	router := chi.NewRouter()
// 	router.Get("/swagger/*any", adaptHandler(swaggerHandler()))
// 	return router
// }

// func adaptHandler(h http.Handler) chi.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 		h.ServeHTTP(w, r)
// 	}
// }

// func swaggerHandler() http.HandlerFunc {
// 	return httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json"))
// }
