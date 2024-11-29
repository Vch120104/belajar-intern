package masteritemrepositoryimpl

import (
	config "after-sales/api/config"
	masteritementities "after-sales/api/entities/master/item"
	exceptions "after-sales/api/exceptions"
	masteritempayloads "after-sales/api/payloads/master/item"
	"after-sales/api/payloads/pagination"
	masteritemrepository "after-sales/api/repositories/master/item"
	"after-sales/api/utils"
	financeserviceapiutils "after-sales/api/utils/finance-service"
	generalserviceapiutils "after-sales/api/utils/general-service"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type PurchasePriceRepositoryImpl struct {
}

func StartPurchasePriceRepositoryImpl() masteritemrepository.PurchasePriceRepository {
	return &PurchasePriceRepositoryImpl{}
}

func (r *PurchasePriceRepositoryImpl) GetAllPurchasePrice(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	var responses []masteritempayloads.PurchasePriceRequest

	newFilterCondition := []utils.FilterCondition{}
	supplierCode := ""
	supplierName := ""
	currencyCode := ""
	effectiveDate := ""
	for _, filter := range filterCondition {
		if strings.Contains(filter.ColumnField, "supplier_code") {
			supplierCode = filter.ColumnValue
			continue
		}
		if strings.Contains(filter.ColumnField, "supplier_name") {
			supplierName = filter.ColumnValue
			continue
		}
		if strings.Contains(filter.ColumnField, "currency_code") {
			currencyCode = filter.ColumnValue
			continue
		}
		if strings.Contains(filter.ColumnField, "purchase_price_effective_date") {
			effectiveDate = filter.ColumnValue
			continue
		}
		newFilterCondition = append(newFilterCondition, filter)
	}

	tableStruct := masteritempayloads.PurchasePriceRequest{}
	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)
	whereQuery := utils.ApplyFilter(joinTable, newFilterCondition)

	if supplierCode != "" || supplierName != "" {
		var supplierIds []int
		supplierParams := generalserviceapiutils.SupplierMasterParams{
			Page:         0,
			Limit:        100000,
			SupplierCode: supplierCode,
			SupplierName: supplierName,
		}
		supplierResponse, supplierError := generalserviceapiutils.GetAllSupplierMaster(supplierParams)
		if supplierError != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error fetching supplier filter data",
				Err:        supplierError.Err,
			}
		}

		for _, supplier := range supplierResponse {
			supplierIds = append(supplierIds, supplier.SupplierId)
		}
		if len(supplierIds) == 0 {
			supplierIds = []int{-1}
		}

		whereQuery = whereQuery.Where("supplier_id IN ?", supplierIds)
	}

	if currencyCode != "" {
		var currencyIds []int
		currencyParams := financeserviceapiutils.CurrencyParams{CurrencyCode: currencyCode}
		currencyCodeResponse, currencyError := financeserviceapiutils.GetAllCurrency(currencyParams)
		if currencyError != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Error fetching currency code filter data",
				Err:        currencyError.Err,
			}
		}

		for _, currency := range currencyCodeResponse {
			currencyIds = append(currencyIds, currency.CurrencyId)
		}

		if len(currencyIds) == 0 {
			currencyIds = []int{-1}
		}

		whereQuery = whereQuery.Where("currency_id IN ?", currencyIds)
	}

	if effectiveDate != "" {
		whereQuery = whereQuery.Where("FORMAT(purchase_price_effective_date, 'd MMM yyyy') LIKE (?)", "%"+effectiveDate+"%")
	}

	rows, err := whereQuery.Find(&responses).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
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
		SupplierURL := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(purchasePriceReq.SupplierId)
		var getSupplierResponse masteritempayloads.PurchasePriceSupplierResponse
		if err := utils.Get(SupplierURL, &getSupplierResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Currency data from external service
		getCurrencyResponse, currencyErr := financeserviceapiutils.GetCurrencyId(purchasePriceReq.CurrencyId)
		if currencyErr != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal server error while fetching currency data",
				Err:        errors.New(currencyErr.Message),
			}
		}

		purchasePriceRes = masteritempayloads.PurchasePriceResponse{
			PurchasePriceId:            purchasePriceReq.PurchasePriceId,
			SupplierId:                 purchasePriceReq.SupplierId,
			SupplierCode:               getSupplierResponse.SupplierCode,
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

	if err := rows.Err(); err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
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
			"identity_sys_number":           response.PurchasePriceId,
		}
		mapResponses = append(mapResponses, responseMap)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(mapResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
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

	entities.IsActive = request.IsActive
	entities.SupplierId = request.SupplierId
	entities.CurrencyId = request.CurrencyId
	entities.PurchasePriceEffectiveDate = request.PurchasePriceEffectiveDate

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

	// Fetch PurchasePrice data
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

	// Fetch Supplier data from external service
	SupplierURL := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(entities.SupplierId)
	var getSupplierResponse masteritempayloads.PurchasePriceSupplierResponse
	if err := utils.Get(SupplierURL, &getSupplierResponse, nil); err != nil {
		return masteritempayloads.PurchasePriceResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error while fetching supplier data",
			Err:        err,
		}
	}

	// Fetch Currency data from external service
	getCurrencyResponse, currencyErr := financeserviceapiutils.GetCurrencyId(entities.CurrencyId)
	if currencyErr != nil {
		return masteritempayloads.PurchasePriceResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error while fetching currency data",
			Err:        errors.New(currencyErr.Message),
		}
	}

	// Fetch Purchase Price Detail data
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

	// Fetch Item data for each purchase price detail
	for i, detail := range purchasepriceDetails {
		ItemURL := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(detail.ItemId)
		var itemResponse masteritempayloads.PurchasePriceItemResponse
		if err := utils.Get(ItemURL, &itemResponse, nil); err != nil {
			return masteritempayloads.PurchasePriceResponse{}, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal server error while fetching item data",
				Err:        err,
			}
		}
		purchasepriceDetails[i].ItemCode = itemResponse.ItemCode
		purchasepriceDetails[i].ItemName = itemResponse.ItemName
	}

	// Construct the payload with pagination information
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

func (r *PurchasePriceRepositoryImpl) GetAllPurchasePriceDetail(tx *gorm.DB, filterCondition []utils.FilterCondition, pages pagination.Pagination) ([]map[string]interface{}, int, int, *exceptions.BaseErrorResponse) {

	tableStruct := masteritempayloads.PurchasePriceDetailRequest{}

	joinTable := utils.CreateJoinSelectStatement(tx, tableStruct)

	whereQuery := utils.ApplyFilter(joinTable, filterCondition)

	rows, err := whereQuery.Find(&tableStruct).Rows()
	if err != nil {
		return nil, 0, 0, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	defer rows.Close()

	var convertedResponses []map[string]interface{}

	for rows.Next() {
		var purchasePriceDetailReq masteritempayloads.PurchasePriceDetailRequest

		if err := rows.Scan(&purchasePriceDetailReq.PurchasePriceDetailId, &purchasePriceDetailReq.PurchasePriceId, &purchasePriceDetailReq.ItemId, &purchasePriceDetailReq.IsActive, &purchasePriceDetailReq.PurchasePrice); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		// Fetch Item data from external service
		ItemURL := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(purchasePriceDetailReq.ItemId)
		var getItemResponse masteritempayloads.PurchasePriceItemResponse
		if err := utils.Get(ItemURL, &getItemResponse, nil); err != nil {
			return nil, 0, 0, &exceptions.BaseErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
		}

		purchasePriceDetailRes := map[string]interface{}{
			"purchase_price_detail_id": purchasePriceDetailReq.PurchasePriceDetailId,
			"purchase_price_id":        purchasePriceDetailReq.PurchasePriceId,
			"item_id":                  purchasePriceDetailReq.ItemId,
			"item_code":                getItemResponse.ItemCode,
			"item_name":                getItemResponse.ItemName,
			"is_active":                purchasePriceDetailReq.IsActive,
			"purchase_price":           purchasePriceDetailReq.PurchasePrice,
		}

		convertedResponses = append(convertedResponses, purchasePriceDetailRes)
	}

	paginatedData, totalPages, totalRows := pagination.NewDataFramePaginate(convertedResponses, &pages)

	return paginatedData, totalPages, totalRows, nil
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
	ItemURL := config.EnvConfigs.AfterSalesServiceUrl + "item/" + strconv.Itoa(entities.ItemId)
	var getItemResponse masteritempayloads.PurchasePriceItemResponse
	if err := utils.Get(ItemURL, &getItemResponse, nil); err != nil {
		return masteritempayloads.PurchasePriceDetailResponses{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	payloads := masteritempayloads.PurchasePriceDetailResponses{
		PurchasePriceDetailId: entities.PurchasePriceDetailId,
		PurchasePriceId:       entities.PurchasePriceId,
		ItemId:                entities.ItemId,
		ItemCode:              getItemResponse.ItemCode,
		ItemName:              getItemResponse.ItemName,
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

	entities.IsActive = request.IsActive
	entities.ItemId = request.ItemId
	entities.PurchasePrice = request.PurchasePrice

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
	err := baseModelQuery.Scopes(pagination.Paginate(&entities, &pages, baseModelQuery)).Scan(&response).Error

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
