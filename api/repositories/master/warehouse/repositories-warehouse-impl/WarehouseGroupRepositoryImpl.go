package masterwarehouserepositoryimpl

import (
	exceptions "after-sales/api/exceptions"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	utils "after-sales/api/utils"
	"errors"
	"net/http"

	// masterwarehousegroupservice "after-sales/api/services/master/warehouse"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	// "after-sales/api/payloads/pagination"

	"gorm.io/gorm"
)

type WarehouseGroupImpl struct {
}

func OpenWarehouseGroupImpl() masterwarehouserepository.WarehouseGroupRepository {
	return &WarehouseGroupImpl{}
}

// GetbyGroupCode implements masterwarehouserepository.WarehouseGroupRepository.
func (r *WarehouseGroupImpl) GetbyGroupCode(tx *gorm.DB, groupCode string) (masterwarehousepayloads.GetWarehouseGroupResponse, *exceptions.BaseErrorResponse) {
	entity := masterwarehouseentities.WarehouseGroup{}
	response := masterwarehousepayloads.GetWarehouseGroupResponse{}

	err := tx.Model(&entity).Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupCode: groupCode}).
		First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return response, nil
}

// GetWarehouseGroupDropdownbyId implements masterwarehouserepository.WarehouseGroupRepository.
func (r *WarehouseGroupImpl) GetWarehouseGroupDropdownbyId(tx *gorm.DB, Id int) (masterwarehousepayloads.GetWarehouseGroupDropdown, *exceptions.BaseErrorResponse) {
	entity := masterwarehouseentities.WarehouseGroup{}
	response := masterwarehousepayloads.GetWarehouseGroupDropdown{}

	err := tx.Model(&entity).Where(masterwarehouseentities.WarehouseGroup{WarehouseGroupId: Id}).
		First(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	return response, nil
}

// GetWarehouseGroupDropdown implements masterwarehouserepository.WarehouseGroupRepository.
func (r *WarehouseGroupImpl) GetWarehouseGroupDropdown(tx *gorm.DB) ([]masterwarehousepayloads.GetWarehouseGroupDropdown, *exceptions.BaseErrorResponse) {
	entity := masterwarehouseentities.WarehouseGroup{}
	response := []masterwarehousepayloads.GetWarehouseGroupDropdown{}

	err := tx.Model(&entity).
		Scan(&response).Error

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if len(response) == 0 {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New(""),
		}
	}

	return response, nil
}

func (r *WarehouseGroupImpl) SaveWarehouseGroup(tx *gorm.DB, request masterwarehousepayloads.GetWarehouseGroupResponse) (bool, *exceptions.BaseErrorResponse) {

	var warehouseGroup = masterwarehouseentities.WarehouseGroup{
		IsActive:           utils.BoolPtr(request.IsActive),
		WarehouseGroupId:   request.WarehouseGroupId,
		WarehouseGroupCode: request.WarehouseGroupCode,
		WarehouseGroupName: request.WarehouseGroupName,
		ProfitCenterId:     request.ProfitCenterId,
	}

	rows, err := tx.Model(&warehouseGroup).
		Save(&warehouseGroup).
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

func (r *WarehouseGroupImpl) GetByIdWarehouseGroup(tx *gorm.DB, warehouseGroupId int) (masterwarehousepayloads.GetWarehouseGroupResponse, *exceptions.BaseErrorResponse) {
	entity := masterwarehouseentities.WarehouseGroup{}
	response := masterwarehousepayloads.GetWarehouseGroupResponse{}

	rows, err := tx.Model(&entity).
		Where("warehouse_group_id = ?", warehouseGroupId).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *WarehouseGroupImpl) GetAllWarehouseGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	entities := []masterwarehouseentities.WarehouseGroup{}

	//ON PROGRESS JOIN WAREHOUSMASTER AND WAREHOUSELOCATION
	baseModelQuery := tx.Model(&entities).Joins("WarehouseMaster")

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if len(entities) == 0 {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (r *WarehouseGroupImpl) ChangeStatusWarehouseGroup(tx *gorm.DB, warehouseGroupId int) (bool, *exceptions.BaseErrorResponse) {
	var entities masterwarehouseentities.WarehouseGroup
	var warehouseGroupPayloads masterwarehousepayloads.GetWarehouseGroupResponse

	rows, err := tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseGroupResponse{
			WarehouseGroupId: warehouseGroupId,
		}).
		Update("is_active", gorm.Expr("1 ^ is_active")).
		Rows()

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	rows, err = tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseGroupResponse{
			WarehouseGroupId: warehouseGroupId,
		}).
		Find(&warehouseGroupPayloads).
		Scan(&warehouseGroupPayloads).
		Rows()

	if err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return false, nil
}
