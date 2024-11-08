package masterwarehouserepositoryimpl

import (
	// masterwarehousepayloads "after-sales/api/payloads/master/warehouse"

	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	pagination "after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	utils "after-sales/api/utils"
	"errors"
	"net/http"
	"strings"

	// utils "after-sales/api/utils"

	// masterwarehousegroupservice "after-sales/api/services/master/warehouse"

	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	// "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseLocationImpl struct {
}

func OpenWarehouseLocationImpl() masterwarehouserepository.WarehouseLocationRepository {
	return &WarehouseLocationImpl{}
}

// ProcessWarehouseLocationTemplate implements masterwarehouserepository.WarehouseLocationRepository.
// CheckIfLocationExist implements masterwarehouserepository.WarehouseLocationRepository.
func (r *WarehouseLocationImpl) CheckIfLocationExist(tx *gorm.DB, warehouseCode string, locationCode string, locationName string) (bool, *exceptions.BaseErrorResponse) {
	entities := masterwarehouseentities.WarehouseMaster{}
	// warehouseGroup := masterwarehouseentities.WarehouseGroup{}
	response := masterwarehouseentities.WarehouseLocation{}

	if err := tx.Model(entities).Select("1").Where(masterwarehouseentities.WarehouseMaster{WarehouseCode: warehouseCode}).
		InnerJoins("WarehouseGroup", tx.Select("1")).
		InnerJoins("WarehouseGroup.WarehouseLocation", tx.Where(masterwarehouseentities.WarehouseLocation{WarehouseLocationCode: locationCode, WarehouseLocationName: locationName})).
		First(&response).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}

	return true, nil
}

func (r *WarehouseLocationImpl) Save(tx *gorm.DB, request masterwarehouseentities.WarehouseLocation) (bool, *exceptions.BaseErrorResponse) {

	err := tx.Save(&request).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}

func (r *WarehouseLocationImpl) GetById(tx *gorm.DB, warehouseLocationId int) (masterwarehousepayloads.GetAllWarehouseLocationResponse, *exceptions.BaseErrorResponse) {

	var entities masterwarehouseentities.WarehouseLocation
	var warehouseLocationResponse masterwarehousepayloads.GetAllWarehouseLocationResponse

	err := tx.Model(entities).
		Select(`"mtr_warehouse_location"."is_active",
		"mtr_warehouse_location"."warehouse_location_id",
		"mtr_warehouse_master"."company_id",
		"mtr_warehouse_master"."warehouse_id",
		"mtr_warehouse_location"."warehouse_group_id",
		"mtr_warehouse_location"."warehouse_location_code",
		"mtr_warehouse_location"."warehouse_location_name",
		"mtr_warehouse_location"."warehouse_location_detail_name",
		"mtr_warehouse_location"."warehouse_location_pick_sequence",
		"mtr_warehouse_location"."warehouse_location_capacity_in_m3",
		mtr_warehouse_group.warehouse_group_id,
        mtr_warehouse_group.warehouse_group_name,
        mtr_warehouse_group.warehouse_group_code,
		mtr_warehouse_master.warehouse_code,
		mtr_warehouse_master.warehouse_name`).
		Joins(" JOIN mtr_warehouse_group ON mtr_warehouse_location.warehouse_group_id = mtr_warehouse_group.warehouse_group_id").
		Joins(" JOIN mtr_warehouse_master ON mtr_warehouse_group.warehouse_group_id = mtr_warehouse_master.warehouse_group_id").
		Where(masterwarehouseentities.WarehouseLocation{WarehouseLocationId: warehouseLocationId}).
		First(&warehouseLocationResponse).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return warehouseLocationResponse, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "warehouse location not found",
				Err:        err,
			}
		}
		return warehouseLocationResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return warehouseLocationResponse, nil
}

func (r *WarehouseLocationImpl) GetByCode(tx *gorm.DB, warehouseLocationCode string) (masterwarehousepayloads.GetAllWarehouseLocationResponse, *exceptions.BaseErrorResponse) {
	entities := masterwarehouseentities.WarehouseLocation{}
	response := masterwarehousepayloads.GetAllWarehouseLocationResponse{}

	err := tx.Model(&entities).
		Select(`
			mtr_warehouse_location.is_active,
			mtr_warehouse_location.warehouse_location_id,
			mtr_warehouse_master.company_id,
			mtr_warehouse_location.warehouse_group_id,
			mtr_warehouse_location.warehouse_location_code,
			mtr_warehouse_location.warehouse_location_name,
			mtr_warehouse_location.warehouse_location_detail_name,
			mtr_warehouse_location.warehouse_location_pick_sequence,
			mtr_warehouse_location.warehouse_location_capacity_in_m3,
			mtr_warehouse_group.warehouse_group_id,
			mtr_warehouse_group.warehouse_group_name,
			mtr_warehouse_group.warehouse_group_code,
			mtr_warehouse_master.warehouse_code,
			mtr_warehouse_master.warehouse_name
		`).
		Joins("INNER JOIN mtr_warehouse_group ON mtr_warehouse_location.warehouse_group_id = mtr_warehouse_group.warehouse_group_id").
		Joins("INNER JOIN mtr_warehouse_master ON mtr_warehouse_group.warehouse_group_id = mtr_warehouse_master.warehouse_group_id").
		Where("mtr_warehouse_location.warehouse_location_code = ?", warehouseLocationCode).
		First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return response, nil
}

func (r *WarehouseLocationImpl) GetAll(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities masterwarehouseentities.WarehouseLocation
	var responses []masterwarehousepayloads.GetAllWarehouseLocationResponse

	query := tx.Model(entities).
		Select(`"mtr_warehouse_location"."is_active",
		"mtr_warehouse_location"."warehouse_location_id",
		mtr_warehouse_master.company_id,
		"mtr_warehouse_location"."warehouse_group_id",
		"mtr_warehouse_location"."warehouse_id",
		"mtr_warehouse_location"."warehouse_location_code",
		"mtr_warehouse_location"."warehouse_location_name",
		"mtr_warehouse_location"."warehouse_location_detail_name",
		"mtr_warehouse_location"."warehouse_location_pick_sequence",
		"mtr_warehouse_location"."warehouse_location_capacity_in_m3",
		mtr_warehouse_group.warehouse_group_id,
        mtr_warehouse_group.warehouse_group_name,
        mtr_warehouse_group.warehouse_group_code,
		mtr_warehouse_master.warehouse_code,
		mtr_warehouse_master.warehouse_name`).
		Joins("JOIN mtr_warehouse_group ON mtr_warehouse_location.warehouse_group_id = mtr_warehouse_group.warehouse_group_id").
		Joins("JOIN mtr_warehouse_master ON mtr_warehouse_location.warehouse_id  = mtr_warehouse_master.warehouse_id")

	filterQuery := utils.ApplyFilter(query, filter)

	err := filterQuery.Scopes(pagination.Paginate(&entities, &pages, filterQuery)).Scan(&responses).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	pages.Rows = responses

	return pages, nil
}

func (r *WarehouseLocationImpl) ChangeStatus(tx *gorm.DB, warehouseLocationId int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterwarehouseentities.WarehouseLocation
	var warehouseLocationPayloads masterwarehousepayloads.GetWarehouseLocationResponse

	rows, err := tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseLocationResponse{
			WarehouseLocationId: warehouseLocationId,
		}).
		Update("is_active", gorm.Expr("1 ^ is_active")).
		Rows()

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	rows, err = tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseLocationResponse{
			WarehouseLocationId: warehouseLocationId,
		}).
		// Find(&warehouseMasterPayloads).
		Scan(&warehouseLocationPayloads).
		Rows()

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return true, nil
}
