package masteritemrepositoryimpl

import (
	config "after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptionsss_test "after-sales/api/expectionsss"
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

type PurchasePriceRepositoryImpl struct {
}

func StartPurchasePriceRepositoryImpl() masteritemrepository.PurchasePriceRepository {
	return &PurchasePriceRepositoryImpl{}
}

func (r *PurchasePriceRepositoryImpl) GetAllPurchasePrice(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	// Define a slice to hold PurchasePriceRequest responses
	var responses []masteritempayloads.PurchasePriceRequest

	// Define table struct
	tableStruct := masteritempayloads.PurchasePriceRequest{}

	// Define join table
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	// Apply filters
	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	// Execute query
	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	// Define a slice to hold PurchasePriceResponse
	var convertedResponses []masteritempayloads.PurchasePriceResponse

	// Iterate over rows
	for rows.Next() {
		// Define variables to hold row data
		var (
			purchasePriceReq masteritempayloads.PurchasePriceRequest
			purchasePriceRes masteritempayloads.PurchasePriceResponse
		)

		// Scan the row into PurchasePriceRequest struct
		if err := rows.Scan(&purchasePriceReq.PurchasePriceId, &purchasePriceReq.SupplierId, &purchasePriceReq.CurrencyId, &purchasePriceReq.PurchasePriceEffectiveDate, &purchasePriceReq.IsActive); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Supplier data from external service
		SupplierURL := config.EnvConfigs.GeneralServiceUrl + "/api/general/supplier-master/" + strconv.Itoa(purchasePriceReq.SupplierId)
		//fmt.Println("Fetching Supplier data from:", SupplierURL)
		var getSupplierResponse masteritempayloads.PurchasePriceSupplierResponse
		if err := utils.Get(SupplierURL, &getSupplierResponse, nil); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Currency data from external service
		CurrencyURL := config.EnvConfigs.FinanceServiceUrl + "/currency-code/" + strconv.Itoa(purchasePriceReq.CurrencyId)
		//fmt.Println("Fetching Currency data from:", CurrencyURL)
		var getCurrencyResponse masteritempayloads.CurrencyResponse
		if err := utils.Get(CurrencyURL, &getCurrencyResponse, nil); err != nil {
			return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
		// Create PurchasePriceResponse
		purchasePriceRes = masteritempayloads.PurchasePriceResponse{
			PurchasePriceId:            purchasePriceReq.PurchasePriceId,
			SupplierId:                 purchasePriceReq.SupplierId,
			SupplierCode:               getSupplierResponse.SupplierCode, // Set SupplierCode from fetched data
			SupplierName:               getSupplierResponse.SupplierName,
			CurrencyId:                 purchasePriceReq.CurrencyId,
			CurrencyCode:               getCurrencyResponse.CurrencyCode,
			CurrencyName:               getCurrencyResponse.CurrencyName,
			PurchasePriceEffectiveDate: purchasePriceReq.PurchasePriceEffectiveDate,
			IsActive:                   purchasePriceReq.IsActive,
		}

		// Append PurchasePriceResponse to the slice
		convertedResponses = append(convertedResponses, purchasePriceRes)
	}

	// Define a slice to hold map responses
	var mapResponses []map[string]interface{}

	// Iterate over convertedResponses and convert them to maps
	for _, response := range convertedResponses {
		responseMap := map[string]interface{}{
			"purchase_price_id":             response.PurchasePriceId,
			"supplier_id":                   response.SupplierId,
			"supplier_code":                 response.SupplierCode,
			"supplier_name":                 response.SupplierName,
			"currency_id":                   response.CurrencyId,
			"currency_code":                 response.CurrencyCode,
			"currency_name":                 response.CurrencyName,
			"purchase_price_effective_date": response.PurchasePriceEffectiveDate,
			"is_active":                     response.IsActive,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	// Paginate the response data
	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *PurchasePriceRepositoryImpl) SavePurchasePrice(tx *gorm.DB, request masteritempayloads.PurchasePriceRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.PurchasePrice{
		IsActive:                   request.IsActive,
		SupplierId:                 request.SupplierId,
		CurrencyId:                 request.CurrencyId,
		PurchasePriceEffectiveDate: request.PurchasePriceEffectiveDate,
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

func (r *PurchasePriceRepositoryImpl) GetPurchasePriceById(tx *gorm.DB, Id int) (masteritempayloads.PurchasePriceRequest, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.PurchasePrice{}
	response := masteritempayloads.PurchasePriceRequest{}

	err := tx.Model(&entities).
		Where(masteritementities.PurchasePrice{
			PurchasePriceId: Id,
		}).
		First(&response).
		Error

	if err != nil {
		return masteritempayloads.PurchasePriceRequest{}, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("data not found"),
		}
	}

	return response, nil
}

func (r *PurchasePriceRepositoryImpl) GetAllPurchasePriceDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptionsss_test.BaseErrorResponse) {
	// Inisialisasi variabel untuk menyimpan respons dari database dan layanan eksternal
	var responses []masteritempayloads.PurchasePriceDetailResponse
	var getItemResponse []masteritempayloads.PurchasePriceItemResponse

	// Mendapatkan struktur dari tipe data PurchasePriceDetailResponse
	responseStruct := reflect.TypeOf(masteritempayloads.PurchasePriceDetailResponse{})

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
	tableStruct := masteritempayloads.PurchasePriceDetailRequest{}
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
	itemUrl := config.EnvConfigs.AfterSalesServiceUrl + "/item/multi-id/" + strings.Join(itemIds, ",")
	if err := utils.Get(itemUrl, &getItemResponse, nil); err != nil {
		return nil, 0, 0, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Melakukan inner join antara respons lokasi item, respons lokasi item eksternal, dan respons item
	joinedData := utils.DataFrameInnerJoin(responses, getItemResponse, "ItemId")

	// Mem-paginate data yang telah di-join
	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *PurchasePriceRepositoryImpl) AddPurchasePrice(tx *gorm.DB, request masteritempayloads.PurchasePriceDetailRequest) (bool, *exceptionsss_test.BaseErrorResponse) {
	entities := masteritementities.PurchasePriceDetail{
		ItemId:          request.ItemId,
		PurchasePriceId: request.PurchasePriceId,
		PurchasePrice:   request.PurchasePrice,
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

// DeletePurchasePrice deletes an item location by ID
func (r *PurchasePriceRepositoryImpl) DeletePurchasePrice(tx *gorm.DB, Id int) *exceptionsss_test.BaseErrorResponse {
	entities := masteritementities.PurchasePriceDetail{}

	// Menghapus data berdasarkan ID
	err := tx.Where("purchase_price_detail_id = ?", Id).Delete(&entities).Error
	if err != nil {
		return &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// Jika data berhasil dihapus, kembalikan nil untuk error
	return nil
}

func (r *PurchasePriceRepositoryImpl) ChangeStatusPurchasePrice(tx *gorm.DB, Id int) (bool, *exceptionsss_test.BaseErrorResponse) {
	var entity masteritementities.PurchasePrice

	// Cari entitas berdasarkan ID
	result := tx.Model(&entity).
		Where("purchase_price_id = ?", Id).
		First(&entity)

	// Periksa apakah entitas ditemukan
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, &exceptionsss_test.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("purchase price with ID %d not found", Id),
			}
		}
		// Jika ada galat lain, kembalikan galat internal server
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	// Ubah status entitas
	entity.IsActive = !entity.IsActive

	// Simpan perubahan
	result = tx.Save(&entity)
	if result.Error != nil {
		return false, &exceptionsss_test.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return true, nil
}
