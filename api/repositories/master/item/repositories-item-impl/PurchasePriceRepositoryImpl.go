package masteritemrepositoryimpl

import (
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	utils "after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type PurchasePriceRepositoryImpl struct {
}

// return false if purchase price detail does not exists
func (r *PurchasePriceRepositoryImpl) CheckPurchasePriceDetailExistence(tx *gorm.DB, Id int, itemId int) (bool, *exceptions.BaseErrorResponse) {
	var entities masteritementities.PurchasePriceDetail
	var count int64

	err := tx.Model(&entities).Where(masteritementities.PurchasePriceDetail{PurchasePriceId: Id, ItemId: itemId}).Count(&count).Error
	if err != nil {
		return true, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error when checking purchase price detail existence",
			Err:        err,
		}
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

func StartPurchasePriceRepositoryImpl() masteritemrepository.PurchasePriceRepository {
	return &PurchasePriceRepositoryImpl{}
}

func (r *PurchasePriceRepositoryImpl) GetAllPurchasePrice(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.PurchasePriceRequest

	internalFilter, externalFilter := utils.DefineInternalExternalFilter(filterCondition, masteritementities.PurchasePrice{})

	baseModelQuery := tx.Model(&masteritementities.PurchasePrice{})

	baseModelQuery = utils.ApplyFilter(baseModelQuery, internalFilter)

	var supplierCode, supplierName, currencyCode, currencyName string

	for _, filter := range externalFilter {
		switch {
		case filter.ColumnField == "supplier_code":
			supplierCode = filter.ColumnValue
		case filter.ColumnField == "supplier_name":
			supplierName = filter.ColumnValue
		case filter.ColumnField == "currency_code":
			currencyCode = filter.ColumnValue
		case filter.ColumnField == "currency_name":
			currencyName = filter.ColumnValue
		}
	}

	if supplierCode != "" || supplierName != "" {
		supplierParams := generalserviceapiutils.SupplierMasterParams{
			Page:         0,
			Limit:        1000,
			SupplierCode: supplierCode,
			SupplierName: supplierName,
		}

		supplierResponse, supplierError := generalserviceapiutils.GetAllSupplierMaster(supplierParams)
		if supplierError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: supplierError.StatusCode,
				Message:    "Error fetching supplier data",
				Err:        supplierError.Err,
			}
		}

		var supplierIds []int
		for _, supplier := range supplierResponse {
			supplierIds = append(supplierIds, supplier.SupplierId)
		}
		if len(supplierIds) > 0 {
			baseModelQuery = baseModelQuery.Where("supplier_id IN ?", supplierIds)
		} else {
			pages.Rows = []map[string]interface{}{}
			return pages, nil
		}
	}

	if currencyCode != "" || currencyName != "" {
		currencyParams := financeserviceapiutils.CurrencyParams{
			CurrencyCode: currencyCode,
			CurrencyName: currencyName,
		}

		currencyResponse, currencyError := financeserviceapiutils.GetAllCurrency(currencyParams)
		if currencyError != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: currencyError.StatusCode,
				Message:    "Error fetching currency data",
				Err:        currencyError.Err,
			}
		}

		var currencyIds []int
		for _, currency := range currencyResponse {
			currencyIds = append(currencyIds, currency.CurrencyId)
		}

		if len(currencyIds) > 0 {
			baseModelQuery = baseModelQuery.Where("currency_id IN ?", currencyIds)
		} else {
			pages.Rows = []map[string]interface{}{}
			return pages, nil
		}
	}

	var totalRows int64
	if err := baseModelQuery.Count(&totalRows).Error; err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	pages.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pages.GetLimit())))
	pages.TotalPages = totalPages

	if pages.SortBy != "" {
		if strings.ToLower(pages.SortOf) != "desc" {
			pages.SortOf = "asc"
		}
		baseModelQuery = baseModelQuery.Order(pages.SortBy + " " + pages.SortOf)
	} else {
		baseModelQuery = baseModelQuery.Order("purchase_price_id asc")
	}

	err := baseModelQuery.Offset(pages.GetOffset()).Limit(pages.GetLimit()).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, response := range responses {
		// Fetch supplier data
		getSupplierResponse, supplierErr := generalserviceapiutils.GetSupplierMasterById(response.SupplierId)
		if supplierErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: supplierErr.StatusCode,
				Err:        errors.New(supplierErr.Message),
			}
		}

		// Fetch currency data
		getCurrencyResponse, currencyErr := financeserviceapiutils.GetCurrencyId(response.CurrencyId)
		if currencyErr != nil {
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: currencyErr.StatusCode,
				Err:        errors.New(currencyErr.Message),
			}
		}

		// Map the result
		result := map[string]interface{}{
			"purchase_price_id":             response.PurchasePriceId,
			"supplier_id":                   response.SupplierId,
			"supplier_code":                 getSupplierResponse.SupplierCode,
			"supplier_name":                 getSupplierResponse.SupplierName,
			"currency_id":                   response.CurrencyId,
			"currency_code":                 getCurrencyResponse.CurrencyCode,
			"currency_name":                 getCurrencyResponse.CurrencyName,
			"purchase_price_effective_date": response.PurchasePriceEffectiveDate,
			"is_active":                     response.IsActive,
		}

		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
}

func (r *PurchasePriceRepositoryImpl) UpdatePurchasePrice(tx *gorm.DB, Id int, request masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PurchasePrice{}

	result := tx.Model(&entities).
		Where("purchase_price_id = ?", Id).
		First(&entities)

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

	updates := make(map[string]interface{})
	if request.IsActive {
		updates["is_active"] = request.IsActive
	}
	if request.SupplierId != 0 {
		updates["supplier_id"] = request.SupplierId
	}
	if request.CurrencyId != 0 {
		updates["currency_id"] = request.CurrencyId
	}
	if !request.PurchasePriceEffectiveDate.IsZero() {
		updates["purchase_price_effective_date"] = request.PurchasePriceEffectiveDate
	}

	if len(updates) == 0 {
		return entities, nil
	}
	result = tx.Save(&entities)
	if result.Error != nil {
		return masteritementities.PurchasePrice{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
}

func (r *PurchasePriceRepositoryImpl) SavePurchasePrice(tx *gorm.DB, request masteritempayloads.PurchasePriceRequest) (masteritementities.PurchasePrice, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PurchasePrice{
		PurchasePriceId:            request.PurchasePriceId,
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

func (r *PurchasePriceRepositoryImpl) GetPurchasePriceById(tx *gorm.DB, Id int, pagination pagination.Pagination) (masteritempayloads.PurchasePriceResponse, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PurchasePrice{}

	err := tx.Model(&masteritementities.PurchasePrice{}).
		Where("purchase_price_id = ?", Id).
		First(&entities).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masteritempayloads.PurchasePriceResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Purchase price data not found",
				Err:        fmt.Errorf("purchase price with ID %d not found", Id),
			}
		}
		return masteritempayloads.PurchasePriceResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error while fetching purchase price",
			Err:        err,
		}
	}

	getSupplierResponse, supplierErr := generalserviceapiutils.GetSupplierMasterById(entities.SupplierId)
	if supplierErr != nil {
		return masteritempayloads.PurchasePriceResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch supplier data",
			Err:        supplierErr,
		}
	}

	getCurrencyResponse, currencyErr := financeserviceapiutils.GetCurrencyId(entities.CurrencyId)
	if currencyErr != nil {
		return masteritempayloads.PurchasePriceResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error while fetching currency data",
			Err:        errors.New(currencyErr.Message),
		}
	}

	var purchasepriceDetails []masteritempayloads.PurchasePriceDetailResponse
	query := tx.Model(&masteritementities.PurchasePriceDetail{}).
		Select("purchase_price_detail_id", "purchase_price_id", "item_id", "is_active", "purchase_price").
		Where("purchase_price_id = ?", Id).
		Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit())

	errpurchasepriceDetails := query.Find(&purchasepriceDetails).Error
	if errpurchasepriceDetails != nil {
		return masteritempayloads.PurchasePriceResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error while fetching purchase price details",
			Err:        errpurchasepriceDetails,
		}
	}

	for i, detail := range purchasepriceDetails {
		var itemResponse masteritementities.Item
		if err := tx.Where("item_id = ?", detail.ItemId).First(&itemResponse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return masteritempayloads.PurchasePriceResponse{}, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item not found",
					Err:        fmt.Errorf("item with ID %d not found", detail.ItemId),
				}
			}
			return masteritempayloads.PurchasePriceResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item",
				Err:        err,
			}
		}
		purchasepriceDetails[i].ItemCode = itemResponse.ItemCode
		purchasepriceDetails[i].ItemName = itemResponse.ItemName
	}

	payloads := masteritempayloads.PurchasePriceResponse{
		PurchasePriceId:            entities.PurchasePriceId,
		SupplierId:                 entities.SupplierId,
		SupplierCode:               getSupplierResponse.SupplierCode,
		SupplierName:               getSupplierResponse.SupplierName,
		CurrencyId:                 entities.CurrencyId,
		CurrencyCode:               getCurrencyResponse.CurrencyCode,
		CurrencyName:               getCurrencyResponse.CurrencyName,
		PurchasePriceEffectiveDate: entities.PurchasePriceEffectiveDate,
		IsActive:                   entities.IsActive,
		IdentitySysNumber:          0,
		PurchasePriceDetails: masteritempayloads.PurchasePriceDetailsResponse{
			Page:       pagination.GetPage(),
			Limit:      pagination.GetLimit(),
			TotalPages: pagination.TotalPages,
			TotalRows:  int(pagination.TotalRows),
			Data:       purchasepriceDetails,
		},
	}

	return payloads, nil
}

func (r *PurchasePriceRepositoryImpl) GetAllPurchasePriceDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var responses []masteritempayloads.PurchasePriceDetailRequest

	baseModelQuery := tx.Model(&masteritementities.PurchasePriceDetail{})
	whereQuery := utils.ApplyFilter(baseModelQuery, filterCondition)

	err := whereQuery.Scopes(pagination.Paginate(&pages, whereQuery)).Find(&responses).Error
	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if len(responses) == 0 {
		pages.Rows = []map[string]interface{}{}
		return pages, nil
	}

	var results []map[string]interface{}
	for _, response := range responses {

		var itemResponse masteritementities.Item
		if err := tx.Where("item_id = ?", response.ItemId).First(&itemResponse).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pages, &exceptions.BaseErrorResponse{
					StatusCode: http.StatusNotFound,
					Message:    "Item not found",
					Err:        fmt.Errorf("item with ID %d not found", response.ItemId),
				}
			}
			return pages, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to fetch Item",
				Err:        err,
			}
		}

		result := map[string]interface{}{
			"purchase_price_detail_id": response.PurchasePriceDetailId,
			"purchase_price_id":        response.PurchasePriceId,
			"item_id":                  response.ItemId,
			"item_code":                itemResponse.ItemCode,
			"item_name":                itemResponse.ItemName,
			"is_active":                response.IsActive,
			"purchase_price":           response.PurchasePrice,
		}

		results = append(results, result)
	}

	pages.Rows = results

	return pages, nil
}

func (r *PurchasePriceRepositoryImpl) GetPurchasePriceDetailById(tx *gorm.DB, Id int) (masteritempayloads.PurchasePriceDetailResponses, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PurchasePriceDetail{}
	err := tx.Model(&masteritementities.PurchasePriceDetail{}).
		Where("purchase_price_detail_id = ?", Id).
		First(&entities).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masteritempayloads.PurchasePriceDetailResponses{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Data not found",
				Err:        fmt.Errorf("purchase price detail with ID %d not found", Id),
			}
		}
		return masteritempayloads.PurchasePriceDetailResponses{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
			Err:        err,
		}
	}

	// Fetch Item data from external service
	var itemResponse masteritementities.Item
	if err := tx.Where("item_id = ?", entities.ItemId).First(&itemResponse).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masteritempayloads.PurchasePriceDetailResponses{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Item not found",
				Err:        fmt.Errorf("item with ID %d not found", entities.ItemId),
			}
		}
		return masteritempayloads.PurchasePriceDetailResponses{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Item",
			Err:        err,
		}
	}

	payloads := masteritempayloads.PurchasePriceDetailResponses{
		PurchasePriceDetailId: entities.PurchasePriceDetailId,
		PurchasePriceId:       entities.PurchasePriceId,
		ItemId:                entities.ItemId,
		ItemCode:              itemResponse.ItemCode,
		ItemName:              itemResponse.ItemName,
		IsActive:              entities.IsActive,
		PurchasePrice:         entities.PurchasePrice,
	}

	return payloads, nil

}

func (r *PurchasePriceRepositoryImpl) UpdatePurchasePriceDetail(tx *gorm.DB, Id int, request masteritempayloads.PurchasePriceDetailRequest) (masteritementities.PurchasePriceDetail, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PurchasePriceDetail{}

	result := tx.Model(&entities).
		Where("purchase_price_detail_id = ?", Id).
		First(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return masteritementities.PurchasePriceDetail{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        fmt.Errorf("purchase price detail with ID %d not found", Id),
			}
		}
		return masteritementities.PurchasePriceDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	updates := make(map[string]interface{})
	if request.IsActive {
		updates["is_active"] = request.IsActive
	}
	if request.ItemId != 0 {
		updates["item_id"] = request.ItemId
	}
	if request.PurchasePrice != 0 {
		updates["purchase_pprice"] = request.PurchasePrice
	}

	if len(updates) == 0 {
		return entities, nil
	}

	result = tx.Save(&entities)
	if result.Error != nil {
		return masteritementities.PurchasePriceDetail{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	return entities, nil
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

func (r *PurchasePriceRepositoryImpl) DeletePurchasePrice(tx *gorm.DB, Id int, iddet []int) (bool, *exceptions.BaseErrorResponse) {
	var entities []masteritementities.PurchasePriceDetail

	result := tx.Where("purchase_price_id = ? AND purchase_price_detail_id IN ?", Id, iddet).Find(&entities)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusNotFound,
				Err:        result.Error,
			}
		}
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        result.Error,
		}
	}

	if err := tx.Delete(&entities).Error; err != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return true, nil
}

func (r *PurchasePriceRepositoryImpl) ActivatePurchasePriceDetail(tx *gorm.DB, Id int, iddet []int) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PurchasePriceDetail{}

	result := tx.Model(&entities).
		Where("purchase_price_id = ? AND purchase_price_detail_id IN (?)", Id, iddet).
		First(&entities)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        result.Error,
		}
	}

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
			Err:        result.Error,
		}
	}

	entities.IsActive = true

	result = tx.Save(&entities)
	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
			Err:        result.Error,
		}
	}

	return true, nil
}

func (r *PurchasePriceRepositoryImpl) DeactivatePurchasePriceDetail(tx *gorm.DB, Id int, iddet []int) (bool, *exceptions.BaseErrorResponse) {
	entities := masteritementities.PurchasePriceDetail{}

	result := tx.Model(&entities).
		Where("purchase_price_id = ? AND purchase_price_detail_id IN (?)", Id, iddet).
		First(&entities)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Data not found",
			Err:        result.Error,
		}
	}

	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
			Err:        result.Error,
		}
	}

	entities.IsActive = false

	result = tx.Save(&entities)
	if result.Error != nil {
		return false, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
			Err:        result.Error,
		}
	}

	return true, nil
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

func (r *PurchasePriceRepositoryImpl) GetPurchasePriceDetailByParam(tx *gorm.DB, curId int, supId int, effectiveDate string, pages pagination.Pagination) (pagination.Pagination, *exceptions.BaseErrorResponse) {
	var entities masteritementities.PurchasePriceDetail
	response := []masteritempayloads.GetPurchasePriceDetailByParamResponse{}

	// Fetch PurchasePrice data
	baseModelQuery := tx.Model(&entities).
		Select(`
			mtr_purchase_price_detail.purchase_price_detail_id,
			mtr_purchase_price_detail.purchase_price_id,
			mi.item_id,
			mi.item_code,
			mi.item_name,
			mtr_purchase_price_detail.is_active,
			mtr_purchase_price_detail.purchase_price
		`).
		Joins("INNER JOIN mtr_purchase_price mpp ON mpp.purchase_price_id = mtr_purchase_price_detail.purchase_price_id").
		Joins("INNER JOIN mtr_item mi ON mi.item_id = mtr_purchase_price_detail.item_id").
		Where("mpp.currency_id = ?", curId).
		Where("mpp.supplier_id = ?", supId).
		Where("mpp.purchase_price_effective_date >= ? AND mpp.purchase_price_effective_date <= ?", effectiveDate+" 00:00:00.000", effectiveDate+" 23:59:59.999")
	err := baseModelQuery.Scopes(pagination.Paginate(&pages, baseModelQuery)).Scan(&response).Error

	if err != nil {
		return pages, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error while fetching purchase price detail",
			Err:        err,
		}
	}

	pages.Rows = response

	return pages, nil
}
