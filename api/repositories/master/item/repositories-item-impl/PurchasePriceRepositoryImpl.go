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
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PurchasePriceRepositoryImpl struct {
}

func StartPurchasePriceRepositoryImpl() masteritemrepository.PurchasePriceRepository {
	return &PurchasePriceRepositoryImpl{}
}

func (r *PurchasePriceRepositoryImpl) GetAllPurchasePrice(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var responses []masteritempayloads.PurchasePriceRequest

	tableStruct := masteritempayloads.PurchasePriceRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}
	defer rows.Close()

	var convertedResponses []masteritempayloads.PurchasePriceResponse

	for rows.Next() {
		var (
			purchasePriceReq masteritempayloads.PurchasePriceRequest
			purchasePriceRes masteritempayloads.PurchasePriceResponse
		)

		if err := rows.Scan(&purchasePriceReq.PurchasePriceId, &purchasePriceReq.SupplierId, &purchasePriceReq.CurrencyId, &purchasePriceReq.PurchasePriceEffectiveDate, &purchasePriceReq.IsActive); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Supplier data from external service
		SupplierURL := config.EnvConfigs.GeneralServiceUrl + "supplier-master/" + strconv.Itoa(purchasePriceReq.SupplierId)
		//fmt.Println("Fetching Supplier data from:", SupplierURL)
		var getSupplierResponse masteritempayloads.PurchasePriceSupplierResponse
		if err := utils.Get(SupplierURL, &getSupplierResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Currency data from external service
		CurrencyURL := config.EnvConfigs.FinanceServiceUrl + "currency-code/" + strconv.Itoa(purchasePriceReq.CurrencyId)
		//fmt.Println("Fetching Currency data from:", CurrencyURL)
		var getCurrencyResponse masteritempayloads.CurrencyResponse
		if err := utils.Get(CurrencyURL, &getCurrencyResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

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
			IdentitySysNumber:          0,
		}

		convertedResponses = append(convertedResponses, purchasePriceRes)
	}

	var mapResponses []map[string]interface{}

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
			"identity_sys_number":           response.IdentitySysNumber,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
}

func (r *PurchasePriceRepositoryImpl) SavePurchasePrice(tx *gorm.DB, request masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PurchasePrice{
		IsActive:                   request.IsActive,
		SupplierId:                 request.SupplierId,
		CurrencyId:                 request.CurrencyId,
		PurchasePriceEffectiveDate: request.PurchasePriceEffectiveDate,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return masteritementities.PurchasePrice{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return masteritementities.PurchasePrice{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return entities, nil
}

func (r *PurchasePriceRepositoryImpl) GetPurchasePriceById(tx *gorm.DB, Id int) (masteritempayloads.PurchasePriceRequest, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PurchasePrice{}
	response := masteritempayloads.PurchasePriceRequest{}

	err := tx.Model(&entities).
		Where(masteritementities.PurchasePrice{
			PurchasePriceId: Id,
		}).
		First(&response).
		Error

	if err != nil {
		return masteritempayloads.PurchasePriceRequest{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        errors.New("data not found"),
		}
	}

	return response, nil
}

func (r *PurchasePriceRepositoryImpl) GetAllPurchasePriceDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.PurchasePriceDetailResponse
	var getItemResponse []masteritempayloads.PurchasePriceItemResponse

	responseStruct := reflect.TypeOf(masteritempayloads.PurchasePriceDetailResponse{})

	var internalServiceFilter []utils.FilterCondition
	for _, condition := range filterCondition {
		for j := 0; j < responseStruct.NumField(); j++ {
			if condition.ColumnField == responseStruct.Field(j).Tag.Get("parent_entity")+"."+responseStruct.Field(j).Tag.Get("json") {
				internalServiceFilter = append(internalServiceFilter, condition)
				break
			}
		}
	}

	tableStruct := masteritempayloads.PurchasePriceDetailRequest{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, internalServiceFilter)

	if err := whereQuery.Scan(&responses).Error; err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
		}
	}

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

	joinedData := utils.DataFrameInnerJoin(responses, getItemResponse, "ItemId")

	dataPaginate, totalPages, totalRows := pagination.NewDataFramePaginate(joinedData, &pages)

	return dataPaginate, totalPages, totalRows, nil
}

func (r *PurchasePriceRepositoryImpl) GetPurchasePriceDetailById(tx *gorm.DB, Id int, pages pagination.Pagination) (map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var getHeaderRequest masteritempayloads.PurchasePriceRequest
	var getDetailRequest []masteritempayloads.PurchasePriceDetailRequest
	var getItemResponse []masteritempayloads.PurchasePriceItemResponse

	if err := tx.Table("mtr_purchase_price").Where("purchase_price_id = ?", Id).First(&getHeaderRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
			}
		}
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if err := tx.Table("mtr_purchase_price_detail").Where("purchase_price_id = ?", Id).Find(&getDetailRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
			}
		}
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	var itemIds []string
	for _, detail := range getDetailRequest {
		itemIds = append(itemIds, strconv.Itoa(detail.ItemId))
	}
	itemUrl := config.EnvConfigs.AfterSalesServiceUrl + "item/multi-id/" + strings.Join(itemIds, ",")
	if err := utils.Get(itemUrl, &getItemResponse, nil); err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	itemMap := make(map[int]masteritempayloads.PurchasePriceItemResponse)
	for _, item := range getItemResponse {
		itemMap[item.ItemId] = item
	}

	var detailResponses []map[string]interface{}
	for _, detail := range getDetailRequest {
		item := itemMap[detail.ItemId]
		detailResponse := map[string]interface{}{
			"purchase_price_detail_id": detail.PurchasePriceDetailId,
			"purchase_price_id":        detail.PurchasePriceId,
			"item_id":                  detail.ItemId,
			"is_active":                detail.IsActive,
			"purchase_price":           detail.PurchasePrice,
			"item_code":                item.ItemCode,
			"item_name":                item.ItemName,
		}
		detailResponses = append(detailResponses, detailResponse)
	}

	totalRows := len(detailResponses)
	page := pages.GetPage()
	limit := pages.GetLimit()
	start := page * limit
	end := start + limit

	if start > totalRows {
		detailResponses = []map[string]interface{}{}
	} else {
		if end > totalRows {
			end = totalRows
		}
		detailResponses = detailResponses[start:end]
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"currency_id":                   getHeaderRequest.CurrencyId,
			"is_active":                     getHeaderRequest.IsActive,
			"purchase_price_effective_date": getHeaderRequest.PurchasePriceEffectiveDate.Format(time.RFC3339),
			"purchase_price_id":             getHeaderRequest.PurchasePriceId,
			"supplier_id":                   getHeaderRequest.SupplierId,
		},
	}

	detailData := map[string]interface{}{
		"data": detailResponses,
	}
	response["data"].(map[string]interface{})["zdetails"] = detailData

	return response, totalPages, totalRows, nil
}

func (r *PurchasePriceRepositoryImpl) AddPurchasePrice(tx *gorm.DB, request masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PurchasePriceDetail{
		ItemId:          request.ItemId,
		PurchasePriceId: request.PurchasePriceId,
		PurchasePrice:   request.PurchasePrice,
		IsActive:        request.IsActive,
	}

	err := tx.Save(&entities).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return masteritementities.PurchasePriceDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusConflict,
				Err:        err,
			}
		} else {

			return masteritementities.PurchasePriceDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}
	}

	return entities, nil
}

func (r *PurchasePriceRepositoryImpl) DeletePurchasePrice(tx *gorm.DB, Id int) *exceptions.BaseErrorResponse {
	entities := masteritementities.PurchasePriceDetail{}

	result := tx.Where("purchase_price_detail_id = ?", Id).First(&entities)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        result.Error,
		}
	} else if result.Error != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	err := tx.Delete(&entities).Error
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

func (r *PurchasePriceRepositoryImpl) ChangeStatusPurchasePrice(tx *gorm.DB, Id int) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {
	var entity masteritementities.PurchasePrice

	result := tx.Model(&entity).
		Where("purchase_price_id = ?", Id).
		First(&entity)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return masteritementities.PurchasePrice{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("purchase price with ID %d not found", Id),
			}
		}
		return masteritementities.PurchasePrice{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	entity.IsActive = !entity.IsActive

	result = tx.Save(&entity)
	if result.Error != nil {
		return masteritementities.PurchasePrice{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entity, nil
}
