package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type ItemLocationRepositoryImpl struct {
}

func StartItemLocationRepositoryImpl() masteritemrepository.ItemLocationRepository {
	return &ItemLocationRepositoryImpl{}
}

func (r *ItemLocationRepositoryImpl) GetAllItemLocationDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {

	entities := []masteritementities.ItemLocationDetail{}
	query := tx.Model(&masteritementities.ItemLocationDetail{})

	whereQuery := utils.ApplyFilter(query, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&entities).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to fetch data from database: %w", err),
		}
	}

	if len(entities) == 0 {
		pages.Rows = []masteritempayloads.ItemLocationDetailResponse{}
		return pages, nil
	}

	pages.Rows = entities
	return pages, nil
}

func (r *ItemLocationRepositoryImpl) PopupItemLocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.ItemLocSourceResponse

	query := tx.Table("mtr_item_location_source")

	for _, condition := range filterCondition {
		query = query.Where(condition.ColumnField+" = ?", condition.ColumnValue)
	}

	err := query.Find(&responses).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(responses, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *ItemLocationRepositoryImpl) AddItemLocation(tx *gorm.DB, ItemlocId int, request masteritempayloads.ItemLocationDetailRequest) (masteritementities.ItemLocationDetail, *exceptions.BaseErrorResponse) {

	var existing masteritementities.ItemLocationDetail
	if err := tx.Where("item_id = ?  AND warehouse_location_id = ?", request.ItemId, request.ItemLocationId).First(&existing).Error; err == nil {
		return masteritementities.ItemLocationDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "combination Item and Warehouse Location already exists",
			Err:        errors.New("combination Item and Warehouse Location already exists"),
		}
	}

	entities := masteritementities.ItemLocationDetail{
		ItemId:                     request.ItemId,
		ItemLocationId:             request.ItemLocationId,
		ItemLocationDetailSourceId: request.ItemLocationSourceId,
	}

	err := tx.Save(&entities).Error
	if err != nil {
		return masteritementities.ItemLocationDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to save data",
			Err:        err,
		}
	}

	return entities, nil
}

// DeleteItemLocation deletes an item location by ID
func (r *ItemLocationRepositoryImpl) DeleteItemLocation(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	entities := masteritementities.ItemLocationDetail{}

	err := tx.Where("item_location_detail_id = ?", Id).Delete(&entities).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

func (r *ItemLocationRepositoryImpl) GetAllItemLoc(tx *gorm.DB, filterConditions []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var results []masteritempayloads.ItemLocationGetAllResponse

	// Join the ItemLocation table with parent entities (Item, Warehouse, WarehouseGroup, WarehouseLocation)
	baseModelQuery := tx.Model(&masteritementities.ItemLocation{}).
		Joins("JOIN mtr_item ON mtr_item.item_id = mtr_location_item.item_id").
		Joins("JOIN mtr_warehouse_master ON mtr_warehouse_master.warehouse_id = mtr_location_item.warehouse_id").
		Joins("JOIN mtr_warehouse_group ON mtr_warehouse_group.warehouse_group_id = mtr_location_item.warehouse_group_id").
		Joins("JOIN mtr_warehouse_location ON mtr_warehouse_location.warehouse_location_id = mtr_location_item.warehouse_location_id")

	whereQuery := utils.ApplyFilter(baseModelQuery, filterConditions)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&results).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to fetch data from database: %w", err),
		}
	}

	if len(results) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var itemLocations []masteritempayloads.ItemLocationGetAllResponse
	for _, result := range results {

		// Ambil data Item berdasarkan ItemId
		var itemResponse masteritementities.Item
		if err := tx.Where("item_id = ?", result.ItemId).First(&itemResponse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item not found",
					Err:        fmt.Errorf("item with ID %d not found", result.ItemId),
				}
			}
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item",
				Err:        err,
			}
		}

		var warehouseResponse masterwarehouseentities.WarehouseMaster
		if err := tx.Where("warehouse_id = ?", result.WarehouseId).First(&warehouseResponse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Warehouse not found",
					Err:        fmt.Errorf("warehouse with ID %d not found", result.WarehouseId),
				}
			}
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Warehouse",
				Err:        err,
			}
		}

		var locationResponse masterwarehouseentities.WarehouseLocation
		if err := tx.Where("warehouse_location_id = ?", result.WarehouseLocationId).First(&locationResponse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Warehouse not found",
					Err:        fmt.Errorf("warehouse with ID %d not found", result.WarehouseLocationId),
				}
			}
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Warehouse Location",
				Err:        err,
			}
		}

		var warehouseGroupResponse masterwarehouseentities.WarehouseGroup
		if err := tx.Where("warehouse_group_id = ?", result.WarehouseGroupId).First(&warehouseGroupResponse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Warehouse group not found",
					Err:        fmt.Errorf("warehouse group with ID %d not found", result.WarehouseGroupId),
				}
			}
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Warehouse group",
				Err:        err,
			}
		}

		itemLocation := masteritempayloads.ItemLocationGetAllResponse{
			ItemLocationId:        result.ItemLocationId,
			ItemId:                result.ItemId,
			ItemCode:              itemResponse.ItemCode,
			ItemName:              itemResponse.ItemName,
			StockOpname:           result.StockOpname,
			WarehouseId:           result.WarehouseId,
			WarehouseCode:         warehouseResponse.WarehouseCode,
			WarehouseName:         warehouseResponse.WarehouseName,
			WarehouseGroupId:      result.WarehouseGroupId,
			WarehouseGroupCode:    warehouseGroupResponse.WarehouseGroupCode,
			WarehouseGroupName:    warehouseGroupResponse.WarehouseGroupName,
			WarehouseLocationId:   result.WarehouseLocationId,
			WarehouseLocationCode: locationResponse.WarehouseLocationCode,
			WarehouseLocationName: locationResponse.WarehouseLocationName,
		}
		itemLocations = append(itemLocations, itemLocation)
	}

	pages.Rows = itemLocations

	return pages, nil
}

func (r *ItemLocationRepositoryImpl) GetByIdItemLoc(tx *gorm.DB, id int) (masteritempayloads.ItemLocationGetByIdResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemLocation{}
	response := masteritempayloads.ItemLocationGetByIdResponse{}

	result := tx.Model(&entities).Select("mtr_location_item.*,mtr_item.item_name,mtr_item.item_code,mtr_warehouse_location.warehouse_location_code,mtr_warehouse_location.warehouse_location_name").
		Where("item_location_id=?", id).
		Joins("Join mtr_item on mtr_item.item_id = mtr_location_item.item_id").
		Joins("Join mtr_warehouse_location on mtr_warehouse_location.warehouse_location_id=mtr_location_item.warehouse_location_id").
		Where("mtr_location_item.item_location_id=?", id).Scan(&response)

	if result.Error != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no data found"),
		}
	}
	return response, nil
}

func (r *ItemLocationRepositoryImpl) SaveItemLoc(tx *gorm.DB, req masteritempayloads.SaveItemlocation) (masteritementities.ItemLocation, *exceptions.BaseErrorResponse) {
	var existing masteritementities.ItemLocation
	if err := tx.Where("item_id = ? AND warehouse_location_id = ?", req.ItemId, req.WarehouseLocationId).First(&existing).Error; err == nil {
		return masteritementities.ItemLocation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Message:    "combination Item, and Warehouse Location already exists",
			Err:        errors.New("combination Item,  and Warehouse Location already exists"),
		}
	}

	entities := masteritementities.ItemLocation{
		ItemLocationId:      req.ItemLocationId,
		WarehouseGroupId:    req.WarehouseGroupId,
		ItemId:              req.ItemId,
		WarehouseId:         req.WarehouseId,
		WarehouseLocationId: req.WarehouseLocationId,
	}
	err := tx.Save(&entities).Error
	if err != nil {
		return masteritementities.ItemLocation{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusConflict,
			Err:        err,
		}
	}
	return entities, nil
}

func (r *ItemLocationRepositoryImpl) DeleteItemLoc(tx *gorm.DB, ids []int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.ItemLocation
	if err := tx.Delete(&entities, ids).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *ItemLocationRepositoryImpl) IsDuplicateItemLoc(tx *gorm.DB, warehouseId int, warehouseLocationId int, itemId int) (bool, error) {
	entities := masteritementities.ItemLocation{}
	responses := []masteritementities.ItemLocation{}

	err := tx.Model(&entities).Where(
		masteritementities.ItemLocation{
			WarehouseId:         warehouseId,
			WarehouseLocationId: warehouseLocationId,
			ItemId:              itemId,
		}).Scan(&responses).Error

	if err != nil {
		return true, err
	}

	if len(responses) > 0 {
		return true, nil
	}
	return false, nil
}
