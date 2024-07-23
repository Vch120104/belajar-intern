package masterwarehouserepositoryimpl

import (
	"after-sales/api/config"
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	utils "after-sales/api/utils"
	"errors"
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

// DropdownbyGroupId implements masterwarehouserepository.WarehouseMasterRepository.
func (r *WarehouseMasterImpl) DropdownbyGroupId(tx *gorm.DB, warehouseGroupId int) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse) {

	var warehouseMasterResponse []masterwarehousepayloads.DropdownWarehouseMasterResponse

	err := tx.Model(&masterwarehouseentities.WarehouseMaster{}).Where(masterwarehouseentities.WarehouseMaster{WarehouseGroupId: warehouseGroupId}).
		Select("warehouse_id", "warehouse_code + ' - ' + warehouse_name as warehouse_code").
		Find(&warehouseMasterResponse)
	if err.Error != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err.Error,
		}
	}

	if len(warehouseMasterResponse) == 0 {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}
	return warehouseMasterResponse, nil
}

func (r *WarehouseMasterImpl) Save(tx *gorm.DB, request masterwarehousepayloads.GetWarehouseMasterResponse) (bool, *exceptions.BaseErrorResponse) {

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
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return true, nil
}

func (r *WarehouseMasterImpl) DropdownWarehouse(tx *gorm.DB) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse) {

	var warehouseMasterResponse []masterwarehousepayloads.DropdownWarehouseMasterResponse

	err := tx.Model(&masterwarehouseentities.WarehouseMaster{}).
		Select("warehouse_id", "warehouse_code + ' - ' + warehouse_name as warehouse_code").
		Find(&warehouseMasterResponse)
	if err.Error != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err.Error,
		}
	}
	return warehouseMasterResponse, nil
}

func (r *WarehouseMasterImpl) GetById(tx *gorm.DB, warehouseId int) (masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	var entities masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponse masterwarehousepayloads.GetAllWarehouseMasterResponse
	var getAddressResponse masterwarehousepayloads.AddressResponse
	var getBrandResponse masterwarehousepayloads.BrandResponse
	var getSupplierResponse masterwarehousepayloads.SupplierResponse
	var getUserResponse masterwarehousepayloads.UserResponse
	var getJobPositionResponse masterwarehousepayloads.JobPositionResponse
	var getVillageResponse masterwarehousepayloads.VillageResponse

	err := tx.Model(&entities).
		Where("warehouse_id = ?", warehouseId).
		First(&warehouseMasterResponse).Error

	if err != nil {
		return masterwarehousepayloads.GetAllWarehouseMasterResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	// Fetch address details
	AddressUrl := config.EnvConfigs.GeneralServiceUrl + "address/" + strconv.Itoa(warehouseMasterResponse.AddressId)
	if err := utils.Get(AddressUrl, &getAddressResponse, nil); err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching address details",
			Err:        err,
		}
	}

	// Fetch brand details
	BrandUrl := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(warehouseMasterResponse.BrandId)
	if err := utils.Get(BrandUrl, &getBrandResponse, nil); err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching brand details",
			Err:        err,
		}
	}

	// Fetch supplier details
	SupplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(warehouseMasterResponse.SupplierId)
	if err := utils.Get(SupplierUrl, &getSupplierResponse, nil); err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching supplier details",
			Err:        err,
		}
	}

	// fetch village details
	VillageUrl := config.EnvConfigs.GeneralServiceUrl + "village/" + strconv.Itoa(getAddressResponse.VillageId)
	if err := utils.Get(VillageUrl, &getVillageResponse, nil); err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching village details",
			Err:        err,
		}
	}

	// Fetch user details
	UserUrl := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(warehouseMasterResponse.UserId)
	if err := utils.Get(UserUrl, &getUserResponse, nil); err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching user details",
			Err:        err,
		}
	}

	// Fetch job position details
	JobPositionUrl := config.EnvConfigs.GeneralServiceUrl + "/job-position/" + strconv.Itoa(getUserResponse.JobPositionId)
	if err := utils.Get(JobPositionUrl, &getJobPositionResponse, nil); err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching job position details",
			Err:        err,
		}
	}

	// Populate the nested fields
	warehouseMasterResponse.AddressDetails = getAddressResponse
	warehouseMasterResponse.BrandDetails = getBrandResponse
	warehouseMasterResponse.SupplierDetails = getSupplierResponse
	warehouseMasterResponse.UserDetails = getUserResponse
	warehouseMasterResponse.JobPositionDetails = getJobPositionResponse
	warehouseMasterResponse.VillageDetails = getVillageResponse

	return warehouseMasterResponse, nil
}

func (r *WarehouseMasterImpl) GetWarehouseWithMultiId(tx *gorm.DB, MultiIds []string) ([]masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse) {

	var entities []masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponse []masterwarehousepayloads.GetAllWarehouseMasterResponse

	rows, err := tx.Model(&entities).
		Where("warehouse_id in ?", MultiIds).
		Scan(&warehouseMasterResponse).
		Rows()

	if err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return warehouseMasterResponse, nil
}

func (r *WarehouseMasterImpl) GetAll(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities masterwarehouseentities.WarehouseMaster
	response := []masterwarehousepayloads.GetLookupWarehouseMasterResponse{}
	query := tx.Model(entities).
		Select("*").
		Joins("LEFT JOIN mtr_warehouse_group on mtr_warehouse_master.warehouse_group_id = mtr_warehouse_group.warehouse_group_id")

	whereQuery := utils.ApplyFilter(query, filter)

	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&response).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(response) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	pages.Rows = response

	return pages, nil

	// 	// Ambil detail alamat dari layanan API
	// 	addressURL := config.EnvConfigs.GeneralServiceUrl + "address/" + strconv.Itoa(entity.AddressId)
	// 	if err := utils.Get(addressURL, &warehouseMasterResponse.AddressDetails, nil); err != nil {
	// 		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
	// 			StatusCode: http.StatusInternalServerError,
	// 			Err:        err,
	// 		}
	// 	}

	// 	// Ambil detail merek dari layanan API
	// 	brandURL := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(entity.BrandId)
	// 	if err := utils.Get(brandURL, &warehouseMasterResponse.BrandDetails, nil); err != nil {
	// 		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
	// 			StatusCode: http.StatusInternalServerError,
	// 			Err:        err,
	// 		}
	// 	}

	// 	// Ambil detail pemasok dari layanan API
	// 	supplierURL := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(entity.SupplierId)
	// 	if err := utils.Get(supplierURL, &warehouseMasterResponse.SupplierDetails, nil); err != nil {
	// 		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
	// 			StatusCode: http.StatusInternalServerError,
	// 			Err:        err,
	// 		}
	// 	}

	// 	// Ambil detail pengguna dari layanan API
	// 	userURL := config.EnvConfigs.GeneralServiceUrl + "user-details/" + strconv.Itoa(entity.UserId)
	// 	if err := utils.Get(userURL, &warehouseMasterResponse.UserDetails, nil); err != nil {
	// 		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
	// 			StatusCode: http.StatusInternalServerError,
	// 			Err:        err,
	// 		}
	// 	}

	// 	// Ambil detail posisi pekerjaan dari layanan API
	// 	jobPositionURL := config.EnvConfigs.GeneralServiceUrl + "job-position/" + strconv.Itoa(warehouseMasterResponse.UserDetails.JobPositionId)
	// 	if err := utils.Get(jobPositionURL, &warehouseMasterResponse.JobPositionDetails, nil); err != nil {
	// 		return pagination.Pagination{}, &exceptions.BaseErrorResponse{
	// 			StatusCode: http.StatusInternalServerError,
	// 			Err:        err,
	// 		}
	// 	}

	// 	// Tambahkan respons ke daftar respons
	// 	warehouseMasterResponses = append(warehouseMasterResponses, warehouseMasterResponse)
	// }

	// // Setel hasil respons dan kembalikan
	// pages.Rows = warehouseMasterResponses
	// return pages, nil
}

func (r *WarehouseMasterImpl) GetAllIsActive(tx *gorm.DB) ([]masterwarehousepayloads.IsActiveWarehouseMasterResponse, *exceptions.BaseErrorResponse) {

	var warehouseMaster []masterwarehouseentities.WarehouseMaster
	response := []masterwarehousepayloads.IsActiveWarehouseMasterResponse{}

	err := tx.Model(&warehouseMaster).Where("is_active = 'true'").Scan(&response).Error

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *WarehouseMasterImpl) GetWarehouseMasterByCode(tx *gorm.DB, Code string) ([]map[string]interface{}, *exceptions.BaseErrorResponse) {

	entities := masterwarehouseentities.WarehouseMaster{}
	warehouseMasterResponse := masterwarehousepayloads.GetWarehouseMasterResponse{}
	var getAddressResponse masterwarehousepayloads.AddressResponse
	var getBrandResponse masterwarehousepayloads.BrandResponse
	var getSupplierResponse masterwarehousepayloads.SupplierResponse
	// var getUserResponse masterwarehousepayloads.UserResponse
	// var getJobPositionResponse masterwarehousepayloads.JobPositionResponse

	rows, err := tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseMasterResponse{
			WarehouseCode: Code,
		}).
		First(&warehouseMasterResponse).
		Rows()

	if err != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	// AddressId                     int    `json:"address_id"` http://172.16.5.101/general-service/api/general/address/
	errUrlAddress := utils.Get(config.EnvConfigs.GeneralServiceUrl+"address/"+strconv.Itoa(warehouseMasterResponse.AddressId), &getAddressResponse, nil)

	if errUrlAddress != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	firstJoin := utils.DataFrameLeftJoin([]masterwarehousepayloads.GetWarehouseMasterResponse{warehouseMasterResponse}, []masterwarehousepayloads.AddressResponse{getAddressResponse}, "AddressId")

	// BrandId                       int    `json:"brand_id"` http://172.16.5.101/sales-service/api/sales/unit-brand/
	errUrlBrand := utils.Get(config.EnvConfigs.SalesServiceUrl+"unit-brand/"+strconv.Itoa(warehouseMasterResponse.AddressId), &getBrandResponse, nil)

	if errUrlBrand != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	secondJoin := utils.DataFrameLeftJoin(firstJoin, []masterwarehousepayloads.BrandResponse{getBrandResponse}, "BrandId")

	// SupplierId                    int    `json:"supplier_id"` http://172.16.5.101/general-service/api/general/supplier-master/
	errUrlSupplier := utils.Get(config.EnvConfigs.GeneralServiceUrl+"supplier-master/"+strconv.Itoa(warehouseMasterResponse.SupplierId), &getSupplierResponse, nil)

	if errUrlSupplier != nil {
		return nil, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	thirdJoin := utils.DataFrameLeftJoin(secondJoin, []masterwarehousepayloads.SupplierResponse{getSupplierResponse}, "SupplierId")

	// ON DEVELOPMENT

	// UserId                        int    `json:"user_id"` http://172.16.5.101/general-service/api/general/user-details/
	// errUrUser := utils.Get(config.EnvConfigs.GeneralServiceUrl+"user-details/"+strconv.Itoa(warehouseMasterResponse.UserId), &getUserResponse, nil)

	// if errUrUser != nil {
	// 	return nil, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Err:        err,
	// 	}
	// }

	// fourthJoin := utils.DataFrameLeftJoin(thirdJoin, []masterwarehousepayloads.UserResponse{getUserResponse}, "UserId")

	// JobPositionId int http://172.16.5.101/general-service/api/general/job-position/
	// errUrlJobPosition := utils.Get(config.EnvConfigs.GeneralServiceUrl+"job-position/"+strconv.Itoa(getUserResponse.JobPositionId), &getJobPositionResponse, nil)

	// if errUrlJobPosition != nil {
	// 	return nil, &exceptions.BaseErrorResponse{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Err:        err,
	// 	}
	// }

	// finalJoin := utils.DataFrameLeftJoin(thirdJoin, []masterwarehousepayloads.JobPositionResponse{getJobPositionResponse}, "JobPositionId")

	return thirdJoin, nil
}

func (r *WarehouseMasterImpl) ChangeStatus(tx *gorm.DB, warehouseId int) (masterwarehousepayloads.GetWarehouseMasterResponse, *exceptions.BaseErrorResponse) {

	var entities masterwarehouseentities.WarehouseMaster
	var warehouseMasterPayloads masterwarehousepayloads.GetWarehouseMasterResponse

	rows, err := tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseMasterResponse{
			WarehouseId: warehouseId,
		}).
		Update("is_active", gorm.Expr("1 ^ is_active")).
		Rows()

	if err != nil {
		return warehouseMasterPayloads, &exceptions.BaseErrorResponse{
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
		return warehouseMasterPayloads, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return warehouseMasterPayloads, nil
}
