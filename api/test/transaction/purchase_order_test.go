package test

import (
	"after-sales/api/config"
	transactionsparepartpayloads "after-sales/api/payloads/transaction/sparepart"
	transactionsparepartrepositoryimpl "after-sales/api/repositories/transaction/sparepart/repositories-sparepart-impl"
	transactionsparepartservice "after-sales/api/services/transaction/sparepart"
	transactionsparepartserviceimpl "after-sales/api/services/transaction/sparepart/services-sparepart-impl"
	"after-sales/api/utils"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

func setupPurchaseOrder() (*gorm.DB, *redis.Client, transactionsparepartservice.PurchaseOrderService) {
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	PurchaseOrderRepository := transactionsparepartrepositoryimpl.NewPurchaseOrderRepositoryImpl()
	PurchaseOrderService := transactionsparepartserviceimpl.NewPurchaseOrderService(PurchaseOrderRepository, db, rdb)
	return db, rdb, PurchaseOrderService
}
func TestGetAllPurchaseOrder(t *testing.T) {
	//t.Error("das")
	db, rdb, _ := setupPurchaseOrder()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()
	db.Begin()
	//err := db.s
}
func float64Ptr(f float64) *float64 {
	return &f
}
func TestNewPurchaseOrderDetail(t *testing.T) {
	//t.Error("das")
	db, rdb, service := setupPurchaseOrder()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()
	//db.Begin()
	createdDate := time.Now()
	payloads := transactionsparepartpayloads.PurchaseOrderDetailPayloads{
		PurchaseOrderSystemNumber:         1025839,
		PurchaseOrderLine:                 1,
		ItemId:                            1,
		ItemUnitOfMeasurement:             "PC",
		UnitOfMeasurementRate:             float64Ptr(1),
		ItemQuantity:                      float64Ptr(1),
		ItemPrice:                         float64Ptr(1),
		ItemTotal:                         float64Ptr(1),
		PurchaseRequestDetailSystemNumber: 0,
		PurchaseRequestSystemNumber:       0,
		PurchaseRequestLineNumber:         0,
		StockOnHand:                       float64Ptr(1),
		ItemRemark:                        "",
		CreatedByUserId:                   0,
		CreatedDate:                       &createdDate,
		UpdatedByUserId:                   1,
		UpdatedDate:                       &createdDate,
		Snp:                               float64Ptr(01),
		ItemDiscountPercentage:            float64Ptr(1),
		ItemDiscountAmount:                float64Ptr(1),
	}
	_, err := service.NewPurchaseOrderDetail(payloads)
	if err != nil {
		t.Errorf("Failed On: %v", err)
	} else {
		assert.Nil(t, err, nil)
		//assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	}
	//err := db.s
}
func TestDeleteDetailMultiId(t *testing.T) {
	db, rdb, service := setupPurchaseOrder()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()

	id := "999891"
	_, err := service.DeletePurchaseOrderDetailMultiId(id)
	if err != nil {
		t.Errorf("Failed On: %v", err)
	} else {
		assert.Nil(t, err, nil)
		//assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	}
}
func TestUpdatePurchaseOrder(t *testing.T) {
	//t.Error("das")
	db, rdb, service := setupPurchaseOrder()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()
	//db.Begin()
	createdDate := time.Now()
	//payloads := transactionsparepartpayloads.PurchaseOrderDetailPayloads{
	//	PurchaseOrderSystemNumber:         1025839,
	//	PurchaseOrderLine:                 1,
	//	ItemId:                            1,
	//	ItemUnitOfMeasurement:             "PC",
	//	UnitOfMeasurementRate:             float64Ptr(1),
	//	ItemQuantity:                      float64Ptr(1),
	//	ItemPrice:                         float64Ptr(1),
	//	ItemTotal:                         float64Ptr(1),
	//	PurchaseRequestDetailSystemNumber: 0,
	//	PurchaseRequestSystemNumber:       0,
	//	PurchaseRequestLineNumber:         0,
	//	StockOnHand:                       float64Ptr(1),
	//	ItemRemark:                        "",
	//	CreatedByUserId:                   0,
	//	CreatedDate:                       &createdDate,
	//	UpdatedByUserId:                   1,
	//	UpdatedDate:                       &createdDate,
	//	Snp:                               float64Ptr(01),
	//	ItemDiscountPercentage:            float64Ptr(1),
	//	ItemDiscountAmount:                float64Ptr(1),
	//}
	payloads := transactionsparepartpayloads.PurchaseOrderSaveDetailPayloads{
		PurchaseOrderDetailSystemNumber:    1,
		PurchaseOrderSystemNumber:          1025839,
		PurchaseOrderLine:                  1,
		ItemId:                             1,
		ItemUnitOfMeasurement:              "PC",
		UnitOfMeasurementRate:              float64Ptr(1),
		ItemQuantity:                       float64Ptr(1),
		ItemPrice:                          float64Ptr(1),
		ItemTotal:                          float64Ptr(2),
		PurchaseRequestDetailSystemNumber:  0,
		PurchaseRequestSystemNumber:        0,
		PurchaseRequestLineNumber:          0,
		StockOnHand:                        float64Ptr(1),
		ChangedItemPurchaseOrderSystemNo:   0,
		ItemRemark:                         "",
		ChangedItemPurchaseOrderLineNumber: 0,
		CreatedByUserId:                    0,
		CreatedDate:                        &createdDate,
		UpdatedByUserId:                    1,
		UpdatedDate:                        &createdDate,
		Snp:                                float64Ptr(01),
		ItemDiscountPercentage:             float64Ptr(1),
		ItemDiscountAmount:                 float64Ptr(1),
	}
	_, err := service.SavePurchaseOrderDetail(payloads)
	if err != nil {
		t.Errorf("Failed On: %v", err)
	} else {
		assert.Nil(t, err, nil)
		//assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	}
	//err := db.s
}

func TestGeneralBenchmark(t *testing.T) {
	db, rdb, _ := setupPurchaseOrder()
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		rdb.Close()
	}()
	var purchaseRequestStatusDesc transactionsparepartpayloads.PurchaseRequestStatusResponse
	StatusURL := config.EnvConfigs.GeneralServiceUrl + "document-status/20"
	err := utils.Get(StatusURL, &purchaseRequestStatusDesc, nil)
	if err != nil {
		assert.Nil(t, err, nil)

	}
}
