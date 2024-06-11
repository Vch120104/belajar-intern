package route

import (
	mastercontroller "after-sales/api/controllers/master"
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	masterwarehousecontroller "after-sales/api/controllers/master/warehouse"
	transactionsparepartcontroller "after-sales/api/controllers/transactions/sparepart"
	transactionworkshopcontroller "after-sales/api/controllers/transactions/workshop"
	"after-sales/api/middlewares"

	_ "after-sales/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

/* Master */

func ItemClassRouter(
	itemClassController masteritemcontroller.ItemClassController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)
	//test
	router.Get("/drop-down", itemClassController.GetItemClassDropdown)
	router.Get("/drop-down/by-group-id/{item_group_id}", itemClassController.GetItemClassDropDownbyGroupId)
	router.Get("/", itemClassController.GetAllItemClass)
	router.Get("/by-code/{item_class_code}", itemClassController.GetItemClassByCode)
	router.Get("/{item_class_id}", itemClassController.GetItemClassbyId)
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
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", unitOfMeasurementController.GetAllUnitOfMeasurement)
	router.Get("/{uom_id}", unitOfMeasurementController.GetUnitOfMeasurementById)
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
	router.Use(middlewares.MetricsMiddleware)

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
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", markupMasterController.GetMarkupMasterList)
	router.Get("/{markup_master_id}", markupMasterController.GetMarkupMasterByID)
	router.Get("/code/{markup_master_code}", markupMasterController.GetMarkupMasterByCode)
	router.Get("/dropdown", markupMasterController.GetAllMarkupMasterIsActive)
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
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", markupRateController.GetAllMarkupRate)
	router.Get("/{markup_rate_id}", markupRateController.GetMarkupRateByID)
	router.Get("/markup-master/{markup_master_id}/order-type/{order_type_id}", markupRateController.GetMarkupRateByMarkupMasterAndOrderType)
	router.Post("/", markupRateController.SaveMarkupRate)
	router.Patch("/{markup_rate_id}", markupRateController.ChangeStatusMarkupRate)

	return router
}

func ItemLevelRouter(
	itemLevelController masteritemcontroller.ItemLevelController,
) chi.Router {
	router := chi.NewRouter()

	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", itemLevelController.GetAll)
	router.Get("/{item_level_id}", itemLevelController.GetById)

	router.Get("/drop-down-item-level/{item_level}", itemLevelController.GetItemLevelDropDown)
	router.Get("/look-up-item-level/{item_class_id}", itemLevelController.GetItemLevelLookUp)

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
	router.Use(middlewares.MetricsMiddleware)
	router.Get("/", itemController.GetAllItem)
	router.Get("/{item_id}", itemController.GetItembyId)
	router.Get("/lookup", itemController.GetAllItemLookup)
	router.Get("/multi-id/{item_ids}", itemController.GetItemWithMultiId)
	router.Get("/by-code/{item_code}", itemController.GetItemByCode)
	router.Get("/uom-type/drop-down", itemController.GetUomTypeDropDown)
	router.Get("/uom/drop-down/{uom_type_id}", itemController.GetUomDropDown)
	router.Post("/", itemController.SaveItem)
	router.Patch("/{item_id}", itemController.ChangeStatusItem)

	router.Get("/detail", itemController.GetAllItemDetail)
	router.Get("/{item_id}/detail/{item_detail_id}", itemController.GetItemDetailById)
	router.Post("/{item_id}/detail", itemController.AddItemDetail)
	router.Delete("/{item_id}/detail/{item_detail_id}", itemController.DeleteItemDetail)

	return router
}

func ItemLocationRouter(
	ItemLocationController masteritemcontroller.ItemLocationController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	//master
	router.Get("/", ItemLocationController.GetAllItemLocation)
	router.Get("/{item_location_id}", ItemLocationController.GetItemLocationById)
	router.Post("/", ItemLocationController.SaveItemLocation)

	//detail
	router.Get("/detail", ItemLocationController.GetAllItemLocationDetail)
	router.Get("/popup-location", ItemLocationController.PopupItemLocation)
	router.Post("/detail", ItemLocationController.AddItemLocation)
	router.Delete("/detail/{item_location_detail_id}", ItemLocationController.DeleteItemLocation)

	// new
	router.Get("/new/get-all/", ItemLocationController.GetAllItemLoc)
	router.Get("/new/get-by-id/{item_location_id}", ItemLocationController.GetByIdItemLoc)
	router.Post("/new/save", ItemLocationController.SaveItemLoc)
	router.Delete("/new/delete/{item_location_id}", ItemLocationController.DeleteItemLoc)

	return router
}

func ItemSubstituteRouter(
	itemSubstituteController masteritemcontroller.ItemSubstituteController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

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

func ItemPackageRouter(
	ItemPackageController masteritemcontroller.ItemPackageController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

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

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/by-package-id/{item_package_id}", ItemPackageDetailController.GetItemPackageDetailByItemPackageId)
	router.Get("/{item_package_detail_id}", ItemPackageDetailController.GetItemPackageDetailById)
	router.Post("/", ItemPackageDetailController.CreateItemPackageDetailByItemPackageId)
	router.Patch("/{item_package_detail_id}", ItemPackageDetailController.ChangeStatusItemPackageDetail)
	router.Put("/", ItemPackageDetailController.UpdateItemPackageDetail)

	return router
}

func ItemImportRouter(
	ItemImportController masteritemcontroller.ItemImportController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", ItemImportController.GetAllItemImport)
	router.Get("/{item_import_id}", ItemImportController.GetItemImportbyId)
	router.Post("/", ItemImportController.SaveItemImport)
	router.Patch("/", ItemImportController.UpdateItemImport)
	// router.Get("/{item_import_id}", ItemImportController.GetItemPackageById)

	return router
}

func ItemModelMappingRouter(
	ItemModelMappingController masteritemcontroller.ItemModelMappingController,
) chi.Router {
	router := chi.NewRouter()

	router.Post("/", ItemModelMappingController.CreateItemModelMapping)
	router.Get("/{item_id}", ItemModelMappingController.GetItemModelMappingByItemId)
	router.Patch("/{item_detail_id}", ItemModelMappingController.UpdateItemModelMapping)
	//router.PanicHandler = exceptions.ErrorHandler

	return router
}

func MovingCodeRouter(
	MovingCodeController mastercontroller.MovingCodeController,
) chi.Router {
	router := chi.NewRouter()

	router.Post("/", MovingCodeController.CreateMovingCode)
	router.Get("/{moving_code_id}", MovingCodeController.GetMovingCodebyId)
	router.Put("/", MovingCodeController.UpdateMovingCode)
	router.Patch("/{moving_code_id}", MovingCodeController.ChangeStatusMovingCode)
	router.Get("/", MovingCodeController.GetAllMovingCode)
	router.Patch("/push-priority/{moving_code_id}", MovingCodeController.PushMovingCodePriority)
	//router.PanicHandler = exceptions.ErrorHandler

	return router
}

func IncentiveGroupRouter(
	incentiveGroupController mastercontroller.IncentiveGroupController,
) chi.Router {
	router := chi.NewRouter()
	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", incentiveGroupController.GetAllIncentiveGroup)
	router.Get("/is-active", incentiveGroupController.GetAllIncentiveGroupIsActive)
	router.Get("/dropdown", incentiveGroupController.GetAllIncentiveGroupDropDown)
	router.Get("/by-id/{incentive_group_id}", incentiveGroupController.GetIncentiveGroupById)
	router.Post("/", incentiveGroupController.SaveIncentiveGroup)
	router.Patch("/{incentive_group_id}", incentiveGroupController.ChangeStatusIncentiveGroup)
	router.Put("/{incentive_group_id}", incentiveGroupController.UpdateIncentiveGroup)
	return router
}

func PriceListRouter(
	priceListController masteritemcontroller.PriceListController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

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
	router.Use(middlewares.MetricsMiddleware)

	//bom master
	router.Get("/", BomController.GetBomMasterList)
	router.Get("/{bom_master_id}", BomController.GetBomMasterById)
	router.Post("/", BomController.SaveBomMaster)
	router.Patch("/{bom_master_id}", BomController.ChangeStatusBomMaster)

	//bom detail
	// Detail
	router.Get("/detail/all", BomController.GetBomDetailList)
	router.Get("/{bom_master_id}/detail", BomController.GetBomDetailById)
	router.Get("/detail/{bom_detail_id}", BomController.GetBomDetailByIds)
	router.Post("/detail", BomController.SaveBomDetail)
	router.Delete("/detail/{bom_detail_id}", BomController.DeleteBomDetail)

	//bom lookup
	router.Get("/popup-item", BomController.GetBomItemList)

	return router
}

func PurchasePriceRouter(
	PurchasePriceController masteritemcontroller.PurchasePriceController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	//master
	router.Get("/", PurchasePriceController.GetAllPurchasePrice)
	router.Get("/by-id/{purchase_price_id}", PurchasePriceController.GetPurchasePriceById)
	router.Post("/", PurchasePriceController.SavePurchasePrice)
	router.Patch("/{purchase_price_id}", PurchasePriceController.ChangeStatusPurchasePrice)

	//detail
	router.Get("/detail", PurchasePriceController.GetAllPurchasePriceDetail)
	router.Get("/{purchase_price_id}/detail", PurchasePriceController.GetPurchasePriceDetailById)
	router.Post("/detail", PurchasePriceController.AddPurchasePrice)
	router.Delete("/detail/{purchase_price_detail_id}", PurchasePriceController.DeletePurchasePrice)

	return router
}

func LandedCostMasterRouter(
	LandedCostMaster masteritemcontroller.LandedCostMasterController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", LandedCostMaster.GetAllLandedCostMaster)
	router.Get("/{landed_cost_id}", LandedCostMaster.GetByIdLandedCost)
	router.Post("/", LandedCostMaster.SaveLandedCostMaster)
	router.Patch("/activate/", LandedCostMaster.ActivateLandedCostMaster)
	router.Patch("/deactivate/", LandedCostMaster.DeactivateLandedCostmaster)
	router.Put("/{landed_cost_id}", LandedCostMaster.UpdateLandedCostMaster)

	return router
}

func OperationGroupRouter(
	operationGroupController masteroperationcontroller.OperationGroupController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

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
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", operationSectionController.GetAllOperationSectionList)
	router.Get("/by-id/{operation_section_id}", operationSectionController.GetOperationSectionByID)
	router.Get("/by-name", operationSectionController.GetOperationSectionName)
	router.Get("/code-by-group-id/{operation_group_id}", operationSectionController.GetSectionCodeByGroupId)
	router.Post("/", operationSectionController.SaveOperationSection)
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
	router.Use(middlewares.MetricsMiddleware)

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
	router.Use(middlewares.MetricsMiddleware)

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
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", operationCodeController.GetAllOperationCode)
	router.Get("/by-id/{operation_id}", operationCodeController.GetByIdOperationCode)
	router.Get("/by-code/{operation_code}", operationCodeController.GetByCodeOperationCode)
	router.Get("/by-code/{operation_code}", operationCodeController.GetByCodeOperationCode)
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
	router.Use(middlewares.MetricsMiddleware)

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
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", warehouseGroupController.GetAllWarehouseGroup)
	router.Get("/by-code/{warehouse_group_code}", warehouseGroupController.GetbyGroupCode)
	router.Get("/{warehouse_group_id}", warehouseGroupController.GetByIdWarehouseGroup)
	router.Get("/drop-down/{warehouse_group_id}", warehouseGroupController.GetWarehouseGroupDropdownbyId)
	router.Get("/drop-down", warehouseGroupController.GetWarehouseGroupDropDown)
	router.Post("/", warehouseGroupController.SaveWarehouseGroup)
	router.Patch("/{warehouse_group_id}", warehouseGroupController.ChangeStatusWarehouseGroup)

	return router
}

func WarehouseLocationDefinitionRouter(
	WarehouseLocationDefinitionController masterwarehousecontroller.WarehouseLocationDefinitionController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", WarehouseLocationDefinitionController.GetAll)
	router.Get("/by-level/{warehouse_location_definition_level_id}/{warehouse_location_definition_id}", WarehouseLocationDefinitionController.GetByLevel)
	router.Get("/by-id/{warehouse_location_definition_id}", WarehouseLocationDefinitionController.GetById)
	router.Get("/popup-level", WarehouseLocationDefinitionController.PopupWarehouseLocationLevel)
	router.Post("/", WarehouseLocationDefinitionController.Save)
	router.Put("/{warehouse_location_definition_id}", WarehouseLocationDefinitionController.SaveData)
	router.Patch("/{warehouse_location_definition_id}", WarehouseLocationDefinitionController.ChangeStatus)

	return router
}

func WarehouseMasterRouter(
	warehouseMasterController masterwarehousecontroller.WarehouseMasterController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", warehouseMasterController.GetAll)
	router.Get("/{warehouse_id}", warehouseMasterController.GetById)
	router.Get("/by-code/{warehouse_code}", warehouseMasterController.GetByCode)
	router.Get("/multi-id/{warehouse_ids}", warehouseMasterController.GetWarehouseWithMultiId)
	router.Get("/is-active", warehouseMasterController.GetAllIsActive)
	router.Get("/drop-down", warehouseMasterController.DropdownWarehouse)
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
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", warehouseLocationController.GetAll)
	router.Get("/by-id/{warehouse_location_id}", warehouseLocationController.GetById)
	router.Post("/", warehouseLocationController.Save)
	router.Patch("/{warehouse_location_id}", warehouseLocationController.ChangeStatus)

	return router
}

func ForecastMasterRouter(
	forecastMasterController mastercontroller.ForecastMasterController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", forecastMasterController.GetAllForecastMaster)
	router.Get("/{forecast_master_id}", forecastMasterController.GetForecastMasterById)
	router.Post("/", forecastMasterController.SaveForecastMaster)
	router.Patch("/{forecast_master_id}", forecastMasterController.ChangeStatusForecastMaster)

	return router
}

func AgreementRouter(
	AgreementController mastercontroller.AgreementController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", AgreementController.GetAllAgreement)
	router.Get("/{agreement_id}", AgreementController.GetAgreementById)
	router.Post("/", AgreementController.SaveAgreement)
	router.Patch("/{agreement_id}", AgreementController.ChangeStatusAgreement)

	router.Get("/{agreement_id}/discount/group", AgreementController.GetAllDiscountGroup)
	router.Get("/{agreement_id}/discount/group/{agreement_discount_group_id}", AgreementController.GetDiscountGroupAgreementById)
	router.Post("/{agreement_id}/discount/group", AgreementController.AddDiscountGroup)
	router.Delete("/{agreement_id}/discount/group/{agreement_discount_group_id}", AgreementController.DeleteDiscountGroup)

	router.Get("/{agreement_id}/discount/item", AgreementController.GetAllItemDiscount)
	router.Get("/{agreement_id}/discount/item/{agreement_item_id}", AgreementController.GetDiscountItemAgreementById)
	router.Post("/{agreement_id}/discount/item", AgreementController.AddItemDiscount)
	router.Delete("/{agreement_id}/discount/item/{agreement_item_id}", AgreementController.DeleteItemDiscount)

	router.Get("/{agreement_id}/discount/value", AgreementController.GetAllDiscountValue)
	router.Get("/{agreement_id}/discount/value/{agreement_discount_id}", AgreementController.GetDiscountValueAgreementById)
	router.Post("/{agreement_id}/discount/value", AgreementController.AddDiscountValue)
	router.Delete("/{agreement_id}/discount/value/{agreement_discount_id}", AgreementController.DeleteDiscountValue)

	return router
}

func SkillLevelRouter(
	SkillLevelController mastercontroller.SkillLevelController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

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
	router.Use(middlewares.MetricsMiddleware)

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
	// Gunakan middleware NotFoundHandler
	// router.Use(middleware.NotFoundHandler)

	router.Get("/", IncentiveMasterController.GetAllIncentiveMaster)
	router.Get("/{incentive_level_id}", IncentiveMasterController.GetIncentiveMasterById)
	router.Post("/", IncentiveMasterController.SaveIncentiveMaster)
	router.Patch("/{incentive_level_id}", IncentiveMasterController.ChangeStatusIncentiveMaster)

	////router.PanicHandler = exceptions.ErrorHandler

	return router
}

func FieldActionRouter(
	FieldActionController mastercontroller.FieldActionController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

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
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", warrantyFreeServiceController.GetAllWarrantyFreeService)
	router.Get("/{warranty_free_services_id}", warrantyFreeServiceController.GetWarrantyFreeServiceByID)
	router.Post("/", warrantyFreeServiceController.SaveWarrantyFreeService)
	router.Patch("/{warranty_free_services_id}", warrantyFreeServiceController.ChangeStatusWarrantyFreeService)

	return router
}

func PackageMasterRouter(
	PackageMasterController mastercontroller.PackageMasterController,
) chi.Router {
	router := chi.NewRouter()
	router.Get("/", PackageMasterController.GetAllPackageMaster)
	router.Get("/detail/{package_id}", PackageMasterController.GetAllPackageMasterDetail)
	router.Get("/header/{package_id}", PackageMasterController.GetByIdPackageMaster)
	router.Get("/detail/{package_id}/{package_detail_id}/{line_type_id}", PackageMasterController.GetByIdPackageMasterDetail)
	router.Get("/copy/{package_id}/{package_name}/{model_id}", PackageMasterController.CopyToOtherModel)

	router.Post("/", PackageMasterController.SavepackageMaster)
	router.Post("/bodyshop/{package_id}", PackageMasterController.SavePackageMasterDetailBodyshop)
	router.Post("/workshop", PackageMasterController.SavePackageMasterDetailWorkshop)

	router.Patch("/{package_id}", PackageMasterController.ChangeStatusPackageMaster)
	router.Patch("/detail/activate/{package_id}/{package_detail_id}", PackageMasterController.ActivateMultiIdPackageMasterDetail)
	router.Patch("/detail/deactivate/{package_id}/{package_detail_id}", PackageMasterController.DeactivateMultiIdPackageMasterDetail)

	return router
}

func DiscountRouter(
	discountController mastercontroller.DiscountController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", discountController.GetAllDiscount)
	router.Get("/drop-down", discountController.GetAllDiscountIsActive)
	router.Get("/by-code", discountController.GetDiscountByCode)
	router.Get("/by-id/{id}", discountController.GetDiscountById)
	router.Post("/", discountController.SaveDiscount)
	router.Patch("/{id}", discountController.ChangeStatusDiscount)

	return router
}

func CampaignMasterRouter(
	campaignmastercontroller mastercontroller.CampaignMasterController,
) chi.Router {
	router := chi.NewRouter()
	//campaign master header
	router.Get("/", campaignmastercontroller.GetAllCampaignMaster)
	router.Get("/{campaign_id}", campaignmastercontroller.GetByIdCampaignMaster)
	router.Get("/history", campaignmastercontroller.GetAllCampaignMasterCodeAndName)
	router.Post("/", campaignmastercontroller.SaveCampaignMaster)
	router.Patch("/{campaign_id}", campaignmastercontroller.ChangeStatusCampaignMaster)

	//campaign master detail
	router.Get("/detail/{campaign_id}", campaignmastercontroller.GetAllCampaignMasterDetail)
	router.Get("/detail/by-id/{campaign_detail_id}", campaignmastercontroller.GetByIdCampaignMasterDetail)
	router.Post("/detail", campaignmastercontroller.SaveCampaignMasterDetail)
	router.Post("/detail/save-from-history/{campaign_id_1}/{campaign_id_2}", campaignmastercontroller.SaveCampaignMasterDetailFromHistory)
	router.Patch("/detail/deactivate/{campaign_detail_id}/{campaign_id}", campaignmastercontroller.DeactivateCampaignMasterDetail)
	router.Patch("/detail/activate/{campaign_detail_id}/{campaign_id}", campaignmastercontroller.ActivateCampaignMasterDetail)
	router.Put("/detail/update/{campaign_detail_id}", campaignmastercontroller.UpdateCampaignMasterDetail)

	//from package master
	router.Get("/package", campaignmastercontroller.GetAllPackageMasterToCopy)
	// router.Get("/package-copy/{package_id}/{campaign_id}",campaignmastercontroller.SelectFromPackageMaster)

	return router
}

func IncentiveGroupDetailRouter(
	incentiveGroupDetailController mastercontroller.IncentiveGroupDetailController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/{id}", incentiveGroupDetailController.GetAllIncentiveGroupDetail)
	router.Get("/by-id/{incentive_group_detail_id}", incentiveGroupDetailController.GetIncentiveGroupDetailById)
	router.Post("/", incentiveGroupDetailController.SaveIncentiveGroupDetail)
	router.Put("/{incentive_group_detail_id}", incentiveGroupDetailController.UpdateIncentiveGroupDetail)

	return router
}

func DeductionRouter(
	DeductionController mastercontroller.DeductionController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", DeductionController.GetAllDeductionList)
	router.Get("/{id}", DeductionController.GetAllDeductionDetail)
	router.Get("/by-detail-id/{id}", DeductionController.GetByIdDeductionDetail)
	router.Get("/by-header-id/{id}", DeductionController.GetDeductionById)
	router.Post("/detail", DeductionController.SaveDeductionDetail)
	router.Post("/", DeductionController.SaveDeductionList)
	router.Patch("/{id}", DeductionController.ChangeStatusDeduction)

	return router
}

func BookingEstimationRouter(
	BookingEstimationController transactionworkshopcontroller.BookingEstimationController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	router.Get("/", BookingEstimationController.GetAll)
	router.Get("/normal", BookingEstimationController.New)
	router.Get("/find/{work_order_system_number}", BookingEstimationController.GetById)
	router.Put("/{id}", BookingEstimationController.Save)
	router.Post("/submit", BookingEstimationController.Submit)
	router.Delete("/{id}", BookingEstimationController.Void)
	router.Put("/close/{id}", BookingEstimationController.CloseOrder)

	return router
}

func WorkOrderRouter(
	WorkOrderController transactionworkshopcontroller.WorkOrderController,
) chi.Router {
	router := chi.NewRouter()

	// Apply the CORS middleware to all routes
	router.Use(middlewares.SetupCorsMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.MetricsMiddleware)

	//add trx normal
	router.Get("/", WorkOrderController.GetAll)
	router.Get("/normal/{work_order_system_number}", WorkOrderController.GetById)
	router.Post("/normal", WorkOrderController.New)
	router.Post("/normal/{work_order_system_number}/submit", WorkOrderController.Submit)
	router.Put("/normal/{work_order_system_number}", WorkOrderController.Save)
	router.Delete("/normal/{work_order_system_number}", WorkOrderController.Void)
	router.Patch("/normal/{work_order_system_number}/close", WorkOrderController.CloseOrder)

	//add post trx sub
	router.Get("/normal/requestservice", WorkOrderController.GetAllRequest)
	router.Get("/normal/{work_order_system_number}/requestservice/{work_order_service_id}", WorkOrderController.GetRequestById)
	router.Post("/normal/{work_order_system_number}/requestservice", WorkOrderController.AddRequest)
	router.Put("/normal/{work_order_system_number}/requestservice/{work_order_service_id}", WorkOrderController.UpdateRequest)
	router.Delete("/normal/{work_order_system_number}/requestservice/{work_order_service_id}", WorkOrderController.DeleteRequest)

	router.Get("/normal/vehicleservice", WorkOrderController.GetAllVehicleService)
	router.Get("/normal/{work_order_system_number}/vehicleservice/{work_order_service_vehicle_id}", WorkOrderController.GetVehicleServiceById)
	router.Put("/normal/{work_order_system_number}/vehicleservice/{work_order_service_vehicle_id}", WorkOrderController.UpdateVehicleService)
	router.Post("/normal/{work_order_system_number}/vehicleservice", WorkOrderController.AddVehicleService)
	router.Delete("/normal/{work_order_system_number}/vehicleservice/{work_order_service_vehicle_id}", WorkOrderController.DeleteVehicleService)

	//add trx detail
	router.Get("/normal/{work_order_system_number}/detail", WorkOrderController.GetAllDetailWorkOrder)
	router.Get("/normal/{work_order_system_number}/detail/{work_order_detail_id}", WorkOrderController.GetDetailByIdWorkOrder)
	router.Post("/normal/{work_order_system_number}/detail", WorkOrderController.AddDetailWorkOrder)
	router.Put("/normal/{work_order_system_number}/detail/{work_order_detail_id}", WorkOrderController.UpdateDetailWorkOrder)
	router.Delete("/normal/{work_order_system_number}/detail/{work_order_detail_id}", WorkOrderController.DeleteDetailWorkOrder)

	//new support function form
	router.Get("/dropdown-status", WorkOrderController.NewStatus)
	router.Get("/dropdown-type", WorkOrderController.NewType)
	router.Get("/dropdown-bill", WorkOrderController.NewBill)
	router.Get("/dropdown-drop-point", WorkOrderController.NewDropPoint)
	router.Get("/dropdown-brand", WorkOrderController.NewVehicleBrand)
	router.Get("/dropdown-model/{brand_id}", WorkOrderController.NewVehicleModel)
	router.Get("/lookup-vehicle", WorkOrderController.VehicleLookup)
	router.Get("/lookup-campaign", WorkOrderController.CampaignLookup)

	router.Post("/normalbooking", WorkOrderController.NewBooking)
	router.Post("/affiliated", WorkOrderController.NewAffiliated)

	return router
}

func SupplySlipRouter(
	SupplySlipController transactionsparepartcontroller.SupplySlipController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/{supply_system_number}", SupplySlipController.GetSupplySlipByID)

	return router
}

func SalesOrderRouter(
	SalesOrderController transactionsparepartcontroller.SalesOrderController,
) chi.Router {
	router := chi.NewRouter()

	router.Get("/{sales_order_system_number}", SalesOrderController.GetSalesOrderByID)

	return router
}

func SwaggerRouter() chi.Router {
	router := chi.NewRouter()

	// Izinkan akses ke Swagger di /aftersales-service/docs
	router.Get("/aftersales-service/docs/v1/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/v1/doc.json"), // Ubah dengan alamat server
	))

	return router
}
