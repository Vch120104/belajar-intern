package masteritemrepositoryimpl

import (
	config "after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ItemLocationRepositoryImpl struct {
}

func StartItemLocationRepositoryImpl() masteritemrepository.ItemLocationRepository {
	return &ItemLocationRepositoryImpl{}
}

func (r *ItemLocationRepositoryImpl) GetAllItemLocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.ItemLocationRequest
	var getWarehouseGroupResponse masteritempayloads.ItemLocWarehouseGroupResponse
	var getWarehouseResponse masteritempayloads.ItemLocationWarehouseResponse
	var getItemResponse []masteritempayloads.ItemLocResponse
	var internalServiceFilter []utils.FilterCondition

	responseStruct := reflect.TypeOf(masteritempayloads.ItemLocationRequest{})

	// Filter internal service conditions
	for _, condition := range filterCondition {
		for j := 0; j < responseStruct.NumField(); j++ {
			if condition.ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, condition)
				break
			}
		}
	}

	// Apply internal service filter conditions
	tableStruct := masteritempayloads.ItemLocationRequest{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)

	// Fetch data from database
	err := whereQuery.Scan(&responses).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to fetch data from database: %w", err),
		}
	}

	// Check if responses are empty
	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no data found"),
		}
	}

	// Define a slice to hold map responses
	var mapResponses []map[string]interface{}

	// Iterate over responses and convert them to maps
	for _, response := range responses {
		responseMap := map[string]interface{}{
			"warehouse_id":       response.WarehouseId,
			"warehouse_group_id": response.WarehouseGroupId,
			"item_id":            response.ItemId,
			"item_location_id":   response.ItemLocationId,
		}

		// Fetch warehouse group data if warehouse group ID is not zero
		if response.WarehouseGroupId != 0 {
			warehouseGroupURL := config.EnvConfigs.AfterSalesServiceUrl + "warehouse-group/by-id/" + strconv.Itoa(response.WarehouseGroupId)
			if err := utils.Get(warehouseGroupURL, &getWarehouseGroupResponse, nil); err != nil {
				return nil, 0, 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			responseMap["warehouse_group_name"] = getWarehouseGroupResponse.WarehouseGroupName
		}

		// Fetch item data if item ID is not zero
		if response.ItemId != 0 {
			itemURL := config.EnvConfigs.AfterSalesServiceUrl + "item/multi-id/" + strconv.Itoa(response.ItemId)
			fmt.Println("Fetching mtr_item data from:", itemURL)
			if err := utils.Get(itemURL, &getItemResponse, nil); err != nil {
				return nil, 0, 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			if len(getItemResponse) > 0 {
				responseMap["item_name"] = getItemResponse[0].ItemName
				responseMap["item_code"] = getItemResponse[0].ItemCode
			}
		}

		// Fetch warehouse data if warehouse ID is not zero
		if response.WarehouseId != 0 {
			warehouseURL := config.EnvConfigs.AfterSalesServiceUrl + "warehouse-master/" + strconv.Itoa(response.WarehouseId)
			fmt.Println("Fetching warehouse_id data from:", warehouseURL)
			if err := utils.Get(warehouseURL, &getWarehouseResponse, nil); err != nil {
				return nil, 0, 0, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusInternalServerError,
					Err:        err,
				}
			}
			responseMap["warehouse_code"] = getWarehouseResponse.WarehouseCode
			responseMap["warehouse_name"] = getWarehouseResponse.WarehouseName
		}

		mapResponses = append(mapResponses, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *ItemLocationRepositoryImpl) SaveItemLocation(tx *gorm.DB, request masteritempayloads.ItemLocationRequest) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemLocation{
		WarehouseGroupId: request.WarehouseGroupId,
		ItemId:           request.ItemId,
		WarehouseId:      request.WarehouseId,
	}

	err := tx.Save(&entities).Error

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

func (r *ItemLocationRepositoryImpl) GetItemLocationById(tx *gorm.DB, Id int) (masteritempayloads.ItemLocationRequest, *exceptions.BaseErrorResponse) {
	entities := masteritementities.ItemLocation{}
	response := masteritempayloads.ItemLocationRequest{}

	err := tx.Model(&entities).
		Where(masteritementities.ItemLocation{
			ItemLocationId: Id,
		}).
		First(&response).
		Error

	if err != nil {
		return masteritempayloads.ItemLocationRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("data not found"),
		}
	}

	return response, nil
}

func (r *ItemLocationRepositoryImpl) GetAllItemLocationDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	// Inisialisasi variabel untuk menyimpan respons dari database dan layanan eksternal
	var responses []masteritempayloads.ItemLocationDetailResponse
	var getItemResponse []masteritempayloads.ItemLocResponse
	var getItemLocResponse []masteritempayloads.ItemLocSourceRequest

	// Mendapatkan struktur dari tipe data ItemLocationDetailResponse
	responseStruct := reflect.TypeOf(masteritempayloads.ItemLocationDetailResponse{})

	// Filter kondisi internal
	var internalServiceFilter []utils.FilterCondition
	for _, condition := range filterCondition {
		for j := 0; j < responseStruct.NumField(); j++ {
			if condition.ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, condition)
				break
			}
		}
	}

	// Menerapkan filter kondisi internal
	tableStruct := masteritempayloads.ItemLocationDetailRequest{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)

	// Mengambil data dari database
	if err := whereQuery.Scan(&responses).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Jika respons dari database kosong, kembalikan error
	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
		}
	}

	// Mengambil data item dari layanan eksternal
	var itemIds []string
	for _, resp := range responses {
		itemIds = append(itemIds, strconv.Itoa(resp.ItemId))
	}
	itemUrl := config.EnvConfigs.AfterSalesServiceUrl + "item/multi-id/" + strings.Join(itemIds, ",")
	if err := utils.Get(itemUrl, &getItemResponse, nil); err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Mengambil data lokasi item dari layanan eksternal
	var itemLocIds []string
	for _, resp := range responses {
		if resp.ItemLocationSourceId != 0 {
			itemLocIds = append(itemLocIds, strconv.Itoa(resp.ItemLocationSourceId))
		}
	}

	// Mengambil data item location source dari layanan eksternal
	for _, id := range itemLocIds {
		itemLocSourceURL := config.EnvConfigs.AfterSalesServiceUrl + "item-location/popup-location?item_location_source_id=" + id
		if err := utils.Get(itemLocSourceURL, &getItemLocResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	// Melakukan inner join antara respons lokasi item, respons lokasi item eksternal, dan respons item
	joinedData, errdf := utils.DataFrameInnerJoin(responses, getItemLocResponse, "ItemLocationSourceId")
	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}
	joinedData, errdf = utils.DataFrameInnerJoin(joinedData, getItemResponse, "ItemId")

	if errdf != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errdf,
		}
	}

	// Mem-paginate data yang telah di-join
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *ItemLocationRepositoryImpl) PopupItemLocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.ItemLocSourceResponse

	// Fetch data from database with joins and conditions
	query := tx.Table("mtr_item_location_source")

	// Apply filter conditions
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

	// Check if responses are empty
	if len(responses) == 0 {
		// notFoundErr := exceptions.NewNotFoundError("No data found")
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	// Perform pagination
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(responses, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *ItemLocationRepositoryImpl) AddItemLocation(tx *gorm.DB, ItemlocId int, request masteritempayloads.ItemLocationDetailRequest) *exceptions.BaseErrorResponse {
	entities := masteritementities.ItemLocationDetail{
		ItemId:                     request.ItemId,
		ItemLocationId:             request.ItemLocationId,
		ItemLocationDetailSourceId: request.ItemLocationSourceId,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

// DeleteItemLocation deletes an item location by ID
func (r *ItemLocationRepositoryImpl) DeleteItemLocation(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	entities := masteritementities.ItemLocationDetail{}

	// Menghapus data berdasarkan ID
	err := tx.Where("item_location_detail_id = ?", Id).Delete(&entities).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Jika data berhasil dihapus, kembalikan nil untuk error
	return nil
}

func (r *ItemLocationRepositoryImpl) GetAllItemLoc(tx *gorm.DB, filtercondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.ItemLocationGetAllResponse

	responseStruct := reflect.TypeOf(masteritempayloads.ItemLocationGetAllResponse{})

	var internalServiceFilter []utils.FilterCondition
	for _, condition := range filtercondition {
		for j := 0; j < responseStruct.NumField(); j++ {
			if condition.ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, condition)
				break
			}
		}
	}

	tableStruct := masteritempayloads.ItemLocationGetAllResponse{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)

	err := whereQuery.Find(&responses).Error
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to fetch data from database: %w", err),
		}
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("no data found"),
		}
	}

	var mapResponses []map[string]interface{}

	// Iterate over responses and convert them to maps
	for _, response := range responses {
		responseMap := map[string]interface{}{
			"item_location_id":        response.ItemLocationId,
			"item_id":                 response.ItemId,
			"item_code":               response.ItemCode,
			"item_name":               response.ItemName,
			"stock_opname":            response.StockOpname,
			"warehouse_id":            response.WarehouseId,
			"warehouse_name":          response.WarehouseName,
			"warehouse_code":          response.WarehouseCode,
			"warehouse_group_id":      response.WarehouseGroupId,
			"warehouse_group_name":    response.WarehouseGroupName,
			"warehouse_group_code":    response.WarehouseGroupCode,
			"warehouse_location_id":   response.WarehouseLocationId,
			"warehouse_location_name": response.WarehouseLocationName,
			"warehouse_location_code": response.WarehouseLocationCode,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
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
