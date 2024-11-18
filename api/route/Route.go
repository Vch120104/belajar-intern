package route

import (
	"after-sales/api/config"
	mastercontroller "after-sales/api/controllers/master"
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	masterwarehousecontroller "after-sales/api/controllers/master/warehouse"
	"after-sales/api/helper"
	masteritemrepositoryimpl "after-sales/api/repositories/master/item/repositories-item-impl"
	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	masterwarehouserepositoryimpl "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"
	masteritemserviceimpl "after-sales/api/services/master/item/services-item-impl"
	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	masterwarehouseserviceimpl "after-sales/api/services/master/warehouse/services-warehouse-impl"

	transactionjpcbcontroller "after-sales/api/controllers/transactions/JPCB"
	transactionbodyshopcontroller "after-sales/api/controllers/transactions/bodyshop"
	transactionsparepartcontroller "after-sales/api/controllers/transactions/sparepart"
	transactionworkshopcontroller "after-sales/api/controllers/transactions/workshop"
	transactionjpcbrepositoryimpl "after-sales/api/repositories/transaction/JPCB/repositories-jpcb-impl"
	transactionbodyshoprepositoryimpl "after-sales/api/repositories/transaction/bodyshop/repositories-bodyshop-impl"
	transactionsparepartrepositoryimpl "after-sales/api/repositories/transaction/sparepart/repositories-sparepart-impl"
	transactionworkshoprepositoryimpl "after-sales/api/repositories/transaction/workshop/repositories-workshop-impl"
	transactionjpcbserviceimpl "after-sales/api/services/transaction/JPCB/services-jpcb-impl"
	transactionbodyshopserviceimpl "after-sales/api/services/transaction/bodyshop/services-bodyshop-impl"
	transactionsparepartserviceimpl "after-sales/api/services/transaction/sparepart/services-sparepart-impl"
	transactionworkshopserviceimpl "after-sales/api/services/transaction/workshop/services-workshop-impl"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

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
	priceListService := masteritemserviceimpl.StartPriceListService(priceListRepository, itemService, db, rdb)
	priceListController := masteritemcontroller.NewPriceListController(priceListService)

	// Item Class
	itemClassRepository := masteritemrepositoryimpl.StartItemClassRepositoryImpl()
	itemClassService := masteritemserviceimpl.StartItemClassService(itemClassRepository, db, rdb)
	itemClassController := masteritemcontroller.NewItemClassController(itemClassService)

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

	// Item Price Code
	ItemPriceCodeRepository := masteritemrepositoryimpl.StartItemPriceCodeImpl()
	ItemPriceCodeServicwe := masteritemserviceimpl.StartItemPriceCodeService(ItemPriceCodeRepository, db, rdb)
	ItemPriceCodeController := masteritemcontroller.NewItemPriceCodeController(ItemPriceCodeServicwe)

	// Purchase Price
	PurchasePriceRepository := masteritemrepositoryimpl.StartPurchasePriceRepositoryImpl()
	PurchasePriceService := masteritemserviceimpl.StartPurchasePriceService(PurchasePriceRepository, db, rdb)
	PurchasePriceController := masteritemcontroller.NewPurchasePriceController(PurchasePriceService)

	// Item Operation
	ItemOperationRepository := masterrepositoryimpl.StartItemOperationRepositoryImpl()
	ItemOperationService := masterserviceimpl.StartItemOperationService(ItemOperationRepository, db, rdb)
	ItemOperationController := mastercontroller.NewItemOperationController(ItemOperationService)

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
	MovingCodeService := masterserviceimpl.StartMovingCodeServiceImpl(MovingCodeRepository, db)
	MovingCodeController := mastercontroller.NewMovingCodeController(MovingCodeService)

	// Order Type
	OrderTypeRepository := masterrepositoryimpl.StartOrderTypeRepositoryImpl()
	OrderTypeService := masterserviceimpl.StartOrderTypeServiceImpl(OrderTypeRepository, db)
	OrderTypeController := mastercontroller.NewOrderTypeControllerImpl(OrderTypeService)

	// ForecastMaster
	forecastMasterRepository := masterrepositoryimpl.StartForecastMasterRepositoryImpl()
	forecastMasterService := masterserviceimpl.StartForecastMasterService(forecastMasterRepository, db, rdb)
	forecastMasterController := mastercontroller.NewForecastMasterController(forecastMasterService)

	// Gmm Price Code
	gmmPriceCodeRepository := masterrepositoryimpl.StartGmmPriceCodeRepositoryImpl()
	gmmPriceCodeService := masterserviceimpl.StartGmmPriceCodeServiceImpl(gmmPriceCodeRepository, db)
	gmmPriceCodeController := mastercontroller.NewGmmPriceCodeControllerImpl(gmmPriceCodeService)

	// Gmm Discount Setting
	gmmDiscountSettingRepository := masterrepositoryimpl.StartGmmDiscountSettingRepositoryImpl()
	gmmDiscountSettingService := masterserviceimpl.StartGmmDiscountSettingServiceImpl(gmmDiscountSettingRepository, db)
	gmmDiscountSettingController := mastercontroller.NewGmmDiscountSettingControllerImpl(gmmDiscountSettingService)

	// Agreement
	AgreementRepository := masterrepositoryimpl.StartAgreementRepositoryImpl()
	AgreementService := masterserviceimpl.StartAgreementService(AgreementRepository, db, rdb)
	AgreementController := mastercontroller.NewAgreementController(AgreementService)

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

	//labour selling price
	labourSellingPriceRepository := masteroperationrepositoryimpl.StartLabourSellingPriceRepositoryImpl()
	laboruSellingPriceService := masteroperationserviceimpl.StartLabourSellingPriceService(labourSellingPriceRepository, db)
	LabourSellingPriceController := masteroperationcontroller.NewLabourSellingPriceController(laboruSellingPriceService)

	//labour selling price detail
	LabourSellingPriceDetailController := masteroperationcontroller.NewLabourSellingPriceDetailController(laboruSellingPriceService)

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

	//stock transaction type
	StockTransactionTypeRepository := masterrepositoryimpl.NewStockTransactionRepositoryImpl()
	StockTransactionTypeService := masterserviceimpl.NewStockTransactionTypeServiceImpl(StockTransactionTypeRepository, db, rdb)
	StockTransactionTypeController := mastercontroller.NewStockTransactionTypeController(StockTransactionTypeService)

	//stock transaction reason
	StockTransactionReasonRepository := masterrepositoryimpl.StartStockTraansactionReasonRepositoryImpl()
	StockTransactionReasonService := masterserviceimpl.StartStockTransactionReasonServiceImpl(StockTransactionReasonRepository, db, rdb)
	StockTransactionReasonController := mastercontroller.StartStockTransactionReasonController(StockTransactionReasonService)
	// Warehouse Group
	warehouseGroupRepository := masterwarehouserepositoryimpl.OpenWarehouseGroupImpl()
	warehouseGroupService := masterwarehouseserviceimpl.OpenWarehouseGroupService(warehouseGroupRepository, db, rdb)
	warehouseGroupController := masterwarehousecontroller.NewWarehouseGroupController(warehouseGroupService)

	// Warehouse Location
	WarehouseLocationDefinitionRepository := masterwarehouserepositoryimpl.OpenWarehouseLocationDefinitionImpl()
	WarehouseLocationDefinitionService := masterwarehouseserviceimpl.OpenWarehouseLocationDefinitionService(WarehouseLocationDefinitionRepository, db, rdb)
	WarehouseLocationDefinitionController := masterwarehousecontroller.NewWarehouseLocationDefinitionController(WarehouseLocationDefinitionService)

	// Warehouse Master
	warehouseMasterRepository := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl()
	warehouseMasterService := masterwarehouseserviceimpl.OpenWarehouseMasterService(warehouseMasterRepository, db, rdb)
	warehouseMasterController := masterwarehousecontroller.NewWarehouseMasterController(warehouseMasterService)

	// Warehouse Location
	warehouseLocationRepository := masterwarehouserepositoryimpl.OpenWarehouseLocationImpl()
	warehouseLocationService := masterwarehouseserviceimpl.OpenWarehouseLocationService(warehouseLocationRepository, warehouseMasterService, db, rdb)
	warehouseLocationController := masterwarehousecontroller.NewWarehouseLocationController(warehouseLocationService)

	// Warehouse Costing Type
	warehouseCostingTypeRepository := masterwarehouserepositoryimpl.NewWarehouseCostingTypeRepositoryImpl()
	warehouseCostingTypeService := masterwarehouseserviceimpl.NewWarehouseCostingTypeServiceImpl(warehouseCostingTypeRepository, db, rdb)
	warehouseCostingTypeController := masterwarehousecontroller.NewWarehouseCostingTypeController(warehouseCostingTypeService)
	// Item Location
	ItemLocationRepository := masteritemrepositoryimpl.StartItemLocationRepositoryImpl()
	ItemLocationService := masteritemserviceimpl.StartItemLocationService(ItemLocationRepository, warehouseMasterRepository, warehouseLocationRepository, itemRepository, db, rdb)
	ItemLocationController := masteritemcontroller.NewItemLocationController(ItemLocationService)

	//location stock master
	LocationStockRepository := masterwarehouserepository.NewLocationStockRepositoryImpl()
	LocationStockService := masterserviceimpl.NewLocationStockServiceImpl(LocationStockRepository, db, rdb)
	LocationStockController := mastercontroller.NewLocationStockController(LocationStockService)

	// Bom Master
	BomRepository := masteritemrepositoryimpl.StartBomRepositoryImpl()
	BomService := masteritemserviceimpl.StartBomService(BomRepository, db, rdb)
	BomController := masteritemcontroller.NewBomController(BomService)

	//package master
	PackageMasterRepository := masterrepositoryimpl.StartPackageMasterRepositoryImpl()
	PackageMasterService := masterserviceimpl.StartPackageMasterService(PackageMasterRepository, db)
	PackageMasterController := mastercontroller.NewPackageMasterController(PackageMasterService)

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

	// Campaign Master
	CampaignMasterRepository := masterrepositoryimpl.StartCampaignMasterRepositoryImpl()
	CampaignMasterService := masterserviceimpl.StartCampaignMasterService(CampaignMasterRepository, db)
	CampaignMasterController := mastercontroller.NewCampaignMasterController(CampaignMasterService)

	//Field Action
	FieldActionRepository := masterrepositoryimpl.StartFieldActionRepositoryImpl()
	FieldActionService := masterserviceimpl.StartFieldActionService(FieldActionRepository, db, rdb)
	FieldActionController := mastercontroller.NewFieldActionController(FieldActionService)

	//Lookup
	LookupRepository := masterrepositoryimpl.StartLookupRepositoryImpl()
	LookupService := masterserviceimpl.StartLookupService(LookupRepository, db, rdb)
	LookupController := mastercontroller.NewLookupController(LookupService)
	//Item Cycle
	ItemCycleRepository := masterrepositoryimpl.NewItemCycleRepositoryImpl()
	ItemCycleService := masterserviceimpl.NewItemCycleServiceImpl(ItemCycleRepository, db, rdb)
	ItemCycleController := mastercontroller.NewItemCycleController(ItemCycleService)

	/* Transaction */
	//Supply Slip
	SupplySlipRepository := transactionsparepartrepositoryimpl.StartSupplySlipRepositoryImpl()
	SupplySlipService := transactionsparepartserviceimpl.StartSupplySlipService(SupplySlipRepository, db, rdb)
	SupplySlipController := transactionsparepartcontroller.NewSupplySlipController(SupplySlipService)

	//Supply Slip Return
	SupplySlipReturnRepository := transactionsparepartrepositoryimpl.StartSupplySlipReturnRepositoryImpl()
	SupplySlipReturnService := transactionsparepartserviceimpl.StartSupplySlipReturnService(SupplySlipReturnRepository, SupplySlipRepository, db, rdb)
	SupplySlipReturnController := transactionsparepartcontroller.NewSupplySlipReturnController(SupplySlipReturnService)

	//Booking Estimation
	BookingEstimationRepository := transactionworkshoprepositoryimpl.OpenBookingEstimationRepositoryImpl()
	BookingEstimationService := transactionworkshopserviceimpl.OpenBookingEstimationServiceImpl(BookingEstimationRepository, db, rdb)
	BookingEstimationController := transactionworkshopcontroller.NewBookingEstimationController(BookingEstimationService)

	//Work order
	WorkOrderRepository := transactionworkshoprepositoryimpl.OpenWorkOrderRepositoryImpl()
	WorkOrderService := transactionworkshopserviceimpl.OpenWorkOrderServiceImpl(WorkOrderRepository, db, rdb)
	WorkOrderController := transactionworkshopcontroller.NewWorkOrderController(WorkOrderService)

	//Sales Order
	SalesOrderRepository := transactionsparepartrepositoryimpl.StartSalesOrderRepositoryImpl()
	SalesOrderService := transactionsparepartserviceimpl.StartSalesOrderService(SalesOrderRepository, db, rdb)
	SalesOrderController := transactionsparepartcontroller.NewSalesOrderController(SalesOrderService)

	//Service Request
	ServiceRequestRepository := transactionworkshoprepositoryimpl.OpenServiceRequestRepositoryImpl()
	ServiceRequestService := transactionworkshopserviceimpl.OpenServiceRequestServiceImpl(ServiceRequestRepository, db, rdb)
	ServiceRequestController := transactionworkshopcontroller.NewServiceRequestController(ServiceRequestService)

	//vehicle history
	VehicleHistoryRepository := transactionworkshoprepositoryimpl.NewVehicleHistoryImpl()
	VehicleHistoryServices := transactionworkshopserviceimpl.NewVehicleHistoryServiceImpl(VehicleHistoryRepository, db, rdb)
	VehicleHistoryController := transactionworkshopcontroller.NewVehicleHistoryController(VehicleHistoryServices)

	//Service Receipt
	ServiceReceiptRepository := transactionworkshoprepositoryimpl.OpenServiceReceiptRepositoryImpl()
	ServiceReceiptService := transactionworkshopserviceimpl.OpenServiceReceiptServiceImpl(ServiceReceiptRepository, db, rdb)
	ServiceReceiptController := transactionworkshopcontroller.NewServiceReceiptController(ServiceReceiptService)

	//Purchase Request
	PurchaseRequestRepository := transactionsparepartrepositoryimpl.NewPurchaseRequestRepositoryImpl()
	PurchaseRequestService := transactionsparepartserviceimpl.NewPurchaseRequestImpl(PurchaseRequestRepository, db, rdb)
	PurchaseRequestController := transactionsparepartcontroller.NewPurchaseRequestController(PurchaseRequestService)
	//Purchase Order
	PurchaseOrderRepository := transactionsparepartrepositoryimpl.NewPurchaseOrderRepositoryImpl()
	PurchaseOrderService := transactionsparepartserviceimpl.NewPurchaseOrderService(PurchaseOrderRepository, db, rdb)
	PurchaseOrderController := transactionsparepartcontroller.NewPurchaseOrderControllerImpl(PurchaseOrderService)

	//goods receive
	GoodsReceiveRepository := transactionsparepartrepositoryimpl.NewGoodsReceiveRepositoryImpl()
	GoodsReceiveService := transactionsparepartserviceimpl.NewGoodsReceiveServiceImpl(GoodsReceiveRepository, db, rdb)
	GoodsReceiveController := transactionsparepartcontroller.NewGoodsReceiveController(GoodsReceiveService)
	//binning list
	BinningListRepository := transactionsparepartrepositoryimpl.NewbinningListRepositoryImpl()
	BinningListService := transactionsparepartserviceimpl.NewBinningListServiceImpl(BinningListRepository, db, rdb)
	BinningListController := transactionsparepartcontroller.NewBinningListControllerImpl(BinningListService)

	//Item Inquiry
	ItemInquiryRepository := transactionsparepartrepositoryimpl.StartItemInquiryRepositoryImpl()
	ItemInquiryService := transactionsparepartserviceimpl.StartItemInquiryService(ItemInquiryRepository, db, rdb)
	ItemInquiryController := transactionsparepartcontroller.NewItemInquiryController(ItemInquiryService)

	//stock transaction
	StockTransactionRepository := transactionsparepartrepositoryimpl.StartStockTransactionRepositoryImpl()
	StockTransactionService := masterserviceimpl.StartStockTransactionServiceImpl(StockTransactionRepository, db, rdb)
	StockTransactionController := transactionsparepartcontroller.StartStockTransactionControllerImpl(StockTransactionService)
	//Work Order Allocation
	WorkOrderAllocationRepository := transactionworkshoprepositoryimpl.OpenWorkOrderAllocationRepositoryImpl()
	WorkOrderAllocationService := transactionworkshopserviceimpl.OpenWorkOrderAllocationServiceImpl(WorkOrderAllocationRepository, db, rdb)
	WorkOrderAllocationController := transactionworkshopcontroller.NewWorkOrderAllocationController(WorkOrderAllocationService)

	//Work order bypass
	WorkOrderBypassRepository := transactionworkshoprepositoryimpl.OpenWorkOrderBypassRepositoryImpl()
	WorkOrderBypassService := transactionworkshopserviceimpl.OpenWorkOrderBypassServiceImpl(WorkOrderBypassRepository, db, rdb)
	WorkOrderBypassController := transactionworkshopcontroller.NewWorkOrderBypassController(WorkOrderBypassService)

	//Setting Technician
	SettingTechnicianRepository := transactionjpcbrepositoryimpl.StartSettingTechnicianRepositoryImpl()
	SettingTechnicianService := transactionjpcbserviceimpl.StartServiceTechnicianService(SettingTechnicianRepository, db, rdb)
	SettingTechnicianController := transactionjpcbcontroller.NewSettingTechnicianController(SettingTechnicianService)

	//Technician Attendance
	TechnicianAttendanceRepository := transactionjpcbrepositoryimpl.StartTechnicianAttendanceRepositoryImpl()
	TechnicianAttendanceService := transactionjpcbserviceimpl.StartTechnicianAttendanceImpl(TechnicianAttendanceRepository, db, rdb)
	TechnicianAttendanceController := transactionjpcbcontroller.NewTechnicianAttendanceController(TechnicianAttendanceService)

	//Job Allocation
	JobAllocationRepository := transactionjpcbrepositoryimpl.StartJobAllocationRepositoryImpl()
	JobAllocationService := transactionjpcbserviceimpl.StartJobAllocationService(JobAllocationRepository, db, rdb)
	JobAllocationController := transactionjpcbcontroller.NewJobAllocationController(JobAllocationService)

	//Outstanding Job Allocation
	OutstandingJobAllocationRepository := transactionjpcbrepositoryimpl.StartOutStandingJobAllocationRepository()
	OutstandingJobAllocationService := transactionjpcbserviceimpl.StartOutstandingJobAllocationService(OutstandingJobAllocationRepository, operationCodeRepository, db, rdb)
	OutstandingJobAllocationController := transactionjpcbcontroller.NewOutstandingJobAllocationController(OutstandingJobAllocationService)

	//Car Wash Bay
	CarWashBayRepository := transactionjpcbrepositoryimpl.NewCarWashBayRepositoryImpl()
	CarWashBayService := transactionjpcbserviceimpl.NewCarWashBayServiceImpl(CarWashBayRepository, db, rdb)
	CarWashBayController := transactionjpcbcontroller.NewCarWashBayController(CarWashBayService)

	//Car Wash
	CarWashRepository := transactionjpcbrepositoryimpl.NewCarWashRepositoryImpl()
	CarWashService := transactionjpcbserviceimpl.NewCarWashServiceImpl(CarWashRepository, db, rdb)
	CarWashController := transactionjpcbcontroller.NewCarWashController(CarWashService)

	//Quality Control
	QualityControlRepository := transactionworkshoprepositoryimpl.OpenQualityControlRepositoryImpl()
	QualityControlService := transactionworkshopserviceimpl.OpenQualityControlServiceImpl(QualityControlRepository, db, rdb)
	QualityControlController := transactionworkshopcontroller.NewQualityControlController(QualityControlService)

	//Quality Control
	QualityControlBodyshopRepository := transactionbodyshoprepositoryimpl.OpenQualityControlBodyshopRepositoryImpl()
	QualityControlBodyshopService := transactionbodyshopserviceimpl.OpenQualityControlBodyshopServiceImpl(QualityControlBodyshopRepository, db, rdb)
	QualityControlBodyshopController := transactionbodyshopcontroller.NewQualityControlBodyshopController(QualityControlBodyshopService)

	//Service Workshop
	ServiceWorkshopRepository := transactionworkshoprepositoryimpl.OpenServiceWorkshopRepositoryImpl()
	ServiceWorkshopService := transactionworkshopserviceimpl.OpenServiceWorkshopServiceImpl(ServiceWorkshopRepository, db, rdb)
	ServiceWorkshopController := transactionworkshopcontroller.NewServiceWorkshopController(ServiceWorkshopService)

	//Service Bodyshop
	ServiceBodyshopRepository := transactionbodyshoprepositoryimpl.OpenServiceBodyshopRepositoryImpl()
	ServiceBodyshopService := transactionbodyshopserviceimpl.OpenServiceBodyshopServiceImpl(ServiceBodyshopRepository, db, rdb)
	ServiceBodyshopController := transactionbodyshopcontroller.NewServiceBodyshopController(ServiceBodyshopService)

	//Contract Service
	ContractServiceRepository := transactionworkshoprepositoryimpl.OpenContractServicelRepositoryImpl()
	ContractServiceService := transactionworkshopserviceimpl.OpenContractServiceServiceImpl(ContractServiceRepository, db, rdb)
	ContractServiceController := transactionworkshopcontroller.NewContractServiceController(ContractServiceService)

	//Contract Service Detail
	ContractServiceDetailRepository := transactionworkshoprepositoryimpl.OpenContractServicelDetailRepositoryImpl()
	ContractServiceDetailService := transactionworkshopserviceimpl.OpenContractServiceDetailServiceImpl(ContractServiceDetailRepository, db, rdb)
	ContractServiceDetailController := transactionworkshopcontroller.NewContractServiceDetailController(ContractServiceDetailService)

	/* Master */
	itemClassRouter := ItemClassRouter(itemClassController)
	itemPackageRouter := ItemPackageRouter(itemPackageController)
	ItemModelMappingRouter := ItemModelMappingRouter(ItemModelMappingController)
	itemPackageDetailRouter := ItemPackageDetailRouter(itemPackageDetailController)
	itemImportRouter := ItemImportRouter(ItemImportController)
	ItemPriceCodeRouter := ItemPriceCodeRouter(ItemPriceCodeController)
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
	LabourSellingPriceRouter := LabourSellingPriceRouter(LabourSellingPriceController)
	LabourSellingPriceDetailRouter := LabourSellingPriceDetailRouter(LabourSellingPriceDetailController)
	MovingCodeRouter := MovingCodeRouter(MovingCodeController)
	OrderTypeRouter := OrderTypeRouter(OrderTypeController)
	ForecastMasterRouter := ForecastMasterRouter(forecastMasterController)
	GmmPriceCodeRouter := GmmPriceCodeRouter(gmmPriceCodeController)
	GmmDiscountSettingRouter := GmmDiscountSettingRouter(gmmDiscountSettingController)
	AgreementRouter := AgreementRouter(AgreementController)
	DiscountPercentRouter := DiscountPercentRouter(discountPercentController)
	DiscountRouter := DiscountRouter(discountController)
	MarkupRateRouter := MarkupRateRouter(markupRateController)
	ItemSubstituteRouter := ItemSubstituteRouter(itemSubstituteController)
	ItemLocationRouter := ItemLocationRouter(ItemLocationController)
	WarehouseGroupRouter := WarehouseGroupRouter(warehouseGroupController)
	WarehouseLocation := WarehouseLocationRouter(warehouseLocationController)
	WarehouseLocationDefinition := WarehouseLocationDefinitionRouter(WarehouseLocationDefinitionController)
	WarehouseMaster := WarehouseMasterRouter(warehouseMasterController)
	WarehouseCostingType := WarehouseCostingTypeMasterRouter(warehouseCostingTypeController)
	StockTransactionTypeRouter := StockTransactionTypeRouter(StockTransactionTypeController)
	StockTransactionReasonRouter := StockTransactionReasonRouter(StockTransactionReasonController)
	StockTransactionRouter := StockTransactionRouter(StockTransactionController)
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
	CampaignMasterRouter := CampaignMasterRouter(CampaignMasterController)
	PackageMasterRouter := PackageMasterRouter(PackageMasterController)
	LocationStockRouter := LocationStockRouter(LocationStockController)
	ItemOperationRouter := ItemOperationRouter(ItemOperationController)
	ItemCycleRouter := ItemCycleRouter(ItemCycleController)
	/* Transaction */
	SupplySlipRouter := SupplySlipRouter(SupplySlipController)
	SupplySlipReturnRouter := SupplySlipReturnRouter(SupplySlipReturnController)
	BookingEstimationRouter := BookingEstimationRouter(BookingEstimationController)
	WorkOrderRouter := WorkOrderRouter(WorkOrderController)
	SalesOrderRouter := SalesOrderRouter(SalesOrderController)
	ServiceRequestRouter := ServiceRequestRouter(ServiceRequestController)
	ServiceReceiptRouter := ServiceReceiptRouter(ServiceReceiptController)
	VehicleHistoryRouter := VehicleHistoryRouter(VehicleHistoryController)
	WorkOrderBypassRouter := WorkOrderBypassRouter(WorkOrderBypassController)
	WorkOrderAllocationRouter := WorkOrderAllocationRouter(WorkOrderAllocationController)
	SettingTechnicianRouter := SettingTechnicianRouter(SettingTechnicianController)
	CarWashBayRouter := CarWashBayRouter(CarWashBayController)
	CarWashRouter := CarWashRouter(CarWashController)
	TechnicianAttendanceRouter := TechnicianAttendanceRouter(TechnicianAttendanceController)
	JobAllocationRouter := JobAllocationRouter(JobAllocationController)
	OutstandingJobAllocationRouter := OutstandingJobAllocationRouter(OutstandingJobAllocationController)
	QualityControlRouter := QualityControlRouter(QualityControlController)
	QualityControlBodyshopRouter := QualityControlBodyshopRouter(QualityControlBodyshopController)
	ServiceWorkshopRouter := ServiceWorkshopRouter(ServiceWorkshopController)
	ServiceBodyshopRouter := ServiceBodyshopRouter(ServiceBodyshopController)
	PurchaseRequestRouter := PurchaseRequestRouter(PurchaseRequestController)
	PurchaseOrderRouter := PurchaseOrderRouter(PurchaseOrderController)
	GoodsReceiveRouter := GoodsReceiveRouter(GoodsReceiveController)
	BinningListRouter := BinningListRouter(BinningListController)
	ItemInquiryRouter := ItemInquiryRouter(ItemInquiryController)
	LookupRouter := LookupRouter(LookupController)
	ContractServiceRouter := ContractServiceRouter(ContractServiceController)
	ContractServiceDetailRouter := ContractServiceDetailRouter(ContractServiceDetailController)

	r := chi.NewRouter()
	// Route untuk setiap versi API
	r.Route("/v1", func(r chi.Router) {
		// Tambahkan routing untuk v1 versi di sini
		/* Master Item */
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
		r.Mount("/item-price-code", ItemPriceCodeRouter)
		r.Mount("/purchase-price", PurchasePriceRouter)

		r.Mount("/landed-cost", LandedCostMasterRouter)
		//r.Mount("/import-duty", ImportDutyRouter)

		/* Master Operation */
		r.Mount("/operation-group", OperationGroupRouter)
		r.Mount("/operation-code", OperationCodeRouter)
		r.Mount("/operation-section", OperationSectionRouter)
		r.Mount("/operation-key", OperationKeyRouter)
		r.Mount("/operation-entries", OperationEntriesRouter)
		r.Mount("/operation-model-mapping", OperationModelMappingRouter)
		r.Mount("/labour-selling-price", LabourSellingPriceRouter)
		r.Mount("/labour-selling-price-detail", LabourSellingPriceDetailRouter)

		/* Master Warehouse */
		r.Mount("/warehouse-group", WarehouseGroupRouter)
		r.Mount("/warehouse-location", WarehouseLocation)
		r.Mount("/warehouse-location-definition", WarehouseLocationDefinition)
		r.Mount("/warehouse-master", WarehouseMaster)
		r.Mount("/warehouse-costing-type", WarehouseCostingType)

		/* Master */
		r.Mount("/moving-code", MovingCodeRouter)
		r.Mount("/order-type", OrderTypeRouter)
		r.Mount("/forecast-master", ForecastMasterRouter)
		r.Mount("/gmm-price-code", GmmPriceCodeRouter)
		r.Mount("/gmm-discount-setting", GmmDiscountSettingRouter)
		r.Mount("/agreement", AgreementRouter)
		r.Mount("/package-master", PackageMasterRouter)
		r.Mount("/skill-level", SkillLevelRouter)
		r.Mount("/shift-schedule", ShiftScheduleRouter)
		r.Mount("/incentive", IncentiveMasterRouter)
		//r.Mount("/work-info-massage", WorkInfoRouter)
		r.Mount("/field-action", FieldActionRouter)
		r.Mount("/warranty-free-service", warrantyFreeServiceRouter)
		r.Mount("/campaign-master", CampaignMasterRouter)
		r.Mount("/discount", DiscountRouter)
		r.Mount("/incentive-group", IncentiveGroupRouter)
		r.Mount("/incentive-group-detail", IncentiveGroupDetailRouter)
		r.Mount("/deduction", DeductionRouter)
		r.Mount("/location-stock", LocationStockRouter)
		r.Mount("/item-operation", ItemOperationRouter)
		r.Mount("/item-cycle", ItemCycleRouter)
		r.Mount("/stock-transaction-type", StockTransactionTypeRouter)
		r.Mount("/stock-transaction-reason", StockTransactionReasonRouter)
		/* Transaction */

		/* Transaction JPCB */
		r.Mount("/bay", CarWashBayRouter)
		r.Mount("/setting-technician", SettingTechnicianRouter)
		r.Mount("/technician-attendance", TechnicianAttendanceRouter)
		r.Mount("/job-allocation", JobAllocationRouter)
		r.Mount("/outstanding-job-allocation", OutstandingJobAllocationRouter)
		r.Mount("/car-wash", CarWashRouter)

		/* Transaction Workshop */
		r.Mount("/booking-estimation", BookingEstimationRouter)
		r.Mount("/work-order", WorkOrderRouter)
		r.Mount("/service-request", ServiceRequestRouter)
		r.Mount("/service-receipt", ServiceReceiptRouter)
		r.Mount("/vehicle-history", VehicleHistoryRouter)
		r.Mount("/work-order-allocation", WorkOrderAllocationRouter)
		r.Mount("/work-order-bypass", WorkOrderBypassRouter)
		r.Mount("/quality-control", QualityControlRouter)
		r.Mount("/service-workshop", ServiceWorkshopRouter)
		r.Mount("/contract-service", ContractServiceRouter)
		r.Mount("/contract-service-detail", ContractServiceDetailRouter)

		r.Mount("/stock-transaction", StockTransactionRouter)
		/* Transaction Bodyshop */
		r.Mount("/service-bodyshop", ServiceBodyshopRouter)
		r.Mount("/quality-control-bodyshop", QualityControlBodyshopRouter)

		/* Transaction Sparepart */
		r.Mount("/supply-slip", SupplySlipRouter)
		r.Mount("/supply-slip-return", SupplySlipReturnRouter)
		r.Mount("/sales-order", SalesOrderRouter)
		r.Mount("/purchase-request", PurchaseRequestRouter)
		r.Mount("/purchase-order", PurchaseOrderRouter)
		r.Mount("/goods-receive", GoodsReceiveRouter)

		r.Mount("/binning-list", BinningListRouter)
		r.Mount("/item-inquiry", ItemInquiryRouter)

		/* Support Func Afs */
		r.Mount("/lookup", LookupRouter)
	})
	// Route untuk Swagger
	r.Mount("/aftersales-service/docs", httpSwagger.WrapHandler)
	// Route untuk Prometheus metrics
	r.Mount("/metrics", promhttp.Handler())

	server := http.Server{
		Addr:    config.EnvConfigs.ClientOrigin,
		Handler: r,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
