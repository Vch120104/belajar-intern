package masterwarehouserepositoryimpl

import (
	"after-sales/api/config"
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	utils "after-sales/api/utils"
	"errors"
	"math"
	"net/http"
	"strconv"
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
func (r *WarehouseMasterImpl) IsWarehouseMasterByCodeAndCompanyIdExist(tx *gorm.DB, companyId int, warehouseCode string) (bool, *exceptions.BaseErrorResponse) {
	entities := masterwarehouseentities.WarehouseMaster{}

	if err := tx.Model(entities).Where(masterwarehouseentities.WarehouseMaster{CompanyId: companyId, WarehouseCode: warehouseCode}).First(&entities).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil

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

	// Save the warehouseMaster
	if err := tx.Save(&warehouseMaster).Error; err != nil {
		return masterwarehouseentities.WarehouseMaster{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
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

func (r *WarehouseMasterImpl) GetById(tx *gorm.DB, warehouseId int, pagination pagination.Pagination) (masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	var entities masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponse masterwarehousepayloads.GetAllWarehouseMasterResponse
	var getAddressResponse masterwarehousepayloads.AddressResponse
	var getBrandResponse masterwarehousepayloads.BrandResponse
	var getSupplierResponse masterwarehousepayloads.SupplierResponse
	var getUserResponse masterwarehousepayloads.UserResponse
	var getJobPositionResponse masterwarehousepayloads.JobPositionResponse
	var getVillageResponse masterwarehousepayloads.VillageResponse

	// Correct the fetching process
	err := tx.Model(&entities).
		Where("warehouse_id = ?", warehouseId).
		First(&entities).Error

	if err != nil {
		return masterwarehousepayloads.GetAllWarehouseMasterResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
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
			Err:        errors.New("warehouse costing type is not found"),
		}
	}
	// Map the entity to the response payload
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
		WarehouseSalesAllow:           *entities.WarehouseSalesAllow,
		WarehouseInTransit:            *entities.WarehouseInTransit,
		WarehouseName:                 entities.WarehouseName,
		WarehouseDetailName:           entities.WarehouseDetailName,
		WarehouseTransitDefault:       entities.WarehouseTransitDefault,
		WarehouseGroupId:              entities.WarehouseGroupId,
		WarehousePhoneNumber:          entities.WarehousePhoneNumber,
		WarehouseFaxNumber:            entities.WarehouseFaxNumber,
		WarehouseCostingTypeCode:      CostingTypeEntities.WarehouseCostingTypeCode,
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
	SupplierUrl := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(warehouseMasterResponse.SupplierId)
	if err := utils.Get(SupplierUrl, &getSupplierResponse, nil); err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching supplier details",
			Err:        err,
		}
	}

	// Fetch village details
	VillageUrl := config.EnvConfigs.GeneralServiceUrl + "village/" + strconv.Itoa(getAddressResponse.VillageId)
	if err := utils.Get(VillageUrl, &getVillageResponse, nil); err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching village details",
			Err:        err,
		}
	}

	// Fetch user details
	UserUrl := config.EnvConfigs.GeneralServiceUrl + "user-detail/" + strconv.Itoa(warehouseMasterResponse.UserId)
	if err := utils.Get(UserUrl, &getUserResponse, nil); err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching user details",
			Err:        err,
		}
	}

	// Fetch job position details
	JobPositionUrl := config.EnvConfigs.GeneralServiceUrl + "job-position/" + strconv.Itoa(getUserResponse.JobPositionId)
	if err := utils.Get(JobPositionUrl, &getJobPositionResponse, nil); err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when fetching job position details",
			Err:        err,
		}
	}

	// Fetch Authorized User details with pagination
	var AuthorizedUserDetails []masterwarehousepayloads.AuthorizedUserResponse

	var totalRows int64
	query := tx.Table("mtr_warehouse_authorize").
		Select("warehouse_authorize_id, mtr_user_details.user_employee_id as employee_id, mtr_user_details.employee_name as employee_name, mtr_user_details.id_number as id_number").
		Joins("JOIN dms_microservices_general_dev.dbo.mtr_user_details ON mtr_warehouse_authorize.employee_id = mtr_user_details.user_employee_id").
		Where("mtr_warehouse_authorize.warehouse_id = ?", warehouseId)
	if err := query.Count(&totalRows).Error; err != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count authorized user records",
			Err:        err,
		}
	}

	query = query.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	if errAuthorizedUserDetails := query.Find(&AuthorizedUserDetails).Error; errAuthorizedUserDetails != nil {
		return warehouseMasterResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve authorized user from the database",
			Err:        errAuthorizedUserDetails,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))

	// Populate the pagination response
	warehouseMasterResponse.AuthorizedDetails = masterwarehousepayloads.AuthorizedUserDetailsResponse{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: totalPages,
		TotalRows:  int(totalRows),
		Data:       AuthorizedUserDetails,
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
		Select("mtr_warehouse_group.*,mtr_warehouse_master.*").
		Joins("LEFT JOIN mtr_warehouse_group on mtr_warehouse_master.warehouse_group_id = mtr_warehouse_group.warehouse_group_id")

	whereQuery := utils.ApplyFilter(query, filter)

	err := whereQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&response).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve warehouse master from the database",
			Err:        err,
		}
	}

	if len(response) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No warehouse master found",
			Err:        err,
		}
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

func (r *WarehouseMasterImpl) GetWarehouseMasterByCode(tx *gorm.DB, Code string) (masterwarehousepayloads.GetAllWarehouseMasterResponse, *exceptions.BaseErrorResponse) {
	var entities masterwarehouseentities.WarehouseMaster
	var warehouseMasterResponse masterwarehousepayloads.GetAllWarehouseMasterResponse
	var getAddressResponse masterwarehousepayloads.AddressResponse
	var getBrandResponse masterwarehousepayloads.BrandResponse
	var getSupplierResponse masterwarehousepayloads.SupplierResponse
	var getUserResponse masterwarehousepayloads.UserResponse
	var getJobPositionResponse masterwarehousepayloads.JobPositionResponse
	var getVillageResponse masterwarehousepayloads.VillageResponse

	err := tx.Model(&entities).
		Where("warehouse_code = ?", Code).
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

func (r *WarehouseMasterImpl) GetAuthorizeUser(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var entities []masterwarehouseentities.WarehouseAuthorize
	var response []map[string]interface{}

	query := tx.Model(&entities).
		Select("mtr_warehouse_authorize.warehouse_authorize_id, mtr_warehouse_authorize.employee_id, mtr_user_details.employee_name, mtr_user_details.id_number").
		Joins("LEFT JOIN dms_microservices_general_dev.dbo.mtr_user_details ON mtr_warehouse_authorize.employee_id = mtr_user_details.user_employee_id")

	whereQuery := utils.ApplyFilter(query, filterCondition)

	var totalRows int64
	if err := whereQuery.Count(&totalRows).Error; err != nil {
		return response, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to count authorized user records",
			Err:        err,
		}
	}

	if totalRows == 0 {
		return response, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No authorized users found",
			Err:        errors.New("no authorized users found"),
		}
	}

	whereQuery = whereQuery.Offset(pages.GetOffset()).Limit(pages.GetLimit())
	if err := whereQuery.Find(&response).Error; err != nil {
		return response, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve authorized users from the database",
			Err:        err,
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pages.Limit)))

	return response, int(totalRows), totalPages, nil
}

func (r *WarehouseMasterImpl) PostAuthorizeUser(tx *gorm.DB, req masterwarehousepayloads.WarehouseAuthorize) (masterwarehousepayloads.WarehouseAuthorize, *exceptions.BaseErrorResponse) {
	var entities = masterwarehouseentities.WarehouseAuthorize{
		EmployeeId:  req.EmployeeId,
		CompanyId:   req.CompanyId,
		WarehouseId: req.WarehouseId,
	}
	err := tx.Save(&entities).Error

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
