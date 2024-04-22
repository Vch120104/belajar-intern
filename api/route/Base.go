package route

import (
	mastercontroller "after-sales/api/controllers/master"
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	masterwarehousecontroller "after-sales/api/controllers/master/warehouse"
	"after-sales/api/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

/* Master */

func ItemClassRouter(
	itemClassController masteritemcontroller.ItemClassController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", itemClassController.GetAllItemClass)
	router.Get("/pop-up", itemClassController.GetAllItemClassLookup)
	router.Post("/", itemClassController.SaveItemClass)
	router.Patch("/{item_class_id}", itemClassController.ChangeStatusItemClass)

	return router
}

func UnitOfMeasurementRouter(
	unitOfMeasurementController masteritemcontroller.UnitOfMeasurementController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", unitOfMeasurementController.GetAllUnitOfMeasurement)
	router.Get("/drop-down", unitOfMeasurementController.GetAllUnitOfMeasurementIsActive)
	router.Get("/code/{uom_code}", unitOfMeasurementController.GetUnitOfMeasurementByCode)
	router.Post("/", unitOfMeasurementController.SaveUnitOfMeasurement)
	router.Patch("/{uom_id}", unitOfMeasurementController.ChangeStatusUnitOfMeasurement)

	return router
}

func DiscountPercentRouter(
	discountPercentController masteritemcontroller.DiscountPercentController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", discountPercentController.GetAllDiscountPercent)
	router.Get("/{discount_percent_id}", discountPercentController.GetDiscountPercentByID)
	router.Post("/", discountPercentController.SaveDiscountPercent)
	router.Patch("/{discount_percent_id}", discountPercentController.ChangeStatusDiscountPercent)

	return router
}

func MarkupMasterRouter(
	markupMasterController masteritemcontroller.MarkupMasterController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", markupMasterController.GetMarkupMasterList)
	router.Get("/code/{markup_master_code}", markupMasterController.GetMarkupMasterByCode)
	router.Post("/", markupMasterController.SaveMarkupMaster)
	router.Patch("/{markup_master_id}", markupMasterController.ChangeStatusMarkupMaster)

	return router
}

func MarkupRateRouter(
	markupRateController masteritemcontroller.MarkupRateController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", markupRateController.GetAllMarkupRate)
	router.Get("/{markup_rate_id}", markupRateController.GetMarkupRateByID)
	router.Post("/", markupRateController.SaveMarkupRate)
	router.Patch("/{markup_rate_id}", markupRateController.ChangeStatusMarkupRate)

	return router
}

func ItemLevelRouter(
	itemLevelController masteritemcontroller.ItemLevelController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", itemLevelController.GetAll)
	router.Get("/{item_level_id}", itemLevelController.GetById)
	router.Post("/", itemLevelController.Save)
	router.Patch("/{item_level_id}", itemLevelController.ChangeStatus)

	return router
}

func ItemRouter(
	itemController masteritemcontroller.ItemController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", itemController.GetAllItem)
	router.Get("/pop-up", itemController.GetAllItemLookup)
	router.Get("/multi-id/{item_ids}", itemController.GetItemWithMultiId)
	router.Get("/by-code/{item_code}", itemController.GetItemByCode)
	router.Post("/", itemController.SaveItem)
	router.Patch("/{item_id}", itemController.ChangeStatusItem)

	return router
}

func ItemLocationRouter(
	ItemLocationController masteritemcontroller.ItemLocationController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	//master
	router.Get("/", ItemLocationController.GetAllItemLocation)
	router.Get("/{item_location_id}", ItemLocationController.GetItemLocationById)
	router.Post("/", ItemLocationController.SaveItemLocation)

	//detail
	router.Get("/all/detail", ItemLocationController.GetAllItemLocationDetail)
	router.Get("/popup-location", ItemLocationController.PopupItemLocation)
	router.Post("/all/detail", ItemLocationController.AddItemLocation)
	router.Delete("/all/detail/{item_location_detail_id}", ItemLocationController.DeleteItemLocation)

	return router
}

func ItemSubstituteRouter(
	itemSubstituteController masteritemcontroller.ItemSubstituteController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", itemSubstituteController.GetAllItemSubstitute)
	router.Get("/header/by-id/{item_substitute_id}", itemSubstituteController.GetByIdItemSubstitute)
	router.Get("/detail/all/by-id/{item_substitute_id}", itemSubstituteController.GetAllItemSubstituteDetail)
	router.Get("/detail/by-id/{item_substitute_detail_id}", itemSubstituteController.GetByIdItemSubstituteDetail)
	router.Post("/", itemSubstituteController.SaveItemSubstitute)
	router.Post("/detail/{item_substitute_id}", itemSubstituteController.SaveItemSubstituteDetail)
	router.Patch("/header/by-id/{item_substitute_id}", itemSubstituteController.ChangeStatusItemSubstitute)
	router.Patch("/detail/activate/by-id/", itemSubstituteController.ActivateItemSubstituteDetail)
	router.Patch("/detail/deactivate/by-id/", itemSubstituteController.DeactivateItemSubstituteDetail)

	return router
}

func ItemPackageRouter(
	ItemPackageController masteritemcontroller.ItemPackageController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", ItemPackageController.GetAllItemPackage)
	router.Post("/", ItemPackageController.SaveItemPackage)
	router.Get("/by-id/{item_package_id}", ItemPackageController.GetItemPackageById)

	return router
}

func ItemPackageDetailRouter(
	ItemPackageDetailController masteritemcontroller.ItemPackageDetailController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/by-package-id/{item_package_id}", ItemPackageDetailController.GetItemPackageDetailByItemPackageId)

	return router
}

func PriceListRouter(
	priceListController masteritemcontroller.PriceListController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", priceListController.GetPriceList)
	router.Get("/pop-up/", priceListController.GetPriceListLookup)
	router.Post("/", priceListController.SavePriceList)
	router.Patch("/{price_list_id}", priceListController.ChangeStatusPriceList)

	return router
}

func BomRouter(
	BomController masteritemcontroller.BomController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	//bom master
	router.Get("/", BomController.GetBomMasterList)
	router.Get("/{bom_master_id}", BomController.GetBomMasterById)
	router.Post("/", BomController.SaveBomMaster)
	router.Patch("/{bom_master_id}", BomController.ChangeStatusBomMaster)

	//bom detail
	router.Get("/all/detail", BomController.GetBomDetailList)
	router.Get("/{bom_master_id}/detail", BomController.GetBomDetailById)
	router.Get("/all/detail/{bom_detail_id}", BomController.GetBomDetailByIds)
	router.Post("/all/detail", BomController.SaveBomDetail)
	router.Delete("/all/detail/{bom_detail_id}", BomController.DeleteBomDetail)

	//bom lookup
	router.Get("/popup-item", BomController.GetBomItemList)

	return router
}

func LandedCostMasterRouter(
	LandedCostMaster masteritemcontroller.LandedCostMasterController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", LandedCostMaster.GetAllLandedCostMaster)
	router.Get("/{landed_cost_id}", LandedCostMaster.GetByIdLandedCost)
	router.Post("/", LandedCostMaster.SaveLandedCostMaster)
	router.Patch("/activate/", LandedCostMaster.ActivateLandedCostMaster)
	router.Patch("/deactivate/", LandedCostMaster.DeactivateLandedCostmaster)

	return router
}

func OperationGroupRouter(
	operationGroupController masteroperationcontroller.OperationGroupController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", operationGroupController.GetAllOperationGroup)
	router.Get("/drop-down", operationGroupController.GetAllOperationGroupIsActive)
	router.Get("/by-code/{operation_group_code}", operationGroupController.GetOperationGroupByCode)
	router.Post("/", operationGroupController.SaveOperationGroup)
	router.Patch("/{operation_group_id}", operationGroupController.ChangeStatusOperationGroup)

	return router
}

func OperationSectionRouter(
	operationSectionController masteroperationcontroller.OperationSectionController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", operationSectionController.GetAllOperationSectionList)
	router.Get("/by-id/{operation_section_id}", operationSectionController.GetOperationSectionByID)
	router.Get("/by-name", operationSectionController.GetOperationSectionName)
	router.Get("/code-by-group-id", operationSectionController.GetSectionCodeByGroupId)
	router.Put("/", operationSectionController.SaveOperationSection)
	router.Patch("/{operation_section_id}", operationSectionController.ChangeStatusOperationSection)

	return router
}

func OperationKeyRouter(
	operationKeyController masteroperationcontroller.OperationKeyController,

) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/{operation_key_id}", operationKeyController.GetOperationKeyByID)
	router.Get("/", operationKeyController.GetAllOperationKeyList)
	router.Get("/operation-key-name/", operationKeyController.GetOperationKeyName)
	router.Post("/", operationKeyController.SaveOperationKey)
	router.Patch("/{operation_key_id}", operationKeyController.ChangeStatusOperationKey)

	return router
}

func OperationEntriesRouter(
	operationEntriesController masteroperationcontroller.OperationEntriesController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", operationEntriesController.GetAllOperationEntries)
	router.Get("/by-id/{operation_entries_id}", operationEntriesController.GetOperationEntriesByID)
	router.Get("/by-name", operationEntriesController.GetOperationEntriesName)
	router.Post("/", operationEntriesController.SaveOperationEntries)
	router.Patch("/{operation_entries_id}", operationEntriesController.ChangeStatusOperationEntries)

	return router
}

func OperationCodeRouter(
	operationCodeController masteroperationcontroller.OperationCodeController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", operationCodeController.GetAllOperationCode)
	router.Get("/by-id/{operation_id}", operationCodeController.GetByIdOperationCode)
	router.Post("/", operationCodeController.SaveOperationCode)
	router.Patch("/{operation_id}", operationCodeController.ChangeStatusOperationCode)

	return router
}

func OperationModelMappingRouter(
	operationModelMappingController masteroperationcontroller.OperationModelMappingController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", operationModelMappingController.GetOperationModelMappingLookup)
	router.Get("/operation-frt/{operation_model_mapping_id}", operationModelMappingController.GetAllOperationFrt)
	router.Get("/operation-document-requirement/{operation_model_mapping_id}", operationModelMappingController.GetAllOperationDocumentRequirement)
	router.Get("/by-id/{operation_model_mapping_id}", operationModelMappingController.GetOperationModelMappingById)
	router.Get("/operation-frt/by-id/{operation_frt_id}", operationModelMappingController.GetOperationFrtById)
	router.Get("/operation-document-requirement/by-id/{operation_document_requirement_id}", operationModelMappingController.GetOperationDocumentRequirementById)
	router.Post("/", operationModelMappingController.SaveOperationModelMapping)
	router.Post("/operation-frt", operationModelMappingController.SaveOperationModelMappingFrt)
	router.Post("/operation-document-requirement", operationModelMappingController.SaveOperationModelMappingDocumentRequirement)
	router.Patch("/{operation_model_mapping_id}", operationModelMappingController.ChangeStatusOperationModelMapping)
	router.Patch("/operation-frt/activate/{operation_frt_id}", operationModelMappingController.ActivateOperationFrt)
	router.Patch("/operation-frt/deactivate/{operation_frt_id}", operationModelMappingController.DeactivateOperationFrt)
	router.Patch("/operation-document-requirement/deactivate/{operation_model_mapping_id}", operationModelMappingController.DeactivateOperationDocumentRequirement)
	router.Patch("/operation-document-requirement/activate/{operation_model_mapping_id}", operationModelMappingController.ActivateOperationDocumentRequirement)

	return router
}

func WarehouseGroupRouter(
	warehouseGroupController masterwarehousecontroller.WarehouseGroupController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", warehouseGroupController.GetAllWarehouseGroup)
	router.Get("/by-id/{warehouse_group_id}", warehouseGroupController.GetByIdWarehouseGroup)
	router.Post("/", warehouseGroupController.SaveWarehouseGroup)
	router.Patch("/{warehouse_group_id}", warehouseGroupController.ChangeStatusWarehouseGroup)

	return router
}

func WarehouseMasterRouter(
	warehouseMasterController masterwarehousecontroller.WarehouseMasterController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", warehouseMasterController.GetAll)
	router.Get("/by-id/{warehouse_id}", warehouseMasterController.GetById)
	router.Get("/by-code/{warehouse_code}", warehouseMasterController.GetByCode)
	router.Get("/multi-id/{warehouse_ids}", warehouseMasterController.GetWarehouseWithMultiId)
	router.Get("/drop-down", warehouseMasterController.GetAllIsActive)
	router.Post("/", warehouseMasterController.Save)
	router.Patch("/{warehouse_id}", warehouseMasterController.ChangeStatus)

	return router
}

func WarehouseLocationRouter(
	warehouseLocationController masterwarehousecontroller.WarehouseLocationController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", warehouseLocationController.GetAll)
	router.Get("/by-id/{warehouse_location_id}", warehouseLocationController.GetById)
	router.Post("/", warehouseLocationController.Save)
	router.Patch("/{warehouse_location_id}", warehouseLocationController.ChangeStatus)

	return router
}

func MovingCodeRouter(
	MovingCodeController mastercontroller.MovingCodeController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", MovingCodeController.GetAllMovingCode)
	router.Post("/", MovingCodeController.SaveMovingCode)
	router.Patch("/priority-increase/{moving_code_id}", MovingCodeController.ChangePriorityMovingCode)
	router.Patch("/activation/{moving_code_id}", MovingCodeController.ChangeStatusMovingCode)

	return router
}

func ForecastMasterRouter(
	forecastMasterController mastercontroller.ForecastMasterController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", forecastMasterController.GetAllForecastMaster)
	router.Get("/{forecast_master_id}", forecastMasterController.GetForecastMasterById)
	router.Post("/", forecastMasterController.SaveForecastMaster)
	router.Patch("/{forecast_master_id}", forecastMasterController.ChangeStatusForecastMaster)

	return router
}

func SkillLevelRouter(
	SkillLevelController mastercontroller.SkillLevelController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", SkillLevelController.GetAllSkillLevel)
	router.Get("/{skill_level_id}", SkillLevelController.GetSkillLevelById)
	router.Post("/", SkillLevelController.SaveSkillLevel)
	router.Patch("/{skill_level_id}", SkillLevelController.ChangeStatusSkillLevel)

	return router
}

func ShiftScheduleRouter(
	ShiftScheduleController mastercontroller.ShiftScheduleController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", ShiftScheduleController.GetAllShiftSchedule)
	// router.Get("/drop-down", ShiftScheduleController.GetAllShiftScheduleIsActive)
	// router.Get("/by-code/{operation_group_code}", ShiftScheduleController.GetShiftScheduleByCode)
	router.Post("/", ShiftScheduleController.SaveShiftSchedule)
	router.Get("/by-id/{shift_schedule_id}", ShiftScheduleController.GetShiftScheduleById)
	router.Patch("/{shift_schedule_id}", ShiftScheduleController.ChangeStatusShiftSchedule)

	return router
}

func IncentiveMasterRouter(
	IncentiveMasterController mastercontroller.IncentiveMasterController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", IncentiveMasterController.GetAllIncentiveMaster)
	router.Get("/{incentive_level_id}", IncentiveMasterController.GetIncentiveMasterById)
	router.Post("/", IncentiveMasterController.SaveIncentiveMaster)
	router.Patch("/{incentive_level_id}", IncentiveMasterController.ChangeStatusIncentiveMaster)

	return router
}

func FieldActionRouter(
	FieldActionController mastercontroller.FieldActionController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", FieldActionController.GetAllFieldAction)
	router.Get("/header/by-id/{field_action_system_number}", FieldActionController.GetFieldActionHeaderById)
	router.Get("/vehicle-detail/all/by-id/{field_action_system_number}", FieldActionController.GetAllFieldActionVehicleDetailById)
	router.Get("/vehicle-detail/by-id/{field_action_eligible_vehicle_system_number}", FieldActionController.GetFieldActionVehicleDetailById)
	router.Get("/item-detail/all/by-id/{field_action_eligible_vehicle_system_number}", FieldActionController.GetAllFieldActionVehicleItemDetailById)
	router.Get("/item-detail/by-id/{field_action_eligible_vehicle_item_system_number}", FieldActionController.GetFieldActionVehicleItemDetailById)
	router.Post("/", FieldActionController.SaveFieldAction)
	router.Post("/vehicle-detail/{field_action_system_number}", FieldActionController.PostFieldActionVehicleDetail)
	router.Post("/multi-vehicle-detail/{field_action_system_number}", FieldActionController.PostMultipleVehicleDetail)
	router.Post("/item-detail/{field_action_eligible_vehicle_system_number}", FieldActionController.PostFieldActionVehicleItemDetail)
	router.Post("/all-item-detail/{field_action_system_number}", FieldActionController.PostVehicleItemIntoAllVehicleDetail)
	router.Patch("/header/by-id/{field_action_system_number}", FieldActionController.ChangeStatusFieldAction)
	router.Patch("/vehicle-detail/by-id/{field_action_eligible_vehicle_system_number}", FieldActionController.ChangeStatusFieldActionVehicle)
	router.Patch("/item-detail/by-id/{field_action_eligible_vehicle_item_system_number}", FieldActionController.ChangeStatusFieldActionVehicleItem)

	return router
}

func WarrantyFreeServiceRouter(
	warrantyFreeServiceController mastercontroller.WarrantyFreeServiceController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", warrantyFreeServiceController.GetAllWarrantyFreeService)
	router.Get("/{warranty_free_services_id}", warrantyFreeServiceController.GetWarrantyFreeServiceByID)
	router.Post("/", warrantyFreeServiceController.SaveWarrantyFreeService)
	router.Patch("/{warranty_free_services_id}", warrantyFreeServiceController.ChangeStatusWarrantyFreeService)

	return router
}

func DiscountRouter(
	discountController mastercontroller.DiscountController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", discountController.GetAllDiscount)
	router.Get("/drop-down", discountController.GetAllDiscountIsActive)
	router.Get("/by-code", discountController.GetDiscountByCode)
	router.Get("/by-id/{id}", discountController.GetDiscountById)
	router.Post("/", discountController.SaveDiscount)
	router.Patch("/{id}", discountController.ChangeStatusDiscount)

	return router
}

func IncentiveGroupRouter(
	incentiveGroupController mastercontroller.IncentiveGroupController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/", incentiveGroupController.GetAllIncentiveGroup)
	router.Get("/drop-down", incentiveGroupController.GetAllIncentiveGroupIsActive)
	router.Get("/by-id/{id}", incentiveGroupController.GetIncentiveGroupById)
	router.Post("/", incentiveGroupController.SaveIncentiveGroup)
	router.Patch("/{id}", incentiveGroupController.ChangeStatusIncentiveGroup)

	return router
}

func IncentiveGroupDetailRouter(
	incentiveGroupDetailController mastercontroller.IncentiveGroupDetailController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)

	router.Get("/{id}", incentiveGroupDetailController.GetAllIncentiveGroupDetail)
	router.Get("/by-id/{incentive_group_detail_id}", incentiveGroupDetailController.GetIncentiveGroupDetailById)
	router.Post("/", incentiveGroupDetailController.SaveIncentiveGroupDetail)

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

	return router
}

func SwaggerRouter() chi.Router {
	router := chi.NewRouter()

	// Use middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Serve Swagger UI index.html
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"), //The url pointing to API definition
	))

	return router
}
