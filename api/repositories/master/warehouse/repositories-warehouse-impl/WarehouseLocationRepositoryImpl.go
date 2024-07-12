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

func (r *WarehouseLocationImpl) Save(tx *gorm.DB, request masterwarehousepayloads.GetWarehouseLocationResponse) (bool, *exceptions.BaseErrorResponse) {

	var warehouseMaster = masterwarehouseentities.WarehouseLocation{
		IsActive:                      utils.BoolPtr(request.IsActive),
		WarehouseLocationId:           request.WarehouseLocationId,
		CompanyId:                     request.CompanyId,
		WarehouseGroupId:              request.WarehouseGroupId,
		WarehouseLocationCode:         request.WarehouseLocationCode,
		WarehouseLocationName:         request.WarehouseLocationName,
		WarehouseLocationDetailName:   request.WarehouseLocationDetailName,
		WarehouseLocationPickSequence: request.WarehouseLocationPickSequence,
		WarehouseLocationCapacityInM3: request.WarehouseLocationCapacityInM3,
	}

	err := tx.Save(&warehouseMaster).Error

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

func (r *WarehouseLocationImpl) GetById(tx *gorm.DB, warehouseLocationId int) (masterwarehousepayloads.GetWarehouseLocationResponse, *exceptions.BaseErrorResponse) {

	var entities masterwarehouseentities.WarehouseLocation
	var warehouseLocationResponse masterwarehousepayloads.GetWarehouseLocationResponse

	rows, err := tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseLocationResponse{
			WarehouseLocationId: warehouseLocationId,
		}).
		First(&warehouseLocationResponse).
		// Find(&warehouseMasterResponse).
		Rows()

	if err != nil {
		return warehouseLocationResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return warehouseLocationResponse, nil
}

func (r *WarehouseLocationImpl) GetAll(tx *gorm.DB, filter []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities masterwarehouseentities.WarehouseLocation
	var responses []masterwarehousepayloads.GetAllWarehouseLocationResponse

	query := tx.Model(entities).
		Select(`"mtr_warehouse_location"."is_active",
		"mtr_warehouse_location"."warehouse_location_id",
		"mtr_warehouse_location"."company_id",
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
		Joins("LEFT OUTER JOIN mtr_warehouse_group ON mtr_warehouse_location.warehouse_group_id = mtr_warehouse_group.warehouse_group_id").
		Joins("LEFT OUTER JOIN mtr_warehouse_master ON mtr_warehouse_group.warehouse_group_id = mtr_warehouse_master.warehouse_group_id")

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
