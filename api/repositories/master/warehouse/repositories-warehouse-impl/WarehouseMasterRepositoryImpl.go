package masterwarehouserepositoryimpl

import (
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	utils "after-sales/api/utils"
	generalserviceapiutils "after-sales/api/utils/general-service"
	salesserviceapiutils "after-sales/api/utils/sales-service"
	"errors"
	"net/http"
	"strings"

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

// GetWarehouseGroupbyCodeandCompanyId implements masterwarehouserepository.WarehouseMasterRepository.
func (r *WarehouseMasterImpl) GetWarehouseGroupAndMasterbyCodeandCompanyId(tx *gorm.DB, companyId int, warehouseCode string) (int, int, *exceptions.BaseErrorResponse) {
	entities := masterwarehouseentities.WarehouseMaster{}

	if err := tx.Model(entities).Where(masterwarehouseentities.WarehouseMaster{CompanyId: companyId, WarehouseCode: warehouseCode}).First(&entities).Error; err != nil {
		return entities.WarehouseGroupId, entities.WarehouseId, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return entities.WarehouseGroupId, entities.WarehouseId, nil
}

// IsWarehouseMasterByCodeAndCompanyIdExist implements masterwarehouserepository.WarehouseMasterRepository.
func (r *WarehouseMasterImpl) IsWarehouseMasterByCodeAndCompanyIdExist(tx *gorm.DB, companyId int, warehouseCodes []string) ([]masterwarehouseentities.WarehouseMaster, *exceptions.BaseErrorResponse) {
	entities := masterwarehouseentities.WarehouseMaster{}
	response := []masterwarehouseentities.WarehouseMaster{}

	if err := tx.Model(&entities).Where("company_id = ? AND warehouse_code IN ?", companyId, warehouseCodes).Scan(&response).Error; err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *WarehouseMasterImpl) InTransitWarehouseCodeDropdown(tx *gorm.DB, companyID int, warehouseGroupID int) ([]masterwarehousepayloads.DropdownWarehouseMasterByCodeResponse, *exceptions.BaseErrorResponse) {

	var warehouses []masterwarehousepayloads.DropdownWarehouseMasterByCodeResponse

	isTrue := true
	err := tx.Model(&masterwarehouseentities.WarehouseMaster{}).
		Select(`
		warehouse_id,
		warehouse_code,
		warehouse_code + ' - ' + warehouse_name + ' - ' + warehouse_detail_name as warehouse_description`).
		Where(masterwarehouseentities.WarehouseMaster{
			CompanyId:          companyID,
			WarehouseGroupId:   warehouseGroupID,
			WarehouseInTransit: &isTrue}).
		Scan(&warehouses).Error

	if err != nil {
		return warehouses, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return warehouses, nil
}

// DropdownbyGroupId implements masterwarehouserepository.WarehouseMasterRepository.
func (r *WarehouseMasterImpl) DropdownbyGroupId(tx *gorm.DB, warehouseGroupId int, companyId int) ([]masterwarehousepayloads.DropdownWarehouseMasterResponse, *exceptions.BaseErrorResponse) {

	var warehouseMasterResponse []masterwarehousepayloads.DropdownWarehouseMasterResponse

	err := tx.Model(&masterwarehouseentities.WarehouseMaster{}).
		Where(masterwarehouseentities.WarehouseMaster{WarehouseGroupId: warehouseGroupId, CompanyId: companyId}).
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

func (r *WarehouseMasterImpl) Save(tx *gorm.DB, request masterwarehousepayloads.GetWarehouseMasterResponse) (masterwarehouseentities.WarehouseMaster, *exceptions.BaseErrorResponse) {
	var warehouseMaster = masterwarehouseentities.WarehouseMaster{
		CompanyId:                     request.CompanyId,
		IsActive:                      utils.BoolPtr(request.IsActive),
		WarehouseId:                   request.WarehouseId,
		WarehouseCostingTypeId:        request.WarehouseCostingTypeId,
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
		WarehouseGroupId:              request.WarehouseGroupId,
		WarehousePhoneNumber:          request.WarehousePhoneNumber,
		WarehouseFaxNumber:            request.WarehouseFaxNumber,
	}

	if err := tx.Save(&warehouseMaster).Error; err != nil {
		return masterwarehouseentities.WarehouseMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save warehouse master",
			Err:        err,
		}
	}

	return warehouseMaster, nil
}

func (r *WarehouseMasterImpl) Update(tx *gorm.DB, warehouseId int, companyId int, request masterwarehousepayloads.UpdateWarehouseMasterRequest) (masterwarehouseentities.WarehouseMaster, *exceptions.BaseErrorResponse) {
	var warehouseMaster = masterwarehouseentities.WarehouseMaster{
		IsActive:                      utils.BoolPtr(request.IsActive),
		WarehouseCostingTypeId:        request.WarehouseCostingTypeId,
		WarehouseKaroseri:             utils.BoolPtr(request.WarehouseKaroseri),
		WarehouseNegativeStock:        utils.BoolPtr(request.WarehouseNegativeStock),
		WarehouseReplishmentIndicator: utils.BoolPtr(request.WarehouseReplishmentIndicator),
		WarehouseContact:              request.WarehouseContact,
		WarehouseCode:                 request.WarehouseCode,
		AddressId:                     request.AddressId,
		BrandId:                       request.BrandId,
		SupplierId:                    request.SupplierId,
		UserId:                        request.UserId,
		CompanyId:                     request.CompanyId,
		WarehouseId:                   request.WarehouseId,
		WarehouseSalesAllow:           utils.BoolPtr(request.WarehouseSalesAllow),
		WarehouseInTransit:            utils.BoolPtr(request.WarehouseInTransit),
		WarehouseName:                 request.WarehouseName,
		WarehouseDetailName:           request.WarehouseDetailName,
		WarehouseTransitDefault:       request.WarehouseTransitDefault,
		WarehouseGroupId:              request.WarehouseGroupId,
		WarehousePhoneNumber:          request.WarehousePhoneNumber,
		WarehouseFaxNumber:            request.WarehouseFaxNumber,
	}

	if err := tx.Model(&masterwarehouseentities.WarehouseMaster{}).
		Where("warehouse_id = ? AND company_id = ?", warehouseId, companyId).
		Updates(&warehouseMaster).Error; err != nil {
		return masterwarehouseentities.WarehouseMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update warehouse master",
			Err:        err,
		}
	}

	return warehouseMaster, nil
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

	err := tx.Model(&entities).
		Where("warehouse_id = ?", warehouseId).
		First(&entities).Error
	if err != nil {
		return masterwarehousepayloads.GetAllWarehouseMasterResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Warehouse not found",
			Err:        err,
		}
	}

	CostingTypeEntities := masterwarehouseentities.WarehouseCostingType{}
	err = tx.Model(&CostingTypeEntities).
		Where("warehouse_costing_type_id = ?", entities.WarehouseCostingTypeId).
		First(&CostingTypeEntities).Error
	if err != nil {
		return masterwarehousepayloads.GetAllWarehouseMasterResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Warehouse costing type not found",
			Err:        errors.New("warehouse costing type is not found"),
		}
	}

	getAddressResponse, addrErr := generalserviceapiutils.GetAddressById(entities.AddressId)
	var addressDetails masterwarehousepayloads.AddressResponse
	if addrErr != nil {
		addressDetails = masterwarehousepayloads.AddressResponse{
			AddressId:      0,
			AddressStreet1: "",
			AddressStreet2: "",
			AddressStreet3: "",
			VillageId:      0,
		}
	} else {
		addressDetails = masterwarehousepayloads.AddressResponse{
			AddressId:      getAddressResponse.AddressId,
			AddressStreet1: getAddressResponse.AddressStreet1,
			AddressStreet2: getAddressResponse.AddressStreet2,
			AddressStreet3: getAddressResponse.AddressStreet3,
			VillageId:      getAddressResponse.VillageId,
		}
	}

	// Brand
	getBrandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(entities.BrandId)
	var brandDetails masterwarehousepayloads.BrandResponse
	if brandErr != nil {
		brandDetails = masterwarehousepayloads.BrandResponse{
			BrandId:   0,
			BrandCode: "",
			BrandName: "",
		}
	} else {
		brandDetails = masterwarehousepayloads.BrandResponse{
			BrandId:   getBrandResponse.BrandId,
			BrandCode: getBrandResponse.BrandCode,
			BrandName: getBrandResponse.BrandName,
		}
	}

	// Supplier
	getSupplierResponse, supplierErr := generalserviceapiutils.GetSupplierMasterById(entities.SupplierId)
	var supplierDetails masterwarehousepayloads.SupplierResponse
	if supplierErr != nil {
		supplierDetails = masterwarehousepayloads.SupplierResponse{
			SupplierId:   0,
			SupplierName: "",
			SupplierCode: "",
		}
	} else {
		supplierDetails = masterwarehousepayloads.SupplierResponse{
			SupplierId:   getSupplierResponse.SupplierId,
			SupplierName: getSupplierResponse.SupplierName,
			SupplierCode: getSupplierResponse.SupplierCode,
		}
	}

	// Village
	getVillageResponse, villageErr := generalserviceapiutils.GetVillageById(addressDetails.VillageId)
	var villageDetails masterwarehousepayloads.VillageResponse
	if villageErr != nil {
		villageDetails = masterwarehousepayloads.VillageResponse{
			VillageId:      0,
			VillageName:    "",
			DistrictCode:   "",
			DistrictName:   "",
			CityName:       "",
			ProvinceName:   "",
			CountryName:    "",
			VillageZipCode: "",
		}
	} else {
		villageDetails = masterwarehousepayloads.VillageResponse{
			VillageId:      getVillageResponse.VillageId,
			VillageName:    getVillageResponse.VillageName,
			DistrictCode:   getVillageResponse.DistrictCode,
			DistrictName:   getVillageResponse.DistrictName,
			CityName:       getVillageResponse.CityName,
			ProvinceName:   getVillageResponse.ProvinceName,
			CountryName:    getVillageResponse.CountryName,
			VillageZipCode: getVillageResponse.VillageZipCode,
		}
	}

	// User
	getUserCompanyResponse, userErr := generalserviceapiutils.GetUserCompanyAccessById(entities.UserId)
	var userDetails masterwarehousepayloads.UserResponse
	if userErr != nil {
		userDetails = masterwarehousepayloads.UserResponse{
			UserId:        0,
			JobPositionId: 0,
		}
	} else {
		userDetails = masterwarehousepayloads.UserResponse{
			UserId:        getUserCompanyResponse.UserId,
			JobPositionId: getUserCompanyResponse.RoleId,
		}
	}

	// Job Position
	getJobPositionResponse, jobPositionErr := generalserviceapiutils.GetRoleById(userDetails.JobPositionId)
	var jobPositionDetails masterwarehousepayloads.JobPositionResponse
	if jobPositionErr != nil {
		jobPositionDetails = masterwarehousepayloads.JobPositionResponse{
			RolePositionId:   0,
			RolePositionCode: "",
			RolePositionName: "",
		}
	} else {
		jobPositionDetails = masterwarehousepayloads.JobPositionResponse{
			RolePositionId:   getJobPositionResponse.RoleId,
			RolePositionCode: getJobPositionResponse.RoleCode,
			RolePositionName: getJobPositionResponse.RoleName,
		}
	}

	// Now assign the warehouse response fields
	warehouseMasterResponse = masterwarehousepayloads.GetAllWarehouseMasterResponse{
		IsActive:                      *entities.IsActive,
		WarehouseId:                   entities.WarehouseId,
		WarehouseCostingTypeId:        entities.WarehouseCostingTypeId,
		WarehouseKaroseri:             *entities.WarehouseKaroseri,
		WarehouseNegativeStock:        *entities.WarehouseNegativeStock,
		WarehouseReplishmentIndicator: *entities.WarehouseReplishmentIndicator,
		WarehouseContact:              entities.WarehouseContact,
		WarehouseCode:                 entities.WarehouseCode,
		AddressId:                     entities.AddressId,
		BrandId:                       entities.BrandId,
		SupplierId:                    entities.SupplierId,
		UserId:                        entities.UserId,
		CompanyId:                     entities.CompanyId,
		WarehouseSalesAllow:           *entities.WarehouseSalesAllow,
		WarehouseInTransit:            *entities.WarehouseInTransit,
		WarehouseName:                 entities.WarehouseName,
		WarehouseDetailName:           entities.WarehouseDetailName,
		WarehouseTransitDefault:       entities.WarehouseTransitDefault,
		WarehouseGroupId:              entities.WarehouseGroupId,
		WarehousePhoneNumber:          entities.WarehousePhoneNumber,
		WarehouseFaxNumber:            entities.WarehouseFaxNumber,
		WarehouseCostingTypeCode:      CostingTypeEntities.WarehouseCostingTypeCode,

		AddressDetails:     addressDetails,
		BrandDetails:       brandDetails,
		VillageDetails:     villageDetails,
		JobPositionDetails: jobPositionDetails,
	}

	if warehouseMasterResponse.WarehouseKaroseri {
		warehouseMasterResponse.SupplierDetails = supplierDetails

	} else {
		warehouseMasterResponse.UserDetails = userDetails

	}

	return warehouseMasterResponse, nil
}

func (r *WarehouseMasterImpl) GetWarehouseWithMultiId(tx *gorm.DB, MultiIds []int) ([]masterwarehousepayloads.GetAllWarehouseMasterCodeResponse, *exceptions.BaseErrorResponse) {
	var warehouseResponses []masterwarehousepayloads.GetAllWarehouseMasterCodeResponse

	for _, warehouseId := range MultiIds {
		warehouseResponse, err := r.GetById(tx, warehouseId)
		if err != nil {
			return nil, err
		}

		// Map the fields from GetById response to the GetAllWarehouseMasterCodeResponse structure
		warehouseCodeResponse := masterwarehousepayloads.GetAllWarehouseMasterCodeResponse{
			IsActive:                      warehouseResponse.IsActive,
			WarehouseId:                   warehouseResponse.WarehouseId,
			WarehouseCostingTypeId:        warehouseResponse.WarehouseCostingTypeId,
			WarehouseKaroseri:             warehouseResponse.WarehouseKaroseri,
			WarehouseNegativeStock:        warehouseResponse.WarehouseNegativeStock,
			WarehouseReplishmentIndicator: warehouseResponse.WarehouseReplishmentIndicator,
			WarehouseContact:              warehouseResponse.WarehouseContact,
			WarehouseCode:                 warehouseResponse.WarehouseCode,
			AddressId:                     warehouseResponse.AddressId,
			BrandId:                       warehouseResponse.BrandId,
			SupplierId:                    warehouseResponse.SupplierId,
			UserId:                        warehouseResponse.UserId,
			WarehouseSalesAllow:           warehouseResponse.WarehouseSalesAllow,
			WarehouseInTransit:            warehouseResponse.WarehouseInTransit,
			WarehouseName:                 warehouseResponse.WarehouseName,
			WarehouseDetailName:           warehouseResponse.WarehouseDetailName,
			WarehouseTransitDefault:       warehouseResponse.WarehouseTransitDefault,
			WarehouseGroupId:              warehouseResponse.WarehouseGroupId,
			WarehousePhoneNumber:          warehouseResponse.WarehousePhoneNumber,
			WarehouseFaxNumber:            warehouseResponse.WarehouseFaxNumber,
			AddressDetails:                warehouseResponse.AddressDetails,
			BrandDetails:                  warehouseResponse.BrandDetails,
			SupplierDetails:               warehouseResponse.SupplierDetails,
			UserDetails:                   warehouseResponse.UserDetails,
			VillageDetails:                warehouseResponse.VillageDetails,
			JobPositionDetails:            warehouseResponse.JobPositionDetails,
		}

		warehouseResponses = append(warehouseResponses, warehouseCodeResponse)
	}

	return warehouseResponses, nil
}

func (r *WarehouseMasterImpl) GetAll(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	var response []masterwarehousepayloads.GetLookupWarehouseMasterResponse

	baseModelQuery := tx.Model(&masterwarehouseentities.WarehouseMaster{}).
		Select("mtr_warehouse_group.*, mtr_warehouse_master.*").
		Joins("LEFT JOIN mtr_warehouse_group ON mtr_warehouse_master.warehouse_group_id = mtr_warehouse_group.warehouse_group_id")

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	var totalCount int64
	err := whereQuery.Count(&totalCount).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count warehouse master from the database",
			Err:        err,
		}
	}

	err = whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&response).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve warehouse master from the database",
			Err:        err,
		}
	}

	if len(response) == 0 {
		pages.Rows = []masterwarehousepayloads.GetLookupWarehouseMasterResponse{}
		pages.TotalRows = 0
		pages.TotalPages = 0
		return pages, nil
	}

	pages.TotalRows = totalCount
	pages.TotalPages = int(totalCount / int64(pages.Limit))
	if totalCount%int64(pages.Limit) != 0 {
		pages.TotalPages++
	}

	pages.Rows = response

	return pages, nil
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

func (r *WarehouseMasterImpl) GetWarehouseMasterByCode(tx *gorm.DB, Code string) (masterwarehousepayloads.GetAllWarehouseMasterCodeResponse, *exceptions.BaseErrorResponse) {
	var entities masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponse masterwarehousepayloads.GetAllWarehouseMasterCodeResponse

	err := tx.Model(&entities).
		Where("warehouse_code = ?", Code).
		First(&entities).Error
	if err != nil {
		return masterwarehousepayloads.GetAllWarehouseMasterCodeResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Warehouse not found",
			Err:        err,
		}
	}

	CostingTypeEntities := masterwarehouseentities.WarehouseCostingType{}
	err = tx.Model(&CostingTypeEntities).
		Where("warehouse_costing_type_id = ?", entities.WarehouseCostingTypeId).
		First(&CostingTypeEntities).Error
	if err != nil {
		return masterwarehousepayloads.GetAllWarehouseMasterCodeResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Warehouse costing type not found",
			Err:        errors.New("warehouse costing type is not found"),
		}
	}

	// Fetch Address Details
	getAddressResponse, addrErr := generalserviceapiutils.GetAddressById(entities.AddressId)
	var addressDetails masterwarehousepayloads.AddressResponse
	if addrErr != nil {
		addressDetails = masterwarehousepayloads.AddressResponse{}
	} else {
		addressDetails = masterwarehousepayloads.AddressResponse{
			AddressId:      getAddressResponse.AddressId,
			AddressStreet1: getAddressResponse.AddressStreet1,
			AddressStreet2: getAddressResponse.AddressStreet2,
			AddressStreet3: getAddressResponse.AddressStreet3,
			VillageId:      getAddressResponse.VillageId,
		}
	}

	// Fetch Brand Details
	getBrandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(entities.BrandId)
	var brandDetails masterwarehousepayloads.BrandResponse
	if brandErr != nil {
		brandDetails = masterwarehousepayloads.BrandResponse{}
	} else {
		brandDetails = masterwarehousepayloads.BrandResponse{
			BrandId:   getBrandResponse.BrandId,
			BrandCode: getBrandResponse.BrandCode,
			BrandName: getBrandResponse.BrandName,
		}
	}

	// Fetch Supplier Details
	getSupplierResponse, supplierErr := generalserviceapiutils.GetSupplierMasterById(entities.SupplierId)
	var supplierDetails masterwarehousepayloads.SupplierResponse
	if supplierErr != nil {
		supplierDetails = masterwarehousepayloads.SupplierResponse{}
	} else {
		supplierDetails = masterwarehousepayloads.SupplierResponse{
			SupplierId:   getSupplierResponse.SupplierId,
			SupplierName: getSupplierResponse.SupplierName,
			SupplierCode: getSupplierResponse.SupplierCode,
		}
	}

	// Fetch Village Details
	getVillageResponse, villageErr := generalserviceapiutils.GetVillageById(addressDetails.VillageId)
	var villageDetails masterwarehousepayloads.VillageResponse
	if villageErr != nil {
		villageDetails = masterwarehousepayloads.VillageResponse{}
	} else {
		villageDetails = masterwarehousepayloads.VillageResponse{
			VillageId:      getVillageResponse.VillageId,
			VillageName:    getVillageResponse.VillageName,
			DistrictCode:   getVillageResponse.DistrictCode,
			DistrictName:   getVillageResponse.DistrictName,
			CityName:       getVillageResponse.CityName,
			ProvinceName:   getVillageResponse.ProvinceName,
			CountryName:    getVillageResponse.CountryName,
			VillageZipCode: getVillageResponse.VillageZipCode,
		}
	}

	// Fetch User Details
	getUserResponse, userErr := generalserviceapiutils.GetUserCompanyAccessById(entities.UserId)
	var userDetails masterwarehousepayloads.UserResponse
	if userErr != nil {
		userDetails = masterwarehousepayloads.UserResponse{}
	} else {
		userDetails = masterwarehousepayloads.UserResponse{
			UserId:        getUserResponse.UserId,
			JobPositionId: getUserResponse.RoleId,
		}
	}

	// Fetch Job Position Details
	getJobPositionResponse, jobPositionErr := generalserviceapiutils.GetRoleById(userDetails.JobPositionId)
	var jobPositionDetails masterwarehousepayloads.JobPositionResponse
	if jobPositionErr != nil {
		jobPositionDetails = masterwarehousepayloads.JobPositionResponse{}
	} else {
		jobPositionDetails = masterwarehousepayloads.JobPositionResponse{
			RolePositionId:   getJobPositionResponse.RoleId,
			RolePositionCode: getJobPositionResponse.RoleCode,
			RolePositionName: getJobPositionResponse.RoleName,
		}
	}

	// Populate Warehouse Response
	warehouseMasterResponse = masterwarehousepayloads.GetAllWarehouseMasterCodeResponse{
		WarehouseId:              entities.WarehouseId,
		WarehouseCode:            entities.WarehouseCode,
		WarehouseCostingTypeId:   entities.WarehouseCostingTypeId,
		WarehouseCostingTypeCode: CostingTypeEntities.WarehouseCostingTypeCode,
		WarehouseKaroseri:        *entities.WarehouseKaroseri,
		AddressDetails:           addressDetails,
		BrandDetails:             brandDetails,
		VillageDetails:           villageDetails,
	}

	if warehouseMasterResponse.WarehouseKaroseri {
		warehouseMasterResponse.SupplierDetails = supplierDetails
	} else {
		warehouseMasterResponse.UserDetails = userDetails
		warehouseMasterResponse.JobPositionDetails = jobPositionDetails
	}

	return warehouseMasterResponse, nil
}

func (r *WarehouseMasterImpl) GetWarehouseMasterByCodeCompany(tx *gorm.DB, warehouseCode string, companyId int) (masterwarehousepayloads.GetAllWarehouseMasterCodeResponse, *exceptions.BaseErrorResponse) {
	var entities masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponse masterwarehousepayloads.GetAllWarehouseMasterCodeResponse

	err := tx.Model(&entities).
		Where("warehouse_code = ? AND company_id = ?", warehouseCode, companyId).
		First(&warehouseMasterResponse).Error
	if err != nil {
		return masterwarehousepayloads.GetAllWarehouseMasterCodeResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Warehouse not found",
			Err:        err,
		}
	}

	CostingTypeEntities := masterwarehouseentities.WarehouseCostingType{}
	err = tx.Model(&CostingTypeEntities).
		Where("warehouse_costing_type_id = ?", entities.WarehouseCostingTypeId).
		First(&CostingTypeEntities).Error
	if err != nil {
		return masterwarehousepayloads.GetAllWarehouseMasterCodeResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Warehouse costing type not found",
			Err:        errors.New("warehouse costing type is not found"),
		}
	}

	// Fetch Address Details
	getAddressResponse, addrErr := generalserviceapiutils.GetAddressById(entities.AddressId)
	var addressDetails masterwarehousepayloads.AddressResponse
	if addrErr != nil {
		addressDetails = masterwarehousepayloads.AddressResponse{}
	} else {
		addressDetails = masterwarehousepayloads.AddressResponse{
			AddressId:      getAddressResponse.AddressId,
			AddressStreet1: getAddressResponse.AddressStreet1,
			AddressStreet2: getAddressResponse.AddressStreet2,
			AddressStreet3: getAddressResponse.AddressStreet3,
			VillageId:      getAddressResponse.VillageId,
		}
	}

	// Fetch Brand Details
	getBrandResponse, brandErr := salesserviceapiutils.GetUnitBrandById(entities.BrandId)
	var brandDetails masterwarehousepayloads.BrandResponse
	if brandErr != nil {
		brandDetails = masterwarehousepayloads.BrandResponse{}
	} else {
		brandDetails = masterwarehousepayloads.BrandResponse{
			BrandId:   getBrandResponse.BrandId,
			BrandCode: getBrandResponse.BrandCode,
			BrandName: getBrandResponse.BrandName,
		}
	}

	// Fetch Supplier Details
	getSupplierResponse, supplierErr := generalserviceapiutils.GetSupplierMasterById(entities.SupplierId)
	var supplierDetails masterwarehousepayloads.SupplierResponse
	if supplierErr != nil {
		supplierDetails = masterwarehousepayloads.SupplierResponse{}
	} else {
		supplierDetails = masterwarehousepayloads.SupplierResponse{
			SupplierId:   getSupplierResponse.SupplierId,
			SupplierName: getSupplierResponse.SupplierName,
			SupplierCode: getSupplierResponse.SupplierCode,
		}
	}

	// Fetch Village Details
	getVillageResponse, villageErr := generalserviceapiutils.GetVillageById(addressDetails.VillageId)
	var villageDetails masterwarehousepayloads.VillageResponse
	if villageErr != nil {
		villageDetails = masterwarehousepayloads.VillageResponse{}
	} else {
		villageDetails = masterwarehousepayloads.VillageResponse{
			VillageId:      getVillageResponse.VillageId,
			VillageName:    getVillageResponse.VillageName,
			DistrictCode:   getVillageResponse.DistrictCode,
			DistrictName:   getVillageResponse.DistrictName,
			CityName:       getVillageResponse.CityName,
			ProvinceName:   getVillageResponse.ProvinceName,
			CountryName:    getVillageResponse.CountryName,
			VillageZipCode: getVillageResponse.VillageZipCode,
		}
	}

	// Fetch User Details
	getUserResponse, userErr := generalserviceapiutils.GetUserCompanyAccessById(entities.UserId)
	var userDetails masterwarehousepayloads.UserResponse
	if userErr != nil {
		userDetails = masterwarehousepayloads.UserResponse{}
	} else {
		userDetails = masterwarehousepayloads.UserResponse{
			UserId:        getUserResponse.UserId,
			JobPositionId: getUserResponse.RoleId,
		}
	}

	// Fetch Job Position Details
	getJobPositionResponse, jobPositionErr := generalserviceapiutils.GetRoleById(userDetails.JobPositionId)
	var jobPositionDetails masterwarehousepayloads.JobPositionResponse
	if jobPositionErr != nil {
		jobPositionDetails = masterwarehousepayloads.JobPositionResponse{}
	} else {
		jobPositionDetails = masterwarehousepayloads.JobPositionResponse{
			RolePositionId:   getJobPositionResponse.RoleId,
			RolePositionCode: getJobPositionResponse.RoleCode,
			RolePositionName: getJobPositionResponse.RoleName,
		}
	}

	// Populate Warehouse Response
	warehouseMasterResponse = masterwarehousepayloads.GetAllWarehouseMasterCodeResponse{
		WarehouseId:              entities.WarehouseId,
		WarehouseCode:            entities.WarehouseCode,
		WarehouseCostingTypeId:   entities.WarehouseCostingTypeId,
		WarehouseCostingTypeCode: CostingTypeEntities.WarehouseCostingTypeCode,
		WarehouseKaroseri:        *entities.WarehouseKaroseri,
		AddressDetails:           addressDetails,
		BrandDetails:             brandDetails,
		VillageDetails:           villageDetails,
	}

	if warehouseMasterResponse.WarehouseKaroseri {
		warehouseMasterResponse.SupplierDetails = supplierDetails
	} else {
		warehouseMasterResponse.UserDetails = userDetails
		warehouseMasterResponse.JobPositionDetails = jobPositionDetails
	}

	return warehouseMasterResponse, nil
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

func (r *WarehouseMasterImpl) GetAuthorizeUser(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities []masterwarehouseentities.WarehouseAuthorize

	baseModelQuery := tx.Model(&masterwarehouseentities.WarehouseAuthorize{}).
		Select("mtr_warehouse_authorize.warehouse_authorize_id, mtr_warehouse_authorize.employee_id, mtr_user_details.employee_name, mtr_user_details.user_id").
		Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_user_details ON mtr_warehouse_authorize.employee_id = mtr_user_details.user_detail_id")

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	var totalRows int64
	if err := whereQuery.Count(&totalRows).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count authorized user records",
			Err:        err,
		}
	}

	if totalRows == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve authorized users from the database",
			Err:        err,
		}
	}

	var results []map[string]interface{}
	for _, entity := range entities {

		userDetails, err := generalserviceapiutils.GetUserDetailsByID(entity.EmployeeId)
		if err != nil {
			return pages, err
		}

		result := map[string]interface{}{
			"warehouse_authorize_id": entity.WarehouseAuthorizedId,
			"employee_id":            entity.EmployeeId,
			"employee_name":          userDetails.EmployeeName,
			"user_id":                userDetails.UserId,
		}

		results = append(results, result)
	}

	pages.Rows = results
	return pages, nil
}

func (r *WarehouseMasterImpl) PostAuthorizeUser(tx *gorm.DB, req masterwarehousepayloads.WarehouseAuthorize) (masterwarehousepayloads.WarehouseAuthorize, *exceptions.BaseErrorResponse) {
	var existingEntity masterwarehouseentities.WarehouseAuthorize
	err := tx.Where("company_id = ? AND warehouse_id = ? AND employee_id = ?", req.CompanyId, req.WarehouseId, req.EmployeeId).First(&existingEntity).Error
	if err == nil {
		return masterwarehousepayloads.WarehouseAuthorize{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "Duplicate entry: combination of company_id, warehouse_id, and employee_id must be unique",
			Err:        errors.New("duplicate entry : combination of company_id, warehouse_id, and employee_id must be unique"),
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return masterwarehousepayloads.WarehouseAuthorize{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check for existing entry",
			Err:        err,
		}
	}

	var entities = masterwarehouseentities.WarehouseAuthorize{
		EmployeeId:  req.EmployeeId,
		CompanyId:   req.CompanyId,
		WarehouseId: req.WarehouseId,
	}
	err = tx.Save(&entities).Error

	if err != nil {
		return masterwarehousepayloads.WarehouseAuthorize{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Failed to save authorized user",
			Err:        err,
		}
	}
	return req, nil
}

func (r *WarehouseMasterImpl) DeleteMultiIdAuthorizeUser(tx *gorm.DB, id string) (bool, *exceptions.BaseErrorResponse) {
	var authorizeuser masterwarehouseentities.WarehouseAuthorize
	ids := strings.Split(id, ",")

	for _, loop := range ids {
		err := tx.Model(&authorizeuser).Where("warehouse_authorize_id = ?", loop).Delete(&authorizeuser).Error
		if err != nil {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Failed to delete authorized user",
				Err:        err,
			}
		}
	}
	return true, nil
}
