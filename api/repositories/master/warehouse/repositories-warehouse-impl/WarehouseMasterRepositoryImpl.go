package masterwarehouserepositoryimpl

import (
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	utils "after-sales/api/utils"
	"strconv"

	// masterwarehousegroupservice "after-sales/api/services/master/warehouse"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	// "after-sales/api/payloads/pagination"

	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WarehouseMasterImpl struct {
	DB *gorm.DB
}

func OpenWarehouseMasterImpl(db *gorm.DB) masterwarehouserepository.WarehouseMasterRepository {
	return &WarehouseMasterImpl{DB: db}
}

func (r *WarehouseMasterImpl) WithTrx(Trxhandle *gorm.DB) masterwarehouserepository.WarehouseMasterRepository {
	if Trxhandle == nil {
		log.Println("Transaction Database Not Found")
		return r
	}
	r.DB = Trxhandle
	return r
}

func (r *WarehouseMasterImpl) Save(request masterwarehousepayloads.GetWarehouseMasterResponse) (bool, error) {

	var warehouseMaster = masterwarehouseentities.WarehouseMaster{
		IsActive:                      utils.BoolPtr(request.IsActive),
		WarehouseId:                   request.WarehouseId,
		WarehouseCostingType:          request.WarehouseCostingType,
		WarehouseKaroseri:             utils.BoolPtr(request.WarehouseKaroseri),
		WarehouseNegativeStock:        utils.BoolPtr(request.WarehouseNegativeStock),
		WarehouseReplishmentIndicator: utils.BoolPtr(request.WarehouseReplishmentIndicator),
		WarehouseContact:              request.WarehouseContact,
		WarehouseCode:                 request.WarehouseCode,
		AddressId:                     request.AddressId,
		BrandId:                       request.BrandId,
		SupplierId:                    request.SupplierId,
		UserId:                        request.UserId,
		WarehouseSalesAllow:           utils.BoolPtr(request.WarehouseSalesAllow),
		WarehouseInTransit:            utils.BoolPtr(request.WarehouseInTransit),
		WarehouseName:                 request.WarehouseName,
		WarehouseDetailName:           request.WarehouseDetailName,
		WarehouseTransitDefault:       request.WarehouseTransitDefault,
	}

	rows, err := r.DB.Model(&warehouseMaster).
		Save(&warehouseMaster).
		Rows()

	if err != nil {
		return false, err
	}

	defer rows.Close()

	return true, nil
}

func (r *WarehouseMasterImpl) GetById(warehouseId int) (masterwarehousepayloads.GetWarehouseMasterResponse, error) {

	var entities masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponse masterwarehousepayloads.GetWarehouseMasterResponse

	rows, err := r.DB.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseMasterResponse{
			WarehouseId: warehouseId,
		}).
		Scan(&warehouseMasterResponse).
		// Find(&warehouseMasterResponse).
		Rows()

	if err != nil {
		return warehouseMasterResponse, err
	}

	defer rows.Close()

	return warehouseMasterResponse, nil
}

func (r *WarehouseMasterImpl) GetWarehouseWithMultiId(MultiIds []string) ([]masterwarehousepayloads.GetAllWarehouseMasterResponse, error) {
	var entities []masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponse []masterwarehousepayloads.GetAllWarehouseMasterResponse

	rows, err := r.DB.Model(&entities).
		Where("warehouse_id in ?", MultiIds).
		Scan(&warehouseMasterResponse).
		Rows()

	if err != nil {
		return warehouseMasterResponse, err
	}

	defer rows.Close()

	return warehouseMasterResponse, nil
}

func (r *WarehouseMasterImpl) GetAll(request masterwarehousepayloads.GetAllWarehouseMasterRequest, pages pagination.Pagination) (pagination.Pagination, error) {
	var entities []masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponse []masterwarehousepayloads.GetAllWarehouseMasterResponse

	tempRows := r.DB.
		Model(&masterwarehouseentities.WarehouseGroup{}).
		Where("warehouse_name like ?", "%"+request.WarehouseName+"%").
		Where("warehouse_code like ?", "%"+request.WarehouseCode+"%")

	if request.IsActive != "" {
		tempRows = tempRows.Where("is_active = ?", request.IsActive)
	}

	rows, err := tempRows.
		Scopes(pagination.Paginate(entities, &pages, tempRows)).
		Scan(&warehouseMasterResponse).
		Rows()

	if err != nil {
		return pages, err
	}

	defer rows.Close()

	pages.Rows = warehouseMasterResponse
	return pages, nil
}

func (r *WarehouseMasterImpl) GetAllIsActive() ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, error) {
	var warehouseMaster []masterwarehouseentities.WarehouseMaster
	response := []masterwarehousepayloads.IsActiveWarehouseMasterResponse{}

	err := r.DB.Model(&warehouseMaster).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *WarehouseMasterImpl) GetWarehouseMasterByCode(Code string) ([]map[string]interface{}, error) {

	entities := masterwarehouseentities.WarehouseMaster{}
	warehouseMasterResponse := masterwarehousepayloads.GetWarehouseMasterResponse{}
	var getAddressResponse masterwarehousepayloads.AddressResponse
	var getBrandResponse masterwarehousepayloads.BrandResponse
	var getSupplierResponse masterwarehousepayloads.SupplierResponse
	var getUserResponse masterwarehousepayloads.UserResponse
	var getJobPositionResponse masterwarehousepayloads.JobPositionResponse
	var c *gin.Context

	rows, err := r.DB.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseMasterResponse{
			WarehouseCode: Code,
		}).
		First(&warehouseMasterResponse).
		Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// AddressId                     int    `json:"address_id"` http://10.1.32.26:8000/general-service/api/general/address/
	errUrlAddress := utils.Get(c, "http://10.1.32.26:8000/general-service/api/general/address/"+strconv.Itoa(warehouseMasterResponse.AddressId), &getAddressResponse, nil)

	if errUrlAddress != nil {
		return nil, errUrlAddress
	}

	firstJoin := utils.DataFrameLeftJoin([]masterwarehousepayloads.GetWarehouseMasterResponse{warehouseMasterResponse}, []masterwarehousepayloads.AddressResponse{getAddressResponse}, "AddressId")

	// BrandId                       int    `json:"brand_id"` http://10.1.32.26:8000/sales-service/api/sales/unit-brand/
	errUrlBrand := utils.Get(c, "http://10.1.32.26:8000/sales-service/api/sales/unit-brand/"+strconv.Itoa(warehouseMasterResponse.AddressId), &getBrandResponse, nil)

	if errUrlBrand != nil {
		return nil, errUrlBrand
	}

	secondJoin := utils.DataFrameLeftJoin(firstJoin, []masterwarehousepayloads.BrandResponse{getBrandResponse}, "BrandId")

	// SupplierId                    int    `json:"supplier_id"` http://10.1.32.26:8000/general-service/api/general/supplier-master/
	errUrlSupplier := utils.Get(c, "http://10.1.32.26:8000/general-service/api/general/supplier-master/"+strconv.Itoa(warehouseMasterResponse.SupplierId), &getSupplierResponse, nil)

	if errUrlSupplier != nil {
		return nil, errUrlSupplier
	}

	thirdJoin := utils.DataFrameLeftJoin(secondJoin, []masterwarehousepayloads.SupplierResponse{getSupplierResponse}, "SupplierId")

	// UserId                        int    `json:"user_id"` http://10.1.32.26:8000/general-service/api/general/user-details/
	errUrUser := utils.Get(c, "http://10.1.32.26:8000/general-service/api/general/user-details/"+strconv.Itoa(warehouseMasterResponse.UserId), &getUserResponse, nil)

	if errUrUser != nil {
		return nil, errUrUser
	}

	fourthJoin := utils.DataFrameLeftJoin(thirdJoin, []masterwarehousepayloads.UserResponse{getUserResponse}, "UserId")

	// JobPositionId int http://10.1.32.26:8000/general-service/api/general/job-position/
	errUrlJobPosition := utils.Get(c, "http://10.1.32.26:8000/general-service/api/general/job-position/"+strconv.Itoa(getUserResponse.JobPositionId), &getJobPositionResponse, nil)

	if errUrlJobPosition != nil {
		return nil, errUrlJobPosition
	}

	finalJoin := utils.DataFrameLeftJoin(fourthJoin, []masterwarehousepayloads.JobPositionResponse{getJobPositionResponse}, "JobPositionId")

	return finalJoin, nil
}

func (r *WarehouseMasterImpl) ChangeStatus(warehouseId int) (masterwarehousepayloads.GetWarehouseMasterResponse, error) {
	var entities masterwarehouseentities.WarehouseMaster
	var warehouseMasterPayloads masterwarehousepayloads.GetWarehouseMasterResponse

	rows, err := r.DB.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseMasterResponse{
			WarehouseId: warehouseId,
		}).
		Update("is_active", gorm.Expr("1 ^ is_active")).
		Rows()

	if err != nil {
		log.Panic((err.Error()))
	}

	defer rows.Close()

	rows, err = r.DB.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseMasterResponse{
			WarehouseId: warehouseId,
		}).
		// Find(&warehouseMasterPayloads).
		Scan(&warehouseMasterPayloads).
		Rows()

	if err != nil {
		return warehouseMasterPayloads, err
	}

	defer rows.Close()

	return warehouseMasterPayloads, nil
}
