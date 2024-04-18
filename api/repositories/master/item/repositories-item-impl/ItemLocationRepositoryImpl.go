package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	"after-sales/api/exceptions"
	exceptionsss_test "after-sales/api/expectionsss"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	"errors"
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

func (r *ItemLocationRepositoryImpl) GetAllItemLocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var responses []masteritempayloads.ItemLocationRequest
	var getWarehouseGroupResponse masteritempayloads.ItemLocWarehouseGroupResponse
	var getItemResponse []masteritempayloads.ItemLocResponse
	var internalServiceFilter []utils.FilterCondition
	var warehouseGroupId int

	responseStruct := reflect.TypeOf(masteritempayloads.ItemLocationResponse{})

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
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Check if responses are empty
	if len(responses) == 0 {
		notFoundErr := exceptions.NewNotFoundError("No data found")
		panic(notFoundErr)
	}

	// Extract warehouse group ID from the first response
	warehouseGroupId = responses[0].WarehouseGroupId

	// Fetch warehouse group data from external service
	warehouseGroupUrl := "http://localhost:8000/warehouse-group/by-id/" + strconv.Itoa(warehouseGroupId)
	err = utils.Get(warehouseGroupUrl, &getWarehouseGroupResponse, nil)
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Fetch item data from external service
	itemIds := make([]string, len(responses))
	for i, resp := range responses {
		itemIds[i] = strconv.Itoa(resp.ItemId)
	}
	itemUrl := "http://localhost:8000/item/multi-id/" + strings.Join(itemIds, ",")
	err = utils.Get(itemUrl, &getItemResponse, nil)
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Perform inner join between item location responses, warehouse group response, and item response
	joinedData := utils.DataFrameInnerJoin(responses, []masteritempayloads.ItemLocWarehouseGroupResponse{getWarehouseGroupResponse}, "WarehouseGroupId")
	joinedData = utils.DataFrameInnerJoin(joinedData, getItemResponse, "ItemId")

	// Paginate the joined data
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *ItemLocationRepositoryImpl) SaveItemLocation(tx *gorm.DB, request masteritempayloads.ItemLocationRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.ItemLocation{
		WarehouseGroupId: request.WarehouseGroupId,
		ItemId:           request.ItemId,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}

func (r *ItemLocationRepositoryImpl) GetItemLocationById(tx *gorm.DB, Id int) (masteritempayloads.ItemLocationRequest, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.ItemLocation{}
	response := masteritempayloads.ItemLocationRequest{}

	err := tx.Model(&entities).
		Where(masteritementities.ItemLocation{
			ItemLocationId: Id,
		}).
		First(&response).
		Error

	if err != nil {
		return masteritempayloads.ItemLocationRequest{}, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("data not found"),
		}
	}

	return response, nil
}

func (r *ItemLocationRepositoryImpl) GetAllItemLocationDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
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
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Jika respons dari database kosong, kembalikan error
	if len(responses) == 0 {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
		}
	}

	// Mengambil data item dari layanan eksternal
	var itemIds []string
	for _, resp := range responses {
		itemIds = append(itemIds, strconv.Itoa(resp.ItemId))
	}
	itemUrl := "http://localhost:8000/item/multi-id/" + strings.Join(itemIds, ",")
	if err := utils.Get(itemUrl, &getItemResponse, nil); err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
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
		itemLocSourceURL := "http://localhost:8000/item-location/popup-location?item_location_source_id=" + id
		if err := utils.Get(itemLocSourceURL, &getItemLocResponse, nil); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	// Melakukan inner join antara respons lokasi item, respons lokasi item eksternal, dan respons item
	joinedData := utils.DataFrameInnerJoin(responses, getItemLocResponse, "ItemLocationSourceId")
	joinedData = utils.DataFrameInnerJoin(joinedData, getItemResponse, "ItemId")

	// Mem-paginate data yang telah di-join
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *ItemLocationRepositoryImpl) PopupItemLocation(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	var responses []masteritempayloads.ItemLocSourceResponse

	// Fetch data from database with joins and conditions
	query := tx.Table("mtr_item_location_source")

	// Apply filter conditions
	for _, condition := range filterCondition {
		query = query.Where(condition.ColumnField+" = ?", condition.ColumnValue)
	}

	err := query.Find(&responses).Error
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Check if responses are empty
	if len(responses) == 0 {
		notFoundErr := exceptions.NewNotFoundError("No data found")
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        notFoundErr,
		}
	}

	// Perform pagination
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(responses, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *ItemLocationRepositoryImpl) AddItemLocation(tx *gorm.DB, request masteritempayloads.ItemLocationDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.ItemLocationDetail{
		ItemId:                     request.ItemId,
		ItemLocationId:             request.ItemLocationId,
		ItemLocationDetailSourceId: request.ItemLocationSourceId,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return true, nil
}
