package route

import (
	"after-sales/api/config"
	mastercontroller "after-sales/api/controllers/master"
	masteritemcontroller "after-sales/api/controllers/master/item"
	masteroperationcontroller "after-sales/api/controllers/master/operation"
	masterwarehousecontroller "after-sales/api/controllers/master/warehouse"
	transactionworkshopcontroller "after-sales/api/controllers/transactions/workshop"

	"after-sales/api/middlewares"
	masteritemrepositoryimpl "after-sales/api/repositories/master/item/repositories-item-impl"
	masteroperationrepositoryimpl "after-sales/api/repositories/master/operation/repositories-operation-impl"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	masterwarehouserepositoryimpl "after-sales/api/repositories/master/warehouse/repositories-warehouse-impl"
	transactionworkshoprepositoryimpl "after-sales/api/repositories/transaction/workshop/repositories-workshop-impl"
	masteritemserviceimpl "after-sales/api/services/master/item/services-item-impl"
	masteroperationserviceimpl "after-sales/api/services/master/operation/services-operation-impl"
	masterwarehouseserviceimpl "after-sales/api/services/master/warehouse/services-warehouse-impl"
	masterserviceimpl "after-sales/api/services/service-impl"
	transactionworkshopserviceimpl "after-sales/api/services/transaction/workshop/services-workshop-impl"

	_ "after-sales/docs"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func CreateHandler(db *gorm.DB, env string, redis *redis.Client) {
	r := gin.New()
	//mtr_operation_group
	operationGroupRepository := masteroperationrepositoryimpl.StartOperationGroupRepositoryImpl(db)
	operationGroupService := masteroperationserviceimpl.StartOperationGroupService(operationGroupRepository)
	//mtr_operation_section
	operationSectionRepository := masteroperationrepositoryimpl.StartOperationSectionRepositoryImpl(db)
	operationSectionService := masteroperationserviceimpl.StartOperationSectionService(operationSectionRepository)
	//mtr_operation_key
	operationKeyRepository := masteroperationrepositoryimpl.StartOperationKeyRepositoryImpl(db)
	operationKeyService := masteroperationserviceimpl.StartOperationKeyService(operationKeyRepository)
	// //mtr_operation_code
	operationCodeRepository := masteroperationrepositoryimpl.StartOperationCodeRepositoryImpl(db)
	operationCodeService := masteroperationserviceimpl.StartOperationCodeService(operationCodeRepository)
	// //mtr_operation_entries
	operationEntriesRepository := masteroperationrepositoryimpl.StartOperationEntriesRepositoryImpl(db)
	operationEntriesService := masteroperationserviceimpl.StartOperationEntriesService(operationEntriesRepository)
	// //mtr_operation_model_mapping
	operationModelMappingRepository := masteroperationrepositoryimpl.StartOperationModelMappingRepositoryImpl(db)
	operationModelMappingService := masteroperationserviceimpl.StartOperationMappingService(operationModelMappingRepository)

	//mtr_markup_master
	markupMasterRepository := masteritemrepositoryimpl.StartMarkupMasterRepositoryImpl(db)
	markupMasterService := masteritemserviceimpl.StartMarkupMasterService(markupMasterRepository)

	//mtr_uom
	UnitOfMeasurementRepository := masteritemrepositoryimpl.StartUnitOfMeasurementRepositoryImpl(db)
	UnitOfMeasurementService := masteritemserviceimpl.StartUnitOfMeasurementService(UnitOfMeasurementRepository)
	//mtr_item
	itemRepository := masteritemrepositoryimpl.StartItemRepositoryImpl(db, redis)
	itemService := masteritemserviceimpl.StartItemService(itemRepository)
	//mtr_discount
	discountRepository := masterrepositoryimpl.StartDiscountRepositoryImpl(db)
	discountService := masterserviceimpl.StartDiscountService(discountRepository)
	//mtr_incentive_group
	incentiveGroupRepository := masterrepositoryimpl.StartIncentiveGroupImpl(db)
	incentiveGroupService := masterserviceimpl.StartIncentiveGroup(incentiveGroupRepository)

	//mtr_item_class
	itemClassRepository := masteritemrepositoryimpl.StartItemClassRepositoryImpl(db)
	itemClassService := masteritemserviceimpl.StartItemClassService(itemClassRepository)
	//mtr_price_list
	priceListRepository := masteritemrepositoryimpl.StartPriceListRepositoryImpl(db)
	priceListService := masteritemserviceimpl.StartPriceListService(priceListRepository)

	//mtr_discount_percent
	discountPercentRepository := masteritemrepositoryimpl.StartDiscountPercentRepositoryImpl(db)
	discountPercentService := masteritemserviceimpl.StartDiscountPercentService(discountPercentRepository)

	warehouseGroupRepository := masterwarehouserepositoryimpl.OpenWarehouseGroupImpl(db)
	warehouseGroupService := masterwarehouseserviceimpl.OpenWarehouseGroupService(warehouseGroupRepository)

	warehouseMasterRepository := masterwarehouserepositoryimpl.OpenWarehouseMasterImpl(db)
	warehouseMasterService := masterwarehouseserviceimpl.OpenWarehouseMasterService(warehouseMasterRepository)

	warehouseLocationRepository := masterwarehouserepositoryimpl.OpenWarehouseLocationImpl(db)
	warehouseLocationService := masterwarehouseserviceimpl.OpenWarehouseLocationService(warehouseLocationRepository)

	bookingEstimationRepository := transactionworkshoprepositoryimpl.OpenBookingEstimationImpl(db)
	bookingEstimationService := transactionworkshopserviceimpl.OpenBookingEstimationServiceImpl(bookingEstimationRepository)

	r.Use(middlewares.SetupCorsMiddleware())

	if env != "prod" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	api := r.Group("aftersales-service/api/aftersales")
	//master
	mastercontroller.StartDiscountRoutes(db, api, discountService)
	mastercontroller.StartIncentiveGroupRoutes(db, api, incentiveGroupService)
	//operation
	masteroperationcontroller.StartOperationGroupRoutes(db, api, operationGroupService)
	masteroperationcontroller.StartOperationSectionRoutes(db, api, operationSectionService)
	masteroperationcontroller.StartOperationKeyRoutes(db, api, operationKeyService)
	masteroperationcontroller.StartOperationCodeRoutes(db, api, operationCodeService)
	masteroperationcontroller.StartOperationEntriesRoutes(db, api, operationEntriesService)
	masteroperationcontroller.StartOperationModelMappingRoutes(db, api, operationModelMappingService)
	masteritemcontroller.StartUnitOfMeasurementRoutes(db, api, UnitOfMeasurementService)
	masteritemcontroller.StartItemRoutes(db, api, itemService)
	masteritemcontroller.StartItemClassRoutes(db, api, itemClassService)
	masteritemcontroller.StartPriceListRoutes(db, api, priceListService)
	masteritemcontroller.StartMarkupMasterRoutes(db, api, markupMasterService)
	masteritemcontroller.StartDiscountPercentRoutes(db, api, discountPercentService)

	masterwarehousecontroller.OpenWarehouseGroupRoutes(db, api, warehouseGroupService)
	masterwarehousecontroller.OpenWarehouseMasterRoutes(db, api, warehouseMasterService)
	masterwarehousecontroller.OpenWarehouseLocationRoutes(db, api, warehouseLocationService)
	// //transaction
	// transactioncontroller.StartSupplySlipRoutes(db, api.Group("/transaction"), supplySlipService)

	transactionworkshopcontroller.OpenBookingEstimationRoutes(db, api, bookingEstimationService)

	server := &http.Server{Handler: r}
	l, err := net.Listen("tcp4", fmt.Sprintf(":%v", config.EnvConfigs.Port))
	err = server.Serve(l)
	if err != nil {
		log.Fatal(err)
		return
	}
}
