package masterwarehouserepositoryimpl

import (
	exceptionsss_test "after-sales/api/expectionsss"
	masterwarehousepayloads "after-sales/api/payloads/master/warehouse"
	"after-sales/api/payloads/pagination"
	masterwarehouserepository "after-sales/api/repositories/master/warehouse"
	utils "after-sales/api/utils"
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

func (r *WarehouseGroupImpl) SaveWarehouseGroup(tx *gorm.DB, request masterwarehousepayloads.GetWarehouseGroupResponse) (bool, *exceptionsss_test.BaseErrorResponse) {

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
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return true, nil
}

func (r *WarehouseGroupImpl) GetByIdWarehouseGroup(tx *gorm.DB, warehouseGroupId int) (masterwarehousepayloads.GetWarehouseGroupResponse, *exceptionsss_test.BaseErrorResponse) {
	entity := masterwarehouseentities.WarehouseGroup{}
	response := masterwarehousepayloads.GetWarehouseGroupResponse{}

	rows, err := tx.Model(&entity).
		Where("warehouse_group_id = ?",warehouseGroupId).
		First(&response).
		Rows()

	if err != nil {
		return response, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	return response, nil
}

func (r *WarehouseGroupImpl) GetAllWarehouseGroup(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptionsss_test.BaseErrorResponse) {
	entities := []masterwarehouseentities.WarehouseGroup{}

	baseModelQuery := tx.Model(&entities)

	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	rows, err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, whereQuery)).Scan(&entities).Rows()

	if err != nil {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if len(entities) == 0 {
		return pages, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	pages.Rows = entities

	return pages, nil
}

func (r *WarehouseGroupImpl) ChangeStatusWarehouseGroup(tx *gorm.DB, warehouseGroupId int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entities masterwarehouseentities.WarehouseGroup
	var warehouseGroupPayloads masterwarehousepayloads.GetWarehouseGroupResponse

	rows, err := tx.Model(&entities).
		Where(masterwarehousepayloads.GetWarehouseGroupResponse{
			WarehouseGroupId: warehouseGroupId,
		}).
		Update("is_active", gorm.Expr("1 ^ is_active")).
		Rows()

	if err != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
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
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	defer rows.Close()

	return false, nil
}
