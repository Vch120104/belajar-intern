package masterwarehouserepositoryimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	utils "after-sales/api/utils"
	"net/http"
	"strconv"

	// masterwarehousegroupservice "after-sales/api/services/master/warehouse"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	// "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseMasterImpl struct {
}

func OpenWarehouseMasterImpl() masterwarehouserepository.WarehouseMasterRepository {
	return &WarehouseMasterImpl{}
}

func (r *WarehouseMasterImpl) Save(tx *gorm.DB, request masterwarehousepayloads.GetWarehouseMasterResponse) (bool, *exceptionsss_test.BaseErrorResponse) {

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

	rows, err := tx.Model(&warehouseMaster).
		Save(&warehouseMaster).
		Rows()

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return true, nil
}

func (r *WarehouseMasterImpl) DropdownWarehouse(tx *gorm.DB) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse) {

	var warehouseMasterResponse []masterwarehousepayloads.DropdownWarehouseMasterResponse

	err := tx.Model(&masterwarehouseentities.WarehouseMaster{}).
		Select("warehouse_id", "warehouse_code + ' - ' + warehouse_name as warehouse_code").
		Find(&warehouseMasterResponse)
	if err.Error != nil {
		return warehouseMasterResponse, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err.Error,
		}
	}
	return warehouseMasterResponse, nil
}

func (r *WarehouseMasterImpl) GetById(tx *gorm.DB, warehouseId int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse) {

	var entities masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponse masterwarehousepayloads.GetWarehouseMasterResponse

	rows, err := tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseMasterResponse{
			WarehouseId: warehouseId,
		}).
		Scan(&warehouseMasterResponse).
		// Find(&warehouseMasterResponse).
		Rows()

	if err != nil {
		return warehouseMasterResponse, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return warehouseMasterResponse, nil
}

func (r *WarehouseMasterImpl) GetWarehouseWithMultiId(tx *gorm.DB, MultiIds []string) ([]masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse) {

	var entities []masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponse []masterwarehousepayloads.GetAllWarehouseMasterResponse

	rows, err := tx.Model(&entities).
		Where("warehouse_id in ?", MultiIds).
		Scan(&warehouseMasterResponse).
		Rows()

	if err != nil {
		return warehouseMasterResponse, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return warehouseMasterResponse, nil
}

func (r *WarehouseMasterImpl) GetAll(tx *gorm.DB, request masterwarehousepayloads.GetAllWarehouseMasterRequest, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	var entities []masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponses []masterwarehousepayloads.GetAllWarehouseMasterResponse

	// Query untuk mengambil semua entitas gudang sesuai permintaan
	err := tx.Where(&request).
		Find(&entities).
		Error
	if err != nil {
		return pagination.Pagination{}, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	for _, entity := range entities {
		var warehouseMasterResponse masterwarehousepayloads.GetAllWarehouseMasterResponse

		// Salin data entitas gudang ke respons
		warehouseMasterResponse.IsActive = *entity.IsActive
		warehouseMasterResponse.WarehouseId = entity.WarehouseId
		warehouseMasterResponse.WarehouseCostingType = entity.WarehouseCostingType
		warehouseMasterResponse.WarehouseKaroseri = *entity.WarehouseKaroseri
		warehouseMasterResponse.WarehouseNegativeStock = *entity.WarehouseNegativeStock
		warehouseMasterResponse.WarehouseReplishmentIndicator = *entity.WarehouseReplishmentIndicator
		warehouseMasterResponse.WarehouseContact = entity.WarehouseContact
		warehouseMasterResponse.WarehouseCode = entity.WarehouseCode
		warehouseMasterResponse.AddressId = entity.AddressId
		warehouseMasterResponse.BrandId = entity.BrandId
		warehouseMasterResponse.SupplierId = entity.SupplierId
		warehouseMasterResponse.UserId = entity.UserId
		warehouseMasterResponse.WarehouseSalesAllow = *entity.WarehouseSalesAllow
		warehouseMasterResponse.WarehouseInTransit = *entity.WarehouseInTransit
		warehouseMasterResponse.WarehouseName = entity.WarehouseName
		warehouseMasterResponse.WarehouseDetailName = entity.WarehouseDetailName
		warehouseMasterResponse.WarehouseTransitDefault = entity.WarehouseTransitDefault

		// Ambil detail alamat dari layanan API
		addressURL := "http://10.1.32.26:8000/general-service/api/general/address/" + strconv.Itoa(entity.AddressId)
		if err := utils.Get(addressURL, &warehouseMasterResponse.AddressDetails, nil); err != nil {
			return pagination.Pagination{}, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Ambil detail merek dari layanan API
		brandURL := "http://10.1.32.26:8000/sales-service/api/sales/unit-brand/" + strconv.Itoa(entity.BrandId)
		if err := utils.Get(brandURL, &warehouseMasterResponse.BrandDetails, nil); err != nil {
			return pagination.Pagination{}, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Ambil detail pemasok dari layanan API
		supplierURL := "http://10.1.32.26:8000/general-service/api/general/supplier-master/" + strconv.Itoa(entity.SupplierId)
		if err := utils.Get(supplierURL, &warehouseMasterResponse.SupplierDetails, nil); err != nil {
			return pagination.Pagination{}, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Ambil detail pengguna dari layanan API
		userURL := "http://10.1.32.26:8000/general-service/api/general/user-details/" + strconv.Itoa(entity.UserId)
		if err := utils.Get(userURL, &warehouseMasterResponse.UserDetails, nil); err != nil {
			return pagination.Pagination{}, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Ambil detail posisi pekerjaan dari layanan API
		jobPositionURL := "http://10.1.32.26:8000/general-service/api/general/job-position/" + strconv.Itoa(warehouseMasterResponse.UserDetails.JobPositionId)
		if err := utils.Get(jobPositionURL, &warehouseMasterResponse.JobPositionDetails, nil); err != nil {
			return pagination.Pagination{}, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Tambahkan respons ke daftar respons
		warehouseMasterResponses = append(warehouseMasterResponses, warehouseMasterResponse)
	}

	// Setel hasil respons dan kembalikan
	pages.Rows = warehouseMasterResponses
	return pages, nil
}

func (r *WarehouseMasterImpl) GetAllIsActive(tx *gorm.DB) ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse) {

	var warehouseMaster []masterwarehouseentities.WarehouseMaster
	response := []masterwarehousepayloads.IsActiveWarehouseMasterResponse{}

	err := tx.Model(&warehouseMaster).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *WarehouseMasterImpl) GetWarehouseMasterByCode(tx *gorm.DB, Code string) ([]map[string]interface{}, *exceptionsss_test.BaseErrorResponse) {

	entities := masterwarehouseentities.WarehouseMaster{}
	warehouseMasterResponse := masterwarehousepayloads.GetWarehouseMasterResponse{}
	var getAddressResponse masterwarehousepayloads.AddressResponse
	var getBrandResponse masterwarehousepayloads.BrandResponse
	var getSupplierResponse masterwarehousepayloads.SupplierResponse
	var getUserResponse masterwarehousepayloads.UserResponse
	var getJobPositionResponse masterwarehousepayloads.JobPositionResponse

	rows, err := tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseMasterResponse{
			WarehouseCode: Code,
		}).
		First(&warehouseMasterResponse).
		Rows()

	if err != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	// AddressId                     int    `json:"address_id"` http://10.1.32.26:8000/general-service/api/general/address/
	errUrlAddress := utils.Get("http://10.1.32.26:8000/general-service/api/general/address/"+strconv.Itoa(warehouseMasterResponse.AddressId), &getAddressResponse, nil)

	if errUrlAddress != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	firstJoin := utils.DataFrameLeftJoin([]masterwarehousepayloads.GetWarehouseMasterResponse{warehouseMasterResponse}, []masterwarehousepayloads.AddressResponse{getAddressResponse}, "AddressId")

	// BrandId                       int    `json:"brand_id"` http://10.1.32.26:8000/sales-service/api/sales/unit-brand/
	errUrlBrand := utils.Get("http://10.1.32.26:8000/sales-service/api/sales/unit-brand/"+strconv.Itoa(warehouseMasterResponse.AddressId), &getBrandResponse, nil)

	if errUrlBrand != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	secondJoin := utils.DataFrameLeftJoin(firstJoin, []masterwarehousepayloads.BrandResponse{getBrandResponse}, "BrandId")

	// SupplierId                    int    `json:"supplier_id"` http://10.1.32.26:8000/general-service/api/general/supplier-master/
	errUrlSupplier := utils.Get("http://10.1.32.26:8000/general-service/api/general/supplier-master/"+strconv.Itoa(warehouseMasterResponse.SupplierId), &getSupplierResponse, nil)

	if errUrlSupplier != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	thirdJoin := utils.DataFrameLeftJoin(secondJoin, []masterwarehousepayloads.SupplierResponse{getSupplierResponse}, "SupplierId")

	// UserId                        int    `json:"user_id"` http://10.1.32.26:8000/general-service/api/general/user-details/
	errUrUser := utils.Get("http://10.1.32.26:8000/general-service/api/general/user-details/"+strconv.Itoa(warehouseMasterResponse.UserId), &getUserResponse, nil)

	if errUrUser != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	fourthJoin := utils.DataFrameLeftJoin(thirdJoin, []masterwarehousepayloads.UserResponse{getUserResponse}, "UserId")

	// JobPositionId int http://10.1.32.26:8000/general-service/api/general/job-position/
	errUrlJobPosition := utils.Get("http://10.1.32.26:8000/general-service/api/general/job-position/"+strconv.Itoa(getUserResponse.JobPositionId), &getJobPositionResponse, nil)

	if errUrlJobPosition != nil {
		return nil, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	finalJoin := utils.DataFrameLeftJoin(fourthJoin, []masterwarehousepayloads.JobPositionResponse{getJobPositionResponse}, "JobPositionId")

	return finalJoin, nil
}

func (r *WarehouseMasterImpl) ChangeStatus(tx *gorm.DB, warehouseId int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptionsss_test.BaseErrorResponse) {

	var entities masterwarehouseentities.WarehouseMaster
	var warehouseMasterPayloads masterwarehousepayloads.GetWarehouseMasterResponse

	rows, err := tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseMasterResponse{
			WarehouseId: warehouseId,
		}).
		Update("is_active", gorm.Expr("1 ^ is_active")).
		Rows()

	if err != nil {
		return warehouseMasterPayloads, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	rows, err = tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseMasterResponse{
			WarehouseId: warehouseId,
		}).
		// Find(&warehouseMasterPayloads).
		Scan(&warehouseMasterPayloads).
		Rows()

	if err != nil {
		return warehouseMasterPayloads, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return warehouseMasterPayloads, nil
}
